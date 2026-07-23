package cmd

import (
	"strings"
	"testing"
)

// TestParseFindNextFlags_HappyPaths covers every accepted form so a
// future refactor that tightens validation can't regress legitimate
// invocations.
func TestParseFindNextFlags_HappyPaths(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		wantID  int64
		wantJSN bool
	}{
		{"empty", []string{}, 0, false},
		{"json only", []string{"--json"}, 0, true},
		{"scan-folder space form", []string{"--scan-folder", "42"}, 42, false},
		{"scan-folder equals form", []string{"--scan-folder=42"}, 42, false},
		{"both flags", []string{"--scan-folder", "7", "--json"}, 7, true},
		{"flags reversed", []string{"--json", "--scan-folder=9"}, 9, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, jsn, err := parseFindNextFlags(tc.args)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if id != tc.wantID || jsn != tc.wantJSN {
				t.Fatalf("got (%d,%v) want (%d,%v)", id, jsn, tc.wantID, tc.wantJSN)
			}
		})
	}
}

// TestParseFindNextFlags_RejectsJSONValue asserts we error out on
// `--json=true` (and any other value) instead of silently accepting
// it. This is the headline behavior change.
func TestParseFindNextFlags_RejectsJSONValue(t *testing.T) {
	for _, arg := range []string{"--json=true", "--json=false", "--json=1"} {
		_, _, err := parseFindNextFlags([]string{arg})
		if err == nil {
			t.Fatalf("expected error for %q, got nil", arg)
		}
		if !strings.Contains(err.Error(), "does not take a value") {
			t.Fatalf("expected boolean-no-value error for %q, got %v", arg, err)
		}
	}
}

// TestParseFindNextFlags_RejectsBadInt asserts that non-integer
// scan-folder values are surfaced (previously silently dropped).
func TestParseFindNextFlags_RejectsBadInt(t *testing.T) {
	cases := [][]string{
		{"--scan-folder", "abc"},
		{"--scan-folder=xyz"},
	}
	for _, args := range cases {
		_, _, err := parseFindNextFlags(args)
		if err == nil {
			t.Fatalf("expected error for %v, got nil", args)
		}
		if !strings.Contains(err.Error(), "expects an integer") {
			t.Fatalf("expected bad-int error for %v, got %v", args, err)
		}
	}
}

// TestParseFindNextFlags_RejectsMissingValue asserts the
// `--scan-folder` flag must be followed by an actual value rather
// than another flag or the end of args.
func TestParseFindNextFlags_RejectsMissingValue(t *testing.T) {
	cases := [][]string{
		{"--scan-folder"},
		{"--scan-folder", "--json"},
	}
	for _, args := range cases {
		_, _, err := parseFindNextFlags(args)
		if err == nil {
			t.Fatalf("expected missing-value error for %v, got nil", args)
		}
		if !strings.Contains(err.Error(), "requires an integer") {
			t.Fatalf("expected missing-value error for %v, got %v", args, err)
		}
	}
}

// TestParseFindNextFlags_UnknownFlagWithSuggestion asserts that
// near-miss typos produce a "did you mean?" hint.
func TestParseFindNextFlags_UnknownFlagWithSuggestion(t *testing.T) {
	cases := []struct {
		arg  string
		want string
	}{
		{"--jsno", "--json"},
		{"--scanfolder", "--scan-folder"},
		{"--scan_folder", "--scan-folder"},
	}
	for _, tc := range cases {
		_, _, err := parseFindNextFlags([]string{tc.arg})
		if err == nil {
			t.Fatalf("expected unknown-flag error for %q, got nil", tc.arg)
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected suggestion %q in %v", tc.want, err)
		}
	}
}

// TestParseFindNextFlags_UnknownFlagWithoutSuggestion asserts that
// totally unrelated flags surface a plain unknown-flag error rather
// than a misleading suggestion.
func TestParseFindNextFlags_UnknownFlagWithoutSuggestion(t *testing.T) {
	_, _, err := parseFindNextFlags([]string{"--quiet"})
	if err == nil {
		t.Fatalf("expected unknown-flag error, got nil")
	}
	if strings.Contains(err.Error(), "did you mean") {
		t.Fatalf("did not expect suggestion in %v", err)
	}
}

// TestParseFindNextFlags_RejectsPositional asserts bare positional
// tokens (almost always a quoting bug) are rejected.
func TestParseFindNextFlags_RejectsPositional(t *testing.T) {
	_, _, err := parseFindNextFlags([]string{"oops"})
	if err == nil {
		t.Fatalf("expected positional error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected positional") {
		t.Fatalf("expected positional error, got %v", err)
	}
}
