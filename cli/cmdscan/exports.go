package cmdscan

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/errreport"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// Delegate hooks
var (
	CreatePendingTaskFn     func(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB)
	CompletePendingTaskFn   func(db *store.DB, taskID int64)
	FailPendingTaskFn       func(db *store.DB, taskID int64, reason string)
	SyncRecordsToVSCodePMFn func(records []model.ScanRecord, noVSCodeSync, noAutoTags bool)
	RunPruneStaleDBFn       func(absDir string, records []model.ScanRecord) error
	CheckHelpFn             func(subcmd string, args []string)
)

// RunScan executes the scan command.
func RunScan(args []string) error {
	return runScan(args)
}

// RunRescan executes the rescan command.
func RunRescan() error {
	return runRescan()
}

// RunRescanSubtree executes the rescan-subtree command.
func RunRescanSubtree(args []string) error {
	return runRescanSubtree(args)
}

// AutoRegisterFirstWorkDir exposes autoRegisterFirstWorkDir.
func AutoRegisterFirstWorkDir(absDir string, quiet bool) bool {
	return autoRegisterFirstWorkDir(absDir, quiet)
}

// ExpandHome exposes expandHome.
func ExpandHome(p string) string {
	return expandHome(p)
}

// ResolveOutFile exposes resolveOutFile.
func ResolveOutFile(outFile, outputDir, defaultName string) string {
	return resolveOutFile(outFile, outputDir, defaultName)
}

// WriteAllOutputs exposes writeAllOutputs.
func WriteAllOutputs(records []model.ScanRecord, outputDir, outFile string, quiet, compact bool) {
	writeAllOutputs(records, outputDir, outFile, quiet, compact)
}

// FinalizeErrorReport exposes finalizeErrorReport.
func FinalizeErrorReport(c *errreport.Collector, quiet bool) {
	finalizeErrorReport(c, quiet)
}
