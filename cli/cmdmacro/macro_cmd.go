package cmdmacro

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func parseExecOptions(flagArgs []string) macro.ExecOptions {
	return ParseExecOptions(flagArgs)
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

func executeMacroByName(macroName string, opts macro.ExecOptions) error {
	loadedMacro, loadErr := macro.LoadMacro(macroName)
	if loadErr != nil {
		fmt.Fprintf(os.Stderr, "%s✖ Error: %v%s\n", constants.ColorRed, loadErr, constants.ColorReset)

		return apperror.WrapSimple(loadErr, "macro.LoadMacro")
	}

	if isStandardDisplay(opts) && len(loadedMacro.Steps) > 0 {
		printMacroStepsTree(loadedMacro)
	}

	if execErr := macro.Execute(context.Background(), loadedMacro, opts); execErr != nil {
		return apperror.WrapSimple(execErr, "macro.Execute")
	}

	return nil
}

func isStandardDisplay(opts macro.ExecOptions) bool {
	return !opts.JSON && !opts.YAML && len(opts.FilePath) == 0 && !opts.IsSummaryOnly && !opts.IsTerminalSuppressed
}

func runExecuteCmd(args []string) error {
	macroName, flagArgs := extractMacroNameAndFlags(args)
	if macroName == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro run <name> [--json] [--yaml] [--file <path>] [--dry-run] [--verbose]\n")

		return apperror.NewValidationError("missing required macro name")
	}

	opts := parseExecOptions(flagArgs)

	return executeMacroByName(macroName, opts)
}

func runMacroCmd(args []string) error {
	if len(args) == 0 {
		printMacroUsage()

		return nil
	}

	return result.AsError(routeMacroSubcommand(args[0], args[1:]))
}

func routeMacroSubcommand(sub string, rest []string) result.ErrorWrapper {
	if isExecSubcommand(sub) {
		return routeExecSubcommand(sub, rest)
	}

	if isModifierSubcommand(sub) && len(rest) > 0 && isExportImportSubcommand(rest[0]) {
		return routeExportImportSubcommand(rest[0], append([]string{"--" + sub}, rest[1:]...))
	}

	if isExportImportSubcommand(sub) {
		return routeExportImportSubcommand(sub, rest)
	}

	return routeManagementSubcommand(sub, rest)
}

func isModifierSubcommand(sub string) bool {
	return sub == "all" || sub == "single"
}

func isExportImportSubcommand(sub string) bool {
	switch sub {
	case "export", "exp", "dump", "export-all", "export-single":
		return true
	case "import", "imp", "load", "restore", "import-all", "import-single":
		return true
	default:
		return false
	}
}

func routeExportImportSubcommand(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "export", "exp", "dump":
		return result.MatchWrapper(runMacroExport(rest))
	case "export-all":
		return result.MatchWrapper(runMacroExport(append([]string{"--all"}, rest...)))
	case "export-single":
		return result.MatchWrapper(runMacroExport(append([]string{"--single"}, rest...)))
	case "import-all":
		return result.MatchWrapper(runMacroImport(append([]string{"--all"}, rest...)))
	case "import-single":
		return result.MatchWrapper(runMacroImport(append([]string{"--single"}, rest...)))
	default:
		return result.MatchWrapper(runMacroImport(rest))
	}
}

func isExecSubcommand(sub string) bool {
	return sub == "run" || sub == "exec" || isRetrySubcommand(sub) || isRunUntilSubcommand(sub)
}

func isRunUntilSubcommand(sub string) bool {
	switch sub {
	case "run-until", "run-until-end", "keep-going", "continue-on-error":
		return true
	default:
		return false
	}
}

func isRetrySubcommand(sub string) bool {
	switch sub {
	case "run-until-succeed", "until-success", "retry", "loop-until-success", "loop", "retry-until-success":
		return true
	default:
		return false
	}
}

func routeExecSubcommand(sub string, rest []string) result.ErrorWrapper {
	if isRetrySubcommand(sub) {
		return result.MatchWrapper(runMacroUntilSuccess(rest))
	}

	if isRunUntilSubcommand(sub) {
		return result.MatchWrapper(runMacroRunUntil(rest))
	}

	return result.MatchWrapper(runExecuteCmd(rest))
}

func runMacroRunUntil(args []string) error {
	macroName, flagArgs := extractMacroNameAndFlags(args)
	if macroName == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro run-until <name> [options]\n")

		return apperror.NewValidationError("missing required macro name")
	}

	opts := ParseExecOptions(flagArgs)
	opts.IsRunUntil = true

	return executeMacroByName(macroName, opts)
}

// MacroSyncRunner is wired by cmd package to run SSH macro sync.
var MacroSyncRunner func([]string) error

func routeManagementSubcommand(sub string, rest []string) result.ErrorWrapper {
	if sub == "startup" {
		return result.MatchWrapper(runMacroStartup(rest))
	}

	if sub == "schedule" || sub == "cron" || sub == "crontab" {
		return result.MatchWrapper(runMacroSchedule(rest))
	}

	if sub == "sync" || sub == "s" {
		return routeSyncSubcommand(rest)
	}

	if isModifySubcommand(sub) {
		return routeModifySubcommand(sub, rest)
	}

	return routeInspectSubcommand(sub, rest)
}

func routeSyncSubcommand(rest []string) result.ErrorWrapper {
	if MacroSyncRunner != nil {
		return result.MatchWrapper(MacroSyncRunner(rest))
	}

	return result.MatchWrapper(apperror.NewSimple("macro sync runner not initialized", "E_MACRO_SYNC"))
}

func isModifySubcommand(sub string) bool {
	switch sub {
	case "add", "create", "new", "edit", "modify", "record", "rec", "rm", "delete":
		return true
	default:
		return false
	}
}

func routeModifySubcommand(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "add", "create", "new":
		return result.MatchWrapper(handleMacroAdd(rest))
	case "edit", "modify":
		return result.MatchWrapper(handleMacroEdit(rest))
	case "record", "rec":
		return result.MatchWrapper(handleMacroRecord(rest))
	default:
		return result.MatchWrapper(handleMacroDelete(rest))
	}
}

func routeInspectSubcommand(sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "list", "ls":
		return result.MatchWrapper(handleMacroList(rest))
	case "show":
		return result.MatchWrapper(handleMacroShow(rest))
	default:
		printMacroUsage()
	}

	return result.SuccessWrapper()
}

func handleMacroRecord(args []string) error {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro record <name>\n")

		return apperror.NewValidationError("missing required macro name")
	}

	if err := macro.RecordInteractive(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return apperror.WrapSimple(err, "macro.RecordInteractive")
	}

	return nil
}

func handleMacroList(args []string) error {
	opts := parseExecOptions(args)
	macroRes := macro.ListMacros()
	if macroRes.IsFailure() {
		fmt.Fprintf(os.Stderr, "Error listing macros: %v\n", macroRes.AppError())

		return macroRes.AppError().WithContext("caller", "handleMacroList")
	}

	if opts.JSON || opts.YAML || len(opts.FilePath) > 0 {
		return outputStructuredData(macroRes.Data, opts)
	}

	renderMacroListTable(macroRes.Data)

	return nil
}

func handleMacroShow(args []string) error {
	name, flagArgs := extractMacroNameAndFlags(args)
	if name == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro show <name> [--json] [--yaml] [--file <path>]\n")

		return apperror.NewValidationError("missing required macro name")
	}

	opts := parseExecOptions(flagArgs)
	m, err := macro.LoadMacro(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return apperror.WrapSimple(err, "macro.LoadMacro")
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

		return apperror.NewValidationError("missing required macro name")
	}

	if err := macro.DeleteMacro(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return apperror.WrapSimple(err, "macro.DeleteMacro")
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
	fmt.Println("  add <name> <steps...>               Create/add a new macro directly from arguments")
	fmt.Println("  edit <name> [--no-exec]             Interactively edit steps of an existing macro")
	fmt.Println("  run <name> [--json] [--yaml]        Replay a macro (optional JSON/YAML & file export)")
	fmt.Println("  run-until-succeed <name|cmd>        Execute macro repeatedly until success")
	fmt.Println("  record <name>                       Record an interactive shell session as a macro")
	fmt.Println("  list [--json] [--yaml]              List all saved macros")
	fmt.Println("  show <name> [--json] [--yaml]       Inspect steps of a macro")
	fmt.Println("  rm <name>                           Delete a saved macro")
	fmt.Println("  startup <subcommand> [name]         Manage macro execution on OS login/reboot")
	fmt.Println("  schedule <subcommand> [name]        Manage recurring scheduled execution of macros")
	fmt.Println("  export [name|all] [options]         Export macro(s) to JSON, YAML, SQLite DB, or ZIP")
	fmt.Println("  export-all [options]                Export all macros (--json, --yaml, --sqlitedb, --zip)")
	fmt.Println("  export-single <name> [options]      Export a single macro to JSON, YAML, or SQLite DB")
	fmt.Println("  sync [flags]                        Synchronize local macros to remote SSH machines")
	fmt.Println("  import <file> [name] [options]      Import macro(s) safely with format auto-inference")
	fmt.Println("  import-all <file> [options]         Import all macros from backup archive or database")
	fmt.Println("  import-single <file> [options]      Import single macro from file with optional --as rename")
}
