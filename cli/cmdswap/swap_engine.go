package cmdswap

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
)

// ExecuteSwap coordinates path resolution, script plan generation, and script execution.
func ExecuteSwap(opts SwapOptions) (*SwapResult, *apperror.AppError) {
	startTime := time.Now()

	resolvedOpts, err := resolveSwapPaths(opts)
	if err != nil {
		return nil, err
	}

	plan, planErr := buildScriptPlan(resolvedOpts)
	if planErr != nil {
		return nil, planErr
	}

	return handlePlanExecution(resolvedOpts, plan, startTime)
}

func handlePlanExecution(opts SwapOptions, plan *ScriptPlan, start time.Time) (*SwapResult, *apperror.AppError) {
	if opts.IsScriptOnly {
		return buildScriptOnlyResult(opts, plan, start), nil
	}

	if err := launchDetachedScript(plan); err != nil {
		return nil, err
	}

	return buildSuccessResult(opts, plan, start), nil
}

func resolveSwapPaths(opts SwapOptions) (SwapOptions, *apperror.AppError) {
	fromPath, err := resolveSourcePath(opts.FromPath)
	if err != nil {
		return opts, err
	}

	toPath, toErr := resolveTargetPath(opts.ToPath)
	if toErr != nil {
		return opts, toErr
	}

	opts.FromPath = fromPath
	opts.ToPath = toPath

	return opts, nil
}

func resolveSourcePath(explicit string) (string, *apperror.AppError) {
	if explicit != "" {
		return verifyFileExists(explicit)
	}

	return autoDetectSourcePath()
}

func autoDetectSourcePath() (string, *apperror.AppError) {
	binName := "gitmap"
	if runtime.GOOS == "windows" {
		binName = "gitmap.exe"
	}

	candidates := []string{
		filepath.Join(".", "cli", binName),
		filepath.Join(".", "bin", binName),
		filepath.Join(".", binName),
	}

	return searchCandidatePaths(candidates)
}

func searchCandidatePaths(candidates []string) (string, *apperror.AppError) {
	for _, cand := range candidates {
		absPath, err := filepath.Abs(cand)
		if err == nil && isExistingRegularFile(absPath) {
			return absPath, nil
		}
	}

	return "", apperror.NewValidationError("source binary not found in candidate locations")
}

func verifyFileExists(path string) (string, *apperror.AppError) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", apperror.WrapSimple(err, "invalid source binary path")
	}

	if !isExistingRegularFile(absPath) {
		return "", apperror.NewValidationError(fmt.Sprintf("source binary does not exist: %s", absPath))
	}

	return absPath, nil
}

func isExistingRegularFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func resolveTargetPath(explicit string) (string, *apperror.AppError) {
	if explicit != "" {
		absPath, err := filepath.Abs(explicit)
		if err != nil {
			return "", apperror.WrapSimple(err, "invalid target path")
		}

		return absPath, nil
	}

	deployed, _ := cmdinstall.ResolveDeployedAndConfigPaths()
	if deployed != "" {
		absDeployed, _ := filepath.Abs(deployed)

		return absDeployed, nil
	}

	execPath, err := os.Executable()
	if err != nil {
		return "", apperror.WrapSimple(err, "unable to determine target executable path")
	}

	absExec, _ := filepath.Abs(execPath)

	return absExec, nil
}

func buildScriptPlan(opts SwapOptions) (*ScriptPlan, *apperror.AppError) {
	pids := collectTargetPIDs(opts)
	isWin := runtime.GOOS == "windows"

	var content string
	var ext string
	if isWin {
		content = GenerateWindowsPowerShellScript(opts, pids)
		ext = ".ps1"
	} else {
		content = GenerateUnixShellScript(opts, pids)
		ext = ".sh"
	}

	return assembleScriptPlan(content, ext, pids, isWin)
}

func collectTargetPIDs(opts SwapOptions) []int {
	pids := make([]int, 0)
	if opts.TargetPID > 0 {
		pids = append(pids, opts.TargetPID)
	}

	return pids
}

func assembleScriptPlan(content, ext string, pids []int, isWin bool) (*ScriptPlan, *apperror.AppError) {
	tempFile, err := os.CreateTemp("", fmt.Sprintf("gitmap-swap-*%s", ext))
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to create temporary swap script")
	}

	scriptPath := tempFile.Name()
	_ = tempFile.Close()

	if writeErr := os.WriteFile(scriptPath, []byte(content), 0755); writeErr != nil {
		return nil, apperror.WrapSimple(writeErr, "failed to write swap script content")
	}

	return createPlanObject(scriptPath, content, pids, isWin), nil
}

func createPlanObject(path, content string, pids []int, isWin bool) *ScriptPlan {
	platform := "unix"
	execArgs := []string{"/bin/sh", path}

	if isWin {
		platform = "windows"
		execArgs = []string{"powershell.exe", "-ExecutionPolicy", "Bypass", "-File", path}
	}

	return &ScriptPlan{
		Platform:      platform,
		ScriptPath:    path,
		ScriptContent: content,
		TargetPIDList: pids,
		IsDetached:    true,
		ExecutionArgs: execArgs,
		EstimatedWait: 2 * time.Second,
	}
}

func launchDetachedScript(plan *ScriptPlan) *apperror.AppError {
	if len(plan.ExecutionArgs) == 0 {
		return apperror.NewValidationError("execution plan has empty arguments")
	}

	cmd := exec.Command(plan.ExecutionArgs[0], plan.ExecutionArgs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return apperror.WrapSimple(err, "failed to start detached swap process")
	}

	return nil
}

func buildScriptOnlyResult(opts SwapOptions, plan *ScriptPlan, start time.Time) *SwapResult {
	return &SwapResult{
		IsSuccess:      true,
		FromPath:       opts.FromPath,
		ToPath:         opts.ToPath,
		TerminatedPIDs: plan.TargetPIDList,
		ScriptPlan:     plan,
		OutputMessage:  fmt.Sprintf("Swap script generated at %s", plan.ScriptPath),
		DurationMs:     time.Since(start).Milliseconds(),
	}
}

func buildSuccessResult(opts SwapOptions, plan *ScriptPlan, start time.Time) *SwapResult {
	return &SwapResult{
		IsSuccess:      true,
		FromPath:       opts.FromPath,
		ToPath:         opts.ToPath,
		TerminatedPIDs: plan.TargetPIDList,
		ScriptPlan:     plan,
		OutputMessage:  fmt.Sprintf("Detached binary swap initiated successfully (script: %s)", plan.ScriptPath),
		DurationMs:     time.Since(start).Milliseconds(),
	}
}

// RenderResultOutput serializes the swap result to JSON or returns formatted text.
func RenderResultOutput(res *SwapResult, isJSON bool) (string, *apperror.AppError) {
	if isJSON {
		bytes, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return "", apperror.WrapSimple(err, "failed to serialize swap result to json")
		}

		return string(bytes), nil
	}

	return res.OutputMessage, nil
}
