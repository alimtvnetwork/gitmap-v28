// Package model — stats.go defines aggregation structs for command history stats.
package model

// CommandStats holds aggregated statistics for a single command.
type CommandStats struct {
	Command      string  `json:"command"`
	TotalRuns    int     `json:"totalRuns"`
	SuccessCount int     `json:"successCount"`
	FailCount    int     `json:"failCount"`
	FailRate     float64 `json:"failRate"`
	AvgDuration  int64   `json:"avgDurationMs"`
	MinDuration  int64   `json:"minDurationMs"`
	MaxDuration  int64   `json:"maxDurationMs"`
	LastUsed     string  `json:"lastUsed"`
}

// OverallStats holds the summary across all commands.
type OverallStats struct {
	TotalCommands   int            `json:"totalCommands"`
	UniqueCommands  int            `json:"uniqueCommands"`
	TotalSuccess    int            `json:"totalSuccess"`
	TotalFail       int            `json:"totalFail"`
	OverallFailRate float64        `json:"overallFailRate"`
	AvgDuration     int64          `json:"avgDurationMs"`
	Commands        []CommandStats `json:"commands"`
}
