package cmdpending
import (
	"encoding/json"
	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"fmt"
	"strings"
)
func renderNewCommandsJSON(payload NewCommandsPayload) error {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return appfault.WrapExecution(err, "failed to serialize new commands JSON")
	}
	fmt.Println(string(data))
	return nil
}

func renderNewCommandsTerminal(payload NewCommandsPayload) error {
	renderTerminalHeader(payload)
	grouped := groupCommandsByCategory(payload.Commands)
	renderGroupedCategories(grouped)
	renderTerminalFooter(payload)
	return nil
}

func renderTerminalHeader(payload NewCommandsPayload) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                 GITMAP RECENT COMMANDS DISCOVERY ENGINE                      ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════════════════╣")
	line := formatHeaderSummaryLine(payload)
	fmt.Println(line)
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func formatHeaderSummaryLine(payload NewCommandsPayload) string {
	info := fmt.Sprintf("║ Total: %d | Returned: %d | Limit: %d", payload.TotalCommands, payload.FilteredCommands, payload.Limit)
	if payload.Category != "" {
		info += " | Cat: " + payload.Category
	}
	if payload.Filter != "" {
		info += " | Query: " + payload.Filter
	}
	for len(info) < 79 {
		info += " "
	}
	return info + "║"
}

type categoryGroup struct {
	category string
	commands []NewCommandEntry
}

func groupCommandsByCategory(commands []NewCommandEntry) []categoryGroup {
	var groups []categoryGroup
	groupMap := make(map[string]int)
	for _, cmd := range commands {
		idx, hasCat := groupMap[cmd.Category]
		if hasCat {
			groups[idx].commands = append(groups[idx].commands, cmd)
			continue
		}
		groupMap[cmd.Category] = len(groups)
		groups = append(groups, categoryGroup{category: cmd.Category, commands: []NewCommandEntry{cmd}})
	}
	return groups
}

func renderGroupedCategories(groups []categoryGroup) {
	for _, grp := range groups {
		renderCategorySection(grp)
	}
}

func renderCategorySection(grp categoryGroup) {
	fmt.Printf("▶ %s (%d commands)\n", strings.ToUpper(grp.category), len(grp.commands))
	fmt.Println(strings.Repeat("─", 80))
	for _, cmd := range grp.commands {
		renderCommandEntry(cmd)
	}
	fmt.Println()
}

func renderCommandEntry(cmd NewCommandEntry) {
	title := formatCommandTitle(cmd)
	fmt.Println(title)
	fmt.Printf("    %s\n", cmd.Description)
	fmt.Printf("    Example: %s\n", cmd.Example)
	fmt.Println()
}

func formatCommandTitle(cmd NewCommandEntry) string {
	res := "  ● gitmap " + cmd.Name
	if cmd.Alias != "" {
		res += " (" + cmd.Alias + ")"
	}
	if cmd.Version != "" {
		res += " [" + cmd.Version + "]"
	}
	return res
}

func renderTerminalFooter(payload NewCommandsPayload) {
	fmt.Println("────────────────────────────────────────────────────────────────────────────────")
	fmt.Printf("Tip: Use 'gitmap nc --category <cat>' or '--filter <query>' to narrow results.\n")
	fmt.Printf("     Use 'gitmap nc --json' for machine-readable JSON output.\n")
}

func printNewCommandsHelp() {
	fmt.Println("GitMap New Commands Discovery Engine")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gitmap new-commands [flags]")
	fmt.Println("  gitmap nc [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --limit, -n <N>          Limit number of commands returned (default: 100)")
	fmt.Println("  --category, -c <cat>     Filter by functional category")
	fmt.Println("  --filter, -f, -q <query> Search commands by name, alias, description, or example")
	fmt.Println("  --json, -j               Emit output as formatted JSON")
	fmt.Println("  -h, --help               Show this help text")
}
