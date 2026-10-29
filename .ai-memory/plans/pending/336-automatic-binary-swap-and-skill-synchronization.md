# Master Plan: Automatic Binary Swap and Skill Synchronization

**Sequence ID:** 336  
**Task Slug:** `automatic-binary-swap-and-skill-synchronization`  
**Specification:** [02-spec/21-app/336-automatic-binary-swap-and-skill-synchronization/01-architecture-spec.md](file;02-spec/21-app/336-automatic-binary-swap-and-skill-synchronization/01-architecture-spec.md)  
**Status:** Pending  
**Created:** 2026-10-27  

---

## 1. Executive Summary & Problem Statement

When compiling or updating GitMap binaries in place on Windows and POSIX systems:
1. Executing binaries are locked by the operating system (`ERROR_SHARING_VIOLATION` on Windows), preventing in-place file replacement while the CLI, daemons, or IDE tools are running.
2. Background sibling processes (e.g., daemon workers, file monitors, or terminal listeners) may hold stale code in memory or contend for file access.
3. Antigravity skills, shared specifications, and prompt libraries across connected repositories require continuous synchronization with the master guidelines.

Plan 336 introduces:
- A new core command `gitmap swap-binary` (aliases: `swap-bin`, `bin-swap`, `replace-bin`, `replace-binary`, `binary swap`, `swap`) that discovers locking processes, generates detached handoff scripts (`.ps1` / `.sh`), executes safe atomic file swaps, verifies the newly installed binary with `--version`, and reports results.
- Comprehensive CLI registration, root dispatcher routing, flag parsing, and themed help documentation.
- Coding guideline rules and skill synchronization across `.agents/skills/` and connected workspace repositories.

---

## 2. Subtasks Roadmap

| Subtask ID | File | Description | Target | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | [01-swap-core-and-script-generator.md](.ai-memory/plans/subtasks/336-automatic-binary-swap-and-skill-synchronization/01-swap-core-and-script-generator.md) | Implement core swap models, PID discovery, detached handoff script generation (`.ps1`/`.sh`), atomic file replacement, and post-swap verification engine. | `cli/cmdswap/`, `cli/cmdupdate/` | Pending |
| **Subtask-02** | [02-cli-dispatch-and-helpdoc.md](.ai-memory/plans/subtasks/336-automatic-binary-swap-and-skill-synchronization/02-cli-dispatch-and-helpdoc.md) | Register `gitmap swap-binary` and aliases in `cli/root.go`, implement CLI command flags, Lipgloss summary cards, and add themed help documentation in `cli/helpdoc/`. | `cli/root.go`, `cli/cmdswap/`, `cli/helpdoc/` | Pending |
| **Subtask-03** | [03-coding-guidelines-skill-sync.md](.ai-memory/plans/subtasks/336-automatic-binary-swap-and-skill-synchronization/03-coding-guidelines-skill-sync.md) | Author new Antigravity skill `gitmap-binary-swap`, update coding guidelines, and synchronize skills and prompts across connected repositories using script 38. | `.agents/skills/`, `.agents/rules/`, `03-ai-scripts/` | Pending |

---

## 3. Architecture & Verification Gate

### Boolean & Style Invariants:
- All booleans must use positive naming starting with `is` or `has` (`isScriptOnly`, `isSuccess`, `hasKill`, `hasRestart`, `hasForce`, `hasTargetPID`). No negative booleans.
- Functions must stay within 8–15 lines with early guard clauses and zero nested `if` statements.
- Return styling requires blank line after `}` and blank line before `return`.
- All paths referenced in specifications and plans must use 100% relative paths.

### Verification Checklist:
- [ ] `gitmap swap-binary --help` renders complete flags and usage documentation.
- [ ] `gitmap swap-binary --script-only` emits valid, platform-appropriate PowerShell or Bash script without terminating current process.
- [ ] Process lock discovery correctly detects sibling processes and respects `--force` / `--kill` flags.
- [ ] Unit tests for `cmdswap` package achieve 100% pass rate.
- [ ] `gitmap-binary-swap` skill is registered in `.agents/skills/` and synchronized.
