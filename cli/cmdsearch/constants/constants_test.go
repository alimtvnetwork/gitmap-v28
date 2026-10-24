package constants

import "testing"

func TestCommandConstants(t *testing.T) {
	commands := []struct {
		name     string
		constant string
		expected string
	}{
		{"CmdSearch", CmdSearch, "search"},
		{"CmdRepoSearch", CmdRepoSearch, "repo-search"},
		{"CmdRepoRegex", CmdRepoRegex, "repo-regex"},
		{"CmdRepoSearchJSON", CmdRepoSearchJSON, "repo-search-json"},
		{"CmdRepoSearchRegexJSON", CmdRepoSearchRegexJSON, "repo-search-regex-json"},
		{"CmdReplaceRegex", CmdReplaceRegex, "replace-regex"},
		{"CmdSearchReplaceAll", CmdSearchReplaceAll, "search-replace-all"},
		{"CmdAUMSearch", CmdAUMSearch, "aum search"},
	}

	for _, c := range commands {
		if c.constant != c.expected {
			t.Errorf("expected %s == %q, got %q", c.name, c.expected, c.constant)
		}
	}
}

func TestFlagConstants(t *testing.T) {
	flags := []string{
		FlagQuery,
		FlagQueryShort,
		FlagPattern,
		FlagDir,
		FlagDirShort,
		FlagRoot,
		FlagPath,
		FlagPathShort,
		FlagExt,
		FlagExtShort,
		FlagRegex,
		FlagRegexShort,
		FlagIgnoreCase,
		FlagIgnoreCaseShort,
		FlagLimit,
		FlagLimitShort,
		FlagAI,
		FlagJSON,
		FlagJSONShort,
		FlagHelp,
		FlagHelpShort,
	}

	for _, f := range flags {
		if len(f) == 0 {
			t.Errorf("expected flag constant to be non-empty")
		}
	}
}

func TestFlagArgConstants(t *testing.T) {
	args := []string{
		FlagArgQuery,
		FlagArgQueryShort,
		FlagArgQueryPrefix,
		FlagArgPattern,
		FlagArgPatternPrefix,
		FlagArgDir,
		FlagArgDirShort,
		FlagArgRoot,
		FlagArgPath,
		FlagArgPathShort,
		FlagArgExt,
		FlagArgExtShort,
		FlagArgRegex,
		FlagArgRegexShort,
		FlagArgIgnoreCase,
		FlagArgIgnoreCaseShort,
		FlagArgLimit,
		FlagArgLimitShort,
		FlagArgLimitPrefix,
		FlagArgAI,
		FlagArgAIShort,
		FlagArgJSON,
		FlagArgJSONShort,
		FlagArgHelp,
		FlagArgHelpShort,
		FlagPrefixDash,
	}

	for _, a := range args {
		if len(a) == 0 {
			t.Errorf("expected flag argument constant to be non-empty")
		}
	}
}

func TestDirectoryAndPathConstants(t *testing.T) {
	if DefaultLocalSearchLimit != 20 {
		t.Errorf("expected DefaultLocalSearchLimit == 20, got %d", DefaultLocalSearchLimit)
	}

	if DirNodeModules != "node_modules" {
		t.Errorf("expected DirNodeModules == 'node_modules', got %q", DirNodeModules)
	}

	if DirVendor != "vendor" {
		t.Errorf("expected DirVendor == 'vendor', got %q", DirVendor)
	}

	if DirBin != "bin" {
		t.Errorf("expected DirBin == 'bin', got %q", DirBin)
	}

	if DirPrefixDot != "." {
		t.Errorf("expected DirPrefixDot == '.', got %q", DirPrefixDot)
	}

	if PathRepoSearchOutput != ".gitmap/output/repo_search" {
		t.Errorf("expected PathRepoSearchOutput == '.gitmap/output/repo_search', got %q", PathRepoSearchOutput)
	}

	if EmptyJSONArray != "[]" {
		t.Errorf("expected EmptyJSONArray == '[]', got %q", EmptyJSONArray)
	}
}

func TestUIStringConstants(t *testing.T) {
	if HelpTitleGlobalSearch != "Global Indexed Search (gitmap search)" {
		t.Errorf("expected HelpTitleGlobalSearch to match, got %q", HelpTitleGlobalSearch)
	}

	if SectionSearchOperations != "Search Operations" {
		t.Errorf("expected SectionSearchOperations to match, got %q", SectionSearchOperations)
	}

	if SectionSQLiteCaching != "High-Speed SQLite Caching" {
		t.Errorf("expected SectionSQLiteCaching to match, got %q", SectionSQLiteCaching)
	}

	if len(TipHotQueryCache) == 0 || len(TipHistoryInspect) == 0 {
		t.Errorf("expected tips to be non-empty")
	}

	if len(MsgSearchCleanSuccess) == 0 || len(MsgRepoSearchClear) == 0 {
		t.Errorf("expected messages to be non-empty")
	}

	if len(ErrRepoNotFound) == 0 || len(ErrRepoSearchFatal) == 0 {
		t.Errorf("expected error constants to be non-empty")
	}

	if len(ExampleSearch1) == 0 || len(ExampleReplaceRegex) == 0 {
		t.Errorf("expected examples to be non-empty")
	}

	if len(HelpEntryDH2DCache) == 0 || len(HelpEntryRAMHotCache) == 0 {
		t.Errorf("expected help entry labels to be non-empty")
	}
}
