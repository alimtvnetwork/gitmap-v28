package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

func parseExecOptions(flagArgs []string) macro.ExecOptions {
	opts := macro.ExecOptions{}
	for i := 0; i < len(flagArgs); i++ {
		arg := flagArgs[i]
		if arg == "--dry-run" {
			opts.DryRun = true
		}
		if arg == "--verbose" || arg == "-v" {
			opts.Verbose = true
		}
		if arg == "--json" {
			opts.JSON = true
		}
		if arg == "--yaml" || arg == "--yml" || arg == "-y" {
			opts.YAML = true
		}
		if isFileFlagWithArg(arg) && i+1 < len(flagArgs) {
			opts.FilePath = flagArgs[i+1]
			i++
		}
		checkInlineFileArg(arg, &opts)
	}
	return opts
}

func checkInlineFileArg(arg string, opts *macro.ExecOptions) {
	prefixes := []string{"--file=", "--filepath=", "--out=", "--output="}
	for _, p := range prefixes {
		if strings.HasPrefix(strings.ToLower(arg), p) {
			opts.FilePath = arg[len(p):]
			return
		}
	}
}

func extractMacroNameAndFlags(args []string) (string, []string) {
	name := ""
	flags := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		isFlag := strings.HasPrefix(args[i], "-")
		if !isFlag && name == "" {
			name = args[i]
			continue
		}
		if !isFlag {
			flags = append(flags, args[i])
			continue
		}
		flags = append(flags, args[i])
		if isFileFlagWithArg(args[i]) && i+1 < len(args) {
			flags = append(flags, args[i+1])
			i++
		}
	}
	return name, flags
}

func isFileFlagWithArg(arg string) bool {
	lower := strings.ToLower(arg)
	return lower == "--file" || lower == "--filepath" || lower == "--out" || lower == "--output" || lower == "-o" || lower == "-f"
}

func executeMacroByName(macroName string, opts macro.ExecOptions) {
	loadedMacro, loadErr := macro.LoadMacro(macroName)
	if loadErr != nil {
		fmt.Fprintf(os.Stderr, "%s✖ Error: %v%s\n", constants.ColorRed, loadErr, constants.ColorReset)
		cliexit.HandleError(nil, 1)
	}
	if isStandardDisplay(opts) && len(loadedMacro.Steps) > 0 {
		printMacroStepsTree(loadedMacro)
	}
	if execErr := macro.Execute(context.Background(), loadedMacro, opts); execErr != nil {
		cliexit.HandleError(nil, 1)
	}
}

func isStandardDisplay(opts macro.ExecOptions) bool {
	return !opts.JSON && !opts.YAML && len(opts.FilePath) == 0
}

func runExecuteCmd(args []string) error {
	macroName, flagArgs := extractMacroNameAndFlags(args)
	if macroName == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro run <name> [--json] [--yaml] [--file <path>] [--dry-run] [--verbose]\n")
		cliexit.HandleError(nil, 1)
	}
	opts := parseExecOptions(flagArgs)
	executeMacroByName(macroName, opts)
	return nil
}

func runMacroCmd(args []string) error {
	if len(args) == 0 {
		printMacroUsage()
		return nil
	}
	return routeMacroSubcommand(args[0], args[1:])
}

func routeMacroSubcommand(sub string, rest []string) error {
	switch sub {
	case "run", "exec":
		return runExecuteCmd(rest)
	case "record", "rec":
		return handleMacroRecord(rest)
	case "list", "ls":
		return handleMacroList(rest)
	case "show":
		return handleMacroShow(rest)
	case "rm", "delete":
		return handleMacroDelete(rest)
	default:
		printMacroUsage()
	}
	return nil
}

func handleMacroRecord(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro record <name>\n")
		cliexit.HandleError(nil, 1)
	}
	if err := macro.RecordInteractive(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	return nil
}

func handleMacroList(args []string) error {
	opts := parseExecOptions(args)
	macros, err := macro.ListMacros()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing macros: %v\n", err)
		return nil
	}
	if opts.JSON || opts.YAML || len(opts.FilePath) > 0 {
		return outputStructuredData(macros, opts)
	}
	renderMacroListTable(macros)
	return nil
}

func handleMacroShow(args []string) error {
	name, flagArgs := extractMacroNameAndFlags(args)
	if name == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro show <name> [--json] [--yaml] [--file <path>]\n")
		cliexit.HandleError(nil, 1)
	}
	opts := parseExecOptions(flagArgs)
	m, err := macro.LoadMacro(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return nil
	}
	if opts.JSON || opts.YAML || len(opts.FilePath) > 0 {
		return outputStructuredData(m, opts)
	}
	renderMacroShow(m)
	return nil
}

func handleMacroDelete(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro rm <name>\n")
		cliexit.HandleError(nil, 1)
	}
	if err := macro.DeleteMacro(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	fmt.Printf("✔ Removed macro %q\n", args[0])
	return nil
}

func outputStructuredData(data interface{}, opts macro.ExecOptions) error {
	isYAML := opts.YAML || strings.HasSuffix(strings.ToLower(opts.FilePath), ".yaml") || strings.HasSuffix(strings.ToLower(opts.FilePath), ".yml")
	formatted := formatStructuredBytes(data, isYAML)
	if len(opts.FilePath) == 0 {
		fmt.Println(formatted)
		return nil
	}
	return saveAndPrintStructuredOutput(opts.FilePath, formatted)
}

func formatStructuredBytes(data interface{}, isYAML bool) string {
	if isYAML {
		bytes, _ := yaml.Marshal(data)
		return string(bytes)
	}
	bytes, _ := json.MarshalIndent(data, "", "  ")
	return string(bytes)
}

func saveAndPrintStructuredOutput(filePath, formatted string) error {
	savedPath, saveErr := macro.SaveReportToFile(filePath, formatted)
	if saveErr != nil {
		return fmt.Errorf("failed saving to %s: %w", filePath, saveErr)
	}
	fmt.Println(formatted)
	fmt.Printf("\n  %s✔ Output saved to:%s %s%s%s\n\n",
		constants.ColorGreen, constants.ColorReset,
		constants.ColorCyan, savedPath, constants.ColorReset)
	return nil
}

func renderMacroListTable(macros []macro.Macro) {
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

func renderMacroShow(m *macro.Macro) {
	fmt.Printf("\n  %sMacro: %s (%d steps)%s\n", constants.ColorCyan, m.Name, len(m.Steps), constants.ColorReset)
	for i, step := range m.Steps {
		fmt.Printf("    %d. %s\n", i+1, step.CommandLine)
	}
	fmt.Println()
}

func printMacroUsage() {
	fmt.Println("Usage: gitmap macro <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  record <name>                  Record an interactive shell session as a macro")
	fmt.Println("  run <name> [--json] [--yaml]   Replay a macro (optional JSON/YAML & file export)")
	fmt.Println("  list [--json] [--yaml]         List all saved macros")
	fmt.Println("  show <name> [--json] [--yaml]  Inspect steps of a macro")
	fmt.Println("  rm <name>                      Delete a saved macro")
}
