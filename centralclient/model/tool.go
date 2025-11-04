package model

type Tool struct {
	BalToolID                    string   `json:"balToolId"`
	Organization                 string   `json:"organization"`
	Name                         string   `json:"name"`
	Version                      string   `json:"version"`
	Platform                     string   `json:"platform"`
	LanguageSpecificationVersion string   `json:"languageSpecificationVersion"`
	IsDeprecated                 *bool    `json:"isDeprecated,omitempty"`
	DeprecateMessage             *bool    `json:"deprecateMessage,omitempty"`
	URL                          string   `json:"URL"`
	BalaVersion                  string   `json:"balaVersion"`
	BalaURL                      string   `json:"balaURL"`
	Readme                       string   `json:"readme"`
	Licenses                     []string `json:"licenses"`
	Authors                      []string `json:"authors"`
	SourceCodeLocation           string   `json:"sourceCodeLocation"`
	Keywords                     []string `json:"keywords"`
	BallerinaVersion             string   `json:"ballerinaVersion"`
	CreatedDate                  *int64   `json:"createdDate,omitempty"`
	Modules                      []Module `json:"modules"`
	Summary                      string   `json:"summary"`
	Icon                         string   `json:"icon"`
}
