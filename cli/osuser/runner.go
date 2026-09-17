package osuser

import "os/exec"

// OSCommandRunner executes an OS command and returns combined output.
type OSCommandRunner func(cmd *exec.Cmd) ([]byte, error)

// defaultOSCommandRunner is the active command runner, customizable for hermetic testing.
var defaultOSCommandRunner OSCommandRunner = func(cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}
