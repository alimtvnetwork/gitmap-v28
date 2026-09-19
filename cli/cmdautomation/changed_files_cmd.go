package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	changedFilesCmd = &cobra.Command{
		Use:     "changed-files",
		Aliases: []string{"git-changes", "diff-files"},
		Short:   "Discovers modified, staged, and untracked files for targeted linter passes",
		RunE:    runChangedFilesCmd,
	}

	changedFilesOpts ChangedFilesOptions
)

func runChangedFilesCmd(cmd *cobra.Command, args []string) error {
	monad := RunChangedFiles(changedFilesOpts)
	if monad.IsFailure() {
		return monad.Err
	}
	renderChangedFilesResult(monad.Value, changedFilesOpts.IsJson)
	return nil
}

func renderChangedFilesResult(res ChangedFilesResult, isJson bool) {
	if isJson {
		printChangedFilesJson(res)
		return
	}
	printChangedFilesTerminal(res)
}

func printChangedFilesJson(res ChangedFilesResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}

func printChangedFilesTerminal(res ChangedFilesResult) {
	printChangedFilesHeader(res)
	if res.TotalFiles == 0 {
		fmt.Printf("%s✅ No changed or untracked files detected in working tree.%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	printChangedFilesList(res.Files)
}

func printChangedFilesHeader(res ChangedFilesResult) {
	fmt.Printf("\n%s[Git Changed Files Discovery]%s\n", constants.ColorBold, constants.ColorReset)
	if len(res.BaseRef) > 0 {
		fmt.Printf("  Base Reference:  %s%s%s\n", constants.ColorCyan, res.BaseRef, constants.ColorReset)
	}
	if len(res.HeadHash) > 0 {
		fmt.Printf("  HEAD Commit:     %s%s%s\n", constants.ColorCyan, truncateHash(res.HeadHash), constants.ColorReset)
	}
	fmt.Printf("  Total Files:     %s%d%s\n", constants.ColorGreen, res.TotalFiles, constants.ColorReset)
	fmt.Printf("  Duration:        %s\n\n", res.Duration)
}

func truncateHash(hash string) string {
	if len(hash) > 10 {
		return hash[:10]
	}
	return hash
}

func formatStatusColor(status string) string {
	switch status {
	case "added", "untracked":
		return constants.ColorGreen
	case "modified":
		return constants.ColorYellow
	case "deleted":
		return constants.ColorRed
	default:
		return constants.ColorCyan
	}
}

func printChangedFilesList(files []ChangedFileItem) {
	for i, f := range files {
		printSingleChangedFile(i+1, f)
	}
	fmt.Println()
}

func printSingleChangedFile(idx int, f ChangedFileItem) {
	color := formatStatusColor(f.Status)
	staged := formatStagedBadge(f.IsStaged)
	missing := formatMissingBadge(f.IsExists)
	fmt.Printf("  %3d. %s[%-9s]%s %s%s%s\n",
		idx, color, f.Status, constants.ColorReset, f.Path, staged, missing)
}

func formatStagedBadge(isStaged bool) string {
	if isStaged {
		return " [staged]"
	}
	return ""
}

func formatMissingBadge(isExists bool) string {
	if isExists == false {
		return " [missing on disk]"
	}
	return ""
}

func init() {
	AutomationCmd.AddCommand(changedFilesCmd)
	initChangedFilesFlags()
}

func initChangedFilesFlags() {
	changedFilesCmd.Flags().StringVar(&changedFilesOpts.Dir, "dir", ".", "Target repository directory")
	changedFilesCmd.Flags().StringVar(&changedFilesOpts.BaseRef, "base", "", "Git base reference or commit (e.g. origin/main, HEAD~1)")
	changedFilesCmd.Flags().IntVarP(&changedFilesOpts.Commits, "commits", "n", 20, "Number of recent commits to evaluate")
	changedFilesCmd.Flags().BoolVar(&changedFilesOpts.IsStagedOnly, "staged-only", false, "Discover staged files only")
	changedFilesCmd.Flags().BoolVar(&changedFilesOpts.IsVerify, "verify", false, "Verify that files currently exist on disk")
	changedFilesCmd.Flags().BoolVar(&changedFilesOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
