// Package cmd — agy_reconcile.go reconciles missing Antigravity projects by searching GitMap repo paths.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	agyReconcileDryRun bool
	agyReconcileYes    bool
)

var agyReconcileCmd = &cobra.Command{
	Use:     "reconcile",
	Aliases: []string{"recon", "reconcile-projects"},
	Short:   "Reconcile and re-link missing Antigravity project paths with tracked repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyReconcile()
	},
}

func init() {
	agyReconcileCmd.Flags().BoolVarP(&agyReconcileDryRun, "dry-run", "d", false, "Preview path reconciliations without saving")
	agyReconcileCmd.Flags().BoolVarP(&agyReconcileYes, "yes", "y", false, "Save reconciliations without interactive prompt")
}

func runAgyReconcile() error {
	dirPath, pathErr := getProjectsDirPath()
	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "path error")
	}
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load projects")
	}
	missing := findMissingAgyProjects(projects, "")
	if len(missing) == 0 {
		fmt.Printf("%s All Antigravity projects are active and paths exist on disk.\n",
			constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}
	return executeAgyReconcile(dirPath, missing)
}

func executeAgyReconcile(dirPath string, missing []AgyProject) error {
	fmt.Printf("\n  %sScanning tracked repositories to reconcile %d missing Antigravity project(s)...%s\n\n",
		constants.ColorCyan, len(missing), constants.ColorReset)
	reconciled := 0
	unresolved := make([]AgyProject, 0)
	for _, p := range missing {
		candidate := findCandidateRepoPath(p.Name, filepath.Base(p.GetPath()))
		if candidate != "" {
			reconcileSingleProject(dirPath, p, candidate)
			reconciled++
			continue
		}
		unresolved = append(unresolved, p)
	}
	printReconcileSummary(reconciled, len(unresolved), len(missing))
	return nil
}

func reconcileSingleProject(dirPath string, p AgyProject, candidate string) {
	printReconcileMatch(p.Name, p.GetPath(), candidate)
	if !agyReconcileDryRun {
		_ = saveReconciledProject(dirPath, p, candidate)
	}
}

func findCandidateRepoPath(name, baseName string) string {
	scanFile := filepath.Join(".gitmap", "output", "gitmap.json")
	records, err := loadStatusRecords(scanFile)
	if err != nil || len(records) == 0 {
		return ""
	}
	for _, r := range records {
		rBase := filepath.Base(r.AbsolutePath)
		matches := strings.EqualFold(rBase, baseName) || strings.EqualFold(r.RepoName, name)
		if matches && checkDirExists(r.AbsolutePath) {
			return r.AbsolutePath
		}
	}
	return ""
}

func printReconcileMatch(name, oldPath, newPath string) {
	fmt.Printf("  %s✓ Found Match:%s %s\n", constants.ColorGreen, constants.ColorReset, name)
	fmt.Printf("    Old Path: %s%s%s\n", constants.ColorRed, oldPath, constants.ColorReset)
	fmt.Printf("    New Path: %s%s%s\n\n", constants.ColorGreen, newPath, constants.ColorReset)
}

func saveReconciledProject(dirPath string, p AgyProject, newPath string) error {
	if p.ProjectResources == nil || len(p.ProjectResources.Resources) == 0 {
		return nil
	}
	gf := p.ProjectResources.Resources[0].GitFolder
	if gf == nil {
		return nil
	}
	gf.FolderURI = "file:///" + filepath.ToSlash(filepath.Clean(newPath))
	p.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	data, marshalErr := json.MarshalIndent(p, "", "  ")
	if marshalErr != nil {
		return apperror.WrapSimple(marshalErr, "marshal")
	}
	targetFile := filepath.Join(dirPath, p.ID+".json")
	return os.WriteFile(targetFile, data, 0644)
}

func printReconcileSummary(reconciled, unresolvedCount, totalMissing int) {
	fmt.Printf("  %s\n", strings.Repeat("─", 68))
	if reconciled > 0 {
		printReconciledStatus(reconciled)
	}
	if unresolvedCount > 0 {
		fmt.Printf("  %s %d project(s) remain unresolved (path not found in tracked repositories).\n",
			constants.ColorYellow+"⚠"+constants.ColorReset, unresolvedCount)
		fmt.Println("    To remove: gitmap agy remove-missing-projects")
		fmt.Println("    To re-scan: gitmap agy scan [path]")
	}
	fmt.Println()
}

func printReconciledStatus(reconciled int) {
	action := "Re-linked"
	if agyReconcileDryRun {
		action = "Identified for re-linking (dry-run)"
	}
	fmt.Printf("  %s %s %d project(s) successfully.\n",
		constants.ColorGreen+"✓"+constants.ColorReset, action, reconciled)
}
