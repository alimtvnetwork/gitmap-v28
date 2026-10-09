# Subtask 04 — stdout writer refactor

> Spec: `02-spec/21-app/250-package-consolidation/03-stdout-and-stacktrace.md`
> ("Stdout writer refactor" section). **Highest-risk subtask in program 250**
> — it touches the output path every command depends on. Sequenced AFTER the
> package merges.

## Checkboxes — phase 1: new `cli/output` package

- [ ] Create `cli/output/output.go` with `FilterWriter`: an `io.Writer`
      wrapping a destination `*os.File` that applies the theme filter chain
      then the glyph filter chain synchronously in `Write`. NO pipes, NO
      goroutines, NO Drain.
- [ ] `output.Install(themeMode, glyphMode)`: the single global-swap
      implementation absorbing the pipe/forwarder/Drain machinery currently
      duplicated in `cli/theme/install.go` and `cli/glyphs/install.go` — one
      pipe registry, one Drain, theme-then-glyphs composition order preserved.
- [ ] `output.Raw()`: returns the unfiltered original `os.Stdout`/`os.Stderr`
      captured at package init, before any swap.
- [ ] Zero-cost passthrough preserved: bright theme + rich glyphs installs
      nothing.

## Checkboxes — phase 2: delegate

- [ ] `theme.Install` / `glyphs.Install` become thin delegates to
      `output.Install` (or are removed in favor of a single call); delete both
      `installedPipe` registries and the two Drain functions.
- [ ] `runDispatch` (`cli/cmd/root.go:168`) defers the single `output.Drain`
      for legacy call sites still on the pipe.

## Checkboxes — phase 3: dispatch context + byte-faithful bypass

- [ ] `Run()` (`cli/cmd/root.go:66`) builds the writer after flag stripping;
      `runDispatch` carries it in the dispatch context.
- [ ] Replace the `byteFaithfulCommands` Install-skip (`root.go:53-64`) with an
      explicit bypass: the cat/view/type implementations (routed via
      `coreDispatchEntries()` in `cli/cmd/rootcore.go`) write file bytes
      through `output.Raw()`.
- [ ] Audit those commands: file bytes MUST go through `output.Raw()` — any
      bare `fmt.Print` of file bytes would still be filtered while the global
      swap is active.

## Checkboxes — phase 4: cliexit migration

- [ ] `cliexit.Reportf` / `HandleError` / `Fail` write via the dispatch-context
      writer (synchronous `FilterWriter`), not the pipe-wrapped `os.Stderr`.
- [ ] Remove the `cliexit.RegisterFlusher(theme.Drain)` /
      `RegisterFlusher(glyphs.Drain)` calls in `Run()` for the error path —
      synchronous writes cannot lose bytes before `os.Exit`.
- [ ] Pass the explicit writer to `RenderErrorSuggestions`
      (`cli/cmd/rootsuggestion.go:33`).

## Checkboxes — phase 5: follow-up documentation

- [ ] Document the follow-up in code comments: new code uses the explicit
      writer from the dispatch context; legacy `fmt.Print*` / `os.Stdout`
      call sites migrate opportunistically; the pipe mechanism in
      `output.Install` is deleted only when zero call sites depend on the
      global swap. NOT built in 250.

## Acceptance criteria (behavior contract)

- Colors work: theme filter applies to UI output in non-bright modes.
- Glyph filtering applies to UI chrome in safe mode (`TERM=dumb`).
- `cat` / `view` / `type` output is byte-identical to the source file (no
  emoji→ASCII rewrites, no SGR injection in file bytes).
- No lost bytes on exit: a failure message written immediately before
  `os.Exit` is fully flushed.
- Bright theme + rich glyphs remains a zero-cost passthrough.
- `cli/cmd/root_no_args_test.go` (the `os.Stdout`-swapping test) still passes.

## Verification rule

**Verify with a fresh binary after every phase** — build, then exercise:
colored UI output, `TERM=dumb` safe-mode glyph filtering, `cat` on a file
containing emoji/ANSI bytes (byte-compare the output against the original),
and an error path that calls `cliexit.Fail` immediately before `os.Exit`.
Do not batch phases without an intervening verification.
