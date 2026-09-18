package cmdagy

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
)

func resolvePromptsForRerun(cwd string, count int) []AgyPromptEntry {
	prompts := CollectPromptsForWorkspace(cwd)
	if len(prompts) == 0 {
		prompts = CollectAllPrompts()
	}
	if len(prompts) > count {
		return prompts[:count]
	}

	return prompts
}

func buildRerunPayload(tplContent string, prompts []AgyPromptEntry) string {
	var parts []string
	if tplContent != "" {
		parts = append(parts, tplContent)
	}
	for i := len(prompts) - 1; i >= 0; i-- {
		parts = append(parts, prompts[i].Content)
	}

	return strings.Join(parts, "\n\n")
}

func renderRerunOutput(payload string, count int) {
	fmt.Printf("\n  %s● Replaying Last %d Prompt(s)%s\n", constants.ColorCyan, count, constants.ColorReset)
	fmt.Println("  ────────────────────────────────────────────────────────────────────────────────")
	termpad.PrintPadded(payload)
	fmt.Printf("\n  ────────────────────────────────────────────────────────────────────────────────\n")
	termpad.EnsureBottomPadding("\n")
}

func handleRerunClipboard(payload string) {
	if rerunNoClipboard || rerunDryRun {
		return
	}
	if err := clipboard.WriteAll(payload); err == nil {
		fmt.Printf("  %s✔ Copied constructed prompt to OS clipboard.%s\n\n", constants.ColorGreen, constants.ColorReset)
	}
}
