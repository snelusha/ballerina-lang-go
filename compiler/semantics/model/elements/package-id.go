package elements

import "ballerina-lang-go/compiler/util"

type PackageID struct {
	Org            util.Name
	PkgName        util.Name
	Name           util.Name
	Version        util.Name
	SourceFileName util.Name
	SourceRoot     string
	IsTestPkg      bool
	SkipTest       bool
}

func NewPackageID(org, name, version util.Name) PackageID {
	return PackageID{
		Org:     org,
		Name:    name,
		Version: version,
	}
}

func NewPackageIDFull(org, pkgName, name, version, sourceFileName util.Name, sourceRoot string, isTestPkg, skipTest bool) PackageID {
	return PackageID{
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
