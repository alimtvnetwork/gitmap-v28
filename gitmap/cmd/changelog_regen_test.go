package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegenChangelogOrdersDescending(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("v1.2.0.json", `{"version":"1.2.0","tag":"v1.2.0","branch":"release/v1.2.0"}`)
	write("v1.10.0.json", `{"version":"1.10.0","tag":"v1.10.0","branch":"release/v1.10.0"}`)
	write("latest.json", `{"version":"1.10.0","tag":"v1.10.0","branch":"release/v1.10.0"}`)

	var buf bytes.Buffer
	if err := RegenChangelog(dir, &buf); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	i10 := strings.Index(got, "## v1.10.0")
	i12 := strings.Index(got, "## v1.2.0")
	if i10 < 0 || i12 < 0 || i10 > i12 {
		t.Fatalf("expected v1.10.0 before v1.2.0:\n%s", got)
	}
	if strings.Contains(got, "latest.json") {
		t.Fatal("latest.json must be skipped")
	}
}

func TestCompareSemverDesc(t *testing.T) {
	if !compareSemverDesc("1.10.0", "1.2.0") {
		t.Fatal("1.10.0 should sort before 1.2.0")
	}
	if compareSemverDesc("1.0.0", "1.0.0") {
		t.Fatal("equal versions must not sort")
	}
}
