# Subtask 01 - Parameter Struct Reduction

## Parent Specification
[42-argument-reduction-audit.md](.lovable/plans/pending/42-argument-reduction-audit.md)

## Acceptance Criteria & Requirements
- Enforce value-based parameter structs (`*Params`) for functions with >2–3 parameters.
- Enforce strict affirmative boolean prefixes (`is` and `has` only) on all struct fields.
- Prohibit loose multi-parameter function signatures.
