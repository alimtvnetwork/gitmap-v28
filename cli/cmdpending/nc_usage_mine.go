// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// usageCommandLine matches `gitmap <parent> [<sub> ...]` inside usage text.
var usageCommandLine = regexp.MustCompile(`gitmap\s+([a-z][a-z0-9_-]*)((?:\s+[^\s\[]+)*)`)

// usageAliasHint matches "(alias: gitmap te all)" style alias declarations.
var usageAliasHint = regexp.MustCompile(`\(alias:\s*(?:gitmap\s+)?([a-z][a-z0-9_-]*(?:\s+[a-z][a-z0-9_-]*)*)`)

// usageSectionHeader matches "Commands:" / "Subcommands:" section headers.
var usageSectionHeader = regexp.MustCompile(`(?i)(?:commands|subcommands):\s*$`)

// usageSubcommandLine matches "  <sub> [flags...]   <description>".
var usageSubcommandLine = regexp.MustCompile(`^ {2}([a-z][a-z0-9_-]*)((?:\s+\[[^\]]*\])*)\s{2,}(\S.*)$`)

// mineSubcommandUsage discovers multi-word commands (e.g. "spec next",
// "te all") from the usage/help text the router's own subsystems print.
// Subcommand descriptions come from the Commands:/Subcommands: sections.
func mineSubcommandUsage(cliDir string) []discoveredCommand {
	var out []discoveredCommand
	for _, dir := range subsystemDirs(cliDir) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			out = append(out, mineUsageFile(filepath.Join(dir, name))...)
		}
	}
	return out
}

// subsystemDirs returns cli/cmd plus every cli/cmd*/ subsystem package dir.
func subsystemDirs(cliDir string) []string {
	dirs := []string{filepath.Join(cliDir, "cmd")}
	entries, err := os.ReadDir(cliDir)
	if err != nil {
		return dirs
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "cmd") && e.Name() != "cmd" {
			dirs = append(dirs, filepath.Join(cliDir, e.Name()))
		}
	}
	return dirs
}

func mineUsageFile(path string) []discoveredCommand {
	src, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(src), "Usage:") {
		return nil
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, 0)
	if err != nil {
		return nil
	}
	// Collect every string literal in source order; usage printers often
	// emit one fmt.Println per line, so section state (Commands:) must
	// persist across literals within the file.
	var literals []string
	ast.Inspect(f, func(n ast.Node) bool {
		lit, isLit := n.(*ast.BasicLit)
		if !isLit || lit.Kind != token.STRING {
			return true
		}
		if text, err := strconv.Unquote(lit.Value); err == nil {
			literals = append(literals, text)
		}
		return true
	})
	return mineUsageLiterals(literals)
}

// minedCompound is one `gitmap <parent> <sub>` invocation found in usage text.
type minedCompound struct {
	parent string
	sub    string
	alias  string
}

// usageScanner carries section state across the string literals of one file.
type usageScanner struct {
	compounds  []minedCompound
	subDocs    map[string]string
	inSection  bool
	pendingSub string
	pendingDes strings.Builder
}

func mineUsageLiterals(literals []string) []discoveredCommand {
	sc := &usageScanner{subDocs: make(map[string]string)}
	for _, text := range literals {
		sc.scanLiteral(text)
	}
	sc.finish()
	var out []discoveredCommand
	for _, c := range sc.compounds {
		out = append(out, sc.buildEntry(c)...)
	}
	return out
}

func (sc *usageScanner) scanLiteral(text string) {
	if strings.Contains(text, "Usage:") {
		// New help block: section docs from a previous printer don't apply.
		sc.flushPendingSub()
		sc.inSection = false
	}
	for _, line := range strings.Split(text, "\n") {
		sc.scanLine(line)
	}
}

func (sc *usageScanner) scanLine(line string) {
	trimmed := strings.TrimSpace(line)
	if usageSectionHeader.MatchString(trimmed) {
		sc.flushPendingSub()
		sc.inSection = true
		return
	}
	if sc.inSection {
		if sub, desc, isSub := parseSubcommandDocLine(line); isSub {
			sc.flushPendingSub()
			sc.pendingSub = sub
			sc.pendingDes.WriteString(desc)
			return
		}
		if isContinuationLine(line, sc.pendingSub) {
			sc.pendingDes.WriteString(" " + trimmed)
			return
		}
		sc.flushPendingSub()
		if trimmed == "" {
			return
		}
		sc.inSection = false
	}
	sc.compounds = append(sc.compounds, scanUsageCommandLine(line, sc.compounds)...)
}

func (sc *usageScanner) flushPendingSub() {
	flushPendingSub(sc.subDocs, &sc.pendingSub, &sc.pendingDes)
}

func (sc *usageScanner) finish() {
	sc.flushPendingSub()
}

func (sc *usageScanner) buildEntry(c minedCompound) []discoveredCommand {
	name := c.parent + " " + c.sub
	entry := discoveredCommand{
		names:    []string{name},
		category: dispatchGroupCategory(c.parent),
		kind:     "usage",
	}
	if desc, hasDesc := sc.subDocs[c.sub]; hasDesc {
		entry.description = desc
	}
	if c.alias != "" && c.alias != name {
		entry.names = append(entry.names, c.alias)
	} else if hyphenAlias := strings.ReplaceAll(name, " ", "-"); hyphenAlias != name {
		entry.names = append(entry.names, hyphenAlias)
	}
	out := []discoveredCommand{entry}
	// A spaced alias ("te all") is itself a real invocation: surface it as
	// its own entry pointing back at the primary compound name.
	if c.alias != "" && strings.Contains(c.alias, " ") {
		out = append(out, discoveredCommand{
			names:       []string{c.alias},
			aliasText:   name,
			description: entry.description,
			category:    entry.category,
			kind:        "usage",
		})
	}
	return out
}

func flushPendingSub(docs map[string]string, pendingSub *string, pendingDesc *strings.Builder) {
	if *pendingSub != "" {
		if _, exists := docs[*pendingSub]; !exists {
			docs[*pendingSub] = strings.TrimSpace(pendingDesc.String())
		}
	}
	*pendingSub = ""
	pendingDesc.Reset()
}
