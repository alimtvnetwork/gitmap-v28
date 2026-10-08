package cmdpipeline

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var ansiRegex = lazyregex.AnsiEscapeRegex

var failureMarkers = []string{
	"##[error]",
	"--- FAIL:",
	"FAIL\t",
	"FAIL:",
	"FAILED",
	"FAIL ",
	"Expected",
	"AssertionError",
	"Error [",
	"Traceback (most recent call last):",
	"FAILED (failures=",
	"fatal error:",
	"syntax error:",
	"exit status",
	"exit code",
	"exited ",
	"panic:",
	"Stack Trace:",
	"Unexpected token",
	"Not Found - ",
	"Error:",
	"error:",
	"[gosec-",
	"gofmt",
	"FAILED TESTS",
	"diff (baseline-diff",
	"Suite:",
	"got \"",
	"want \"",
	"lint error",
	"Process completed with exit code",
	"failed to bundle",
	"Error failed to bundle",
	"Failed to copy binary",
	"failed to copy",
	"failed to build",
	"failed to compile",
	"does not exist",
	"error[E",
	"Error [E",
	"command not found",
	"is not recognized as an internal or external command",
	"no such file or directory",
	"cannot find module",
}
