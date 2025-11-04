package model

import "net/url"

type ToolResolutionMode string

const (
	ToolResolutionModeSoft   ToolResolutionMode = "soft"
	ToolResolutionModeMedium ToolResolutionMode = "medium"
	ToolResolutionModeHard   ToolResolutionMode = "hard"
	ToolResolutionModeLocked ToolResolutionMode = "locked"
)

type ToolResolutionCentralRequest struct {
	Tools []ToolResolutionCentralRequestTool `json:"tools"`
}

type ToolResolutionCentralRequestTool struct {
	ID      string              `json:"id"`
	Version string              `json:"version"`
	Mode    *ToolResolutionMode `json:"mode,omitempty"`
}

func (r *ToolResolutionCentralRequest) AddTool(id, version string, mode *ToolResolutionMode) {
	if r.Tools == nil {
		r.Tools = []ToolResolutionCentralRequestTool{}
	}
	encodedVersion := url.QueryEscape(version)
	r.Tools = append(r.Tools, ToolResolutionCentralRequestTool{
		ID:      id,
		Version: encodedVersion,
		Mode:    mode,
	})
}
