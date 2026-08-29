# Subtask 03: Cloner and Scanner Parameter Structs

- **Status:** `PENDING`
- **Target Files:** `gitmap/cloner/`, `gitmap/scanner/`

## Instructions

1. Audit multi-parameter functions in `gitmap/cloner/` and `gitmap/scanner/`.
2. Group loose arguments into parameter structs (`ScanParams`, `ClonerBatchParams`).
3. Enforce `is`/`has` boolean prefixing on all struct fields.
4. Eliminate void functions where error reporting is required by returning `*apperror.AppError`.
