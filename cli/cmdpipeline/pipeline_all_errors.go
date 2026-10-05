package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AllPipelineSummary aggregates pipeline error statuses across all repositories.
type AllPipelineSummary struct {
	TotalPipelines int                     `json:"totalPipelines"`
	CleanCount     int                     `json:"cleanCount"`
	FailedCount    int                     `json:"failedCount"`
	FailedItems    []AllPipelineFailedItem `json:"failedItems,omitempty"`
	CleanRepos     []string                `json:"cleanRepos,omitempty"`
	PullErrors     []store.PullErrorRecord `json:"pullErrors,omitempty"`
}

// AllPipelineFailedItem represents a failed pipeline run discovered across repositories.
type AllPipelineFailedItem struct {
	RepoSlug     string `json:"repoSlug"`
	WorkflowName string `json:"workflowName"`
	RunId        uint64 `json:"runId"`
	Branch       string `json:"branch"`
	Sha          string `json:"sha"`
	ErrorSummary string `json:"errorSummary"`
	UpdatedAt    string `json:"updatedAt"`
}

func executeAllPipelineErrorLogs(flags PipelineErrorFlags, args []string) error {
	summary := collectAllPipelineSummary()

	if flags.IsJSON {
		return emitAllPipelineSummaryJSON(summary)
	}

	renderAllPipelineSummaryTerminal(summary)

	return nil
}

func emitAllPipelineSummaryJSON(summary AllPipelineSummary) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(summary)
}

func collectAllPipelineSummary() AllPipelineSummary {
	slugs := discoverPipelineRepoSlugs()
	summary := AllPipelineSummary{TotalPipelines: len(slugs)}

	for _, slug := range slugs {
		failedItem, isClean := inspectSinglePipelineRepo(slug)

		if isClean {
			summary.CleanCount++
			summary.CleanRepos = append(summary.CleanRepos, slug)
			continue
		}

		if failedItem != nil {
			summary.FailedCount++
			summary.FailedItems = append(summary.FailedItems, *failedItem)
		}
	}

	summary.PullErrors = queryPullErrorsForAll()

	return summary
}

func discoverPipelineRepoSlugs() []string {
	var slugs []string
	pipelineDir := filepath.Join(store.BinaryDataDir(), "pipeline")
	entries, err := os.ReadDir(pipelineDir)

	if err != nil {
		return slugs
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		dbFile := filepath.Join(pipelineDir, e.Name(), "sql.db")

		if isFileExisting(dbFile) {
			slugs = append(slugs, e.Name())
		}
	}

	return slugs
}

func inspectSinglePipelineRepo(slug string) (*AllPipelineFailedItem, bool) {
	db, err := pipelinedb.OpenPipelineSplitDb(slug)

	if err != nil {
		return nil, false
	}
	defer db.Close()

	run, runErr := db.QueryRunByNegativeOffset(-1)

	if runErr != nil || run == nil {
		return nil, false
	}

	if run.IsSuccess || strings.EqualFold(run.Conclusion, "success") {
		return nil, true
	}

	errSummary := extractRunErrorSummary(db, run.RunId)

	return &AllPipelineFailedItem{
		RepoSlug:     slug,
		WorkflowName: run.WorkflowName,
		RunId:        run.RunId,
		Branch:       run.Branch,
		Sha:          run.Sha,
		ErrorSummary: errSummary,
		UpdatedAt:    run.UpdatedAt,
	}, false
}

func extractRunErrorSummary(db *pipelinedb.PipelineSplitDb, runId uint64) string {
	compactRes := db.QueryCompactErrorLogsByRunId(runId)

	if compactRes.IsSuccess() && len(compactRes.Data) > 0 {
		return compactRes.Data[0].ErrorText
	}

	return "Pipeline workflow run failed"
}

func queryPullErrorsForAll() []store.PullErrorRecord {
	pullDb, err := store.OpenPullSplitDB()

	if err != nil {
		return nil
	}
	defer pullDb.Close()

	records, _ := pullDb.QueryLatestPullErrors("all", 10)

	return records
}

func renderAllPipelineSummaryTerminal(s AllPipelineSummary) {
	fmt.Printf("\n  %s● Pipeline Error Inspector: ALL REPOSITORIES%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    Total Pipelines Tracked: %d\n", s.TotalPipelines)
	fmt.Printf("    Passing / Clean:         %d\n", s.CleanCount)
	fmt.Printf("    Failing:                 %d\n\n", s.FailedCount)

	if s.FailedCount == 0 {
		fmt.Printf("  %s✔ All tracked CI/CD pipelines across all repositories are PASSING (100%% green).%s\n\n",
			constants.ColorGreen, constants.ColorReset)
	}

	if s.FailedCount > 0 {
		renderFailingPipelinesList(s.FailedItems)
	}

	if len(s.PullErrors) > 0 {
		renderPullErrorsNotice(len(s.PullErrors))
	}
}

func renderFailingPipelinesList(items []AllPipelineFailedItem) {
	fmt.Printf("  %sFailing Repositories (%d):%s\n", constants.ColorRed, len(items), constants.ColorReset)

	for _, item := range items {
		fmt.Printf("    • %s%s%s (Run #%d):\n", constants.ColorYellow, item.RepoSlug, constants.ColorReset, item.RunId)
		fmt.Printf("      └── ✖ %s\n", item.WorkflowName)
		fmt.Printf("          Error: %s\n", item.ErrorSummary)
		fmt.Printf("          ↳ Inspect: gitmap pe %s\n\n", item.RepoSlug)
	}
}

func renderPullErrorsNotice(count int) {
	fmt.Printf("  %s● Recorded Pull Errors (%d failure(s)):%s\n", constants.ColorYellow, count, constants.ColorReset)
	fmt.Printf("    ↳ To inspect pull stack traces: gitmap pull-error all (or: gitmap pulle all)\n\n")
}
