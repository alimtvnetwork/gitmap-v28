# Rule: Boolean Principles & Naming

1. **Explicit True Checks (TOTAL BAN):** Never evaluate a boolean explicitly against `true` (e.g., `if isReady == true` or `if (isValid === true)`). Always evaluate implicitly: `if isReady { ... }`.
2. **Boolean Prefixes:** Identifiers must strictly use `is` or `has` prefixes (e.g. `isValid`, `hasPermission`). Prefixes such as `can`, `should`, `was`, `will`, `did` are banned.
3. **No Inverted Success Checks:** Never invert positive success checks (e.g., `!response.isSuccess`). Use positive failure states (e.g., `response.isFail`, `isError`).
4. **No Mixed Polarity:** Never combine a positive check and a negative check in the same condition (`if isA && !isB`). Split into separate guard clauses.
