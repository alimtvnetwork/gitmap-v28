// Package gitutil — dirty_inspect.go inspects uncommitted files and unstaged changes.
package gitutil

import (
	"os/exec"
	"strconv"
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
	DeletedFiles   []string
	StagedFiles    []string
	AllFiles       []string
}

func recordUntracked(diagnosis *DirtyDiagnosis, filePath string) {
	diagnosis.UntrackedCount++
	diagnosis.UntrackedFiles = append(diagnosis.UntrackedFiles, filePath)
	diagnosis.AllFiles = append(diagnosis.AllFiles, "untracked: "+filePath)
}

func recordModified(diagnosis *DirtyDiagnosis, filePath string) {
	diagnosis.ModifiedCount++
	diagnosis.ModifiedFiles = append(diagnosis.ModifiedFiles, filePath)
	diagnosis.AllFiles = append(diagnosis.AllFiles, "modified: "+filePath)
}

func recordDeleted(diagnosis *DirtyDiagnosis, filePath string) {
	diagnosis.DeletedCount++
	diagnosis.DeletedFiles = append(diagnosis.DeletedFiles, filePath)
	diagnosis.AllFiles = append(diagnosis.AllFiles, "deleted: "+filePath)
}

func recordStaged(diagnosis *DirtyDiagnosis, filePath string) {
	diagnosis.StagedCount++
	diagnosis.StagedFiles = append(diagnosis.StagedFiles, filePath)
	diagnosis.AllFiles = append(diagnosis.AllFiles, "staged: "+filePath)
}

func classifyDirtyFile(prefix, filePath string, diagnosis *DirtyDiagnosis) {
	if strings.HasPrefix(prefix, "??") {
		recordUntracked(diagnosis, filePath)

		return
	}
	if strings.Contains(prefix, "M") {
		recordModified(diagnosis, filePath)

		return
	}
	if strings.Contains(prefix, "D") {
		recordDeleted(diagnosis, filePath)

		return
	}
	recordStaged(diagnosis, filePath)
}

func parseDirtyLine(line string, diagnosis *DirtyDiagnosis) {
	if len(line) < 3 {
		return
	}
	prefix := line[:2]
	filePath := strings.TrimSpace(line[3:])
	classifyDirtyFile(prefix, filePath, diagnosis)
}

func collectReasonParts(diagnosis *DirtyDiagnosis) []string {
	var parts []string
	if diagnosis.ModifiedCount > 0 {
		parts = append(parts, "+"+strconv.Itoa(diagnosis.ModifiedCount)+" modified")
	}
	if diagnosis.UntrackedCount > 0 {
		parts = append(parts, "+"+strconv.Itoa(diagnosis.UntrackedCount)+" untracked")
	}
	if diagnosis.DeletedCount > 0 {
		parts = append(parts, "-"+strconv.Itoa(diagnosis.DeletedCount)+" deleted")
	}

	return parts
}

func buildSummaryReason(diagnosis *DirtyDiagnosis) string {
	parts := collectReasonParts(diagnosis)
	if len(parts) <= 0 {
		return "uncommitted changes"
	}

	return strings.Join(parts, ", ")
}

func populateDirtyDiagnosis(output string) DirtyDiagnosis {
	diagnosis := DirtyDiagnosis{IsDirty: true}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		parseDirtyLine(line, &diagnosis)
	}
	diagnosis.SummaryReason = buildSummaryReason(&diagnosis)

	return diagnosis
}

// InspectDirtyState analyzes the working directory to diagnose why a repo is dirty.
func InspectDirtyState(repoPath string) DirtyDiagnosis {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	outputBytes, err := cmd.Output()
	if err != nil {
		return DirtyDiagnosis{IsDirty: false}
	}
	trimmedOutput := strings.TrimSpace(string(outputBytes))
	if len(trimmedOutput) <= 0 {
		return DirtyDiagnosis{IsDirty: false}
	}

	return populateDirtyDiagnosis(trimmedOutput)
}
