package cmdrm

import (
	"os"
	"time"
)

// RmFileRecord represents an individual file staged for removal.
type RmFileRecord struct {
	RelativePath string      `json:"relativePath"`
	StagedPath   string      `json:"stagedPath"`
	SizeBytes    int64       `json:"sizeBytes"`
	Sha256       string      `json:"sha256"`
	FileMode     os.FileMode `json:"fileMode"`
	DeletedAt    time.Time   `json:"deletedAt"`
}

// RmManifest represents the complete metadata for a task-based removal session.
type RmManifest struct {
	TaskId         string         `json:"taskId"`
	CreatedAt      time.Time      `json:"createdAt"`
	RepoRoot       string         `json:"repoRoot"`
	Reason         string         `json:"reason,omitempty"`
	FileCount      int            `json:"fileCount"`
	TotalSizeBytes int64          `json:"totalSizeBytes"`
	Files          []RmFileRecord `json:"files"`
}

// RmOptions configures the safe removal execution.
type RmOptions struct {
	Patterns []string
	TaskId   string
	Reason   string
	DryRun   bool
	Force    bool
	RepoRoot string
}

// RmUndoOptions configures restoration of a previous removal session.
type RmUndoOptions struct {
	TaskId   string
	RepoRoot string
}

// RmPurgeOptions configures manual purging of OS temporary backups.
type RmPurgeOptions struct {
	All       bool
	OlderThan time.Duration
	TaskId    string
}
