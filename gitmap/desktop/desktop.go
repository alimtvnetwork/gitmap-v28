// Package desktop integrates with GitHub Desktop application.
package desktop

import (
	"fmt"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// AddRepos registers discovered repositories with GitHub Desktop.
func AddRepos(records []model.ScanRecord) DesktopSummary {
	summary := DesktopSummary{}
	cli := ResolveCLI()
	if cli != "" {
		return addAll(records, summary, cli)
	}
	fmt.Println(constants.MsgDesktopNotFound)

	return summary
}

// addAll iterates records and adds each to GitHub Desktop.
func addAll(records []model.ScanRecord, summary DesktopSummary, cli string) DesktopSummary {
	for _, rec := range records {
		err := addOne(rec.AbsolutePath, cli)
		summary = updateSummary(summary, rec.RepoName, err)
	}

	return summary
}

// addOne opens a single repo in GitHub Desktop.
func addOne(repoPath, cli string) error {
	cmd := exec.Command(cli, repoPath)
	_, err := cmd.Output()

	return err
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
