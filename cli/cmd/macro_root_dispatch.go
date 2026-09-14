// Package cmd — macro_root_dispatch.go enables dynamic root-level execution of saved macros.
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func handleMacroLoadFailure(command string, loadErr error) bool {
	if macro.IsMacroNotFound(loadErr) {
		return false
	}

	cliexit.Reportf(command, "macro-load", "", loadErr)
	cliexit.HandleError(loadErr, 1)

	return true
}

func dispatchMacroDynamic(command string, shouldAudit bool, auditID int64, auditStart time.Time) bool {
	if isExcludedRootCommand(command) {
		return false
	}

	m, loadErr := macro.LoadMacroPolymorphic(command)
	if loadErr != nil {
		return handleMacroLoadFailure(command, loadErr)
	}

	if m == nil {
		return false
	}

	cmdmacro.ExecuteDynamicMacroWithArgs(command, os.Args[2:])
	finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

	return true
}

func isExcludedRootCommand(command string) bool {
	return strings.HasPrefix(command, "-") || command == ""
}

func isMacroRunHelp(args []string) bool {
	return len(args) == 0 || args[0] == "-h" || args[0] == "--help"
}

// runMacroRootRun handles 'gitmap run <macro-name> [flags]'.
func runMacroRootRun(args []string) error {
	if !isMacroRunHelp(args) {
		cmdmacro.ExecuteDynamicMacroWithArgs(args[0], args[1:])

		return nil
	}

	fmt.Fprintf(os.Stderr, "Usage: gitmap run <macro-name> [--run-until] [--no-terminal] [--summary] [--no-tree]\n")
	if len(args) == 0 {
		return apperror.NewValidationError("missing required macro name for run command")
	}

	return nil
}

// runMacroRootRunUntil handles 'gitmap run-until <macro-name> [flags]'.
func runMacroRootRunUntil(args []string) error {
	if len(args) == 0 {
		return runMacroRootRun(args)
	}

	forwardArgs := append([]string{args[0], "--run-until"}, args[1:]...)

	return runMacroRootRun(forwardArgs)
}
