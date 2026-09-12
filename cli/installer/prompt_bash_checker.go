// Package installer — prompt_bash_checker.go verifies bash and curl on PATH.
package installer

import "os/exec"

func HasBashAndCurl() bool {
	_, errBash := exec.LookPath("bash")
	_, errCurl := exec.LookPath("curl")

	return errBash == nil && errCurl == nil
}
