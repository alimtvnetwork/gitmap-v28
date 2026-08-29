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

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
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
//nolint:unused
func runPull(args []string) error {
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
		return nil
	}
	opts := parsePullFlags(args)
	if opts.verbose {
		initVerboseLog()
	}
	if shouldPullCWD(opts) {
		fmt.Println("  ↳ cwd is a git repo — running plain `git pull` here")
		runPullCWD()
		return nil
	}
	var records []model.ScanRecord
	if opts.slug == "" && opts.group == "" && !opts.all && !HasAlias() {
		cwd, _ := os.Getwd()
		records = ResolvePullDirectoryTargets(cwd)
		if len(records) == 0 {
			records = findChildrenOfCWD(cwd)
		}
		if len(records) == 0 {
			fmt.Println("  ↳ nothing to pull: no tracked repositories found in or under this directory.")
			return nil
		}
	} else {
		records = resolvePullTargets(opts.slug, opts.group, opts.all)
	}
	fmt.Printf("  ↳ resolved %d repo(s) to pull\n", len(records))
	if opts.onlyAvailable == true {
		records = filterByAvailableUpdates(records)
	}

	isAvailableEmpty := opts.onlyAvailable == true && len(records) == 0
	if isAvailableEmpty == true {
		fmt.Print(constants.MsgPullNoAvailable)
		return nil
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

	var tableRows []model.PullTableRow
	for _, rec := range records {
		sha := gitutil.GetLastCommitSHA(rec.AbsolutePath)
		branch := gitutil.GetActiveBranch(rec.AbsolutePath)
		pr := gitutil.DetectPRStatus(rec.AbsolutePath)
		diag := gitutil.InspectDirtyState(rec.AbsolutePath)
		status := "UP_TO_DATE"
		if diag.IsDirty {
			status = "DIRTY"
		}
		latestBranch := gitutil.GetLatestRemoteBranch(rec.AbsolutePath)
		tableRows = append(tableRows, model.PullTableRow{
			RepoName:     rec.RepoName,
			Branch:       branch,
			LatestBranch: latestBranch,
			LastSHA:      sha,
			PRStatus:     pr,
			PullStatus:   status,
			Duration:     "1.0s",
			IsDirty:      diag.IsDirty,
			Reason:       diag.SummaryReason,
		})
	}
	RenderPullBatchTable(tableRows)

	for _, rec := range records {
		diag := gitutil.InspectDirtyState(rec.AbsolutePath)
		if diag.IsDirty {
			PrintRemediationBox(rec.RepoName, rec.AbsolutePath, diag)
		}
	}

	if code := prog.ExitCodeForBatch(); code != 0 {
		failPendingTask(taskDB, taskID, fmt.Sprintf("pull batch failed with exit code %d", code))
		exitWith(code)
	}

	completePendingTask(taskDB, taskID)
	return nil
}

// shouldPullCWD reports whether `gitmap pull` was invoked with no
// targeting flags AND the current working directory is itself a git
// repo. In that case we short-circuit to a plain `git pull` so the
// command behaves like the muscle-memory `git pull` users expect.
//nolint:unused
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
//nolint:unused
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
//nolint:unused
func isGitRepoCWD() bool {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// runPullCWD streams `git pull` in the current directory, forwarding
// stdout/stderr/stdin and propagating the underlying exit code.
//nolint:unused
func runPullCWD() error {
	runPullCWDWithTransport(false, false, nil)
	return nil
}

// runPullCWDWithTransport is the shared cwd runner. Applies the
// optional transport rewrite (persisting via `git remote set-url`)
// before invoking `git pull`. Extra positional args after the
// transport flags are forwarded verbatim.
//nolint:unused
func runPullCWDWithTransport(useSSH, useHTTPS bool, extraArgs []string) error {
	cwd, _ := os.Getwd()
	isNonGitRepoCWD := !isGitRepoCWD()
	if isNonGitRepoCWD {
		if childRepos, err := fsutil.DiscoverChildGitRepos(cwd); err == nil && len(childRepos) > 0 {
			fmt.Printf("→ Discovered %d child repositories in %s for pull:\n", len(childRepos), cwd)
			for _, r := range childRepos {
				fmt.Printf("  • %s\n", filepath.Base(r))
				gitArgs := append([]string{"-C", r, "pull"}, extraArgs...)
				cmd := exec.Command("git", gitArgs...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
			}
			return nil
		}
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap pull` inside a repo)")
		exitWith(1)
		return nil
	}
	if _, _, _, err := ApplyTransportFlag(cwd, useSSH, useHTTPS); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		exitWith(1)
		return nil
	}
	gitArgs := append([]string{"pull"}, extraArgs...)
	fmt.Printf("→ Running: git %s (cwd: %s)\n", joinForLog(gitArgs), cwd)
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	isExitErr := err != nil && errors.As(err, &exitErr) == true

	if isExitErr == true {
		exitWith(exitErr.ExitCode())
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "git pull failed: %v\n", err)
		exitWith(1)
	}
	return nil
}

// extractTransportFlags scans args for --ssh/-ssh/--sh/--https/-https/--ht
// and returns (useSSH, useHTTPS, remaining-args-with-those-removed).
// Used to detect the cwd-transport intent BEFORE handing args to the
// existing parsePullFlags (which doesn't know about transport flags).
//nolint:unused
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
//nolint:unused
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
//nolint:unused
func executePull(records []model.ScanRecord, prog *cloner.BatchProgress, opts pullOptions) {
	workers, ok := cloneconcurrency.Resolve(opts.parallel)
	if !ok {
		panic("invalid concurrency")
	}
	opts.parallel = workers

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
//nolint:unused
func parsePullFlags(args []string) pullOptions {
	fs := flag.NewFlagSet(constants.CmdPull, flag.ExitOnError)
	vFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	sFlag := fs.Bool(constants.FlagStopOnFail, false, constants.FlagDescStopOnFail)
	pFlag := fs.Int("parallel", 0, constants.FlagDescPullParallel)
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
//nolint:unused
func initVerboseLog() {
	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

		return
	}
	log.Close()
}

// resolvePullTargets returns records based on alias, group, all, or slug lookup.
//nolint:unused
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
		return nil
	}

	return lookupBySlugDBFirst(slug)
}

// lookupBySlugDBFirst tries the database first, then falls back to JSON.
//nolint:unused
func lookupBySlugDBFirst(slug string) []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		return lookupBySlugJSON(slug)
	}

	defer db.Close()
	repos, dbErr := db.FindBySlug(strings.ToLower(slug))

	foundRepos := dbErr == nil && len(repos) > 0
	if foundRepos == true {
		return repos
	}

	return lookupBySlugJSON(slug)
}

// lookupBySlugJSON loads gitmap.json and matches by repo name.
//nolint:unused
func lookupBySlugJSON(slug string) []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadJSONRecords(jsonPath)
	if err != nil {
		return nil
	}

	return findBySlug(records, slug)
}

// loadJSONRecords reads ScanRecords from a JSON file.
//nolint:unused
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
//nolint:unused
func findBySlug(records []model.ScanRecord, slug string) []model.ScanRecord {
	slugLower := strings.ToLower(slug)
	exact, partial := partitionBySlug(records, slugLower)

	if len(exact) > 0 {
		return exact
	}

	return partial
}

// partitionBySlug separates records into exact and partial matches.
//nolint:unused
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
//nolint:unused
func pullOneRepo(rec model.ScanRecord) {
	fmt.Printf(constants.MsgPullStarting, rec.RepoName, rec.AbsolutePath)

	if cloner.IsMissingRepo(rec.AbsolutePath) {
		fmt.Fprintf(os.Stderr, constants.ErrPullNotRepo, rec.AbsolutePath)

		return
	}

	result := cloner.SafePullOne(rec, rec.AbsolutePath)
	if result.IsSuccess {
		fmt.Printf(constants.MsgPullSuccess, rec.RepoName)
	} else {
		fmt.Fprintf(os.Stderr, constants.MsgPullFailed, rec.RepoName, result.Error)
	}
}

// pullOneRepoTracked runs safe-pull with progress tracking.
//nolint:unused
func pullOneRepoTracked(rec model.ScanRecord, prog *cloner.BatchProgress) {
	if cloner.IsMissingRepo(rec.AbsolutePath) {
		prog.Skip(rec.RepoName)

		return
	}

	result := cloner.SafePullOne(rec, rec.AbsolutePath)
	isUpToDate := result.IsSuccess == true && result.Notes == "up-to-date"
	isSucceed := result.IsSuccess == true && result.Notes != "up-to-date"

	if isUpToDate == true {
		prog.UpToDate(rec.RepoName)
	}
	if isSucceed == true {
		prog.Succeed(rec.RepoName)
	}
	if result.IsSuccess == false {
		prog.FailWithError(rec.RepoName, result.Error)
	}
}

// findChildrenOfCWD returns all tracked repos that are inside the given directory.
//nolint:unused
func findChildrenOfCWD(cwd string) []model.ScanRecord {
	all := loadAllRecordsDB()
	var children []model.ScanRecord
	prefix := cwd
	if !strings.HasSuffix(prefix, string(os.PathSeparator)) {
		prefix += string(os.PathSeparator)
	}
	for _, r := range all {
		if strings.HasPrefix(r.AbsolutePath, prefix) || r.AbsolutePath == cwd {
			children = append(children, r)
		}
	}
	return children
}
