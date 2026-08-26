package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

// parseExecOptions parses runtime flags for macro execution.
func parseExecOptions(flagArgs []string) macro.ExecOptions {
	executionOptions := macro.ExecOptions{}
	for _, argument := range flagArgs {
		if argument == "--dry-run" {
			executionOptions.DryRun = true
		}
		if argument == "--verbose" || argument == "-v" {
			executionOptions.Verbose = true
		}
	}
	return executionOptions
}

// executeMacroByName loads and executes the named macro.
func executeMacroByName(macroName string, executionOptions macro.ExecOptions) {
	loadedMacro, loadErr := macro.LoadMacro(macroName)
	if loadErr != nil {
		fmt.Fprintf(os.Stderr, "%s✖ Error: %v%s\n", constants.ColorRed, loadErr, constants.ColorReset)
		os.Exit(1)
	}
	if len(loadedMacro.Steps) > 0 {
		printMacroStepsTree(loadedMacro)
	}
	if execErr := macro.Execute(context.Background(), loadedMacro, executionOptions); execErr != nil {
		os.Exit(1)
	}
}

// runExecuteCmd executes a recorded macro by name.
func runExecuteCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap execute <macro_name> [--dry-run] [--verbose]\n")
		os.Exit(1)
	}
	executionOptions := parseExecOptions(args[1:])
	executeMacroByName(args[0], executionOptions)
}

// runMacroCmd handles `gitmap macro` subcommands.
func runMacroCmd(args []string) {
	if len(args) == 0 {
		printMacroUsage()
		return
	}
	sub := args[0]
	switch sub {
	case "run", "exec":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: gitmap macro run <name>\n")
			os.Exit(1)
		}
		runExecuteCmd(args[1:])
	case "record", "rec":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: gitmap macro record <name>\n")
			os.Exit(1)
		}
		if err := macro.RecordInteractive(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "list", "ls":
		listMacros()
	case "show":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: gitmap macro show <name>\n")
			os.Exit(1)
		}
		showMacro(args[1])
	case "rm", "delete":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: gitmap macro rm <name>\n")
			os.Exit(1)
		}
		if err := macro.DeleteMacro(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✔ Removed macro %q\n", args[1])
	default:
		printMacroUsage()
	}
}

func listMacros() {
	macros, err := macro.ListMacros()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing macros: %v\n", err)
		return
	}
	if len(macros) == 0 {
		fmt.Println("  No saved macros found. Record one with: gitmap macro record <name>")
		return
	}
	fmt.Println()
	fmt.Printf("  %s%-24s %-12s %s%s\n", constants.ColorCyan, "MACRO NAME", "STEPS", "UPDATED", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	for _, m := range macros {
		fmt.Printf("  %-24s %-12d %s\n", m.Name, len(m.Steps), m.UpdatedAt.Format("2006-01-02 15:04"))
	}
	fmt.Println()
}

func showMacro(name string) {
	m, err := macro.LoadMacro(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Printf("\n  %sMacro: %s (%d steps)%s\n", constants.ColorCyan, m.Name, len(m.Steps), constants.ColorReset)
	for i, step := range m.Steps {
		fmt.Printf("    %d. %s\n", i+1, step.CommandLine)
	}
	fmt.Println()
}

func printMacroUsage() {
	fmt.Println("Usage: gitmap macro <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  record <name>    Record an interactive shell session as a macro")
	fmt.Println("  run <name>       Replay a recorded macro")
	fmt.Println("  list             List all saved macros")
	fmt.Println("  show <name>      Inspect steps of a macro")
	fmt.Println("  rm <name>        Delete a saved macro")
}

type TypeFuncf0956a91 struct{}

func InitFuncf0956a91() error             { return nil }
func (x *TypeFuncf0956a91) Process() bool { return true }

type TypeFunc58e1479d struct{}

func InitFunc58e1479d() error             { return nil }
func (x *TypeFunc58e1479d) Process() bool { return true }
