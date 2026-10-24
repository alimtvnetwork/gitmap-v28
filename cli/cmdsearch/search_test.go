package cmdsearch

import (
	"os"
	"path/filepath"
	"testing"

	searchconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch/constants"
)

func TestParseLimit(t *testing.T) {
	limit, clean := parseLimit([]string{"--limit", "42", "keyword"})
	if limit != 42 || len(clean) != 1 || clean[0] != "keyword" {
		t.Errorf("unexpected parseLimit result: %d, %v", limit, clean)
	}

	limitShort, cleanShort := parseLimit([]string{"-l", "15", "test"})
	if limitShort != 15 || len(cleanShort) != 1 {
		t.Errorf("unexpected short limit result: %d", limitShort)
	}

	limitPrefix, cleanPrefix := parseLimit([]string{"--limit=99"})
	if limitPrefix != 99 || len(cleanPrefix) != 0 {
		t.Errorf("unexpected prefix limit result: %d", limitPrefix)
	}
}

func TestParseSearchAiFlag(t *testing.T) {
	isAi, remaining := parseSearchAiFlag([]string{"--ai", "my-query"})
	if !isAi || len(remaining) != 1 || remaining[0] != "my-query" {
		t.Errorf("expected isAi=true, got %v, %v", isAi, remaining)
	}

	isAiShort, remainingShort := parseSearchAiFlag([]string{"-ai", "query2"})
	if !isAiShort || len(remainingShort) != 1 {
		t.Errorf("expected isAiShort=true, got %v", isAiShort)
	}

	isAiNone, remainingNone := parseSearchAiFlag([]string{"query3"})
	if isAiNone || len(remainingNone) != 1 {
		t.Errorf("expected isAiNone=false, got %v", isAiNone)
	}
}

func TestParseSearchInvocation(t *testing.T) {
	q, limit, isAi, hasQuery := parseSearchInvocation([]string{"--ai", "--limit", "10", "needle"})
	if !hasQuery || q != "needle" || limit != 10 || !isAi {
		t.Errorf("unexpected invocation: q=%s limit=%d isAi=%v hasQuery=%v", q, limit, isAi, hasQuery)
	}

	qFlag, _, _, hasQueryFlag := parseSearchInvocation([]string{"--query=direct"})
	if !hasQueryFlag || qFlag != "direct" {
		t.Errorf("unexpected query flag parsing: %s", qFlag)
	}

	_, _, _, hasEmpty := parseSearchInvocation([]string{})
	if hasEmpty {
		t.Errorf("expected hasEmpty=false for empty args")
	}
}

func TestIsSearchQueryAndDirArgs(t *testing.T) {
	if !isSearchQueryArg("--query") || !isSearchQueryArg("-q") || !isSearchQueryArg("--pattern") {
		t.Errorf("expected valid query args to return true")
	}

	if isSearchQueryArg("--limit") {
		t.Errorf("expected non-query arg to return false")
	}

	if !isSearchDirArg("--root") || !isSearchDirArg("--dir") || !isSearchDirArg("-d") {
		t.Errorf("expected valid dir args to return true")
	}

	if isSearchDirArg("--query") {
		t.Errorf("expected non-dir arg to return false")
	}
}

func TestIsIgnoredSearchDir(t *testing.T) {
	if !isIgnoredSearchDir(".git") || !isIgnoredSearchDir("node_modules") || !isIgnoredSearchDir("vendor") {
		t.Errorf("expected ignored directories to return true")
	}

	if isIgnoredSearchDir("src") || isIgnoredSearchDir("cmd") {
		t.Errorf("expected non-ignored directory to return false")
	}
}

func TestFormatRelOrAbsPath(t *testing.T) {
	root := string(filepath.Separator) + "my" + string(filepath.Separator) + "project"
	nested := filepath.Join(root, "pkg", "search.go")
	rel := formatRelOrAbsPath(root, nested)
	if rel != filepath.Join("pkg", "search.go") {
		t.Errorf("expected relative path, got %q", rel)
	}
}

func TestHelpMenuRender(t *testing.T) {
	menu := buildSearchHelpMenu()
	if menu.Title != searchconstants.HelpTitleGlobalSearch {
		t.Errorf("unexpected menu title: %q", menu.Title)
	}

	if len(menu.UsageLines) != 4 || len(menu.Sections) != 2 {
		t.Errorf("unexpected menu structure: %v", menu)
	}
}

func TestRunSearch_NoArgs(t *testing.T) {
	err := RunSearch([]string{})
	if err != nil {
		t.Errorf("expected nil error on empty search, got %v", err)
	}
}

func TestRunReplaceRegex_NoArgs(t *testing.T) {
	err := RunReplaceRegex([]string{})
	if err != nil {
		t.Errorf("expected nil error on empty replace-regex, got %v", err)
	}
}

func TestRunRepoSearch_NoArgs(t *testing.T) {
	err := RunRepoSearch([]string{})
	if err != nil {
		t.Errorf("expected nil error on empty repo-search, got %v", err)
	}
}

func TestRunRepoRegex_NoArgs(t *testing.T) {
	err := RunRepoRegex([]string{})
	if err != nil {
		t.Errorf("expected nil error on empty repo-regex, got %v", err)
	}
}

func TestRunRepoSearchJson_NoArgs(t *testing.T) {
	err := RunRepoSearchJson([]string{})
	if err != nil {
		t.Errorf("expected nil error on empty repo-search-json, got %v", err)
	}
}

func TestRunRepoSearchRegexJson_NoArgs(t *testing.T) {
	err := RunRepoSearchRegexJson([]string{})
	if err != nil {
		t.Errorf("expected nil error on empty repo-search-regex-json, got %v", err)
	}
}

func TestRunSearchReplaceAll_Clear(t *testing.T) {
	tmpDir := filepath.Join(".gitmap", "output", "repo_search")
	_ = os.MkdirAll(tmpDir, 0o755)

	err := RunSearchReplaceAll([]string{"clear"})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}
