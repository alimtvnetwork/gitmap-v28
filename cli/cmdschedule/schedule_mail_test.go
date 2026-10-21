// Package cmdschedule manages scheduled tasks, privacy masking, and mailbox inspection.
package cmdschedule

import (
	"fmt"
	"os"
	"path/filepath"
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
		{"a@domain.com", "a*@d***n.com"},
		{"user@x.org", "u***r@x*.org"},
		{"a@x.org", "a*@x*.org"},
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
	testPrivacyMaskerBasics(t)
	testPrivacyMaskerListAndEnv(t)
}

func testPrivacyMaskerBasics(t *testing.T) {
	maskerEnabled := NewPrivacyMasker(true)
	if got := maskerEnabled.Mask("user@domain.com"); got != "u***r@d***n.com" {
		t.Errorf("PrivacyMasker(true) = %q; want \"u***r@d***n.com\"", got)
	}

	maskerDisabled := NewPrivacyMasker(false)
	if got := maskerDisabled.Mask("user@domain.com"); got != "user@domain.com" {
		t.Errorf("PrivacyMasker(false) = %q; want \"user@domain.com\"", got)
	}

	if !HasAt("a@b.com") || HasAt("invalid") {
		t.Errorf("HasAt validation failed")
	}
}

func testPrivacyMaskerListAndEnv(t *testing.T) {
	masker := NewPrivacyMasker(true)
	list := masker.MaskList([]string{"a@domain.com", "user@x.org"})
	if len(list) != 2 || list[0] != "a*@d***n.com" || list[1] != "u***r@x*.org" {
		t.Errorf("MaskList masked failed: %v", list)
	}

	os.Setenv("GITMAP_ALLOW_UNMASKED_PII", "1")
	if ResolveIsMasked(false) {
		t.Errorf("ResolveIsMasked expected false with GITMAP_ALLOW_UNMASKED_PII=1")
	}
	os.Unsetenv("GITMAP_ALLOW_UNMASKED_PII")
}

func TestRunScheduleMailCheck(t *testing.T) {
	flags := [][]string{
		{"-a"}, {"-u"}, {"-j"}, {"-m", "primary"}, {"--account=primary"},
		{"--all"}, {"--json"}, {"--unread-only"}, {"--show-full-email"},
	}
	for _, f := range flags {
		if err := RunScheduleMailCheck(f); err != nil {
			t.Errorf("RunScheduleMailCheck(%v) failed: %v", f, err)
		}
	}
	if err := RunScheduleMailCheck([]string{"--account", "unknown-alias-999"}); err == nil {
		t.Fatalf("RunScheduleMailCheck(unknown) expected error, got nil")
	}
}

func TestRunScheduleMailCheck_EmptyTable(t *testing.T) {
	renderMailCheckTable([]MailboxCheckResult{})
	if err := RunScheduleMailCheck([]string{"-m", "alerts", "-u"}); err != nil {
		t.Errorf("RunScheduleMailCheck empty query failed: %v", err)
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

	if err := RunScheduleEmail([]string{"--to", "u@d.com", "--subject", "Test", "--every", "invalid_val"}); err == nil {
		t.Errorf("RunScheduleEmail invalid interval expected error, got nil")
	}

	if err := RunScheduleEmail([]string{"--to", "u@d.com", "--subject", "Test", "--every", "-5m"}); err == nil {
		t.Errorf("RunScheduleEmail negative interval expected error, got nil")
	}

	if err := RunScheduleEmail([]string{"--help"}); err != nil {
		t.Errorf("RunScheduleEmail(--help) expected nil, got %v", err)
	}
}

func TestRunScheduleEmailShorthandFlags(t *testing.T) {
	name := fmt.Sprintf("test-short-%d", time.Now().UnixNano())
	args := []string{name, "-t", "developer@example.com", "-s", "Short", "-b", "Body", "-i", "15m", "-j"}
	if err := RunScheduleEmail(args); err != nil {
		t.Fatalf("RunScheduleEmail shorthands failed: %v", err)
	}

	db, _ := store.OpenDefault()
	defer db.Close()
	defer func() { _ = db.DeleteSchedule(name) }()

	task, err := db.GetSchedule(name)
	if err != nil || task.IntervalVal != "15m" {
		t.Errorf("expected interval 15m, got %s (err: %v)", task.IntervalVal, err)
	}
}

func TestRunScheduleEmailFileFlag(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "body.txt")
	_ = os.WriteFile(tmpFile, []byte("file payload"), 0644)
	name := fmt.Sprintf("test-f-%d", time.Now().UnixNano())
	args := []string{name, "-t", "a@b.com", "-s", "File", "-f", tmpFile, "-i", "2h"}
	if err := RunScheduleEmail(args); err != nil {
		t.Fatalf("RunScheduleEmail -f failed: %v", err)
	}

	db, _ := store.OpenDefault()
	defer db.Close()
	_ = db.DeleteSchedule(name)
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
