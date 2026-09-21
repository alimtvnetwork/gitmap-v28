package cmdai

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ResolveRepoRoot walks up the directory hierarchy to find the repository root.
func ResolveRepoRoot() (string, *apperror.AppError) {
	curr, getErr := os.Getwd()
	hasGetErr := getErr != nil
	if hasGetErr {
		return "", apperror.WrapSimple(getErr, "get_repo_root")
	}

	root, isFound := findRootInAncestors(curr)
	if isFound {
		return root, nil
	}

	ctx := map[string]any{"startDir": curr}

	return "", apperror.New("resolve_root", "E_REPO_ROOT_NOT_FOUND", ctx)
}

func findRootInAncestors(start string) (string, bool) {
	dir := start
	for {
		isRoot := hasRepoMarkers(dir)
		if isRoot {
			return dir, true
		}

		parent := filepath.Dir(dir)
		isTop := parent == dir
		if isTop {
			return "", false
		}

		dir = parent
	}
}

func hasRepoMarkers(dir string) bool {
	scriptsDir := filepath.Join(dir, "03-ai-scripts")
	scriptsInfo, scriptsErr := os.Stat(scriptsDir)
	hasScripts := scriptsErr == nil && scriptsInfo.IsDir()
	if hasScripts {
		return true
	}

	gitDir := filepath.Join(dir, ".git")
	gitInfo, gitErr := os.Stat(gitDir)
	hasGit := gitErr == nil && gitInfo.IsDir()

	return hasGit
}

// ResolveScriptPath finds the absolute file path for a script filename.
func ResolveScriptPath(filename string) (string, *apperror.AppError) {
	repoRoot, rootErr := ResolveRepoRoot()
	hasRootErr := rootErr != nil
	if hasRootErr {
		return "", rootErr
	}

	fullPath := filepath.Join(repoRoot, "03-ai-scripts", filename)

	return verifyScriptFile(fullPath, filename)
}

func verifyScriptFile(fullPath, filename string) (string, *apperror.AppError) {
	info, statErr := os.Stat(fullPath)
	hasStatErr := statErr != nil
	if hasStatErr {
		ctx := map[string]any{"filename": filename, "path": fullPath}

		return "", apperror.New("resolve_path", "E_SCRIPT_NOT_FOUND_ON_DISK", ctx)
	}

	isDir := info.IsDir()
	if isDir {
		ctx := map[string]any{"filename": filename, "path": fullPath}

		return "", apperror.New("resolve_path", "E_SCRIPT_PATH_IS_DIR", ctx)
	}

	return fullPath, nil
}

// ExecuteScriptStreaming runs a python script with unbuffered terminal I/O.
func ExecuteScriptStreaming(ctx context.Context, scriptPath string, args []string) *apperror.AppError {
	rt, rtErr := DetectPythonRuntime()
	hasRtErr := rtErr != nil
	if hasRtErr {
		return rtErr
	}

	repoRoot, rootErr := ResolveRepoRoot()
	hasRootErr := rootErr != nil
	if hasRootErr {
		return rootErr
	}

	return runStreamingCommand(ctx, rt.ExecutablePath, repoRoot, scriptPath, args)
}

func runStreamingCommand(ctx context.Context, pythonExe, repoRoot, scriptPath string, args []string) *apperror.AppError {
	cmdArgs := buildCmdArgs(scriptPath, args)
	cmd := exec.CommandContext(ctx, pythonExe, cmdArgs...)
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = buildStreamingEnv()

	start := time.Now()
	runErr := cmd.Run()
	durationMs := int(time.Since(start).Milliseconds())
	recordExecutionHistory(scriptPath, args, repoRoot, durationMs, runErr)
	if runErr != nil {
		return wrapProcessError(runErr, scriptPath)
	}

	return nil
}

func recordExecutionHistory(scriptPath string, args []string, repoRoot string, durationMs int, err error) {
	argsJSON, _ := json.Marshal(args)
	exitCode := 0
	errMsg := ""
	isSuccess := err == nil
	if err != nil {
		exitCode = extractExitCode(err)
		errMsg = err.Error()
	}
	category := resolveScriptCategory(scriptPath)
	_ = store.RecordAiExecution(
		category,
		scriptPath,
		string(argsJSON),
		repoRoot,
		"127.0.0.1",
		durationMs,
		exitCode,
		"",
		errMsg,
		isSuccess,
	)
}

func resolveScriptCategory(scriptPath string) string {
	base := filepath.Base(scriptPath)
	if strings.Contains(base, "fix") {
		return "code_gen"
	}
	if strings.Contains(base, "test") || strings.Contains(base, "audit") || strings.Contains(base, "check") {
		return "audit"
	}
	if strings.Contains(base, "search") || strings.Contains(base, "find") || strings.Contains(base, "grep") {
		return "search"
	}

	return "instruction"
}

func buildCmdArgs(scriptPath string, args []string) []string {
	cmdArgs := []string{"-u", scriptPath}

	return append(cmdArgs, args...)
}

func buildStreamingEnv() []string {
	env := os.Environ()

	return append(env, "PYTHONUNBUFFERED=1")
}

func wrapProcessError(err error, scriptPath string) *apperror.AppError {
	exitCode := extractExitCode(err)
	ctx := map[string]any{
		"script":   scriptPath,
		"exitCode": exitCode,
	}

	return apperror.WrapWithDetails(
		err,
		"execute_script",
		"E_SCRIPT_EXEC_FAILED",
		"script execution failed",
		"ai_exec",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		ctx,
	)
}

func extractExitCode(err error) int {
	exitErr, isExitErr := err.(*exec.ExitError)
	if isExitErr {
		return exitErr.ExitCode()
	}

	return 1
}
