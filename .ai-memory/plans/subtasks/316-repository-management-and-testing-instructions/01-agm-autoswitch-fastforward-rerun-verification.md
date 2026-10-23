# Subtask Plan 01: AGM Auto-Switch, Fast-Forward, and Single-Prompt Rerun Verification

- **Subtask Slug:** `01-agm-autoswitch-fastforward-rerun-verification`
- **Parent Task:** `316-repository-management-and-testing-instructions`
- **Target Repository:** Antigravity-Manager (research read-only; fixes by lead only)

---

## 1. Overview & Objectives

Verify, with code-level evidence, that Antigravity-Manager:

1. Auto-switches accounts/instances under the discussed trigger conditions.
2. Fast-forwards following the modern logic (clean fast-forward only; no silent
   merge/rebase of published history).
3. After a switch or a fast-forward, reruns exactly **one** running prompt
   (restore/rerun semantics from the v4.190.0 fix must not regress: unscoped restore
   is FIFO, explicit instance scope is strict, `agm rerun` takes the most recent N
   and actually re-sends).

## 2. Method

- Search exclusively via GitMap (`gitmap aum search`, `gitmap find`, `gitmap cat`)
  inside the Antigravity-Manager repo. Total ban on rg/grep/git grep.
- Headless verification via the `agm` CLI where the logic is reachable without the
  Tauri GUI / Antigravity IDE (this VM has no display server).
- Read the governing specs under `02-spec/` in Antigravity-Manager and compare the
  implementation against them.

## 3. Acceptance Criteria

- Auto-switch trigger conditions and selection logic documented with relative paths.
- Fast-forward call path traced end to end; divergence behavior stated.
- Post-switch / post-fast-forward rerun path traced; proof it dispatches a single
  prompt, or a proven defect report.
- Gap list with severity, ready for the fix phase.
