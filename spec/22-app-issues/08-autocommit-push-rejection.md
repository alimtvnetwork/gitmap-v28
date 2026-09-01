# Auto-commit Push Rejection After Release

## Symptom

`gitmap-v28 release` could complete the tag/asset/metadata flow, create the final
`Release vX.Y.Z` commit on the original branch, and then fail on the last push
with a non-fast-forward error such as `main -> main (fetch first)`.

## Root Cause

The post-release auto-commit flow committed locally and then executed a plain
`git push origin <branch>`.

That worked only when the original branch had not changed remotely during the
release pipeline. If another commit landed on the remote branch before the
metadata push, Git rejected the push because the local branch was behind.

There was no recovery path, so the release ended in a partially-finished state:
metadata was committed locally but not pushed remotely.

## Fix

The auto-commit push flow now:

1. Detects non-fast-forward push rejection messages.
2. Runs `git pull --rebase origin <branch>` to replay the metadata commit on top
   of the updated remote branch.
3. Retries `git push origin <branch>` once after a successful rebase.
4. Aborts an in-progress rebase if the sync step fails, avoiding a stuck repo.

## Safety Notes

- No force-push is used.
- Automatic recovery is limited to the non-fast-forward case.
- Other push failures still surface normally with the original Git error text.

## Validation

- Added unit coverage for non-fast-forward detection.
- Added unit coverage to ensure Git stderr/stdout is preserved in returned
  auto-commit errors.
- Manual expectation: when the remote branch advances during release, metadata
  push now self-heals with rebase + retry instead of stopping after commit.
