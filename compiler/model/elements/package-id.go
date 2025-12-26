package elements

import (
	"fmt"
	"strings"

	"ballerina-lang-go/compiler/common"
)

type PackageID interface {
	GetOrgName() common.Name
	GetPkgName() common.Name
	GetName() common.Name
	GetVersion() common.Name
	GetNameComp(index int) common.Name
	GetNameComps() []common.Name
	GetSourceFileName() common.Name
	GetSourceRoot() string
	IsUnnamed() bool
	SkipTests() bool
	IsTestPkg() bool
	String() string
}

type packageIDImpl struct {
	orgName        common.Name
	pkgName        common.Name
	name           common.Name
	version        common.Name
	nameComps      []common.Name
	sourceFileName common.Name
	sourceRoot     string
	isUnnamed      bool
	skipTests      bool
	isTestPkg      bool
}

func NewPackageID(orgName, name, version common.Name) PackageID {
	return &packageIDImpl{
		orgName:   orgName,
		name:      name,
		pkgName:   name,
		version:   version,
		nameComps: createNameComps(name),
		skipTests: true,
	}
}

func NewPackageIDWithComponents(orgName common.Name, nameComps []common.Name, version common.Name) PackageID {
	nameStrs := make([]string, len(nameComps))
	for i, n := range nameComps {
		nameStrs[i] = n.GetValue()
	}
	name := common.NewName(strings.Join(nameStrs, "."))

	return &packageIDImpl{
		orgName:   orgName,
		name:      name,
		pkgName:   name,
		version:   version,
		nameComps: nameComps,
		skipTests: true,
	}
}

func NewPackageIDWithSourceFile(orgName, pkgName, name, version, sourceFileName common.Name) PackageID {
	return &packageIDImpl{
		orgName:        orgName,
		pkgName:        pkgName,
		name:           name,
		version:        version,
		nameComps:      createNameComps(name),
		sourceFileName: sourceFileName,
		skipTests:      true,
	}
}

func NewPackageIDFull(orgName, pkgName, name, version, sourceFileName common.Name, sourceRoot string, isTestPkg, skipTest bool) PackageID {
	return &packageIDImpl{
		orgName:        orgName,
		pkgName:        pkgName,
		name:           name,
		version:        version,
		nameComps:      createNameComps(name),
		sourceFileName: sourceFileName,
		sourceRoot:     sourceRoot,
		isTestPkg:      isTestPkg,
		skipTests:      skipTest,
	}
}

func NewUnnamedPackageID(orgName common.Name, sourceFileName string, version common.Name) PackageID {
	defaultPkg := common.NewName(".")
	return &packageIDImpl{
		orgName:        orgName,
		name:           defaultPkg,
		pkgName:        defaultPkg,
		version:        version,
		nameComps:      []common.Name{defaultPkg},
		sourceFileName: common.NewName(sourceFileName),
		isUnnamed:      true,
		skipTests:      true,
	}
}

func createNameComps(name common.Name) []common.Name {
	// TODO: replace with names
	nameValue := name.GetValue()
	if nameValue == "." {
		return []common.Name{name}
	}

	parts := strings.Split(nameValue, ".")
	comps := make([]common.Name, len(parts))
	for i, part := range parts {
		comps[i] = common.NewName(part)
	}
	return comps
}

func (p *packageIDImpl) GetOrgName() common.Name {
	return p.orgName
}

func (p *packageIDImpl) GetPkgName() common.Name {
	return p.pkgName
}

func (p *packageIDImpl) GetName() common.Name {
	return p.name
}

func (p *packageIDImpl) GetVersion() common.Name {
	return p.version
}

func (p *packageIDImpl) GetNameComp(index int) common.Name {
	return p.nameComps[index]
}

func (p *packageIDImpl) GetNameComps() []common.Name {
	return p.nameComps
}

func (p *packageIDImpl) GetSourceFileName() common.Name {
	return p.sourceFileName
}

func (p *packageIDImpl) GetSourceRoot() string {
	return p.sourceRoot
}

func (p *packageIDImpl) IsUnnamed() bool {
	return p.isUnnamed
}

func (p *packageIDImpl) SkipTests() bool {
	return p.skipTests
}

func (p *packageIDImpl) IsTestPkg() bool {
	return p.isTestPkg
}

func (p *packageIDImpl) String() string {
	// TODO: replace with names
	if p.name.GetValue() == "." {
		return "."
	}

	org := ""
	if p.orgName != nil && p.orgName.GetValue() != "$anon" {
		org = p.orgName.GetValue() + "/"
	}

	if p.version.GetValue() == "" {
		return org + p.name.GetValue()
	}

	return fmt.Sprintf("%s%s:%s", org, p.name.GetValue(), p.version.GetValue())
}
