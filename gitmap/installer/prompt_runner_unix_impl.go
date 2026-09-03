// Package installer — prompt_runner_unix.go executes the remote bash installer on Unix.
package installer

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func RunUnixPromptInstaller(targetDir string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmdStr := "curl -sL " + constants.PromptArchitectBashURL + " | bash"
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
