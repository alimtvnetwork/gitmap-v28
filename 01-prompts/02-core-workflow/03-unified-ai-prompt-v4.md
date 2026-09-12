# Unified AI Autonomous Execution Protocol — Core Protocol (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## Part 1 - Repository Analysis, Memory Reconstruction, and Implementation Readiness

### Proofread instruction

Read and synthesize existing repository context from the Lovable memory folder and the full specification set, then produce a reliability risk report before any implementation work begins. Do not implement anything. Only produce a report and specification-side artifacts for memory, suggestions, and planning.

### Mandatory pre-analysis steps

Before producing any report or analysis, the AI must:

1. Scan the repository tree rapidly using `python 03-ai-scripts/17-fast-file-reader.py` (which leverages pre-warmed `tmp/cache/` in <15ms) instead of slow manual shell scans. Do not read contents inside folders marked skipped, ignored, deprecated, generated, archived, or otherwise excluded.
   - **If the script is missing:** Immediately recreate it (`03-ai-scripts/17-fast-file-reader.py`) using Python standard libraries. It MUST support flags: `--list-folder <path> [--ext .md,.ts]`, `--read-file <path> [--max-bytes N]`, and `--search-pattern "<regex>" [--path <dir>]`. Ensure strict UTF-8 output (`sys.stdout.reconfigure(encoding="utf-8")`) and implement local caching.
2. Read workflow memory - specifically `.lovable/plan.md` - to understand what has been done and what is pending. This avoids repeated work.
3. Read all relevant memory files under `.lovable/memory/`, including workflow, suggestions, rules, decisions, history, issue references, and any protocol or process files present.

/goal 1. Reconstruct project requirements by reading:

   1. the .lovable memory content
   2. the existing spec files and idea files across all projects
2. Produce a detailed risk and failure-chance report for handing the current specs to another AI.
3. Establish a disciplined workflow for Lovable suggestions tracking and future planning so another AI can continue work reliably.

### Inputs to read

1. .lovable/
   1. memories/
   2. memory/
   3. memory/suggestions/
   4. any other Lovable state folders present
   5. What todo and what not to do - remember.
   6. Folders marked skipped, ignored, deprecated, generated, or archived must not be read or modified - they may be listed structurally but their contents must not be opened.
2. Spec folder content for all projects:
   1. ideas
   2. backend and frontend specs
   3. specs
   4. instruction builder specs
   5. seeding and configuration specs
   6. data model specs
   7. acceptance criteria specs
   8. Read root `spec/` folder or get a general idea of files.

### Deliverable 1: Reliability and failure-chance report

1. Success probability estimates by module complexity tier (simple, medium, complex agentic workflows, end-to-end), with explicit assumptions.
2. Failure map - where (module + workflow), why (missing constraints, ambiguity, cross-file inconsistency), how it manifests.
3. Corrective actions - prioritized list of spec fixes; for each: what to change, where, expected reliability gain.
4. Readiness decision - classify each major area as: ready / ready with assumptions / blocked by ambiguity / blocked by contradiction / blocked by missing acceptance criteria. State what must be fixed before implementation and what can be deferred safely.

### Deliverable 2: Lovable suggestions workflow (filesystem contract)

All suggestions must be tracked in a single file: `.lovable/suggestions.md`. If the file grows beyond manageable size (50+ items), suggestions may be split per project.

Suggestion entry fields: suggestionId, createdAt, source (Lovable), affectedProject, description, rationale, proposed change, acceptance criteria, status (open, inProgress, done), completion notes.

Completion handling - When a suggestion is completed, update its status to done. Optionally move completed items to `completed/` subfolder.

### Deliverable 3: .lovable/plan.md future work roadmap

`.lovable/plan.md` is the canonical workflow tracker. Root `.lovable/plan.md` (if created) is a summarized AI handoff roadmap only. It must not contradict the canonical plan.

.lovable/plan.md requirements:

1. A prioritized backlog of tasks
2. Grouping by phase and by project
3. For each task: objective, dependencies, expected outputs (spec updates, UI changes, API changes), acceptance criteria
4. A section titled Next task selection where the next implementable items are listed so the user can pick what to implement next.

### Required Part 1 artifacts

- `.lovable/plan.md`
- `.lovable/suggestions.md`
- `.lovable/memory/01-index.md` - create the `history/` folder if it does not exist
- `.lovable/strictly-avoid.md` - if new rules or constraints are discovered
- Root `.lovable/plan.md` - only if a handoff roadmap is needed
- Update memory issue references if analysis uncovers prior unresolved issue patterns

### Interaction rule

After producing the report and creating the memory and plan artifacts, ask which specific task should be implemented next since the specs define what to build.

---

## Part 2 - Specification Fix Workflow and Issue Documentation

### Original input (verbatim)

update spec properly so that the mistake doesn't appear and update memory and also write the details how you fixed it, and every time we fix it, add to do /spec/02-app/issues/01-{issue slug name}.md explain the issue first, then root cause analysis, how you fixed and how not to repeat it again, and if iterations required, then write all the iterations And put all the spec files in 01-app Keep it in your memory to update all the time so that mistakes don't happen this is the most important part the many times i have remind the mistakes make sure to update in the

### Objectives

1. Update the relevant spec files so the mistake cannot happen again.
2. Update memory documentation to record the mistake, the fix, and the prevention rule.
3. Every time a fix is made, create a dedicated issue documentation file under the required path.
4. Consolidate all application spec files under a single folder.

### Required folder structure

- All application spec files: `/spec/01-app/`
- All issue write-ups: `/spec/02-app/issues/`
- File naming format: `{seq}-{issueSlugName}.md` (sequential numbering: 01, 02, 03…)

### issueSlugName rules

- lowercase only
- hyphen-separated
- short, descriptive, and stable
- no spaces or special characters

### Issue numbering

Issues use sequential numbering across the entire issues folder. Before creating a new issue, check the highest existing sequence number and increment by one.

```
/spec/02-app/issues/01-auth-timeout.md
/spec/02-app/issues/02-cache-race-condition.md
/spec/02-app/issues/03-missing-default-config.md
```

### Issue write-up file requirements

Create an issue file at `/spec/02-app/issues/{seq}-{issueSlugName}.md`. Sections in this exact order:

Issue summary - what happened, where (feature/module + paths), symptoms and impact, how discovered.

Root cause analysis - direct cause, contributing factors, triggering conditions, why existing spec did not prevent it.

Fix description - spec changes (no code), new rules/constraints, why fix resolves root cause, config/default changes, logging or diagnostics required.

Iterations history (only if multiple attempts) - Iteration 1: what was tried and why it failed. Iteration 2: ... continue until final resolution.

Prevention and non-regression - prevention rule, acceptance criteria/test scenarios, guardrails or linting policies, references to exact spec sections updated (by file path), explicit testable regression prevention rule.

TODO and follow-ups - remaining tasks, owners or roles if applicable.

Done checklist

- [ ] Spec updated under /spec/01-app/
- [ ] Issue write-up created under /spec/02-app/issues/
- [ ] Memory updated with summary and prevention rule
- [ ] Acceptance criteria updated or added
- [ ] Iterations recorded if applicable
- [ ] Plan status updated in workflow tracker

### Spec update requirements

Update the relevant spec files under /spec/01-app/ to include: corrected behavior, explicit constraints to prevent the old mistake, failure modes and debugging guidance, acceptance criteria updates that make regression testable, and a Known pitfalls and prevention section that references the issue file path.

### Memory update requirements

Maintain a memory record updated every time a fix is made: short description of the mistake, prevention rule, reference to the issue write-up file path.

Memory update is mandatory. If memory is not updated the fix is incomplete.

### Decision logging

All important decisions must be written to `.lovable/memory/01-index.md`. If the `history/` folder does not exist, create it and use this file as the canonical decision log.

Required entries: architecture changes, spec interpretation decisions, rejected approaches and why, trade-off resolutions.

### Output requirements

1. A concise process checklist to follow after every fix.
2. A copy-paste template for `/spec/02-app/issues/{seq}-{issueSlugName}.md`.
3. A brief note stating all specs live under `/spec/01-app/`.

Formatting rule: ensure there is a blank line after every Markdown header.

---

## Part 3 - Unit Test Failure Investigation and Documentation

### Original input (verbatim)

Fix these and when fixing failing tests: 1. check code, 2. Method code actual one, 3. Logical implementation of the test, 4. Check Testcase, 5. Logically fix it either actual or the test depending on the logical discussion and write it.

### Unit Test Failure Logic

1. Check code - read the production code under test.
2. Method code actual implementation - understand what the method actually does.
3. Logical implementation of the test - read the test and understand what it asserts.
4. Check the testcase logic - verify whether the test expectation is logically correct, including fixtures, mocks, seed data, and expected outputs.
5. Decide - fix either the implementation or the test depending on which is logically wrong.

### Documentation requirement for failing tests

Every failing test resolution must be documented at: `/spec/05-failing-tests/{seq}-failing-test-name.md`. Include:

1. Root cause analysis
2. Solution description
3. Whether the issue was caused by incorrect implementation or incorrect test logic
4. Any corrections made to test logic or implementation logic
5. Prevention guidance so similar failures do not occur again
6. Reference to the relevant spec section that governs the expected behavior

---

## Specification Authority

The specification is the source of truth.

Priority order (highest to lowest):

1. Specification files under `/spec/01-app/`
2. Issue corrections under `/spec/02-app/issues/`
3. Failing test documentation under `/spec/05-failing-tests/`
4. Memory and decision logs
5. Existing implementation code

If implementation contradicts the specification, the specification takes precedence unless the spec is proven incorrect. If the spec is proven incorrect, document the correction as an issue before changing the spec.

---

## Required Execution Order

The AI must follow this sequence strictly. Steps must not be skipped or reordered.

1. Scan the entire repository tree.
2. Read Lovable memory folders.
3. Read workflow tracker `.lovable/plan.md`.
4. Read specification folders.
5. Reconstruct project context.
6. Produce the reliability and failure-chance report.
7. Propose spec corrections if required.
8. Update memory artifacts.
9. Update workflow plan.
10. Ask the user which task to implement next.

---

## Context Preservation

After any significant analysis, correction, or decision, the AI must summarize the key conclusions into Lovable memory files. Summaries must include: what was learned, what was decided, what constraints now exist, what must not be repeated.

If the AI does not persist conclusions to memory, the work is considered incomplete.

---

## Task Selection Protocol

When asking for the next task, present:

1. The top 3 next implementable tasks from the plan
2. Their dependencies (what must be done first)
3. Their estimated complexity (simple / medium / complex)
4. The spec files involved

Then ask the user to select the task number.

---

## Blocker Handling

If a blocker prevents reliable implementation or specification updates:

1. Record the blocker in `.lovable/plan.md`
2. Document it in the relevant spec or issue file
3. Explain the minimum information or change required to unblock progress
4. Avoid guessing past the blocker

The AI must not silently work around blockers. Blockers must be surfaced to the user.

---

## Allowed Actions Before Implementation

Analysis, reporting, memory updates, planning, spec corrections, issue documentation, and test diagnosis documentation are allowed before implementation. Application code changes are not allowed until the user explicitly selects a task and authorizes implementation.

---

## Global Rules (apply to all parts)

### Spec before code

Always write or update specs before any implementation. Never implement until the user explicitly says to start a specific phase or task.

### Ambiguity handling

If the specification is ambiguous, the AI must document the ambiguity in the relevant spec file and in `.lovable/memory/01-index.md` before implementing a solution. Do not silently resolve ambiguity.

### Repository scan requirement

Before implementation analysis, scan the entire repository tree at the directory level. Do not read contents inside folders marked excluded. Reading only the spec folder is insufficient.

### Skipped folders

Folders marked skipped, ignored, deprecated, generated, or archived must not be read or modified. They may be listed structurally but their contents must not be opened. This overrides any other instruction.

### Code style (GitMap enforced)

- All `if` conditions must be positive (no `!`, no negation).
- Functions: 8-15 lines.
- Files: 100-200 lines max.
- Small, focused packages - one responsibility per package.
- Mandatory file path & variable context in all errors: Any file, directory, or validation error MUST embed the resolved file path (`.WithPath(p)`) and variable name (`.WithVar(name, val)`). Generic error messages without target path/variable metadata are strictly banned.

### Version bumping

Any changes to code must bump at least the minor version. The `.release` folder is off-limits - do not read, modify, or reference it.

### File naming

- Use stable canonical filenames such as `.lovable/plan.md`, `01-suggestions.md`, `01-decisions.md` for singleton tracker files.
- Use `{seq}-{slug}.md` for repeating records such as issues and failing test write-ups.
- Keep folder file counts small.
- Plans and suggestions are tracked in single files and updated in place unless explicitly split by scale.

### Regression prevention

Every fix - specs, code, or tests - must include an explicit, testable regression prevention rule.

### Definition of Done

A task is done only when:

1. Spec updated (if applicable)
2. Issue documented (if applicable)
3. Memory updated
4. Acceptance criteria added or verified
5. Plan status updated in `.lovable/plan.md`
6. Decision log updated (if a decision was made)

---

## Final Instruction

Implementation must not begin until readiness analysis and specification validation are completed and the user explicitly selects the next task.

Use:

1. Required Execution Order for sequencing
2. Specification Authority for conflict resolution
3. Context Preservation for memory persistence
4. Blocker Handling for unresolved situations

*Prompt v4.0. Trigger phrase: "unified ai prompt".*

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] **Error Context Verification:** Verify that every error created or wrapped embeds exact target file paths, variable names, and multi-path metadata without generic messages.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.

- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.
