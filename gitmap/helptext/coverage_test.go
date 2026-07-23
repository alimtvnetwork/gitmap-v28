package helptext

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestEveryCmdIDHasHelpFile parses constants/constants_cli.go and
// asserts that every primary command ID (i.e. not an alias and not
// in the documented exemption list) has a matching <id>.md inside
// this package's embedded helptext file system.
//
// Aliases (CmdFooAlias, CmdFooAlias2) and group subcommands map to
// their primary command's help file, so they're skipped here. The
// goal is: a future PR that adds `CmdFoo = "foo"` MUST also add
// `helptext/foo.md` or this test fails CI.
func TestEveryCmdIDHasHelpFile(t *testing.T) {
	ids, err := parsePrimaryCmdIDs()
	if err != nil {
		t.Fatalf("parse constants_cli.go: %v", err)
	}
	if len(ids) == 0 {
		t.Fatal("parsed zero Cmd* constants — locator broke?")
	}
	var missing []string
	for _, id := range ids {
		if _, lookupErr := files.ReadFile(id.value + ".md"); lookupErr != nil {
			missing = append(missing, fmt.Sprintf("%s (constant %s = %q)", id.value+".md", id.name, id.value))
		}
	}
	if len(missing) > 0 {
		t.Fatalf("commands missing helptext/<id>.md (%d):\n  - %s\n\n"+
			"Add the file or, if intentional (subcommand/internal/runner), append the constant name to helptextExemptConstants in helptext_coverage_test.go.",
			len(missing), strings.Join(missing, "\n  - "))
	}
}

// cmdID is one parsed `CmdFoo = "foo"` declaration.
type cmdID struct {
	name  string // Go identifier, e.g. "CmdScan"
	value string // string literal value, e.g. "scan"
}

// parsePrimaryCmdIDs returns every primary command constant from
// constants/constants_cli.go (aliases and exempt entries skipped).
func parsePrimaryCmdIDs() ([]cmdID, error) {
	path := locateCLIConstantsFile()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	var out []cmdID
	for _, decl := range file.Decls {
		out = append(out, collectFromDecl(decl)...)
	}

	return filterPrimaries(out), nil
}

// collectFromDecl extracts every `CmdX = "y"` constant from one decl.
func collectFromDecl(decl ast.Decl) []cmdID {
	gen, ok := decl.(*ast.GenDecl)
	if !ok || gen.Tok != token.CONST {
		return nil
	}
	var out []cmdID
	for _, spec := range gen.Specs {
		valSpec, ok2 := spec.(*ast.ValueSpec)
		if !ok2 {
			continue
		}
		out = append(out, extractFromSpec(valSpec)...)
	}

	return out
}

// extractFromSpec pulls every `Name = "value"` pair from one ValueSpec.
func extractFromSpec(spec *ast.ValueSpec) []cmdID {
	var out []cmdID
	for i, name := range spec.Names {
		if !strings.HasPrefix(name.Name, "Cmd") || i >= len(spec.Values) {
			continue
		}
		lit, ok := spec.Values[i].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			continue
		}
		out = append(out, cmdID{name: name.Name, value: strings.Trim(lit.Value, `"`)})
	}

	return out
}

// filterPrimaries keeps only IDs that should have their own help file.
func filterPrimaries(in []cmdID) []cmdID {
	exempt := map[string]struct{}{}
	for _, n := range helptextExemptConstants {
		exempt[n] = struct{}{}
	}
	out := make([]cmdID, 0, len(in))
	for _, id := range in {
		if isAliasName(id.name) || isExempt(id.name, exempt) {
			continue
		}
		out = append(out, id)
	}

	return out
}

// isAliasName matches CmdFooAlias / CmdFooAlias2 / CmdFooAlias3 ...
func isAliasName(name string) bool {
	return strings.Contains(name, "Alias")
}

// isExempt reports whether name is in the documented exemption list.
func isExempt(name string, exempt map[string]struct{}) bool {
	_, ok := exempt[name]

	return ok
}

// locateCLIConstantsFile resolves constants/constants_cli.go relative
// to this test file (works regardless of `go test` cwd).
func locateCLIConstantsFile() string {
	_, thisFile, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(thisFile)

	return filepath.Join(pkgDir, "..", "constants", "constants_cli.go")
}

// helptextExemptConstants lists every Cmd* constant whose value is
// intentionally NOT a top-level command and so doesn't need its own
// helptext/<value>.md file. Each entry MUST have a one-line reason.
//
// Add to this list ONLY when shipping a subcommand, internal runner,
// or alias that is documented inside another command's help page.
var helptextExemptConstants = []string{
	// Subcommands of `gitmap group` — documented inside helptext/group.md.
	"CmdGroupCreate", "CmdGroupAdd", "CmdGroupRemove",
	"CmdGroupList", "CmdGroupShow", "CmdGroupDelete",
	// Internal updater runners — never invoked by users directly.
	"CmdUpdateRunner", "CmdRevertRunner",
	// `changelog.md` is a file-format alias for `changelog`.
	"CmdChangelogMD",
	// Maintenance command surfaced only via `gitmap doctor`.
	"CmdDBReset",
	// Pending-review wrappers — covered by helptext/do-pending.md.
	"CmdPending",
	// Internal source-repo override — set via env, not a help page.
	"CmdSetSourceRepo",
	// `help` is the help dispatcher itself — printed by the binary.
	"CmdHelp",
	// `r` is a release short-form alias of `release` — covered by release.md.
	"CmdReleaseShort",
	// `releases` is plural alias of `list-releases` — covered by list-releases.md.
	"CmdReleases",
	// `default` is a `gitmap branch` subcommand verb — covered by branch.md.
	"CmdBranchSubDefault",
}
