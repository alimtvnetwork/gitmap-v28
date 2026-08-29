# 03. Generic Error Handling Guidelines & Centralized Dispatch

> **Canonical Mirror:** [centralized-error-handling-architecture.md](../../.lovable/coding-guidelines/centralized-error-handling-architecture.md)

Please refer to the root architectural specification at [`.lovable/coding-guidelines/centralized-error-handling-architecture.md`](../../.lovable/coding-guidelines/centralized-error-handling-architecture.md).

## Quick Summary of Non-Negotiable Rules

1. **Never Bare Exit or Panic**: `os.Exit(1)` and `panic("...")` at arbitrary call sites are prohibited.
2. **Always Wrap in Domain `AppError`**: Every error MUST include `Op`, `Code`, `Type`, `Severity`, `Creator`, and `Ctx`.
3. **Dispatch via Central Handler**: Route all error paths through `cliexit.HandleError(err, ...)` or the language-equivalent central dispatcher.
4. **Never Be Silent**: Swallowed errors or silent terminations are auto-reject violations.
