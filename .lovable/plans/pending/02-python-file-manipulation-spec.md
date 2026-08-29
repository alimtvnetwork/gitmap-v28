# AI Implementation Spec: Python File Manipulation CLI (`ai-fix-scripts`)

## Overview

You are an expert Python Developer AI. Your task is to write a standalone, reusable Python script that handles mass file renaming (lowercasing) and sequence fixing. This script will act as an autonomous tool for other AIs and developers to organize files without needing a compiled binary.

**Target Path:** `.lovable/ai-fix-scripts/01-file-manipulator.py`

## Non-Negotiable Rules for the Python Script

1. **Zero Dependencies**: The script MUST use only Python standard libraries (e.g., `os`, `sys`, `argparse`, `shutil`, `subprocess`, `pathlib`).
2. **Robust CLI**: Use `argparse` to provide a professional, CLI-like experience with complete `--help` documentation and examples.
3. **Windows Long Paths**: The script must normalize paths and safely handle Windows `MAX_PATH` limitations (e.g., prefixing absolute paths with `\\?\` on Windows environments).
4. **Git Awareness**: Whenever renaming a file, the script must attempt to use `git mv` via `subprocess` first. If the file is untracked or the command fails, gracefully fallback to standard `os.rename`.
5. **Update Index**: After generating the script, you MUST document its usage in `.lovable/ai-fix-scripts/index.md`.

---

## Core Feature 1: Lowercase Renamer

**Command Pattern**:
`python 01-file-manipulator.py lowercase <target_directory> [flags]`

**Requirements**:
1. Recursively convert all files matching a target pattern to lowercase.
2. **Default Ignores**: By default, the script MUST silently ignore `node_modules` and `.git` folders. Do not traverse them.
3. **Extendable Ignores**: Provide an `--except` flag accepting a comma-separated list of additional files, folders, or wildcard patterns to ignore (e.g., `--except "docs/*, temp.md"`).

**Example Output in `--help`**:
- `python 01-file-manipulator.py lowercase ./src` (Ignores node_modules/.git by default)
- `python 01-file-manipulator.py lowercase ./src --except "vendor/*, build/*"`

---

## Core Feature 2: Fix File Sequencing (`fix-seq-files`)

**Command Pattern**:
`python 01-file-manipulator.py fix-seq-files <target_directory> [flags]`

**Requirements**:
1. Scan the specified directory for sequenced files (e.g., `01-draft.md`, `02-notes.md`).
2. **Ordering Flags**:
   - `--order-by-time`: Re-sequence files sequentially based on their filesystem modification time.
   - `--order-by-az`: Re-sequence files alphabetically based on the string following the sequence number.
3. **Tie-Breaker / Preservation**:
   - `--keep-old-order`: Preserve existing numeric ordering as much as possible. Only assign new sequence numbers to unnumbered files or resolve direct conflicts using time/alphabetization.
4. **Fixated / Pinned Sequences**:
   - `--pin "<mapping>"`: Allow users to explicitly lock specific files to a sequence number. (e.g., `--pin "readme=00,draft=01"`). The script must increment other files around these locked sequences.

**Example Output in `--help`**:
- `python 01-file-manipulator.py fix-seq-files ./docs --order-by-time`
- `python 01-file-manipulator.py fix-seq-files ./docs --order-by-az --keep-old-order`
- `python 01-file-manipulator.py fix-seq-files ./docs --pin "readme=00,intro=01"`

---

## Execution Checklist for the AI

Before completing this task, you MUST verify:
- [ ] I saved the script precisely to `.lovable/ai-fix-scripts/01-file-manipulator.py`.
- [ ] I used `argparse` to handle subcommands (`lowercase` and `fix-seq-files`) and provided detailed help text.
- [ ] `node_modules` and `.git` are hardcoded into the default ignore list.
- [ ] Renames use `git mv` where applicable to preserve history.
- [ ] I implemented the pinning (`--pin`) logic for sequences.
- [ ] I handled Windows long paths properly via path normalization.
- [ ] I updated `.lovable/ai-fix-scripts/index.md` with instructions on how to use this new script.
- [ ] I did NOT leave any `TODO` placeholders in the generated Python code.
