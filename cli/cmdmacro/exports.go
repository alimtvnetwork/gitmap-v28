package cmdmacro

import (
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

// MacroCmd dispatches the main macro command.
func MacroCmd(args []string) error {
	return runMacroCmd(args)
}

// RunMacroCmd dispatches macro commands.
func RunMacroCmd(args []string) error {
	return runMacroCmd(args)
}

// RunExecuteCmd executes a macro.
func RunExecuteCmd(args []string) error {
	return runExecuteCmd(args)
}

// HandleMacroAdd handles adding a macro.
func HandleMacroAdd(args []string) error {
	return handleMacroAdd(args)
}

// HandleMacroEdit handles editing a macro.
func HandleMacroEdit(args []string) error {
	return handleMacroEdit(args)
}

// HandleMacroList lists macros.
func HandleMacroList(args []string) error {
	return handleMacroList(args)
}

// HandleMacroRecord records a macro.
func HandleMacroRecord(args []string) error {
	return handleMacroRecord(args)
}

// HandleMacroShow displays a macro.
func HandleMacroShow(args []string) error {
	return handleMacroShow(args)
}

// HandleMacroDelete deletes a macro.
func HandleMacroDelete(args []string) error {
	return handleMacroDelete(args)
}

// RunMacroExport exports macros.
func RunMacroExport(args []string) error {
	return runMacroExport(args)
}

// RunMacroImport imports macros.
func RunMacroImport(args []string) error {
	return runMacroImport(args)
}

// RunMacroUntilSuccess retries a command until success.
func RunMacroUntilSuccess(args []string) error {
	return runMacroUntilSuccess(args)
}

// ExecuteDynamicMacro executes a dynamic macro by name.
func ExecuteDynamicMacro(name string) {
	opts := parseExecOptions(os.Args[2:])
	execErr := executeMacroByName(name, opts)

	if execErr != nil {
		cliexit.HandleError(execErr, 1)
	}
}

// RunCatCmd runs the cat command.
func RunCatCmd(args []string) error {
	return runCatCmd(args)
}

// RunTouchCmd runs the touch command.
func RunTouchCmd(args []string) error {
	return runTouchCmd(args)
}

// RunMkfileCmd runs the mkfile command.
func RunMkfileCmd(args []string) error {
	return runMkfileCmd(args)
}

// ParseExecOptions parses execution options.
func ParseExecOptions(args []string) macro.ExecOptions {
	return parseExecOptions(args)
}

// OutputStructuredData outputs data as JSON/YAML or saves to file.
func OutputStructuredData(data interface{}, opts macro.ExecOptions) error {
	return outputStructuredData(data, opts)
}

// ExtractMacroNameAndFlags extracts macro name and trailing flags.
func ExtractMacroNameAndFlags(args []string) (string, []string) {
	return extractMacroNameAndFlags(args)
}

// ParseDurationArg parses duration string with fallback.
func ParseDurationArg(val string, fallback time.Duration) time.Duration {
	return parseDurationArg(val, fallback)
}
