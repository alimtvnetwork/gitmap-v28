# Plan 325: GitMap and Antigravity CI/CD and Code Refactoring

Slug: gitmap-and-antigravity-cicd-and-code-refactoring
Spec: 02-spec/21-app/325-gitmap-and-antigravity-cicd-and-code-refactoring/
Status: active — overnight run, user asleep, no questions
Coordinator decisions: see 01-architecture-spec.md (ambiguity decisions 1–9)

## State at plan time (2026-10-25 ~02:55 +08)

- gitmap binary updated to v6.608.0 (manual placement; `gitmap update` placement step fails in this container, E1151; SHA256 verified; old binary backed up).
- gitmap pa: 17 pulled; 2 failures, one was gitmap-v28 (merge conflict from the pending local gosec fix eac563f2) — resolved by completing the merge and pushing d886982d; the second failure is an out-of-scope repo and is not chased.
- gitmap-v28 tip: d886982d, clean, pushed. CI runs in flight on it.
- Antigravity-Manager tip: cfd4b5258 (version 4.266.0), clean, CI green on main (CI success on the last 5 runs).
- coding-guidelines-v24 tip: 181aa98d, clean.
- Other session actively pushing to gitmap-v28 and AGM. Rule: fetch, work on latest tip, never fight it.

## Subtasks

| ID | Subtask | Owner | Owned files / repos | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | Spec + plan (this file + spec folder) | Coordinator | 02-spec/21-app/325-.../, .ai-memory/plans/pending/325-...md, .ai-memory/plans/subtasks/325-.../ | DONE | spec files written |
| Task-02 | AGM CI-green + minor-bump release | Worker AGM | Antigravity-Manager only | QUEUED | — |
| Task-03 | gitmap CI-green on latest tip | Worker gitmap-ci | gitmap-v28 CI-blocking files only | QUEUED | — |
| Task-04 | Enum/constants audit (whole repo, grouped by package) | Worker audit | read-only, gitmap-v28 cli/ | QUEUED | — |
| Task-05 | Enum/constants refactor waves (per audit groups) | Workers refactor (waves) | gitmap-v28 cli/ package groups, disjoint per wave | QUEUED | — |
| Task-06 | Prompt update: enum-per-package rule + cargo_exe example | Worker prompt | coding-guidelines-v24 enum CG execute prompt + mirrors | DONE | commit 709c1722: 01-prompts/15-cg-execute/06-constants-and-enums.md new §4.3 (hard rule) + .agents/.cursor cg-enum-standards mirrors (byte-identical); constants-and-enums pointer skills deliberately untouched (they point at prompt 06) |
| Task-07 | gitmap minor-bump release | Coordinator + worker | gitmap-v28 | QUEUED | — |
| Task-08 | Morning report | Coordinator | — | QUEUED | — |

## Wave plan

- Wave 1 (parallel, now): Task-02 (AGM release), Task-03 (gitmap CI), Task-04 (audit).
- Wave 2: Task-05 refactor waves from the audit groups, Task-06 prompt update (after the cargo_exe exemplar lands).
- Wave 3: Task-07 gitmap release, Task-08 morning report.

## Audit results (Worker audit, 2026-10-25 ~03:05 +08) — Task-04 DONE

Corrections to the original order (verified, not assumed):
- cli_cmd_cargo/cargo_exe.go does NOT exist. Real file: cli/cmdcargo/cargo_exec.go (4 `+` sites, 3 enum-derived).
- cli/cmdpool/pool_cli_input.go does NOT exist (no cli/cmdpool package). Real file: cli/cmdpull/pull_cli_input.go, flag chain starts line 76 (13 comparisons, 12 constants, already constant-based; remaining step is set lookup). One suspect inside it: `"--" + pullconstants.FlagWorkerShort` (short flag with long prefix) — verify against flag registration before copying.
- ApplyDefaultPragmas (cli/store/constants/constants.go:17) ALREADY returns *apperror.AppError. Remaining work there is context enrichment (driver, path) at the OpenSQLiteDB caller, not conversion. Also fix three call sites that discard the error outright (cli/cmdagy/agy_clean_cache_purge.go:22, agy_prompts_collector.go:100, agy_undo_cache.go:57) and delete the dbengine pragma duplicate (cli/dbengine/wrapper.go applySQLitePragmas, same 4 PRAGMAs as literals).

Totals: 27 packages have their own enums/ (49 files); 48 packages still hold stray enums (incl. partially-done cmdclone, cmdcommit). 25 packages have own constants/. String `+` on enum values: 9 open sites (4 cmdcargo, 4 cmdrelease/release, 1 dbengine adjacent). `==` flag chains: 517 `arg ==` lines tree-wide, overwhelmingly literal. Regex: 166 eager regexp.MustCompile hits (package-level var), 18 regexp.Compile mostly legit dynamic (keep); lazy pattern to copy is cli/lazyregex (lazyregex.New, pattern from constants + lazy compile, e.g. cli/cmdcommit/commit_formatter.go). Wrapped returns: convention is cli/result containers + flat per-payload aliases (cli/store/repository_types.go: ScanRecordWrap etc.); OpenSQLiteDB is the 1 of 2 sqlite_config.go functions still two-valued; dozens of cli/store two-value functions remain; every done enum's Parse still returns plain error.
Name collisions to never merge: ModeType (glyphs, termout, render PrettyModeType), DirectionType/PreferPolicyType (committransfer, movemerge), AlignmentType (termgeneric, termout).

Wave 0 decisions (coordinator):
- (a) DB vs Db vs Database: DB wins (Go initialism rule). Rename DbWrapper → DBWrapper, OpenDb → OpenDB, WrapDb → WrapDB, DbSQLite → DBSQLite. String values unchanged.
- (b) cli/apperror ErrorType/SeverityType: LEAVE in place (error infra; moving it risks everything). Intentional exception, recorded.
- (c) cli/pipelinedb/enums and cli/repodb/enums record structs: LEAVE (already there; moving them is churn). Mixed-but-accepted, recorded.
- (d) FlagWorkerShort long-prefix suspect: exemplar worker verifies against registration before copying the pattern.
- Wave order: Wave 1 store/dbengine core (largest blast radius, 75+ OpenSQLiteDB call sites — never parallelize with other waves) → Wave 2 done-package cleanups (cmdcargo, cmdpull, cmdrelease/release, then remaining done packages' Parse conversions) → Wave 3 remaining small packages → Wave 4 remaining large packages → Wave 5 shared cli/constants + lazyregex sweep. Wave 1 is HELD until the gitmap-ci worker finishes, to avoid colliding on CI-blocking files.

## Added phase: gitmap commit-suite --dry-run actually dry-runs (added 2026-10-25 ~03:15 +08)

Incident (verified in ~/workspace/repos/img-pdf, the img-pdf-v2 clone): `gitmap cpf --dry-run` did NOT dry-run. It committed and pushed 67b7b21 "Feature: test dry run", sweeping unrelated uncommitted sync files and .gitmap/backup into main. Recovery was cleanup commit 9a7840b (no force-push). The flag was accepted/ignored and the command mutated anyway.
Scope: gitmap-v28 commit-suite commands (cpf/cpb/cpc/cpr/pcp) + every other mutating command advertising dry-run (push, release, sync, anything with a dry-run flag). Work: root-cause how --dry-run is handled today (parsed but ignored, never threaded through, shadowed); implement a REAL dry-run that prints exactly what it would do (files it would stage, commit message, branch, remote/push target, tag/release actions) and exits WITHOUT mutating anything (no commit, no push, no staging changes left behind, no files swept in; the file list must make a sweep visible before a real run); add tests proving zero mutations (no commit, no push, clean tree after); keep files under the 300-line cap. Then update the gitmap SKILL in coding-guidelines-v24 (.agents/skills/gitmap + .cursor/skills/ mirror): document the flag, what it prints, and the standing rule — verify with --dry-run plus git status/diff before any gitmap commit — including the caution that pre-fix --dry-run did not actually dry-run. Commits via the gitmap suite as the user, incremental, CI verification per standing rules. Coordination: cli/cmdcommit is refactor territory — the dry-run worker owns commit-suite files for this phase; other refactor waves must not touch them concurrently (RC heatmap tracked). Task-09 in the table below.

| Task-09 | gitmap --dry-run real dry-run (root cause, fix, tests, skill doc) | Worker dryrun | gitmap-v28 commit-suite files + coding-guidelines-v24 gitmap skill mirrors | QUEUED | — |

## Assumptions

- The other session's releases (AGM 4.266.0, gitmap v6.608.0) are legitimate user-side work; this run builds on top of them.
- AGM release uses the repo's bump script + `gitmap release`, same discipline as the v4.256.0 release.
- If the enum refactor cannot finish in one night, CI-green + both releases + a precise remaining-groups list is the acceptable fallback.

## Conflicts

- execute-parent-task-with-n-steps-v6 vs parent task on atomic commit / zero CI / zero releases: parent wins (see spec decision 7).
- Skill worker-git-ban vs parent worker execution: workers may use the gitmap suite inside their own repo only (see spec decision 8).
