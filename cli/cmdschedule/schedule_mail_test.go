// Package cmdschedule manages scheduled tasks, privacy masking, and mailbox inspection.
package cmdschedule

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user@domain.com", "u***r@d***n.com"},
		{"al@ex.com", "a*@e*.com"},
		{"devops@company.org", "d***s@c***y.org"},
		{"invalid", "***"},
		{"", "***"},
		{"@domain.com", "***"},
		{"user@", "***"},
		{"user@domain", "***"},
		{"user@.com", "***"},
		{"me@domain.org", "m*@d***n.org"},
		{"bob@site.io", "b*b@s***e.io"},
		{"test@sub.example.com", "t***t@s*b.e***e.com"},
	}

	for _, tc := range tests {
		actual := MaskEmail(tc.input)
		if actual != tc.expected {
			t.Errorf("MaskEmail(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestMaskEmailList(t *testing.T) {
	inputs := []string{"user@domain.com", "al@ex.com"}
	expected := []string{"u***r@d***n.com", "a*@e*.com"}

	actual := MaskEmailList(inputs)
	if len(actual) != len(expected) {
		t.Fatalf("MaskEmailList length = %d; want %d", len(actual), len(expected))
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("MaskEmailList[%d] = %q; want %q", i, actual[i], expected[i])
		}
	}

	emptyRes := MaskEmailList(nil)
	if len(emptyRes) != 0 {
		t.Errorf("MaskEmailList(nil) returned non-empty: %v", emptyRes)
	}
}

func TestMaskString(t *testing.T) {
	if got := MaskString(""); got != "" {
		t.Errorf("MaskString(\"\") = %q; want \"\"", got)
	}
	if got := MaskString("abc"); got != "****" {
		t.Errorf("MaskString(\"abc\") = %q; want \"****\"", got)
	}
	if got := MaskString("12345678"); got != "1****8" {
		t.Errorf("MaskString(\"12345678\") = %q; want \"1****8\"", got)
	}
	if got := MaskString("secret-token-1234"); !strings.HasPrefix(got, "se****") {
		t.Errorf("MaskString long secret expected prefix se****, got %q", got)
	}
}

func TestPrivacyMasker(t *testing.T) {
	maskerEnabled := NewPrivacyMasker(true)
	if got := maskerEnabled.Mask("user@domain.com"); got != "u***r@d***n.com" {
		t.Errorf("PrivacyMasker(true) = %q; want \"u***r@d***n.com\"", got)
	}

	maskerDisabled := NewPrivacyMasker(false)
	if got := maskerDisabled.Mask("user@domain.com"); got != "user@domain.com" {
		t.Errorf("PrivacyMasker(false) = %q; want \"user@domain.com\"", got)
	}

	if !HasAt("a@b.com") {
		t.Errorf("HasAt(\"a@b.com\") expected true")
	}
	if HasAt("invalid") {
		t.Errorf("HasAt(\"invalid\") expected false")
	}
}

func TestRunScheduleMailCheck(t *testing.T) {
	if err := RunScheduleMailCheck([]string{"--all"}); err != nil {
		t.Fatalf("RunScheduleMailCheck(--all) failed: %v", err)
	}

	if err := RunScheduleMailCheck([]string{"--json"}); err != nil {
		t.Fatalf("RunScheduleMailCheck(--json) failed: %v", err)
	}

	if err := RunScheduleMailCheck([]string{"--unread-only"}); err != nil {
		t.Fatalf("RunScheduleMailCheck(--unread-only) failed: %v", err)
	}

	if err := RunScheduleMailCheck([]string{"--show-full-email"}); err != nil {
		t.Fatalf("RunScheduleMailCheck(--show-full-email) failed: %v", err)
	}

	if err := RunScheduleMailCheck([]string{"--account", "unknown-alias-999"}); err == nil {
		t.Fatalf("RunScheduleMailCheck(unknown) expected error, got nil")
	}
}

func TestRunScheduleEmailValidation(t *testing.T) {
	if err := RunScheduleEmail([]string{"--subject", "Test", "--every", "1h"}); err == nil {
		t.Errorf("RunScheduleEmail missing --to expected error, got nil")
	}

	if err := RunScheduleEmail([]string{"--to", "u@d.com", "--every", "1h"}); err == nil {
		t.Errorf("RunScheduleEmail missing --subject expected error, got nil")
	}

	if err := RunScheduleEmail([]string{"--to", "u@d.com", "--subject", "Test"}); err == nil {
		t.Errorf("RunScheduleEmail missing --every expected error, got nil")
	}

	if err := RunScheduleEmail([]string{"--help"}); err != nil {
		t.Errorf("RunScheduleEmail(--help) expected nil, got %v", err)
	}
}

func TestRunScheduleEmailSuccess(t *testing.T) {
	taskName := fmt.Sprintf("test-email-%d", time.Now().UnixNano())
	args := []string{
		taskName,
		"--to", "developer@example.com",
		"--subject", "Build Alert",
		"--every", "30m",
		"--body", "System is healthy",
	}

	if err := RunScheduleEmail(args); err != nil {
		t.Fatalf("RunScheduleEmail failed: %v", err)
	}

	db, err := store.OpenDefault()
	if err != nil {
		t.Fatalf("store.OpenDefault failed: %v", err)
	}
	defer db.Close()

	task, getErr := db.GetSchedule(taskName)
	if getErr != nil {
		t.Fatalf("GetSchedule(%q) failed: %v", taskName, getErr)
	}
	defer func() { _ = db.DeleteSchedule(taskName) }()

	if task.IntervalVal != "30m" {
		t.Errorf("task.IntervalVal = %q; want \"30m\"", task.IntervalVal)
	}
	if !task.IsEnabled {
		t.Errorf("task.IsEnabled expected true")
	}
}

func TestScheduleDispatchMailOps(t *testing.T) {
	if err := dispatchScheduleSubcommand("mail-check", []string{"--all"}); err != nil {
		t.Errorf("dispatchScheduleSubcommand(mail-check) failed: %v", err)
	}

	if err := dispatchScheduleSubcommand("watch-folder", []string{"--help"}); err != nil {
		t.Errorf("dispatchScheduleSubcommand(watch-folder) failed: %v", err)
	}

	if err := dispatchScheduleSubcommand("ui", []string{"--help"}); err != nil {
		t.Errorf("dispatchScheduleSubcommand(ui) failed: %v", err)
	}
}
