package cmdpull

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func TestResolveCommandHint(t *testing.T) {
	hintPe := resolveCommandHint("pe")
	if !strings.Contains(hintPe, "gitmap pe") {
		t.Fatalf("expected hint for 'pe' to contain 'gitmap pe', got %q", hintPe)
	}

	hintPd := resolveCommandHint("pd")
	if !strings.Contains(hintPd, "gitmap pd") {
		t.Fatalf("expected hint for 'pd' to contain 'gitmap pd', got %q", hintPd)
	}

	hintUnknown := resolveCommandHint("random-unknown-repo")
	if len(hintUnknown) > 0 {
		t.Fatalf("expected empty hint for unknown slug, got %q", hintUnknown)
	}
}

func TestHandlePullSlugNotFound_ExitCode(t *testing.T) {
	var capturedCode int
	prevExit := cliexit.SetExitFunc(func(code int) {
		capturedCode = code
	})
	defer cliexit.SetExitFunc(prevExit)

	opts := pullOptions{slug: "nonexistent-repo"}
	handlePullTargetNotFound(opts)

	if capturedCode != int(cliexit.ExitCodeNotFound) {
		t.Fatalf("expected exit code %d, got %d", cliexit.ExitCodeNotFound, capturedCode)
	}
}

func TestHandlePullGroupNotFound_ExitCode(t *testing.T) {
	var capturedCode int
	prevExit := cliexit.SetExitFunc(func(code int) {
		capturedCode = code
	})
	defer cliexit.SetExitFunc(prevExit)

	opts := pullOptions{group: "nonexistent-group"}
	handlePullTargetNotFound(opts)

	if capturedCode != int(cliexit.ExitCodeNotFound) {
		t.Fatalf("expected exit code %d, got %d", cliexit.ExitCodeNotFound, capturedCode)
	}
}

func TestHandlePushSlugNotFound_ExitCode(t *testing.T) {
	var capturedCode int
	prevExit := cliexit.SetExitFunc(func(code int) {
		capturedCode = code
	})
	defer cliexit.SetExitFunc(prevExit)

	opts := pushOptions{slug: "pe"}
	handlePushTargetNotFound(opts)

	if capturedCode != int(cliexit.ExitCodeNotFound) {
		t.Fatalf("expected exit code %d, got %d", cliexit.ExitCodeNotFound, capturedCode)
	}
}
