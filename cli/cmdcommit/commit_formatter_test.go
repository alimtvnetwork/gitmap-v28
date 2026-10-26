package cmdcommit

import (
	"strings"
	"testing"
)

func TestParseAndFormatCommitMessage_HyphenSyntax(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		args     []string
		expected string
		module   string
	}{
		{
			name:     "cpf hyphen auth",
			verb:     "cpf",
			args:     []string{"auth", "-", "support", "jwt", "refresh", "tokens"},
			expected: "Feature: auth - support jwt refresh tokens",
			module:   "auth",
		},
		{
			name:     "cbf hyphen sqlite",
			verb:     "cbf",
			args:     []string{"sqlite", "-", "handle", "connection", "pool", "lock"},
			expected: "Bug: sqlite - handle connection pool lock",
			module:   "sqlite",
		},
		{
			name:     "cpb hyphen db",
			verb:     "cpb",
			args:     []string{"db", "-", "fix", "deadlock"},
			expected: "Bug: db - fix deadlock",
			module:   "db",
		},
		{
			name:     "cpc hyphen skills",
			verb:     "cpc",
			args:     []string{"skills", "-", "sync", "gitmap", "skill"},
			expected: "Chore: skills - sync gitmap skill",
			module:   "skills",
		},
		{
			name:     "cpr hyphen version",
			verb:     "cpr",
			args:     []string{"v6.600.0", "-", "release", "notes"},
			expected: "Release: v6.600.0 - release notes",
			module:   "v6.600.0",
		},
		{
			name:     "cfr hyphen refactor",
			verb:     "cfr",
			args:     []string{"toolchain", "-", "extract", "llvm", "resolver"},
			expected: "Refactor: toolchain - extract llvm resolver",
			module:   "toolchain",
		},
		{
			name:     "module with path slashes",
			verb:     "cpf",
			args:     []string{"cli/cmdworkgroup", "-", "add", "resolver"},
			expected: "Feature: cli/cmdworkgroup - add resolver",
			module:   "cli/cmdworkgroup",
		},
		{
			name:     "colon separator without brackets",
			verb:     "cpf",
			args:     []string{"auth:", "add", "jwt"},
			expected: "Feature: auth - add jwt",
			module:   "auth",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ParseAndFormatCommitMessage(tc.verb, tc.args)
			if res.FormattedMessage != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, res.FormattedMessage)
			}
			if !res.HasModule || res.Module != tc.module {
				t.Fatalf("expected module '%s', got '%s'", tc.module, res.Module)
			}
			if !res.HasHyphenSeparator {
				t.Fatalf("expected HasHyphenSeparator=true")
			}
			if strings.Contains(res.FormattedMessage, "[") || strings.Contains(res.FormattedMessage, "]") {
				t.Fatalf("forbidden square bracket found in formatted message: %s", res.FormattedMessage)
			}
		})
	}
}

func TestParseAndFormatCommitMessage_BracketSanitization(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		args     []string
		expected string
		module   string
	}{
		{
			name:     "cpf bracket auth sanitization",
			verb:     "cpf",
			args:     []string{"[auth]:", "support", "jwt", "refresh", "tokens"},
			expected: "Feature: auth - support jwt refresh tokens",
			module:   "auth",
		},
		{
			name:     "cbf bracket sqlite sanitization",
			verb:     "cbf",
			args:     []string{"[sqlite]:", "handle", "connection", "pool", "lock"},
			expected: "Bug: sqlite - handle connection pool lock",
			module:   "sqlite",
		},
		{
			name:     "cpc bracket skills sanitization",
			verb:     "cpc",
			args:     []string{"[skills]:", "sync", "gitmap", "skill"},
			expected: "Chore: skills - sync gitmap skill",
			module:   "skills",
		},
		{
			name:     "cpr bracket version sanitization",
			verb:     "cpr",
			args:     []string{"[v6.600.0]:", "release", "notes"},
			expected: "Release: v6.600.0 - release notes",
			module:   "v6.600.0",
		},
		{
			name:     "bracket module with slashes",
			verb:     "cpf",
			args:     []string{"[cli/cmdworkgroup]:", "add", "resolver"},
			expected: "Feature: cli/cmdworkgroup - add resolver",
			module:   "cli/cmdworkgroup",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ParseAndFormatCommitMessage(tc.verb, tc.args)
			if res.FormattedMessage != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, res.FormattedMessage)
			}
			if !res.HasModule || res.Module != tc.module {
				t.Fatalf("expected module '%s', got '%s'", tc.module, res.Module)
			}
			if !res.HasHyphenSeparator {
				t.Fatalf("expected HasHyphenSeparator=true")
			}
			if strings.Contains(res.FormattedMessage, "[") || strings.Contains(res.FormattedMessage, "]") {
				t.Fatalf("forbidden square bracket found in formatted message: %s", res.FormattedMessage)
			}
		})
	}
}

func TestParseAndFormatCommitMessage_ConventionalSyntax(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		args     []string
		expected string
		module   string
	}{
		{
			name:     "feat with bracketed scope sanitized",
			verb:     "cpf",
			args:     []string{"feat([workgroup]):", "add", "folder", "ssh", "key", "binding"},
			expected: "feat(workgroup): add folder ssh key binding",
			module:   "workgroup",
		},
		{
			name:     "fix with plain scope",
			verb:     "cbf",
			args:     []string{"fix(ssh):", "handle", "missing", "private", "key"},
			expected: "fix(ssh): handle missing private key",
			module:   "ssh",
		},
		{
			name:     "chore with plain scope",
			verb:     "cpc",
			args:     []string{"chore(skills):", "cleanup", "unused", "files"},
			expected: "chore(skills): cleanup unused files",
			module:   "skills",
		},
		{
			name:     "cbf with bracketed module syntax",
			verb:     "cbf",
			args:     []string{"(constants):", "restore", "untyped", "OrderMode", "aliases"},
			expected: "fix(constants): restore untyped OrderMode aliases",
			module:   "constants",
		},
		{
			name:     "cfr with bracketed module syntax",
			verb:     "cfr",
			args:     []string{"(guidelines):", "update", "v6", "execute", "prompt"},
			expected: "refactor(guidelines): update v6 execute prompt",
			module:   "guidelines",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ParseAndFormatCommitMessage(tc.verb, tc.args)
			if res.FormattedMessage != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, res.FormattedMessage)
			}
			if !res.IsConventional {
				t.Fatalf("expected isConventional=true")
			}
			if strings.Contains(res.FormattedMessage, "[") || strings.Contains(res.FormattedMessage, "]") {
				t.Fatalf("forbidden square bracket found in formatted message: %s", res.FormattedMessage)
			}
		})
	}
}

func TestParseAndFormatCommitMessage_TrimmingAndQuotes(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		args     []string
		expected string
	}{
		{
			name:     "double quotes wrapping",
			verb:     "cpf",
			args:     []string{`"[auth]: add jwt"`},
			expected: "Feature: auth - add jwt",
		},
		{
			name:     "single quotes wrapping",
			verb:     "cbf",
			args:     []string{`'[db]: fix lock'`},
			expected: "Bug: db - fix lock",
		},
		{
			name:     "duplicate Feature prefix stripping",
			verb:     "cpf",
			args:     []string{"Feature:", "[auth]:", "add", "jwt"},
			expected: "Feature: auth - add jwt",
		},
		{
			name:     "duplicate Bug prefix stripping",
			verb:     "cbf",
			args:     []string{"Bug:", "[db]:", "fix", "lock"},
			expected: "Bug: db - fix lock",
		},
		{
			name:     "legacy hyphen syntax conversion",
			verb:     "cpf",
			args:     []string{"auth", "-", "add", "jwt", "validation"},
			expected: "Feature: auth - add jwt validation",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatted := FormatCommitMessage(tc.verb, tc.args)
			if formatted != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, formatted)
			}
			if strings.Contains(formatted, "[") || strings.Contains(formatted, "]") {
				t.Fatalf("forbidden square bracket found: %s", formatted)
			}
		})
	}
}

func TestFormatCommitMessage(t *testing.T) {
	tests := []struct {
		verb     string
		args     []string
		expected string
	}{
		{"cpf", []string{"auth - add jwt"}, "Feature: auth - add jwt"},
		{"cbf", []string{"db - fix deadlock"}, "Bug: db - fix deadlock"},
		{"cpc", []string{"[skills]: cleanup"}, "Chore: skills - cleanup"},
		{"cpr", []string{"v6.605.0 - release"}, "Release: v6.605.0 - release"},
	}

	for _, tc := range tests {
		got := FormatCommitMessage(tc.verb, tc.args)
		if got != tc.expected {
			t.Fatalf("verb %s: expected %q, got %q", tc.verb, tc.expected, got)
		}
	}
}

func TestTryParseLegacyHyphenSyntax(t *testing.T) {
	res, ok := tryParseHyphenSyntax("cpf", "auth - add jwt validation")
	if !ok {
		t.Fatalf("expected tryParseHyphenSyntax ok=true")
	}

	if res.FormattedMessage != "Feature: auth - add jwt validation" {
		t.Fatalf("unexpected formatted message: %s", res.FormattedMessage)
	}

	if res.Module != "auth" {
		t.Fatalf("unexpected module: %s", res.Module)
	}

	if !res.HasHyphenSeparator {
		t.Fatalf("expected HasHyphenSeparator=true")
	}
}

func TestParseAndFormatCommitMessage_RedundantPrefixes(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		args     []string
		expected string
		module   string
	}{
		{
			name:     "cpf Feature prefix",
			verb:     "cpf",
			args:     []string{"Feature:", "CLI", "-", "add", "flag"},
			expected: "Feature: cli - add flag",
			module:   "cli",
		},
		{
			name:     "cpf FEAT prefix",
			verb:     "cpf",
			args:     []string{"FEAT:", "CLI", "-", "add", "flag"},
			expected: "Feature: cli - add flag",
			module:   "cli",
		},
		{
			name:     "cpr Release prefix",
			verb:     "cpr",
			args:     []string{"Release:", "V6.605.0", "-", "release", "notes"},
			expected: "Release: v6.605.0 - release notes",
			module:   "v6.605.0",
		},
		{
			name:     "cbf Fix prefix mapped to Bug",
			verb:     "cbf",
			args:     []string{"Fix:", "STORE", "-", "handle", "nil", "db"},
			expected: "Bug: store - handle nil db",
			module:   "store",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ParseAndFormatCommitMessage(tc.verb, tc.args)
			if res.FormattedMessage != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, res.FormattedMessage)
			}
			if !res.HasModule || res.Module != tc.module {
				t.Fatalf("expected module '%s', got '%s'", tc.module, res.Module)
			}
		})
	}
}

func TestParseAndFormatCommitMessage_PseudoModuleStripping(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		args     []string
		expected string
		module   string
		hasMod   bool
	}{
		{
			name:     "cpf FEATURE stripped without nested module",
			verb:     "cpf",
			args:     []string{"FEATURE", "-", "add", "flag"},
			expected: "Feature: add flag",
			module:   "",
			hasMod:   false,
		},
		{
			name:     "cpf FEATURE stripped with nested CLI module promoted",
			verb:     "cpf",
			args:     []string{"FEATURE", "-", "CLI", "-", "add", "flag"},
			expected: "Feature: cli - add flag",
			module:   "cli",
			hasMod:   true,
		},
		{
			name:     "cpr RELEASE stripped without nested module",
			verb:     "cpr",
			args:     []string{"RELEASE", "-", "v6.605.0"},
			expected: "Release: v6.605.0",
			module:   "",
			hasMod:   false,
		},
		{
			name:     "cpf SPEC stripped without nested module",
			verb:     "cpf",
			args:     []string{"SPEC", "-", "author", "specs"},
			expected: "Feature: author specs",
			module:   "",
			hasMod:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ParseAndFormatCommitMessage(tc.verb, tc.args)
			if res.FormattedMessage != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, res.FormattedMessage)
			}
			if res.HasModule != tc.hasMod || res.Module != tc.module {
				t.Fatalf("expected hasMod=%v, module '%s', got hasMod=%v, module '%s'", tc.hasMod, tc.module, res.HasModule, res.Module)
			}
		})
	}
}

func TestParseAndFormatCommitMessage_CanonicalBehavior(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		args     []string
		expected string
		module   string
	}{
		{
			name:     "cpf HEATMAP visual overhaul",
			verb:     "cpf",
			args:     []string{"HEATMAP", "-", "visual", "overhaul"},
			expected: "Feature: heatmap - visual overhaul",
			module:   "heatmap",
		},
		{
			name:     "cbf FIX prevent non-git commit",
			verb:     "cbf",
			args:     []string{"FIX", "-", "prevent", "non-git", "commit"},
			expected: "Bug: fix - prevent non-git commit",
			module:   "fix",
		},
		{
			name:     "cfr STORE optimize query",
			verb:     "cfr",
			args:     []string{"STORE", "-", "optimize", "query"},
			expected: "Refactor: store - optimize query",
			module:   "store",
		},
		{
			name:     "cpc DOCS update readme",
			verb:     "cpc",
			args:     []string{"DOCS", "-", "update", "readme"},
			expected: "Chore: docs - update readme",
			module:   "docs",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ParseAndFormatCommitMessage(tc.verb, tc.args)
			if res.FormattedMessage != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, res.FormattedMessage)
			}
			if !res.HasModule || res.Module != tc.module {
				t.Fatalf("expected module '%s', got '%s'", tc.module, res.Module)
			}
		})
	}
}
