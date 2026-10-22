// Package desktop integrates with GitHub Desktop application.
package desktop

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// AddRepos registers discovered repositories with GitHub Desktop.
func AddRepos(records []model.ScanRecord) DesktopSummary {
	summary := DesktopSummary{}
	cli := ResolveCLI()
	if cli != "" {
		return addAll(records, summary, cli)
	}

	fmt.Println(constants.MsgDesktopNotFound)
	PrintInstallSuggestions()

	return summary
}

// addAll iterates records and adds each to GitHub Desktop.
func addAll(records []model.ScanRecord, summary DesktopSummary, cli string) DesktopSummary {
	for i, rec := range records {
		err := addOne(rec.AbsolutePath, cli)
		summary = updateSummary(summary, rec.RepoName, err)

		if i < len(records)-1 {
			time.Sleep(250 * time.Millisecond)
		}
	}

	return summary
}

// addOne opens a single repo in GitHub Desktop.
func addOne(repoPath, cli string) error {
	cmd := exec.Command(cli, constants.GitHubDesktopVerbOpen, repoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, "open repo in GitHub Desktop").WithContext("output", strings.TrimSpace(string(output)))
	}

	return nil
}

// updateSummary tracks success/failure for each repo.
func updateSummary(s DesktopSummary, name string, err error) DesktopSummary {
	if err == nil {
		s.Added++
		fmt.Printf(constants.MsgDesktopAdded, name)

		return s
	}

	s.Failed++
	fmt.Printf(constants.MsgDesktopFailed, name, err)

	return s
}

// DesktopSummary tracks GitHub Desktop registration results.
type DesktopSummary struct {
	Added  int
	Failed int
}
