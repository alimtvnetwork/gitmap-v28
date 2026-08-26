# Release Architecture Map

## Overview
Releases in this repository are managed primarily via Git tags (`vMAJOR.MINOR.PATCH`) pushing to GitHub, which trigger Goreleaser via GitHub Actions (`.github/workflows/goreleaser.yml` or `release.yml`). 
However, the single source of truth for the *current* application version is **`version.json`** located at the repository root.

## Source of Truth
1. **`version.json`**: This file contains the authoritative version string (e.g., `{"version": "6.117.0"}`). It MUST be updated manually or via script during every release.
2. **`changelog.md`**: Contains release notes mapped to the specific version. Must be updated with a matching header `## [vX.Y.Z] - YYYY-MM-DD`.

## The Release Process (Agent Guidelines)
When instructed to "bump the MINOR version" and trigger a release at the end of a tunnel:
1. **Update `version.json`**: Increment the MINOR part and reset PATCH to `0` (e.g., `6.117.0` -> `6.118.0`).
2. **Update Badges / Links**: Run a find-and-replace to update hardcoded versions in:
   - `readme.md`
   - `.lovable/what-to-read.md`
   (e.g., replacing `gitmap-v27` with `gitmap-v28` or pinned versions).
   **CRITICAL RCA WARNING (Stale Version Contamination):**
   Do not blindly trust that the repository is perfectly synced to the immediate previous version. Sloppy past releases may have left severely outdated versions (e.g., `v6.111.0` when you expect `v6.120.0`) in the `readme.md` version matrix or install snippets. 
   - You MUST scan `readme.md` using a regex like `[vV]?6\.\d+\.\d+` to hunt down ALL old versions and obliterate them in favor of the new version.
   - You MUST run a post-rewrite audit to prove zero stale versions remain.

3. **Update `changelog.md`**: Append the latest fixes to a new `## [vX.Y.Z]` block.
4. **DO NOT CREATE GIT TAGS VIA CLI**: The user explicitly stated in previous instructions: "you should not create the tag... I will release it, uh, using the Git map." You should only prepare the commit. Wait for the user to trigger the actual tag/release process unless explicitly asked to run `git tag`.
5. **DO NOT modify test files**: Never modify `*_test.*` or `*.spec.*` during version scanning.

## GitHub Actions Triggers
- Commit builds: `ci.yml` triggers on pushes to branches matching `**`, ignoring tags `v*`.
- Releases: `release.yml` triggers on pushes to tags matching `v*`.
