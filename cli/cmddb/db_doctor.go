package cmddb

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const dbDoctorHelp = `db doctor — database health check across the split-db tiers.

Description:
  Reports sizes and staleness for gitmap.db, installation.db,
  pipeline.db, and the per-repo repodb/ caches. Flags anything not
  touched in 30+ days. With --repair it prints the planned prune
  actions (stub: prints only, makes no changes).

Usage:
  gitmap db doctor [--repair] [--dry-run]
  gitmap dbdoc [--repair]

Flags:
  --repair      show the planned repair/prune actions
  --dry-run     with --repair: preview only (default: true)
  --help, -h    show this help

Example:
  gitmap db doctor
`

// doctorItem is one database file or directory under inspection.
type doctorItem struct {
	name     string
	path     string
	size     int64
	isDir    bool
	isFound  bool
	modified time.Time
}

// staleAfter marks files untouched this long as stale.
const staleAfter = 30 * 24 * time.Hour

// RunDbDoctor implements `gitmap db doctor` (alias: dbdoc).
func RunDbDoctor(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(dbDoctorHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("db doctor", flag.ContinueOnError)
	repair := fs.Bool("repair", false, "show planned repair actions")
	dryRun := fs.Bool("dry-run", true, "preview repair actions without applying")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(dbDoctorHelp)
		return apperror.WrapSimple(parseErr, "db.doctor.flags")
	}

	items := collectDoctorItems()
	fmt.Printf("  %-18s %-10s %-8s %s\n", "DATABASE", "SIZE", "STATUS", "PATH")
	fmt.Printf("  %-18s %-10s %-8s %s\n", "--------", "----", "------", "----")
	staleCount := 0
	for _, it := range items {
		status := doctorStatus(it)
		if status == "STALE" {
			staleCount++
		}
		size := "-"
		if it.isFound {
			size = formatDoctorSize(it.size)
		}
		fmt.Printf("  %-18s %-10s %-8s %s\n", it.name, size, status, it.path)
	}
	fmt.Println()
	if staleCount > 0 {
		fmt.Printf("db doctor: %d item(s) stale (>30d untouched) — re-run with --repair to see prune actions.\n", staleCount)
	} else {
		fmt.Println("db doctor: all databases healthy.")
	}

	if *repair {
		printRepairPlan(items, *dryRun)
	}

	return nil
}

// collectDoctorItems gathers the known database tiers.
func collectDoctorItems() []doctorItem {
	var items []doctorItem
	if mainInfo, isFound := collectMainDBInfo(); isFound {
		items = append(items, doctorItemFromInfo("gitmap.db", mainInfo.Path, mainInfo.Size, false))
	} else {
		items = append(items, doctorItem{name: "gitmap.db", path: "not registered", isFound: false})
	}
	binDir := store.BinaryDataDir()
	for _, name := range []string{"installation.db", "pipeline.db"} {
		full := filepath.Join(binDir, name)
		items = append(items, statDoctorItem(name, full, false))
	}
	for _, split := range collectSplitDBs() {
		items = append(items, doctorItemFromInfo(split.Name, split.Path, split.Size, false))
	}
	for _, dir := range findSplitDbDirs() {
		items = append(items, statDoctorItem("repodb/", dir, true))
	}

	return items
}

// doctorItemFromInfo builds an item from an already-known DBFileInfo.
func doctorItemFromInfo(name, path string, size int64, isDir bool) doctorItem {
	info, statErr := os.Stat(path)
	item := doctorItem{name: name, path: path, size: size, isDir: isDir, isFound: statErr == nil}
	if statErr == nil {
		item.modified = info.ModTime()
	}

	return item
}

// statDoctorItem stats a path that may or may not exist.
func statDoctorItem(name, path string, isDir bool) doctorItem {
	info, statErr := os.Stat(path)
	if statErr != nil {
		return doctorItem{name: name, path: path, isDir: isDir, isFound: false}
	}
	size := info.Size()
	if isDir {
		size = dirDoctorSize(path)
	}

	return doctorItem{name: name, path: path, size: size, isDir: isDir, isFound: true, modified: info.ModTime()}
}

// dirDoctorSize sums file sizes one level deep (bounded, no recursion storms).
func dirDoctorSize(dir string) int64 {
	var total int64
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return 0
	}
	for _, e := range entries {
		info, infoErr := e.Info()
		if infoErr != nil {
			continue
		}
		total += info.Size()
	}

	return total
}

// doctorStatus classifies an item as OK, STALE, or MISSING.
func doctorStatus(it doctorItem) string {
	if !it.isFound {
		return "MISSING"
	}
	if time.Since(it.modified) > staleAfter {
		return "STALE"
	}

	return "OK"
}

// formatDoctorSize renders bytes as a short human size.
func formatDoctorSize(size int64) string {
	const mb = 1024 * 1024
	const kb = 1024
	switch {
	case size >= mb:
		return fmt.Sprintf("%.1f MB", float64(size)/float64(mb))
	case size >= kb:
		return fmt.Sprintf("%.1f KB", float64(size)/float64(kb))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// printRepairPlan prints the stub prune plan; it never mutates anything.
func printRepairPlan(items []doctorItem, isDryRun bool) {
	fmt.Println()
	fmt.Println("repair plan (stub — no changes are made by this command):")
	for _, it := range items {
		if !it.isFound {
			continue
		}
		status := doctorStatus(it)
		if status == "STALE" || it.isDir {
			fmt.Printf("  - prune %s (%s, %s)\n", it.path, formatDoctorSize(it.size), strings.ToLower(status))
		} else if !it.isDir {
			fmt.Printf("  - vacuum %s (%s)\n", it.path, formatDoctorSize(it.size))
		}
	}
	fmt.Println("  - remove orphaned -wal/-shm journals next to repaired DBs")
	if isDryRun {
		fmt.Println("dry-run: listed actions only, nothing was executed.")
	}
}
