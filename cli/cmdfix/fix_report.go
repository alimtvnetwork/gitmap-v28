package cmdfix

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// Per-repository result reporting for parallel fix runs, split out of
// fix_cmd.go for the 300-line file-size gate.

func reportFixBatchResult(mu *sync.Mutex, completed *int, total int, res FixStepResult, failedRepos *[]string) {
	mu.Lock()
	defer mu.Unlock()
	*completed++
	branch := "├──"
	if *completed == total {
		branch = "└──"
	}
	repoCol := formatRepoColumn(res.RepoName)
	durStr := fmt.Sprintf("%.1fs", res.Duration.Seconds())
	if res.Err != nil {
		*failedRepos = append(*failedRepos, res.RepoName)
		fmt.Printf("  %s [%d/%d] %s  %s %s (%s)\n",
			branch, *completed, total, repoCol,
			constants.ColorRed+"✗"+constants.ColorReset,
			shortenErrMsg(res.Err.Error()), durStr)
		return
	}
	if res.IsClean {
		fmt.Printf("  %s [%d/%d] %s  %s already clean (%s)\n",
			branch, *completed, total, repoCol,
			constants.ColorYellow+"•"+constants.ColorReset, durStr)
		return
	}
	fmt.Printf("  %s [%d/%d] %s  %s %s (%s)\n",
		branch, *completed, total, repoCol,
		constants.ColorGreen+"✔"+constants.ColorReset,
		res.ActionSummary, durStr)
}

func formatRepoColumn(name string) string {
	if len(name) >= 26 {
		return name
	}

	return fmt.Sprintf("%-26s", name)
}

func shortenErrMsg(msg string) string {
	lines := strings.Split(strings.TrimSpace(msg), "\n")
	first := lines[0]
	if len(first) > 50 {
		return first[:47] + "..."
	}

	return first
}
