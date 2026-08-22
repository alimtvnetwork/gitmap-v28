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
			if fn, ok := n.(*ast.FuncDecl); ok {
				start := fset.Position(fn.Pos()).Line
				end := fset.Position(fn.End()).Line
				lines := end - start
				if lines > 15 {
					fmt.Printf("%s:%d: Monolithic function %s exceeds 15 lines (%d lines)\n", path, start, fn.Name.Name, lines)
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
