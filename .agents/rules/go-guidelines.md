# Rule: Go Language Guidelines

1. **No Underscores:** Go identifiers, structs, fields, and variables must use `camelCase` or `PascalCase`. Snake_case is strictly banned.
2. **Error Types:** Functions returning errors must return `*apperror.AppError`.
3. **Parameter Structs:** Avoid functions with more than 2-3 parameters; encapsulate into `*Params` structs.
