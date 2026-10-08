package cmdai

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AiTaskStatusType defines valid task statuses.
type AiTaskStatusType string

const (
	AiTaskStatusInProgress AiTaskStatusType = "in_progress"
	AiTaskStatusCompleted  AiTaskStatusType = "completed"
	AiTaskStatusFailed     AiTaskStatusType = "failed"
	AiTaskStatusCancelled  AiTaskStatusType = "cancelled"
)

// AiOperationType defines the file interaction performed.
type AiOperationType string

const (
	AiOperationRead   AiOperationType = "read"
	AiOperationEdit   AiOperationType = "edit"
	AiOperationDelete AiOperationType = "delete"
)

// AiTaskRecord represents an AI analysis session record for JSON serialization and CLI display.
type AiTaskRecord struct {
	AiTaskId    int64  `json:"aiTaskId"`
	TaskUuid    string `json:"taskUuid"`
	Goal        string `json:"goal"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	Reasoning   string `json:"reasoning"`
	TotalFiles  int    `json:"totalFiles"`
	TotalLines  int    `json:"totalLines"`
	IsActive    bool   `json:"isActive"`
	HasFailed   bool   `json:"hasFailed"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
	Comments    string `json:"comments"`
	StartedAt   string `json:"startedAt"`
	CompletedAt string `json:"completedAt"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// AiTaskFileRecord represents an investigated, modified, or removed file.
type AiTaskFileRecord struct {
	AiTaskFileId int64  `json:"aiTaskFileId"`
	AiTaskId     int64  `json:"aiTaskId"`
	RelPath      string `json:"relPath"`
	AbsPath      string `json:"absPath"`
	Action       string `json:"action"`
	BeforeSha256 string `json:"beforeSha256"`
	AfterSha256  string `json:"afterSha256"`
	BackupPath   string `json:"backupPath"`
	LineCount    int    `json:"lineCount"`
	Reasoning    string `json:"reasoning"`
	IsModified   bool   `json:"isModified"`
	IsRemoved    bool   `json:"isRemoved"`
	Notes        string `json:"notes"`
	Comments     string `json:"comments"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// AiTaskLineRecord represents a granular line change, diff, and rule violation explanation.
type AiTaskLineRecord struct {
	AiTaskLineId    int64  `json:"aiTaskLineId"`
	AiTaskFileId    int64  `json:"aiTaskFileId"`
	LineNumber      int    `json:"lineNumber"`
	OriginalContent string `json:"originalContent"`
	ProposedContent string `json:"proposedContent"`
	DiffKind        string `json:"diffKind"`
	RuleViolation   string `json:"ruleViolation"`
	Reasoning       string `json:"reasoning"`
	IsApplied       bool   `json:"isApplied"`
	Notes           string `json:"notes"`
	Comments        string `json:"comments"`
	CreatedAt       string `json:"createdAt"`
}

// AiAnalysisSummary represents aggregated statistics across AI analysis tasks.
type AiAnalysisSummary struct {
	TotalTasks     int `json:"totalTasks"`
	ActiveTasks    int `json:"activeTasks"`
	CompletedTasks int `json:"completedTasks"`
	FailedTasks    int `json:"failedTasks"`
	TotalFiles     int `json:"totalFiles"`
	TotalLines     int `json:"totalLines"`
}

// AiAnalysisTask represents a high-level AI analysis session.
type AiAnalysisTask struct {
	TaskId             string           `json:"taskId" db:"task_id"`
	RepoPath           string           `json:"repoPath" db:"repo_path"`
	TaskDescription    string           `json:"taskDescription" db:"task_description"`
	ModelName          string           `json:"modelName" db:"model_name"`
	Status             AiTaskStatusType `json:"status" db:"status"`
	StartedAt          time.Time        `json:"startedAt" db:"started_at"`
	CompletedAt        *time.Time       `json:"completedAt,omitempty" db:"completed_at"`
	TotalFilesRead     int              `json:"totalFilesRead" db:"total_files_read"`
	TotalFilesModified int              `json:"totalFilesModified" db:"total_files_modified"`
	TotalFilesDeleted  int              `json:"totalFilesDeleted" db:"total_files_deleted"`
	ReasoningSummary   string           `json:"reasoningSummary" db:"reasoning_summary"`
	CreatedAt          time.Time        `json:"createdAt" db:"created_at"`
}

// AiAnalysisLine represents a granular file observation or mutation record.
type AiAnalysisLine struct {
	Id             int64           `json:"id" db:"id"`
	TaskId         string          `json:"taskId" db:"task_id"`
	FilePath       string          `json:"filePath" db:"file_path"`
	OperationType  AiOperationType `json:"operationType" db:"operation_type"`
	StartLine      int             `json:"startLine" db:"start_line"`
	EndLine        int             `json:"endLine" db:"end_line"`
	LineCount      int             `json:"lineCount" db:"line_count"`
	ContentSnippet string          `json:"contentSnippet" db:"content_snippet"`
	Reasoning      string          `json:"reasoning" db:"reasoning"`
	CreatedAt      time.Time       `json:"createdAt" db:"created_at"`
}

// StartTaskOptions configures task initiation.
type StartTaskOptions struct {
	TaskId      string
	Description string
	ModelName   string
	RepoRoot    string
}

// RecordLineOptions configures recording an individual line observation.
type RecordLineOptions struct {
	TaskId    string
	FilePath  string
	Lines     string
	Operation string
	Reasoning string
	Snippet   string
	RepoRoot  string
}

// ParseLineRange parses raw line range strings like "42" or "10-50".
func ParseLineRange(raw string) (int, int, error) {
	clean := strings.TrimSpace(raw)
	isEmpty := clean == ""
	if isEmpty {
		return 0, 0, nil
	}

	parts := strings.Split(clean, "-")
	isSingle := len(parts) == 1
	if isSingle {
		return parseSingleLine(parts[0], clean)
	}

	isRange := len(parts) == 2
	if isRange {
		return parsePairLineRange(parts[0], parts[1], clean)
	}

	return 0, 0, apperror.NewValidationError(fmt.Sprintf("malformed line range %q", clean))
}

func parseSingleLine(part, original string) (int, int, error) {
	num, err := strconv.Atoi(strings.TrimSpace(part))
	isValid := err == nil && num > 0
	if !isValid {
		return 0, 0, apperror.NewValidationError(fmt.Sprintf("invalid line number %q", original))
	}

	return num, num, nil
}

func parsePairLineRange(partStart, partEnd, original string) (int, int, error) {
	start, err1 := strconv.Atoi(strings.TrimSpace(partStart))
	end, err2 := strconv.Atoi(strings.TrimSpace(partEnd))
	isValid := err1 == nil && err2 == nil && start > 0 && end >= start
	if !isValid {
		return 0, 0, apperror.NewValidationError(fmt.Sprintf("invalid line range %q", original))
	}

	return start, end, nil
}
