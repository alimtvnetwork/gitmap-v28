package osclean

// DevCleanOptions configures developer tools cache sweeping.
type DevCleanOptions struct {
	IsDryRun       bool     `json:"isDryRun"`
	HasAutoYes     bool     `json:"hasAutoYes"`
	IsJSON         bool     `json:"isJson"`
	IsVerbose      bool     `json:"isVerbose"`
	OnlyCategories []string `json:"onlyCategories,omitempty"`
}

// CategoryCleanStats records metrics for a single developer cache category.
type CategoryCleanStats struct {
	Category     string   `json:"category"`
	Label        string   `json:"label"`
	ItemsRemoved int      `json:"itemsRemoved"`
	DirsRemoved  int      `json:"dirsRemoved"`
	BytesFreed   int64    `json:"bytesFreed"`
	Notes        []string `json:"notes,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// DevCleanSummary aggregates results across all cleaned categories.
type DevCleanSummary struct {
	TotalBytesFreed   int64                `json:"totalBytesFreed"`
	TotalItemsRemoved int                  `json:"totalItemsRemoved"`
	TotalDirsRemoved  int                  `json:"totalDirsRemoved"`
	Categories        []CategoryCleanStats `json:"categories"`
	IsDryRun          bool                 `json:"isDryRun"`
	DurationMs        int64                `json:"durationMs"`
}
