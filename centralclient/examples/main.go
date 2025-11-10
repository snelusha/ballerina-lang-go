package main

import (
	"ballerina-lang-go/centralclient"
)

var (
	baseUrl           = "https://api.central.ballerina.io/2.0/registry"
	supportedPlatform = "java21,java17,java11"
	ballerinaVersion  = "2201.13.0-m3"

	orgName     = "ballerina"
	packageName = "url"
	// version     = "2.6.1"
	version = ""

	packagePathInBalaCache = "/Users/sithi/.ballerina/repositories/central.ballerina.io/bala/ballerina/url"

	isBuild = false
)

func main() {
	client := centralclient.NewCentralAPIClient(baseUrl, nil, "")

	err := client.PullPackage(orgName, packageName, version, packagePathInBalaCache, supportedPlatform, ballerinaVersion, isBuild)
	if err != nil {
		panic(err)
	}
	println("Package pulled successfully")

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
