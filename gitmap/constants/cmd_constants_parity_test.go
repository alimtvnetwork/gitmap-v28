package constants

// TestTopLevelCmdRegistryMatchesAST is the §5 future-hardening guard from
// spec/01-app/99-cli-cmd-uniqueness-ci-guard.md. It walks every
// `constants_*.go` file with go/parser, collects every Cmd* string constant
// declared inside a const block marked `// gitmap:cmd top-level` (minus
// per-spec `// gitmap:cmd skip` lines), and asserts the resulting set is
// exactly equal to the manual `topLevelCmds()` registry.
//
// This makes registry drift impossible: adding a new top-level Cmd* constant
// without updating topLevelCmds() (or vice versa) fails CI before any
// collision check can be bypassed by a stale registry.

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
	parityMarkerTopLevel = "gitmap:cmd top-level"
	parityMarkerSkip     = "gitmap:cmd skip"
)

func TestTopLevelCmdRegistryMatchesAST(t *testing.T) {
	astNames := collectTopLevelCmdNamesFromAST(t)

	registryNames := map[string]struct{}{}
	for name := range topLevelCmds() {
		registryNames[name] = struct{}{}
	}

	missingFromRegistry := diffNameSets(astNames, registryNames)
	extraInRegistry := diffNameSets(registryNames, astNames)

	if len(missingFromRegistry) > 0 {
		t.Errorf("AST has %d top-level Cmd constant(s) missing from topLevelCmds() registry:\n  %s\n"+
			"Add them to topLevelCmds() in cmd_constants_test.go, or mark them `// gitmap:cmd skip`.",
			len(missingFromRegistry), strings.Join(missingFromRegistry, "\n  "))
	}
	if len(extraInRegistry) > 0 {
		t.Errorf("topLevelCmds() registry has %d entry/entries with no matching AST declaration:\n  %s\n"+
			"Remove them from topLevelCmds(), or restore the constant under a `// gitmap:cmd top-level` block.",
			len(extraInRegistry), strings.Join(extraInRegistry, "\n  "))
	}
}

// TestNoUnmarkedTopLevelCmdBlocks mirrors the gencommands missed-marker
// guard inside the standard test suite. Any const block declaring one or
// more Cmd* string constants without `// gitmap:cmd top-level` on its doc
// comment is rejected here, so `go test ./...` catches the drift even when
// nobody runs `go generate`. This is what would have caught the v3.31.0
// `has-change` / `hc` regression at PR time.
func TestNoUnmarkedTopLevelCmdBlocks(t *testing.T) {
	dir := constantsDirForParityTest(t)
	files, err := filepath.Glob(filepath.Join(dir, "constants_*.go"))
	if err != nil {
		t.Fatalf("glob constants dir: %v", err)
	}

	var violations []string

	for _, path := range files {
		violations = append(violations, findUnmarkedCmdBlocksForTest(t, path)...)
	}

	if len(violations) == 0 {
		return
	}

	t.Errorf("found %d Cmd* constant(s) in const block(s) lacking `// %s`:\n  %s\n"+
		"Add `// %s` to the block's doc comment, or mark individual specs with `// %s` inside an opted-in block.",
		len(violations), parityMarkerTopLevel,
		strings.Join(violations, "\n  "),
		parityMarkerTopLevel, parityMarkerSkip)
}

func findUnmarkedCmdBlocksForTest(t *testing.T, path string) []string {
	t.Helper()

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	var out []string

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		if parityCommentHas(gen.Doc, parityMarkerTopLevel) {
			continue
		}

		names := unmarkedCmdNamesForTest(gen)
		if len(names) == 0 {
			continue
		}

		pos := fset.Position(gen.Pos())
		out = append(out, filepath.Base(path)+":"+strconv.Itoa(pos.Line)+"  "+strings.Join(names, ", "))
	}

	return out
}

func unmarkedCmdNamesForTest(gen *ast.GenDecl) []string {
	var names []string

	for _, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if parityCommentHas(vs.Comment, parityMarkerSkip) || parityCommentHas(vs.Doc, parityMarkerSkip) {
			continue
		}
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
			names = append(names, name.Name)
		}
	}

	return names
}

// collectTopLevelCmdNamesFromAST returns the set of Cmd* constant names
// declared in any opted-in const block across the constants package.
func collectTopLevelCmdNamesFromAST(t *testing.T) map[string]struct{} {
	t.Helper()

	dir := constantsDirForParityTest(t)
	files, err := filepath.Glob(filepath.Join(dir, "constants_*.go"))
	if err != nil {
		t.Fatalf("glob constants dir: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("no constants_*.go files found under %s", dir)
	}

	names := map[string]struct{}{}
	for _, path := range files {
		collectNamesFromFile(t, path, names)
	}

	return names
}

func collectNamesFromFile(t *testing.T, path string, names map[string]struct{}) {
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
		if !parityCommentHas(gen.Doc, parityMarkerTopLevel) {
			continue
		}
		collectSpecNames(gen, names)
	}
}

func collectSpecNames(gen *ast.GenDecl, names map[string]struct{}) {
	for _, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if parityCommentHas(vs.Comment, parityMarkerSkip) || parityCommentHas(vs.Doc, parityMarkerSkip) {
			continue
		}
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
			names[name.Name] = struct{}{}
		}
	}
}

func parityCommentHas(cg *ast.CommentGroup, needle string) bool {
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

// diffNameSets returns sorted keys present in a but absent from b.
func diffNameSets(a, b map[string]struct{}) []string {
	var out []string
	for k := range a {
		if _, ok := b[k]; !ok {
			out = append(out, k)
		}
	}
	sort.Strings(out)

	return out
}

// constantsDirForParityTest resolves this package's directory so the test
// runs identically from `go test ./...` and from inside an IDE.
func constantsDirForParityTest(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed; cannot locate constants directory")
	}

	return filepath.Dir(thisFile)
}
