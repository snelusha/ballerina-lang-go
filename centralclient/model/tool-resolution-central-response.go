package model

type ToolResolutionCentralResponse struct {
	Resolved   []ToolResolutionCentralResponseResolvedTool   `json:"resolved"`
	Unresolved []ToolResolutionCentralResponseUnresolvedTool `json:"unresolved"`
}

type ToolResolutionCentralResponseUnresolvedTool struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type ToolResolutionCentralResponseResolvedTool struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
	OrgName string `json:"orgName"`
}
