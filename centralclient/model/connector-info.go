package model

type ConnectorInfo struct {
	OrgName     string
	PackageName string
	ModuleName  string
	Version     string
	Name        string
}

func NewConnectorInfo(orgName, packageName, moduleName, version, name string) *ConnectorInfo {
	return &ConnectorInfo{
		OrgName:     orgName,
		PackageName: packageName,
		ModuleName:  moduleName,
		Version:     version,
		Name:        name,
	}
}
