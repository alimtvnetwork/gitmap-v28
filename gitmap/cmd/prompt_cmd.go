// Package cmd — prompt_cmd.go defines cobra commands for prompt templates.
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// PromptCmd is the root command for managing AI prompt templates.
var PromptCmd = &cobra.Command{
	Use:     "prompt",
	Aliases: []string{"prompts", "pmt"},
	Short:   "Manage, import, export, and inject AI prompt templates",
}

var promptLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all available prompt templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPromptList()
	},
}

var promptShowCmd = &cobra.Command{
	Use:   "show [slug]",
	Short: "Show details and body of a prompt template",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: gitmap prompt show <slug>")
		}

		return runPromptShow(args[0])
	},
}

var promptAddCmd = &cobra.Command{
	Use:   "add [slug] [file.md]",
	Short: "Add or update a prompt template from a markdown file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: gitmap prompt add <slug> <file.md>")
		}

		return runPromptAdd(args[0], args[1])
	},
}

var promptRmCmd = &cobra.Command{
	Use:     "rm [slug]",
	Aliases: []string{"delete", "del"},
	Short:   "Remove a prompt template",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: gitmap prompt rm <slug>")
		}

		return runPromptRm(args[0])
	},
}

func initPromptCmd() {
	PromptCmd.AddCommand(promptLsCmd)
	PromptCmd.AddCommand(promptShowCmd)
	PromptCmd.AddCommand(promptAddCmd)
	PromptCmd.AddCommand(promptRmCmd)
	PromptCmd.AddCommand(promptExportCmd)
	PromptCmd.AddCommand(promptImportCmd)
	PromptCmd.AddCommand(promptInjectCmd)
}

func dispatchPrompt(ctx context.Context, args []string, root *cobra.Command) error {
	if len(args) > 0 && (args[0] == "prompt" || args[0] == "prompts" || args[0] == "pmt") {
		args = args[1:]
	}

	PromptCmd.SetArgs(args)
	return PromptCmd.ExecuteContext(ctx)
}

func init() {
	initPromptCmd()
}
