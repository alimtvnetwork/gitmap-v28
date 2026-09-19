package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// InjectAgyFixTask attempts dispatch of the fix prompt to Antigravity CLI or active IDE.
func InjectAgyFixTask(repoDir, absPayloadPath string, isSkipInject bool) AgyInjectionResult {
	if isSkipInject {
		return makeSkipInjectionResult()
	}

	repoRoot := resolveTargetRepoRoot(repoDir)
	promptPath := stageActivePromptInRepo(repoRoot, absPayloadPath)
	ideProcRes := DetectRunningAntigravityIDE()
	cliRes := ResolveAntigravityCLI()

	if cliRes.IsSuccess() {
		return launchAgyBackgroundRunner(cliRes.Value, repoRoot, promptPath, ideProcRes)
	}

	if ideProcRes.IsSuccess() {
		return makeIDESuccessResult(ideProcRes.Value.PID, repoRoot, promptPath)
	}

	return makeNoTargetResult(repoRoot, promptPath)
}

func resolveTargetRepoRoot(repoDir string) string {
	startPath := resolveInitialPath(repoDir)
	root, err := gitutil.RepoRoot(startPath)
	if err == nil && len(root) > 0 {
		return root
	}

	return toAbsPath(startPath)
}

func resolveInitialPath(repoDir string) string {
	if len(repoDir) > 0 {
		return repoDir
	}

	return resolveProjectRootDir()
}

func stageActivePromptInRepo(repoRoot, absPayloadPath string) string {
	targetPrompt := filepath.Join(repoRoot, activeAgyPromptRelativePath)
	if isSamePath(targetPrompt, absPayloadPath) {
		return targetPrompt
	}

	content := readPromptContentOrDefault(absPayloadPath)
	writePromptFile(targetPrompt, content)

	return targetPrompt
}

func isSamePath(pathA, pathB string) bool {
	cleanA := filepath.Clean(pathA)
	cleanB := filepath.Clean(pathB)

	return strings.EqualFold(cleanA, cleanB)
}

func readPromptContentOrDefault(path string) string {
	if len(path) == 0 {
		return ""
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(data)
}

func writePromptFile(path, content string) {
	if len(content) == 0 {
		return
	}

	_ = os.MkdirAll(filepath.Dir(path), 0755)
	_ = os.WriteFile(path, []byte(content), 0644)
}

func buildAgyPromptArg(absPayloadPath string) string {
	return fmt.Sprintf("Autonomous CI/CD pipeline fix: follow all directives in %s", absPayloadPath)
}

func attachInjectionLog(cmd *exec.Cmd, repoDir string) {
	if len(repoDir) == 0 {
		return
	}

	cmd.Dir = repoDir
	logPath := filepath.Join(repoDir, ".ai-memory", "temp", "agy-injection.log")
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}

	cmd.Stdout = logFile
	cmd.Stderr = logFile
}

func launchAgyBackgroundRunner(
	binPath string,
	repoDir string,
	promptPath string,
	ideRes result.Result[AgyProcessInfo],
) AgyInjectionResult {
	promptArg := buildAgyPromptArg(promptPath)
	cmd := exec.Command(binPath, "--dangerously-skip-permissions", "-p", promptArg)
	attachInjectionLog(cmd, repoDir)
	configureBackgroundProcess(cmd)

	if err := cmd.Start(); err != nil {
		return makeLaunchFailureResult(err, repoDir, promptPath)
	}

	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	return makeLaunchSuccessResult(pid, repoDir, promptPath, ideRes)
}

func makeLaunchSuccessResult(
	pid int,
	repoDir string,
	promptPath string,
	ideRes result.Result[AgyProcessInfo],
) AgyInjectionResult {
	msg := formatLaunchSuccessMessage(pid, repoDir, ideRes)

	return AgyInjectionResult{
		IsSuccess:  true,
		Mode:       AgyInjectionModeCLI,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func formatLaunchSuccessMessage(
	pid int,
	repoDir string,
	ideRes result.Result[AgyProcessInfo],
) string {
	if ideRes.IsSuccess() {
		return fmt.Sprintf("Started Antigravity CLI background runner (PID: %d) with active IDE detected (PID: %d)", pid, ideRes.Value.PID)
	}

	return fmt.Sprintf("Started Antigravity CLI background runner (PID: %d) in %s", pid, repoDir)
}

func makeIDESuccessResult(pid int, repoDir, promptPath string) AgyInjectionResult {
	msg := fmt.Sprintf("Active Antigravity IDE detected (PID: %d); staged fix prompt in %s", pid, promptPath)

	return AgyInjectionResult{
		IsSuccess:  true,
		Mode:       AgyInjectionModeIDE,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func makeLaunchFailureResult(err error, repoDir, promptPath string) AgyInjectionResult {
	return AgyInjectionResult{
		IsSuccess:  false,
		Mode:       AgyInjectionModeCLI,
		Message:    fmt.Sprintf("failed to launch agy CLI process: %v", err),
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func makeSkipInjectionResult() AgyInjectionResult {
	return AgyInjectionResult{
		IsSuccess: false,
		Mode:      AgyInjectionModeNone,
		Message:   "direct injection skipped by flag (--no-inject)",
	}
}

func makeNoTargetResult(repoDir, promptPath string) AgyInjectionResult {
	return AgyInjectionResult{
		IsSuccess:  false,
		Mode:       AgyInjectionModeNone,
		Message:    "antigravity binary (agy) not detected and no active IDE process found",
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}
