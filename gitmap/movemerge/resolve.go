package movemerge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ResolveEndpoint takes a raw CLI token and returns a fully resolved
// Endpoint with WorkingDir on disk. URL endpoints clone or reuse;
// folder endpoints validate existence per side semantics.
func ResolveEndpoint(raw string, isLeft bool, opts Options) (Endpoint, error) {
	kind, url, branch, display := ClassifyEndpoint(raw)
	ep := Endpoint{Raw: raw, DisplayName: display, Kind: kind, URL: url, Branch: branch}
	if kind == EndpointURL {
		return resolveURLEndpoint(ep, opts)
	}

	return resolveFolderEndpoint(ep, isLeft, opts)
}

// resolveFolderEndpoint validates a folder path per LEFT/RIGHT rules.
func resolveFolderEndpoint(ep Endpoint, isLeft bool, opts Options) (Endpoint, error) {
	abs, err := filepath.Abs(ep.DisplayName)
	if err != nil {
		return ep, fmt.Errorf("abs %s: %w", ep.DisplayName, err)
	}
	ep.WorkingDir = abs
	exists, err := FolderExists(abs)
	if err != nil {
		return ep, err
	}
	ep.Existed = exists
	if !exists && isLeft {
		return ep, fmt.Errorf(constants.ErrMMSrcMissingFmt, ep.DisplayName)
	}
	if exists && opts.PullFolder && IsGitRepo(abs) {
		if pullErr := PullFFOnly(abs); pullErr != nil {
			return ep, pullErr
		}
	}
	ep.IsGitRepo = IsGitRepo(abs)

	return ep, nil
}

// resolveURLEndpoint clones or reuses the working folder for a URL.
func resolveURLEndpoint(ep Endpoint, opts Options) (Endpoint, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return ep, fmt.Errorf("getwd: %w", err)
	}
	dir := MapURLToFolder(cwd, ep.URL)
	ep.WorkingDir = dir
	exists, err := FolderExists(dir)
	if err != nil {
		return ep, err
	}
	if !exists {
		if cloneErr := CloneURL(ep.URL, ep.Branch, dir); cloneErr != nil {
			return ep, cloneErr
		}
		ep.IsGitRepo = true

		return ep, nil
	}

	return reuseExistingURLFolder(ep, dir, opts)
}

// reuseExistingURLFolder verifies the folder's origin matches and pulls.
func reuseExistingURLFolder(ep Endpoint, dir string, opts Options) (Endpoint, error) {
	origin, err := GetOriginURL(dir)
	if err != nil {
		return ep, fmt.Errorf("read origin in %s: %w", dir, err)
	}
	if !originMatches(origin, ep.URL) {
		if !opts.ForceFolder {
			return ep, fmt.Errorf(constants.ErrMMOriginFmt, dir, origin, ep.URL)
		}
		if rmErr := os.RemoveAll(dir); rmErr != nil {
			return ep, fmt.Errorf("force-folder remove %s: %w", dir, rmErr)
		}
		if cloneErr := CloneURL(ep.URL, ep.Branch, dir); cloneErr != nil {
			return ep, cloneErr
		}
		ep.IsGitRepo, ep.Existed = true, false

		return ep, nil
	}
	if pullErr := PullFFOnly(dir); pullErr != nil {
		return ep, pullErr
	}
	ep.IsGitRepo, ep.Existed = true, true

	return ep, nil
}

// originMatches compares two URLs ignoring trailing .git and case.
func originMatches(a, b string) bool {
	norm := func(s string) string {
		return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(s)), ".git")
	}

	return norm(a) == norm(b)
}
