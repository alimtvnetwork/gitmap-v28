package gitutil

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// IndexEntry represents a parsed record from git ls-files -s.
type IndexEntry struct {
	Mode      string
	BlobSHA   string
	Stage     string
	Path      string
	LowerPath string
}

// CaseCollision represents a detected case conflict between two or more index paths.
type CaseCollision struct {
	CanonicalPath   string
	CollidingPaths  []string
	ActiveEntries   []IndexEntry
	HasDivergentSHA bool
}

// DetectCaseCollisions scans the repository index for case collisions.
func DetectCaseCollisions(repoPath string) ([]CaseCollision, error) {
	out, err := runLsFilesStage(repoPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "detect case collisions")
	}

	entries := parseGitLsFilesOutput(out)
	buckets := groupEntriesByLowerPath(entries)

	return filterCollidingBuckets(buckets), nil
}

func runLsFilesStage(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "ls-files", "-s")
	cmd.Env = BuildSafeGitEnv()
	outBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(outBytes), nil
}

func parseGitLsFilesOutput(output string) []IndexEntry {
	lines := strings.Split(output, "\n")
	entries := make([]IndexEntry, 0, len(lines))
	for _, line := range lines {
		entry, hasEntry := parseSingleIndexEntry(line)
		if hasEntry {
			entries = append(entries, entry)
		}
	}

	return entries
}

func parseSingleIndexEntry(line string) (IndexEntry, bool) {
	tabIdx := strings.Index(line, "\t")
	if tabIdx == -1 {
		return IndexEntry{}, false
	}
	metaPart := strings.TrimSpace(line[:tabIdx])
	filePath := normalizeIndexPath(line[tabIdx+1:])
	fields := strings.Fields(metaPart)
	if len(fields) < 3 {
		return IndexEntry{}, false
	}

	return buildParsedIndexEntry(fields, filePath), true
}

func normalizeIndexPath(raw string) string {
	clean := strings.TrimSpace(raw)
	clean = strings.Trim(clean, "\"")

	return filepath.ToSlash(clean)
}

func buildParsedIndexEntry(fields []string, filePath string) IndexEntry {
	return IndexEntry{
		Mode:      fields[0],
		BlobSHA:   fields[1],
		Stage:     fields[2],
		Path:      filePath,
		LowerPath: strings.ToLower(filePath),
	}
}

func groupEntriesByLowerPath(entries []IndexEntry) map[string][]IndexEntry {
	buckets := make(map[string][]IndexEntry)
	for _, entry := range entries {
		buckets[entry.LowerPath] = append(buckets[entry.LowerPath], entry)
	}

	return buckets
}

func filterCollidingBuckets(buckets map[string][]IndexEntry) []CaseCollision {
	var collisions []CaseCollision
	for _, bucket := range buckets {
		collision, hasCollision := evaluateBucketCollision(bucket)
		if hasCollision {
			collisions = append(collisions, collision)
		}
	}

	return collisions
}

func evaluateBucketCollision(entries []IndexEntry) (CaseCollision, bool) {
	if len(entries) <= 1 {
		return CaseCollision{}, false
	}
	if !hasCasingDifference(entries) {
		return CaseCollision{}, false
	}

	return buildCollisionRecord(entries), true
}

func hasCasingDifference(entries []IndexEntry) bool {
	firstPath := entries[0].Path
	for _, e := range entries[1:] {
		if e.Path != firstPath {
			return true
		}
	}

	return false
}

func buildCollisionRecord(entries []IndexEntry) CaseCollision {
	canonical := resolveCanonicalPath(entries)
	paths := extractUniquePaths(entries)
	hasDivergent := checkDivergentSHA(entries)

	return CaseCollision{
		CanonicalPath:   canonical,
		CollidingPaths:  paths,
		ActiveEntries:   entries,
		HasDivergentSHA: hasDivergent,
	}
}

func resolveCanonicalPath(entries []IndexEntry) string {
	for _, e := range entries {
		if e.Path == e.LowerPath {
			return e.Path
		}
	}

	return entries[0].LowerPath
}

func extractUniquePaths(entries []IndexEntry) []string {
	seen := make(map[string]bool)
	var paths []string
	for _, e := range entries {
		if !seen[e.Path] {
			seen[e.Path] = true
			paths = append(paths, e.Path)
		}
	}

	return paths
}

func checkDivergentSHA(entries []IndexEntry) bool {
	firstSHA := entries[0].BlobSHA
	for _, e := range entries[1:] {
		if e.BlobSHA != firstSHA {
			return true
		}
	}

	return false
}
