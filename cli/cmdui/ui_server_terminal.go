package cmdui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type fullTermExecResp struct {
	Success    bool   `json:"success"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	Output     string `json:"output"`
	ExitCode   int    `json:"exitCode"`
	Error      string `json:"error,omitempty"`
	NodeAlias  string `json:"nodeAlias"`
	DurationMs int64  `json:"durationMs"`
}

func handleAPITerminalExec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(TerminalExecResp{Success: false, Error: "Method not allowed, POST required"})
		return
	}

	var req TerminalExecReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(TerminalExecResp{Success: false, Error: "Invalid JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(req.Command) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(TerminalExecResp{Success: false, Error: "Command cannot be empty"})
		return
	}

	timeoutSec := req.TimeoutSec
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	if timeoutSec > 120 {
		timeoutSec = 120
	}

	var resp TerminalExecResp
	isRemote := req.NodeAlias != "" && !strings.EqualFold(req.NodeAlias, "local")
	if isRemote {
		resp = executeRemoteTerminal(req.NodeAlias, req.Command, timeoutSec)
	} else {
		resp = executeLocalTerminal(req.Command, req.Cwd, timeoutSec)
	}

	output := combineTerminalOutput(resp.Stdout, resp.Stderr)

	_ = json.NewEncoder(w).Encode(fullTermExecResp{
		Success:    resp.Success,
		Stdout:     resp.Stdout,
		Stderr:     resp.Stderr,
		Output:     output,
		ExitCode:   resp.ExitCode,
		Error:      resp.Error,
		NodeAlias:  resp.NodeAlias,
		DurationMs: resp.DurationMs,
	})
}

func combineTerminalOutput(stdout, stderr string) string {
	if stderr == "" {
		return stdout
	}
	if stdout == "" {
		return stderr
	}
	return stdout + "\n" + stderr
}

func executeRemoteTerminal(nodeAlias, command string, timeoutSec int) TerminalExecResp {
	conn, hasConn := findSSHNodeByAlias(nodeAlias)
	if !hasConn {
		return TerminalExecResp{
			Success:   false,
			NodeAlias: nodeAlias,
			ExitCode:  1,
			Error:     "remote node not found: " + nodeAlias,
		}
	}

	client, isConnected := cmdssh.ConnectSSHClient(conn, fmt.Sprintf("[%s]", nodeAlias))
	if !isConnected {
		return TerminalExecResp{
			Success:   false,
			NodeAlias: nodeAlias,
			ExitCode:  1,
			Error:     "ssh connection failed to " + nodeAlias,
		}
	}
	defer client.Close()

	session, sessErr := client.NewSession()
	if sessErr != nil {
		return TerminalExecResp{
			Success:   false,
			NodeAlias: nodeAlias,
			ExitCode:  1,
			Error:     sessErr.Error(),
		}
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	start := time.Now()
	done := make(chan error, 1)

	go func() {
		done <- session.Run(command)
	}()

	var runErr error
	select {
	case <-time.After(time.Duration(timeoutSec) * time.Second):
		_ = session.Close()
		return TerminalExecResp{
			Success:    false,
			NodeAlias:  nodeAlias,
			ExitCode:   124,
			Error:      fmt.Sprintf("command timed out after %d seconds", timeoutSec),
			DurationMs: time.Since(start).Milliseconds(),
		}
	case runErr = <-done:
	}

	durationMs := time.Since(start).Milliseconds()
	exitCode := resolveSSHExitCode(runErr)
	errMsg := ""
	if runErr != nil {
		errMsg = runErr.Error()
	}

	return TerminalExecResp{
		Success:    runErr == nil && exitCode == 0,
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
		ExitCode:   exitCode,
		Error:      errMsg,
		NodeAlias:  nodeAlias,
		DurationMs: durationMs,
	}
}

func executeLocalTerminal(command, cwd string, timeoutSec int) TerminalExecResp {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd.exe", "/c", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	if cwd != "" {
		cmd.Dir = cwd
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	start := time.Now()
	runErr := cmd.Run()
	durationMs := time.Since(start).Milliseconds()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return TerminalExecResp{
			Success:    false,
			ExitCode:   124,
			Error:      fmt.Sprintf("command timed out after %d seconds", timeoutSec),
			NodeAlias:  "local",
			DurationMs: durationMs,
		}
	}

	exitCode := resolveExitCode(runErr)
	errMsg := ""
	if runErr != nil {
		errMsg = runErr.Error()
	}

	return TerminalExecResp{
		Success:    runErr == nil && exitCode == 0,
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
		ExitCode:   exitCode,
		Error:      errMsg,
		NodeAlias:  "local",
		DurationMs: durationMs,
	}
}

func resolveExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return 1
}

func resolveSSHExitCode(err error) int {
	if err == nil {
		return 0
	}
	var sshExitErr *ssh.ExitError
	if errors.As(err, &sshExitErr) {
		return sshExitErr.ExitStatus()
	}
	return 1
}
