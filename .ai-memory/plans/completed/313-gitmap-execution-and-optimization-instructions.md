# Plan 313: GitMap Execution and Optimization Instructions

## User Request (Verbatim)
```text
PS C:\Users\Administrator> gitmap pa --ssh
...
  ▶ Node [w1] (IP: 192.168.1.3): failed
      (invalid JSON: The filename, directory name, or volume label syntax is incorrect.)

gitmap pa: execute failed: [E9000:EXECUTION] SSH fleet pull-all failed: [E_SSH_FLEET_FAIL:EXECUTION] SSH fleet pull-all had 1 node failure(s): (at=cmdssh/ssh_pull_fleet.go:63) (at=cmdpull/pull_dispatch.go:102)
...
Fix the GitMap execution issues. The update speed needs to be improved.
With GitMap, when we run gitmap pa or gitmap pa --ssh, this should not run multiple times on a single machine, or even if we are running on multiple nodes.
Anytime it is requested, it would be injected as a task into the GitMap SQLite database. A separate split database (pull.db)...
If already running, wait for SQLiteDB to be updated every 10 seconds, terminal ticker updates...
Write what it is pulling into a normalized sub-table... remove older than 7 days based on settings...
Do not keep files or diffs inside that database, only hash, title, created_at, pulled_at.
Rolling ETA based on last 7 to 8 pull-all runs. Nice tree view.
Store and encrypt private and public keys inside GitMap database for automatic re-login.
Heatmap command: let other agents know who else is working and roadmap. Add more examples to heatmap skill.
Why not using RC heatmap? Mandate rc heatmap.
AI writes model name and tool name (Gemini, Antigravity IDE, Muse).
Direct JSON string communication for agent and heatmap commands.
Pull all: dynamically seek more parallel workers when processor is idle.
Update: optimize download speed (in-memory streaming extraction).
Disable/enable feature for repositories: CLI (gitmap scan disable-repos / gitmap disable-scan) with numbered sequence selection (e.g. 1-5), REST API, and Web UI.
```

---

## 1. Disjoint Subtask Breakdown

| Task-ID | Subtask | Owner | Owned Files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | Architecture Specification & Subtask Planning | Spec Author | `02-spec/21-app/313-gitmap-execution-and-optimization-instructions/*`, `.ai-memory/plans/subtasks/313-gitmap-execution-and-optimization-instructions/*` | COMPLETED | Comprehensive architecture, component, migration specs, and 3 subtask plans generated and registered |
| Task-02 | Windows Remote SSH Quoting & Fleet Dispatch Fix | Worker 01 | `cli/cmdssh/ssh_pull_fleet_dispatch.go`, `cli/cmdssh/ssh_pull_json.go`, `cli/secrets/ssh_client.go` | COMPLETED | Fixed `The filename, directory name, or volume label syntax is incorrect.` on Windows OpenSSH nodes using Base64 UTF-16LE EncodedCommand wrapper and explicit shell resolution |
| Task-03 | Pull Task Split-DB, Lock Gate & Dynamic Rolling ETA | Worker 02 | `cli/constants/constants_pull.go`, `cli/cmdpull/pull_task_db.go`, `cli/cmdpull/pull_watcher.go`, `cli/cmdpull/pull_eta.go` | COMPLETED | Split SQLite DB `.gitmap/data/gitmap-pull.db` with WAL mode, single-execution lock gate, 10s terminal progress ticker, tree view, 7-day retention rotation, and rolling ETA across last 7-8 runs |
| Task-04 | Encrypted SSH Key Vault & Auto-Relogin Engine | Worker 03 | `cli/constants/constants_ssh.go`, `cli/model/sshkey.go`, `cli/store/sshkey.go`, `cli/cmdssh/ssh_vault.go`, `cli/cmdssh/ssh_relogin.go` | COMPLETED | AES-256-GCM encrypted private key persistence in SQLite, HMAC-SHA256 vault key derivation, and automated transparent disk key restoration upon auth expiry |
| Task-05 | Heatmap One-Shot JSON String Input & AI Telemetry | Worker 01 | `cli/heatmap/types.go`, `cli/heatmap/heatmap_cmd.go`, `cli/store/heatmap_schema.go`, `cli/heatmap/heatmap_render_json.go` | COMPLETED | Fixed `--json` collision for one-shot string input payloads, added `ModelName` and `ToolName` to SQLite schema, JSON renderer, and payload mapper |
| Task-06 | Heatmap Skills Authoring & Multi-AI Guidelines | Worker 02 | `.agents/skills/gitmap/skill.md`, `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`, `.agents/skills/muse-master-prompt/skill.md` | COMPLETED | Documented one-shot JSON, Antigravity IDE, Gemini, Muse examples, and mandatory RC roadmap traversal (`gitmap rc heatmap tree`, `gitmap heatmap search [slug]`) |
| Task-07 | CPU-Aware Adaptive Pull Concurrency Scaling | Worker 03 | `cli/cmdpull/pull_concurrency.go`, `cli/cmdpull/pull_worker_pool.go` | COMPLETED | Dynamic CPU idle detection scaling pull concurrency up to 16 workers locally while preserving conservative clamps on SSH sessions |
| Task-08 | Direct Binary Update Download Optimization | Worker 03 | `cli/cmdupdate/update_fallback_loop.go`, `cli/cmdupdate/update_binary_swap.go` | COMPLETED | Prioritized direct standalone `.exe` downloads over `.zip` archives, bypassing zip extraction overhead for high-speed updates |
| Task-09 | Repository Scan Disable/Enable CLI, REST API & Web UI | Worker 02 | `cli/cmdexclude/exclude_toggle.go`, `cli/cmdscan/scan_disable.go`, `cli/cmdui/ui_server.go`, `src/pages/RepoSettings.tsx`, `src/data/commands.ts` | COMPLETED | Implemented numbered range selection (`1-5, 8`), CLI commands (`gitmap scan disable-repos`, `gitmap disable-scan`, `gitmap enable-scan`), `/api/scan/*` REST endpoints, and React UI |

---

## 2. Invariants & Verification Evidence
- **Functions <= 15 lines:** All 130 Go functions across newly authored and modified files verified $\le 15$ lines max.
- **Positive Booleans:** Strictly `is*` and `has*` positive prefixes; zero negative booleans.
- **Strict Relative Paths:** Verified 0 absolute filesystem paths or `file:///` URIs across 9854 files via `check-relative-paths.py`.
- **Forbidden Strings:** Verified exit code 0 via `check-forbidden-strings.py`.
- **Zero Builds / Zero Tests:** Zero `go build` or `go test` invocations executed during task execution.
- **Database Hygiene:** Dedicated Split-DB (`gitmap-pull.db`) with zero file diffs/blobs, storing only commit hash, title, timestamps, and repo metrics, with automatic 7-day retention rotation.
