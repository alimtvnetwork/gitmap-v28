# Project Architecture, Memory, and Skills Synthesis

**Updated:** 2026-09-02
**Status:** Active
**Scope:** Full-repository memory load, codebase architecture, and skill extraction

---

## 1. Executive Summary

A comprehensive, deep-dive read of all institutional memory files (132 files), consolidated guidelines (36 files), spec authoring guides (17 files), CI/CD issues (41 files), ambiguities, and core codebase structures was performed. From this institutional memory, 5 dedicated skills were synthesized and codified in `.agents/skills/` to provide autonomous agents with precision execution capabilities for Gitmap workflows.

---

## 2. Ingested Repository Surface

- **Memory Files:** 132 files absorbed from `.lovable/memory/` across architecture, constraints, features, tech, and workflow domains.
- **Consolidated Guidelines:** 36 policy files absorbed from `spec/17-consolidated-guidelines/` (mirroring authoritative standards from `31-compiled-simple-coding-guidelines.md` and `.lovable/coding-guidelines/`).
- **Spec Authoring:** 17 files absorbed from `spec/01-spec-authoring-guide/` establishing the standard 3-tier scoring, file naming, required files (`00-overview.md`, `99-consistency-report.md`), and cross-reference rules.
- **CI/CD Issues & RCA:** 41 documented failure modes in `.lovable/cicd-issues/` absorbed to prevent repeating past failures (including race conditions, AST parity breaks, version drift, and linter regressions).
- **Plans & Ambiguities:** Verified 0 pending plans in `.lovable/plans/index.md`, 1 documented open ambiguity (`01-spec-gaps.md`), and 0 resolved ambiguities.

---

## 3. Core Architectural Pillars

### 3.1 Go CLI Toolchain (`gitmap/`)

- Go 1.24.13 with strict linting (`golangci-lint` v1.64.8, `govulncheck` v1.1.4).
- Over 60 subcommands organized in `gitmap/cmd/` with AST parity enforced via `TestTopLevelCmdRegistryMatchesAST` against `gitmap/constants/constants_cli.go`.
- Maximum 120 lines for command help markdown files (`gitmap/helptext/`) with 3–8 line realistic simulations and mandatory fenced code blocks.
- Centralized error and exit handling via `apperror` and `cliexit.Reportf` / `cliexit.Fail`. Bare `fmt.Fprintln(os.Stderr, err)` is strictly prohibited.

### 3.2 Database Conventions (`store/`, SQLite)

- CGo-free SQLite engine via `modernc.org/sqlite`.
- Singular PascalCase table names, `{TableName}Id INTEGER PRIMARY KEY AUTOINCREMENT`.
- Strict connection limit: `db.SetMaxOpenConns(1)` preventing database locks and concurrency races.
- Database location anchored to binary execution location via `filepath.EvalSymlinks(os.Executable())`.
- Mandatory context columns: `Description TEXT NULL` for entity tables, `Notes TEXT NULL` / `Comments TEXT NULL` for transaction tables.

### 3.3 Versioning & Release (SSoT)

- Single source of truth is `version.json` at repository root.
- Go binaries receive version via `-ldflags` string injection.
- Release pipeline workflows (`.github/workflows/release.yml`) are strictly untouchable; missing release assets are always investigated as upstream test/build gate failures.
- Non-gitmap repositories are gated against receiving gitmap installation snippets via `ShouldPrintInstallHint`.

---

## 4. Skills Extracted and Created

The following 5 new skills were extracted from the codebase and added to `.agents/skills/`:

1. **`gitmap-cli-command`** (`.agents/skills/gitmap-cli-command/skill.md`):
   Authoring, registering, and testing Gitmap CLI subcommands with AST parity (`// gitmap:cmd top-level`), help text formatting, and `go generate` drift prevention.

2. **`db-sqlite-conventions`** (`.agents/skills/db-sqlite-conventions/skill.md`):
   Designing and querying SQLite schema adhering to PascalCase naming, `SetMaxOpenConns(1)`, binary anchoring, and DB rules 10–12.

3. **`spec-authoring`** (`.agents/skills/spec-authoring/skill.md`):
   Authoring specification modules adhering to `spec/01-spec-authoring-guide/` with 3-tier scoring (AI Confidence, Ambiguity, Health Score) and cross-reference validation.

4. **`release-and-versioning`** (`.agents/skills/release-and-versioning/skill.md`):
   Managing version bumps across `version.json`, `changelog.md`, `package.json`, and preserving release pipeline boundaries.

5. **`powershell-cross-platform-scripting`** (`.agents/skills/powershell-cross-platform-scripting/skill.md`):
   Authoring cross-platform PowerShell and Bash scripts with UTF-8 BOM encoding, exit code contracts, and direct Python invocation.
