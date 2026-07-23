package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// pullOptions holds parsed pull flags.
type pullOptions struct {
	slug          string
	group         string
	all           bool
	verbose       bool
	stopOnFail    bool
	parallel      int
	onlyAvailable bool
}

// runPull handles the "pull" subcommand.
func runPull(args []string) {
	checkHelp("pull", args)
	cwd, _ := os.Getwd()
	fmt.Printf("→ gitmap pull (cwd: %s)\n", cwd)
	requireOnline()
	// Transport flags (--ssh/--sh/--https/--ht) are only meaningful
	// for the cwd short-circuit; when present we MUST take the cwd
	// path regardless of other flags so the rewrite actually applies.
	useSSH, useHTTPS, rest := extractTransportFlags(args)
	if useSSH || useHTTPS {
		runPullCWDWithTransport(useSSH, useHTTPS, rest)
		return
	}
	opts := parsePullFlags(args)
	if opts.verbose {
		initVerboseLog()
	}
	if shouldPullCWD(opts) {
		fmt.Println("  ↳ cwd is a git repo — running plain `git pull` here")
		runPullCWD()
		return
	}
	if pullNoTargetsHint(opts) {
		return
	}
	records := resolvePullTargets(opts.slug, opts.group, opts.all)
	fmt.Printf("  ↳ resolved %d repo(s) to pull\n", len(records))
	if opts.onlyAvailable {
		records = filterByAvailableUpdates(records)
		if len(records) == 0 {
			fmt.Print(constants.MsgPullNoAvailable)
			return
		}
	}

	taskID, taskDB := beginPullTask(records)
	if taskDB != nil {
		defer taskDB.Close()
	}

	prog := cloner.NewBatchProgress(len(records), "Pull", false)
	prog.SetStopOnFail(opts.stopOnFail)
	executePull(records, prog, opts)
	prog.PrintSummary()
	prog.PrintFailureReport()

	if code := prog.ExitCodeForBatch(); code != 0 {
		failPendingTask(taskDB, taskID, fmt.Sprintf("pull batch failed with exit code %d", code))
		exitWith(code)
	}

	completePendingTask(taskDB, taskID)
}

// shouldPullCWD reports whether `gitmap pull` was invoked with no
// targeting flags AND the current working directory is itself a git
// repo. In that case we short-circuit to a plain `git pull` so the
// command behaves like the muscle-memory `git pull` users expect.
func shouldPullCWD(opts pullOptions) bool {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() {
		return false
	}
	return isGitRepoCWD()
}

// pullNoTargetsHint prints an actionable message when the user runs
// bare `gitmap pull` from a directory that is NOT a git repo and
// without any targeting flag (slug/group/all/alias). Without this
// the command would exit via resolvePullTargets' stderr error, which
// is easy to miss in some terminals — leaving the user staring at a
// blank prompt. Returns true when the hint was printed (caller stops).
func pullNoTargetsHint(opts pullOptions) bool {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() {
		return false
	}
	if isGitRepoCWD() {
		return false
	}
	fmt.Println("  ↳ nothing to pull:")
	fmt.Println("     • current directory is not a git repository")
	fmt.Println("     • no <repo-name>, --group, --all, or -A alias provided")
	fmt.Println("  Try one of:")
	fmt.Println("     gitmap pull <repo-name>")
	fmt.Println("     gitmap pull --all")
	fmt.Println("     gitmap pull --group <group>")
	fmt.Println("     cd <repo> && gitmap pull")
	return true
}

// isGitRepoCWD returns true when the cwd (or an ancestor) is inside a
// git work tree. Uses `git rev-parse --is-inside-work-tree` so worktrees
// and submodules are honoured.
func isGitRepoCWD() bool {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// runPullCWD streams `git pull` in the current directory, forwarding
// stdout/stderr/stdin and propagating the underlying exit code.
func runPullCWD() {
	runPullCWDWithTransport(false, false, nil)
}

// runPullCWDWithTransport is the shared cwd runner. Applies the
// optional transport rewrite (persisting via `git remote set-url`)
// before invoking `git pull`. Extra positional args after the
// transport flags are forwarded verbatim.
func runPullCWDWithTransport(useSSH, useHTTPS bool, extraArgs []string) {
	cwd, _ := os.Getwd()
	if !isGitRepoCWD() {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap pull` inside a repo)")
		exitWith(1)
		return
	}
	if _, _, _, err := ApplyTransportFlag(cwd, useSSH, useHTTPS); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		exitWith(1)
		return
	}
	gitArgs := append([]string{"pull"}, extraArgs...)
	fmt.Printf("→ Running: git %s (cwd: %s)\n", joinForLog(gitArgs), cwd)
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitWith(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "git pull failed: %v\n", err)
		exitWith(1)
	}
}

// extractTransportFlags scans args for --ssh/-ssh/--sh/--https/-https/--ht
// and returns (useSSH, useHTTPS, remaining-args-with-those-removed).
// Used to detect the cwd-transport intent BEFORE handing args to the
// existing parsePullFlags (which doesn't know about transport flags).
func extractTransportFlags(args []string) (bool, bool, []string) {
	var useSSH, useHTTPS bool
	rest := make([]string, 0, len(args))
	for _, a := range args {
		switch a {
		case "--ssh", "-ssh", "--sh", "-sh":
			useSSH = true
		case "--https", "-https", "--ht", "-ht":
			useHTTPS = true
		default:
			rest = append(rest, a)
		}
	}
	return useSSH, useHTTPS, rest
}

// beginPullTask records the pending task entry for this pull batch.
func beginPullTask(records []model.ScanRecord) (int64, *store.DB) {
	workDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine working directory: %v\n", wdErr)
	}
	cmdArgs := buildCommandArgs(append([]string{"pull"}, os.Args[2:]...))
	targetPath := workDir
	if len(records) == 1 {
		targetPath = records[0].AbsolutePath
	}

	return createPendingTask(constants.TaskTypePull, targetPath, workDir, "pull", cmdArgs)
}

// executePull dispatches to either the serial or parallel runner.
func executePull(records []model.ScanRecord, prog *cloner.BatchProgress, opts pullOptions) {
	if opts.parallel > 1 {
		runPullParallel(records, prog, opts.parallel, opts.stopOnFail)
		return
	}
	for _, rec := range records {
		if prog.Stopped() {
			break
		}
		prog.BeginItem(rec.RepoName)
		pullOneRepoTracked(rec, prog)
	}
}

// parsePullFlags parses flags for the pull command.
func parsePullFlags(args []string) pullOptions {
	fs := flag.NewFlagSet(constants.CmdPull, flag.ExitOnError)
	vFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	sFlag := fs.Bool(constants.FlagStopOnFail, false, constants.FlagDescStopOnFail)
	pFlag := fs.Int("parallel", 1, constants.FlagDescPullParallel)
	oFlag := fs.Bool("only-available", false, constants.FlagDescPullOnlyAvailable)
	fs.Parse(args)

	opts := pullOptions{
		group: *gFlag, all: *aFlag, verbose: *vFlag, stopOnFail: *sFlag,
		parallel: *pFlag, onlyAvailable: *oFlag,
	}
	if fs.NArg() > 0 {
		opts.slug = fs.Arg(0)
	}

	return opts
}

// initVerboseLog sets up verbose logging, warning on failure.
func initVerboseLog() {
	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

		return
	}
	log.Close()
}

// resolvePullTargets returns records based on alias, group, all, or slug lookup.
func resolvePullTargets(slug, groupName string, all bool) []model.ScanRecord {
	if HasAlias() {
		return []model.ScanRecord{{
			RepoName:     GetAliasSlug(),
			Slug:         GetAliasSlug(),
			AbsolutePath: GetAliasPath(),
		}}
	}
	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}
	if all {
		return loadAllRecordsDB()
	}
	if len(slug) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrPullSlugRequired)
		fmt.Fprintln(os.Stderr, constants.ErrPullUsage)
		os.Exit(1)
	}

	return lookupBySlugDBFirst(slug)
}

// lookupBySlugDBFirst tries the database first, then falls back to JSON.
func lookupBySlugDBFirst(slug string) []model.ScanRecord {
	db, err := openDB()
	if err == nil {
		defer db.Close()
		repos, dbErr := db.FindBySlug(strings.ToLower(slug))
		if dbErr == nil && len(repos) > 0 {
			return repos
		}
	}

	return lookupBySlugJSON(slug)
}

// lookupBySlugJSON loads gitmap.json and matches by repo name.
func lookupBySlugJSON(slug string) []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadJSONRecords(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrPullLoadFailed, jsonPath, err)
		os.Exit(1)
	}

	return findBySlug(records, slug)
}

// loadJSONRecords reads ScanRecords from a JSON file.
func loadJSONRecords(path string) ([]model.ScanRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []model.ScanRecord
	err = json.NewDecoder(file).Decode(&records)

	return records, err
}

// findBySlug finds records matching the slug (case-insensitive, partial match).
func findBySlug(records []model.ScanRecord, slug string) []model.ScanRecord {
	slugLower := strings.ToLower(slug)
	exact, partial := partitionBySlug(records, slugLower)

	if len(exact) > 0 {
		return exact
	}

	return partial
}

// partitionBySlug separates records into exact and partial matches.
func partitionBySlug(records []model.ScanRecord, slugLower string) ([]model.ScanRecord, []model.ScanRecord) {
	var exact, partial []model.ScanRecord

	for _, r := range records {
		nameLower := strings.ToLower(r.RepoName)
		if nameLower == slugLower {
			exact = append(exact, r)
		} else if strings.Contains(nameLower, slugLower) {
			partial = append(partial, r)
		}
	}

	return exact, partial
}

// pullOneRepo runs safe-pull on a single repo using its absolute path.
func pullOneRepo(rec model.ScanRecord) {
	fmt.Printf(constants.MsgPullStarting, rec.RepoName, rec.AbsolutePath)

	if cloner.IsMissingRepo(rec.AbsolutePath) {
		fmt.Fprintf(os.Stderr, constants.ErrPullNotRepo, rec.AbsolutePath)

		return
	}

	result := cloner.SafePullOne(rec, rec.AbsolutePath)
	if result.Success {
		fmt.Printf(constants.MsgPullSuccess, rec.RepoName)
	} else {
		fmt.Fprintf(os.Stderr, constants.MsgPullFailed, rec.RepoName, result.Error)
	}
}

// pullOneRepoTracked runs safe-pull with progress tracking.
func pullOneRepoTracked(rec model.ScanRecord, prog *cloner.BatchProgress) {
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		prog.Skip()

		return
	}

	result := cloner.SafePullOne(rec, rec.AbsolutePath)
	if result.Success {
		prog.Succeed()
	} else {
		prog.FailWithError(rec.RepoName, result.Error)
	}
}
