package cmd

import (
	"flag"
	"fmt"
	"os"

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

// runInstall handles the "install" command.
func runInstall(args []string) {
	checkHelp("install", args)
	opts, list := parseInstallFlags(args)
	if list || opts.Tool == "ls" || opts.Tool == "list" {
		printInstallListGrouped()
		return
	}
	if opts.Tool == "" {
		fmt.Fprint(os.Stderr, constants.ErrInstallToolRequired)
		os.Exit(1)
	}
	validateToolName(opts.Tool)
	executeInstall(opts)
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

// validateToolName checks if the tool is supported.
func validateToolName(tool string) {
	if isCleanCodeAlias(tool) {
		return
	}
	if _, exists := constants.InstallToolDescriptions[tool]; exists {
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrInstallUnknownTool, tool)
	os.Exit(1)
}

// executeInstall runs the install flow for a tool.
func executeInstall(opts installOptions) {
	if handler := specialInstallHandler(opts.Tool); handler != nil {
		handler(opts)
		return
	}
	executeGenericInstall(opts)
}

func specialSyncHandler(tool string) func(installOptions) {
	return map[string]func(installOptions){
		constants.ToolScripts:     func(installOptions) { runInstallScripts() },
		constants.ToolNppSettings: func(installOptions) { runNppSettingsOnly() },
		constants.ToolVSCodeSync:  func(installOptions) { runVSCodeSettingsOnly() },
		constants.ToolOBSSync:     func(installOptions) { runOBSSettingsOnly() },
		constants.ToolWTSync:      func(installOptions) { runWTSettingsOnly() },
	}[tool]
}

func specialToolHandler(tool string) func(installOptions) {
	return map[string]func(installOptions){
		constants.ToolVSCodeCtx:      func(installOptions) { runVSCodeContextMenu() },
		constants.ToolPwshCtx:        func(installOptions) { runPwshContextMenu() },
		constants.ToolCtx:            func(opts installOptions) { runInstallCtx(opts.Explain) },
		constants.ToolAllDevTools:    func(opts installOptions) { runAllDevTools(opts) },
		constants.ToolGitmapOneliner: func(installOptions) { runInstallGitmapOneliner() },
		constants.ToolScriptsFixer:     func(installOptions) { runInstallCustomTool("scripts-fixer") },
		constants.ToolCodingGuidelines: func(installOptions) { runInstallCustomTool("coding-guidelines") },
		constants.ToolMacroAhk:         func(installOptions) { runInstallCustomTool("macro-ahk") },
	}[tool]
}

func specialInstallHandler(tool string) func(installOptions) {
	if isCleanCodeAlias(tool) {
		return func(installOptions) { runInstallCleanCode() }
	}
	if h := specialSyncHandler(tool); h != nil {
		return h
	}
	return specialToolHandler(tool)
}

func shouldProceedInstall(opts installOptions, installName string) bool {
	if alreadyInstalled(installName) {
		return false
	}
	if opts.Check {
		fmt.Printf(constants.MsgInstallNotFound, installName)
		return false
	}
	return true
}

func runToolInstallation(opts installOptions, originalTool, installName string) {
	opts.Tool = installName
	manager := resolvePackageManager(opts.Manager, opts.Tool)
	announceInstallPlan(opts.Version, manager)
	if confirmInstallIfNeeded(opts, installName, manager) {
		installTool(opts)
		postInstallSettingsSync(originalTool)
	}
}

// executeGenericInstall runs the standard package-manager-backed install pipeline.
func executeGenericInstall(opts installOptions) {
	originalTool := opts.Tool
	installName := resolveNppInstallName(opts.Tool)
	if !shouldProceedInstall(opts, installName) {
		return
	}
	runToolInstallation(opts, originalTool, installName)
}

// alreadyInstalled prints the "checking" line and probes the installed version.
func alreadyInstalled(installName string) bool {
	fmt.Printf(constants.MsgInstallChecking, installName)
	existingVersion := detectInstalledVersion(installName)
	if existingVersion == "" {
		return false
	}
	fmt.Printf(constants.MsgInstallFound, installName, existingVersion)
	return true
}

// announceInstallPlan prints the version and manager that will be used.
func announceInstallPlan(version, manager string) {
	if version != "" {
		fmt.Printf(constants.MsgInstallVersion, version)
	} else {
		fmt.Print(constants.MsgInstallVersionLabel)
	}
	fmt.Printf(constants.MsgInstallManager, manager)
}

// confirmInstallIfNeeded honors the --yes and --dry-run flags.
func confirmInstallIfNeeded(opts installOptions, installName, manager string) bool {
	if opts.Yes || opts.DryRun {
		return true
	}
	if confirmInstall(installName, opts.Version, manager) {
		return true
	}
	fmt.Print(constants.MsgInstallAborted)
	return false
}

// postInstallSettingsSync handles the npp-specific post-install behavior.
func postInstallSettingsSync(originalTool string) {
	switch originalTool {
	case constants.ToolNpp:
		runNppSettings()
	case constants.ToolNppInstall:
		fmt.Print(constants.MsgInstallNppSkipSet)
	}
}

// confirmInstall prompts the user for install confirmation.
func confirmInstall(tool, version, manager string) bool {
	if version != "" {
		fmt.Printf(constants.MsgInstallPrompt, tool, version, manager)
	} else {
		fmt.Printf(constants.MsgInstallPromptNoVer, tool, manager)
	}
	var answer string
	fmt.Scanln(&answer)
	return answer == "y" || answer == "Y"
}
