package cmdcommit

import (
	commitconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit/constants"
)

// CommitMetadata encapsulates operational flags and ambient AI provenance.
type CommitMetadata struct {
	TaskSlug  string
	AgentID   string
	HeatmapID string
	AIName    string
	SpecPath  string
	Author    string
	IsAIMode  bool
	HasPush   bool
	IsDryRun  bool
}

const (
	FlagAgentID     = commitconstants.FlagAgentID
	FlagAgent       = "--agent"
	FlagHeatmapID   = commitconstants.FlagHeatmapID
	FlagHeatmap     = "--heatmap"
	FlagAI          = commitconstants.FlagAI
	FlagSpec        = commitconstants.FlagSpec
	FlagTask        = commitconstants.FlagTask
	FlagTaskShort   = "-t"
	FlagAuthor      = commitconstants.FlagAuthor
	FlagAuthorShort = "-a"
	FlagPushLong    = commitconstants.FlagPush
	FlagPushShort   = "-p"
	FlagDryRunLong  = commitconstants.FlagDryRun
	FlagDryRunShort = "-n"

	EnvGitmapAgentID          = commitconstants.EnvGitmapAgentID
	EnvAntigravityAgentID     = commitconstants.EnvAntigravityAgentID
	EnvAgentID                = "AGENT_ID"
	EnvGitmapHeatmapID        = commitconstants.EnvGitmapHeatmapID
	EnvGitmapHeatmapSessionID = "GITMAP_HEATMAP_SESSION_ID"
	EnvGitmapAIName           = "GITMAP_AI_NAME"
	EnvAIModel                = commitconstants.EnvAIModel
	EnvGeminiModel            = "GEMINI_MODEL"
	EnvMuseModel              = "MUSE_MODEL"
	EnvGitmapTaskSlug         = "GITMAP_TASK_SLUG"
	EnvGitmapTask             = "GITMAP_TASK"
	EnvGitmapSpecPath         = "GITMAP_SPEC_PATH"
)
