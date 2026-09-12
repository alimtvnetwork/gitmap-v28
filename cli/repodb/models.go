package repodb

// RepoFile represents an indexed file tracked in the repository database.
type RepoFile struct {
	RepoFileId   uint64 `db:"RepoFileId" json:"repoFileId"`
	RelativePath string `db:"RelativePath" json:"relativePath"`
	AbsolutePath string `db:"AbsolutePath" json:"absolutePath"`
	Content      string `db:"Content" json:"content"`
	IsBig        bool   `db:"IsBig" json:"isBig"`
	WriteTime    int64  `db:"WriteTime" json:"writeTime"`
	CreatedAt    int64  `db:"CreatedAt" json:"createdAt"`
	UpdatedAt    int64  `db:"UpdatedAt" json:"updatedAt"`
}

// SearchCache represents a cached search query and its JSON payload.
type SearchCache struct {
	SearchCacheId uint64 `db:"SearchCacheId" json:"searchCacheId"`
	Query         string `db:"Query" json:"query"`
	Hits          int    `db:"Hits" json:"hits"`
	ResultJson    string `db:"ResultJson" json:"resultJson"`
	CreatedAt     int64  `db:"CreatedAt" json:"createdAt"`
	UpdatedAt     int64  `db:"UpdatedAt" json:"updatedAt"`
}

// FileSequence represents sequence ordering assigned to a file in a directory.
type FileSequence struct {
	FileSequenceId uint64 `db:"FileSequenceId" json:"fileSequenceId"`
	Directory      string `db:"Directory" json:"directory"`
	Filename       string `db:"Filename" json:"filename"`
	SequenceNumber int    `db:"SequenceNumber" json:"sequenceNumber"`
	BaseName       string `db:"BaseName" json:"baseName"`
	UpdatedAt      int64  `db:"UpdatedAt" json:"updatedAt"`
}

// SequenceHistory represents historical sequence operations performed on a directory.
type SequenceHistory struct {
	SequenceHistoryId uint64 `db:"SequenceHistoryId" json:"sequenceHistoryId"`
	Directory         string `db:"Directory" json:"directory"`
	OperationsJson    string `db:"OperationsJson" json:"operationsJson"`
	CreatedAt         int64  `db:"CreatedAt" json:"createdAt"`
}

// RepoScanLog represents an audit log entry for repository scans and sync activity.
type RepoScanLog struct {
	RepoScanLogId uint64 `db:"RepoScanLogId" json:"repoScanLogId"`
	RepoId        int64  `db:"RepoId" json:"repoId"`
	RepoSlug      string `db:"RepoSlug" json:"repoSlug"`
	Action        string `db:"Action" json:"action"`
	Status        string `db:"Status" json:"status"`
	ErrorMessage  string `db:"ErrorMessage" json:"errorMessage,omitempty"`
	Details       string `db:"Details" json:"details,omitempty"`
	Notes         string `db:"Notes" json:"notes,omitempty"`
	Comments      string `db:"Comments" json:"comments,omitempty"`
	CreatedAt     string `db:"CreatedAt" json:"createdAt"`
}

// IndexedRepo represents metadata about an indexed repository in root SplitDB.
type IndexedRepo struct {
	IndexedRepoId   uint64 `db:"IndexedRepoId" json:"indexedRepoId"`
	Path            string `db:"Path" json:"path"`
	Slug            string `db:"Slug" json:"slug"`
	MigratedVersion int    `db:"MigratedVersion" json:"migratedVersion"`
	CreatedAt       int64  `db:"CreatedAt" json:"createdAt"`
	UpdatedAt       int64  `db:"UpdatedAt" json:"updatedAt"`
}
