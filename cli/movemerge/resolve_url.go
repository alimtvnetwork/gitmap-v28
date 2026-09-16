package movemerge

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// resolveURLEndpoint clones or reuses the working folder for a URL.
func resolveURLEndpoint(ep Endpoint, opts Options) (Endpoint, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return ep, fmt.Errorf("getwd: %w", err)
	}

	dir := MapURLToFolder(cwd, ep.URL)
	ep.WorkingDir = dir
	isFound, err := IsFolderExisting(dir)
	if err != nil {
		return ep, err
	}

	if isFound {
		return reuseExistingURLFolder(ep, dir, opts)
	}

	if cloneErr := CloneURL(ep.URL, ep.Branch, dir); cloneErr != nil {
		return ep, cloneErr
	}

	ep.IsGitRepo = true

	return ep, nil
}

// reuseExistingURLFolder verifies the folder's origin matches and pulls.
func reuseExistingURLFolder(ep Endpoint, dir string, opts Options) (Endpoint, error) {
	origin, err := GetOriginURL(dir)
	if err != nil {
		return ep, fmt.Errorf("read origin in %s: %w", dir, err)
	}

	if IsOriginMatching(origin, ep.URL) {
		return handleMatchingOrigin(ep, dir)
	}

	if opts.IsForceFolder {
		return handleForceFolder(ep, dir)
	}

	return ep, fmt.Errorf(constants.ErrMMOriginFmt, dir, origin, ep.URL)
}

func handleMatchingOrigin(ep Endpoint, dir string) (Endpoint, error) {
	if pullErr := PullFFOnly(dir); pullErr != nil {
		return ep, pullErr
	}

	ep.IsGitRepo, ep.IsExisted = true, true

	return ep, nil
}

func handleForceFolder(ep Endpoint, dir string) (Endpoint, error) {
	if rmErr := os.RemoveAll(dir); rmErr != nil {
		return ep, fmt.Errorf("force-folder remove %s: %w", dir, rmErr)
	}

	if cloneErr := CloneURL(ep.URL, ep.Branch, dir); cloneErr != nil {
		return ep, cloneErr
	}

	ep.IsGitRepo, ep.IsExisted = true, false

	return ep, nil
}

// IsOriginMatching compares two URLs ignoring trailing .git and case.
func IsOriginMatching(a, b string) bool {
	norm := func(s string) string {
		return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(s)), ".git")
	}

	return norm(a) == norm(b)
}
