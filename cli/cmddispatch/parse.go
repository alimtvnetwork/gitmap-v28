package cmddispatch

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
)

// loadConstantStrings maps "constants.CmdDoctor" -> "doctor" by parsing plain
// string const/var declarations in cli/constants/*.go. Only simple string
// literals resolve; anything else is left unresolvable (skipped by callers).
func loadConstantStrings(root string) (map[string]string, error) {
	out := make(map[string]string)
	files, err := filepath.Glob(filepath.Join(root, "cli", "constants", "*.go"))
	if err != nil {
		return nil, err
	}
	for _, path := range files {
		if isTestFile(path) {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			vs, ok := n.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 {
				return true
			}
			lit, ok := vs.Values[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			val, err := strconv.Unquote(lit.Value)
			if err != nil {
				return true
			}
			out["constants."+vs.Names[0].Name] = val
			return false
		})
	}
	return out, nil
}

// collectFromFile extracts dispatch names from one cli/cmd/*.go source file.
// It recognizes two patterns: []dispatchEntry{...} (or bare dispatchEntry{...})
// composite literals, and switch statements on cmd/command inside dispatch*
// functions.
func collectFromFile(path, rel string, consts map[string]string) ([]DispatchName, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	var names []DispatchName
	// Package-level var initializers holding tables (uncommon but possible).
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, val := range vs.Values {
				if comp, ok := val.(*ast.CompositeLit); ok && isDispatchEntrySlice(comp.Type) {
					names = append(names, sliceNames(comp, varName(vs), fset, rel, consts)...)
				}
			}
		}
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		fname := fn.Name.Name
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			comp, ok := n.(*ast.CompositeLit)
			if ok {
				if isDispatchEntrySlice(comp.Type) {
					names = append(names, sliceNames(comp, fname, fset, rel, consts)...)
					return false
				}
				if isDispatchEntry(comp.Type) {
					names = append(names, entryNames(comp, fname, fset, rel, consts)...)
					return false
				}
				return true
			}
			sw, ok := n.(*ast.SwitchStmt)
			if !ok || !strings.HasPrefix(fname, "dispatch") {
				return true
			}
			tag, ok := sw.Tag.(*ast.Ident)
			if !ok || (tag.Name != "cmd" && tag.Name != "command") {
				return true
			}
			names = append(names, switchNames(sw, fname, fset, rel, consts)...)
			return false
		})
	}
	return names, nil
}

// sliceNames extracts names from a []dispatchEntry{...} literal.
func sliceNames(comp *ast.CompositeLit, table string, fset *token.FileSet, rel string, consts map[string]string) []DispatchName {
	var names []DispatchName
	for _, elt := range comp.Elts {
		if entry, ok := elt.(*ast.CompositeLit); ok {
			names = append(names, entryNames(entry, table, fset, rel, consts)...)
		}
	}
	return names
}

// entryNames extracts the names list from one dispatchEntry{...} literal.
// It handles both positional ({[]string{...}, handler}) and keyed
// ({names: []string{...}, handler: ...}) field forms.
func entryNames(entry *ast.CompositeLit, table string, fset *token.FileSet, rel string, consts map[string]string) []DispatchName {
	var slice *ast.CompositeLit
	if len(entry.Elts) == 0 {
		return nil
	}
	if kv, ok := entry.Elts[0].(*ast.KeyValueExpr); ok && kv.Key != nil {
		for _, elt := range entry.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == "names" {
				slice, _ = kv.Value.(*ast.CompositeLit)
				break
			}
		}
	} else {
		slice, _ = entry.Elts[0].(*ast.CompositeLit)
	}
	if slice == nil {
		return nil
	}
	var names []DispatchName
	for _, elt := range slice.Elts {
		name, ok := resolveName(elt, consts)
		if !ok {
			continue
		}
		pos := fset.Position(elt.Pos())
		names = append(names, DispatchName{
			Name:  name,
			Owner: rel + ":" + strconv.Itoa(pos.Line),
			Table: table,
		})
	}
	return names
}

// switchNames extracts case-clause string names from one switch statement.
func switchNames(sw *ast.SwitchStmt, table string, fset *token.FileSet, rel string, consts map[string]string) []DispatchName {
	var names []DispatchName
	for _, stmt := range sw.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		for _, expr := range cc.List {
			name, ok := resolveName(expr, consts)
			if !ok {
				continue
			}
			pos := fset.Position(expr.Pos())
			names = append(names, DispatchName{
				Name:  name,
				Owner: rel + ":" + strconv.Itoa(pos.Line),
				Table: table,
			})
		}
	}
	return names
}

// resolveName resolves a name expression to its string value: either a string
// literal or a constants.Xxx reference resolvable via consts.
func resolveName(expr ast.Expr, consts map[string]string) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", false
		}
		v, err := strconv.Unquote(e.Value)
		if err != nil {
			return "", false
		}
		return v, true
	case *ast.SelectorExpr:
		ident, ok := e.X.(*ast.Ident)
		if !ok {
			return "", false
		}
		if v, ok := consts[ident.Name+"."+e.Sel.Name]; ok {
			return v, true
		}
		return "", false
	}
	return "", false
}

func isDispatchEntrySlice(t ast.Expr) bool {
	arr, ok := t.(*ast.ArrayType)
	if !ok || arr.Len != nil {
		return false
	}
	ident, ok := arr.Elt.(*ast.Ident)
	return ok && ident.Name == "dispatchEntry"
}

func isDispatchEntry(t ast.Expr) bool {
	ident, ok := t.(*ast.Ident)
	return ok && ident.Name == "dispatchEntry"
}

func varName(vs *ast.ValueSpec) string {
	if len(vs.Names) == 1 {
		return "var " + vs.Names[0].Name
	}
	return "var"
}
