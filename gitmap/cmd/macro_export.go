package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

type macroExportOpts struct {
	TargetName string
	FilePath   string
	Format     string
	ExceptList []string
	IsAll      bool
}

func runMacroExport(args []string) error {
	opts := parseMacroExportOpts(args)
	if opts.TargetName == "" && !opts.IsAll {
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
	if isAllExportTarget(opts.TargetName) {
		opts.IsAll = true
		opts.TargetName = ""
	}

	return opts
}

func isAllExportTarget(name string) bool {
	return name == "*" || name == "all" || name == "export-all"
}

func processExportFlag(arg string, args []string, index *int, opts *macroExportOpts) {
	switch {
	case matchFlagWithVal(arg, "-f", "--file", "--filepath", "-o", "--output", "--out"):
		opts.FilePath = extractFlagValue(index, args)
	case matchFlagWithVal(arg, "-except", "--except", "--exclude"):
		opts.ExceptList = parseExceptTokens(extractFlagValue(index, args))
	case matchFlagWithVal(arg, "--format"):
		opts.Format = normalizeMacroFormat(extractFlagValue(index, args))
	case arg == "--all":
		opts.IsAll = true
	case arg == "--json":
		opts.Format = constants.OutputJSON
	case arg == "--yaml" || arg == "--yml" || arg == "-y":
		opts.Format = constants.OutputYAML
	case arg == "--sqlite" || arg == "--db" || arg == "--sqlitedb":
		opts.Format = "sqlite"
	case arg == "--zip":
		opts.Format = "zip"
	case !strings.HasPrefix(arg, "-") && opts.TargetName == "":
		opts.TargetName = arg
	}
}

func normalizeMacroFormat(raw string) string {
	lower := strings.ToLower(strings.TrimSpace(raw))
	if lower == "sqlitedb" || lower == "db" {
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

	if ext == ".db" || ext == ".sqlite" || ext == ".sqlite3" {
		opts.Format = "sqlite"
	}

	if ext == ".zip" {
		opts.Format = "zip"
	}
}

func collectMacrosForExport(opts macroExportOpts) ([]macro.Macro, error) {
	if !opts.IsAll && opts.TargetName != "" {
		m, err := macro.LoadMacro(opts.TargetName)
		if err != nil {
			return nil, apperror.WrapSimple(err, "load macro: "+opts.TargetName)
		}

		return []macro.Macro{*m}, nil
	}

	list, err := macro.ListMacros()
	if err != nil {
		return nil, apperror.WrapSimple(err, "list macros for export")
	}

	filtered := macro.FilterMacrosForExport(list, macro.ExportOptions{
		IsAll:      true,
		ExceptList: opts.ExceptList,
	})

	return filtered, nil
}

func executeMacroExport(macros []macro.Macro, opts macroExportOpts) error {
	switch opts.Format {
	case "sqlite", "db":
		return exportMacrosSQLiteOutput(macros, opts)
	case "zip":
		return exportMacrosZIPOutput(macros, opts)
	default:
		return exportMacrosTextOutput(macros, opts)
	}
}

func exportMacrosSQLiteOutput(macros []macro.Macro, opts macroExportOpts) error {
	outPath := opts.FilePath
	if outPath == "" {
		outPath = "macros_export.db"
	}

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
