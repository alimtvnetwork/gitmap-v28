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
			expected: "Feature: AUTH - support jwt refresh tokens",
			module:   "AUTH",
		},
		{
			name:     "cbf hyphen sqlite",
			verb:     "cbf",
			args:     []string{"sqlite", "-", "handle", "connection", "pool", "lock"},
			expected: "Bug: SQLITE - handle connection pool lock",
			module:   "SQLITE",
		},
		{
			name:     "cpb hyphen db",
			verb:     "cpb",
			args:     []string{"db", "-", "fix", "deadlock"},
			expected: "Bug: DB - fix deadlock",
			module:   "DB",
		},
		{
			name:     "cpc hyphen skills",
			verb:     "cpc",
			args:     []string{"skills", "-", "sync", "gitmap", "skill"},
			expected: "Chore: SKILLS - sync gitmap skill",
			module:   "SKILLS",
		},
		{
			name:     "cpr hyphen version",
			verb:     "cpr",
			args:     []string{"v6.600.0", "-", "release", "notes"},
			expected: "Release: V6.600.0 - release notes",
			module:   "V6.600.0",
		},
		{
			name:     "cfr hyphen refactor",
			verb:     "cfr",
			args:     []string{"toolchain", "-", "extract", "llvm", "resolver"},
			expected: "Refactor: TOOLCHAIN - extract llvm resolver",
			module:   "TOOLCHAIN",
		},
		{
			name:     "module with path slashes",
			verb:     "cpf",
			args:     []string{"cli/cmdworkgroup", "-", "add", "resolver"},
			expected: "Feature: CLI/CMDWORKGROUP - add resolver",
			module:   "CLI/CMDWORKGROUP",
		},
		{
			name:     "colon separator without brackets",
			verb:     "cpf",
			args:     []string{"auth:", "add", "jwt"},
			expected: "Feature: AUTH - add jwt",
			module:   "AUTH",
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
			expected: "Feature: AUTH - support jwt refresh tokens",
			module:   "AUTH",
		},
		{
			name:     "cbf bracket sqlite sanitization",
			verb:     "cbf",
			args:     []string{"[sqlite]:", "handle", "connection", "pool", "lock"},
			expected: "Bug: SQLITE - handle connection pool lock",
			module:   "SQLITE",
		},
		{
			name:     "cpc bracket skills sanitization",
			verb:     "cpc",
			args:     []string{"[skills]:", "sync", "gitmap", "skill"},
			expected: "Chore: SKILLS - sync gitmap skill",
			module:   "SKILLS",
		},
		{
			name:     "cpr bracket version sanitization",
			verb:     "cpr",
			args:     []string{"[v6.600.0]:", "release", "notes"},
			expected: "Release: V6.600.0 - release notes",
			module:   "V6.600.0",
		},
		{
			name:     "bracket module with slashes",
			verb:     "cpf",
			args:     []string{"[cli/cmdworkgroup]:", "add", "resolver"},
			expected: "Feature: CLI/CMDWORKGROUP - add resolver",
			module:   "CLI/CMDWORKGROUP",
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
			expected: "FEAT(WORKGROUP): add folder ssh key binding",
			module:   "WORKGROUP",
		},
		{
			name:     "fix with plain scope",
			verb:     "cbf",
			args:     []string{"fix(ssh):", "handle", "missing", "private", "key"},
			expected: "FIX(SSH): handle missing private key",
			module:   "SSH",
		},
		{
			name:     "chore with plain scope",
			verb:     "cpc",
			args:     []string{"chore(skills):", "cleanup", "unused", "files"},
			expected: "CHORE(SKILLS): cleanup unused files",
			module:   "SKILLS",
		},
		{
			name:     "cbf with bracketed module syntax",
			verb:     "cbf",
			args:     []string{"(constants):", "restore", "untyped", "OrderMode", "aliases"},
			expected: "FIX(CONSTANTS): restore untyped OrderMode aliases",
			module:   "CONSTANTS",
		},
		{
			name:     "cfr with bracketed module syntax",
			verb:     "cfr",
			args:     []string{"(guidelines):", "update", "v6", "execute", "prompt"},
			expected: "REFACTOR(GUIDELINES): update v6 execute prompt",
			module:   "GUIDELINES",
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
			expected: "Feature: AUTH - add jwt",
		},
		{
			name:     "single quotes wrapping",
			verb:     "cbf",
			args:     []string{`'[db]: fix lock'`},
			expected: "Bug: DB - fix lock",
		},
		{
			name:     "duplicate Feature prefix stripping",
			verb:     "cpf",
			args:     []string{"Feature:", "[auth]:", "add", "jwt"},
			expected: "Feature: AUTH - add jwt",
		},
		{
			name:     "duplicate Bug prefix stripping",
			verb:     "cbf",
			args:     []string{"Bug:", "[db]:", "fix", "lock"},
			expected: "Bug: DB - fix lock",
		},
		{
			name:     "legacy hyphen syntax conversion",
			verb:     "cpf",
			args:     []string{"auth", "-", "add", "jwt", "validation"},
			expected: "Feature: AUTH - add jwt validation",
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
		{"cpf", []string{"auth - add jwt"}, "Feature: AUTH - add jwt"},
		{"cbf", []string{"db - fix deadlock"}, "Bug: DB - fix deadlock"},
		{"cpc", []string{"[skills]: cleanup"}, "Chore: SKILLS - cleanup"},
		{"cpr", []string{"v6.605.0 - release"}, "Release: V6.605.0 - release"},
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

	if res.FormattedMessage != "Feature: AUTH - add jwt validation" {
		t.Fatalf("unexpected formatted message: %s", res.FormattedMessage)
	}

	if res.Module != "AUTH" {
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
			expected: "Feature: CLI - add flag",
			module:   "CLI",
		},
		{
			name:     "cpf FEAT prefix",
			verb:     "cpf",
			args:     []string{"FEAT:", "CLI", "-", "add", "flag"},
			expected: "Feature: CLI - add flag",
			module:   "CLI",
		},
		{
			name:     "cpr Release prefix",
			verb:     "cpr",
			args:     []string{"Release:", "V6.605.0", "-", "release", "notes"},
			expected: "Release: V6.605.0 - release notes",
			module:   "V6.605.0",
		},
		{
			name:     "cbf Fix prefix mapped to Bug",
			verb:     "cbf",
			args:     []string{"Fix:", "STORE", "-", "handle", "nil", "db"},
			expected: "Bug: STORE - handle nil db",
			module:   "STORE",
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
			expected: "Feature: CLI - add flag",
			module:   "CLI",
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
			expected: "Feature: HEATMAP - visual overhaul",
			module:   "HEATMAP",
		},
		{
			name:     "cbf FIX prevent non-git commit",
			verb:     "cbf",
			args:     []string{"FIX", "-", "prevent", "non-git", "commit"},
			expected: "Bug: FIX - prevent non-git commit",
			module:   "FIX",
		},
		{
			name:     "cfr STORE optimize query",
			verb:     "cfr",
			args:     []string{"STORE", "-", "optimize", "query"},
			expected: "Refactor: STORE - optimize query",
			module:   "STORE",
		},
		{
			name:     "cpc DOCS update readme",
			verb:     "cpc",
			args:     []string{"DOCS", "-", "update", "readme"},
			expected: "Chore: DOCS - update readme",
			module:   "DOCS",
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
