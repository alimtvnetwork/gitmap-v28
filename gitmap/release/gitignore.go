package release

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// gitignoreEntries lists paths that must be present in .gitignore.
var gitignoreEntries = []string{
	constants.AssetsStagingDir,
	"release-assets",
}

// EnsureGitignore appends missing release-related entries to .gitignore.
func EnsureGitignore() {
	const path = ".gitignore"

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	existing := make(map[string]bool, len(lines))
	for _, l := range lines {
		existing[strings.TrimSpace(l)] = true
	}

	var toAdd []string
	for _, entry := range gitignoreEntries {
		if !existing[entry] && !existing["/"+entry] {
			toAdd = append(toAdd, entry)
		}
	}

	if len(toAdd) == 0 {
		return
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("gitignore: appending %d entries", len(toAdd))
	}

	// Ensure trailing newline before appending.
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	content += strings.Join(toAdd, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write .gitignore at %s: %v\n", path, err)
	}
}
