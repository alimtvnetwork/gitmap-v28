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
	if err != nil || mock.sendCount != 1 || mock.lastOpts.Subject != "Build Succeeded" {
		t.Fatalf("send failed: err=%v count=%d", err, mock.sendCount)
	}
	if !strings.Contains(string(mock.lastRaw), "Subject: Build Succeeded") {
		t.Errorf("raw message missing Subject header")
	}
}

func TestMailSend_FileAndHTML(t *testing.T) {
	mock := &mockTransmitter{}
	CurrentTransmitter = mock
	tmpFile := filepath.Join(t.TempDir(), "body.html")
	_ = os.WriteFile(tmpFile, []byte("<h1>Report</h1>"), 0644)
	err := RunMail([]string{"send", "--to", "dev@example.com", "--subject", "HTML Report", "--file", tmpFile, "--html"})
	if err != nil || mock.sendCount != 1 || !mock.lastOpts.IsHTML || mock.lastOpts.Body != "<h1>Report</h1>" {
		t.Errorf("send --file --html failed: err=%v count=%d isHTML=%t", err, mock.sendCount, mock.lastOpts.IsHTML)
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
		input, expected string
	}{
		{"user.name@example.com", "u***e@e***e.com"},
		{"devops@company.org", "d***s@c***y.org"},
		{"me@i.org", "m***@i***.org"},
		{"single", "s***e"},
		{"", "***"},
	}
	for _, c := range cases {
		if actual := MaskEmail(c.input); actual != c.expected {
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
	if err != nil || summary.TotalAccounts == 0 {
		t.Fatalf("eval mailboxes failed: err=%v total=%d", err, summary.TotalAccounts)
	}
	table := FormatMailCheckTable(summary, false)
	if !strings.Contains(table, "MAILBOX ACCOUNT") {
		t.Error("expected table header in output")
	}
}

func TestMailTemplate_CRUD(t *testing.T) {
	store, _ := OpenMailStore(":memory:")
	defer store.Close()
	tmpl := MailTemplateRecord{TemplateID: "t1", TemplateName: "release-alert", SubjectTemplate: "Release v{{.Version}} Ready", BodyTemplate: "New release ready.", IsHTML: false}
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

func TestMailStore_RulesAndLogsCRUD(t *testing.T) {
	s, _ := OpenMailStore(":memory:")
	defer s.Close()
	rule := FolderWatcherRuleRecord{RuleID: "r1", RuleName: "watch-src", FolderPath: "./src", IsEnabled: true}
	_ = s.SaveWatcherRule(rule)
	rules, _ := s.ListWatcherRules()
	if len(rules) != 1 || rules[0].RuleName != "watch-src" || !rules[0].IsEnabled {
		t.Errorf("watcher rule CRUD failed: %v", rules)
	}
	log := NotificationLogRecord{LogID: "l1", RuleID: "r1", Subject: "Alert", IsDelivered: true}
	_ = s.RecordNotificationLog(log)
	logs, _ := s.ListNotificationLogs(10)
	if len(logs) != 1 || logs[0].Subject != "Alert" || !logs[0].IsDelivered {
		t.Errorf("notification log CRUD failed: %v", logs)
	}
	prof := MailboxProfileRecord{ProfileID: "p1", Alias: "sec", EmailAddress: "s@corp.io", IsActive: true, IsMasked: true}
	_ = s.SaveMailboxProfile(prof)
	profs, _ := s.ListMailboxProfiles()
	if len(profs) != 1 || profs[0].Alias != "sec" || !profs[0].IsActive || !profs[0].IsMasked {
		t.Errorf("mailbox profile CRUD failed: %v", profs)
	}
}

func TestIntervalGuard_Clamping(t *testing.T) {
	g1 := NewIntervalGuard(2 * time.Second)
	if !g1.WasClamped || g1.EffectiveInterval != MinPollingInterval {
		t.Errorf("expected 2s clamped to 5s, got %v", g1.EffectiveInterval)
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
	store, _ := OpenMailStore(":memory:")
	defer store.Close()
	mock := &mockTransmitter{}
	baseline, _ := ScanFolderTree(tempDir, true)
	testFile := filepath.Join(tempDir, "mutation.log")
	_ = os.WriteFile(testFile, []byte("log entry 1\n"), 0644)
	cfg := FolderWatcherConfig{RuleID: "rule-test", TargetDirectory: tempDir, Recipients: []string{"notifier@test.org"}, IsRecursive: true, RunOnce: true}
	batch, _, err := ExecuteWatchCycle(cfg, baseline, store, mock)
	if err != nil || batch.TotalEvents != 1 || batch.Events[0].EventType != EventCreate || mock.sendCount != 1 {
		t.Fatalf("single pass detect failed: err=%v events=%d count=%d", err, batch.TotalEvents, mock.sendCount)
	}
}

func TestFolderWatcher_ZeroMutations(t *testing.T) {
	dir := t.TempDir()
	mock := &mockTransmitter{}
	base, _ := ScanFolderTree(dir, true)
	cfg := FolderWatcherConfig{TargetDirectory: dir, RunOnce: true}
	batch, _, _ := ExecuteWatchCycle(cfg, base, nil, mock)
	if batch.HasEvents || mock.sendCount != 0 {
		t.Errorf("expected 0 events, got hasEvents=%t count=%d", batch.HasEvents, mock.sendCount)
	}
}

func TestFolderWatcher_CompoundMutations(t *testing.T) {
	prev := map[string]FileItemSnapshot{
		"m.txt": {RelativePath: "m.txt", SizeBytes: 10, ModTimeUnix: 100},
		"d.txt": {RelativePath: "d.txt", SizeBytes: 20, ModTimeUnix: 100},
	}
	curr := map[string]FileItemSnapshot{
		"m.txt": {RelativePath: "m.txt", SizeBytes: 25, ModTimeUnix: 200},
		"c.txt": {RelativePath: "c.txt", SizeBytes: 30, ModTimeUnix: 200},
	}
	events := DetectFileMutations(prev, curr)
	verifyCompoundEvents(t, events)
}

func verifyCompoundEvents(t *testing.T, events []FileEventRecord) {
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	hasMod, hasDel, hasCreate := false, false, false
	for _, e := range events {
		switch e.EventType {
		case EventModify:
			hasMod = e.Path == "m.txt"
		case EventDelete:
			hasDel = e.Path == "d.txt"
		case EventCreate:
			hasCreate = e.Path == "c.txt"
		}
	}
	if !hasMod || !hasDel || !hasCreate {
		t.Errorf("compound events missing: %+v", events)
	}
}

func TestMailCLI_TemplateRouting(t *testing.T) {
	_ = RunMail([]string{"template", "add", "ci-build", "--subject", "CI Alert", "--body", "Failed", "--html"})
	_ = RunMail([]string{"template", "ls"})
	_ = RunMail([]string{"template", "view", "ci-build"})
	_ = RunMail([]string{"template", "rm", "ci-build"})
	if err := RunMail([]string{"template", "unknown"}); err == nil {
		t.Error("expected error for unknown template subcommand")
	}
}

func TestMailCLI_NotifyRouting(t *testing.T) {
	tempDir := t.TempDir()
	_ = RunMail([]string{"notify", "watch", "--path", tempDir, "--once"})
	_ = RunMail([]string{"notify", "ls"})
	_ = RunMail([]string{"notify", "trigger", "sample-rule"})
	if err := RunMail([]string{"notify", "unknown"}); err == nil {
		t.Error("expected error for unknown notify subcommand")
	}
}
