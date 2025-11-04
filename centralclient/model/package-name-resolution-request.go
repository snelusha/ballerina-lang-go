package model

type PackageNameResolutionRequest struct {
	Modules []PackageNameResolutionRequestModule `json:"modules"`
}

type PackageNameResolutionRequestModule struct {
	Organization     string                                        `json:"organization"`
	ModuleName       string                                        `json:"moduleName"`
	PossiblePackages []PackageNameResolutionRequestPossiblePackage `json:"possiblePackages,omitempty"`
	Mode             *PackageResolutionMode                        `json:"mode,omitempty"`
}

type PackageNameResolutionRequestPossiblePackage struct {
	Org     string `json:"org"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (r *PackageNameResolutionRequest) AddModule(orgName, moduleName string, possiblePackages []PackageNameResolutionRequestPossiblePackage, mode *PackageResolutionMode) {
	if r.Modules == nil {
		r.Modules = []PackageNameResolutionRequestModule{}
	}
	r.Modules = append(r.Modules, PackageNameResolutionRequestModule{
		Organization:     orgName,
		ModuleName:       moduleName,
		PossiblePackages: possiblePackages,
		Mode:             mode,
	})
}

func (r *PackageNameResolutionRequest) AddModuleSimple(orgName, moduleName string) {
	if r.Modules == nil {
		r.Modules = []PackageNameResolutionRequestModule{}
	}
	r.Modules = append(r.Modules, PackageNameResolutionRequestModule{
		Organization: orgName,
		ModuleName:   moduleName,
	})
}
