# 126 — Cargo Command & Rust Toolchain Runner Specification

## Overview

**Module Number:** 126
**Version:** 1.0.0
**Updated:** 2026-09-19
**Status:** Production-Ready
**AI Confidence:** Production-Ready
**Ambiguity Score:** None

---

## 1. Purpose & Architectural Vision

In polyglot engineering and desktop application development (e.g. Tauri, Rust backends, WebAssembly modules), developers and automated pipelines frequently invoke Cargo commands (`cargo fmt -- --check`, `cargo build`, `cargo check`, `cargo test`).

When Cargo is not installed on the system, standard shells fail with generic, unhelpful errors:
```
cargo: The term 'cargo' is not recognized as a name of a cmdlet, function, script file, or executable program.
```

The GitMap `cargo` command solves this by providing:
1. **Intelligent Binary Discovery & Fallback Resolution:** Inspects system `PATH`, `$CARGO_HOME/bin`, and `%USERPROFILE%\.cargo\bin`. If found in user fallback directories, it automatically enriches the process `PATH` so Cargo sub-tools (`cargo-fmt`, `clippy-driver`, `rustc`) execute without path issues.
2. **Proactive Remediation Suggestions:** If Cargo is missing, GitMap intercepts immediately and provides clear, formatted installation instructions pointing to `gitmap install cargo`.
3. **One-Shot Auto-Installation Flag:** `--install` / `-i` allows developers to install Cargo automatically on-demand and immediately execute the command.
4. **First-Class Installer Surface:** `cargo` is registered as a first-class tool under `gitmap install cargo` and `gitmap in cargo`.

---

## 2. CLI Command Syntax & Grammar

### 2.1 Direct Cargo Delegation

```bash
# Execute any standard cargo command through GitMap
gitmap cargo fmt -- --check
gitmap cargo build --release
gitmap cargo check
gitmap cargo test
gitmap cargo clippy

# Alias: crg
gitmap crg build
```

### 2.2 Toolchain Inspection

```bash
# Display installed versions and resolved paths for cargo, rustc, and rustup
gitmap cargo status
gitmap cargo st
gitmap cargo info
```

### 2.3 Installation & Auto-Installation

```bash
# Dedicated install command
gitmap install cargo
gitmap in cargo

# Run cargo with auto-install if missing
gitmap cargo --install fmt -- --check
gitmap cargo -i build
```

---

## 3. Remediation & Missing Tool Guidance

When Cargo is missing and `--install` is not supplied, GitMap returns a structured `*apperror.AppError` (`E7100:RUNTIME_MISSING`) and outputs formatted guidance:

```
  ✗ Cargo (Rust toolchain) is not installed or not found in PATH.

  Recommended (Install via GitMap):
    gitmap install cargo
    (or: gitmap in cargo)
    (or: gitmap install rust)

  Package Manager Fallback:
    Windows (winget) : winget install Rustlang.Rustup
    Windows (choco)  : choco install rust
    macOS (brew)     : brew install rust
    Linux (Ubuntu)   : sudo apt install -y cargo
    Official (curl)  : curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

  Tip: Automatically install and run:
    gitmap cargo --install <command...>
```

---

## 4. Package Manager Mapping

| Package Manager | Platform | Package ID / Command |
|---|---|---|
| **Winget** | Windows | `Rustlang.Rustup` (installs `cargo`, `rustc`, `rustup`) |
| **Chocolatey** | Windows | `rust` |
| **APT** | Ubuntu/Debian | `cargo` |
| **Homebrew** | macOS | `rust` |
| **Snap** | Linux | `rustup` |

---

## 5. Acceptance Criteria

### Scenario 1: Missing Cargo Triggers Formatted Suggestion
- **Given** an environment where Cargo is not in `PATH` or `~/.cargo/bin`
- **When** the user executes `gitmap cargo fmt -- --check`
- **Then** GitMap outputs the remediation block suggesting `gitmap install cargo` and exits with `E7100`.

### Scenario 2: Cargo Execution Passes Transparently
- **Given** Cargo is installed on the host machine
- **When** the user executes `gitmap cargo build --release`
- **Then** GitMap delegates directly to `cargo`, forward stdio live, and exits with Cargo's exit code.

### Scenario 3: Auto-Installation with `--install` Flag
- **Given** Cargo is not installed
- **When** the user executes `gitmap cargo --install fmt -- --check`
- **Then** GitMap invokes `gitmap install cargo --yes`, discovers the newly installed binary, and executes `cargo fmt -- --check`.

---

## 6. Cross-References

- Developer Tool Installer: [`./81-install.md`](./81-install.md)
- Polyglot Worker Orchestrator: [`./124-polyglot-worker-orchestrator-and-automation-runner.md`](./124-polyglot-worker-orchestrator-and-automation-runner.md)
- LLM Orchestration Playbook: [`./125-automation-llm-orchestration-guide.md`](./125-automation-llm-orchestration-guide.md)
