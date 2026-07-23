package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// goModOpts holds parsed flags for the gomod command.
type goModOpts struct {
	newPath string
	dryRun  bool
	noMerge bool
	noTidy  bool
	verbose bool
	exts    []string
}

// runGoMod is the entry point for the gomod command.
func runGoMod(args []string) {
	checkHelp("gomod", args)
	opts := parseGoModFlags(args)

	if len(opts.newPath) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrGoModUsage)
		os.Exit(1)
	}

	oldPath := readModulePath()
	validateGoModPreconditions(oldPath, opts.newPath)

	if opts.dryRun {
		runGoModDryRun(oldPath, opts.newPath, opts.exts)

		return
	}

	originalBranch := goModCurrentBranch()
	slug := deriveSlug(opts.newPath)
	backupBranch, featureBranch := createGoModBranches(slug)

	fileCount := replaceModulePath(oldPath, opts.newPath, opts.verbose, opts.exts)
	runGoModTidy(opts.noTidy)
	commitGoModChanges(oldPath, opts.newPath, fileCount)

	if opts.noMerge {
		printGoModSummaryNoMerge(oldPath, opts.newPath, fileCount, backupBranch, featureBranch)

		return
	}

	mergeGoModBranch(originalBranch, featureBranch, opts.newPath)
	printGoModSummary(oldPath, opts.newPath, fileCount, backupBranch, featureBranch, originalBranch)
}

// parseGoModFlags parses flags for the gomod command.
func parseGoModFlags(args []string) goModOpts {
	fs := flag.NewFlagSet(constants.CmdGoMod, flag.ExitOnError)
	dryRun := fs.Bool(constants.FlagGoModDryRun, false, constants.FlagDescGoModDryRun)
	noMerge := fs.Bool(constants.FlagGoModNoMerge, false, constants.FlagDescGoModNoMerge)
	noTidy := fs.Bool(constants.FlagGoModNoTidy, false, constants.FlagDescGoModNoTidy)
	verbose := fs.Bool("verbose", false, constants.FlagDescVerbose)
	extFlag := fs.String(constants.FlagGoModExt, "", constants.FlagDescGoModExt)
	fs.Parse(args)

	newPath := ""
	if fs.NArg() > 0 {
		newPath = fs.Arg(0)
	}

	exts := parseExtFlag(*extFlag)

	return goModOpts{
		newPath: newPath,
		dryRun:  *dryRun,
		noMerge: *noMerge,
		noTidy:  *noTidy,
		verbose: *verbose,
		exts:    exts,
	}
}

// validateGoModPreconditions checks all prerequisites before starting.
func validateGoModPreconditions(oldPath, newPath string) {
	if oldPath == newPath {
		fmt.Printf(constants.MsgGoModNothingRename, oldPath)
		os.Exit(0)
	}

	requireInsideWorkTree()

	if isWorkTreeDirty() {
		fmt.Fprint(os.Stderr, constants.ErrGoModDirtyTree)
		os.Exit(1)
	}
}

// runGoModDryRun previews changes without modifying files.
func runGoModDryRun(oldPath, newPath string, exts []string) {
	files := findFilesWithPath(oldPath, exts)

	fmt.Print(constants.MsgGoModDryHeader)
	fmt.Printf(constants.MsgGoModDryOld, oldPath)
	fmt.Printf(constants.MsgGoModDryNew, newPath)
	fmt.Printf(constants.MsgGoModDryFiles, len(files))

	for _, f := range files {
		fmt.Printf(constants.MsgGoModDryFile, f)
	}
}

// printGoModSummary prints the final summary after merge.
func printGoModSummary(oldPath, newPath string, fileCount int, backup, feature, merged string) {
	fmt.Print(constants.MsgGoModSummary)
	fmt.Printf(constants.MsgGoModOld, oldPath)
	fmt.Printf(constants.MsgGoModNew, newPath)
	fmt.Printf(constants.MsgGoModFiles, fileCount)
	fmt.Printf(constants.MsgGoModBackupBranch, backup)
	fmt.Printf(constants.MsgGoModFeatureBranch, feature)
	fmt.Printf(constants.MsgGoModMergedInto, merged)
}

// printGoModSummaryNoMerge prints the summary when --no-merge is used.
func printGoModSummaryNoMerge(oldPath, newPath string, fileCount int, backup, feature string) {
	fmt.Print(constants.MsgGoModSummary)
	fmt.Printf(constants.MsgGoModOld, oldPath)
	fmt.Printf(constants.MsgGoModNew, newPath)
	fmt.Printf(constants.MsgGoModFiles, fileCount)
	fmt.Printf(constants.MsgGoModBackupBranch, backup)
	fmt.Printf(constants.MsgGoModLeftOn, feature)
}

// runGoModTidy runs go mod tidy unless --no-tidy is set.
func runGoModTidy(noTidy bool) {
	if noTidy {
		return
	}

	err := goModTidy()
	if err != nil {
		fmt.Printf(constants.MsgGoModTidyWarn, err)
	}
}
