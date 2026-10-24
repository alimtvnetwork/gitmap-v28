package cmdcommit

import (
	"strings"
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
			expected: "FEAT(AUTH): support jwt refresh tokens",
			module:   "AUTH",
		},
		{
			name:     "cbf bracket sqlite",
			verb:     "cbf",
			args:     []string{"[sqlite]:", "handle", "connection", "pool", "lock"},
			expected: "FIX(SQLITE): handle connection pool lock",
			module:   "SQLITE",
		},
		{
			name:     "cpb bracket db",
			verb:     "cpb",
			args:     []string{"[db]:", "fix", "deadlock"},
			expected: "FIX(DB): fix deadlock",
			module:   "DB",
		},
		{
			name:     "cfr bracket engine",
			verb:     "cfr",
			args:     []string{"[engine]:", "optimize", "credential", "lookup", "cascade"},
			expected: "REFACTOR(ENGINE): optimize credential lookup cascade",
			module:   "ENGINE",
		},
		{
			name:     "cpc bracket skills",
			verb:     "cpc",
			args:     []string{"[skills]:", "sync", "gitmap", "skill"},
			expected: "CHORE(SKILLS): sync gitmap skill",
			module:   "SKILLS",
		},
		{
			name:     "cpr bracket version",
			verb:     "cpr",
			args:     []string{"[v6.600.0]:", "release", "notes"},
			expected: "RELEASE(V6.600.0): release notes",
			module:   "V6.600.0",
		},
		{
			name:     "module with path slashes",
			verb:     "cpf",
			args:     []string{"[cli/cmdworkgroup]:", "add", "resolver"},
			expected: "FEAT(CLI/CMDWORKGROUP): add resolver",
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
			if !res.HasBracketSyntax {
				t.Fatalf("expected hasBracketSyntax=true")
			}
			if strings.Contains(res.FormattedMessage, "[") || strings.Contains(res.FormattedMessage, "]") {
				t.Fatalf("commit message contains forbidden square brackets: %s", res.FormattedMessage)
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
			name:     "refactor with plain scope",
			verb:     "cfr",
			args:     []string{"refactor(engine):", "optimize", "cascade"},
			expected: "REFACTOR(ENGINE): optimize cascade",
			module:   "ENGINE",
		},
		{
			name:     "chore with plain scope",
			verb:     "cpc",
			args:     []string{"chore(skills):", "cleanup", "unused", "files"},
			expected: "CHORE(SKILLS): cleanup unused files",
			module:   "SKILLS",
		},
		{
			name:     "uppercase conventional with brackets",
			verb:     "cpf",
			args:     []string{"FEAT([AUTH]):", "add", "jwt", "rotation"},
			expected: "FEAT(AUTH): add jwt rotation",
			module:   "AUTH",
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
				t.Fatalf("commit message contains forbidden square brackets: %s", res.FormattedMessage)
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
			expected: "FEAT(AUTH): add jwt",
		},
		{
			name:     "single quotes wrapping",
			verb:     "cbf",
			args:     []string{`'[db]: fix lock'`},
			expected: "FIX(DB): fix lock",
		},
		{
			name:     "cfr double quotes wrapping",
			verb:     "cfr",
			args:     []string{`"[engine]: optimize cascade"`},
			expected: "REFACTOR(ENGINE): optimize cascade",
		},
		{
			name:     "duplicate Feature prefix stripping",
			verb:     "cpf",
			args:     []string{"Feature:", "[auth]:", "add", "jwt"},
			expected: "FEAT(AUTH): add jwt",
		},
		{
			name:     "duplicate Bug prefix stripping",
			verb:     "cbf",
			args:     []string{"Bug:", "[db]:", "fix", "lock"},
			expected: "FIX(DB): fix lock",
		},
		{
			name:     "legacy hyphen syntax conversion",
			verb:     "cpf",
			args:     []string{"auth", "-", "add", "jwt", "validation"},
			expected: "FEAT(AUTH): add jwt validation",
		},
		{
			name:     "square bracket sanitization in description",
			verb:     "cpf",
			args:     []string{"[auth]:", "support", "[rfc7519]", "jwt"},
			expected: "FEAT(AUTH): support (rfc7519) jwt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatted := FormatCommitMessage(tc.verb, tc.args)
			if formatted != tc.expected {
				t.Fatalf("expected '%s', got '%s'", tc.expected, formatted)
			}
			if strings.Contains(formatted, "[") || strings.Contains(formatted, "]") {
				t.Fatalf("commit message contains forbidden square brackets: %s", formatted)
			}
		})
	}
}
