# Issue Domain 09: Coding Guideline and Linter Violations

- **Domain:** Boolean Naming, Function Sizing, and Blank Line Hygiene
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Negative Boolean Polarity
- **Symptoms:** Linter flagged negative variable names like `isNotValid`, `hasNoErrors`, `isDisabled`.
- **Root Cause:** Legacy code inverted logical checks directly in variable names.
- **Resolution:** Converted all identifiers to affirmative prefixes (`isValid`, `hasErrors`, `isEnabled`).

## 2. Functions Exceeding Size Caps
- **Symptoms:** CI linter failed on functions exceeding 15 lines of executable logic.
- **Root Cause:** Complex orchestration and error handling nested within a single function block.
- **Resolution:** Decomposed monoliths into specialized helper functions adhering to <= 8–15 line limits.
