# Plan 102: Code Hygiene & Universal Encoding Standards Audit

## Executive Summary
This master architectural plan establishes repo-wide compliance with universal file hygiene, encoding standards, and line ending specifications defined in `spec/02-coding-guidelines/08-file-folder-naming/`, `spec/02-coding-guidelines/01-cross-language/04-code-style/`, and `.lovable/coding-guidelines.md`.

## Core Objectives
1. **Unix LF (`\n`) Line Endings:** Total ban on Windows CRLF line endings; enforce pure LF line feeds across all text and code files.
2. **UTF-8 (No BOM) Encoding:** Enforce strict UTF-8 without Byte Order Mark headers across all source files, schemas, and markdown documents.
3. **Newline Styling & Spacing:** Exactly one terminating newline at EOF; zero double blank lines in code; single blank line before `return`/`throw` and after closing `}`.
4. **Automated Hygiene Verification:** Verify with `check-newline-styling.py` and `10-encoding-normalizer.py`.

## Verification Results
- `python 03-ai-scripts/10-encoding-normalizer.py --fix`: PASS (11,599 text files verified UTF-8 LF)
- `python 03-ai-scripts/04-newline-fixer.py --fix`: PASS (252 files normalized, trailing whitespace and final newlines ensured)
- `python linter-scripts/check-newline-styling.py`: PASS (all files meet newline styling standards)
