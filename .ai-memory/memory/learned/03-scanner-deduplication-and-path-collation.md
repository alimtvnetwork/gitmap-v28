# 03 — Scanner Deduplication and Path Collation

- **Subsystem:** Repository Discovery & Filesystem Ingestion
- **Status:** Authoritative Reference

## 1. High-Speed Repository Discovery
- Scans multi-tier directory trees across local filesystems and mounted network shares.
- Employs parallel directory walker pools with channel-based queuing.
- Detection heuristics identify git repositories via `.git` directories and worktrees.

## 2. Case-Insensitive Path Collation
- Cross-platform path normalization utilizes `filepath.ToSlash` and lowercase collation (`strings.EqualFold`) to prevent duplicate tracking on case-insensitive filesystems (Windows NTFS, macOS APFS).
- Deduplicates repository candidates across bookmarks, aliases, and scan output.

## 3. Ignore Engine (`.gitmapignore`)
- Hierarchical ignore rules bypass heavy dependencies (`node_modules`, `vendor`, `.venv`, `target`).
