# 263 — Pipeline Errors All: Aggregate + `--file` + `te` Shortcut

Status: spec (approved) · Date: 2026-10-09 · Owner task: `263-pipeline-errors-all`

## Owner request (verbatim)

"add a new command `gitmap pipeline errors all` — scans ALL repositories we have, combines all repositories' errors, shows them in the terminal; with `--file` the combined report goes to a single file which AI models can read and fix. Also shortcut `te` so `gitmap te all` works. Implement, minor bump, release."

## Current behavior (as-is, verified in tree)

- Top-level `pe` / `pipeline-errors` / `pipeline_errors` / `ee` dispatches in `cli/cmd/rootutility.go:453` (`utilityPipelineEntries`) to `cmdpipeline.RunPipelineErrors(argsTail())`.
- `RunPipelineErrors` (`cli/cmdpipeline/exports.go:13`) calls `runPipeline(append([]string{"errors"}, args...))`.
- `runPipeline` (`cli/cmdpipeline/pipeline.go:143`) routes via `checkErrorLogsSubcmd` → `checkErrorLogsOrClear` → `isErrorLogsSubcmd` (`pipeline.go:202`) → `handlePipelineErrorLogs` (`cli/cmdpipeline/pipeline_logs_entry.go:14`) → `executePipelineErrorLogs` (`pipeline_logs_entry.go:51`).
- `PipelineErrorFlags` (`cli/cmdpipeline/pipeline_flags.go:13`): `IsAll` set from `all` / `--all` positional (`pipeline_flags.go:62`); `FilePath` parsed from `--file` (`pipeline_flags.go:62`) but **not wired into the all-repos path**.
- When `flags.IsAll`, `executeAllPipelineErrorLogs` runs (`cli/cmdpipeline/pipeline_all_errors.go:36`): builds `AllPipelineSummary` (`TotalPipelines`, `CleanCount`, `FailedCount`, `FailedItems []AllPipelineFailedItem`, `CleanRepos`, `PullErrors []store.PullErrorRecord`); `--json` prints indented JSON; otherwise prints the terminal view.
- Repo discovery today: `discoverPipelineRepoSlugs` (`pipeline_all_errors.go:77`) walks the **filesystem** — subdirs of `<bindir>/data/pipeline/` that contain `sql.db`. It does NOT use the repo catalog.
- `cli/cmdpipeline/pipeline_all_errors_test.go` covers flag parsing (`TestParsePipelineErrorFlags_All`) and the all-path entry (`TestExecuteAllPipelineErrorLogs`).

## Goals (to-be)

1. `gitmap pipeline errors all` scans ALL repositories from the repo catalog, combines every repo's pipeline errors, and renders them in the terminal grouped by repo.
2. `--file <path>` writes the combined report to ONE file as AI-readable markdown (models can read it and fix the errors). Terminal still prints the summary.
3. New top-level shortcut `te` (alias of `pe`) so `gitmap te all` works.
4. Minor version bump + release after implementation (release step owned by lead).

## Decisions

- **D1:** "ALL" = all repos in the catalog (`store.OpenDefault()` + `db.ListRepos()` → `[]model.ScanRecord`, `cli/store/repo.go:91`). Repos without pipeline data appear as **"no pipeline data" rows** — they are not silently dropped.
- **D2:** Terminal view is grouped by repo: `repo → error count → workflow/step, message snippet, timestamp`, followed by a combined summary: **repos scanned, repos with errors, total errors**.
- **D3:** Implementation stays inside `cli/cmdpipeline/` (no new top-level package). New formatter file `pipeline_all_errors_format.go`, ≤300 lines.
- **D4:** `--file` + `--json` combined → the JSON report is written to the file; terminal keeps the human summary. `--file` alone → AI-readable markdown file; terminal keeps the human summary.
- **D5:** `executeAllPipelineErrorLogs` keeps its `error` signature (callers unchanged); internal failures are wrapped in `*apperror.AppError` per repo code rules.
- **D6:** "No pipeline data" detection must `os.Stat` the file at `pipelinedb.ResolvePipelineDbPath(slug)` **before** opening, because `pipelinedb.OpenPipelineSplitDb` **creates** the DB as a side effect (`cli/pipelinedb/pipeline_split_conn.go:14-29`: `MkdirAll` + `sql.Open` initializes). Blindly opening every catalog repo would litter empty pipeline DBs.
- **D7:** Repo discovery switches from filesystem scan to catalog: `store.OpenDefault()` → `db.ListRepos()`; each `ScanRecord.Slug` passes through `pipelinedb.SanitizeRepoSlug` (`cli/pipelinedb/pipeline_split_db.go:21`) before use as the pipeline DB key (filesystem walkers already saw sanitized names; catalog slugs are raw).
- **D8:** Per-repo `pipeline errors` / `pe` behavior is untouched. Only the `all` path changes.

## Data model changes

`AllPipelineSummary` (in `pipeline_all_errors.go`) gains:

```go
ReposScanned   int      `json:"reposScanned"`          // catalog repos visited
NoDataCount    int      `json:"noDataCount"`           // repos with no pipeline DB
NoDataRepos    []string `json:"noDataRepos,omitempty"` // slugs with no pipeline data
TotalErrorEntries int   `json:"totalErrorEntries"`     // sum of error entries across failed runs
```

Repo inspection becomes tri-state instead of `(item, isClean)`: **has pipeline data / clean / failed**. Booleans use positive prefixes (`hasPipelineData`, `isClean`, `hasErrors`) per the positive-boolean rule.

## Terminal view (D2)

Keep the existing header style, extended:

```
● Pipeline Error Inspector: ALL REPOSITORIES
  Repos Scanned:            <N>      (from catalog ListRepos)
  Passing / Clean:          <N>
  Failing:                  <N>
  No pipeline data:         <N>

  Repo: <slug>
    Errors: <count>
    └── ✖ <workflow> (Run #<id>, branch, sha, updated)
        Error: <message snippet>
        ↳ Inspect: gitmap pe <slug>

  ... per failed repo ...

  No pipeline data (<N>): slug-a, slug-b, …

  Pull errors notice (existing renderPullErrorsNotice)
```

JSON output keeps the existing shape, extended with the new fields above.

## `--file` report format (AI-readable markdown)

Written byte-faithful via `os.WriteFile(path, []byte(markdown), 0644)` — raw bytes, **no glyph filtering/transliteration** (per the `gitmap cat` TERM=dumb lesson). Suggested skeleton (worker may refine):

```markdown
# Pipeline Errors — All Repositories

- Generated: <RFC3339 UTC>
- Repos scanned: N
- Repos with errors: N
- Total error entries: N
- Repos without pipeline data: N

## <repo-slug>
- Workflow: <name> | Run: #<id> | Branch: <branch> | SHA: <sha> | Updated: <ts>
- Error count: <N>
- Suggested inspect: `gitmap pe <slug>`

### Errors
1. <message snippet>
...

## No pipeline data
- <slug>
- ...

## Pull errors
- ...
```

On write failure the command returns an `*apperror.AppError` and prints nothing misleading to the terminal. On success, terminal prints the summary and a confirmation line with the file path (same style as `writeContentToFile` in `cli/cmdpipeline/pipeline_helpers.go:28`).

## `te` shortcut

Add `"te"` to the alias list of the pipeline-errors dispatch entry in `cli/cmd/rootutility.go:453`:

```go
{[]string{"te", "pe", "pipeline-errors", "pipeline_errors", "ee"}, func() error { return cmdpipeline.RunPipelineErrors(argsTail()) }},
```

Then `gitmap te all` flows: `te` → `RunPipelineErrors(["all"])` → `runPipeline(["errors", "all"])` → `IsAll` → all-repos path.

## Acceptance criteria

1. `gitmap pipeline errors all` scans every catalog repo (count in header equals `ListRepos()` count), shows per-repo error groups + combined summary.
2. Repos with no pipeline DB appear in a "no pipeline data" section, and no new `sql.db` files are created for them.
3. `gitmap pipeline errors all --file /tmp/errors.md` writes a single AI-readable markdown file; terminal still shows the summary.
4. `gitmap pipeline errors all --file /tmp/errors.json --json` writes the JSON report to the file.
5. `gitmap te all` behaves exactly like `gitmap pe all`.
6. Per-repo `gitmap pipeline errors` / `gitmap pe <slug>` output unchanged.
7. Existing tests in `pipeline_all_errors_test.go` still pass (updated only if the discovery change breaks their assumptions).
8. New code: ≤300-line files, lowercase filenames, positive booleans, `*apperror.AppError` errors.

## Out of scope

- Per-repo error rendering changes.
- Auto-fix / agy fix integration (existing `--fix` flows untouched).
- Pushing the release (lead owns the minor bump + release step).
