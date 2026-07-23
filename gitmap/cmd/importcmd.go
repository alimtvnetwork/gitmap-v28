package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runImport handles the "import" subcommand.
func runImport(args []string) {
	checkHelp("import", args)
	inFile, confirm := parseImportFlags(args)
	if !confirm {
		fmt.Fprint(os.Stderr, constants.ErrImportNoConfirm)
		os.Exit(1)
	}

	data := readImportFile(inFile)
	executeImport(data)
	printImportSummary(inFile, data)
}

// parseImportFlags parses the optional file arg and --confirm flag.
func parseImportFlags(args []string) (string, bool) {
	fs := flag.NewFlagSet(constants.CmdImport, flag.ExitOnError)
	confirmFlag := fs.Bool("confirm", false, constants.FlagDescConfirm)
	fs.Parse(args)

	file := constants.DefaultExportFile
	if fs.NArg() > 0 {
		file = fs.Arg(0)
	}

	return file, *confirmFlag
}

// readImportFile reads and parses the export JSON file.
func readImportFile(path string) model.DatabaseExport {
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgImportReadFailed, err)
		os.Exit(1)
	}

	var data model.DatabaseExport

	err = json.Unmarshal(raw, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgImportParseFailed, err)
		os.Exit(1)
	}

	return data
}

// executeImport restores all data into the database.
func executeImport(data model.DatabaseExport) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgImportFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.ImportAll(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgImportFailed, err)
		os.Exit(1)
	}
}

// printImportSummary prints the import result summary.
func printImportSummary(path string, e model.DatabaseExport) {
	fmt.Printf(constants.MsgImportDone, path,
		len(e.Repos), len(e.Groups), len(e.Releases),
		len(e.History), len(e.Bookmarks))
}
