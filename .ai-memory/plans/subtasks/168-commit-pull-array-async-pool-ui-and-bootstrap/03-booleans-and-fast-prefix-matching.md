# Subtask 03: Boolean Normalization & Fast Prefix Matching over Regex
Traceability ID: Task-04
Spec Reference: [02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md](../../../02-spec/21-app/168-commit-pull-array-async-pool-ui-and-bootstrap.md)
Target Files: cli/cmd/commitin/lineskipper/line_skipper.go, cli/cmd/commitin/parse_types.go, cli/cmd/commitin/config_json.go
Action:
- Audit all booleans across `cli/cmd/commitin/` to ensure compliance with positive boolean conventions (`is...`/`has...`, e.g. `isDryRun`, `isFinalSync`, `isRecreate`, `isTree`, `hasCache`, `isFound`, `isReady`).
- Optimize `line_skipper.go` to evaluate zero-allocation `strings.HasPrefix` (for mode `starts_with`) and `strings.Contains` (for mode `contains`) before falling back to regex.
- Pre-compile any required regex patterns in memory into an in-memory cache/map upon config load so zero regexes are compiled inside tight line-by-line loops.
Acceptance Criteria:
- Zero raw booleans without `is` or `has` prefix in commit-in types.
- Line skipper prioritizes `strings.HasPrefix` and pre-compiled regex cache.
Targeted Verification: python linter-scripts/check-enum-and-boolean.py
