package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func resolveTargetRepoAndRunId(args []string) (string, uint64) {
	repo := resolveCurrentRepoSlug()
	var runId uint64

	for _, a := range args {
		runId, repo = parseStageArg(a, runId, repo)
	}
	if runId == 0 {
		runId = fetchLatestRunIdForRepo(repo)
	}
	return repo, runId
}

func parseStageArg(arg string, currId uint64, currRepo string) (uint64, string) {
	if strings.HasPrefix(arg, "-") {
		return currId, currRepo
	}
	if id, err := strconv.ParseUint(arg, 10, 64); err == nil && id > 0 {
		return id, currRepo
	}
	if strings.Contains(arg, "/") {
		return currId, arg
	}
	return currId, currRepo
}

func fetchLatestRunIdForRepo(repo string) uint64 {
	runs := queryWorkflowRuns(repo)
	if len(runs) > 0 {
		return runs[0].DatabaseId
	}
	return 0
}

func renderStagesJSON(summary *pipelinedb.PipelineStageSummary) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
