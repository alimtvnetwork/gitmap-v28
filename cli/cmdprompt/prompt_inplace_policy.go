// Package cmdprompt — prompt_inplace_policy.go checks in-place update requirements.
package cmdprompt

import "github.com/alimtvnetwork/gitmap-v28/cli/installer"

func ShouldPromptForOverwrite(targetDir string) bool {
	return installer.HasExistingPrompts(targetDir)
}
