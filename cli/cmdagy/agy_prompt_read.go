package cmdagy

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunAgyPromptRead reads and renders user prompts for a resolved conversation.
func RunAgyPromptRead(args []string) error {
	convID, err := resolvePromptReadConvID(args)
	if err != nil {
		return apperror.WrapSimple(err, "resolve conversation id")
	}

	brainDir, err := GetBrainLogsDirPath()
	if err != nil {
		return apperror.WrapSimple(err, "get brain logs directory")
	}

	prompts := readConvPrompts(brainDir, convID)
	renderConversationPrompts(convID, prompts)

	return nil
}

// RunAgyPromptLs lists conversation prompts across workspaces up to the given limit.
func RunAgyPromptLs(args []string) error {
	limit := parsePromptLimit(args)
	prompts := CollectAllPrompts()
	if len(prompts) == 0 {
		fmt.Println("No conversation prompts found in Antigravity brain.")

		return nil
	}

	displayPrompts := limitPrompts(prompts, limit)
	renderPromptsTable(displayPrompts)

	return nil
}

func limitPrompts(prompts []AgyPromptEntry, limit int) []AgyPromptEntry {
	if len(prompts) > limit {
		return prompts[:limit]
	}

	return prompts
}

func renderConversationPrompts(convID string, prompts []AgyPromptEntry) {
	fmt.Printf("\n  %s● Antigravity Conversation: %s%s\n", constants.ColorCyan, convID, constants.ColorReset)
	if len(prompts) == 0 {
		fmt.Printf("  No user prompts recorded in this conversation.\n\n")

		return
	}

	fmt.Printf("  Total Prompts: %d\n\n", len(prompts))
	for _, p := range prompts {
		renderSinglePromptEntry(p)
	}
}

func renderSinglePromptEntry(p AgyPromptEntry) {
	relTime := formatRelativeTime(p.CreatedAt.Format(time.RFC3339))
	fmt.Printf("  %s[Step %d]%s %s (%s)\n", constants.ColorGreen, p.StepIndex, constants.ColorReset, p.CreatedAt.Format("2006-01-02 15:04:05"), relTime)
	fmt.Printf("  %s\n\n", p.Content)
}
