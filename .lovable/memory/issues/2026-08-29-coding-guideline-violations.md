# Issue: Coding Guideline Violations & Boolean Anti-Patterns (Audit v4)

- date: 2026-08-29
- status: identified
- scope: cross-package

## 1. Why it happened
During rapid iterative feature development across CLI commands, installers, cluster orchestration, and git transports, developer focus was centered on functionality and passing end-to-end flows rather than strict adherence to low-level structural coding standards. This led to the gradual accumulation of monolithic functions exceeding line caps, unchecked negative booleans, and missing type suffixes.

## 2. How it happened
- Functions were written inline without decomposing secondary concerns into discrete helpers, pushing function lengths beyond the 15-line hard cap (reaching 40-100+ lines in legacy modules).
- Boolean conditions were negated with `!` operators (`!isReady`, `!found`, `!ok`) directly inside complex conditional expressions, violating positive framing and explicit boolean comparison standards.
- Custom domain types defined as string/int aliases (acting as enums) omitted the required `Type` suffix.
- Error suppression patterns like `_ = err` were occasionally used in non-critical cleanup tasks without wrapping or logging context.

## 3. Root Cause
- Absence of strict custom AST linter gates in older CI iterations allowing functions > 15 lines and inverted boolean identifiers to enter the repository.
- File and symbol evolution without scheduled automated cleanup cycles.

## 4. Code Fix
- Execute the 150-step plan defined in `.lovable/plans/pending/01-coding-guideline-fixes.md`.
- Extract positive guard variables (`isMissingEntry := !ok; if isMissingEntry == true { ... }`).
- Append `Type` suffixes to all enum aliases (`Format` -> `FormatType`).
- Decompose oversized functions into <= 15-line helpers using shell-and-wire dispatch patterns.
- Wrap all errors with `apperror.Wrap` and contextual operation parameters.
