package cmdmacro

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseCatOptions(t *testing.T) {
	opts, err := ParseCatOptions([]string{"file.go", "-L", "140:220", "--num"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.TargetFile != "file.go" {
		t.Errorf("expected target 'file.go', got %q", opts.TargetFile)
	}

	if opts.StartLine != 140 || opts.EndLine != 220 {
		t.Errorf("expected range 140:220, got %d:%d", opts.StartLine, opts.EndLine)
	}

	if !opts.HasLineNumbers {
		t.Errorf("expected HasLineNumbers true")
	}
}

func TestParseCatOptionsHeadAndQuery(t *testing.T) {
	opts, err := ParseCatOptions([]string{"main.go", "--head", "25", "-q", "search_term", "-l", "10"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.HeadCount != 25 {
		t.Errorf("expected HeadCount 25, got %d", opts.HeadCount)
	}

	if opts.Query != "search_term" {
		t.Errorf("expected Query 'search_term', got %q", opts.Query)
	}

	if opts.Limit != 10 {
		t.Errorf("expected Limit 10, got %d", opts.Limit)
	}
}

func TestExecuteCatWithOptionsRange(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample.txt")

	lines := make([]string, 50)
	for i := 0; i < 50; i++ {
		lines[i] = "line content " + string(rune('A'+(i%26)))
	}
	_ = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)

	opts := CatOptions{
		TargetFile:     filePath,
		StartLine:      10,
		EndLine:        20,
		HasLineRange:   true,
		HasLineNumbers: true,
		IsRaw:          true,
	}

	if err := ExecuteCatWithOptions(opts); err != nil {
		t.Fatalf("ExecuteCatWithOptions failed: %v", err)
	}
}

func TestExecuteCatWithOptionsQuery(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample_query.txt")
	content := "first line\nsecond keyword line\nthird line\nfourth keyword line\n"
	_ = os.WriteFile(filePath, []byte(content), 0644)

	opts := CatOptions{
		TargetFile:     filePath,
		Query:          "keyword",
		Limit:          5,
		HasLineNumbers: true,
		IsRaw:          true,
	}

	if err := ExecuteCatWithOptions(opts); err != nil {
		t.Fatalf("ExecuteCatWithOptions failed: %v", err)
	}
}
