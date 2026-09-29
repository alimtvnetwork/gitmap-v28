package osclean

// TerminalCleanOptions configures terminal history and suggestion sweeping.
type TerminalCleanOptions struct {
	IsDryRun     bool     `json:"isDryRun"`
	HasAutoYes   bool     `json:"hasAutoYes"`
	IsJSON       bool     `json:"isJson"`
	IsVerbose    bool     `json:"isVerbose"`
	OnlyShells   []string `json:"onlyShells,omitempty"` // powershell, bash, zsh, fish, sh
	ShouldReseed bool     `json:"shouldReseed"`         // default: true
}

// ShellCleanStats records metrics for a single shell's cleaned history and suggestions.
type ShellCleanStats struct {
	Shell        string   `json:"shell"`
	Label        string   `json:"label"`
	FilesCleared int      `json:"filesCleared"`
	BytesFreed   int64    `json:"bytesFreed"`
	IsReseeded   bool     `json:"isReseeded"`
	ReseedCount  int      `json:"reseedCount"`
	ClearedPaths []string `json:"clearedPaths,omitempty"`
	Notes        []string `json:"notes,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// TerminalCleanSummary aggregates results across all cleaned terminal shells.
type TerminalCleanSummary struct {
	TotalFilesCleared int               `json:"totalFilesCleared"`
	TotalBytesFreed   int64             `json:"totalBytesFreed"`
	TotalReseedCount  int               `json:"totalReseedCount"`
	Shells            []ShellCleanStats `json:"shells"`
	IsDryRun          bool              `json:"isDryRun"`
	DurationMs        int64             `json:"durationMs"`
}
