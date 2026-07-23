package message

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func fixedPick(_ int) int { return 0 }

func TestStripRulesRemovesAndCollapses(t *testing.T) {
	in := "feat: x\n\nSigned-off-by: a\n\n\nbody\n"
	rules := []profile.MessageRule{{Kind: constants.CommitInMessageRuleKindStartsWith, Value: "Signed-off-by:"}}
	got := stripRules(in, rules)
	if strings.Contains(got, "Signed-off-by") {
		t.Fatalf("rule not applied: %q", got)
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
	if !r.IsEmpty {
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
