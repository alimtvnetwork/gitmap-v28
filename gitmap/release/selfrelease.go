package release

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ExecuteSelf resolves the gitmap source repo, switches to that directory,
// runs Execute, then returns to the original dir.
func ExecuteSelf(opts Options) error {
	srcRoot, err := resolveSourceRepo()
	if err != nil {
		return err
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not determine current directory: %w", err)
	}

	if filepath.Clean(srcRoot) == filepath.Clean(originalDir) {
		fmt.Printf(constants.MsgSelfReleaseSameDir, srcRoot)

		return Execute(opts)
	}

	fmt.Printf(constants.MsgSelfReleaseSwitch, srcRoot)
	if err := os.Chdir(srcRoot); err != nil {
		return fmt.Errorf("could not switch to source repo: %w", err)
	}

	releaseErr := Execute(opts)
	if cdErr := os.Chdir(originalDir); cdErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not return to %s: %v\n", originalDir, cdErr)
	} else {
		fmt.Printf(constants.MsgSelfReleaseReturn, originalDir)
	}

	return releaseErr
}

// IsInsideGitRepo checks if the current directory is inside a Git repository.
func IsInsideGitRepo() bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}

	return findGitRoot(dir) != ""
}
