# Plan 252 — DB-backed spec number issuance + final build and release

## User request (verbatim)
Owner 2026-10-09: spec numbers are hand-picked by agents scanning the filesystem, so two parallel streams both created spec "250". The existing collision is NOT a problem — do not renumber retroactively. Fix forward: a gitmap command that ISSUES numbers from a central authority. "The command should give us the number; the database tells us where it is coming from for whichever repository we are working on." Then final build and release.

## Scope
- Repo: `~/workspace/repos/gitmap-v28`, main.
- Backup branch `backup/252-spec-number-command` created from HEAD and pushed BEFORE any changes.
- New command: `gitmap spec next [--repo <path>] [--json]` — DB-backed atomic spec-number issuance.
- Wave 4: version bump (minor, new feature) → tag → GitHub release → rebuild `~/.local/bin/gitmap`.

## Non-goals
- No retroactive renumbering of the 250 collision (owner explicit).
- `next` is the only required surface — no spec-management suite.

## Checkboxes
- [x] Step 0: git pull (already up to date), backup branch `backup/252-spec-number-command` pushed.
- [x] Research 01 (DB infrastructure) DONE — central DB `.ai-memory/temp-agents/ai_agents.db` (Tier-1, WAL); canonical access via `cli/store` (`ResolveMasterAgentDbPath("")` + `InitMasterAgentDB`); table pattern = idempotent `CREATE TABLE IF NOT EXISTS` at open (copy `ensureCollisionSchema`); write = `INSERT OR IGNORE` + RowsAffected check; repo identity is cwd-based (DB already per-repo, no --repo flag needed).
- [x] Research 02 (command registration) DONE — follow the `space` precedent (constants in constants_cli.go, DispatchSpec in new cli/cmdspec/, register in root.go dispatch chain); help via hand-written cli/helpdoc/spec.md (auto-embedded); max spec on disk = 251 → 252 next; output conventions from cmdstats (plain value + --json).
- [x] Spec 252 written — Spec Agent 1 DONE (01-overview.md + 02-command-design.md); Spec Agent 2 DONE (03-db-and-issuance.md + 2 subtask files). Notable: disk already has manual collisions (240-x2, 250-x2) — DB claim is authoritative.
- [ ] Wave 2: `gitmap spec next` implemented — Worker 1 DONE (cli/store/spec_numbers.go, 92 lines, builds); Worker 2 running (command layer).
- [ ] Spec 252 written (`02-spec/21-app/252-<slug>/` + subtasks).
- [ ] Wave 2: `gitmap spec next` implemented (DB table, atomic issuance, race-safe).
- [ ] Wave 3: verified — `go build ./...` exit 0 (lead-run); N-parallel concurrency test → N distinct numbers; output format checks.
- [ ] Wave 4: committed + pushed; version bumped (minor); tag pushed; GitHub release published; `~/.local/bin/gitmap` rebuilt and version verified.

## Decisions
- D1: No retroactive renumbering (owner explicit).
- D2: `next` returns a plain number for script-friendliness; `--json` for machines.
- D3: Issuance = max(disk spec numbers, DB issued) + 1, atomic INSERT with UNIQUE(repo_path, number), retry on conflict.
