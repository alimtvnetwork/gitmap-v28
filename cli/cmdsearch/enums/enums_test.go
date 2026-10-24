package enums

import "testing"

func TestSearchScopeType_Valid(t *testing.T) {
	scopes := AllSearchScopes()
	if len(scopes) != 5 {
		t.Fatalf("expected 5 search scopes, got %d", len(scopes))
	}

	for _, s := range scopes {
		if !s.IsValid() {
			t.Errorf("expected %s to be valid", s)
		}

		if s.String() != string(s) {
			t.Errorf("expected String() == %s, got %s", string(s), s.String())
		}
	}

	if SearchScopeType("unknown").IsValid() {
		t.Errorf("expected invalid search scope to be false")
	}
}

func TestSearchScopeType_Booleans(t *testing.T) {
	if !SearchScopeContent.IsContent() || SearchScopeSymbol.IsContent() {
		t.Errorf("unexpected IsContent behavior")
	}

	if !SearchScopeSymbol.IsSymbol() || SearchScopePath.IsSymbol() {
		t.Errorf("unexpected IsSymbol behavior")
	}

	if !SearchScopePath.IsPath() || SearchScopeCommit.IsPath() {
		t.Errorf("unexpected IsPath behavior")
	}

	if !SearchScopeCommit.IsCommit() || SearchScopeRepo.IsCommit() {
		t.Errorf("unexpected IsCommit behavior")
	}

	if !SearchScopeRepo.IsRepo() || SearchScopeContent.IsRepo() {
		t.Errorf("unexpected IsRepo behavior")
	}
}

func TestSearchScopeType_Parse(t *testing.T) {
	parsed, err := ParseSearchScope("content")
	if err != nil || parsed != SearchScopeContent {
		t.Errorf("expected content, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseSearchScopeType(" REPO ")
	if errAlias != nil || parsedAlias != SearchScopeRepo {
		t.Errorf("expected repo, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseSearchScope("invalid-scope")
	if invalidErr == nil {
		t.Errorf("expected error for invalid scope")
	}

	if len(AllSearchScopeTypes()) != 5 {
		t.Errorf("expected AllSearchScopeTypes count to match")
	}
}

func TestSearchMatchModeType_Valid(t *testing.T) {
	modes := AllSearchMatchModes()
	if len(modes) != 4 {
		t.Fatalf("expected 4 match modes, got %d", len(modes))
	}

	for _, m := range modes {
		if !m.IsValid() {
			t.Errorf("expected %s to be valid", m)
		}

		if m.String() != string(m) {
			t.Errorf("expected String() == %s, got %s", string(m), m.String())
		}
	}

	if SearchMatchModeType("invalid").IsValid() {
		t.Errorf("expected invalid match mode to be false")
	}
}

func TestSearchMatchModeType_Booleans(t *testing.T) {
	if !SearchMatchModeKeyword.IsKeyword() || SearchMatchModeRegex.IsKeyword() {
		t.Errorf("unexpected IsKeyword behavior")
	}

	if !SearchMatchModeRegex.IsRegex() || SearchMatchModeExact.IsRegex() {
		t.Errorf("unexpected IsRegex behavior")
	}

	if !SearchMatchModeExact.IsExact() || SearchMatchModeFuzzy.IsExact() {
		t.Errorf("unexpected IsExact behavior")
	}

	if !SearchMatchModeFuzzy.IsFuzzy() || SearchMatchModeKeyword.IsFuzzy() {
		t.Errorf("unexpected IsFuzzy behavior")
	}
}

func TestSearchMatchModeType_Parse(t *testing.T) {
	parsed, err := ParseSearchMatchMode("regex")
	if err != nil || parsed != SearchMatchModeRegex {
		t.Errorf("expected regex, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseSearchMatchModeType(" KEYWORD ")
	if errAlias != nil || parsedAlias != SearchMatchModeKeyword {
		t.Errorf("expected keyword, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseSearchMatchMode("invalid-match-mode")
	if invalidErr == nil {
		t.Errorf("expected error for invalid match mode")
	}

	if len(AllSearchMatchModeTypes()) != 4 {
		t.Errorf("expected AllSearchMatchModeTypes count to match")
	}
}

func TestSearchOutputFormatType_Valid(t *testing.T) {
	formats := AllSearchOutputFormats()
	if len(formats) != 4 {
		t.Fatalf("expected 4 output formats, got %d", len(formats))
	}

	for _, f := range formats {
		if !f.IsValid() {
			t.Errorf("expected %s to be valid", f)
		}

		if f.String() != string(f) {
			t.Errorf("expected String() == %s, got %s", string(f), f.String())
		}
	}

	if SearchOutputFormatType("xml").IsValid() {
		t.Errorf("expected invalid format to be false")
	}
}

func TestSearchOutputFormatType_Booleans(t *testing.T) {
	if !SearchOutputFormatText.IsText() || SearchOutputFormatJSON.IsText() {
		t.Errorf("unexpected IsText behavior")
	}

	if !SearchOutputFormatJSON.IsJSON() || SearchOutputFormatTable.IsJSON() {
		t.Errorf("unexpected IsJSON behavior")
	}

	if !SearchOutputFormatTable.IsTable() || SearchOutputFormatCompact.IsTable() {
		t.Errorf("unexpected IsTable behavior")
	}

	if !SearchOutputFormatCompact.IsCompact() || SearchOutputFormatText.IsCompact() {
		t.Errorf("unexpected IsCompact behavior")
	}
}

func TestSearchOutputFormatType_Parse(t *testing.T) {
	parsed, err := ParseSearchOutputFormat("json")
	if err != nil || parsed != SearchOutputFormatJSON {
		t.Errorf("expected json, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseSearchOutputFormatType(" COMPACT ")
	if errAlias != nil || parsedAlias != SearchOutputFormatCompact {
		t.Errorf("expected compact, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseSearchOutputFormat("invalid-format")
	if invalidErr == nil {
		t.Errorf("expected error for invalid format")
	}

	if len(AllSearchOutputFormatTypes()) != 4 {
		t.Errorf("expected AllSearchOutputFormatTypes count to match")
	}
}

func TestSearchCategoryType_Valid(t *testing.T) {
	cats := AllSearchCategories()
	if len(cats) != 2 {
		t.Fatalf("expected 2 search categories, got %d", len(cats))
	}

	for _, c := range cats {
		if !c.IsValid() {
			t.Errorf("expected %s to be valid", c)
		}

		if c.String() != string(c) {
			t.Errorf("expected String() == %s, got %s", string(c), c.String())
		}
	}

	if SearchCategoryType("bot").IsValid() {
		t.Errorf("expected invalid category to be false")
	}
}

func TestSearchCategoryType_Booleans(t *testing.T) {
	if !SearchCategoryUser.IsUser() || SearchCategoryAI.IsUser() {
		t.Errorf("unexpected IsUser behavior")
	}

	if !SearchCategoryAI.IsAI() || SearchCategoryUser.IsAI() {
		t.Errorf("unexpected IsAI behavior")
	}
}

func TestSearchCategoryType_Parse(t *testing.T) {
	parsed, err := ParseSearchCategory("user")
	if err != nil || parsed != SearchCategoryUser {
		t.Errorf("expected user, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseSearchCategoryType(" AI ")
	if errAlias != nil || parsedAlias != SearchCategoryAI {
		t.Errorf("expected ai, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseSearchCategory("invalid-category")
	if invalidErr == nil {
		t.Errorf("expected error for invalid category")
	}

	if len(AllSearchCategoryTypes()) != 2 {
		t.Errorf("expected AllSearchCategoryTypes count to match")
	}
}

func TestSearchActionType_Valid(t *testing.T) {
	actions := AllSearchActions()
	if len(actions) != 9 {
		t.Fatalf("expected 9 search actions, got %d", len(actions))
	}

	for _, a := range actions {
		if !a.IsValid() {
			t.Errorf("expected %s to be valid", a)
		}

		if a.String() != string(a) {
			t.Errorf("expected String() == %s, got %s", string(a), a.String())
		}
	}

	if SearchActionType("drop").IsValid() {
		t.Errorf("expected invalid action to be false")
	}
}

func TestSearchActionType_Booleans(t *testing.T) {
	if !SearchActionQuery.IsQuery() || SearchActionHistory.IsQuery() {
		t.Errorf("unexpected IsQuery behavior")
	}

	if !SearchActionHistory.IsHistory() || SearchActionTop.IsHistory() {
		t.Errorf("unexpected IsHistory behavior")
	}

	if !SearchActionTop.IsTop() || SearchActionStats.IsTop() {
		t.Errorf("unexpected IsTop behavior")
	}

	if !SearchActionStats.IsStats() || SearchActionDH2D.IsStats() {
		t.Errorf("unexpected IsStats behavior")
	}

	if !SearchActionDH2D.IsDH2D() || SearchActionClean.IsDH2D() {
		t.Errorf("unexpected IsDH2D behavior")
	}

	if !SearchActionClean.IsClean() || SearchActionClear.IsClean() {
		t.Errorf("unexpected IsClean behavior")
	}

	if !SearchActionClear.IsClear() || SearchActionPurge.IsClear() {
		t.Errorf("unexpected IsClear behavior")
	}

	if !SearchActionPurge.IsPurge() || SearchActionReset.IsPurge() {
		t.Errorf("unexpected IsPurge behavior")
	}

	if !SearchActionReset.IsReset() || SearchActionQuery.IsReset() {
		t.Errorf("unexpected IsReset behavior")
	}
}

func TestSearchActionType_CompoundBooleans(t *testing.T) {
	historyActions := []SearchActionType{
		SearchActionHistory,
		SearchActionTop,
		SearchActionStats,
		SearchActionDH2D,
	}
	for _, ha := range historyActions {
		if !ha.IsHistoryAction() {
			t.Errorf("expected %s to report IsHistoryAction() == true", ha)
		}
	}

	if SearchActionClean.IsHistoryAction() || SearchActionQuery.IsHistoryAction() {
		t.Errorf("unexpected IsHistoryAction for non-history actions")
	}

	cleanActions := []SearchActionType{
		SearchActionClean,
		SearchActionClear,
		SearchActionPurge,
	}
	for _, ca := range cleanActions {
		if !ca.IsCleanAction() {
			t.Errorf("expected %s to report IsCleanAction() == true", ca)
		}
	}

	if SearchActionQuery.IsCleanAction() || SearchActionStats.IsCleanAction() {
		t.Errorf("unexpected IsCleanAction for non-clean actions")
	}

	resetActions := []SearchActionType{
		SearchActionReset,
		SearchActionClear,
	}
	for _, ra := range resetActions {
		if !ra.IsResetAction() {
			t.Errorf("expected %s to report IsResetAction() == true", ra)
		}
	}

	if SearchActionClean.IsResetAction() || SearchActionTop.IsResetAction() {
		t.Errorf("unexpected IsResetAction for non-reset actions")
	}
}

func TestSearchActionType_Parse(t *testing.T) {
	parsed, err := ParseSearchAction("history")
	if err != nil || parsed != SearchActionHistory {
		t.Errorf("expected history, got %s, err: %v", parsed, err)
	}

	parsedAlias, errAlias := ParseSearchActionType(" PURGE ")
	if errAlias != nil || parsedAlias != SearchActionPurge {
		t.Errorf("expected purge, got %s, err: %v", parsedAlias, errAlias)
	}

	_, invalidErr := ParseSearchAction("invalid-action")
	if invalidErr == nil {
		t.Errorf("expected error for invalid action")
	}

	if len(AllSearchActionTypes()) != 9 {
		t.Errorf("expected AllSearchActionTypes count to match")
	}
}
