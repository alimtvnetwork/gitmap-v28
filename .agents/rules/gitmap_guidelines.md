---
description: Enforces Gitmap's strict error management architecture and coding style guidelines.
trigger:
  - ".*"

---
  **Gitmap Strict Error Architecture and Coding Guidelines**

  When working on the Gitmap repository, you MUST adhere to the following rules at all times. Failure to follow these rules will break the DB audit logging and violate the project's standards.

  ### 1. Error Management & Architecture (CRITICAL)
  - **Error Return Sovereignty**: Leaf functions that declare an `error` return type MUST return the error directly (`return err` or `return apperror.Wrap(...)`). Leaf functions MUST NEVER call exit handlers (`cliexit.HandleError`) or `os.Exit` and then return `nil`.
  - **Outer Caller Handles Exits**: Only the top-level orchestrator or root command dispatcher receives the error and decides how to terminate or report.
  - **No Magic Integer Exit Codes**: Never pass raw integers (`1`, `2`) to exit handlers. Always use strongly-typed enums ending in `Type` (`ExitCodeType`, `constants.ExitGeneralError`, `constants.ExitUsageError`).
  - **Parameter Reduction via Specialized Helpers**: Repeated handler calls with constant codes should use specialized helpers (`cliexit.HandleValidationError`, `cliexit.HandleUsageError`).
  - **Database Audit Capture**: Errors bubble up to `dispatch` handlers in `root.go` so they are captured by `finishCommandAudit`.

  ### 2. Coding Style & Naming Rules
  - **Boolean Naming**: All booleans MUST begin with `is` or `has` ONLY (`can`, `should`, `was`, `will` are banned). NO negative booleans.
  - **15-Line Function Limit**: Functions <= 8 lines preferred, <= 15 lines max.
  - **Zero Nested Ifs**: Flatten all nested conditionals with guard clauses and early returns.
  - **Return Styling**: Blank line after `}` and blank line before `return`.

  ### 3. Automated Refactoring Warning
  - **No Blind Regex Mass Refactoring**: Never blindly regex-replace exit calls without matching function return signatures.
