---
name: gitmap-scanner-and-cloner
description: Autonomously audit, extend, test, and modify GitMap's repository discovery scanner, output formatters (CSV, JSON, Markdown), and high-speed parallel cloner.
---

# GitMap Scanner and Cloner Skill (`gitmap-scanner-and-cloner`)

## Mission & Purpose
This skill provides authoritative architectural guidance, code navigation, and execution rules for developing, auditing, and enhancing GitMap's repository scanning, formatting, and re-cloning engine.

---

## 1. Key Architectural Components & Code Map

| Component | Primary Location | Key Responsibilities |
|---|---|---|
| **Scanner Engine** | `cli/scanner/`, `cli/cmdscan/` | Recursive BFS/DFS traversal, `.git` repository identification, worktree detection, remote URL extraction (HTTPS/SSH), and active branch resolution. |
| **Output Formatters** | `cli/formatter/`, `cli/render/` | Generates standardized output artifacts: `gitmap.csv`, `gitmap.json`, and `folder-structure.md`. Aligned terminal glyphs and colors. |
| **Cloner Engine** | `cli/cloner/`, `cli/cmdclone/` | High-speed concurrent repository cloning, directory structure reproduction, retry logic, and branch checkout. |
| **Concurrency Pool** | `cli/cloneconcurrency/`, `cli/worker/` | Controlled worker pools for multi-repo cloning without resource exhaustion or rate-limiting. |
| **Git Utilities** | `cli/gitutil/` | Native Git CLI invocation, config parsing, credential helpers, and subprocess management. |

---

## 2. Essential Commands

```bash
# Basic scan with default outputs (terminal, CSV, JSON, folder-structure.md)
gitmap scan [path]
gitmap s [path]

# Scan with SSH remote URLs prioritized
gitmap scan [path] --mode ssh

# Scan and automatically add discovered repositories to GitHub Desktop
gitmap scan [path] --github-desktop

# High-speed clone execution using generated manifest
gitmap clone [manifest.json] --concurrency 8

# Re-clone preserving full directory hierarchy
gitmap clone --preserve-hierarchy --dest [target-dir]
```

---

## 3. Core Invariants & Engineering Guardrails

1. **Output Contract Integrity:** Every scan must deterministically generate identical JSON and CSV schemas adhering to `02-spec/08-json-schemas/`. Never alter schema keys without migrating JSON Schema contract tests (`*_contract_test.go`).
2. **Deterministic Folder Structure:** `folder-structure.md` must render clean Unicode directory trees with package glyphs (`📦`) and branch indicators without trailing whitespace.
3. **Cross-Platform Pathspec Normalization:** Normalize all Windows backslashes (`\`) to forward slashes (`/`) internally to ensure consistent hashes, CSV paths, and git remotes across Windows, macOS, and Linux.
4. **Zero Error Swallowing:** All filesystem or Git subprocess errors must be captured as `*apperror.AppError` with clear error codes (`[E1000:CLI]`, `[E9000:EXECUTION]`) and full stack traces.
5. **No Routine Unit Test / Build Execution:** When developing scanner or cloner features, rely on targeted static analysis and file inspections. Do not run uncommanded test suites or builds during routine turns.
