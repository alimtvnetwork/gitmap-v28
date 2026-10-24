package cmddb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/cli/repodb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func optimizeMasterDB() (int64, int) {
	masterPath := store.DefaultDBPath()
	if _, err := os.Stat(masterPath); err == nil {
		return tryOptimizeDBFile(masterPath, "Master Database")
	}

	return 0, 0
}

func optimizeSplitRepoDBs() (int64, int) {
	var totalReclaimed int64
	var optimizedCount int
	for _, s := range collectSplitDBs() {
		rec, cnt := tryOptimizeDBFile(s.Path, "Split Repo DB")
		totalReclaimed += rec
		optimizedCount += cnt
	}

	return totalReclaimed, optimizedCount
}

func optimizePipelineDBs() (int64, int) {
	var totalReclaimed int64
	var optimizedCount int
	pipeDir := pipelinedb.PipelineDBDir()
	pipeEntries, err := os.ReadDir(pipeDir)
	if err != nil {
		return 0, 0
	}

	for _, p := range pipeEntries {
		rec, cnt := tryOptimizePipelineFile(pipeDir, p)
		totalReclaimed += rec
		optimizedCount += cnt
	}

	return totalReclaimed, optimizedCount
}

func runDBOptimize(args []string) error {
	fmt.Printf("%s● Optimizing all Gitmap SQLite databases (VACUUM & PRAGMA optimize)...%s\n\n",
		constants.ColorCyan, constants.ColorReset)

	mRec, mCnt := optimizeMasterDB()
	sRec, sCnt := optimizeSplitRepoDBs()
	pRec, pCnt := optimizePipelineDBs()

	totalReclaimed := mRec + sRec + pRec
	optimizedCount := mCnt + sCnt + pCnt

	fmt.Println()
	fmt.Printf("%s✓ Optimization complete across %d database file(s). Total disk space reclaimed: %s%s\n",
		constants.ColorGreen, optimizedCount, formatBytes(totalReclaimed), constants.ColorReset)

	return nil
}

func tryOptimizePipelineFile(dir string, p os.DirEntry) (int64, int) {
	if p.IsDir() || filepath.Ext(p.Name()) != ".db" {
		return 0, 0
	}

	pipePath := filepath.Join(dir, p.Name())

	return tryOptimizeDBFile(pipePath, "Split Pipeline DB")
}

func tryOptimizeDBFile(path, label string) (int64, int) {
	reclaimed, err := optimizeSingleDBFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Failed to optimize %s (%s): %v\n", label, filepath.Base(path), err)

		return 0, 0
	}

	fmt.Printf("  ✔ %-20s %-40s (reclaimed: %s)\n", label+":", filepath.Base(path), formatBytes(reclaimed))

	return reclaimed, 1
}

func optimizeSingleDBFile(path string) (int64, error) {
	openRes := store.OpenSQLiteDB(path)
	if openRes.IsFailure() {
		return 0, openRes.AppError()
	}
	db := openRes.Value

	defer db.Close()

	return repodb.OptimizeRepoDB(context.Background(), db, path)
}
