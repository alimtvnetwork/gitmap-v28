package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// collectVersionTags reads git tags, parses, sorts descending.
func collectVersionTags() []release.Version {
	cmd := exec.Command(constants.GitBin, constants.GitTag,
		constants.GitTagListFlag, constants.GitTagGlob)
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, constants.ErrListVersionsNoTags)
		os.Exit(1)
	}

	versions := parseVersionTags(strings.TrimSpace(string(out)))
	if len(versions) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrListVersionsNoTags)
		os.Exit(1)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].GreaterThan(versions[j])
	})

	return versions
}

// parseVersionTags parses lines into valid versions.
func parseVersionTags(output string) []release.Version {
	lines := strings.Split(output, "\n")
	var versions []release.Version

	for _, line := range lines {
		tag := strings.TrimSpace(line)
		if len(tag) == 0 {
			continue
		}
		v, err := release.Parse(tag)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}

	return versions
}

// loadChangelogMap reads CHANGELOG.md into a version→notes map.
func loadChangelogMap() map[string][]string {
	entries, err := release.ReadChangelog()
	if err != nil {
		return map[string][]string{}
	}

	m := make(map[string][]string, len(entries))
	for _, e := range entries {
		m[e.Version] = e.Notes
	}

	return m
}

// printVersionEntriesTerminal prints versions with source and changelog sub-points.
func printVersionEntriesTerminal(entries []versionEntry) {
	for _, e := range entries {
		if e.Source != "" {
			fmt.Printf("%s  [%s]\n", e.Version.String(), e.Source)
		} else {
			fmt.Println(e.Version.String())
		}
		for _, note := range e.Notes {
			fmt.Printf("  - %s\n", note)
		}
	}
}

// printVersionEntriesJSON prints versions as stable JSON via the
// stablejson encoder in listversionsrender.go. Key order is a
// compile-time decision (version, source, changelog) rather than a
// reflection accident; source/changelog remain effectively omitempty.
func printVersionEntriesJSON(entries []versionEntry) {
	if err := encodeListVersionsJSON(os.Stdout, entries); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to marshal versions to JSON: %v\n", err)
	}
}
