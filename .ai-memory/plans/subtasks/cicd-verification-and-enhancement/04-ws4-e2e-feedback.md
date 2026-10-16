# Subtask 04 — WS4: E2E of recent features

**Version:** 1.0.0
**Updated:** 2026-10-17
**Owner:** E2E worker
**Status:** pending (may run in parallel with WS2/WS3)

---

## Goal

CLI-level end-to-end verification of the recent features from programs 272
and 273, per `05-e2e-test-plan.md` — honest evidence, no fixes, no side
effects.

## Inputs

- `05-e2e-test-plan.md` — E2E plan (normative).
- Built `gitmap` binary at the version under test.

## Steps

1. Record the program version under test (`version.json` Version).
2. Probe feature 1 — `mr` / `merge write`: `--help` on both spellings;
   `--dry-run` if supported. NEVER create a real PR.
3. Probe feature 2 — `gitmap help example <cmd>`: run for `mr`, `scan`,
   and 3+ more verbs; verify all three parts render (what-it-does +
   invocation + expected output).
4. Probe feature 3 — four themes: render a help/table command under
   `--theme navy-purple`, `tokyo-night`, `vscode-dark-plus`, `monokai-pro`;
   verify rendering actually changes; pipe output and verify no ANSI
   corruption.
5. Probe feature 4 — `gitmap nc` codegen catalog: `--help` + catalog
   listing; confirm entries are build-time generated, not hardcoded.
6. Per probe record: exact command, exit code, 3-line output digest,
   PASS/FAIL/SKIPPED. Failures recorded as-is — never smoothed over, never
   re-run until passing. Unsafe probes marked `SKIPPED — unsafe` with reason.
7. Write the run log under `.ai-memory/plans/subtasks/` (WS4 evidence file)
   with the version, the probe table, and a "surprises" list.

## Acceptance criteria

- [ ] All four features probed; every probe has command + exit code +
      output digest + verdict.
- [ ] Zero live side effects (no real PRs, pushes, or external mutations).
- [ ] No `go test` / test suites / builds run.
- [ ] Failures recorded honestly with exact output; nothing hidden.
- [ ] `06-e2e-critical-feedback.md` NOT created (reserved for a later wave).

## Out of scope

Fixing found failures (lead decision); writing the critical feedback;
any probe that would cause a side effect.
