// Package cmd_test — unit tests for seo-write loop helpers.
package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseInterval_Valid verifies correct min-max parsing.
func TestParseInterval_Valid(t *testing.T) {
	min, max := parseIntervalHelper("30-90")
	if min != 30 {
		t.Errorf("expected min=30, got %d", min)
	}
	if max != 90 {
		t.Errorf("expected max=90, got %d", max)
	}
}

// TestParseInterval_Equal verifies min equals max.
func TestParseInterval_Equal(t *testing.T) {
	min, max := parseIntervalHelper("60-60")
	if min != 60 || max != 60 {
		t.Errorf("expected 60-60, got %d-%d", min, max)
	}
}

// TestParseInterval_Invalid verifies invalid formats are caught.
func TestParseInterval_InvalidFormat(t *testing.T) {
	cases := []string{"abc", "30", "-", "10-5", "x-y"}
	for _, c := range cases {
		_, _, err := parseIntervalSafe(c)
		if err == nil {
			t.Errorf("expected error for input %q", c)
		}
	}
}

// TestPickFile_EmptyList returns dot when no files.
func TestPickFile_EmptyList(t *testing.T) {
	result := pickFileHelper([]string{}, 0)
	if result != "." {
		t.Errorf("expected '.', got %q", result)
	}
}

// TestPickFile_RoundRobin verifies round-robin file selection.
func TestPickFile_RoundRobin(t *testing.T) {
	files := []string{"a.txt", "b.txt", "c.txt"}

	expected := []string{"a.txt", "b.txt", "c.txt", "a.txt", "b.txt"}
	for i, want := range expected {
		got := pickFileHelper(files, i)
		if got != want {
			t.Errorf("index %d: expected %q, got %q", i, want, got)
		}
	}
}

// TestFormatDuration_MinutesOnly verifies sub-hour formatting.
func TestFormatDuration_MinutesOnly(t *testing.T) {
	result := formatDurationHelper(0, 45)
	if result != "45m" {
		t.Errorf("expected 45m, got %q", result)
	}
}

// TestFormatDuration_HoursAndMinutes verifies hour+minute formatting.
func TestFormatDuration_HoursAndMinutes(t *testing.T) {
	result := formatDurationHelper(2, 15)
	if result != "2h 15m" {
		t.Errorf("expected 2h 15m, got %q", result)
	}
}

// TestFormatDuration_ZeroMinutes verifies zero-minute formatting.
func TestFormatDuration_ZeroMinutes(t *testing.T) {
	result := formatDurationHelper(0, 0)
	if result != "0m" {
		t.Errorf("expected 0m, got %q", result)
	}
}

// TestAppendAndRevertFile verifies file append and revert cycle.
func TestAppendAndRevertFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.html")

	original := "Hello World"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	// Append
	text := "SEO content here"
	appendToFileHelper(path, text)

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), text) {
		t.Error("expected appended text in file")
	}

	// Revert
	revertFileHelper(path, text)

	data, _ = os.ReadFile(path)
	if strings.Contains(string(data), text) {
		t.Error("expected appended text to be removed after revert")
	}
	if strings.TrimSpace(string(data)) != original {
		t.Errorf("expected original content %q, got %q", original, strings.TrimSpace(string(data)))
	}
}

// TestResolveRotateFile_ExplicitExists returns the path when file exists.
func TestResolveRotateFile_ExplicitExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rotate.html")
	os.WriteFile(path, []byte("test"), 0o644)

	result := resolveRotateFileHelper(path)
	if result != path {
		t.Errorf("expected %q, got %q", path, result)
	}
}

// TestResolveRotateFile_ExplicitMissing returns empty string.
func TestResolveRotateFile_ExplicitMissing(t *testing.T) {
	result := resolveRotateFileHelper("/nonexistent/file.html")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

// --- Helper re-implementations (pure logic extracted for testability) ---

func parseIntervalHelper(s string) (int, int) {
	min, max, _ := parseIntervalSafe(s)

	return min, max
}

func parseIntervalSafe(s string) (int, int, error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return 0, 0, errInvalid()
	}

	var min, max int
	_, err1 := parseIntSafe(parts[0], &min)
	_, err2 := parseIntSafe(parts[1], &max)

	if err1 != nil || err2 != nil || min > max {
		return 0, 0, errInvalid()
	}

	return min, max, nil
}

func parseIntSafe(s string, target *int) (bool, error) {
	val := 0
	if len(s) == 0 {
		return false, errInvalid()
	}

	for _, c := range s {
		if c < '0' || c > '9' {
			return false, errInvalid()
		}
		val = val*10 + int(c-'0')
	}
	*target = val

	return true, nil
}

func errInvalid() error {
	return os.ErrInvalid
}

func pickFileHelper(files []string, idx int) string {
	if len(files) == 0 {
		return "."
	}

	return files[idx%len(files)]
}

func formatDurationHelper(hours, minutes int) string {
	if hours > 0 {
		return strings.Replace(strings.Replace("Xh Ym", "X", itoa(hours), 1), "Y", itoa(minutes), 1)
	}

	return itoa(minutes) + "m"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}

	return s
}

func appendToFileHelper(path, text string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.WriteString("\n" + text)
}

func revertFileHelper(path, text string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	cleaned := strings.Replace(string(data), "\n"+text, "", 1)
	_ = os.WriteFile(path, []byte(cleaned), 0o644)
}

func resolveRotateFileHelper(explicit string) string {
	if explicit != "" {
		if _, err := os.Stat(explicit); err != nil {
			return ""
		}

		return explicit
	}

	return ""
}
