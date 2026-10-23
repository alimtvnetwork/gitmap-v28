# Subtask 311-05: Release and Green Verification

## Objective

Ship subtasks 311-01..04 (plus the already-merged pe false-green fix) as one
minor release and prove CI/CD is fully green on the release commit.

## Steps

1. Read the current version from version.json (v6.590.0 at plan time) and
   bump the minor component by one in version.json and every synced version
   file (locate all of them via gitmap find / aum search first).
2. Write the release notes file for the new version covering the complete
   change set: CI/CD root-cause fixes, e2e method coverage, installation
   download-speed improvement, and the gitmap pe false-green fix
   (db665c939, f723012c5). Relative paths only — no absolute paths, no
   filesystem URLs, including on the release page text.
3. One atomic commit for the whole task (stage only this task's files),
   then push the annotated release tag.
4. Verify via GitHub check-runs on the release commit: all runs completed,
   zero failures. Do not use the locally installed gitmap pe as the oracle
   (it predates the false-green fix).
5. Move the parent plan to .ai-memory/plans/completed/ with the evidence
   filled in, and update the spec evidence section.

## Exit Criteria

- Tag pushed for the new minor version; release notes present and complete.
- Check-runs for the release commit: 0 failing checks.
- Parent plan closed with evidence.
