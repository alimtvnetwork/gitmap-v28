# Prompt Architect: Canonical AI Prompts Library

This directory hosts the canonical, production-grade prompts library (V4 architecture) for the Prompt Architect meta-repository and all connected repositories.

> **Historical Archive Note:** Legacy versioned tiers (`v1/`, `v2/`, and `v3/`) have been archived in [`06-old-prompts/`](../06-old-prompts/) inside the `coding-guidelines` meta-repository and are never synchronized to downstream repositories.

---

## Core Architecture & Capabilities

1. **Antigravity Slash Command Links + GitMap AUM Engine:**
   - Interactive slash command links (`[/goal](slashCommand:goal)`, `[/learn](slashCommand:learn)`) combined with **GitMap AUM Engine** (`gitmap` CLI) as primary and Python scripts (`03-ai-scripts/`) as fallback.

2. **High-Speed File & Content Discovery:**
   - Universal Wildcard Search: `gitmap find "<pattern>" [-ext <ext>]` (<10ms across 10,000+ files)
   - Zero-Disk Streaming: `gitmap cat <filepath>`
   - Instant Filesystem Walk: `gitmap search "<symbol>"`
   - Substring Filename Search: `gitmap find-files-any "<str>"` (alias: `gitmap ffa`)
   - Directory Indexing: `gitmap list-files [pattern]` (alias: `gitmap lf`)

3. **Pipeline AI & Zero-Credit-Waste Dynamic Waiting:**
   - Telemetry Query: `gitmap pipeline-ai status --json`
   - Adaptive Sleep: `gitmap pipeline-ai status -t <etaSeconds>` (replaces expensive busy-polling loops)
   - Failure Log Snippets: `gitmap pipeline errors` (alias: `gitmap pe`)
   - Runner Target Diagnostics: `gitmap pipeline details` (alias: `gitmap pd`)

4. **Autonomous Automated Release:**
   - Single-Command Release: `gitmap release --bump <patch|minor> -y` (alias: `gitmap r -y`)
   - Release History: `gitmap changelog` (alias: `gitmap cl`), `gitmap list-versions` (alias: `gitmap lv`)

5. **Instruction Precedence Mandates:**
   - **Standard Execute & Audit Prompts:** Enforce **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence)** — instructions provided ABOVE the prompt take highest priority.
   - **Below-Steps Prompts (`*-in-below-steps.md`, `32-cg-follow-other-prompts.md`):** Enforce **Bottom-Instruction Priority Mandate (Below Precedence)** — instructions appended after the trailing `--` and `## 🚨 High Priority Instructions Below` header take highest priority.

---

## Directory Index

All 22 canonical prompt categories reside directly at the root of `01-prompts/`:

```text
01-prompts/
├── 00-folder-structure/
├── 01-prompt-library-setup/
├── 02-core-workflow/
├── 03-read-write/
├── 04-coding-standards/
├── 05-coding-guidelines/
├── 06-testing-and-qa/
├── 07-bug-fix/
├── 08-dry-code/
├── 09-commit-and-multi-agent-code-fix/
├── 10-ui-and-design/
├── 11-content-and-seo/
├── 12-old-plan-prompts/
├── 13-plan-audit/
├── 14-execute/
├── 15-cg-execute/
├── 16-ci-cd/
├── 17-release-management/
├── 18-insults/
├── 19-old-execute-prompts/
├── 20-ai-fix-script-prompts/
└── 21-temp-end-to-end-tests/
```
