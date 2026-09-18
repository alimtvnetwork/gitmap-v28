package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func displayPromptListing(limit int) error {
	prompts := selectPromptsForDisplay(limit)
	if listPromptsJSON {
		return outputPromptsJSON(prompts)
	}
	renderPromptsTable(prompts)

	return nil
}

func selectPromptsForDisplay(limit int) []AgyPromptEntry {
	all := CollectAllPrompts()
	if listPromptsProjectPrefix != "" {
		return filterPromptsByPrefix(all, listPromptsProjectPrefix, limit)
	}
	if listPromptsAllProjects {
		return truncatePromptSlice(all, limit)
	}

	return resolveCurrentWorkspacePrompts(all, limit)
}

func resolveCurrentWorkspacePrompts(all []AgyPromptEntry, limit int) []AgyPromptEntry {
	cwd, _ := os.Getwd()
	matched := CollectPromptsForWorkspace(cwd)
	if len(matched) > 0 {
		return truncatePromptSlice(matched, limit)
	}

	return truncatePromptSlice(all, limit)
}

func filterPromptsByPrefix(all []AgyPromptEntry, prefix string, limit int) []AgyPromptEntry {
	var out []AgyPromptEntry
	norm := strings.ToLower(prefix)
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Workspace), norm) {
			out = append(out, p)
		}
		if len(out) >= limit {
			break
		}
	}

	return out
}

func truncatePromptSlice(list []AgyPromptEntry, limit int) []AgyPromptEntry {
	if len(list) > limit {
		return list[:limit]
	}

	return list
}

func outputPromptsJSON(list []AgyPromptEntry) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	return nil
}
