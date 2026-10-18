package cmdlemon

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type mockTransmitter struct {
	to      string
	subject string
	body    string
	err     error
}

func (m *mockTransmitter) Transmit(to, subject, body string) error {
	m.to = to
	m.subject = subject
	m.body = body

	return m.err
}

func TestAckReplier_DispatchAck(t *testing.T) {
	mock := &mockTransmitter{}
	replier := NewAckReplier(mock)
	fixedTime := time.Date(2026, 10, 19, 10, 0, 0, 0, time.UTC)
	env := &CommandEnvelope{
		Command:   "status",
		Args:      []string{"--json"},
		RequestID: "req-456",
		Sender:    "dev@example.com",
		Timestamp: fixedTime,
	}

	if err := replier.DispatchAck(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mock.to != "dev@example.com" {
		t.Errorf("to = %s, want dev@example.com", mock.to)
	}

	expectedBody := "ACK Received: Command 'status --json' received at 2026-10-19T10:00:00Z (Request ID: req-456). Executing..."
	if mock.body != expectedBody {
		t.Errorf("body = %q, want %q", mock.body, expectedBody)
	}

	if !strings.Contains(mock.subject, "status") {
		t.Errorf("subject %q missing status", mock.subject)
	}
}

func TestAckReplier_Errors(t *testing.T) {
	replier := NewAckReplier(nil)
	if err := replier.DispatchAck(&CommandEnvelope{}); err == nil {
		t.Error("expected error for nil transmitter")
	}

	mockErr := &mockTransmitter{err: errors.New("smtp connection failed")}
	replierErr := NewAckReplier(mockErr)
	if err := replierErr.DispatchAck(&CommandEnvelope{Command: "pull"}); err == nil {
		t.Error("expected propagated error from transmitter")
	}
}

func TestEmailReporter_SendCompletionReport(t *testing.T) {
	mock := &mockTransmitter{}
	reporter := NewEmailReporter(mock)
	outcome := ExecutionOutcome{
		Command:     "status",
		RequestID:   "req-789",
		ExitCode:    0,
		Duration:    150 * time.Millisecond,
		Stdout:      "all repositories clean",
		Stderr:      "",
		CompletedAt: time.Now().UTC(),
		IsSuccess:   true,
	}

	if err := reporter.SendCompletionReport("dev@example.com", outcome); err != nil {
		t.Fatalf("unexpected reporter error: %v", err)
	}

	if !strings.Contains(mock.subject, "[SUCCESS]") {
		t.Errorf("subject %q missing [SUCCESS]", mock.subject)
	}

	if !strings.Contains(mock.body, "all repositories clean") {
		t.Errorf("body %q missing stdout output", mock.body)
	}
}
