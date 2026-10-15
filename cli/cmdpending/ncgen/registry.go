package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
)

// registeredCommand is one command surface discovered in the dispatcher
// registration: every name/alias routing to the same handler.
type registeredCommand struct {
	canonical string   // first name in the registration
	names     []string // all registered names/aliases, canonical first
	short     string   // cobra Short help text, when available
}

// registry maps every registered name AND alias to its command entry.
type registry map[string]*registeredCommand

// buildRegistry parses every non-test Go file in cmdDir and extracts the
// command surface from dispatchEntry routing tables and cobra.Command
// literals. consts resolves constants.CmdFoo selector references.
func buildRegistry(cmdDir string, consts map[string]string) (registry, error) {
	files, err := filepath.Glob(filepath.Join(cmdDir, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("glob cmd dir: %w", err)
	}

	reg := registry{}

	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		if err := scanDispatchFile(path, consts, reg); err != nil {
			return nil, err
		}
	}

	return reg, nil
}

// scanDispatchFile parses one file and registers every dispatchEntry table
// and cobra.Command literal it contains.
func scanDispatchFile(path string, consts map[string]string, reg registry) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		switch {
		case isDispatchTable(lit):
			addDispatchEntries(lit, consts, reg)
		case isCobraCommand(lit):
			addCobraCommand(lit, reg)
		}

		return true
	})

	return nil
}

// isDispatchTable reports whether lit is a []dispatchEntry{...} literal.
func isDispatchTable(lit *ast.CompositeLit) bool {
	arr, ok := lit.Type.(*ast.ArrayType)
	if !ok {
		return false
	}

	id, ok := arr.Elt.(*ast.Ident)
	return ok && id.Name == "dispatchEntry"
}

// addDispatchEntries registers every {[]string{...}, handler} element of a
// dispatch table. The first element is the names list; the rest (handler
// closures) are ignored.
func addDispatchEntries(lit *ast.CompositeLit, consts map[string]string, reg registry) {
	for _, e := range lit.Elts {
		el, ok := e.(*ast.CompositeLit)
		if !ok || len(el.Elts) == 0 {
			continue
		}

		namesLit, ok := el.Elts[0].(*ast.CompositeLit)
		if !ok {
			continue
		}

		addRegistered(reg, resolveNameList(namesLit, consts), "")
	}
}

// resolveNameList converts a []string{...} literal to values, resolving
// constants.CmdFoo selectors through consts.
func resolveNameList(lit *ast.CompositeLit, consts map[string]string) []string {
	var out []string

	for _, e := range lit.Elts {
		switch v := e.(type) {
		case *ast.BasicLit:
			if s := stringLit(v); s != "" {
				out = append(out, s)
			}
		case *ast.SelectorExpr:
			id, ok := v.X.(*ast.Ident)
			if !ok || id.Name != "constants" {
				continue
			}

			if s, ok := consts[v.Sel.Name]; ok && s != "" {
				out = append(out, s)
			}
		}
	}

	return out
}

// isCobraCommand reports whether lit is a cobra.Command{...} literal
// (including the &cobra.Command{...} form: the address-of applies outside
// the composite literal itself).
func isCobraCommand(lit *ast.CompositeLit) bool {
	sel, ok := lit.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == "cobra" && sel.Sel.Name == "Command"
}

// addCobraCommand registers one cobra.Command literal. The command name is
// the first token of Use ("rerun [target]" -> "rerun").
func addCobraCommand(lit *ast.CompositeLit, reg registry) {
	var use, short string

	var aliases []string

	for _, e := range lit.Elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		switch key.Name {
		case "Use":
			use = firstToken(stringLit(kv.Value))
		case "Short":
			short = stringLit(kv.Value)
		case "Aliases":
			if al, ok := kv.Value.(*ast.CompositeLit); ok {
				aliases = append(aliases, resolveNameList(al, nil)...)
			}
		}
	}

	if use == "" {
		return
	}

	addRegistered(reg, append([]string{use}, aliases...), short)
}

// addRegistered inserts one command surface. Every name and alias maps to
// the same entry. The same command is sometimes registered in more than one
// dispatch table with partial name lists (e.g. "replace" appears once bare
// and once with its alias), so collisions merge the name lists instead of
// dropping the later one; the first-seen canonical name is kept.
func addRegistered(reg registry, names []string, short string) {
	uniq := dedupeStrings(names)
	if len(uniq) == 0 {
		return
	}

	rc := &registeredCommand{canonical: uniq[0], names: uniq, short: short}

	for _, n := range uniq {
		if existing, exists := reg[n]; exists {
			mergeRegistered(existing, rc)
			rc = existing
			continue
		}

		reg[n] = rc
	}
}

// mergeRegistered unions src's names into dst, preserving dst's canonical
// order and repointing every name at the merged entry.
func mergeRegistered(dst, src *registeredCommand) {
	for _, n := range src.names {
		if !containsString(dst.names, n) {
			dst.names = append(dst.names, n)
		}
	}

	if dst.short == "" {
		dst.short = src.short
	}
}

func containsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}

	return false
}

func dedupeStrings(in []string) []string {
	seen := map[string]bool{}

	var out []string

	for _, s := range in {
		if s == "" || seen[s] {
			continue
		}

		seen[s] = true
		out = append(out, s)
	}

	return out
}

// stringLit unquotes a string BasicLit, returning "" for anything else.
func stringLit(e ast.Expr) string {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}

	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return ""
	}

	return s
}

func firstToken(s string) string {
	if i := strings.IndexAny(s, " \t"); i >= 0 {
		return s[:i]
	}

	return s
}
