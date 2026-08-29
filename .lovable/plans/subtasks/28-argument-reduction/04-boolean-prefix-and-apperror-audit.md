# Subtask 04: Boolean Prefix and AppError Audit

- **Status:** `PENDING`
- **Target Files:** `gitmap/archive/`, `gitmap/apperror/`, `gitmap/cliexit/`

## Instructions

1. Verify all boolean parameters and struct fields use `is` or `has` exclusively.
2. Ensure no negative booleans exist (`isNot*`, `hasNo*`).
3. Ensure no bare `panic` calls exist.
4. Wrap external framework errors into `*apperror.AppError`.
