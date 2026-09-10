package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func bindInstallFlags(fs *flag.FlagSet, opts *installOptions, list *bool) {
	fs.StringVar(&opts.Manager, constants.FlagInstallManager, "", constants.FlagDescInstallManager)
	fs.StringVar(&opts.Version, constants.FlagInstallVersion, "", constants.FlagDescInstallVersion)
	fs.BoolVar(&opts.Verbose, constants.FlagInstallVerbose, false, constants.FlagDescInstallVerbose)
	fs.BoolVar(&opts.DryRun, constants.FlagInstallDryRun, false, constants.FlagDescInstallDryRun)
	fs.BoolVar(&opts.Check, constants.FlagInstallCheck, false, constants.FlagDescInstallCheck)
	fs.BoolVar(list, constants.FlagInstallList, false, constants.FlagDescInstallList)
	fs.BoolVar(&opts.Yes, constants.FlagInstallYes, false, constants.FlagDescInstallYes)
	fs.BoolVar(&opts.Yes, "y", false, constants.FlagDescInstallYes)
	fs.BoolVar(&opts.Explain, constants.FlagInstallExplain, false, constants.FlagDescInstallExplain)
}

func parseInstallFlags(args []string) (installOptions, bool) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	var opts installOptions
	var list bool
	bindInstallFlags(fs, &opts, &list)
	fs.Parse(reorderFlagsBeforeArgs(args))
	opts.Tool = fs.Arg(0)

	return opts, list
}

func runInstall(args []string) error {
	if isInstallLogsCommand(args) {

		return runInstallLogs(extractInstallLogsArgs(args))
	}
	if isInstallAddCommand(args) {

		return runInstallAdd(args[1:])
	}
	if isInstallExportCommand(args) {

		return runInstallExport(args[1:])
	}
	if isInstallImportCommand(args) {

		return runInstallImport(args[1:])
	}
	checkHelp("install", args)
	opts, list := parseInstallFlags(args)
	if list || opts.Tool == "ls" || opts.Tool == "list" {
		printInstallListGrouped()

		return nil
	}
	if opts.Tool == "" {

		return handleMissingInstallTool()
	}
	opts.Tool = resolveToolAlias(opts.Tool)
	if customScript := findCustomInstaller(opts.Tool); customScript != nil {

		return executeCustomInstaller(customScript, opts)
	}
	validateToolName(opts.Tool)
	executeInstall(opts)

	return nil
}

func handleMissingInstallTool() error {
	fmt.Fprintf(os.Stderr, "%s\n", constants.ErrInstallToolRequired)
	fmt.Fprintf(os.Stderr, "Usage:\n  gitmap install <tool> [flags]\n  gitmap in <tool> [flags]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n  --list, ls, list       List all available developer tools\n  add, create            Create or update custom installer interactively\n  export                 Export installer(s) to JSON or ZIP\n  import                 Import installer(s) from JSON or ZIP\n  --logs, logs           View installation execution logs\n  --help                 Show detailed install help and examples\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n  $ gitmap install vscode\n  $ gitmap install node\n  $ gitmap install add \"my-tool\" v1.0\n  $ gitmap install export my-tool -o my-tool.json\n  $ gitmap install import my-tool.json\n  $ gitmap install logs\n  $ gitmap install --logs --failed\n  $ gitmap install --list\n\n")

	return nil
}

// installOptions holds parsed install flags.
type installOptions struct {
	Tool    string
	Manager string
	Version string
	Verbose bool
	DryRun  bool
	Check   bool
	Yes     bool
	Explain bool
}

func isKnownInstallTool(tool string) bool {
	if isCleanCodeAlias(tool) || tool == "ag-m" || isBuildEssentialAlias(tool) {

		return true
	}
	if hasCustomInstaller(tool) {

		return true
	}
	_, exists := constants.InstallToolDescriptions[tool]

	return exists
}

// validateToolName checks if the tool is supported.
func validateToolName(tool string) {
	canonical := resolveToolAlias(tool)
	if isKnownInstallTool(canonical) {

		return
	}
	fmt.Fprintf(os.Stderr, "\n  %s✗ Unknown tool: '%s'%s\n\n", constants.ColorRed, tool, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "  Use 'gitmap install --list' to see all supported tools.\n")
	fmt.Fprintf(os.Stderr, "  Use 'gitmap install --help' for usage examples.\n\n")
	errMsg := fmt.Sprintf("unknown tool '%s'. Use 'gitmap install --list' to see available tools", tool)
	cliexit.HandleError(apperror.NewSimple(errMsg, "E9000"), 1)
}

// executeInstall runs the install flow for a tool.
func executeInstall(opts installOptions) {
	opts.Tool = resolveToolAlias(opts.Tool)
	if opts.Tool == "ag-m" {
		opts.Tool = constants.ToolAgManager
	}
	if handler := specialInstallHandler(opts.Tool); handler != nil {
		handler(opts)

		return
	}
	executeGenericInstall(opts)
}
