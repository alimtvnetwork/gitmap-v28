package cmdpy

import (
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func recordPyTelemetry(binary string, args []string, dir string, durationMs int64, exitCode int, err error) {
	cmdLine := "py " + strings.Join(args, " ")
	recordHistoryTelemetry(cmdLine, exitCode, durationMs)
	recordAiTelemetry(cmdLine, args, dir, durationMs, exitCode, err)
}

func recordHistoryTelemetry(cmdLine string, exitCode int, durationMs int64) {
	histDB, err := store.OpenCommandHistorySplitDB("")
	if err != nil {
		return
	}
	defer histDB.Close()
	_ = histDB.InsertCommandRecord(cmdLine, "py", exitCode, durationMs)
}

func resolveErrMsg(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func recordAiTelemetry(cmdLine string, args []string, dir string, durationMs int64, exitCode int, err error) {
	argsJSON, _ := json.Marshal(args)
	isSuccess := exitCode == 0
	_ = store.RecordAiExecution(
		"python_runner",
		cmdLine,
		string(argsJSON),
		dir,
		"127.0.0.1",
		int(durationMs),
		exitCode,
		"",
		resolveErrMsg(err),
		isSuccess,
	)
}
