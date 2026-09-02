// Package cmd — agy_scan.go scans directory trees and cross-references Antigravity projects.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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
		return runAgyScan(args)
	},
}

func runAgyScan(args []string) error {
	rootPath := resolveAgyScanRoot(args)
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	repos := discoverGitRepos(rootPath)
	results := matchReposWithAgyProjects(repos, projects)
	renderAgyScanResults(rootPath, results)
	return nil
}

func resolveAgyScanRoot(args []string) string {
	if len(args) > 0 && args[0] != "" {
		abs, err := filepath.Abs(args[0])
		if err == nil {
			return abs
		}
	}
	cwd, _ := os.Getwd()
	return cwd
}

func discoverGitRepos(root string) []string {
	var repos []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		name := info.Name()
		if shouldSkipScanDir(name) {
			return filepath.SkipDir
		}
		if checkDirExists(filepath.Join(path, ".git")) {
			repos = append(repos, path)
			return filepath.SkipDir
		}
		return nil
	})
	return repos
}

func shouldSkipScanDir(name string) bool {
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
		normRepo := strings.ToLower(filepath.Clean(r))
		matchedIDs := make([]string, 0)
		for _, p := range projects {
			normProj := strings.ToLower(filepath.Clean(p.GetPath()))
			if normProj == normRepo {
				matchedIDs = append(matchedIDs, shortProjectId(p.ID))
			}
		}
		results = append(results, agyScanRepoResult{
			Name:       filepath.Base(r),
			Path:       r,
			MatchCount: len(matchedIDs),
			ProjectIDs: matchedIDs,
		})
	}
	return results
}

func renderAgyScanResults(root string, results []agyScanRepoResult) {
	fmt.Printf("\n  %s╔══════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║       antigravity repo scan          ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Scanned: %s\n", root)
	fmt.Println("  ──────────────────────────────────────────────────────────────────────────────────────")
	fmt.Printf("  %-30s  %-15s  %s\n", "REPOSITORY", "AGY STATUS", "TARGET PATH")
	fmt.Println("  ──────────────────────────────────────────────────────────────────────────────────────")
	added, repeated, missing := 0, 0, 0
	for _, res := range results {
		statusStr := formatAgyScanStatus(res.MatchCount)
		fmt.Printf("  %-30s  %-24s  %s\n", res.Name, statusStr, res.Path)
		tallyScanResult(res.MatchCount, &added, &repeated, &missing)
	}
	printAgyScanSummary(len(results), added, repeated, missing)
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
	if count == 1 {
		*added++
	} else if count > 1 {
		*repeated++
	} else {
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
	fmt.Println()
}
