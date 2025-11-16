package main

import (
	"ballerina-lang-go/centralclient"
	"ballerina-lang-go/common/bfs"
	"io/fs"
)

var (
	baseUrl           = "https://api.central.ballerina.io/2.0/registry"
	supportedPlatform = "java21,java17,java11"
	ballerinaVersion  = "2201.13.0-m3"

	orgName     = "ballerina"
	packageName = "url"
	version     = "2.6.1"
	// version = ""

	packagePathInBalaCache = "/Users/sithi/.ballerin sfma/repositories/central.ballerina.io/bala/ballerina/url"

	isBuild = false
)

func main() {
	client := centralclient.NewCentralAPIClient(baseUrl, nil, "")

	memFS := bfs.NewMemFS()

	err := client.PullPackage(orgName, packageName, version, memFS, ".ballerina", supportedPlatform, ballerinaVersion, isBuild)
	if err != nil {
		panic(err)
	}

	bfs.PrintFiles(memFS)

	file, err := fs.ReadFile(memFS, ".ballerina/2.6.1/java21/bala.json")
	if err != nil {
		panic(err)
	}

	println(string(file))

	// connectors, err := client.GetConnectors(map[string]string{
	// 	"q": "paypal.orders",
	// }, supportedPlatform, ballerinaVersion)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("Connectors: %+v\n", connectors)

	// trigger, err := client.GetTrigger("90", supportedPlatform, ballerinaVersion)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("Trigger: %+v\n", trigger["displayName"])

	// triggers, err := client.GetTriggers(map[string]string{
	// 	"moduleName": "twilio",
	// }, supportedPlatform, ballerinaVersion)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("Triggers: %+v\n", triggers)

	// err := client.PullPackage(orgName, packageName, version, packagePathInBalaCache, supportedPlatform, ballerinaVersion, isBuild)
	// if err != nil {
	// 	panic(err)
	// }
	// println("Package pulled successfully")
	//
	// connector, err := client.GetConnector("10573", supportedPlatform, ballerinaVersion)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("Connector: %+v\n", connector["name"])
	//
	// connector, err := client.GetConnectorByInfo(&models.ConnectorInfo{
	// 	OrgName:     "ballerinax",
	// 	ModuleName:  "paypal.orders",
	// 	PackageName: "paypal.orders",
	// 	Version:     "2.0.0",
	// 	Name:        "Client",
	// }, supportedPlatform, ballerinaVersion)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("Connector: %+v\n", connector["id"])

	// versions, err := client.GetPackageVersions(orgName, packageName, supportedPlatform, ballerinaVersion)
	// if err != nil {
	// 	panic(err)
	// }
	// println("Versions:")
	// for _, v := range versions {
	// 	println(v)
	// }

	// pkg, err := client.GetPackage(orgName, packageName, version, supportedPlatform, ballerinaVersion)
	// if err != nil {
	// 	fmt.Printf("Error while fetching package: %v\n", err)
	// 	return
	// }
	//
	// fmt.Printf("Package: %+v\n", pkg)
}
