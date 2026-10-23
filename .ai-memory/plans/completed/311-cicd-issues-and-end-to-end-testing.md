# Plan 311: CI/CD Issues and End-to-End Testing

## 1. Context & Goal

`main` (gitmap-v28) has been red since at least 2026-10-21: 11 failing checks
at the current baseline (Test (unit), Full Suite Guard, Legacy Refs Check,
Relative Path Check, Lint Baseline Diff, Lint Baseline Guard, Lint Script
Unit Tests, go build + test x3, go test -race). This task restores a fully
green CI/CD, proves every gitmap method end-to-end, improves the installation
download speed, and ships everything as one minor release with complete
release notes. Full text of the user request is captured verbatim in
02-spec/21-app/311-cicd-issues-and-end-to-end-testing/01-architecture-spec.md.

## 2. Invariants & Guardrails

- Pull before work (done: main level with origin at abb35f01d).
- Search exclusively via GitMap (aum search / find / cat / ps / py); total ban
  on rg, ripgrep, grep, git grep, Select-String.
- Relative Git paths only, including the release page and release notes.
- Tests via CI only; no local builds. CI truth comes from GitHub check-runs
  (the locally installed gitmap pe predates the false-green fix and is not a
  valid oracle until this release is installed).
- One atomic commit per task; stage only this task's files.
- Collision control: the spec-310 stream has ~15 uncommitted files in this
  tree (cli/cmddb/, cli/cmdssh/, cli/fsutil/, cli/store/,
  cli/cmdpipeline/pipeline_errorlogs_test.go, cli/cmd/cmd_invocation_test.go,
  cli/release/workflow_test.go, others). Stay out of its dirty files;
  re-inventory with git status before every edit wave. RC heatmap commands
  are the preferred tracker when the gitmap CLI is responsive.

## 3. Subtask Breakdown

| Task-ID | Subtask | Plan file | Status |
|---|---|---|---|
| Task-01 | CI/CD diagnosis to root cause for every failing check at current HEAD | .ai-memory/plans/subtasks/311-cicd-issues-and-end-to-end-testing/01-ci-diagnosis.md | PENDING |
| Task-02 | Fix every CI/CD root cause (dead code, docs, tests, baselines) | .ai-memory/plans/subtasks/311-cicd-issues-and-end-to-end-testing/02-ci-fixes.md | PENDING |
| Task-03 | End-to-end testing for all gitmap methods via repo e2e/smoke harnesses | .ai-memory/plans/subtasks/311-cicd-issues-and-end-to-end-testing/03-e2e-all-methods.md | PENDING |
| Task-04 | Improve installation download speed (profile, optimize, measure) | .ai-memory/plans/subtasks/311-cicd-issues-and-end-to-end-testing/04-install-download-speed.md | PENDING |
| Task-05 | Minor release bump, release notes, final green verification | .ai-memory/plans/subtasks/311-cicd-issues-and-end-to-end-testing/05-release-and-green-verification.md | PENDING |

## 4. Verification Gates

- GitHub check-runs on the release commit: 0 failing checks (all ~51 runs).
- Installer Smoke (linux + windows) green — covers the download-speed change.
- Release notes file present for the new version, full change set, relative
  paths only; tag pushed.
- Download timing before/after recorded in the component spec evidence.
