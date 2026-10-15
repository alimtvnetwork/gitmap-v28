# Subtask 04 — Missing-command audit: dispatcher vs help (program 273, WS3)

## File box

- Audit report ONLY — no code changes. Write the report as a dated markdown
  file inside
  `.ai-memory/plans/subtasks/gitmap-help-text-and-ui-enhancement/`
  (e.g. `audit-report-2026-10-16.md`). Read-only work elsewhere.

## Context

Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/04-missing-command-audit.md`.
Procedure: (1) enumerate dispatcher verbs+aliases from `cli/cmd/root*.go`
dispatch tables; (2) enumerate help entries (`cli/cmd/help.go`,
`cli/helpdoc/*.md`, subtask-01 example registry); (3) Diff A = verbs without
help; Diff B = help without verbs; Diff C = verbs without examples;
(4) Lovable-migration scan (PR/merge surface = P0); (5) decision per gap:
implement / add help / alias / remove promise. Known headline gap:
`merge write` / `mr` (spec §05 — implemented in subtask 05). Rules: ONLY
`gitmap aum search` / `gitmap find` / `gitmap cat` (TOTAL BAN on grep/rg);
relative paths; report ≤300 lines.

## Checklist

- [ ] Full verb→alias map published in the report.
- [ ] Diff A/B/C tables with a decision per gap.
- [ ] Lovable-migration surface ranked P0 with rationale.
- [ ] Non-`mr` gaps routed: help-only gaps → subtask 01/03 backlog; code
  gaps → WS4 proposal candidates (spec §06).
- [ ] Commit: `gitmap cpf "audit - missing-command dispatcher vs help"` + push
  (verify `git ls-remote`).
