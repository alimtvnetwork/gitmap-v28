# Repository & .lovable Folder Architecture Map

This document establishes the official relationship and mapping between the **GitMap Repository Structure** and the **`.lovable` AI Agent System Architecture**.

---

## 1. High-Level Architecture Overview

GitMap is designed with a strict separation between four core domains:
1. **Source Code & CLI Engine** (`gitmap/`): Pure Go implementation adhering to strict size limits and error policies.
2. **Canonical Specifications** (`spec/`): The authoritative specifications (modules 01 through 21).
3. **User Documentation & Visual Assets** (`docs/commands/` & `docs/assets/`): Modular, lowercase command guides and animated terminal SVG demos.
4. **Agent Intelligence & Memory Core** (`.lovable/`): The AI operating environment governing plans, memory, constraints, and operational guidelines.

```text
gitmap-v28/
├── .lovable/                 # 🧠 AI Agent Intelligence, Plans & Memory Core
├── spec/                     # 📜 Canonical System Specifications & Standards
├── docs/                     # 📖 User Documentation & Visual Media
│   ├── assets/               #    Vector animated terminal SVGs (.svg)
│   └── commands/             #    Modular command & subcommand reference
├── gitmap/                   # ⚙️ Go CLI Source Engine
├── src/                      # 🌐 Web UI Documentation Dashboard (React + Vite)
├── bin/                      # 📦 Compiled Executables & Local Databases
└── readme.md                 # 🏠 Root Gateway & Project README
```

---

## 2. `.lovable/` Internal Taxonomy & Role Map

The `.lovable/` directory functions as the cognitive workspace for autonomous AI agents collaborating on the codebase:

| `.lovable/` Directory / File | Functional Purpose | Upstream Code / Doc Mapping |
|---|---|---|
| **`overview.md`** | High-level project identity, versioning, stack summary, and directory index | Roots `readme.md` & `gitmap/constants/constants.go` |
| **`folder-structure.md`** | Canonical architectural blueprint (this file) | Maps repo topology to agent execution workflows |
| **`plans/`** | Multi-phase task planning lifecycle (`pending/`, `active/`, `completed/`) | Drives agent task execution engines |
| **`memory/`** | Long-term cross-session knowledge storage |
| ├── `features/` | Verified feature capabilities and CLI registries | Maps to `gitmap/cmd/` and `docs/commands/` |
| ├── `constraints/` | Invariants, hard limits, and non-negotiables | Maps to `spec/02-coding-guidelines/` |
| ├── `specs/` | Spec synchronization and schema caches | Maps to `spec/21-app/` |
| └── `issues/` | Solved regressions and root-cause post-mortems | Maps to `spec/02-app-issues/` |
| **`coding-guidelines/`** | Coding standards (AST parity, 200-line limit, 15-line functions, positive booleans, zero swallowed errors) | Strictly enforced on all Go files in `gitmap/` |
| **`spec/`** | Agent-scoped mirrors and indices of specifications | Mirror of `spec/` |
| **`prompts/`** | Agent personas, subagent templates, and prompt chains | Referenced by agent execution skills |
| **`audits/`** | Audit logs verifying compliance with repository standards | Output of continuous CI and verification sweeps |
| **`strictly-avoid.md`** | Critical anti-patterns, dangerous commands, and forbidden techniques | Safeguards database, git state, and build tags |
| **`suggestions.md`** | Backlog of architectural proposals and optimizations | Backlog for future release planning |

---

## 3. Detailed Mapping: Code, Specs, Docs to `.lovable/`

### A. CLI Commands: Source $\leftrightarrow$ Docs $\leftrightarrow$ `.lovable`

- **Implementation:** `gitmap/cmd/<command>.go`
- **Terminal Help:** `gitmap/helptext/<command>.md`
- **User Documentation:** `docs/commands/<category>/<subcommand>.md`
- **Animated SVG Visuals:** `docs/assets/<category>.svg`
- **AI Agent Tracking:** `.lovable/memory/features/cli-commands.md`

### B. Architectural Decisions & Database

- **Master & Split DB Engine:** `gitmap/db/` & `gitmap/cmd/cmddb_*.go`
- **Formal Spec:** `spec/21-app/120-database-suite-and-start-fresh.md`
- **User Guide:** `docs/commands/db/readme.md` & `docs/commands/db/start-fresh.md`
- **Agent Guidelines:** `.lovable/coding-guidelines/` & `.lovable/memory/constraints/`

### C. Antigravity (AGY) Integration

- **Implementation:** `gitmap/cmd/agy_*.go`
- **Formal Spec:** `spec/21-app/122-antigravity-empty-conversations-pruner.md`
- **User Guide:** `docs/commands/agy/readme.md` & `docs/commands/agy/ls.md`
- **Visual Demo:** `docs/assets/agy.svg`
- **Agent Prompts:** `.lovable/prompts/`

---

## 4. Maintenance Invariants for `.lovable/`

1. **Lowercase File Names:** All documentation files and subdirectories created for features must use lowercase names (e.g. `readme.md`, `folder-structure.md`).
2. **Modular Granularity:** Do not combine unrelated topics into monolithic files. Separate into topic-specific markdown files and route them via `readme.md` or index files.
3. **Traceability:** Every plan created under `.lovable/plans/pending/` must cite relevant specs under `spec/` and user documentation under `docs/commands/`.
