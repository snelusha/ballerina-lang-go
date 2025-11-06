package models

type PackageSearchResult struct {
	Packages []Package `json:"packages"`
	Count    int       `json:"count"`
}
