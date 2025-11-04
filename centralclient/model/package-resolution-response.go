package model

type PackageResolutionResponse struct {
	Resolved   []PackageResolutionResponsePackage `json:"resolved"`
	Unresolved []PackageResolutionResponsePackage `json:"unresolved"`
}

type PackageResolutionResponsePackage struct {
	Org              string                                `json:"org"`
	Name             string                                `json:"name"`
	Version          string                                `json:"version"`
	DependencyGraph  []PackageResolutionResponseDependency `json:"dependencyGraph"`
	IsDeprecated     *bool                                 `json:"isDeprecated,omitempty"`
	DeprecateMessage string                                `json:"deprecateMessage,omitempty"`
}

type PackageResolutionResponseDependency struct {
	Org          string                                `json:"org"`
	Name         string                                `json:"name"`
	Version      string                                `json:"version"`
	Dependencies []PackageResolutionResponseDependency `json:"dependencies"`
}
