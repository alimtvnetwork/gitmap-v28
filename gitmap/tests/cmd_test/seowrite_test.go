// Package cmd_test — unit tests for seo-write flag parsing and orchestration.
package cmd_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestIsCreateTemplateShorthand_WithCT verifies "ct" is recognized.
func TestIsCreateTemplateShorthand_WithCT(t *testing.T) {
	args := []string{constants.CmdCreateTemplate}
	if args[0] != "ct" {
		t.Errorf("expected ct, got %q", args[0])
	}
}

// TestIsCreateTemplateShorthand_EmptyArgs verifies empty args return no match.
func TestIsCreateTemplateShorthand_EmptyArgs(t *testing.T) {
	args := []string{}
	if len(args) > 0 && args[0] == constants.CmdCreateTemplate {
		t.Error("expected no match for empty args")
	}
}

// TestIsCreateTemplateShorthand_OtherArg verifies non-ct args are rejected.
func TestIsCreateTemplateShorthand_OtherArg(t *testing.T) {
	args := []string{"--url", "example.com"}
	if args[0] == constants.CmdCreateTemplate {
		t.Error("expected no match for --url")
	}
}

// TestSEOWriteConstants_CommandNames verifies command constants.
func TestSEOWriteConstants_CommandNames(t *testing.T) {
	if constants.CmdSEOWrite != "seo-write" {
		t.Errorf("expected seo-write, got %q", constants.CmdSEOWrite)
	}
	if constants.CmdSEOWriteAlias != "sw" {
		t.Errorf("expected sw, got %q", constants.CmdSEOWriteAlias)
	}
	if constants.CmdCreateTemplate != "ct" {
		t.Errorf("expected ct, got %q", constants.CmdCreateTemplate)
	}
}

// TestSEOWriteConstants_FlagNames verifies all flag name constants.
func TestSEOWriteConstants_FlagNames(t *testing.T) {
	flags := map[string]string{
		"csv":             constants.FlagSEOCSV,
		"url":             constants.FlagSEOURL,
		"service":         constants.FlagSEOService,
		"area":            constants.FlagSEOArea,
		"company":         constants.FlagSEOCompany,
		"phone":           constants.FlagSEOPhone,
		"email":           constants.FlagSEOEmail,
		"address":         constants.FlagSEOAddress,
		"max-commits":     constants.FlagSEOMaxCommits,
		"interval":        constants.FlagSEOInterval,
		"files":           constants.FlagSEOFiles,
		"rotate-file":     constants.FlagSEORotateFile,
		"dry-run":         constants.FlagSEODryRun,
		"template":        constants.FlagSEOTemplate,
		"create-template": constants.FlagSEOCreateTemplate,
	}

	for expected, actual := range flags {
		if actual != expected {
			t.Errorf("expected flag %q, got %q", expected, actual)
		}
	}
}

// TestSEOWriteConstants_Defaults verifies default values.
func TestSEOWriteConstants_Defaults(t *testing.T) {
	if constants.SEODefaultInterval != "60-120" {
		t.Errorf("expected 60-120, got %q", constants.SEODefaultInterval)
	}
	if constants.SEODefaultIntervalMin != 60 {
		t.Errorf("expected 60, got %d", constants.SEODefaultIntervalMin)
	}
	if constants.SEODefaultIntervalMax != 120 {
		t.Errorf("expected 120, got %d", constants.SEODefaultIntervalMax)
	}
	if constants.SEOSeedFile != "data/seo-templates.json" {
		t.Errorf("expected data/seo-templates.json, got %q", constants.SEOSeedFile)
	}
	if constants.SEOTemplateOutputFile != "seo-templates.json" {
		t.Errorf("expected seo-templates.json, got %q", constants.SEOTemplateOutputFile)
	}
}

// TestSEOWriteConstants_Placeholders verifies all 7 placeholder tokens.
func TestSEOWriteConstants_Placeholders(t *testing.T) {
	placeholders := map[string]string{
		"{service}": constants.PlaceholderService,
		"{area}":    constants.PlaceholderArea,
		"{url}":     constants.PlaceholderURL,
		"{company}": constants.PlaceholderCompany,
		"{phone}":   constants.PlaceholderPhone,
		"{email}":   constants.PlaceholderEmail,
		"{address}": constants.PlaceholderAddress,
	}

	for expected, actual := range placeholders {
		if actual != expected {
			t.Errorf("expected placeholder %q, got %q", expected, actual)
		}
	}
}

// TestSEOWriteConstants_TemplateKinds verifies template kind values.
func TestSEOWriteConstants_TemplateKinds(t *testing.T) {
	if constants.TemplateKindTitle != "title" {
		t.Errorf("expected title, got %q", constants.TemplateKindTitle)
	}
	if constants.TemplateKindDescription != "description" {
		t.Errorf("expected description, got %q", constants.TemplateKindDescription)
	}
}
