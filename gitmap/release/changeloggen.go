package release

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// GenerateChangelog returns commit subjects between two tags (or from a tag to HEAD).
func GenerateChangelog(fromTag, toRef string) ([]string, error) {
	if len(toRef) == 0 {
		toRef = constants.GitHEAD
	}

	rangeArg := fmt.Sprintf("%s..%s", fromTag, toRef)
	cmd := exec.Command(constants.GitBin, constants.GitLog,
		constants.ChangelogGenFormat, constants.ChangelogGenNoMerges, rangeArg)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrChangelogGenCommits, fromTag, toRef, err)
	}

	return parseCommitLines(strings.TrimSpace(string(out))), nil
}

// parseCommitLines splits raw output into non-empty trimmed lines.
func parseCommitLines(output string) []string {
	if len(output) == 0 {
		return nil
	}

	lines := strings.Split(output, "\n")
	var results []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 {
			results = append(results, trimmed)
		}
	}

	return results
}

// FormatChangelogSection formats commits as a markdown changelog section.
func FormatChangelogSection(version string, commits []string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("## %s\n\n", NormalizeVersion(version)))

	for _, c := range commits {
		b.WriteString(fmt.Sprintf("- %s\n", c))
	}

	return b.String()
}

// ListTags returns all version tags sorted descending by version.
func ListTags() ([]string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitTag,
		constants.GitTagListFlag, constants.GitTagGlob,
		constants.ChangelogGenSortFlag, constants.ChangelogGenSortVersion)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrChangelogGenTags, err)
	}

	return parseCommitLines(strings.TrimSpace(string(out))), nil
}

// ResolveTagRange determines the from/to refs based on user input.
func ResolveTagRange(fromTag, toTag string) (string, string, error) {
	if len(fromTag) > 0 {
		from := NormalizeVersion(fromTag)
		if !TagExistsLocally(from) {
			return "", "", fmt.Errorf(constants.ErrChangelogGenTagNotFound, from)
		}

		to := constants.GitHEAD
		if len(toTag) > 0 {
			to = NormalizeVersion(toTag)
			if !TagExistsLocally(to) {
				return "", "", fmt.Errorf(constants.ErrChangelogGenTagNotFound, to)
			}
		}

		return from, to, nil
	}

	tags, err := ListTags()
	if err != nil || len(tags) < 1 {
		return "", "", fmt.Errorf(constants.ErrChangelogGenNoTags)
	}

	if len(tags) == 1 {
		return tags[0], constants.GitHEAD, nil
	}

	return tags[1], tags[0], nil
}
