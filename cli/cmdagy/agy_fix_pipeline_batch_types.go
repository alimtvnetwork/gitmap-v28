package cmdagy

// PipelineFixBatchCursor tracks multi-project batch pagination state.
type PipelineFixBatchCursor struct {
	LastIndex      int      `json:"lastIndex"`
	BatchLimit     int      `json:"batchLimit"`
	ProcessedRepos []string `json:"processedRepos"`
	FailingRepos   []string `json:"failingRepos"`
	UpdatedAt      string   `json:"updatedAt"`
}

// ProjectFixCandidate represents a candidate repository scanned for pipeline failures.
type ProjectFixCandidate struct {
	Name        string
	Path        string
	RepoSlug    string
	HasFailures bool
	RunID       uint64
	SHA         string
	ErrorReport string
}
