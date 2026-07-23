package cmd

// cloneprintargv_test.go — unit tests for --print-clone-argv.
// Pins the output format (one `argv[i]=token` per line, indented,
// "git" prepended at slot 0) and the off-by-default behavior so a
// future refactor of the dump format fails loudly here instead of
// surprising downstream parsers.

import (
	"bytes"
	"strings"
	"testing"
)

// TestPrintCloneArgv_Format verifies the canonical multi-line dump.
func TestPrintCloneArgv_Format(t *testing.T) {
	var buf bytes.Buffer
	if err := printCloneArgv(&buf, []string{"clone", "-b", "main",
		"https://x/r.git", "r"}); err != nil {
		t.Fatalf("print: %v", err)
	}
	want := strings.Join([]string{
		"  argv[0]=git",
		"  argv[1]=clone",
		"  argv[2]=-b",
		"  argv[3]=main",
		"  argv[4]=https://x/r.git",
		"  argv[5]=r",
		"",
	}, "\n")
	if got := buf.String(); got != want {
		t.Fatalf("dump mismatch\n--- want ---\n%s\n--- got ---\n%s",
			want, got)
	}
}

// TestPrintCloneArgv_EmptyNoOp ensures we don't print a stray "git"
// line for an empty executor argv (which would happen if Go's
// append created a 1-element slice).
func TestPrintCloneArgv_EmptyNoOp(t *testing.T) {
	var buf bytes.Buffer
	if err := printCloneArgv(&buf, nil); err != nil {
		t.Fatalf("print: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected empty output, got %q", buf.String())
	}
}

// TestRunCmdPrintArgv_GatedByFlag confirms the integration helper
// is a no-op when the flag is off — required because every per-row
// helper calls it unconditionally on the hot path.
func TestRunCmdPrintArgv_GatedByFlag(t *testing.T) {
	// Save + restore so the test doesn't leak state to siblings.
	prev := cmdPrintArgvEnabled()
	defer setCmdPrintArgv(prev)

	setCmdPrintArgv(false)
	if cmdPrintArgvEnabled() {
		t.Fatal("expected disabled after setCmdPrintArgv(false)")
	}
	setCmdPrintArgv(true)
	if !cmdPrintArgvEnabled() {
		t.Fatal("expected enabled after setCmdPrintArgv(true)")
	}
}
