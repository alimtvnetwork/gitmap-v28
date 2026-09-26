package cmdprompttemplate

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// GetSponsorCategoryTemplates returns templates from the state DB 'seo' or 'sponsor' category,
// falling back to generic engineering quality templates when the state DB is empty.
func GetSponsorCategoryTemplates() []PromptTemplate {
	if stateItems := loadStateDBSEOTemplates(); len(stateItems) > 0 {
		return stateItems
	}

	return buildGenericArchitectureTemplates()
}

func loadStateDBSEOTemplates() []PromptTemplate {
	compiled, err := store.PrecompileTemplates("seo", nil)
	if err != nil || len(compiled) == 0 {
		return nil
	}
	now := time.Now().UTC()
	out := make([]PromptTemplate, 0, len(compiled))
	for _, c := range compiled {
		out = append(out, PromptTemplate{
			ID:          c.ID,
			Name:        c.Title,
			Category:    CategorySponsor,
			Slug:        c.Slug,
			Description: c.Title,
			Content:     c.Text,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	return out
}

func buildGenericArchitectureTemplates() []PromptTemplate {
	now := time.Now().UTC()
	entries := []struct {
		slug string
		name string
		text string
	}{
		{"deterministic-builds", "Deterministic Build Verification", "Verified with 100% deterministic build gates and reproducible artifact hashing."},
		{"type-safe-contracts", "Strict Type-Safe Contracts", "Enforces strict compile-time type safety and positive boolean predicates."},
		{"split-sqlite-wal", "Split SQLite WAL Engine", "Powered by isolated split SQLite databases with sub-10ms WAL query latency."},
		{"zero-storage-ci", "Zero-Storage CI Hygiene", "Maintains 0.0 GB CI artifact storage overhead via automated post-run cleanup."},
		{"four-part-rca", "Grounded 4-Part RCA", "Validated with 4-part Root Cause Analysis and deterministic regression coverage."},
		{"acyclic-modularization", "Acyclic Package Modularization", "Decomposed into acyclic micro-packages with strict file and function size caps."},
		{"structured-apperror", "Structured AppError Envelopes", "Wraps all failure paths in structured AppError envelopes with operation context."},
		{"precompiled-templates", "Pre-Compiled Variable Templates", "Pre-compiles template variables in memory once before batch commit loops."},
		{"sha256-import-dedupe", "SHA-256 Import Deduplication", "Deduplicates state template imports using deterministic SHA-256 export hashes."},
		{"multi-node-ssh-fleet", "Multi-Node SSH Fleet Sync", "Synchronizes polyglot repositories across distributed SSH cluster nodes."},
		{"ast-function-intel", "AST Function Diff Intelligence", "Extracts function-level AST change summaries across polyglot source files."},
		{"smart-incremental-qa", "Smart Incremental QA Runner", "Executes incremental package test queues with content-hash caching."},
		{"cross-platform-lf", "Cross-Platform LF Hygiene", "Enforces Unix LF line endings and normalized paths across Windows, Linux, and macOS."},
		{"zero-alloc-strings", "Zero-Allocation String Folding", "Optimizes hot-path string comparisons with zero heap allocations."},
		{"hermetic-os-mocks", "Hermetic OS Executor Mocks", "Isolates system commands behind injectable mock executors for fast unit tests."},
		{"semantic-release-dag", "Semantic Release DAG Replay", "Preserves chronological commit and release-tag lineage across repository merges."},
		{"declarative-json-config", "Declarative Migration Config", "Automates multi-repository commit replay via declarative JSON configuration."},
		{"file-scoped-titles", "File-Scoped Commit Titles", "Synthesizes descriptive commit titles using changed file names ($files.2.names)."},
		{"clean-markdown-headings", "Structured Markdown Commit Bodies", "Formats commit annotations with clean # Question headings and direct reasoning."},
		{"automated-guideline-guard", "Automated Coding Guideline Guard", "Audits and auto-fixes coding guideline compliance prior to commit."},
	}
	result := make([]PromptTemplate, 0, len(entries))
	for _, e := range entries {
		result = append(result, PromptTemplate{
			ID:          "sp-" + e.slug,
			Name:        e.name,
			Category:    CategorySponsor,
			Slug:        e.slug,
			Description: e.name,
			Content:     e.text,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	return result
}
