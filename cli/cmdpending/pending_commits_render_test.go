package cmdpending

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

func sampleRemediationRecord() RepoPendingCommitRecord {
	return RepoPendingCommitRecord{
		RepoName:           "gitmap",
		RelativePath:       "gitmap",
		IsDirty:            true,
		RemediationOptions: createRemediationOptionsList("gitmap"),
	}
}

func TestRenderTreeRemediationSubtree(t *testing.T) {
	rec := sampleRemediationRecord()
	output := captureStdout(func() {
		renderTreeRemediationHints(rec, "", "")
	})
	verifyRemediationTitleHierarchy(t, output)
	verifyRemediationChildHierarchy(t, output)
}

func verifyRemediationTitleHierarchy(t *testing.T, output string) {
	titles := []string{
		"├── Option 1: Commit & Push WIP",
		"├── Option 2: Stash WIP",
		"└── Option 3: All Repos (Commit & Push)",
	}
	for _, title := range titles {
		if !strings.Contains(output, title) {
			t.Errorf("missing expected title in output: %q", title)
		}
	}
}

func verifyRemediationChildHierarchy(t *testing.T, output string) {
	children := []string{
		"   │   └── gitmap sends cp gitmap \"wip: save changes\"",
		"   │   └── gitmap stash gitmap",
		"       └── gitmap sends cp all \"wip: save changes\"",
	}
	for _, child := range children {
		if !strings.Contains(output, child) {
			t.Errorf("missing expected child in output: %q", child)
		}
	}
}

func TestRenderSingleRemediationSubtree_Individual(t *testing.T) {
	opts := createRemediationOptionsList("test-repo")
	outNonTerminal := captureStdout(func() {
		renderSingleRemediationSubtree(opts[0], false, "", "")
	})
	if !strings.Contains(outNonTerminal, "├── Option 1: Commit & Push WIP") {
		t.Errorf("expected non-terminal title, got %q", outNonTerminal)
	}
	if !strings.Contains(outNonTerminal, "   │   └── gitmap sends cp") {
		t.Errorf("expected non-terminal child, got %q", outNonTerminal)
	}
}

func TestRenderSingleRemediationSubtree_Terminal(t *testing.T) {
	opts := createRemediationOptionsList("test-repo")
	outTerminal := captureStdout(func() {
		renderSingleRemediationSubtree(opts[2], true, "", "")
	})
	if !strings.Contains(outTerminal, "└── Option 3: All Repos (Commit & Push)") {
		t.Errorf("expected terminal title, got %q", outTerminal)
	}
	if !strings.Contains(outTerminal, "       └── gitmap sends cp all") {
		t.Errorf("expected terminal child, got %q", outTerminal)
	}
}

func TestSubtreeLineDimensions(t *testing.T) {
	opts := createRemediationOptionsList("gitmap")
	for idx, opt := range opts {
		isLast := idx == len(opts)-1
		verifySubtreeOptionDimensions(t, opt, isLast)
	}
}

func verifySubtreeOptionDimensions(t *testing.T, opt RemediationOption, isLast bool) {
	prefix := "├──"
	childPrefix := "   │   └── "
	if isLast {
		prefix = "└──"
		childPrefix = "       └── "
	}
	titleLine := formatTreeLine("   " + prefix + " Option " + string(rune('0'+opt.OptionNumber)) + ": " + opt.Label)
	childLine := formatTreeLine(childPrefix + opt.Command)
	assertExactRuneLength(t, titleLine, 78)
	assertExactRuneLength(t, childLine, 78)
}

func assertExactRuneLength(t *testing.T, s string, expected int) {
	actual := len([]rune(s))
	if actual != expected {
		t.Errorf("expected exactly %d runes, got %d for %q", expected, actual, s)
	}
}

func TestBuildRepoRemediationOptions_IncludesUnifiedOption(t *testing.T) {
	dirtyRec := RepoPendingCommitRecord{
		RepoName:     "seo-packages",
		RelativePath: "seo-packages",
		IsDirty:      true,
	}
	opts := buildRepoRemediationOptions(dirtyRec)
	if len(opts) != 3 {
		t.Fatalf("expected 3 remediation options, got %d", len(opts))
	}
	verifyUnifiedOptionProperties(t, opts[2])
	verifyIndividualOptionProperties(t, opts[0], opts[1])
}

func verifyUnifiedOptionProperties(t *testing.T, opt RemediationOption) {
	if opt.OptionNumber != 3 {
		t.Errorf("expected option number 3, got %d", opt.OptionNumber)
	}
	if opt.Label != "All Repos (Commit & Push)" {
		t.Errorf("expected label 'All Repos (Commit & Push)', got %q", opt.Label)
	}
	expectedCmd := `gitmap sends cp all "wip: save changes"`
	if opt.Command != expectedCmd {
		t.Errorf("expected command %q, got %q", expectedCmd, opt.Command)
	}
}

func verifyIndividualOptionProperties(t *testing.T, opt1, opt2 RemediationOption) {
	if opt1.OptionNumber != 1 || opt1.Label != "Commit & Push WIP" {
		t.Errorf("invalid option 1 properties: %+v", opt1)
	}
	if opt2.OptionNumber != 2 || opt2.Label != "Stash WIP" {
		t.Errorf("invalid option 2 properties: %+v", opt2)
	}
}

func TestBuildRepoRemediationOptions_CleanRepoEmpty(t *testing.T) {
	cleanRec := RepoPendingCommitRecord{
		RepoName: "clean-repo",
		IsClean:  true,
		IsDirty:  false,
	}
	opts := buildRepoRemediationOptions(cleanRec)
	if len(opts) != 0 {
		t.Errorf("expected 0 options for clean repo, got %d", len(opts))
	}
}
