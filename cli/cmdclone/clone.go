package cmdclone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

func lookupSSHKeyPath(name string) (string, error) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)

		return "", err
	}

	defer db.Close()

	return findKeyPrivatePath(db, name)
}

func findKeyPrivatePath(db *store.DB, name string) (string, error) {
	key, err := db.FindSSHKeyByName(name)
	if err != nil {
		return "", err
	}

	return key.PrivatePath, nil
}

func setSSHKeyEnv(name, privatePath string) {
	sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", privatePath)
	os.Setenv("GIT_SSH_COMMAND", sshCmd)
	fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, privatePath)
}

// applySSHKey sets GIT_SSH_COMMAND if an SSH key name is provided.
func applySSHKey(name string) error {
	if len(name) == 0 {
		return nil
	}

	privatePath, err := lookupSSHKeyPath(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, name)

		return err
	}

	setSSHKeyEnv(name, privatePath)

	return nil
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

func resolveCloneManifest(source string) (string, error) {
	if len(source) > 0 {
		return resolveExplicitCloneSource(source), nil
	}

	discovered := discoverDefaultCloneManifest()
	if len(discovered) > 0 {
		return discovered, nil
	}

	fmt.Fprintln(os.Stderr, "Error: no clone manifest found (searched *.json, gitmap.json, .gitmap/gitmap.json)")
	fmt.Fprintln(os.Stderr, constants.ErrCloneUsage)
	printRepoSlugSuggestions("available")

	return "", apperror.NewValidationError(constants.ErrCloneUsage)
}

func resolveExplicitCloneSource(source string) string {
	if fsutil.FileExists(source) {
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

	source, err := resolveCloneManifest(cf.Source)
	if err != nil {
		return err
	}

	cf.Source = source
	if cf.IsListOnly {
		return runCloneListTable(cf.Source, cf.TargetDir, cf.UseSSH, cf.OnlyFilter, cf.ExcludeFilter)
	}

	return executeClonePipeline(cf)
}

func executeClonePipeline(cf CloneFlags) error {
	prepareCloneEnv(cf)
	if cf.Audit {
		runCloneAudit(cf)
		maybeExitOnCmdFaithfulMismatch()

		return nil
	}

	cf, err := prepareCloneTransport(cf)
	if err != nil {
		return err
	}

	return executeAndSuggestClone(cf)
}

func executeAndSuggestClone(cf CloneFlags) error {
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

func prepareCloneTransport(cf CloneFlags) (CloneFlags, error) {
	requireOnline()
	if err := applySSHKey(cf.SSHKeyName); err != nil {
		return cf, err
	}

	applyCloneAssumeYesEnv(cf.IsAssumeYes || cf.UseSSH)
	if len(cf.AccessToken) > 0 {
		SetGlobalAccessToken(cf.AccessToken)
	}

	return applyURLSchemeFlags(cf), nil
}

func dispatchCloneExecution(cf CloneFlags) error {
	if isMultiCloneEnabled(cf) {
		runCloneMulti(cf)
		maybeExitOnCmdFaithfulMismatch()

		return nil
	}

	if isDirectURL(cf.Source) {
		return runDirectClone(cf)
	}

	if fsutil.FileExists(cf.Source) || isRecognizedCloneShorthand(cf.Source) {
		return executeManifestCloneFlow(cf)
	}

	return handleBareSlugClone(cf)
}

func runDirectClone(cf CloneFlags) error {
	err := executeDirectClone(buildDirectCloneParams(cf))
	maybeExitOnCmdFaithfulMismatch()

	return err
}

func handleBareSlugClone(cf CloneFlags) error {
	resolved := ResolveRepoSlugForProfile(cf.Source, cf.Profile)
	if isDirectURL(resolved) {
		if cf.UseSSH {
			if sshURL, ok := ConvertURLToSSH(resolved); ok {
				resolved = sshURL
			}
		}

		cf.Source = resolved

		return runDirectClone(cf)
	}

	printRepoSlugSuggestions(cf.Source)

	return apperror.NewExecutionError(fmt.Sprintf("repository %q not found or resolved to a valid git url", cf.Source))
}

func isRecognizedCloneShorthand(source string) bool {
	low := strings.ToLower(strings.TrimSpace(source))

	return low == constants.ShorthandJSON || low == constants.ShorthandCSV || low == constants.ShorthandText
}

func buildDirectCloneParams(cf CloneFlags) DirectCloneParams {
	return DirectCloneParams{
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
	}
}

func executeManifestCloneFlow(cf CloneFlags) error {
	shorthand, err := resolveCloneShorthand(cf.Source)
	if err != nil {
		return err
	}

	cf.Source = shorthand
	err = executeClone(cf)
	maybeExitOnCmdFaithfulMismatch()

	return err
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
func resolveCloneShorthand(source string) (string, error) {
	shorthandMap := map[string]string{
		constants.ShorthandJSON: filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
		constants.ShorthandCSV:  filepath.Join(constants.DefaultOutputFolder, constants.DefaultCSVFile),
		constants.ShorthandText: filepath.Join(constants.DefaultOutputFolder, constants.DefaultTextFile),
	}

	resolved, ok := shorthandMap[strings.ToLower(source)]
	if !ok {
		return source, nil
	}

	return validateShorthandPath(resolved)
}

// validateShorthandPath checks that the resolved shorthand file exists.
func validateShorthandPath(resolved string) (string, error) {
	_, err := os.Stat(resolved)
	if err == nil {
		return resolved, nil
	}

	fmt.Fprintf(os.Stderr, constants.ErrShorthandNotFound, resolved)

	return "", apperror.NewNotFoundError(fmt.Sprintf(constants.ErrShorthandNotFound, resolved))
}
