package commitin

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// AC #1 — every separator/quoting form yields the same input list.
// Locked by the spec acceptance matrix; do not collapse into one case.
func TestParseSeparatorEquivalenceAC1(t *testing.T) {
	want := []string{"a", "b", "c"}
	cases := [][]string{
		{"src", "a", "b", "c"},
		{"src", "a,b,c"},
		{"src", "a, b, c"},
		{"src", `"a, b, c"`},
		{"src", `"a"`, `"b"`, `"c"`},
		{"src", `"a, b"`, "c"},
	}
	for _, argv := range cases {
		got, perr := Parse(argv)
		if perr != nil {
			t.Fatalf("argv %v: unexpected error: %v", argv, perr)
		}
		if !reflect.DeepEqual(got.Inputs, want) {
			t.Errorf("argv %v: got inputs %v, want %v", argv, got.Inputs, want)
		}
	}
}

// AC #4 — `all` and `-N` are recognized as keywords; mixing them with
// explicit inputs is a hard error. Also covers the spec §2.4 grammar
// rule that a KEYWORD must appear ALONE.
func TestParseKeywordsAC4(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		got, perr := Parse([]string{"src", "all"})
		if perr != nil {
			t.Fatalf("unexpected error: %v", perr)
		}
		if got.Keyword != constants.CommitInInputKeywordAll || got.KeywordTail != 0 {
			t.Errorf("got %q tail=%d, want all/0", got.Keyword, got.KeywordTail)
		}
	})
	t.Run("tail-N", func(t *testing.T) {
		got, perr := Parse([]string{"src", "-5"})
		if perr != nil {
			t.Fatalf("unexpected error: %v", perr)
		}
		if got.Keyword != "-5" || got.KeywordTail != 5 {
			t.Errorf("got %q tail=%d, want -5/5", got.Keyword, got.KeywordTail)
		}
	})
	t.Run("-0 rejected", func(t *testing.T) {
		_, perr := Parse([]string{"src", "-0"})
		if perr == nil || perr.ExitCode != constants.CommitInExitBadArgs {
			t.Errorf("want BadArgs for -0, got %v", perr)
		}
	})
	t.Run("-abc rejected", func(t *testing.T) {
		_, perr := Parse([]string{"src", "-abc"})
		// Leading-`-` non-numeric reaches the flag parser (unknown flag).
		// Either way it must surface as BadArgs.
		if perr == nil || perr.ExitCode != constants.CommitInExitBadArgs {
			t.Errorf("want BadArgs for -abc, got %v", perr)
		}
	})
}

// Missing positionals must always fail with BadArgs and a helpful msg.
func TestParseMissingPositionals(t *testing.T) {
	cases := [][]string{
		{},
		{"only-source"},
	}
	for _, argv := range cases {
		_, perr := Parse(argv)
		if perr == nil || perr.ExitCode != constants.CommitInExitBadArgs {
			t.Errorf("argv %v: want BadArgs, got %v", argv, perr)
		}
	}
}

// Author flags must be supplied as a pair (both or neither) per §2.5.
func TestParseAuthorPair(t *testing.T) {
	cases := []struct {
		name string
		argv []string
		ok   bool
	}{
		{"both", []string{"--author-name", "Jane", "--author-email", "j@x.io", "src", "a"}, true},
		{"neither", []string{"src", "a"}, true},
		{"name only", []string{"--author-name", "Jane", "src", "a"}, false},
		{"email only", []string{"--author-email", "j@x.io", "src", "a"}, false},
	}
	for _, tc := range cases {
		_, perr := Parse(tc.argv)
		if (perr == nil) != tc.ok {
			t.Errorf("%s: want ok=%v, got err=%v", tc.name, tc.ok, perr)
		}
	}
}

// Validators reject anything outside the spec's enum sets.
func TestParseEnumValidators(t *testing.T) {
	cases := []struct {
		name    string
		argv    []string
		wantErr bool
	}{
		{"conflict ForceMerge ok", []string{"--conflict", "ForceMerge", "src", "a"}, false},
		{"conflict Prompt ok", []string{"--conflict", "Prompt", "src", "a"}, false},
		{"conflict bogus", []string{"--conflict", "MergeAll", "src", "a"}, true},
		{"function-intel on", []string{"--function-intel", "on", "src", "a"}, false},
		{"function-intel off", []string{"--function-intel", "off", "src", "a"}, false},
		{"function-intel bogus", []string{"--function-intel", "maybe", "src", "a"}, true},
		{"languages Go,Rust ok", []string{"--languages", "Go,Rust", "src", "a"}, false},
		{"languages Cobol bad", []string{"--languages", "Cobol", "src", "a"}, true},
	}
	for _, tc := range cases {
		_, perr := Parse(tc.argv)
		if (perr != nil) != tc.wantErr {
			t.Errorf("%s: wantErr=%v, got %v", tc.name, tc.wantErr, perr)
		}
	}
}

// --message-exclude requires `Kind:Value` with Kind in the enum set.
func TestParseMessageRules(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got, perr := Parse([]string{
			"--message-exclude", "StartsWith:Signed-off-by:,Contains:[skip ci]",
			"src", "a",
		})
		if perr != nil {
			t.Fatalf("unexpected error: %v", perr)
		}
		want := []MessageRuleArg{
			{Kind: "StartsWith", Value: "Signed-off-by:"},
			{Kind: "Contains", Value: "[skip ci]"},
		}
		if !reflect.DeepEqual(got.MessageRules, want) {
			t.Errorf("got %+v, want %+v", got.MessageRules, want)
		}
	})
	t.Run("bad kind", func(t *testing.T) {
		_, perr := Parse([]string{"--message-exclude", "Equals:foo", "src", "a"})
		if perr == nil || perr.ExitCode != constants.CommitInExitBadArgs {
			t.Errorf("want BadArgs, got %v", perr)
		}
	})
	t.Run("missing colon", func(t *testing.T) {
		_, perr := Parse([]string{"--message-exclude", "no-colon", "src", "a"})
		if perr == nil {
			t.Error("want BadArgs for missing colon")
		}
	})
}

// Flags after positionals must still be picked up (reorder logic).
func TestParseFlagsAfterPositionals(t *testing.T) {
	got, perr := Parse([]string{"src", "a", "b", "--dry-run", "--conflict", "Prompt"})
	if perr != nil {
		t.Fatalf("unexpected error: %v", perr)
	}
	if !got.IsDryRun {
		t.Error("--dry-run not picked up after positionals")
	}
	if got.ConflictMode != "Prompt" {
		t.Errorf("got conflict %q, want Prompt", got.ConflictMode)
	}
	if !reflect.DeepEqual(got.Inputs, []string{"a", "b"}) {
		t.Errorf("got inputs %v, want [a b]", got.Inputs)
	}
}

// Bool short alias `-d` toggles --default.
func TestParseDefaultShortAlias(t *testing.T) {
	got, perr := Parse([]string{"-d", "src", "a"})
	if perr != nil {
		t.Fatalf("unexpected error: %v", perr)
	}
	if !got.UseDefaultProfile {
		t.Error("-d did not toggle UseDefaultProfile")
	}
}

// Sanity: the error message for an unknown language enumerates the
// supported set so users can self-correct without reading the spec.
func TestParseLanguageErrorListsKnownSet(t *testing.T) {
	_, perr := Parse([]string{"--languages", "Cobol", "src", "a"})
	if perr == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(perr.Message, "Go") || !strings.Contains(perr.Message, "CSharp") {
		t.Errorf("error %q should list supported languages", perr.Message)
	}
}

// --no-release-branch defaults to OFF (i.e. release branches ARE
// created). When the flag IS passed, IsNoReleaseBranch flips to true
// so downstream replay-mapping (runlog.ResolveReleaseBranchName) can
// suppress branch creation. Spec §08 §2.5 + §09 R4.
func TestParseNoReleaseBranchDefaultsOff(t *testing.T) {
	got, perr := Parse([]string{"src", "a"})
	if perr != nil {
		t.Fatalf("unexpected error: %v", perr)
	}
	if got.IsNoReleaseBranch {
		t.Error("default should be OFF — branches ON")
	}
}

func TestParseNoReleaseBranchFlagFlipsToTrue(t *testing.T) {
	got, perr := Parse([]string{"--no-release-branch", "src", "a"})
	if perr != nil {
		t.Fatalf("unexpected error: %v", perr)
	}
	if !got.IsNoReleaseBranch {
		t.Error("--no-release-branch should set IsNoReleaseBranch=true")
	}
}

// Flag must reorder past positionals like every other bool flag.
func TestParseNoReleaseBranchReordersPastPositionals(t *testing.T) {
	got, perr := Parse([]string{"src", "a", "--no-release-branch"})
	if perr != nil {
		t.Fatalf("unexpected error: %v", perr)
	}
	if !got.IsNoReleaseBranch {
		t.Error("--no-release-branch after positionals not picked up")
	}
}
