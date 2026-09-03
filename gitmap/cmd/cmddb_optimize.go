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

	// 1. Optimize Primary Master DB
	masterPath := store.DefaultDBPath()
	if _, err := os.Stat(masterPath); err == nil {
		reclaimed, optErr := optimizeSingleDBFile(masterPath)
		if optErr == nil {
			totalReclaimed += reclaimed
			optimizedCount++
			fmt.Printf("  ✔ Master Database:     %-40s (reclaimed: %s)\n", filepath.Base(masterPath), formatBytes(reclaimed))
		}
	}

	// 2. Optimize Split Repo DBs
	splitDBs := collectSplitDBs()
	for _, s := range splitDBs {
		reclaimed, optErr := optimizeSingleDBFile(s.Path)
		if optErr == nil {
			totalReclaimed += reclaimed
			optimizedCount++
			fmt.Printf("  ✔ Split Repo DB:       %-40s (reclaimed: %s)\n", s.Name, formatBytes(reclaimed))
		}
	}

	// 3. Optimize Split Pipeline DBs
	pipeDir := pipelinedb.PipelineDBDir()
	pipeEntries, _ := os.ReadDir(pipeDir)
	for _, p := range pipeEntries {
		if !p.IsDir() && filepath.Ext(p.Name()) == ".db" {
			pipePath := filepath.Join(pipeDir, p.Name())
			reclaimed, optErr := optimizeSingleDBFile(pipePath)
			if optErr == nil {
				totalReclaimed += reclaimed
				optimizedCount++
				fmt.Printf("  ✔ Split Pipeline DB:   %-40s (reclaimed: %s)\n", p.Name(), formatBytes(reclaimed))
			}
		}
	}

	fmt.Println()
	fmt.Printf("%s✓ Optimization complete across %d database file(s). Total disk space reclaimed: %s%s\n",
		constants.ColorGreen, optimizedCount, formatBytes(totalReclaimed), constants.ColorReset)
	return nil
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
