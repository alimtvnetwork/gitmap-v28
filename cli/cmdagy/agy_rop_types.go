package cmdagy

import "time"

// AGYConversationBackup represents a backed-up Antigravity conversation record.
type AGYConversationBackup struct {
	ConversationID string    `json:"conversationId"`
	ProjectSlug    string    `json:"projectSlug"`
	ProjectPath    string    `json:"projectPath"`
	Title          string    `json:"title"`
	StepCount      int       `json:"stepCount"`
	TranscriptJSON string    `json:"transcriptJson"`
	CreatedAt      time.Time `json:"createdAt"`
}

// AGYROPProjectResult captures per-project reread and optimization outcome.
type AGYROPProjectResult struct {
	ProjectName       string `json:"projectName"`
	ProjectSlug       string `json:"projectSlug"`
	ProjectPath       string `json:"projectPath"`
	ConvsBackedUp     int    `json:"convsBackedUp"`
	NewConversationID string `json:"newConversationId,omitempty"`
	Status            string `json:"status"`
	ErrorMsg          string `json:"errorMsg,omitempty"`
}

// AGYROPResult represents the overall batch result for reread-optimize-project.
type AGYROPResult struct {
	RequestedN    int                   `json:"requestedN"`
	ActiveFound   int                   `json:"activeFound"`
	CacheCleared  bool                  `json:"cacheCleared"`
	TotalBackedUp int                   `json:"totalBackedUp"`
	Projects      []AGYROPProjectResult `json:"projects"`
}
