# Specification 318-01: GitMap First-Class Alternative Commands, Skill Recipes & Ad-Hoc Script Eradication

> **Spec ID:** `SPEC-318-01`  
> **Parent Program:** `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/00-overview.md`  
> **Master Plan:** `.ai-memory/plans/pending/318-gitmap-command-and-coding-guideline-synchronization.md`  
> **Subtask Reference:** `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/subtask-01-alternative-commands-and-skills.md`  
> **Target Skill:** `.agents/skills/gitmap/skill.md`  
> **Subsystem:** CLI Tooling, Developer Skills & Process Orchestration  
> **Status:** Draft / Ready for Execution  

---

## 1. Executive Summary & Problem Formulation

In multi-agent collaborative workflows across Cursor, Muse, and Google Antigravity, autonomous coding agents frequently encounter friction points during low-level operating system operations. When facing tasks involving:
1. Disk capacity queries and folder size calculations,
2. Native Rust/C++ compilation failures (such as BoringSSL, Ring, and Tauri dependencies),
3. Process contention and port collisions,
4. Missing toolchain variables such as `LIBCLANG_PATH` during `cargo test` execution,

AI agents frequently degrade into executing raw, unbuffered ad-hoc scripts (e.g. ad-hoc Python inline scripts, raw Windows `taskkill`, hardcoded PowerShell environment mutations, and shell directory traversals).

This behavior represents the **Ad-Hoc Script Anti-Pattern**. Running ambient ad-hoc scripts introduces several severe operational hazards:
- **Portability Breakdown:** Windows commands (`taskkill /F /PID`, `tasklist`) fail immediately on Linux/macOS fleet nodes.
- **Audit & Telemetry Bypass:** Ad-hoc scripts bypass GitMap's SQLite Split-DB tracking (`gitmap.db`), losing diagnostic history and crash context.
- **Resource Exhaustion & Concurrency Collisions:** Ambient scripts do not participate in GitMap's file locking, worker concurrency throttling, or timeout governance.
- **Terminal Freezing:** Unbuffered recursive shell traversals freeze CLI communication sessions.

This specification details why AI agents (specifically Cursor) reverted to raw scripts, and defines the official, first-class GitMap alternative commands to be added to `.agents/skills/gitmap/skill.md`.

---

## 2. Root Cause Analysis: Why Cursor Executed Ad-Hoc Scripts

A comprehensive audit of recent agent interaction sessions reveals four distinct operational triggers that led to raw script execution:

### 2.1 The Folder Sizing & Disk Telemetry Incident
- **Observed Behavior:** The agent executed `gitmap py -c "import os; ..."` or ad-hoc Python loops with `os.walk` to calculate the aggregate byte size of build directories (`target/`, `node_modules/`, `.git/objects`).
- **Root Cause:**
  1. The agent understood that PowerShell `Get-ChildItem -Recurse | Measure-Object -Sum Length` was strictly banned in the skill matrix.
  2. The skill matrix prominently listed `gitmap py` as the sanctioned Python runner.
  3. However, `.agents/skills/gitmap/skill.md` did not adequately document `gitmap disk [--json]`, `gitmap os storage`, or `gitmap os dev-clean`.
  4. The LLM reasoned that writing an inline Python script was the most compliant path, failing to realize that GitMap already possesses sub-millisecond, compiled Go disk analysis commands.
- **Architectural Remedy:** Document `gitmap disk [--json]` and `gitmap os storage` in the non-negotiable command replacement table and provide concrete command recipes.

### 2.2 The BoringSSL, Ring & Tauri Dependency Incident
- **Observed Behavior:** When building or testing Rust/Tauri projects (such as `src-tauri` or native Go/Rust FFI modules), the compiler emitted errors regarding missing native toolchains (`cmake`, `nasm`, `perl`, `vcvarsall.bat`). Cursor attempted to download precompiled archives or run ad-hoc shell installations.
- **Root Cause:**
  1. Compiling cryptographic dependencies (e.g., BoringSSL in `boring-sys`, OpenSSL bindings, or Ring) on Windows requires MSVC environment flags and assembly toolchains.
  2. The agent was unaware that GitMap includes automated developer toolchain locators (`gitmap aum locate vcvars`, `gitmap aum locate llvm`) and package installers (`gitmap install <tool>`, `gitmap cargo --install build`).
- **Architectural Remedy:** Document first-class Rust/Cargo build flags (`gitmap cargo build [--install]`, `gitmap crg check`) and tool locators in the primary developer skill.

### 2.3 The Process Termination & Port Collision Incident
- **Observed Behavior:** When Vite dev servers (port 1420), local HTTP listeners (port 8080), or zombie `node.exe` / `gitmap.exe` instances locked database or source files, the agent ran raw `taskkill /F /PID <pid>` or `Stop-Process -Id <pid> -Force`.
- **Root Cause:**
  1. The agent encountered `EADDRINUSE` or file lock exceptions (`The process cannot access the file because it is being used by another process`).
  2. The agent was unaware of `gitmap ports [port]`, `gitmap os user kill <user>`, and `gitmap vscode repair --kill`.
  3. Raw `taskkill` calls emitted noisy non-zero exit codes when processes had already exited or lacked administrative elevation, confusing agent error handling.
- **Architectural Remedy:** Promote `gitmap ports [port]`, `gitmap os user kill`, and `gitmap vscode repair --kill` as the mandatory alternatives, with `gitmap pwsh "Stop-Process -Id <pid> -Force"` as the only permitted low-level wrapper when explicit PID termination is required.

### 2.4 The `LIBCLANG_PATH` Cargo Test Incident
- **Observed Behavior:** When running `cargo test` in crates with `bindgen` (e.g. SQLite C bindings or system API bridges), tests failed with:
  ```text
  thread 'main' panicked at 'Unable to find libclang: "couldn't find any valid shared libraries, check LIBCLANG_PATH"'
  ```
  Cursor responded by executing raw inline PowerShell string sequences:
  ```powershell
  $env:LIBCLANG_PATH = "C:\Program Files\LLVM\bin"; cargo test
  ```
- **Root Cause:**
  1. `cargo test` requires `libclang.dll` (Windows) or `libclang.so` (Linux) to generate Rust bindings from C headers at runtime.
  2. Hardcoding `C:\Program Files\LLVM\bin` broke portability across worker machines that installed LLVM via Chocolatey (`C:\ProgramData\chocolatey\lib\llvm`), Winget, or custom directories.
  3. Setting `$env:LIBCLANG_PATH` in a raw shell session did not persist across subagent process boundaries.
- **Architectural Remedy:** Provide a unified three-step pattern:
  1. Locate the tool dynamically via `gitmap aum locate llvm`.
  2. Persist or configure the environment variable via `gitmap env set LIBCLANG_PATH <path>`.
  3. Execute tests via native GitMap wrappers: `gitmap cargo test` or `gitmap pwsh "$env:LIBCLANG_PATH='...'; gitmap cargo test"`.

---

## 3. First-Class GitMap Alternative Commands Specification

The following first-class GitMap commands are specified as mandatory replacements. These commands must be added to `.agents/skills/gitmap/skill.md`:

### 3.1 Cargo & Rust Ecosystem Tooling
GitMap provides native wrappers around `cargo` that enrich the process environment with discovered toolchain paths (`CARGO_HOME`, `RUSTUP_HOME`, MSVC paths) and offer automated installation fallback.

| Command Syntax | Operational Semantics | Fallback Behavior |
| :--- | :--- | :--- |
| `gitmap cargo build [--install]` | Resolves `cargo` binary from PATH or fallback roots (`~/.cargo/bin`). Compiles the active Rust project. | If `cargo` is missing and `--install` is supplied, invokes `gitmap install cargo` automatically before build. |
| `gitmap cargo test [flags]` | Executes project test suites through the enriched environment with unified exit code reporting. | Traps missing toolchain and renders installation instructions. |
| `gitmap crg check` | Alias for `gitmap cargo check`. Fast syntax and type checking without producing binary artifacts. | Sub-second execution for CI validation gates. |
| `gitmap cargo status` | Inspects and displays Rust compiler version, Cargo package manager state, and active toolchain targets. | Emits actionable install advice if unconfigured. |

### 3.2 LLVM & libclang Path Management
To eliminate hardcoded system paths when compiling or testing crates dependent on `bindgen`:

| Command Syntax | Operational Semantics | Execution Context |
| :--- | :--- | :--- |
| `gitmap aum locate llvm` | Ultra-fast locator (<15ms) searching standard LLVM install locations, PATH, and registry entries. | Returns the verified bin directory containing `libclang.dll` / `libclang.so`. |
| `gitmap aum locate clang` | Locates Clang compiler executable across standard locations. | Emits exact executable path. |
| `gitmap env set LIBCLANG_PATH <path>` | Writes `LIBCLANG_PATH` to GitMap's managed environment registry or user environment. | Persists across process invocations. |
| `gitmap pwsh "$env:LIBCLANG_PATH='<path>'; gitmap cargo test"` | Scoped PowerShell execution with `-NoProfile` injecting `LIBCLANG_PATH` for a single run. | Cross-platform, deterministic environment propagation. |

### 3.3 Tauri Directory Inspection & Log Filtering
To inspect Tauri applications (`src-tauri/`) without slow shell file scans:

| Command Syntax | Operational Semantics | Prohibited Alternative |
| :--- | :--- | :--- |
| `gitmap lf src-tauri` | High-speed directory and file listing of the Tauri workspace. | ❌ `Get-ChildItem src-tauri`, `dir src-tauri` |
| `gitmap aum search "<pattern>" src-tauri` | Multi-core streaming live search across Rust and configuration files in `src-tauri`. | ❌ `Select-String`, `grep`, `ripgrep` |
| `gitmap pe -f tauri` | Filters CI/CD pipeline and runtime diagnostic logs specifically for Tauri build errors. | ❌ Manual log scraping or unbuffered stdout pipes |

### 3.4 Process & Port Termination
To resolve locked files, stale servers, and port conflicts without raw `taskkill`:

| Command Syntax | Operational Semantics | Safety Invariants |
| :--- | :--- | :--- |
| `gitmap ports [port]` | Lists active listening ports, associated PIDs, and process names. | Read-only inspection; displays port owners safely. |
| `gitmap os user kill <user>` | Terminates processes belonging to the specified OS user. | Clean teardown; supports `--force` flag. |
| `gitmap vscode repair --kill` | Safely terminates hung VS Code / Antigravity editor processes holding file locks. | Preserves workspace configuration and avoids data corruption. |
| `gitmap pwsh "Stop-Process -Id <pid> -Force"` | Deterministic process termination by explicit PID via GitMap's managed PowerShell runner. | Captures exceptions cleanly without crashing caller scripts. |

### 3.5 Disk, Storage & Developer Cache Telemetry
To inspect storage consumption and clean build artifacts without ad-hoc Python scripts:

| Command Syntax | Operational Semantics | Performance Characteristics |
| :--- | :--- | :--- |
| `gitmap disk [--json]` | Emits total, used, and free capacity for all visible physical and mounted drives. | Sub-millisecond (<1ms) compiled Go execution. |
| `gitmap os storage` | Inspects workspace drive status, working directory usage, and system temporary folder allocation. | Instant console display with colorized metrics. |
| `gitmap os dev-clean` | Safely purges developer cache directories (Cargo registry, Go build cache, npm store, temporary test files). | Dry-run enabled (`--dry-run`), selective category cleaning. |

---

## 4. Non-Negotiable Command Replacement Matrix Expansion

The following rows MUST be appended to the **Non-Negotiable Command Replacement Matrix** in `.agents/skills/gitmap/skill.md`:

```markdown
| Prohibited Legacy Tool / Pattern | GitMap Mandatory Command | Why GitMap is Required |
| :--- | :--- | :--- |
| ❌ Raw `cargo build`, `cargo test`, `cargo check` | ✅ `gitmap cargo build [--install]`, `gitmap cargo test`, `gitmap crg check` | Automatically enriches environment with discovered toolchain paths (`CARGO_HOME`, MSVC) and offers auto-install if missing. |
| ❌ Hardcoded `$env:LIBCLANG_PATH = "..."` or manual LLVM registry searching | ✅ `gitmap aum locate llvm` + `gitmap env set LIBCLANG_PATH <path>` or `gitmap pwsh "$env:LIBCLANG_PATH=...; gitmap cargo test"` | Dynamically resolves LLVM/Clang bin directories in <15ms across machines without breaking portability. |
| ❌ Raw `taskkill /F /PID <pid>`, `tasklist`, `kill -9` | ✅ `gitmap ports [port]`, `gitmap os user kill <user>`, `gitmap vscode repair --kill`, `gitmap pwsh "Stop-Process -Id <pid> -Force"` | Cross-platform, audit-logged process management that avoids spurious non-zero exit crashes and file locks. |
| ❌ Raw Python folder sizing (`os.walk`, `sys.getsizeof`) or `Get-ChildItem -Recurse \| Measure-Object` | ✅ `gitmap disk [--json]`, `gitmap os storage`, `gitmap os dev-clean` | Native compiled Go drive and storage metrics (<1ms). Prevents terminal freezing and avoids spawning unnecessary Python runtimes. |
| ❌ Raw directory inspection of Tauri apps (`dir src-tauri`, `Get-ChildItem src-tauri`) | ✅ `gitmap lf src-tauri`, `gitmap aum search "<pattern>" src-tauri`, `gitmap pe -f tauri` | High-speed directory inventory and pipeline error filtering optimized for Tauri/Rust desktop projects. |
```

---

## 5. Cheat Sheet Additions for `.agents/skills/gitmap/skill.md`

The following command clusters MUST be incorporated into the **Essential Command Cheat Sheet** in `.agents/skills/gitmap/skill.md`:

```markdown
### 6. Rust, Cargo & C Toolchain Management
- `gitmap cargo build [--install]` — Builds active Rust project with auto-toolchain resolution and installation fallback.
- `gitmap cargo test` — Executes project test suites in enriched environment with exit code normalization.
- `gitmap crg check` — Ultra-fast project check and validation (alias for `cargo check`).
- `gitmap cargo status` — Inspects Rust compiler and Cargo toolchain status.
- `gitmap aum locate llvm` (or `clang`) — Sub-millisecond locator for LLVM/Clang binaries (resolves `LIBCLANG_PATH`).
- `gitmap env set <key> <value>` — Persists environment variables across GitMap executions.
- `gitmap env get <key>` — Reads persisted environment variables.

### 7. Process, Port & Developer Environment Health
- `gitmap ports [port]` — Inspects listening network ports, PIDs, and process names.
- `gitmap os user kill <user> [--force]` — Terminates all processes owned by the given OS user.
- `gitmap vscode repair --kill` — Safely terminates hung VS Code / Antigravity processes locking workspace files.
- `gitmap pwsh "Stop-Process -Id <pid> -Force"` — Deterministic PID termination via managed PowerShell wrapper.

### 8. Hardware, Disk & Storage Telemetry
- `gitmap disk [--json]` — Instant physical drive metrics, capacity, and available space.
- `gitmap os storage` — Workspace drive utilization and system temporary folder diagnostics.
- `gitmap os dev-clean [--dry-run]` — Cleans developer caches (Cargo, Go, npm, pip) across the system.
```

---

## 6. Positive Boolean & Path Hygiene Contracts

In accordance with GitMap coding standards (`02-spec/02-coding-guidelines/`):
1. **Positive Boolean Requirement:** All boolean parameters, struct fields, and flags MUST use positive prefixes (`is`, `has`). Negative prefixes (`no`, `not`, `non`, `without`) are strictly prohibited in new schemas and documentation.
   - Allowed: `hasInstall`, `isJson`, `isForce`, `hasToolchain`, `isResolved`.
   - Prohibited: `noInstall`, `nonVerbose`, `notFound`.
2. **Strict Relative Git Paths:** All file and directory references in specifications, subtasks, and documentation MUST use forward-slash relative Git paths (`02-spec/...`, `.agents/...`, `cli/...`). Never use Windows backslash drive prefixes (`d:\...`) or file URIs (`file:///...`).
3. **Zero Test/Build Invariant (Rule R1):** When applying documentation and skill enhancements, subagents must never execute full workspace test suites (`go test ./...`, `cargo test`) or project builds. Pre-commit validation must be performed exclusively via lightweight static linters (`python 03-ai-scripts/50-fastgate.py`).

---

## 7. Downstream Traceability

This specification provides the architectural foundation for:
- `.ai-memory/plans/subtasks/318-gitmap-command-and-coding-guideline-synchronization/subtask-01-alternative-commands-and-skills.md`
- `.agents/skills/gitmap/skill.md` (active update)
- `cli/cmd/llm/llm_skill.go` (runtime training output synchronization)
- Fleet synchronization across all 43 connected repositories via `gitmap sync`.
