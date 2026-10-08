package helpdisplay

import (
	"fmt"
	"strings"
)

// HelpDisplay is the concrete help display, bound to Displayer via the
// NewHelpDisplay constructor.
type HelpDisplay struct {
	header      string
	groups      []CommandHelpGroup
	suggestions []Suggestion
	theme       *Theme
}

// NewHelpDisplay binds HelpDisplay to the Displayer interface.
func NewHelpDisplay(header string, groups []CommandHelpGroup, suggestions []Suggestion, theme *Theme) Displayer {
	return &HelpDisplay{header: header, groups: groups, suggestions: suggestions, theme: theme}
}

// Render builds the full help text: header, command groups, hints, suggestions.
func (d *HelpDisplay) Render(ctx RenderContext) string {
	var out strings.Builder
	d.writeHeader(&out)
	for i := range d.groups {
		d.writeGroup(&out, &d.groups[i])
	}
	d.writeHints(&out)
	d.writeSuggestions(&out, ctx)
	return out.String()
}

// Print writes the rendered help text to stdout.
func (d *HelpDisplay) Print(ctx RenderContext) {
	fmt.Print(d.Render(ctx))
}

// colorize applies a theme color via Color.Wrap; a nil theme yields plain text.
func (d *HelpDisplay) colorize(getColor func(*Theme) Color, text string) string {
	if d.theme == nil {
		return text
	}
	return getColor(d.theme).Wrap(text)
}

func (d *HelpDisplay) writeHeader(out *strings.Builder) {
	out.WriteString(d.colorize((*Theme).HeaderColor, d.header) + "\n\n")
}

func (d *HelpDisplay) writeGroup(out *strings.Builder, g *CommandHelpGroup) {
	out.WriteString(d.colorize((*Theme).HeaderColor, g.header) + "\n")
	for i := range g.commands {
		d.writeCommand(out, &g.commands[i])
	}
}

func (d *HelpDisplay) writeCommand(out *strings.Builder, c *CommandHelper) {
	out.WriteString(d.commandLine(c) + "\n")
	d.writeExample(out, c)
}

// commandLine renders one `  <command>  <description>` line.
func (d *HelpDisplay) commandLine(c *CommandHelper) string {
	line := "  " + d.colorize((*Theme).CommandColor, c.command)
	if len(c.description) == 0 {
		return line
	}
	return line + "  " + d.colorize((*Theme).DescriptionColor, c.description)
}

func (d *HelpDisplay) writeExample(out *strings.Builder, c *CommandHelper) {
	if len(c.example) == 0 {
		return
	}
	out.WriteString("    e.g. " + d.colorize((*Theme).DescriptionColor, c.example) + "\n")
}

// hasHints reports whether any group carries hints.
func (d *HelpDisplay) hasHints() bool {
	for i := range d.groups {
		if len(d.groups[i].hints) > 0 {
			return true
		}
	}
	return false
}

func (d *HelpDisplay) writeHints(out *strings.Builder) {
	if !d.hasHints() {
		return
	}
	out.WriteString("\n")
	for i := range d.groups {
		d.writeGroupHints(out, &d.groups[i])
	}
}

func (d *HelpDisplay) writeGroupHints(out *strings.Builder, g *CommandHelpGroup) {
	for _, h := range g.hints {
		out.WriteString("  " + d.colorize((*Theme).HintColor, h) + "\n")
	}
}

func (d *HelpDisplay) writeSuggestions(out *strings.Builder, ctx RenderContext) {
	if len(d.suggestions) == 0 {
		return
	}
	out.WriteString("\n")
	for i := range d.suggestions {
		out.WriteString("* " + d.suggestions[i].Render(ctx) + "\n")
	}
}
