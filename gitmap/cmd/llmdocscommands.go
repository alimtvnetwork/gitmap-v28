// Package cmd — llmdocscommands.go writes the command reference tables.
package cmd

import (
	"fmt"
	"strings"
)

// llmCmdEntry represents a command for LLM doc generation.
type llmCmdEntry struct {
	name    string
	alias   string
	desc    string
	example string
}

// llmCmdGroup represents a group of commands.
type llmCmdGroup struct {
	title    string
	commands []llmCmdEntry
}

// writeLLMCommands writes all command reference tables.
func writeLLMCommands(sb *strings.Builder) {
	sb.WriteString("## Complete Command Reference\n\n")

	groups := buildCommandGroups()
	for _, g := range groups {
		writeLLMCommandGroup(sb, g)
	}
}

// writeLLMCommandGroup writes a single command group table with examples.
func writeLLMCommandGroup(sb *strings.Builder, g llmCmdGroup) {
	fmt.Fprintf(sb, "### %s\n\n", g.title)
	sb.WriteString("| Command | Alias | What it does |\n")
	sb.WriteString("|---------|-------|--------------|\n")

	for _, c := range g.commands {
		fmt.Fprintf(sb, "| `%s` | `%s` | %s |\n", c.name, c.alias, c.desc)
	}

	sb.WriteString("\n**Examples:**\n```bash\n")

	for _, c := range g.commands {
		if c.example != "" {
			sb.WriteString(c.example + "\n")
		}
	}

	sb.WriteString("```\n\n")
}

// buildCommandGroups returns all command groups dynamically.
func buildCommandGroups() []llmCmdGroup {
	return []llmCmdGroup{
		buildScanningGroup(),
		buildCloningGroup(),
		buildGitOpsGroup(),
		buildNavigationGroup(),
		buildReleaseGroup(),
		buildReleaseInfoGroup(),
		buildDataGroup(),
		buildHistoryGroup(),
		buildAmendGroup(),
		buildProjectGroup(),
		buildSSHGroup(),
		buildZipGroup(),
		buildEnvToolsGroup(),
		buildTaskGroup(),
		buildUtilityGroup(),
	}
}
