// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
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

