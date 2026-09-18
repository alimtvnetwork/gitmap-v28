# 121 — Cross-Platform Duplicate Audit & Immediate Remediation

## Overview

**Module Number:** 121
**Version:** 1.0.0
**Updated:** 2026-09-03
**Status:** Production-Ready
**AI Confidence:** Production-Ready
**Ambiguity Score:** None

---

## Purpose

During repository scanning, cloning, or tool synchronization, identical repositories or projects often get registered repeatedly across developer tooling (Antigravity workspaces, VS Code Project Manager, Chrome browser profiles, and GitMap database tracking). This specification formalizes the centralized duplicate detection engine (`gitmap find-duplicates`) and the contract to output copy-pasteable remediation CLI commands directly beneath findings.

---

## Detection Scope & Subcommands

### 1. Unified Dispatch: `gitmap find-duplicates [platform] [flags]`

When executed with no platform argument, runs a full multi-platform audit across:
- **Antigravity Workspaces (AGY)**: `~/.gemini/config/projects/*.json`
- **VS Code Project Manager**: `%APPDATA%/Code/User/globalStorage/alefragnani.project-manager/projects.json`
- **Chrome Profiles**: `%LOCALAPPDATA%/Google/Chrome/User Data/`
- **Git Repositories**: Internal GitMap database `Repo` records

### 2. Platform-Specific Subcommands

- `gitmap agy find-duplicates`
- `gitmap vscode find-duplicates`
- `gitmap chrome find-duplicates`
- `gitmap git find-duplicates`

---

## Output Contract: Findings & Immediate Remediation

The terminal output renders two distinct blocks:

### 1. Duplicate Findings Grouping

```text
── Duplicate Antigravity Projects Found: 2 Group(s) ──

[Group 1] Path: D:\wp-work\riseup-asia\macro-ahk
  • 56f4c903-1658-492f-b31d-a9c25cce4c0d  macro-ahk-v55
  • b0474ef4-d701-4a1c-8299-b65568151ccd  macro-ahk-broken
```

### 2. Immediate Remediation Recipes

Directly beneath the table, the CLI outputs ready-to-run copy-pasteable commands:

```text
Remediation Commands:
────────────────────────────────────────────────────────────────────────────────
● Batch optimize (deduplicate keeping newest):
  gitmap agy optimize-projects

● Remove specific duplicate project:
  gitmap agy rm b0474ef4-d701-4a1c-8299-b65568151ccd

● Clear all except preserved project:
  gitmap agy clear --except "56f4c903-1658-492f-b31d-a9c25cce4c0d"

● Fix clone directory repeats:
  gitmap clone --fix
```

---

## Prevention During Scanning & Cloning

1. **Scan Invariant**: Before registering a discovered repository path into the GitMap database or external tooling, the scanner verifies existing entries. If the normalized path exists, registration is skipped as a byte-stable no-op.
2. **Clone Invariant**: If `gitmap clone` targets a repository whose path is already cloned locally, cloning aborts with a duplicate notification and suggests `gitmap pull` or `gitmap clone --fix`.

---

## Acceptance Criteria

### Scenario 1: Duplicate Audit

- **Given** multiple project configurations pointing to `D:\wp-work\riseup-asia\macro-ahk`
- **When** `gitmap agy find-duplicates` is executed
- **Then** the terminal renders the duplicate cluster and lists the exact `optimize-projects` and single `rm <id>` remediation commands.

### Scenario 2: Batch Optimization

- **Given** detected duplicates across Antigravity projects
- **When** `gitmap agy optimize-projects` is executed
- **Then** duplicate entries are pruned while preserving the newest configuration file, leaving zero duplicates.

---

## Cross-References

- VS Code Project Manager Sync: [`01-vscode-project-manager-sync/00-overview.md`](./01-vscode-project-manager-sync/00-overview.md)
- Cloner Specification: [`05-cloner.md`](./05-cloner.md)
- Scanner Specification: [`03-scanner.md`](./03-scanner.md)
