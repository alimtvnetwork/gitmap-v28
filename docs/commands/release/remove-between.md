# `gitmap release remove-between` (`rrb`)

Prune intermediate Git tags and GitHub releases between two release bounds while anchoring milestone versions.

<div align="center">

```
  [4.162]  ───────────────── ( intermediate tags purged ) ─────────────────  [4.190]
  Anchor (KEPT)           v4.162.1 .. v4.189.0 deleted            Anchor (KEPT)
```

</div>

---

## Overview

When developing in high-frequency CI/CD environments, repositories accumulate many intermediate pre-release, beta, or test tags. When preparing clean releases, `gitmap release remove-between` (`rrb`) safely prunes the intermediate tag clutter while keeping:
1. The **lower bound** anchor milestone (e.g. `v4.162.0`).
2. The **upper bound** anchor milestone (e.g. `v4.190.0`).
3. Any tag registered in the **pin registry** (`.gitmap/release/pinned-tags.json`) or passed via `--pin`.

---

## Syntax & Aliases

```bash
# Canonical long form
gitmap release remove-between <start> <end> [flags]

# High-speed alias
gitmap rrb <start> <end> [flags]
```

### Argument Formats
- Comma-separated: `gitmap rrb 4.162, 4.190`
- Space-separated: `gitmap rrb 4.162 4.190`
- SemVer prefixed: `gitmap rrb v4.162.0 v4.190.0`

---

## Flags

| Flag | Description |
|---|---|
| `--pin <list>` / `--keep <list>` | Comma-separated list of additional tags to protect |
| `--keep-remote` | Only delete local tags; preserve remote origin tags |
| `--dry-run` | Preview candidate deletions without modifying the repo |
| `-y, --yes` | Skip user confirmation prompt |

---

## Operational Pipeline

1. **AST/SemVer Validation:** Validates `<start>` and `<end>`. Inverted inputs are automatically normalized (`min .. max`).
2. **Anchor & Pin Protection:** Excludes `<start>`, `<end>`, and all pinned tags from deletion.
3. **GitHub Release Deletion:** Deletes corresponding GitHub Release objects via `gh release delete` if configured.
4. **Remote Git Deletion:** Pushes batch deletion requests to `origin`.
5. **Local Git Deletion:** Deletes local tags via `git tag -d`.
6. **Sidecar Cleanup:** Removes `.gitmap/release/<tag>.json` metadata files.
7. **Repository Optimization:** Executes `git fetch --prune --prune-tags origin` and `git gc --prune=now`.

---

## Examples

```bash
# Keep only 4.162 and 4.190, removing 4.162.1..4.189.0
gitmap rrb 4.162, 4.190

# Non-interactive batch removal
gitmap rrb 4.162 4.190 -y

# Preview candidates without modifying repo
gitmap rrb 4.162, 4.190 --dry-run

# Preserve an additional build version
gitmap rrb 4.162, 4.190 --pin 4.175.0 -y
```
