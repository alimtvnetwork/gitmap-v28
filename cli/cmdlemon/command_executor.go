package cmdlemon

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"time"
)

// ExecutionResult encapsulates the outcome of a command execution.
type ExecutionResult struct {
	ExitCode    int    `json:"exitCode"`
	Stdout      string `json:"stdout"`
	Stderr      string `json:"stderr"`
	DurationMs  int64  `json:"durationMs"`
	IsSuccess   bool   `json:"isSuccess"`
	IsTruncated bool   `json:"isTruncated"`
}

// CommandExecutor abstracts command execution for production and testing.
type CommandExecutor interface {
	Execute(ctx context.Context, rawCmd string) (ExecutionResult, error)
	ExecuteArgs(ctx context.Context, args []string) (ExecutionResult, error)
}

// SystemCommandExecutor runs commands against the local GitMap binary.
type SystemCommandExecutor struct {
	Timeout time.Duration
	ExePath string
}

func (e *SystemCommandExecutor) resolveBinary() string {
	if e.ExePath != "" {
		return e.ExePath
	}
	if exe, err := os.Executable(); err == nil && exe != "" {
		return exe
	}
	return "gitmap"
}

// Execute parses, authorizes, and runs the raw command string.
func (e *SystemCommandExecutor) Execute(ctx context.Context, raw string) (ExecutionResult, error) {
	args, err := ParseAndAuthorizeCommand(raw)
	if err != nil {
		return ExecutionResult{ExitCode: 1, Stderr: err.Error()}, err
	}
	return e.ExecuteArgs(ctx, args)
}

// ExecuteArgs validates and executes tokenized arguments.
func (e *SystemCommandExecutor) ExecuteArgs(ctx context.Context, args []string) (ExecutionResult, error) {
	validArgs, ok := ValidateTokens(args)
	if !ok {
		return ExecutionResult{ExitCode: 1, Stderr: "unauthorized"}, errors.New("unauthorized command")
	}
	execCtx, cancel := context.WithTimeout(ctx, ResolveTimeout(e.Timeout))
	defer cancel()
	cmd := exec.CommandContext(execCtx, e.resolveBinary(), validArgs...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdoutBuf, &stderrBuf
	start := time.Now()
	runErr := cmd.Run()
	out, outTrunc := TruncateBuffer(stdoutBuf.Bytes(), MaxBufferBytes)
	errStr, errTrunc := TruncateBuffer(stderrBuf.Bytes(), MaxBufferBytes)
	exitCode := ResolveExitCode(runErr)
	return ExecutionResult{
		ExitCode: exitCode, Stdout: out, Stderr: errStr,
		DurationMs:  time.Since(start).Milliseconds(),
		IsSuccess:   runErr == nil && exitCode == 0,
		IsTruncated: outTrunc || errTrunc,
	}, runErr
}
