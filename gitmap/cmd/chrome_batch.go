// Package cmd — chrome_batch.go: batch operations for Chrome profiles
// (copy-all, export-all, import-all).
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runChromeCopyAll copies all discovered Chrome profiles to a destination directory.
func runChromeCopyAll(args []string) error {
	checkHelp(constants.SubCmdChromeCopyAll, args)
	dstRoot := ".gitmap/chrome/all-profiles"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		dstRoot = args[0]
	}
	names := availableChromeProfileNames()
	if len(names) == 0 {
		fmt.Println("chrome copy-all: no Chrome profiles found.")
		return nil
	}
	return executeCopyAllProfiles(names, dstRoot)
}

func executeCopyAllProfiles(names []string, dstRoot string) error {
	maxW := maxChromeProfileLabelWidth(names)
	fmt.Printf("\n\033[1;96m▸ chrome copy-all\033[0m  %d profile(s) → \033[1m%s\033[0m\n", len(names), dstRoot)
	copiedCount := 0
	for _, name := range names {
		if copySingleProfileToRoot(name, dstRoot, maxW) {
			copiedCount++
		}
	}
	fmt.Printf("\n\033[1;92m✓ copy-all complete\033[0m  %d/%d profiles copied successfully (gitmap v%s).\n", copiedCount, len(names), constants.Version)
	return nil
}

func copySingleProfileToRoot(name, dstRoot string, maxLabelWidth int) bool {
	srcPath := chromeProfilePath(name)
	dstPath := filepath.Join(dstRoot, name)
	label := formatChromeProfileLabel(name, nil)
	pad := calculateLabelPadding(maxLabelWidth, label)
	fmt.Printf("  • copying %s%s → %s\n", label, strings.Repeat(" ", pad), dstPath)
	fileCount, copyErr := copyChromeProfile(srcPath, dstPath)
	if copyErr != nil {
		fmt.Fprintf(os.Stderr, "    \033[1;91m✗ error:\033[0m %v\n", copyErr)
		return false
	}
	fmt.Printf("    \033[1;92m✓\033[0m %d files copied\n", fileCount)
	return true
}

// runChromeExportAll exports all discovered Chrome profiles.
func runChromeExportAll(args []string) error {
	checkHelp(constants.SubCmdChromeExportAll, args)
	format, positional := parseChromeExportArgs(args)
	outDir := ""
	if len(positional) > 0 {
		outDir = positional[0]
	}
	return runExportAllProfilesToPath(outDir, format)
}

func runExportAllProfilesToPath(outPath, format string) error {
	names := availableChromeProfileNames()
	if len(names) == 0 {
		fmt.Println("chrome export: no Chrome profiles found.")
		return nil
	}
	if outPath == "" {
		outPath = filepath.Join(constants.GitMapDir, "chrome-profiles.zip")
	}
	format = inferExportFormatFromPath(outPath, format)
	if format == constants.OutputZIP && !strings.HasSuffix(strings.ToLower(outPath), constants.ExtZIP) {
		outPath += constants.ExtZIP
	}
	fmt.Printf("\n\033[1;96m▸ chrome export\033[0m  %d profile(s) (format=%s) → \033[1m%s\033[0m\n", len(names), format, outPath)
	return dispatchAllProfilesExport(names, format, outPath)
}

func dispatchAllProfilesExport(names []string, format, outPath string) error {
	switch format {
	case constants.OutputSQLite:
		return handleAllProfilesSQLite(names, outPath)
	case constants.OutputYAML:
		return handleAllProfilesYAML(names, outPath)
	case constants.OutputJSON:
		return handleAllProfilesJSON(names, outPath)
	default:
		return handleAllProfilesZIP(names, outPath)
	}
}

func handleAllProfilesSQLite(names []string, outPath string) error {
	size, err := writeAllChromeProfilesSQLite(names, outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export failed: %v\n", err)
		return err
	}
	rec := chromeExportRecord{SQLitePath: outPath, SQLiteSize: size}
	printChromeArtifacts(rec)
	fmt.Printf("\n\033[1;92m✓ export complete\033[0m  %d profiles saved to SQLite database %s (%d bytes) — gitmap v%s\n", len(names), outPath, size, constants.Version)
	return nil
}

func handleAllProfilesYAML(names []string, outPath string) error {
	size, err := writeAllChromeProfilesYAML(names, outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export failed: %v\n", err)
		return err
	}
	printChromeArtifacts(chromeExportRecord{JSONPath: outPath, JSONSize: size})
	fmt.Printf("\n\033[1;92m✓ export complete\033[0m  %d profiles saved to YAML %s (%d bytes) — gitmap v%s\n", len(names), outPath, size, constants.Version)
	return nil
}

func handleAllProfilesZIP(names []string, outPath string) error {
	size, err := writeAllChromeProfilesZIP(names, outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export failed: %v\n", err)
		return err
	}
	printChromeArtifacts(chromeExportRecord{ZIPPath: outPath, ZIPSize: size})
	fmt.Printf("\n\033[1;92m✓ export complete\033[0m  %d profiles saved to ZIP %s (%d bytes) — gitmap v%s\n", len(names), outPath, size, constants.Version)
	return nil
}

func handleAllProfilesJSON(names []string, outPath string) error {
	if strings.HasSuffix(strings.ToLower(outPath), constants.ExtJSON) {
		return handleSingleFileJSON(names, outPath)
	}
	return exportProfileNameList(names, constants.OutputJSON, outPath)
}

func handleSingleFileJSON(names []string, outPath string) error {
	size, err := writeAllChromeProfilesJSON(names, outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "export failed: %v\n", err)
		return err
	}
	printChromeArtifacts(chromeExportRecord{JSONPath: outPath, JSONSize: size})
	fmt.Printf("\n\033[1;92m✓ export complete\033[0m  %d profiles saved to JSON %s (%d bytes) — gitmap v%s\n", len(names), outPath, size, constants.Version)
	return nil
}

func exportProfileNameList(names []string, format, outDir string) error {
	maxW := maxChromeProfileLabelWidth(names)
	exportedCount := 0
	for _, name := range names {
		if exportSingleProfileToDir(name, format, outDir, maxW) {
			exportedCount++
		}
	}
	fmt.Printf("\n\033[1;92m✓ export-all complete\033[0m  %d/%d profiles exported successfully (gitmap v%s).\n", exportedCount, len(names), constants.Version)
	return nil
}

func exportSingleProfileToDir(name, format, outDir string, maxLabelWidth int) bool {
	srcPath, hasDir := resolveChromeProfileDir(name)
	if !hasDir {
		return false
	}
	outPath := filepath.Join(outDir, name+"."+formatExt(format))
	rec, err := exportChromeFormat(format, srcPath, name, outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  \033[1;91m✗ %s:\033[0m %v\n", name, err)
		return false
	}
	persistChromeProfileSilent(name, srcPath, rec)
	label := formatChromeProfileLabel(name, nil)
	pad := calculateLabelPadding(maxLabelWidth, label)
	fmt.Printf("  \033[1;92m✓\033[0m %s%s → %s\n", label, strings.Repeat(" ", pad), outPath)
	return true
}

func formatExt(format string) string {
	if format == constants.OutputZIP {
		return "zip"
	}
	if format == constants.OutputSQLite {
		return "sqlite"
	}
	if format == constants.OutputYAML {
		return "yaml"
	}
	return "json"
}

// runChromeImportAll imports all profile snapshots from a directory.
func runChromeImportAll(args []string) error {
	checkHelp(constants.SubCmdChromeImportAll, args)
	if len(args) == 0 {
		return apperror.NewSimple("chrome import-all: source directory path required", "E4001")
	}
	srcDir := args[0]
	entries, readErr := os.ReadDir(srcDir)
	if readErr != nil {
		return apperror.WrapSimple(readErr, "chrome import-all: cannot read directory")
	}
	return processImportEntries(entries, srcDir)
}

func processImportEntries(entries []os.DirEntry, srcDir string) error {
	importCount := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(srcDir, e.Name())
		if isImportableSnapshot(path) {
			_ = runChromeProfileImport([]string{path})
			importCount++
		}
	}
	fmt.Printf("\n\033[1;92m✓ import-all complete\033[0m  %d file(s) imported from %s\n", importCount, srcDir)
	return nil
}

func isImportableSnapshot(path string) bool {
	return strings.HasSuffix(path, ".json") || strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".sqlite")
}
