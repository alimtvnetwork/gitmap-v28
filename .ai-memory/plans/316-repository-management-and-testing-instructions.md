# Master Execution Plan: Repository Management and Testing Instructions

- **Plan Path:** `.ai-memory/plans/316-repository-management-and-testing-instructions.md`
- **Slug:** `repository-management-and-testing-instructions`
- **Specification:** `02-spec/21-app/316-repository-management-and-testing-instructions/01-architecture-spec.md`
- **Status:** Active
- **Lifecycle Budget:** 400 Steps (A=3, H=2)
- **Target Releases:** Antigravity-Manager v4.254.0, gitmap v6.600.0

---

## 1. Executive Strategy

One coordinated program across the four in-scope repositories (user scope lock,
2026-10-24): verify behavior first (research wave of 3 read-only subagents), fix only
proven gaps, gate every repo on CI, then ship two minor-bump releases with tags —
Antigravity-Manager v4.254.0 and gitmap v6.600.0 — plus a root-cause report for
img-pdf-v2 and a functional verdict for yt-manage.

Preflight already satisfied (2026-10-24): all four repos pulled and clean;
img-pdf-v2 cloned via `gitmap cfr` into local folder `img-pdf` (remote verified);
dedicated `img-pdf-v2 repo` chat created; SQLite task DB initialized at
`.ai-memory/temp-agents/repository-management-and-testing-instructions/agent-task.db`.

## 2. Subtask Breakdown

| # | Subtask | Repo | File |
|---|---|---|---|
| 01 | AGM auto-switch, fast-forward, rerun verification | Antigravity-Manager | `.ai-memory/plans/subtasks/316-repository-management-and-testing-instructions/01-agm-autoswitch-fastforward-rerun-verification.md` |
| 02 | AGM CI/CD finalization and minor release | Antigravity-Manager | `.ai-memory/plans/subtasks/316-repository-management-and-testing-instructions/02-agm-cicd-and-minor-release.md` |
| 03 | gitmap recent commands review | gitmap-v28 | `.ai-memory/plans/subtasks/316-repository-management-and-testing-instructions/03-gitmap-recent-commands-review.md` |
| 04 | gitmap fast-forward logic, CI, minor release | gitmap-v28 | `.ai-memory/plans/subtasks/316-repository-management-and-testing-instructions/04-gitmap-fastforward-cicd-and-release.md` |
| 05 | img-pdf right-click and YouTube thumbnail investigation | img-pdf-v2 | `.ai-memory/plans/subtasks/316-repository-management-and-testing-instructions/05-img-pdf-right-click-and-youtube-thumbnail.md` |
| 06 | yt-manage functional verification | yt-manage | `.ai-memory/plans/subtasks/316-repository-management-and-testing-instructions/06-yt-manage-functional-verification.md` |

## 3. Execution Order

1. Research wave (parallel, read-only): subtasks 01 (research half), 03, 05, 06.
2. Implementation: fixes for proven gaps in AGM and gitmap; CI gates.
3. Releases: AGM v4.254.0, then gitmap v6.600.0, each with tags and relative-path
   release notes.
4. Final report to the user with evidence (commit SHAs, CI run IDs, release URLs).
