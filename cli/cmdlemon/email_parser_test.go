package cmdlemon

import (
	"testing"
	"time"
)

func TestParseEmailCommand_Formats(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantCmd     string
		wantArgsLen int
		wantFormat  EmailFormatType
		wantErr     bool
	}{
		{"plain cmd", "cmd: status", "status", 0, EmailFormatTypeText, false},
		{"plain command args", "command: pull --force", "pull", 1, EmailFormatTypeText, false},
		{"plain gitmap prefix", "gitmap run --interval 8s", "run", 2, EmailFormatTypeText, false},
		{"plain slash status", "/status", "status", 0, EmailFormatTypeText, false},
		{"plain slash agm", "/agm test 123", "agm", 2, EmailFormatTypeText, false},
		{"plain slash lap", "/lap", "lap", 0, EmailFormatTypeText, false},
		{"plain slash rp", "/rp", "rp", 0, EmailFormatTypeText, false},
		{"raw json", `{"command": "status", "args": ["-v"], "requestId": "r1"}`, "status", 1, EmailFormatTypeJSON, false},
		{"fenced json", "```json\n{\"command\": \"pull\"}\n```", "pull", 0, EmailFormatTypeJSON, false},
		{"md header", "# Command: status --json", "status", 1, EmailFormatTypeYAML, false},
		{"fenced yaml", "```yaml\ncommand: pull\nargs: [\"--all\"]\n```", "pull", 1, EmailFormatTypeYAML, false},
		{"empty body", "", "", 0, EmailFormatTypeUnknown, true},
		{"unrecognized", "hello world", "", 0, EmailFormatTypeUnknown, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := EmailMessage{From: "admin@example.com", Body: tc.body, ReceivedAt: time.Now()}
			env, err := ParseEmailCommand(msg)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %s", tc.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if env.Command != tc.wantCmd {
				t.Errorf("command = %s, want %s", env.Command, tc.wantCmd)
			}
			if len(env.Args) != tc.wantArgsLen {
				t.Errorf("args count = %d, want %d", len(env.Args), tc.wantArgsLen)
			}
			if env.Format != tc.wantFormat {
				t.Errorf("format = %s, want %s", env.Format, tc.wantFormat)
			}
			if env.Sender != "admin@example.com" {
				t.Errorf("sender = %s, want admin@example.com", env.Sender)
			}
			if len(env.RequestID) == 0 {
				t.Error("expected non-empty requestId")
			}
		})
	}
}

func TestParseEmailCommand_MetadataOverride(t *testing.T) {
	body := `{"command": "status", "requestId": "custom-id", "sender": "bot@example.com"}`
	msg := EmailMessage{From: "admin@example.com", Body: body}
	env, err := ParseEmailCommand(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.RequestID != "custom-id" {
		t.Errorf("requestId = %s, want custom-id", env.RequestID)
	}
	if env.Sender != "bot@example.com" {
		t.Errorf("sender = %s, want bot@example.com", env.Sender)
	}
}
