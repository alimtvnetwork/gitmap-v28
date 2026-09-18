# Code Review Guidelines: Aspect Checking

These rules **MUST** be adhered to across the entire codebase to maintain consistency and prevent regressions. Next.ai and all other agents must follow this strictly.

## 1. Boolean Variables and States

- **NO INVERTED LOGIC**: NEVER use inverted properties for success checks (e.g., avoid `!response.isSuccess` or `!isValid`).
- **USE EXPLICIT POSITIVE CHECKS**: Always define and use explicit boolean state properties, such as `response.isFail`, `isInvalid`, or `isError`.

## 2. TypeScript Enums vs Strings

- **NO MAGIC STRINGS OR UNION TYPES**: In TypeScript, rather than using strings as sub-items or comparing string union types (e.g., `"pass" | "fail" | "fallback"`), you **MUST** use Enums. Enums provide better type safety and refactoring capabilities.
- **SUFFIX**: Every single Enum name must end with the suffix `Type` (e.g., `StatusType` rather than `Status`).
- **NO MAGIC STRINGS**: Do not introduce any magic strings or magic numbers anywhere in the codebase (Go, TS, or Python), unless explicitly required for a logger message, and mention that in the typing/comments.

## 3. Python and Cross-language Guidelines

- Python scripts MUST utilize the `query_wrapper` in `.github/scripts/query_wrapper.py` for any operations that can fail, completely eliminating scattered `try/catch` statements and ensuring standardized error logging to `stderr`.

## 3. Go Branching (`switch` vs `if/else`)

- **PREFER SWITCH**: When matching against multiple states, enums, or multi-condition cases, always prefer `switch` statements over chained `if/else if` blocks.
- **NO NESTED IFS**: Avoid deeply nested `if` statements. Extract to functions if necessary.
