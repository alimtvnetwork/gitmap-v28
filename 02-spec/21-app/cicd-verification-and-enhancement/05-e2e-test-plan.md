# 05 — E2E Test Plan

**Version:** 1.0.0
**Updated:** 2026-10-17

---

## 1. Purpose

CLI-level end-to-end verification of the **recent features** shipped by
programs 272 and 273 — run by WS4 against the built binary. This plan covers
**what to exercise and how to record evidence**. No code changes.

## 2. Scope — features under test

| # | Feature | Origin | Probe commands (examples) |
|---|---------|--------|---------------------------|
| 1 | `mr` / `merge write` command | program 273 | `gitmap mr --help`, `gitmap merge write --help`, `gitmap mr --dry-run` (if supported) |
| 2 | `gitmap help example <cmd>` 3-part format | program 273 | `gitmap help example mr`, `gitmap help example scan`, spot-check 5+ verbs |
| 3 | Four new themes | program 273 | `--theme navy-purple` / `tokyo-night` / `vscode-dark-plus` / `monokai-pro` on a help/table command |
| 4 | `gitmap nc` build-time codegen catalog | program 272 | `gitmap nc --help`, catalog listing; confirm entries are generated (not hardcoded) |

## 3. Method — dry-run / `--help` only

- **No live side effects.** No real PRs created (`mr` is exercised with
  `--help` / `--dry-run` only — never an actual `gh pr create`), no pushes,
  no external mutations, no network writes.
- **No test suites.** WS4 does not run `go test` or any suite; this is
  black-box CLI probing, not unit testing.
- Per probe, record: exact command run, exit code, and a 3-line digest of
  output (first line, a representative middle line, last line) — enough to
  prove it ran and rendered.
- Check each probe for: exit code 0; help text renders (no template errors,
  no `<no value>` leaks); piped output stays clean (no ANSI corruption in
  pipes — dumb-terminal safety); theme actually changes rendering across the
  four themes; `help example` shows all three parts (what-it-does +
  invocation + expected output).

## 4. Failure handling

- A probe that fails (non-zero exit, garbled output, missing example) is
  recorded as-is with the exact command and output digest — never smoothed
  over, never re-run until it passes.
- WS4 does **not** fix failures. It records them. Fixes are a lead decision
  (possible WS2-adjacent loop or a follow-up program).
- If a probe cannot run safely (would cause a side effect), mark it
  `SKIPPED — unsafe` with the reason; do not force it.

## 5. Evidence output (WS4)

WS4 writes its run log under `.ai-memory/plans/subtasks/` (evidence file
owned by WS4's worker, named in its subtask file) containing:

1. Program version under test (`version.json` Version at run time).
2. One row per probe: command, exit code, output digest, PASS/FAIL/SKIPPED.
3. A short "surprises" list — anything that behaved unexpectedly, even if
   technically passing.

## 6. Critical feedback — reserved for a later wave

The owner's program includes **critical feedback** on the E2E results. That
feedback is written to:

`02-spec/21-app/cicd-verification-and-enhancement/06-e2e-critical-feedback.md`

**Do NOT create that file in Wave 1** (per program brief — it belongs to a
later worker who will read WS4's evidence and write the critique). WS4's job
ends at honest evidence; the critique is a separate pass with fresh eyes.
