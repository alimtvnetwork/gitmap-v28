# Coding Guidelines

## General Project Stack
- **Languages**: Golang (Backend), TypeScript (Frontend/Scripts), PowerShell (Build/CI Scripts).
- **Formatters**: `gofmt` and `goimports` (Golang), `prettier` (TypeScript).
- **Linters**: `golangci-lint` (strict, `--issues-exit-code=1`).

## Strictly Enforced Rules (TypeScript/PHP/Python)
- **Enums Over Unions**: TypeScript string unions (`"pass" | "fail"`) are strictly banned. You must use robust Enums.
- **Enum Naming**: Every Enum must end with the `Type` suffix (e.g., `ResultType`).
- **Boolean States**: Always use explicit boolean state variables (e.g., `isFail`) instead of inverting positive ones (`!isSuccess`).
- **Query Wrappers**: All database queries must be funneled through a centralized `queryWrapper` that implements robust automatic failure logging in accordance with error management specs.
- **Magic Strings**: Magic strings/numbers are completely banned across the codebase unless used directly for error logging labels.

## Golang Rules
- No nested `if` branches; flatten logic using early returns.
- No single-character variables (`s`, `x`, `d`).
- Every error must be explicitly logged to `os.Stderr`. No silent swallows.
- `constants.go` handles all static values; versions are injected via `-ldflags` using `var` instead of `const`.
- **No State Mutation in Switch Statements**: When computing state/objects across multiple conditions (e.g., inside a switch statement), do not mutate an existing object. Instead, extract the logic into a separate function that creates and returns the object for each case, and assign the result.
- **No Nested Ternaries**: Avoid complex or nested ternary operators (?:). Favor clear, explicit conditional blocks (if/else) or early returns for readability.
