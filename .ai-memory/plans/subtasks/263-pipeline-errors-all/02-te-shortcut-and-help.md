# Subtask 02 — `te` shortcut dispatch + help coverage (263)

## Goal
`gitmap te` routes to the pipeline-errors command group so that
`te all` ≡ `pe all` ≡ `pipeline errors all`, flags included
(`te all --file <path>`, `te all --json`, `te --help`, `te -1`, …).
Existing `pe` / `pipeline errors` / `pipeline_errors` / `ee` behavior stays untouched.

## Context the worker must know
- `te` is FREE (verified 2026-10-09: `gitmap te` → "Command 'te' is not recognized").
- Registration pattern at `cli/cmd/rootutility.go:453`:
  `{[]string{"pe", "pipeline-errors", "pipeline_errors", "ee"}, func() error { return cmdpipeline.RunPipelineErrors(argsTail()) }}`
- `RunPipelineErrors(args)` = `runPipeline(append([]string{"errors"}, args...))`
  (`cli/cmdpipeline/exports.go:17`), so any new dispatch entry pointing at
  `cmdpipeline.RunPipelineErrors(argsTail())` inherits `--help`, `all`, `--file`,
  `--json` handling automatically — no flag plumbing needed.
- `pe` help is a **Go builder** (`cli/cmdpipeline/pipeline_help_menu.go`:
  `RenderPipelineHelp()` → `buildPipelineHelpMenu()`, rendered via `termout.HelpMenu`).
  There is no `helptext/*.md` directory for pipeline help — do not create one.
- Shell completion: `makeTopLevelPECmd()` in `cli/cmd/root_cobra_completion.go:376`
  registers `Use: "pe [path|alias|url] [flags]"` with
  `Aliases: []string{"pipeline-errors", "pipeline_errors", "ee"}` (completion-only;
  real dispatch is `rootutility.go`).

## Step 1 — Add the `te` dispatch entry (one line)
In `cli/cmd/rootutility.go`, immediately below the `pe` entry (~line 453), add:

```go
{[]string{"te"}, func() error { return cmdpipeline.RunPipelineErrors(argsTail()) }},
```

Do NOT edit the existing `pe` entry. Do NOT add `"te"` to the existing alias list —
the mirror-entry form is what the owner approved (keeps `te` independently removable
and greppable). Do NOT touch any other dispatch entries.

## Step 2 — Verify identical routing
Rebuild (`go build ./...` is lead-authorized per the parent plan), then run each
command pair against the local catalog and confirm byte-identical behavior
(same stdout modulo terminal glyphs, same file output, same exit code):

| # | `pe` form | `te` form | Must match |
| - | --------- | --------- | ---------- |
| 1 | `gitmap pe --help` | `gitmap te --help` | same help menu |
| 2 | `gitmap pe all` | `gitmap te all` | same terminal summary |
| 3 | `gitmap pe all --file /tmp/pe-all.md` | `gitmap te all --file /tmp/te-all.md` | files identical (run after `--file` wiring from subtask 01 lands; `diff` the two files) |
| 4 | `gitmap pe all --json` | `gitmap te all --json` | same JSON |

Also sanity-run `gitmap te -1` (routes to the per-repo path, not `all`) to prove the
shortcut carries the full argument tail, not just the `all` subcommand.

## Step 3 — Help coverage (DRY: reference, never copy)
`pe` help is the Go builder in `cli/cmdpipeline/pipeline_help_menu.go` — add `te`
coverage there the same way, WITHOUT duplicating the menu:

- In `buildPipelineHelpMenu().UsageLines`, add one line referencing `pe`, e.g.:
  `"gitmap te all [flags]  (shortcut: same as 'gitmap pe all [flags]')"`
- Optionally one `Tips` entry, e.g.:
  `"Use shortcut 'gitmap te' anywhere 'gitmap pe' works — identical behavior."`
- Do NOT copy command tables, footer flags, or descriptions into a `te` variant.
- The `te --help` output comes from the same `RenderPipelineHelp()` call, so no
  separate help handler is needed — verify this in Step 2 row 1 instead of
  writing new help code.

Shell completion parity (one-word change): in
`cli/cmd/root_cobra_completion.go` `makeTopLevelPECmd()`, extend the alias list to
`[]string{"pipeline-errors", "pipeline_errors", "ee", "te"}` — same pattern the
existing `ee` alias already follows. This is completion-only; confirm it does not
alter dispatch behavior (re-run Step 2 row 2 after the change).

## Step 4 — Regression guard: `pe` untouched
- `git diff` (review-only, never commit — the parent handles commits) must show no
  changes to the existing `pe` dispatch entry, `RunPipelineErrors`,
  `renderAllPipelineSummaryTerminal`, or `buildPipelineHelpMenu` beyond the additive
  lines from Steps 1 and 3.
- Re-run `gitmap pe all` and `gitmap pe --help` after all edits: output identical
  to before (apart from the new `te` usage/tip lines).
- Do not add new test files for this subtask unless the parent plan's wave 3
  authorizes them; routing is verified by the live invocations in Step 2.

## Done when
- `te` dispatches in `rootutility.go` (one added line).
- All four verification pairs in Step 2 match.
- Help menu mentions `te` by reference; cobra aliases include `te`.
- `pe`/`pipeline errors` behavior byte-identical to before.
