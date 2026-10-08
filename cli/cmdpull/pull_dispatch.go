package cmdpull

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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
