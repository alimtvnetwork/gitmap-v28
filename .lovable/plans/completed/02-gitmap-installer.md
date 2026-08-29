---
status: pending
---

# Gitmap Installer Commands & Scripts Implementation

## Context

Goal: Implement the gitmap installer create/export/import commands and sqlite persistence (40 steps).
Reference spec: `.lovable/spec/commands/01-gitmap-installer.md`.

Execution model:
One step per run. Exactly one step is executed per run. Self-loop after Verify.
At most 2 spawned agents, max 3 threads.

## Tasks

- [ ] Task 001 — Define Installer Models
- [ ] Task 002 — Add SQLite Migrations for Installers
- [ ] Task 003 — Store CreateInstaller
- [ ] Task 004 — Store GetInstallerBySlug
- [ ] Task 005 — Store SaveVersion
- [ ] Task 006 — Store ListInstallers
- [ ] Task 007 — Store ResetInstallers
- [ ] Task 008 — Store DeleteInstaller
- [ ] Task 009 — Root Installer Command
- [ ] Task 010 — Installer Create Command
- [ ] Task 011 — Installer Update Command
- [ ] Task 012 — Installer Update Win Command
- [ ] Task 013 — Installer Install Win Command
- [ ] Task 014 — Installer Export Commands
- [ ] Task 015 — Installer Import Command
- [ ] Task 016 — Installer Reset Command
- [ ] Task 017 — Installer Revert Commands
- [ ] Task 018 — Installer List Command
- [ ] Task 019 — Path Normalization Utilities
- [ ] Task 020 — OS Targets Constants
- [ ] Task 021 — Manager Struct
- [ ] Task 022 — Manager Create Logic
- [ ] Task 023 — Manager Update Logic
- [ ] Task 024 — Semantic Versioning Helper
- [ ] Task 025 — Manager Undo Logic
- [ ] Task 026 — Manager Redo Logic
- [ ] Task 027 — Manager Exact Revert Logic
- [ ] Task 028 — Export Single ZIP
- [ ] Task 029 — Export All Installers
- [ ] Task 030 — Export Global State
- [ ] Task 031 — Import ZIP Dispatcher
- [ ] Task 032 — Import JSON Content
- [ ] Task 033 — Import Global State
- [ ] Task 034 — Import Conflict Resolver
- [ ] Task 035 — Installer Execution Dispatch
- [ ] Task 036 — Instruction Parser
- [ ] Task 037 — OS Command Runner
- [ ] Task 038 — Format Installer Table
- [ ] Task 039 — Detailed Help Printer
- [ ] Task 040 — Composer Example Printer
