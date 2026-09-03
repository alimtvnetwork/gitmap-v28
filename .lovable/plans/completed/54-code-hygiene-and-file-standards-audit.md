# 33 - Code Hygiene & Project Architecture Audit Specification

## 1. Verbatim Acceptance Criteria Echo (from spec/02-coding-guidelines/01-cross-language/97-acceptance-criteria.md)

### AC-01: Guideline Coverage
- [ ] Boolean principles define naming, evaluation, and composition patterns
- [ ] Casting elimination patterns cover type-safe alternatives to type assertions
- [ ] Code style defines formatting, naming, and structural conventions

### AC-02: Enforcement
- [ ] All guidelines include ❌ (forbidden) and ✅ (compliant) code examples
- [ ] ESLint/linter rules are documented for automated enforcement
- [ ] Master guidelines document consolidates all standards for AI reference

---

## 2. Task-Specific Rule Set (Domain Rules)

1. **Rule CH-1 (Strict Unix LF & UTF-8 Encoding):** Every file must use Unix LF (`\n`, `0x0A`) line endings. Windows CRLF (`\r\n`) and UTF-8 BOM (`\xef\xbb\xbf`) are strictly forbidden.
2. **Rule CH-2 (Mandatory Single Trailing Newline):** Every file must terminate with exactly one newline (`\n`) at EOF. Zero trailing blank lines.
3. **Rule CH-3 (No Function Starts with Blank Line):** Executable code in any function must begin immediately on line 1 after the opening brace `{`.
4. **Rule CH-4 (Zero Double Blank Lines):** Never use two or more consecutive blank lines (`\n\n\n`) in source code or markdown documents.
5. **Rule CH-5 (Markdown Heading Spacing):** Exactly one blank line before and after every markdown heading (no leading blank line on line 1).

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | Repository-wide | Various | CRLF line endings in 342 files | Converted 342 files from CRLF to Unix LF via automated script | FIXED |
| V-02 | Repository-wide | Various | Trailing whitespace and missing single newline at EOF | Cleaned and normalized 620 files via `03-ai-scripts/04-newline-fixer.py` | FIXED |
| V-03 | spec/ & docs/ | Various | Markdown heading spacing violations | Fixed 75 markdown files with `linter-scripts/check-markdown-headings.py --fix` | FIXED |
| V-04 | linter-scripts/check-newline-styling.py | N/A | Newline styling linter | Verified exit code 0 | VERIFIED |
| V-05 | 03-ai-scripts/06-cicd-local-runner.py | 34 | Newline styling quality gate | Verified in Batch 1 of local runner | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/33-hygiene/01-line-endings-and-encoding.md`)**:
   Verify strict Unix LF (`\n`), UTF-8 (no BOM), and single terminating EOF newline across all files.
2. **Subtask 02 (`.lovable/plans/subtasks/33-hygiene/02-newline-and-heading-spacing.md`)**:
   Verify zero double blank lines (`\n\n\n`) and enforce markdown heading spacing (MD022/MD032).
3. **Subtask 03 (`.lovable/plans/subtasks/33-hygiene/03-ci-runner-verification.md`)**:
   Execute `python 03-ai-scripts/06-cicd-local-runner.py` with exit code 0 across all 15 quality gates.
