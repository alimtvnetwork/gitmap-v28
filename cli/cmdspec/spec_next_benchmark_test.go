package cmdspec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// BenchmarkNextSpecCandidate measures spec-number issuance (disk scan +
// DB max) over a synthetic repo root. Baseline for program 265 (WS6);
// no optimization in this program.
func BenchmarkNextSpecCandidate(b *testing.B) {
	tempDir := b.TempDir()
	specRoot := filepath.Join(tempDir, "02-spec", "21-app")
	if mkdirErr := os.MkdirAll(specRoot, 0755); mkdirErr != nil {
		b.Fatal(mkdirErr)
	}
	for _, n := range []string{"260-foo", "261-bar", "262-baz"} {
		if mkdirErr := os.MkdirAll(filepath.Join(specRoot, n), 0755); mkdirErr != nil {
			b.Fatal(mkdirErr)
		}
	}
	conn, openErr := store.InitMasterAgentDB(filepath.Join(tempDir, "agent.db"))
	if openErr != nil {
		b.Fatal(openErr)
	}
	defer conn.Close()
	if schemaErr := store.EnsureSpecNumbersSchema(conn); schemaErr != nil {
		b.Fatal(schemaErr)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, candErr := nextSpecCandidate(conn, tempDir); candErr != nil {
			b.Fatalf("nextSpecCandidate failed: %v", candErr)
		}
	}
}
