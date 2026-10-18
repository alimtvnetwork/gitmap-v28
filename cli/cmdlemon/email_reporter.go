package cmdlemon

import (
	"fmt"
	"strings"
	"time"
)

// EmailReporter formats and transmits command completion reports.
type EmailReporter struct {
	transmitter EmailTransmitter
}

// NewEmailReporter constructs a new completion report dispatcher.
func NewEmailReporter(transmitter EmailTransmitter) *EmailReporter {
	return &EmailReporter{transmitter: transmitter}
}

// SendCompletionReport sends the formatted execution report to recipient.
func (r *EmailReporter) SendCompletionReport(to string, outcome ExecutionOutcome) error {
	if r.transmitter == nil {
		return fmt.Errorf("nil email transmitter")
	}

	subject := FormatReportSubject(outcome)
	body := FormatReportBody(outcome)

	return r.transmitter.Transmit(to, subject, body)
}

// FormatReportSubject produces status indicator subject line.
func FormatReportSubject(outcome ExecutionOutcome) string {
	statusTag := resolveStatusTag(outcome)

	return fmt.Sprintf("[%s] Command '%s' (Req: %s)", statusTag, outcome.Command, outcome.RequestID)
}

func resolveStatusTag(outcome ExecutionOutcome) string {
	if outcome.IsSuccess {
		return "SUCCESS"
	}

	return fmt.Sprintf("FAILED (Exit %d)", outcome.ExitCode)
}

// FormatReportBody generates a structured markdown-compatible execution summary.
func FormatReportBody(outcome ExecutionOutcome) string {
	completedStr := outcome.CompletedAt.Format(time.RFC3339)
	if outcome.CompletedAt.IsZero() {
		completedStr = time.Now().UTC().Format(time.RFC3339)
	}

	stdoutPreview := TruncatePreview(outcome.Stdout, 25)
	stderrPreview := TruncatePreview(outcome.Stderr, 25)

	return fmt.Sprintf("Command: %s\nRequest ID: %s\nExit Code: %d\nDuration: %v\nCompleted: %s\n\n--- Standard Output ---\n%s\n\n--- Standard Error ---\n%s\n",
		outcome.Command, outcome.RequestID, outcome.ExitCode, outcome.Duration, completedStr, stdoutPreview, stderrPreview)
}

// TruncatePreview restricts output text to maxLines for email delivery readability.
func TruncatePreview(text string, maxLines int) string {
	clean := strings.TrimSpace(text)
	if len(clean) == 0 {
		return "(empty)"
	}

	lines := strings.Split(clean, "\n")
	if len(lines) <= maxLines {
		return clean
	}

	preview := strings.Join(lines[:maxLines], "\n")

	return fmt.Sprintf("%s\n... [truncated %d lines]", preview, len(lines)-maxLines)
}
