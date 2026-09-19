package cmdautomation

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// CleanArtifactsOptions configures build artifact and cache purges.
type CleanArtifactsOptions struct {
	Dir             string   `json:"dir"`
	Targets         []string `json:"targets,omitempty"`
	Exclude         []string `json:"exclude,omitempty"`
	IsDryRun        bool     `json:"isDryRun"`
	IsAll           bool     `json:"isAll"`
	IsVerbose       bool     `json:"isVerbose"`
	IsJson          bool     `json:"isJson"`
	IsCleanPycache  bool     `json:"isCleanPycache"`
	IsCleanTemp     bool     `json:"isCleanTemp"`
	IsCleanBinaries bool     `json:"isCleanBinaries"`
	IsPermanent     bool     `json:"isPermanent"`
}

// CleanArtifactItem records an identified or deleted file or folder.
type CleanArtifactItem struct {
	Path         string `json:"path"`
	SizeBytes    int64  `json:"sizeBytes"`
	Category     string `json:"category"`
	IsDeleted    bool   `json:"isDeleted"`
	IsGitTracked bool   `json:"isGitTracked"`
}

// CleanArtifactsResult aggregates the outcome of artifact cleanup.
type CleanArtifactsResult struct {
	TotalFound   int                 `json:"totalFound"`
	TotalDeleted int                 `json:"totalDeleted"`
	FreedBytes   int64               `json:"freedBytes"`
	FreedMB      float64             `json:"freedMB"`
	Items        []CleanArtifactItem `json:"items"`
	Duration     time.Duration       `json:"duration"`
	IsSuccess    bool                `json:"isSuccess"`
}

// ChangedFilesOptions configures git changed files extraction.
type ChangedFilesOptions struct {
	Dir          string `json:"dir"`
	BaseRef      string `json:"baseRef"`
	Commits      int    `json:"commits"`
	IsStagedOnly bool   `json:"isStagedOnly"`
	IsVerify     bool   `json:"isVerify"`
	IsJson       bool   `json:"isJson"`
}

// ChangedFileItem records a single changed file from git history or tree.
type ChangedFileItem struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	Extension string `json:"extension"`
	IsStaged  bool   `json:"isStaged"`
	IsExists  bool   `json:"isExists"`
}

// ChangedFilesResult stores git changed and staged files.
type ChangedFilesResult struct {
	BaseRef    string            `json:"baseRef"`
	HeadHash   string            `json:"headHash"`
	TotalFiles int               `json:"totalFiles"`
	Files      []ChangedFileItem `json:"files"`
	Duration   time.Duration     `json:"duration"`
	IsSuccess  bool              `json:"isSuccess"`
}

// PurgeHistoryOptions configures historical blob identification and purging.
type PurgeHistoryOptions struct {
	Dir           string  `json:"dir"`
	MinSizeMB     float64 `json:"minSizeMb"`
	PathPattern   string  `json:"pathPattern"`
	IsDryRun      bool    `json:"isDryRun"`
	IsJson        bool    `json:"isJson"`
	IsAutoConfirm bool    `json:"isAutoConfirm"`
}

// PurgeHistoryItem represents a large or deleted historical git blob.
type PurgeHistoryItem struct {
	Path       string  `json:"path"`
	CommitHash string  `json:"commitHash"`
	SizeBytes  int64   `json:"sizeBytes"`
	SizeMB     float64 `json:"sizeMb"`
	Author     string  `json:"author,omitempty"`
	Date       string  `json:"date,omitempty"`
	IsPurged   bool    `json:"isPurged"`
}

// PurgeHistoryResult captures the findings and actions of history purging.
type PurgeHistoryResult struct {
	TotalFound  int                `json:"totalFound"`
	TotalPurged int                `json:"totalPurged"`
	TotalBytes  int64              `json:"totalBytes"`
	TotalMB     float64            `json:"totalMb"`
	Items       []PurgeHistoryItem `json:"items"`
	Duration    time.Duration      `json:"duration"`
	IsSuccess   bool               `json:"isSuccess"`
}

// FormatGoOptions configures Go AST code formatting and import organization.
type FormatGoOptions struct {
	Dir      string   `json:"dir"`
	Paths    []string `json:"paths,omitempty"`
	IsWrite  bool     `json:"isWrite"`
	IsCheck  bool     `json:"isCheck"`
	IsStaged bool     `json:"isStaged"`
	IsJson   bool     `json:"isJson"`
}

// FormatGoItem records the formatting status of a single Go file.
type FormatGoItem struct {
	Path               string `json:"path"`
	IsFormatted        bool   `json:"isFormatted"`
	HasBom             bool   `json:"hasBom"`
	HasCrlf            bool   `json:"hasCrlf"`
	HasReorderImports bool   `json:"hasReorderImports"`
	ErrorMessage       string `json:"errorMessage,omitempty"`
}

// FormatGoResult aggregates outcomes of Go formatting passes.
type FormatGoResult struct {
	TotalFiles     int            `json:"totalFiles"`
	FormattedCount int            `json:"formattedCount"`
	ViolationCount int            `json:"violationCount"`
	Items          []FormatGoItem `json:"items"`
	Duration       time.Duration  `json:"duration"`
	IsClean        bool           `json:"isClean"`
}

type (
	// CleanArtifactsResultMonad wraps CleanArtifactsResult with AppError.
	CleanArtifactsResultMonad = result.Result[CleanArtifactsResult]

	// ChangedFilesResultMonad wraps ChangedFilesResult with AppError.
	ChangedFilesResultMonad = result.Result[ChangedFilesResult]

	// PurgeHistoryResultMonad wraps PurgeHistoryResult with AppError.
	PurgeHistoryResultMonad = result.Result[PurgeHistoryResult]

	// FormatGoResultMonad wraps FormatGoResult with AppError.
	FormatGoResultMonad = result.Result[FormatGoResult]
)
