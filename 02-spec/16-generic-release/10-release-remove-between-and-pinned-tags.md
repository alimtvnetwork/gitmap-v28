# 10 — Release Tag Removal Between Bounds & Pinned Tags

- **Status:** Approved
- **Domain:** Git Operations & Release Hygiene
- **CLI Commands:**
  - `gitmap release remove-between` (alias: `rrb`, `release-remove-between`)
  - `gitmap release pin-not-to-remove` (alias: `pnr`, `release-pin-not-to-remove`)

---

## 1. Purpose

During high-velocity development and iterative CI/CD testing, repositories often accumulate dozens of intermediate tags (e.g. `v4.162.1`..`v4.189.0`). When stabilizing at a production milestone (e.g. `v4.190.0`), maintainers require a safe, deterministic mechanism to prune obsolete intermediate tags without accidentally removing historical baseline milestones (e.g. `v4.162.0`) or explicitly pinned versions.

This specification defines two complementary commands:
1. **`release remove-between` (`rrb`)**: Evaluates all tags between a lower and upper bound, preserves both endpoint milestones, skips pinned tags, and purges intermediate tags locally, on GitHub, and in local metadata.
2. **`release pin-not-to-remove` (`pnr`)**: Manages a persistent registry of protected tags at `.gitmap/release/pinned-tags.json` that are immune to range deletion and rollback automation.

---

## 2. Command Synopses & Aliases

### A. Range Deletion (`remove-between` / `rrb`)

```bash
gitmap release remove-between <start> <end> [flags]
gitmap rrb <start> <end> [flags]
gitmap release-remove-between <start> <end> [flags]
```

#### Syntax Flexibility:
- Comma-separated or space-separated arguments:
  ```bash
  gitmap rrb 4.162, 4.190
  gitmap rrb 4.162 4.190
  gitmap rrb v4.162.0 v4.190.0
  ```
- Flags:
  - `--pin <v1,v2...>`: Ad-hoc additional tags to protect for this execution.
  - `--keep-remote`: Delete local tags only; leave `origin` intact.
  - `--dry-run`: Preview candidate deletions and protected tags without executing.
  - `-y, --yes`: Auto-confirm deletion prompt.

### B. Pin Registry (`pin-not-to-remove` / `pnr`)

```bash
gitmap release pin-not-to-remove [tags...] [flags]
gitmap pnr [tags...] [flags]
gitmap release-pin-not-to-remove [tags...] [flags]
```

#### Flags:
- `--list`: Display currently pinned tags.
- `--remove <tag>` / `--unpin <tag>`: Remove a tag from the pin registry.
- `--clear`: Clear all pins after confirmation.

---

## 3. SemVer Range Mechanics & Safety Rules

1. **Inclusive Endpoint Immunity:**
   - Both `<start>` and `<end>` versions are strictly preserved.
   - For range `4.162` to `4.190`:
     - Retains `v4.162.0` (or `4.162`)
     - Retains `v4.190.0` (or `4.190`)
     - Candidates evaluated: strictly `start < candidate < end`.
2. **Pinned Version Dominance:**
   - Any tag listed in `.gitmap/release/pinned-tags.json` is unconditionally skipped with a `[PINNED]` indicator.
   - Any tag passed via `--pin` or `--keep` is unconditionally skipped.
3. **Prefix Tolerance & Normalization:**
   - Inputs `4.162`, `v4.162`, `4.162.0`, `v4.162.0` resolve to the same normalized SemVer struct `Major=4, Minor=162, Patch=0`.
4. **Endpoint Ordering Invariant:**
   - If user supplies inverted bounds (e.g. `gitmap rrb 4.190, 4.162`), the system automatically normalizes `lo = min(a, b)` and `hi = max(a, b)` with an informational notice.

---

## 4. Execution Pipeline

When executing `gitmap rrb <start> <end>`:

```text
1. Parse Arguments & Flags
   ├── Parse bounds (start, end)
   ├── Load persistent pins from .gitmap/release/pinned-tags.json
   └── Merge ad-hoc pins from --pin flags
2. Discover Git Tags
   ├── Query local tags: git tag -l
   └── Query remote tags: git ls-remote --tags origin
3. Filter Candidates
   ├── For each tag: parse to SemVer
   ├── Check: start < tag < end
   └── Exclude if tag in PinnedList or tag in {start, end}
4. Preview or Prompt
   ├── If --dry-run: print breakdown & exit 0
   └── If interactive: prompt "Delete N tags between start and end? [y/N]: "
5. Purge Pipeline (Sequential or Batch)
   ├── 5a. GitHub Releases: gh release delete <tag> -y --cleanup-tag=false
   ├── 5b. Remote Git Refs: git push origin --delete <tags...> (in batches)
   ├── 5c. Local Git Tags: git tag -d <tags...>
   └── 5d. Local Metadata: remove .gitmap/release/<tag>.json
6. Optimization & Compaction
   ├── git fetch --prune --prune-tags origin
   └── git gc --prune=now
7. Output Summary Line
   └── "✅ release remove-between complete: N tags deleted, M kept"
```

---

## 5. Persistent Pinned Tags Schema

File path: `.gitmap/release/pinned-tags.json`

```json
{
  "pinnedTags": [
    "v4.162.0",
    "v4.190.0"
  ],
  "updatedAt": "2026-10-12T11:45:00Z"
}
```

---

## 6. Examples

```bash
# Pin critical production milestone releases
gitmap pnr 4.162.0, 4.190.0

# Inspect active pins
gitmap pnr --list

# Preview range removal between 4.162 and 4.190
gitmap rrb 4.162, 4.190 --dry-run

# Execute range removal with auto-confirmation
gitmap rrb 4.162, 4.190 -y

# Ad-hoc protection of a specific build during range purge
gitmap rrb 4.162, 4.190 --pin 4.175.0 -y
```
