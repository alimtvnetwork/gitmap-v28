# 123 — Cross-Platform Python Tooling & Automation Standard

## Overview

**Module Number:** 123
**Version:** 1.0.0
**Updated:** 2026-09-03
**Status:** Production-Ready
**AI Confidence:** Production-Ready
**Ambiguity Score:** None

---

## Purpose

To ensure developer automation, maintenance scripts, and CI/CD verification workflows execute seamlessly regardless of the host operating system (Windows, macOS, or Linux), all operational bash/shell scripts (`.sh`) are maintained as cross-platform standard-library Python scripts (`.py`).

---

## Invariants & Design Principles

### 1. Pure Standard Library (Zero External Dependencies)

All converted scripts must use only the Python standard library:
- Process management: `subprocess`
- Filesystem & paths: `os`, `sys`, `pathlib`, `glob`
- Serialization & Network: `json`, `urllib.request`, `urllib.parse`
- Storage: `sqlite3`
No `pip install` or virtual environment requirements are permitted.

### 2. Windows UTF-8 Terminal Safety

Windows PowerShell and command prompts frequently corrupt Unicode glyphs (checkmarks `✓`, boxes `─`, emojis `🚀`) when stdout defaults to legacy code pages (CP1252/CP437). Every script must enforce UTF-8 stream re-configuration:

```python
import sys
if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")
```

### 3. Exit Code Contract

- Exit `0`: Successful execution.
- Exit `1` or non-zero: Fatal error, with standard error output formatted cleanly.

---

## Converted Script Inventory

All scripts under `linter-scripts/` and operational toolchains have been converted and validated:
- `check-spec-cross-links.py`: Relative markdown link and anchor verifier.
- `audit-booleans.py`: Boolean prefix (`is`, `has`) and naming auditor.
- `verify-go-limits.py`: 200-line file and 15-line function linter.
- `release-version-v*.py`: Release packaging and asset generator.

---

## Acceptance Criteria

### Scenario 1: Windows Execution Without WSL/Bash

- **Given** a pure Windows environment without Git Bash or WSL installed
- **When** any verification or linter script is executed via `python <script>.py`
- **Then** the script executes natively with complete UTF-8 console output and zero dependency errors.

### Scenario 2: Cross-Platform Path Handling

- **Given** file paths with mixed forward and backslashes
- **When** path operations are executed
- **Then** `os.path.normpath` and `pathlib.Path` normalize paths consistently across all platforms.

---

## Cross-References

- Consolidated Guidelines: [`../17-consolidated-guidelines/00-overview.md`](../17-consolidated-guidelines/00-overview.md)
- PowerShell Integration: [`../11-powershell-integration/00-overview.md`](../11-powershell-integration/00-overview.md)
