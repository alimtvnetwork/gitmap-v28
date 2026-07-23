package cmd

// Test helpers for cmdconstants_unique_test.go. Kept in a separate _test.go
// file so the main test reads top-down without auxiliary noise and so the
// 200-line per-file guideline holds for both files.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const (
	uniqMarkerTopLevel = "gitmap:cmd top-level"
	uniqMarkerSkip     = "gitmap:cmd skip"
)

// cmdConstantOccurrence records where a single Cmd* constant was declared.
type cmdConstantOccurrence struct {
	Name string
	File string
	Line int
}

// constantsDirForTest resolves the absolute path of gitmap/constants relative
// to this test file so the test runs the same way from `go test ./...` as
// it does from inside an IDE.
func constantsDirForTest(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed; cannot locate constants directory")
	}

	return filepath.Join(filepath.Dir(thisFile), "..", "constants")
}

// collectTopLevelCmdConstants parses one constants file and appends every
// qualifying Cmd* constant occurrence to byValue, keyed by string value.
func collectTopLevelCmdConstants(t *testing.T, path string, byValue map[string][]cmdConstantOccurrence) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		if !commentGroupHas(gen.Doc, uniqMarkerTopLevel) {
			continue
		}

		appendBlockOccurrences(fset, gen, byValue)
	}
}

// appendBlockOccurrences walks every spec in an opted-in const block and
// records each Cmd* string constant that is not marked `gitmap:cmd skip`.
func appendBlockOccurrences(fset *token.FileSet, gen *ast.GenDecl, byValue map[string][]cmdConstantOccurrence) {
	for _, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if commentGroupHas(vs.Comment, uniqMarkerSkip) || commentGroupHas(vs.Doc, uniqMarkerSkip) {
			continue
		}

		recordValueSpec(fset, vs, byValue)
	}
}

// recordValueSpec extracts every Cmd* identifier in the spec and records its
// declared string value plus source location.
func recordValueSpec(fset *token.FileSet, vs *ast.ValueSpec, byValue map[string][]cmdConstantOccurrence) {
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

		pos := fset.Position(name.Pos())
		byValue[val] = append(byValue[val], cmdConstantOccurrence{
			Name: name.Name,
			File: filepath.Base(pos.Filename),
			Line: pos.Line,
		})
	}
}

// commentGroupHas reports whether any comment line in cg contains needle.
func commentGroupHas(cg *ast.CommentGroup, needle string) bool {
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

// reportDuplicateValues fails the test (with a sorted, human-readable list)
// for every string value claimed by more than one distinct Cmd* identifier.
func reportDuplicateValues(t *testing.T, byValue map[string][]cmdConstantOccurrence) {
	t.Helper()

	values := make([]string, 0, len(byValue))
	for v := range byValue {
		values = append(values, v)
	}
	sort.Strings(values)

	hasDup := false
	for _, v := range values {
		occ := byValue[v]
		if !hasDistinctNames(occ) {
			continue
		}

		hasDup = true
		t.Errorf("duplicate top-level Cmd value %q claimed by multiple constants:\n%s",
			v, formatOccurrences(occ))
	}

	if !hasDup {
		t.Logf("verified %d unique top-level Cmd values across %d files",
			len(values), countFiles(byValue))
	}
}

// hasDistinctNames reports whether the slice contains two occurrences with
// different constant names. Identical-name occurrences would already fail at
// `go build` time as redeclarations, so we surface only cross-name collisions.
func hasDistinctNames(occ []cmdConstantOccurrence) bool {
	if len(occ) < 2 {
		return false
	}
	first := occ[0].Name
	for _, o := range occ[1:] {
		if o.Name != first {
			return true
		}
	}

	return false
}

// formatOccurrences renders the duplicate list as one indented line per site.
func formatOccurrences(occ []cmdConstantOccurrence) string {
	var b strings.Builder
	for _, o := range occ {
		b.WriteString("  - ")
		b.WriteString(o.Name)
		b.WriteString(" (")
		b.WriteString(o.File)
		b.WriteString(":")
		b.WriteString(strconv.Itoa(o.Line))
		b.WriteString(")\n")
	}

	return strings.TrimRight(b.String(), "\n")
}

// countFiles returns the number of unique source files that contributed at
// least one Cmd constant — used purely for the success log line.
func countFiles(byValue map[string][]cmdConstantOccurrence) int {
	files := map[string]struct{}{}
	for _, occ := range byValue {
		for _, o := range occ {
			files[o.File] = struct{}{}
		}
	}

	return len(files)
}
