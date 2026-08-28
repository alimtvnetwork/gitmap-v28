package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cmdDir := "gitmap/cmd"
	files, err := filepath.Glob(filepath.Join(cmdDir, "*.go"))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		modified := false

		ast.Inspect(node, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if strings.HasPrefix(fn.Name.Name, "run") && fn.Name.Name != "runDispatchTable" && fn.Name.Name != "runDispatch" {
				ast.Inspect(fn.Body, func(n2 ast.Node) bool {
					ret, ok := n2.(*ast.ReturnStmt)
					if ok && len(ret.Results) == 0 {
						ret.Results = []ast.Expr{&ast.Ident{Name: "nil"}}
						modified = true
					}
					return true
				})
			}
			return true
		})

		if modified {
			var buf bytes.Buffer
			if err := format.Node(&buf, fset, node); err != nil {
				panic(err)
			}
			os.WriteFile(file, buf.Bytes(), 0644)
			fmt.Println("Refactored bare returns in:", file)
		}
	}
}
