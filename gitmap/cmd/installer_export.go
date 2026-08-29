// Package cmd — installer_export.go defines installer export CLI subcommands.
package cmd

import (
	"archive/zip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// ExportInstallerFlags encapsulates options for installer exports.
type ExportInstallerFlags struct {
	Slug       string
	OutputPath string
	ExportAll  bool
}

// installerExportCmd represents the 'gitmap installer export' subcommand.
var installerExportCmd = &cobra.Command{
	Use:                "export <slug>",
	Short:              "Export an installer script into a zip bundle",
	Long:               "Packages installer definitions and versions as JSON within a .zip archive.",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerExport(cmd, args, false)
	},
}

// installerExportAllCmd represents the 'gitmap installer export-all' subcommand.
var installerExportAllCmd = &cobra.Command{
	Use:                "export-all",
	Short:              "Export all installer scripts into a zip bundle",
	Long:               "Packages all registered installer scripts as JSON files inside a .zip archive.",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallerExport(cmd, args, true)
	},
}

// InstallerExportCmd is an exported alias for installerExportCmd.
var InstallerExportCmd = installerExportCmd

// InstallerExportAllCmd is an exported alias for installerExportAllCmd.
var InstallerExportAllCmd = installerExportAllCmd

func init() {
	if installerCmd != nil {
		installerCmd.AddCommand(installerExportCmd)
		installerCmd.AddCommand(installerExportAllCmd)
	}
}

// parseExportFlags parses CLI options for export commands.
func parseExportFlags(args []string, isAll bool) (*ExportInstallerFlags, error) {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	out := fs.String("output", "", "Output zip file path")
	fs.StringVar(out, "o", "", "Output shorthand")

	var flagArgs []string
	var positional []string
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			flagArgs = append(flagArgs, args[i])
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flagArgs = append(flagArgs, args[i+1])
				i++
			}
		} else {
			positional = append(positional, args[i])
		}
	}

	if err := fs.Parse(flagArgs); err != nil {
		appErr := apperror.Wrap(err, "parseExportFlags", map[string]any{"args": args})
		appErr.Code = "E_INSTALLER_INVALID_FLAGS"
		return nil, appErr
	}

	targetPath := *out
	if targetPath == "" {
		targetPath = "gitmap-export.zip"
	}

	slug := ""
	if !isAll {
		if len(positional) > 0 {
			slug = positional[0]
		}
		if strings.TrimSpace(slug) == "" {
			return nil, apperror.New("parseExportFlags", "E_INSTALLER_INVALID_INPUT", map[string]any{
				"error": "installer slug is required for single export",
			})
		}
	}

	return &ExportInstallerFlags{
		Slug:       strings.TrimSpace(slug),
		OutputPath: targetPath,
		ExportAll:  isAll,
	}, nil
}

// writeZipEntry writes a JSON-marshalled installer script into a zip archive.
func writeZipEntry(zw *zip.Writer, script model.InstallerScript) error {
	data, errMarshal := json.MarshalIndent(script, "", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	filename := fmt.Sprintf("%s.json", script.Slug)
	w, errCreate := zw.Create(filename)
	if errCreate != nil {
		return errCreate
	}

	_, errWrite := w.Write(data)
	return errWrite
}

// executeExport writes matching installer scripts to the specified zip file.
//nolint:revive
func executeExport(ctx context.Context, db *store.DB, flags *ExportInstallerFlags) error {
	if db == nil {
		return apperror.New("executeExport", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "db cannot be nil"})
	}

	var scripts []model.InstallerScript
	if flags.ExportAll {
		list, errList := db.ListInstallers()
		if errList != nil {
			return errList
		}
		scripts = list
	} else {
		item, errGet := db.GetInstallerBySlug(flags.Slug)
		if errGet != nil {
			return errGet
		}
		scripts = append(scripts, *item)
	}

	if err := os.MkdirAll(filepath.Dir(flags.OutputPath), 0755); err != nil && filepath.Dir(flags.OutputPath) != "." {
		appErr := apperror.Wrap(err, "executeExport", map[string]any{"path": flags.OutputPath})
		appErr.Code = "E_INSTALLER_EXPORT_FAILED"
		return appErr
	}

	outFile, errCreate := os.Create(flags.OutputPath)
	if errCreate != nil {
		appErr := apperror.Wrap(errCreate, "executeExport", map[string]any{"path": flags.OutputPath})
		appErr.Code = "E_INSTALLER_EXPORT_FAILED"
		return appErr
	}
	defer outFile.Close()

	zw := zip.NewWriter(outFile)
	defer zw.Close()

	for _, s := range scripts {
		if errWrite := writeZipEntry(zw, s); errWrite != nil {
			appErr := apperror.Wrap(errWrite, "executeExport", map[string]any{"slug": s.Slug})
			appErr.Code = "E_INSTALLER_EXPORT_FAILED"
			return appErr
		}
	}

	fmt.Printf("Exported %d installer script(s) to %s successfully.\n", len(scripts), flags.OutputPath)
	return nil
}

// runInstallerExport coordinates flag parsing, database connection, and export execution.
func runInstallerExport(cmd *cobra.Command, args []string, isAll bool) error {
	if cmd == nil {
		return apperror.New("runInstallerExport", "E_INSTALLER_NIL_COMMAND", map[string]any{"error": "command is nil"})
	}

	flags, errParse := parseExportFlags(args, isAll)
	if errParse != nil {
		return errParse
	}

	db, errDB := store.OpenDefault()
	if errDB != nil {
		appErr := apperror.Wrap(errDB, "runInstallerExport", map[string]any{"action": "open_db"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}
	defer db.Close()

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		appErr := apperror.Wrap(errMigrate, "runInstallerExport", map[string]any{"action": "migrate_installers"})
		appErr.Code = "E_INSTALLER_DB_ERROR"
		return appErr
	}

	return executeExport(cmd.Context(), db, flags)
}
