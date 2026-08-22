package main

import (
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
			// Find assignments like _ = err
			if assign, ok := n.(*ast.AssignStmt); ok {
				for i, lhs := range assign.Lhs {
					if ident, isIdent := lhs.(*ast.Ident); isIdent && ident.Name == "_" {
						if i < len(assign.Rhs) {
							if rhsIdent, isRhsIdent := assign.Rhs[i].(*ast.Ident); isRhsIdent && rhsIdent.Name == "err" {
								pos := fset.Position(assign.Pos())
								fmt.Printf("%s:%d: Ignored error (_ = err)\n", pos.Filename, pos.Line)
							}
						}
					}
				}
			}

			// Find empty if err != nil blocks
			if ifStmt, ok := n.(*ast.IfStmt); ok {
				if binOp, isBinOp := ifStmt.Cond.(*ast.BinaryExpr); isBinOp {
					if binOp.Op == token.NEQ {
						if ident, ok := binOp.X.(*ast.Ident); ok && ident.Name == "err" {
							if len(ifStmt.Body.List) == 0 {
								pos := fset.Position(ifStmt.Pos())
								fmt.Printf("%s:%d: Empty err check (if err != nil {})\n", pos.Filename, pos.Line)
							}
						}
					}
				}
			}
			return true
		})
		return nil
	})

	if err != nil {
		panic(err)
	}
}
