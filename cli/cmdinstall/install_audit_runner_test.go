package cmdinstall

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteCommandWithAudit_Success(t *testing.T) {
	args := []string{"go", "version"}
	res := executeCommandWithAudit(args, false)

	if res.IsFailed() {
		t.Fatalf("expected successful execution, got err: %v", res.Err)
	}

	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}

	hasGo := strings.Contains(res.Stdout, "go version")
	if !hasGo {
		t.Errorf("expected stdout to contain 'go version', got: %q", res.Stdout)
	}

	if res.CommandLine != "go version" {
		t.Errorf("expected command line 'go version', got %q", res.CommandLine)
	}
}

func TestExecuteCommandWithAudit_Failure(t *testing.T) {
	args := []string{"go", "definitely-invalid-subcommand-12345"}
	res := executeCommandWithAudit(args, false)

	if res.IsSuccess {
		t.Fatalf("expected command to fail, but it succeeded")
	}

	if res.ExitCode == 0 {
		t.Errorf("expected non-zero exit code on failure, got %d", res.ExitCode)
	}

	if res.Err == nil {
		t.Errorf("expected non-nil error on failure")
	}
}

func TestSelectAuditWriter(t *testing.T) {
	var buf, out bytes.Buffer

	w1 := selectAuditWriter(&buf, &out, false)
	_, _ = w1.Write([]byte("quiet"))

	if buf.String() != "quiet" || out.String() != "" {
		t.Errorf("unexpected quiet output: buf=%q out=%q", buf.String(), out.String())
	}

	buf.Reset()
	out.Reset()

	w2 := selectAuditWriter(&buf, &out, true)
	_, _ = w2.Write([]byte("loud"))

	if buf.String() != "loud" || out.String() != "loud" {
		t.Errorf("unexpected verbose output: buf=%q out=%q", buf.String(), out.String())
	}
}

func TestExtractRunExitCode_Nil(t *testing.T) {
	code := extractRunExitCode(nil)
	if code != 0 {
		t.Errorf("expected exit code 0 for nil error, got %d", code)
	}
}
