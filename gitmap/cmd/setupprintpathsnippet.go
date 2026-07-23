package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/setup"
)

// runPrintPathSnippet handles `gitmap setup print-path-snippet`.
//
// Renders the canonical marker-block PATH snippet for the requested
// shell to stdout. Shell scripts (run.sh, gitmap/scripts/install.sh)
// shell out to this command to obtain byte-identical snippet text,
// guaranteeing single-source-of-truth across all three drivers.
//
// Spec: spec/04-generic-cli/21-post-install-shell-activation/02-snippets.md
func runPrintPathSnippet(args []string) {
	shell, dir, manager := parsePrintPathSnippetFlags(args)
	out, err := setup.RenderPathSnippet(shell, dir, manager)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	// Trailing newline so callers can `>>` straight into a profile file.
	fmt.Println(out)
}

// parsePrintPathSnippetFlags parses --shell --dir --manager.
func parsePrintPathSnippetFlags(args []string) (shell, dir, manager string) {
	fs := flag.NewFlagSet("setup print-path-snippet", flag.ExitOnError)
	shellFlag := fs.String("shell", constants.PathSnippetShellBash, constants.FlagDescPathSnippetShell)
	dirFlag := fs.String("dir", "", constants.FlagDescPathSnippetDir)
	managerFlag := fs.String("manager", "gitmap setup", constants.FlagDescPathSnippetManager)
	fs.Parse(args)

	return *shellFlag, *dirFlag, *managerFlag
}
