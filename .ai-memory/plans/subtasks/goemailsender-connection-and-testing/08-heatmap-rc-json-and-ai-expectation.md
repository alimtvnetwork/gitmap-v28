# Subtask 08: Heatmap RC JSON Response & AI Expectation Tracking

## 1. Goal
Implement structured JSON output for `gitmap rc heatmap log` and `gitmap heatmap log`, including `where_saved`, `what_saved`, `db_info`, and bidirectional AI expectation lifecycle (`ExpectedOutcome`, `ActualOutcome`, `IsExpectationMet`, and status), with support for entering or updating outcome status post-execution via `update` / `outcome`.

## 2. Implementation Checklist
- [x] Create `cli/store/heatmap_rc_outcome.go` with `GetActionLogById` and `UpdateActionOutcome` for both `RcHeatmapStore` and `HeatmapDB`.
- [x] Define `RcLogResponseEnvelope` with `where_saved`, `what_saved`, `db_info`, and `ai_expectation` in `cli/heatmap/heatmap_rc_envelope.go`.
- [x] Implement `ResolveCompanionRepoSecretsUrl` and `computeRelativeDbPath` helpers.
- [x] Create `cli/heatmap/heatmap_rc_update.go` with `runRcUpdate` and `runHeatmapUpdate` handlers.
- [x] Extract rendering and help functions to `cli/heatmap/heatmap_cmd_render.go` to maintain strict <= 300 line cap across all changed files.
- [x] Add `-e` (`--expected`), `-o` (`--actual`), `-i` (`--id`) short flags in `cli/heatmap/heatmap_cmd.go`.
- [x] Support `--json` output in `runRcLog` and `runHeatmapLog`.
- [x] Add comprehensive unit tests in `cli/heatmap/heatmap_rc_outcome_test.go` verifying envelope keys, JSON marshaling, expectation transitions (PENDING -> MET -> FAILED), flag resolution, and action ID parsing.
- [x] Author technical specification in `02-spec/21-app/goemailsender-connection-and-testing/05-heatmap-rc-json-and-ai-expectation-spec.md`.

## 3. Verification
- `go test -v ./heatmap/...` passed (11/11 tests green).
- `go test ./store/...` passed.
- Manual execution of `gitmap rc heatmap log --json` verified envelope output.
- Manual execution of `gitmap rc heatmap update --id <id> --actual "..." --met true --json` verified outcome transitions.
