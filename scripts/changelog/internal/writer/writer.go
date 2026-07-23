// Package writer splices a freshly-rendered Entry into CHANGELOG.md
// (after the leading "# Changelog" header) and src/data/changelog.ts
// (at the top of the exported array).
package writer

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/scripts/changelog/internal/render"
)

const (
	mdHeader      = "# Changelog\n"
	tsArrayMarker = "export const changelog: ChangelogEntry[] = [\n"
)

// PrependBoth writes the new Markdown entry to mdPath and the TypeScript
// entry to tsPath, leaving every existing entry intact.
func PrependBoth(mdPath, tsPath string, entry render.Entry) error {
	err := prependFile(mdPath, mdHeader+"\n", render.Markdown(entry))
	if err != nil {
		return fmt.Errorf("CHANGELOG.md: %w", err)
	}

	err = spliceAfterMarker(tsPath, tsArrayMarker, render.TypeScript(entry))
	if err != nil {
		return fmt.Errorf("src/data/changelog.ts: %w", err)
	}

	return nil
}

func prependFile(path, header, body string) error {
	current, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	rest := strings.TrimPrefix(string(current), header)
	out := header + body + strings.TrimLeft(rest, "\n")

	return os.WriteFile(path, []byte(out), 0o644)
}

func spliceAfterMarker(path, marker, fragment string) error {
	current, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	idx := strings.Index(string(current), marker)
	if idx < 0 {
		return fmt.Errorf("marker %q not found", marker)
	}

	insertAt := idx + len(marker)
	out := string(current[:insertAt]) + fragment + string(current[insertAt:])

	return os.WriteFile(path, []byte(out), 0o644)
}
