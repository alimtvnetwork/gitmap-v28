// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// routerConstInfo holds a resolved string constant plus its doc comment.
type routerConstInfo struct {
	value string
	doc   string
}

// locateDispatchSourceDir finds the cli/ source tree that owns the command
// router at runtime. It honors the baked-in constants.RepoPath first (the
// repo's own convention for locating its source), then GITMAP_DISPATCH_SRC,
// then walks up from the running binary and the working directory. Returns
// "" when the source tree is unavailable (e.g. installed binary), in which
// case discovery falls back to the generated completion registry.
func locateDispatchSourceDir() string {
	roots := []string{strings.TrimSpace(constants.RepoPath), strings.TrimSpace(os.Getenv("GITMAP_DISPATCH_SRC"))}
	for _, root := range roots {
		if root == "" {
			continue
		}
		if dir := dispatchDirFromRepoRoot(root); dir != "" {
			return dir
		}
	}
	if exe, err := os.Executable(); err == nil {
		if dir := searchUpwardForDispatch(filepath.Dir(exe), 5); dir != "" {
			return dir
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		if dir := searchUpwardForDispatch(cwd, 5); dir != "" {
			return dir
		}
	}
	return ""
}

// dispatchDirFromRepoRoot resolves a repo-root-style path to the cli/ dir
// that holds the router, accepting either the repo root or cli/ itself.
func dispatchDirFromRepoRoot(root string) string {
	if hasDispatchTable(root) {
		return root
	}
	if cliDir := filepath.Join(root, "cli"); hasDispatchTable(cliDir) {
		return cliDir
	}
	return ""
}

func searchUpwardForDispatch(start string, maxUp int) string {
	dir := start
	for i := 0; i <= maxUp; i++ {
		if found := dispatchDirFromRepoRoot(dir); found != "" {
			return found
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
	return ""
}

func hasDispatchTable(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "cmd", "rootdispatch.go"))
	return err == nil && !info.IsDir()
}

// parseCommandConstants resolves string constants in cli/constants/*.go so
// dispatch-table references like constants.CmdScan become "scan". Doc
// comments are captured for subcommand description enrichment.
func parseCommandConstants(dir string) map[string]routerConstInfo {
	out := make(map[string]routerConstInfo)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		parseConstantsFile(filepath.Join(dir, name), out)
	}
	return out
}

func parseConstantsFile(path string, out map[string]routerConstInfo) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return
	}
	for _, decl := range f.Decls {
		gd, isGen := decl.(*ast.GenDecl)
		if !isGen || gd.Tok != token.CONST {
			continue
		}
		for _, spec := range gd.Specs {
			vs, isVal := spec.(*ast.ValueSpec)
			if !isVal {
				continue
			}
			collectConstSpecs(vs, out)
		}
	}
}

func collectConstSpecs(vs *ast.ValueSpec, out map[string]routerConstInfo) {
	doc := strings.TrimSpace(vs.Doc.Text())
	for i, ident := range vs.Names {
		if i >= len(vs.Values) {
			break
		}
		val, hasVal := constStringValue(vs.Values[i])
		if hasVal {
			out[ident.Name] = routerConstInfo{value: val, doc: doc}
		}
	}
}

// constStringValue unquotes a string literal expression.
func constStringValue(expr ast.Expr) (string, bool) {
	lit, isLit := expr.(*ast.BasicLit)
	if !isLit || lit.Kind != token.STRING {
		return "", false
	}
	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return val, true
}

// forEachRouterGoFile parses every non-test Go file under dir whose raw
// bytes contain marker, invoking visit with the parsed AST.
func forEachRouterGoFile(dir, marker string, visit func(path string, f *ast.File)) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(dir, name)
		src, err := os.ReadFile(path)
		if err != nil || !strings.Contains(string(src), marker) {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			continue
		}
		visit(path, f)
	}
}

// isSliceOfComposite reports whether cl is a composite literal of the named
// slice element type: []dispatchEntry, []string, or cobra.Command via selector.
