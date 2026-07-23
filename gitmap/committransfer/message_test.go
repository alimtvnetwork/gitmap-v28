package committransfer

import (
	"strings"
	"testing"
	"time"
)

func TestCleanMessageDropPattern(t *testing.T) {
	policy := MessagePolicy{DropPatterns: DefaultDropPatterns}
	res := CleanMessage("Merge branch 'main' into dev", "", policy, "abc123", time.Now())
	if res.Final != "" {
		t.Fatalf("expected drop, got %q", res.Final)
	}
	if !strings.HasPrefix(res.Skipped, "drop-pattern") {
		t.Fatalf("expected drop-pattern reason, got %q", res.Skipped)
	}
}

func TestCleanMessageStripThenConventional(t *testing.T) {
	policy := MessagePolicy{
		StripPatterns: []string{`^\[WIP\]\s*`, `\s*\(#\d+\)$`},
		Conventional:  true,
	}
	res := CleanMessage("[WIP] Add login form (#42)", "", policy, "abc123", time.Now())
	if res.Skipped != "" {
		t.Fatalf("unexpected skip: %q", res.Skipped)
	}
	if !strings.HasPrefix(res.Final, "feat: ") {
		t.Fatalf("expected feat: prefix, got %q", res.Final)
	}
}

func TestCleanMessagePreservesConventional(t *testing.T) {
	policy := MessagePolicy{Conventional: true}
	res := CleanMessage("fix(auth): handle expired tokens", "", policy, "abc123", time.Now())
	if res.Final != "fix(auth): handle expired tokens" {
		t.Fatalf("conventional subject mutated: %q", res.Final)
	}
}

func TestCleanMessageProvenanceFooter(t *testing.T) {
	when := time.Date(2026, 4, 23, 12, 0, 0, 0, time.UTC)
	policy := MessagePolicy{
		Provenance:        true,
		SourceDisplayName: "repo-A",
		CommandName:       "commit-right",
	}
	res := CleanMessage("feat: add OAuth", "body line", policy, "a3f2c1d", when)
	if !strings.Contains(res.Final, "gitmap-replay: from repo-A a3f2c1d") {
		t.Fatalf("missing provenance footer: %q", res.Final)
	}
	if !strings.Contains(res.Final, "gitmap-replay-cmd: commit-right") {
		t.Fatalf("missing cmd footer: %q", res.Final)
	}
}

func TestAlreadyReplayedDetection(t *testing.T) {
	log := "feat: foo\n\ngitmap-replay: from repo-A a3f2c1d\n---commit-sep---\n"
	if !AlreadyReplayed(log, "repo-A", "a3f2c1d") {
		t.Fatal("expected match")
	}
	if AlreadyReplayed(log, "repo-A", "deadbee") {
		t.Fatal("false positive on different sha")
	}
}

func TestBuildReplayedSet(t *testing.T) {
	log := "feat: foo\n\ngitmap-replay: from repo-A a3f2c1d\ngitmap-replay: from repo-B b4e5g6h\n---commit-sep---\n"
	set := BuildReplayedSet(log)
	if len(set) != 2 {
		t.Fatalf("expected 2 entries, got %d: %+v", len(set), set)
	}
	if _, ok := set["repo-A a3f2c1d"]; !ok {
		t.Fatal("missing repo-A a3f2c1d")
	}
	if _, ok := set["repo-B b4e5g6h"]; !ok {
		t.Fatal("missing repo-B b4e5g6h")
	}
}

func TestSetHasReplayed(t *testing.T) {
	set := map[string]struct{}{
		"repo-A a3f2c1d": {},
	}
	if !SetHasReplayed(set, "repo-A", "a3f2c1d") {
		t.Fatal("expected match")
	}
	if SetHasReplayed(set, "repo-A", "deadbee") {
		t.Fatal("false positive on different sha")
	}
	if SetHasReplayed(set, "repo-B", "a3f2c1d") {
		t.Fatal("false positive on different repo")
	}
}

func TestCleanMessageEmptyAfterStrip(t *testing.T) {
	policy := MessagePolicy{StripPatterns: []string{`.*`}}
	res := CleanMessage("anything", "", policy, "abc", time.Now())
	if res.Skipped != "cleaned-empty" {
		t.Fatalf("expected cleaned-empty, got %q / final %q", res.Skipped, res.Final)
	}
}
