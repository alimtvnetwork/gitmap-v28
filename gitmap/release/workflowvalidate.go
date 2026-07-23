package release

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// resolveBump reads latest.json or falls back to git tags, then increments.
func resolveBump(level string) (Version, error) {
	current, err := resolveLatestVersion()
	if err != nil {
		return Version{}, err
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("version: current baseline: %s", current.String())
	}

	bumped, err := Bump(current, level)
	if err != nil {
		return Version{}, err
	}

	fmt.Printf(constants.MsgReleaseBumpResult, current.String(), bumped.String())

	return bumped, nil
}

// resolveLatestVersion tries latest.json first, then falls back to git tags.
func resolveLatestVersion() (Version, error) {
	return ResolveLatestVersion()
}

// ResolveLatestVersion is the public, single source of truth for "what
// version is the repo currently pinned at?". Read order:
//
//  1. .gitmap/release/latest.json — the canonical pin written by every
//     successful release (see WriteLatest).
//  2. Local git tags matching constants.GitTagGlob — the highest stable
//     semver wins. Used for fresh clones / repos where the JSON has been
//     deleted but the tag history is intact.
//
// Returns the parsed Version on success, or an error describing which
// sources were tried. Used by:
//   - release.resolveBump (the `gitmap release --bump <level>` path)
//   - cmd.peekNextMinorVersion (the bare `gitmap release` auto-bump
//     prompt + the right-click "Release next" context-menu entry)
//   - any future caller that needs the same source-of-truth ordering.
func ResolveLatestVersion() (Version, error) {
	latest, err := ReadLatest()
	if err == nil {
		v, parseErr := Parse(latest.Tag)
		if parseErr == nil {
			if verbose.IsEnabled() {
				verbose.Get().Log("version: baseline from latest.json: %s", v.String())
			}
			return v, nil
		}
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("version: latest.json unavailable, falling back to git tags")
	}

	return latestFromGitTags()
}

// resolveFromFile reads version.json.
func resolveFromFile() (Version, error) {
	raw, err := ReadVersionFile()
	if err != nil {
		return Version{}, fmt.Errorf(constants.ErrReleaseVersionRequired)
	}

	fmt.Printf(constants.MsgReleaseVersionRead, constants.DefaultVersionFile, raw)

	return Parse(raw)
}

// checkDuplicate verifies the version hasn't been released.
// If a release JSON exists but no tag or branch, prompts to remove it.
func checkDuplicate(v Version) error {
	if ReleaseExists(v) {
		tagExists := TagExistsLocally(v.String()) || TagExistsRemote(v.String())
		branchName := constants.ReleaseBranchPrefix + v.String()
		branchExists := BranchExists(branchName)

		if !tagExists && !branchExists {
			return handleOrphanedMeta(v)
		}

		return fmt.Errorf(constants.ErrReleaseAlreadyExists, v.String(), v.String())
	}
	if TagExistsLocally(v.String()) || TagExistsRemote(v.String()) {
		return fmt.Errorf(constants.ErrReleaseTagExists, v.String())
	}

	return nil
}

// handleOrphanedMeta detects a release JSON with no matching tag/branch
// and prompts the user to remove it before proceeding.
func handleOrphanedMeta(v Version) error {
	fmt.Printf(constants.MsgReleaseOrphanedMeta, v.String())
	fmt.Print(constants.MsgReleaseOrphanedPrompt)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return fmt.Errorf(constants.ErrReleaseAlreadyExists, v.String(), v.String())
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer != "y" && answer != "yes" {
		return fmt.Errorf(constants.ErrReleaseAborted)
	}

	filename := v.String() + constants.ExtJSON
	path := filepath.Join(constants.DefaultReleaseDir, filename)

	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf(constants.ErrReleaseOrphanedRemove, path, err)
	}

	fmt.Printf(constants.MsgReleaseOrphanedRemoved, v.String())

	return nil
}
