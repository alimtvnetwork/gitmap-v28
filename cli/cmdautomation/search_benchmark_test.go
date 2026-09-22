package cmdautomation

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkRunSearch_Literal(b *testing.B) {
	tempDir := b.TempDir()
	createSearchBenchmarkFiles(tempDir, 100)
	opts := SearchOptions{
		Pattern: "unique_search_token_12345",
		Dir:     tempDir,
		Workers: 4,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := RunSearch(opts)
		if err != nil || res.TotalHits != 100 {
			b.Fatalf("search failed or incorrect hits: %v, hits=%d", err, res.TotalHits)
		}
	}
}

func BenchmarkRunSearch_Regex(b *testing.B) {
	tempDir := b.TempDir()
	createSearchBenchmarkFiles(tempDir, 100)
	opts := SearchOptions{
		Pattern: "unique_search_[a-z]+_[0-9]+",
		Dir:     tempDir,
		Workers: 4,
		IsRegex: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := RunSearch(opts)
		if err != nil || res.TotalHits != 100 {
			b.Fatalf("regex search failed: %v, hits=%d", err, res.TotalHits)
		}
	}
}

func TestRunSearch_Correctness(t *testing.T) {
	tempDir := t.TempDir()
	createSearchBenchmarkFiles(tempDir, 10)
	opts := SearchOptions{
		Pattern: "unique_search_token_12345",
		Dir:     tempDir,
		Workers: 2,
	}

	res, err := RunSearch(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotalHits != 10 {
		t.Fatalf("expected 10 hits, got %d", res.TotalHits)
	}

	for _, match := range res.Matches {
		if match.LineNumber != 5 {
			t.Errorf("expected line number 5, got %d for file %s", match.LineNumber, match.Path)
		}
		if match.LineContent != "const token = \"unique_search_token_12345\"" {
			t.Errorf("unexpected line content: %s", match.LineContent)
		}
	}
}

func createSearchBenchmarkFiles(dir string, count int) {
	for i := 0; i < count; i++ {
		filePath := filepath.Join(dir, fmt.Sprintf("bench_%04d.go", i))
		content := "package bench\n\nimport \"fmt\"\n\nconst token = \"unique_search_token_12345\"\n\nfunc Run() {\n\tfmt.Println(token)\n}\n"
		_ = os.WriteFile(filePath, []byte(content), 0644)
	}
}
