# RCA 23: Linter Guideline Conflict (gosimple & misspell vs Coding Guidelines)

## 1. Issue Description

Running `golangci-lint run --fix` indiscriminately on the codebase caused 188 files to be modified in a highly destructive manner. Specifically, the linters `gosimple` and `misspell` auto-applied standard Go idioms and US spellings that directly conflict with the project's custom coding guidelines and JSON API payload schemas.

## 2. Root Cause Analysis

- **`gosimple` conflict**: The linter `gosimple` forces the conversion of explicit boolean checks (like `hasFlag(...) == false`) into negated boolean operations (like `!hasFlag(...)`). This directly violates the project's strict boolean naming rules (`spec/02-coding-guidelines/03-golang/02-boolean-standards.md`), which prohibit negated booleans combined with other expressions to avoid mixed polarity, and forbid the `!fn()` pattern (requiring explicit positive function names like `isInvalid()` instead of `!isValid()`, or explicit explicit value checks like `== false`).
- **`misspell` conflict**: The linter `misspell` enforces US spelling standards. It automatically converted `canceled` to `canceled`. In `clonepick/clonepick.go`, this broke the `StatusCancelled` JSON API payload value (`"canceled"` became `"canceled"`), causing a regression in downstream API consumers.
- **Unreachable Code damage**: When fixing `panic("error")` without returning a proper value, the linter caused `govet` unreachable code errors or dropped expected returns entirely.

## 3. Resolution

- **Disable Destructive Linters**: We have explicitly disabled `gosimple`, `misspell`, and `nolintlint` (which strips `nolint` comments for these rules) in `.golangci.yml`.
- **Reverted Damage**: We reverted the destructive `!hasFlag` occurrences back to explicit `== false` checks to comply with `isX && isY == false` strict boolean logic.
- **Fixed API Payload**: Restored `StatusCancelled = "canceled"` in `clonepick/clonepick.go` to fix the API contract.
- **Fixed Panics**: Replaced `panic("error")` statements in `update.go` and `taskops.go` with proper `apperror.NewSimple()` wrapping to maintain robust error management. Fixed the trailing `return` statements that caused `govet` unreachable code failures.

## 4. Preventative Measures

- **Never run `golangci-lint run --fix` blindly.**
- Custom `.golangci.yml` must explicitly `disable` linters that conflict with internal spec guidelines.
- Always run `python .lovable/ai-fix-scripts/02-cicd-local-runner.py` before committing.
