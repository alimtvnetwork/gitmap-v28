package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// runImport handles the "import" subcommand.
func runImport(args []string) error {
	checkHelp(constants.CmdImport, args)
	if len(args) == 0 {
		printBareImportGuidance()
		return nil
	}
	inFile, isConfirm := parseImportFlags(args)
	if !isConfirm {
		printImportConfirmRequired(inFile)
		return nil
	}

	data := readImportFile(inFile)
	if err := executeImport(data); err != nil {
		return err
	}

	printImportSummary(inFile, data)

	return nil
}

func printBareImportGuidance() {
	fmt.Printf("\n  %s📦 GitMap Import Guidance & Formats%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  SQLite Database Export Import:")
	fmt.Println("    gitmap import <file.json> --confirm")
	fmt.Println()
	fmt.Println("  Typed JSON Envelope & Configuration Commands:")
	fmt.Println("    • Inspect format & commands:  gitmap which-format [file.json]   (alias: gitmap wc)")
	fmt.Println("    • Bulk import any format:     gitmap import-all-json *.json")
	fmt.Println("    • SSH fleet nodes & keys:     gitmap sj import [file.json] -y   (or: gitmap ssh nodes import-json)")
	fmt.Println("    • Terminal / Web UI layout:   gitmap cmdui import-settings [file.json]")
	fmt.Println("    • Recorded macros:            gitmap macro import [file.json] -y")
	fmt.Println("    • State templates:            gitmap templates import [file.json]")
	fmt.Println()
	fmt.Println("  💡 Tip: To bypass prompts on envelope imports, pass the '-y' flag.")
	fmt.Println()
}

func printImportConfirmRequired(file string) {
	fmt.Printf("\n  %s⚠️  Safety Confirmation Required%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  Database import will overwrite or merge SQLite tables with '%s'.\n", file)
	fmt.Printf("  To execute this import, re-run with the --confirm flag:\n\n")
	fmt.Printf("    gitmap import %s --confirm\n\n", file)
}

// parseImportFlags parses the optional file arg and --confirm flag.
func parseImportFlags(args []string) (string, bool) {
	fs := flag.NewFlagSet(constants.CmdImport, flag.ExitOnError)
	var isConfirm bool
	fs.BoolVar(&isConfirm, constants.FlagConfirm, false, constants.FlagDescConfirm)
	_ = fs.Parse(args)

	file := constants.DefaultExportFile
	if fs.NArg() > 0 {
		file = fs.Arg(0)
	}

	return file, isConfirm || hasConfirmFlag(args)
}

// readImportFile reads and parses the export JSON file.
func readImportFile(path string) model.DatabaseExport {
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, apperror.WrapSimple(err, constants.MsgImportReadFailed).Error())
		cliexit.HandleError(err, 1)
	}

	var data model.DatabaseExport

	err = json.Unmarshal(raw, &data)
	if err != nil {
		fmt.Fprintln(os.Stderr, apperror.WrapSimple(err, constants.MsgImportParseFailed).Error())
		cliexit.HandleError(err, 1)
	}

	return data
}

// executeImport restores all data into the database.
func executeImport(data model.DatabaseExport) error {
	db, err := openDb()
	if err != nil {
		return apperror.WrapSimple(err, constants.MsgImportFailed)
	}

	defer db.Close()

	err = db.ImportAll(data)
	if err != nil {
		return apperror.WrapSimple(err, "import all")
	}

	return nil
}

// printImportSummary prints the import result summary.
func printImportSummary(path string, e model.DatabaseExport) {
	fmt.Printf(constants.MsgImportDone, path,
		len(e.Repos), len(e.Groups), len(e.Releases),
		len(e.History), len(e.Bookmarks))
}
