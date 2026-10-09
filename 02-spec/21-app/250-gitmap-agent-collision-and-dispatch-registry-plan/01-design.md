# Spec 250 — Agent Collision Handling, Dispatch Registry, A08, Version Truth

## 1. Problem statement

Three structural defects, one theme — the system does not know about itself:

1. **Dispatch collisions are invisible.** Command/alias names are claimed in scattered
   dispatch tables (`cli/cmd/rootcore.go`, `cli/cmd/root.go:dispatchAgySubsystem`, …)
   with no uniqueness validation. The `agm`→gitignore hijack (task 247) shipped and
   made `gitmap agm update` unreachable. Nothing would have caught the next one.
2. **Agents collide blindly.** Concurrent AI workstreams (246/248/249) modify the same
   tree with no shared knowledge of who owns which files, and `gitmap cpf` runs
   `git add -A` (`cli/cmdcommit/commit_push.go:executeStageAllStep`), sweeping every
   workstream's files into whichever commit runs first. The task DB
   (`cli/cmdagent/agent_config.go`: `ParentTask`/`Subtask`) already has an
   `OwnedFilesJson` column — it is never populated.
3. **Swallowed errors (A08).** `writeLastErrorFile` (`cli/cmd/root.go:404`) discards
   `os.MkdirAll` and `os.WriteFile` errors with `_ =`, violating the never-swallow
   rule. Failures are invisible.
4. **Scattered version truth.** `version.json` carries two version fields (`Version`
   and `version`), `cli/constants/constants.go` carries a third, and `version.json`
   lagged the `v6.516.0` tag. "Single source of truth" is currently aspirational.

## 2. Dispatch registry (design)

**Goal:** a duplicate command/alias name anywhere in the dispatch tables fails fast.

- New small package `cli/cmddispatch/registry.go`: `CollectNames() []DispatchName`
  aggregates every claimed name+alias from each dispatch table
  (core entries in `rootcore.go`, `dispatchAgySubsystem` cases in `root.go`,
  extra/general commands), each tagged with its owner file:line.
- `gitmap doctor --check-dispatch` runs the aggregation and reports duplicates with
  both owners; exit 1 on any collision.
- Go test `TestDispatchRegistryNoCollisions` asserts the same (runs in CI).
- Regression proof: temporarily re-adding `"agm"` to the gitignore aliases must fail
  the check (test with a synthetic duplicate, not by touching production tables).
- No refactor of the dispatch flow itself in this task — aggregation only. A full
  table-driven registry (single `Register()` choke point) is a documented follow-up.

## 3. Agent collision feature (design)

Built on the existing 3-tier task engine (`gitmap agent task|subtask`, SQLite at
`.ai-memory/temp-agents/<slug>/agent-task.db`, already gitignored).

### 3.1 Look-ahead: file-box claims (the feed-before-work step)

- `gitmap agent subtask claim --task-id <slug> --files <rel,paths…>` (new; also
  accepted on `subtask add` via `--files`). The lead declares each worker's file box
  BEFORE dispatching. Stored in `Subtask.OwnedFilesJson` (column exists).
- Workers read their own file box via `gitmap agent subtask ls --task-id <slug>`
  (extend output with files). An agent knows what it owns and what siblings own.

### 3.2 Collision visibility (not locks)

- `gitmap agent collisions [--parent <slug>]`: reports overlapping write claims
  across active subtasks (same file, ≥2 owners), with owner subtasks and statuses.
- Mindset (per owner): overlaps are REPORTED, never hard-blocked — parallelism is the
  point. The lead plans boxes to minimize overlap; whatever collides is caught by
  the end-of-task verify-and-fix pass (build + e2e + `collisions` re-check).

### 3.3 gitmap maintains the adding (no agent-side `git add`, ever)

- `executeStageAllStep()` (`cli/cmdcommit/commit_push.go:445`): when a task context
  is active (`--task <slug>` flag on `cpf/cpb/cpc/cpr`, or `GITMAP_TASK` env), stage
  ONLY the union of that task's `OwnedFilesJson` paths (`git add -- <paths>`) instead
  of `git add -A`. Everything else stays out via `.gitignore` (temp-agents already
  ignored) and unstaged.
- Agents are FORBIDDEN from running `git add`/`git commit` directly — staging and
  committing happen only through `gitmap <cpf|cpb|cpc|cpr> --task <slug>`.
  (This retires the lead-side targeted-`git add` workaround; the tool owns it.)

### 3.4 Communication model

Agents communicate through the shared task DB, not direct messaging: claims,
statuses, and the action log (`AgentActionLog.TargetFile` already exists) are the
bus. `gitmap agent task progress --slug` shows file boxes alongside statuses.

### 3.5 Agent observability methods (added per owner review)

`gitmap agent` gains read-only observability subcommands (methods, not flags),
backed by the central task registry DB:

- `gitmap agent stats [--days N]` — tasks recently completed, subtasks completed,
  per-agent completion counts, total collisions detected.
- `gitmap agent heatmap [--parent <slug>]` — files ranked by write-claim count
  (from `FileClaim` + `AgentActionLog.TargetFile`).
- `gitmap agent ps` — running tasks: slug, status, active agent count, subtask
  progress rollup.
- `gitmap agent collisions` (§3.2) reports live overlaps; collision HISTORY comes
  from the `CollisionEvent` table (written at claim-time overlap detection and on
  each `collisions` run).

**Schema** (central registry DB, keyed by parent task slug; implementation confirms
the central DB path — `ParentTaskRegistry` table location is the lead candidate;
per-task DBs keep task details):

```sql
CREATE TABLE IF NOT EXISTS FileClaim (
  ClaimId INTEGER PRIMARY KEY AUTOINCREMENT,
  ParentTaskSlug TEXT NOT NULL,
  SubtaskId INTEGER NOT NULL,
  AgentRole TEXT NULL,
  FilePath TEXT NOT NULL,
  ClaimKind TEXT NOT NULL DEFAULT 'write',
  Status TEXT NOT NULL DEFAULT 'active',
  CreatedAt TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS CollisionEvent (
  EventId INTEGER PRIMARY KEY AUTOINCREMENT,
  ParentTaskSlug TEXT NOT NULL,
  FilePath TEXT NOT NULL,
  OwnerA TEXT NOT NULL,
  OwnerB TEXT NOT NULL,
  DetectedAt TEXT NOT NULL,
  Resolved INTEGER NOT NULL DEFAULT 0
);
```

File boxes declared up front go stale — agents MUST re-claim (`subtask claim`
again) when their box grows mid-work; the latest claim wins.

## 4. A08 swallowed-error fix (design)

- `writeLastErrorFile` (`cli/cmd/root.go:404`): replace `_ = os.MkdirAll` /
  `_ = os.WriteFile` discards with visible failure — concise stderr warning
  (`warning: could not persist last-error log: <err>`) via the existing error
  pipeline. The original error being handled is never obscured; the persistence
  failure is no longer silent.
- Audit `persistToErrorsDB` on the same path for discards; fix any found the same way.
- No signature changes required if the warning approach is taken (keep the diff small).

## 5. Version truth unification (design)

- `version.json`: exactly ONE version field (`Version`). Remove the lowercase
  duplicate; enumerate and migrate every reader (implementation phase lists them).
- Write rule: human edits ONLY `version.json` → `Version`. `37-bump-version.py`
  remains the only writer for everything else (`constants.go`, `package.json`,
  changelog, readme pins).
- Enforcement: wire the existing `03-ai-scripts/14-version-sync-checker.py` into the
  pre-tag release gate — the tag is refused if any copy diverges.
- `cli/constants/constants.go`: keep the ldflags override mechanism; default value
  synced by the bump script (no hand edits).

## 6. File boxes (implementation phase, post-review)

- Subtask 01 (dispatch registry): `cli/cmddispatch/` (new), `cli/cmd/doctor*.go`
  (check wiring) — new files only, no dispatch flow changes.
- Subtask 02 (collision feature): `cli/cmdagent/` (claim/collisions CLI),
  `cli/cmdcommit/commit_push.go` (`executeStageAllStep` scoping only).
- Subtask 03 (A08): `cli/cmd/root.go` (`writeLastErrorFile` only + audit).
- Subtask 04 (version truth): `version.json`, readers, release gate wiring.
- Explicitly out: `gitmap update` self-updater, command taxonomy changes
  (subcommands stay, per owner), Letterly prompt edits (idea shared, not implemented).

## 7. Release

Minor bump to **6.518.0** via `03-ai-scripts/37-bump-version.py` (proven in 247),
changelog entry, annotated tag `v6.518.0` (`Release v6.518.0`), push main + tag.
Test-suite gates skipped per standing rule; e2e is the evidence.

## 8. Test plan

`go build ./...` only (no `go test` per standing rule). E2E per acceptance spec 02,
all artifacts outside the repo.

## 9. As-built notes (lead verification 2026-10-09)

- Claim command is `gitmap agent subtask claim-files` (not `claim`): `subtask claim`
  already exists for work-queue claiming; the spec's `claim --files` form was
  ambiguous without a subtask target.
- Completion counting uses `Status IN ('DONE','COMPLETED')`, not `HasCompleted`:
  the Tier-2 `Subtask` table has no `HasCompleted` column (older schema); the
  `_ =` discards in the original code masked this.
- Scoped staging partitions claimed files: ignored-but-tracked files stage via
  `git add -f`, ignored-and-untracked are skipped. `git check-ignore` needs
  `--no-index` to see tracked files (without it, one ignored path poisons the
  whole `git add`). Found during the release's own dogfood commit.
- `agent ps` originally scanned `ParentTaskId` as int64; it is TEXT (string task
  IDs). Fixed during verification.
