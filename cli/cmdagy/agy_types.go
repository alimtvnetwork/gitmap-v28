// Package cmd — agy_types.go provides data structures and helpers for Antigravity projects.
package cmdagy

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type AgyProject struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	ProjectResources *AgyProjectResources   `json:"projectResources,omitempty"`
	Settings         map[string]interface{} `json:"settings,omitempty"`
	UpdatedAt        string                 `json:"updatedAt,omitempty"`
	IsWorkspaceOnly  bool                   `json:"isWorkspaceOnly,omitempty"`
}

type AgyProjectResources struct {
	Resources []AgyResource `json:"resources,omitempty"`
}

type AgyResource struct {
	GitFolder *AgyGitFolder `json:"gitFolder,omitempty"`
}

type AgyGitFolder struct {
	FolderURI     string `json:"folderUri,omitempty"`
	DefaultBranch string `json:"defaultBranch,omitempty"`
}

func (p *AgyProject) GetPath() string {
	if p.ProjectResources == nil || len(p.ProjectResources.Resources) == 0 {
		return ""
	}

	gf := p.ProjectResources.Resources[0].GitFolder
	if gf == nil {
		return ""
	}

	return parseFolderUri(gf.FolderURI)
}

func (p *AgyProject) GetBranch() string {
	if p.ProjectResources == nil || len(p.ProjectResources.Resources) == 0 {
		return ""
	}

	gf := p.ProjectResources.Resources[0].GitFolder
	if gf == nil {
		return ""
	}

	return gf.DefaultBranch
}

func parseFolderUri(rawUri string) string {
	if rawUri == "" {
		return ""
	}

	cleanUri := strings.TrimPrefix(rawUri, "file:///")
	cleanUri = strings.TrimPrefix(cleanUri, "file://")
	decoded, err := url.PathUnescape(cleanUri)
	if err != nil {
		decoded = cleanUri
	}

	if len(decoded) >= 2 && decoded[1] == ':' {
		return filepath.FromSlash(decoded)
	}

	if strings.HasPrefix(decoded, "/") {
		return filepath.FromSlash(decoded)
	}

	if strings.HasPrefix(rawUri, "file:///") {
		return filepath.FromSlash("/" + decoded)
	}

	return filepath.FromSlash(decoded)
}

func shortProjectId(id string) string {
	if len(id) <= 8 {
		return id
	}

	return id[:8]
}

func formatRelativeTime(rfc3339Str string) string {
	if rfc3339Str == "" {
		return "—"
	}

	t, isTimestampValid := parseTimestampString(rfc3339Str)
	if !isTimestampValid {
		return "—"
	}

	return formatTimeDuration(time.Since(t), t)
}

func parseTimestampString(s string) (time.Time, bool) {
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, true
	}

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}

	return time.Time{}, false
}

func formatTimeDuration(d time.Duration, t time.Time) string {
	if d < time.Minute {
		return "just now"
	}

	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}

	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}

	days := int(d.Hours() / 24)
	if days < 30 {
		return fmt.Sprintf("%dd ago", days)
	}

	return t.Format("2006-01-02")
}

// AGYRestoreSnapshot contains serialized metadata of projects and conversations for full restore.
type AGYRestoreSnapshot struct {
	CreatedAt     time.Time            `json:"createdAt"`
	Timestamp     int64                `json:"timestamp"`
	FilePath      string               `json:"filePath,omitempty"`
	TotalProjects int                  `json:"totalProjects"`
	TotalConvs    int                  `json:"totalConversations"`
	Projects      []AGYSnapshotProject `json:"projects"`
	Conversations []AGYSnapshotConv    `json:"conversations"`
}

// AGYSnapshotProject records a single registered workspace.
type AGYSnapshotProject struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
	Branch    string `json:"branch,omitempty"`
}

// AGYSnapshotConv records conversation identifiers and titles.
type AGYSnapshotConv struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	ProjectID     string `json:"projectId,omitempty"`
	WorkspacePath string `json:"workspacePath,omitempty"`
	StepCount     int    `json:"stepCount,omitempty"`
}

// AgyInstanceInfo encapsulates metadata for a running or configured Antigravity instance.
type AgyInstanceInfo struct {
	InstanceID       string   `json:"instanceId"`
	InstanceName     string   `json:"instanceName"`
	ProcessID        int      `json:"processId"`
	LanguageServer   string   `json:"languageServer"`
	ConfigDir        string   `json:"configDir"`
	BrainDir         string   `json:"brainDir"`
	SummariesDBPath  string   `json:"summariesDbPath"`
	ActiveWorkspaces []string `json:"activeWorkspaces"`
	IsPrimary        bool     `json:"isPrimary"`
	IsRunning        bool     `json:"isRunning"`
	LastActiveAt     string   `json:"lastActiveAt"`
}

// AgyInstancePromptQueryOptions provides filtering criteria for multi-instance prompt retrieval.
type AgyInstancePromptQueryOptions struct {
	InstanceID   string `json:"instanceId,omitempty"`
	IsAll        bool   `json:"isAll"`
	Status       string `json:"status,omitempty"` // "running", "queued", "all"
	Limit        int    `json:"limit"`
	MaxWords     int    `json:"maxWords"`
	IncludeConvs bool   `json:"includeConvs"`
}

// AgyMultiInstancePromptResponse provides the top-level API envelope.
type AgyMultiInstancePromptResponse struct {
	IsSuccess      bool                       `json:"isSuccess"`
	TotalInstances int                        `json:"totalInstances"`
	TotalPrompts   int                        `json:"totalPrompts"`
	Instances      []AgyInstancePromptPayload `json:"instances"`
	Error          string                     `json:"error,omitempty"`
}

// AgyInstancePromptPayload groups prompt snapshots by instance.
type AgyInstancePromptPayload struct {
	InstanceID   string                   `json:"instanceId"`
	InstanceName string                   `json:"instanceName"`
	RunningCount int                      `json:"runningCount"`
	QueuedCount  int                      `json:"queuedCount"`
	Running      []AgyPromptSnapshotItem  `json:"running"`
	Queued       []AgyPromptQueueEntry    `json:"queued"`
	RecentConvs  []AgyConversationPreview `json:"recentConvs,omitempty"`
}

// AgyConversationPreview captures recent conversation dialogue details.
type AgyConversationPreview struct {
	ConversationID string `json:"conversationId"`
	Title          string `json:"title"`
	WorkspacePath  string `json:"workspacePath"`
	LastPromptText string `json:"lastPromptText"`
	StepCount      int    `json:"stepCount"`
	UpdatedAt      string `json:"updatedAt"`
	IsActive       bool   `json:"isActive"`
}
