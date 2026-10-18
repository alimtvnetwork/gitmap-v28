# Subtask 04: Relative DB Paths & Repo-Secrets URL Exposure in Heatmap and Agent

- **Task-ID**: Task-04
- **Subtask**: 04-db-path-relative-and-repo-secrets-url
- **Owner**: Worker 01
- **Status**: PENDING

## Owned Files (Exclusive)
- `cli/heatmap/heatmap_view.go`
- `cli/heatmap/heatmap_query.go`
- `cli/cmdagent/agent_task_sqlite.go`
- `cli/heatmap/heatmap_db_paths_test.go`

## Read-Only Context Files
- `02-spec/21-app/goemailsender-connection-and-testing/02-component-spec.md`
- `cli/store/specialrepos.go`

## Objectives
1. Update `cli/heatmap/heatmap_view.go` and `cli/heatmap/heatmap_query.go` to compute relative DB paths (`filepath.Rel(repoRoot, dbPath)`) and inject `repoSecretsUrl` into JSON output envelopes.
2. Update `cli/cmdagent/agent_task_sqlite.go` (`runTaskSQLiteStatus`, `runTaskSQLiteInit`) to emit `relativeDbPath`, `absoluteDbPath`, and `repoSecretsUrl`.
3. Add unit test `cli/heatmap/heatmap_db_paths_test.go` asserting relative paths and `repoSecretsUrl` in JSON responses.
