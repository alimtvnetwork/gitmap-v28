package suggestion

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func findSuggestionByCommand(items []Suggestion, cmd string) (Suggestion, bool) {
	for _, s := range items {
		if s.Command == cmd {
			return s, true
		}
	}
	return Suggestion{}, false
}

func TestTier1CanonicalAliases(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("s")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for alias 's'")
	}
	item, found := findSuggestionByCommand(group.Suggestions, "gitmap status")
	if !found {
		t.Fatalf("expected 'gitmap status' suggestion")
	}
	if item.Confidence != 1.0 {
		t.Errorf("expected 1.0 confidence, got %f", item.Confidence)
	}
}

func TestTier1AliasPE(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("pe")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for alias 'pe'")
	}
	_, found := findSuggestionByCommand(group.Suggestions, "gitmap pull-error")
	if !found {
		t.Fatalf("expected 'gitmap pull-error' suggestion for 'pe'")
	}
}

func TestTier2LevenshteinScan(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("scann")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for typo 'scann'")
	}
	item, found := findSuggestionByCommand(group.Suggestions, "gitmap scan")
	if !found {
		t.Fatalf("expected 'gitmap scan' suggestion for 'scann'")
	}
	if item.Confidence < 0.70 {
		t.Errorf("expected high confidence for 1-edit distance, got %f", item.Confidence)
	}
}

func TestTier2LevenshteinStatus(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("statuss")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for typo 'statuss'")
	}
	_, found := findSuggestionByCommand(group.Suggestions, "gitmap status")
	if !found {
		t.Fatalf("expected 'gitmap status' suggestion for 'statuss'")
	}
}

func TestTier3Prefix(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("rele")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for prefix 'rele'")
	}
	_, found := findSuggestionByCommand(group.Suggestions, "gitmap release")
	if !found {
		t.Fatalf("expected 'gitmap release' suggestion for 'rele'")
	}
}

func TestTier4SynonymDocker(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("docker")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for synonym 'docker'")
	}
	item, found := findSuggestionByCommand(group.Suggestions, "gitmap cluster")
	if !found {
		t.Fatalf("expected 'gitmap cluster' suggestion for 'docker'")
	}
	if item.Category != CategoryIntent {
		t.Errorf("expected CategoryIntent, got %s", item.Category)
	}
}

func TestTier4SynonymGrep(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("grep")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for synonym 'grep'")
	}
	_, found := findSuggestionByCommand(group.Suggestions, "gitmap aum search")
	if !found {
		t.Fatalf("expected 'gitmap aum search' suggestion for 'grep'")
	}
}

func TestTier4SynonymRm(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("rm")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for synonym 'rm'")
	}
	_, found := findSuggestionByCommand(group.Suggestions, "gitmap clean")
	if !found {
		t.Fatalf("expected 'gitmap clean' suggestion for 'rm'")
	}
}

func TestResolveFlagHelp(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveFlag("scan", "--hlp")
	if !group.HasSuggestions() {
		t.Fatalf("expected suggestions for flag '--hlp'")
	}
	_, found := findSuggestionByCommand(group.Suggestions, "--help")
	if !found {
		t.Fatalf("expected '--help' suggestion for '--hlp'")
	}
}

func TestResolveRemediationDirty(t *testing.T) {
	engine := NewEngine()
	err := errors.New("cannot switch branch: dirty working tree has modified files")
	group := engine.ResolveRemediation(err, nil)
	if !group.HasSuggestions() {
		t.Fatalf("expected remediation suggestions for dirty tree")
	}
	_, found := findSuggestionByCommand(group.Suggestions, "gitmap status")
	if !found {
		t.Fatalf("expected 'gitmap status' remediation for dirty tree")
	}
}

func TestResolveRemediationConflict(t *testing.T) {
	engine := NewEngine()
	err := errors.New("automatic merge failed: resolve conflicts before committing")
	group := engine.ResolveRemediation(err, nil)
	if !group.HasSuggestions() {
		t.Fatalf("expected remediation suggestions for conflict")
	}
	_, found := findSuggestionByCommand(group.Suggestions, "gitmap reconcile")
	if !found {
		t.Fatalf("expected 'gitmap reconcile' remediation for conflict")
	}
}

func TestRenderBox(t *testing.T) {
	var buf bytes.Buffer
	group := SuggestionGroup{
		Title:       "Did you mean?",
		Reason:      "Command 'scann' is not recognized",
		Suggestions: []Suggestion{{Command: "gitmap scan", Description: "Fast repository scanner", Confidence: 0.95}},
	}
	if err := RenderBox(&buf, group); err != nil {
		t.Fatalf("RenderBox returned unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "gitmap scan") {
		t.Errorf("expected output to contain 'gitmap scan'")
	}
}

func TestRenderCompact(t *testing.T) {
	var buf bytes.Buffer
	group := SuggestionGroup{
		Suggestions: []Suggestion{{Command: "gitmap scan", Confidence: 0.95}},
	}
	if err := RenderCompact(&buf, group); err != nil {
		t.Fatalf("RenderCompact returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Did you mean: gitmap scan (95%)") {
		t.Errorf("unexpected compact output: %s", out)
	}
}

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	group := SuggestionGroup{
		Title:       "Test",
		Suggestions: []Suggestion{{Command: "gitmap scan", Confidence: 0.95}},
	}
	if err := RenderJSON(&buf, group); err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}
	var decoded SuggestionGroup
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if len(decoded.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion in JSON, got %d", len(decoded.Suggestions))
	}
}

func TestEmptyTokenHandling(t *testing.T) {
	engine := NewEngine()
	group := engine.ResolveCommand("")
	if group.HasSuggestions() {
		t.Errorf("expected no suggestions for empty token")
	}
}
