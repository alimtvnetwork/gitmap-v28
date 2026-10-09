# Task 265 — stdout CI gates

> Spec: task 265-stdout-ci-gates (program follow-up to program 250,
> `02-spec/21-app/250-package-consolidation/03-stdout-and-stacktrace.md`).

## User Request (Verbatim)

The 7 workstreams the user approved:

1. **stdout refactor** — remove the os.Stdout/os.Stderr pipe interception
   (the "goroutine forwarder" design) and replace it with a synchronous
   writer path that can never lose bytes before `os.Exit`.
2. **CI gates** — add CI checks that fail the build when new os.Stdout /
   os.Stderr mutation, new pipe-based output interception, or new
   background goroutine forwarders are introduced.
3. **Tests in CI** — stdout behavior tests (color application, safe-mode
   glyph filtering, byte-faithful cat, no-lost-bytes on error exit) run as
   part of the CI pipeline, not only locally.
4. **Deprecation policy** — a documented rule for retiring the legacy
   `termout.Install` / `glyphs.Install` pipe machinery and the
   `cliexit.RegisterFlusher(termout.Drain)` / `RegisterFlusher(glyphs.Drain)`
   calls once the new path is proven.
5. **Credential audit** — audit the output path for accidental credential /
   secret printing; the refactor must not change which bytes reach the
   terminal in a way that leaks anything (see also the `cli/secrets` package).
6. **Benchmarks** — measure output throughput before and after the refactor;
   the synchronous FilterWriter must not be slower than the pipe+goroutine
   design for large output (e.g. `cat` of a multi-MB file).
7. **pe-all parallel** — after the stdout path is stable, re-run the full
   parallel-execution (`pe-all`) test suite against a fresh binary built
   from the new path, because `cli/cmdnodes/nodes_clone.go:348`
   (`captureOutput`) mutates os.Stdout/os.Stderr in production and
   interacts with whatever the new Install does.

This spec (task 265) covers workstream 1 plus the design surface the other
workstreams depend on. Workstreams 2–7 are tracked as follow-ups in
`02-stdout-refactor.md` § "Follow-ups for workstreams 2–7"; each gets its
own spec file and subtask set when started.

## Decisions D1–D5

- **D1 — Phased migration, reusing the program-250 design.** The design in
  `.ai-memory/plans/subtasks/250-package-consolidation/04-stdout-writer-refactor.md`
  is sound and is reused with corrections (see "Corrections to the 250
  design" below). Five phases, verified one at a time with a fresh binary.
  Do not batch phases without an intervening verification.

- **D2 — Synchronous `FilterWriter`; no pipes, no goroutines, no Drain.**
  Filtering (theme SGR rewrite, then glyph filtering) happens inline in
  `Write`. This eliminates the entire lost-bytes-on-exit class: there is no
  forwarder goroutine that can be starved before `os.Exit`, so the
  `cliexit.RegisterFlusher` machinery becomes unnecessary for the error path.

- **D3 — Byte-faithful commands go through `output.Raw()`.** `cat`/`view`/
  `type` (and the shadowed raw `cmdmacro.RunCat` if ever re-activated) write
  file bytes through the pre-swap `*os.File` handles captured at package
  init — never through the filter. The current `byteFaithfulCommands`
  Install-skip (`cli/cmd/root.go:58-60`) is replaced by an explicit bypass
  at the call site: `printCatContent` in
  `cli/cmdmacro/macro_add_file_ops.go` writes the file payload via
  `output.Raw()`.

- **D4 — Deprecation = stderr warning pattern.** Deprecated entry points
  (`termout.Install`, `glyphs.Install`, the two `Drain` functions, and the
  `cliexit.RegisterFlusher` calls) keep working as thin delegates during the
  migration window and emit a one-time stderr warning
  (`output: termout.Install is deprecated; use output.Install`) pointing at
  the replacement. Hard removal waits for a follow-up task once zero call
  sites remain; the warning gives us a grep-able migration signal for the
  CI gate in workstream 2.

- **D5 — Dispatch context carries the writer.** `Run()` builds the writer
  after flag stripping (`--theme`/`--glyphs` stripped first, as today);
  `runDispatch` carries it in the dispatch context so `cliexit.Reportf`/
  `Fail`/`HandleError` write through the explicit writer instead of the
  globally swapped `os.Stderr`.

## Corrections to the 250 design (verified facts)

1. The theme package is `cli/termout` (not `cli/theme`).
   `termout.Install()` runs at `cli/cmd/root.go:84`; `glyphs.Install()` at
   `cli/cmd/root.go:93`, guarded by `byteFaithfulCommands` at
   `cli/cmd/root.go:58-60`.
2. Flush registration: `cliexit.RegisterFlusher(termout.Drain)` at
   `cli/cmd/root.go:101`, `cliexit.RegisterFlusher(glyphs.Drain)` at
   `cli/cmd/root.go:102`; `defer termout.Drain()` / `defer glyphs.Drain()`
   at `cli/cmd/root.go:172-173` inside `runDispatch` (`cli/cmd/root.go:166`).
3. `cat`/`view`/`type` route via `utilityDesktopEntries()` at
   `cli/cmd/rootutility.go:474` → `cmdautofix.RunCatCmd` →
   `runCatCmd`/`printCatContent` in
   `cli/cmdmacro/macro_add_file_ops.go` (decorated with a banner — NOT
   byte-faithful today). The raw `RunCat` (`cli/cmdmacro/cat.go:10`,
   `io.Copy(os.Stdout)`) is registered at `cli/cmd/roottooling.go:335` but
   is SHADOWED (dispatchUtility runs before dispatchTooling) — dead code
   path. Do not design around it; either keep it dead or remove it in a
   follow-up.
4. `cli/cmdnodes/nodes_clone.go:348` also mutates os.Stdout/os.Stderr in
   production (`captureOutput` capture/restore) — the new `output.Install`
   must not fight it; see `02-stdout-refactor.md` Phase 3.
5. ~15 test files swap os.Stdout — test-only; leave them alone. The spec
   notes `cli/cmd/root_no_args_test.go` as the canonical os.Stdout-swapping
   test that must keep passing.

## Architecture diagram

### Current output path (pipe + goroutine forwarder)

```
subcommand code
    fmt.Print* / os.Stdout writes
        │
        ▼
os.Stdout ──▶ [pipe writer] ──(goroutine)──▶ termout.Filter (SGR rewrite)
        │                                            │
        │ (glyphs wraps AFTER theme: two stacked pipes)▼
        ▼                                    glyphs.Filter (emoji→ASCII in safe mode)
os.Stderr ──▶ [pipe writer] ──(goroutine)──▶ original os.Stdout/os.Stderr fds
                                                        │
problems:                                               ▼
  • lost bytes before os.Exit on Windows (goroutine never   terminal
    scheduled) → needs Drain + RegisterFlusher
  • two duplicated installedPipe registries (termout + glyphs)
  • cat/view/type rely on a fragile Install-skip map
```

### Target output path (synchronous writer)

```
Run()
  strip --theme / --glyphs
  output.Install(themeMode, glyphMode)      ← single global swap
        │
        ▼
dispatch context carries output.Writer    ← explicit writer
        │
        ├── UI output ──▶ FilterWriter{dest: os.Stdout}
        │                   Write(): theme filter THEN glyph filter,
        │                   synchronously. No pipes, no goroutines, no Drain.
        │
        ├── byte-faithful (cat/view/type)
        │                   ──▶ output.Raw() → pre-swap *os.File handles
        │
        └── cliexit.Reportf/Fail/HandleError
                            ──▶ dispatch-context writer (sync write,
                                no flusher needed before os.Exit)
        │
        ▼
terminal

invariants:
  • bright theme + rich glyphs = zero-cost passthrough (installs nothing)
  • colors still apply to UI output in non-bright modes
  • glyph filtering still applies to UI chrome in safe mode (TERM=dumb)
  • cat byte-identical to the source file
  • failure message written just before os.Exit is fully flushed
```

## File map (this task writes ONLY these)

| File | Contents |
| ---- | -------- |
| `02-spec/21-app/265-stdout-ci-gates/01-overview.md` | this file |
| `02-spec/21-app/265-stdout-ci-gates/02-stdout-refactor.md` | phases 1–5, behavior contract, workstream 2–7 follow-ups |
| `.ai-memory/plans/subtasks/265-stdout-ci-gates/01-stdout-output-package.md` | subtask checkboxes, phases 1+2 |
| `.ai-memory/plans/subtasks/265-stdout-ci-gates/02-stdout-dispatch-migration.md` | subtask checkboxes, phases 3–5 + verification |
