package elements

type PackageID struct {
	Org            string
	PkgName        string
	Name           string
	Version        string
	SourceFileName string
	SourceRoot     string
	IsTestPkg      bool
	SkipTest       bool
}

func NewPackageID(org, pkgName, name, version, sourceFileName, sourceRoot string, isTestPkg, skipTest bool) *PackageID {
	return &PackageID{
		Org:            org,
		PkgName:        pkgName,
		Name:           name,
		Version:        version,
		SourceFileName: sourceFileName,
		SourceRoot:     sourceRoot,
		IsTestPkg:      isTestPkg,
		SkipTest:       skipTest,
	}
}

func NewPackageIDShort(org, name, version string) *PackageID {
	return &PackageID{
		Org:     org,
		Name:    name,
		Version: version,
	}
}
