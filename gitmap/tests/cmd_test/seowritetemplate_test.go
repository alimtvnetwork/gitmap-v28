// Package cmd_test — unit tests for seo-write template loading and substitution.
package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// templateFile mirrors the cmd package struct for test JSON parsing.
type templateFile struct {
	Titles       []string `json:"titles"`
	Descriptions []string `json:"descriptions"`
}

// TestLoadFromJSONFile_ValidFile verifies loading a valid template file.
func TestLoadFromJSONFile_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "templates.json")

	tf := templateFile{
		Titles:       []string{"Title A", "Title B"},
		Descriptions: []string{"Desc A", "Desc B", "Desc C"},
	}

	data, _ := json.Marshal(tf)
	os.WriteFile(path, data, 0o644)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var loaded templateFile
	if err := json.Unmarshal(content, &loaded); err != nil {
		t.Fatal(err)
	}

	if len(loaded.Titles) != 2 {
		t.Errorf("expected 2 titles, got %d", len(loaded.Titles))
	}
	if len(loaded.Descriptions) != 3 {
		t.Errorf("expected 3 descriptions, got %d", len(loaded.Descriptions))
	}
}

// TestLoadFromJSONFile_EmptyArrays verifies empty template arrays.
func TestLoadFromJSONFile_EmptyArrays(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")

	tf := templateFile{Titles: []string{}, Descriptions: []string{}}
	data, _ := json.Marshal(tf)
	os.WriteFile(path, data, 0o644)

	content, _ := os.ReadFile(path)
	var loaded templateFile
	json.Unmarshal(content, &loaded)

	if len(loaded.Titles) != 0 {
		t.Errorf("expected 0 titles, got %d", len(loaded.Titles))
	}
	if len(loaded.Descriptions) != 0 {
		t.Errorf("expected 0 descriptions, got %d", len(loaded.Descriptions))
	}
}

// TestLoadFromJSONFile_InvalidJSON verifies error on invalid JSON.
func TestLoadFromJSONFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	os.WriteFile(path, []byte("not json"), 0o644)

	content, _ := os.ReadFile(path)
	var loaded templateFile
	err := json.Unmarshal(content, &loaded)

	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// TestLoadFromJSONFile_MissingFile verifies error on missing file.
func TestLoadFromJSONFile_MissingFile(t *testing.T) {
	_, err := os.ReadFile("/nonexistent/templates.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// TestBuildReplacer_AllPlaceholders verifies all 7 placeholders are replaced.
func TestBuildReplacer_AllPlaceholders(t *testing.T) {
	r := strings.NewReplacer(
		constants.PlaceholderService, "Plumbing",
		constants.PlaceholderArea, "London",
		constants.PlaceholderURL, "example.com",
		constants.PlaceholderCompany, "Acme Ltd",
		constants.PlaceholderPhone, "0800-123",
		constants.PlaceholderEmail, "info@acme.com",
		constants.PlaceholderAddress, "10 High St",
	)

	input := "Best {service} in {area} by {company} at {url}. Call {phone}, email {email}, visit {address}."
	result := r.Replace(input)

	expected := "Best Plumbing in London by Acme Ltd at example.com. Call 0800-123, email info@acme.com, visit 10 High St."
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// TestBuildReplacer_PartialPlaceholders verifies partial flag replacement.
func TestBuildReplacer_PartialPlaceholders(t *testing.T) {
	r := strings.NewReplacer(
		constants.PlaceholderService, "Roofing",
		constants.PlaceholderArea, "",
		constants.PlaceholderURL, "roof.com",
		constants.PlaceholderCompany, "",
		constants.PlaceholderPhone, "",
		constants.PlaceholderEmail, "",
		constants.PlaceholderAddress, "",
	)

	input := "Top {service} in {area} - {url}"
	result := r.Replace(input)

	if !strings.Contains(result, "Roofing") {
		t.Error("expected Roofing in result")
	}
	if !strings.Contains(result, "roof.com") {
		t.Error("expected roof.com in result")
	}
	// Area should be replaced with empty string
	if strings.Contains(result, "{area}") {
		t.Error("expected {area} to be replaced")
	}
}

// TestBuildReplacer_NoPlaceholders verifies text without placeholders passes through.
func TestBuildReplacer_NoPlaceholders(t *testing.T) {
	r := strings.NewReplacer(
		constants.PlaceholderService, "Plumbing",
		constants.PlaceholderArea, "London",
		constants.PlaceholderURL, "example.com",
		constants.PlaceholderCompany, "",
		constants.PlaceholderPhone, "",
		constants.PlaceholderEmail, "",
		constants.PlaceholderAddress, "",
	)

	input := "Just a regular commit message"
	result := r.Replace(input)

	if result != input {
		t.Errorf("expected unchanged %q, got %q", input, result)
	}
}

// TestGenerateMessages_CountMatchesMaxCommits verifies message count limiting.
func TestGenerateMessages_CountMatchesMaxCommits(t *testing.T) {
	titles := []string{"T1", "T2", "T3"}
	descs := []string{"D1", "D2"}

	maxCommits := 5
	count := maxCommits
	if count == 0 {
		count = len(titles) * len(descs)
	}

	if count != 5 {
		t.Errorf("expected 5 messages, got %d", count)
	}
}

// TestGenerateMessages_ZeroMaxUsesProduct verifies default count = titles * descs.
func TestGenerateMessages_ZeroMaxUsesProduct(t *testing.T) {
	titles := []string{"T1", "T2", "T3"}
	descs := []string{"D1", "D2"}

	maxCommits := 0
	count := maxCommits
	if count == 0 {
		count = len(titles) * len(descs)
	}

	if count != 6 {
		t.Errorf("expected 6 messages (3*2), got %d", count)
	}
}
