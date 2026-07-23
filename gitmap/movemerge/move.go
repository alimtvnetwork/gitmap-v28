package movemerge

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// RunMove executes `gitmap mv LEFT RIGHT`: copy LEFT into RIGHT,
// then delete LEFT entirely.
func RunMove(left, right Endpoint, opts Options) error {
	if err := GuardEndpoints(left, right); err != nil {
		return err
	}
	if err := ensureRightExists(right, opts); err != nil {
		return err
	}
	logf(opts.LogPrefix, "copying files LEFT -> RIGHT (excluding .git/) ...")
	count, err := copyOrDryRun(left.WorkingDir, right.WorkingDir, opts)
	if err != nil {
		return err
	}
	logIndent(opts.LogPrefix, "copied %d files", count)
	if delErr := deleteLeftFolder(left, opts); delErr != nil {
		return delErr
	}

	return finalizeURLSides(left, right, DirRightOnly, opts)
}

// ensureRightExists creates RIGHT when it's a folder endpoint that
// didn't exist; honors --init.
func ensureRightExists(right Endpoint, opts Options) error {
	if right.Existed || right.Kind == EndpointURL {
		return nil
	}
	if opts.DryRun {
		logIndent(opts.LogPrefix, "[dry-run] mkdir %s", right.WorkingDir)

		return nil
	}
	if err := os.MkdirAll(right.WorkingDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", right.WorkingDir, err)
	}
	if opts.InitNewRight {
		if _, err := runGit(right.WorkingDir, "init"); err != nil {
			return fmt.Errorf("git init %s: %w", right.WorkingDir, err)
		}
	}

	return nil
}

// copyOrDryRun honors --dry-run while still indexing the source.
func copyOrDryRun(src, dst string, opts Options) (int, error) {
	if !opts.DryRun {
		return CopyTree(src, dst, opts)
	}
	idx, err := IndexTree(src, opts)
	if err != nil {
		return 0, err
	}
	for rel := range idx {
		logIndent(opts.LogPrefix, "[dry-run] copy %s", rel)
	}

	return len(idx), nil
}

// deleteLeftFolder removes LEFT recursively (mv semantic).
func deleteLeftFolder(left Endpoint, opts Options) error {
	logf(opts.LogPrefix, "deleting LEFT (%s) ...", left.DisplayName)
	if opts.DryRun {
		logIndent(opts.LogPrefix, "[dry-run] rm -rf %s", left.WorkingDir)

		return nil
	}
	if err := os.RemoveAll(left.WorkingDir); err != nil {
		return fmt.Errorf("delete LEFT %s: %w", left.WorkingDir, err)
	}
	logIndent(opts.LogPrefix, "deleted")

	return nil
}

// silence unused-import warning for filepath under MSVC-style builds.
var _ = filepath.Join
var _ = constants.GitBin
