# Subtask 02 — stdout dispatch migration (phases 3+4+5)

> Spec: `02-spec/21-app/265-stdout-ci-gates/02-stdout-refactor.md` (phases 3–5).
> Depends on subtask 01 (`cli/output` package exists). Wires the writer
> through dispatch, adds the byte-faithful bypass, migrates the error path,
> and verifies with a fresh binary.

## Checkboxes — phase 3: dispatch context + byte-faithful bypass

- [ ] `Run()` (`cli/cmd/root.go:66`) builds the writer via `output.Install`
      after flag stripping (same position as today's `termout.Install` at
      root.go:84 / `glyphs.Install` at root.go:93); `runDispatch`
      (`cli/cmd/root.go:166`) carries it in the dispatch context.
- [ ] Remove the `byteFaithfulCommands` Install-skip map
      (`cli/cmd/root.go:58-60`); replace with an explicit bypass at the call
      site: `printCatContent` in `cli/cmdmacro/macro_add_file_ops.go` writes
      the banner through the filtered writer and the file bytes through
      `output.Raw()`. `view`/`type` share `runCatCmd` — one change covers
      all three. Audit: no bare `fmt.Print` of file bytes remains.
- [ ] Do NOT touch the shadowed raw `RunCat` (`cli/cmdmacro/cat.go:10`,
      registered at `cli/cmd/roottooling.go:335`, unreachable — dispatchUtility
      shadows it). Removal is a follow-up, not this task.
- [ ] Document the ordering contract with `cli/cmdnodes/nodes_clone.go:348`
      (`captureOutput` swaps os.Stdout/os.Stderr in production): `Install`
      runs once at process start; captures restore the handles `Install`
      swapped. Do not re-swap inside a capture.
- [ ] Leave the ~15 test-only os.Stdout-swapping test files alone.

## Checkboxes — phase 4: cliexit migration

- [ ] `cliexit.Reportf` / `HandleError` / `Fail` write via the
      dispatch-context writer (synchronous `FilterWriter`), not the
      pipe-wrapped `os.Stderr`.
- [ ] Remove `cliexit.RegisterFlusher(termout.Drain)` /
      `RegisterFlusher(glyphs.Drain)` (`cli/cmd/root.go:101-102`) for the
      error path — synchronous writes cannot lose bytes before `os.Exit`.
- [ ] Pass the explicit writer to `RenderErrorSuggestions`
      (`cli/cmd/rootsuggestion.go:33`).

## Checkboxes — phase 5: docs

- [ ] Code comments record the direction: new code uses the explicit writer
      from the dispatch context; legacy `fmt.Print*` / `os.Stdout` call
      sites migrate opportunistically; the pipe mechanism in
      `output.Install` is deleted only when zero call sites depend on the
      global swap. NOT built in 265.

## Verification (fresh binary after every phase — do not batch)

1. `TERM=dumb cat` byte-compare: create a fixture file with emoji + ANSI
   bytes, run the fresh binary's `cat`, `diff` output payload against the
   fixture — must be byte-identical.
2. Color check: run a UI-heavy command in a non-bright theme; confirm SGR
   colors apply to UI output.
3. `--glyphs safe` UI filter check: run a UI-heavy command with safe glyphs
   on a dumb terminal; confirm emoji in UI chrome are rewritten while file
   payloads are untouched.
4. `cliexit.Fail` flush check: trigger an error path that calls
   `cliexit.Fail` immediately before `os.Exit`; confirm the full failure
   message reaches the terminal (the Windows lost-bytes case).
5. `cli/cmd/root_no_args_test.go` still passes; `go build ./...` clean.
