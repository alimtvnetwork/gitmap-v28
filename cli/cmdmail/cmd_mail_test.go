package cmdmail

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type mockTransmitter struct {
	lastOpts  MailSendOptions
	lastRaw   []byte
	sendCount int
	failErr   error
}

func (m *mockTransmitter) SendMail(opts MailSendOptions, rawMessage []byte) error {
	m.lastOpts = opts
	m.lastRaw = rawMessage
	m.sendCount++
	return m.failErr
}

func TestMailSend_Success(t *testing.T) {
	mock := &mockTransmitter{}
	CurrentTransmitter = mock

	err := RunMail([]string{"send", "--to", "dev@example.com", "--subject", "Build Succeeded", "--body", "All tests passed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.sendCount != 1 {
		t.Fatalf("expected 1 send call, got %d", mock.sendCount)
	}
	if mock.lastOpts.Subject != "Build Succeeded" {
		t.Errorf("expected subject 'Build Succeeded', got '%s'", mock.lastOpts.Subject)
	}
	if !strings.Contains(string(mock.lastRaw), "Subject: Build Succeeded") {
		t.Errorf("raw message missing Subject header")
	}
}

func TestMailSend_BCCStripped(t *testing.T) {
	opts := MailSendOptions{
		From:    "sender@company.com",
		To:      []string{"user@company.com"},
		Bcc:     []string{"hidden@secret.org"},
		Subject: "Private Report",
		Body:    "Secret content",
	}
	raw := string(BuildMIMEMessage(opts))
	if strings.Contains(raw, "hidden@secret.org") {
		t.Errorf("BCC leaked into MIME headers: %s", raw)
	}
}

func TestMailSend_Validation(t *testing.T) {
	mock := &mockTransmitter{}
	CurrentTransmitter = mock

	if err := RunMail([]string{"send", "--subject", "No Recipient"}); err == nil {
		t.Error("expected validation error for missing --to")
	}
	if err := RunMail([]string{"send", "--to", "user@test.org"}); err == nil {
		t.Error("expected validation error for missing --subject")
	}
	_ = RunMail([]string{"send", "--to", "user@test.org", "--subject", "Dry", "--dry-run"})
	if mock.sendCount != 0 {
		t.Errorf("dry-run should not call send, got %d sends", mock.sendCount)
	}
}

func TestPrivacyMasking(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"user.name@example.com", "u***e@e***e.com"},
		{"devops@company.org", "d***s@c***y.org"},
		{"me@i.org", "m***@i***.org"},
		{"single", "s***e"},
		{"", "***"},
	}
	for _, c := range cases {
		actual := MaskEmail(c.input)
		if actual != c.expected {
			t.Errorf("MaskEmail(%q) = %q, expected %q", c.input, actual, c.expected)
		}
	}
	csvOut := MaskEmailCSV("a@b.com, dev@corp.io")
	if !strings.Contains(csvOut, "a***@b***.com") || !strings.Contains(csvOut, "d***v@c***p.io") {
		t.Errorf("MaskEmailCSV failed: %s", csvOut)
	}
}

func TestMailCheck_Evaluation(t *testing.T) {
	store, err := OpenMailStore(":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory store: %v", err)
	}
	defer store.Close()

	summary, err := EvaluateMailboxes(store, MailboxCheckFilter{CheckAll: true})
	if err != nil {
		t.Fatalf("unexpected error evaluating mailboxes: %v", err)
	}
	if summary.TotalAccounts == 0 {
		t.Error("expected at least 1 account evaluated")
	}
	table := FormatMailCheckTable(summary, false)
	if !strings.Contains(table, "MAILBOX ACCOUNT") {
		t.Error("expected table header in output")
	}
}

func TestMailTemplate_CRUD(t *testing.T) {
	store, err := OpenMailStore(":memory:")
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	defer store.Close()

	tmpl := MailTemplateRecord{
		TemplateID:      "t1",
		TemplateName:    "release-alert",
		SubjectTemplate: "Release v{{.Version}} Ready",
		BodyTemplate:    "The new release is available.",
		IsHTML:          false,
	}
	if err := store.SaveTemplate(tmpl); err != nil {
		t.Fatalf("failed to save template: %v", err)
	}
	fetched, err := store.GetTemplate("release-alert")
	if err != nil || fetched.SubjectTemplate != tmpl.SubjectTemplate {
		t.Fatalf("failed to retrieve saved template: %v", err)
	}
	list, _ := store.ListTemplates()
	if len(list) != 1 {
		t.Errorf("expected 1 template in list, got %d", len(list))
	}
	_ = store.DeleteTemplate("release-alert")
	if _, err := store.GetTemplate("release-alert"); err == nil {
		t.Error("expected error after deleting template")
	}
}

func TestIntervalGuard_Clamping(t *testing.T) {
	g1 := NewIntervalGuard(2 * time.Second)
	if !g1.WasClamped || g1.EffectiveInterval != MinPollingInterval {
		t.Errorf("expected 2s clamped to 5s, got %v (clamped: %t)", g1.EffectiveInterval, g1.WasClamped)
	}
	g2 := NewIntervalGuard(10 * time.Second)
	if g2.WasClamped || g2.EffectiveInterval != 10*time.Second {
		t.Errorf("expected 10s preserved, got %v", g2.EffectiveInterval)
	}
	g3 := NewIntervalGuard(0)
	if g3.WasClamped || g3.EffectiveInterval != DefaultPollingInterval {
		t.Errorf("expected 0 to default to 8s, got %v", g3.EffectiveInterval)
	}
}

func TestFolderWatcher_SinglePassDetect(t *testing.T) {
	tempDir := t.TempDir()
	store, err := OpenMailStore(":memory:")
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	defer store.Close()

	mock := &mockTransmitter{}
	baseline, _ := ScanFolderTree(tempDir, true)

	// Create a new file
	testFile := filepath.Join(tempDir, "mutation.log")
	if err := os.WriteFile(testFile, []byte("log entry 1\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cfg := FolderWatcherConfig{
		RuleID:          "rule-test",
		TargetDirectory: tempDir,
		Recipients:      []string{"notifier@test.org"},
		IsRecursive:     true,
		RunOnce:         true,
	}

	batch, _, err := ExecuteWatchCycle(cfg, baseline, store, mock)
	if err != nil {
		t.Fatalf("unexpected watch cycle error: %v", err)
	}
	if batch.TotalEvents != 1 {
		t.Fatalf("expected 1 event, got %d", batch.TotalEvents)
	}
	if batch.Events[0].EventType != EventCreate {
		t.Errorf("expected EventCreate, got %s", batch.Events[0].EventType)
	}
	if mock.sendCount != 1 {
		t.Errorf("expected 1 notification email sent, got %d", mock.sendCount)
	}
}
