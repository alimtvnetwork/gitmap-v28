package release

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TagEntry holds a parsed git tag name and its creation date.
type TagEntry struct {
	Tag       string
	CreatedAt string
}

// ListVersionTags returns all semver tags with their creation dates.
func ListVersionTags() []TagEntry {
	cmd := exec.Command(constants.GitBin,
		constants.GitForEachRef,
		constants.GitForEachRefTagFmt,
		constants.GitRefsTagsPrefix,
	)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	return parseTagEntries(strings.TrimSpace(string(out)))
}

// parseTagEntries parses for-each-ref output lines into TagEntry slices.
func parseTagEntries(output string) []TagEntry {
	if len(output) == 0 {
		return nil
	}

	lines := strings.Split(output, "\n")
	entries := make([]TagEntry, 0, len(lines))

	for _, line := range lines {
		entry, ok := parseTagLine(line)
		if ok {
			entries = append(entries, entry)
		}
	}

	return entries
}

// parseTagLine parses a single "tag|date" line into a TagEntry.
func parseTagLine(line string) (TagEntry, bool) {
	parts := strings.SplitN(strings.TrimSpace(line), "|", 2)
	if len(parts) != 2 {
		return TagEntry{}, false
	}

	tag := parts[0]
	_, err := Parse(tag)
	if err != nil {
		return TagEntry{}, false
	}

	return TagEntry{Tag: tag, CreatedAt: parts[1]}, true
}
