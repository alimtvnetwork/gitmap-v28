# Prompt Architect: Canonical AI Prompts Library

This directory hosts the canonical, production-grade prompts library (V4 architecture) for the Prompt Architect meta-repository and all connected repositories.

> **Historical Archive Note:** Legacy versioned tiers (`v1/`, `v2/`, and `v3/`) have been archived in [`06-old-prompts/`](../06-old-prompts/) inside the `coding-guidelines` meta-repository and are never synchronized to downstream repositories.

---

## Core Architecture & Capabilities

1. **Antigravity Slash Command Links + GitMap AUM Engine:**
   - Interactive slash command links (`[/goal](slashCommand;goal)`, `[/learn](slashCommand;learn)`) combined with **GitMap AUM Engine** (`gitmap` CLI) as primary and Python scripts (`03-ai-scripts/`) as fallback.

2. **High-Speed File & Content Discovery (TOTAL BAN ON `Select-String` & `git grep`):**
   - Multi-Core Streaming Live Search: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` (alias: `gitmap aum grep`) — streaming live disk text/regex scanner (replaces `Select-String`, `git grep`)
   - Instant Indexed Symbol Search: `gitmap search "<query>" [--limit <n>]` — cached SQLite symbol & keyword search
   - Universal Wildcard Search: `gitmap find "<pattern>" [-ext <ext>]` (<10ms across 10,000+ files)
   - Zero-Disk Streaming: `gitmap cat <filepath>`
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

6. **V3 Parent Task N-Steps Continuous Loop (`14-execute/10-execute-parent-task-with-n-steps-v3.md`):**
   - Top-header editable step budget (`N = 300`, `PHASE_1_STEPS = 150`, `PHASE_2_STEPS = 150`) and concurrency parameters (`A = 2`, `H = 2`).
   - **Mandatory Subagent Spawning Gate (`invoke_subagent`, Zero Solo Execution):** Requires spawning `A = 2` concurrent subagents (`H = 2` disjoint tasks per subagent, `TypeName: "self"` in Phase 2) across both Phase 1 discovery/spec generation and Phase 2 code execution.

7. **V4 Antigravity-Native Parent Task N-Steps (`14-execute/11-execute-parent-task-with-n-steps-v4.md`):**
   - Rewrite of V3 for Google Antigravity 2.0: every rule stated once (R1 to R15), a capability preflight, and a resumable gitignored ledger.
   - Lead verification of every worker report, explicit-path staging, and evidence-based confidence.

8. **V5 Antigravity-Native Ultra-Orchestrator (`14-execute/12-execute-parent-task-with-n-steps-v5.md`):**
   - Ultimate synthesis of V4's rule-indexed efficiency (R1–R16), resumable ledger (`ledger.md`), explicit path staging (R8), and evidence-gating with 100% GitMap command primacy, in-brief coding guideline injection (positive booleans, `*appfault.AppError`, <=8-15 lines), corrected Antigravity 2.0 tool schemas (`Model: "inherit"`), repo-secrets default work directory governance (R16), and non-negotiable wake-up urgency.

9. **V6 Parameter-Driven Ultra-Orchestrator (`14-execute/13-execute-parent-task-with-n-steps-v6.md`):**
   - Pure parameterization driven entirely by header variables (`N`, `A`, `H`, `C`, `PHASE_1_BUDGET`, `PHASE_2_BUDGET`, `WAVES`) with zero hardcoded literal step or agent counts in the body.
   - 100% GitMap commit primacy via atomic `gitmap cpf` / `gitmap cpb` (eliminating manual `git add` and `git commit`), upstream `.gitignore` hygiene gate (R8) with automatic untracking of ignored files (`git rm --cached`), root task JSON manifest (`task.json`) and subagent JSON contracts, self-contained worker briefs with language-specific rules, secrets gate, and <= 3,200-word footprint.

---

## Directory Index

All 24 canonical prompt categories reside directly at the root of `01-prompts/`:

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
├── 21-temp-end-to-end-tests/
├── 22-letterly/
└── 23-sync/
```


