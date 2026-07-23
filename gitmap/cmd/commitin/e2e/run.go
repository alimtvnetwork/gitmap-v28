package e2e

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/orchestrator"
)

// RunResult bundles every output channel an E2E test needs to assert
// on: the orchestrator exit code plus the captured stdout/stderr.
type RunResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

// Run drives orchestrator.Run with the given RawArgs, capturing both
// streams. Defaults to IsNoPrompt=true so tests never block on stdin.
// Callers can pre-populate raw with profile names, exclude rules, etc.
func Run(t *testing.T, raw *commitin.RawArgs) RunResult {
	t.Helper()
	if raw == nil {
		t.Fatal("e2e.Run: raw must not be nil")
	}
	raw.IsNoPrompt = true
	var outBuf, errBuf bytes.Buffer
	code := orchestrator.Run(raw, &outBuf, &errBuf)
	return RunResult{ExitCode: code, Stdout: outBuf.String(), Stderr: errBuf.String()}
}

// NewRawArgs returns a baseline RawArgs with `source` and `inputs`
// set and every behavior flag at its spec-default. Tests mutate the
// returned pointer to flip individual flags.
func NewRawArgs(source string, inputs ...string) *commitin.RawArgs {
	return &commitin.RawArgs{
		Source: source,
		Inputs: inputs,
	}
}
