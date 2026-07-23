package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestPrettyFixtures is a table-driven loop over every paired
// testdata/pretty/*.in.md / *.want.txt fixture. Adding a new edge case is
// done by dropping a new pair of files into the directory — no test code
// changes required.
func TestPrettyFixtures(t *testing.T) {
	dir := filepath.Join("testdata", "pretty")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read fixture dir: %v", err)
	}
	cases := 0
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".in.md") {
			continue
		}
		cases++
		name := strings.TrimSuffix(e.Name(), ".in.md")
		t.Run(name, func(t *testing.T) {
			runFixture(t, dir, name)
		})
	}
	if cases == 0 {
		t.Fatal("no fixtures found in testdata/pretty/")
	}
}

func runFixture(t *testing.T, dir, name string) {
	t.Helper()
	in, err := os.ReadFile(filepath.Join(dir, name+".in.md"))
	if err != nil {
		t.Fatalf("read input: %v", err)
	}
	want, err := os.ReadFile(filepath.Join(dir, name+".want.txt"))
	if err != nil {
		t.Fatalf("read want: %v", err)
	}
	got := Render(string(in))
	if got != string(want) {
		t.Errorf("fixture %s mismatch\n--- want ---\n%s--- got ---\n%s",
			name, string(want), got)
	}
}

// TestRenderANSISwapsTokens guards the swap layer that turns sentinel
// tokens into real ANSI escape codes. Token leakage in CLI output would be
// embarrassing, so this is locked in.
func TestRenderANSISwapsTokens(t *testing.T) {
	out := RenderANSI(`mention "thing"`)
	if strings.Contains(out, TokCyanOpen) || strings.Contains(out, TokCyanClose) {
		t.Fatalf("RenderANSI leaked tokens: %q", out)
	}
	if !strings.Contains(out, constants.ColorCyan) || !strings.Contains(out, constants.ColorReset) {
		t.Fatalf("RenderANSI missing ANSI codes: %q", out)
	}
}

// TestUnterminatedQuoteClosedDefensively prevents a stray ANSI sequence
// from bleeding into the rest of the terminal session.
func TestUnterminatedQuoteClosedDefensively(t *testing.T) {
	out := Render(`"oops`)
	if !strings.Contains(out, TokCyanOpen) || !strings.Contains(out, TokCyanClose) {
		t.Fatalf("unterminated quote should still close: %q", out)
	}
}
