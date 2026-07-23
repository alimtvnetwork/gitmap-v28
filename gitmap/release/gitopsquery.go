package release

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TagExistsLocally checks if a git tag exists in the local repo.
func TagExistsLocally(tag string) bool {
	cmd := exec.Command(constants.GitBin, constants.GitTag, constants.GitTagListFlag, tag)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// TagExistsRemote checks if a git tag exists on the remote.
func TagExistsRemote(tag string) bool {
	cmd := exec.Command(constants.GitBin,
		constants.GitLsRemote, constants.GitLsRemoteTags, constants.GitOrigin, tag)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// BranchExists checks if a local branch exists.
func BranchExists(branch string) bool {
	cmd := exec.Command(constants.GitBin, constants.GitBranch, constants.GitBranchListFlag, branch)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// CurrentCommitSHA returns the full SHA of HEAD.
func CurrentCommitSHA() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, constants.GitHEAD)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// CurrentBranchName returns the current branch name.
func CurrentBranchName() (string, error) {
	cmd := exec.Command(constants.GitBin,
		constants.GitRevParse, constants.GitAbbrevRef, constants.GitHEAD)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// CommitExists checks if a commit SHA is valid.
func CommitExists(sha string) bool {
	cmd := exec.Command(constants.GitBin, constants.GitCatFile, constants.GitCatFileTypeFlag, sha)
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == constants.GitCommitType
}

// latestFromGitTags scans all local git tags for the highest stable semver.
func latestFromGitTags() (Version, error) {
	cmd := exec.Command(constants.GitBin, constants.GitTag, constants.GitTagListFlag, constants.GitTagGlob)
	out, err := cmd.Output()
	if err != nil {
		return Version{}, fmt.Errorf("no git tags found and no .gitmap/release/latest.json exists")
	}

	return findHighestVersion(strings.TrimSpace(string(out)))
}

// findHighestVersion parses tag lines and returns the highest stable version.
func findHighestVersion(output string) (Version, error) {
	lines := strings.Split(output, "\n")
	var highest Version
	matched := false

	for _, line := range lines {
		v, ok := parseStableTag(line)
		if ok {
			highest, matched = updateHighest(highest, v, matched)
		}
	}

	if matched {
		fmt.Printf("  → Detected latest version from git tags: %s\n", highest.String())

		return highest, nil
	}

	return Version{}, fmt.Errorf("no version tags found. Create an initial release first")
}

// updateHighest returns the higher of current and candidate versions.
func updateHighest(current, candidate Version, hasCurrent bool) (Version, bool) { //nolint:unparam // bool clarifies caller intent
	if hasCurrent {
		if candidate.GreaterThan(current) {
			return candidate, true
		}

		return current, true
	}

	return candidate, true
}

// parseStableTag attempts to parse a tag as a stable (non-prerelease) version.
func parseStableTag(line string) (Version, bool) {
	tag := strings.TrimSpace(line)
	if len(tag) == 0 {
		return Version{}, false
	}

	v, err := Parse(tag)
	if err != nil {
		return Version{}, false
	}

	if v.IsPreRelease() {
		return Version{}, false
	}

	return v, true
}
