package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/pipelinedb"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/repodb"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runDBOptimize(args []string) error {
	fmt.Printf("%s● Optimizing all Gitmap SQLite databases (VACUUM & PRAGMA optimize)...%s\n\n",
		constants.ColorCyan, constants.ColorReset)

	var totalReclaimed int64
	var optimizedCount int

	masterPath := store.DefaultDBPath()
	if _, err := os.Stat(masterPath); err == nil {
		rec, cnt := tryOptimizeDBFile(masterPath, "Master Database")
		totalReclaimed += rec
		optimizedCount += cnt
	}

	splitDBs := collectSplitDBs()
	for _, s := range splitDBs {
		rec, cnt := tryOptimizeDBFile(s.Path, "Split Repo DB")
		totalReclaimed += rec
		optimizedCount += cnt
	}

	pipeDir := pipelinedb.PipelineDBDir()
	pipeEntries, _ := os.ReadDir(pipeDir)
	for _, p := range pipeEntries {
		rec, cnt := tryOptimizePipelineFile(pipeDir, p)
		totalReclaimed += rec
		optimizedCount += cnt
	}

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
		return 0, 0
	}
	fmt.Printf("  ✔ %-20s %-40s (reclaimed: %s)\n", label+":", filepath.Base(path), formatBytes(reclaimed))
	return reclaimed, 1
}

func optimizeSingleDBFile(path string) (int64, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	return repodb.OptimizeRepoDB(context.Background(), db, path)
}
