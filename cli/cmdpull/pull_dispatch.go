package cmdpull

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"os"
	"strings"
)

// runPull handles the "pull" subcommand.
func runPull(args []string) error {
	args = NormalizePullArgs(args)
	isPullAll := isPullAllInvocation(args)
	checkPullHelp(isPullAll, args)
	if handled, err := checkEfficientSubcommand(args); handled {
		return err
	}
	if isPullAll && hasSSHFleetFlag(args) {
		return handleSSHFleetPullAll(args)
	}

	return runPullStandardFlow(args, isPullAll)
}

func runPullStandardFlow(args []string, isPullAll bool) error {
	args = NormalizePullArgs(args)
	printHeaderUnlessJSON(isPullAll, args)
	requireOnline()
	if db, err := openDB(); err == nil {
		_, _ = OptimizeRedundantRepos(db, hasJSONArg(args))
		db.Close()
	}
	useSSH, useHTTPS, restArgs := ExtractTransportFlags(args)
	if isCWDTransportEligible(isPullAll, useSSH, useHTTPS) {
		return runPullCWDWithTransport(useSSH, useHTTPS, restArgs)
	}
	opts := resolveParsedPullOptions(restArgs, useSSH, useHTTPS)

	return executePullWithResolvedOptions(opts)
}

func printHeaderUnlessJSON(isPullAll bool, args []string) {
	if !hasJSONArg(args) {
		printPullInvocationHeader(isPullAll)
	}
}

func isCWDTransportEligible(isPullAll, useSSH, useHTTPS bool) bool {
	if isPullAll {
		return false
	}
	hasTransport := useSSH || useHTTPS
	if !hasTransport {
		return false
	}

	return isGitRepoCWD()
}

func checkPullHelp(isPullAll bool, args []string) {
	helpCmd := "pull"
	if isPullAll {
		helpCmd = "pull-all"
	}
	checkHelp(helpCmd, args)
}

func handleSSHFleetPullAll(args []string) error {
	cleanArgs := stripSSHFleetFlags(args)
	hasHook := RunRemoteSSHPullAllFleetFn != nil
	if !hasHook {
		return apperror.NewSimple("ssh fleet pull-all is not wired", "E_SSH_PULLALL")
	}

	return executeSSHFleetPullTask(cleanArgs)
}

func executeSSHFleetPullTask(cleanArgs []string) error {
	queueId, tDB := enqueueTaskQueue("pull-all-ssh", "fleet")
	defer closeTaskDB(tDB)
	updateTaskQueue(tDB, queueId, "running")

	errFleet := RunRemoteSSHPullAllFleetFn(cleanArgs)
	hasErr := errFleet != nil
	if hasErr {
		return handleSSHFleetFailure(tDB, queueId, errFleet)
	}

	updateTaskQueue(tDB, queueId, "completed")
	return nil
}

func handleSSHFleetFailure(tDB *store.TasksSplitDB, queueId string, errFleet error) error {
	store.LogInternalError("PULL_ALL_SSH", "SSH_DELEGATION_ERROR", errFleet.Error(), "", "")
	updateTaskQueue(tDB, queueId, "failed")

	return apperror.WrapSimple(errFleet, "SSH fleet pull-all failed")
}

func hasSSHFleetFlag(args []string) bool {
	for _, a := range args {
		if isSSHFleetToken(a) {
			return true
		}
	}
	return false
}

func isSSHFleetToken(token string) bool {
	low := strings.ToLower(token)
	return low == "--ssh" || low == "-ssh" || low == "--sh" || low == "-sh" || low == "ssh" || low == "pas" || low == "pull-all-ssh"
}

func stripSSHFleetFlags(args []string) []string {
	var clean []string
	for _, a := range args {
		if isSSHFleetToken(a) {
			continue
		}
		clean = append(clean, a)
	}
	return clean
}

func executePullWithResolvedOptions(opts pullOptions) error {
	if opts.isProbe {
		executePullProbeHook(opts)
	}
	if opts.verbose {
		initVerboseLog()
	}

	return dispatchPullExecution(opts)
}

func executePullProbeHook(opts pullOptions) {
	if RunSpecialRepoProbeOnPullFn == nil {
		return
	}
	if !opts.isJSON {
		fmt.Printf("\n  %s●%s %sProbing companion repositories (repo-secrets & repo-cache)...%s\n",
			constants.ColorCyan, constants.ColorReset, constants.ColorBold, constants.ColorReset)
	}
	_ = RunSpecialRepoProbeOnPullFn("", opts.yes)
}

func hasJSONArg(args []string) bool {
	for _, a := range args {
		if strings.EqualFold(a, "--json") || strings.EqualFold(a, "-json") {
			return true
		}
	}

	return false
}

func resolveParsedPullOptions(args []string, useSSH, useHTTPS bool) pullOptions {
	opts := parsePullFlags(args)
	opts = applyTransportOptions(opts, useSSH, useHTTPS)
	if isPullAllTableRootCmd() {
		opts.all = true
		opts.showStatus = true
	}

	return opts
}

func applyTransportOptions(opts pullOptions, useSSH, useHTTPS bool) pullOptions {
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
	if opts.parallel < 0 || opts.workers < 0 || opts.hands < 0 {
		cliexit.HandleError(apperror.NewSimple("invalid concurrency", "E9000"), 1)
	}
	opts = clampPullConcurrencyForSSH(opts)
	if opts.parallel > 1 {
		runPullParallel(records, bar, opts.parallel)
		return
	}
	runSerialPull(records, bar)
}

func clampPullConcurrencyForSSH(opts pullOptions) pullOptions {
	isSSH := cloneconcurrency.IsSSHSession() || opts.useSSH
	if isSSH {
		opts.workers = 1
		opts.hands = 1
		opts.parallel = 1
		return opts
	}
	effectiveWorkers := resolveEffectiveWorkers(opts.workers, opts.parallel)
	workers, hands := cloneconcurrency.ResolveWorkerHands(effectiveWorkers, opts.hands, opts.isWWOH, false)
	opts.workers = workers
	opts.hands = hands
	opts.parallel = workers * hands
	return opts
}

func resolveEffectiveWorkers(workers, parallel int) int {
	if workers <= 0 {
		return parallel
	}
	return workers
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
	if len(remItems) == 0 || isConcisePullOutput(opts.all, opts.showStatus) {
		return
	}
	dispatchRemediationSummary(remItems, opts)
}

func dispatchRemediationSummary(remItems []RemediationItem, opts pullOptions) {
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
	sshFlag, httpsFlag, statusFlag, jsonFlag, probeFlag            *bool
	wwohFlag                                                       *bool
	autoScaleFlag, highPerfFlag, lowCPUFlag                        *bool
	gFlag                                                          *string
	pFlag, wFlag, handFlag                                         *int
}
