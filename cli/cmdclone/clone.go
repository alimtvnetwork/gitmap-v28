package cmdclone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

// applySSHKey sets GIT_SSH_COMMAND if an SSH key name is provided.
func applySSHKey(name string) {
	if len(name) == 0 {
		return
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		cliexit.HandleError(err, 1)
	}

	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, name)
		cliexit.HandleError(err, 1)
	}

	sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", key.PrivatePath)
	os.Setenv("GIT_SSH_COMMAND", sshCmd)
	fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)
}

func handleCloneFixFlag(cf CloneFlags) bool {
	if !cf.Fix {
		return false
	}

	runCloneFixOptimization()

	return len(cf.Source) == 0
}

func runCloneFixOptimization() {
	fmt.Printf("%s Running workspace project optimization and repeat removal...\n", constants.ColorCyan+"▸"+constants.ColorReset)
	_ = cmdagy.RunOptimize()
	summary, _ := vscodepm.OptimizeProjects(nil, false)
	printVSCodeOptimizeResult(summary, false)
	_ = RunGitHubDesktopOptimize(nil)
}

func resolveCloneManifest(source string) string {
	if len(source) > 0 {
		return resolveExplicitCloneSource(source)
	}
	discovered := discoverDefaultCloneManifest()
	if len(discovered) > 0 {
		return discovered
	}
	fmt.Fprintln(os.Stderr, "Error: no clone manifest found (searched *.json, gitmap.json, .gitmap/gitmap.json)")
	fmt.Fprintln(os.Stderr, constants.ErrCloneUsage)
	cliexit.HandleError(nil, 1)
	return ""
}

func resolveExplicitCloneSource(source string) string {
	if isRegularFile(source) {
		return source
	}
	if resolved := secrets.ResolveRepoSecretsManifest(source, "gitmap.json"); resolved != "" {
		return resolved
	}
	return source
}

func isCloneSubcommand(arg string) string {
	low := strings.TrimPrefix(strings.ToLower(arg), "--")
	switch low {
	case "only-missing", "onlymissing":
		return "only-missing"
	case "next":
		return "next"
	case "now":
		return "now"
	case "sync":
		return "sync"
	case "pick":
		return "pick"
	case "fix-repo", "fixrepo", "fr", "cfr":
		return "fix-repo"
	case "rc", "repo-cache":
		return "rc"
	default:
		return ""
	}
}

func dispatchCloneSubcommand(subcmd string, subArgs []string) error {
	switch subcmd {
	case "only-missing":
		return RunCloneOnlyMissing(subArgs)
	case "next":
		return RunCloneNext(subArgs)
	case "now":
		return RunCloneNow(subArgs)
	case "sync":
		return RunCloneSync()
	case "pick":
		return RunClonePick(subArgs)
	case "fix-repo":
		return RunCloneFixRepo(subArgs)
	case "rc":
		return RunCloneRC(subArgs, CloneModeStandard)
	default:
		return nil
	}
}

func resolveCloneSubcommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}
	sub := isCloneSubcommand(args[0])
	if len(sub) == 0 {
		return "", nil
	}
	return sub, args[1:]
}

// runClone handles the "clone" subcommand.
func runClone(args []string) error {
	checkHelp("clone", args)
	subcmd, subArgs := resolveCloneSubcommand(args)
	if len(subcmd) > 0 {
		return dispatchCloneSubcommand(subcmd, subArgs)
	}
	if hasRCFlagOrToken(args) {
		return RunCloneRC(args, CloneModeStandard)
	}
	cf := parseCloneFlags(args)

	return executeParsedClone(cf)
}

func executeParsedClone(cf CloneFlags) error {
	if handleCloneFixFlag(cf) {
		return nil
	}
	cf.Source = resolveCloneManifest(cf.Source)
	if cf.IsListOnly {
		return runCloneListTable(cf.Source, cf.TargetDir, cf.UseSSH, cf.OnlyFilter, cf.ExcludeFilter)
	}
	prepareCloneEnv(cf)
	if cf.Audit {
		runCloneAudit(cf)
		maybeExitOnCmdFaithfulMismatch()

		return nil
	}
	cf = prepareCloneTransport(cf)

	err := dispatchCloneExecution(cf)
	if err == nil && !cf.DryRun && !cf.IsListOnly && !cf.Audit && !cf.IsJSON {
		MaybePrintFleetCloneSuggestion(cf.Source)
	}
	return err
}

func prepareCloneEnv(cf CloneFlags) {
	initCloneVerbose(cf.Verbose)
	SetCloneDryRun(cf.DryRun)
	SetCloneAssumeYes(cf.IsAssumeYes)
	setCmdFaithfulVerify(cf.VerifyCmdFaithful)
	setCmdFaithfulExitOnMismatch(cf.VerifyCmdFaithfulExitOnMismatch)
	setCmdPrintArgv(cf.PrintCloneArgv)
}

func prepareCloneTransport(cf CloneFlags) CloneFlags {
	requireOnline()
	applySSHKey(cf.SSHKeyName)
	applyCloneAssumeYesEnv(cf.IsAssumeYes || cf.UseSSH)
	if len(cf.AccessToken) > 0 {
		SetGlobalAccessToken(cf.AccessToken)
	}

	return applyURLSchemeFlags(cf)
}

func dispatchCloneExecution(cf CloneFlags) error {
	if isMultiCloneEnabled(cf) {
		runCloneMulti(cf)
		maybeExitOnCmdFaithfulMismatch()

		return nil
	}
	if isDirectURL(cf.Source) {
		executeDirectClone(DirectCloneParams{
			URL:          cf.Source,
			FolderName:   cf.FolderName,
			GHDesktop:    cf.GHDesktop,
			NoReplace:    cf.NoReplace,
			Output:       cf.Output,
			NoVSCodeSync: cf.NoVSCodeSync,
			IsClean:      cf.Clean,
			IsJSON:       cf.IsJSON,
			UseSSH:       cf.UseSSH,
			Profile:      cf.Profile,
		})
		maybeExitOnCmdFaithfulMismatch()

		return nil
	}
	cf.Source = resolveCloneShorthand(cf.Source)
	executeClone(cf)
	maybeExitOnCmdFaithfulMismatch()

	return nil
}

// initCloneVerbose sets up verbose logging if enabled.
func initCloneVerbose(enabled bool) {
	if !enabled {
		return
	}

	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

		return
	}

	defer log.Close()
}

// resolveCloneShorthand maps "json", "csv", and "text" to default output paths.
func resolveCloneShorthand(source string) string {
	shorthandMap := map[string]string{
		constants.ShorthandJSON: filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
		constants.ShorthandCSV:  filepath.Join(constants.DefaultOutputFolder, constants.DefaultCSVFile),
		constants.ShorthandText: filepath.Join(constants.DefaultOutputFolder, constants.DefaultTextFile),
	}

	resolved, ok := shorthandMap[strings.ToLower(source)]
	if ok {
		return validateShorthandPath(resolved)
	}

	return source
}

// validateShorthandPath checks that the resolved shorthand file exists.
func validateShorthandPath(resolved string) string {
	_, err := os.Stat(resolved)
	if err == nil {
		return resolved
	}

	fmt.Fprintf(os.Stderr, constants.ErrShorthandNotFound, resolved)
	cliexit.HandleError(nil, 1)

	return ""
}
