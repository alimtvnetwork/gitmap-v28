package store

// HistoryPurgeOperation represents a top-level history rewrite operation.
type HistoryPurgeOperation struct {
	OperationId          int64  `json:"operationId"`
	RepoSlug             string `json:"repoSlug"`
	TargetType           string `json:"targetType"` // "folder", "file", "commit", "commit-range"
	TargetPath           string `json:"targetPath"` // Relative path
	CommitListJson       string `json:"commitListJson"`
	TotalCommitsScanned  int    `json:"totalCommitsScanned"`
	AffectedCommitsCount int    `json:"affectedCommitsCount"`
	BackedUpFilesCount   int    `json:"backedUpFilesCount"`
	BackupVaultPath      string `json:"backupVaultPath"`
	IsDryRun             bool   `json:"isDryRun"`
	IsVerified           bool   `json:"isVerified"`
	HasPushed            bool   `json:"hasPushed"`
	IsUndone             bool   `json:"isUndone"`
	IsSuccess            bool   `json:"isSuccess"`
	CreatedAt            int64  `json:"createdAt"` // Unix epoch seconds
	CompletedAt          int64  `json:"completedAt"`
	Notes                string `json:"notes"`
}

// HistoryPurgeCommitMap records 1:1 mapping between pre-rewrite and post-rewrite commits.
type HistoryPurgeCommitMap struct {
	MapId              int64  `json:"mapId"`
	OperationId        int64  `json:"operationId"`
	OriginalCommitSha  string `json:"originalCommitSha"`
	RewrittenCommitSha string `json:"rewrittenCommitSha"`
	ParentOriginalSha  string `json:"parentOriginalSha"`
	ParentRewrittenSha string `json:"parentRewrittenSha"`
	AuthorName         string `json:"authorName"`
	AuthorEmail        string `json:"authorEmail"`
	CommitTimestamp    int64  `json:"commitTimestamp"`
	CommitMessage      string `json:"commitMessage"`
	IsPurged           bool   `json:"isPurged"`
	CreatedAt          int64  `json:"createdAt"`
}

// HistoryPurgeFile catalogs an excised blob and its temp backup location.
type HistoryPurgeFile struct {
	FileId            int64  `json:"fileId"`
	OperationId       int64  `json:"operationId"`
	OriginalCommitSha string `json:"originalCommitSha"`
	RelativePath      string `json:"relativePath"`
	BlobSha           string `json:"blobSha"`
	FileSizeBytes     int64  `json:"fileSizeBytes"`
	FileMode          uint32 `json:"fileMode"`
	BackupRelPath     string `json:"backupRelPath"`
	IsPurged          bool   `json:"isPurged"`
	CreatedAt         int64  `json:"createdAt"`
}

// HistoryUndoOperation records an execution of the undo engine.
type HistoryUndoOperation struct {
	UndoId              int64  `json:"undoId"`
	OperationId         int64  `json:"operationId"`
	RepoSlug            string `json:"repoSlug"`
	RestoredCommitCount int    `json:"restoredCommitCount"`
	RestoredFileCount   int    `json:"restoredFileCount"`
	BackupVaultPath     string `json:"backupVaultPath"`
	IsVerified          bool   `json:"isVerified"`
	IsSuccess           bool   `json:"isSuccess"`
	CreatedAt           int64  `json:"createdAt"`
	CompletedAt         int64  `json:"completedAt"`
	ErrorMessage        string `json:"errorMessage"`
}

// CreatePurgeOptions wraps parameters for starting a purge operation.
type CreatePurgeOptions struct {
	RepoSlug   string
	TargetType string
	TargetPath string
	CommitList []string
	IsDryRun   bool
	HasPushed  bool
}

// UpdatePurgeStatusOptions wraps completion flags for an operation.
type UpdatePurgeStatusOptions struct {
	OperationId int64
	IsSuccess   bool
	IsVerified  bool
	Notes       string
}
