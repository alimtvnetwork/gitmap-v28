package visibility

import (
	"testing"
)

func TestParsePatternValid(t *testing.T) {
	p, err := ParsePattern("foo*bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Raw != "foo*bar" || !p.anchorL || !p.anchorR {
		t.Fatalf("unexpected pattern: %+v", p)
	}
}

func TestParsePatternErrors(t *testing.T) {
	if _, err := ParsePattern(""); err == nil {
		t.Fatal("expected error for empty pattern")
	}
	if _, err := ParsePattern("*"); err == nil {
		t.Fatal("expected error for bare * pattern")
	}
	if _, err := ParsePattern("***"); err == nil {
		t.Fatal("expected error for only wildcards")
	}
}

func TestPatternMatches(t *testing.T) {
	p, _ := ParsePattern("gitmap-*")
	if !p.Matches("gitmap-v28") {
		t.Fatal("expected gitmap-v28 to match gitmap-*")
	}
	if p.Matches("other-repo") {
		t.Fatal("expected other-repo not to match gitmap-*")
	}
	if p.Matches("") {
		t.Fatal("expected empty name not to match")
	}
}

func TestParsePatternList(t *testing.T) {
	pats, err := ParsePatternList("repo1, repo2*, repo1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pats) != 2 {
		t.Fatalf("expected 2 deduped patterns, got %d", len(pats))
	}
	if _, err := ParsePatternList(""); err == nil {
		t.Fatal("expected error for blank pattern list")
	}
	if _, err := ParsePatternList("a,,b"); err == nil {
		t.Fatal("expected error for empty token in pattern list")
	}
}

func TestMatchOwnerRepos(t *testing.T) {
	pats, _ := ParsePatternList("gitmap-*, *worker*")
	repos := []string{"gitmap-v28", "bg-worker-service", "unknown-tool"}
	matched := MatchOwnerRepos(repos, pats)
	if len(matched) != 2 {
		t.Fatalf("expected 2 matched repos, got %d", len(matched))
	}
	if matched[0].RepoName != "gitmap-v28" || matched[0].MatchedPattern != "gitmap-*" {
		t.Fatalf("unexpected match 0: %+v", matched[0])
	}
	if matched[1].RepoName != "bg-worker-service" || matched[1].MatchedPattern != "*worker*" {
		t.Fatalf("unexpected match 1: %+v", matched[1])
	}
}

func TestVersionHelpers(t *testing.T) {
	base, ver, ok := ParseRepoNameMeta("gitmap-v28")
	if !ok || base != "gitmap" || ver != 28 {
		t.Fatalf("ParseRepoNameMeta failed: %q %d %v", base, ver, ok)
	}
	base2, digits, ok2 := SplitTrailingDigits("macro-51")
	if !ok2 || base2 != "macro" || digits != "51" {
		t.Fatalf("SplitTrailingDigits failed: %q %q %v", base2, digits, ok2)
	}
	bestName, bestVer, ok3 := HighestVersionedMatch([]string{"repo-v1", "repo-v3", "other-v10"}, "repo")
	if !ok3 || bestName != "repo-v3" || bestVer != 3 {
		t.Fatalf("HighestVersionedMatch failed: %q %d %v", bestName, bestVer, ok3)
	}
}
