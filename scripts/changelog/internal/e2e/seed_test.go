package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// formatUnix renders a unix timestamp the way `git --date` accepts.
func formatUnix(unix int64) string {
	return fmt.Sprintf("%d +0000", unix)
}

// seedFiles writes the on-disk CHANGELOG.md / changelog.ts skeletons the
// writer expects to find before it splices a new entry. Returns the two
// absolute paths.
func seedFiles(t *testing.T, root string) (string, string) {
	t.Helper()

	mdPath := filepath.Join(root, "CHANGELOG.md")
	tsDir := filepath.Join(root, "src", "data")
	tsPath := filepath.Join(tsDir, "changelog.ts")

	err := os.MkdirAll(tsDir, 0o755)
	if err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	writeFile(t, mdPath, mdSeed)
	writeFile(t, tsPath, tsSeed)

	return mdPath, tsPath
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()

	err := os.WriteFile(path, []byte(body), 0o644)
	if err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// readFile is a tiny convenience used by both the runner and the
// fixture-diff assertion.
func readFile(t *testing.T, path string) string {
	t.Helper()

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(body)
}

const (
	mdSeed = "# Changelog\n\n"
	tsSeed = "export const changelog: ChangelogEntry[] = [\n];\n"
)
