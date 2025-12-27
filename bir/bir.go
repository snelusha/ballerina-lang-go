package bir

import (
	"fmt"
	"os"
	"path/filepath"

	"ballerina-lang-go/bir/generated"
	"ballerina-lang-go/compiler/bir/model"
	"ballerina-lang-go/compiler/common"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func Parse(path string) error {
	file, err := os.Open(filepath.Join("bir", "testdata", "sample.bir"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bir := generated.NewBir()

	err = bir.Read(kaitai.NewStream(file), nil, bir)
	if err != nil {
		panic(err)
	}

	pkgId, err := GetPackage(bir.ConstantPool.ConstantPoolEntries, bir.Module.IdCpIndex)
	if err != nil {
		panic(err)
	}

	org, err := GetString(bir.ConstantPool.ConstantPoolEntries, pkgId.OrgIndex)
	if err != nil {
		panic(err)
	}

	name, err := GetString(bir.ConstantPool.ConstantPoolEntries, pkgId.NameIndex)
	if err != nil {
		panic(err)
	}

	version, err := GetString(bir.ConstantPool.ConstantPoolEntries, pkgId.VersionIndex)
	if err != nil {
		panic(err)
	}

	Package := model.NewBIRPackage(nil, common.NewName(org), common.NewName(name), common.NewName(name), common.NewName(version), nil, "", false, false)
	fmt.Println(Package.GetPackageID().GetOrgName())

	return nil
}
