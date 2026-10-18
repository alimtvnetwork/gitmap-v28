package cmdspace

import (
	"os"
	"strings"
	"testing"
)

func TestParseSendEmailUIArgs(t *testing.T) {
	args := []string{
		"--to", "alpha@example.com",
		"--subject", "Release 1.0",
		"--body", "Build complete",
		"--html",
		"--task", "release",
		"--dry-run",
		"--non-interactive",
	}
	opts, err := parseSendEmailUIArgs(args)
	if err != nil {
		t.Fatalf("parseSendEmailUIArgs unexpected error: %v", err)
	}
	assertUIOptionsParsed(t, opts)
}

func assertUIOptionsParsed(t *testing.T, opts sendEmailUIOptions) {
	t.Helper()
	if opts.to != "alpha@example.com" || opts.subject != "Release 1.0" {
		t.Fatalf("unexpected to/subject: %+v", opts)
	}
	if opts.body != "Build complete" || opts.task != "release" {
		t.Fatalf("unexpected body/task: %+v", opts)
	}
	if !opts.isHTML {
		t.Fatal("expected isHTML true")
	}
	if !opts.isDryRun {
		t.Fatal("expected isDryRun true")
	}
	if opts.isInteractive {
		t.Fatal("expected isInteractive false")
	}
}

func TestParseSendEmailUIArgsEqualsForm(t *testing.T) {
	args := []string{"--to=user@example.com", "--subject=Hello", "-n"}
	opts, err := parseSendEmailUIArgs(args)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if opts.to != "user@example.com" || opts.subject != "Hello" {
		t.Fatalf("unexpected parsed to/subject: %+v", opts)
	}
	if !opts.isDryRun {
		t.Fatal("expected isDryRun true")
	}
}

func TestParseSendEmailUIArgsUnknownFlag(t *testing.T) {
	_, err := parseSendEmailUIArgs([]string{"--bogus"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestRenderEmailUIBanner(t *testing.T) {
	banner := RenderEmailUIBanner()
	if !strings.Contains(banner, "GitMap Space") {
		t.Fatalf("banner missing GitMap Space: %s", banner)
	}
	if !strings.Contains(banner, "Email & Trigger Studio") {
		t.Fatalf("banner missing Email & Trigger Studio: %s", banner)
	}
}

func TestRenderAccountTableWithAccounts(t *testing.T) {
	accounts := []EmailAccountInfo{
		{
			Alias: "DevOps", Email: "dev@company.com",
			ImapHost: "imap.company.com", ImapPort: 993, ImapTLS: true,
			SmtpHost: "smtp.company.com", SmtpPort: 465, SmtpTLSMode: "implicit",
		},
	}
	rendered := RenderAccountTable(accounts)
	if !strings.Contains(rendered, "DevOps") {
		t.Fatalf("rendered table missing DevOps: %s", rendered)
	}
	if !strings.Contains(rendered, "dev@company.com") {
		t.Fatalf("rendered table missing dev@company.com: %s", rendered)
	}
	if !strings.Contains(rendered, "TLS=true") {
		t.Fatalf("rendered table missing TLS=true: %s", rendered)
	}
	if !strings.Contains(rendered, "implicit") {
		t.Fatalf("rendered table missing implicit: %s", rendered)
	}
}

func TestRenderAccountTableEmpty(t *testing.T) {
	rendered := RenderAccountTable(nil)
	if !strings.Contains(rendered, "No accounts stored") {
		t.Fatalf("expected 'No accounts stored' row in empty table, got: %s", rendered)
	}
}

func TestRenderNotifierLine(t *testing.T) {
	withVal := RenderNotifierLine("notifier@example.com")
	if !strings.Contains(withVal, "notifier@example.com") {
		t.Fatalf("expected notifier address in output: %s", withVal)
	}
	withoutVal := RenderNotifierLine("")
	if !strings.Contains(withoutVal, "none") {
		t.Fatalf("expected 'none' when notifier empty: %s", withoutVal)
	}
}

func TestRenderDraftPreview(t *testing.T) {
	opts := sendEmailUIOptions{
		to: "admin@corp.internal", subject: "Nightly backup",
		body: "Status OK", task: "backup", isHTML: true,
	}
	preview := RenderDraftPreview(opts)
	assertDraftPreviewContains(t, preview)
}

func assertDraftPreviewContains(t *testing.T, preview string) {
	t.Helper()
	for _, expected := range []string{"admin@corp.internal", "Nightly backup", "Status OK", "HTML", "backup"} {
		if !strings.Contains(preview, expected) {
			t.Fatalf("preview missing %q: %s", expected, preview)
		}
	}
}

func TestBuildSendEmailArgsFromUI(t *testing.T) {
	opts := sendEmailUIOptions{
		to: "test@example.com", subject: "Ping", body: "Pong", isHTML: true,
	}
	argv := buildSendEmailArgsFromUI(opts)
	assertArgvContains(t, argv, "--to", "test@example.com")
	assertArgvContains(t, argv, "--subject", "Ping")
	assertArgvContains(t, argv, "--body", "Pong")
	if !hasArg(argv, "--html") {
		t.Fatal("expected --html in argv")
	}
}

func assertArgvContains(t *testing.T, argv []string, flag, expectedVal string) {
	t.Helper()
	for i, arg := range argv {
		if arg == flag && i+1 < len(argv) && argv[i+1] == expectedVal {
			return
		}
	}
	t.Fatalf("argv missing %s %s: %+v", flag, expectedVal, argv)
}

func hasArg(argv []string, target string) bool {
	for _, a := range argv {
		if a == target {
			return true
		}
	}
	return false
}

func TestRunSpaceSendEmailUI_DryRun(t *testing.T) {
	args := []string{
		"--dry-run", "--non-interactive",
		"--to", "user@test.org",
		"--subject", "CI Finished",
		"--body", "All checks green",
		"--task", "release",
	}
	err := RunSpaceSendEmailUI(args)
	if err != nil {
		t.Fatalf("RunSpaceSendEmailUI dry-run returned error: %v", err)
	}
}

func TestRunSpaceSendEmailUI_Help(t *testing.T) {
	err := RunSpaceSendEmailUI([]string{"--help"})
	if err != nil {
		t.Fatalf("RunSpaceSendEmailUI --help returned error: %v", err)
	}
}

func TestLoadEmailAccountsAndNotifierFallback(t *testing.T) {
	accounts := LoadEmailAccounts()
	if len(accounts) > 0 {
		first := accounts[0]
		if len(first.Email) == 0 {
			t.Fatalf("expected non-empty email for first account: %+v", first)
		}
	}
	notifier := LoadDefaultNotifier()
	if len(notifier) == 0 {
		t.Log("notifier is empty or unconfigured")
	}
}

func TestDispatchSpaceEmailUIAndSendEmailUI(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"gitmap", "space", "email", "ui", "--help"}
	handled, err := DispatchSpace("space")
	if err != nil {
		t.Fatalf("DispatchSpace('space') email ui --help error: %v", err)
	}
	if !handled {
		t.Fatal("expected handled true for email ui --help")
	}

	os.Args = []string{"gitmap", "space", "send-email", "--ui", "--help"}
	handledSend, errSend := DispatchSpace("space")
	if errSend != nil {
		t.Fatalf("DispatchSpace('space') send-email --ui --help error: %v", errSend)
	}
	if !handledSend {
		t.Fatal("expected handled true for send-email --ui --help")
	}
}
