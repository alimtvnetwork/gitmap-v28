# Specification: 311-cicd-issues-and-end-to-end-testing

## User Request (Verbatim)

```text
# CI/CD Issues and End-to-End Testing: high priority instruction, non-negotiable task

Can you please check for the CI/CD issues? If everything is all right, do the end-to-end testing for all the methods in GitLab and make sure there is a release bump with a minor bump, and all the release notes are there in the last release, and CI/CD is green.

## task: CI/CD Issues and End-to-End Testing
## slug: cicd-issues-and-end-to-end-testing

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-cicd-issues-and-end-to-end-testing/ and enqueue plan task in .ai-memory/plans/pending/XX-cicd-issues-and-end-to-end-testing.md (subtasks in .ai-memory/plans/subtasks/XX-cicd-issues-and-end-to-end-testing/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Check for CI/CD issues and ensure everything is functioning correctly.
6. Perform end-to-end testing for all methods in GitLab.
7. Ensure there is a release bump with a minor version increment.
8. Verify that all release notes are included in the last release.
9. Confirm that CI/CD status is green.

## Must follow and spawn agent using

[execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster execution for AI and coding without errors and following Coding Guidelines.
- [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Track in-flight file touches and prevent write collisions using GitMap RC Heatmap (`gitmap rc heatmap status` / `gitmap rc heatmap tree` / `gitmap rc heatmap live` / `gitmap rc heatmap sync`).
- Only add relative paths, never add absolute paths during work; this must be respected on release page and in release notes as well.

Follow-up (same session): "Make sure you improve the installation download speed"
```

## 1. Objective

Restore gitmap-v28 `main` to a fully green CI/CD state, prove the full gitmap
method surface end-to-end, measurably improve the installation download speed,
and ship all of it as one minor release with complete release notes.

## 2. Baseline Evidence (verified 2026-10-23)

- `main` has been red since at least 2026-10-21. At `f723012c5` the GitHub
  check-runs show 51 runs with exactly 11 failures, unchanged from the
  pre-existing baseline at `1f96b906`:
  Test (unit), Full Suite Guard, Legacy Refs Check, Relative Path Check,
  Lint Baseline Diff, Lint Baseline Guard, Lint Script Unit Tests,
  go build + test (ubuntu/macos/windows), go test -race (hot packages).
- The `gitmap pe` false-green defect (empty data fabricated into a CLEAN
  verdict when the gh CLI is unauthenticated) was root-caused and fixed in
  `db665c939` + `f723012c5`; that fix ships in this release, so `pe` becomes
  trustworthy only after this release is installed.
- Known concrete findings behind the lint failures (from check-run
  annotations): unused functions `recordVacuumTelemetry`, `handleVacuumError`,
  `executeVacuumOnDB`, `vacuumSingleTarget` (cli/cmddb/db_vacuum.go),
  `requireOnlineCheck` (cli/cmd/helpcheck.go), `hasCredentialPersistenceFailure`
  (cli/cloner/safe_pull.go), `formatFolderFileURI`, `hasSkipIDEsFlag`,
  `parseGitHubSlugFromURL`; misspell `DOUB` -> `DOUBT`; Legacy Refs hits for
  /gitmap-v[567]\b/ in 02-spec/01-spec-authoring-guide/13-root-readme-conventions.md
  (lines 20 and 257).
- A second worker stream (spec 310, strict relative paths) is in flight in the
  same working tree with ~15 uncommitted files (cli/cmddb/, cli/cmdssh/,
  cli/fsutil/, cli/store/, cli/cmdpipeline/pipeline_errorlogs_test.go,
  cli/release/workflow_test.go, others). Its HEAD commit `abb35f01d` already
  targets the Relative Path / Legacy Refs area. This task's file boxes must
  stay disjoint from that stream's dirty files; re-inventory before editing.

## 3. Workstreams

- **A. CI/CD root-cause fixes** — every failing check traced to its root cause
  and fixed at the source (dead code removed, docs corrected, tests repaired,
  baselines updated only where the baseline itself is the defect). No
  suppression of a check to reach green.
- **B. End-to-end testing of all methods** — every command/method in the
  gitmap registry exercised end-to-end through the repo's own e2e and smoke
  harnesses (CI Installer Smoke, history-rewrite smoke, command invocation
  coverage), with the full matrix green in CI.
- **C. Installation download speed** — profile the install/update download
  path (installer scripts + cli/cmdupdate/), remove its bottlenecks
  (single-connection sequential download, redundant round-trips, no resume),
  building on the direction of specs 242 (aria2c multi-connection download)
  and 286 (download speedup). Measure before/after download time.
- **D. Release** — minor version bump, tag, and release notes containing the
  complete change set, written with relative paths only.

## 4. Non-Negotiable Constraints

- Pull before work (done: main level with origin at `abb35f01d`).
- Codebase search exclusively via GitMap (aum search / find / cat / ps / py).
  Total ban on rg, ripgrep, grep, git grep, Select-String.
- Relative Git paths only, everywhere, including the release page and release
  notes. No absolute paths in any produced file.
- Tests via CI only; no local builds. CI state is verified through GitHub
  check-runs, never through the pre-fix `gitmap pe` on this machine.
- One atomic commit per task; stage only this task's files (other streams are
  active in the same tree).
- Collision control: RC heatmap when the CLI is responsive; otherwise the
  dirty-file inventory from git status is the coordination source of truth.

## 5. Acceptance Criteria

1. GitHub check-runs for the release commit: zero failing checks.
2. E2E coverage: every registered gitmap method exercised by the repo's
   e2e/smoke harnesses in CI, all passing.
3. Installation download: timed before/after measurement recorded in
   02-component-spec.md evidence; after is faster than before.
4. version.json (and every synced version file) bumped by exactly one minor
   version; annotated tag pushed; release notes file lists the full change
   set with relative paths only.
5. Spec + plan moved to completed state with evidence links.
