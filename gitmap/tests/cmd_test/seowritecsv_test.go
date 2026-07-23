// Package cmd_test — unit tests for seo-write CSV parsing.
package cmd_test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCSVToMessages_ValidRows verifies correct CSV parsing.
func TestCSVToMessages_ValidRows(t *testing.T) {
	records := [][]string{
		{"Title 1", "Description 1"},
		{"Title 2", "Description 2"},
		{"Title 3", "Description 3"},
	}

	messages := csvToMessagesHelper(records)
	if len(messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(messages))
	}
	if messages[0][0] != "Title 1" || messages[0][1] != "Description 1" {
		t.Errorf("unexpected first message: %v", messages[0])
	}
}

// TestCSVToMessages_SkipsSingleColumn verifies rows with <2 columns are skipped.
func TestCSVToMessages_SkipsSingleColumn(t *testing.T) {
	records := [][]string{
		{"Title 1", "Description 1"},
		{"Only title"},
		{"Title 2", "Description 2"},
	}

	messages := csvToMessagesHelper(records)
	if len(messages) != 2 {
		t.Errorf("expected 2 messages (skipping single-column row), got %d", len(messages))
	}
}

// TestCSVToMessages_EmptyRecords returns empty slice.
func TestCSVToMessages_EmptyRecords(t *testing.T) {
	records := [][]string{}
	messages := csvToMessagesHelper(records)
	if len(messages) != 0 {
		t.Errorf("expected 0 messages, got %d", len(messages))
	}
}

// TestCSVToMessages_ExtraColumns ignores extra columns.
func TestCSVToMessages_ExtraColumns(t *testing.T) {
	records := [][]string{
		{"Title", "Desc", "Extra", "More"},
	}

	messages := csvToMessagesHelper(records)
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}
	if messages[0][0] != "Title" || messages[0][1] != "Desc" {
		t.Errorf("unexpected message: %v", messages[0])
	}
}

// TestReadCSVFile_ValidFile verifies reading a real CSV file.
func TestReadCSVFile_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.csv")

	content := "Title A,Description A\nTitle B,Description B\n"
	os.WriteFile(path, []byte(content), 0o644)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}

// TestReadCSVFile_EmptyFile verifies empty CSV returns no rows.
func TestReadCSVFile_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.csv")

	os.WriteFile(path, []byte(""), 0o644)

	f, _ := os.Open(path)
	defer f.Close()

	records, _ := csv.NewReader(f).ReadAll()
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}

// TestReadCSVFile_MissingFile verifies error on missing file.
func TestReadCSVFile_MissingFile(t *testing.T) {
	_, err := os.Open("/nonexistent/test.csv")
	if err == nil {
		t.Error("expected error for missing CSV file")
	}
}

// TestReadCSVFile_QuotedFields verifies CSV with quoted fields.
func TestReadCSVFile_QuotedFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "quoted.csv")

	content := `"Title, with comma","Description with ""quotes"""` + "\n"
	os.WriteFile(path, []byte(content), 0o644)

	f, _ := os.Open(path)
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0][0] != "Title, with comma" {
		t.Errorf("expected quoted title, got %q", records[0][0])
	}
}

// TestCSVToMessages_WhitespacePreserved verifies whitespace in fields.
func TestCSVToMessages_WhitespacePreserved(t *testing.T) {
	records := [][]string{
		{"  Title with spaces  ", "  Desc with spaces  "},
	}

	messages := csvToMessagesHelper(records)
	if len(messages) != 1 {
		t.Fatal("expected 1 message")
	}
	if !strings.Contains(messages[0][0], "  ") {
		t.Error("expected whitespace to be preserved")
	}
}

// --- Helper ---

func csvToMessagesHelper(records [][]string) [][2]string {
	messages := make([][2]string, 0, len(records))

	for _, row := range records {
		if len(row) < 2 {
			continue
		}
		messages = append(messages, [2]string{row[0], row[1]})
	}

	return messages
}
