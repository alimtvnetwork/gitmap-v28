# 268 — pe/te help routing + `all` verb documentation

## 1. Problem statement (from the pe/te error audit)

1. `gitmap pe --help` and `gitmap te --help` rendered a generic auto-generated
   placeholder ("Execute, manage, and automate GitMap pe/te operations") via
   the root dynamic help renderer, teaching users nothing about `-t`, `all`,
   `-f`, or the format system. (`pe -t --help` reached the real help, proving
   the routing gap.)
2. The `all` aggregator verb (program 263: `pe all` / `te all` /
   `pipeline errors all`) appeared nowhere in the pipeline help — undiscoverable.
3. `--timeout`, `-w`, `--watch` silently trigger timeline mode but are
   undocumented (left as undocumented compat shims; `-t`/`--timeline` is the
   documented form).
4. Edge (documented, not changed): a bare `all` anywhere in args triggers the
   aggregator, so a repo literally named "all" cannot be inspected single-repo.

## 2. Design

- Root help interceptor (`cli/cmd/root.go: tryInterceptCommandHelp`) now
  detects the pipeline-errors family (`pe`, `te`, `ee`, `pipeline`,
  `pipeline-errors`, `pipeline_errors`) on a help flag and delegates to the
  canonical printer `cmdpipeline.PrintPipelineErrorsHelp()` (new export in
  `cli/cmdpipeline/exports.go`), exiting 0 — before the generic dynamic
  renderer.
- `printPipelineErrorLogsUsage` gains the `pe all` usage line;
  the Commands section documents `all [--json] [--file <path>]`.

## 3. Acceptance

- `go build ./...` exit 0; gofmt clean.
- `gitmap pe --help` and `gitmap te --help` print the real pipeline help
  (usage lines, Commands incl. `all`, Targeting, Flags).
- `gitmap pe all --help`, `gitmap pipeline errors --help` unchanged behavior.
