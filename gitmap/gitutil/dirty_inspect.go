// Package gitutil — dirty_inspect.go inspects uncommitted files and unstaged changes.
package gitutil

import (
	"os/exec"
	"strings"
)

type DirtyDiagnosis struct {
	IsDirty        bool
	ModifiedCount  int
	UntrackedCount int
	DeletedCount   int
	StagedCount    int
	SummaryReason  string
	ModifiedFiles  []string
	UntrackedFiles []string
}

// InspectDirtyState analyzes the working directory to diagnose why a repo is dirty.
func InspectDirtyState(repoPath string) DirtyDiagnosis {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return DirtyDiagnosis{IsDirty: false}
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return DirtyDiagnosis{IsDirty: false}
	}

	var d DirtyDiagnosis
	d.IsDirty = true

	for _, l := range lines {
		if len(l) < 3 {
			continue
		}
		prefix := l[:2]
		file := strings.TrimSpace(l[3:])

		switch {
		case strings.HasPrefix(prefix, "??"):
			d.UntrackedCount++
			d.UntrackedFiles = append(d.UntrackedFiles, file)
		case strings.Contains(prefix, "M"):
			d.ModifiedCount++
			d.ModifiedFiles = append(d.ModifiedFiles, file)
		case strings.Contains(prefix, "D"):
			d.DeletedCount++
		case prefix[0] != ' ' && prefix[0] != '?':
			d.StagedCount++
		}
	}

	var parts []string
	if d.ModifiedCount > 0 {
		parts = append(parts, "+"+string(rune('0'+d.ModifiedCount))+" modified")
	}
	if d.UntrackedCount > 0 {
		parts = append(parts, "+"+string(rune('0'+d.UntrackedCount))+" untracked")
	}
	if d.DeletedCount > 0 {
		parts = append(parts, "-"+string(rune('0'+d.DeletedCount))+" deleted")
	}
	if len(parts) == 0 {
		d.SummaryReason = "uncommitted changes"
	} else {
		d.SummaryReason = strings.Join(parts, ", ")
	}

	return d
}
