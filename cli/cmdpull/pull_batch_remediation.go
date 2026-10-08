package cmdpull

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

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
		RepoPath:      filepath.ToSlash(filepath.Clean(rec.AbsolutePath)),
		SummaryReason: diag.SummaryReason,
		Recipes:       recipes,
		Files:         diag.AllFiles,
	}
}

func finalizePullBatchTask(taskDB *store.DB, taskID int64, failCount int, failures []PullFailureSummary) error {
	if failCount > 0 {
		errMsg := fmt.Sprintf("pull batch finished with %d failure(s)", failCount)
		failPendingTask(taskDB, taskID, errMsg)
		fmt.Fprintf(os.Stderr, "\n  %s%s%s\n\n", constants.ColorRed, errMsg, constants.ColorReset)
		handleBatchFailuresRemediation(failures)
		cliexit.Exit(1)

		return nil
	}
	completePendingTask(taskDB, taskID)

	return nil
}

func handleBatchFailuresRemediation(failures []PullFailureSummary) {
	hasFailures := len(failures) > 0
	if hasFailures {
		_ = PromptInteractiveBatchFix(failures)
	}
}
