// Package cmd — agy_scan.go scans directory trees and cross-references Antigravity projects.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type agyScanRepoResult struct {
	Name       string
	Path       string
	MatchCount int
	ProjectIDs []string
}

var agyScanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan path recursively for git repos and check Antigravity status",
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgyScan(args)
		if appErr != nil {
			return appErr
		}

		return nil
	},
}

func runAgyScan(args []string) *apperror.AppError {
	rootPath := resolveAgyScanRoot(args)
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}

	return executeScanOnPath(rootPath, dirPath)
}

func executeScanOnPath(rootPath, dirPath string) *apperror.AppError {
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}

	repos := discoverGitRepos(rootPath)
	results := matchReposWithAgyProjects(repos, projects)
	renderAgyScanResults(rootPath, results)
	backupScannedPrompts()

	return nil
}

func resolveAgyScanRoot(args []string) string {
	if len(args) > 0 && args[0] != "" {
		return resolveTargetAbs(args[0])
	}

	cwd, _ := os.Getwd()

	return cwd
}

func resolveTargetAbs(raw string) string {
	abs, err := filepath.Abs(raw)
	if err == nil {
		return abs
	}

	cwd, _ := os.Getwd()

	return cwd
}

func discoverGitRepos(root string) []string {
	var repos []string
	_ = filepath.Walk(root, makeScanWalkFunc(&repos))

	return repos
}

func makeScanWalkFunc(repos *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		return handleScanWalkEntry(path, info, err, repos)
	}
}

func handleScanWalkEntry(path string, info os.FileInfo, err error, repos *[]string) error {
	if err != nil || info == nil {
		return nil
	}
	if !info.IsDir() {
		return nil
	}
	if isSkippableScanDir(info.Name()) {
		return filepath.SkipDir
	}

	return checkGitRepoEntry(path, repos)
}

func checkGitRepoEntry(path string, repos *[]string) error {
	gitDir := filepath.Join(path, ".git")
	if checkDirExists(gitDir) {
		*repos = append(*repos, path)

		return filepath.SkipDir
	}

	return nil
}

func isSkippableScanDir(name string) bool {
	skipNames := []string{".git", "node_modules", "vendor", ".gemini", "dist", "bin", "temp", "obj"}
	for _, skip := range skipNames {
		if strings.EqualFold(name, skip) {
			return true
		}
	}

	return false
}

func matchReposWithAgyProjects(repos []string, projects []AgyProject) []agyScanRepoResult {
	results := make([]agyScanRepoResult, 0, len(repos))
	for _, r := range repos {
		results = append(results, matchSingleRepoWithProjects(r, projects))
	}

	return results
}

func matchSingleRepoWithProjects(repo string, projects []AgyProject) agyScanRepoResult {
	normRepo := strings.ToLower(filepath.Clean(repo))
	matchedIDs := findMatchedProjectIDs(normRepo, projects)

	return agyScanRepoResult{
		Name:       filepath.Base(repo),
		Path:       repo,
		MatchCount: len(matchedIDs),
		ProjectIDs: matchedIDs,
	}
}

func findMatchedProjectIDs(normRepo string, projects []AgyProject) []string {
	matched := make([]string, 0)
	for _, p := range projects {
		normProj := strings.ToLower(filepath.Clean(p.GetPath()))
		if normProj == normRepo {
			matched = append(matched, shortProjectId(p.ID))
		}
	}

	return matched
}

func renderAgyScanResults(root string, results []agyScanRepoResult) {
	printAgyScanBanner(root)
	added, repeated, missing := 0, 0, 0
	for _, res := range results {
		renderScanRow(res)
		tallyScanResult(res.MatchCount, &added, &repeated, &missing)
	}
	printAgyScanSummary(len(results), added, repeated, missing)
}

func printAgyScanBanner(root string) {
	fmt.Printf("\n  %s╔══════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║       antigravity repo scan          ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Scanned: %s\n", root)
	fmt.Println("  ──────────────────────────────────────────────────────────────────────────────────────")
	fmt.Printf("  %-30s  %-15s  %s\n", "REPOSITORY", "AGY STATUS", "TARGET PATH")
	fmt.Println("  ──────────────────────────────────────────────────────────────────────────────────────")
}

func renderScanRow(res agyScanRepoResult) {
	statusStr := formatAgyScanStatus(res.MatchCount)
	fmt.Printf("  %-30s  %-24s  %s\n", res.Name, statusStr, res.Path)
}

func formatAgyScanStatus(count int) string {
	if count == 1 {
		return constants.ColorGreen + "✔ added (1)" + constants.ColorReset
	}
	if count > 1 {
		return fmt.Sprintf("%s⚠ repeated (%d)%s", constants.ColorYellow, count, constants.ColorReset)
	}

	return constants.ColorRed + "✖ not added" + constants.ColorReset
}

func tallyScanResult(count int, added, repeated, missing *int) {
	switch {
	case count == 1:
		*added++
	case count > 1:
		*repeated++
	default:
		*missing++
	}
}

func printAgyScanSummary(total, added, repeated, missing int) {
	fmt.Println("  ──────────────────────────────────────────────────────────────────────────────────────")
	fmt.Printf("  %d git repos found · %s%d added%s · %s%d repeated%s · %s%d not added%s\n",
		total,
		constants.ColorGreen, added, constants.ColorReset,
		constants.ColorYellow, repeated, constants.ColorReset,
		constants.ColorRed, missing, constants.ColorReset,
	)
	if repeated > 0 {
		fmt.Println("  Tip: Run 'gitmap agy optimize-projects' to remove duplicate projects.")
	}
	printPromptScanSummary(computePromptScanStats())
	fmt.Println()
}

func backupScannedPrompts() {
	allPrompts := CollectAllPrompts()
	if len(allPrompts) == 0 {
		return
	}

	backupPath := resolvePromptBackupPath()
	if backupPath == "" {
		return
	}

	data, err := json.MarshalIndent(allPrompts, "", "  ")
	if err != nil {
		return
	}

	writePromptBackupFile(backupPath, data, len(allPrompts))
}

func writePromptBackupFile(backupPath string, data []byte, count int) {
	writeErr := os.WriteFile(backupPath, data, 0644)
	if writeErr == nil {
		fmt.Printf("  %s✔ Backed up %d project prompts to %s%s\n", constants.ColorGreen, count, backupPath, constants.ColorReset)
	}
}

func resolvePromptBackupPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	dir := filepath.Join(home, ".gemini", "config", "backup", "prompts")
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		return ""
	}

	return filepath.Join(dir, fmt.Sprintf("prompts-backup-%d.json", time.Now().Unix()))
}
