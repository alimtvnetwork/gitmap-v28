package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// loadConstants maps every Cmd* string constant identifier in dir to its
// value, so dispatch tables referencing constants.CmdFoo can be resolved.
func loadConstants(dir string) (map[string]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("glob constants dir: %w", err)
	}

	out := map[string]string{}

	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		if err := scanConstantsFile(path, out); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// scanConstantsFile parses one constants file and records every Cmd*
// string constant's identifier and value.
func scanConstantsFile(path string, out map[string]string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		vs, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}

		for i, name := range vs.Names {
			if !strings.HasPrefix(name.Name, "Cmd") || i >= len(vs.Values) {
				continue
			}

			if s := stringLit(vs.Values[i]); s != "" {
				out[name.Name] = s
			}
		}

		return true
	})

	return nil
}
