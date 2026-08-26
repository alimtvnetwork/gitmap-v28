// Package cmd — prompt_inplace_policy.go checks in-place update requirements.
package cmd

import "github.com/alimtvnetwork/gitmap-v28/gitmap/installer"

func ShouldPromptForOverwrite(targetDir string) bool {
	return installer.HasExistingPrompts(targetDir)
}
