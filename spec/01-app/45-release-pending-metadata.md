# Release Pending: Metadata-Based Discovery

## Overview

Enhances `gitmap-v28 release-pending` (alias `rp`) to discover unreleased
versions from **two sources**:

1. **Git branches** (existing behaviour) — local `release/v*` branches
   missing their corresponding `vX.Y.Z` tag.
2. **Release metadata files** (new) — `.gitmap/release/vX.Y.Z.json` files
   where neither the Git tag nor the Git branch exists.

When a `.gitmap/release/vX.Y.Z.json` file contains a `commit` SHA but the
corresponding `release/vX.Y.Z` branch and `vX.Y.Z` tag are both
missing, `release-pending` creates the branch from the stored SHA,
tags it, pushes both, and writes updated metadata.

---

## Discovery Flow

```
1. Scan local Git branches for release/v* patterns   → existing
2. Scan .gitmap/release/v*.json files on disk                → new
3. For each candidate:
   a. Parse version from branch name or JSON filename
   b. Check if vX.Y.Z tag exists (local or remote)
   c. Check if release/vX.Y.Z branch exists
   d. If BOTH missing (metadata-only) → create branch from stored SHA, then tag + push
   e. If branch exists but tag missing → tag + push (existing behaviour)
4. Report summary
```

---

## Metadata-Based Release Steps

When a `.gitmap/release/vX.Y.Z.json` has a valid `commit` field and no
branch/tag exists:

1. Parse the `commit` SHA from the JSON file.
2. Verify the commit exists in the local repository (`git cat-file`).
3. Create branch `release/vX.Y.Z` at the stored SHA.
4. Checkout the branch.
5. Create tag `vX.Y.Z`.
6. Push branch and tag to origin.

**Important:** `release-pending` (and `release-branch`) **skip**
writing `.gitmap/release/` metadata JSON files and committing them. These
commands process already-existing branches or metadata — they only
create the tag, push, and optionally upload assets. The `.gitmap/release/`
JSON writing and committing is exclusive to the primary `release`
command's metadata-first workflow.

If the commit SHA is missing or invalid, skip the version with a
warning and continue processing remaining candidates.

---

## Deduplication

A version discovered from both a Git branch AND a metadata file
should only be processed once. Branch-based discovery takes priority
(existing behaviour). Metadata-based discovery only applies when
**neither** the branch nor the tag exists.

---

## Error Handling

| Condition | Behaviour |
|-----------|-----------|
| `.gitmap/release/` directory missing | Skip metadata scan silently |
| JSON parse error | Warn and skip file |
| Empty or missing `commit` field | Warn and skip version |
| Commit SHA not in local repo | Warn and skip version |
| Push failure | Print error, continue to next |

---

## CLI Output

### New messages

```
  → Found %d unreleased version(s) from .gitmap/release/ metadata
  → Creating release from metadata: %s (commit: %.7s)
  ⚠ Skipping %s: commit %s not found in repository
  ⚠ Skipping %s: no commit SHA in metadata
```

### Combined summary

The existing `Found %d pending release branch(es)` message is updated
to reflect the combined total from both sources.

---

## Flags

No new flags. Existing flags apply to metadata-based releases:

- `--assets` — attach assets to each release
- `--draft` — mark as draft
- `--dry-run` — preview without executing
- `--verbose` — detailed output

---

## Package Layout

### Modified files

| File | Change |
|------|--------|
| `release/workflowbranch.go` | Add `discoverMetadataPending()`, `releaseFromMetadata()` |
| `release/metadata.go` | Add `ListReleaseMetaFiles()` to glob `.gitmap/release/v*.json` |
| `constants/constants_messages.go` | Add pending-metadata messages |

### New functions

| Function | File | Responsibility |
|----------|------|----------------|
| `discoverMetadataPending` | `release/workflowbranch.go` | Glob `.gitmap/release/v*.json`, filter to unreleased |
| `releaseFromMetadata` | `release/workflowbranch.go` | Create branch+tag from stored SHA |
| `ListReleaseMetaFiles` | `release/metadata.go` | Return parsed metadata for all `.gitmap/release/v*.json` |
| `CommitExists` | `release/gitops.go` | Verify SHA exists via `git cat-file` |

---

## Acceptance Criteria

1. `gitmap-v28 rp` with only `.gitmap/release/v1.0.0.json` (no branch/tag) creates
   `release/v1.0.0` branch, `v1.0.0` tag, and pushes both.
2. `gitmap-v28 rp` with both a pending branch AND a metadata-only version
   processes both without duplicates.
3. `gitmap-v28 rp --dry-run` shows metadata-based releases without executing.
4. Missing/invalid commit SHA in metadata prints a warning and continues.
5. `gitmap-v28 rp` with no pending branches and no metadata files prints
   "No pending release branches found."
6. Already-released versions in `.gitmap/release/` are skipped silently.
