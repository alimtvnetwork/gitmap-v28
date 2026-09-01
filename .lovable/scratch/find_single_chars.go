package main

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"

	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rootDir := `d:\wp-work\riseup-asia\gitmap`
	fset := token.NewFileSet()

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor") || strings.Contains(path, "node_modules") {
			return nil
		}

		node, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}

		ast.Inspect(node, func(n ast.Node) bool {
			// Find variables with 1 character
			if ident, ok := n.(*ast.Ident); ok {
				if len(ident.Name) == 1 && ident.Name != "_" {
					// We must make sure it's a variable declaration or function parameter, not just usage
					// Actually, checking its Pos is tricky without context, let's just inspect ValueSpec
				}
			}

			if vspec, ok := n.(*ast.ValueSpec); ok {
				for _, name := range vspec.Names {
					if len(name.Name) == 1 && name.Name != "_" {
						pos := fset.Position(name.Pos())
						fmt.Printf("%s:%d: Single-character variable '%s'\n", pos.Filename, pos.Line, name.Name)
					}
				}
			}

			if assign, ok := n.(*ast.AssignStmt); ok && assign.Tok == token.DEFINE {
				for _, expr := range assign.Lhs {
					if ident, ok := expr.(*ast.Ident); ok {
						if len(ident.Name) == 1 && ident.Name != "_" {
							pos := fset.Position(ident.Pos())
							fmt.Printf("%s:%d: Single-character variable '%s'\n", pos.Filename, pos.Line, ident.Name)
						}
					}
				}
			}

			if rangeStmt, ok := n.(*ast.RangeStmt); ok {
				if ident, ok := rangeStmt.Key.(*ast.Ident); ok {
					if len(ident.Name) == 1 && ident.Name != "_" {
						pos := fset.Position(ident.Pos())
						fmt.Printf("%s:%d: Single-character variable '%s'\n", pos.Filename, pos.Line, ident.Name)
					}
				}
				if ident, ok := rangeStmt.Value.(*ast.Ident); ok {
					if len(ident.Name) == 1 && ident.Name != "_" {
						pos := fset.Position(ident.Pos())
						fmt.Printf("%s:%d: Single-character variable '%s'\n", pos.Filename, pos.Line, ident.Name)
					}
				}
			}
			return true
		})
		return nil
	})

	if err != nil {
		panic(apperror.New("scratch failure", "ERR_SCRATCH", map[string]any{"err": err}))
	}
}
