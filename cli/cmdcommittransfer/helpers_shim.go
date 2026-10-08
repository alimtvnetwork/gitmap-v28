package cmdcommittransfer

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcreate"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"path/filepath"
	"os"
)

// EnsureOrProvisionDestinationRepo ensures a destination exists; if not, provisions it.
func EnsureOrProvisionDestinationRepo(rawTarget string, isLocal bool) (string, error) {
	if cmdcreate.IsRemoteGitURL(rawTarget) {
		return rawTarget, nil
	}

	abs, err := filepath.Abs(rawTarget)
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve destination:")
	}

	info, statErr := os.Stat(abs)
	isFound := statErr == nil && info.IsDir()
	if isFound {
		return abs, nil
	}

	return cmdcreate.ProvisionMissingDestination(rawTarget, abs, isLocal)
}

func argsTail() []string {
	if len(os.Args) > 2 {
		return os.Args[2:]
	}
	return nil
}

// runCommitIn is the top-level entry point for `gitmap commit-in` / `gitmap cin`.

// runCommitPull executes `gitmap commit-pull` / `cpull` / `pull-commits`.
// Replays commits chronologically while creating Pull Requests for all merges and releases.

