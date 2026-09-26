package cmdprompttemplate

import (
	"testing"
)

func TestCatalogTemplatesCount(t *testing.T) {
	def := GetDefaultCategoryTemplates()
	if len(def) != 20 {
		t.Fatalf("expected 20 default templates, got %d", len(def))
	}

	ui := GetUIUXCategoryTemplates()
	if len(ui) != 20 {
		t.Fatalf("expected 20 UI/UX templates, got %d", len(ui))
	}

	sp := GetSponsorCategoryTemplates()
	if len(sp) < 20 {
		t.Fatalf("expected at least 20 sponsor/seo templates, got %d", len(sp))
	}

	pr := GetPRDescriptionTemplates()
	if len(pr) != 20 {
		t.Fatalf("expected 20 PR templates, got %d", len(pr))
	}

	all := GetAllCatalogTemplates()
	if len(all) < 80 {
		t.Fatalf("expected at least 80 total templates, got %d", len(all))
	}
}

func TestResolveTemplateContent(t *testing.T) {
	content := ResolveTemplateContent("default", 0)
	if content != "Fix the below code." {
		t.Fatalf("unexpected content for default[0]: %s", content)
	}

	csvPicked := ResolveTemplateContent("alpha, beta, gamma", 1)
	if csvPicked != "beta" {
		t.Fatalf("expected beta, got %s", csvPicked)
	}

	slugPicked := ResolveTemplateContent("fix-below-code", 0)
	if slugPicked != "Fix the below code." {
		t.Fatalf("expected 'Fix the below code.', got %s", slugPicked)
	}

	randomContent := ResolveRandomTemplateContent("sponsor")
	if len(randomContent) == 0 {
		t.Fatalf("expected non-empty sponsor content")
	}
}

func TestGenerateCatalogMarkdown(t *testing.T) {
	md := GenerateCatalogMarkdown()
	if len(md) < 500 {
		t.Fatalf("expected catalog markdown > 500 chars, got %d", len(md))
	}
}
