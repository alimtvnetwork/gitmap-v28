package cmdautomation

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunPlanConsolidate clusters completed plans and synchronizes memory plans index.
func RunPlanConsolidate(opts PlanConsolidateOptions) PlanConsolidateResultMonad {
	start := time.Now()
	plansDir := resolvePlansDir(opts.Dir)
	completedDir := filepath.Join(plansDir, "completed")
	pendingDir := filepath.Join(plansDir, "pending")
	subtasksDir := filepath.Join(plansDir, "subtasks")

	completedFiles := collectPlanFiles(completedDir)
	pendingFiles := collectPlanFiles(pendingDir)
	subtaskDirs := collectSubtaskDirs(subtasksDir)

	threshold := resolveThreshold(opts.Threshold)
	clusters := clusterCompletedPlans(completedFiles, threshold)
	res := buildConsolidateResult(pendingFiles, completedFiles, subtaskDirs, clusters)
	appErr := applyConsolidation(plansDir, pendingFiles, completedFiles, opts.IsDryRun)
	if appErr != nil {
		return result.Fail[PlanConsolidateResult](appErr)
	}
	res.Duration = time.Since(start)
	res.IsSuccess = true
	return result.Ok(res)
}

func resolveThreshold(threshold int) int {
	if threshold > 0 {
		return threshold
	}
	return 5
}

func resolvePlansDir(dir string) string {
	if len(dir) > 0 {
		return dir
	}
	return filepath.Join(findRepoRoot(), ".ai-memory", "plans")
}

func collectPlanFiles(dir string) []string {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files
	}
	for _, e := range entries {
		isMd := strings.HasSuffix(strings.ToLower(e.Name()), ".md")
		if !e.IsDir() && isMd {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	return files
}

func collectSubtaskDirs(dir string) []string {
	var dirs []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return dirs
	}
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	return dirs
}

func clusterCompletedPlans(files []string, threshold int) []PlanCluster {
	var clusters []PlanCluster
	hasEnoughFiles := len(files) >= threshold
	if !hasEnoughFiles {
		return clusters
	}
	clusterMap := groupPlansByDomain(files)
	for domain, pList := range clusterMap {
		clusters = append(clusters, PlanCluster{
			Id:           fmt.Sprintf("cluster-%02d", len(clusters)+1),
			Title:        fmt.Sprintf("Milestone: %s", formatDomainTitle(domain)),
			Domain:       domain,
			MergedPlans:  pList,
			SubtaskCount: len(pList),
		})
	}
	return clusters
}

func formatDomainTitle(domain string) string {
	if domain == "" {
		return ""
	}
	r := []rune(domain)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func groupPlansByDomain(files []string) map[string][]string {
	grouped := make(map[string][]string)
	for _, f := range files {
		domain := extractDomainFromPlanName(f)
		grouped[domain] = append(grouped[domain], f)
	}
	return grouped
}

func extractDomainFromPlanName(name string) string {
	clean := strings.TrimSuffix(name, ".md")
	parts := strings.Split(clean, "-")
	hasPrefix := len(parts) > 1 && isNumericPrefix(parts[0])
	if hasPrefix {
		return parts[1]
	}
	if len(parts) > 0 {
		return parts[0]
	}
	return "general"
}

func isNumericPrefix(s string) bool {
	for _, c := range s {
		isDigit := c >= '0' && c <= '9'
		if !isDigit {
			return false
		}
	}
	return len(s) > 0
}

func buildConsolidateResult(pending, completed, subtasks []string, clusters []PlanCluster) PlanConsolidateResult {
	return PlanConsolidateResult{
		TotalPlans:      len(pending) + len(completed),
		CompletedPlans:  len(completed),
		ClustersCount:   len(clusters),
		SubtasksCleaned: len(subtasks),
		Clusters:        clusters,
	}
}

func applyConsolidation(plansDir string, pending, completed []string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}
	return syncPlansIndexFile(plansDir, pending, completed)
}

func syncPlansIndexFile(plansDir string, pending, completed []string) *apperror.AppError {
	indexPath := filepath.Join(plansDir, "01-index.md")
	lines := buildIndexLines(pending, completed)
	content := strings.Join(lines, "\n")
	err := os.WriteFile(indexPath, []byte(content), 0o644)
	if err != nil {
		return apperror.WrapSimple(err, "write 01-index.md")
	}
	return nil
}

func buildIndexLines(pending, completed []string) []string {
	lines := []string{
		"# Plans Index",
		"",
		"Master directory of architectural and execution plans.",
		"",
		"## Pending Plans",
		"",
	}
	for _, p := range pending {
		lines = append(lines, fmt.Sprintf("- [%s](pending/%s)", p, p))
	}
	lines = append(lines, "", "## Completed Plans", "")
	for _, c := range completed {
		lines = append(lines, fmt.Sprintf("- [%s](completed/%s)", c, c))
	}
	lines = append(lines, "")
	return lines
}
