package cmdrun

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai"
)

const maxSnippetBytes = 2048

// ExecuteTarget executes a resolved target script, streaming output while capturing tail buffers.
func ExecuteTarget(target *RunTarget, opts RunOptions) (*RunResult, error) {
	binary, finalArgs := buildCommandInvocation(target, opts.Args)
	target.Interpreter = binary

	ctx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, binary, finalArgs...)
	if opts.WorkingDir != "" {
		cmd.Dir = opts.WorkingDir
	}

	cmd.Stdin = os.Stdin
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	start := time.Now()
	runErr := cmd.Run()
	durationMs := time.Since(start).Milliseconds()
	exitCode := extractExitCode(runErr)

	return &RunResult{
		Target:     *target,
		ExitCode:   exitCode,
		DurationMs: durationMs,
		Stdout:     captureTail(stdoutBuf.Bytes(), maxSnippetBytes),
		Stderr:     captureTail(stderrBuf.Bytes(), maxSnippetBytes),
		Err:        runErr,
	}, runErr
}

func captureTail(b []byte, maxBytes int) string {
	if len(b) <= maxBytes {
		return string(b)
	}

	return string(b[len(b)-maxBytes:])
}

func extractExitCode(err error) int {
	if err == nil {
		return 0
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}

	return 1
}

func buildCommandInvocation(target *RunTarget, userArgs []string) (string, []string) {
	switch target.Extension {
	case ".py":
		return resolvePythonExe(), append([]string{target.ResolvedPath}, userArgs...)
	case ".ps1":
		flags := []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", target.ResolvedPath}
		return resolvePowerShellExe(), append(flags, userArgs...)
	case ".sh", ".bash":
		return resolveBashExe(), append([]string{target.ResolvedPath}, userArgs...)
	case ".js":
		return resolveNodeExe(), append([]string{target.ResolvedPath}, userArgs...)
	case ".ts":
		return resolveTSExe(), append([]string{target.ResolvedPath}, userArgs...)
	case ".go":
		return "go", append([]string{"run", target.ResolvedPath}, userArgs...)
	default:
		return target.ResolvedPath, userArgs
	}
}

func resolvePythonExe() string {
	if venvPath := resolveVenvPython(); venvPath != "" {
		return venvPath
	}

	if rt, err := cmdai.DetectPythonRuntime(); err == nil && rt.IsAvailable && rt.ExecutablePath != "" {
		return rt.ExecutablePath
	}

	for _, cand := range []string{"python3", "python", "py"} {
		if p, err := exec.LookPath(cand); err == nil && p != "" {
			return p
		}
	}

	return "python"
}

func resolveVenvPython() string {
	venv := os.Getenv("VIRTUAL_ENV")
	if venv == "" {
		return ""
	}

	sub := buildVenvPythonSubpath(venv)
	if _, err := os.Stat(sub); err == nil {
		return sub
	}

	return ""
}

func buildVenvPythonSubpath(venv string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(venv, "Scripts", "python.exe")
	}

	return filepath.Join(venv, "bin", "python")
}

func resolvePowerShellExe() string {
	if p, err := exec.LookPath("pwsh"); err == nil && p != "" {
		return p
	}

	if p, err := exec.LookPath("powershell"); err == nil && p != "" {
		return p
	}

	return "powershell"
}

func resolveBashExe() string {
	if p, err := exec.LookPath("bash"); err == nil && p != "" {
		return p
	}

	if p, err := exec.LookPath("sh"); err == nil && p != "" {
		return p
	}

	return "bash"
}

func resolveNodeExe() string {
	for _, cand := range []string{"node", "bun", "deno"} {
		if p, err := exec.LookPath(cand); err == nil && p != "" {
			return p
		}
	}

	return "node"
}

func resolveTSExe() string {
	for _, cand := range []string{"bun", "tsx", "ts-node", "deno", "node"} {
		if p, err := exec.LookPath(cand); err == nil && p != "" {
			return p
		}
	}

	return "node"
}
