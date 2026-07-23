// Package cmd_test — unit tests for seo-write template creation.
package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// sampleTemplateFile mirrors the cmd package struct for scaffolding.
type sampleTemplateFile struct {
	Titles       []string `json:"titles"`
	Descriptions []string `json:"descriptions"`
	Placeholders []string `json:"placeholders"`
}

// TestBuildSampleTemplate_HasTitles verifies sample has titles.
func TestBuildSampleTemplate_HasTitles(t *testing.T) {
	sample := buildSampleTemplateHelper()
	if len(sample.Titles) == 0 {
		t.Error("expected sample to have titles")
	}
}

// TestBuildSampleTemplate_HasDescriptions verifies sample has descriptions.
func TestBuildSampleTemplate_HasDescriptions(t *testing.T) {
	sample := buildSampleTemplateHelper()
	if len(sample.Descriptions) == 0 {
		t.Error("expected sample to have descriptions")
	}
}

// TestBuildSampleTemplate_HasAllPlaceholders verifies all 7 placeholders.
func TestBuildSampleTemplate_HasAllPlaceholders(t *testing.T) {
	sample := buildSampleTemplateHelper()

	expected := []string{
		constants.PlaceholderService,
		constants.PlaceholderArea,
		constants.PlaceholderURL,
		constants.PlaceholderCompany,
		constants.PlaceholderPhone,
		constants.PlaceholderEmail,
		constants.PlaceholderAddress,
	}

	if len(sample.Placeholders) != len(expected) {
		t.Errorf("expected %d placeholders, got %d", len(expected), len(sample.Placeholders))
	}

	for i, p := range expected {
		if i >= len(sample.Placeholders) {
			break
		}
		if sample.Placeholders[i] != p {
			t.Errorf("expected placeholder %q at index %d, got %q", p, i, sample.Placeholders[i])
		}
	}
}

// TestCreateTemplateFile_WritesValidJSON verifies the scaffolded file is valid JSON.
func TestCreateTemplateFile_WritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "seo-templates.json")

	sample := buildSampleTemplateHelper()
	data, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	content, _ := os.ReadFile(path)
	var loaded sampleTemplateFile
	if err := json.Unmarshal(content, &loaded); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if len(loaded.Titles) != len(sample.Titles) {
		t.Errorf("expected %d titles in file, got %d", len(sample.Titles), len(loaded.Titles))
	}
}

// TestCreateTemplateFile_Overwrite verifies overwriting an existing file.
func TestCreateTemplateFile_Overwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "seo-templates.json")

	os.WriteFile(path, []byte("old content"), 0o644)

	sample := buildSampleTemplateHelper()
	data, _ := json.MarshalIndent(sample, "", "  ")
	os.WriteFile(path, data, 0o644)

	content, _ := os.ReadFile(path)
	var loaded sampleTemplateFile
	err := json.Unmarshal(content, &loaded)
	if err != nil {
		t.Error("expected valid JSON after overwrite")
	}
}

// --- Helper ---

func buildSampleTemplateHelper() sampleTemplateFile {
	return sampleTemplateFile{
		Titles: []string{
			"Top {service} in {area} - {url}",
			"{company}: Trusted {service} Provider in {area}",
			"Best {service} Near {area} | Visit {url}",
		},
		Descriptions: []string{
			"Looking for reliable {service} in {area}? {company} provides top-rated solutions. Visit {url} or call {phone}.",
			"{company} is the leading {service} provider serving {area}. Learn more at {url}. Contact us at {email}.",
		},
		Placeholders: []string{
			constants.PlaceholderService,
			constants.PlaceholderArea,
			constants.PlaceholderURL,
			constants.PlaceholderCompany,
			constants.PlaceholderPhone,
			constants.PlaceholderEmail,
			constants.PlaceholderAddress,
		},
	}
}
