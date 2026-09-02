// Package cmd — agy_cmd.go is the root command for Antigravity workspace management.
package cmd

import (
	"context"
	"os"
	"path/filepath"

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
	if len(args) > 0 && (args[0] == "--repeat-fix" || args[0] == "-r") {
		args[0] = "optimize-projects"
	}
	AgyCmd.SetArgs(args)
	return AgyCmd.ExecuteContext(ctx)
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
	initPlugins()
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
