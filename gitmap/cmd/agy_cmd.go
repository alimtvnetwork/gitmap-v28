// Package cmd — agy_cmd.go is the root command for Antigravity workspace management.
package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// AgyCmd is the root agy command
var AgyCmd = &cobra.Command{
	Use:   "agy",
	Short: "Antigravity CLI Management",
}

func dispatchAgy(ctx context.Context, args []string, root *cobra.Command) error {
	if len(args) > 0 && (args[0] == "agy" || args[0] == "ag" || args[0] == "antigravity") {
		args = args[1:]
	}
	if len(args) > 0 {
		args[0] = normalizeAgySubcommand(args[0])
	}
	if len(args) > 0 && isAgyFindDuplicatesArg(args[0]) {
		return runFindDuplicates("agy", args[1:])
	}
	if isAgyLsEmptyConvsArg(args) {
		return runAgyLsEmptyConvs(args[1:])
	}
	AgyCmd.SetArgs(args)
	return AgyCmd.ExecuteContext(ctx)
}

func normalizeAgySubcommand(sub string) string {
	low := strings.ToLower(sub)
	if isCureDupsAlias(low) {
		return "optimize-projects"
	}
	if isRemoveMissingAlias(low) {
		return "remove-missing-projects"
	}
	if isReadMemoryAlias(low) {
		return "all-projects-read-memory-prompt"
	}
	if low == "reconcile" || low == "recon" || low == "reconcile-projects" {
		return "reconcile"
	}
	if low == "find-duplicate-projects" || low == "fdp" {
		return "find-duplicate-projects"
	}
	return sub
}

func isCureDupsAlias(low string) bool {
	return low == "--repeat-fix" || low == "-r" ||
		low == "cure-duplicate-projects" || low == "cdp" ||
		low == "cure-duplicates" || low == "cure-duplicate"
}

func isRemoveMissingAlias(low string) bool {
	return low == "remove-misisng-projects" || low == "remove-missing-projects" ||
		low == "rm-missing-projects" || low == "rm-missing" || low == "clean-missing"
}

func isReadMemoryAlias(low string) bool {
	return low == "all-projects-read-memory-prompt" || low == "aprmp" ||
		low == "read-memory-all" || low == "rm-all-prompt"
}

func isAgyFindDuplicatesArg(sub string) bool {
	low := strings.ToLower(sub)
	return low == "find-duplicates" || low == "duplicates" || low == "dups" || low == "find-dups"
}

func isAgyLsEmptyConvsArg(args []string) bool {
	if len(args) == 0 {
		return false
	}
	if args[0] == "show-projects-with-empty-conversations" || args[0] == "show-proects-with-empty-conversations" {
		return true
	}
	if args[0] == "ls" && len(args) > 1 {
		sub := strings.ToLower(args[1])
		return sub == "show-projects-with-empty-conversations" ||
			sub == "show-proects-with-empty-conversations" ||
			sub == "empty-conversations" ||
			sub == "--empty-conversations" ||
			sub == "empty-convs"
	}
	return false
}

func init() {
	AgyCmd.AddCommand(agyAddCmd)
	AgyCmd.AddCommand(agyRmCmd)
	AgyCmd.AddCommand(agyLsCmd)
	AgyCmd.AddCommand(agyStatusCmd)
	AgyCmd.AddCommand(agyOptimizeCmd)
	AgyCmd.AddCommand(agyScanCmd)
	AgyCmd.AddCommand(agyStatsCmd)
	AgyCmd.AddCommand(agyUpdateCmd)
	AgyCmd.AddCommand(agyClearCmd)
	AgyCmd.AddCommand(agyOpenCmd)
	AgyCmd.AddCommand(agyPromptCmd)
	AgyCmd.AddCommand(agyRwCmd)
	AgyCmd.AddCommand(agySyncCmd)
	AgyCmd.AddCommand(agyPapCmd)
	AgyCmd.AddCommand(agyExportCmd)
	AgyCmd.AddCommand(agyImportCmd)
	AgyCmd.AddCommand(agyPluginsCmd)
	AgyCmd.AddCommand(agyRemoveEmptyConvsCmd)
	AgyCmd.AddCommand(agyFindDupsCmd)
	AgyCmd.AddCommand(agyRemoveMissingCmd)
	AgyCmd.AddCommand(agyReconcileCmd)
	AgyCmd.AddCommand(agyAllProjectsReadMemoryCmd)
	AgyCmd.AddCommand(agyGroupCmd)
	AgyCmd.AddCommand(agyUndoCmd)
	AgyCmd.AddCommand(agyRedoCmd)
	AgyCmd.AddCommand(agySettingsCmd)
	initPlugins()
	initAgyGroup()
	initAgySettings()
}

func getProjectsDirPath() (string, error) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", homeErr
	}
	return filepath.Join(homeDir, ".gemini", "config", "projects"), nil
}

func ensureDirExists(dirPath string) bool {
	mkErr := os.MkdirAll(dirPath, 0755)
	return mkErr == nil
}
