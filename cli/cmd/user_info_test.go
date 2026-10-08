package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

func TestExecuteUserInfo_TerminalCard(t *testing.T) {
	var buf bytes.Buffer
	err := executeUserInfo(&buf, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Git & GitHub Identity Context") {
		t.Errorf("expected card title in output, got: %s", out)
	}
}

func TestExecuteUserInfo_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := executeUserInfo(&buf, []string{"--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var summary usercontext.UserSummary
	if err := json.Unmarshal(buf.Bytes(), &summary); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v\nOutput: %s", err, buf.String())
	}
}

func TestExecuteUserInfo_JSONShortFlag(t *testing.T) {
	var buf bytes.Buffer
	err := executeUserInfo(&buf, []string{"-j"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var summary usercontext.UserSummary
	if err := json.Unmarshal(buf.Bytes(), &summary); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v\nOutput: %s", err, buf.String())
	}
}
