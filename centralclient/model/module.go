package model

type Module struct {
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	Readme     string `json:"readme"`
	APIDocURL  string `json:"apiDocURL"`
	Executable *bool  `json:"executable,omitempty"`
	PackageURL string `json:"packageUrl"`
}
