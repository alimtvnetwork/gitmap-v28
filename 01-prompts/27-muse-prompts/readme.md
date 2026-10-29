# Muse Prompts (`27-muse-prompts`) — Index & Catalog

> [!IMPORTANT]
> Category: `27-muse-prompts`
> Architecture: Self-Bootstrapping Autonomous Onboarding & Fleet Execution Directive (V6 Engine)
> Scope: Muse AI Session Bootstrapping, Authentication (SSH Primacy), Fleet Scan, Auto-Cloning, Pull-All, Operating Contract, and Multi-Agent Execution

---

## Overview

The `27-muse-prompts` prompt category houses the unified master onboarding and execution prompt for
Muse AI sessions. Pasting the master prompt into a fresh Muse chat boots a
fully-onboarded autonomous agent:
1. **Authentication Bootstrap with SSH Primacy:** Fast-path bypass if already logged in; preferred SSH configuration (`gh auth login -p ssh -w`) or browser OAuth fallback.
2. **GitMap Installation, Self-Update & Fleet Scan:** Auto-installs GitMap via `install-quick.sh`, updates GitMap first, and executes `gitmap scan` across the workspace to register all repositories and branches into SQLite (`gitmap.db`).
3. **Repository Fleet Orchestration & Cloning:** Canonical starting point URL `https://github.com/alimtvnetwork/coding-guidelines-v24.git` (with explicit ban on stale `coding-guidelines`), 14+ repository registry, auto-cloning missing repos, pull-all via `gitmap pa` ensuring conflict-free trees, and establishing dedicated side chats (`<repo> repo` and "all summary").
4. **Memory Operating Contract Checklist:** Turbo mode (act, don't ask), strict no-delete guards (repos are sacred, file deletion always confirms), CI-only test verification (`gitmap pe -t`).
5. **Coding Guidelines, Formatting & Large File Split:** Booleans (`is*`/`has*`, no `== true`), early returns, `*appfault.AppError`, lowercase filenames, relative paths, large-file 15 MB 7-Zip multi-volume split (`7z a -v3m -mx=0`), and title convention (`# <title>: high priority instruction, non-negotiable task`, `## slug: <task-slug>`).
6. **Task Confirmation & Multi-Agent Execution Protocol (V6 Engine):** Verbatim capture, confirmed task breakdown FIRST (`Task-01`, `Task-02`, ... each with `Understood: [YES]`), RUNNING + ETA declaration, 5-minute status pings, SQLite task DB init (`gitmap agent task init`), pull first (`gitmap pull` / `git pull`), mandatory multi-agent dispatch (`invoke_subagent`, A=3, H=2, capacity 6 concurrent operations, dynamic respawn, disjoint file boxes, Worker Git Ban).
7. **Verification, Pre-Commit Secrets Gate & Atomic Push:** Targeted checks only (R1 zero build/test ban), pre-commit secrets gate (`check-forbidden-strings.py` and `gitmap aum search` -> `gitmap rs`), atomic commit via GitMap hyphen format (`gitmap cpf "<module> - <summary>"`), immediate push, and CI telemetry (`gitmap pe -t`).
8. **Phase 7.5 Git Repository Troubleshooting & Repo-Cache Forensics:** Non-destructive Git remediation (total ban on force-push/rebase, Lovable history protection, Worker Git Ban), divergence reconciliation, mechanical conflict solvers, index lock clearance, debris isolation, and conditional forensics persistence in `repo-cache/muse/troubleshooting/<task-slug>/` via `gitmap rc`.
9. **Muse Memory IO (Memory Re-sync Protocol):** 7-day memory digest under `12-muse-ai/` in `repo-cache`, re-sync via `gitmap llm train muse --days 7`, chat-seeds recreation, and standing SMTP pre-approval for `go-email-reader`.

The canonical execution engine it embodies is
[`02-execute-parent-task-with-n-steps-v6.md`](../14-execute/02-execute-parent-task-with-n-steps-v6.md)
(Canonical V6: `N = 300`, `A = 3`, `H = 2`, `C = 30`).

Its installable skill is `muse-master-prompt`
(`.agents/skills/muse-master-prompt/skill.md`, mirrored in `.cursor/skills/`).

---

## Directory Index & Catalog

| Prompt File | Version | Scope | Key Capabilities | Status |
| :--- | :--- | :--- | :--- | :--- |
| [`01-muse-master-prompt.md`](01-muse-master-prompt.md) | 6.5.0 | Unified Master Prompt | Consolidated V6 directive: SSH auth primacy + fast-path bypass, GitMap quick-install & update first, `gitmap scan`, 14+ repo fleet clone & pull-all (`gitmap pa`), conflict prevention, dedicated side chats (`<repo> repo`, `all summary`), repo-cache sequence, task confirmation first, RUNNING+ETA, 5-minute pings, multi-agent dispatch (A=3, H=2, capacity 6 operations, dynamic respawn), Worker Git Ban, Phase 7.5 non-destructive Git troubleshooting & repo-cache forensics, 15 MB 7-Zip split, GitMap hyphen atomic push, CI telemetry (`gitmap pe -t`), and Muse Memory IO. | **Active (Canonical)** |
| [`02-muse-generic.md`](02-muse-generic.md) | 6.5.0 | Generic Master Prompt (GitMap Suite) | Clean, public, generic master prompt for any codebase: zero hardcoded repos/secrets, dynamic repository onboarding gate, cross-platform GitMap installation & update suite, full Muse capabilities (code fixing, bulk operations, scheduling, email alerts, content, image/video asset generation), canonical directory architecture (`WORKSPACE_ROOT = ~/workspace/repos/`, `02-spec/21-app/XX-<slug>/`), and companion first-time/daily runner scripts. | **Active (Canonical)** |
| [`03-muse-standalone-master-prompt.md`](03-muse-standalone-master-prompt.md) | 6.5.0 | Standalone Universal Master Prompt (Pure Native Git) | Pure native Git standalone master prompt for general users outside the company: zero `gitmap` dependency, zero hardcoded repos/secrets, interactive Git connection ritual (prompt for repo URL, clone or init), Full Turbo Mode ("Act, don't ask", write-and-execute), recommended yet adaptable folder structure, standard Git commit workflow, and companion Bash/PowerShell runner scripts. | **Active (Canonical)** |
| [`04-muse-universal-master-prompt.md`](04-muse-universal-master-prompt.md) | 6.5.0 | Universal Master Prompt for Everyone (Following 01 Architecture) | Comprehensive 9-step universal onboarding prompt following the full architecture of `01-muse-master-prompt.md`: dynamic Git repository onboarding gate (asks for repo URLs/local path, clones/connects, and auto-starts), V6 autonomous execution engine, mandatory subagent swarm (A=3, H=2), TURBO-01..05 & GUARD-01..09 operating contract, adaptive folder layout, and Phase 7.5 Git troubleshooting runbooks. | **Active (Canonical)** |

---

## Directory Tree

```text
01-prompts/27-muse-prompts/
├── 01-muse-master-prompt.md            # UNIFIED master prompt (internal fleet + GitMap)
├── 02-muse-generic.md                  # GENERIC master prompt (GitMap suite)
├── 03-muse-standalone-master-prompt.md # STANDALONE universal master prompt (pure native Git, no GitMap)
├── 04-muse-universal-master-prompt.md  # UNIVERSAL 9-step master prompt for everyone (following 01 + dynamic Git auto-start)
└── readme.md                           # this index
```
