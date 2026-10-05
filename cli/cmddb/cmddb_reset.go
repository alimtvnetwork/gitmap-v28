package cmddb

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ResetOptions holds configuration flags for database reset operations.
type ResetOptions struct {
	IsConfirm bool
	IsDryRun  bool
	IsRescan  bool
}

// ParseResetOptions parses command-line flags for database reset.
func ParseResetOptions(args []string) ResetOptions {
	fs := flag.NewFlagSet("reset", flag.ContinueOnError)
	var isConfirm bool
	var isDryRun bool
	var isRescan bool

	fs.BoolVar(&isConfirm, "confirm", false, "Confirm database deletion without prompt")
	fs.BoolVar(&isConfirm, "y", false, "Confirm database deletion without prompt")
	fs.BoolVar(&isConfirm, "yes", false, "Confirm database deletion without prompt")
	fs.BoolVar(&isDryRun, "dry-run", false, "Preview databases to delete without deleting")
	fs.BoolVar(&isDryRun, "n", false, "Preview databases to delete without deleting")
	fs.BoolVar(&isRescan, "rescan", false, "Automatically trigger scan after reseed")

	_ = fs.Parse(args)

	for _, a := range args {
		trimmed := strings.TrimSpace(a)
		if trimmed == "-y" || trimmed == "--yes" || trimmed == "--confirm" || trimmed == "-f" || trimmed == "--force" {
			isConfirm = true
		}
		if trimmed == "--dry-run" || trimmed == "-n" {
			isDryRun = true
		}
		if trimmed == "--rescan" {
			isRescan = true
		}
	}

	return ResetOptions{
		IsConfirm: isConfirm,
		IsDryRun:  isDryRun,
		IsRescan:  isRescan,
	}
}

func runDbResetAction(args []string) error {
	opts := ParseResetOptions(args)
	needsPrompt := !opts.IsDryRun && !opts.IsConfirm
	if needsPrompt && !confirmOrSkip(constants.ColorYellow+"Are you sure you want to reset all databases? All tracked repository records and split databases will be cleared. [y/N]: "+constants.ColorReset, args) {
		fmt.Println(constants.ColorDim + "Database reset canceled." + constants.ColorReset)

		return nil
	}

	return PerformComprehensiveReset(opts)
}

// PerformComprehensiveReset discovers, purges all database files, and reseeds clean baseline schemas.
func PerformComprehensiveReset(opts ResetOptions) error {
	targets := collectAllResetTargets()
	if len(targets) == 0 {
		fmt.Println("  No databases found to reset.")

		return nil
	}

	var purgedCount int
	var reclaimedBytes int64

	for _, dbPath := range targets {
		fi, err := os.Stat(dbPath)
		if err != nil {
			continue
		}
		totalSize := fi.Size()
		companions := []string{"-wal", "-shm", ".wal", ".shm"}
		var compSizes int64
		for _, ext := range companions {
			if cfi, cerr := os.Stat(dbPath + ext); cerr == nil {
				compSizes += cfi.Size()
			}
		}
		totalSize += compSizes

		if opts.IsDryRun {
			fmt.Printf("  [dry-run] Would remove database: %s (%s)\n", dbPath, formatBytes(totalSize))
			purgedCount++
			reclaimedBytes += totalSize

			continue
		}

		errRem := os.Remove(dbPath)
		for _, ext := range companions {
			_ = os.Remove(dbPath + ext)
		}
		parentDir := filepath.Dir(dbPath)
		_ = os.Remove(filepath.Join(parentDir, "gitmap.lock"))

		if errRem != nil && !os.IsNotExist(errRem) {
			fmt.Printf("  %s⚠ Could not remove database:%s %s (%v)\n", constants.ColorYellow, constants.ColorReset, dbPath, errRem)

			continue
		}

		fmt.Printf("  %s✔ Removed database:%s %s (%s)\n", constants.ColorGreen, constants.ColorReset, dbPath, formatBytes(totalSize))
		purgedCount++
		reclaimedBytes += totalSize
	}

	if opts.IsDryRun {
		fmt.Printf("\n  %s[dry-run] Total %d database(s) found (%s would be reclaimed).%s\n",
			constants.ColorCyan, purgedCount, formatBytes(reclaimedBytes), constants.ColorReset)

		return nil
	}

	fmt.Printf("\n  %s✔ Total %d database(s) purged (%s reclaimed).%s\n",
		constants.ColorGreen, purgedCount, formatBytes(reclaimedBytes), constants.ColorReset)

	db, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "E9001")
	}
	defer db.Close()

	if migErr := db.Migrate(); migErr != nil {
		return apperror.WrapSimple(migErr, "E9002")
	}
	_ = db.SeedProjectTypes()

	fmt.Printf("  %s✔ All databases reseeded with clean baseline schemas.%s\n",
		constants.ColorGreen, constants.ColorReset)

	return nil
}

func collectAllResetTargets() []string {
	entries := store.CollectAllDatabaseEntries()
	seen := make(map[string]bool)
	var targets []string

	addPath := func(p string) {
		if p == "" {
			return
		}
		clean := filepath.Clean(p)
		key := strings.ToLower(clean)
		if seen[key] {
			return
		}
		seen[key] = true
		if _, err := os.Stat(clean); err == nil {
			targets = append(targets, clean)
		}
	}

	for _, e := range entries {
		addPath(e.DatabasePath)
	}

	addPath(store.DefaultDBPath())
	if bin := store.BinaryDataDir(); bin != "" {
		addPath(filepath.Join(bin, constants.DBFile))
	}
	if global := store.GlobalUserDataDir(); global != "" {
		addPath(filepath.Join(global, constants.DBFile))
	}
	addPath(filepath.Join(".gitmap", constants.DBFile))
	addPath(filepath.Join(constants.DefaultOutputFolder, constants.DBDir, constants.DBFile))

	if searchMatches, err := filepath.Glob(filepath.Join(".gitmap", "output", "repo_search", "*.db")); err == nil {
		for _, m := range searchMatches {
			addPath(m)
		}
	}
	if cacheMatches, err := filepath.Glob(filepath.Join(".gitmap", "cache", "repos", "*", "*.db")); err == nil {
		for _, m := range cacheMatches {
			addPath(m)
		}
	}

	return targets
}

