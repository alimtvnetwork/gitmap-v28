package cmdssh

import (
	"strings"
	"testing"
)

func TestSuggestSSHSubcommand_Proof(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"hosts", "nodes"},
		{"host", "nodes"},
		{"machines", "nodes"},
		{"vms", "nodes"},
		{"updat", "update"},
		{"upgrade", "update"},
		{"exe", "exec"},
		{"cmd", "exec"},
		{"inst", "install-exec"},
		{"setup", "install-exec"},
		{"unknown-target", ""},
	}

	for _, tc := range cases {
		got := suggestSSHSubcommand(tc.input)
		if got != tc.expected {
			t.Errorf("suggestSSHSubcommand(%q) = %q, expected %q", tc.input, got, tc.expected)
		}
	}
}

func TestFormatAliasNotFoundMessage_Proof(t *testing.T) {
	msg := formatAliasNotFoundMessage("hosts", nil)
	if !strings.Contains(msg, "SSH host alias 'hosts' not found in registry.") {
		t.Errorf("expected missing alias message in output, got:\n%s", msg)
	}
	if !strings.Contains(msg, "Did you mean: gitmap ssh nodes?") {
		t.Errorf("expected suggestion 'gitmap ssh nodes' in output, got:\n%s", msg)
	}
}
