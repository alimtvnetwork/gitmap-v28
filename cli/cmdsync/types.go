package cmdsync

// SyncOptions holds parsed command line arguments for gitmap sync.
type SyncOptions struct {
	Projects  string
	Repos     []string
	Workers   int
	DryRun    bool
	NoPush    bool
	NoRelease bool
	Verbose   bool
}

// ProjectConfig defines a single target project configuration.
type ProjectConfig struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	URL    string `json:"url"`
	Folder string `json:"folder"`
	Mode   string `json:"mode"`
}

// RepoSyncResult records the synchronization outcome for a single repository.
type RepoSyncResult struct {
	Repo         string `json:"repo"`
	Status       string `json:"status"` // OK, FAIL, SKIPPED
	PreTag       string `json:"preTag"`
	PostTag      string `json:"postTag"`
	FilesAdded   int    `json:"filesAdded"`
	FilesRemoved int    `json:"filesRemoved"`
	FilesUpdated int    `json:"filesUpdated"`
	Action       string `json:"action"`
	Error        string `json:"error,omitempty"`
}
