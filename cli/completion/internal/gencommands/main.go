// Command gencommands generates allcommands_generated.go by extracting every
// Cmd* string constant from constants files in the sibling constants/ package
// that have explicitly opted into completion via marker comments.
//
// Usage (invoked via go:generate from completion.go):
//
//	go run ./internal/gencommands
//
// # Marker-comment opt-in
//
// The generator scans every ../constants/*.go file automatically — no source
// file allowlist is maintained in the generator. Inclusion is controlled
// locally by domain owners via two marker comments:
//
//   - "gitmap:cmd top-level" placed in the doc comment of a `const (...)`
//     block opts every Cmd* string constant in that block into top-level
//     shell tab-completion.
//
//   - "gitmap:cmd skip" placed as the line comment of a single ValueSpec
//     inside an opted-in block excludes that one entry — used for
//     subcommand IDs that share the Cmd* prefix but should not appear at
//     the top level (e.g. "create" / "add" used by `gitmap group`).
//
// Const blocks without "gitmap:cmd top-level" in their doc comment are
// silently ignored. This is the inverse of the previous design (an explicit
// sourceFiles list + skipNames map) and lets domain owners control inclusion
// without ever touching this generator.
//
// Example:
//
//	// gitmap:cmd top-level
//	// CLI commands.
//	const (
//		CmdScan         = "scan"
//		CmdScanAlias    = "s"
//		CmdGroupCreate  = "create" // gitmap:cmd skip
//	)
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	constantsDir = "../constants"

	markerTopLevel = "gitmap:cmd top-level"
	markerSkip     = "gitmap:cmd skip"
)

func main() {
	if err := checkUnmarkedBlocks(); err != nil {
		fmt.Fprintln(os.Stderr, "gencommands:", err)
		os.Exit(1)
	}

	values, err := collect()
	if err != nil {
		fmt.Fprintln(os.Stderr, "gencommands:", err)
		os.Exit(1)
	}

	if err := writeOutput(values); err != nil {
		fmt.Fprintln(os.Stderr, "gencommands:", err)
		os.Exit(1)
	}
}

// checkUnmarkedBlocks fails the build when a const block declares one or
// more Cmd* string constants without the `gitmap:cmd top-level` marker on
// the block's doc comment. This is the "missed-marker drift" guard: the
// previous design silently skipped such blocks, so a contributor adding a
// new command in an unmarked block would ship completion files out of sync
// with the actual CLI surface (the v3.31.0 `has-change`/`hc` regression is
// the canonical example). Authors must either annotate the block with
// `// gitmap:cmd top-level` or, if the constant intentionally is not a
// top-level command, add `// gitmap:cmd skip` to that individual spec
// (which also requires the enclosing block to be opted in).
func checkUnmarkedBlocks() error {
	files, err := filepath.Glob(filepath.Join(constantsDir, "*.go"))
	if err != nil {
		return fmt.Errorf("glob constants dir: %w", err)
	}

	var violations []string

	for _, rel := range files {
		found, err := findUnmarkedCmdBlocks(rel)
		if err != nil {
			return fmt.Errorf("scan %s: %w", rel, err)
		}

		violations = append(violations, found...)
	}

	if len(violations) == 0 {
		return nil
	}

	return fmt.Errorf("found %d Cmd* constant(s) in const block(s) lacking the `// %s` marker:\n  %s\n"+
		"Add `// %s` to the block's doc comment, or mark individual specs with `// %s` inside an opted-in block.",
		len(violations), markerTopLevel,
		strings.Join(violations, "\n  "),
		markerTopLevel, markerSkip)
}

// findUnmarkedCmdBlocks parses one constants file and returns a slice of
// "file:line  Name1, Name2" descriptions for every const block that holds
// at least one Cmd* string constant but lacks the top-level marker.
func findUnmarkedCmdBlocks(rel string) ([]string, error) {
	abs, err := filepath.Abs(rel)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, abs, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var out []string

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}

		if hasMarker(gen.Doc, markerTopLevel) {
			continue
		}

		names := unmarkedCmdNames(gen)
		if len(names) == 0 {
			continue
		}

		pos := fset.Position(gen.Pos())
		out = append(out, fmt.Sprintf("%s:%d  %s", filepath.Base(abs), pos.Line, strings.Join(names, ", ")))
	}

	return out, nil
}

// unmarkedCmdNames returns Cmd* string-constant names declared in an
// unmarked const block. Any single spec carrying `// gitmap:cmd skip` is
// excluded: that marker is the explicit way to acknowledge a Cmd* name
// that is not a top-level command, but it is only honored alongside a
// block-level top-level marker; outside that context it would be
// indistinguishable from a typo, so we still report the rest of the block.
func unmarkedCmdNames(gen *ast.GenDecl) []string {
	var names []string

	for _, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		if hasMarker(vs.Comment, markerSkip) || hasMarker(vs.Doc, markerSkip) {
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

// collect walks every ../constants/*.go file and returns the deduplicated,
