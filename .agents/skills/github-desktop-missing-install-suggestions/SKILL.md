---
name: github-desktop-missing-install-suggestions
description: Autonomously implement and verify missing GitHub Desktop CLI detection, formatted install suggestions with exact commands and multi-platform options, --install flag support, and clean validation error handling across GitMap.
---

# GitHub Desktop Missing Install Suggestions Suite

## Overview
Guides the implementation of user-friendly missing dependency remediation for GitHub Desktop across `gitmap github-desktop` (alias `gd`) and `gitmap desktop-sync` (alias `ds`).
When the GitHub Desktop executable/CLI is absent, GitMap provides exact installation commands across GitMap Installer, APT, Snap, Winget, Chocolatey, and Homebrew, provides `--install` flag auto-triggering, and avoids dumping raw execution stack traces.

## Core Invariants
1. **Clear Diagnosis & No Raw Stack Traces**:
   - When `desktop.ResolveCLI()` returns empty, emit a clear diagnostic and return `apperror.NewValidationError(...)` or structured exit without triggering raw execution stack traces.
2. **Exact Installation Commands & Options**:
   - Show Option 1 (Recommended): `gitmap install github-desktop` (or `gitmap in github-desktop`).
   - Show Option 2 (Native Package Managers): Tailored to current OS (`runtime.GOOS`) with exact copy-pasteable commands for Linux/Ubuntu (APT shiftkey repo, Snap), Windows (Winget, Chocolatey), and macOS (Homebrew cask).
   - Show Option 3: Official GUI installer URL `https://desktop.github.com`.
3. **`--install` Flag Auto-Remediation**:
   - Support `gitmap github-desktop --install` / `-i` to invoke `gitmap install github-desktop` directly when the tool is missing.
4. **Coding Guidelines Invariants**:
   - Functions strictly <= 15 lines (target <= 8 lines).
   - Affirmative booleans only (`is*`, `has*`).
   - Zero nested ifs (guard clauses and early returns).
   - Universal `*apperror.AppError` returns.
   - Strict Unix LF line endings.
