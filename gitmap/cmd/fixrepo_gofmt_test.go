package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestChunkPathsForGofmt(t *testing.T) {
	t.Run("empty input yields nil", func(t *testing.T) {
		if got := chunkPathsForGofmt(nil, 100); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("all fits in one chunk", func(t *testing.T) {
		paths := []string{"a.go", "b.go", "c.go"}
		got := chunkPathsForGofmt(paths, 1000)
		if len(got) != 1 || len(got[0]) != 3 {
			t.Fatalf("expected 1 chunk of 3, got %v", got)
		}
	})

	t.Run("overflow splits into multiple chunks", func(t *testing.T) {
		long := strings.Repeat("x", 100)
		paths := make([]string, 50) // 50 * ~101 = 5050 bytes
		for i := range paths {
			paths[i] = long
		}
		got := chunkPathsForGofmt(paths, 500)
		if len(got) < 2 {
			t.Fatalf("expected multiple batches, got %d", len(got))
		}
		total := 0
		for _, b := range got {
			batchLen := 0
			for _, p := range b {
				batchLen += len(p) + 1
			}
			// Allow single-path overflow (documented behavior).
			if len(b) > 1 && batchLen > 500 {
				t.Fatalf("batch exceeds budget: %d > 500", batchLen)
			}
			total += len(b)
		}
		if total != 50 {
			t.Fatalf("expected 50 total paths preserved, got %d", total)
		}
	})

	t.Run("single path over budget emitted alone", func(t *testing.T) {
		huge := strings.Repeat("z", 2000)
		got := chunkPathsForGofmt([]string{huge, "small.go"}, 500)
		if len(got) != 2 {
			t.Fatalf("expected 2 batches, got %d", len(got))
		}
		if len(got[0]) != 1 || got[0][0] != huge {
			t.Fatalf("expected huge path alone in first batch, got %v", got[0])
		}
	})

	t.Run("zero or negative budget falls back to default", func(t *testing.T) {
		paths := []string{"a.go", "b.go"}
		got := chunkPathsForGofmt(paths, 0)
		if len(got) != 1 {
			t.Fatalf("expected 1 chunk under default budget, got %d", len(got))
		}
		if constants.FixRepoGofmtMaxCmdLen < 1000 {
			t.Fatalf("default budget suspiciously small: %d", constants.FixRepoGofmtMaxCmdLen)
		}
	})
}
