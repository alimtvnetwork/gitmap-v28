package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runTaskRun starts the file-sync watch loop for a named task.
func runTaskRun(args []string) {
	fs := flag.NewFlagSet("task-run", flag.ExitOnError)

	var interval int
	var verbose, dryRun bool

	fs.IntVar(&interval, constants.FlagTaskInterval, constants.TaskDefaultInterval, constants.FlagDescTaskInterval)
	fs.BoolVar(&verbose, constants.FlagTaskVerbose, false, constants.FlagDescTaskVerbose)
	fs.BoolVar(&dryRun, constants.FlagTaskDryRun, false, constants.FlagDescTaskDryRun)
	fs.Parse(args)

	name := fs.Arg(0)
	if name == "" {
		fmt.Fprint(os.Stderr, constants.ErrTaskNameRequired)
		os.Exit(1)
	}

	interval = enforceMinInterval(interval)
	tasks := loadTaskFile()
	entry := findTaskByName(tasks, name)

	fmt.Printf(constants.MsgTaskRunning, name, interval)
	runSyncLoop(entry, interval, verbose, dryRun)
}

// enforceMinInterval clamps the interval to the minimum.
func enforceMinInterval(interval int) int {
	if interval < constants.TaskMinInterval {
		return constants.TaskMinInterval
	}

	return interval
}

// runSyncLoop runs the sync cycle on a timer until interrupted.
func runSyncLoop(entry model.TaskEntry, interval int, verbose, dryRun bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	syncOnce(entry, verbose, dryRun)

	for {
		select {
		case <-sigChan:
			fmt.Printf(constants.MsgTaskStopped, entry.Name)

			return
		case <-ticker.C:
			syncOnce(entry, verbose, dryRun)
		}
	}
}

// syncOnce performs a single sync pass from source to destination.
func syncOnce(entry model.TaskEntry, verbose, dryRun bool) {
	ignorePatterns := loadGitignorePatterns(entry.Source)
	syncCount := 0

	err := filepath.Walk(entry.Source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath := relativePath(entry.Source, path)
		if isIgnored(relPath, info.IsDir(), ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if info.IsDir() {
			return nil
		}

		synced := syncSingleFile(entry.Source, entry.Dest, relPath, info, dryRun, verbose)
		if synced {
			syncCount++
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTaskSyncFailed, entry.Source, err)
	}

	if syncCount == 0 && verbose {
		fmt.Print(constants.MsgTaskUpToDate)
	}
}

// syncSingleFile compares timestamps and copies if source is newer.
func syncSingleFile(srcRoot, destRoot, relPath string, srcInfo os.FileInfo, dryRun, verbose bool) bool {
	destPath := filepath.Join(destRoot, relPath)
	destInfo, err := os.Stat(destPath)

	isNew := err != nil
	isNewer := isNew || srcInfo.ModTime().After(destInfo.ModTime())
	if isNewer {
		srcPath := filepath.Join(srcRoot, relPath)

		return handleSyncCopy(srcPath, destPath, relPath, dryRun, verbose)
	}

	return false
}

// handleSyncCopy performs the copy or dry-run print.
func handleSyncCopy(srcPath, destPath, relPath string, dryRun, verbose bool) bool {
	if dryRun {
		fmt.Printf(constants.MsgTaskDrySync, relPath)

		return true
	}

	err := ensureDestDir(destPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTaskDestCreate, destPath, err)

		return false
	}

	err = copyFileContent(srcPath, destPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTaskSyncFailed, relPath, err)

		return false
	}

	if verbose {
		fmt.Printf(constants.MsgTaskSynced, relPath)
	}

	return true
}

// ensureDestDir creates parent directories for a destination file.
func ensureDestDir(destPath string) error {
	return os.MkdirAll(filepath.Dir(destPath), constants.DirPermission)
}

// relativePath returns a path relative to the base directory.
func relativePath(base, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}

	return rel
}
