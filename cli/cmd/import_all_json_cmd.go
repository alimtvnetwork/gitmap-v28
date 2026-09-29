package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdui"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonenvelope"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

// ImportItemResult records the execution outcome of importing a single JSON file.
type ImportItemResult struct {
	FilePath string
	Type     string
	TypeName string
	Status   string // "SUCCESS", "DRY-RUN", "SKIPPED", "FAILED"
	Details  string
}

// runImportAllJSONCLI executes `gitmap import-all-json [globs...] [flags]`.
func runImportAllJSONCLI(args []string) error {
	opts := parseWhichFormatArgs(args)
	files, err := resolveInspectionFiles(opts.FileArgs, opts.TargetDir)
	if err != nil {
		fmt.Printf("Error resolving JSON files: %v\n", err)
		return nil
	}

	if len(files) == 0 {
		printImportAllNoFilesBanner(opts.TargetDir)
		return nil
	}

	isDryRun := hasDryRunFlag(args)
	printImportAllStartBanner(len(files), isDryRun)
	results := executeImportAllJSONBatch(files, isDryRun)
	printImportAllSummary(results, isDryRun)
	return nil
}

func hasDryRunFlag(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "--dry-run" || low == "-n" {
			return true
		}
	}
	return false
}

func printImportAllNoFilesBanner(dir string) {
	fmt.Printf("\nNo JSON files found to import in: %s\n", dir)
	fmt.Printf("Usage: gitmap import-all-json [patterns...] [-y] [--dry-run]\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  gitmap import-all-json *\n")
	fmt.Printf("  gitmap import-all-json *.json\n")
	fmt.Printf("  gitmap import-all-json specific.json\n\n")
}

func printImportAllStartBanner(total int, isDryRun bool) {
	modeLabel := "EXECUTE"
	if isDryRun {
		modeLabel = "DRY-RUN"
	}
	fmt.Printf("\n╔══════════════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║ GITMAP BATCH JSON IMPORTER (IMPORT-ALL-JSON)                                 ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("  Discovered: %d JSON file(s) | Mode: %s\n\n", total, modeLabel)
}

func executeImportAllJSONBatch(files []string, isDryRun bool) []ImportItemResult {
	results := make([]ImportItemResult, 0, len(files))
	for idx, f := range files {
		res := executeImportSingleFile(idx+1, len(files), f, isDryRun)
		results = append(results, res)
	}
	return results
}

func executeImportSingleFile(current, total int, path string, isDryRun bool) ImportItemResult {
	normPath := filepath.ToSlash(path)
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("  [%d/%d] ✖ Failed reading %s: %v\n", current, total, normPath, err)
		return ImportItemResult{FilePath: normPath, Status: "FAILED", Details: fmt.Sprintf("read error: %v", err)}
	}

	desc, attrs, isMatched := jsonenvelope.DetectFormat(content)
	if !isMatched {
		fmt.Printf("  [%d/%d] ✖ Skipped %s (unsupported or non-GitMap JSON format)\n", current, total, normPath)
		return ImportItemResult{FilePath: normPath, Status: "SKIPPED", Details: "Could not match known GitMap JSON format"}
	}

	fmt.Printf("  [%d/%d] Importing %s (%s v%s)... ", current, total, desc.Name, normPath, attrs.Version)
	if isDryRun {
		fmt.Printf("%s[PREVIEW]%s\n", constants.ColorYellow, constants.ColorReset)
		return ImportItemResult{
			FilePath: normPath,
			Type:     desc.Type,
			TypeName: desc.Name,
			Status:   "DRY-RUN",
			Details:  fmt.Sprintf("Would execute: %s", fmt.Sprintf(desc.SuggestedImportCmd, normPath)),
		}
	}

	return runSubsystemImport(normPath, desc)
}

func runSubsystemImport(normPath string, desc jsonenvelope.TypeDescriptor) ImportItemResult {
	res := ImportItemResult{FilePath: normPath, Type: desc.Type, TypeName: desc.Name}

	switch desc.Type {
	case jsonenvelope.TypeSSHNodes:
		importErr := cmdssh.RunSSHNodesImportJSON([]string{normPath})
		return finalizeSubsystemResult(res, importErr, "Enrolled SSH nodes into installation.db")
	case jsonenvelope.TypeMacro:
		importErr := cmdmacro.RunMacroImport([]string{normPath, "--force"})
		return finalizeSubsystemResult(res, importErr, "Registered macros in repository macros table")
	case jsonenvelope.TypeUISettings:
		importErr := cmdui.ImportSettingsFromFile(normPath)
		return finalizeSubsystemResult(res, importErr, "Updated UI settings (~/.gitmap/ui_settings.json)")
	case jsonenvelope.TypeTemplates:
		count, _, _, importErr := store.ImportTemplatesFromFile(normPath, true)
		return finalizeSubsystemResult(res, importErr, fmt.Sprintf("Imported %d project templates", count))
	case jsonenvelope.TypeCommitPullConfig:
		fmt.Printf("%s✔ VALIDATED%s\n", constants.ColorGreen, constants.ColorReset)
		res.Status = "SUCCESS"
		res.Details = "Validated commit-in & pull configuration manifest"
		return res
	default:
		fmt.Printf("%s✔ PROCESSED%s\n", constants.ColorGreen, constants.ColorReset)
		res.Status = "SUCCESS"
		res.Details = desc.SystemImpact
		return res
	}
}

func finalizeSubsystemResult(res ImportItemResult, err error, successMsg string) ImportItemResult {
	if err != nil {
		fmt.Printf("%s✖ FAILED%s (%v)\n", constants.ColorRed, constants.ColorReset, err)
		res.Status = "FAILED"
		res.Details = err.Error()
		return res
	}
	fmt.Printf("%s✔ SUCCESS%s\n", constants.ColorGreen, constants.ColorReset)
	res.Status = "SUCCESS"
	res.Details = successMsg
	return res
}

func printImportAllSummary(results []ImportItemResult, isDryRun bool) {
	fmt.Printf("\n%s================================================================================%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" %sBATCH IMPORT EXECUTION SUMMARY%s\n", constants.ColorBold, constants.ColorReset)
	succeeded, skipped, failed := tallyImportResults(results)
	fmt.Printf(" Total: %d | Succeeded: %d | Skipped: %d | Failed: %d\n", len(results), succeeded, skipped, failed)
	fmt.Printf("%s--------------------------------------------------------------------------------%s\n", constants.ColorDim, constants.ColorReset)

	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "FILE", Align: termtable.AlignLeft, MinWidth: 25},
			{Title: "TYPE", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 10},
			{Title: "DETAILS", Align: termtable.AlignLeft, MinWidth: 30},
		},
		Rows: buildImportSummaryRows(results),
	}
	termtable.PrintTable(cfg)
	fmt.Printf("%s================================================================================%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func tallyImportResults(results []ImportItemResult) (int, int, int) {
	succeeded := 0
	skipped := 0
	failed := 0
	for _, r := range results {
		if r.Status == "SUCCESS" || r.Status == "DRY-RUN" {
			succeeded++
			continue
		}
		if r.Status == "SKIPPED" {
			skipped++
			continue
		}
		failed++
	}
	return succeeded, skipped, failed
}

func buildImportSummaryRows(results []ImportItemResult) []termtable.Row {
	rows := make([]termtable.Row, 0, len(results))
	for _, r := range results {
		statusStr := resolveImportRowStatus(r.Status)
		typeLabel := r.Type
		if typeLabel == "" {
			typeLabel = "-"
		}
		rows = append(rows, termtable.Row{
			Cells: []string{
				r.FilePath,
				typeLabel,
				statusStr,
				r.Details,
			},
		})
	}
	return rows
}

func resolveImportRowStatus(status string) string {
	switch status {
	case "SUCCESS":
		return constants.ColorGreen + "✔ SUCCESS" + constants.ColorReset
	case "DRY-RUN":
		return constants.ColorYellow + "○ PREVIEW" + constants.ColorReset
	case "SKIPPED":
		return constants.ColorDim + "✖ SKIPPED" + constants.ColorReset
	default:
		return constants.ColorRed + "✖ FAILED" + constants.ColorReset
	}
}
