package cmdpull

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspace"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// RunRemoteSSHPullFn delegates remote SSH pull all-efficient to cmdssh package.
var RunRemoteSSHPullFn func(target string, isTable bool) error

// RunRemoteSSHPullAllFleetFn delegates SSH fleet pull-all to cmdssh package.
var RunRemoteSSHPullAllFleetFn func(cleanArgs []string) error

// RunRemoteSSHPullFleetRemediateFn delegates multi-machine remediation to cmdssh package.
var RunRemoteSSHPullFleetRemediateFn func(cleanArgs []string) error

// RunPull handles the "pull" subcommand.
func RunPull(args []string) error {
	err := runPull(args)
	cmdspace.FireEmailTaskTriggers(cmdspace.EmailTaskPull, cmdspace.InferTaskOutcome(cmdspace.EmailTaskPull, err))
	return err
}

// RunPullAll is the explicit batch-pull entry point.
func RunPullAll(args []string) error {
	return runPullAll(args)
}

// RunPullReleaseCD implements `gitmap pull-release-cd` (alias `prc`).
func RunPullReleaseCD(args []string) error {
	return runPullReleaseCD(args)
}

// RunPush handles the "push" subcommand.
func RunPush(args []string) error {
	err := runPush(args)
	cmdspace.FireEmailTaskTriggers(cmdspace.EmailTaskPush, cmdspace.InferTaskOutcome(cmdspace.EmailTaskPush, err))
	return err
}

// FindBySlug finds records matching the slug (case-insensitive, partial match).
func FindBySlug(records []model.ScanRecord, slug string) []model.ScanRecord {
	return findBySlug(records, slug)
}

// IsGitRepoCWD returns true when cwd is inside a git work tree.
func IsGitRepoCWD() bool {
	return isGitRepoCWD()
}

// PullOneRepo pulls a single repository.
func PullOneRepo(rec model.ScanRecord) {
	pullOneRepo(rec)
}
