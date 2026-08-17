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
	rootDir := `d:\work\gitmap\gitmap`
	fset := token.NewFileSet()

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor") {
			return nil
		}

		node, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil // ignore parse errors
		}

		ast.Inspect(node, func(n ast.Node) bool {
			// Find IfStmt
			if ifStmt, ok := n.(*ast.IfStmt); ok {
				// Check inside the Body of the IfStmt
				ast.Inspect(ifStmt.Body, func(innerNode ast.Node) bool {
					if innerIfStmt, innerOk := innerNode.(*ast.IfStmt); innerOk {
						// We found a nested if statement!
						pos := fset.Position(innerIfStmt.Pos())
						fmt.Printf("%s:%d: Nested if statement found\n", pos.Filename, pos.Line)
						// don't need to traverse further inside this nested if to find more, but we could
						return false
					}
					// Traverse block statements
					if _, ok := innerNode.(*ast.BlockStmt); ok {
						return true
					}
					// Stop traversal for other nodes, wait actually, an if could be inside a for loop which is inside the if block.
					// We should traverse down anything EXCEPT function declarations so we don't catch an if inside a closure defined inside an if.
					if _, ok := innerNode.(*ast.FuncLit); ok {
						return false
					}
					return true
				})

				// Also check the Else part
				if ifStmt.Else != nil {
					// If the else is another if statement (i.e. `else if`), it's NOT a nested if in the bad sense,
					// but its BODY could contain nested ifs, which will be caught when that `else if` is visited by the outer ast.Inspect anyway.
					// If the else is a block `else { ... }`, we should inspect it.
					if elseBlock, ok := ifStmt.Else.(*ast.BlockStmt); ok {
						ast.Inspect(elseBlock, func(innerNode ast.Node) bool {
							if innerIfStmt, innerOk := innerNode.(*ast.IfStmt); innerOk {
								pos := fset.Position(innerIfStmt.Pos())
								fmt.Printf("%s:%d: Nested if statement found (in else block)\n", pos.Filename, pos.Line)
								return false
							}
							if _, ok := innerNode.(*ast.BlockStmt); ok {
								return true
							}
							if _, ok := innerNode.(*ast.FuncLit); ok {
								return false
							}
							return true
						})
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
