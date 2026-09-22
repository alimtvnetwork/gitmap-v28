// Package cmd — lowercasefix.go defines the lower-case-fix CLI command.
package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	lcfDryRun  bool
	lcfCommit  bool
	lcfYes     bool
	lcfMessage string
)

var lowerCaseFixCmd = &cobra.Command{
	Use:     "lower-case-fix [patterns...]",
	Aliases: []string{"lowercase-fix", "lcf", "lower-case-readme", "lc-fix"},
	Short:   "Rename uppercase files (e.g. README.md -> readme.md) across the repo and commit",
	Long: `Scans the repository for uppercase or mixed-case files matching specified patterns
(defaults to *.md / README.md) and safely renames them to lowercase using a two-step git mv,
preventing case-collision on case-insensitive filesystems (Windows/macOS), and optionally commits.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		patterns := resolveLcfPatterns(args)
		opts := LowerCaseFixOptions{
			Patterns:      patterns,
			IsDryRun:      lcfDryRun,
			IsCommit:      lcfCommit,
			IsYes:         lcfYes,
			CommitMessage: lcfMessage,
		}
		return ExecuteLowerCaseFix(opts)
	},
}

func init() {
	lowerCaseFixCmd.Flags().BoolVarP(&lcfDryRun, "dry-run", "d", false, "Preview matching files without renaming")
	lowerCaseFixCmd.Flags().BoolVarP(&lcfCommit, "commit", "c", true, "Automatically commit renames with git (default: true)")
	lowerCaseFixCmd.Flags().BoolVarP(&lcfYes, "yes", "y", false, "Proceed without interactive confirmation")
	lowerCaseFixCmd.Flags().StringVarP(&lcfMessage, "message", "m", "", "Custom git commit message")
}

func runLowerCaseFixCLI(args []string) error {
	lowerCaseFixCmd.SetArgs(args)
	return lowerCaseFixCmd.ExecuteContext(context.Background())
}

func resolveLcfPatterns(args []string) []string {
	if len(args) == 0 {
		return []string{"*.md"}
	}
	return args
}
