package model

type PackageNameResolutionResponse struct {
	ResolvedModules   []PackageNameResolutionResponseModule `json:"resolvedModules"`
	UnresolvedModules []PackageNameResolutionResponseModule `json:"unresolvedModules"`
}

type PackageNameResolutionResponseModule struct {
	Organization string `json:"organization"`
	ModuleName   string `json:"moduleName"`
	Version      string `json:"version"`
	PackageName  string `json:"packageName"`
	Reason       string `json:"reason"`
}
