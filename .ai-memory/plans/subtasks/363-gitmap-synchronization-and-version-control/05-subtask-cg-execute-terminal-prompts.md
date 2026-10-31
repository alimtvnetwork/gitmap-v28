# Subtask 05: On-The-Fly Terminal Execution Prompts (cg execute v6 / v7)

- **Parent Plan:** [363-gitmap-synchronization-and-version-control](../../pending/363-gitmap-synchronization-and-version-control.md)
- **Subtask ID:** `TASK-05`
- **Slug:** `gitmap-synchronization-and-version-control - Task 05`
- **Target Package:** `cli/cmdcg/`, `cli/cmdprompt/`
- **Status:** `pending`
- **Assigned Worker:** Prompt Architecture Subagent

---

## 1. Goal

Implement the on-the-fly execution prompt generator within the GitMap CLI. Enable developers and AI agents to instantly generate, view, and clipboard-copy standardized V6 (Tree Swarm Hierarchy) and V7 (Time-Restricted & Autobot Governed) prompt envelopes using `gitmap cg execute v6` and `gitmap cg execute v7`.

---

## 2. In Scope

1. **Subcommand Dispatcher (`cg_execute_cmd.go`):**
   - Register `gitmap cg execute [v6|v7]` under `cmdcg`.
   - Support flag equivalents: `--v6`, `--v7`, `-6`, `-7`.
   - Default to V7 when version flag is omitted.

2. **Prompt Template Engine (`prompt_template_engine.go`):**
   - Ingest canonical execution prompts:
     - V6 Engine: Tree Swarm Hierarchy (Architect=3, Hunter=2, Capacity=6), micro-batches 5-8 files, GitMap hyphen-separated atomic commits, SQLite task manager.
     - V7 Engine: Dual operating modes (Option 1: Time-Restricted with timezone & weekday `< 08:00 AM` cutoff; Option 2: Unrestricted), single-user bot identity, secrets gate, ledger resume.
   - Dynamic parameter interpolation:
     - `--task <task-slug>`: Inject current task name into headers and target files.
     - `--batch <size>`: Configure micro-batch size (defaults to 6).
     - `--agent <name>`: Pre-configure designated AI agent identity.

3. **Output & Clipboard Modes (`prompt_emitter.go`):**
   - Default Mode: Render high-contrast Lipgloss framed card in the terminal with token statistics and usage instructions.
   - Raw Mode (`--raw`): Stream unformatted markdown text suitable for shell redirection (`> prompt.md`) or piping.
   - Clipboard Mode (`--clipboard` / `-c`): Automatically copy generated execution prompt to the system clipboard (cross-platform via `atotto/clipboard` or native powershell/xclip adapters) and report confirmation card.

---

## 3. Out of Scope

- Auto-invoking AI agent APIs or launching subagents directly (this command provides the instructions to feed to agents).
- Full builds or test suite execution.

---

## 4. Architectural Invariants

- **Positive Booleans:** `isClipboard`, `isRaw`, `hasTaskName`, `hasCustomBatch`.
- **Leaf Error Sovereignty:** Clipboard write failures return structured `*apperror.AppError` without calling `os.Exit`.
- **Function Line Limit:** Keep template rendering and interpolation under 15 lines per function.
- **Relative Paths:** All prompt templates and links inside generated prompts must use 100% relative paths.

---

## 5. Verification Checklist

- [ ] `gitmap cg execute v6` emits complete V6 Tree Swarm execution instructions.
- [ ] `gitmap cg execute v7` emits complete V7 Time-Restricted execution instructions.
- [ ] `--clipboard` copies the full prompt string to the system clipboard and prints confirmation.
- [ ] `--raw` outputs bare markdown without ANSI escape sequences.
- [ ] `--task "demo-task"` correctly injects task slug into headers and checklists.
