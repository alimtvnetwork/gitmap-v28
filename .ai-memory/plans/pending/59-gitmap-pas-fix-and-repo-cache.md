# Plan 59: GitMap PAS Fix, CPAR, Repo Cache, and Nodes Clone UI Fix

## Status: Pending
- **Plan ID:** 59 (Tracking Spec 197)
- **Spec Reference:** [02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md](../../../02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md)
- **Scope:** CLI Subsystems, Async Worker Pool, Split-DB Cache Engine, Remote Nodes Telemetry
- **Created At:** 2026-10-01
- **Parent Goal:** Consolidate GitMap PAS Formula execution, asynchronous ignore scanning with non-blocking prompt resolution, GitMap CPAR commit-push suite, split-DB repository cache engine with full-text/regex search, and fleet nodes clone reachability telemetry.

---

## 1. Immediate Pre-Requisite Fixes
- [ ] Task 0.1: **Fix `nodes clone` alignment and `w3` reachability misreporting** (as per Spec 196 / RCA 59). Ensure `DURATION` uses `%-12s` and properly times local executions. Update `(machine is off)` string to reflect unreachable port 22.
- [ ] Task 0.2: **Address `v6.436.0` Credential Store Error on Fleet**: Ensure nodes update to the new version deployed in the previous session so the `pa --ssh` credential bugs vanish.

---

## 2. Phase 1: GitMap Pull All & PAS Formula Optimization
- [ ] Task 1.1: Refactor `gitmap pull all (pa)` to do raw pulling first. Move ignore duplicate checking to an async worker group (handling up to 5 repos at a time).
- [ ] Task 1.2: Implement `gitmap pull all ssh (pas)` utilizing the "GitMap PAS Formula". Highly restrictive concurrent workers (max 1 or 2 workers doing 2 async operations) on SSH nodes. Enqueue into task server and return JSON responses.

---

## 3. Phase 2: GitMap Ignore Subsystem & Groups
- [ ] Task 2.1: Create `gitmap ignore` subcommands skeleton (`add`, `ls`, `rm-grp`, `agtd`, `cgwp`, `apply`, `export`, `import`).
- [ ] Task 2.2: Implement SQLite-backed Gitignore Group persistence. Define default group automatically including `.gitmap/` and `.gitmap/backup/`.
- [ ] Task 2.3: Implement duplicate rejection and validation rules (only flag if tracked in current commit state).

---

## 4. Phase 3: Global Operation Commands
- [ ] Task 3.1: Build `gitmap fix ignore all [-y] (fia)` and its SSH variant `fias`, using global prompts or single-stepping fix sessions.
- [ ] Task 3.2: Implement `gitmap commit-push-all-repos (cpar) [-y] [--review|-r] [--commit-only|-co]`.
- [ ] Task 3.3: Implement `gitmap see (c)` subcommands (`commit pending`, `gitignore issues (ig)`, `errors`).

---

## 5. Phase 4: Split-DB Repo Cache & Search Engine
- [ ] Task 4.1: Architect `gitmap cache create <path>` split-DB logic (`sql.db` at root, `[slug].db` per folder).
- [ ] Task 4.2: Build file scanner avoiding >200KB files, images, binaries. Track `last_modified` (no hashes).
- [ ] Task 4.3: Implement `gitmap cache search`, `multi-search`, `search-multi-grep` directly querying SQLite DBs.
- [ ] Task 4.4: Add async cache reconciliation during searches (detecting outdated `last_modified`).

---

## 6. Phase 5: UI, Help Text, and Final Verification
- [ ] Task 5.1: Write expansive CLI Help, terminal Help, UI texts, and update markdown docs for all new commands.
- [ ] Task 5.2: Final PE Test and Build step to ensure stability.
- [ ] Task 5.3: Automated Release execution (Semantic Versioning + git push + tag).
