# Issue 11 — Automatic Legacy Directory Migration

## Summary

Legacy repo-local directories (`gitmap-output/`, `.release/`, `.deployed/`) are
automatically migrated to `.gitmap/` subdirectories when detected.

## Key Follow-up Fix

Release workflows can temporarily restore tracked legacy `.release/` files when
checking out the original branch after tagging from a release branch. To prevent
old folders from persisting after `gitmap-v28 release`, migration now runs in two
places:

1. On CLI startup (except `version`)
2. Again after the release workflow returns to the original branch

This guarantees legacy folders are merged into `.gitmap/` and removed before
release auto-commit inspects the working tree.

## Migration Map

| Legacy Directory   | Target              |
|--------------------|---------------------|
| `gitmap-output/`   | `.gitmap/output/`   |
| `.release/`        | `.gitmap/release/`  |
| `.deployed/`       | `.gitmap/deployed/` |

## Rules

1. Detect legacy directory at working directory root.
2. Create `.gitmap/` if missing.
3. If target does **not** exist → rename (move) legacy directory to target.
4. If target **already exists** → merge files from legacy into target
   (skip files that already exist), then **remove the legacy directory**.
5. Re-run migration after `release` returns to the original branch.
6. Database (`data/`) is **not affected**.

## Status

Complete — startup migration plus post-release cleanup in place.
