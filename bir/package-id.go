package bir

import (
	"fmt"
	"strings"
)

type PackageID struct {
	OrgName        *Name
	PkgName        *Name
	Name           *Name
	Version        *Name
	NameComps      []*Name
	IsUnnamed      bool
	SkipTests      bool
	IsTestPkg      bool
	SourceFileName *Name
	SourceRoot     string
}

func NewPackageID(orgName, pkgName, name, version *Name) *PackageID {
	return &PackageID{
		OrgName:   orgName,
		PkgName:   pkgName,
		Name:      name,
		Version:   version,
		NameComps: createNameComps(name),
		IsUnnamed: false,
		SkipTests: true,
		IsTestPkg: false,
	}
}

func NewPackageIDWithNameComps(orgName *Name, nameComps []*Name, version *Name) *PackageID {
	var parts []string
	for _, nc := range nameComps {
		parts = append(parts, nc.Value)
	}
	name := NewName(strings.Join(parts, "."))

	return &PackageID{
		OrgName:   orgName,
		PkgName:   name,
		Name:      name,
		Version:   version,
		NameComps: nameComps,
		IsUnnamed: false,
		SkipTests: true,
		IsTestPkg: false,
	}
}

func NewPackageIDWithSourceFileName(orgName, name, version, sourceFileName *Name) *PackageID {
	return &PackageID{
		OrgName:        orgName,
		PkgName:        name,
		Name:           name,
		Version:        version,
		NameComps:      createNameComps(name),
		SourceFileName: sourceFileName,
		IsUnnamed:      false,
		SkipTests:      true,
		IsTestPkg:      false,
	}
}

func NewPackageIDFull(orgName, pkgName, name, version, sourceFileName *Name, sourceRoot string, isTestPkg, skipTest bool) *PackageID {
	return &PackageID{
		OrgName:        orgName,
		PkgName:        pkgName,
		Name:           name,
		Version:        version,
		NameComps:      createNameComps(name),
		SourceFileName: sourceFileName,
		SourceRoot:     sourceRoot,
		IsTestPkg:      isTestPkg,
		SkipTests:      skipTest,
		IsUnnamed:      false,
	}
}

func NewUnnamedPackageID(orgName *Name, sourceFileName string, version *Name) *PackageID {
	defaultPackage := NewName(".")
	return &PackageID{
		OrgName:        orgName,
		Name:           defaultPackage,
		PkgName:        defaultPackage,
		Version:        version,
		NameComps:      []*Name{defaultPackage},
		IsUnnamed:      true,
		SourceFileName: NewName(sourceFileName),
		SkipTests:      true,
		IsTestPkg:      false,
	}
}

func NewUnnamedPackageIDSimple(sourceFileName string) *PackageID {
	anonOrg := NewName("$anon")
	defaultPackage := NewName(".")
	defaultVersion := NewName("0.0.0")

	return &PackageID{
		OrgName:        anonOrg,
		Name:           defaultPackage,
		PkgName:        defaultPackage,
		Version:        defaultVersion,
		NameComps:      []*Name{defaultPackage},
		IsUnnamed:      true,
		SourceFileName: NewName(sourceFileName),
		SkipTests:      true,
		IsTestPkg:      false,
	}
}

func createNameComps(name *Name) []*Name {
	if name.Value == "." {
		return []*Name{NewName(".")}
	}

	parts := strings.Split(name.Value, ".")
	nameComps := make([]*Name, len(parts))
	for i, part := range parts {
		nameComps[i] = NewName(part)
	}
	return nameComps
}

func (p *PackageID) GetPkgName() *Name {
	return p.PkgName
}

func (p *PackageID) GetName() *Name {
	return p.Name
}

func (p *PackageID) GetNameComp(index int) *Name {
	return p.NameComps[index]
}

func (p *PackageID) GetNameComps() []*Name {
	return p.NameComps
}

func (p *PackageID) GetPackageVersion() *Name {
	return p.Version
}

func (p *PackageID) GetOrgName() *Name {
	return p.OrgName
}

func (p *PackageID) Equals(other *PackageID) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}

	samePkg := false
	if p.IsUnnamed == other.IsUnnamed {
		if !p.IsUnnamed {
			samePkg = true
		} else {
			samePkg = p.SourceFileName.Equals(other.SourceFileName)
		}
	}

	return samePkg &&
		p.OrgName.Equals(other.OrgName) &&
		p.PkgName.Equals(other.PkgName) &&
		p.Name.Equals(other.Name) &&
		p.Version.Equals(other.Version)
}

func (p *PackageID) String() string {
	if p.Name.Value == "." {
		return p.Name.Value
	}

	org := ""
	anonOrg := NewName("$anon")
	if p.OrgName != nil && !p.OrgName.Equals(anonOrg) {
		org = p.OrgName.Value + "/"
	}

	emptyVersion := NewName("")
	if p.Version.Equals(emptyVersion) {
		return org + p.Name.Value
	}

	return fmt.Sprintf("%s%s:%s", org, p.Name.Value, p.Version.Value)
}

func IsLangLibPackageID(packageID *PackageID) bool {
	ballerinaOrg := NewName("ballerina")
	if !packageID.GetOrgName().Equals(ballerinaOrg) {
		return false
	}

	langName := NewName("lang")
	javaName := NewName("java")

	return (len(packageID.NameComps) > 1 && packageID.NameComps[0].Equals(langName)) ||
		packageID.Name.Equals(javaName)
}
