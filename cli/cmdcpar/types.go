package cmdcpar

import "github.com/alimtvnetwork/gitmap-v28/cli/gitutil"

type cparOptions struct {
	isAutoYes    bool
	isReview     bool
	isCommitOnly bool
	commitMsg    string
}

// DirtyRepoSummary captures uncommitted changes for review.
type DirtyRepoSummary struct {
	RepoName  string
	RepoPath  string
	Diagnosis gitutil.DirtyDiagnosis
}

// CPARRunState captures the outcome of a batch CPAR run.
type CPARRunState struct {
	HasFailures  bool
	SuccessCount int
	FailureCount int
}
