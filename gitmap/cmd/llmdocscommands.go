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

// writeLLMCommandTable writes the markdown table header and rows for a command group.
func writeLLMCommandTable(sb *strings.Builder, g llmCmdGroup) {
	fmt.Fprintf(sb, "### %s\n\n", g.title)
	sb.WriteString("| Command | Alias | What it does |\n")
	sb.WriteString("|---------|-------|--------------|\n")
	for _, c := range g.commands {
		fmt.Fprintf(sb, "| `%s` | `%s` | %s |\n", c.name, c.alias, c.desc)
	}
}

// writeLLMCommandExamples writes the markdown code block containing command examples.
func writeLLMCommandExamples(sb *strings.Builder, commands []llmCmdEntry) {
	sb.WriteString("\n**Examples:**\n```bash\n")
	for _, c := range commands {
		if c.example != "" {
			sb.WriteString(c.example + "\n")
		}
	}
	sb.WriteString("```\n\n")
}

// writeLLMCommandGroup writes a single command group table with examples.
func writeLLMCommandGroup(sb *strings.Builder, g llmCmdGroup) {
	writeLLMCommandTable(sb, g)
	writeLLMCommandExamples(sb, g.commands)
}

// buildCoreCommandGroups returns the core workflow command groups.
func buildCoreCommandGroups() []llmCmdGroup {
	return []llmCmdGroup{
		buildScanningGroup(),
		buildCloningGroup(),
		buildGitOpsGroup(),
		buildNavigationGroup(),
		buildReleaseGroup(),
		buildReleaseInfoGroup(),
		buildDataGroup(),
		buildHistoryGroup(),
	}
}

// buildExtendedCommandGroups returns the extended utility command groups.
func buildExtendedCommandGroups() []llmCmdGroup {
	return []llmCmdGroup{
		buildAmendGroup(),
		buildProjectGroup(),
		buildSSHGroup(),
		buildZipGroup(),
		buildEnvToolsGroup(),
		buildTaskGroup(),
		buildVisibilityGroup(),
		buildUtilityGroup(),
	}
}

// buildCommandGroups returns all command groups dynamically.
func buildCommandGroups() []llmCmdGroup {
	groups := buildCoreCommandGroups()

	return append(groups, buildExtendedCommandGroups()...)
}
