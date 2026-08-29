# Specification 17: Move Command, Robust Path Resolver, Removal Sync & Replace Engine

## 1. Executive Summary

This specification defines the architectural design, algorithmic contracts, and CLI interfaces for four tightly integrated systems in GitMap:
1. **`gitmap mv` (Move)**: Safe physical filesystem relocation of repositories (supporting `..`, relative paths, absolute paths, and Windows long paths `\\?\`), atomic SQLite database record re-pointing, and synchronization with VS Code Project Manager and GitHub Desktop.
2. **Unified Project Resolver**: A single, DRY, robust repository resolution engine that seamlessly maps slugs (`prompt-architect`), relative paths (`.\prompt-architect`, `./prompt-architect/`, `../repo`), aliases (`Alias` table), normalized absolute paths, and working directories (`PWD`) into authoritative `Repo` records.
3. **Enhanced `gitmap rm` (Remove & Delete)**: Path-aware repository untracking and physical deletion with interactive prompts (or `-y` / `--yes` bypass), cascading removals across SQLite, VS Code Project Manager (`projects.json`), and GitHub Desktop.
4. **Enhanced `gitmap replace` Engine**: Robust literal and versioned token substitution with proper path normalization, Windows forward/backward slash reconciliation, binary sniffing protections, and multi-file transactional safety.

## 2. Core Pillars & Architecture

```
                                  ┌────────────────────────────────────────┐
                                  │       Unified Project Resolver         │
                                  │  (Path/Slug/Alias/PWD/Glob Normalizer) │
                                  └───────────────────┬────────────────────┘
                                                      │
                       ┌──────────────────────────────┼──────────────────────────────┐
                       ▼                              ▼                              ▼
          ┌──────────────────────────┐   ┌──────────────────────────┐   ┌──────────────────────────┐
          │      gitmap mv           │   │      gitmap rm / del     │   │      gitmap replace      │
          │  - Path relocation (..)  │   │  - Path-aware targeting  │   │  - Exact & regex token   │
          │  - Windows long-path     │   │  - SQLite cascade wipe   │   │  - Slash normalization   │
          │  - SQLite re-pointing    │   │  - VSCode PM untrack     │   │  - Zero false-skips      │
          └────────────┬─────────────┘   └────────────┬─────────────┘   └──────────────────────────┘
                       │                              │
                       └──────────────┬───────────────┘
                                      ▼
                        ┌──────────────────────────┐
                        │   External Integrations  │
                        │  - VSCode projects.json  │
                        │  - GitHub Desktop DB     │
                        │  - Interactive [-y] Guard│
                        └──────────────────────────┘
```

## 3. Mandatory Engineering Constraints

- Every Go file must stay under 200 lines; functions under 15 lines.
- Zero error swallowing; all errors explicitly handled or surfaced to `os.Stderr`.
- Positive conditional logic (`isTargetValid`, `hasMatch`), avoiding complex inverted expressions.
