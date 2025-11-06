package models

import "net/url"

type PackageResolutionMode string

const (
	PackageResolutionModeSoft   PackageResolutionMode = "soft"
	PackageResolutionModeMedium PackageResolutionMode = "medium"
	PackageResolutionModeHard   PackageResolutionMode = "hard"
	PackageResolutionModeLocked PackageResolutionMode = "locked"
)

type PackageResolutionRequest struct {
	Packages []PackageResolutionRequestPackage `json:"packages"`
}

type PackageResolutionRequestPackage struct {
	Org     string                 `json:"org"`
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Mode    *PackageResolutionMode `json:"mode,omitempty"`
}

func (r *PackageResolutionRequest) AddPackage(orgName, name, version string, mode *PackageResolutionMode) {
	if r.Packages == nil {
		r.Packages = []PackageResolutionRequestPackage{}
	}
	encodedVersion := url.QueryEscape(version)
	r.Packages = append(r.Packages, PackageResolutionRequestPackage{
		Org:     orgName,
		Name:    name,
		Version: encodedVersion,
		Mode:    mode,
	})
}
