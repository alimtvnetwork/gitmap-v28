package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
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

func renderPromptsTable(prompts []AgyPromptEntry) {
	fmt.Printf("\n  %s● AGY Prompts (%d entries)%s\n\n", constants.ColorCyan, len(prompts), constants.ColorReset)
	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "DATE", Align: termtable.AlignLeft, MinWidth: 16},
			{Title: "WORKSPACE", Align: termtable.AlignLeft, MinWidth: 20},
			{Title: "PROMPT PREVIEW", Align: termtable.AlignLeft, MinWidth: 45},
		},
		Rows: buildPromptRows(prompts),
	}
	termtable.PrintTable(cfg)
	termpad.EnsureBottomPadding("\n")
}

func buildPromptRows(prompts []AgyPromptEntry) []termtable.Row {
	var rows []termtable.Row
	for _, p := range prompts {
		rows = append(rows, termtable.Row{
			Cells: []string{
				p.CreatedAt.Format("2006-01-02 15:04"),
				truncatePromptString(p.Workspace, 20),
				truncatePromptString(p.Content, 45),
			},
		})
	}

	return rows
}

func truncatePromptString(s string, max int) string {
	clean := strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
	if len(clean) <= max {
		return clean
	}

	return clean[:max-3] + "..."
}
