package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmdchromeprofile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runImportExport provides a unified entry point for export and import commands.
func runImportExport(args []string) error {
	hasArgs := len(args) > 0
	if hasArgs && (args[0] == "import" || args[0] == "im") {
		runImport(args[1:])

		return nil
	}

	if hasArgs && (args[0] == "export" || args[0] == "ex") {
		runExport(args[1:])

		return nil
	}

	runExport(args)

	return nil
}

// runExport handles the "export" subcommand.
func runExport(args []string) error {
	checkHelp("export", args)
	outFile := resolveExportFile(args)
	export := loadExportData()

	if err := writeExportFile(outFile, export); err != nil {
		return err
	}

	printExportSummary(outFile, export)

	return nil
}

// resolveExportFile determines the output file from args or default.
func resolveExportFile(args []string) string {
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		return args[0]
	}

	return constants.DefaultExportFile
}

// loadExportData fetches the full database export.
func loadExportData() model.DatabaseExport {
	db, err := openDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, apperror.WrapSimple(err, constants.MsgExportFailed).Error())
		cliexit.HandleError(nil, 1)
	}

	defer db.Close()

	export, err := db.ExportAll()
	if err != nil && isLegacyDataError(err) {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
		cliexit.HandleError(apperror.NewSimple("legacy project data", "E9000"), 1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, apperror.WrapSimple(err, constants.MsgExportFailed).Error())
		cliexit.HandleError(nil, 1)
	}

	return export
}

// writeExportFile marshals the export data to a JSON file using the
// stablejson-backed encoder so the top-level key order is contractual.
func writeExportFile(path string, export model.DatabaseExport) *apperror.AppError {
	var buf bytes.Buffer
	if err := encodeDatabaseExportJSON(&buf, export); err != nil {
		return apperror.WrapSimple(err, constants.MsgExportFailed)
	}

	if err := os.WriteFile(path, buf.Bytes(), constants.DirPermission); err != nil {
		return apperror.WrapSimple(err, constants.MsgExportFailed)
	}

	return nil
}

// printExportSummary prints the export result summary.
func printExportSummary(path string, e model.DatabaseExport) {
	fmt.Printf(constants.MsgExportDone, path,
		len(e.Repos), len(e.Groups), len(e.Releases),
		len(e.History), len(e.Bookmarks))
}

// runExportAll handles the "export-all" command.
func runExportAll(args []string) error {
	checkHelp("export-all", args)

	return cmdchromeprofile.RunExportAll(args)
}

// runImportAll handles the "import-all" command.
func runImportAll(args []string) error {
	checkHelp("import-all", args)

	return cmdchromeprofile.RunImportAll(args)
}

// runExportOnly handles the "export-only" command.
func runExportOnly(args []string) error {
	checkHelp("export", args)

	return runExport(args)
}
