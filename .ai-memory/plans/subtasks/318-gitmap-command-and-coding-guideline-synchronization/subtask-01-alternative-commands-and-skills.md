# Subtask 01: GitMap First-Class Alternative Commands & Skills Enhancement

## Status: COMPLETED
- **Subtask Slug:** `subtask-01-alternative-commands-and-skills`
- **Parent Plan:** `.ai-memory/plans/pending/318-gitmap-command-and-coding-guideline-synchronization.md`
- **Architecture Spec:** `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/01-alternative-commands-and-skills-spec.md`
- **Owner:** Skills & CLI Tooling Subagent
- **Target Subsystem:** `.agents/skills/gitmap/`, `cli/cmd/llm/`
- **Exclusive Owned Files:**
  - `.agents/skills/gitmap/skill.md`
  - `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/01-alternative-commands-and-skills-spec.md`
  - `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/subtask-01-alternative-commands-and-skills.md`

---

## 1. Context & Objectives

During recent pair programming and autonomous task execution sessions, AI coding agents (specifically within Cursor and Muse environments) frequently resorted to running raw, unbuffered ad-hoc shell scripts:
1. **Raw Python for folder sizing:** `gitmap py -c "import os..."` or writing ad-hoc scripts using `os.walk` to compute folder sizes.
2. **Ad-hoc compilation scripts for BoringSSL / Ring / Tauri dependencies:** Attempting raw shell downloads or package installs when C/C++ build tools were missing.
3. **Raw Windows `taskkill`:** Running `taskkill /F /PID <pid>` or `Stop-Process` to resolve port conflicts or terminate hung processes.
4. **Ad-hoc `LIBCLANG_PATH` cargo test invocations:** Running inline `$env:LIBCLANG_PATH="..."; cargo test` when bindgen failed during Rust tests.

This subtask addresses the root cause of this behavior: **documentation gaps in agent skill instructions**. While GitMap already possesses high-speed, native Go implementations for these operations, they were not adequately presented in `.agents/skills/gitmap/skill.md`.

This subtask executes the following objectives:
- Author the definitive architectural specification in `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/01-alternative-commands-and-skills-spec.md`.
- Expand `.agents/skills/gitmap/skill.md` with first-class alternative commands across Cargo/Rust, LLVM/libclang path resolution, Tauri workflows, cross-platform process termination, and disk/storage telemetry.
- Enforce strict relative Git paths and positive boolean naming conventions across all deliverables.
- Prepare the updated skill for repository fleet synchronization across all connected workspaces.

---

## 2. Root Cause Analysis Summary

| Incident Pattern | Observed Ad-Hoc Script | Root Cause in Agent Reasoning | GitMap First-Class Alternative |
| :--- | :--- | :--- | :--- |
| **Folder Sizing** | `gitmap py -c "import os; ..."` or `Get-ChildItem -Recurse \| Measure-Object` | PowerShell `Get-ChildItem` was banned; agent knew `gitmap py` was allowed, but did not know `gitmap disk` or `gitmap os storage` existed. | `gitmap disk [--json]`, `gitmap os storage`, `gitmap os dev-clean` |
| **BoringSSL / Tauri Deps** | Raw curl/winget downloads for NASM, CMake, or C compilers | Build failures occurred during Rust FFI compilation; agent was unaware of GitMap's tool locator and automatic tool installer. | `gitmap cargo build [--install]`, `gitmap aum locate vcvars`, `gitmap install <tool>` |
| **Process Termination** | `taskkill /F /PID <pid>`, `taskkill /IM node.exe /F` | Dev server or port lock exceptions (Vite 1420, HTTP 8080); agent lacked knowledge of GitMap port and process managers. | `gitmap ports [port]`, `gitmap os user kill <user>`, `gitmap vscode repair --kill`, `gitmap pwsh "Stop-Process -Id <pid> -Force"` |
| **`LIBCLANG_PATH` Cargo Tests** | `$env:LIBCLANG_PATH = "C:\Program Files\LLVM\bin"; cargo test` | Rust bindgen panicked when `libclang.dll` was missing from PATH; agent hardcoded host-specific Windows paths. | `gitmap aum locate llvm`, `gitmap env set LIBCLANG_PATH <path>`, `gitmap pwsh "$env:LIBCLANG_PATH=...; gitmap cargo test"` |

---

## 3. Alternative Commands Specification

### 3.1 Cargo & Rust Toolchains
- **`gitmap cargo build [--install]`**: Compiles the current Rust project with automatic `CARGO_HOME` environment enrichment and auto-installation fallback via `gitmap install cargo`.
- **`gitmap cargo test [flags]`**: Runs Rust test suites with cross-platform exit code normalization.
- **`gitmap crg check`**: Fast, lightweight syntax and type checking without artifact emission.
- **`gitmap cargo status`**: Displays active Rust compiler, Cargo version, and target toolchain health.

### 3.2 LLVM & libclang Path Management
- **`gitmap aum locate llvm`** (or `clang`): Sub-millisecond tool finder resolving LLVM bin directories containing `libclang.dll` / `libclang.so`.
- **`gitmap env set LIBCLANG_PATH <path>`**: Sets and persists `LIBCLANG_PATH` in GitMap's environment registry across executions.
- **`gitmap pwsh "$env:LIBCLANG_PATH='<path>'; gitmap cargo test"`**: Scoped execution wrapper ensuring deterministic environment propagation without contaminating global shell state.

### 3.3 Tauri Directory Inspection & Log Filtering
- **`gitmap lf src-tauri`**: Instant sub-millisecond file listing and directory inventory of Tauri application roots (banning `Get-ChildItem`).
- **`gitmap aum search "<pattern>" src-tauri`**: Streaming multi-core live search scoped to Tauri Rust and frontend configurations.
- **`gitmap pe -f tauri`**: Pipeline and diagnostic log error extractor filtered specifically for Tauri build failures.

### 3.4 Process & Port Termination
- **`gitmap ports [port]`**: Displays listening ports, bound PIDs, and process names.
- **`gitmap os user kill <user> [--force]`**: Terminates all processes owned by the specified OS user.
- **`gitmap vscode repair --kill`**: Safely terminates orphaned VS Code / Antigravity editor processes holding file locks.
- **`gitmap pwsh "Stop-Process -Id <pid> -Force"`**: Sanctioned low-level PowerShell wrapper when explicit PID termination is mandatory.

### 3.5 Disk & Storage Telemetry
- **`gitmap disk [--json]`**: Sub-millisecond drive telemetry reporting total capacity, used space, and free space across all mounted drives.
- **`gitmap os storage`**: Workspace storage diagnostics and system temporary folder monitoring.
- **`gitmap os dev-clean [--dry-run]`**: Purges developer caches (Cargo, Go build cache, npm, pip) safely without data loss.

---

## 4. Execution Plan & Action Steps

### Step 1: Author Architecture Spec (Completed)
- Create `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/01-alternative-commands-and-skills-spec.md`.
- Detail the 4 ad-hoc script triggers, root causes, and official alternative commands.

### Step 2: Update `.agents/skills/gitmap/skill.md` (Completed)
- Appended the 5 prohibited tool rows to the **Non-Negotiable Command Replacement Matrix**:
  - Raw `cargo` commands -> `gitmap cargo build [--install]`, `gitmap cargo test`, `gitmap crg check`.
  - Hardcoded `$env:LIBCLANG_PATH` -> `gitmap aum locate llvm` + `gitmap env set LIBCLANG_PATH <path>`.
  - Raw `taskkill` / `kill -9` -> `gitmap ports [port]`, `gitmap os user kill <user>`, `gitmap vscode repair --kill`.
  - Raw Python folder sizing -> `gitmap disk [--json]`, `gitmap os storage`, `gitmap os dev-clean`.
  - Raw Tauri directory inspection -> `gitmap lf src-tauri`, `gitmap aum search ...`, `gitmap pe -f tauri`.
- Added dedicated cheat sheet sections in `## Essential Command Cheat Sheet`:
  - Section 18: Rust, Cargo & C Toolchain Management
  - Section 19: Tauri Directory Inspection & Log Diagnostics
  - Section 20: Process, Port & Developer Environment Health
  - Section 21: Hardware, Disk & Storage Telemetry
- Added real-world copy-pasteable examples in `## Practical Real-World Command Examples & Workflows`:
  - Section 15: Rust, Cargo & C Toolchain Management
  - Section 16: Tauri Workspace Inspection & Build Diagnostics
  - Section 17: Process Termination, Port Unlocking & IDE Health
  - Section 18: Hardware, Disk & Storage Telemetry
- Added Operational Guardrails 13–16 covering total bans on raw cargo, taskkill, Python folder sizing, and raw Tauri directory scans.

### Step 3: Validate Compliance & Hygiene (Completed)
- Verified relative Git paths across all modified documents with zero absolute paths or drive letters.
- Enforced positive booleans (`is*`, `has*`).
- Verified zero build and zero test runs (Rule R1).
- Executed `gitmap py 03-ai-scripts/50-fastgate.py` with 100% PASS (0.054s).
- Executed `gitmap py linter-scripts/check-relative-paths.py` with 100% PASS (0 absolute paths across 9449 files).

---

## 5. Acceptance Criteria

- **AC-1:** The specification file `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/01-alternative-commands-and-skills-spec.md` is fully authored and structured with executive summary, 4-incident root cause analysis, command specifications, and skill deltas.
- **AC-2:** The subtask file `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/subtask-01-alternative-commands-and-skills.md` is authored with clear metadata, ownership, task steps, and verification contracts.
- **AC-3:** All 5 first-class alternative command categories (Cargo/Rust, LLVM/libclang, Tauri workflows, Process termination, Disk/Storage telemetry) are covered with concrete syntax, flags, and operational semantics.
- **AC-4:** Strict relative Git paths are used across all markdown documents; zero absolute paths or drive letters (`d:\`, `C:\`, `file:///`) exist.
- **AC-5:** Positive booleans (`is`, `has`) are strictly maintained; zero negative booleans are introduced.
- **AC-6:** Rule R1 (Zero Build / Zero Test) is strictly observed during authoring.
- **AC-7:** FastGate static linter confirms 100% clean repository state.

---

## 6. Verification Evidence Contract

1. **Path Hygiene Linter:**
   ```bash
   python 03-ai-scripts/50-fastgate.py
   ```
   Must exit with code 0 in < 1.5s.
2. **File Existence & Non-Empty Check:**
   - `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/01-alternative-commands-and-skills-spec.md` (exists, > 5 KB).
   - `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/subtask-01-alternative-commands-and-skills.md` (exists, > 4 KB).
3. **Audit Protocol:**
   Deliverable files are reviewed for clarity, accuracy against existing GitMap CLI code, and ready for execution by downstream subagents.
