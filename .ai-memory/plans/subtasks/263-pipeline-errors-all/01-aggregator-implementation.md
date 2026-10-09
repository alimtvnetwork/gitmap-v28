# Subtask 263-pipeline-errors-all · 01 — Aggregator Implementation

Scope: worker implementer of `gitmap pipeline errors all` aggregation + `--file` wiring + `te` shortcut.
Spec home: `02-spec/21-app/263-pipeline-errors-all/01-overview.md` (read it first — it holds decisions D1–D8).
Rules for this subtask: files ≤300 lines, lowercase filenames, positive-prefixed booleans (`hasX`, `isX`), `*apperror.AppError` for errors, relative paths only. Do NOT run `go build`, `go test`, or any git commands — the lead verifies the build and owns commits/releases.

## Step 1 — Switch repo discovery from filesystem to catalog

File: `cli/cmdpipeline/pipeline_all_errors.go` (keep ≤300 lines; move formatting helpers out per Step 2 if needed).

1. Replace `discoverPipelineRepoSlugs()` body (currently walks `<bindir>/data/pipeline/` subdirs, `pipeline_all_errors.go:77`) with a catalog-backed version:
   - `db, err := store.OpenDefault()` (`cli/store/location.go:71`); on error return an empty slice (all-repos must not crash when the catalog is unavailable; a later summary field can note the failure).
   - `defer db.Close()`.
   - `records, err := db.ListRepos()` (`cli/store/repo.go:91`); on error return an empty slice.
   - For each `model.ScanRecord` append `pipelinedb.SanitizeRepoSlug(r.Slug)` (`cli/pipelinedb/pipeline_split_db.go:21`). This is required: catalog slugs are raw, but pipeline DB dirs are keyed by sanitized slugs.
2. Keep the function name or rename to `discoverCatalogRepoSlugs` — either is fine; update `collectAllPipelineSummary()` accordingly.
3. Delete the now-unused `os`/`path/filepath`/`store.BinaryDataDir` references if the compiler complains about unused imports (lead verifies the build; keep imports tidy).

## Step 2 — Tri-state repo inspection (D1, D6)

`OpenPipelineSplitDb` (`cli/pipelinedb/pipeline_split_conn.go:14`) **creates** an empty DB via `MkdirAll` + `sql.Open` when none exists. Since discovery now comes from the full catalog, opening every repo blindly would create empty `sql.db` files everywhere. Prevent this:

1. Change `inspectSinglePipelineRepo(slug)` so it FIRST stats the file:
   ```go
   dbPath := pipelinedb.ResolvePipelineDbPath(slug)
   if _, statErr := os.Stat(dbPath); statErr != nil {
       return nil, false, false // failedItem, isClean, hasPipelineData=false
   }
   ```
2. Make inspection tri-state: `hasPipelineData` (DB file exists), `isClean` (latest run success), else failed with an `*AllPipelineFailedItem`. Boolean names stay positive (`hasPipelineData`, `isClean`) per the repo rule.
3. In `collectAllPipelineSummary()`, route each slug:
   - `hasPipelineData == false` → `summary.NoDataCount++`, append slug to `summary.NoDataRepos`.
   - `isClean` → `summary.CleanCount++`, append to `CleanRepos`.
   - failed item → `summary.FailedCount++`, append item to `FailedItems`; also accumulate the error count for that run into `summary.TotalErrorEntries` (count `db.QueryCompactErrorLogsByRunId(run.RunId).Data` entries — the same query `extractRunErrorSummary` already uses).
4. Add the new fields to `AllPipelineSummary` (JSON tags included):
   - `ReposScanned int` (set to number of catalog slugs visited),
   - `NoDataCount int`, `NoDataRepos []string \`json:"noDataRepos,omitempty"\``,
   - `TotalErrorEntries int`.

## Step 3 — Wire `PipelineErrorFlags.FilePath` into the all-repos path (D4)

File: new `cli/cmdpipeline/pipeline_all_errors_format.go` (≤300 lines; holds ALL new rendering: markdown builder + extended terminal render + file write helper).

1. In `executeAllPipelineErrorLogs`, after `collectAllPipelineSummary()`:
   - If `flags.IsJSON` AND `flags.FilePath == ""` → existing behavior: print indented JSON to stdout, return.
   - If `flags.FilePath != ""`:
     - If `flags.IsJSON` → marshal the summary (indented, same as `emitAllPipelineSummaryJSON`) to bytes and write to `flags.FilePath`.
     - Else → build the AI-readable markdown report (see skeleton below) and write it to `flags.FilePath`.
     - Write via `os.WriteFile(path, content, 0644)` — **BYTE-FAITHFUL: write raw bytes, no glyph filtering or transliteration** (the `gitmap cat` TERM=dumb lesson: never map Unicode to ASCII fallbacks; the bytes that build the string are the bytes that land in the file).
     - Create parent dirs first (`os.MkdirAll(filepath.Dir(path), 0755)`), mirroring `writeContentToFile` (`cli/cmdpipeline/pipeline_helpers.go:28`).
     - On write failure return `apperror.WrapSimple(err, "write all pipeline errors report")`.
     - On success print the confirmation line `✓ Output written to <path>` (same style as `writeContentToFile`) — then fall through to the terminal render.
   - Always run `renderAllPipelineSummaryTerminal(summary)` afterwards (even when `--file` was given): **terminal still prints the summary** — `--file` only adds the file output.
2. Markdown report skeleton (refine wording freely, keep the structure AI-fixable):
   ```markdown
   # Pipeline Errors — All Repositories
   - Generated: <UTC RFC3339>
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
   ...

   ## Pull errors
   - <pull error records, if any>
   ```

## Step 4 — Keep `--json` and the existing terminal render working (D2)

1. `--json` to stdout stays byte-identical in shape, only extended by the new summary fields (Step 2.4).
2. Extend `renderAllPipelineSummaryTerminal` (move it into the new format file if `pipeline_all_errors.go` risks exceeding 300 lines):
   - Header adds: `Repos Scanned: <N>` and `No pipeline data: <N>`.
   - Per failed repo print the error count: `Repo: <slug> · Errors: <count>` then the existing workflow/error/inspect lines (keep `gitmap pe <slug>` hint).
   - After the failing list, print the "no pipeline data" slugs (truncate the list display if very long; the full list lives in `--json` and `--file`).
   - Keep the existing pull-errors notice (`renderPullErrorsNotice`) unchanged.

## Step 5 — Register the `te` shortcut

File: `cli/cmd/rootutility.go`, line ~453 in `utilityPipelineEntries()`. Change:

```go
{[]string{"pe", "pipeline-errors", "pipeline_errors", "ee"}, func() error { return cmdpipeline.RunPipelineErrors(argsTail()) }},
```

to:

```go
{[]string{"te", "pe", "pipeline-errors", "pipeline_errors", "ee"}, func() error { return cmdpipeline.RunPipelineErrors(argsTail()) }},
```

`gitmap te all` then flows: `te` → `RunPipelineErrors(["all"])` → `runPipeline(["errors", "all"])` → `flags.IsAll` → the all-repos path.

## Step 6 — Touch-up checklist

- [ ] No changes to any per-repo error rendering path (`executePipelineErrorLogs` non-`IsAll` branch untouched).
- [ ] No new top-level package; everything inside `cli/cmdpipeline/`.
- [ ] Existing tests `cli/cmdpipeline/pipeline_all_errors_test.go` updated only if the discovery change invalidates their assumptions.
- [ ] Errors surfaced as `*apperror.AppError` (use `apperror.WrapSimple(err, "op")`, `cli/apperror/apperror.go:454`); `executeAllPipelineErrorLogs` keeps its plain `error` signature.
- [ ] No `go build` / `go test` / git commands in this subtask — lead verifies build, runs tests, and owns the minor bump + release.
- [ ] File sizes: `pipeline_all_errors.go` ≤300 lines, `pipeline_all_errors_format.go` ≤300 lines.
