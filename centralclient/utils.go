package centralclient

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ballerina-lang-go/common/bfs"

	"github.com/Masterminds/semver/v3"
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
	formatLog(message string) string
}

type logFormatterImpl struct {
	isBuild bool
}

func (l *logFormatterImpl) formatLog(message string) string {
	if l.isBuild {
		return fmt.Sprintf("\t%s", message)
	}
	return message
}

func NewLogFormatter(isBuild bool) LogFormatter {
	return &logFormatterImpl{
		isBuild: isBuild,
	}
}

func createBalaInHomeRepo(balaDownloadResponse *http.Response, fsys fs.FS, pkgPathInBalaCache, pkgOrg, pkgName string, isNightlyBuild bool, deprecationMsg, newUrl, contentDisposition string, outStream io.Writer, logFormatter LogFormatter, trueDigest string) error {
	responseContentLength := balaDownloadResponse.ContentLength
	if responseContentLength <= 0 {
		return NewCentralClientError(logFormatter.formatLog("invalid response from the server, please try again!"))
	}

	resolvedURI := balaDownloadResponse.Header.Get(ResolvedRequestedURI)
	if resolvedURI == "" {
		resolvedURI = newUrl
	}

	uriParts := strings.Split(resolvedURI, "/")
	pkgVersion := uriParts[len(uriParts)-2]

	validPkgVersion, err := validatePackageVersion(pkgVersion, logFormatter)
	if err != nil {
		return err
	}

	balaFile := getBalaFileName(contentDisposition, uriParts[len(uriParts)-1])
	platform := getPlatformFromBala(balaFile, pkgName, pkgVersion)

	// <user.home>.ballerina/bala_cache/<org-name>/<pkg-name>/<pkg-version>
	balaCacheWithPkgPath := filepath.Join(pkgPathInBalaCache, validPkgVersion, platform)

	info, err := fs.Stat(fsys, balaCacheWithPkgPath)
	if err == nil && info.IsDir() {
		entries, err := os.ReadDir(balaCacheWithPkgPath)
		if err != nil {
			return NewPackageAlreadyExistsError(logFormatter.formatLog(fmt.Sprintf("error accessing bala : %s", balaCacheWithPkgPath)), validPkgVersion)
		}

		if len(entries) > 0 {
			deprecatedFilePath := filepath.Join(balaCacheWithPkgPath, DeprecatedMetaFileName)
			if _, err := os.Stat(deprecatedFilePath); err == nil && deprecationMsg == "" {
				if err := os.Remove(deprecatedFilePath); err != nil {
					return NewPackageAlreadyExistsError(logFormatter.formatLog(fmt.Sprintf("error accessing bala : %s", balaCacheWithPkgPath)), validPkgVersion)
				}
			} else if deprecationMsg != "" {
				if err := os.WriteFile(deprecatedFilePath, []byte(deprecationMsg), 0o644); err != nil {
					return NewPackageAlreadyExistsError(logFormatter.formatLog(fmt.Sprintf("error accessing bala : %s", balaCacheWithPkgPath)), validPkgVersion)
				}
			}

			return NewPackageAlreadyExistsError(logFormatter.formatLog(fmt.Sprintf("package already exists in the home repository: %s", balaCacheWithPkgPath)), validPkgVersion)
		}
	}

	// Create the following temp path
	// bala/<org-name>/<pkg-name>/<pkg-version_temp/<platform>
	tempPath := filepath.Join(pkgPathInBalaCache, validPkgVersion+"_temp", platform)
	if err := createBalaFileDirectory(fsys, tempPath, logFormatter); err != nil {
		return err
	}

	if err := writeBalaFile(balaDownloadResponse, fsys, filepath.Join(tempPath, balaFile), fmt.Sprintf("%s/%s:%s", pkgOrg, pkgName, validPkgVersion), outStream, logFormatter, pkgPathInBalaCache, trueDigest); err != nil {
		return err
	}

	tempDir := filepath.Dir(tempPath)
	platformDir := filepath.Dir(balaCacheWithPkgPath)

	if err := bfs.Rename(fsys, tempDir, platformDir); err != nil {
		return NewCentralClientError(logFormatter.formatLog("error creating directory for bala file"))
	}

	if err := handleNightlyBuild(isNightlyBuild, balaCacheWithPkgPath, logFormatter); err != nil {
		return err
	}

	if err := handlePackageDeprecation(deprecationMsg, balaCacheWithPkgPath, logFormatter); err != nil {
		return err
	}

	return nil
}

func validatePackageVersion(pkgVersion string, logFormatter LogFormatter) (string, error) {
	if pkgVersion == "" {
		return "", NewCentralClientError(logFormatter.formatLog("Version cannot be empty"))
	}

	version, err := semver.StrictNewVersion(pkgVersion)
	if err != nil {
		return "", NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("Invalid version: '%s'. %s", pkgVersion, err.Error())))
	}

	return version.String(), nil
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

func getPlatformFromBala(balaName, packageName, version string) string {
	parts := strings.SplitN(balaName, packageName+"-", 2)
	if len(parts) < 2 {
		return ""
	}
	parts = strings.SplitN(parts[1], "-"+version, 2)
	return parts[0]
}

func createBalaFileDirectory(fsys fs.FS, fullPathToStoreBala string, logFormatter LogFormatter) error {
	if err := bfs.MkdirAll(fsys, fullPathToStoreBala, 0o755); err != nil {
		return NewCentralClientError(logFormatter.formatLog("error creating directory for bala file"))
	}
	return nil
}

func writeBalaFile(balaDownloadResponse *http.Response, fsys fs.FS, balaPath, fullPkgName string, outStream io.Writer, logFormatter LogFormatter, homeRepo, trueDigest string) error {
	// file, err := os.Create(balaPath)
	// if err != nil {
	// 	return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("error occurred copying bala file: %s", err.Error())))
	// }
	// defer file.Close()

	balaDownloadResponseBody := balaDownloadResponse.Body
	if balaDownloadResponseBody == nil {
		return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("error occurred extracting bytes of bala file: %s", fullPkgName)))
	}

	if outStream == nil {
		content, err := io.ReadAll(balaDownloadResponseBody)
		if err != nil {
			return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("error occurred reading bala file content: %s", err.Error())))
		}
		err = bfs.WriteFile(fsys, balaPath, content, 0o644)
		if err != nil {
			fmt.Println("Failed to write bala file to the file system:", err)
			return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("error occurred writing bala file content: %s", err.Error())))
		}
	} else {
		// if err := writeAndHandleProgress(balaDownloadResponse.Body, file, fullPkgName, outStream, logFormatter, homeRepo); err != nil {
		// 	return err
		// }
	}

	if err := extractBala(fsys, balaPath, filepath.Dir(balaPath), trueDigest, fullPkgName, outStream); err != nil {
		fmt.Println("Failed to extract bala file:", err)
		return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("error occurred extracting bala file: %s", err.Error())))
	}

	if err := bfs.Remove(fsys, balaPath); err != nil {
		fmt.Println("Failed to remove bala file after extraction:", err)
		return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("error occurred extracting bala file: %s", err.Error())))
	}

	return nil
}

func writeAndHandleProgress(inputStream io.Reader, outputStream io.Writer,
	fullPkgName string, outStream io.Writer, logFormatter LogFormatter, homeRepo string,
) error {
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

	fmt.Fprintln(outStream, logFormatter.formatLog(fmt.Sprintf(" %s pulled from central successfully", fullPkgName)))
	return nil
}

func extractBala(fsys fs.FS, balaFilePath, balaFileDestPath, trueDigest, packageName string, outStream io.Writer) error {
	if err := bfs.MkdirAll(fsys, balaFileDestPath, 0o755); err != nil {
		return err
	}

	actualDigest := SHA256 + checkHashInternal(balaFilePath)
	if trueDigest != "" && trueDigest != actualDigest {
		warning := fmt.Sprintf(`*************************************************************
* WARNING: Certain packages may have originated from sources other than the official distributors. *
*************************************************************

* Verification failed: The hash value of the following package could not be confirmed. 
%s
`, packageName)
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
			// os.MkdirAll(path, file.Mode())
			bfs.MkdirAll(fsys, path, file.Mode())
			continue
		}

		if err := bfs.MkdirAll(fsys, filepath.Dir(path), 0o755); err != nil {
			return err
		}

		outFile, err := bfs.OpenFile(fsys, path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		// _, err = io.Copy(outFile, rc)
		err = bfs.WriteFile(fsys, path, func() []byte {
			data, _ := io.ReadAll(rc)
			return data
		}(), file.Mode())
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func handleNightlyBuild(isNightlyBuild bool, balaCacheWithPkgPath string, logFormatter LogFormatter) error {
	if isNightlyBuild {
		nightlyBuildMetaFile := filepath.Join(balaCacheWithPkgPath, "nightly.build")
		if _, err := os.Stat(nightlyBuildMetaFile); os.IsNotExist(err) {
			errMsg := "error occurred while creating nightly.build file."
			return createMetaFile(nightlyBuildMetaFile, logFormatter, errMsg)
		}
	}
	return nil
}

func handlePackageDeprecation(deprecateMsg, balaCacheWithPkgPath string, logFormatter LogFormatter) error {
	if deprecateMsg != "" {
		deprecateMsgFile := filepath.Join(balaCacheWithPkgPath, DeprecatedMetaFileName)
		if _, err := os.Stat(deprecateMsgFile); os.IsNotExist(err) {
			errMsg := fmt.Sprintf("error occurred while creating the file '%s'.", DeprecatedMetaFileName)
			if err := createMetaFile(deprecateMsgFile, logFormatter, errMsg); err != nil {
				return err
			}
		}
		return writeDeprecatedMsg(deprecateMsgFile, logFormatter, deprecateMsg)
	}
	return nil
}

func writeDeprecatedMsg(metaFilePath string, logFormatter LogFormatter, message string) error {
	if _, err := os.Stat(metaFilePath); err == nil {
		if err := os.WriteFile(metaFilePath, []byte(message), 0o644); err != nil {
			return NewCentralClientError(
				logFormatter.formatLog(fmt.Sprintf("error occurred while writing deprecation message to the file '%s': %s", DeprecatedMetaFileName, err.Error())))
		}
	}
	return nil
}

func createMetaFile(metaFilePath string, logFormatter LogFormatter, errMsg string) error {
	file, err := os.Create(metaFilePath)
	if err != nil {
		return NewCentralClientError(logFormatter.formatLog(errMsg))
	}
	defer file.Close()
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

func getAsList(arrayString string) ([]string, error) {
	var list []string
	if err := json.Unmarshal([]byte(arrayString), &list); err != nil {
		return nil, err
	}
	return list, nil
}
