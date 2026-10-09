# 263 — AI Report Format: `pipeline errors all --file`

## 1. Purpose

`gitmap pipeline errors all --file <path>` (also `pe all --file <path>`, `te all --file <path>`)
writes the combined cross-repo pipeline error report to a SINGLE markdown file that an AI
model can read and fix all reported issues from. The terminal still prints the usual
summary table when `--file` is used (flag only *adds* file output; it never replaces it).

This differs deliberately from `pe history-ai`: that command writes dual md+json
(`pipeline_history_ai.go` → `resolveHistoryAIOutputPaths`). For 263 the owner asked for ONE
file — markdown only. No companion `.json` sidecar.

## 2. Invariants (non-negotiable)

### 2.1 Traceback preservation
`02-spec/21-app/07-pipeline-and-diagnostics/01-architecture-spec.md` §3.1:
a report must NEVER be just `FAIL: Step #4`. Every per-error block MUST include:
1. the failing test name (e.g. `FAIL: test_gofmt_check_clean_repo`),
2. source file + line coordinates (`File "...", line 171`),
3. the exact assertion/exception message (`AssertionError: 1 != 0`).

If a failing run yields no extractable traceback, the block MUST say so explicitly
(`- **Traceback:** not available — raw log excerpt only`) rather than printing a bare
summary line.

### 2.2 One file, not two
`flags.FilePath` (parsed in `pipeline_flags.go:62` via `extractFlagVal(args, "--file")`) is
written exactly as given. The worker does NOT append `.json` or derive a second path.
(Contrast `history-ai`, which emits `<path>` + `<path>.json`.)

### 2.3 Byte-faithful write
Write with `os.WriteFile(path, []byte(content), 0644)` on the raw built string.
The content must NOT pass through glyph/emoji transliteration (the `gitmap cat`
`TERM=dumb` lesson: `→`→`->`, `⚠`→`▲` style mangling must never reach the file).
If the terminal summary uses colorized glyphs, strip or bypass them in the file
content — build the markdown separately from terminal-rendered strings.

### 2.4 Terminal still shows the summary
When `--file` is present: write the file first (return the write error if it fails),
THEN render the existing terminal summary (`renderAllPipelineSummaryTerminal`) unchanged.
`--json` + `--file` together: `--json` wins for terminal output semantics per existing
flag precedence; the file is still written (terminal summary is skipped only under
`--json`, which already exists today).

## 3. Data sources

- **Repo universe (D1):** `store.OpenDefault()` + `db.ListRepos()` → `[]model.ScanRecord`
  (Slug + AbsolutePath). Every catalog repo appears in the report; repos without
  pipeline data appear as "no pipeline data" rows in the trailing table (§4.4) —
  never silently skipped.
- **Per-repo run:** `pipelinedb.OpenPipelineSplitDb(slug)` →
  `QueryRunByNegativeOffset(-1)` for the latest run (reuses `inspectSinglePipelineRepo`).
- **Error records:** `QueryDetailedErrorLogsByRunId(runId)` → `PipelineErrorRecord`
  fields: `RunId`, `RepoSlug`, `WorkflowName`, `StepName`, `ErrorText`, `RawLogs`,
  `CreatedAt`. `ErrorText` = full error text; `RawLogs` = log excerpt source.
- **Run meta:** `PipelineRunRecord` → `WorkflowName`, `RunId`, `Branch`, `Sha`,
  `UpdatedAt`. Repo absolute path comes from the catalog record, not the DB slug.

If a catalog repo's split DB is missing or unreadable, its section records
`- **Status:** pipeline DB unavailable` (not an error block, not a silent skip).

## 4. Exact markdown structure

### 4.1 Document header (always first)

```markdown
# Pipeline Errors — All Repositories (AI Fix Dossier)

- **Generated At:** <UTC RFC3339, e.g. 2026-10-09T12:34:56Z>
- **Command:** `gitmap pipeline errors all --file`
- **Repos Scanned:** <N>
- **Repos With Errors:** <N>
- **Total Errors:** <N>

---
```

### 4.2 Per-repo sections

One `##` section per repo WITH errors (repos sorted alphabetically by slug):

```markdown
## Repo: <slug>

- **Repo Slug:** `<slug>`
- **Repo Path:** `<absolute path from catalog>`
- **Branch:** `<branch>`
- **Latest Run:** `<run id>` (<workflow name>, SHA `<sha>`, <UTC RFC3339>)
- **Status:** FAILING (<N> error(s))

### Error <n>

- **Workflow:** `<workflow name>`
- **Run ID:** `<run id>`
- **Branch:** `<branch>`
- **SHA:** `<sha>`
- **Failed Step:** `<step name>` (job: `<job name>` where available)
- **Timestamp:** `<UTC RFC3339>`
- **Error Summary:** `<one-line summary>`

**Full Error Text:**

```text
<complete ErrorText — never truncated mid-traceback>
```

**Log Excerpt:**

```text
<bounded RawLogs window around the failure; truncation marked
[log excerpt truncated: showing first/last <k> lines of <total>]>
```

---
```

- Repos appear in alphabetical slug order; errors within a repo in record order.
- Field labels are fixed bold strings (`**Repo Path:**`, `**Failed Step:**`, …) so the
  file is both human-readable and machine-parseable (see §6).

### 4.3 Field rules per error block

- **Repo Path:** absolute path from the catalog `ScanRecord` (worker has it from §3).
- **Step:** `PipelineErrorRecord.StepName`; include the parent job name when
  `QueryPipelineJobs(runId)` has it (`Failed Step: <step> (job: <job>)`).
- **Timestamp:** prefer the error record's `CreatedAt`; fall back to the run's
  `UpdatedAt`. Always UTC RFC3339.
- **Error Summary:** one line — the same `extractRunErrorSummary` value shown in
  the terminal.
- **Full Error Text:** complete `ErrorText`. If it is empty, print
  `(no extracted error text — see log excerpt below)` instead of an empty fence.
- **Log Excerpt:** bounded window from `RawLogs` (worker picks a sane cap, e.g.
  200 lines), with an explicit truncation marker when cut. Never cut mid-traceback:
  prefer cutting from the head, keeping the failure at the tail.

### 4.4 Trailing summary table

After the last error section, list every repo WITHOUT errors so nothing is hidden:

```markdown
## Repos Without Errors

| Repo | Path | Status |
| ---- | ---- | ------ |
| <slug> | <absolute path> | clean |
| <slug> | <absolute path> | no pipeline data |
| <slug> | <absolute path> | pipeline DB unavailable |
```

`clean` = latest run success. `no pipeline data` = DB opened but no latest run.
`pipeline DB unavailable` = DB missing/unreadable.

### 4.5 Pull errors appendix (only when present)

`queryPullErrorsForAll()` results get a trailing section only if non-empty:

```markdown
## Recorded Pull Errors (<N>)

- `<repo or scope>` — <record summary> (<timestamp>)
```

## 5. Example block

```markdown
# Pipeline Errors — All Repositories (AI Fix Dossier)

- **Generated At:** 2026-10-09T12:34:56Z
- **Command:** `gitmap pipeline errors all --file`
- **Repos Scanned:** 9
- **Repos With Errors:** 1
- **Total Errors:** 2

---

## Repo: gitmap-v28

- **Repo Slug:** `gitmap-v28`
- **Repo Path:** `/home/hatch/workspace/repos/gitmap-v28`
- **Branch:** `main`
- **Latest Run:** `18422001107` (ci.yml, SHA `8c9259e`, 2026-10-09T11:45:02Z)
- **Status:** FAILING (2 error(s))

### Error 1

- **Workflow:** `ci.yml`
- **Run ID:** `18422001107`
- **Branch:** `main`
- **SHA:** `8c9259e`
- **Failed Step:** `go test ./...` (job: `test`)
- **Timestamp:** 2026-10-09T11:45:02Z
- **Error Summary:** `FAIL: test_gofmt_check_clean_repo`

**Full Error Text:**

```text
--- FAIL: test_gofmt_check_clean_repo (0.01s)
    main_test.go:171: files not gofmt-clean: cli/cmd/root.go
FAIL
```

**Log Excerpt:**

```text
=== RUN   test_gofmt_check_clean_repo
--- FAIL: test_gofmt_check_clean_repo (0.01s)
    main_test.go:171: files not gofmt-clean: cli/cmd/root.go
FAIL	github.com/alimtvnetwork/gitmap-v28/cli	0.042s
```

---

## Repos Without Errors

| Repo | Path | Status |
| ---- | ---- | ------ |
| alim-seo-writing | /home/hatch/workspace/repos/alim-seo-writing | clean |
| wp-exam-v2 | /home/hatch/workspace/repos/wp-exam-v2 | no pipeline data |
```

## 6. Machine-parseability rules

- Heading levels are fixed: `#` title once, `## Repo: <slug>` per failing repo,
  `## Repos Without Errors` and `## Recorded Pull Errors (<N>)` trailing,
  `### Error <n>` per block.
- Field labels are fixed bold strings terminated by a colon; values follow on the
  same line. Slugs and paths are backticked.
- Fenced ` ```text ` blocks delimit full error text and log excerpts — an AI
  parser can split on fences without heuristics.
- Timestamps are UTC RFC3339 everywhere.
- `---` rules separate the header, each repo section, and each error block.

## 7. Error cases

- `--file` path unwritable: return the `os.WriteFile` error immediately; do NOT
  print the terminal summary as if all were well (fail closed on the requested file).
- Zero repos with errors: file still written — header shows `Repos With Errors: 0`,
  `Total Errors: 0`, no `## Repo:` sections, only the trailing table + pull-error
  appendix (if any). Terminal keeps its existing all-green message.
- Formatter lives in its own new file under `cli/cmdpipeline/` (≤300 lines per
  repo convention), invoked from `executeAllPipelineErrorLogs` when
  `flags.FilePath` is non-empty (§2.4 ordering applies).
