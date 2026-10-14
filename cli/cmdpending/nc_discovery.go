// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// discoveredCommand is one raw hit from router discovery: every registered
// name for the entry (first = shortest, treated as the primary name),
// an optional description, the dispatch group it was found in, and an
// explicit alias display for usage-mined alias entries.
type discoveredCommand struct {
	names       []string
	description string
	category    string
	kind        string // dispatch | cobra | usage | registry
	aliasText   string
}

// discoverRouterCommands enumerates commands from the live dispatcher
// registration instead of a hardcoded list:
//  1. []dispatchEntry routing tables in cli/cmd (names + aliases, grouped
//     exactly as the router groups them),
//  2. &cobra.Command{Use, Aliases, Short} registrations (with canonical help text),
//  3. subsystem usage/help text (multi-word commands like "spec next"),
//  4. the generated completion registry as a safety net (works even when
//     the source tree is unavailable).
//
// Descriptions resolve from cobra Short text, then the helpdoc topic
// catalog, then constant doc comments captured during parsing.
func discoverRouterCommands() []NewCommandEntry {
	var found []discoveredCommand
	if srcDir := locateDispatchSourceDir(); srcDir != "" {
		consts := parseCommandConstants(srcDir + "/constants")
		cmdDir := srcDir + "/cmd"
		found = append(found, parseDispatchTables(cmdDir, consts)...)
		found = append(found, parseCobraCommands(cmdDir)...)
		found = append(found, mineSubcommandUsage(srcDir)...)
		enrichSubcommandDescriptions(found, consts)
	}
	known := make(map[string]bool, len(found)*2)
	for _, dc := range found {
		for _, n := range dc.names {
			known[n] = true
		}
	}
	found = append(found, registrySafetyNet(known)...)
	return buildDiscoveredEntries(found)
}

// enrichSubcommandDescriptions fills blank usage-mined descriptions from the
// doc comments captured on the matching constants (e.g. CmdSpecNext documents
// `spec next`).
func enrichSubcommandDescriptions(found []discoveredCommand, consts map[string]routerConstInfo) {
	for i, dc := range found {
		if dc.kind != "usage" || dc.description != "" || len(dc.names) == 0 {
			continue
		}
		if desc := constDocForCompound(dc.names[0], consts); desc != "" {
			found[i].description = desc
		}
	}
}

// constDocForCompound finds a constant whose doc comment describes the
// compound command, e.g. "spec next" via CmdSpecNext's comment.
func constDocForCompound(compound string, consts map[string]routerConstInfo) string {
	for _, info := range consts {
		doc := info.doc
		if doc == "" || !strings.Contains(doc, "`"+compound+"`") {
			continue
		}
		if desc := strings.TrimSpace(stripDocPrefix(doc, compound)); desc != "" {
			return desc
		}
	}
	return ""
}

func stripDocPrefix(doc, compound string) string {
	lower := strings.ToLower(doc)
	marker := "is the `" + compound + "` subcommand"
	if idx := strings.Index(lower, marker); idx >= 0 {
		rest := strings.TrimSpace(doc[idx+len(marker):])
		return strings.TrimPrefix(rest, "— ")
	}
	return doc
}

// registrySafetyNet unions the generated completion registry so commands the
// table/cobra/usage passes miss still appear. Names already attached to a
// discovered entry are skipped.
func registrySafetyNet(known map[string]bool) []discoveredCommand {
	var out []discoveredCommand
	for _, name := range completion.AllCommands() {
		if known[name] {
			continue
		}
		out = append(out, discoveredCommand{
			names:    []string{name},
			category: dispatchGroupCategory(name),
			kind:     "registry",
		})
	}
	return out
}

// buildDiscoveredEntries converts raw hits into catalog entries: deduped by
// primary name, shortest name first (matching the CLI's short-alias-first
// convention), router discovery order preserved. `byName` tracks emitted
// primary names; `known` tracks every attached name so the registry safety
// net never re-adds a covered alias as a flat entry.
func buildDiscoveredEntries(found []discoveredCommand) []NewCommandEntry {
	known := make(map[string]bool, len(found)*2)
	byName := make(map[string]int, len(found))
	entries := make([]NewCommandEntry, 0, len(found))
	version := "v" + constants.Version
	parentCategory := make(map[string]string)
	for _, dc := range found {
		if dc.kind == "usage" {
			entries = appendUsageEntry(entries, dc, version, parentCategory, known, byName)
			continue
		}
		name, alias := splitDiscoveredNames(dc.names)
		if name == "" {
			continue
		}
		if idx, hasEntry := byName[name]; hasEntry {
			enrichEntryDescription(entries, idx, dc)
			continue
		}
		byName[name] = len(entries)
		markNamesKnown(known, dc.names)
		category := discoveredCategory(dc)
		entries = append(entries, NewCommandEntry{
			Name:        name,
			Alias:       alias,
			Version:     version,
			Category:    category,
			Description: resolveDescription(dc, name),
			Example:     "gitmap " + name,
		})
		parentCategory[name] = category
	}
	return entries
}

func appendUsageEntry(entries []NewCommandEntry, dc discoveredCommand, version string,
	parentCategory map[string]string, known map[string]bool, byName map[string]int) []NewCommandEntry {
	if len(dc.names) == 0 {
		return entries
	}
	name := dc.names[0]
	if idx, hasEntry := byName[name]; hasEntry {
		enrichEntryDescription(entries, idx, dc)
		return entries
	}
	alias := dc.aliasText
	if alias == "" && len(dc.names) > 1 {
		alias = strings.Join(dc.names[1:], ", ")
	}
	byName[name] = len(entries)
	markNamesKnown(known, dc.names)
	category := dc.category
	if parent := parentOfCompound(name); parent != "" {
		if pc, hasParent := parentCategory[parent]; hasParent {
			category = pc
		}
	}
	entries = append(entries, NewCommandEntry{
		Name:        name,
		Alias:       alias,
		Version:     version,
		Category:    category,
		Description: resolveDescription(dc, name),
		Example:     "gitmap " + name,
	})
	return entries
}

// enrichEntryDescription upgrades a generic placeholder description with a
// real one discovered later (e.g. a cobra Short shadowing a dispatch entry).
func enrichEntryDescription(entries []NewCommandEntry, idx int, dc discoveredCommand) {
	if strings.TrimSpace(dc.description) == "" {
		return
	}
	if !isGenericDescription(entries[idx].Description) {
		return
	}
	entries[idx].Description = strings.TrimSpace(dc.description)
}

func isGenericDescription(desc string) bool {
	return strings.HasPrefix(desc, "Run 'gitmap ") && strings.HasSuffix(desc, " --help' for usage details.")
}

func parentOfCompound(name string) string {
	if idx := strings.Index(name, " "); idx > 0 {
		return name[:idx]
	}
	return ""
}

func markNamesKnown(known map[string]bool, names []string) {
	for _, n := range names {
		known[n] = true
	}
}

// splitDiscoveredNames picks the shortest registered name as the primary
// (the CLI's short-alias-first convention: "cpar" over "commit-push-all-repos")
// and joins the rest as the alias value.
func splitDiscoveredNames(names []string) (string, string) {
	if len(names) == 0 {
		return "", ""
	}
	ordered := append([]string(nil), names...)
	sort.SliceStable(ordered, func(i, j int) bool { return len(ordered[i]) < len(ordered[j]) })
	name := ordered[0]
	var rest []string
	for _, n := range ordered[1:] {
		if n != name {
			rest = append(rest, n)
		}
	}
	return name, strings.Join(rest, ", ")
}

// discoveredCategory applies the commit-family name override first, then the
// dispatch-group mapping recorded during parsing.
