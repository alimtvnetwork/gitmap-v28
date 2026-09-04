# Full Codebase Architecture and Skill Expansion

**Updated:** 2026-09-04
**Status:** Active
**Scope:** Repository-wide codebase audit, institutional memory absorption, and skill expansion

---

## 1. Executive Summary

A comprehensive, end-to-end reading of the Gitmap codebase was performed in accordance with the `read-memory-enhanced` protocol. This included 133 institutional memory files, 71 consolidated guidelines, 31 spec authoring documents, 41 CI/CD failure post-mortems, active and completed plans, open ambiguities, and the complete Go CLI backend and React frontend architectures.

To codify operational expertise for previously unrepresented domains, 2 new skills were authored in `.agents/skills/`:
1. `react-docs-frontend`: Dedicated React 18, Tailwind CSS, TypeScript enum, and Result envelope frontend architecture.
2. `cluster-ssh-delegation`: Multi-node cluster orchestration, SSH config templating, and distributed command delegation.

---

## 2. Ingested Repository Surface

- **Memory Files Read:** 133 files in `.lovable/memory/` across architecture, constraints, features, tech, and workflow domains.
- **Consolidated Guidelines Read:** 71 files in `spec/17-consolidated-guidelines/` (with `31-compiled-simple-coding-guidelines.md` as canonical source of truth).
- **Spec Authoring Documents Read:** 31 files in `spec/01-spec-authoring-guide/` defining 3-tier scoring, naming conventions, and required spec templates.
- **CI/CD Failure Post-Mortems Absorbed:** 41 failure analyses in `.lovable/cicd-issues/` spanning AST parity drift, race conditions, Go 1.25 compatibility, and linter regressions.
- **Plans & Tasks:** 0 active pending plans in `.lovable/plans/index.md` (all 28 historical plans completed).
- **Ambiguities:** 1 documented open ambiguity (`01-spec-gaps.md`), 0 resolved ambiguities.

---

## 3. Architectural Domains Covered

### 3.1 Go CLI Core (`gitmap/`)
- Over 60 subcommands in `gitmap/cmd/` with strict AST parity enforced by `TestTopLevelCmdRegistryMatchesAST` against `gitmap/constants/constants_cli.go`.
- Maximum 120 lines for command help markdown in `gitmap/helptext/` with 3–8 line simulations.
- Mandatory `cliexit.Reportf` / `cliexit.Fail` error handling; zero bare `os.Stderr` prints.
- SQLite connection pooling locked to `SetMaxOpenConns(1)` with binary path anchoring via `filepath.EvalSymlinks`.

### 3.2 React Frontend & Documentation (`src/`)
- React 18 SPA with Vite and Tailwind CSS.
- Semantic design tokens with amber gold `--primary` (`38 92% 50%` light, `41 96% 56%` dark).
- Strict TypeScript rules: Enum suffix `*Type`, custom hook object returns (no tuple returns), Result envelopes, and component size cap <= 100 lines.

### 3.3 Cluster & Multi-Host SSH (`gitmap/cluster/`, `gitmap/cmd/ssh.go`)
- Distributed cluster node management, SSH key discovery, and config generation.
- Transport preservation across reclone/clone operations.
- Worker pool delegation with explicit TLS dial timeouts.

---

## 4. Skills Ecosystem

The repository now maintains 15 specialized skills in `.agents/skills/`:
- `read-memory-enhanced`
- `coding-guidelines`
- `cg-boolean-and-naming`
- `cg-error-management`
- `db-sqlite-conventions`
- `gitmap-cli-command`
- `powershell-cross-platform-scripting`
- `spec-authoring`
- `release-and-versioning`
- `ci-cd-fix`
- `execute-pending-tasks`
- `execute-batched-loop-wor`
- `execute-parent-task`
- `react-docs-frontend` *(new)*
- `cluster-ssh-delegation` *(new)*
