// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

func isSliceOf(cl *ast.CompositeLit, elem string) bool {
	arr, isArr := cl.Type.(*ast.ArrayType)
	if !isArr {
		return false
	}
	ident, isIdent := arr.Elt.(*ast.Ident)
	return isIdent && ident.Name == elem
}

// stringListLiteral extracts string values from a []string{...} literal,
// resolving constants.X selectors through the parsed constants table.
func stringListLiteral(expr ast.Expr, consts map[string]routerConstInfo) []string {
	cl, isLit := expr.(*ast.CompositeLit)
	if !isLit || !isSliceOf(cl, "string") {
		return nil
	}
	var out []string
	for _, elt := range cl.Elts {
		if s, hasStr := literalStringValue(elt, consts); hasStr {
			out = append(out, s)
		}
	}
	return out
}

func literalStringValue(expr ast.Expr, consts map[string]routerConstInfo) (string, bool) {
	if val, hasVal := constStringValue(expr); hasVal {
		return val, true
	}
	sel, isSel := expr.(*ast.SelectorExpr)
	if !isSel {
		return "", false
	}
	// constants.CmdXxx selector — resolve via the parsed constants.
	pkg, isIdent := sel.X.(*ast.Ident)
	if !isIdent || pkg.Name != "constants" {
		return "", false
	}
	info, hasInfo := consts[sel.Sel.Name]
	if !hasInfo {
		return "", false
	}
	return info.value, true
}

// parseDispatchTables extracts every {names, handler} entry from the
// []dispatchEntry routing tables in cli/cmd, recording the enclosing
// dispatch function for category derivation.
func parseDispatchTables(cmdDir string, consts map[string]routerConstInfo) []discoveredCommand {
	var out []discoveredCommand
	forEachRouterGoFile(cmdDir, "dispatchEntry", func(_ string, f *ast.File) {
		for _, decl := range f.Decls {
			fn, isFn := decl.(*ast.FuncDecl)
			if !isFn || fn.Body == nil {
				continue
			}
			category := dispatchGroupCategory(fn.Name.Name)
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				cl, isLit := n.(*ast.CompositeLit)
				if !isLit || !isSliceOf(cl, "dispatchEntry") {
					return true
				}
				for _, elt := range cl.Elts {
					entry, isEntry := elt.(*ast.CompositeLit)
					if !isEntry || len(entry.Elts) == 0 {
						continue
					}
					names := stringListLiteral(entry.Elts[0], consts)
					if len(names) == 0 {
						continue
					}
					out = append(out, discoveredCommand{names: names, category: category, kind: "dispatch"})
				}
				return false
			})
		}
	})
	return out
}

// parseCobraCommands extracts &cobra.Command{Use, Aliases, Short} registrations
// in cli/cmd; Short is the canonical one-line help text for the command.
func parseCobraCommands(cmdDir string) []discoveredCommand {
	var out []discoveredCommand
	forEachRouterGoFile(cmdDir, "cobra.Command", func(_ string, f *ast.File) {
		ast.Inspect(f, func(n ast.Node) bool {
			cl, isLit := n.(*ast.CompositeLit)
			if !isLit || !isCobraCommandLit(cl) {
				return true
			}
			use, aliases, short := cobraCommandFields(cl)
			name := firstUseToken(use)
			if name == "" {
				return false
			}
			names := append([]string{name}, aliases...)
			out = append(out, discoveredCommand{names: names, description: short, category: "tooling", kind: "cobra"})
			return false
		})
	})
	return out
}

func isCobraCommandLit(cl *ast.CompositeLit) bool {
	sel, isSel := cl.Type.(*ast.SelectorExpr)
	if !isSel {
		return false
	}
	pkg, isIdent := sel.X.(*ast.Ident)
	return isIdent && pkg.Name == "cobra" && sel.Sel.Name == "Command"
}

func cobraCommandFields(cl *ast.CompositeLit) (use string, aliases []string, short string) {
	for _, elt := range cl.Elts {
		kv, isKV := elt.(*ast.KeyValueExpr)
		if !isKV {
			continue
		}
		key, isIdent := kv.Key.(*ast.Ident)
		if !isIdent {
			continue
		}
		switch key.Name {
		case "Use":
			if lit, isLit := kv.Value.(*ast.BasicLit); isLit && lit.Kind == token.STRING {
				use, _ = strconv.Unquote(lit.Value)
			}
		case "Short":
			if lit, isLit := kv.Value.(*ast.BasicLit); isLit && lit.Kind == token.STRING {
				short, _ = strconv.Unquote(lit.Value)
			}
		case "Aliases":
			if list, isList := kv.Value.(*ast.CompositeLit); isList {
				for _, item := range list.Elts {
					if lit, isLit := item.(*ast.BasicLit); isLit && lit.Kind == token.STRING {
						if s, err := strconv.Unquote(lit.Value); err == nil {
							aliases = append(aliases, s)
						}
					}
				}
			}
		}
	}
	return use, aliases, short
}

// firstUseToken takes "pe [path|alias|url] [flags]" to "pe".
func firstUseToken(use string) string {
	fields := strings.Fields(use)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
