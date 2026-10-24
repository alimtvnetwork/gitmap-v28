package cmdsync

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/enums"
)

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

// Mode returns the operational sync mode represented by these options.
func (o SyncOptions) Mode() enums.SyncModeType {
	if o.DryRun {
		return enums.SyncModeDryRun
	}

	return enums.SyncModeFastForward
}

// Target returns the target sync scope represented by these options.
func (o SyncOptions) Target() enums.SyncTargetType {
	if len(o.Repos) == 0 {
		return enums.SyncTargetAll
	}

	return enums.SyncTargetDirty
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

// OutcomeStatus returns the typed repo status for this sync result.
func (r RepoSyncResult) OutcomeStatus() enums.SyncRepoStatusType {
	if r.Status == "FAIL" {
		return enums.SyncRepoStatusError
	}
	if r.Status == "SKIPPED" {
		return enums.SyncRepoStatusSkipped
	}

	return enums.SyncRepoStatusSynced
}
