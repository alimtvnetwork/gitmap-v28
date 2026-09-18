package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func checkAgyFixDuplicate(opts AgyFixOptions, payload cmdpipeline.PipelineErrorLogsPayload, errorReport string) (bool, string, string, string) {
	sig, errHash := ComputeErrorSignature(payload.Repo, payload.RunId, payload.Sha, errorReport)
	storePath := sentAgyErrorsStorePath()
	store := LoadSentAgyErrorsStore(storePath)
	isDuplicate, existingRecord := CheckSentErrorDuplicate(sig, store, opts.IsForce)
	if isDuplicate {
		renderDuplicateNotice(payload.Repo, payload.RunId, payload.Sha, existingRecord)

		return true, storePath, sig, errHash
	}

	return false, storePath, sig, errHash
}

func renderDuplicateNotice(repo string, runID uint64, sha string, rec *SentAgyErrorRecord) {
	desc := formatErrorRunDesc(repo, runID, sha)
	fmt.Printf("\n  %s⚠ Pipeline errors for %s have already been sent to Antigravity!%s\n",
		constants.ColorYellow, desc, constants.ColorReset)
	if rec != nil && len(rec.SentAt) > 0 {
		fmt.Printf("    Previously sent at: %s (dispatch count: %d)\n", rec.SentAt, rec.SentCount)
	}
	fmt.Printf("    %sDo you want to send again? Use --force or -f to send again.%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func formatErrorRunDesc(repo string, runID uint64, sha string) string {
	if runID > 0 && len(sha) > 0 {
		return fmt.Sprintf("run #%d (commit %s)", runID, truncateSHA(sha))
	}
	if runID > 0 {
		return fmt.Sprintf("run #%d", runID)
	}
	if len(sha) > 0 {
		return fmt.Sprintf("commit %s", truncateSHA(sha))
	}

	return formatDisplayRepo(repo)
}

func truncateSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}

	return sha
}
