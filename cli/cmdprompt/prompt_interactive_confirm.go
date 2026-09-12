// Package cmdprompt — prompt_interactive_confirm.go prompts user before multi-repo installations.
package cmdprompt

func ShouldConfirmPromptBatch(count int) bool {
	return count > 5
}
