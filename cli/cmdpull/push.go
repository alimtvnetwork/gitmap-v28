package cmdpull

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// pushOptions holds parsed push flags.
type pushOptions struct {
	slug       string
	group      string
	all        bool
	verbose    bool
	stopOnFail bool
	parallel   int
}

// runPush is the entry point for `gitmap push`.
func runPush(args []string) error {
	checkHelp(constants.CmdPush, args)
	cwd, _ := os.Getwd()
	fmt.Printf("→ gitmap push (cwd: %s)\n", cwd)
	requireOnline()

	useSSH, useHTTPS, rest := ExtractTransportFlags(args)
	if useSSH || useHTTPS {
		runPushCWDWithTransport(useSSH, useHTTPS, rest)

		return nil
	}

	opts := parsePushFlags(args)

	return handlePushExecution(opts, rest)
}

func handlePushExecution(opts pushOptions, rest []string) error {
	if opts.verbose {
		initVerboseLog()
	}

	if isPushCWDEnabled(opts) {
		fmt.Printf("  %s cwd is a git repo — running plain `git push` here\n", resolveSubArrow())
		runPushCWD(rest)

		return nil
	}

	if pushNoTargetsHint(opts) {
		return nil
	}

	return dispatchPushRecords(opts, rest)
}

func dispatchPushRecords(opts pushOptions, rest []string) error {
	records := resolvePullTargets(opts.slug, opts.group, opts.all) // Reusing target resolver from pull.go
	if len(records) == 0 {
		handlePushTargetNotFound(opts)

		return nil
	}
	fmt.Printf("  %s resolved %d repo(s) to push\n", resolveSubArrow(), len(records))

	taskID, taskDB := beginPushTask(records, rest)
	if taskDB != nil {
		defer taskDB.Close()
	}

	return runPushBatchWithProgress(records, taskDB, taskID, opts)
}

func runPushBatchWithProgress(records []model.ScanRecord, taskDB *store.DB, taskID int64, opts pushOptions) error {
	prog := cloner.NewBatchProgress(len(records), "Push", false)
	prog.SetStopOnFail(opts.stopOnFail)
	executePush(records, prog, opts)
	prog.PrintSummary()
	prog.PrintFailureReport()

	if code := prog.ExitCodeForBatch(); code != 0 {
		failPendingTask(taskDB, taskID, fmt.Sprintf("push batch failed with exit code %d", code))
		cliexit.HandleError(apperror.NewExecutionError(fmt.Sprintf("push batch failed with exit code %d", code)), code)

		return nil
	}

	completePendingTask(taskDB, taskID)
	runStatus(buildPushStatusArgs(opts))

	return nil
}

func buildPushStatusArgs(opts pushOptions) []string {
	var statusArgs []string
	if opts.group != "" {
		statusArgs = append(statusArgs, "--group", opts.group)
	} else if opts.all {
		statusArgs = append(statusArgs, "--all")
	}

	return statusArgs
}

// isPushCWDEnabled reports whether `gitmap push` was invoked with no targeting flags
// AND the current working directory is a git repo.
func isPushCWDEnabled(opts pushOptions) bool {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() {
		return false
	}

	return isGitRepoCWD()
}

// pushNoTargetsHint prints a hint if nothing to push.
func pushNoTargetsHint(opts pushOptions) bool {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() || isGitRepoCWD() {
		return false
	}

	cwd, _ := os.Getwd()
	if pushDiscoveredChildRepos(cwd) {
		return true
	}

	printPushNoTargetsHelp()

	return true
}

func pushDiscoveredChildRepos(cwd string) bool {
	childRepos, err := fsutil.DiscoverChildGitRepos(cwd)
	if err != nil || len(childRepos) == 0 {
		return false
	}

	fmt.Printf("→ Discovered %d child repositories in %s for push:\n", len(childRepos), cwd)
	for _, r := range childRepos {
		fmt.Printf("  • %s\n", filepath.Base(r))
		cmd := exec.Command("git", "-C", r, "push")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}

	return true
}

func printPushNoTargetsHelp() {
	fmt.Printf("  %s nothing to push:\n", resolveSubArrow())
	fmt.Println("     • current directory is not a git repository")
	fmt.Println("     • no <repo-name>, --group, --all, or -A alias provided")
	fmt.Println("  Try one of:")
	fmt.Println("     gitmap push <repo-name>")
	fmt.Println("     gitmap push --all")
	fmt.Println("     gitmap push --group <group>")
	fmt.Println("     cd <repo> && gitmap push")
}

func parsePushFlags(args []string) pushOptions {
	fs := flag.NewFlagSet(constants.CmdPush, flag.ExitOnError)
	vFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	sFlag := fs.Bool(constants.FlagStopOnFail, false, constants.FlagDescStopOnFail)
	pFlag := fs.Int("parallel", 0, constants.FlagDescPullParallel)
	fs.Parse(args)

	opts := pushOptions{
		group: *gFlag, all: *aFlag, verbose: *vFlag, stopOnFail: *sFlag,
		parallel: *pFlag,
	}

	if fs.NArg() > 0 {
		opts.slug = fs.Arg(0)
	}

	return opts
}

func runPushCWDWithTransport(useSSH, useHTTPS bool, extraArgs []string) error {
	cwd, _ := os.Getwd()
	if !isGitRepoCWD() {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap push` inside a repo)")
		cliexit.HandleValidationError(apperror.NewValidationError("not a git repository (run `gitmap push` inside a repo)"))

		return nil
	}

	if _, _, _, err := ApplyTransportFlag(cwd, useSSH, useHTTPS); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		cliexit.HandleGeneralError(apperror.WrapSimple(err, "apply transport flag"))

		return nil
	}

	pushWithAutoRebase(cwd, extraArgs)

	return nil
}

func runPushCWD(extraArgs []string) error {
	cwd, _ := os.Getwd()
	pushWithAutoRebase(cwd, extraArgs)

	return nil
}

func beginPushTask(records []model.ScanRecord, rest []string) (int64, *store.DB) {
	workDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine working directory: %v\n", wdErr)
	}

	cmdArgs := buildCommandArgs(append([]string{"push"}, os.Args[2:]...))
	targetPath := workDir
	if len(records) == 1 {
		targetPath = records[0].AbsolutePath
	}

	return createPendingTask("push", targetPath, workDir, "push", cmdArgs)
}
