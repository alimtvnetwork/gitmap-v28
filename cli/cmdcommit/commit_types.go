package cmdcommit

// CommitMetadata encapsulates operational flags and ambient AI provenance.
type CommitMetadata struct {
	TaskSlug  string
	AgentID   string
	HeatmapID string
	AIName    string
	SpecPath  string
	IsAIMode  bool
	HasPush   bool
	IsDryRun  bool
}

const (
	FlagAgentID     = "--agent-id"
	FlagAgent       = "--agent"
	FlagHeatmapID   = "--heatmap-id"
	FlagHeatmap     = "--heatmap"
	FlagAI          = "--ai"
	FlagSpec        = "--spec"
	FlagTask        = "--task"
	FlagTaskShort   = "-t"
	FlagPushLong    = "--push"
	FlagPushShort   = "-p"
	FlagDryRunLong  = "--dry-run"
	FlagDryRunShort = "-n"

	EnvGitmapAgentID          = "GITMAP_AGENT_ID"
	EnvAntigravityAgentID     = "ANTIGRAVITY_AGENT_ID"
	EnvAgentID                = "AGENT_ID"
	EnvGitmapHeatmapID        = "GITMAP_HEATMAP_ID"
	EnvGitmapHeatmapSessionID = "GITMAP_HEATMAP_SESSION_ID"
	EnvGitmapAIName           = "GITMAP_AI_NAME"
	EnvAIModel                = "AI_MODEL"
	EnvGeminiModel            = "GEMINI_MODEL"
	EnvMuseModel              = "MUSE_MODEL"
	EnvGitmapTaskSlug         = "GITMAP_TASK_SLUG"
	EnvGitmapTask             = "GITMAP_TASK"
	EnvGitmapSpecPath         = "GITMAP_SPEC_PATH"
)
