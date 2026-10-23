# Subtask Plan 02: AGM CI/CD Finalization and Minor Release

- **Subtask Slug:** `02-agm-cicd-and-minor-release`
- **Parent Task:** `316-repository-management-and-testing-instructions`
- **Target Repository:** Antigravity-Manager

---

## 1. Overview & Objectives

1. Finalize CI/CD: every workflow on main green at release time. Distinguish
   regressions (fix) from pre-existing failures (report with evidence). Note the
   known state from 2026-10-24: the latest completed main run before this program
   had failures and several runs were cancelled; PR #6 "Sync v5" is conflicting and
   stale — out of scope unless it blocks CI.
2. Release **v4.254.0** (minor bump from v4.253.0) with tags, using the `gitmap
   release` command where it fits.
3. Release notes use relative paths only and attribute the release to
   Antigravity-Manager (v4.x line — never present as a gitmap release).

## 2. Acceptance Criteria

- CI run IDs on main recorded; zero unexplained failures at the release commit.
- Tag `v4.254.0` exists on the remote; GitHub release published (not draft).
- Any post-v4.253.0 fixes from subtask 01 included and listed in the notes.
