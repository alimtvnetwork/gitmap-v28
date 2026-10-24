package cmdcpar

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdcpar/enums"

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

// CPARReviewActionType represents user choice in review flow.
// CPARReviewActionType aliases the strongly-typed enum from the enums package.
type CPARReviewActionType = enums.CPARReviewActionType

const (
	ActionCommitAllFeature = enums.ActionCommitAllFeature
	ActionCommitAllBug     = enums.ActionCommitAllBug
	ActionSingleRepo       = enums.ActionSingleRepo
	ActionCommitAllChore   = enums.ActionCommitAllChore
	ActionAbort            = enums.ActionAbort
)

func parseCPAROptions(args []string) cparOptions {
	opts := cparOptions{commitMsg: "chore: commit pending changes"}
	for _, a := range args {
		opts = applyOption(opts, strings.ToLower(a))
	}
	return opts
}

func applyOption(opts cparOptions, a string) cparOptions {
	switch {
	case isYesFlag(a):
		opts.isAutoYes = true
	case isReviewFlag(a):
		opts.isReview = true
	case isCommitOnlyFlag(a):
		opts.isCommitOnly = true
	case isCompoundFlag(a):
		opts = applyCompoundFlags(opts, a)
	case strings.HasPrefix(a, "-m="):
		opts.commitMsg = a[3:]
	}
	return opts
}

func isYesFlag(a string) bool {
	return a == "-y" || a == "--yes"
}

func isReviewFlag(a string) bool {
	return a == "-r" || a == "--review"
}

func isCommitOnlyFlag(a string) bool {
	return a == "-co" || a == "--co" || a == "--commit-only"
}

func isCompoundFlag(a string) bool {
	return a == "-rco" || a == "-ry" || a == "-yr"
}

func applyCompoundFlags(opts cparOptions, a string) cparOptions {
	if a == "-rco" {
		opts.isReview = true
		opts.isCommitOnly = true
	} else if a == "-ry" || a == "-yr" {
		opts.isReview = true
		opts.isAutoYes = true
	}
	return opts
}
