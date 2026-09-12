package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runImport handles the "import" subcommand.
func runImport(args []string) error {
	checkHelp(constants.CmdImport, args)
	inFile, isConfirm := parseImportFlags(args)
	if !isConfirm {
		return apperror.NewSimple(constants.ErrImportNoConfirm, "E9000")
	}

	data := readImportFile(inFile)
	if err := executeImport(data); err != nil {
		return err
	}
	printImportSummary(inFile, data)

	return nil
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
		cliexit.HandleError(nil, 1)
	}

	var data model.DatabaseExport

	err = json.Unmarshal(raw, &data)
	if err != nil {
		fmt.Fprintln(os.Stderr, apperror.WrapSimple(err, constants.MsgImportParseFailed).Error())
		cliexit.HandleError(nil, 1)
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
