# Specification: 174 — Antigravity Decision Log DB, Comprehensive Project Discovery, and Active Rerun Recency

## Status
`active`

## Domain
`cli / agy / decision audit db / project discovery / brain session scan / remote ssh aggregation / rerun recency`

## User Request (Verbatim)
```text
PS C:\Users\Administrator> gitmap rerun help

Antigravity IDE Rerun Suite:
  • Target Project:   #37 - strhelper
  • Workspace Path:   ./03-aukgo\strhelper
  • Conversation:     default


  ✓ Injected prompt into active Antigravity session (418bd745-3afa-4ea4-8614-3254f8f860e8) via agentapi!
  ✓ Completed prompt replay for project #37 (strhelper)!

PS C:\Users\Administrator> gitmap rerun 1

Antigravity IDE Rerun Suite:
  • Target Project:   #1 - white-presentation-v1
  • Workspace Path:   ./presentations-repos\white-presentation-v1
  • Conversation:     8e34ebc6-e0e2-4182-96d3-809d5d5065ac


  ✓ Injected prompt into active Antigravity session (8e34ebc6-e0e2-4182-96d3-809d5d5065ac) via agentapi!
  ✓ Completed prompt replay for project #1 (white-presentation-v1)!

PS C:\Users\Administrator> gitmap rerun 2

Antigravity IDE Rerun Suite:
  • Target Project:   #2 - rasia-logo
  • Workspace Path:   ./presentations-repos\rasia-logo
  • Conversation:     default


  ✓ Injected prompt into active Antigravity session (80a04ace-9539-4cdf-9433-5c4bb025a9fc) via agentapi!
  ✓ Completed prompt replay for project #2 (rasia-logo)!

PS C:\Users\Administrator>

PS C:\Users\Administrator> gitmap agy ls

  ╔══════════════════════════════════════╗
  ║         antigravity projects         ║
  ╚══════════════════════════════════════╝

  8 projects from ~/.gemini/config/projects
  ──────────────────────────────────────────────────────────────────────

   ~/Documents/antigravity (1 projects)
  SEQ   CONV NAME          ID           PROJECT               PATH
  ──────────────────────────────────────────────────────────────────────
  048   —                  9a5f0c7c     adventurous-galileo   ~/Documents/antigravity/adventurous-galileo

   d:\work (5 projects)
  SEQ   CONV NAME          ID           PROJECT               PATH
  ──────────────────────────────────────────────────────────────────────
  002   cv                 e714aa58     alim-cv               ./alim-cv
  003   profile            2862a47e     alim-karim-profile    ./alim-karim-profile
  004   status             96fb39c3     alim-status-sample    ./alim-status-sample
  005   —                  57f0b96a     Antigravity-Manager   ./Antigravity-Manager
  007   cat                14e40fc0     cat-my-ui-v12         ./cat-my

   ./02-prompts (1 projects)
  SEQ   CONV NAME          ID           PROJECT                   PATH
  ──────────────────────────────────────────────────────────────────────
  001   —                  1a65a4ff     ai-empathy-prompt-tuner   ./02-prompts\ai-empathy-prompt-tuner

   ./presentations-repos (1 projects)
  SEQ   CONV NAME                  ID           PROJECT                    PATH
  ──────────────────────────────────────────────────────────────────────
  006   Memory Ret...t Ingestion   d2dcb704     bsrm-prese...ion-hiltrax   ./presentatio...presentation-hiltrax


  ──────────────────────────────────────────────────────────────────────
  8 projects · 8 active · 0 missing

  ✓ Antigravity (AGY): No duplicate project paths found.

PS C:\Users\Administrator>

So this was completely wrong what you have run here. So it shows the latest project is the R Asia, and the presentation has running. So when I say rerun, it should basically going to rerun on the active projects and projects which has recent activities, right? But it turns out, even after saying all this stuff, it actually is running the project into the git map, which is the very old conversation. I'm not sure if you could understand and get it. The reason I'm saying this because even the git map LS does not show all the projects that we have in the system. It's very terrible. I think you need to work on it. Git map LS, AGY LS looks really bad. It needs improvement. It actually is coming from the W1 machine. So you can do remote access using SSH and try to use the git map to get its log and stuff like that. Currently, for each one of these stage, like AGY, what it does, how it does it, I want you to create some logs. You can have some AGY DB where you can actually keep the logs based on why you are taking each decision. So if anything goes wrong, we can say git map AGY log or logs, and it will show the logs from the SSH or terminal or as a JSON as well. So all these other flags needs to be there. So I think this is what you need to integrate to understand this and also try to fix it if it possible. Do you understand the task? Can you please fix it for me? And make a minor bump, please, and release
```

## Architectural Design & Invariants

### 1. Antigravity Decision Audit SQLite Split-DB
- Store logs in `gitmap-agy-log.db` located under `store.BinaryDataDir()`.
- Table `AgyDecisionLog`:
  - `LogId` (TEXT PRIMARY KEY)
  - `Command` (TEXT NOT NULL) — e.g. `rerun`, `ls`, `inject`, `watch`
  - `TargetProject` (TEXT NOT NULL) — e.g. `white-presentation-v1`, `rasia-logo`
  - `ProjectPath` (TEXT NOT NULL) — Workspace directory path
  - `ConversationId` (TEXT NOT NULL) — Target conversation UUID
  - `DecisionReason` (TEXT NOT NULL) — Human-readable reason (e.g. "most recent conversation 8e34ebc6 updated 2026-09-27 14:10:12")
  - `Status` (TEXT NOT NULL) — `success`, `dry-run`, `error`, `skipped`
  - `Node` (TEXT NOT NULL) — Hostname or remote IP
  - `CreatedAt` (TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP)
- CLI Commands:
  - `gitmap agy log` / `gitmap agy logs`
  - Flags: `--json`, `--limit/-n 20`, `--project/-P <name>`, `--ssh` (retrieves remote logs from `w1` and other cluster nodes).

### 2. Multi-Source Comprehensive Project Discovery
- Current discovery is restricted to `~/.gemini/config/projects/*.json` which misses projects like `white-presentation-v1` and `rasia-logo`.
- Expand discovery to inspect:
  1. `~/.gemini/config/projects/*.json`
  2. `~/.gemini/antigravity/brain/*/` transcript headers (`workspaceUris`)
  3. GitMap tracked repositories matching active Antigravity conversations
  4. Remote node aggregation via SSH (`gitmap agy ls --ssh` queries `w1` and cluster nodes).

### 3. Active Conversation Recency Sorting
- Sort projects by true latest conversation modification timestamp (`updated_at` / latest transcript line).
- The most recently active project is ranked `#1` so `gitmap rerun` from any directory targets the genuinely active project.
