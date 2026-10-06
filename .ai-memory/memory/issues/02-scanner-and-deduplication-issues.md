# Issue Domain 02: Scanner and Deduplication Issues

- **Domain:** Repository Scanning, Collation, and Path Collisions
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Case Collisions on Windows NTFS
- **Symptoms:** Scanning reported duplicate repository entries with different casing (e.g., `RiseUp-Asia` vs `riseup-asia`).
- **Root Cause:** NTFS preserves case but is case-insensitive; Go map keys performed exact string comparisons.
- **Resolution:** Normalized all scan paths via `filepath.ToSlash` and lowercase comparison (`strings.EqualFold`).

## 2. Infinite Symlink Recursion
- **Symptoms:** File walker hung when traversing circular symbolic links in node environments.
- **Root Cause:** `filepath.Walk` followed symlinks without visited inode/path tracking.
- **Resolution:** Added visited directory set with cycle detection and depth caps.
