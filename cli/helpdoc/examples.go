// Package helpdoc — examples.go: canonical help-example registry (program 273).
//
// Single source of truth for `gitmap help example <cmd>`: every entry is
// what-it-does + concrete invocation + expected output (spec
// 03-help-example-format.md). Three renderers consume this registry and
// nothing else:
//
//   - `gitmap help example <cmd>` (terminal, theme-aware)
//   - the "Examples" section appended to `gitmap help <cmd>`
//   - the generated helpdoc/*.md topics (via the emitter)
//
// Hand-maintaining the same example in two formats is banned: curated
// entries live in examples_<domain>.go; every other documented topic is
// derived at init from the topic summaries and the embedded help files.
package helpdoc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// Example is one concrete invocation plus its expected output.
// Invocations use real-looking arguments (never <foo>); unstable values
// in expected output are normalized to <sha>, <date>, <machine>,
// <count>, <path>.
type Example struct {
	Invocation     string
	ExpectedOutput string
}

// HelpExample is the canonical 3-part help example for one command.
type HelpExample struct {
	WhatItDoes string
	Examples   []Example
}

// helpExampleRegistry is the single registry, keyed by command verb.
var helpExampleRegistry = buildHelpExampleRegistry()

// buildHelpExampleRegistry composes curated entries (examples_<domain>.go)
// with derived entries so every documented topic is covered.
func buildHelpExampleRegistry() map[string]HelpExample {
	reg := make(map[string]HelpExample, 512)
	registerCoreExamples(reg)
	registerFleetExamples(reg)
	registerDevExamples(reg)
	registerOpsExamples(reg)
	registerMergeWriteExamples(reg)
	register275Examples(reg)
	deriveSummaryExamples(reg)
	deriveFileExamples(reg)

	return reg
}

// normalizeExampleKey lowercases and trims the verb for lookup.
func normalizeExampleKey(verb string) string {
	return strings.ToLower(strings.TrimSpace(verb))
}

// LookupHelpExample returns the canonical example for a command verb,
// resolving help aliases (pc → pending-commits, …). ok=false when the
// verb has no entry.
func LookupHelpExample(verb string) (HelpExample, bool) {
	key := normalizeExampleKey(verb)
	if ex, ok := helpExampleRegistry[key]; ok {
		return ex, true
	}

	aliased := resolveHelpAlias(key)
	if aliased == key {
		return HelpExample{}, false
	}

	ex, ok := helpExampleRegistry[aliased]
	if ok {
		return ex, true
	}

	return HelpExample{}, false
}

// HelpExampleVerbs returns the sorted registry keys (verbs and aliases).
func HelpExampleVerbs() []string {
	verbs := make([]string, 0, len(helpExampleRegistry))
	for verb := range helpExampleRegistry {
		verbs = append(verbs, verb)
	}

	sort.Strings(verbs)

	return verbs
}

// renderHelpExample renders the canonical 3-part layout. colorize=false
// forces plain output (pipes, --no-pretty, generated markdown).
func renderHelpExample(verb string, colorize bool) string {
	accent := ""
	accentBold := ""
	primaryBold := ""
	okColor := ""
	reset := ""

	if colorize {
		accent = termout.RoleSGR(termout.RoleAccent)
		accentBold = termout.RoleSGRBold(termout.RoleAccent)
		primaryBold = termout.RoleSGRBold(termout.RolePrimary)
		okColor = termout.RoleSGR(termout.RoleOK)
		reset = termout.Reset()
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("%sgitmap %s — help example%s\n\n", accentBold, verb, reset))

	ex, found := LookupHelpExample(verb)
	if !found {
		out.WriteString(fmt.Sprintf("%sno help example for %q yet.%s\n\n", okColor, verb, reset))
		out.WriteString(fmt.Sprintf("Try %sgitmap help %s%s for the command's help screen.\n", accent, verb, reset))

		return out.String()
	}

	out.WriteString(fmt.Sprintf("%sWHAT IT DOES%s\n", primaryBold, reset))
	out.WriteString(indentBlock(ex.WhatItDoes) + "\n\n")

	for i, example := range ex.Examples {
		invTitle := "INVOCATION"
		outTitle := "EXPECTED OUTPUT"

		if len(ex.Examples) > 1 {
			invTitle = fmt.Sprintf("INVOCATION (%d)", i+1)
			outTitle = fmt.Sprintf("EXPECTED OUTPUT (%d)", i+1)
		}

		out.WriteString(fmt.Sprintf("%s%s%s\n", primaryBold, invTitle, reset))
		out.WriteString(fmt.Sprintf("  %s$%s %s\n\n", okColor, reset, example.Invocation))
		out.WriteString(fmt.Sprintf("%s%s%s\n", primaryBold, outTitle, reset))
		out.WriteString(indentBlock(example.ExpectedOutput) + "\n")

		if i < len(ex.Examples)-1 {
			out.WriteString("\n")
		}
	}

	return out.String()
}

// indentBlock indents every line of s by two spaces.
func indentBlock(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, line := range lines {
		lines[i] = "  " + line
	}

	return strings.Join(lines, "\n")
}

// colorizeExamples reports whether example output may carry ANSI:
// stdout must be a real TTY and the theme must not be mono.
// prettyOff forces plain output (honors --no-pretty / --pretty=false).
func colorizeExamples(prettyOff bool) bool {
	if prettyOff {
		return false
	}

	if _, strip := termout.ResolveTheme(); strip {
		return false
	}

	return termout.IsStdoutTTY()
}

// PrintHelpExample renders `gitmap help example <verb>` to the dispatch
// writer (theme filter applies on top). With no verb it lists the
// curated major commands.
func PrintHelpExample(verb string, prettyOff bool) {
	fmt.Fprint(output.UI(), renderHelpExample(verb, colorizeExamples(prettyOff)))
}

// PrintHelpExampleIndex lists the major (curated) example verbs for
// bare `gitmap help example`.
func PrintHelpExampleIndex(prettyOff bool) {
	colorize := colorizeExamples(prettyOff)
	accentBold := ""
	primaryBold := ""
	reset := ""

	if colorize {
		accentBold = termout.RoleSGRBold(termout.RoleAccent)
		primaryBold = termout.RoleSGRBold(termout.RolePrimary)
		reset = termout.Reset()
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("%shelp example — available commands%s\n\n", accentBold, reset))
	out.WriteString("Usage:\n  gitmap help example <command>\n\n")
	out.WriteString(fmt.Sprintf("%sMajor commands%s\n", primaryBold, reset))

	for _, verb := range curatedExampleVerbs() {
		out.WriteString("  " + verb + "\n")
	}

	total := len(helpExampleRegistry)
	out.WriteString(fmt.Sprintf("\n%d commands and aliases have examples — try any command name.\n", total))
	fmt.Fprint(output.UI(), out.String())
}

// curatedExampleVerbs lists the hand-written major verbs for the index.
func curatedExampleVerbs() []string {
	verbs := []string{
		"scan", "clone", "pull", "pull-all", "commit", "nodes", "ssh",
		"cluster", "agm", "agy", "install", "update", "prompt", "llm-train",
		"task", "schedule", "fix", "db", "pipeline", "tree", "summary",
		"version", "help", "theme",
	}

	sort.Strings(verbs)

	return verbs
}

// PrintExampleSection appends the registry's Examples section to
// `gitmap help <cmd>` output. Silent when the verb has no entry.
func PrintExampleSection(verb string, prettyOff bool) {
	ex, found := LookupHelpExample(verb)
	if !found || len(ex.Examples) == 0 {
		return
	}

	var out strings.Builder
	if colorizeExamples(prettyOff) {
		out.WriteString(fmt.Sprintf("%sExamples%s\n", termout.RoleSGRBold(termout.RolePrimary), termout.Reset()))
	} else {
		out.WriteString("Examples\n")
	}

	for _, example := range ex.Examples {
		out.WriteString("  $ " + example.Invocation + "\n")
	}

	out.WriteString(fmt.Sprintf("  (full example: gitmap help example %s)\n\n", verb))
	fmt.Fprint(output.UI(), out.String())
}

// RenderHelpExampleMarkdown renders the plain (never ANSI) canonical
// block for the generated helpdoc/*.md topics.
func RenderHelpExampleMarkdown(verb string) string {
	return renderHelpExample(verb, false)
}
