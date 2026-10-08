// Package cmd — macro_root_dispatch.go enables dynamic root-level execution of saved macros.
package cmdmacro

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrun"
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

func DispatchMacroDynamic(command string, shouldAudit bool, auditID int64, auditStart time.Time) bool {
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

	ExecuteDynamicMacroWithArgs(command, os.Args[2:])
	finishCommandAudit(shouldAudit, auditID, auditStart, 0, "", 0)

	return true
}

func isExcludedRootCommand(command string) bool {
	return strings.HasPrefix(command, "-") || command == ""
}

func isMacroHelpToken(token string) bool {
	return token == "-h" || token == "--help" || token == "help"
}

// runMacroRootRun handles 'gitmap run <file-or-macro> [args...]'.
func RunMacroRootRun(args []string) error {
	if len(args) == 0 {
		cmdrun.PrintRunHelp()

		return apperror.NewValidationError("missing required target file or macro name for run command")
	}

	if isMacroHelpToken(args[0]) {
		cmdrun.PrintRunHelp()

		return nil
	}

	target := args[0]
	if isRunSubcommand(target) || cmdrun.CanResolveTarget(target) {
		return cmdrun.Run(args)
	}

	if isSavedMacro(target) {
		ExecuteDynamicMacroWithArgs(target, args[1:])

		return nil
	}

	cmdrun.PrintRunHelp()

	return apperror.NewValidationError(fmt.Sprintf("target '%s' is neither an executable script nor a saved macro", target))
}

func isRunSubcommand(target string) bool {
	lower := strings.ToLower(target)

	return lower == "errors" || lower == "err" || lower == "history" || lower == "clear-errors"
}

func isSavedMacro(target string) bool {
	m, err := macro.LoadMacroPolymorphic(target)

	return err == nil && m != nil
}

// runMacroRootRunUntil handles 'gitmap run-until <macro-name> [flags]'.
func RunMacroRootRunUntil(args []string) error {
	if len(args) == 0 {
		return RunMacroRootRun(args)
	}

	forwardArgs := append([]string{args[0], "--run-until"}, args[1:]...)

	return RunMacroRootRun(forwardArgs)
}
