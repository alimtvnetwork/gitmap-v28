// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunNewCommands executes the new-commands discovery tool.
func RunNewCommands(args []string) error {
	opts, appErr := parseNewCommandsOptions(args)
	if appErr != nil {
		return appErr
	}
	if opts.IsHelp {
		printNewCommandsHelp()
		return nil
	}
	return executeNewCommands(opts)
}

func executeNewCommands(opts NewCommandsOptions) error {
	catalog := buildNewCommandsCatalog()
	payload := buildNewCommandsPayload(catalog, opts)
	if opts.IsJSON {
		return renderNewCommandsJSON(payload)
	}
	return renderNewCommandsTerminal(payload)
}

func parseNewCommandsOptions(args []string) (NewCommandsOptions, *appfault.AppError) {
	opts := DefaultNewCommandsOptions()
	for i := 0; i < len(args); i++ {
		nextIdx, appErr := parseNextCommandOption(args, i, &opts)
		if appErr != nil {
			return opts, appErr
		}
		i = nextIdx
	}
	clampNewCommandsLimit(&opts)
	return opts, nil
}

func parseNextCommandOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if isHelpFlag(arg) {
		opts.IsHelp = true
		return idx, nil
	}
	if isJSONFlag(arg) {
		opts.IsJSON = true
		return idx, nil
	}
	return parseValuedOption(args, idx, opts)
}

func parseValuedOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if isCategoryFlag(arg) {
		return parseCategoryOption(args, idx, opts)
	}
	if isFilterFlag(arg) {
		return parseFilterOption(args, idx, opts)
	}
	if isLimitFlag(arg) {
		return parseLimitOption(args, idx, opts)
	}
	return idx, nil
}

func parseCategoryOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if strings.HasPrefix(arg, "--category=") {
		opts.Category = strings.TrimSpace(strings.TrimPrefix(arg, "--category="))
		return idx, nil
	}
	if strings.HasPrefix(arg, "-c=") {
		opts.Category = strings.TrimSpace(strings.TrimPrefix(arg, "-c="))
		return idx, nil
	}
	if idx+1 >= len(args) {
		return idx, appfault.NewValidation("new-commands", "E9003", "missing value for --category")
	}
	opts.Category = strings.TrimSpace(args[idx+1])
	return idx + 1, nil
}

func parseFilterOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if strings.HasPrefix(arg, "--filter=") {
		opts.Filter = strings.TrimSpace(strings.TrimPrefix(arg, "--filter="))
		return idx, nil
	}
	if strings.HasPrefix(arg, "-f=") {
		opts.Filter = strings.TrimSpace(strings.TrimPrefix(arg, "-f="))
		return idx, nil
	}
	if strings.HasPrefix(arg, "-q=") {
		opts.Filter = strings.TrimSpace(strings.TrimPrefix(arg, "-q="))
		return idx, nil
	}
	if idx+1 >= len(args) {
		return idx, appfault.NewValidation("new-commands", "E9003", "missing value for --filter")
	}
	opts.Filter = strings.TrimSpace(args[idx+1])
	return idx + 1, nil
}

func parseLimitOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	val, nextIdx, appErr := resolveLimitValue(args, idx, arg)
	if appErr != nil {
		return idx, appErr
	}
	lim, err := strconv.Atoi(val)
	if err != nil {
		return idx, appfault.NewValidation("new-commands", "E9003", "invalid numeric value for --limit: "+val)
	}
	opts.Limit = lim
	return nextIdx, nil
}

func resolveLimitValue(args []string, idx int, arg string) (string, int, *appfault.AppError) {
	if strings.HasPrefix(arg, "--limit=") {
		return strings.TrimPrefix(arg, "--limit="), idx, nil
	}
	if strings.HasPrefix(arg, "-n=") {
		return strings.TrimPrefix(arg, "-n="), idx, nil
	}
	if idx+1 >= len(args) {
		return "", idx, appfault.NewValidation("new-commands", "E9003", "missing value for --limit")
	}
	return args[idx+1], idx + 1, nil
}

func clampNewCommandsLimit(opts *NewCommandsOptions) {
	if opts.Limit <= 0 {
		opts.Limit = 100
	}
}

func isHelpFlag(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func isJSONFlag(arg string) bool {
	return arg == "--json" || arg == "-j"
}

func isCategoryFlag(arg string) bool {
	return arg == "--category" || arg == "-c" || strings.HasPrefix(arg, "--category=") || strings.HasPrefix(arg, "-c=")
}

func isFilterFlag(arg string) bool {
	return arg == "--filter" || arg == "-f" || arg == "-q" ||
		strings.HasPrefix(arg, "--filter=") || strings.HasPrefix(arg, "-f=") || strings.HasPrefix(arg, "-q=")
}

func isLimitFlag(arg string) bool {
	return arg == "--limit" || arg == "-n" || strings.HasPrefix(arg, "--limit=") || strings.HasPrefix(arg, "-n=")
}

func buildNewCommandsPayload(catalog []NewCommandEntry, opts NewCommandsOptions) NewCommandsPayload {
	filtered := filterNewCommands(catalog, opts)
	return NewCommandsPayload{
		Timestamp:        time.Now().UTC(),
		TotalCommands:    len(catalog),
		FilteredCommands: len(filtered),
		Limit:            opts.Limit,
		Category:         opts.Category,
		Filter:           opts.Filter,
		Commands:         filtered,
	}
}

func filterNewCommands(catalog []NewCommandEntry, opts NewCommandsOptions) []NewCommandEntry {
	res := make([]NewCommandEntry, 0, len(catalog))
	for _, entry := range catalog {
		if isCommandMatched(entry, opts) {
			res = append(res, entry)
		}
		if len(res) >= opts.Limit {
			break
		}
	}
	return res
}

func isCommandMatched(entry NewCommandEntry, opts NewCommandsOptions) bool {
	if !isCategoryMatch(entry.Category, opts.Category) {
		return false
	}
	if !isTextFilterMatch(entry, opts.Filter) {
		return false
	}
	return true
}

func isCategoryMatch(entryCat, filterCat string) bool {
	trimmed := strings.TrimSpace(filterCat)
	if trimmed == "" {
		return true
	}
	return strings.EqualFold(entryCat, trimmed)
}

func isTextFilterMatch(entry NewCommandEntry, filter string) bool {
	query := strings.ToLower(strings.TrimSpace(filter))
	if query == "" {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Name), query) {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Alias), query) {
		return true
	}
	return isEntryDetailMatch(entry, query)
}

func isEntryDetailMatch(entry NewCommandEntry, query string) bool {
	if strings.Contains(strings.ToLower(entry.Description), query) {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Example), query) {
		return true
	}
	return false
}

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

// buildNewCommandsCatalog discovers commands from the live dispatcher
// registration (dispatch tables, cobra commands, subsystem usage text, and
// the generated completion registry). The hardcoded seed catalog is used
// ONLY when discovery finds nothing.
func buildNewCommandsCatalog() []NewCommandEntry {
	if discovered := discoverRouterCommands(); len(discovered) > 0 {
		return discovered
	}
	return seedNewCommandsCatalog()
}
