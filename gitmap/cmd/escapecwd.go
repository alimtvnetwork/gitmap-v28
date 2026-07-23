package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// escapeCwdIfInside chdirs out of `target` when the current working
// directory IS `target` or sits under it. This releases the directory
// handle that Windows holds on the cwd, allowing a subsequent
// os.RemoveAll(target) to succeed.
//
// Returns the directory we ended up in (parent of `target` when an
// escape happened; the original cwd otherwise). A non-nil error means
// we attempted to escape but Chdir failed — callers should treat that
// as fatal because the follow-up remove will deadlock on Windows.
//
// Used by spec/01-app/113 to make every clone-family command tolerate
// "I'm already inside the folder I'm about to re-clone".
func escapeCwdIfInside(target string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", nil
	}
	if isPathInside(cleanExistingPath(cwd), cleanExistingPath(target)) {
		return escapeCwdToParent(cwd, target)
	}

	return cwd, nil
}

func escapeCwdToParent(cwd, target string) (string, error) {
	tgtClean := cleanExistingPath(target)
	parent := filepath.Dir(tgtClean)
	fmt.Printf("↪ cwd is inside %q; chdir → %q to release handle\n",
		tgtClean, parent)

	if chErr := os.Chdir(parent); chErr != nil {
		return cwd, fmt.Errorf("escapeCwdIfInside: chdir %q: %w", parent, chErr)
	}

	return parent, nil
}

func cleanExistingPath(path string) string {
	cleaned := filepath.Clean(path)
	resolved, err := filepath.EvalSymlinks(cleaned)
	if err != nil {
		return cleaned
	}

	return filepath.Clean(resolved)
}

// isPathInside reports whether `child` equals `parent` or is a
// descendant. Case-insensitive on every OS (Windows safety; harmless
// on Linux/macOS for typical repo paths).
func isPathInside(child, parent string) bool {
	if strings.EqualFold(child, parent) {
		return true
	}

	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	if rel == "." {
		return true
	}

	return !strings.HasPrefix(rel, "..")
}
