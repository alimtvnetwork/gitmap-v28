## Problem

CI failure in `gitmap/helptext/examples_golden_test.go` → `TestEveryHelpFileHasExamples`. The gate requires each command help file to contain a `## Examples` heading followed somewhere below by a fenced code block (```` ``` ````). Four files fail:

- `hd.md`: no `## Examples` heading at all.
- `list-update.md`, `update-all.md`, `update-apply.md`: have `## Examples` headings, but their code samples are 4-space-indented blocks; the gate specifically checks for triple-backtick fences after the heading, so they fail even though example content exists.

## Fix

Bring all four files into compliance with the gate, without rewriting their content or restructuring the docs site.

1. `gitmap/helptext/hd.md`: add a new `## Examples` section (just before `## Exit Codes`) with two short fenced snippets covering the human view and `--json` output. Reuse literals already shown elsewhere in the file so the doc stays consistent.

2. `gitmap/helptext/list-update.md`, `update-all.md`, `update-apply.md`: convert the first code sample under each existing `## Examples` heading from a 4-space indent to a triple-backtick fence (```` ```text ```` for shell transcripts, ```` ```json ```` for JSON payloads). One fenced block per file is sufficient to satisfy the gate; leave the remaining indented samples untouched to minimize diff noise.

## Validation

- Re-run `go test ./gitmap/helptext/... -count=1` and confirm `TestEveryHelpFileHasExamples` passes.
- Run `go test ./... -count=1` at repo root to confirm no collateral regressions.

## Non-goals

- Not converting every indented sample to fenced style project-wide (separate cleanup).
- No version bump or changelog entry: this is a doc/test fix inside the already-shipped v6.80.1 cycle, not a user-visible change.

## Files touched

```text
gitmap/helptext/hd.md              (add ## Examples with fenced blocks)
gitmap/helptext/list-update.md     (fence first example)
gitmap/helptext/update-all.md      (fence first example)
gitmap/helptext/update-apply.md    (fence first example)
```
