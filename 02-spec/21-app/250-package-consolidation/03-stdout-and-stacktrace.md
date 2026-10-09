# Spec 250 — Package consolidation: stdout writer refactor and stack-trace setting

> Part of task 250-package-consolidation. Research: 02 (verified).

## Stack-trace setting

### Problem

`handleGlobalError` (`cli/cmd/root.go:215-268`) prints a stack trace for every
unhandled error whenever `ErrorDisplay` is `"full"` (the default). There is no
independent toggle: the only way to suppress traces today is the coarse
`errorDisplay: "simple"` switch, which also changes the whole error format.

### Design

1. Add one boolean to `model.Config` (`cli/model/record.go`), directly after
   the `ErrorDisplay` field:

   ```go
   ErrorDisplay  string `json:"errorDisplay"`
   ShowStackTrace bool  `json:"showStackTrace"`
   ```

   Naming precedent: `CommitReplayKeepUrl` (`record.go:66`, default `false`
   at line 82).

2. Default `true` in `DefaultConfig()` (`record.go:74-93`). Default ON preserves
   current behavior exactly — `parseConfig` unmarshals onto `DefaultConfig()`,
   so existing config files without the key keep today's behavior with no
   migration.

3. Gate the trace print in `handleGlobalError` (`root.go:261-264`), composing
   with — not replacing — the existing `errorDisplay == "simple"` early exit
   (which stays as-is):

   ```go
   showStack := true
   if cfgErr == nil {
       showStack = cfg.ShowStackTrace
   }
   ...
   cliexit.Reportf(command, "execute", "", err)

   stack := resolveErrorStackTrace(err)
   if stack != "" && showStack {
       fmt.Fprintf(os.Stderr, "Stack Trace:%s\n", stack)
   }
   ```

   Resolve from `cfg` only when `cfgErr == nil`, mirroring how `display` is
   already resolved in the same function — a missing/unreadable config must
   fall back to today's behavior (traces on), not silently disable them.

### Decisions

- Default ON: preserves current behavior; opt-out is explicit.
- No `validate.go` change: a bool needs no validation, and the
  unmarshal-onto-defaults in `parseConfig` covers missing keys.
- The `errorDisplay == "simple"` early exit is untouched — the two controls
  are orthogonal (error format vs. trace detail).
- `RenderErrorSuggestions` (`cli/cmd/rootsuggestion.go:33`, the only
  writer-parameterized output site) is unaffected.

### Caveat — cwd-relative config path

`handleGlobalError` reads `./data/config.json` (via
`constants.DefaultConfigPath`) relative to the process working directory. The
new setting only applies when the invocation's cwd actually contains that
file; otherwise the default (traces on) applies. Do not "fix" the path in
this program — that is a separate concern.

## Stdout writer refactor

### Upfront honesty

Full elimination of the global `os.Stdout`/`os.Stderr` mutation means
threading an explicit writer through ~8,500 `fmt.Print*` call sites and
~2,000 direct `os.Stdout`/`os.Stderr` references across `cli/` — every one of
them resolves the globals at call time today. That is NOT feasible in one
program. What IS achievable in program 250:

1. Contain the interception mechanism in ONE package instead of two duplicated
   pipe implementations.
2. Give new and migrated code a synchronous, pipe-free explicit writer.
3. Put the writer in the dispatch context so call sites CAN take it explicitly.
4. Migrate the highest-value writer — `cliexit` error reporting — off the pipe.
5. Document the follow-up: gradual call-site migration.

The global swap stays for legacy call sites. `os.Stdout` is `*os.File` — a
concrete type, so a custom writer can never be assigned to it; the pipe is the
only way `fmt.Print*` output can be filtered. This program does not pretend
otherwise.

### Phased design

#### Phase 1 — new `cli/output` package

- `FilterWriter`: an `io.Writer` wrapping a destination `*os.File` that applies
  the theme filter chain then the glyph filter chain synchronously inside
  `Write`. NO pipes, NO goroutines, NO Drain. New code and migrated call sites
  use `fmt.Fprintf(w, ...)` against it.
- `Install(themeMode, glyphMode)`: the single global-swap implementation. The
  pipe/forwarder/Drain mechanism currently duplicated in
  `cli/theme/install.go` and `cli/glyphs/install.go` moves here — one pipe
  registry, one Drain, one filter-chain composition point (theme first, then
  glyphs, matching today's nesting order).
- `Raw()`: returns the unfiltered originals (`os.Stdout`/`os.Stderr` captured
  at package init, before any swap) for byte-faithful output.
- Keep the zero-cost passthrough: when theme mode is bright AND glyph mode is
  rich, install nothing (today's behavior, still skipped entirely).

#### Phase 2 — `theme.Install` / `glyphs.Install` delegate

- Both become thin wrappers over `output.Install` (or are removed in favor of
  a single call). The duplicated `installedPipe` registries and the two
  separate Drain functions are deleted; one `output.Drain` remains.
- `runDispatch` (`cli/cmd/root.go:168`) keeps its deferred Drain call(s), now
  against the single registry.

#### Phase 3 — writer in the dispatch context

- `Run()` (`cli/cmd/root.go:66`) builds the writer after flag stripping and
  stores it where dispatch can reach it (dispatch context passed to
  `runDispatch`).
- `byteFaithfulCommands` (`root.go:53-64`): replace the "skip
  `glyphs.Install`" special case with an explicit bypass — the cat/view/type
  implementations (routed via `coreDispatchEntries()` in
  `cli/cmd/rootcore.go`) write file bytes through `output.Raw()`. UI chrome in
  those commands may keep the filtered path; file bytes must not.
- Constraint: any byte-faithful command MUST route file bytes through
  `output.Raw()`; bare `fmt.Print` of file bytes in those commands would still
  be filtered while the global swap is active.

#### Phase 4 — migrate `cliexit` error reporting

- `cliexit.Reportf` / `HandleError` / `Fail` write through the
  dispatch-context writer (synchronous `FilterWriter`) instead of the
  pipe-wrapped `os.Stderr`.
- The `cliexit.RegisterFlusher(theme.Drain)` /
  `RegisterFlusher(glyphs.Drain)` calls in `Run()` go away for the error path
  — synchronous writes cannot lose bytes before `os.Exit`, so the Windows
  lost-bytes class is eliminated on migrated paths. `runDispatch`'s deferred
  `output.Drain` stays for legacy call sites still on the pipe.
- `RenderErrorSuggestions` (`cli/cmd/rootsuggestion.go:33`) already takes an
  `io.Writer` — pass it the explicit writer.

#### Phase 5 — document the follow-up (not built in 250)

- New code uses the explicit writer from the dispatch context.
- Legacy `fmt.Print*` / `os.Stdout` call sites migrate opportunistically.
- Only when zero call sites depend on the global swap can `output.Install`'s
  pipe mechanism be deleted. That is a separate program.

### Behavior contract (acceptance)

- Colors work: theme filter applies to UI output in non-bright modes.
- Glyph filtering applies to UI chrome in safe mode (`TERM=dumb`).
- `cat` / `view` / `type` output is byte-identical to the file (no emoji→ASCII
  rewrites, no SGR injection in file bytes).
- No lost bytes on exit: a failure message written immediately before `os.Exit`
  is fully flushed (synchronous on migrated paths; drained on legacy paths).
- Bright theme + rich glyphs = zero-cost passthrough, unchanged.
- Existing tests that swap `os.Stdout` directly
  (`cli/cmd/root_no_args_test.go:35-40`) keep passing.

### Files touched

- NEW `cli/output/output.go` — package, `FilterWriter`, `Install`, `Raw`,
  `Drain`, filter-chain composition.
- `cli/theme/install.go` — delete pipe machinery, delegate to `cli/output`.
- `cli/glyphs/install.go` — delete pipe machinery, delegate to `cli/output`.
- `cli/cmd/root.go` — build writer in `Run()`, dispatch context in
  `runDispatch`, `byteFaithfulCommands` bypass via `output.Raw()`, plus the
  stack-trace gate above.
- `cli/cliexit/` — error reporting via the explicit writer.
- `cli/cmd/rootsuggestion.go` — pass the explicit writer to
  `RenderErrorSuggestions`.
- cat/view/type command implementations (via `coreDispatchEntries()` in
  `cli/cmd/rootcore.go`) — route file bytes through `output.Raw()`.
