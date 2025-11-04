package centralclient

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	BufferSize             = 1024
	DeprecatedMetaFileName = "deprecated.txt"
)

var (
	SetBallerinaStageCentral = os.Getenv(BallerinaStageCentral) == "true"
	SetBallerinaDevCentral   = os.Getenv(BallerinaDevCentral) == "true"
	SetTestModeActive        = os.Getenv(TestModeActive) == "true"
)

type LogFormatter interface {
	FormatLog(message string) string
}

type logFormatterImpl struct{}

func NewLogFormatter() LogFormatter {
	return &logFormatterImpl{}
}

func (l *logFormatterImpl) FormatLog(message string) string {
	return message
}

type BuildLogFormatter interface {
	LogFormatter
}

type buildLogFormatterImpl struct {
	logFormatterImpl
}

func NewBuildLogFormatter() BuildLogFormatter {
	return &buildLogFormatterImpl{}
}

func createBalaInHomeRepo(balaDownloadResponse *http.Response, pkgPathInBalaCache, pkgOrg, pkgName string,
	isNightlyBuild bool, deprecationMsg, newURL, contentDisposition string,
	outStream io.Writer, logFormatter LogFormatter, trueDigest string) error {

	responseContentLength := balaDownloadResponse.ContentLength
	if responseContentLength <= 0 {
		return NewCentralClientError(logFormatter.FormatLog("invalid response from the server, please try again"))
	}

	resolvedURI := balaDownloadResponse.Header.Get("Resolved-Requested-URI")
	if resolvedURI == "" {
		resolvedURI = newURL
	}

	uriParts := strings.Split(resolvedURI, "/")
	pkgVersion := uriParts[len(uriParts)-2]

	validPkgVersion, err := validatePackageVersion(pkgVersion, logFormatter)
	if err != nil {
		return err
	}

	balaFile := getBalaFileName(contentDisposition, uriParts[len(uriParts)-1])
	platform := getPlatformFromBala(balaFile, pkgName, validPkgVersion)
	balaCacheWithPkgPath := filepath.Join(pkgPathInBalaCache, validPkgVersion, platform)

	info, err := os.Stat(balaCacheWithPkgPath)
	if err == nil && info.IsDir() {
		entries, err := os.ReadDir(balaCacheWithPkgPath)
		if err != nil {
			return NewPackageAlreadyExistsError(
				logFormatter.FormatLog("error accessing bala: "+balaCacheWithPkgPath),
				validPkgVersion)
		}

		if len(entries) > 0 {
			deprecatedFilePath := filepath.Join(balaCacheWithPkgPath, DeprecatedMetaFileName)
			if _, err := os.Stat(deprecatedFilePath); err == nil && deprecationMsg == "" {
				if err := os.Remove(deprecatedFilePath); err != nil {
					return NewCentralClientErrorWithCause("error removing deprecated file", err)
				}
			} else if deprecationMsg != "" {
				if err := os.WriteFile(deprecatedFilePath, []byte(deprecationMsg), 0644); err != nil {
					return NewCentralClientErrorWithCause("error writing deprecation message", err)
				}
			}

			return NewPackageAlreadyExistsError(
				logFormatter.FormatLog("package already exists in the home repository: "+balaCacheWithPkgPath),
				validPkgVersion)
		}
	}

	tempPath := filepath.Join(pkgPathInBalaCache, validPkgVersion+"_temp", platform)
	if err := createBalaFileDirectory(tempPath, logFormatter); err != nil {
		return err
	}

	if err := writeBalaFile(balaDownloadResponse, filepath.Join(tempPath, balaFile),
		pkgOrg+"/"+pkgName+":"+validPkgVersion, responseContentLength,
		outStream, logFormatter, filepath.Join(pkgPathInBalaCache, validPkgVersion), trueDigest); err != nil {
		return err
	}

	tempDir := filepath.Dir(tempPath)
	platformDir := filepath.Dir(balaCacheWithPkgPath)
	if err := os.Rename(tempDir, platformDir); err != nil {
		return NewCentralClientError(logFormatter.FormatLog("error creating directory for bala file"))
	}

	if err := handleNightlyBuild(isNightlyBuild, balaCacheWithPkgPath, logFormatter); err != nil {
		return err
	}

	return handlePackageDeprecation(deprecationMsg, balaCacheWithPkgPath, logFormatter)
}

func validatePackageVersion(pkgVersion string, logFormatter LogFormatter) (string, error) {
	if pkgVersion == "" {
		return "", NewCentralClientError(logFormatter.FormatLog("Version cannot be empty"))
	}
	return pkgVersion, nil
}

func getBalaFileName(contentDisposition, balaFile string) string {
	if contentDisposition != "" {
		prefix := "attachment; filename="
		if strings.HasPrefix(contentDisposition, prefix) {
			return contentDisposition[len(prefix):]
		}
	}
	return balaFile
}

func createBalaFileDirectory(fullPathToStoreBala string, logFormatter LogFormatter) error {
	if err := os.MkdirAll(fullPathToStoreBala, 0755); err != nil {
		return NewCentralClientError(logFormatter.FormatLog("error creating directory for bala file"))
	}
	return nil
}

func writeBalaFile(balaDownloadResponse *http.Response, balaPath, fullPkgName string, resContentLength int64,
	outStream io.Writer, logFormatter LogFormatter, homeRepo, trueDigest string) error {

	file, err := os.Create(balaPath)
	if err != nil {
		return NewCentralClientError(logFormatter.FormatLog("error occurred copying the bala file: " + err.Error()))
	}
	defer file.Close()

	if outStream == nil {
		if _, err := io.Copy(file, balaDownloadResponse.Body); err != nil {
			return NewCentralClientError(logFormatter.FormatLog("error occurred copying the bala file: " + err.Error()))
		}
	} else {
		if err := writeAndHandleProgress(balaDownloadResponse.Body, file, resContentLength/1024,
			fullPkgName, outStream, logFormatter, homeRepo); err != nil {
			return err
		}
	}

	if err := extractBala(balaPath, filepath.Dir(balaPath), trueDigest, fullPkgName, outStream); err != nil {
		return NewCentralClientError(logFormatter.FormatLog("error occurred extracting the bala file: " + err.Error()))
	}

	if err := os.Remove(balaPath); err != nil {
		return NewCentralClientError(logFormatter.FormatLog("error occurred removing the bala file: " + err.Error()))
	}

	return nil
}

func handleNightlyBuild(isNightlyBuild bool, balaCacheWithPkgPath string, logFormatter LogFormatter) error {
	if isNightlyBuild {
		nightlyBuildMetaFile := filepath.Join(balaCacheWithPkgPath, "nightly.build")
		if _, err := os.Stat(nightlyBuildMetaFile); os.IsNotExist(err) {
			return createMetaFile(nightlyBuildMetaFile, logFormatter, "error occurred while creating nightly.build file.")
		}
	}
	return nil
}

func handlePackageDeprecation(deprecateMsg, balaCacheWithPkgPath string, logFormatter LogFormatter) error {
	if deprecateMsg != "" {
		deprecateMsgFile := filepath.Join(balaCacheWithPkgPath, DeprecatedMetaFileName)
		if _, err := os.Stat(deprecateMsgFile); os.IsNotExist(err) {
			if err := createMetaFile(deprecateMsgFile, logFormatter,
				"error occurred while creating the file '"+DeprecatedMetaFileName+"'"); err != nil {
				return err
			}
		}
		return writeDeprecatedMsg(deprecateMsgFile, logFormatter, deprecateMsg)
	}
	return nil
}

func writeAndHandleProgress(inputStream io.Reader, outputStream io.Writer, totalSizeInKB int64,
	fullPkgName string, outStream io.Writer, logFormatter LogFormatter, homeRepo string) error {

	buffer := make([]byte, BufferSize)
	remoteRepo := getRemoteRepo()
	progressBarTask := fmt.Sprintf("%s [%s -> %s] ", fullPkgName, remoteRepo, homeRepo)

	fmt.Fprint(outStream, progressBarTask)

	for {
		count, err := inputStream.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if _, err := outputStream.Write(buffer[:count]); err != nil {
			return err
		}
	}

	fmt.Fprintln(outStream, logFormatter.FormatLog(fullPkgName+" pulled from central successfully"))
	return nil
}

func createMetaFile(metaFilePath string, logFormatter LogFormatter, message string) error {
	file, err := os.Create(metaFilePath)
	if err != nil {
		return NewCentralClientError(logFormatter.FormatLog(message))
	}
	defer file.Close()
	return nil
}

func writeDeprecatedMsg(metaFilePath string, logFormatter LogFormatter, message string) error {
	if _, err := os.Stat(metaFilePath); err == nil {
		if err := os.WriteFile(metaFilePath, []byte(message), 0644); err != nil {
			return NewCentralClientError(
				logFormatter.FormatLog("error occurred while writing deprecation message to the file '" +
					DeprecatedMetaFileName + "'"))
		}
	}
	return nil
}

func getAsList(arrayString string) ([]string, error) {
	var list []string
	if err := json.Unmarshal([]byte(arrayString), &list); err != nil {
		return nil, err
	}
	return list, nil
}

func getPlatformFromBala(balaName, packageName, version string) string {
	parts := strings.Split(balaName, packageName+"-")
	if len(parts) < 2 {
		return "any"
	}
	platformParts := strings.Split(parts[1], "-"+version)
	if len(platformParts) < 1 {
		return "any"
	}
	return platformParts[0]
}

func extractBala(balaFilePath, balaFileDestPath, trueDigest, packageName string, outStream io.Writer) error {
	if err := os.MkdirAll(balaFileDestPath, 0755); err != nil {
		return err
	}

	actualDigest := "SHA-256=" + checkHashInternal(balaFilePath)
	if trueDigest != "" && trueDigest != actualDigest {
		warning := "*************************************************************\n" +
			"* WARNING: Certain packages may have originated from sources other than the official distributors. *\n" +
			"*************************************************************\n\n" +
			"* Verification failed: The hash value of the following package could not be confirmed. \n" +
			packageName + "\n"
		if outStream != nil {
			fmt.Fprint(outStream, warning)
		}
	}

	reader, err := zip.OpenReader(balaFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(balaFileDestPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func checkHashInternal(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return ""
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func getRemoteRepo() string {
	if SetBallerinaStageCentral {
		return StagingRepo
	} else if SetBallerinaDevCentral {
		return DevRepo
	}
	return ProductionRepo
}
