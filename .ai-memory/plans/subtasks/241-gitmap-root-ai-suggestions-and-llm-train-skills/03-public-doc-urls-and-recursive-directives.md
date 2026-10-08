# Subtask Plan 03: Public Markdown Specification Links & Recursive AI Directives

## Overview
Implement public documentation URL discovery and recursive Git network learning directives within `cli/cmd/llm/`. Provide direct raw GitHub URLs to authoritative Markdown docs (`llm.md`, `what-to-read.md`, `10-anti-pattern-replacements.md`, etc.) for LLM memory ingestion. Embed recursive directives for traversing multi-repo topologies, leveraging `gitmap pe history-ai` for CI failure prevention, and syncing memory ledgers.

---

## Owned Files
- `cli/cmd/llm/llm_urls.go` (new)
- `cli/cmd/llm/llm_recursive.go` (new)
- `cli/cmd/llm/llm.go`
- `llm.md`

---

## Step-by-Step Implementation

1. **Create `cli/cmd/llm/llm_urls.go`:**
   - Define `DocLink` type with `Title`, `URL`, `Description`, `Category`, `IsMandatory`.
   - Implement `GetPublicDocLinks() []DocLink` listing all authoritative Markdown documentation:
     * `llm.md`: Core LLM primer & tool catalog
     * `SKILL.md`: Native Antigravity skill
     * `what-to-read.md`: Architectural memory & spec manifest
     * `10-anti-pattern-replacements.md`: Non-negotiable command substitutions
     * `02-anti-hallucination-rules.md`: Anti-hallucination protocols
     * `04-common-ai-mistakes.md`: Common AI mistakes & guardrails
   - Implement `RenderPublicDocLinksText() string` rendering Markdown table / list with mandatory ingestion instructions.

2. **Create `cli/cmd/llm/llm_recursive.go`:**
   - Implement `RenderRecursiveInstructions() string`:
     * Multi-Repo Traversal: use `gitmap st`, `gitmap commit-in`, `gitmap commit-pull` across the 40+ connected repositories.
     * CI Error Prevention: query `gitmap pe history-ai` prior to code modifications to avoid repeated failures.
     * Memory Ledger Synchronization: update `.ai-memory/what-to-read.md` and `.ai-memory/plans/` post-task.
   - Implement `RenderSkillCreationDirective() string` instructing LLMs to formulate their own skills from training data.

3. **Update `cli/cmd/llm/llm.go` & `llm.md`:**
   - Export public doc link constants and embed them in root `llm.md`.
   - Update `llm.md` to document the new `--urls` flag and recursive directives.

---

## Acceptance Criteria
- [ ] `cli/cmd/llm/llm_urls.go` provides verified raw GitHub links to core specs.
- [ ] `cli/cmd/llm/llm_recursive.go` renders multi-repo, CI telemetry, and memory directives.
- [ ] Root `llm.md` updated with public doc registry and recursive instructions.
- [ ] No nested ifs or guideline violations.
