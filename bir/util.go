package bir

import (
	"errors"

	"ballerina-lang-go/bir/generated"
)

func GetString(cpEntries []*generated.Bir_ConstantPoolEntry, index int32) (string, error) {
	if name, ok := cpEntries[index].CpInfo.(*generated.Bir_StringCpInfo); ok {
		return name.Value, nil
	}
	return "", errors.New("Invalid constant pool entry type for string")
}

func GetPackage(cpEntries []*generated.Bir_ConstantPoolEntry, index int32) (*generated.Bir_PackageCpInfo, error) {
	if pkg, ok := cpEntries[index].CpInfo.(*generated.Bir_PackageCpInfo); ok {
		return pkg, nil
	}
	return nil, errors.New("Invalid constant pool entry type for package")
}
