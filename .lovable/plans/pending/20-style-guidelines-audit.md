# Plan 20: Coding Style, Formatting & Line-Gaps Remediation

## Overview

Comprehensive repository-wide audit and enforcement of Return New Line rules (R13-R16), mandatory blank lines before control structures (`if`, `for`, `switch`, `while`), blank lines after closing braces `}`, blank lines before `return`/`throw`, zero nested `if` statements, and sizing tier compliance.

## Key Guidelines Checked

1. **Rule 1:** Mandatory blank line before control structures (`if`, `for`, `switch`, `while`).
2. **Rule 2:** Mandatory blank line after closing brace `}` when followed by code.
3. **Rule 3:** Mandatory blank line before `return` / `throw` in multi-line blocks.
4. **Rule 4:** Zero clumping of consecutive guard clauses (separated by blank lines).
5. **Rule 5:** Zero nested `if` statements (nesting depth <= 1).
6. **Rule 8:** No double blank lines (`\n\n\n`) and no blank lines at function body start.
7. **Rule 9:** Sizing tiers (functions <= 8-15 lines, files <= 100 lines).

## Subtasks Breakdown

- **Subtask 20.01:** Enforce blank lines before `if` and before `return` across `src/` (`.lovable/plans/subtasks/20-style-guidelines/01-src-newlines.md`).
- **Subtask 20.02:** Enforce blank lines after `}` and guard clause spacing across `gitmap/cmd/` (`.lovable/plans/subtasks/20-style-guidelines/02-cmd-newlines.md`).
- **Subtask 20.03:** Function length and file sizing audits in helper scripts (`.lovable/plans/subtasks/20-style-guidelines/03-scripts-sizing.md`).
- **Subtask 20.04:** Automated style linter verification (`check-newline-styling.py`, `check-nested-ifs.py`, `check-function-lengths.py`).

## Acceptance Criteria

- [ ] `check-newline-styling.py` exits with 0 errors.
- [ ] `check-nested-ifs.py` exits with 0 errors across all repository files.
- [ ] Zero double blank lines anywhere in the codebase.
