# Consolidated Completed Plan: 174 — Antigravity Decision Log DB, Comprehensive Project Discovery, and Active Rerun Recency

## Header & Provenance
- **Canonical Spec Reference:** [02-spec/21-app/174-agy-decision-log-db-project-discovery-and-recent-rerun.md](../../../02-spec/21-app/174-agy-decision-log-db-project-discovery-and-recent-rerun.md)
- **Status:** `completed`
- **Total Execution Steps / Loops:** 1 continuous dual-subagent iteration (A=2, H=2, 0 blockers)

---

## Consolidated Subtasks & Implemented Deliverables

### Subtask 01: Antigravity Decision Audit SQLite Split-DB & CLI (`gitmap agy log`)
- **Traceability ID:** Task-03, Task-04
- **Target Files:**
  - `cli/store/agy_log_split_db.go`
  - `cli/cmdagy/agy_log_cmd.go`
  - `cli/cmdagy/agy_log_render.go`
  - `cli/helptext/agy-log.md`
  - `cli/helptext/catalog.go`
  - `cli/helptext/agy.md`
- **Accomplished Changes:**
  - Implemented `OpenAgyLogSplitDB()` managing SQLite Split-DB `gitmap-agy-log.db` located under `store.BinaryDataDir()`.
  - Created table `AgyDecisionLog` (`LogId`, `Command`, `TargetProject`, `ProjectPath`, `ConversationId`, `DecisionReason`, `Status`, `Node`, `CreatedAt`) with indexes on `CreatedAt DESC`, `TargetProject`, and `Command`.
  - Implemented `InsertAgyDecisionLog()` and `QueryAgyDecisionLogs()` supporting filtering by project, command, node, and limit.
  - Implemented `RecordAgyDecision(cmd, targetProj, projPath, convID, reason, status)` helper and hooked into `RerunProject`, `processAgyLsProjects`, and `runAgyInjectPrompts`.
  - Implemented `AgyLogCmd` (`log`, `logs`, `decision-log`, `decision-logs`) supporting `--json`, `--limit/-n 20`, `--project/-P`, `--command/-c`, and `--ssh` to query remote cluster nodes via SSH and merge logs chronologically with deduplication.
  - Authored comprehensive documentation in `cli/helptext/agy-log.md` and registered in `cli/helptext/catalog.go`.

### Subtask 02: Multi-Source Project Discovery & Active Conversation Recency Sorting
- **Traceability ID:** Task-02
- **Target Files:**
  - `cli/cmdagy/agy_multi_discovery.go`
  - `cli/cmdagy/agy_ls.go`
  - `cli/cmdagy/agy_rerun_project_resolve.go`
  - `cli/cmdagy/agy_rerun_restart.go`
- **Accomplished Changes:**
  - Implemented `DiscoverUnifiedAgyProjects(dirPath)` aggregating projects across:
    1. `~/.gemini/config/projects/*.json`
    2. `conversation_summaries.db` (`workspace_uris`, `title`, `last_modified_time`)
    3. Candidate workspaces and running prompt files (`active-agy-pipeline-fix-prompt.txt`)
  - Changed `gitmap agy ls` default count from 8 to 0 (unlimited/all projects) so all projects on the system are visible.
  - Added `--ssh` flag to `gitmap agy ls` to query remote cluster nodes (e.g. `w1`) and merge remote projects into table/JSON output.
  - Updated `sortProjectsByActivityAndPins` and `compareProjectEntries` in `cli/cmdagy/agy_rerun_project_resolve.go`:
    - Projects hosting active prompt files (`active-agy-pipeline-fix-prompt.txt`) or active running status rank highest.
    - True latest activity timestamp (`lastActTime`) is prioritized over inactive working directories so `gitmap rerun` targets genuinely active projects (such as `white-presentation-v1` and `rasia-logo`) rather than stale CWDs.
    - Removed `"1"` from `isCwdTarget` so `gitmap rerun 1` unambiguously selects `projects[0]` without being hijacked by CWD.
    - Added `isProjectActiveOrRunning` check in `findCwdProjectOrDefault` so idle CWD does not hijack `gitmap rerun` when other active projects exist.

### Subtask 03: Rerun Help Routing & Substring Guard Hardening
- **Traceability ID:** Task-01
- **Target Files:**
  - `cli/cmdagy/agy_rerun_help.go`
  - `cli/cmdagy/agy_rerun.go`
  - `cli/cmdagy/agy_rerun_project_resolve.go`
  - `cli/cmd/root.go`
- **Accomplished Changes:**
  - Exported `IsRerunHelpRequested` and `IsRerunHelpToken` in `cli/cmdagy/agy_rerun_help.go`.
  - Added top-level guard in `cli/cmd/root.go` dispatch to immediately invoke `RenderAgyRerunHelp()` whenever `args[1:]` requests help.
  - Hardened `resolveClosestActiveProject`, `findProjectByFlexibleTarget`, and `isProjectTargetMatch` to strictly reject help tokens and never substring-match projects (e.g. `strhelper`).

### Subtask 04: SSH Cluster Pre-Flight TCP Probing
- **Traceability ID:** Task-03, Task-04
- **Target Files:**
  - `cli/cmdagy/agy_multi_discovery.go`
  - `cli/cmdagy/agy_log_cmd.go`
- **Accomplished Changes:**
  - Implemented `probeNodeOnline(host, port, timeout)` with a 400ms TCP probe.
  - Wired into `querySingleNodeProjects` and `querySingleNodeAgyLogs` so unreachable nodes (e.g. `alpha-win`, `beta-linux`, `gamma-mac` at 10.20.0.x) are skipped immediately, preventing 15-30s dial timeouts during `--ssh` discovery and log queries.
