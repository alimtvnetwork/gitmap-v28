package cmdai

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// GenerateScriptContent builds the complete Python script based on options.
func GenerateScriptContent(opts CreateScriptOptions) result.Result[string] {
	slug := sanitizeScriptSlug(opts.Name)
	desc := resolveScriptDesc(opts.Description, slug)
	head := buildScriptHeader(slug, desc)
	body := buildScriptBody(opts)
	full := head + body

	return result.Ok(full)
}

func sanitizeScriptSlug(name string) string {
	clean := strings.TrimSuffix(name, ".py")
	clean = strings.ReplaceAll(clean, "_", "-")
	clean = strings.ToLower(strings.TrimSpace(clean))

	return clean
}

func resolveScriptDesc(desc string, slug string) string {
	isCustom := strings.TrimSpace(desc) != ""
	if isCustom {
		return strings.TrimSpace(desc)
	}

	return fmt.Sprintf("AI automation script for %s", slug)
}

func buildScriptHeader(slug string, desc string) string {
	lines := []string{
		"#!/usr/bin/env python3",
		`"""`,
		fmt.Sprintf("%s - %s", slug, desc),
		`"""`,
		"",
		"import argparse",
		"from importlib import import_module",
		"from pathlib import Path",
		"import sys",
		"",
		"sys.path.insert(0, str(Path(__file__).parent))",
		`engine = import_module("02-shared-engine")`,
		"",
		"DEFAULT_ENCODING = engine.DEFAULT_ENCODING",
		"LINE_SEPARATOR = engine.LINE_SEPARATOR",
		"",
	}

	return strings.Join(lines, "\n")
}

func buildScriptBody(opts CreateScriptOptions) string {
	switch opts.Type {
	case TemplateFixer:
		return buildFixerBody(opts)
	case TemplateAuditor:
		return buildAuditorBody(opts)
	case TemplateGenerator:
		return buildGeneratorBody(opts)
	case TemplateLinter, TemplateChecker, TemplateUtil:
		return buildLinterBody(opts)
	default:
		return buildLinterBody(opts)
	}
}
