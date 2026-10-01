package release

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
)

// gitignoreEntries lists paths that must be present in .gitignore.
var gitignoreEntries = []string{
	constants.AssetsStagingDir,
	"release-assets",
	".antigravity_resume_task.json",
	"antigravity-resume_task.json",
}

// EnsureGitignore sanitizes and appends missing release-related entries to .gitignore.
func EnsureGitignore() {
	const path = ".gitignore"

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	content := string(data)
	cleaned, isModified := gitignoreagm.DeduplicateAndSanitizeGitignore(content, gitignoreEntries...)
	if !isModified {
		return
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("gitignore: sanitized and updated %s", path)
	}

	if writeErr := os.WriteFile(path, []byte(cleaned), constants.FilePermission); writeErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write .gitignore at %s: %v\n", path, writeErr)
	}
}
