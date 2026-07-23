# 2026-05-06 — commit-in spec authored

**Status:** All phases complete (v4.18.0). End-to-end orchestration, pipeline polish, and CLI integration all landed.

## Produced
- `spec/03-commit-in/` — README + 7 iteration files (overview / CLI / pipeline / DB+ERD / profiles+JSON / message+fn-intel / acceptance).
- `.lovable/plan.md` — appended 7-phase gated implementation plan.
- `.lovable/memory/features/commit-in.md` + index entry.
- Core one-liner in `mem://index.md`.
- Strictly-prohibited entries #3, #4.

## Decisions (resolves prompt's ambiguity list)
1. `all` / `-N` scope: parent dir of `<source>`.
2. Plain `<base>` walks first as `v0`, then ascending `v1..vK`.
3. `--exclude` per-commit BEFORE staging; existing tracked files untouched.
4. `Prompt` mode without VS Code → hard-fail `CommitInExitConflictAborted`.
5. Renamed/moved fn detection: out of scope v1.
6. Pre-existing source commit sharing SHA: skip via `ShaMap`.
7. Profile binding key: absolute symlink-resolved path.
8. `--save-profile` overwrite refused unless `--save-profile-overwrite`.

## Source-repo auto-init precedence (frozen, no flag, no prompt)
1. URL → `git clone` into `CWD/<basename>`.
2. Existing repo → reuse.
3. Existing non-repo folder → `git init` in place.
4. Missing path → `mkdir -p && git init`.

## Verbatim user prompt
The original 2026-05-06 user message ("Complete it in 7 iterations…") is the single source of truth — see git history of this file rather than duplicating the prose to avoid drift.

## Progress
- 2026-05-06 — **Phase 1 ✅** Constants + typed enums + parity tests landed.
  Files: gitmap/constants/constants_commitin.go, gitmap/cmd/commitin/enums.go,
  gitmap/cmd/commitin/enums_test.go; edits to constants_cli.go and
  cmd_constants_test.go.
- Next phases (in order): 2 DB migrations · 3 CLI parsing · 4 Workspace+source
  resolution · 5 Walk+dedupe+replay · 6 Profiles+message pipeline ·
  7 Function-intel+finalize.
- 2026-05-06 — **Phase 2 ✅** DB migrations + enum-mirror seeds landed.
  Files: gitmap/constants/constants_commitin_sql.go,
  gitmap/store/migrate_commitin.go, gitmap/store/migrate_commitin_test.go;
  edits to gitmap/store/store.go (wire-in) and
  gitmap/constants/constants_settings.go (SchemaVersionCurrent 23→24).
  Tables: 18 (8 enum mirrors + Profile + 2 profile children + CommitInRun,
  InputRepo, SourceCommit, SourceCommitFile, RewrittenCommit, SkipLog,
  ShaMap). Tests: presence, seed parity, idempotence.
- 2026-05-06 — **Phase 3 ✅** Pure CLI parser landed (5 files under
  gitmap/cmd/commitin/parse*.go + parse_test.go). RawArgs/ParseError,
  separator+quote split, -N keyword classifier, CSV/enum/author-pair
  validators, flag re-orderer that treats `-N` as positional. Tests
  cover AC #1, AC #4, author-pair, enum rejects, message-rule shape,
  flags-after-positionals.
- 2026-05-06 — **Phase 4 ✅** Workspace + source resolution landed
  under gitmap/cmd/commitin/workspace/ (paths.go, lock.go, source.go,
  expand.go, clone.go, runner.go + workspace_test.go). EnsureWorkspace
  is idempotent; AcquireLock reclaims stale-PID locks; EnsureSource
  implements all four §2.3 cases via a swappable gitRunner;
  ExpandInputs sorts versioned siblings ascending (plain base = v0)
  and supports `-N` truncation; CloneInputs stages each input under
  <TempRoot>/<runId>/<idx>-<basename> with local folders reused in
  place. Hermetic tests (no real git) cover all branches.
- 2026-05-06 — **Phase 5 ✅** Walk + dedupe + replay + runlog landed
  as four sibling packages under gitmap/cmd/commitin/.
  walk/: first-parent oldest→newest via rev-list, \x1f-delimited
  hydrate (author+committer dates + files), empty-repo path returns
  nil. dedupe/: ShaMap lookup; miss is non-error. replay/: byte-perfect
  date replication via plumbing pipeline (cat-file blob → hash-object
  → update-index --cacheinfo → write-tree → commit-tree -p HEAD →
  update-ref). dryRun short-circuits all hooks. runlog/: enum-id
  lookups + StartRun/FinishRun/InsertInputRepo/InsertSourceCommit
  (tx-wrapped) / RecordRewritten (auto-seeds ShaMap on Created) /
  RecordSkip. All hooks are swappable; tests use in-memory SQLite +
  fake git runners — no real git or filesystem required.
- 2026-05-06 — **Phase 6 ✅** Profile + message + prompt packages
  under gitmap/cmd/commitin/. profile/: strict JSON decode (rejects
  unknown fields, gates SchemaVersion=1), canonical encode with §5.2
  key order + trailing newline, atomic SaveToDisk with overwrite
  refusal, ProfilePath under <root>/.gitmap/commit-in/profiles/,
  Resolve() applies four-layer precedence defaults<profile<CLI.
  message/: Build() runs §6.1 stages in order — strip rules + blank
  collapse, override gate (honours OverrideOnlyWeak via §6.2 first-
  word lowercase + punctuation strip), title affix on first line,
  body affix random-pick wrap, function-intel block append, IsEmpty
  flag for EmptyAfterMessageRules skip. prompt/: Asker honours
  --no-prompt by emitting standardized stderr line + ErrNoPrompt for
  exit-code mapping; AskEnum loops until valid. 18 tests across the
  three packages; PickIndex injection keeps message tests determini-
  stic.
- 2026-05-06 — **Phase 7 ✅** Function-intel detectors + finalize +
  dispatcher landed. funcintel/: per-language regex detectors in
  isolated files (Go, JS, TS reusing JS+arrow, Rust, Python, PHP,
  Java shared with C#) self-registering via init() into a central
  registry; render.go emits §6.3 per-file block sorted ascending,
  includes newly-added files even when no functions detected.
  finalize/: Counters + Outcome (PartiallyFailed when Failed>0),
  PrintSummary using CommitInMsgSummaryLine, CleanupTemp honours
  --keep-temp, Resolve maps ConflictMode→ConflictDecision with
  standardized abort banner. Dispatcher: runCommitIn parses argv +
  exits BadArgs on parse error + emits "orchestration loop pending"
  stub; rootcore.go registers CmdCommitIn/CmdCommitInAlias.
  helptext/commit-in.md (105 lines, 5 examples, flag + exit-code
  tables). Version 4.17.0 → 4.18.0. 13 tests added. Deferred
  (non-blocking): end-to-end orchestration glue inside runCommitIn,
  // gitmap:cmd top-level marker on CmdCommitIn const block,
  CHANGELOG v4.18.0 entry.
- 2026-05-06 — **Step 2 ✅** End-to-end orchestration glue landed.
  New gitmap/cmd/commitin/orchestrator/ package (7 files, all <200
  lines, all funcs ≤15 lines): run.go owns top-level Run + setUp +
  finalRunStatus mapping; setup.go threads resolveSource → workspace
  → lock → store.OpenAt+Migrate → loadProfile via new exported
  store.OpenAt(dbPath) so the SQLite anchors at <source>/.gitmap/
  gitmap.db per spec; cli_overrides.go projects RawArgs onto
  profile.CliOverrides; pipeline.go expands+stages inputs then walks
  each one with a per-run rand.New seeded picker (spec §3.4
  determinism); commit.go runs the per-commit dedupe → message build
  → replay → record loop with separate skip/fail/created paths and
  routes IsDryRun to a SkipReasonDryRun branch; context.go bundles
  handles + idempotent Cleanup (CleanupTemp respects --keep-temp);
  input_cache.go ensures one InputRepo row per staged input via
  OrderIndex cache. runCommitIn delegates to orchestrator.Run and
  propagates exit codes. go build ./... clean; go test ./cmd/
  commitin/... ./store/... all green.
- 2026-05-06 — **Step 3 ✅** Per-commit pipeline polish: exclusions
  filter + function-intel block + ExcludedAllFiles skip path. Two new
  files in orchestrator/: exclude.go (applyExclusions /
  matchesExclusion / matchesFolder — PathFile = exact rel match,
  PathFolder = prefix or any path segment match, POSIX-normalized via
  filepath.ToSlash), funcintel_block.go (renderFunctionIntel best-
  effort builder that runs `git show <sha>:path` + `git show
  <sha>^:path` per file, dispatches via funcintel.LanguageForPath +
  EnabledLanguages, returns "" on any failure so a parser glitch never
  aborts the commit per spec §6.3). commit.go now: filters c.Files
  through applyExclusions, emits SkipReasonExcludedAllFiles when the
  filter empties a non-empty file list, materializes the §6.3 block
  via renderFunctionIntel, threads it into message.Build via the
  FunctionIntel input field. New exclude_test.go covers 4 cases (pass-
  through, folder+file mix, nested-segment folder match, exact-only
  file match). go build clean; ./cmd/commitin/... ./store/... all
  green.
