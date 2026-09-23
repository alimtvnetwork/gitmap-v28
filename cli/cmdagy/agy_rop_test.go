package cmdagy

import (
	"testing"
	"time"
)

func TestResolveRepoSlug(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		expected string
	}{
		{"Antigravity-Manager", "d:/work/antigravity-manager", "antigravity-manager"},
		{"GitMap Core", "c:/repos/gitmap", "gitmap-core"},
		{"", "d:/work/my-project", "my-project"},
		{"", "", "default"},
	}

	for _, c := range cases {
		actual := ResolveRepoSlug(c.name, c.path)
		if actual != c.expected {
			t.Errorf("for (%s, %s): expected %s, got %s", c.name, c.path, c.expected, actual)
		}
	}
}

func TestAGYBackupDB_Lifecycle(t *testing.T) {
	slug := "test-repo-slug"
	db, err := OpenAGYBackupDB(slug)
	if err != nil {
		t.Fatalf("failed to open agy backup db: %v", err)
	}
	defer db.Close()

	backup := AGYConversationBackup{
		ConversationID: "conv-12345",
		ProjectSlug:    slug,
		ProjectPath:    "/path/to/test",
		Title:          "Test Conversation",
		StepCount:      15,
		TranscriptJSON: `{"step": 1, "content": "hello"}`,
		CreatedAt:      time.Now().UTC(),
	}

	if err := BackupConversation(db, backup); err != nil {
		t.Fatalf("failed to save conversation backup: %v", err)
	}

	// Verify query
	var readTitle string
	var readStepCount int
	err = db.QueryRow("SELECT Title, StepCount FROM AGYConversationBackup WHERE ConversationID = ?", "conv-12345").Scan(&readTitle, &readStepCount)
	if err != nil {
		t.Fatalf("failed to read back saved conversation: %v", err)
	}
	if readTitle != "Test Conversation" || readStepCount != 15 {
		t.Errorf("mismatched read values: title=%s, steps=%d", readTitle, readStepCount)
	}

	// Close db before purge to release file handle on Windows
	_ = db.Close()

	// Test purge
	count, _, err := ClearAGYBackups(slug)
	if err != nil {
		t.Fatalf("failed to clear agy backups: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 dir purged, got %d", count)
	}
}

func TestParseROPArgs(t *testing.T) {
	n, isDry := parseROPArgs([]string{"3", "--dry-run"})
	if n != 3 {
		t.Errorf("expected n=3, got %d", n)
	}
	if !isDry {
		t.Error("expected isDry true")
	}

	nDef, isDryDef := parseROPArgs([]string{})
	if nDef != 5 {
		t.Errorf("expected default n=5, got %d", nDef)
	}
	if isDryDef {
		t.Error("expected default isDry false")
	}
}
