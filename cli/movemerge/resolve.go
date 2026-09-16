package movemerge

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	isFound, err := IsFolderExisting(abs)
	if err != nil {
		return ep, err
	}

	ep.IsExisted = isFound
	if !isFound && isLeft {
		return ep, fmt.Errorf(constants.ErrMMSrcMissingFmt, ep.DisplayName)
	}

	pullErr := checkAndPullFolder(abs, isFound, opts.IsPullFolder)
	if pullErr != nil {
		return ep, pullErr
	}

	ep.IsGitRepo = IsGitRepo(abs)

	return ep, nil
}

func checkAndPullFolder(abs string, isFound, isPullFolder bool) error {
	if isFound && isPullFolder && IsGitRepo(abs) {
		return PullFFOnly(abs)
	}

	return nil
}
