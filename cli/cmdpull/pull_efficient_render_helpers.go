package cmdpull

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// ResolveConciseRepoColWidth dynamically calculates the repo column width to prevent overflow.
func ResolveConciseRepoColWidth(states []*PullRepoState, allRecords ...[]model.ScanRecord) int {
	colWidth := minConciseRepoColWidth
	var records []model.ScanRecord
	if len(allRecords) > 0 {
		records = allRecords[0]
	}
	for _, r := range records {
		if len(r.RepoName) > colWidth {
			colWidth = len(r.RepoName)
		}
	}
	collisions := detectRepoCollisions(states, allRecords...)
	for _, s := range states {
		name := resolveRepoDisplayName(s, collisions)
		if len(name) > colWidth {
			colWidth = len(name)
		}
	}

	termWidth := detectTerminalWidth()
	maxAllowed := termWidth - conciseReservedColGap
	if maxAllowed > minConciseRepoColWidth && colWidth > maxAllowed {
		colWidth = maxAllowed
	}

	return colWidth
}

// ResolveRepoStatusLabel normalizes empty or synced changes to up-to-date.
func ResolveRepoStatusLabel(changes string) string {
	if changes == "" || changes == "synced" {
		return "up-to-date"
	}
	return changes
}

// StyleRepoStatusLabel decorates status labels with ANSI colors for fast visual triage.
func StyleRepoStatusLabel(statusLabel string) string {
	switch statusLabel {
	case "up-to-date":
		return constants.ColorGreen + "up-to-date" + constants.ColorReset
	case "dirty":
		return constants.ColorYellow + "dirty" + constants.ColorReset
	case "failed":
		return constants.ColorRed + "failed" + constants.ColorReset
	}

	if strings.HasPrefix(statusLabel, "+") {
		return formatDiffStatsColor(statusLabel)
	}

	return statusLabel
}

func formatDiffStatsColor(statusLabel string) string {
	parts := strings.SplitN(statusLabel, " ", 2)
	diffPart := parts[0]
	slashIdx := strings.Index(diffPart, "/")
	if slashIdx == -1 {
		return constants.ColorGreen + statusLabel + constants.ColorReset
	}

	ins := diffPart[:slashIdx]
	del := diffPart[slashIdx:]
	coloredDiff := constants.ColorGreen + ins + constants.ColorReset + constants.ColorRed + del + constants.ColorReset
	if len(parts) > 1 {
		return coloredDiff + " " + constants.ColorDim + parts[1] + constants.ColorReset
	}

	return coloredDiff
}

// FormatConciseActiveResultLine formats a single row with dynamic padding and clean spacing.
func FormatConciseActiveResultLine(colWidth int, repoName, statusLabel string) string {
	displayName := repoName
	if len(displayName) > colWidth && colWidth > 5 {
		displayName = displayName[:colWidth-3] + "..."
	}
	styledStatus := StyleRepoStatusLabel(statusLabel)
	return fmt.Sprintf("    • %-*s  %s", colWidth, displayName, styledStatus)
}

// FormatWrappedInactiveList formats repository names wrapped cleanly at commas without splitting tokens.
func FormatWrappedInactiveList(names []string, indent string, maxLineLen int) string {
	if len(names) == 0 {
		return ""
	}
	if maxLineLen <= len(indent)+10 {
		maxLineLen = 120
	}

	var sb strings.Builder
	currentLineLen := len(indent)
	sb.WriteString(indent)

	for i, name := range names {
		token := name
		if i < len(names)-1 {
			token += ", "
		}

		if currentLineLen+len(token) > maxLineLen && currentLineLen > len(indent) {
			sb.WriteString("\n")
			sb.WriteString(indent)
			currentLineLen = len(indent)
		}

		sb.WriteString(token)
		currentLineLen += len(token)
	}

	return sb.String()
}

func detectRepoCollisions(states []*PullRepoState, allRecords ...[]model.ScanRecord) map[string]bool {
	pathsByName := make(map[string]map[string]bool)
	collectStatePaths(pathsByName, states)
	if len(allRecords) > 0 {
		collectRecordPaths(pathsByName, allRecords[0])
	}
	return filterCollidingNames(pathsByName)
}

func collectStatePaths(pathsByName map[string]map[string]bool, states []*PullRepoState) {
	for _, s := range states {
		if s == nil || s.RepoName == "" {
			continue
		}
		pathKey := CanonicalRepoPathKey(s.RepoPath)
		if pathKey == "" {
			pathKey = strings.ToLower(s.RepoName)
		}
		if pathsByName[s.RepoName] == nil {
			pathsByName[s.RepoName] = make(map[string]bool)
		}
		pathsByName[s.RepoName][pathKey] = true
	}
}

func collectRecordPaths(pathsByName map[string]map[string]bool, records []model.ScanRecord) {
	for _, r := range records {
		if r.RepoName == "" {
			continue
		}
		pathKey := CanonicalRepoPathKey(r.AbsolutePath)
		if pathKey == "" {
			pathKey = strings.ToLower(r.RepoName)
		}
		if pathsByName[r.RepoName] == nil {
			pathsByName[r.RepoName] = make(map[string]bool)
		}
		pathsByName[r.RepoName][pathKey] = true
	}
}

func filterCollidingNames(pathsByName map[string]map[string]bool) map[string]bool {
	hasCollision := make(map[string]bool, len(pathsByName))
	for name, paths := range pathsByName {
		if len(paths) > 1 {
			hasCollision[name] = true
		}
	}
	return hasCollision
}

func resolveRepoDisplayName(s *PullRepoState, collisions map[string]bool) string {
	if s == nil {
		return ""
	}
	hasCollision := collisions[s.RepoName]
	if hasCollision && s.RepoPath != "" {
		rel := formatRelativeOrCleanPath(s.RepoPath)
		return fmt.Sprintf("%s (%s)", s.RepoName, rel)
	}
	return s.RepoName
}

func formatRelativeOrCleanPath(repoPath string) string {
	if repoPath == "" {
		return ""
	}

	cwd, err := os.Getwd()
	if err != nil || cwd == "" {
		return filepath.Clean(repoPath)
	}

	rel, relErr := filepath.Rel(cwd, repoPath)
	isDescendant := relErr == nil && !strings.HasPrefix(rel, "..")
	if isDescendant {
		return filepath.ToSlash(rel)
	}

	return filepath.Clean(repoPath)
}
