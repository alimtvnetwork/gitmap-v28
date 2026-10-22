---
name: cg-enum-standards
description: Autonomously audits and enforces cross-language enum standards (Go, TS, PHP, Rust), ensuring PascalCase formatting, Type suffixes, and generator workflow usage.
---

# Enum Standards Coding Guidelines (`cg-enum-standards`)

This skill provides autonomous audit and validation of repository-wide enum definitions based on `02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines/03-boolean-and-enum.md` and `02-spec/17-consolidated-guidelines/07-enum-standards.md`.

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

---

## 2. Core Architectural Invariants

1. **Source of Truth**: Enums are defined ONCE in YAML specifications (`02-spec/<module>/enums/<EnumName>.yaml`).
2. **Generator Workflow**: Must use `03-ai-scripts/30-enum-generator.py` to generate enums for Go, TypeScript, PHP, Rust. Never hand-edit `*_generated.*` files directly. Generated enum output lands in the shared types/constants module; enums are never defined inline in a handler or test.
3. **Enum Naming & Type Suffix**:
   - Enums and members must strictly use **PascalCase**.
   - All enum type identifiers across TypeScript, Go, and PHP MUST carry the `Type` suffix (e.g. `UserRoleType`, `CommandStatusType`, `PaymentMethodType`).
4. **String-Backed Enums**: Prefer string-backed enums over numeric enums for cross-language compatibility, serialization clarity, and debugging ergonomics. Serialized values are explicit and stable — never change a shipped value. Where the language needs it, an explicit unknown/zero sentinel exists and is never a valid action.
5. **Elimination of Magic Literals**: Zero inline magic strings or raw integers in business logic. All known values must reference the typed enum constant.
6. **Rust Derives**: Rust enums must include `Debug, Clone, PartialEq, Serialize, Deserialize`.
7. **MIME & Media Type Enum Standard (Spec 36)**: Zero raw magic string literals for MIME or media types (e.g. `'text/html'`, `'text/plain'`, `'application/json'`, `'image/png'`). All MIME types used in Blob constructors (`new Blob([data], { type: MimeType.TextHtml })`), ClipboardItem APIs (`new ClipboardItem({ [MimeType.TextHtml]: htmlBlob })`), HTTP headers (`Content-Type`, `Accept`), or API responses MUST use centralized `MimeType` enum members or typed constant dictionaries.
8. **Map Dispatch for Extension Resolution (Spec 35)**: Extension-to-MIME lookups MUST use declarative map dispatch tables (`Record<string, MimeType>`) rather than `if-else` or `switch` statements. (Data lookups only; action dispatch follows item 9.)
9. **Enum-Driven Action Dispatch**: Every action, operation, command kind, and event kind is a member of a named enum. Dispatch ONLY by exhaustive switch/match on the enum — never by raw-string or magic-number comparison. The `default` branch raises a typed error (see `cg-error-management`); never a silent no-op. Adding an action = adding one enum member + one switch case; no other registration step.
10. **IsValid Before Dispatch**: Before dispatch, check the action with the enum's `IsValid` / membership util. Unchecked external input never reaches the switch. External input enters ONLY through the enum's parse util (item 11).
11. **Util-Methods Checkup (mandatory, never deferred)**: Every enum ships `String()` / label, `Parse` / `TryParse` (typed error on unknown input), `IsValid`, and `All` / `Values` — or a documented exception. Any enum missing a util gets it added in the same change that touches the enum.

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
gitmap agent heatmap --json
gitmap heatmap status

# Cognitive file exploration and editing telemetry
gitmap heatmap log -t <task-slug> -a "READ" -f "<path>" -r "Reasoning: <explanation>" -s "EXPLORATION"
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <refactor-rationale>" -s "EXECUTION"

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

## 5. Fast File Discovery & Reading Toolchain (GitMap Primacy)

All symbol and literal searches MUST use GitMap. Shell search tools (`rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`) are strictly banned:

- **Discover Un-typed Constants:** `gitmap aum search "const (" [dir] -e .go`
- **Find Raw String Literals:** `gitmap aum search -r "== \"[A-Za-z0-9_-]+\"" [dir] -e .go`
- **Find Enum Declarations:** `gitmap aum search "type .*Type string" [dir] -e .go`
- **Find Raw MIME & Media Type Literals:** `gitmap aum search -r "(?:'|\")(text|application|image|audio|video)\/(html|plain|json|xml|markdown|csv|png|jpeg|octet-stream)(?:'|\")" [dir]`
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]`
- **Stream File Content:** `gitmap cat <filepath>`

---

## 6. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine turns. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run test suites (`go test ./...`, `pytest`) during routine execution.
- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 7. Verification Checklist

- [ ] All enums carry `Type` suffix across TypeScript, Go, and PHP.
- [ ] Enums use PascalCase naming for types and member constants.
- [ ] Magic literals replaced with typed enum constants.
- [ ] MIME & Media Types extracted into centralized `MimeType` enums for Blob constructors, ClipboardItem APIs, HTTP headers, and API responses (Spec 36).
- [ ] File extension to MIME resolution uses declarative map dispatch tables (Spec 35).
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap agent heatmap (or gitmap heatmap status) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
- [ ] Every action/operation/command/event kind is an enum member, dispatched by exhaustive switch/match; `default` raises a typed error.
- [ ] `IsValid` checked before dispatch; external input enters only via `Parse`/`TryParse`.
- [ ] Util-methods checkup passed: `String()`, `Parse`/`TryParse`, `IsValid`, `All`/`Values` present (or documented exception) for every touched enum.
