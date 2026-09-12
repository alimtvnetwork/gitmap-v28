# Audit Planning for Coding Guidelines — Planning Spec (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Deeply audit the entire codebase for coding guideline violations, boolean anti-patterns, missing enums, cyclomatic complexity, and error-handling flaws. Structure all findings into actionable, fine-grained tasks in .lovable/plans/pending/ and subtasks before stopping.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## Variables — Auto-Discovered at Runtime

```text
N = 150 (Default number of steps the planning AI should take to generate the audit plan. The user may override this when triggering the prompt)
```

/learn Ingest, analyze, and internalize all coding guidelines, boolean principles, function size limits, and error handling architectures across the codebase and specs.

Autonomously self-loop and read:

- /learn the master cross-language coding guidelines in `spec/02-coding-guidelines/01-cross-language/01-index.md` through `06-advanced-patterns.md`.
- /learn the code style, braces, spacing, and multi-line rules in `spec/02-coding-guidelines/01-cross-language/01-index.md` through `06-comments-and-documentation.md`.
- /learn the strict function and type size caps (8 lines preferred, 15 lines max) in `spec/02-coding-guidelines/01-cross-language/01-index.md`.
- /learn the boolean principles, prefixing rules (is and has only (can, should, was, etc. are banned)), and guard extraction in `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md` through `05-exemptions-and-api.md`.
- /learn the absolute prohibition against negative booleans and inverted logic in `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`.
- /learn the strict identifier and file naming conventions in `spec/02-coding-guidelines/01-cross-language/01-index.md`.
- /learn the DRY principles and duplication extraction patterns in `spec/02-coding-guidelines/01-cross-language/08-dry-principles.md`.
- /learn the error management architecture and logging diagnostics in `spec/03-error-manage/01-index.md` and `spec/03-error-manage/02-error-architecture/02-error-handling-reference.md` as an AI skill checklist.
- /learn the language-specific standards in `spec/02-coding-guidelines/` (TypeScript, Go, PHP, Rust, C#, Python, PowerShell).
- /learn the anti-hallucination rules and common AI mistakes in `spec/02-coding-guidelines/01-cross-language/01-index.md` and `03-common-ai-mistakes.md`.
- Read `.lovable/plans/01-index.md` and `.lovable/memory/01-index.md`.

## 2. Planning Loop (Deep N-Step Analysis)

This is not a quick glance. You must deeply read the codebase, looping yourself as much as needed (taking exactly `N` steps of internal planning and reading). Each one of these steps MUST be followed properly using your highest processing capacity, proper memory retention of prior files, and careful multi-agent cognitive logic.

You must dedicate this immense processing power to uncover:

- Every inverted boolean (`!isSuccess`).
- Every magic string or number.
- Every swallowed error or generic `catch {}`.
- Every error without file path or variable context (`WithPath`, `WithVar`) — Rule R7.
- Every missing Enum or `Type` suffix.
- Every monolithic function exceeding 15 lines.
- Every nested `if` statement.

If there are NO discrepancies, explicitly state: "There are no coding guideline issues or discrepancies." However, assume the codebase is a mess until proven otherwise.

## 3. Root Cause & Fallout Analysis

For every issue found:

- What is the root cause? Why was it written this way?
- How many places does it need to be fixed?
- Fallout Check: If we change this, what else breaks? Will it break the CI/CD pipeline? Will it break tests? Map the entire blast radius.

## 4. Enqueueing Tasks for Sub-Agents (Bounded 5–8 Files Batches)

Your final output must be a massively detailed plan stored at `.lovable/plans/01-index.md` and granular subtask batch files written to `.lovable/plans/01-index.md`, `batch-02.md`, etc.
The plan must partition all discovered violations into **bounded micro-batches of strictly 5–8 files each** so that 2 concurrent sub-agents (max 2 threads each) can safely execute without context exhaustion or truncation:

1. [ ] **Batch Partitioning:** Group all codebase files needing fixes into 5–8 file chunks in `.lovable/plans/subtasks/01-coding-guideline-fixes/`.
2. [ ] **Exact Line Locations:** List exact file, exact line, exact boolean to rename, exact parameter struct to introduce, and exact enum to extract.
3. [ ] **Strict Planning Mode:** Do NOT fix the code in this turn. Your job is ONLY to plan, audit, and enqueue.
4. [ ] **Auto-Looping Hand-Off:** Once Phase 1 planning subtasks are written, the master orchestrator MUST self-loop directly into execution mode.
5. [ ] **Anti-Hallucination:** If referenced guidelines or files are missing, ask clarifying questions rather than guessing.

---

## 5. Coding Guidelines Strict Adherence

/learn You MUST internalize the master coding guidelines located at `01-prompts/04-coding-standards/01-coding-guidelines.md`. It contains the ultimate source of truth for Boolean rules, Function limits, Error handling, and language-specific React/Go/Python paradigms. Do not hallucinate rules; enforce exactly what is in that file.

## Metadata

- slug: plan-coding-guideline-audit-v4
- status: active

## The 4-Part RCA Requirement (Mandatory Memory File)

Before you write any code to fix the problem, you MUST document the issue in `.lovable/memory/issues/xx-<slug>.md` (where XX is the next available sequential number). The file MUST contain these exact four sections:

1. **Why it happened:** The high-level business, logical, or architectural breakdown of the failure.
2. **How it happened:** The technical execution flow that triggered the bug.
3. **Root Cause:** The exact file, line, and dependency responsible for the failure.
4. **Code Fix:** The exact code snippets showing what needed to be changed to fix the root cause.
