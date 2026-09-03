# Subtask 02 - AppError & Result Returns

## Parent Specification
[42-argument-reduction-audit.md](.lovable/plans/pending/42-argument-reduction-audit.md)

## Acceptance Criteria & Requirements
- Mandate `*apperror.AppError` returns for side-effect operations in Go domain logic (eliminating bare "void" functions).
- Wrap external and standard library errors into `*apperror.AppError` context wrappers.
- Return single `Result[T]` envelopes for data retrieval functions.
