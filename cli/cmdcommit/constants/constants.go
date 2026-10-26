package constants

// CLI command flag literals for commit commands.
const (
	FlagAgentID   = "--agent-id"
	FlagHeatmapID = "--heatmap-id"
	FlagAI        = "--ai"
	FlagSpec      = "--spec"
	FlagTask      = "--task"
	FlagAuthor    = "--author"
	FlagPush      = "--push"
	FlagDryRun    = "--dry-run"
)

// Environment variable names for ambient commit metadata.
const (
	EnvGitmapAgentID      = "GITMAP_AGENT_ID"
	EnvAntigravityAgentID = "ANTIGRAVITY_AGENT_ID"
	EnvGitmapHeatmapID    = "GITMAP_HEATMAP_ID"
	EnvAIModel            = "AI_MODEL"
)

// Regular expression patterns for commit message parsing.
const (
	PatternHyphen          = `^([A-Za-z0-9_\-\./]+)\s*[:-]\s*(.+)$`
	PatternBracketedModule = `^\(([A-Za-z0-9_\-\./]+)\)\s*[:-]\s*(.+)$`
	PatternConventional    = `^([A-Za-z]+)\(([A-Za-z0-9_\-\./]+)\):\s*(.+)$`
)
