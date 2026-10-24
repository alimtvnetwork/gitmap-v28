package cmd

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type commandExecReq struct {
	Command string `json:"command"`
}

type commandExecResp struct {
	Status     string   `json:"status,omitempty"`
	Success    bool     `json:"success"`
	Output     string   `json:"output"`
	Stdout     string   `json:"stdout,omitempty"`
	Stderr     string   `json:"stderr,omitempty"`
	ExitCode   int      `json:"exitCode"`
	Errors     []string `json:"errors,omitempty"`
	Error      string   `json:"error,omitempty"`
	DurationMs int64    `json:"durationMs,omitempty"`
}

func termCommandExecHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeTermJSONError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var req commandExecReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeTermJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := executeCLIForAPI(req.Command)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func handleCommandExec(w http.ResponseWriter, r *http.Request) {
	termCommandExecHandler(w, r)
}

func executeCLIForAPI(cmdStr string) commandExecResp {
	cmd := buildPlatformCmd(cmdStr)
	start := time.Now()
	out, err := cmd.CombinedOutput()
	durationMs := time.Since(start).Milliseconds()

	if err == nil {
		return commandExecResp{
			Status:     "success",
			Success:    true,
			Output:     string(out),
			Stdout:     string(out),
			ExitCode:   0,
			DurationMs: durationMs,
		}
	}

	exitCode := resolveExitCode(err)

	return commandExecResp{
		Status:     "error",
		Success:    false,
		Output:     string(out),
		Stderr:     string(out),
		ExitCode:   exitCode,
		Errors:     []string{err.Error()},
		Error:      err.Error(),
		DurationMs: durationMs,
	}
}

func buildPlatformCmd(cmdStr string) *exec.Cmd {
	if constants.IsWindowsPlatform() {
		return exec.Command("cmd.exe", "/C", cmdStr)
	}

	return exec.Command("sh", "-c", cmdStr)
}

func resolveExitCode(err error) int {
	exitErr, isExit := err.(*exec.ExitError)
	if isExit {
		return exitErr.ExitCode()
	}

	return 1
}
