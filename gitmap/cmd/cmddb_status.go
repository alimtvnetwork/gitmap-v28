package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type dbUnifiedStatus struct {
	MasterPath     string `json:"masterPath"`
	MasterSize     int64  `json:"masterSize"`
	MasterRepos    int    `json:"masterRepos"`
	SplitRepoDir   string `json:"splitRepoDir"`
	SplitRepoCount int    `json:"splitRepoCount"`
	SplitRepoSize  int64  `json:"splitRepoSize"`
	PipelineDir    string `json:"pipelineDir"`
	PipelineCount  int    `json:"pipelineCount"`
	PipelineSize   int64  `json:"pipelineSize"`
}

func runDBStatus(args []string) error {
	var status dbUnifiedStatus
	status.MasterPath = store.DefaultDBPath()
	if info, err := os.Stat(status.MasterPath); err == nil {
		status.MasterSize = info.Size()
	}

	mainDB, err := store.OpenDefault()
	if err == nil {
		defer mainDB.Close()
		repos, _ := mainDB.ListRepos()
		status.MasterRepos = len(repos)
	}

	status.SplitRepoDir = filepath.Join(store.BinaryDataDir(), "repo_search")
	repoDBs := collectSplitDBs()
	status.SplitRepoCount = len(repoDBs)
	for _, r := range repoDBs {
		status.SplitRepoSize += r.Size
	}

	status.PipelineDir = pipelinedb.PipelineDBDir()
	pipeFiles, _ := os.ReadDir(status.PipelineDir)
	for _, f := range pipeFiles {
		accumulatePipelineFileStats(&status, f)
	}

	if hasArgFlag(args, "--json") {
		return printJSON(status)
	}

	printUnifiedStatus(status)
	return nil
}

func printUnifiedStatus(status dbUnifiedStatus) {
	fmt.Println(constants.ColorCyan + "● Gitmap SQLite Unified Database Status:" + constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %s1. Primary Master Database (gitmap.db):%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("     • %-18s %s\n", "Location:", filepath.Dir(status.MasterPath))
	fmt.Printf("     • %-18s %s\n", "Full Path:", status.MasterPath)
	fmt.Printf("     • %-18s %s\n", "File Size:", formatBytes(status.MasterSize))
	fmt.Printf("     • %-18s %d registered repository records\n", "Master Records:", status.MasterRepos)
	fmt.Println()
	fmt.Printf("  %s2. Split Repository Databases (repo_search/):%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("     • %-18s %s\n", "Directory:", status.SplitRepoDir)
	fmt.Printf("     • %-18s %d isolated database files\n", "Database Count:", status.SplitRepoCount)
	fmt.Printf("     • %-18s %s\n", "Total Disk Size:", formatBytes(status.SplitRepoSize))
	fmt.Println()
	fmt.Printf("  %s3. Split Pipeline Databases (pipeline_db/):%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("     • %-18s %s\n", "Directory:", status.PipelineDir)
	fmt.Printf("     • %-18s %d isolated pipeline database files\n", "Database Count:", status.PipelineCount)
	fmt.Printf("     • %-18s %s\n", "Total Disk Size:", formatBytes(status.PipelineSize))
	fmt.Println()
	fmt.Printf("  %sTip: Run 'gitmap db optimize' to vacuum and reclaim disk space across all databases.%s\n",
		constants.ColorDim, constants.ColorReset)
}

func accumulatePipelineFileStats(status *dbUnifiedStatus, f os.DirEntry) {
	if f.IsDir() || filepath.Ext(f.Name()) != ".db" {
		return
	}
	status.PipelineCount++
	info, err := f.Info()
	if err == nil {
		status.PipelineSize += info.Size()
	}
}
