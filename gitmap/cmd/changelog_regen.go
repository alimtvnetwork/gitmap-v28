// changelog_regen.go — derives CHANGELOG.md entries from the
// canonical per-release JSON files under `.gitmap/release/` (#18).
//
// Until v6.60.0 the changelog was hand-edited in lockstep with the
// version bump, drifting easily. The release JSONs already capture
// `{version, tag, branch}` per release; this helper enumerates them,
// sorts by semver descending, and prints a freshly-rendered block
// that can be diffed against CHANGELOG.md (or piped in).
//
// Invoked via `gitmap changelog regen` (read-only stdout). The
// existing scripts/changelog/ Go module remains the authoritative
// generator for CI; this helper is for ad-hoc developer use without
// leaving the binary.
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ReleaseMeta mirrors the on-disk shape of `.gitmap/release/*.json`.
type ReleaseMeta struct {
	Version string `json:"version"`
	Tag     string `json:"tag"`
	Branch  string `json:"branch"`
}

// RegenChangelog writes a Markdown changelog skeleton (newest first)
// derived from every release JSON under `releaseDir`. Each entry is
// a stub the developer fills in with notes; the version + tag are
// the source of truth so drift between files is impossible.
func RegenChangelog(releaseDir string, w io.Writer) error {
	entries, err := os.ReadDir(releaseDir)
	if err != nil {
		return err
	}
	metas := loadReleaseMetas(releaseDir, entries)
	sortMetasDesc(metas)

	fmt.Fprintln(w, "# Changelog")
	fmt.Fprintln(w, "")
	for _, m := range metas {
		fmt.Fprintf(w, "## %s\n\n", m.Tag)
		fmt.Fprintf(w, "- (fill in)\n\n")
	}
	return nil
}

func loadReleaseMetas(dir string, entries []os.DirEntry) []ReleaseMeta {
	out := make([]ReleaseMeta, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".json") || name == "latest.json" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		var m ReleaseMeta
		if err := json.Unmarshal(raw, &m); err != nil || m.Version == "" {
			continue
		}
		out = append(out, m)
	}
	return out
}

func sortMetasDesc(m []ReleaseMeta) {
	sort.Slice(m, func(i, j int) bool {
		return compareSemverDesc(m[i].Version, m[j].Version)
	})
}

// compareSemverDesc returns true when a > b. Tolerant of missing
// patch numbers; non-numeric components compare as 0.
func compareSemverDesc(a, b string) bool {
	pa, pb := splitSemver(a), splitSemver(b)
	for i := 0; i < 3; i++ {
		if pa[i] != pb[i] {
			return pa[i] > pb[i]
		}
	}
	return false
}

func splitSemver(v string) [3]int {
	parts := strings.SplitN(strings.TrimPrefix(v, "v"), ".", 3)
	var out [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n, _ := strconv.Atoi(parts[i])
		out[i] = n
	}
	return out
}
