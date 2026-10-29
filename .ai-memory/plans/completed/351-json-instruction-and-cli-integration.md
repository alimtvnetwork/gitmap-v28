# Master Plan: 351-json-instruction-and-cli-integration

**Sequence ID:** 351  
**Task Slug:** `351-json-instruction-and-cli-integration`  
**Specification:** [02-spec/21-app/351-json-instruction-and-cli-integration/readme.md](file;02-spec/21-app/351-json-instruction-and-cli-integration/readme.md)  
**Status:** PENDING  
**Category:** CLI & Automation / Instruction Ingestion Engine  
**Target Version:** `v6.631.0`  
**Created:** 2026-10-29  

---

## 1. Executive Summary & Problem Statement

Modern multi-agent workflows and command orchestration require machine-readable, schema-validated JSON instruction formats, seamless link expansion, structured CLI output envelopes, and lean prompt formatters.

This plan addresses five interrelated engineering challenges:
1. **JSON Instruction Format:** Unifying instruction files and payloads into a strongly typed, deterministic schema supporting task titles, slugs, actionable non-negotiable items, model/skill requirements, and metadata envelopes.
2. **GitMap CLI Integration:** Equipping GitMap commands (`gitmap space send-email`, `gitmap email expand`) with `--json` output flags, `--file` JSON ingestion, and robust link expansion via `cli/linkexpander`.
3. **Cross-Binary Parity (`go-email-sender`):** Synchronizing link expansion capabilities and structured JSON CLI output into the standalone `go-email-sender` binary and service.
4. **Prompt Condensation (< 5,000 Chars):** Pruning 12 bloated prompt formatters across `01-prompts/22-letterly/` and `01-prompts/23-cursor-prompts/` to strictly under 5,000 characters by eliminating double-declarations between Section 5 and Output Format, replacing inlined SQL/enum paragraphs in `03` with `@sql-repository-pattern-and-enum-actions`, and replacing verbose engine paragraphs in `14` with `@execute-parent-task-with-n-steps-v7` and `@gitmap-time-restriction`.
5. **Multi-Repository Fleet Synchronization:** Executing `03-ai-scripts/38-sync-prompts-skills-scripts.py --force` to mirror condensed prompts and shared skills safely across all 43 connected repositories with strict boundary protections.

---

## 2. User Request & Non-Negotiable Directives

```markdown
# JSON Instruction and CLI Integration: high priority instruction, non-negotiable task

Implement canonical JSON instruction formats and integration across GitMap CLI and go-email-sender with link expansion and JSON output options. Condense all 12 bloated prompt files across Letterly and Cursor suites to strictly under 5,000 characters, eliminate double-declarations, standardize skill links, and synchronize changes across all 43 connected repositories via Python script 38.

## task: JSON Instruction and CLI Integration
## slug: 351-json-instruction-and-cli-integration

# Actionable Items Must Follow Non-Negotiable:
1. Git pull BEFORE starting any work (first gitmap pull, fallback git pull).
2. Author specs under 02-spec/21-app/351-json-instruction-and-cli-integration/ and enqueue plan in .ai-memory/plans/pending/351-json-instruction-and-cli-integration.md (subtasks in .ai-memory/plans/subtasks/351-json-instruction-and-cli-integration/).
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String.
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cli/...); never add absolute paths.
5. Standardize JSON instruction envelope, schema validator, and ingestion pipeline.
6. Integrate --json and link expansion into GitMap CLI (gitmap space send-email, gitmap email expand).
7. Integrate --json and link expansion into go-email-sender binary with test parity.
8. Prune all 12 oversized prompt files across 01-prompts/22-letterly/ and 01-prompts/23-cursor-prompts/ to strictly under 5,000 characters; replace items 5-8 in prompt 03 with sql-repository-pattern-and-enum-actions reference; replace items 8-9 in prompt 14 with execute-parent-task-with-n-steps-v7 and gitmap-time-restriction references; eradicate Section 5 vs Output Format double-declarations.
9. Synchronize fleet via python 03-ai-scripts/38-sync-prompts-skills-scripts.py --force across all 43 connected repositories.
10. Execute terminal E2E verification, test suites to 100% green, and atomic minor bump release ceremony.
```

---

## 3. High-Level Execution Phases

### Phase 1: Planning, Discovery & Specifications
- Task-01: Deep discovery across GitMap CLI, `go-email-sender`, and prompt suites.
- Task-02: Author comprehensive specs in `02-spec/21-app/351-json-instruction-and-cli-integration/`:
  - `01-architecture-spec.md` (Schema, Link Expansion Pipeline, Invariant Rules)
  - `02-component-spec.md` (GitMap & Go-Email-Sender Component Contracts)
  - `03-prompt-condensation-and-sync-spec.md` (Prompt Pruning, Link Hygiene, Sync Protocol)
  - `readme.md` (Suite Index)
- Task-03: Author pending plan and 5 subtasks in `.ai-memory/plans/subtasks/351-json-instruction-and-cli-integration/`.

### Phase 2: JSON Instruction Schema & Expansion Engine
- Task-04: Implement schema definition and validation models (`cli/linkexpander/types.go`).
- Task-05: Implement URL expansion and metadata resolver (`cli/linkexpander/expander.go`).

### Phase 3: GitMap CLI Integration & Error Management
- Task-06: Enhance `cli/cmdspace/send_email.go` and `cli/cmdspace/space_email.go` with `--json` and `--expand`.
- Task-07: Implement leaf error sovereignty and strongly typed exit codes for email commands.

### Phase 4: Go-Email-Sender Cross-Repository Integration
- Task-08: Integrate link expander and `--json` format into `d:\work\go-email-sender`.
- Task-09: Author unit tests achieving 100% pass for email link expansion and CLI emission.

### Phase 5: Prompt Condensation & Redundancy Elimination
- Task-10: Prune 6 Letterly prompts (`02`, `03`, `06`, `09`, `10`, `14`) to strictly < 5,000 characters.
- Task-11: Prune 6 Cursor prompts (`02`, `03`, `06`, `09`, `10`, `14`) to strictly < 5,000 characters.
- Task-12: Replace items 5-8 in `03` with 2-line directive to `sql-repository-pattern-and-enum-actions`.
- Task-13: Replace items 8-9 in `14` with references to `execute-parent-task-with-n-steps-v7` and `gitmap-time-restriction`.
- Task-14: Eradicate Section 5 double-declarations and enforce relative skill links.

### Phase 6: Fleet Synchronization, E2E Verification & Release Ceremony
- Task-15: Run `03-ai-scripts/38-sync-prompts-skills-scripts.py --force` across all 43 repositories.
- Task-16: Execute verification matrix (V-01 through V-10).
- Task-17: GitMap atomic commit (`gitmap cpr "v6.631.0 json instruction and cli integration with prompt condensation"`) and tag push.

---

## 4. Subtasks Manifest

| Subtask ID | Subtask File | Description | Target Subsystems | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | [01-json-instruction-format.md](file;.ai-memory/plans/subtasks/351-json-instruction-and-cli-integration/01-json-instruction-format.md) | JSON instruction envelope, schema validation, payload serialization, and error diagnostics. | `cli/instruction/`, `cli/schema/` | PENDING |
| **Subtask-02** | [02-gitmap-cli-integration.md](file;.ai-memory/plans/subtasks/351-json-instruction-and-cli-integration/02-gitmap-cli-integration.md) | GitMap CLI integration (`space send-email`, `email expand`), link expander pipeline, `--json` flag, typed exit codes. | `cli/cmdspace/`, `cli/linkexpander/` | PENDING |
| **Subtask-03** | [03-go-email-sender-integration.md](file;.ai-memory/plans/subtasks/351-json-instruction-and-cli-integration/03-go-email-sender-integration.md) | Standalone `go-email-sender` binary parity, link expander module integration, `--json` stdout output, unit tests. | `d:\work\go-email-sender` | PENDING |
| **Subtask-04** | [04-prompt-condensation-under-5000.md](file;.ai-memory/plans/subtasks/351-json-instruction-and-cli-integration/04-prompt-condensation-under-5000.md) | Condensation of all 12 bloated prompts in `01-prompts/22-letterly/` and `01-prompts/23-cursor-prompts/` to < 5,000 characters. | `01-prompts/22-letterly/`, `01-prompts/23-cursor-prompts/` | PENDING |
| **Subtask-05** | [05-fleet-sync-and-release.md](file;.ai-memory/plans/subtasks/351-json-instruction-and-cli-integration/05-fleet-sync-and-release.md) | Fleet synchronization across all 43 repositories via Python script 38, verification matrix V-01..V-10, and release ceremony. | `03-ai-scripts/`, 43 Connected Repositories | PENDING |

---

## 5. Architectural & Quality Guardrails

### 5.1 Coding Guidelines Invariants
- **Boolean Naming:** Positive naming exclusively using `is*` or `has*` (`isJSON`, `hasExpander`, `isForce`, `isDryRun`). Total ban on `can`, `should`, `was`, `will` and negative booleans.
- **Function Limit:** <= 15 lines maximum, <= 8 lines preferred.
- **Zero Nested Ifs:** Flatten all conditional branching via early guard returns.
- **Return & Brace Spacing:** Blank line after `}` and blank line before `return`.
- **Leaf Error Sovereignty:** Leaf functions must return wrapped errors directly. Never call `cliexit.HandleError` or `os.Exit` inside leaf helpers.
- **100% Relative Paths:** Strictly relative paths (`02-spec/...`, `.ai-memory/...`, `cli/...`); zero drive letters or absolute paths.

---

## 6. Verification Acceptance Criteria Matrix (V-01 through V-10)

| ID | Gate / Acceptance Criteria | Target Deliverable | Verification Tool / Command |
| :--- | :--- | :--- | :--- |
| **V-01** | JSON Schema Validation: Complete validation for task envelope and actionable items | `cli/instruction/` | Unit test with valid and malformed JSON payloads |
| **V-02** | GitMap Link Expansion: URL unshortening and metadata resolution | `cli/linkexpander/` | Unit test checking URL resolution and timeout handling |
| **V-03** | GitMap Email CLI JSON Output: `--json` flag emits structured output envelope | `cli/cmdspace/send_email.go` | CLI invocation checking JSON stdout format |
| **V-04** | Go-Email-Sender Parity: Identical link expander and `--json` format | `d:\work\go-email-sender` | Go test suite in `go-email-sender` |
| **V-05** | Prompt Character Ceiling: All 12 target files strictly < 5,000 characters | `01-prompts/22-letterly/`, `01-prompts/23-cursor-prompts/` | Python character length audit script |
| **V-06** | Prompt 03 Directive: Items 5-8 replaced with 2-line skill reference | `03-execute-n-steps-letterly.md`, cursor counterpart | Pattern verification for `sql-repository-pattern-and-enum-actions` |
| **V-07** | Prompt 14 Directive: Items 8-9 replaced with V7 & time-restriction refs | `14-execute-recommit-time-restricted-letterly.md`, cursor counterpart | Pattern verification for `execute-parent-task-with-n-steps-v7` |
| **V-08** | Double-Declaration Elimination: Section 5 does not duplicate Output Format | All 12 modified prompt files | Line diff audit against original baseline |
| **V-09** | Fleet Synchronization: `38-sync-prompts-skills-scripts.py --force` succeeds across 43 repos | `03-ai-scripts/38-sync-prompts-skills-scripts.py` | Python execution verifying exit code 0 |
| **V-10** | Release Ceremony: Atomic commit and tag push for GitMap | GitMap repository | `gitmap cpr "v6.631.0 json instruction and cli integration with prompt condensation"` |
