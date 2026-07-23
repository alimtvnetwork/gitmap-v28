package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestEffectiveGofmtBudget(t *testing.T) {
	if got := effectiveGofmtBudget(fixRepoOptions{}); got != constants.FixRepoGofmtMaxCmdLen {
		t.Fatalf("default budget: want %d, got %d", constants.FixRepoGofmtMaxCmdLen, got)
	}
	if got := effectiveGofmtBudget(fixRepoOptions{gofmtMaxCmdLen: 5000}); got != 5000 {
		t.Fatalf("override budget: want 5000, got %d", got)
	}
	if got := effectiveGofmtBudget(fixRepoOptions{gofmtMaxCmdLen: 0}); got != constants.FixRepoGofmtMaxCmdLen {
		t.Fatalf("zero override should fall back, got %d", got)
	}
}

func TestBatchCmdLen(t *testing.T) {
	got := batchCmdLen([]string{"a", "bb", "ccc"})
	// 1+1 + 2+1 + 3+1 + overhead
	want := 1 + 1 + 2 + 1 + 3 + 1 + gofmtArgvOverhead
	if got != want {
		t.Fatalf("want %d, got %d", want, got)
	}
}

func TestParseFixRepoGofmtMaxCmdLen(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    int
		wantErr bool
	}{
		{"long form space", []string{"--gofmt-max-cmd-len", "5000"}, 5000, false},
		{"long form equals", []string{"--gofmt-max-cmd-len=5000"}, 5000, false},
		{"below floor errors", []string{"--gofmt-max-cmd-len", "100"}, 0, true},
		{"non-integer errors", []string{"--gofmt-max-cmd-len", "abc"}, 0, true},
		{"missing value errors", []string{"--gofmt-max-cmd-len"}, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := parseFixRepoArgs(tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got opts=%+v", opts)
				}

				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opts.gofmtMaxCmdLen != tc.want {
				t.Fatalf("want %d, got %d", tc.want, opts.gofmtMaxCmdLen)
			}
		})
	}
}

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
		paths := make([]string, 50)
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

func TestProbeChunkerSelfTest(t *testing.T) {
	r := probeChunkerSelfTest(constants.FixRepoGofmtMaxCmdLen)
	if !r.OK {
		t.Fatalf("selftest failed: %s", r.Detail)
	}
	// Small budget still safe: chunker never emits empty batches.
	r2 := probeChunkerSelfTest(constants.FixRepoGofmtMinCmdLen)
	if !r2.OK {
		t.Fatalf("selftest at min budget failed: %s", r2.Detail)
	}
}
