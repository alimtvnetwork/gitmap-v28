# Specification: GitMap MV, Enhanced RM, Unified Project Resolver, and Replace Engine

## 1. Overview & Objectives

This specification codifies the new `gitmap mv` command, the path-aware `gitmap rm` enhancements, the unified project resolver utility, and the `gitmap replace` diagnostics for GitMap.

## 2. Key Commands & Signatures

### `gitmap mv <source> <destination> [-y|--yes] [--dry-run]`

- **Source**: Slug, alias, relative path (`.\prompt-architect`, `./subfolder`), or absolute path.
- **Destination**: `..` (parent directory), relative path, or absolute path.
- **Features**: Windows long path safety (`\\?\`), atomic database updates, VS Code Project Manager sync, GitHub Desktop sync.

### `gitmap rm <target> [-y|--yes] [--db-only]`

- **Target**: Resolves paths (`.\folder`, `./folder/`, `folder`), aliases, slugs, or PWD `.`.
- **Features**: Cleans SQLite DB, removes folder from disk, cascades removal to VS Code PM and GitHub Desktop.

### `gitmap replace <old> <new> [-y|--yes] [--dry-run]`

- **Features**: Zero false-skips, normalized Windows path traversal, accurate multi-file token replacement with visual diffs.

## 3. Unified Project Resolver Utility

- Consolidated in `gitmap/cmd/resolver.go` and `gitmap/fsutil/path_normalize.go`.
- Maps any target string to an authoritative `model.ScanRecord`.
