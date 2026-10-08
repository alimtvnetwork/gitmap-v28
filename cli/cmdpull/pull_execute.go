package cmdpull

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

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
