package cmdmacro

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

type macroExportOpts struct {
	TargetName string
	FilePath   string
	Format     string
	ExceptList []string
	IsAll      bool
	IsSingle   bool
}

func runMacroExport(args []string) error {
	opts := parseMacroExportOpts(args)
	if opts.TargetName == "" && !opts.IsAll && !opts.IsSingle {
		opts.IsAll = true
	}

	macros, err := collectMacrosForExport(opts)
	if err != nil {
		return err
	}

	if len(macros) == 0 {
		fmt.Println("No matching macros found to export.")

		return nil
	}

	return executeMacroExport(macros, opts)
}

func parseMacroExportOpts(args []string) macroExportOpts {
	opts := macroExportOpts{Format: constants.OutputJSON}
	for i := 0; i < len(args); i++ {
		a := args[i]
		processExportFlag(a, args, &i, &opts)
	}

	inferMacroExportFormat(&opts)
	reconcileExportTargets(&opts)

	return opts
}

func reconcileExportTargets(opts *macroExportOpts) {
	if isAllExportTarget(opts.TargetName) {
		opts.IsAll = true
		opts.TargetName = ""
	}

	if opts.TargetName == "single" {
		opts.IsSingle = true
		opts.TargetName = ""
	}

	if opts.TargetName != "" && opts.FilePath == "" && isLikelyFilePath(opts.TargetName) {
		opts.FilePath = opts.TargetName
		opts.TargetName = ""
		inferMacroExportFormat(opts)
	}
}

func isLikelyFilePath(token string) bool {
	return hasKnownExportExtension(token) || strings.ContainsAny(token, "/\\")
}

func hasKnownExportExtension(token string) bool {
	ext := strings.ToLower(filepath.Ext(token))
	switch ext {
	case ".json", ".yaml", ".yml", ".db", ".sqlite", ".sqlite3", ".sqlitedb", ".zip":
		return true
	default:
		return false
	}
}

func isAllExportTarget(name string) bool {
	return name == "*" || name == "all" || name == "export-all"
}

func processExportFlag(arg string, args []string, index *int, opts *macroExportOpts) {
	if processExportFormatFlag(arg, args, index, opts) {
		return
	}

	if processExportScopeFlag(arg, opts) {
		return
	}

	processExportTargetFlag(arg, args, index, opts)
}

func processExportFormatFlag(arg string, args []string, index *int, opts *macroExportOpts) bool {
	switch {
	case matchFlagWithVal(arg, "--format"):
		opts.Format = normalizeMacroFormat(extractFlagValue(index, args))
		return true
	case arg == "--json":
		opts.Format = constants.OutputJSON
		return true
	case arg == "--yaml" || arg == "--yml" || arg == "-y":
		opts.Format = constants.OutputYAML
		return true
	case arg == "--sqlite" || arg == "--db" || arg == "--sqlitedb":
		opts.Format = "sqlite"
		return true
	case arg == "--zip":
		opts.Format = "zip"
		return true
	default:
		return false
	}
}

func processExportScopeFlag(arg string, opts *macroExportOpts) bool {
	switch {
	case arg == "--all" || arg == "all" || arg == "*":
		opts.IsAll = true
		opts.IsSingle = false
		return true
	case arg == "--single" || arg == "single":
		opts.IsSingle = true
		opts.IsAll = false
		return true
	default:
		return false
	}
}

func processExportTargetFlag(arg string, args []string, index *int, opts *macroExportOpts) {
	switch {
	case matchFlagWithVal(arg, "-f", "--file", "--filepath", "-o", "--output", "--out"):
		opts.FilePath = extractFlagValue(index, args)
	case matchFlagWithVal(arg, "-except", "--except", "--exclude"):
		opts.ExceptList = parseExceptTokens(extractFlagValue(index, args))
	case matchFlagWithVal(arg, "--name", "--target", "--macro"):
		opts.TargetName = extractFlagValue(index, args)
	case !strings.HasPrefix(arg, "-") && opts.TargetName == "":
		opts.TargetName = arg
	case !strings.HasPrefix(arg, "-") && opts.FilePath == "":
		opts.FilePath = arg
	}
}

func normalizeMacroFormat(raw string) string {
	lower := strings.ToLower(strings.TrimSpace(raw))
	if lower == "sqlitedb" || lower == "db" || lower == "sqlite3" {
		return "sqlite"
	}

	if lower == "yml" {
		return constants.OutputYAML
	}

	return lower
}

func inferMacroExportFormat(opts *macroExportOpts) {
	if opts.FilePath == "" {
		return
	}

	ext := strings.ToLower(filepath.Ext(opts.FilePath))
	if ext == ".yaml" || ext == ".yml" {
		opts.Format = constants.OutputYAML
	}

	if ext == ".db" || ext == ".sqlite" || ext == ".sqlite3" || ext == ".sqlitedb" {
		opts.Format = "sqlite"
	}

	if ext == ".zip" {
		opts.Format = "zip"
	}
}

func collectSingleMacro(name string) ([]macro.Macro, error) {
	m, err := macro.LoadMacro(name)
	if err != nil {
		return nil, apperror.WrapSimple(err, "load macro: "+name)
	}

	return []macro.Macro{*m}, nil
}

func collectMacrosForExport(opts macroExportOpts) ([]macro.Macro, error) {
	if opts.IsSingle && opts.TargetName == "" {
		return nil, apperror.NewValidationError("macro name required for single export (usage: gitmap macro export single <name>)")
	}

	if !opts.IsAll && opts.TargetName != "" {
		return collectSingleMacro(opts.TargetName)
	}

	listRes := macro.ListMacros()
	if listRes.IsFailure() {
		return nil, listRes.AppError()
	}

	return macro.FilterMacrosForExport(listRes.Data, macro.ExportOptions{
		IsAll:      true,
		ExceptList: opts.ExceptList,
	}), nil
}

func executeMacroExport(macros []macro.Macro, opts macroExportOpts) error {
	switch opts.Format {
	case "sqlite", "db", "sqlitedb", "sqlite3":
		return exportMacrosSQLiteOutput(macros, opts)
	case "zip":
		return exportMacrosZIPOutput(macros, opts)
	default:
		return exportMacrosTextOutput(macros, opts)
	}
}

func resolveSQLiteOutPath(opts macroExportOpts) string {
	if opts.FilePath != "" {
		return opts.FilePath
	}

	if opts.TargetName != "" && !opts.IsAll {
		return opts.TargetName + ".db"
	}

	return "macros_export.db"
}

func exportMacrosSQLiteOutput(macros []macro.Macro, opts macroExportOpts) error {
	outPath := resolveSQLiteOutPath(opts)

	if err := macro.ExportMacrosToSQLite(macros, outPath); err != nil {
		return err
	}

	printMacroExportSuccessBanner(outPath, len(macros), "sqlite")

	return nil
}

func exportMacrosZIPOutput(macros []macro.Macro, opts macroExportOpts) error {
	outPath := opts.FilePath
	if outPath == "" {
		outPath = "macros_export.zip"
	}

	if err := macro.ExportToZIP(macros, outPath); err != nil {
		return err
	}

	printMacroExportSuccessBanner(outPath, len(macros), "zip")

	return nil
}

func exportMacrosTextOutput(macros []macro.Macro, opts macroExportOpts) error {
	payload, err := macro.SerializeSingleOrAll(macros, !opts.IsAll, opts.Format)
	if err != nil {
		return err
	}

	if opts.FilePath == "" {
		fmt.Println(string(payload))

		return nil
	}

	if err := macro.WriteExportPayload(opts.FilePath, payload); err != nil {
		return err
	}

	printMacroExportSuccessBanner(opts.FilePath, len(macros), opts.Format)

	return nil
}

func printMacroExportSuccessBanner(filePath string, count int, format string) {
	fmt.Printf("\n  %s✔ Exported %d macro(s)%s to %s%s%s (%s%s%s)\n\n",
		constants.ColorGreen, count, constants.ColorReset,
		constants.ColorCyan, filePath, constants.ColorReset,
		constants.ColorYellow, format, constants.ColorReset)
}
