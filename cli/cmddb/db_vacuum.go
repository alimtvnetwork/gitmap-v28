package cmddb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/repodb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func executeVacuumOnDB(path string) (int64, error) {
	db, err := store.OpenSQLiteDB(path)
	if err != nil {
		return 0, err
	}

	defer db.Close()

	return repodb.OptimizeRepoDB(context.Background(), db, path)
}

func recordVacuumTelemetry(path string, reclaimed int64, label string) {
	baseName := filepath.Base(path)
	formatted := formatBytes(reclaimed)
	fmt.Printf("  ✔ %-20s %-40s (reclaimed: %s)\n", label+":", baseName, formatted)
}

func handleVacuumError(path string, label string, err error) error {
	baseName := filepath.Base(path)
	fmt.Fprintf(os.Stderr, "  ⚠ Failed to vacuum %s (%s): %v\n", label, baseName, err)

	return err
}

func vacuumSingleTarget(path, label string) int64 {
	reclaimed, err := executeVacuumOnDB(path)
	if err != nil {
		_ = handleVacuumError(path, label, err)

		return 0
	}

	recordVacuumTelemetry(path, reclaimed, label)

	return reclaimed
}

// runDBVacuum performs SQLite VACUUM and PRAGMA optimize across databases.
func runDBVacuum(args []string) error {
	return runDBOptimize(args)
}
