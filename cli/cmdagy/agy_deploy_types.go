// Package cmdagy — agy_deploy_types.go defines data structures for Antigravity fleet deployment.
package cmdagy

import "time"

// Component names for Antigravity deployment telemetry.
const (
	ComponentTheme        = "theme"
	ComponentPreset       = "preset"
	ComponentPlugins      = "plugins"
	ComponentSkills       = "skills"
	ComponentBinaries     = "binaries"
	ComponentProjects     = "projects"
	ComponentInstances    = "instances"
	ComponentSanitization = "sanitization"
)

// Status codes for individual deployment items.
const (
	DeployStatusSuccess = "SUCCESS"
	DeployStatusSkipped = "SKIPPED"
	DeployStatusFailed  = "FAILED"
)

// AgyDeployOptions holds parsed CLI parameters for Antigravity deployment.
type AgyDeployOptions struct {
	TargetNode  string `json:"targetNode"`
	Preset      string `json:"preset"`
	Theme       string `json:"theme"`
	IsAll       bool   `json:"isAll"`
	HasPlugins  bool   `json:"hasPlugins"`
	HasSkills   bool   `json:"hasSkills"`
	HasBinaries bool   `json:"hasBinaries"`
	HasProjects bool   `json:"hasProjects"`
	IsDryRun    bool   `json:"isDryRun"`
	IsJSON      bool   `json:"isJson"`
	IsForce     bool   `json:"isForce"`
	IsRestart   bool   `json:"isRestart"`
}

// AgyDeployItemResult represents the status of an individual deployment component.
type AgyDeployItemResult struct {
	Component string `json:"component"`
	Status    string `json:"status"` // SUCCESS, SKIPPED, FAILED
	Details   string `json:"details"`
}

// AgyDeployMetrics holds numerical counts and timings.
type AgyDeployMetrics struct {
	PluginsCount        int   `json:"pluginsCount"`
	SkillsCount         int   `json:"skillsCount"`
	ProjectsUpdated     int   `json:"projectsUpdated"`
	SanitizedPathsCount int   `json:"sanitizedPathsCount"`
	BytesTransferred    int64 `json:"bytesTransferred"`
	DurationMs          int64 `json:"durationMs"`
}

// AgyDeployResultJSON defines the structured machine-readable output contract.
type AgyDeployResultJSON struct {
	IsSuccess          bool                  `json:"success"`
	Node               string                `json:"node"`
	IP                 string                `json:"ip"`
	Timestamp          time.Time             `json:"timestamp"`
	DeployedComponents map[string]bool       `json:"deployedComponents"`
	Items              []AgyDeployItemResult `json:"items"`
	Metrics            AgyDeployMetrics      `json:"metrics"`
	Errors             []string              `json:"errors"`
}
