package llm

import (
	"fmt"
	"strings"
)

// DocLink represents an authoritative public documentation URL for LLM ingestion.
type DocLink struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsMandatory bool   `json:"is_mandatory"`
}

const (
	PublicSkillURL             = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/.agents/skills/gitmap/SKILL.md"
	PublicWhatToReadURL        = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/what-to-read.md"
	PublicAiMemoryURL          = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/.ai-memory/what-to-read.md"
	PublicAntiPatternURL       = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md"
	PublicAntiHallucinationURL = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/02-spec/02-coding-guidelines/06-ai-optimization/02-anti-hallucination-rules.md"
	PublicCommonMistakesURL    = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md"
)

// GetPublicDocLinks returns the authoritative public documentation registry.
func GetPublicDocLinks() []DocLink {
	return []DocLink{
		{
			Title:       "Core LLM Specification & AI Guidelines",
			URL:         PublicLlmSpecURL,
			Description: "Authoritative 6-phase execution lifecycle, CLI tool reference, and pair programming standards.",
			Category:    "Specification",
			IsMandatory: true,
		},
		{
			Title:       "GitMap Native Antigravity Skill",
			URL:         PublicSkillURL,
			Description: "Comprehensive cheatsheet, regex search rules, replacement matrices, and toolchain commands.",
			Category:    "Skill",
			IsMandatory: true,
		},
		{
			Title:       "What-To-Read Memory Manifest",
			URL:         PublicWhatToReadURL,
			Description: "Canonical reading order, code-red principles, active specifications, and persistent ledger.",
			Category:    "Memory",
			IsMandatory: true,
		},
		{
			Title:       "AI Memory Ledger & Issue History",
			URL:         PublicAiMemoryURL,
			Description: "Historical context, resolved issues, regression prevention dossiers, and plan indices.",
			Category:    "Memory",
			IsMandatory: true,
		},
		{
			Title:       "Anti-Pattern Command Replacement Matrix",
			URL:         PublicAntiPatternURL,
			Description: "Mandatory substitutions banning Select-String, Get-ChildItem -Filter, git grep, and findstr.",
			Category:    "Guideline",
			IsMandatory: true,
		},
		{
			Title:       "Anti-Hallucination Rules & Verification",
			URL:         PublicAntiHallucinationURL,
			Description: "Strict zero-assumption policy, real ground truth verification, and relative path hygiene.",
			Category:    "Guideline",
			IsMandatory: true,
		},
		{
			Title:       "Common AI Mistakes & Guardrails",
			URL:         PublicCommonMistakesURL,
			Description: "Catalog of recurring AI failure modes and compiler/runtime crash prevention guardrails.",
			Category:    "Guideline",
			IsMandatory: true,
		},
	}
}

// RenderPublicDocLinksText renders formatted documentation links for terminal and LLM ingestion.
func RenderPublicDocLinksText() string {
	links := GetPublicDocLinks()
	var b strings.Builder

	b.WriteString("======================================================================\n")
	b.WriteString("STAGE 3: AUTHORITATIVE PUBLIC DOC REGISTRY (FOR LLM MEMORY INGESTION)\n")
	b.WriteString("======================================================================\n")
	b.WriteString("Instruction: AI models must fetch and read these documents to expand memory:\n\n")

	for i, link := range links {
		mandTag := ""
		if link.IsMandatory {
			mandTag = " [MANDATORY]"
		}
		b.WriteString(fmt.Sprintf("%d. %s%s\n", i+1, link.Title, mandTag))
		b.WriteString(fmt.Sprintf("   URL:  %s\n", link.URL))
		b.WriteString(fmt.Sprintf("   Desc: %s\n\n", link.Description))
	}

	return b.String()
}
