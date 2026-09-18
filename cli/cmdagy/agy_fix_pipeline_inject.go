package cmdagy

import (
	"fmt"
	"os/exec"
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

func launchAgyBackground(binPath, repoDir, absPayloadPath string) (bool, string) {
	promptArg := buildAgyPromptArg(absPayloadPath)
	cmd := exec.Command(binPath, "-p", promptArg)
	if len(repoDir) > 0 {
		cmd.Dir = repoDir
	}

	if err := cmd.Start(); err != nil {
		return false, fmt.Sprintf("failed to launch agy process: %v", err)
	}

	pid := cmd.Process.Pid
	go func() { _ = cmd.Wait() }()

	return true, fmt.Sprintf("Injected fix task directly into Antigravity IDE (PID: %d)", pid)
}
