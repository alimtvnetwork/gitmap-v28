package cmdclone

import (
	"flag"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CloneFlags holds all parsed clone-command flags and positional args.
type CloneFlags struct {
	Source                          string
	FolderName                      string
	TargetDir                       string
	SSHKeyName                      string
	DefaultBranch                   string
	Positional                      []string
	SafePull                        bool
	GHDesktop                       bool
	NoReplace                       bool
	Verbose                         bool
	Audit                           bool
	MaxConcurrency                  int
	Output                          string
	VerifyCmdFaithful               bool
	VerifyCmdFaithfulExitOnMismatch bool
	PrintCloneArgv                  bool
	NoVSCodeSync                    bool
	UseSSH                          bool
	UseHTTPS                        bool
	DryRun                          bool
	IsAssumeYes                     bool
	Clean                           bool
	MissingOnly                     bool
	Fix                             bool
}

type cloneFlagPointers struct {
	targetFlag        *string
	safePullFlag      *bool
	ghDesktopFlag     *bool
	verboseFlag       *bool
	noReplaceFlag     *bool
	cleanFlag         *bool
	missingOnlyFlag   *bool
	auditFlag         *bool
	maxConcFlag       *int
	sshKeyFlag        *string
	defaultBranchFlag *string
	outputFlag        *string
	verifyFlag        *bool
	verifyExitFlag    *bool
	printArgvFlag     *bool
	noVSCodeSyncFlag  *bool
	debugPathsFlag    *bool
	sshFlag           *bool
	httpsFlag         *bool
	dryRunFlag        *bool
	yesFlag           *bool
	fixFlag           *bool
}

func registerCloneStringFlags(fs *flag.FlagSet, flagPtrs *cloneFlagPointers) {
	flagPtrs.targetFlag = fs.String("target-dir", constants.DefaultDir, constants.FlagDescTargetDir)
	flagPtrs.sshKeyFlag = fs.String("ssh-key", "", "SSH key name for clone")
	fs.StringVar(flagPtrs.sshKeyFlag, "K", "", "SSH key name (short)")
	flagPtrs.defaultBranchFlag = fs.String(constants.FlagScanDefaultBranch, "", constants.FlagDescScanDefaultBranch)
	flagPtrs.outputFlag = fs.String(constants.FlagCloneTermOutput, "", constants.FlagDescCloneTermOutput)
	flagPtrs.maxConcFlag = fs.Int(constants.CloneFlagMaxConcurrency,
		constants.CloneDefaultMaxConcurrency, constants.FlagDescCloneMaxConcurrency)
}

func registerCloneToggles(fs *flag.FlagSet, flagPtrs *cloneFlagPointers) {
	flagPtrs.safePullFlag = fs.Bool("safe-pull", false, constants.FlagDescSafePull)
	flagPtrs.ghDesktopFlag = fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)
	flagPtrs.verboseFlag = fs.Bool("verbose", false, constants.FlagDescVerbose)
	flagPtrs.noReplaceFlag = fs.Bool("no-replace", false, constants.FlagDescCloneNoReplace)
	flagPtrs.cleanFlag = fs.Bool("clean", false, "Forcefully delete the local folder and re-clone")
	flagPtrs.missingOnlyFlag = fs.Bool("missing-only", false, "Skip existing directories entirely")
	flagPtrs.auditFlag = fs.Bool(constants.CloneFlagAudit, false, constants.FlagDescCloneAudit)
	flagPtrs.noVSCodeSyncFlag = fs.Bool(constants.FlagNoVSCodeSync, false, constants.FlagDescNoVSCodeSync)
	flagPtrs.debugPathsFlag = fs.Bool(constants.FlagDebugPaths, false, constants.FlagDescDebugPaths)
	flagPtrs.fixFlag = fs.Bool("fix", false, "Remove repeated projects across tools")
	fs.BoolVar(flagPtrs.fixFlag, "repeat-fix", false, "Alias for --fix")
}

func registerCloneExecutionFlags(fs *flag.FlagSet, flagPtrs *cloneFlagPointers) {
	flagPtrs.verifyFlag = fs.Bool(constants.FlagCloneVerifyCmdFaithful, false,
		constants.FlagDescCloneVerifyCmdFaithful)
	flagPtrs.verifyExitFlag = fs.Bool(constants.FlagCloneVerifyCmdFaithfulExitOnMismatch,
		false, constants.FlagDescCloneVerifyCmdFaithfulExitOnMismatch)
	flagPtrs.printArgvFlag = fs.Bool(constants.FlagClonePrintArgv, false,
		constants.FlagDescClonePrintArgv)
	flagPtrs.sshFlag = fs.Bool("ssh", false,
		"Force every clone URL into `git@host:owner/repo.git` SSH-shorthand form before git runs (auto-converts HTTPS / `ssh://` URLs)")
	fs.BoolVar(flagPtrs.sshFlag, "sh", false, "Short alias for --ssh")
	flagPtrs.httpsFlag = fs.Bool("https", false,
		"Force every clone URL into `https://host/owner/repo.git` form (auto-converts SSH-shorthand / `ssh://` URLs)")
	fs.BoolVar(flagPtrs.httpsFlag, "ht", false, "Short alias for --https")
	flagPtrs.dryRunFlag = fs.Bool(constants.FlagCloneDryRun, false, constants.FlagDescCloneDryRun)
	fs.BoolVar(flagPtrs.dryRunFlag, constants.FlagCloneDryRunShort, false, "Short alias for --dry-run")
	flagPtrs.yesFlag = fs.Bool(constants.FlagCloneYes, false, constants.FlagDescCloneYes)
	fs.BoolVar(flagPtrs.yesFlag, constants.FlagCloneYesShort, false, constants.FlagDescCloneYes)
}

func newCloneFlagSet(fs *flag.FlagSet) *cloneFlagPointers {
	flagPtrs := &cloneFlagPointers{}
	registerCloneStringFlags(fs, flagPtrs)
	registerCloneToggles(fs, flagPtrs)
	registerCloneExecutionFlags(fs, flagPtrs)

	return flagPtrs
}

func populateCloneToggles(cloneOpts *CloneFlags, flagPtrs *cloneFlagPointers) {
	cloneOpts.SafePull = *flagPtrs.safePullFlag
	cloneOpts.GHDesktop = *flagPtrs.ghDesktopFlag
	cloneOpts.NoReplace = *flagPtrs.noReplaceFlag
	cloneOpts.Verbose = *flagPtrs.verboseFlag
	cloneOpts.Audit = *flagPtrs.auditFlag
	cloneOpts.MaxConcurrency = *flagPtrs.maxConcFlag
	cloneOpts.NoVSCodeSync = *flagPtrs.noVSCodeSyncFlag
	cloneOpts.Clean = *flagPtrs.cleanFlag
	cloneOpts.MissingOnly = *flagPtrs.missingOnlyFlag
	cloneOpts.Fix = *flagPtrs.fixFlag
}

func populateCloneExecutionFlags(cloneOpts *CloneFlags, flagPtrs *cloneFlagPointers) {
	cloneOpts.VerifyCmdFaithful = *flagPtrs.verifyFlag
	cloneOpts.VerifyCmdFaithfulExitOnMismatch = *flagPtrs.verifyExitFlag
	cloneOpts.PrintCloneArgv = *flagPtrs.printArgvFlag
	cloneOpts.UseSSH = *flagPtrs.sshFlag
	cloneOpts.UseHTTPS = *flagPtrs.httpsFlag
	cloneOpts.DryRun = *flagPtrs.dryRunFlag
	cloneOpts.IsAssumeYes = *flagPtrs.yesFlag
}

func buildCloneFlags(fs *flag.FlagSet, flagPtrs *cloneFlagPointers) CloneFlags {
	cloneOpts := CloneFlags{
		Source:        resolveCloneSource(fs),
		FolderName:    resolveCloneFolderName(fs),
		TargetDir:     *flagPtrs.targetFlag,
		SSHKeyName:    *flagPtrs.sshKeyFlag,
		DefaultBranch: *flagPtrs.defaultBranchFlag,
		Positional:    fs.Args(),
		Output:        *flagPtrs.outputFlag,
	}

	populateCloneToggles(&cloneOpts, flagPtrs)
	populateCloneExecutionFlags(&cloneOpts, flagPtrs)

	return cloneOpts
}

// ParseCloneFlags parses flags for the clone command.
func ParseCloneFlags(args []string) CloneFlags {
	fs := flag.NewFlagSet(constants.CmdClone, flag.ExitOnError)
	flagPtrs := newCloneFlagSet(fs)
	fs.Parse(reorderFlagsBeforeArgs(args))
	cmdvscode.ApplyDebugPathsEnv(*flagPtrs.debugPathsFlag)

	return buildCloneFlags(fs, flagPtrs)
}

func parseCloneFlags(args []string) CloneFlags {
	return ParseCloneFlags(args)
}

func resolveCloneSource(fs *flag.FlagSet) string {
	if fs.NArg() > 0 {
		return fs.Arg(0)
	}

	return ""
}

func resolveCloneFolderName(fs *flag.FlagSet) string {
	if fs.NArg() <= 1 {
		return ""
	}

	secondArg := fs.Arg(1)
	if IsLikelyURL(secondArg) {
		return ""
	}

	return secondArg
}

// IsLikelyURL checks if rawURL is a recognized URL shape.
func IsLikelyURL(rawURL string) bool {
	lowerURL := strings.ToLower(strings.TrimSpace(rawURL))

	return strings.HasPrefix(lowerURL, "https://") ||
		strings.HasPrefix(lowerURL, "http://") ||
		strings.HasPrefix(lowerURL, "ssh://") ||
		strings.HasPrefix(lowerURL, "git@")
}
