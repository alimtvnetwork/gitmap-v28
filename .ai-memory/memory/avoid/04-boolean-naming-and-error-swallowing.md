# Avoid 04: Boolean Naming and Error Swallowing

- **Category:** Anti-Pattern & Strict Constraint
- **Status:** Mandatory Enforcement

## 1. Core Rule
NEVER use negative boolean prefixes (`isNot*`, `hasNo*`) and NEVER swallow Go error tuples using blank identifiers (`_ = err`).

## 2. Rationale
- Negative booleans cause mental overhead and lead to double-negative bugs (`!isNotReady`).
- Swallowed errors conceal silent failures, turning recoverable errors into fatal crashes.

## 3. Enforcement
- Enforce affirmative prefixes: `isValid`, `isReady`, `hasData`.
- All errors must be wrapped via `*appfault.AppError` and propagated.
