package cmdinstall

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	fs.BoolVar(&opts.Tree, "tree", false, "Preview full tool hierarchy of a profile before installing")
	fs.BoolVar(&opts.Tree, "t", false, "Preview tree")
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

	if isInstallProfileCommand(args) {
		return runInstallProfileCommand(args)
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

	if opts.Tree && opts.Tool == "" {
		renderAllProfilesTree(loadInstalledLookup())

		return nil
	}

	if opts.Tool == "" {
		return handleMissingInstallTool()
	}

	if IsInstallProfile(opts.Tool) {
		return runInstallProfile(opts.Tool, opts)
	}

	opts.Tool = resolveToolAlias(opts.Tool)
	validateToolName(opts.Tool)
	executeInstall(opts)

	return nil
}

func handleMissingInstallTool() error {
	printInstallListGrouped()
	printInstallUsageHints()

	return nil
}

func printInstallUsageHints() {
	fmt.Fprintf(os.Stderr, "Usage:\n  gitmap install <tool|profile> [flags]\n  gitmap in <tool|profile> [flags]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n  --list, ls, list       List all available developer tools & profiles\n  profile <name>         Run an installation profile (dev, ubuntu, ai, minimal, base, fullstack...)\n  --tree, -t             Preview full tool tree of a profile before installing\n  --logs, logs           View installation execution logs\n  --help                 Show detailed install help and examples\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n  $ gitmap install qtorrent\n  $ gitmap install utorrent\n  $ gitmap export-config qtorrent\n  $ gitmap import-config utorrent\n  $ gitmap install profile dev --tree\n  $ gitmap in dev\n  $ gitmap in logs\n\n")
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
	Tree    bool
}

func isKnownInstallTool(tool string) bool {
	canonical := resolveToolAlias(tool)
	if isCleanCodeAlias(canonical) || isBuildEssentialAlias(canonical) || IsInstallProfile(canonical) {
		return true
	}

	_, exists := constants.InstallToolDescriptions[canonical]

	return exists
}

func isInstallProfileCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}

	low := strings.ToLower(args[0])

	return low == "profile" || low == "profiles"
}

func isProfileListKeyword(arg string) bool {
	low := strings.ToLower(strings.TrimSpace(arg))

	return low == "" || low == "--list" || low == "list" || low == "ls"
}

func runInstallProfileCommand(args []string) error {
	if len(args) <= 1 || isProfileListKeyword(args[1]) {
		printInstallProfilesOnly()

		return nil
	}

	if isProfileTreeKeyword(args[1]) {
		return handleProfileTreeCommand(args[1:])
	}

	return dispatchInstallProfile(args[1:])
}

func dispatchInstallProfile(args []string) error {
	opts, list := parseInstallFlags(args)
	if list || opts.Tool == "" || isProfileListKeyword(opts.Tool) {
		printInstallProfilesOnly()

		return nil
	}

	return runInstallProfile(opts.Tool, opts)
}

func isProfileTreeKeyword(arg string) bool {
	low := strings.ToLower(arg)

	return low == "tree" || low == "--tree" || low == "-t"
}

func handleProfileTreeCommand(args []string) error {
	installed := loadInstalledLookup()
	if len(args) <= 1 {
		renderAllProfilesTree(installed)

		return nil
	}

	opts, _ := parseInstallFlags(args[1:])
	opts.Tree = true

	return runInstallProfile(opts.Tool, opts)
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
