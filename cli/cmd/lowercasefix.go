// Package cmd — lowercasefix.go defines the lowercase CLI command and aliases.
package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	lcfDryRun   bool
	lcfNoCommit bool
	lcfYes      bool
	lcfReadme   bool
	lcfMessage  string
)

var lowerCaseFixCmd = &cobra.Command{
	Use:     "lowercase [patterns...]",
	Aliases: []string{"lower-case-fix", "lowercase-fix", "lcf", "lower", "lower-case-readme", "lowercase-readme", "readme-lower", "readme-lowercase", "lcr", "lc-fix"},
	Short:   "Rename uppercase files (e.g. README.md -> readme.md) across the repo and commit",
	Long: `Scans the repository for uppercase or mixed-case files matching specified patterns
(e.g. *.md, *md, *, SKILL*) and safely renames them to lowercase using a two-step git mv,
preventing case-collision on case-insensitive filesystems (Windows/macOS), and optionally commits.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		patterns, isReadme := resolveLcfPatterns(args)
		opts := LowerCaseFixOptions{
			Patterns:      patterns,
			IsDryRun:      lcfDryRun,
			IsNoCommit:    lcfNoCommit,
			IsYes:         lcfYes,
			IsReadmeOnly:  lcfReadme || isReadme,
			CommitMessage: lcfMessage,
		}

		return ExecuteLowerCaseFix(opts)
	},
}

func init() {
	lowerCaseFixCmd.Flags().BoolVarP(&lcfDryRun, "dry-run", "d", false, "Preview matching files without renaming")
	lowerCaseFixCmd.Flags().BoolVar(&lcfNoCommit, "no-commit", false, "Do not commit renames to git")
	lowerCaseFixCmd.Flags().BoolVarP(&lcfYes, "yes", "y", false, "Proceed without interactive confirmation")
	lowerCaseFixCmd.Flags().BoolVarP(&lcfReadme, "readme", "r", false, "Target root README files only")
	lowerCaseFixCmd.Flags().StringVarP(&lcfMessage, "message", "m", "", "Custom git commit message")
}

func runLowerCaseFixCLI(args []string) error {
	lcfDryRun = false
	lcfNoCommit = false
	lcfYes = false
	lcfReadme = false
	lcfMessage = ""

	if len(os.Args) > 1 && isReadmeAlias(strings.ToLower(os.Args[1])) {
		lcfReadme = true
	}

	subArgs, isReadme := extractSubArgs(args)
	if isReadme {
		lcfReadme = true
	}
	checkHelp("lowercase", subArgs)
	lowerCaseFixCmd.SetArgs(subArgs)

	return lowerCaseFixCmd.ExecuteContext(context.Background())
}

func extractSubArgs(args []string) ([]string, bool) {
	if len(args) == 0 {
		return args, false
	}
	subcmd := strings.ToLower(args[0])
	if isReadmeAlias(subcmd) {
		return args[1:], true
	}
	if isLcfCmdToken(subcmd) {
		return args[1:], false
	}

	return args, false
}

func isReadmeAlias(s string) bool {
	return s == "lowercase-readme" || s == "lower-case-readme" ||
		s == "readme-lower" || s == "readme-lowercase" || s == "lcr"
}

func isLcfCmdToken(s string) bool {
	return s == "lowercase" || s == "lower" || s == "lower-case-fix" ||
		s == "lowercase-fix" || s == "lcf" || s == "lc-fix"
}

func resolveLcfPatterns(args []string) ([]string, bool) {
	if len(args) == 0 {
		return []string{"*"}, false
	}
	if len(args) == 1 && strings.EqualFold(args[0], "readme") {
		return []string{"readme*"}, true
	}

	return args, false
}
