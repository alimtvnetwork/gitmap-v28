---
name: cg-clean-repo-build-and-caches
description: Autonomously clean repository build issues, purge temporary directories, OS temp folders, Golang build/test caches, pnpm/npm stores, Vite/TS incremental caches, and polyglot caches while enforcing lowercase hygiene and .gitignore sync.
---

# Skill: Clean Repository Build Issues & Polyglot Caches (`cg-clean-repo-build-and-caches`)

> **Invocation:** `@[cg-clean-repo-build-and-caches](file;.agents/skills/cg-clean-repo-build-and-caches)`

This skill governs autonomous diagnosis and systematic cleanup of repository build issues, stale build artifacts, temporary directories, and polyglot caches across OS, Go, pnpm/npm, Vite, Python, and Rust toolchains.

## 1. Concurrency Invariants & Header Parameters

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Fleet Sync, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

[/goal](slashCommand;goal) Autonomously diagnose and clean the repository to resolve build issues, purge stale build artifacts, and systematically clean all temporary directories and multi-language caches — including repository temp folders, OS temp directories, Golang build/test caches, pnpm/npm/Vite caches, Python/Rust/PHP caches, and untracked artifact pollution — with strict no-build and no-test execution (NEVER run heavy build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine turns; compilation and testing are verified later in CI/CD). Spawn autonomous subagents (A = 3, H = 2, Capacity = 6) for parallel inspection and cache purging, leverage GitMap high-speed tooling as primary, synchronize `.gitignore` rules, and finalize any repository hygiene fixes with an atomic push via GitMap hyphen commit.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom cache targets, build issue symptoms, or user instructions are provided ABOVE this prompt (in the user preamble or header above) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Execute the multi-layer cache and temp cleanup protocol across repository, OS, Go, pnpm/npm, and polyglot toolchains, and persist hygiene findings into `.ai-memory/plans/`.

---

## 2. Multi-Layer Cache & Temp Cleanup Protocol

1. **Layer 1 — Repository & Work Directory Build Artifacts:**
   - Run `python 03-ai-scripts/44-work-and-system-cache-cleaner.py --dry-run` (Plan preview) and `python 03-ai-scripts/44-work-and-system-cache-cleaner.py -y` to clean work build elements, Go/npm/pnpm/Node/DevTools caches, OS temp, `SoftwareDistribution\Download`, Recycle Bin, and Git caches.
   - Remove `.tmp/`, `tmp/`, `__pycache__/`, `.pytest_cache/`, `.mypy_cache/`, `.ruff_cache/`, `.vite/`, `.turbo/`, `.parcel-cache/`, `*.pyc`, `*.pyo`, `*.tmp`, `*.test`, `coverage.out`, `tsconfig.tsbuildinfo`, `.eslintcache`.
   - Run `python 03-ai-scripts/19-artifact-remover.py --execute` (if present).
2. **Layer 2 — OS Temporary Directories, SoftwareDistribution & Recycle Bin:**
   - Clean stale `$env:TEMP\go-build*`, `$env:TEMP\vite*`, `$env:TEMP\npm-*`, `/tmp/go-build*`, `%WINDIR%\SoftwareDistribution\Download`, and Recycle Bin (`Clear-RecycleBin -Force`).
   - Run `gitmap update-cleanup`.
3. **Layer 3 — Golang Caches (`GOCACHE` & `GOMODCACHE`):**
   - Execute `go clean -cache -testcache -fuzzcache -modcache` and `golangci-lint cache clean`.
4. **Layer 4 — pnpm, npm, Node.js & DevTools Caches (`node_modules` preserved):**
   - Execute `pnpm store prune`, `npm cache clean --force`, and remove `node_modules/.cache`, `node_modules/.vite`, `.pnpm-store`, and Browser/VSCode DevTools caches.
5. **Layer 5 — Git Hygiene & Lowercase Enforcement:**
   - Run `gitmap commons` to sync `.gitignore` and `.gitattributes`.
   - Run `gitmap lcf` and `gitmap lowercase-readme` to eliminate case-sensitivity build blockers.
   - Run `gitmap storage` (`gitmap stor`) to verify disk space recovery.

---

## 3. Step 0 Preflight with Data Schema Verification

Before enqueuing subtasks or claiming work, initialize and verify the SQLite task database schema:

```bash
gitmap agent task init --name "<task name>" --budget 400
gitmap agent task status --ddl
gitmap agent task status --json
```

- **DDL & Schema Verification Protocol:** Inspect tables (`ParentTask`, `Subtask`, `AgentActionLog`), affirmative boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and foreign key linkages.
- **Forensic Autopsy Check:** Check `diagnostics.hasCrashesDetected`. If true, inspect `diagnostics.crashedAgents` to identify failed subtasks, targeting files, and last logged actions to facilitate clean resume.

---

## 4. 3-Tier Database Location Transparency & Diagnostic Suite

Whenever executing or logging subtask telemetry, fulfill the Database Location Transparency Contract:

| Database Tier | Physical File Location | Access Mode & Concurrency Protocol | Data Retained & Ownership |
| :--- | :--- | :--- | :--- |
| **Tier 1: Machine-Local Root Database** | `~/.gitmap/agent-root.db` | Machine sessions, cross-repo leases | Machine Sessions & Cross-Repo Leases |
| **Tier 2: Repo-Cache Companion Storage** | `repo-cache/.heatmap/<ai-slug>/<repo-slug>/heatmap.db` | Isolated Split-DB with periodic sync (`gitmap heatmap rc sync`) | Cross-repository multi-agent actions, tool usage traces, commit checkpoints |
| **Tier 3: In-Flight Local Run Database** | `.ai-memory/temp-agents/<slug>/agent-task.db` | WAL Mode (`PRAGMA journal_mode=WAL; busy_timeout=5000`) | Task status, worker claims, subtasks, in-flight action logs (<1ms Go engine) |
| **Tier 3: Local Heatmap Cache** | `.ai-memory/heatmap.db` | WAL Mode + Read-Only Queries (`mode=ro`) | Cognitive action telemetry, file exploration reads/edits, contention leases |

### Executable Diagnostic Commands:

```bash
# Pre-edit contention and lease check
gitmap heatmap discuss [--window 15m] [--threshold 2] [--json]

# Cognitive file exploration and cleanup telemetry
gitmap heatmap log -t <task-slug> -a "READ" -f "<path>" -r "Reasoning: <explanation>" -s "EXPLORATION"
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <purge-rationale>" -s "EXECUTION"

# Real-time continuous activity stream & tree
gitmap heatmap tree [-t <task-slug>]
gitmap heatmap live

# Read-only companion query and schema inspection
gitmap heatmap rc query "SELECT ActionType, FilePath, ReasoningText FROM AgentActionLog ORDER BY AgentActionLogId DESC LIMIT 5" --format table
gitmap heatmap rc tables
gitmap heatmap rc view AgentActionLog --limit 10
gitmap heatmap rc sync
```

---

## 5. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine turns. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run test suites (`go test ./...`, `pytest`) during routine execution.
- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`) after individual file edits.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 6. Verification Checklist

- [ ] All 5 layers of cache and temp targets systematically audited and purged.
- [ ] Strict lowercase file hygiene maintained (`gitmap lcf`).
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
