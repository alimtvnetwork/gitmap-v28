# Error Management Guidelines

These rules dictate how errors should be caught and logged across the codebase to reduce scattered try/catch blocks and guarantee unified logging formatting.

## 1. Centralized Query Wrapper
- **MANDATORY USAGE**: All queries or operations that can fail in TypeScript/PHP/Python **MUST** be wrapped using a centralized wrapper (e.g., `queryWrapper`).
- **NO SCATTERED TRY/CATCH**: Do not write raw `try/catch` blocks scattered across application logic. Use the wrapper to handle the exception, explicitly logging the failure, and returning a structured state.

## 2. Structured Error State
- The wrapper should return a defined type containing explicit boolean properties such as `isFail: true` instead of relying on consumers to check `!isSuccess` or `!data`.

## 3. Explicit Logging
- All caught errors MUST be explicitly logged. The wrapper is responsible for formatting this log.
- Do not introduce random magic strings for logging. If a specific context string is required, pass it to the wrapper and declare it clearly.
