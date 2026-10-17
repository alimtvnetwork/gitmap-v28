// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// openTestTriggerDBPath creates a migrated temp-file settings DB and
// returns its path. File-backed (not in-memory) because the command
// under test closes its handle after every call — each step opens a
// fresh handle, mirroring production open/use/close. Never touches the
// real user DB.
func openTestTriggerDBPath(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "triggers.db")
	db, err := store.OpenAt(path)
	if err != nil {
		t.Fatalf("open temp db: %v", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("migrate temp db: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("close temp db: %v", err)
	}

	return path
}

// stubTriggerStore points openTriggerStore at the temp file DB for the
// test's duration, opening a fresh handle per call.
func stubTriggerStore(t *testing.T, path string) {
	t.Helper()

	original := openTriggerStore
	openTriggerStore = func() (*store.DB, error) { return store.OpenAt(path) }
	t.Cleanup(func() { openTriggerStore = original })
}

// loadTestTriggers reads the trigger list through a fresh handle.
func loadTestTriggers(t *testing.T, path string) []EmailTrigger {
	t.Helper()

	db, err := store.OpenAt(path)
	if err != nil {
		t.Fatalf("open temp db: %v", err)
	}
	defer db.Close()

	triggers, err := LoadEmailTriggers(db)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	return triggers
}

// saveTestTrigger registers one trigger through a fresh handle.
func saveTestTrigger(t *testing.T, path string, trigger EmailTrigger) {
	t.Helper()

	db, err := store.OpenAt(path)
	if err != nil {
		t.Fatalf("open temp db: %v", err)
	}
	defer db.Close()

	if err := RegisterEmailTrigger(db, trigger); err != nil {
		t.Fatalf("register: %v", err)
	}
}

func TestParseEmailTaskTypeValid(t *testing.T) {
	cases := map[string]EmailTaskType{
		"release": EmailTaskRelease,
		"pull":    EmailTaskPull,
		"clone":   EmailTaskClone,
		"scan":    EmailTaskScan,
		"push":    EmailTaskPush,
		"backup":  EmailTaskBackup,
		"RELEASE": EmailTaskRelease,
		" Pull ":  EmailTaskPull,
	}

	for input, want := range cases {
		got, err := ParseEmailTaskType(input)
		if err != nil {
			t.Fatalf("ParseEmailTaskType(%q): %v", input, err)
		}

		if got != want {
			t.Fatalf("ParseEmailTaskType(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestParseEmailTaskTypeInvalidListsValidValues(t *testing.T) {
	_, err := ParseEmailTaskType("nope")
	if err == nil {
		t.Fatal("expected error for unknown task, got nil")
	}

	for _, task := range AllEmailTaskTypes() {
		if !strings.Contains(err.Error(), string(task)) {
			t.Fatalf("error %q does not list valid value %q", err.Error(), task)
		}
	}
}

func TestParseScheduleSendEmailArgsMissing(t *testing.T) {
	missingTo := []string{"--task", "release"}
	if err := RunScheduleSendEmail(missingTo); err == nil {
		t.Fatal("expected RunScheduleSendEmail error for missing --to, got nil")
	}

	missingTask := []string{"--to", "me@example.com"}
	if err := RunScheduleSendEmail(missingTask); err == nil {
		t.Fatal("expected RunScheduleSendEmail error for missing --task, got nil")
	}
}

func TestParseScheduleSendEmailArgsUnknownFlag(t *testing.T) {
	if err := RunScheduleSendEmail([]string{"--to", "a@b.c", "--task", "pull", "--bogus"}); err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestParseScheduleSendEmailArgsForms(t *testing.T) {
	flags, err := parseScheduleSendEmailArgs([]string{"--to=me@example.com", "--task", "pull", "--name", "nightly"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if flags.to != "me@example.com" || flags.task != "pull" || flags.name != "nightly" {
		t.Fatalf("unexpected flags: %+v", flags)
	}

	listFlags, err := parseScheduleSendEmailArgs([]string{"--list"})
	if err != nil || !listFlags.isList {
		t.Fatalf("expected --list to parse, got %+v / %v", listFlags, err)
	}

	removeFlags, err := parseScheduleSendEmailArgs([]string{"--remove", "--name", "nightly"})
	if err != nil || !removeFlags.isRemove || removeFlags.name != "nightly" {
		t.Fatalf("expected --remove --name to parse, got %+v / %v", removeFlags, err)
	}
}

func TestTriggerRegisterListRemoveRoundTrip(t *testing.T) {
	path := openTestTriggerDBPath(t)
	stubTriggerStore(t, path)

	registerArgs := []string{"--to", "me@example.com", "--task", "release", "--name", "rel-alert"}
	if err := RunScheduleSendEmail(registerArgs); err != nil {
		t.Fatalf("register: %v", err)
	}

	triggers := loadTestTriggers(t, path)
	if len(triggers) != 1 || triggers[0].To != "me@example.com" || triggers[0].Task != EmailTaskRelease {
		t.Fatalf("unexpected triggers: %+v", triggers)
	}

	if err := RunScheduleSendEmail([]string{"--list"}); err != nil {
		t.Fatalf("list: %v", err)
	}

	if err := RunScheduleSendEmail([]string{"--remove", "--name", "rel-alert"}); err != nil {
		t.Fatalf("remove: %v", err)
	}

	if remaining := loadTestTriggers(t, path); len(remaining) != 0 {
		t.Fatalf("expected no triggers after remove, got %+v", remaining)
	}
}

func TestTriggerRegisterDefaultsNameToTask(t *testing.T) {
	path := openTestTriggerDBPath(t)
	stubTriggerStore(t, path)

	if err := RunScheduleSendEmail([]string{"--to", "me@example.com", "--task", "pull"}); err != nil {
		t.Fatalf("register: %v", err)
	}

	triggers := loadTestTriggers(t, path)
	if len(triggers) != 1 || triggers[0].Name != "pull" {
		t.Fatalf("expected default name \"pull\", got %+v", triggers)
	}
}

func TestTriggerRegisterInvalidTask(t *testing.T) {
	path := openTestTriggerDBPath(t)
	stubTriggerStore(t, path)

	err := RunScheduleSendEmail([]string{"--to", "me@example.com", "--task", "nope"})
	if err == nil {
		t.Fatal("expected error for invalid task, got nil")
	}
}

func TestRemoveMissingTrigger(t *testing.T) {
	path := openTestTriggerDBPath(t)
	stubTriggerStore(t, path)

	if err := RunScheduleSendEmail([]string{"--remove", "--name", "ghost"}); err == nil {
		t.Fatal("expected error removing unknown trigger, got nil")
	}

	if err := RunScheduleSendEmail([]string{"--remove"}); err == nil {
		t.Fatal("expected error for --remove without --name, got nil")
	}
}

func TestBuildTriggerEmailHTMLEscapesAndRows(t *testing.T) {
	info := MachineInfo{Host: "h<ost>", User: "u&ser", IP: "10.0.0.1"}
	outcome := TaskOutcome{
		Repo: "<script>alert(1)</script>", Branch: "main", Operation: "release",
		IsSuccess: true, Summary: "ok & done", StartedAt: "s", FinishedAt: "f",
	}

	html := BuildTriggerEmailHTML(EmailTaskRelease, outcome, info)

	for _, want := range []string{
		"&lt;script&gt;alert(1)&lt;/script&gt;", "h&lt;ost&gt;", "u&amp;ser",
		"10.0.0.1", "main", "release", "ok &amp; done", "s", "f",
		"Host", "User", "Outbound IP", "Repo", "Branch", "Operation", "Outcome", "Started", "Finished",
		bannerOKColor, statusOKText,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("html missing %q", want)
		}
	}

	if strings.Contains(html, "<script>") {
		t.Fatal("html contains unescaped script tag")
	}
}

func TestBuildTriggerEmailHTMLFailureBanner(t *testing.T) {
	info := CollectMachineInfo()
	outcome := InferTaskOutcome(EmailTaskRelease, errors.New("boom"))

	html := BuildTriggerEmailHTML(EmailTaskRelease, outcome, info)

	if !strings.Contains(html, bannerFailColor) || !strings.Contains(html, statusFailText) {
		t.Fatal("expected red FAILED banner")
	}

	if !strings.Contains(html, "boom") {
		t.Fatal("expected outcome summary in html")
	}
}

func TestInferTaskOutcome(t *testing.T) {
	ok := InferTaskOutcome(EmailTaskPull, nil)
	if !ok.IsSuccess || ok.Summary != triggerSuccessSummary || ok.Operation != "pull" {
		t.Fatalf("unexpected success outcome: %+v", ok)
	}

	failed := InferTaskOutcome(EmailTaskClone, errors.New("network down"))
	if failed.IsSuccess || failed.Summary != "network down" {
		t.Fatalf("unexpected failure outcome: %+v", failed)
	}
}

func TestFireEmailTaskTriggers(t *testing.T) {
	path := openTestTriggerDBPath(t)
	stubTriggerStore(t, path)

	saveTestTrigger(t, path, EmailTrigger{Name: "rel-alert", To: "me@example.com", Task: EmailTaskRelease})
	saveTestTrigger(t, path, EmailTrigger{Name: "pull-alert", To: "other@example.com", Task: EmailTaskPull})

	var gotTo, gotSubject, gotHTML string
	var callCount int
	original := sendTriggerEmail
	sendTriggerEmail = func(to, subject, htmlBody string) error {
		callCount++
		gotTo, gotSubject, gotHTML = to, subject, htmlBody

		return nil
	}
	t.Cleanup(func() { sendTriggerEmail = original })

	outcome := TaskOutcome{
		Repo: "demo", Branch: "main", Operation: "release",
		IsSuccess: true, Summary: "v1.2.3 shipped", StartedAt: "s", FinishedAt: "f",
	}
	FireEmailTaskTriggers(EmailTaskRelease, outcome)

	if callCount != 1 {
		t.Fatalf("expected 1 send, got %d", callCount)
	}

	if gotTo != "me@example.com" {
		t.Fatalf("expected to me@example.com, got %q", gotTo)
	}

	if !strings.Contains(gotSubject, "release") || !strings.Contains(gotHTML, "v1.2.3 shipped") {
		t.Fatalf("subject/html missing outcome: %q / %.80s", gotSubject, gotHTML)
	}
}

func TestFireEmailTaskTriggersSendFailureIsSilent(t *testing.T) {
	path := openTestTriggerDBPath(t)
	stubTriggerStore(t, path)

	saveTestTrigger(t, path, EmailTrigger{Name: "x", To: "a@b.c", Task: EmailTaskScan})

	original := sendTriggerEmail
	sendTriggerEmail = func(to, subject, htmlBody string) error { return errors.New("smtp down") }
	t.Cleanup(func() { sendTriggerEmail = original })

	FireEmailTaskTriggers(EmailTaskScan, InferTaskOutcome(EmailTaskScan, nil))
}
