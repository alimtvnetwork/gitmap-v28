# Completed: 244-e2e-verify-program-243-changes

Status: COMPLETED
Date completed: 2026-10-09 (Asia/Kuala_Lumpur)
Slug: 244-e2e-verify-program-243-changes
Spec: `02-spec/21-app/244-e2e-verify-program-243-changes/`
Skill: execute-parent-task-with-n-steps-v6 (A=2, H=2)
Consolidates: `.ai-memory/plans/244-e2e-verify-program-243-changes.md` (deleted after merge)
and `.ai-memory/plans/subtasks/244-e2e-verify-program-243-changes/` (01, 02 — deleted after merge)

## Objective

End-to-end verify program 243's COMMITTED changes (main through `f127f7a`):
`space backup-branch`, help-displayer core + DRY migration (scan/pull-all/space),
unknown-command suggestion engine, enforcement + regression. Build a binary from
pristine HEAD, run 27 scripted assertions, report PASS/FAIL with evidence and a verdict.

## Execution log

- **Phase 1A** — Task breakdown shown in chat (`Understood: [YES]` × 6). Preflight:
  `gitmap pa` pulled all 5 registered repos (7.3s; Antigravity-Manager `8a69451..3658199`,
  rest up-to-date). Task DB initialized at
  `.ai-memory/temp-agents/01-244-e2e-verify-program-243-changes/` (subtask populate
  skipped: installed `gitmap task add` only accepts `-db`/`-parent-id`, not the
  skill-documented `--code/--title/--files/--role` — logged, non-blocking).
- **Phase 1 research** — Research 01 mapped the `space backup-branch` CLI contract
  (flags `--no-push`/`--force`/`--help`, no `--dry-run`, slug algorithm, error paths);
  Research 02 mapped the help displayer (`cli/helpdisplay/`, migrated builders for
  scan/pull-all/space), the suggestion engine, and wave-D enforcement state.
- **Key pivot (verified)** — The working tree is mid-243b-consolidation (~629
  uncommitted changes, another coordinator actively working) and does NOT compile
  (import cycles, e.g. `cli/cmdspace/space.go` self-import; confirmed by lead
  repro). Per the out-of-scope boundary, e2e targeted the COMMITTED state:
  pristine scratch clone `~/workspace/gitmap-e2e-244/src` at `f127f7a`
  (0 dirty), `go build ./...` from `cli/` exit 0, binary
  `~/workspace/gitmap-e2e-244/gitmap` (`gitmap v6.515.0`,
  sha256 `6caafbdecdc46889d20c01f4dc47ad3617c337859d2fe11bf6618720ec222f10`).
- **Spec** — Spec Agent 1: `01-test-scope-and-plan.md` (groups A–E, harness design)
  + subtask `01-author-e2e-scripts.md`. Spec Agent 2:
  `02-acceptance-criteria-and-evidence.md` (IDs A01–A10, B01–B07, C01–C04, D01–D05, E01)
  + subtask `02-run-and-verify.md`. Lead corrections before the run: D02/D05
  `300`→`500` (committed reality at `f127f7a`; verified against pristine source —
  Research 02 had read the dirty 243b tree), C04 wording → `unrecognized command`,
  B04 suggestion check → WARN-only, Group A/B tables aligned to the implemented scripts.
- **Phase 2 wave 1** — Worker 01 authored 6 scripts under `~/workspace/gitmap-e2e-244/`
  (outside the repo, never committed): `common.sh`, `e2e-backup-branch.sh`,
  `e2e-help-displayer.sh`, `e2e-suggestions.sh`, `e2e-enforcement-regression.sh`,
  `run-all.sh` — 27 assertions, all `bash -n` clean, matchers grounded against live
  binary runs. Worker 02 (parallel): binary smoke tests (3/3 pass), `env.json`,
  `verification-report.md` skeleton, `work/` scratch parent. Lead relayed Worker 02's
  ANSI-color caveat to Worker 01 mid-wave (scripts strip ANSI before matching).
- **Phase 2 wave 2** — Worker 02 ran `run-all.sh` (suite log:
  `~/workspace/gitmap-e2e-244/run-all.log`, two runs, identical):
  **RESULT: 26 passed, 1 failed** (1 informational warning). Report:
  `~/workspace/gitmap-e2e-244/verification-report.md`.

## Verdict

`NOT VERIFIED — see RCA notes` — 1 failing assertion: **A08**.

## The A08 defect (RCA verified by lead repro)

`gitmap space backup-branch`'s **refusal path dirties the repo it runs in**, defeating
the command's own clean-tree guard on the next invocation:

1. `space backup-branch --no-push "dup case"` with pre-existing `backup/dup-case`
   → correctly refuses (`branch backup/dup-case already exists …`, exit 1) — but the
   dispatch error handler calls `writeLastErrorFile` (`cli/cmd/root.go`), which does
   `os.MkdirAll(".gitmap", 0755)` and writes `.gitmap/last_error.log` **relative to
   the user's CWD**. The scratch repo now has untracked `.gitmap/`.
2. `space backup-branch --no-push --force "dup case"` → the clean-tree guard
   (`git status --porcelain` non-empty) refuses with `Working tree has uncommitted
   changes…` (exit 1) **before the `--force` recreate logic is reached**.

Lead-verified: the `--force` implementation itself is correct (clean tree → exit 0,
branch recreated at HEAD); the defect is purely the refusal-path side effect.
Blast radius is wider than A08: **any** failing `space backup-branch` invocation
writes `.gitmap/` into the target repo (lead-verified for both the branch-exists
refusal and the missing-task refusal). The `scna` unknown-command path does not
write it (different error path). So every failed `backup-branch` run leaves the
user's repo dirty — and any subsequent `backup-branch` run in that repo will be
refused by the clean-tree guard until `.gitmap/` is removed/ignored.

## Other findings (informational)

- **500-vs-300**: at committed `f127f7a` the file-size gate enforces 500
  (`52-file-size-check.py` default; `50-fastgate.py` docstring/message agree). The
  300-line cap is docs-only at this commit. The 243b working tree already flips the
  code default to 300 — that alignment is 243b's scope, NOT a 243 failure.
- **B04**: `space --help common` renders groups correctly but carries no trailing
  `* ` suggestion line (WARN-only).
- **C04**: the unknown-command message is `unrecognized command` (E1001), not
  "Unknown command".
- **A10**: the real program-243 backup branch `backup/2026-10-08-pre-help-displayer`
  exists on origin — the command worked in production.

## Evidence (all outside the repo, never committed)

- Scripts: `~/workspace/gitmap-e2e-244/{common.sh,e2e-backup-branch.sh,e2e-help-displayer.sh,e2e-suggestions.sh,e2e-enforcement-regression.sh,run-all.sh}`
- Report: `~/workspace/gitmap-e2e-244/verification-report.md` · Log: `~/workspace/gitmap-e2e-244/run-all.log`
- Binary: `~/workspace/gitmap-e2e-244/gitmap` (from pristine `f127f7a`, sha256 above)
- Pristine source: `~/workspace/gitmap-e2e-244/src`

## Follow-ups (NOT fixed in this run — out of scope by design)

1. **A08 defect** — `writeLastErrorFile` (`cli/cmd/root.go`) writing `.gitmap/`
   into the user's repo on failure paths. Fix options: (a) write `last_error.log`
   outside the repo (user-global dir, e.g. alongside the errors DB); (b) make the
   clean-tree guard ignore `.gitmap/`; (c) skip the error-file write for validation
   errors. Also verify `persistToErrorsDB`'s write location for the same class of
   side effect. Do NOT fix inside the active 243b consolidation without coordinating
   with its coordinator (it is rewriting `cli/cmdspace/` right now).
2. **300-line cap alignment** — owned by the in-flight 243b consolidation.
3. **Spec index hygiene** — spec 243 (`02-spec/21-app/243-cli-help-displayer-and-backup-branch/`)
   is not registered in `02-spec/21-app/readme.md`; its coordinator should register it.

## Assumptions / conflicts log

- "space PA" read as `gitmap pa` (pull-all). E2E non-destructive: no real backup
  branches pushed; disposable scratch repos only. 243b uncommitted work out of scope.
- Skill R1 (zero builds/tests) vs the user's explicit order to write and run e2e
  scripts: user wins per the Top-Instruction Priority Mandate; scoped to this task.
- Skill R8/R9 vs dirty tree: `gitmap cpf` runs `git add -A`, which would have swept
  the 243b coordinator's ~629 uncommitted changes into this task's commit, violating
  R8's own "strictly this task's files". Committed instead via targeted
  `git add <paths>` + `git commit -m "Feature: …"` + `git push` — one atomic commit,
  `Feature: ` prefix preserved, only this task's 5 files. Logged here as the
  reasoned deviation.
- Subtask consolidation followed the skill (merge → delete); content preserved
  above and in the canonical spec; deletions are this task's own transient files.

## Files changed by this run (the atomic commit)

- `02-spec/21-app/244-e2e-verify-program-243-changes/01-test-scope-and-plan.md` (new)
- `02-spec/21-app/244-e2e-verify-program-243-changes/02-acceptance-criteria-and-evidence.md` (new)
- `.ai-memory/plans/completed/244-e2e-verify-program-243-changes.md` (new — this file)
- `.ai-memory/plans/readme.md` (entry `22-e2e-verify-program-243-changes.md`)
- `02-spec/21-app/readme.md` (entry `244-e2e-verify-program-243-changes`, completed)
