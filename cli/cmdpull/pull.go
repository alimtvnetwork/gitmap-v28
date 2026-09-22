package cmdpull

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
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
	autoFix       bool
	yes           bool
	noFix         bool
	isRaw         bool
	useSSH        bool
	useHTTPS      bool
}

// NormalizePullArgs converts positional "all", "pa", or "pull-all" argument into "--all" flag.
func NormalizePullArgs(args []string) []string {
	normalized := make([]string, 0, len(args))
	for _, a := range args {
		if a == "all" || a == "pa" || a == "pull-all" {
			normalized = append(normalized, "--all")
		} else {
			normalized = append(normalized, a)
		}
	}

	return normalized
}

// runPull handles the "pull" subcommand.
func runPull(args []string) error {
	checkHelp("pull", args)
	isPullAll := isPullAllInvocation(args)
	args = NormalizePullArgs(args)
	printPullInvocationHeader(isPullAll)
	requireOnline()
	useSSH, useHTTPS, restArgs := ExtractTransportFlags(args)
	isCWDTransport := !isPullAll && (useSSH || useHTTPS)
	if isCWDTransport {
		runPullCWDWithTransport(useSSH, useHTTPS, restArgs)

		return nil
	}
	opts := resolveParsedPullOptions(restArgs, useSSH, useHTTPS)
	if opts.verbose {
		initVerboseLog()
	}

	return dispatchPullExecution(opts)
}

func resolveParsedPullOptions(args []string, useSSH, useHTTPS bool) pullOptions {
	opts := parsePullFlags(args)
	if useSSH {
		opts.useSSH = true
	}
	if useHTTPS {
		opts.useHTTPS = true
	}

	return opts
}

func isPullAllInvocation(args []string) bool {
	if isPullAllRootCmd() {
		return true
	}

	return hasPullAllArg(args)
}

func isPullAllRootCmd() bool {
	if len(os.Args) <= 1 {
		return false
	}
	first := strings.ToLower(os.Args[1])

	return first == "pull-all" || first == "pa"
}

func hasPullAllArg(args []string) bool {
	for _, a := range args {
		if a == "--all" || a == "-all" || a == "-a" || a == "all" || a == "pa" || a == "pull-all" {
			return true
		}
	}

	return false
}

func printPullInvocationHeader(isPullAll bool) {
	cwd, _ := os.Getwd()
	cmdName := "pull"
	if isPullAll {
		cmdName = "pull-all"
	}
	fmt.Printf("\n  %s→%s %sgitmap %s%s %s(cwd: %s)%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, cmdName, constants.ColorReset,
		constants.ColorDim, cwd, constants.ColorReset)
}

func resolveSubArrow() string {
	if glyphs.Resolve() == glyphs.ModeSafe {
		return "->"
	}

	return "→"
}

func dispatchPullExecution(opts pullOptions) error {
	if isPullCWDEnabled(opts) {
		arrow := resolveSubArrow()
		fmt.Printf("    %s%s%s %scwd is a git repo — running plain `git pull` here%s\n\n",
			constants.ColorCyan, arrow, constants.ColorReset,
			constants.ColorDim, constants.ColorReset)

		return runPullCWD(opts.isRaw)
	}

	return runPullBatch(opts)
}

func runPullBatch(opts pullOptions) error {
	records, isFound := resolvePullBatchRecords(opts)
	if !isFound {
		return nil
	}
	printResolvedPullRepos(len(records))
	filtered := applyPullAvailableFilter(records, opts.onlyAvailable)
	if opts.onlyAvailable && len(filtered) == 0 {
		fmt.Print(constants.MsgPullNoAvailable)

		return nil
	}

	return executePullBatchLifecycle(filtered, opts)
}

func printResolvedPullRepos(count int) {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s resolved %s%d%s repo(s) to pull\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorBold, count, constants.ColorReset)
}

func applyPullAvailableFilter(records []model.ScanRecord, isOnlyAvailable bool) []model.ScanRecord {
	if isOnlyAvailable {
		return filterByAvailableUpdates(records)
	}

	return records
}

func executePullBatchLifecycle(records []model.ScanRecord, opts pullOptions) error {
	taskID, taskDB := beginPullTask(records)
	if taskDB != nil {
		defer taskDB.Close()
	}
	maybeApplyTransportToRecords(records, opts.useSSH, opts.useHTTPS)
	bar := NewPullProgressBar(len(records), false, opts.stopOnFail)
	bar.Start()
	executePull(records, bar, opts)
	bar.Stop()
	sortedStates := sortStatesAlphabetically(bar.States())
	renderPullBatchResults(sortedStates)
	handlePullRemediationForRecords(records, opts)

	return finalizePullBatchTask(taskDB, taskID, bar.Failed())
}

func maybeApplyTransportToRecords(records []model.ScanRecord, useSSH, useHTTPS bool) {
	hasTransport := useSSH || useHTTPS
	if !hasTransport {
		return
	}
	for _, rec := range records {
		applyTransportIfPathExists(rec.AbsolutePath, useSSH, useHTTPS)
	}
}

func applyTransportIfPathExists(path string, useSSH, useHTTPS bool) {
	hasPath := len(path) > 0
	if !hasPath {
		return
	}

	ApplyTransportFlag(path, useSSH, useHTTPS)
}

func sortStatesAlphabetically(states []*PullRepoState) []*PullRepoState {
	sorted := make([]*PullRepoState, len(states))
	copy(sorted, states)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].RepoName) < strings.ToLower(sorted[j].RepoName)
	})

	return sorted
}

func renderPullBatchResults(states []*PullRepoState) {
	var tableRows []model.PullTableRow
	for _, state := range states {
		tableRows = append(tableRows, buildPullTableRowFromState(state))
	}
	RenderPullBatchTable(tableRows)
}

func buildPullTableRowFromState(state *PullRepoState) model.PullTableRow {
	row := initBasePullTableRow(state)
	row.LatestBranch = gitutil.GetLatestRemoteBranch(state.RepoPath)
	row.PRStatus = gitutil.DetectPRStatus(state.RepoPath)
	row.PullStatus = resolveStateStatus(state)
	row.Release = resolveRepoRelease(state.RepoPath, row.LatestBranch)
	row.LastSHA = resolveRowSHA(state)

	return row
}

func resolveRowSHA(state *PullRepoState) string {
	sha := state.NewSHA
	if len(sha) == 0 {
		sha = state.OldSHA
	}
	if len(sha) == 0 && len(state.RepoPath) > 0 {
		sha = gitutil.GetLastCommitSHA(state.RepoPath)
	}

	return ShortenSHA(sha)
}

func resolveRepoRelease(repoPath, latestBranch string) string {
	hasPath := len(repoPath) > 0
	if hasPath == false {
		return "-"
	}

	rel := queryReleaseIdentifier(repoPath, latestBranch)
	if len(rel) > 0 {
		return rel
	}

	return "-"
}

func queryReleaseIdentifier(repoPath, latestBranch string) string {
	tag := gitutil.GetLatestTag(repoPath)
	if len(tag) > 0 {
		return tag
	}

	if rel := extractReleaseFromBranch(latestBranch); len(rel) > 0 {
		return rel
	}

	return readRepoManifestVersion(repoPath)
}

func extractReleaseFromBranch(branch string) string {
	lower := strings.ToLower(branch)
	isReleasePrefix := strings.HasPrefix(lower, "release/")
	trimmed := strings.TrimPrefix(branch, "release/")
	hasRelease := isReleasePrefix && len(trimmed) > 0
	if hasRelease {
		return trimmed
	}

	isVPrefix := strings.HasPrefix(lower, "v")
	if isVPrefix && isSemverLike(branch) {
		return branch
	}

	return ""
}

func isSemverLike(s string) bool {
	raw := strings.TrimPrefix(s, "v")
	parts := strings.Split(raw, ".")
	hasEnoughParts := len(parts) >= 2
	if hasEnoughParts == false {
		return false
	}

	_, err := strconv.Atoi(parts[0])

	return err == nil
}

func readRepoManifestVersion(repoPath string) string {
	vJson := filepath.Join(repoPath, "version.json")
	v := readJsonVersion(vJson)
	if len(v) > 0 {
		return ensureVPrefix(v)
	}

	pkgJson := filepath.Join(repoPath, "package.json")
	pkgV := readJsonVersion(pkgJson)
	if len(pkgV) > 0 {
		return ensureVPrefix(pkgV)
	}

	return ""
}

func readJsonVersion(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) == 0 {
		return ""
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}

	if v, ok := m["Version"].(string); ok && len(v) > 0 {
		return v
	}
	if v, ok := m["version"].(string); ok && len(v) > 0 {
		return v
	}

	return ""
}

func ensureVPrefix(v string) string {
	trimmed := strings.TrimSpace(v)
	hasV := strings.HasPrefix(trimmed, "v")
	if hasV {
		return trimmed
	}

	return "v" + trimmed
}

func initBasePullTableRow(state *PullRepoState) model.PullTableRow {
	return model.PullTableRow{
		RepoName:    state.RepoName,
		Branch:      state.Branch,
		LastSHA:     ShortenSHA(state.NewSHA),
		CommitRange: state.CommitRange,
		Changes:     state.Changes,
		Duration:    FormatPullDuration(state.Duration),
		IsDirty:     state.IsDirty,
		Reason:      state.ErrorMsg,
	}
}

func resolveStateStatus(state *PullRepoState) string {
	if state.IsDirty {
		return "DIRTY"
	}

	return string(state.Step)
}

func handlePullRemediationForRecords(records []model.ScanRecord, opts pullOptions) {
	var remItems []RemediationItem
	for _, rec := range records {
		diag := gitutil.InspectDirtyState(rec.AbsolutePath)
		if diag.IsDirty {
			remItems = append(remItems, buildRemediationItem(rec, diag))
		}
	}
	handlePullRemediation(remItems, opts)
}

func buildRemediationItem(rec model.ScanRecord, diag gitutil.DirtyDiagnosis) RemediationItem {
	recipes := gitutil.GenerateRemediationRecipes(rec.AbsolutePath, diag)

	return RemediationItem{
		RepoName:      rec.RepoName,
		RepoPath:      rec.AbsolutePath,
		SummaryReason: diag.SummaryReason,
		Recipes:       recipes,
		Files:         diag.AllFiles,
	}
}

func finalizePullBatchTask(taskDB *store.DB, taskID int64, failCount int) error {
	if failCount > 0 {
		errMsg := fmt.Sprintf("pull batch failed with %d failure(s)", failCount)
		failPendingTask(taskDB, taskID, errMsg)
		cliexit.HandleError(apperror.NewExecutionError(errMsg), 1)

		return nil
	}
	completePendingTask(taskDB, taskID)

	return nil
}

func isPullCWDEnabled(opts pullOptions) bool {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() {
		return false
	}

	return isGitRepoCWD()
}

func isGitRepoCWD() bool {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

// runPullCWD streams pull in CWD, using progress bar unless isRaw is true.
func runPullCWD(isRaw ...bool) error {
	if len(isRaw) > 0 && isRaw[0] {
		return runPullCWDWithTransport(false, false, nil)
	}

	return runPullCWDTracked()
}

func runPullCWDTracked() error {
	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "getwd")
	}

	state := executeCWDTrackedPull(cwd)
	row := buildPullTableRowFromState(state)
	RenderPullBatchTable([]model.PullTableRow{row})

	return nil
}

func executeCWDTrackedPull(cwd string) *PullRepoState {
	rec := model.ScanRecord{RepoName: filepath.Base(cwd), AbsolutePath: cwd}
	bar := NewPullProgressBar(1, false, false)
	bar.Start()
	state := ExecuteTrackedPull(rec, bar)
	bar.Stop()

	return state
}

func runPullCWDWithTransport(useSSH, useHTTPS bool, extraArgs []string) error {
	cwd, _ := os.Getwd()
	if !isGitRepoCWD() {
		return handleNonGitPull(cwd, extraArgs)
	}
	if !applyTransportSafe(cwd, useSSH, useHTTPS) {
		return nil
	}

	return executeGitPullCommand(cwd, extraArgs)
}

func applyTransportSafe(cwd string, useSSH, useHTTPS bool) bool {
	if _, _, _, err := ApplyTransportFlag(cwd, useSSH, useHTTPS); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		cliexit.HandleGeneralError(apperror.WrapSimple(err, "apply transport flag"))

		return false
	}

	return true
}

func executeGitPullCommand(cwd string, extraArgs []string) error {
	gitArgs := append([]string{"pull"}, extraArgs...)
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s Running: git %s (cwd: %s)\n", constants.ColorCyan, arrow, constants.ColorReset, joinForLog(gitArgs), cwd)
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return handleGitExecResult(cmd.Run())
}

func handleGitExecResult(err error) error {
	var exitErr *exec.ExitError
	if err != nil && errors.As(err, &exitErr) {
		cliexit.HandleError(apperror.WrapSimple(exitErr, "git pull"), exitErr.ExitCode())

		return nil
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "git pull failed: %v\n", err)
		cliexit.HandleGeneralError(apperror.WrapSimple(err, "git pull failed"))

		return nil
	}

	return nil
}

func ExtractTransportFlags(args []string) (bool, bool, []string) {
	var useSSH, useHTTPS bool
	rest := make([]string, 0, len(args))
	for _, a := range args {
		switch strings.ToLower(a) {
		case "--ssh", "-ssh", "--sh", "-sh", "ssh":
			useSSH = true
		case "--https", "-https", "--ht", "-ht", "https":
			useHTTPS = true
		default:
			rest = append(rest, a)
		}
	}

	return useSSH, useHTTPS, rest
}

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

func executePull(records []model.ScanRecord, bar *PullProgressBar, opts pullOptions) {
	workers, isResolved := cloneconcurrency.Resolve(opts.parallel)
	if !isResolved {
		cliexit.HandleError(apperror.NewSimple("invalid concurrency", "E9000"), 1)
	}

	opts.parallel = workers
	if opts.parallel > 1 {
		runPullParallel(records, bar, opts.parallel)

		return
	}

	runSerialPull(records, bar)
}

func runSerialPull(records []model.ScanRecord, bar *PullProgressBar) {
	for _, rec := range records {
		if bar != nil && bar.IsStopped() {
			break
		}
		ExecuteTrackedPull(rec, bar)
	}
}

func handlePullRemediation(remItems []RemediationItem, opts pullOptions) {
	if len(remItems) == 0 {
		return
	}
	if opts.noFix {
		PrintRemediationSummaryNoPrompt(remItems)

		return
	}
	if opts.yes || opts.autoFix {
		PrintRemediationSummaryAutoFix(remItems)

		return
	}
	PrintRemediationSummary(remItems)
}

type pullFlagHolders struct {
	vFlag, aFlag, sFlag, oFlag, fixFlag, yFlag, noFixFlag, rawFlag *bool
	sshFlag, httpsFlag                                             *bool
	gFlag                                                          *string
	pFlag                                                          *int
}

func initPullFlagSet() (*flag.FlagSet, *pullFlagHolders) {
	fs := flag.NewFlagSet(constants.CmdPull, flag.ExitOnError)
	h := &pullFlagHolders{}
	registerPullCoreFlags(fs, h)
	registerPullRemediationFlags(fs, h)

	return fs, h
}

func registerPullCoreFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.vFlag = fs.Bool("verbose", false, constants.FlagDescVerbose)
	h.gFlag = fs.String("group", "", constants.FlagDescGroup)
	h.aFlag = fs.Bool("all", false, constants.FlagDescAll)
	h.sFlag = fs.Bool(constants.FlagStopOnFail, false, constants.FlagDescStopOnFail)
	h.pFlag = fs.Int("parallel", 0, constants.FlagDescPullParallel)
	h.oFlag = fs.Bool("only-available", false, constants.FlagDescPullOnlyAvailable)
	h.rawFlag = fs.Bool("raw", false, constants.FlagDescPullRaw)
	h.sshFlag = fs.Bool("ssh", false, "Pull using SSH transport")
	h.httpsFlag = fs.Bool("https", false, "Pull using HTTPS transport")
	fs.StringVar(h.gFlag, "g", "", constants.FlagDescGroup)
}

func registerPullRemediationFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.fixFlag = fs.Bool("fix", false, "Auto-remediate dirty repos")
	h.yFlag = fs.Bool("yes", false, "Remediate without prompt")
	h.noFixFlag = fs.Bool("no-fix", false, "Skip remediation prompt")
	fs.BoolVar(h.yFlag, "y", false, "Remediate without prompt")
}

func buildPullOptions(h *pullFlagHolders) pullOptions {
	return pullOptions{
		group:         *h.gFlag,
		all:           *h.aFlag,
		verbose:       *h.vFlag,
		stopOnFail:    *h.sFlag,
		parallel:      *h.pFlag,
		onlyAvailable: *h.oFlag,
		autoFix:       *h.fixFlag,
		yes:           *h.yFlag,
		noFix:         *h.noFixFlag,
		isRaw:         *h.rawFlag,
		useSSH:        *h.sshFlag,
		useHTTPS:      *h.httpsFlag,
	}
}

func parsePullFlags(args []string) pullOptions {
	fs, h := initPullFlagSet()
	fs.Parse(args)
	opts := buildPullOptions(h)
	if fs.NArg() > 0 {
		opts.slug = fs.Arg(0)
	}

	return opts
}

func initVerboseLog() {
	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

		return
	}
	log.Close()
}

func resolvePullTargets(slug, groupName string, all bool) []model.ScanRecord {
	if HasAlias() {
		return resolveAliasRecord()
	}
	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}
	if all {
		return loadAllRecordsDB()
	}

	return resolveSlugTarget(slug)
}

func resolveAliasRecord() []model.ScanRecord {
	return []model.ScanRecord{{
		RepoName:     GetAliasSlug(),
		Slug:         GetAliasSlug(),
		AbsolutePath: GetAliasPath(),
	}}
}

func resolveSlugTarget(slug string) []model.ScanRecord {
	if len(slug) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrPullSlugRequired)

		return nil
	}

	return lookupBySlugDBFirst(slug)
}

func lookupBySlugDBFirst(slug string) []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		return lookupBySlugJSON(slug)
	}
	defer db.Close()
	repos, dbErr := db.FindBySlug(strings.ToLower(slug))
	foundRepos := dbErr == nil && len(repos) > 0
	if foundRepos {
		return repos
	}

	return lookupBySlugJSON(slug)
}

func lookupBySlugJSON(slug string) []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadJSONRecords(jsonPath)
	if err != nil {
		return nil
	}

	return findBySlug(records, slug)
}

func loadJSONRecords(path string) ([]model.ScanRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open json records")
	}
	defer file.Close()
	var records []model.ScanRecord
	err = json.NewDecoder(file).Decode(&records)
	if err != nil {
		return nil, apperror.WrapSimple(err, "decode json records")
	}

	return records, nil
}

func findBySlug(records []model.ScanRecord, slug string) []model.ScanRecord {
	slugLower := strings.ToLower(slug)
	exact, partial := partitionBySlug(records, slugLower)
	if len(exact) > 0 {
		return exact
	}

	return partial
}

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

func resolvePullBatchRecords(opts pullOptions) ([]model.ScanRecord, bool) {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() {
		return resolvePullTargets(opts.slug, opts.group, opts.all), true
	}
	cwd, _ := os.Getwd()
	records := ResolvePullDirectoryTargets(cwd)
	if len(records) == 0 {
		records = findChildrenOfCWD(cwd)
	}
	if len(records) == 0 {
		arrow := resolveSubArrow()
		fmt.Printf("    %s%s%s %snothing to pull: no tracked repositories found in or under this directory.%s\n\n",
			constants.ColorCyan, arrow, constants.ColorReset,
			constants.ColorDim, constants.ColorReset)

		return nil, false
	}

	return records, true
}

func handleNonGitPull(cwd string, extraArgs []string) error {
	childRepos, err := fsutil.DiscoverChildGitRepos(cwd)
	if err == nil && len(childRepos) > 0 {
		return pullDiscoveredChildren(cwd, childRepos, extraArgs)
	}
	fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap pull` inside a repo)")
	cliexit.HandleValidationError(apperror.NewValidationError("not a git repository (run `gitmap pull` inside a repo)"))

	return nil
}

func pullDiscoveredChildren(cwd string, childRepos []string, extraArgs []string) error {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s Discovered %d child repositories in %s for pull:\n",
		constants.ColorCyan, arrow, constants.ColorReset, len(childRepos), cwd)
	for _, r := range childRepos {
		fmt.Printf("      • %s\n", filepath.Base(r))
		gitArgs := make([]string, 0, 3+len(extraArgs))
		gitArgs = append(gitArgs, "-C", r, "pull")
		gitArgs = append(gitArgs, extraArgs...)
		cmd := exec.Command("git", gitArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}

	return nil
}
