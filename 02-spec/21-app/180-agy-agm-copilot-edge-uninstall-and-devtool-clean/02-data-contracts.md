# Spec 180: Data Contracts, Schemas & API Envelopes

Spec Reference: [01-overview.md](01-overview.md)

---

## 1. AGY Snapshot & Restore Contracts

Before purging Antigravity from the machine, the system serializes a complete state snapshot:

```go
package cmdagy

import "time"

// AGYRestoreSnapshot captures project and conversation mappings.
type AGYRestoreSnapshot struct {
	Version       string                 `json:"version"`
	Timestamp     time.Time              `json:"timestamp"`
	HostOS        string                 `json:"host_os"`
	SnapshotPath  string                 `json:"snapshot_path"`
	Projects      []AGYProjectEntry      `json:"projects"`
	Conversations []AGYConversationEntry `json:"conversations"`
	Metadata      map[string]string      `json:"metadata,omitempty"`
}

// AGYProjectEntry captures a registered project workspace.
type AGYProjectEntry struct {
	WorkspaceURI string `json:"workspace_uri"`
	CorpusName   string `json:"corpus_name"`
	BasePath     string `json:"base_path"`
	ActiveBranch string `json:"active_branch,omitempty"`
}

// AGYConversationEntry captures an active conversation record.
type AGYConversationEntry struct {
	ConversationId string    `json:"conversation_id"`
	Title          string    `json:"title"`
	WorkspaceURI   string    `json:"workspace_uri"`
	LastModified   time.Time `json:"last_modified"`
	MessageCount   int       `json:"message_count"`
}
```

---

## 2. DevTool Clean Contracts

```go
package osclean

// DevToolCategory represents a distinct developer tool cache domain.
type DevToolCategory string

const (
	CategoryGoBuild     DevToolCategory = "go-build"
	CategoryGoMod       DevToolCategory = "go-mod"
	CategoryNodeModules DevToolCategory = "node-cache"
	CategoryPip         DevToolCategory = "pip-cache"
	CategoryVite        DevToolCategory = "vite-cache"
	CategoryAntigravity DevToolCategory = "antigravity"
	CategorySystemTemp  DevToolCategory = "system-temp"
	CategoryVSCode      DevToolCategory = "vscode-cache"
	CategoryChromeDev   DevToolCategory = "chrome-dev"
	CategoryGitCaches   DevToolCategory = "git-caches"
)

// DevToolCategoryResult contains execution metrics for a category.
type DevToolCategoryResult struct {
	Category    DevToolCategory `json:"category"`
	Description string          `json:"description"`
	BytesFreed  int64           `json:"bytes_freed"`
	ItemCount   int             `json:"item_count"`
	IsSuccess   bool            `json:"is_success"`
	Error       string          `json:"error,omitempty"`
}

// DevToolCleanSummary aggregates all results.
type DevToolCleanSummary struct {
	TotalBytesFreed int64                   `json:"total_bytes_freed"`
	TotalItemsClean int                     `json:"total_items_clean"`
	DurationMs      int64                   `json:"duration_ms"`
	Categories      []DevToolCategoryResult `json:"categories"`
}
```

---

## 3. Windows Copilot & Edge Removal Contracts

```go
package cmdwinutil

// WinRemovalResult captures the outcome of component de-installation.
type WinRemovalResult struct {
	Component   string   `json:"component"`
	AppxRemoved []string `json:"appx_removed"`
	RegKeysSet  []string `json:"reg_keys_set"`
	FilesPurged []string `json:"files_purged"`
	IsSuccess   bool     `json:"is_success"`
	Message     string   `json:"message"`
}
```
