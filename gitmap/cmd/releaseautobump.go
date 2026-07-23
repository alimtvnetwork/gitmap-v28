package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// peekNextMinorVersion reads .gitmap/release/latest.json (falling back to
// the highest local git tag) and returns the auto-incremented MINOR version.
// Returns ok=false when no prior release exists in either source.
func peekNextMinorVersion() (current, next release.Version, ok bool) {
	cur, err := readCurrentReleaseVersion()
	if err != nil {
		return release.Version{}, release.Version{}, false
	}

	bumped, bumpErr := release.Bump(cur, constants.BumpMinor)
	if bumpErr != nil {
		return release.Version{}, release.Version{}, false
	}

	return cur, bumped, true
}

// readCurrentReleaseVersion returns the repo's current pinned version
// using the same JSON-first / git-tag-fallback resolver the explicit
// `--bump` path uses. This keeps every release entry-point (release,
// release --bump, release-pending, release-pull, and the right-click
// "Release next" menu item) reading from one single source of truth.
func readCurrentReleaseVersion() (release.Version, error) {
	return release.ResolveLatestVersion()
}

// confirmAutoBump prompts the user to confirm an auto-bump.
// Returns true when the user accepts (or when -y was passed).
func confirmAutoBump(current, next release.Version, yes bool) bool {
	fmt.Printf(constants.MsgReleaseAutoBumpHeader, current.String(), next.String())
	if yes {
		fmt.Print(constants.MsgReleaseAutoBumpYes)

		return true
	}

	fmt.Print(constants.MsgReleaseAutoBumpPrompt)

	return readYesNo()
}

// readYesNo reads one line from stdin and returns true for "y"/"yes".
func readYesNo() bool {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))

	return answer == "y" || answer == "yes"
}

// shouldAutoBumpMinor reports whether the bare-release auto-bump path should
// fire: no explicit version, no --bump, no commit/branch override.
func shouldAutoBumpMinor(version, bump, commit, branch string) bool {
	if len(version) > 0 || len(bump) > 0 {
		return false
	}
	if len(commit) > 0 || len(branch) > 0 {
		return false
	}

	return true
}
