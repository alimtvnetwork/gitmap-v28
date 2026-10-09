# Subtask 01 — stdout output package (phases 1+2)

> Spec: `02-spec/21-app/265-stdout-ci-gates/02-stdout-refactor.md` (phases 1–2).
> Builds the new `cli/output/` package and demotes `termout`/`glyphs` Install
> to thin delegates. No behavior change for users yet — same filters, new
> machinery.

## Checkboxes — phase 1: new `cli/output` package

- [ ] Create `cli/output/output.go` with `FilterWriter`: an `io.Writer`
      wrapping a destination `*os.File` that applies the theme filter chain
      then the glyph filter chain synchronously in `Write`. NO pipes, NO
      goroutines, NO Drain.
- [ ] `output.Install(themeMode, glyphMode)`: single global-swap
      implementation absorbing the pipe/forwarder/Drain machinery currently
      duplicated in `cli/termout/install.go` and `cli/glyphs/install.go` —
      one registry, one Drain, theme-then-glyphs composition order preserved
      (matches `cli/cmd/root.go:84-93` ordering).
- [ ] `output.Raw()`: returns the unfiltered original `os.Stdout`/`os.Stderr`
      captured at package init, before any swap.
- [ ] `output.Writer()`: returns the installed `*FilterWriter` for the
      dispatch context.
- [ ] Zero-cost passthrough preserved: bright theme + rich glyphs installs
      nothing (same early-return as `termout.Install` today).
- [ ] Unit test: `Write` applies theme SGR rewrite then glyph safe-mode
      rewrite in order; bright+rich mode writes bytes untouched.

## Checkboxes — phase 2: delegate

- [ ] `termout.Install` / `glyphs.Install` become thin delegates to
      `output.Install` (or collapse to a single call from `Run()`); delete
      both `installedPipe` registries and both `Drain` functions.
- [ ] Deprecated entry points emit a one-time stderr warning, e.g.
      `output: termout.Install is deprecated; use output.Install`
      (grep-able signal for the future CI gate, workstream 2).
- [ ] `runDispatch` (`cli/cmd/root.go:166`) defers the single
      `output.Drain` for legacy call sites still on the pipe during the
      migration window.
- [ ] Benchmark (workstream 6): commit before/after throughput numbers for
      large writes through `FilterWriter` vs the old pipe+goroutine path.

## Acceptance

- Colors still apply to UI output in non-bright modes; glyph filtering still
  applies in safe mode (`TERM=dumb`); bright+rich is a zero-cost passthrough.
- `go build ./...` clean; `go test ./cli/output/...` green.
- `cli/cmd/root_no_args_test.go` still passes.
