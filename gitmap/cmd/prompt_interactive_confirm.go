// Package cmd — prompt_interactive_confirm.go prompts user before multi-repo installations.
package cmd

func ShouldConfirmPromptBatch(count int) bool {
	return count > 5
}
