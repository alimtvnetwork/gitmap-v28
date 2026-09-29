# Plan 190: Universal Command Help Modernization, LLM Train Self-Loop, and Polyglot Benchmarks

> **Plan Number:** 190  
> **Status:** Complete  
> **Lead Architect:** MD ALIM UL KARIM  
> **Sponsor:** RISEUP ASIA LLC  
> **Date:** 2026-09-29  
> **Specification:** [02-spec/21-app/180-universal-command-help-modernization-llm-train-loop-and-polyglot-benchmarks.md](../../02-spec/21-app/180-universal-command-help-modernization-llm-train-loop-and-polyglot-benchmarks.md)

---

## 1. Master Goal & Scope

Deliver complete repository-wide command help uniformity, modern display frame restructuring, an autonomous LLM training self-loop with direct specification URLs, and verified polyglot search performance benchmarks across File Search, Project Search, and Grid Search.

---

## 2. Subtask Breakdown

### Subtask 190.1: Universal Help Interception & Command Audit
- Audit all commands across `cli/cmd/` to ensure `gitmap <cmd> help` is intercepted before executing live logic or parsing positional parameters.
- Fix commands that fail on `help`: `commit`, `stash`, `branch`, `aum`, `templates`, `ui`, `revert`, `reinstall`, `user`.
- Ensure recursive subcommand help (e.g. `deploy config ssh help`, `pipeline errors help`) functions cleanly.

### Subtask 190.2: Modern Box Help Display Restructuring
- Restructure command help output to use the standardized double-line box frame (`╔═╗`, `║ ║`, `╚═╝`) with consistent sections (`Usage:`, `Commands / Subcommands:`, `Flags:`, `Examples:`, `[tip]`).
- Upgrade markdown-based command help files in `cli/helptext/` and `render` pipeline to format headers, flag tables, and code snippets cleanly in modern terminal style.

### Subtask 190.3: LLM Train Command Enhancement (`--loop`, `--url`, `--json`)
- Implement `gitmap llm-train` and `gitmap train` aliases directly in root dispatch.
- Add `--loop` / `--self-loop [N]` flag to run the autonomous 5-phase training cycle with simulated or live self-healing iterations.
- Add `--url` flag to output the public raw GitHub instruction spec URL.
- Update `cli/helptext/llm.md` and embedded spec with full coverage of all modern GitMap suites (`wpr`, `folder-tree`, `deploy`, `ssh`, `aum`, `pipeline-ai`).

### Subtask 190.4: Polyglot Search Benchmark Suite (File, Project, Grid)
- Author `03-ai-scripts/43-run-search-benchmarks.py` to benchmark:
  1. **File Search:** GitMap `find` vs PowerShell `Get-ChildItem -Recurse` vs Python `os.walk`/`rglob`.
  2. **Project Search:** GitMap `search`/`aum search` vs PowerShell `Select-String` vs Python grep.
  3. **Grid Search:** GitMap multi-filter AUM search vs PowerShell pipeline vs Python filtered scan.
- Measure latency (ms), throughput (files/sec), memory overhead (MB), and speedup multipliers.

### Subtask 190.5: Benchmark Documentation & Root README Integration
- Author root `benchmark.md` containing full methodology, environment specs, and detailed data matrices.
- Update root `readme.md` with three side-by-side benchmark comparison tables.

### Subtask 190.6: CI/AI Verification Prompt & Release
- Create `01-prompts/21-verify-universal-help-llm-and-benchmarks.md` with precise verification criteria.
- Execute cache & temp cleanup using `03-ai-scripts/42-clean-test-and-build-caches.py`.
- Bump minor version, cut release, and verify CI/CD is 100% green.

---

## 3. Progress Tracker

- [x] Spec 180 written (`02-spec/21-app/180-universal-command-help-modernization-llm-train-loop-and-polyglot-benchmarks.md`)
- [x] Plan 190 written
- [ ] Subtask 190.1: Universal Help Interception & Audit
- [ ] Subtask 190.2: Modern Box Help Restructuring
- [ ] Subtask 190.3: LLM Train Self-Loop & Machine URL Corpus
- [ ] Subtask 190.4: Polyglot Benchmark Suite Execution
- [ ] Subtask 190.5: `benchmark.md` & Root `readme.md` Update
- [ ] Subtask 190.6: CI Verification Prompt & Release Orchestration
