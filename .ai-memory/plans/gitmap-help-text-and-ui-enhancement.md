# Execution Plan — gitmap-help-text-and-ui-enhancement (Program 273)

References the program ledger: `.ai-memory/plans/273-gitmap-help-text-and-ui-enhancement.md`
(governs waves + hard rules). This plan governs HOW the feature waves execute.
Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/`.

## Wave order

1. **WS2 — Help text + theme overhaul.** Subtasks 01–03 (parallel-safe):
   theme registry (spec §02), help-example data registry (spec §03),
   UI help screens rendered from the same data. No new commands.
2. **WS3 — Missing commands.** Subtask 04 audit FIRST, then subtask 05
   implements ONLY `merge write` / `mr` (spec §05). Other gaps → WS4.
3. **WS4 — 20 proposals.** Subtask 06: proposals ONLY in `.ai-memory/plans/`,
   zero code, zero version bump.
4. **WS5 — Virus RCA.** Subtask 07: document in `.ai-memory/` grounded on
   commit 385aba38e (spec §07), written for a continuing AI.
5. **WS6 — Lead verify.** Subtask 08: `go build ./...` from `cli/`, acceptance
   criteria pass (spec §97), atomic commits pushed. NO version bump, NO release.

## Subtask map

| # | Subtask file | Wave | File box |
|---|--------------|------|----------|
| 01 | `subtasks/gitmap-help-text-and-ui-enhancement/01-help-text-format.md` | WS2 | `cli/helpdoc/examples.go` (new), renderers |
| 02 | `subtasks/gitmap-help-text-and-ui-enhancement/02-theme-integration.md` | WS2 | `cli/termout/themes.go` (new), `themeflag.go` (surgical) |
| 03 | `subtasks/gitmap-help-text-and-ui-enhancement/03-help-example-data.md` | WS2 | example entries for every command |
| 04 | `subtasks/gitmap-help-text-and-ui-enhancement/04-missing-command-audit.md` | WS3 | audit report (no code changes) |
| 05 | `subtasks/gitmap-help-text-and-ui-enhancement/05-mr-implementation.md` | WS3 | `cli/cmdmerge/` (new, ≤300 lines) |
| 06 | `subtasks/gitmap-help-text-and-ui-enhancement/06-twenty-command-proposals.md` | WS4 | `.ai-memory/plans/` proposals file |
| 07 | `subtasks/gitmap-help-text-and-ui-enhancement/07-virus-rca.md` | WS5 | `.ai-memory/` RCA document |
| 08 | `subtasks/gitmap-help-text-and-ui-enhancement/08-verification.md` | WS6 | build + acceptance verification |

## Constraints (repeat in every subtask)

- Search ONLY: `gitmap aum search` / `gitmap find` / `gitmap cat` /
  `gitmap py` / `gitmap llm train`. TOTAL BAN on grep/rg/git grep/Select-String.
- RELATIVE paths only. NEVER commit binaries. NEVER delete files.
- Every file ≤300 lines. No `go test` locally. Aliases kept.
- Commits: `gitmap cpf "<module> - <summary>"` (hyphen, NO colons) + push,
  verified with `git ls-remote`. Parallel streams: `git status` before AND
  `git show --stat HEAD` after every commit.

## Definition of Done

- All 33 acceptance criteria (spec `97-acceptance-criteria.md`) verifiable.
- `go build ./...` green from `cli/`.
- Every wave committed atomically and pushed. No release cut.
