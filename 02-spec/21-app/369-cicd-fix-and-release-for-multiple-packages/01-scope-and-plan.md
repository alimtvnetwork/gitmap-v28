# 369 - CI/CD fix and release for multiple packages

Date: 31-Oct-26
Repo: gitmap-v28 (Go CLI, gitmap v6.x line only - never the Antigravity-Manager v4.x line)
Base at start: main 58292e98, Release v6.643.0 (pulled via gitmap pull from 49142cfe, Release v6.642.0, +1164/-175, 37 files).
Base moved during the run: another agent released v6.644.0 (main
d88c139c, heatmap file auto-read auto-touch and high-speed modtime
fleet sync) before any fix was committed, so all fixes below are
rebased onto d88c139c and the v6.644.0 CI runs (CI 41199726604,
Cross-Platform 41199726630, both completed failure) become the
operative before-state. Their failed jobs are the same set as
v6.643.0 minus Lint Baseline Diff (which passed on v6.644.0):
Lint Script Unit Tests, File Size Gate, Boolean and Enum Linter,
Relative Path Check, Boolean Guidelines Linter, Full Suite Guard,
Legacy Refs Check, Nested If Linter, Test (unit), plus all three
Cross-Platform go build + test jobs.
Task type: autopilot CI/CD fix, minor version bump, release. Solo execution:
the parent task explicitly forbids spawning subagents, so the v6
multi-agent dispatch gate is deviated from by parent instruction, not
by choice. That deviation is recorded here, not hidden.

## 1. Baseline CI state (before, recorded 31-Oct-26)

Live runs on the v6.643.0 tip, workflow CI and Cross-Platform Build,
were still in progress when first read, with these jobs already failed:

- CI run 41198431481 (https://github.com/alimtvnetwork/gitmap-v28/actions/runs/41198431481):
  Lint Baseline Diff, Lint Script Unit Tests, Legacy Refs Check,
  File Size Gate (300-line cap), Boolean and Enum Linter,
  Relative Path Check, Boolean Guidelines Linter, Nested If Linter.
- Cross-Platform Build run 41198431561 (https://github.com/alimtvnetwork/gitmap-v28/actions/runs/41198431561):
  ubuntu-latest go build + test failed at `go test ./...`
  (build itself passed); macos and windows jobs still running.

The previous tip (Release v6.642.0) was also red, which shows the
failures are not a one-off flake on the new tip:

- CI run 41196146952: Lint Script Unit Tests, Relative Path Check,
  Legacy Refs Check, Full Suite Guard, Test (unit) failed.
- Cross-Platform Build run 41196146948: ubuntu, windows and macos
  go build + test all failed.
- Known failure details from the completed v6.642.0 run:
  test_ci_scripts.py failed test_gofmt_check_clean_repo and
  test_legacy_refs_runs_clean_on_repo (both returncode 1, not 0);
  unit tests failed in cli/release and in
  TestShouldPrintInstallHint_NonGitmapRepos/Partial_prefix.

## 2. Scope

In scope:

- Reproduce each failing job locally where possible (gofmt, go test,
  the policy-check linters under .github/actions/policy-check, and
  test_ci_scripts.py) and fix only what CI actually fails on.
- Fix source, test, or linter-script files outside .github/workflows/.
- Minor version bump via the repo's own mechanism, then a release
  tag via gitmap release, verified on the remote with git ls-remote
  and by a green release workflow run with assets present.
- Final verification: working tree clean, local == origin on main,
  latest CI run on the release tip green.

Out of scope (deferred, reported honestly, never stalled on):

- Pre-existing linter findings that CI does not fail on.
- Gold-plating, refactors unrelated to a failing job, prompt/skill
  changes (those live in coding-guidelines only, and none are needed).
- Any change under .github/workflows/: the VM gh token lacks the
  workflow scope, so GitHub rejects such pushes. If a CI fix needs a
  workflow change, it is committed locally only and reported as
  BLOCKED with its exact SHA and the unblock step (user runs
  `gh auth refresh -s workflow` interactively, then push).

## 3. Rules applied

- Search exclusively via gitmap (aum search, find, cat, ps, py,
  llm train). Total ban on rg, ripgrep, grep, git grep, Select-String.
- Strictly relative git paths everywhere (02-spec/..., .ai-memory/...,
  cli/...). No absolute path in code, docs, release page, or notes.
- Comment all changes briefly with WHY. Commit all modified files
  together and push immediately via the gitmap commit suite
  (hyphen-format messages, no colons in the message argument).
  Never force-push, never rewrite history under existing tags,
  never delete repos or files unless explicitly asked.
- Git identity: pre-commit check passed at start - local and global
  are `gitmap automation bot - muse - runs autonomously
  <info@rgitmap.com>`; never author as Muse <muse@local>, never
  invent an identity.
- Doubtful code gets a strict DOUBT comment at that spot plus a
  gitmap rc heatmap log entry with file, line, and what needs
  verification.
- Dates DD-MMM-YY. Code time UTC+0, saved as UTC. UI time is the
  machine-local timezone, never the user's. Branch names date only.
- Headless VM: Go 1.26.0 at the toolchain under the workspace tools
  directory (no Go on default PATH); never build under /tmp
  (512M tmpfs) - use workspace tool caches. The
  `git: 'credential-manager' is not a git command` push warning is
  benign and ignored.

## 4. Method

1. Pull first (done: gitmap pull to 58292e98, clean tree).
2. Spec + plan (this folder, the pending plan, and subtasks).
3. Diagnose: local gofmt, go test for the failing packages, and
   the policy-check linters, in that order - fix what is actually
   red, smallest diff that makes the job green.
4. Bump minor version, commit + push, watch CI on the new tip.
5. Release tag via gitmap release, verify tag via git ls-remote,
   verify the release workflow run is green and the release is
   published (not draft, not prerelease) with assets.
6. Final verification and report: CI before -> after with run
   IDs/URLs, every fix with commit SHAs, version old -> new,
   tag + release URL + asset count, tree/origin status, BLOCKED
   section if anything is held locally.

## 5. Risks

- Another agent may push to main concurrently (as happened between
  v6.642.0 and v6.643.0). Mitigation: pull again before the release
  commit and before tagging; merge, never force-push.
- The in-progress v6.643.0 CI runs may finish with more failing jobs
  than the eight already seen. The fix list must be re-checked
  against the final run state before the release commit.
- The workflow-scope token blocker may hold a workflow fix locally;
  that is a reportable blocker, not a reason to bypass GitHub.
