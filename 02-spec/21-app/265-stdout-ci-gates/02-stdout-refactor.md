# Task 265 — stdout refactor spec

> Reuses `.ai-memory/plans/subtasks/250-package-consolidation/04-stdout-writer-refactor.md`
> with the corrections listed in `01-overview.md` (package is `cli/termout`,
> corrected line numbers, shadowed `RunCat`, `nodes_clone.go:348` mutation).
> Each phase is verified with a fresh binary before the next begins.

## Phase 1 — new `cli/output/` package

- `cli/output/output.go` — `FilterWriter`: an `io.Writer` wrapping a
  destination `*os.File` that applies the theme filter chain then the glyph
  filter chain **synchronously in `Write`**. NO pipes, NO goroutines, NO
  Drain. Composition order preserved: theme first (SGR rewrite), glyphs
  second (safe-mode ASCII rewrite), exactly as `cli/cmd/root.go:84-93`
  installs them today.
- `output.Install(themeMode, glyphMode)` — the single global-swap
  implementation absorbing the pipe/forwarder/Drain machinery currently
  duplicated in `cli/termout/install.go` and `cli/glyphs/install.go`: one
  registry, one Drain, one composition order. Kept as a global swap for
  legacy call sites (`fmt.Print*` everywhere) during the migration window.
- `output.Raw()` — returns the unfiltered original `os.Stdout`/`os.Stderr`
  captured at package init, before any swap. Byte-faithful commands write
  file bytes through these handles.
- `output.Writer()` — returns the installed `*FilterWriter` for the
  dispatch context (Phase 3).
- Zero-cost passthrough preserved: bright theme + rich glyphs installs
  nothing (same early-return as `termout.Install` today).

## Phase 2 — `termout` / `glyphs` become delegates

- `termout.Install()` and `glyphs.Install()` become thin delegates to
  `output.Install` (or collapse to a single call from `Run()`).
- Delete both `installedPipe` registries and both `Drain` functions —
  replaced by one registry + one Drain in `cli/output`.
- Emit a one-time stderr deprecation warning from the legacy entry points
  (D4 in `01-overview.md`): `output: termout.Install is deprecated; use
  output.Install`. The warning is the grep-able signal for the future CI
  gate (workstream 2).
- `runDispatch` (`cli/cmd/root.go:166`) defers the single `output.Drain`
  for any legacy call sites still on the pipe during the window.

## Phase 3 — `Run()` builds the writer; dispatch context + byte-faithful bypass

- `Run()` (`cli/cmd/root.go:66`) builds the writer after flag stripping
  (same position as today's `termout.Install` at root.go:84 and
  `glyphs.Install` at root.go:93). `runDispatch` carries it in the dispatch
  context so downstream code can take the explicit writer.
- Replace the `byteFaithfulCommands` Install-skip (`cli/cmd/root.go:58-60`)
  with an explicit bypass at the call site: `printCatContent` in
  `cli/cmdmacro/macro_add_file_ops.go` writes the decorative banner through
  the filtered writer (banner is UI chrome — glyph filtering may apply to
  it) and writes the **file bytes** through `output.Raw()`. Note the banner
  means `cat` is NOT byte-faithful today; the Raw bypass covers the payload.
  `view`/`type` share `runCatCmd`, so one change covers all three.
- Audit: file bytes MUST go through `output.Raw()` — any bare `fmt.Print`
  of file bytes would still be filtered while the global swap is active.
- Do not touch the shadowed raw `RunCat` (`cli/cmdmacro/cat.go:10`,
  registered at `cli/cmd/roottooling.go:335`, unreachable because
  dispatchUtility shadows it) — dead code; removal is a follow-up, not
  this task.
- Interaction with `cli/cmdnodes/nodes_clone.go:348`: `captureOutput`
  swaps `os.Stdout`/`os.Stderr` around a `fn()` in production
  (pe-all runs). The new `output.Install` must remain idempotent and must
  not re-swap while a capture is active; document the ordering contract
  (`Install` once at process start; captures restore the handles
  `Install` swapped). The pe-all parallel verification (workstream 7)
  exercises this.

## Phase 4 — `cliexit` migration

- `cliexit.Reportf` / `HandleError` / `Fail` write via the dispatch-context
  writer (the synchronous `FilterWriter`), not the pipe-wrapped `os.Stderr`.
- Remove the `cliexit.RegisterFlusher(termout.Drain)` /
  `RegisterFlusher(glyphs.Drain)` calls (`cli/cmd/root.go:101-102`) for the
  error path — synchronous writes cannot lose bytes before `os.Exit`, which
  is the whole point of D2.
- Pass the explicit writer to `RenderErrorSuggestions`
  (`cli/cmd/rootsuggestion.go:33`).

## Phase 5 — docs

- Code comments record the direction: new code uses the explicit writer
  from the dispatch context; legacy `fmt.Print*` / `os.Stdout` call sites
  migrate opportunistically; the pipe mechanism in `output.Install` is
  deleted only when zero call sites depend on the global swap. NOT built
  in 265 — tracked as a follow-up.
- ~15 test files that swap os.Stdout are test-only: leave them.

## Behavior contract (acceptance criteria)

1. Colors work: theme filter applies to UI output in non-bright modes.
2. Glyph filtering applies to UI chrome in safe mode (`TERM=dumb`).
3. `cat`/`view`/`type` file payloads are byte-identical to the source file
   (no emoji→ASCII rewrites, no SGR injection in file bytes). Banner lines
   are UI chrome and may be filtered.
4. No lost bytes on exit: a failure message written immediately before
   `os.Exit` is fully flushed.
5. Bright theme + rich glyphs remains a zero-cost passthrough.
6. `cli/cmd/root_no_args_test.go` (the os.Stdout-swapping test) still passes.
7. Benchmark: synchronous `FilterWriter` throughput ≥ pipe+goroutine design
   on large output (workstream 6 gates this).

## Follow-ups for workstreams 2–7 (not in this task)

- **WS2 CI gates**: fail CI on new `os.Stdout`/`os.Stderr` assignments in
  non-test code, new `os.Pipe()` in output paths, new forwarder goroutines,
  and new call sites of the deprecated `termout.Install`/`glyphs.Install`/
  `*.Drain` (grep for the D4 deprecation warning string).
- **WS3 tests in CI**: color, safe-mode glyph, byte-faithful cat, and
  fail-before-exit tests run in the CI pipeline.
- **WS4 deprecation policy**: hard-removal task once zero call sites
  remain; the D4 warning is the migration signal.
- **WS5 credential audit**: audit the new output path against the
  `cli/secrets` package rules; no behavior change in what reaches the
  terminal beyond filtering.
- **WS6 benchmarks**: before/after throughput numbers committed with the
  Phase 1 change.
- **WS7 pe-all parallel**: re-run the full `pe-all` suite on a fresh binary
  after Phase 3, covering the `nodes_clone.go:348` capture interaction.
