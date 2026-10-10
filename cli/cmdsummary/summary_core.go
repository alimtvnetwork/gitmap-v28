// Package cmdsummary implements single repository release summaries and heated file churn analytics.
package cmdsummary

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// gitExec runs a git command in the target directory and returns the trimmed stdout.
func gitExec(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// RunSummary executes single repository release summarization with Split-DB caching.
func RunSummary(args []string) error {
	targetDir, releasesLimit := parseSummaryArgs(args)

	absPath, errAbs := filepath.Abs(targetDir)
	if errAbs != nil {
		absPath = targetDir
	}

	gitPath := filepath.Join(absPath, ".git")
	if _, errStat := os.Stat(gitPath); errStat != nil {
		return fmt.Errorf("directory '%s' is not a valid git repository", absPath)
	}

	remoteURL, _ := gitExec(absPath, "remote", "get-url", "origin")
	branch, _ := gitExec(absPath, "rev-parse", "--abbrev-ref", "HEAD")
	headHash, _ := gitExec(absPath, "rev-parse", "HEAD")
	repoName := filepath.Base(absPath)
	if remoteURL == "" {
		remoteURL = repoName
	}

	startTime := time.Now()

	releases, cacheHits := extractRepoReleaseSummaries(absPath, remoteURL, releasesLimit)

	renderSingleRepoSummaryTerminal(repoName, remoteURL, branch, headHash, releases, cacheHits, time.Since(startTime))
	return nil
}

func parseSummaryArgs(args []string) (string, int) {
	targetDir := "."
	releasesLimit := 8

	for _, arg := range args {
		if arg == "" || strings.HasPrefix(arg, "-") {
			continue
		}
		if n, err := strconv.Atoi(arg); err == nil && n > 0 {
			releasesLimit = n
			continue
		}
		targetDir = arg
	}

	return targetDir, releasesLimit
}

func extractRepoReleaseSummaries(dir, remoteURL string, limit int) ([]ReleaseSummaryRecord, int) {
	tagOutput, _ := gitExec(dir, "tag", "--sort=-creatordate")
	rawTags := strings.Split(tagOutput, "\n")

	var validTags []string
	for _, t := range rawTags {
		clean := strings.TrimSpace(t)
		if clean != "" {
			validTags = append(validTags, clean)
		}
	}

	if len(validTags) == 0 {
		return extractFallbackCommitSummaries(dir, limit), 0
	}

	if len(validTags) > limit {
		validTags = validTags[:limit]
	}

	var results []ReleaseSummaryRecord
	cacheHits := 0

	for i, tag := range validTags {
		tagCommitHash, _ := gitExec(dir, "rev-list", "-n", "1", tag)
		if tagCommitHash == "" {
			tagCommitHash = tag
		}

		if cached, isHit := GetCachedRelease(remoteURL, tag, tagCommitHash); isHit {
			results = append(results, *cached)
			cacheHits++
			continue
		}

		var prevTag string
		if i+1 < len(validTags) {
			prevTag = validTags[i+1]
		}

		rec := computeReleaseSummary(dir, tag, prevTag, tagCommitHash)
		_ = SaveCachedRelease(remoteURL, filepath.Base(dir), dir, rec)
		results = append(results, rec)
	}

	return results, cacheHits
}

func extractFallbackCommitSummaries(dir string, limit int) []ReleaseSummaryRecord {
	logOutput, errLog := gitExec(dir, "log", fmt.Sprintf("-%d", limit), "--pretty=format:%h|%cs|%s")
	if errLog != nil || logOutput == "" {
		return nil
	}

	var recs []ReleaseSummaryRecord
	lines := strings.Split(logOutput, "\n")
	for _, l := range lines {
		parts := strings.SplitN(l, "|", 3)
		if len(parts) < 3 {
			continue
		}
		recs = append(recs, ReleaseSummaryRecord{
			TagName:       parts[0],
			TagCommitHash: parts[0],
			ReleaseDate:   parts[1],
			SummaryGist:   parts[2],
			WordCount:     len(strings.Fields(parts[2])),
		})
	}
	return recs
}

func computeReleaseSummary(dir, tag, prevTag, tagCommitHash string) ReleaseSummaryRecord {
	releaseDate, _ := gitExec(dir, "log", "-1", "--format=%cs", tag)
	if releaseDate == "" {
		releaseDate = time.Now().Format("2006-01-02")
	}

	rangeExpr := tag
	if prevTag != "" {
		rangeExpr = fmt.Sprintf("%s..%s", prevTag, tag)
	}

	commitsOut, _ := gitExec(dir, "log", rangeExpr, "--oneline", "-n", "30")
	diffStatOut, _ := gitExec(dir, "diff", "--numstat", rangeExpr)

	gist := synthesizeReleaseGist(commitsOut)
	heatedFiles := extractHeatedFilesFromDiff(diffStatOut)

	return ReleaseSummaryRecord{
		TagName:       tag,
		TagCommitHash: tagCommitHash,
		ReleaseDate:   releaseDate,
		SummaryGist:   gist,
		WordCount:     len(strings.Fields(gist)),
		HeatedFiles:   heatedFiles,
		IsCacheHit:    false,
	}
}

func synthesizeReleaseGist(commitLog string) string {
	if commitLog == "" {
		return "Maintenance updates and codebase stability improvements."
	}

	lines := strings.Split(commitLog, "\n")
	var subjects []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		parts := strings.SplitN(trimmed, " ", 2)
		if len(parts) == 2 {
			subjects = append(subjects, parts[1])
		}
	}

	if len(subjects) == 0 {
		return "Maintenance and build updates."
	}

	gist := strings.Join(subjects, "; ")
	words := strings.Fields(gist)
	if len(words) > 180 {
		words = words[:180]
		gist = strings.Join(words, " ") + "..."
	}

	return gist
}

func extractHeatedFilesFromDiff(numstat string) []HeatedFileMetric {
	if numstat == "" {
		return nil
	}

	lines := strings.Split(numstat, "\n")
	var metrics []HeatedFileMetric

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		ins, _ := strconv.Atoi(parts[0])
		del, _ := strconv.Atoi(parts[1])
		path := parts[2]
		totalChanges := ins + del

		metrics = append(metrics, HeatedFileMetric{
			Path:         path,
			ChangesCount: totalChanges,
			Insertions:   ins,
			Deletions:    del,
			Description:  summarizeFileChurn(path, ins, del),
		})
	}

	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].ChangesCount > metrics[j].ChangesCount
	})

	if len(metrics) > 5 {
		metrics = metrics[:5]
	}

	return metrics
}

func summarizeFileChurn(path string, ins, del int) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".go", ".ts", ".rs", ".py":
		if ins > del*2 {
			return "Feature expansion and new logic implementation"
		} else if del > ins*2 {
			return "Refactoring, pruning, and dead code elimination"
		}
		return "Iterative feature enhancements and logic adjustments"
	case ".md":
		return "Documentation and architecture updates"
	case ".json", ".yaml", ".yml":
		return "Configuration, schema, and dependency updates"
	default:
		return "Codebase maintenance"
	}
}

func renderSingleRepoSummaryTerminal(name, url, branch, head string, releases []ReleaseSummaryRecord, cacheHits int, dur time.Duration) {
	fmt.Printf("[GitMap] Repository: %s (%s)\n", name, url)
	shortHead := head
	if len(shortHead) > 8 {
		shortHead = shortHead[:8]
	}
	fmt.Printf("Branch: %s | HEAD: %s\n\n", branch, shortHead)

	if len(releases) == 0 {
		fmt.Println("No releases or tags discovered in repository.")
		return
	}

	for _, rel := range releases {
		checkMark := pterm.Green("✓")
		cacheTag := ""
		if rel.IsCacheHit {
			cacheTag = " (cached)"
		}

		fmt.Printf("%s Release %s (%s)%s\n", checkMark, pterm.Bold.Sprint(rel.TagName), rel.ReleaseDate, cacheTag)
		fmt.Printf("  Summary: %s\n", rel.SummaryGist)

		if len(rel.HeatedFiles) > 0 {
			fmt.Println("  Heated Files:")
			for _, h := range rel.HeatedFiles {
				fmt.Printf("    - %s (+%d, -%d): %s\n", h.Path, h.Insertions, h.Deletions, h.Description)
			}
		}
		fmt.Println()
	}

	fmt.Printf("[GitMap] Completed in %.2fs (Cache: %d/%d hits).\n", dur.Seconds(), cacheHits, len(releases))
}
