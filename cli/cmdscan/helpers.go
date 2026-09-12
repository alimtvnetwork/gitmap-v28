package cmdscan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/probe"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func createPendingTask(typeName, targetPath, workDir, sourceCmd, cmdArgs string) (int64, *store.DB) {
	if CreatePendingTaskFn != nil {
		return CreatePendingTaskFn(typeName, targetPath, workDir, sourceCmd, cmdArgs)
	}

	return 0, nil
}

func completePendingTask(db *store.DB, taskID int64) {
	if CompletePendingTaskFn != nil {
		CompletePendingTaskFn(db, taskID)
	}
}

func failPendingTask(db *store.DB, taskID int64, reason string) {
	if FailPendingTaskFn != nil {
		FailPendingTaskFn(db, taskID, reason)
	}
}

func syncRecordsToVSCodePM(records []model.ScanRecord, noVSCodeSync, noAutoTags bool) {
	if SyncRecordsToVSCodePMFn != nil {
		SyncRecordsToVSCodePMFn(records, noVSCodeSync, noAutoTags)
	}
}

func runPruneStaleDB(absDir string, records []model.ScanRecord) error {
	if RunPruneStaleDBFn != nil {
		return RunPruneStaleDBFn(absDir, records)
	}

	return nil
}

func checkHelp(subcmd string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(subcmd, args)
	}
}

func buildCommandArgs(args []string) string {
	if len(args) <= 1 {
		return ""
	}

	return strings.Join(args[1:], " ")
}

func pickProbeURL(r model.ScanRecord) string {
	if r.Transport == constants.ScanTransportSSH && r.SSHUrl != "" {
		return r.SSHUrl
	}

	if r.HTTPSUrl != "" {
		return r.HTTPSUrl
	}

	return r.SSHUrl
}

func recordProbeResult(db *store.DB, repo model.ScanRecord, result probe.Result) {
	if err := db.RecordVersionProbe(result.AsModel(repo.ID)); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func resolveBinaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return filepath.Dir(exe)
	}

	return filepath.Dir(resolved)
}
