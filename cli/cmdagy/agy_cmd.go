// Package cmd — agy_cmd.go is the root command for Antigravity workspace management.
package cmdagy

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

// DispatchAgy routes CLI arguments to agy commands.
func DispatchAgy(ctx context.Context, args []string, root *cobra.Command) error {
	args = stripAgyPrefix(args)
	if len(args) > 0 && isAgyOpenPathArg(args[0]) {
		return RunAgyOpen(args[0])
	}
	args = normalizeAgyArgs(args)
	if len(args) > 0 && isAgyFindDuplicatesArg(args[0]) {
		return RunFindDuplicates()
	}
	if isAgyLsEmptyConvsArg(args) {
		return runAgyLsEmptyConvs(args[1:])
	}
	AgyCmd.SetArgs(args)

	return AgyCmd.ExecuteContext(ctx)
}

func stripAgyPrefix(args []string) []string {
	if len(args) > 0 && (args[0] == "agy" || args[0] == "ag" || args[0] == "antigravity") {
		return args[1:]
	}

	return args
}

func normalizeAgyArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	if isCompoundAgyFix(args) {
		return rewriteCompoundAgyFix(args)
	}

	args[0] = normalizeAgySubcommand(args[0])

	return args
}

func isCompoundAgyFix(args []string) bool {
	if len(args) < 2 {
		return false
	}

	first := strings.ToLower(args[0])
	second := strings.ToLower(args[1])
	if first == "errors" || first == "error" || first == "err" {
		return second == "fix" || second == "aef"
	}
	if first == "fix" {
		return second == "errors" || second == "error" || second == "pipeline" || second == "agy"
	}

	return false
}

func rewriteCompoundAgyFix(args []string) []string {
	return append([]string{"fix-pipeline"}, args[2:]...)
}

func isAgyOpenPathArg(arg string) bool {
	if arg == "." || arg == ".." || strings.HasPrefix(arg, "./") || strings.HasPrefix(arg, ".\\") {
		return true
	}
	if strings.HasPrefix(arg, "/") || strings.HasPrefix(arg, "\\") || strings.Contains(arg, string(filepath.Separator)) {
		return true
	}
	if info, err := os.Stat(arg); err == nil && info.IsDir() {
		return true
	}

	return false
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

	if isRprpAlias(low) {
		return "read-all-projects-with-read-prompts"
	}

	if low == "reconcile" || low == "recon" || low == "reconcile-projects" {
		return "reconcile"
	}

	if low == "find-duplicate-projects" || low == "fdp" {
		return "find-duplicate-projects"
	}

	if low == "pin-projects" || low == "pin-project" || low == "pinned-projects" || low == "pinned" || low == "pins" {
		return "pin-projects"
	}

	if low == "install" || low == "in" || low == "i" {
		return "install"
	}

	if low == "clean-cache" || low == "cleancache" || low == "clean_cache" || low == "cc" {
		return "clean-cache"
	}

	if isFixPipelineAlias(low) {
		return "fix-pipeline"
	}

	if low == "list-prompts" || low == "listprompts" || low == "lp" || low == "list-prompt" {
		return "list-prompts"
	}

	if low == "rerun" || low == "replay" || low == "rr" {
		return "rerun"
	}

	return sub
}

func isFixPipelineAlias(low string) bool {
	return low == "fix-pipeline" || low == "fix" || low == "fp" ||
		low == "pipeline-fix" || low == "fixpipeline" ||
		low == "aef" || low == "agy-errors-fix" || low == "errors-fix" ||
		low == "fix-errors" || low == "fix-agy"
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

func isRprpAlias(low string) bool {
	return low == "read-all-projects-with-read-prompts" || low == "rprp" ||
		low == "rapwrp" || low == "read-all-with-prompts"
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
	AgyCmd.AddCommand(agyReadAllProjectsWithReadPromptsCmd)
	AgyCmd.AddCommand(agyGroupCmd)
	AgyCmd.AddCommand(agyUndoCmd)
	AgyCmd.AddCommand(agyRedoCmd)
	AgyCmd.AddCommand(agySettingsCmd)
	AgyCmd.AddCommand(agyPinProjectsCmd)
	AgyCmd.AddCommand(agyCleanCacheCmd)
	AgyCmd.AddCommand(agyFixPipelineCmd)
	AgyCmd.AddCommand(agyRerunCmd)
	AgyCmd.AddCommand(agyListPromptsCmd)
	initPlugins()
	initAgyGroup()
	initAgySettings()
	initAgyPinProjects()
	AgyCmd.SetHelpFunc(renderAgyHelp)
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
