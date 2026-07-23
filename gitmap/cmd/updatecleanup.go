package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// runUpdateCleanup handles the "update-cleanup" subcommand.
// Removes leftover temp binaries and .old backup files.
func runUpdateCleanup() {
	dumpDebugWindowsHeader("update-cleanup (deployed binary)")
	defer dumpDebugWindowsFooter()

	selfPath := resolveCleanupSelfPath()
	fmt.Println(constants.MsgUpdateCleanStart)
	if len(selfPath) > 0 {
		fmt.Printf(constants.MsgUpdateCleanBinary, selfPath)
	}
	logUpdateCleanup(constants.UpdateCleanupLogStart, selfPath)
	logHandoffEvent("cleanup", "start",
		map[string]string{"self": selfPath})
	dumpDebugWindowsNote("cleanup self path resolved: %s", selfPath)
	delayUpdateCleanupIfNeeded()

	ctx := loadUpdateCleanupContext()
	dumpDebugWindowsCleanupPlan(ctx)
	total := cleanupTempArtifacts(ctx)
	total += cleanupBackupArtifacts(ctx)
	total += cleanupDriveRootShim(ctx)
	total += cleanupCloneSwapDirs(ctx)
	printUpdateCleanupResult(total)
	logUpdateCleanup(constants.UpdateCleanupLogDone, total)
	logHandoffEvent("cleanup", "done",
		map[string]string{"removed": fmt.Sprintf("%d", total)})
	dumpDebugWindowsNote("cleanup finished, removed=%d", total)
}

// delayUpdateCleanupIfNeeded gives the just-exited handoff/update process time
// to release Windows file handles before deletion begins.
func delayUpdateCleanupIfNeeded() {
	raw := os.Getenv(constants.EnvUpdateCleanupDelayMS)
	if len(raw) == 0 {
		return
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms <= 0 {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateCleanDelayInvalid, raw)
		logUpdateCleanup(constants.UpdateCleanupLogDelayInvalid, raw)
		logHandoffEvent("cleanup", "delay_invalid",
			map[string]string{"raw": raw})

		return
	}
	fmt.Printf(constants.MsgUpdateCleanDelay, ms)
	logHandoffEvent("cleanup", "delay",
		map[string]string{"ms": fmt.Sprintf("%d", ms)})
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// printUpdateCleanupResult reports the cleanup result summary.
func printUpdateCleanupResult(total int) {
	if total > 0 {
		fmt.Printf(constants.MsgUpdateCleanDone, total)

		return
	}

	fmt.Println(constants.MsgUpdateCleanNone)
}

// logUpdateCleanup writes cleanup diagnostics to the shared verbose logger.
func logUpdateCleanup(format string, args ...interface{}) {
	log := verbose.Get()
	if log != nil {
		log.Log(format, args...)
	}
}
