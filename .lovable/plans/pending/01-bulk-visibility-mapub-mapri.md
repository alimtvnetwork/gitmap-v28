# Bulk wildcard visibility flips — `make-all-public` / `make-all-private` / `MAPUB` / `MAPRI`

**Slug:** bulk-visibility-mapub-mapri
**Steps:** 40
**Status:** pending
**Created:** 2026-06-06

## Context

Add four new CLI commands that flip GitHub/GitLab repo visibility in bulk
against an owner/org, selected by comma-separated wildcard patterns
(`exact`, `prefix*`, `*contains*`, `prefix*suffix`). Interactive numbered
list with per-index exclusion; `-Y` bypasses confirmation. Every run +
per-repo result persisted to SQLite in new `GitMapRun` / `GitMapRepoResult`
tables. Help MD + changelog + version bump at the end.

**Confirmed unique** — existing `make-public` / `make-private`
(`gitmap/constants/constants_cli.go:173-174`, `gitmap/cmd/visibility.go`,
`visibilitybulk.go`) target ONE repo (or `vN..vN-count+1`); they do NOT
operate owner-wide or accept wildcards. `MAPUB` / `MAPRI` do not exist.

**Captured inputs:**
- Full verbatim user spec: this file's "Verbatim user spec" section below.
- Command: `.lovable/spec/commands/01-make-public-url-or-folder.md`
  (URL ↔ folder ↔ `.` interchangeable for ALL four commands).
- Command: `.lovable/spec/commands/02-bulk-visibility-mapub-mapri.md`
  (high-level command capture).

**Files involved:**
- `gitmap/constants/constants_cli.go` — add 4 CLI IDs
- `gitmap/cmd/dispatch.go` (or equivalent) — wire commands
- `gitmap/cmd/visibilityresolve.go` — extend owner-only resolver
- `gitmap/cmd/visibilityallbulk.go` (new) — orchestrator
- `gitmap/cmd/visibilitybulkprompt.go` (new) — interactive UX
- `gitmap/visibility/pattern.go` (new) — wildcard engine
- `gitmap/db/migrations/0NN_gitmap_run_repo_result.sql` (new)
- `gitmap/db/enums.go` + `gitmap/db/gitmaprunrepo.go` (new)
- Help: `gitmap/cmd/help/make-all-public.md`, `make-all-private.md`, `MAPUB.md`, `MAPRI.md`
- `CHANGELOG.md` + docs-site changelog mirror
- `gitmap/constants/constants_version.go` — bump

**Verbatim user spec** (for traceability — DO NOT lose):
> see chat message 2026-06-06 with header "GitMap Bulk Repo Visibility Commands Instruction".
> Key clauses: 4 commands (long+short), URL/bare/folder/`.` accepted, `*` wildcard
> only, comma-separated patterns, numbered list + numeric exclusion, `-Y` bypass,
> SQLite `GitMapRun` (1→n) `GitMapRepoResult` with int-backed enums, PascalCase PKs,
> help MD + changelog + version bump.

## Steps

1. Re-confirm uniqueness: `rg -n "make-all-public|make-all-private|MAPUB|MAPRI"` over `gitmap/` returns zero hits.
2. Read coding guidelines: `spec/05-coding-guidelines/**/*.md` (every file) + any `spec/*-error-manage/` folder; apply Boolean (`is`/`has` prefix), Enum (small-int), and zero-swallow error rules to every code step below.
3. Write the formal spec at `spec/01-app/NNN-bulk-visibility-mapub-mapri.md` (next free `NNN`): symptoms, command grammar, pattern rules, interactive flow, SQL DDL, ERD (Mermaid), acceptance checklist.
4. Add four CLI ID constants to `gitmap/constants/constants_cli.go`: `CmdMakeAllPublic = "make-all-public"`, `CmdMakeAllPrivate = "make-all-private"`, `CmdMAPUB = "MAPUB"`, `CmdMAPRI = "MAPRI"` (per no-magic-strings rule, IDs exclusively in this file).
5. Register the four commands in the dispatcher with a single shared handler `runMakeAllVisibility(target VisibilityTarget, args []string)`; add `// gitmap:cmd top-level` marker so the completion generator picks them up (per Core memory).
6. Implement `ResolveOwnerOnly(arg string) (Provider, Owner, error)` in `gitmap/cmd/visibilityresolve.go`: accept full URL, bare `host/owner`, folder path, `.`; folder/`.` reads `.git/config` origin → owner. See `.lovable/spec/commands/01-make-public-url-or-folder.md`.
7. Implement pattern engine. See `./subtasks/01-bulk-visibility-mapub-mapri/01-pattern-engine.md`.
8. Implement comma-list parser `ParsePatternList(raw string) ([]Pattern, error)`: split on `,`, trim spaces, dedupe, reject empty tokens, surface `ParsePattern` errors with token index.
9. Implement repo lister + provider API mapping. See `./subtasks/01-bulk-visibility-mapub-mapri/04-provider-api-mapping.md`.
10. Compose: `MatchOwnerRepos(owner, patterns) []MatchedRepo` returning ordered `{RepoName, MatchedPattern}` pairs, dedupe by `RepoName` keeping first matcher.
11. Implement numbered-list renderer `renderMatchedTable(matches []MatchedRepo, action string) string` with 1-based indexes, right-aligned counts, color-aware via existing `constants.ColorBold`.
12. Implement interactive prompt + exclusion parser. See `./subtasks/01-bulk-visibility-mapub-mapri/03-exclusion-parser.md`.
13. Add `-Y` / `--yes` flag parsing reusing existing `visibilityFlags` struct; route through `opts.autoConfirm`; ensure it skips BOTH the y/n prompt AND the exclusion prompt (per spec "Important #3").
14. Implement apply loop `applyAllVisibility(ctx, matches, target, runId, db)`: per repo call existing `applyVisibility`, update `GitMapRepoResult.ResultStatus` immediately after each call (Succeeded/Failed), continue on failure.
15. Build error aggregation: end-of-run summary table `OK n  FAIL m  SKIPPED k`; exit code 0 if all OK, 2 if any FAIL (matches existing visibility convention — verify in `visibility.go`).
16. Define enums in `gitmap/db/enums.go`: `CommandKindEnum` (MakeAllPublic=1, MakeAllPrivate=2), `ResultStatusEnum` (Pending=1, Succeeded=2, Failed=3, Skipped=4); add `String()` method per enum; export `Parse*` for tests.
17. Create SQLite migration + repository layer. See `./subtasks/01-bulk-visibility-mapub-mapri/02-sqlite-schema.md`.
18. Register migration in `gitmap/db/migrations.go` registry (or whatever the existing registration mechanism is — read `gitmap/db/` first).
19. Wrap run-row insert + per-result inserts in one BEGIN/COMMIT tx; on tx-open failure, log to stderr and continue without persistence (so a corrupted DB does NOT block visibility flips per zero-swallow rationale — but log loudly).
20. Persist `GitMapRun` row at start of `runMakeAllVisibility` BEFORE the prompt (so even an aborted run is recorded with rows marked `Skipped`).
21. After exclusion phase, set `IsExcluded=1` on excluded repos and `ResultStatus=Skipped`; non-excluded start as `Pending`.
22. After each `applyVisibility` call, `UPDATE GitMapRepoResult SET ResultStatus=? WHERE GitMapRepoResultId=?`.
23. Centralize all log/error messages in `gitmap/constants/constants_visibilitybulk.go` (no magic strings rule); include `MsgMakeAllPrompt`, `MsgMakeAllExcludePrompt`, `ErrMakeAllNoMatches`, etc.
24. Apply error-management rules: `fmt.Fprintf(os.Stderr, "✗ %s: %v\n", repo, err)` for per-repo failures; use `errors.Is(err, context.Canceled)` to detect Ctrl-C and persist remaining rows as `Skipped`.
25. Add pre-flight auth check: call existing `mustEnsureProviderCLI(provider, verbose)` before the repo list.
26. Add GitHub repo-list pagination cap warning: if `gh repo list` returns exactly `--limit`, log a stderr warning that more repos may exist beyond the cap.
27. Honor `make-all-* .` and `make-all-* <folder>`: when arg resolves to a folder, extract its origin owner and operate owner-wide on that owner (per commands/01).
28. Tests: `gitmap/cmd/visibilityresolve_test.go` — table tests for URL / bare host / folder / `.` → same owner.
29. Tests: `gitmap/visibility/pattern_test.go` — all six rules from subtask 01 (exact, `macro*`, `*mate*`, `lotus*v1-4`, ordering `a*b*c`, bare `*` error).
30. Tests: `ParsePatternList` — whitespace, trailing comma, empty token, duplicate token dedup.
31. Tests: exclusion parser per subtask 03 (`"1,3-5"`, `"none"`, `"all"`, out-of-range, descending range, non-numeric).
32. Tests: `runMakeAllVisibility` with `-Y` skips BOTH prompts (use a stub stdin that would panic if read).
33. Tests: DB layer — open `:memory:` SQLite, run migration, insert run + results, read back; assert FK cascade by deleting run row.
34. Tests: Golden command-help fixtures — run `gitmap regoldens` after the new help MDs land (per Regoldens Command memory).
35. Update help MD files: create `gitmap/cmd/help/make-all-public.md`, `make-all-private.md`, `MAPUB.md`, `MAPRI.md` — each ≤120 lines, with a 3-8 line realistic simulation (per Command Help System memory). `MAPUB.md` / `MAPRI.md` may be 5-line stubs pointing at the long form.
36. Add `CHANGELOG.md` entry under a new minor version section: list the four commands, the two new tables, the migration filename; mirror to docs-site React changelog (per Changelog System memory).
37. Bump version constant (locate via `rg 'CurrentVersion|VersionString' gitmap/constants/`); run `gitmap fix-repo --strict` over the touched packages to auto-rewrite `{base}-vN` references and re-run package tests (per Fix-Repo Strict Mode memory).
38. Run `gofmt -l .`, `go vet ./...`, `golangci-lint run ./...` (golangci-lint v1.64.8 per Core memory). Zero findings required.
39. Run targeted tests: `go test ./gitmap/cmd/... ./gitmap/visibility/... ./gitmap/db/... -count=1 -race`. All green.
40. Manual smoke: against a throwaway personal owner with at least one repo matching each pattern shape, run all four commands twice (once interactive with exclusions, once with `-Y`); verify `GitMapRun` + `GitMapRepoResult` rows; verify visibility flipped on GitHub UI; verify exit codes 0/2.

## Verification

- Step-by-step exit signals are listed inline (build/vet/test/lint commands, DB read-back, GitHub UI check).
- Final smoke (step 40) is the only acceptance gate that requires a live provider; everything else runs in CI.
- After step 39, move this file to `.lovable/plans/completed/01-bulk-visibility-mapub-mapri.md` and flip `Status:` to `completed`.
- Subtask files (`./subtasks/01-bulk-visibility-mapub-mapri/0{1..4}-*.md`) flip individually as each is finished.

## Appended from prior pending tasks

none — `.lovable/plans/pending/` and `.lovable/plans/completed/` were empty before this plan; no prior subtasks to absorb. (Other `.lovable/` content — `pending-issues/`, `cicd-issues/`, `solved-issues/`, `prompts/`, `memory/` — was scanned and contains no actionable steps for this bulk-visibility scope.)
