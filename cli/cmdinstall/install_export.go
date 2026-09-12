// Package cmd - install_export.go handles JSON and ZIP export/import for custom installers.
package cmdinstall

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
)

// ExportOptions holds options for exporting installers to JSON or ZIP.
type ExportOptions struct {
	Slug       string
	OutputPath string
	Format     string
	ExportAll  bool
}

// isInstallExportCommand checks if args invoke export or export-all.
func isInstallExportCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}

	subCmd := strings.ToLower(strings.TrimSpace(args[0]))

	return subCmd == "export" || subCmd == "export-all"
}

// isInstallImportCommand checks if args invoke import.
func isInstallImportCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}

	subCmd := strings.ToLower(strings.TrimSpace(args[0]))

	return subCmd == "import"
}

// runInstallExport handles installer export in JSON or ZIP format.
func runInstallExport(args []string) error {
	opts, errParse := parseExportOptions(args)
	if errParse != nil {
		return errParse
	}

	db, errDB := cmdinstaller.OpenAndMigrateInstallerDB()
	if errDB != nil {
		return errDB
	}

	defer db.Close()

	return executeInstallerExport(db, opts)
}

func parseExportOptions(args []string) (*ExportOptions, error) {
	fs := flag.NewFlagSet("install-export", flag.ContinueOnError)
	opts := &ExportOptions{}
	fs.StringVar(&opts.OutputPath, "output", "", "Output file path")
	fs.StringVar(&opts.OutputPath, "o", "", "Output shorthand")
	fs.StringVar(&opts.Format, "format", "", "Export format: json or zip")
	fs.BoolVar(&opts.ExportAll, "all", false, "Export all installers")
	fs.BoolVar(&opts.ExportAll, "a", false, "Export all shorthand")

	return parseExportPositional(fs, args, opts)
}

func parseExportPositional(fs *flag.FlagSet, args []string, opts *ExportOptions) (*ExportOptions, error) {
	flagArgs, positional := cmdinstaller.SeparateFlagAndPositionalArgs(args)
	if err := fs.Parse(flagArgs); err != nil {
		appErr := apperror.Wrap(err, "parseExportOptions", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"

		return nil, appErr
	}

	if len(positional) > 0 {
		opts.Slug = strings.TrimSpace(positional[0])
	}

	if opts.Slug == "" && !opts.ExportAll {
		appErr := apperror.NewValidationError("installer slug or --all required")
		appErr.Code = "E_INSTALLER_INVALID_INPUT"

		return nil, appErr
	}

	resolveExportDefaults(opts)

	return opts, nil
}

func resolveExportDefaults(opts *ExportOptions) {
	if opts.OutputPath != "" {
		return
	}

	ext := pickExportExtension(opts)
	opts.OutputPath = pickExportFilename(opts, ext)
}

func pickExportExtension(opts *ExportOptions) string {
	isJSON := strings.ToLower(opts.Format) == "json" || strings.HasSuffix(strings.ToLower(opts.OutputPath), ".json")
	if isJSON {
		return ".json"
	}

	return ".zip"
}

func pickExportFilename(opts *ExportOptions, ext string) string {
	if opts.ExportAll {
		return "gitmap-installers" + ext
	}

	return opts.Slug + ext
}

func executeInstallerExport(db *store.DB, opts *ExportOptions) error {
	scripts, errScripts := loadScriptsForExport(db, opts)
	if errScripts != nil {
		return errScripts
	}

	isJSON := strings.ToLower(opts.Format) == "json" || strings.HasSuffix(strings.ToLower(opts.OutputPath), ".json")
	if isJSON {
		return writeJSONExport(scripts, opts)
	}

	return writeZipExport(scripts, opts)
}

func loadScriptsForExport(db *store.DB, opts *ExportOptions) ([]model.InstallerScript, error) {
	if opts.ExportAll {
		return db.ListInstallers()
	}

	single, err := db.GetInstallerBySlug(opts.Slug)
	if err != nil {
		return nil, err
	}

	return []model.InstallerScript{*single}, nil
}

func writeJSONExport(scripts []model.InstallerScript, opts *ExportOptions) error {
	var payload []byte
	var err error
	if opts.ExportAll {
		payload, err = json.MarshalIndent(scripts, "", "  ")
	} else if len(scripts) > 0 {
		payload, err = json.MarshalIndent(scripts[0], "", "  ")
	}

	if err != nil {
		return apperror.WrapSimple(err, "marshal json export")
	}

	if errWrite := os.WriteFile(opts.OutputPath, payload, 0644); errWrite != nil {
		return apperror.WrapSimple(errWrite, "write json export file")
	}

	printExportSuccess(len(scripts), opts.OutputPath, "JSON")

	return nil
}

func writeZipExport(scripts []model.InstallerScript, opts *ExportOptions) error {
	flags := &cmdinstaller.ExportInstallerFlags{
		Slug:       opts.Slug,
		OutputPath: opts.OutputPath,
		ExportAll:  opts.ExportAll,
	}

	return executeExportArchive(scripts, flags)
}

func executeExportArchive(scripts []model.InstallerScript, flags *cmdinstaller.ExportInstallerFlags) error {
	outFile, errCreate := os.Create(flags.OutputPath)
	if errCreate != nil {
		return apperror.WrapSimple(errCreate, "create zip output")
	}

	defer outFile.Close()

	zw := zip.NewWriter(outFile)
	defer zw.Close()

	for _, s := range scripts {
		if errWrite := cmdinstaller.WriteZipEntry(zw, s); errWrite != nil {
			return apperror.WrapSimple(errWrite, "write zip entry")
		}
	}

	printExportSuccess(len(scripts), flags.OutputPath, "ZIP")

	return nil
}

func printExportSuccess(count int, path, format string) {
	msg := fmt.Sprintf("\n  %s✔ Exported %d installer(s) to %s (%s format) successfully!%s\n\n",
		constants.ColorGreen, count, path, format, constants.ColorReset)
	fmt.Print(msg)
}

// runInstallImport handles importing installers from JSON or ZIP format.
func runInstallImport(args []string) error {
	targetPath := parseImportPath(args)
	if targetPath == "" {
		appErr := apperror.NewValidationError("target file path is required for import")
		appErr.Code = "E_INSTALLER_INVALID_INPUT"

		return appErr
	}

	db, errDB := cmdinstaller.OpenAndMigrateInstallerDB()
	if errDB != nil {
		return errDB
	}

	defer db.Close()

	return dispatchImportByExtension(db, targetPath)
}

func parseImportPath(args []string) string {
	fs := flag.NewFlagSet("install-import", flag.ContinueOnError)
	fileFlag := fs.String("file", "", "File to import")
	fs.StringVar(fileFlag, "f", "", "File shorthand")
	flagArgs, positional := cmdinstaller.SeparateFlagAndPositionalArgs(args)
	fs.Parse(flagArgs)
	if *fileFlag != "" {
		return strings.TrimSpace(*fileFlag)
	}

	if len(positional) > 0 {
		return strings.TrimSpace(positional[0])
	}

	return ""
}

func dispatchImportByExtension(db *store.DB, targetPath string) error {
	if _, err := os.Stat(targetPath); err != nil {
		appErr := apperror.WrapSimple(err, "import file not found")
		appErr.Code = "E_INSTALLER_FILE_NOT_FOUND"

		return appErr
	}

	if strings.HasSuffix(strings.ToLower(targetPath), ".json") {
		return importFromJSONFilePath(db, targetPath)
	}

	return importFromZipArchiveFile(db, targetPath)
}

func importFromJSONFilePath(db *store.DB, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return apperror.WrapSimple(err, "read json import file")
	}

	return importRawJSONPayload(db, data, path)
}

func importRawJSONPayload(db *store.DB, data []byte, path string) error {
	trimmed := strings.TrimSpace(string(data))
	if strings.HasPrefix(trimmed, "[") {
		return importJSONArray(db, data, path)
	}

	return importJSONObject(db, data, path)
}

func importJSONArray(db *store.DB, data []byte, path string) error {
	var list []model.InstallerScript
	if err := json.Unmarshal(data, &list); err != nil {
		return apperror.WrapSimple(err, "unmarshal json array")
	}

	count := upsertScriptList(db, list)
	printImportSuccess(count, path)

	return nil
}

func importJSONObject(db *store.DB, data []byte, path string) error {
	var script model.InstallerScript
	if err := json.Unmarshal(data, &script); err != nil {
		return apperror.WrapSimple(err, "unmarshal json object")
	}

	count := upsertScriptList(db, []model.InstallerScript{script})
	printImportSuccess(count, path)

	return nil
}

func importFromZipArchiveFile(db *store.DB, path string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return apperror.WrapSimple(err, "open zip archive")
	}

	defer r.Close()

	var list []model.InstallerScript
	for _, f := range r.File {
		script := extractZipEntryScript(f)
		if script != nil {
			list = append(list, *script)
		}
	}

	count := upsertScriptList(db, list)
	printImportSuccess(count, path)

	return nil
}

func extractZipEntryScript(f *zip.File) *model.InstallerScript {
	if !strings.HasSuffix(strings.ToLower(f.Name), ".json") {
		return nil
	}

	rc, errOpen := f.Open()
	if errOpen != nil {
		return nil
	}

	defer rc.Close()

	var s model.InstallerScript
	if errDec := json.NewDecoder(rc).Decode(&s); errDec != nil {
		return nil
	}

	return &s
}

func upsertScriptList(db *store.DB, scripts []model.InstallerScript) int {
	importedCount := 0
	for _, s := range scripts {
		if s.Slug == "" {
			s.Slug = cmdinstaller.Slugify(s.Name)
		}

		if saveOrUpdateScript(db, &s) {
			importedCount++
		}
	}

	return importedCount
}

func saveOrUpdateScript(db *store.DB, s *model.InstallerScript) bool {
	existing, _ := db.GetInstallerBySlug(s.Slug)
	if existing != nil {
		s.ID = existing.ID

		return db.UpdateInstaller(s) == nil
	}

	return db.CreateInstaller(s) == nil
}

func printImportSuccess(count int, path string) {
	msg := fmt.Sprintf("\n  %s✔ Successfully imported %d installer(s) from %s!%s\n\n",
		constants.ColorGreen, count, filepath.Base(path), constants.ColorReset)
	fmt.Print(msg)
}
