// Package cmdcg provides version drift detection and comparison against canonical standards.
package cmdcg

import (
	"fmt"
	"path/filepath"
	"strings"
)

// CanonicalCGVersion defines the reference coding guidelines release.
const CanonicalCGVersion = "v2.24.0"

// CGDiffEntry records version comparison status for a single repository.
type CGDiffEntry struct {
	RepoPath       string
	CurrentVersion string
	Status         string
	HasDrift       bool
	IsInstalled    bool
}

func isVersionEqual(v1, v2 string) bool {
	clean1 := strings.TrimPrefix(v1, "v")
	clean2 := strings.TrimPrefix(v2, "v")

	return clean1 == clean2
}

func buildDriftEntry(repoPath, currentVersion string) CGDiffEntry {
	return CGDiffEntry{
		RepoPath:       repoPath,
		CurrentVersion: currentVersion,
		Status:         fmt.Sprintf("Drift detected: %s -> %s", currentVersion, CanonicalCGVersion),
		HasDrift:       true,
		IsInstalled:    true,
	}
}

func evaluateCGDiff(repoPath, currentVersion string) CGDiffEntry {
	if isVersionEqual(currentVersion, CanonicalCGVersion) {
		return CGDiffEntry{
			RepoPath:       repoPath,
			CurrentVersion: currentVersion,
			Status:         "Up to date",
			HasDrift:       false,
			IsInstalled:    true,
		}
	}

	return buildDriftEntry(repoPath, currentVersion)
}

func diffSingleRepo(repoPath string) CGDiffEntry {
	meta, err := ReadCGMetadata(repoPath)
	if err != nil || meta == nil || meta.Version == "" {
		return CGDiffEntry{
			RepoPath:       repoPath,
			CurrentVersion: "none",
			Status:         "Not installed",
			HasDrift:       true,
			IsInstalled:    false,
		}
	}

	return evaluateCGDiff(repoPath, meta.Version)
}

// DiffSingleRepo computes the version drift entry for a given repository.
func DiffSingleRepo(repoPath string) CGDiffEntry {
	return diffSingleRepo(repoPath)
}

func renderDiffStatus(entry CGDiffEntry) string {
	if !entry.IsInstalled {
		return cgFailStyle.Render("✖ " + entry.Status)
	}

	if entry.HasDrift {
		return cgWarnStyle.Render("▲ " + entry.Status)
	}

	return cgPassStyle.Render("✔ " + entry.Status)
}

func renderDiffRow(entry CGDiffEntry) {
	name := filepath.Base(entry.RepoPath)
	statusStr := renderDiffStatus(entry)
	fmt.Printf("  %-25s %-12s %s\n", name, entry.CurrentVersion, statusStr)
}

func classifyDiffEntry(e CGDiffEntry, notInstalled, drifted, upToDate *int) {
	if !e.IsInstalled {
		*notInstalled++

		return
	}

	if e.HasDrift {
		*drifted++

		return
	}

	*upToDate++
}

func renderDiffSummary(entries []CGDiffEntry) {
	upToDateCount, driftCount, notInstalledCount := 0, 0, 0
	for _, e := range entries {
		classifyDiffEntry(e, &notInstalledCount, &driftCount, &upToDateCount)
	}

	fmt.Printf("\n  Total: %d | Up to date: %d | Drifted: %d | Not installed: %d\n", len(entries), upToDateCount, driftCount, notInstalledCount)
}

func printDiffTableHeader() {
	fmt.Printf("%s\n", cgAuditTitleStyle.Render("=== Coding Guidelines Version Drift & Diff ==="))
	fmt.Printf("Canonical Version: %s\n\n", cgPassStyle.Render(CanonicalCGVersion))
	fmt.Printf("  %-25s %-12s %s\n", "REPOSITORY", "INSTALLED", "STATUS")
	fmt.Printf("  %-25s %-12s %s\n", "----------", "---------", "------")
}

// RunCGDiff compares local installed coding guidelines versions against the canonical release.
func RunCGDiff(repos []string) error {
	printDiffTableHeader()
	var entries []CGDiffEntry

	for _, repo := range repos {
		entry := diffSingleRepo(repo)
		entries = append(entries, entry)
		renderDiffRow(entry)
	}

	renderDiffSummary(entries)

	return nil
}
