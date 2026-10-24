package cmdcommit

import (
	"testing"
)

func TestParseAndFormatCommitMessage_BracketSyntax(t *testing.T) {
	tests := []struct {
		name     string
		verb     string
		args     []string
		expected string
		module   string
	}{
		{
			name:     "cpf bracket auth",
			verb:     "cpf",
			args:     []string{"[auth]:", "support", "jwt", "refresh", "tokens"},
			expected: "Feature: [auth]: support jwt refresh tokens",
			module:   "auth",
		},
		{
			name:     "cbf bracket sqlite",
			verb:     "cbf",
			args:     []string{"[sqlite]:", "handle", "connection", "pool", "lock"},
			expected: "Bug: [sqlite]: handle connection pool lock",
			module:   "sqlite",
		},
		{
			name:     "cpb bracket sqlite",
			verb:     "cpb",
			args:     []string{"[db]:", "fix", "deadlock"},
			expected: "Bug: [db]: fix deadlock",
			module:   "db",
		},
		{
			name:     "cpc bracket skills",
			verb:     "cpc",
			args:     []string{"[skills]:", "sync", "gitmap", "skill"},
			expected: "Chore: [skills]: sync gitmap skill",
			module:   "skills",
		},
		{
			name:     "cpr bracket version",
			verb:     "cpr",
			args:     []string{"[v6.600.0]:", "release", "notes"},
			expected: "Release: [v6.600.0]: release notes",
			module:   "v6.600.0",
		},
		{
			name:     "module with path slashes",
			verb:     "cpf",
			args:     []string{"[cli/cmdworkgroup]:", "add", "resolver"},
			expected: "Feature: [cli/cmdworkgroup]: add resolver",
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
			if !res.HasBracketSyntax {
				t.Fatalf("expected hasBracketSyntax=true")
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
			name:     "feat with bracketed scope",
			verb:     "cpf",
			args:     []string{"feat([workgroup]):", "add", "folder", "ssh", "key", "binding"},
			expected: "feat([workgroup]): add folder ssh key binding",
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
			expected: "Feature: [auth]: add jwt",
		},
		{
			name:     "single quotes wrapping",
			verb:     "cbf",
			args:     []string{`'[db]: fix lock'`},
			expected: "Bug: [db]: fix lock",
		},
		{
			name:     "duplicate Feature prefix stripping",
			verb:     "cpf",
			args:     []string{"Feature:", "[auth]:", "add", "jwt"},
			expected: "Feature: [auth]: add jwt",
		},
		{
			name:     "duplicate Bug prefix stripping",
			verb:     "cbf",
			args:     []string{"Bug:", "[db]:", "fix", "lock"},
			expected: "Bug: [db]: fix lock",
		},
		{
			name:     "legacy hyphen syntax conversion",
			verb:     "cpf",
			args:     []string{"auth", "-", "add", "jwt", "validation"},
			expected: "Feature: [auth]: add jwt validation",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatted := FormatCommitMessage(tc.verb, tc.args)
			if formatted != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, formatted)
			}
		})
	}
}
