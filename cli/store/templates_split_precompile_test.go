package store

import (
	"strings"
	"testing"
)

func TestExpandTemplateVariables_PascalCaseAndArray(t *testing.T) {
	vars := map[string]string{
		"CompanyName":   "RISEUP ASIA LLC",
		"AlimPortfolio": "https://alimkarim.com",
		"AlimName":      `["MD. Alim Ul Karim", "MD Alim Ul Karim", "Alim Karim"]`,
	}

	content := "Welcome to ${CompanyName}! Created by ${AlimName[0]} (${AlimPortfolio}). Also known as ${AlimName[2]}."
	expanded := ExpandTemplateVariables(content, vars)

	expected := "Welcome to RISEUP ASIA LLC! Created by MD. Alim Ul Karim (https://alimkarim.com). Also known as Alim Karim."
	if expanded != expected {
		t.Fatalf("expected %q, got %q", expected, expanded)
	}

	// Test random selection: ${AlimName}
	randomContent := "Author: ${AlimName}"
	randomExpanded := ExpandTemplateVariables(randomContent, vars)
	validOptions := []string{"Author: MD. Alim Ul Karim", "Author: MD Alim Ul Karim", "Author: Alim Karim"}
	matched := false
	for _, opt := range validOptions {
		if randomExpanded == opt {
			matched = true
			break
		}
	}
	if !matched {
		t.Fatalf("expected one of %v, got %q", validOptions, randomExpanded)
	}

	// Test normalized fallback (lowercase in vars, PascalCase in template)
	fallbackVars := map[string]string{
		"company_name": "RISEUP ASIA LLC",
	}
	normContent := "Partner: ${CompanyName}"
	normExpanded := ExpandTemplateVariables(normContent, fallbackVars)
	if !strings.Contains(normExpanded, "RISEUP ASIA LLC") {
		t.Fatalf("expected normalized match, got %q", normExpanded)
	}
}
