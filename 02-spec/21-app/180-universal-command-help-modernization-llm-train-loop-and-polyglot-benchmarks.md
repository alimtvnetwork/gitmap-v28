# Spec 180: Universal Command Help Modernization, LLM Train Self-Loop, and Polyglot Benchmarks

> **Version:** 1.0.0  
> **Status:** Draft  
> **Author:** MD ALIM UL KARIM  
> **Sponsor:** RISEUP ASIA LLC  
> **Date:** 2026-09-29  
> **Scope:** Repository-wide CLI Help Uniformity, Modern Banner Display, LLM Train Self-Loop Engine, File/Project/Grid Search Benchmarks, and Verification Prompt

---

## 1. Executive Summary

This specification formalizes three core architectural enhancements across GitMap:
1. **Universal `<command> help` Interception & Modern Box Styling:** Ensures every command and subcommand registered in GitMap natively supports `gitmap <cmd> help` (in addition to `--help` and `-h`) without failing or executing live operations, and standardizes help output into the modern double-line framed banner with aligned columns, flags, and concrete examples.
2. **LLM Train Self-Loop & Machine Corpus (`gitmap llm-train` / `gitmap train`):** Introduces a continuous self-looping curriculum (`--loop`, `--self-loop`, `--max-steps`) and canonical raw/JSON spec endpoints (`--url`, `--json`) enabling autonomous AI agents to internalize the full GitMap primitive catalog for zero-disk-thrashing operations.
3. **Polyglot Search Benchmark Suite (File Search, Project Search, Grid Search):** Implements a comprehensive benchmarking harness comparing GitMap against native PowerShell and Python across three distinct workloads (file discovery, content grep, and multi-parameter grid search), publishing verified latency, speedup multipliers, and memory overhead tables to both root `readme.md` and `benchmark.md`.

---

## 2. Universal `<command> help` Architecture

### 2.1 Problem Analysis
Prior to this specification:
- Commands like `gitmap commit help`, `gitmap stash help`, `gitmap branch help`, `gitmap aum help`, `gitmap templates help`, and `gitmap ui help` failed or triggered real actions because `help` was parsed as an argument (e.g. commit message, target branch, or UI page) instead of triggering help menus.
- Commands relying solely on raw markdown files dumped raw markdown headers (`# command`, `## Flags`) into the terminal without the visual polish, framed banners, and aligned formatting seen in modern commands (`gitmap clone help`, `gitmap deploy help`, `gitmap wpr help`, `gitmap folder-tree help`).

### 2.2 Global Help Interceptor (`cli/cmd/helpcheck.go` & `cli/cmd/root.go`)
- The root dispatcher intercepts any invocation where the first argument after the command is `"help"` (`os.Args[2] == "help"` or `-h` / `--help`).
- If an individual command does not register a custom rich renderer, the centralized help router renders the command's documentation using a standardized `ModernBox` frame:
  ```
  ╔═════════════════════════════════════════════════════════════╗
  ║   <Command Title> (gitmap <command>)                        ║
  ╚═════════════════════════════════════════════════════════════╝

  Usage:
    gitmap <command> [subcommand] [flags]

  Description:
    <Concise architectural description of the command>

  Commands / Subcommands:
    <Subcommand 1>            <Aligned description>
    <Subcommand 2>            <Aligned description>

  Flags:
    -f, --flag <type>         <Description>

  Examples:
    gitmap <command> <example 1>
    gitmap <command> <example 2>
  ```
- Subcommands (e.g. `gitmap deploy config ssh help`, `gitmap pipeline errors help`, `gitmap folder-tree export help`) inherit recursive help detection via `hasHelpFlag(args)`.

---

## 3. LLM Train Self-Loop & Machine Corpus

### 3.1 CLI Surface
```bash
gitmap llm train [flags]
gitmap llm-train [flags]
gitmap train [flags]
```

### 3.2 Feature Matrix
1. **`--loop`, `--self-loop [N]`**: Initiates an autonomous self-loop through the 5 phases of AI operation:
   - Phase 1: High-Speed Discovery (`find`, `ff`, `ffa`, `folder-tree`)
   - Phase 2: Targeted Refactoring (`replace`, `replace-regex`)
   - Phase 3: Verification & Local Linting (`go-format-check`, linters)
   - Phase 4: Atomic Semantic Commits (`cpf`, `cpb`, `cpr`)
   - Phase 5: CI/CD Live Telemetry & ETA Monitoring (`pipeline-ai`, `pe`, `eta`)
2. **`--url`**: Returns the canonical raw specification URL for direct HTTP ingestion by LLM models (`https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md`).
3. **`--json`**: Emits a structured machine-readable JSON schema defining every command, alias, purpose, and optimal AI usage pattern.
4. **`--skill-path`**: Generates and synchronizes `.agents/skills/gitmap/SKILL.md`.

---

## 4. Polyglot Benchmarks: File Search, Project Search & Grid Search

### 4.1 Benchmark Workloads
1. **File Search Workload**:
   - Query: Wildcard patterns (e.g. `*_test.go`, `*runner*.py`, `*.json`) across 3,000+ files.
   - Engines: GitMap Native Find (`gitmap find` / `ff`) vs PowerShell `Get-ChildItem -Recurse` vs Python `Path.rglob`.
2. **Project Search Workload (Content Grep)**:
   - Query: High-frequency tokens (`AppError`, `SSHConnection`, `Resolve-Version`) across all repository source files.
   - Engines: GitMap In-Memory / SQLite Hot-Cache (`gitmap search` / `aum search`) vs PowerShell `Select-String` pipeline vs Python multiprocessing grep (`12-fast-cached-grep.py`).
3. **Grid Search Workload (Multi-Filter & Matrix Search)**:
   - Query: Multi-extension + scoped subdirectory + token combinations (e.g. `cli/` directory, `.go` extension, pattern `"func Run"`).
   - Engines: GitMap Scoped AUM Grid Search vs PowerShell complex Where-Object pipeline vs Python filtered search.

### 4.2 Benchmark Artifacts
- **Harness**: `03-ai-scripts/43-run-search-benchmarks.py` and `gitmap aum benchmark all`.
- **Primary Report**: `benchmark.md` at repository root.
- **Root README Summary**: Three comparison tables in `readme.md`.

---

## 5. Verification Prompt Specification

A standalone prompt file will be authored at `01-prompts/21-verify-universal-help-llm-and-benchmarks.md` containing:
- Quantitative pass/fail criteria for `<command> help` across all primary command families.
- Verification checklist for `gitmap llm train --loop` and `--url`.
- Table formatting and numerical integrity checks for `benchmark.md` and `readme.md`.
- Cache cleanliness verification ensuring zero leftover artifacts after test runs.
