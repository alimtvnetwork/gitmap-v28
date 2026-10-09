// Package termhelp provides terminal help layout engines, Catppuccin cards,
// and data-driven Markdown help document parsing.
package termout

import (
	"errors"
	"strings"
)

// FromMarkdown parses raw Markdown bytes into a structured HelpMenu.
func FromMarkdown(mdBytes []byte) (*HelpMenu, error) {
	if len(mdBytes) == 0 {
		return nil, errors.New("empty markdown input")
	}
	menu := parseMarkdown(string(mdBytes))
	return menu, nil
}

// RenderMarkdown parses markdown bytes and prints the styled help card.
func RenderMarkdown(mdBytes []byte) error {
	menu, err := FromMarkdown(mdBytes)
	if err != nil {
		return err
	}
	RenderMenu(*menu)
	return nil
}

type markdownParseState struct {
	menu           *HelpMenu
	curSection     *HelpSection
	activeSection  string
	inCodeFence    bool
	descriptionAcc []string
}

func newEmptyMenu() *HelpMenu {
	return &HelpMenu{
		UsageLines:  make([]string, 0),
		Sections:    make([]HelpSection, 0),
		FooterFlags: make([]CommandEntry, 0),
		Tips:        make([]string, 0),
	}
}

func parseMarkdown(text string) *HelpMenu {
	state := &markdownParseState{
		menu:           newEmptyMenu(),
		descriptionAcc: make([]string, 0),
	}
	lines := strings.Split(text, "\n")
	for _, rawLine := range lines {
		processMarkdownLine(rawLine, state)
	}
	return finalizeParsedMenu(state)
}

func processMarkdownLine(line string, state *markdownParseState) {
	trimmed := strings.TrimSpace(line)
	if checkToggleCodeFence(trimmed, state) {
		return
	}
	if strings.HasPrefix(line, "# ") {
		handleH1Line(line, state)
		return
	}
	if strings.HasPrefix(line, "## ") {
		handleH2Line(line, state)
		return
	}
	handleBodyLine(line, trimmed, state)
}

func checkToggleCodeFence(trimmed string, state *markdownParseState) bool {
	if isCodeFence(trimmed) {
		state.inCodeFence = !state.inCodeFence
		return true
	}
	return false
}

func isCodeFence(line string) bool {
	return strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~")
}

func handleH1Line(line string, state *markdownParseState) {
	title := strings.TrimPrefix(line, "# ")
	state.menu.Title = strings.TrimSpace(title)
	state.activeSection = "description"
}

func handleH2Line(line string, state *markdownParseState) {
	secTitle := strings.TrimPrefix(line, "## ")
	secTitle = strings.TrimSpace(secTitle)
	state.activeSection = strings.ToLower(secTitle)
	if isStandardSection(state.activeSection) {
		return
	}
	ensureActiveSection(state, secTitle)
}

func isStandardSection(name string) bool {
	return name == "usage" || name == "flags" || name == "tips"
}

func handleBodyLine(line, trimmed string, state *markdownParseState) {
	if len(trimmed) == 0 {
		return
	}
	if state.activeSection == "usage" {
		handleUsageLine(line, trimmed, state)
		return
	}
	if state.activeSection == "flags" {
		handleFlagsLine(trimmed, state)
		return
	}
	dispatchBodyLine(trimmed, state)
}

func dispatchBodyLine(trimmed string, state *markdownParseState) {
	if state.activeSection == "description" {
		handleDescriptionLine(trimmed, state)
		return
	}
	if state.activeSection == "tips" {
		handleTipsLine(trimmed, state)
		return
	}
	handleSectionLine(trimmed, state)
}

func handleDescriptionLine(trimmed string, state *markdownParseState) {
	clean := trimmed
	if strings.HasPrefix(clean, ">") {
		clean = strings.TrimPrefix(clean, ">")
		clean = strings.TrimSpace(clean)
	}
	state.descriptionAcc = append(state.descriptionAcc, clean)
}

func handleUsageLine(line, trimmed string, state *markdownParseState) {
	if isTableSeparator(trimmed) {
		return
	}
	clean := strings.TrimPrefix(trimmed, "$ ")
	state.menu.UsageLines = append(state.menu.UsageLines, clean)
}

func handleFlagsLine(trimmed string, state *markdownParseState) {
	if isTableSeparator(trimmed) {
		return
	}
	entry, hasEntry := parseTableOrBullet(trimmed)
	if hasEntry {
		state.menu.FooterFlags = append(state.menu.FooterFlags, entry)
	}
}

func handleTipsLine(trimmed string, state *markdownParseState) {
	clean := strings.TrimLeft(trimmed, "-* ")
	clean = strings.TrimPrefix(clean, "💡 ")
	clean = strings.TrimSpace(clean)
	if len(clean) > 0 {
		state.menu.Tips = append(state.menu.Tips, clean)
	}
}

func handleSectionLine(trimmed string, state *markdownParseState) {
	if state.curSection == nil || isTableSeparator(trimmed) {
		return
	}
	entry, hasEntry := parseTableOrBullet(trimmed)
	if hasEntry {
		state.curSection.Entries = append(state.curSection.Entries, entry)
	}
}

func parseTableOrBullet(trimmed string) (CommandEntry, bool) {
	if strings.HasPrefix(trimmed, "|") {
		return parseTableEntry(trimmed)
	}
	if strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "*") {
		return parseBulletEntry(trimmed)
	}
	return CommandEntry{}, false
}

func parseTableEntry(trimmed string) (CommandEntry, bool) {
	parts := splitCleanRow(trimmed)
	if len(parts) < 2 {
		return CommandEntry{}, false
	}
	cmd := cleanMarkdownTokens(parts[0])
	if strings.EqualFold(cmd, "flag") || strings.EqualFold(cmd, "command") {
		return CommandEntry{}, false
	}
	desc := cleanMarkdownTokens(parts[len(parts)-1])
	hasSubs := strings.Contains(cmd, "➔") || strings.Contains(cmd, "->")
	return CommandEntry{
		Command:        cmd,
		Description:    desc,
		HasSubcommands: hasSubs,
	}, true
}

func parseBulletEntry(trimmed string) (CommandEntry, bool) {
	clean := strings.TrimLeft(trimmed, "-* ")
	clean = strings.TrimSpace(clean)
	cmd, desc, hasSep := splitCommandDesc(clean)
	if !hasSep {
		return CommandEntry{Command: clean}, true
	}
	hasSubs := strings.Contains(cmd, "➔") || strings.Contains(cmd, "->")
	return CommandEntry{
		Command:        cleanMarkdownTokens(cmd),
		Description:    cleanMarkdownTokens(desc),
		HasSubcommands: hasSubs,
	}, true
}

func splitCommandDesc(s string) (string, string, bool) {
	for _, sep := range []string{" — ", " - ", ": ", "\t"} {
		idx := strings.Index(s, sep)
		if idx != -1 {
			return strings.TrimSpace(s[:idx]), strings.TrimSpace(s[idx+len(sep):]), true
		}
	}
	return s, "", false
}

func splitCleanRow(trimmed string) []string {
	trimmed = strings.Trim(trimmed, "|")
	rawParts := strings.Split(trimmed, "|")
	parts := make([]string, 0, len(rawParts))
	for _, p := range rawParts {
		cleaned := strings.TrimSpace(p)
		if len(cleaned) > 0 {
			parts = append(parts, cleaned)
		}
	}
	return parts
}

func isTableSeparator(line string) bool {
	clean := strings.ReplaceAll(line, "|", "")
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, ":", "")
	return len(strings.TrimSpace(clean)) == 0
}

func cleanMarkdownTokens(raw string) string {
	res := strings.ReplaceAll(raw, "`", "")
	res = strings.ReplaceAll(res, "**", "")
	res = strings.ReplaceAll(res, "*", "")
	res = strings.ReplaceAll(res, "\\<", "<")
	res = strings.ReplaceAll(res, "\\>", ">")
	return strings.TrimSpace(res)
}

func ensureActiveSection(state *markdownParseState, title string) {
	sec := HelpSection{
		Title:   title,
		Entries: make([]CommandEntry, 0),
	}
	state.menu.Sections = append(state.menu.Sections, sec)
	state.curSection = &state.menu.Sections[len(state.menu.Sections)-1]
}

func finalizeParsedMenu(state *markdownParseState) *HelpMenu {
	if len(state.descriptionAcc) > 0 {
		descText := strings.Join(state.descriptionAcc, " ")
		descSec := HelpSection{
			Title: "Description",
			Entries: []CommandEntry{
				{Command: state.menu.Title, Description: descText},
			},
		}
		state.menu.Sections = append([]HelpSection{descSec}, state.menu.Sections...)
	}
	return state.menu
}
