package elements

import (
	"fmt"
	"strings"

	"ballerina-lang-go/compiler/bir/model"
)

type PackageID interface {
	GetOrgName() model.Name
	GetPkgName() model.Name
	GetName() model.Name
	GetVersion() model.Name
	GetNameComp(index int) model.Name
	GetNameComps() []model.Name
	GetSourceFileName() model.Name
	GetSourceRoot() string
	IsUnnamed() bool
	SkipTests() bool
	IsTestPkg() bool
	String() string
}

type packageIDImpl struct {
	orgName        model.Name
	pkgName        model.Name
	name           model.Name
	version        model.Name
	nameComps      []model.Name
	sourceFileName model.Name
	sourceRoot     string
	isUnnamed      bool
	skipTests      bool
	isTestPkg      bool
}

func NewPackageID(orgName, name, version model.Name) PackageID {
	return &packageIDImpl{
		orgName:   orgName,
		name:      name,
		pkgName:   name,
		version:   version,
		nameComps: createNameComps(name),
		skipTests: true,
	}
}

func NewPackageIDWithComponents(orgName model.Name, nameComps []model.Name, version model.Name) PackageID {
	nameStrs := make([]string, len(nameComps))
	for i, n := range nameComps {
		nameStrs[i] = n.GetValue()
	}
	name := model.NewName(strings.Join(nameStrs, "."))

	return &packageIDImpl{
		orgName:   orgName,
		name:      name,
		pkgName:   name,
		version:   version,
		nameComps: nameComps,
		skipTests: true,
	}
}

func NewPackageIDWithSourceFile(orgName, pkgName, name, version, sourceFileName model.Name) PackageID {
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

func NewPackageIDFull(orgName, pkgName, name, version, sourceFileName model.Name, sourceRoot string, isTestPkg, skipTest bool) PackageID {
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

func NewUnnamedPackageID(orgName model.Name, sourceFileName string, version model.Name) PackageID {
	defaultPkg := model.NewName(".")
	return &packageIDImpl{
		orgName:        orgName,
		name:           defaultPkg,
		pkgName:        defaultPkg,
		version:        version,
		nameComps:      []model.Name{defaultPkg},
		sourceFileName: model.NewName(sourceFileName),
		isUnnamed:      true,
		skipTests:      true,
	}
}

func createNameComps(name model.Name) []model.Name {
	nameValue := name.GetValue()
	if nameValue == "." {
		return []model.Name{name}
	}

	parts := strings.Split(nameValue, ".")
	comps := make([]model.Name, len(parts))
	for i, part := range parts {
		comps[i] = model.NewName(part)
	}
	return comps
}

func (p *packageIDImpl) GetOrgName() model.Name {
	return p.orgName
}

func (p *packageIDImpl) GetPkgName() model.Name {
	return p.pkgName
}

func (p *packageIDImpl) GetName() model.Name {
	return p.name
}

func (p *packageIDImpl) GetVersion() model.Name {
	return p.version
}

func (p *packageIDImpl) GetNameComp(index int) model.Name {
	return p.nameComps[index]
}

func (p *packageIDImpl) GetNameComps() []model.Name {
	return p.nameComps
}

func (p *packageIDImpl) GetSourceFileName() model.Name {
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
