# 369 - Acceptance criteria and evidence

Date: 31-Oct-26

## Acceptance criteria

1. The latest CI workflow run on the release tip is green, including
   every linter job that was red on the v6.643.0 tip (Lint Baseline
   Diff, Lint Script Unit Tests, Legacy Refs Check, File Size Gate,
   Boolean and Enum Linter, Relative Path Check, Boolean Guidelines
   Linter, Nested If Linter) and the Cross-Platform Build jobs.
2. If any failing job cannot be fixed without a workflow-file change,
   that fix is committed locally, its SHA is named in the final
   report under BLOCKED, and no token trick, scope bypass, or
   force-push is attempted.
3. Version bump is a minor bump from the base version via the repo's
   own mechanism, and the change is visible in version.json.
4. The release tag exists on the remote (proven by git ls-remote),
   the release workflow run is green, and the release is published
   (not draft, not prerelease) with assets present if the release
   workflow produces them.
5. Working tree clean and local == origin on main at the end.
6. No absolute path, no invented git identity, no rewritten tag
   history, no deleted file introduced by this run.

## Evidence to attach to the final report

- CI run IDs and URLs before (41198431481, 41198431561) and after.
- Commit SHA for every fix batch, the bump commit, and the release
  commit/tag.
- Release URL and asset count.
- git status and git ls-remote output at the end.
- BLOCKED section (empty if nothing is held locally, otherwise the
  exact SHA, the file, why it is held, and the exact unblock step).
