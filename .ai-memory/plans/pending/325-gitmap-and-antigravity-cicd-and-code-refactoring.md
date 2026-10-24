# Plan 325: GitMap and Antigravity CI/CD and Code Refactoring

Slug: gitmap-and-antigravity-cicd-and-code-refactoring
Spec: 02-spec/21-app/325-gitmap-and-antigravity-cicd-and-code-refactoring/
Status: active — overnight run, user asleep, no questions
Coordinator decisions: see 01-architecture-spec.md (ambiguity decisions 1–9)

## State at plan time (2026-10-25 ~02:55 +08)

- gitmap binary updated to v6.608.0 (manual placement; `gitmap update` placement step fails in this container, E1151; SHA256 verified; old binary backed up).
- gitmap pa: 17 pulled; 2 failures, one was gitmap-v28 (merge conflict from the pending local gosec fix eac563f2) — resolved by completing the merge and pushing d886982d; the second failure is an out-of-scope repo and is not chased.
- gitmap-v28 tip: d886982d, clean, pushed. CI runs in flight on it.
- Antigravity-Manager tip: cfd4b5258 (version 4.266.0), clean, CI green on main (CI success on the last 5 runs).
- coding-guidelines-v24 tip: 181aa98d, clean.
- Other session actively pushing to gitmap-v28 and AGM. Rule: fetch, work on latest tip, never fight it.

## Subtasks

| ID | Subtask | Owner | Owned files / repos | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | Spec + plan (this file + spec folder) | Coordinator | 02-spec/21-app/325-.../, .ai-memory/plans/pending/325-...md, .ai-memory/plans/subtasks/325-.../ | DONE | spec files written |
| Task-02 | AGM CI-green + minor-bump release | Worker AGM | Antigravity-Manager only | QUEUED | — |
| Task-03 | gitmap CI-green on latest tip | Worker gitmap-ci | gitmap-v28 CI-blocking files only | QUEUED | — |
| Task-04 | Enum/constants audit (whole repo, grouped by package) | Worker audit | read-only, gitmap-v28 cli/ | QUEUED | — |
| Task-05 | Enum/constants refactor waves (per audit groups) | Workers refactor (waves) | gitmap-v28 cli/ package groups, disjoint per wave | QUEUED | — |
| Task-06 | Prompt update: enum-per-package rule + cargo_exe example | Worker prompt | coding-guidelines-v24 enum CG execute prompt + mirrors | QUEUED | — |
| Task-07 | gitmap minor-bump release | Coordinator + worker | gitmap-v28 | QUEUED | — |
| Task-08 | Morning report | Coordinator | — | QUEUED | — |

## Wave plan

- Wave 1 (parallel, now): Task-02 (AGM release), Task-03 (gitmap CI), Task-04 (audit).
- Wave 2: Task-05 refactor waves from the audit groups, Task-06 prompt update (after the cargo_exe exemplar lands).
- Wave 3: Task-07 gitmap release, Task-08 morning report.

## Assumptions

- The other session's releases (AGM 4.266.0, gitmap v6.608.0) are legitimate user-side work; this run builds on top of them.
- AGM release uses the repo's bump script + `gitmap release`, same discipline as the v4.256.0 release.
- If the enum refactor cannot finish in one night, CI-green + both releases + a precise remaining-groups list is the acceptable fallback.

## Conflicts

- execute-parent-task-with-n-steps-v6 vs parent task on atomic commit / zero CI / zero releases: parent wins (see spec decision 7).
- Skill worker-git-ban vs parent worker execution: workers may use the gitmap suite inside their own repo only (see spec decision 8).
