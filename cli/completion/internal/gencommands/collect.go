package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func collect() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(constantsDir, "*.go"))
	if err != nil {
		return nil, apperror.WrapSimple(err, "glob constants dir")
	}

	if len(files) == 0 {
		return nil, apperror.NewNotFoundError(fmt.Sprintf("no constants files found in %s", constantsDir))
	}

	seen := map[string]bool{}

	for _, rel := range files {
		if err := scanFile(rel, seen); err != nil {
			return nil, apperror.WrapSimple(err, "scan "+rel)
		}
	}

	out := make([]string, 0, len(seen))
	for v := range seen {
		out = append(out, v)
	}

	sort.Strings(out)

	return out, nil
}

// scanFile parses a single Go source file and adds every qualifying Cmd*
// string constant value to seen. Only const blocks whose doc comment contains
// the markerTopLevel sentinel are considered; specs with a markerSkip line
// comment inside an opted-in block are excluded.
func scanFile(rel string, seen map[string]bool) error {
	abs, err := filepath.Abs(rel)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, abs, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}

		if !hasMarker(gen.Doc, markerTopLevel) {
			continue
		}

		collectFromConstBlock(gen, seen)
	}

	return nil
}

// hasMarker reports whether any line of the comment group contains needle.
func hasMarker(cg *ast.CommentGroup, needle string) bool {
	if cg == nil {
		return false
	}

	for _, c := range cg.List {
		if strings.Contains(c.Text, needle) {
			return true
		}
	}

	return false
}

// collectFromConstBlock walks a single const ( ... ) block.
func collectFromConstBlock(gen *ast.GenDecl, seen map[string]bool) {
	for _, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		if hasMarker(vs.Comment, markerSkip) || hasMarker(vs.Doc, markerSkip) {
			continue
		}

		extractValueSpec(vs, seen)
	}
}

// extractValueSpec records the string literal value of every Cmd* constant
// in the spec whose value is a non-empty string literal.
func extractValueSpec(vs *ast.ValueSpec, seen map[string]bool) {
	for i, name := range vs.Names {
		if !strings.HasPrefix(name.Name, "Cmd") || i >= len(vs.Values) {
			continue
		}

		lit, ok := vs.Values[i].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			continue
		}

		val, err := strconv.Unquote(lit.Value)
		if err != nil || val == "" {
			continue
		}

		seen[val] = true
	}
}
