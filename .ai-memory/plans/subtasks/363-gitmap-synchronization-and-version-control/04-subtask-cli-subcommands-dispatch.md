# Subtask 04: CLI Subcommand Dispatch (gitmap cg sync, prompt sync, skill sync)

- **Parent Plan:** [363-gitmap-synchronization-and-version-control](../../pending/363-gitmap-synchronization-and-version-control.md)
- **Subtask ID:** `TASK-04`
- **Slug:** `gitmap-synchronization-and-version-control - Task 04`
- **Target Package:** `cli/root.go`, `cli/cmdcg/`, `cli/cmdprompt/`, `cli/cmdsync/`
- **Status:** `pending`
- **Assigned Worker:** CLI Dispatcher Subagent

---

## 1. Goal

Integrate first-class CLI routing and argument handling for `gitmap cg sync`, `gitmap prompt sync`, and `gitmap skill sync` into `cli/root.go`. Support targeted repository selection, dry-run simulation, JSON output, manifest-only generation, and universal alias dispatching (`gitmap sync`).

---

## 2. In Scope

1. **Dispatcher Routing & Alias Registration (`cli/root.go`):**
   - Register subcommand routes:
     - `gitmap cg sync` (and aliases `gitmap cg-sync`, `gitmap sync cg`)
     - `gitmap prompt sync` (and aliases `gitmap prompts sync`, `gitmap sync prompts`, `gitmap prompt-sync`)
     - `gitmap skill sync` (and aliases `gitmap skills sync`, `gitmap sync skills`, `gitmap skill-sync`)
     - `gitmap sync [category]` (universal umbrella command supporting `prompts`, `skills`, `cg`, `specs`, `all`)
   - Preserve existing CLI behavior and prevent command routing collisions.

2. **Flag Parsing & Options Struct (`sync_flags.go`):**
   - Implement `SyncCLIOptions`:
     - `isDryRun` (`--dry-run`): Scan and report actions without mutating disk files.
     - `isForce` (`--force`): Override human conflict protection and overwrite local files.
     - `isManifestOnly` (`--manifest-only`): Generate/update `version.json` sync manifests without copying file contents.
     - `isJSON` (`--json`): Emit structured telemetry payload conforming to standard GitMap JSON envelope.
     - `hasBackup` (`--backup`): Create non-destructive `.bak` sidecar files before modifications.
     - `targetRepos` (`--target <slug|dir>` or positional arguments).

3. **Target Repository Resolution (`target_resolver.go`):**
   - If positional arguments are passed, resolve paths relative to current workspace.
   - If no arguments are provided, use GitMap's repository cache to discover all active repositories in the current fleet/workspace.
   - Filter out ignored, hidden, or non-git folders (`.git`, `node_modules`, `vendor`).

4. **Command Execution Handlers:**
   - `cli/cmdcg/cg_sync_cmd.go`: Orchestrates full coding guidelines and prompt synchronization.
   - `cli/cmdprompt/prompt_sync_cmd.go`: Scoped synchronization for `01-prompts/`.
   - `cli/cmdprompt/skill_sync_cmd.go`: Scoped synchronization for `.agents/skills/`.

---

## 3. Out of Scope

- Modifying existing git commit or push commands.
- Visual Lipgloss card styling (implemented in Task 07).
- Full builds or test suite execution.

---

## 4. Architectural Invariants

- **Leaf Error Sovereignty:** All dispatch handlers must return `error` or `*apperror.AppError`. Do not call `cliexit.HandleError` inside command action handlers.
- **Positive Booleans:** `isDryRun`, `isForce`, `isManifestOnly`, `isJSON`, `hasBackup`.
- **Function Limit:** 8–15 lines max per dispatcher function. Flatten all argument validation.
- **Blank Lines:** Blank line after `}` and blank line before `return`.

---

## 5. Verification Checklist

- [ ] `gitmap cg sync --help` renders complete flags, aliases, and usage examples.
- [ ] `gitmap prompt sync` routes specifically to prompt assets.
- [ ] `gitmap skill sync` routes specifically to skill bundles.
- [ ] `--dry-run` executes full comparison without touching disk assets.
- [ ] `--json` emits clean, parseable JSON output matching the envelope spec.
