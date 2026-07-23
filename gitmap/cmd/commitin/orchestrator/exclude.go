package orchestrator

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// applyExclusions filters `files` against the resolved Exclusion list.
// Returns the kept subset preserving input order. Rules:
//   - PathFile: exact relative-path match (case-sensitive, POSIX).
//   - PathFolder: file is excluded when any path segment equals the
//     folder name OR the file path begins with "<folder>/".
func applyExclusions(files []string, rules []profile.Exclusion) []string {
	if len(files) == 0 || len(rules) == 0 {
		return files
	}
	out := make([]string, 0, len(files))
	for _, f := range files {
		if !isExcluded(f, rules) {
			out = append(out, f)
		}
	}
	return out
}

func isExcluded(rel string, rules []profile.Exclusion) bool {
	clean := filepath.ToSlash(rel)
	for _, r := range rules {
		if matchesExclusion(clean, r) {
			return true
		}
	}
	return false
}

func matchesExclusion(rel string, r profile.Exclusion) bool {
	switch r.Kind {
	case constants.CommitInExclusionKindPathFile:
		return rel == r.Value
	case constants.CommitInExclusionKindPathFolder:
		return matchesFolder(rel, r.Value)
	}
	return false
}

func matchesFolder(rel, folder string) bool {
	folder = strings.Trim(folder, "/")
	if folder == "" {
		return false
	}
	if strings.HasPrefix(rel, folder+"/") {
		return true
	}
	for _, seg := range strings.Split(rel, "/") {
		if seg == folder {
			return true
		}
	}
	return false
}
