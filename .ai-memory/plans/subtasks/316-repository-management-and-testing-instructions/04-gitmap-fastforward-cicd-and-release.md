# Subtask Plan 04: gitmap Fast-Forward Logic, CI, and Minor Release

- **Subtask Slug:** `04-gitmap-fastforward-cicd-and-release`
- **Parent Task:** `316-repository-management-and-testing-instructions`
- **Target Repository:** gitmap-v28

---

## 1. Overview & Objectives

1. Verify the fast-forward logic in the pull engine against the governing specs:
   `02-spec/21-app/03-git-operations-and-pull/`,
   `02-spec/21-app/243-pull-forward-slash-paths-pre-pull-and-interactive-remediation-prompt/`,
   `02-spec/21-app/270-gitmap-nodes-pa-fleet-pull-command/`.
   The modern logic: fast-forward only when clean; surface remediation instead of
   silently merging or rebasing published history.
2. Resolve CI failures on main attributable to recent work; prove any remaining
   failure pre-existing with run evidence.
3. Release **v6.600.0** (minor bump from v6.599.0) with tag, using `gitmap release`
   where it fits, so the user can test in the morning.
4. Release notes use relative paths only and attribute the release to gitmap
   (v6.x line — never present as an Antigravity-Manager release).

## 2. Acceptance Criteria

- Fast-forward verdict: conforms to spec, or defects listed with relative paths
  and fixed before release.
- CI green on the release commit, or every red job proven pre-existing and
  unrelated, with run IDs.
- Tag `v6.600.0` on the remote; GitHub release published (not draft).
