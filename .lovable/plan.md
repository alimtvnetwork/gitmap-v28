# Rejog the Memory v1 — Execution Plan

Read-only synthesis pass. No app code, no spec content changes. Only new artifacts under `.lovable/memory/` and a root `plan.md`.

## Scope of reading (inputs)

1. `.lovable/` — `overview.md`, `strictly-avoid.md`, `prompt.md`, `suggestions.md`, `cicd-index.md`, `issues/`, `cicd-issues/`, `solved-issues/`, `pending-issues/`, `question-and-ambiguity/`, `audits/`, `plans/`, `spec/`, `archive/` (skim only).
2. `.lovable/memory/` — `index.md`, `project/`, `tech/`, `style/`, `constraints/`, `workflow/`, `features/` (60+ files), `plans/`, `issues/`, `suggestions/`, session notes `01-…` `02-…` `03-…`.
3. `spec/` — `01-app/` (116+ numbered specs), `02-app-issues/`, `03-commit-in/`, `03-general/`, `03-tasks/`, `04-generic-cli/`, `05-coding-guidelines/`, `06-design-system/`, `07-generic-release/`, `08-generic-update/`, `08-json-schemas/`, `09-pipeline/`, `12-consolidated-guidelines/`, `15-research/`, `README.md`.
4. Explicitly skipped: any `skipped/` folder, `.gitmap/release/` payloads, generated `test-output.txt`, and Go source under `gitmap/` / `gitmap-updater/` (only referenced, not audited).

Reading is delegated to parallel subagent explorations grouped by area (memory, spec/01-app, spec generic-cli + coding-guidelines, spec pipeline + release + update, design-system + docs-site) to keep context small.

## Deliverable 1 — Reliability & Failure-Chance Report

File: `.lovable/memory/reports/YYYYMMDD-rejog-reliability.md`

Sections:
1. Corpus inventory — counts per folder, notable gaps, stale/duplicate specs (e.g. duplicated `108-…`, `109-…`, `110-…`, `111-…` numeric prefixes in `spec/01-app/`).
2. Success probability by tier:
   - Simple (single-file spec, isolated CLI flag): estimate + assumptions.
   - Medium (multi-file feature, e.g. clone-next, ssh flag).
   - Complex agentic (release/self-update, history-rewrite, cross-platform install).
   - End-to-end (v6 migration, bulk visibility, docs site parity).
3. Failure map — table of module x likely failure mode x symptom, drawn from `.lovable/issues/`, `02-app-issues/`, `cicd-issues/`, `question-and-ambiguity/`.
4. Corrective actions — prioritized spec fixes (what, where, expected reliability gain).
5. Readiness decision — go / fix-first, with the blocking list.

## Deliverable 2 — Suggestions workflow contract

File: `.lovable/memory/suggestions/README.md` documenting the filesystem contract:
- Location: `.lovable/memory/suggestions/`
- Naming: `YYYYMMDD-HHMMSS-suggestion-<slug>.md`
- Frontmatter fields: `suggestionId`, `createdAt`, `source`, `affectedProject`, `status` (`open|inProgress|done`).
- Body sections: description, rationale, proposed change, acceptance criteria, completion notes.
- Completion: flip `status: done`, add completion notes; archival policy = move to `.lovable/memory/suggestions/archive/` after done (kept, not deleted).
- Template file: `.lovable/memory/suggestions/_template.md`.

Existing `01-suggestions.md` is left in place (legacy log); the README notes it as historical.

## Deliverable 3 — Root `plan.md`

File: `plan.md` at repo root. Sections:
1. Prioritized backlog grouped by phase (P0 stabilize, P1 feature completion, P2 polish) and by project (`gitmap` CLI, `gitmap-updater`, docs site under `src/`, CI/CD pipeline, spec hygiene).
2. Per-task block: objective, dependencies, expected outputs (spec files / UI / API), acceptance criteria.
3. "Next task selection" section listing 3-5 ready-to-pick items with a one-line rationale each, sourced from the readiness decision in Deliverable 1.

## Non-goals

- No edits to any existing spec file, source file, or CI config.
- No implementation of any listed task.
- No touching `skipped/` folders or `.gitmap/release/` payloads.

## Acceptance criteria

- Three artifacts exist: reliability report under `.lovable/memory/reports/`, suggestions README + template under `.lovable/memory/suggestions/`, root `plan.md`.
- Report cites concrete file paths for every claim.
- `plan.md` next-task list is short enough for the user to pick one in a single reply.
- After artifacts land, the assistant asks the user which task to implement next.

## Technical notes

- Use parallel `acp_subagent--explore` runs per area to avoid pulling every file into the main context.
- Keep each artifact under ~500 lines; link out to source files instead of quoting them wholesale.
- Timestamp uses UTC at generation time.
