// Package diff implements `gitmap diff LEFT RIGHT`: a read-only
// preview of what `gitmap merge-both / merge-left / merge-right`
// would change. No files are written, no commits, no pushes.
//
// Spec: companion to spec/01-app/97-move-and-merge.md
package diff

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Endpoint is a fully classified LEFT or RIGHT argument.
//
// `gitmap diff` only supports local folders. URL endpoints are
// rejected with a hint to clone them first; this keeps `diff`
// strictly read-only and side-effect-free.
type Endpoint struct {
	Raw         string
	DisplayName string
	WorkingDir  string
}

// ResolveEndpoint validates that raw points at an existing folder
// and returns its absolute path.
func ResolveEndpoint(raw string) (Endpoint, error) {
	trimmed := strings.TrimRight(raw, "/\\")
	ep := Endpoint{Raw: raw, DisplayName: trimmed}

	if isURLLike(trimmed) {
		return ep, fmt.Errorf(constants.ErrDiffNotFolderFmt, trimmed)
	}
	abs, err := filepath.Abs(trimmed)
	if err != nil {
		return ep, fmt.Errorf("resolve abs path %s: %w", trimmed, err)
	}
	ep.WorkingDir = abs

	return validateFolder(ep)
}

// validateFolder enforces existence + directory-ness.
func validateFolder(ep Endpoint) (Endpoint, error) {
	info, err := os.Stat(ep.WorkingDir)
	if os.IsNotExist(err) {
		return ep, fmt.Errorf(constants.ErrDiffMissingFmt, ep.DisplayName)
	}
	if err != nil {
		return ep, fmt.Errorf("stat %s: %w", ep.WorkingDir, err)
	}
	if !info.IsDir() {
		return ep, fmt.Errorf(constants.ErrDiffNotFolderFmt, ep.DisplayName)
	}

	return ep, nil
}

// isURLLike returns true when raw smells like a remote git URL.
func isURLLike(raw string) bool {
	lower := strings.ToLower(raw)
	for _, p := range []string{"https://", "http://", "ssh://", "git@"} {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}

	return false
}
