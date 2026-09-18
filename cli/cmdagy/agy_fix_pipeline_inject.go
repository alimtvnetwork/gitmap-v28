package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// InjectAgyFixTask attempts direct background dispatch of the fix prompt to Antigravity CLI.
func InjectAgyFixTask(repoDir, absPayloadPath string, skipInject bool) (bool, string) {
	if skipInject {
		return false, "direct injection skipped by flag (--no-inject)"
	}

	binPath, hasAgy := resolveAntigravityBinary()
	if !hasAgy {
		return false, "antigravity binary (agy) not detected on PATH or system locations"
	}

	return launchAgyBackground(binPath, repoDir, absPayloadPath)
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

func launchAgyBackground(binPath, repoDir, absPayloadPath string) (bool, string) {
	promptArg := buildAgyPromptArg(absPayloadPath)
	cmd := exec.Command(binPath, "--dangerously-skip-permissions", "-p", promptArg)
	attachInjectionLog(cmd, repoDir)
	configureBackgroundProcess(cmd)

	if err := cmd.Start(); err != nil {
		return false, fmt.Sprintf("failed to launch agy process: %v", err)
	}

	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	return true, fmt.Sprintf("Injected fix task directly into Antigravity IDE (PID: %d)", pid)
}
