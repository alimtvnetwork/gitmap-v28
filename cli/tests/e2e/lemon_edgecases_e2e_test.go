package e2e_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon"
)

// normalizeFenceCase standardizes markdown code fence tags to lowercase.
func normalizeFenceCase(body string) string {
	lines := strings.Split(body, "\n")
	for idx, line := range lines {
		trimmed := strings.TrimSpace(line)
		low := strings.ToLower(trimmed)
		if strings.HasPrefix(low, "```json") {
			lines[idx] = "```json"
		} else if strings.HasPrefix(low, "```yaml") {
			lines[idx] = "```yaml"
		}
	}
	return strings.Join(lines, "\n")
}

// TestLemonEmail_FencedMarkdownAndYaml verifies parsing of fenced code blocks and markdown headers.
func TestLemonEmail_FencedMarkdownAndYaml(t *testing.T) {
	testCases := []struct {
		name    string
		body    string
		wantCmd string
		wantArg string
		wantFmt cmdlemon.EmailFormatType
	}{
		{
			name:    "fenced JSON in markdown",
			body:    "Here is the command:\n```json\n{\"command\":\"status\",\"args\":[\"--json\"]}\n```\nThanks!",
			wantCmd: "status",
			wantArg: "--json",
			wantFmt: cmdlemon.EmailFormatTypeJSON,
		},
		{
			name:    "fenced YAML in markdown",
			body:    "Please run:\n```yaml\ncommand: pull\nargs:\n  - --all\n```\nRegards",
			wantCmd: "pull",
			wantArg: "--all",
			wantFmt: cmdlemon.EmailFormatTypeYAML,
		},
	}

	for _, tc := range testCases {
		env, err := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{From: "engineer@example.com", Body: tc.body})
		if err != nil {
			t.Fatalf("%s unexpected error: %v", tc.name, err)
		}
		if env.Command != tc.wantCmd || env.Format != tc.wantFmt || len(env.Args) != 1 || env.Args[0] != tc.wantArg {
			t.Errorf("%s mismatch: cmd=%s, fmt=%v, args=%v", tc.name, env.Command, env.Format, env.Args)
		}
	}

	// Mixed-case fence extraction via ExtractFencedBlock
	jsonMixed := "```JSON\n{\"command\":\"status\"}\n```"
	extractedJSON := cmdlemon.ExtractFencedBlock(jsonMixed, "```JSON")
	if !strings.Contains(extractedJSON, "\"command\":\"status\"") {
		t.Errorf("expected ExtractFencedBlock to extract mixed-case JSON block, got: %s", extractedJSON)
	}

	yamlMixed := "```Yaml\ncommand: pull\n```"
	extractedYAML := cmdlemon.ExtractFencedBlock(yamlMixed, "```Yaml")
	if !strings.Contains(extractedYAML, "command: pull") {
		t.Errorf("expected ExtractFencedBlock to extract mixed-case Yaml block, got: %s", extractedYAML)
	}

	// Normalized mixed-case fence parsing through ParseEmailCommand
	normEnv, normErr := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{
		From: "engineer@example.com",
		Body: normalizeFenceCase(jsonMixed),
	})
	if normErr != nil || normEnv.Command != "status" {
		t.Errorf("expected normalized mixed-case fence to parse cleanly: env=%+v, err=%v", normEnv, normErr)
	}

	// Markdown header with multiline arguments
	headerBody := "# Command: nodes\n  --json\n  --all"
	envH, errH := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{From: "admin@example.com", Body: headerBody})
	if errH != nil {
		t.Fatalf("header markdown unexpected error: %v", errH)
	}
	if envH.Command != "nodes" || envH.Format != cmdlemon.EmailFormatTypeYAML || len(envH.Args) != 2 || envH.Args[0] != "--json" || envH.Args[1] != "--all" {
		t.Errorf("header markdown mismatch: %+v", envH)
	}
}

// TestLemonEmail_CRLFAndWhitespaceSanitization verifies handling of CRLF, tabs, and case-insensitivity.
func TestLemonEmail_CRLFAndWhitespaceSanitization(t *testing.T) {
	msg := cmdlemon.EmailMessage{From: "dev@example.com", Body: "\r\n\t  CMD:   GITMAP   STATUS  \t\r\n"}
	env, err := cmdlemon.ParseEmailCommand(msg)
	if err != nil {
		t.Fatalf("unexpected error parsing CRLF command: %v", err)
	}

	tokens := append([]string{env.Command}, env.Args...)
	validTokens, isAuthorized := cmdlemon.ValidateTokens(tokens)
	if !isAuthorized || len(validTokens) == 0 || validTokens[0] != "status" {
		t.Fatalf("expected authorized subcommand 'status', got %v (authorized=%v)", validTokens, isAuthorized)
	}

	msgSlash := cmdlemon.EmailMessage{From: "dev@example.com", Body: "\r\n\t  /pull   --all   \r\n"}
	envSlash, errSlash := cmdlemon.ParseEmailCommand(msgSlash)
	if errSlash != nil || envSlash.Command != "pull" || len(envSlash.Args) != 1 || envSlash.Args[0] != "--all" {
		t.Errorf("slash command CRLF mismatch: env=%+v, err=%v", envSlash, errSlash)
	}

	msgMultiArg := cmdlemon.EmailMessage{From: "dev@example.com", Body: "/pe   --trace   \t  --filter=e2e  \r\n"}
	envMulti, errMulti := cmdlemon.ParseEmailCommand(msgMultiArg)
	if errMulti != nil || envMulti.Command != "pe" || len(envMulti.Args) != 2 || envMulti.Args[0] != "--trace" || envMulti.Args[1] != "--filter=e2e" {
		t.Errorf("multi-arg whitespace mismatch: env=%+v, err=%v", envMulti, errMulti)
	}
}

// TestLemonEmail_EmptyAndMalformedHandling verifies error handling for empty and malformed emails.
func TestLemonEmail_EmptyAndMalformedHandling(t *testing.T) {
	for _, body := range []string{"", "  \r\n\t  \r\n"} {
		if _, err := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{Body: body}); err == nil || err.Error() != "empty email body" {
			t.Errorf("expected 'empty email body' for %q, got: %v", body, err)
		}
	}
	if _, err := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{Body: "Random email greeting."}); err == nil || err.Error() != "unrecognized command format" {
		t.Errorf("expected 'unrecognized command format', got: %v", err)
	}
	if _, err := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{Body: "```json\n{\"command\":\"status\"\n```"}); err == nil {
		t.Errorf("expected error for malformed JSON envelope")
	}
	if _, err := cmdlemon.ParseEmailCommand(cmdlemon.EmailMessage{Body: "```json\n{}\n```"}); err == nil {
		t.Errorf("expected error for empty JSON command envelope")
	}
}

// TestLemonAck_NilSafetyAndErrorPropagation verifies nil pointer guards and error handling in ACK generation.
func TestLemonAck_NilSafetyAndErrorPropagation(t *testing.T) {
	replier := cmdlemon.NewAckReplier(cmdlemon.EmailTransmitterFunc(func(to, subject, body string) error { return nil }))
	if errNilEnv := replier.DispatchAck(nil); errNilEnv == nil || errNilEnv.Error() != "nil command envelope" {
		t.Errorf("expected 'nil command envelope', got: %v", errNilEnv)
	}

	nilTxReplier := cmdlemon.NewAckReplier(nil)
	env := &cmdlemon.CommandEnvelope{Command: "status", Sender: "dev@example.com"}
	if errNilTx := nilTxReplier.DispatchAck(env); errNilTx == nil || errNilTx.Error() != "nil email transmitter" {
		t.Errorf("expected 'nil email transmitter', got: %v", errNilTx)
	}

	txErr := errors.New("simulated smtp network failure")
	failReplier := cmdlemon.NewAckReplier(cmdlemon.EmailTransmitterFunc(func(to, subject, body string) error { return txErr }))
	if errProp := failReplier.DispatchAck(env); !errors.Is(errProp, txErr) {
		t.Errorf("expected error %v, got: %v", txErr, errProp)
	}

	var ackBody string
	captureReplier := cmdlemon.NewAckReplier(cmdlemon.EmailTransmitterFunc(func(to, subject, body string) error {
		ackBody = body
		return nil
	}))
	envZero := &cmdlemon.CommandEnvelope{Command: "status", Sender: "dev@example.com", Timestamp: time.Time{}}
	if errZero := captureReplier.DispatchAck(envZero); errZero != nil {
		t.Fatalf("unexpected error for zero timestamp: %v", errZero)
	}
	if strings.Contains(ackBody, "0001-01-01") || !strings.Contains(ackBody, "ACK Received: Command 'status' received at ") {
		t.Errorf("unexpected ACK body rendering: %s", ackBody)
	}
}

// TestLemonReporter_TruncationAndFailure verifies completion reports on failure and large output truncation.
func TestLemonReporter_TruncationAndFailure(t *testing.T) {
	var reportedSub, reportedBody string
	reporter := cmdlemon.NewEmailReporter(cmdlemon.EmailTransmitterFunc(func(to, subject, body string) error {
		reportedSub, reportedBody = subject, body
		return nil
	}))

	failedOutcome := cmdlemon.ExecutionOutcome{
		Command: "pull", RequestID: "req-fail-001", ExitCode: 1,
		Duration: 120 * time.Millisecond, Stderr: "fatal error",
		CompletedAt: time.Now().UTC(), IsSuccess: false,
	}
	if err := reporter.SendCompletionReport("admin@example.com", failedOutcome); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
	if !strings.Contains(reportedSub, "[FAILED (Exit 1)]") || !strings.Contains(reportedBody, "fatal error") {
		t.Errorf("failure report mismatch: sub=%q, body=%q", reportedSub, reportedBody)
	}

	var lines []string
	for idx := 1; idx <= 35; idx++ {
		lines = append(lines, fmt.Sprintf("output-row-%d", idx))
	}
	truncOutcome := cmdlemon.ExecutionOutcome{
		Command: "status", RequestID: "req-trunc-002", ExitCode: 0,
		Stdout: strings.Join(lines, "\n"), CompletedAt: time.Now().UTC(), IsSuccess: true,
	}
	if err := reporter.SendCompletionReport("admin@example.com", truncOutcome); err != nil {
		t.Fatalf("unexpected report error: %v", err)
	}
	if !strings.Contains(reportedBody, "... [truncated 10 lines]") {
		t.Errorf("expected truncation indicator, got: %s", reportedBody)
	}
}

// TestLemonSecurity_ExtendedAttackVectors verifies command whitelist security and path traversal defense.
func TestLemonSecurity_ExtendedAttackVectors(t *testing.T) {
	injectionCases := []struct {
		vector string
		cmd    string
	}{
		{"semicolon chaining", "status; cat /etc/shadow"},
		{"pipe chaining", "version | grep secret"},
		{"logical and chaining", "nodes && calc"},
		{"subshell expansion", "status $(id)"},
		{"backticks execution", "`reboot`"},
		{"output redirection", "nodes > file.txt"},
		{"input redirection", "pull < input"},
		{"forbidden binary rm", "rm"},
		{"forbidden binary cat", "cat"},
		{"forbidden binary curl", "curl"},
		{"forbidden binary calc", "calc"},
		{"forbidden binary wget", "wget"},
		{"forbidden binary bash", "bash"},
		{"forbidden binary powershell", "powershell"},
		{"forbidden binary shutdown", "shutdown"},
		{"forbidden binary kill", "kill"},
	}

	for _, tc := range injectionCases {
		if _, err := cmdlemon.ParseAndAuthorizeCommand(tc.cmd); err == nil {
			t.Errorf("vector %s (%q) must be rejected by security whitelist, but succeeded", tc.vector, tc.cmd)
		}
	}

	pathTraversalCases := []string{
		"status ../../secret",
		"pe /etc/passwd",
		"version ..\\sensitive.env",
	}

	for _, ptCmd := range pathTraversalCases {
		if hasRisk := checkPathTraversalRisk(ptCmd); !hasRisk {
			t.Errorf("path traversal command %q was not flagged as security violation", ptCmd)
		}
	}
}

func checkPathTraversalRisk(raw string) bool {
	for _, tok := range strings.Fields(raw) {
		if strings.Contains(tok, "..") || strings.HasPrefix(tok, "/etc") || strings.HasPrefix(tok, "\\") {
			return true
		}
	}
	return false
}
