package cmdagy

import (
	"reflect"
	"testing"
)

func TestStripAgyPrefix(t *testing.T) {
	cases := []struct {
		input []string
		want  []string
	}{
		{[]string{}, []string{}},
		{[]string{"agy", "fix-pipeline"}, []string{"fix-pipeline"}},
		{[]string{"ag", "lp"}, []string{"lp"}},
		{[]string{"antigravity", "rerun"}, []string{"rerun"}},
		{[]string{"status"}, []string{"status"}},
	}

	for _, tc := range cases {
		got := stripAgyPrefix(tc.input)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("stripAgyPrefix(%v) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestIsCompoundAgyFix(t *testing.T) {
	trueCases := [][]string{
		{"errors", "fix"},
		{"error", "aef"},
		{"err", "fix"},
		{"fix", "errors"},
		{"fix", "error"},
		{"fix", "pipeline"},
		{"fix", "agy"},
	}

	for _, tc := range trueCases {
		if !isCompoundAgyFix(tc) {
			t.Errorf("expected isCompoundAgyFix(%v) to be true", tc)
		}
	}

	falseCases := [][]string{
		{"errors"},
		{"fix"},
		{"status", "all"},
		{"clean", "cache"},
	}

	for _, tc := range falseCases {
		if isCompoundAgyFix(tc) {
			t.Errorf("expected isCompoundAgyFix(%v) to be false", tc)
		}
	}
}

func TestRewriteCompoundAgyFix(t *testing.T) {
	input := []string{"errors", "fix", "--repo", "my-repo"}
	got := rewriteCompoundAgyFix(input)
	want := []string{"fix-pipeline", "--repo", "my-repo"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("rewriteCompoundAgyFix(%v) = %v, want %v", input, got, want)
	}
}

func TestNormalizeAgySubcommand_Aliases(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"fix", "fix-pipeline"},
		{"fp", "fix-pipeline"},
		{"aef", "fix-pipeline"},
		{"fixpipeline", "fix-pipeline"},
		{"pipeline-fix", "fix-pipeline"},
		{"list-prompts", "list-prompts"},
		{"listprompts", "list-prompts"},
		{"lp", "list-prompts"},
		{"list-prompt", "list-prompts"},
		{"rerun", "rerun"},
		{"replay", "rerun"},
		{"rr", "rerun"},
		{"pinned", "pin-projects"},
		{"pins", "pin-projects"},
		{"clean-cache", "clean-cache"},
		{"cc", "clean-cache"},
		{"cleancache", "clean-cache"},
		{"unknown-sub", "unknown-sub"},
	}

	for _, tc := range cases {
		got := normalizeAgySubcommand(tc.input)
		if got != tc.want {
			t.Errorf("normalizeAgySubcommand(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestNormalizeAgyArgs_CompoundAndSingle(t *testing.T) {
	if len(normalizeAgyArgs([]string{})) != 0 {
		t.Error("expected empty slice for empty args")
	}

	compound := normalizeAgyArgs([]string{"errors", "fix", "--dry-run"})
	if compound[0] != "fix-pipeline" || compound[1] != "--dry-run" {
		t.Errorf("unexpected normalized compound args: %v", compound)
	}

	single := normalizeAgyArgs([]string{"lp"})
	if single[0] != "list-prompts" {
		t.Errorf("unexpected normalized single arg: %v", single)
	}
}

func TestIsAgyOpenPathArg_Extended(t *testing.T) {
	if !isAgyOpenPathArg(".") || !isAgyOpenPathArg("..") {
		t.Error("expected dot paths to be recognized as open paths")
	}

	if isAgyOpenPathArg("not-a-path-xyz-123456") {
		t.Error("expected non-path argument to be false")
	}
}
