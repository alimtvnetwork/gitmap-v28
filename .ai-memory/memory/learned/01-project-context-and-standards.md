# 01 — Project Context and Architecture Standards

- **Subsystem:** Core Architecture & Project Governance
- **Status:** Authoritative Reference

## 1. Project Identity & Governance
GitMap is an autonomous developer companion and high-performance command-line utility built in Go, designed for multi-repository discovery, fast git operations, multi-node SSH fleet management, CI/CD telemetry diagnosis, and AI coding agent orchestration.

## 2. Multi-Agent Task Orchestration
- **Agent Roles:** Parent orchestrator dispatches disjoint work packages across subagents (Worker 01 for specifications/linters, Worker 02 for memory/completed plans consolidation).
- **Execution Budget:** Continuous self-loops bounded to designated steps (N) with non-destructive progress tracking.
- **Git State Hygiene:** Subagents are subject to a strict ban on git commands; only the parent lead agent stages, commits, and tags releases.

## 3. Core Directory Architecture
- `01-prompts/`: Canonical prompts and workflow directives.
- `02-spec/`: Architectural, component, and coding guideline specifications.
- `03-ai-scripts/`: Toolchain helper scripts and cross-platform automation.
- `.ai-memory/`: Persistent architectural memory, plans, and issues.
- `cli/`: Go source code implementing Cobra commands, handlers, and renderers.
