package model

// ScanFolder is a root path that `gitmap scan` was invoked against.
// Repos record the ScanFolderId of their most recent scan via the
// nullable Repo.ScanFolderId column.
type ScanFolder struct {
	ID            int64  `json:"id"`
	AbsolutePath  string `json:"absolutePath"`
	Label         string `json:"label"`
	Notes         string `json:"notes"`
	LastScannedAt string `json:"lastScannedAt"`
	CreatedAt     string `json:"createdAt"`
}

// VersionProbe stores the result of a single HEAD-then-clone version
// probe for a repo. Empty in Phase 2.1; populated from Phase 2.3 onward.
type VersionProbe struct {
	ID             int64  `json:"id"`
	RepoID         int64  `json:"repoId"`
	ProbedAt       string `json:"probedAt"`
	NextVersionTag string `json:"nextVersionTag"`
	NextVersionNum int64  `json:"nextVersionNum"`
	Method         string `json:"method"`
	IsAvailable    bool   `json:"isAvailable"`
	Error          string `json:"error"`
}
