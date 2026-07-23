package vscodepm

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// DetectTags inspects rootPath's top-level filesystem and returns a sorted,
// de-duplicated list of marker tags (e.g. "git", "node", "go"). The set is
// intentionally small and deterministic so projects.json diffs stay stable
// across runs.
//
// Detection is shallow (top-level only) and read-only — no recursion, no
// network, no shelling out. Unreadable / missing rootPath returns nil so
// callers can treat detection as best-effort.
//
// The returned slice respects the canonical order in constants.AutoTagOrder
// so two runs on the same repo always produce identical JSON output.
func DetectTags(rootPath string) []string {
	if rootPath == "" {
		return nil
	}

	info, err := os.Stat(rootPath)
	if err != nil || !info.IsDir() {
		return nil
	}

	hits := map[string]struct{}{}
	for marker, tag := range constants.AutoTagMarkers {
		if markerExists(rootPath, marker) {
			hits[tag] = struct{}{}
		}
	}

	return tagsInCanonicalOrder(hits)
}

// markerExists reports whether <root>/<marker> exists. Both files and
// directories count (e.g. ".git" can be either; "node_modules" is a dir;
// "package.json" is a file).
func markerExists(root, marker string) bool {
	_, err := os.Stat(filepath.Join(root, marker))

	return err == nil
}

// tagsInCanonicalOrder returns hits ordered by constants.AutoTagOrder so
// the JSON output stays diff-stable across runs and machines.
func tagsInCanonicalOrder(hits map[string]struct{}) []string {
	if len(hits) == 0 {
		return nil
	}

	out := make([]string, 0, len(hits))
	for _, tag := range constants.AutoTagOrder {
		if _, ok := hits[tag]; ok {
			out = append(out, tag)
		}
	}

	return out
}

// unionTags merges existing on-disk tags with auto-detected tags, preserving
// the order of `existing` first then appending any new auto tags. This means
// user-edited tags in the VS Code UI are NEVER removed by gitmap — auto tags
// are purely additive.
func unionTags(existing, incoming []string) []string {
	seen := make(map[string]struct{}, len(existing)+len(incoming))
	out := make([]string, 0, len(existing)+len(incoming))

	for _, t := range existing {
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}

	for _, t := range incoming {
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}

	return out
}
