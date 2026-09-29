// Package store — watch_prompts_types.go defines data types for watch-prompts-running Split-DB storage.
package store

// WatchPromptRecord models a backed-up prompt record inside a repo-scoped watch-prompts Split-DB.
type WatchPromptRecord struct {
	RecordId       string   `json:"RecordId"`
	RepoSlug       string   `json:"RepoSlug"`
	ProjectName    string   `json:"ProjectName"`
	ProjectPath    string   `json:"ProjectPath"`
	ProjectId      string   `json:"ProjectId"`
	SequenceId     string   `json:"SequenceId"`
	ConversationId string   `json:"ConversationId"`
	PromptText     string   `json:"PromptText"`
	PromptStatus   string   `json:"PromptStatus"`
	WordCount      int      `json:"WordCount"`
	MediaPaths     []string `json:"MediaPaths"`
	IsActive       bool     `json:"IsActive"`
	UpdatedAt      string   `json:"UpdatedAt"`
	CreatedAt      string   `json:"CreatedAt"`
}

// WatchLogRecord logs watcher events, heartbeats, and recovery executions.
type WatchLogRecord struct {
	LogId     int64  `json:"LogId"`
	RepoSlug  string `json:"RepoSlug"`
	Event     string `json:"Event"`
	Message   string `json:"Message"`
	Status    string `json:"Status"`
	CreatedAt string `json:"CreatedAt"`
}

// WatchConfigRecord holds configuration key-value pairs stored in the watch split DB.
type WatchConfigRecord struct {
	ConfigKey   string `json:"ConfigKey"`
	ConfigValue string `json:"ConfigValue"`
	UpdatedAt   string `json:"UpdatedAt"`
}

// WatchPromptsSummary models a consolidated summary of a watched project.
type WatchPromptsSummary struct {
	RepoSlug     string   `json:"RepoSlug"`
	ProjectName  string   `json:"ProjectName"`
	ProjectPath  string   `json:"ProjectPath"`
	ProjectId    string   `json:"ProjectId"`
	SequenceId   string   `json:"SequenceId"`
	PromptStatus string   `json:"PromptStatus"`
	WordCount    int      `json:"WordCount"`
	MediaCount   int      `json:"MediaCount"`
	MediaPaths   []string `json:"MediaPaths"`
	IsWatching   bool     `json:"IsWatching"`
	DatabasePath string   `json:"DatabasePath"`
	UpdatedAt    string   `json:"UpdatedAt"`
}
