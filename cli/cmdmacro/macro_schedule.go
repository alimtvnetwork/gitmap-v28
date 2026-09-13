package cmdmacro

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

// RunScheduleFn hooks into the scheduler subsystem without import cycles.
var RunScheduleFn func(args []string) error

func invokeSchedule(args []string) error {
	if RunScheduleFn != nil {
		return RunScheduleFn(args)
	}

	return fmt.Errorf("schedule runner not initialized")
}

func dispatchMacroSchedule(sub string, tail, all []string) error {
	switch sub {
	case "add", "create", "new":
		return runMacroScheduleAdd(tail)
	case "ls", "list":
		return invokeSchedule([]string{"ls"})
	case "rm", "remove", "delete":
		return invokeSchedule(append([]string{"rm"}, tail...))
	case "debug", "info":
		return invokeSchedule(append([]string{"debug"}, tail...))
	default:
		return invokeSchedule(all)
	}
}

func runMacroSchedule(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		printMacroScheduleHelp()

		return nil
	}

	return dispatchMacroSchedule(strings.ToLower(args[0]), args[1:], args)
}

func runMacroScheduleAdd(args []string) error {
	if len(args) == 0 {
		return apperror.New("macro", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify macro name to schedule: gitmap macro schedule add <macro-name> --every <interval>",
		})
	}

	macroName := args[0]
	if !macro.MacroExists(macroName) {
		return apperror.New("macro", "E_NOT_FOUND", map[string]any{
			"msg": fmt.Sprintf("Macro %q does not exist. Create it first with 'gitmap macro add %s'", macroName, macroName),
		})
	}

	scheduleArgs := []string{"add", macroName, "--macro", macroName}
	scheduleArgs = append(scheduleArgs, args[1:]...)

	return invokeSchedule(scheduleArgs)
}

func printMacroScheduleHelp() {
	fmt.Println("Usage: gitmap macro schedule <subcommand> [macro-name] [flags]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  add <macro-name> --every <interval>  Schedule a macro to run on a recurring interval")
	fmt.Println("  ls, list                             List all scheduled tasks and macros")
	fmt.Println("  rm <name>                            Unregister a scheduled macro or task")
	fmt.Println("  debug <name>                         Display execution diagnostics and split DB info")
	fmt.Println("\nSupported Script & Task Formats:")
	fmt.Println("  - Macros:        Saved GitMap macros (via --macro <name>)")
	fmt.Println("  - PowerShell:    PowerShell scripts (.ps1) or inline cmdlets")
	fmt.Println("  - Bash/Shell:    Bash scripts (.sh) or POSIX shell commands")
	fmt.Println("  - JavaScript:    Node.js scripts (.js)")
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap macro schedule add backup-daily --every 1d")
	fmt.Println("  gitmap macro schedule add sync-repos --every 2h")
	fmt.Println("  gitmap macro schedule ls")
	fmt.Println("  gitmap macro schedule debug backup-daily")
}
