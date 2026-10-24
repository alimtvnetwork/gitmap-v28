package enums

import (
	"testing"
)

func TestDocsCategoryType(t *testing.T) {
	testDocsCategoryAllAndString(t)
	testDocsCategoryParse(t)
	testDocsCategoryPredicates(t)
}

func testDocsCategoryAllAndString(t *testing.T) {
	all := AllDocsCategoryTypes()
	if len(all) != 7 {
		t.Errorf("expected 7 category variants, got %d", len(all))
	}

	for _, category := range all {
		if !category.IsValid() {
			t.Errorf("expected valid for %s", category)
		}
		if category.String() != string(category) {
			t.Errorf("expected string representation to match %s", category)
		}
	}

	if len(AllDocsCategories()) != 7 {
		t.Error("expected AllDocsCategories alias to match length 7")
	}
}

func testDocsCategoryParse(t *testing.T) {
	parsed, err := ParseDocsCategoryType("  ARCHITECTURE ")
	if err != nil || parsed != DocsCategoryArchitecture {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	aliasParsed, aliasErr := ParseDocsCategory("spec")
	if aliasErr != nil || aliasParsed != DocsCategorySpec {
		t.Errorf("unexpected alias parse result: %v, err: %v", aliasParsed, aliasErr)
	}

	if _, err := ParseDocsCategoryType("invalid"); err == nil {
		t.Error("expected error parsing invalid category")
	}

	if DocsCategoryType("unknown").IsValid() {
		t.Error("expected unknown category to not be valid")
	}
}

func testDocsCategoryPredicates(t *testing.T) {
	if !DocsCategoryCLI.IsCLI() || DocsCategorySpec.IsCLI() {
		t.Error("IsCLI predicate failure")
	}
	if !DocsCategoryArchitecture.IsArchitecture() || DocsCategoryPrompt.IsArchitecture() {
		t.Error("IsArchitecture predicate failure")
	}
	if !DocsCategorySpec.IsSpec() || DocsCategoryGuide.IsSpec() {
		t.Error("IsSpec predicate failure")
	}
	if !DocsCategoryPrompt.IsPrompt() || DocsCategoryAPI.IsPrompt() {
		t.Error("IsPrompt predicate failure")
	}
	if !DocsCategoryGuide.IsGuide() || DocsCategoryCLI.IsGuide() {
		t.Error("IsGuide predicate failure")
	}
	if !DocsCategoryAPI.IsAPI() || DocsCategoryAll.IsAPI() {
		t.Error("IsAPI predicate failure")
	}
	if !DocsCategoryAll.IsAll() || DocsCategoryCLI.IsAll() {
		t.Error("IsAll predicate failure")
	}
	if !DocsCategoryCLI.HasCategoryFilter() || DocsCategoryAll.HasCategoryFilter() {
		t.Error("HasCategoryFilter predicate failure")
	}
}

func TestDocsViewerModeType(t *testing.T) {
	testDocsViewerModeAllAndString(t)
	testDocsViewerModeParse(t)
	testDocsViewerModePredicates(t)
}

func testDocsViewerModeAllAndString(t *testing.T) {
	all := AllDocsViewerModeTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 viewer mode variants, got %d", len(all))
	}

	for _, mode := range all {
		if !mode.IsValid() {
			t.Errorf("expected valid for %s", mode)
		}
		if mode.String() != string(mode) {
			t.Errorf("expected string representation to match %s", mode)
		}
	}

	if len(AllDocsViewerModes()) != 4 {
		t.Error("expected AllDocsViewerModes alias to match length 4")
	}
}

func testDocsViewerModeParse(t *testing.T) {
	parsed, err := ParseDocsViewerModeType("  BROWSER ")
	if err != nil || parsed != DocsViewerModeBrowser {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	aliasParsed, aliasErr := ParseDocsViewerMode("server")
	if aliasErr != nil || aliasParsed != DocsViewerModeServer {
		t.Errorf("unexpected alias parse result: %v, err: %v", aliasParsed, aliasErr)
	}

	if _, err := ParseDocsViewerModeType("invalid"); err == nil {
		t.Error("expected error parsing invalid viewer mode")
	}

	if DocsViewerModeType("unknown").IsValid() {
		t.Error("expected unknown viewer mode to not be valid")
	}
}

func testDocsViewerModePredicates(t *testing.T) {
	if !DocsViewerModeBrowser.IsBrowser() || DocsViewerModeTerminal.IsBrowser() {
		t.Error("IsBrowser predicate failure")
	}
	if !DocsViewerModeTerminal.IsTerminal() || DocsViewerModeBrowser.IsTerminal() {
		t.Error("IsTerminal predicate failure")
	}
	if !DocsViewerModeServer.IsServer() || DocsViewerModeWeb.IsServer() {
		t.Error("IsServer predicate failure")
	}
	if !DocsViewerModeWeb.IsWeb() || DocsViewerModeServer.IsWeb() {
		t.Error("IsWeb predicate failure")
	}
	if !DocsViewerModeServer.HasServerTarget() || DocsViewerModeBrowser.HasServerTarget() {
		t.Error("HasServerTarget predicate failure for server/browser")
	}
	if !DocsViewerModeWeb.HasServerTarget() {
		t.Error("HasServerTarget predicate failure for web")
	}
}

func TestDocsFormatType(t *testing.T) {
	testDocsFormatAllAndString(t)
	testDocsFormatParse(t)
	testDocsFormatPredicates(t)
}

func testDocsFormatAllAndString(t *testing.T) {
	all := AllDocsFormatTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 format variants, got %d", len(all))
	}

	for _, format := range all {
		if !format.IsValid() {
			t.Errorf("expected valid for %s", format)
		}
		if format.String() != string(format) {
			t.Errorf("expected string representation to match %s", format)
		}
		if !format.HasFormatting() {
			t.Errorf("expected HasFormatting to be true for %s", format)
		}
	}

	if len(AllDocsFormats()) != 4 {
		t.Error("expected AllDocsFormats alias to match length 4")
	}
}

func testDocsFormatParse(t *testing.T) {
	parsed, err := ParseDocsFormatType("  MARKDOWN ")
	if err != nil || parsed != DocsFormatMarkdown {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	aliasParsed, aliasErr := ParseDocsFormat("json")
	if aliasErr != nil || aliasParsed != DocsFormatJSON {
		t.Errorf("unexpected alias parse result: %v, err: %v", aliasParsed, aliasErr)
	}

	if _, err := ParseDocsFormatType("invalid"); err == nil {
		t.Error("expected error parsing invalid format")
	}

	if DocsFormatType("unknown").IsValid() {
		t.Error("expected unknown format to not be valid")
	}
}

func testDocsFormatPredicates(t *testing.T) {
	if !DocsFormatMarkdown.IsMarkdown() || DocsFormatHTML.IsMarkdown() {
		t.Error("IsMarkdown predicate failure")
	}
	if !DocsFormatHTML.IsHTML() || DocsFormatMarkdown.IsHTML() {
		t.Error("IsHTML predicate failure")
	}
	if !DocsFormatText.IsText() || DocsFormatJSON.IsText() {
		t.Error("IsText predicate failure")
	}
	if !DocsFormatJSON.IsJSON() || DocsFormatText.IsJSON() {
		t.Error("IsJSON predicate failure")
	}
}
