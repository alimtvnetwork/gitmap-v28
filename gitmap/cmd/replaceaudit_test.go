package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestBuildAuditNeedles verifies the dual-form contract for the audit
// scanner: every target version produces exactly two needles
// (`<base>-vN` and `<base>/vN`) in deterministic order. Expected
// needles are derived from the targets slice (NOT hard-coded) so a
// width-crossing bump (v9 -> v10/v12) cannot silently desync the
// fixture from the code under test.
func TestBuildAuditNeedles(t *testing.T) {
	targets := []int{4, 12}
	got := buildAuditNeedles("gitmap", targets)
	want := make([][]byte, 0, len(targets)*2)
	for _, n := range targets {
		want = append(want,
			[]byte(fmt.Sprintf("gitmap-v%d", n)),
			[]byte(fmt.Sprintf("gitmap/v%d", n)),
		)
	}
	// Diagnostic log: print both slices so any future failure here
	// surfaces the full got/want context even when CI truncates the
	// per-assertion error line.
	t.Logf("buildAuditNeedles(gitmap, %v) -> got=%q want=%q", targets, got, want)
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (got=%q want=%q)", len(got), len(want), got, want)
	}
	for i := range got {
		if !bytes.Equal(got[i], want[i]) {
			t.Errorf("needle[%d] = %q, want %q (targets=%v)", i, got[i], want[i], targets)
		}
	}
}

// TestBuildAuditNeedlesWidthCrossing exercises the v9 -> v10/v12 width
// boundary directly. Audit needles are OLD tokens, so the boundary
// affects targets[i] (not current). Regression guard against the
// historical desync where fixtures hard-coded `gitmap-v27` next to
// a single-digit input.
//
// IMPORTANT: every needle string is derived from the same `targets`
// int slice via fmt.Sprintf so fix-repo cannot rewrite one half of
// the pair without the other. See .lovable/memory/issues/2026-05-02-
// fixrepo-paired-literal-desync.md.
func TestBuildAuditNeedlesWidthCrossing(t *testing.T) {
	targets := []int{8, 9, 10, 12}
	got := buildAuditNeedles("gitmap", targets)
	want := buildExpectedNeedles("gitmap", targets)
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if string(got[i]) != want[i] {
			t.Errorf("needle[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

// buildExpectedNeedles derives the expected needle slice from the
// same int slice the production code consumes — the single source of
// truth that prevents paired-literal desync after a version bump.
func buildExpectedNeedles(base string, targets []int) []string {
	out := make([]string, 0, len(targets)*2)
	for _, n := range targets {
		out = append(out,
			fmt.Sprintf("%s-v%d", base, n),
			fmt.Sprintf("%s/v%d", base, n))
	}

	return out
}

// TestLineContainsAny is the primitive the audit scanner uses to decide
// whether a line is reportable. Tested independently so future
// optimizations cannot regress the contract.
func TestLineContainsAny(t *testing.T) {
	needles := [][]byte{[]byte("gitmap-v27"), []byte("gitmap/v9")}
	if !lineContainsAny([]byte("import gitmap-v27/foo"), needles) {
		t.Error("expected dash-form match")
	}
	if !lineContainsAny([]byte("module github.com/x/gitmap/v9"), needles) {
		t.Error("expected slash-form match")
	}
	if lineContainsAny([]byte("nothing relevant"), needles) {
		t.Error("unexpected match on clean line")
	}
}

// TestScanAuditFileFormatting locks the printed format to the spec's
// `path:line: matched-text` shape and confirms the per-file hit count.
func TestScanAuditFileFormatting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	body := "intro line\nsee gitmap-v27 for details\nclean line\nuse gitmap/v9 too\n"
	mustWriteFile(t, path, []byte(body))

	needles := [][]byte{[]byte("gitmap-v27"), []byte("gitmap/v9")}

	stdout, hits := captureStdout(t, func() int {
		return scanAuditFile(path, needles)
	})

	if hits != 2 {
		t.Fatalf("hits = %d, want 2", hits)
	}

	want2 := fmt.Sprintf(constants.MsgReplaceAuditMatch, path, 2, "see gitmap-v27 for details")
	want4 := fmt.Sprintf(constants.MsgReplaceAuditMatch, path, 4, "use gitmap/v9 too")
	if !strings.Contains(stdout, want2) {
		t.Errorf("missing line-2 hit\n got: %q\nwant substring: %q", stdout, want2)
	}
	if !strings.Contains(stdout, want4) {
		t.Errorf("missing line-4 hit\n got: %q\nwant substring: %q", stdout, want4)
	}
}

// captureStdout temporarily redirects os.Stdout so we can assert on
// fmt.Fprintf(os.Stdout, ...) output without coupling tests to a
// global writer.
func captureStdout(t *testing.T, fn func() int) (string, int) {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	result := fn()

	w.Close()
	os.Stdout = orig

	buf, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return string(buf), result
}
