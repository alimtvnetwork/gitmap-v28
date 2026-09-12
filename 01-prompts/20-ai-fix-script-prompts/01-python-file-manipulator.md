# Python File Manipulator CLI Specification — Tooling Spec (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously generate and maintain a robust, dependency-free Python CLI tool to handle mass file renaming, sequencing, and encoding normalization across any specified folder.

## Overview

You are an expert Python Developer AI. Your task is to write and maintain a standalone, reusable Python script that handles mass file renaming (lowercasing), sequence fixing, and encoding normalization across any target directory. This script will act as an autonomous tool for other AIs and developers to organize files without needing a compiled binary.

**Target Path:** `03-ai-scripts/03-file-manipulator.py`

---

## 🔍 Pre-Flight AI Checklist (Before Modifying or Rewriting Any Script)

Before modifying or creating any Python tool in `03-ai-scripts/` or `.agents/scripts/`, the AI agent MUST inspect these two files first:
1. **`03-ai-scripts/02-shared-engine.py`**: The canonical library containing all centralized constants, Enums (`UPPER_CASE` members matching string values), file locking, two-phase caching, and lazy regex memoization.
2. **`03-ai-scripts/01-index.md`**: The master registry of existing tools, parameters, and benchmark speeds.

### Canonical Library Import Pattern:
```python
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
try:
    from importlib import import_module
    engine = import_module("02-shared-engine")
    process_repository_files = engine.process_repository_files
    read_file_safe = engine.read_file_safe
    write_file_lf = engine.write_file_lf
    normalize_rel_path = engine.normalize_rel_path
    normalize_extensions = engine.normalize_extensions
    is_ignored_directory = engine.is_ignored_directory
    is_ignored_path = engine.is_ignored_path
    is_binary_file = engine.is_binary_file
    ExitCodeType = engine.ExitCodeType
    RegexPatternType = engine.RegexPatternType
    get_compiled_regex = engine.get_compiled_regex
    DEFAULT_TEXT_EXTENSIONS = engine.DEFAULT_TEXT_EXTENSIONS
except Exception:
    # Graceful fallback definitions if standalone
    ...
```

---

## Non-Negotiable Rules for Python Scripts

1. **Zero Dependencies**: The script MUST use only Python standard libraries (`os`, `sys`, `argparse`, `shutil`, `subprocess`, `pathlib`, `enum`, `json`, `time`, `re`, `threading`).
2. **Python Enum Generation Standard**:
   - Class Name: `PascalCase` singular noun ending with `Type` (e.g. `RegexPatternType`, `ExitCodeType`, `ScanModeType`).
   - Enum Members: `UPPER_CASE` with underscores (e.g. `UPPERCASE`, `SEQ_PREFIX`, `SUCCESS`).
   - Enum Values: String literals exactly matching the member name in `UPPER_CASE` (e.g. `UPPERCASE = "UPPERCASE"`).
3. **Implicit Booleans Only**: Never write `if is_valid == True:` (FORBIDDEN). Always write `if is_valid:`.
4. **Small Decomposed Functions**: Code must be decomposed into small, testable functions under 25 lines each following the Single Responsibility Principle.
5. **DRY Shared Engine**: Import common filesystem scanning, caching, regexes, and line-ending utilities from `00-shared-engine.py`.
6. **Thread-Safe Lazy Regex Memoization**: Consume regex patterns via `get_compiled_regex(RegexPatternType.<MEMBER>)` or `RegexRegistry.get(...)`.
7. **Multi-Folder Scoping**: Support target directory specification (`<path>` or `--path` / `--dir`) so the tool can run on any directory, repository root, or nested subfolder.
8. **Customizable File Extensions**: Support `--ext` for filtering specific file extensions (e.g. `.md,.ts,.py`), normalizing extensions to lowercase with leading dots.
9. **Nested Ignore Pruning**: Prune `.git`, `.gitmap`, `node_modules`, `dist`, `build`, `.venv`, `.gemini`, `tmp`, `.system_generated`, and `release-artifacts` at all subtree depths including nested subprojects.
10. **Windows Long Paths & Case-Hop**: Normalize paths and safely handle Windows NTFS case-only renames with an intermediate hop.
11. **Git Awareness**: Attempt `git mv` via subprocess first, gracefully falling back to standard `shutil.move` / `os.rename`.
12. **Update Index & Mirror**: Document usage in `03-ai-scripts/01-index.md` and sync changes to `.agents/scripts/`.

---

## Core Feature 1: Lowercase Renamer

**Command Pattern**:
```bash
python 03-ai-scripts/03-file-manipulator.py lowercase <target_directory> [flags]
```

**Requirements**:
1. Recursively convert all files matching a target pattern to lowercase using `RegexPatternType.UPPERCASE`.
2. **Extension Enforcement**: Strictly ensure all markdown and code file extensions are lowercased (e.g., converting `.MD` to `.md`).
3. **Default Ignores**: Silently ignore `node_modules`, `.git`, `.gitmap`, and excluded folders from `00-shared-engine.py`.
4. **Extendable Ignores**: Provide `--except` accepting a comma-separated list of additional files or patterns.

---

## Core Feature 2: Fix File Sequencing (`fix-seq-files`)

**Command Pattern**:
```bash
python 03-ai-scripts/03-file-manipulator.py fix-seq-files <target_directory> [flags]
```

**Requirements**:
1. Scan the specified directory for sequenced files using `RegexPatternType.SEQ_PREFIX`.
2. **Ordering Flags**:
   - `--order-by-time`: Re-sequence files sequentially based on modification time.
   - `--order-by-az`: Re-sequence files alphabetically.
3. **Tie-Breaker / Preservation**:
   - `--keep-old-order`: Preserve existing numeric ordering where possible.
4. **Fixated / Pinned Sequences**:
   - `--pin "<mapping>"`: Allow users to explicitly lock specific files (e.g., `--pin "readme=00,intro=01"`).

---

## Core Feature 3: Fix Encoding & Line Endings (`fix-encoding`)

**Command Pattern**:
```bash
python 03-ai-scripts/03-file-manipulator.py fix-encoding <target_directory> [flags]
```

**Requirements**:
1. Scan the specified directory and aggressively normalize encoding for all text files (customizable via `--ext`, defaulting to `DEFAULT_TEXT_EXTENSIONS`).
2. **BOM Stripping**: Detect and strip any UTF-8 Byte Order Marks (BOM), standardizing strictly to UTF-8 without BOM.
3. **Line Ending Normalization**: Automatically convert Windows CRLF (`\r\n`) to Unix LF (`\n`) using `RegexPatternType.CRLF` or `RegexPatternType.UNIVERSAL_LINE_ENDING`.

---

## Execution Checklist for the AI

Before completing this task, you MUST verify:

- [ ] Checked `00-shared-engine.py` and `01-index.md` before writing code.
- [ ] Saved the script precisely to `03-ai-scripts/03-file-manipulator.py` and synced to `.agents/scripts/`.
- [ ] Used `RegexPatternType` with UPPER_CASE members (`UPPERCASE`, `SEQ_PREFIX`, `CRLF`).
- [ ] Used `argparse` for subcommands (`lowercase`, `fix-seq-files`, and `fix-encoding`) with detailed examples.
- [ ] Pruned `node_modules`, `.git`, `.gitmap` via `is_ignored_directory()`.
- [ ] Attempted `git mv` where applicable to preserve history.
- [ ] Handled Windows long paths via `normalize_rel_path()`.
- [ ] Ensured all 18 CI gates pass via `python 03-ai-scripts/06-cicd-local-runner.py`.

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.
