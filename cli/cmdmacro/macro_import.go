package cmdmacro

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

type macroImportOpts struct {
	TargetName string
	FilePath   string
	Format     string
	RenameAs   string
	IsForce    bool
	IsDryRun   bool
	IsAll      bool
	IsSingle   bool
	ExceptList []string
}

func runMacroImport(args []string) error {
	opts := parseMacroImportOpts(args)
	if opts.FilePath == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro import <file> [name] [--all|--single] [--as <name>] [--force] [--dry-run]\n")

		return apperror.NewSimple("import file path required", "E6021")
	}

	macros, err := macro.ParseImportFile(opts.FilePath, opts.Format)
	if err != nil {
		return err
	}

	res, err := macro.ImportMacros(macros, macro.ImportOptions{
		TargetName: opts.TargetName,
		Format:     opts.Format,
		FilePath:   opts.FilePath,
		IsForce:    opts.IsForce,
		IsDryRun:   opts.IsDryRun,
		IsSingle:   opts.IsSingle,
		IsAll:      opts.IsAll,
		RenameAs:   opts.RenameAs,
		ExceptList: opts.ExceptList,
	})
	if err != nil {
		return err
	}

	printMacroImportSummary(res, opts.IsDryRun)

	return nil
}

func parseMacroImportOpts(args []string) macroImportOpts {
	opts := macroImportOpts{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		processImportFlag(a, args, &i, &opts)
	}

	reconcileImportPaths(&opts)

	return opts
}

func reconcileImportPaths(opts *macroImportOpts) {
	if opts.FilePath == "" && opts.TargetName != "" && fileExists(opts.TargetName) {
		opts.FilePath = opts.TargetName
		opts.TargetName = ""

		return
	}

	if !fileExists(opts.FilePath) && fileExists(opts.TargetName) {
		opts.FilePath, opts.TargetName = opts.TargetName, opts.FilePath
	}

	disambiguateImportFileExtensions(opts)
	cleanImportTargetKeywords(opts)
}

func isImportPathSwapNeeded(opts *macroImportOpts) bool {
	neitherFileExists := !fileExists(opts.FilePath) && !fileExists(opts.TargetName)
	hasTargetExtOnly := hasKnownImportExtension(opts.TargetName) && !hasKnownImportExtension(opts.FilePath)

	return neitherFileExists && hasTargetExtOnly
}

func disambiguateImportFileExtensions(opts *macroImportOpts) {
	if isImportPathSwapNeeded(opts) {
		opts.FilePath, opts.TargetName = opts.TargetName, opts.FilePath
	}
}

func hasKnownImportExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".yaml", ".yml", ".db", ".sqlite", ".sqlite3", ".sqlitedb", ".zip":
		return true
	default:
		return false
	}
}

func cleanImportTargetKeywords(opts *macroImportOpts) {
	if opts.TargetName == "all" || opts.TargetName == "*" {
		opts.IsAll = true
		opts.TargetName = ""
	}

	if opts.TargetName == "single" {
		opts.IsSingle = true
		opts.TargetName = ""
	}
}

func processImportFlag(arg string, args []string, index *int, opts *macroImportOpts) {
	if processImportFormatFlag(arg, args, index, opts) {
		return
	}

	if processImportModifierFlag(arg, opts) {
		return
	}

	processImportParamFlag(arg, args, index, opts)
}

func processImportFormatFlag(arg string, args []string, index *int, opts *macroImportOpts) bool {
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

func processImportModifierFlag(arg string, opts *macroImportOpts) bool {
	switch {
	case arg == "--all" || arg == "all":
		opts.IsAll = true
		opts.IsSingle = false
		return true
	case arg == "--single" || arg == "single":
		opts.IsSingle = true
		opts.IsAll = false
		return true
	case arg == "--force" || arg == "--overwrite":
		opts.IsForce = true
		return true
	case arg == "--dry-run":
		opts.IsDryRun = true
		return true
	default:
		return false
	}
}

func processImportParamFlag(arg string, args []string, index *int, opts *macroImportOpts) {
	switch {
	case matchFlagWithVal(arg, "-f", "--file", "--filepath", "-i", "--in"):
		opts.FilePath = extractFlagValue(index, args)
	case matchFlagWithVal(arg, "-except", "--except", "--exclude"):
		opts.ExceptList = parseExceptTokens(extractFlagValue(index, args))
	case matchFlagWithVal(arg, "--name", "--target"):
		opts.TargetName = extractFlagValue(index, args)
	case matchFlagWithVal(arg, "--as", "--rename"):
		opts.RenameAs = extractFlagValue(index, args)
	case !strings.HasPrefix(arg, "-") && opts.FilePath == "":
		opts.FilePath = arg
	case !strings.HasPrefix(arg, "-") && opts.TargetName == "":
		opts.TargetName = arg
	}
}

func printMacroImportSummary(res *macro.ImportResult, isDryRun bool) {
	if isDryRun {
		fmt.Printf("\n  %sℹ [Dry Run]%s Found %d macro(s): %d would import, %d would overwrite, %d skipped\n",
			constants.ColorYellow, constants.ColorReset,
			res.TotalFound, res.Imported, res.Overwritten, res.Skipped)
		printSkippedDetails(res.SkippedNames)

		return
	}

	totalSaved := res.Imported + res.Overwritten
	fmt.Printf("\n  %s✔ Successfully imported %d macro(s)%s into store (%d overwritten, %d skipped)\n",
		constants.ColorGreen, totalSaved, constants.ColorReset,
		res.Overwritten, res.Skipped)
	printSkippedDetails(res.SkippedNames)
}

func printSkippedDetails(skippedNames []string) {
	if len(skippedNames) > 0 {
		fmt.Printf("  %sSkipped existing:%s %s (use --force to overwrite)\n\n",
			constants.ColorDim, constants.ColorReset, strings.Join(skippedNames, ", "))
	} else {
		fmt.Println()
	}
}
