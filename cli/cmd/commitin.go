package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/orchestrator"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func init() {
	commitin.SyncStateTemplatesHook = syncStateTemplatesFromFiles
	commitin.PrecompileStateCategoryHook = precompileStateTemplatesByCategory
}

func syncStateTemplatesFromFiles(importPaths []string) error {
	for _, p := range importPaths {
		_, _, _, _ = store.ImportTemplatesFromFile(p, false)
	}

	return nil
}

func precompileStateTemplatesByCategory(categoryOrSlug string, extraVars map[string]string) []string {
	compiled, err := store.PrecompileTemplates(categoryOrSlug, extraVars)
	if err != nil || len(compiled) == 0 {
		return nil
	}
	out := make([]string, 0, len(compiled))
	for _, item := range compiled {
		out = append(out, formatStateCompiledBlock(item))
	}

	return out
}

func formatStateCompiledBlock(item store.CompiledTemplate) string {
	title := strings.TrimSpace(item.Title)
	text := strings.TrimSpace(item.Text)
	if strings.HasPrefix(text, "#") || title == "" {
		return text
	}
	cleanTitle := strings.TrimSpace(strings.TrimLeft(title, "#"))

	return "# " + cleanTitle + "\n" + text
}

// runCommitIn is the top-level entry point for `gitmap commit-in` / `gitmap cin`.
func runCommitIn(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			fmt.Println(commitin.PrintCommitInHelp())
			cliexit.HandleError(nil, 0)
		}
	}

	raw, perr := commitin.Parse(args)
	if perr != nil {
		fmt.Fprintf(os.Stderr, constants.CommitInErrBadArgs, perr.Message)
		cliexit.HandleError(perr, constants.CommitInExitBadArgs)
	}

	if raw.IsTree {
		printCommitPullTree()
	}

	exitCode := orchestrator.Run(raw, os.Stdout, os.Stderr)
	cliexit.HandleError(nil, exitCode)

	return nil
}
