package cluster

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"strings"
)

const (
	bashCmd = "bash"
	shCmd   = "sh"
	cFlag   = "-c"
)

// ExecBash executes a bash command on the local node.
func ExecBash(
	ctx context.Context,
	node ClusterNode,
	command string,
) (stdout, stderr string, exitCode int, err error) {
	cmdPath := resolveBashPath()
	if cmdPath == "" {
		return "", "bash not found", 1, exec.ErrNotFound
	}

	cmd := exec.CommandContext(ctx, cmdPath, cFlag, command)

	return executeCmdBuffered(cmd)
}

// ExecShell executes a POSIX shell command on the local node.
func ExecShell(
	ctx context.Context,
	node ClusterNode,
	command string,
) (stdout, stderr string, exitCode int, err error) {
	cmdPath := resolveShellPath()
	if cmdPath == "" {
		return "", "sh not found", 1, exec.ErrNotFound
	}

	cmd := exec.CommandContext(ctx, cmdPath, cFlag, command)

	return executeCmdBuffered(cmd)
}

func resolveBashPath() string {
	path, err := lookPathFuncVar(bashCmd)
	if err == nil && path != "" {
		return path
	}

	return resolveWindowsGitBash()
}

func resolveShellPath() string {
	path, err := lookPathFuncVar(shCmd)
	if err == nil && path != "" {
		return path
	}

	return resolveBashPath()
}

func resolveWindowsGitBash() string {
	if runtime.GOOS != winOS {
		return ""
	}

	candidates := []string{
		`C:\Program Files\Git\bin\bash.exe`,
		`C:\Program Files\Git\usr\bin\bash.exe`,
	}
	if gitPath, err := lookPathFuncVar("git"); err == nil && gitPath != "" {
		candidates = append(candidates, deriveBashFromGit(gitPath)...)
	}
	for _, cand := range candidates {
		if path, err := lookPathFuncVar(cand); err == nil && path != "" {
			return path
		}
	}

	return ""
}

func deriveBashFromGit(gitPath string) []string {
	norm := strings.ReplaceAll(gitPath, "\\", "/")
	idx := strings.LastIndex(norm, "/cmd/git.exe")
	if idx >= 0 {
		root := gitPath[:idx]
		return []string{
			root + `\bin\bash.exe`,
			root + `\usr\bin\bash.exe`,
		}
	}
	return nil
}

func executeCmdBuffered(cmd *exec.Cmd) (string, string, int, error) {
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := runCmdFunc(cmd)
	exitCode := extractExitCode(err)

	return outBuf.String(), errBuf.String(), exitCode, err
}
