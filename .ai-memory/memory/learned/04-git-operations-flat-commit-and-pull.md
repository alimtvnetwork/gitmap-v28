# 04 — Git Operations, Flat Commit, and Pull Architecture

- **Subsystem:** Git Automation & Status Synchronization
- **Status:** Authoritative Reference

## 1. Semantic Flat Commit Suite (`gitmap commit` / `cm`)
- Single-command atomic commit creation with automatic dirty detection.
- Fast alias `gitmap cm` executes with sub-50ms overhead.
- Auto-stage flag (`-a` / `--all`) stages tracked and untracked changes before committing.
- Pre-flight guard blocks execution in non-git directories with actionable suggestions.

## 2. Pull-All Status (PAS) Formula
- Calibrates eligible repositories against current local/remote branch divergence.
- Worker concurrency scales to `runtime.NumCPU() * 2`.
- Fast-forward auto-rebase eliminates merge commit noise across clean branches.

## 3. Commit-Pull-Array-Runner (CPAR)
- Orchestrates multi-repo commit and pull pipelines with fail-fast recovery.
