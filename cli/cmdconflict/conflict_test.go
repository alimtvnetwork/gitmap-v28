package cmdconflict

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict/constants"
)

func TestStatusCodeParsing(t *testing.T) {
	testUnmergedCodes(t)
	testNonUnmergedCodes(t)
}

func testUnmergedCodes(t *testing.T) {
	unmerged := []string{"UU", "AA", "UD", "DU", "DD", "AU", "UA"}
	for _, code := range unmerged {
		if !isUnmergedStatusCode(code) {
			t.Errorf("expected isUnmergedStatusCode to be true for unmerged code %q", code)
		}
	}
}

func testNonUnmergedCodes(t *testing.T) {
	normal := []string{"M ", " M", "A ", "??", " D", "D ", "", "XX"}
	for _, code := range normal {
		if isUnmergedStatusCode(code) {
			t.Errorf("expected isUnmergedStatusCode to be false for code %q", code)
		}
	}
}

func TestHasString(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}
	if !hasString(items, "banana") {
		t.Error("expected hasString to find banana")
	}
	if hasString(items, "dragonfruit") {
		t.Error("expected hasString to not find dragonfruit")
	}
}

func TestFormatMergeEditorGuide(t *testing.T) {
	guide := FormatMergeEditorGuide(".", []string{"file.txt"})
	if !strings.Contains(guide, constants.GuideTitle) {
		t.Error("expected guide to contain title")
	}
	if !strings.Contains(guide, constants.GuideStep1) {
		t.Error("expected guide to contain step 1")
	}
	if !strings.Contains(guide, constants.GuideStep6) {
		t.Error("expected guide to contain step 6")
	}
}

func TestParseResolveFlags(t *testing.T) {
	testResolveFlagsCanonical(t)
	testResolveFlagsControlAliases(t)
	testResolveFlagsRerereAlias(t)
	testResolveFlagsTargetOnly(t)
}

func testResolveFlagsCanonical(t *testing.T) {
	target, strategy := parseResolveFlags([]string{"my-repo", "--theirs"})
	if target != "my-repo" || strategy != "theirs" {
		t.Errorf("unexpected parse result: target=%q, strategy=%q", target, strategy)
	}

	target, strategy = parseResolveFlags([]string{"--ours", "repo-x"})
	if target != "repo-x" || strategy != "ours" {
		t.Errorf("unexpected parse result: target=%q, strategy=%q", target, strategy)
	}
}

func testResolveFlagsControlAliases(t *testing.T) {
	_, strategy := parseResolveFlags([]string{"-i"})
	if strategy != "ide" {
		t.Errorf("expected ide, got %q", strategy)
	}

	_, strategy = parseResolveFlags([]string{"--abort"})
	if strategy != "abort" {
		t.Errorf("expected abort, got %q", strategy)
	}
}

func testResolveFlagsRerereAlias(t *testing.T) {
	_, strategy := parseResolveFlags([]string{"rerere"})
	if strategy != "rerere" {
		t.Errorf("expected rerere, got %q", strategy)
	}
}

func testResolveFlagsTargetOnly(t *testing.T) {
	target, strategy := parseResolveFlags([]string{"just-target"})
	if target != "just-target" || strategy != "" {
		t.Errorf("unexpected parse result: target=%q, strategy=%q", target, strategy)
	}
}

func TestResolveConflictActionUnknown(t *testing.T) {
	handled, err := ResolveConflictAction(".", "completely_bogus_action")
	if handled || err != nil {
		t.Errorf("expected false, nil for bogus action, got %v, %v", handled, err)
	}
}

func TestPrintQuickActionsAndTree(t *testing.T) {
	printQuickActions("sample-repo")
	printConflictTree("sample-repo", "/path/to/sample", []string{"a.go", "b.go"})
}

func TestDetectInstalledIDE(t *testing.T) {
	bin, name := DetectInstalledIDE()
	_ = bin
	_ = name
}
