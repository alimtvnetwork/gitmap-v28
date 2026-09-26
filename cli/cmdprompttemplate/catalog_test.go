package cmdprompttemplate

import (
	"os"
	"path/filepath"
	"strings"
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
	if len(sp) != 20 {
		t.Fatalf("expected 20 sponsor templates, got %d", len(sp))
	}

	pr := GetPRDescriptionTemplates()
	if len(pr) != 20 {
		t.Fatalf("expected 20 PR templates, got %d", len(pr))
	}

	all := GetAllCatalogTemplates()
	if len(all) != 80 {
		t.Fatalf("expected 80 total templates, got %d", len(all))
	}
}

func TestResolveTemplateContent(t *testing.T) {
	// 1. Resolve by category
	content := ResolveTemplateContent("default", 0)
	if content != "Fix the below code." {
		t.Fatalf("unexpected content for default[0]: %s", content)
	}

	// 2. Resolve by CSV
	csvPicked := ResolveTemplateContent("alpha, beta, gamma", 1)
	if csvPicked != "beta" {
		t.Fatalf("expected beta, got %s", csvPicked)
	}

	// 3. Resolve by slug
	slugPicked := ResolveTemplateContent("fix-below-code", 0)
	if slugPicked != "Fix the below code." {
		t.Fatalf("expected 'Fix the below code.', got %s", slugPicked)
	}

	// 4. Resolve random
	randomContent := ResolveRandomTemplateContent("sponsor")
	if !strings.Contains(randomContent, "RISEUP ASIA") {
		t.Fatalf("expected sponsor content with RISEUP ASIA, got %s", randomContent)
	}
}

func TestGenerateCatalogMarkdown(t *testing.T) {
	md := GenerateCatalogMarkdown()
	if len(md) < 1000 {
		t.Fatalf("expected catalog markdown > 1000 chars, got %d", len(md))
	}
	if !strings.Contains(md, "RISEUP ASIA LLC") {
		t.Fatalf("expected markdown to contain RISEUP ASIA LLC")
	}
	if !strings.Contains(md, "Marek Flejszman") {
		t.Fatalf("expected markdown to contain Marek Flejszman")
	}
	if !strings.Contains(md, "Alim Ul Karim") {
		t.Fatalf("expected markdown to contain Alim Ul Karim")
	}

	tempPath := filepath.Join("..", "..", ".ai-memory", "temp", "riseup-asia-templates.md")
	_ = os.MkdirAll(filepath.Dir(tempPath), 0755)
	_ = os.WriteFile(tempPath, []byte(md), 0644)

	specPath := filepath.Join("..", "..", "02-spec", "21-app", "148-riseup-asia-templates-and-prompt-enhancement.md")
	_ = os.MkdirAll(filepath.Dir(specPath), 0755)
	_ = os.WriteFile(specPath, []byte(md), 0644)
}
