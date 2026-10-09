package cmdautofix

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkScan_AllCategories measures the fix scan hot path over a
// synthetic tree with violations. Baseline for program 265 (WS6);
// no optimization in this program.
func BenchmarkScan_AllCategories(b *testing.B) {
	tempDir := b.TempDir()
	for i := 0; i < 50; i++ {
		name := fmt.Sprintf("file%03d.txt", i)
		content := "line with trailing space   \nsecond line\n"
		if writeErr := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0600); writeErr != nil {
			b.Fatal(writeErr)
		}
	}
	opts := Options{Root: tempDir, Workers: 4}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, scanErr := Scan(opts); scanErr != nil {
			b.Fatalf("scan failed: %v", scanErr)
		}
	}
}
