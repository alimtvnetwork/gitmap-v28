package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func executeBatchSlice(batch []ProjectFixCandidate, opts AgyFixOptions) {
	for i, cand := range batch {
		fmt.Printf("\n  [%d/%d] Processing %s (%s)...\n", i+1, len(batch), cand.Name, cand.RepoSlug)
		candOpts := opts
		candOpts.Repo = cand.RepoSlug
		_ = dispatchCandidateFix(cand, candOpts)
	}
}

func handleBatchResetIfRequested(opts AgyFixOptions, cursorPath string) {
	if opts.IsResetBatch {
		_ = ResetPipelineFixBatchCursor(cursorPath)
	}
}

func checkBatchComplete(cursor PipelineFixBatchCursor, total int) bool {
	if cursor.LastIndex < total {
		return false
	}
	fmt.Printf("\n  %sℹ All %d failing projects have already been processed in previous batches.%s\n",
		constants.ColorYellow, total, constants.ColorReset)
	fmt.Printf("    Run with %s--reset-batch%s to restart from the beginning.\n\n",
		constants.ColorCyan, constants.ColorReset)

	return true
}

func processBatchRange(failing []ProjectFixCandidate, cursor *PipelineFixBatchCursor, cursorPath string, limit int, opts AgyFixOptions) {
	end := cursor.LastIndex + limit
	if end > len(failing) {
		end = len(failing)
	}

	printBatchHeader(cursor.LastIndex, end, len(failing))
	executeBatchSlice(failing[cursor.LastIndex:end], opts)
	cursor.LastIndex = end
	_ = SavePipelineFixBatchCursor(cursorPath, *cursor)
	printBatchProgress(end, len(failing))
}

func handleNoFailingCandidates() bool {
	fmt.Printf("\n  %s✔ All repositories have passing pipelines (100%% green)!%s\n\n",
		constants.ColorGreen, constants.ColorReset)

	return true
}

func dispatchBatchExecution(failing []ProjectFixCandidate, cursorPath string, opts AgyFixOptions) {
	limit := resolveBatchLimit(opts)
	cursor := LoadPipelineFixBatchCursor(cursorPath, limit)
	if checkBatchComplete(cursor, len(failing)) {
		return
	}
	processBatchRange(failing, &cursor, cursorPath, limit, opts)
}

// RunAgyFixMultiProjectBatch executes multi-project batch scanning and AGY injection.
func RunAgyFixMultiProjectBatch(opts AgyFixOptions) error {
	cursorPath := resolveBatchCursorPath()
	handleBatchResetIfRequested(opts, cursorPath)
	failing := ScanFailingProjectCandidates(opts.IsDetailed)
	if len(failing) == 0 && handleNoFailingCandidates() {
		return nil
	}
	dispatchBatchExecution(failing, cursorPath, opts)

	return nil
}
