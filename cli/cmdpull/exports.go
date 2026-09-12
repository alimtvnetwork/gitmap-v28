package cmdpull

import "github.com/alimtvnetwork/gitmap-v28/cli/model"

// RunPull handles the "pull" subcommand.
func RunPull(args []string) error {
	return runPull(args)
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
	return runPush(args)
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
