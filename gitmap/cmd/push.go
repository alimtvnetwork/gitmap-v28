package cmd

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
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

	useSSH, useHTTPS, rest := extractTransportFlags(args)
	if useSSH || useHTTPS {
		runPushCWDWithTransport(useSSH, useHTTPS, rest)
		return nil
	}
	opts := parsePushFlags(args)
	if opts.verbose {
		initVerboseLog()
	}

	if shouldPushCWD(opts) {
		fmt.Println("  ↳ cwd is a git repo — running plain `git push` here")
		runPushCWD(rest)
		return nil
	}
	if pushNoTargetsHint(opts) {
		return nil
	}
	records := resolvePullTargets(opts.slug, opts.group, opts.all) // Reusing target resolver from pull.go
	fmt.Printf("  ↳ resolved %d repo(s) to push\n", len(records))

	taskID, taskDB := beginPushTask(records, rest)
	if taskDB != nil {
		defer taskDB.Close()
	}

	prog := cloner.NewBatchProgress(len(records), "Push", false)
	prog.SetStopOnFail(opts.stopOnFail)
	executePush(records, prog, opts)
	prog.PrintSummary()
	prog.PrintFailureReport()

	if code := prog.ExitCodeForBatch(); code != 0 {
		failPendingTask(taskDB, taskID, fmt.Sprintf("push batch failed with exit code %d", code))
		exitWith(code)
	}

	completePendingTask(taskDB, taskID)

	var statusArgs []string
	if opts.group != "" {
		statusArgs = append(statusArgs, "--group", opts.group)
	} else if opts.all {
		statusArgs = append(statusArgs, "--all")
	}
	runStatus(statusArgs)
	return nil
}

// shouldPushCWD reports whether `gitmap push` was invoked with no targeting flags
// AND the current working directory is a git repo.
func shouldPushCWD(opts pushOptions) bool {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() {
		return false
	}
	return isGitRepoCWD()
}

// pushNoTargetsHint prints a hint if nothing to push.
func pushNoTargetsHint(opts pushOptions) bool {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() {
		return false
	}
	if isGitRepoCWD() {
		return false
	}
	cwd, _ := os.Getwd()
	if childRepos, err := fsutil.DiscoverChildGitRepos(cwd); err == nil && len(childRepos) > 0 {
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
	fmt.Println("  ↳ nothing to push:")
	fmt.Println("     • current directory is not a git repository")
	fmt.Println("     • no <repo-name>, --group, --all, or -A alias provided")
	fmt.Println("  Try one of:")
	fmt.Println("     gitmap push <repo-name>")
	fmt.Println("     gitmap push --all")
	fmt.Println("     gitmap push --group <group>")
	fmt.Println("     cd <repo> && gitmap push")
	return true
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
	isNonGitRepoCWD := !isGitRepoCWD()
	if isNonGitRepoCWD {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap push` inside a repo)")
		exitWith(1)
		return nil
	}
	if _, _, _, err := ApplyTransportFlag(cwd, useSSH, useHTTPS); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		exitWith(1)
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

// Actually let's use the real pending task logic.
// For that I'll add a helper beginPushTask
// ... wait, I will implement it at the bottom.

func executePush(records []model.ScanRecord, prog *cloner.BatchProgress, opts pushOptions) {
	workers, ok := cloneconcurrency.Resolve(opts.parallel)
	if !ok {
		return apperror.Wrap(opts.parallel, "constants.ErrCloneMaxConcurrencyInvalid", nil)
	}
	opts.parallel = workers

	if opts.parallel > 1 {
		runPushParallel(records, prog, opts.parallel, opts.stopOnFail)
		return
	}
	for _, rec := range records {
		if prog.Stopped() {
			break
		}
		prog.BeginItem(rec.RepoName)
		pushOneRepoTracked(rec, prog)
	}
}

func pushOneRepoTracked(rec model.ScanRecord, prog *cloner.BatchProgress) {
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		prog.Skip(rec.RepoName)
		return
	}

	result := cloner.SafePushOne(rec, rec.AbsolutePath)
	if result.IsSuccess == false {
		prog.FailWithError(rec.RepoName, result.Error)
		return
	}
	if result.Notes == "up-to-date" {
		prog.UpToDate(rec.RepoName)
		return
	}
	prog.Succeed(rec.RepoName)
}

// pushWithAutoRebase runs `git push`, and on non-fast-forward
// rejection auto-runs `git pull --rebase` then retries once.
func pushWithAutoRebase(cwd string, rest []string) {
	gitArgs := append([]string{"push"}, rest...)
	fmt.Printf("→ Running: git %s (cwd: %s)\n", joinForLog(gitArgs), cwd)
	runErr, stderr := runGitCapturingStderr(gitArgs)
	if runErr == nil {
		return
	}
	isNonNonFastForwardRejection := !isNonFastForwardRejection(stderr)
	if isNonNonFastForwardRejection {
		handleGitExit("git push", runErr)
		return
	}

	fmt.Fprintln(os.Stderr, "↻ push rejected (non-fast-forward) — auto-running `git pull --rebase` and retrying")
	if pullErr := runGitInherit([]string{"pull", "--rebase"}); pullErr != nil {
		fmt.Fprintln(os.Stderr, "✗ auto pull --rebase failed — resolve conflicts then re-run `gitmap push`")
		handleGitExit("git pull --rebase", pullErr)
		return
	}
	fmt.Printf("→ Retrying: git %s (cwd: %s)\n", joinForLog(gitArgs), cwd)
	if retryErr := runGitInherit(gitArgs); retryErr != nil {
		handleGitExit("git push (retry)", retryErr)
	}
}

func runGitCapturingStderr(gitArgs []string) (error, string) {
	var buf bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
	return cmd.Run(), buf.String()
}

func runGitInherit(gitArgs []string) error {
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isNonFastForwardRejection(stderr string) bool {
	lower := strings.ToLower(stderr)
	if !strings.Contains(lower, "[rejected]") && !strings.Contains(lower, "failed to push some refs") {
		return false
	}
	return strings.Contains(lower, "fetch first") || strings.Contains(lower, "non-fast-forward")
}

func handleGitExit(label string, runErr error) {
	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		exitWith(exitErr.ExitCode())
		return
	}
	fmt.Fprintf(os.Stderr, "%s failed: %v\n", label, runErr)
	exitWith(1)
}

func joinForLog(args []string) string {
	out := ""
	for i, a := range args {
		if i > 0 {
			out += " "
		}
		out += a
	}
	return out
}
