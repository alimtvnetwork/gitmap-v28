// Package cmdprompt — prompt_version_cmd.go prints prompt architect version.
package cmdprompt

import (
	"fmt"
	"path/filepath"
)

func RunPromptVersion(targets []string) error {
	if len(targets) == 0 {
		targets, _ = ResolvePromptTarget("")
	}

	for _, t := range targets {
		name := filepath.Base(t)
		meta, _ := ReadPromptArchitectMetadata(t)
		if IsPromptArchitectInstalled(meta) {
			fmt.Printf("%s: %s (%s)\n", name, meta.Version, meta.Status)
		} else {
			fmt.Printf("%s: not installed\n", name)
		}
	}

	return nil
}
