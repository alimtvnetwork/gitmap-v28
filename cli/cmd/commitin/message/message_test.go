package message

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func fixedPick(_ int) int { return 0 }

func TestStripRulesRemovesAndCollapses(t *testing.T) {
	in := "feat: x\n\nSigned-off-by: a\nCo-authored-by: alimtvnetwork <253777419+alimtvnetwork@users.noreply.github.com>\n\n\nbody\n"
	rules := []profile.MessageRule{{Kind: constants.CommitInMessageRuleKindStartsWith, Value: "Signed-off-by:"}}
	got := stripRules(in, rules)
	if strings.Contains(got, "Signed-off-by") {
		t.Fatalf("rule not applied: %q", got)
	}
	if strings.Contains(got, "Co-authored-by:") {
		t.Fatalf("noreply co-author not stripped: %q", got)
	}
	if strings.Contains(got, "\n\n\n") {
		t.Fatalf("blank lines not collapsed: %q", got)
	}
}

func TestWeakWordMatching(t *testing.T) {
	weak := []string{"change", "update", "updates"}
	cases := map[string]bool{
		"Update README":      true,
		"Updates: bump deps": true,
		"Refactor parser":    false,
		"":                   false,
		"UPDATE the thing":   true,
	}

	for in, want := range cases {
		if got := matchesWeak(in, weak); got != want {
			t.Errorf("matchesWeak(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestOverrideOnlyFiresWhenWeak(t *testing.T) {
	res := profile.Resolved{
		OverrideMessages: []string{"Refine implementation"},
		OverrideOnlyWeak: true,
		WeakWords:        []string{"update"},
	}

	weakIn := Inputs{OriginalMessage: "Update foo", Resolved: res, PickIndex: fixedPick}
	if got := Build(weakIn).Message; got != "Refine implementation" {
		t.Fatalf("weak override missed: %q", got)
	}

	strongIn := Inputs{OriginalMessage: "Refactor foo", Resolved: res, PickIndex: fixedPick}
	if got := Build(strongIn).Message; got != "Refactor foo" {
		t.Fatalf("strong should not override: %q", got)
	}
}

func TestTitleAffixOnlyTouchesFirstLine(t *testing.T) {
	res := profile.Resolved{TitlePrefix: "[x] ", TitleSuffix: " <-"}
	out := Build(Inputs{OriginalMessage: "title\nbody1\nbody2", Resolved: res}).Message
	lines := strings.Split(out, "\n")
	if lines[0] != "[x] title <-" {
		t.Fatalf("title affix wrong: %q", lines[0])
	}

	if lines[1] != "body1" || lines[2] != "body2" {
		t.Fatalf("body mutated: %v", lines)
	}
}

func TestBodyAffixWraps(t *testing.T) {
	res := profile.Resolved{
		MessagePrefix: []string{"chore:"},
		MessageSuffix: []string{"--end--"},
	}

	out := Build(Inputs{OriginalMessage: "title", Resolved: res, PickIndex: fixedPick}).Message
	if !strings.HasPrefix(out, "chore:\n") || !strings.HasSuffix(out, "\n--end--") {
		t.Fatalf("body affix wrong: %q", out)
	}
}

func TestTitleReplacementWithFiles2NamesAndBlankGap(t *testing.T) {
	res := profile.Resolved{
		TitleReplacements: []profile.TitleReplacementRule{
			{MatchMode: "equals", Match: "Changes", Replacement: "$files.2.names: $seo.title"},
		},
		MessageSuffix: []string{"# Why is deterministic architecture critical?\nBecause it guarantees 99.98% build reliability."},
	}
	outSingle := Build(Inputs{
		OriginalMessage: "Changes\n\nCo-authored-by: user <123+user@users.noreply.github.com>",
		Files:           []string{"cli/cliexit/cliexit.go"},
		Resolved:        res,
		PickIndex:       fixedPick,
	}).Message
	wantSingle := "cliexit.go: Why is deterministic architecture critical?\n\n# Why is deterministic architecture critical?\nBecause it guarantees 99.98% build reliability."
	if outSingle != wantSingle {
		t.Fatalf("single file replacement mismatch:\ngot:  %q\nwant: %q", outSingle, wantSingle)
	}
	assertTwoFileTitleReplacement(t, res)
}

func assertTwoFileTitleReplacement(t *testing.T, res profile.Resolved) {
	outMulti := Build(Inputs{
		OriginalMessage: "Changes",
		Files:           []string{"cli/cliexit/cliexit.go", "cli/cmd/root.go", "cli/cmd/pull.go"},
		Resolved:        res,
		PickIndex:       fixedPick,
	}).Message
	if !strings.HasPrefix(outMulti, "cliexit.go, root.go: Why is deterministic architecture critical?\n\n# Why") {
		t.Fatalf("two-file $files.2.names replacement mismatch: %q", outMulti)
	}
}

func TestTitleReplacementWithoutSuffixTemplate(t *testing.T) {
	res := profile.Resolved{
		TitleReplacements: []profile.TitleReplacementRule{
			{MatchMode: "equals", Match: "Changes", Replacement: "$files.2.names: $seo.title"},
		},
	}
	out := Build(Inputs{
		OriginalMessage: "Changes",
		Files:           []string{"cli/cliexit/cliexit.go", "cli/cmd/root.go"},
		Resolved:        res,
		PickIndex:       fixedPick,
	}).Message
	if out != "cliexit.go, root.go" {
		t.Fatalf("expected clean filenames without sponsor fallback, got %q", out)
	}
}

func TestFunctionIntelAppended(t *testing.T) {
	out := Build(Inputs{
		OriginalMessage: "title",
		FunctionIntel:   "- src/x.go\n  - added: A",
		Resolved:        profile.Resolved{},
	}).Message
	if !strings.Contains(out, "title\n\n- src/x.go") {
		t.Fatalf("intel block not appended: %q", out)
	}
}

func TestEmptyAfterStripFlagged(t *testing.T) {
	res := profile.Resolved{MessageRules: []profile.MessageRule{
		{Kind: constants.CommitInMessageRuleKindStartsWith, Value: "x"},
	}}
	r := Build(Inputs{OriginalMessage: "x line 1\nx line 2", Resolved: res})
	if r.IsDefined() {
		t.Fatalf("expected IsEmpty=true, got %q", r.Message)
	}
}

func TestPipelineOrderStripBeforeOverride(t *testing.T) {
	res := profile.Resolved{
		MessageRules:     []profile.MessageRule{{Kind: constants.CommitInMessageRuleKindStartsWith, Value: "Update"}},
		OverrideMessages: []string{"Refined"},
		OverrideOnlyWeak: true,
		WeakWords:        []string{"update"},
	}

	out := Build(Inputs{OriginalMessage: "Update foo\nReal body", Resolved: res, PickIndex: fixedPick}).Message
	if out == "Refined" {
		t.Fatalf("override fired after strip removed weak title")
	}
}
