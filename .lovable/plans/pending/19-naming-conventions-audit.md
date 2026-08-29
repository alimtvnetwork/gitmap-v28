# Plan 19: Variable & Boolean Naming Conventions, Anti-`ok` & Positive Framing

## Overview
Comprehensive repository-wide refactoring of all variable and boolean naming violations, eliminating bare `ok` variables, replacing negative boolean names (`hasNo*`, `isNot*`), enforcing affirmative prefixes (`is`, `has`, `can`, `should`), and applying positive framing with inverted `if` guard clauses.

## Key Audit Inventory
- **Bare `ok` Variables:** 114 instances across Go type assertions, map lookups, and channel receives.
- **Negative Booleans:** 0 remaining across core runtime (previously cleaned up `hasNoColors`, `hasNoPayload`).
- **Acronym Casing:** Normalized to PascalCase (`Id`, `Url`, `Api`).

## Subtasks Breakdown
1. **Subtask 19.01:** Refactor bare `ok` in `gitmap/cmd/` and `gitmap/tui/`.
2. **Subtask 19.02:** Refactor bare `ok` in `gitmap/store/`, `gitmap/gitutil/`, and `gitmap/release/`.
3. **Subtask 19.03:** Refactor bare `ok` in `gitmap/fixtureversion/`, `gitmap/helptext/`, `gitmap/startup/`, and `gitmap/vscodepm/`.
4. **Subtask 19.04:** Refactor bare `ok` in `gitmap/constants/`, `scripts/`, and Go tests.
5. **Subtask 19.05:** Linter verification & local CI runner gate check.

## Acceptance Criteria
- [ ] Zero bare `ok` variables in any Go source file (`isFound`, `isAppErr`, `isKeyMsg`, `isSuccess`, `hasTag`).
- [ ] Affirmative boolean naming enforced repository-wide.
- [ ] `check-boolean-guidelines.py` and `check-enum-and-boolean.py` exit 0.
- [ ] Go unit tests pass cleanly (`go test ./...`).
