package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/macro"
)

type macroImportOpts struct {
	TargetName string
	FilePath   string
	Format     string
	IsForce    bool
	IsDryRun   bool
	ExceptList []string
}

func runMacroImport(args []string) error {
	opts := parseMacroImportOpts(args)
	if opts.FilePath == "" {
		fmt.Fprintf(os.Stderr, "Usage: gitmap macro import <file> [name] [--format json|yaml|sqlite|zip] [--force] [--dry-run]\n")

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

	return opts
}

func processImportFlag(arg string, args []string, index *int, opts *macroImportOpts) {
	switch {
	case matchFlagWithVal(arg, "-f", "--file", "--filepath", "-i", "--in"):
		opts.FilePath = extractFlagValue(index, args)
	case matchFlagWithVal(arg, "-except", "--except", "--exclude"):
		opts.ExceptList = parseExceptTokens(extractFlagValue(index, args))
	case matchFlagWithVal(arg, "--name", "--target"):
		opts.TargetName = extractFlagValue(index, args)
	case matchFlagWithVal(arg, "--format"):
		opts.Format = normalizeMacroFormat(extractFlagValue(index, args))
	case arg == "--sqlite" || arg == "--db" || arg == "--sqlitedb":
		opts.Format = "sqlite"
	case arg == "--force" || arg == "--overwrite":
		opts.IsForce = true
	case arg == "--dry-run":
		opts.IsDryRun = true
	case !strings.HasPrefix(arg, "-") && opts.FilePath == "":
		opts.FilePath = arg
	case !strings.HasPrefix(arg, "-") && opts.TargetName == "":
		opts.TargetName = arg
	}
}

func printMacroImportSummary(res *macro.ImportResult, isDryRun bool) {
	if isDryRun {
		fmt.Printf("\n  %sℹ [Dry Run]%s Found %d macro(s): %d would import, %d would overwrite, %d skipped\n\n",
			constants.ColorYellow, constants.ColorReset,
			res.TotalFound, res.Imported, res.Overwritten, res.Skipped)

		return
	}

	totalSaved := res.Imported + res.Overwritten
	fmt.Printf("\n  %s✔ Successfully imported %d macro(s)%s into store (%d overwritten, %d skipped)\n\n",
		constants.ColorGreen, totalSaved, constants.ColorReset,
		res.Overwritten, res.Skipped)
}
