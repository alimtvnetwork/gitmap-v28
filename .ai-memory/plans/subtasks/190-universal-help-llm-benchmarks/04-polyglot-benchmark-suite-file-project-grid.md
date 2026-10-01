# Subtask 190.4: Polyglot Search Benchmark Suite (File, Project, Grid)

> **Parent Plan:** [Plan 190](../../completed/190-universal-command-help-modernization-llm-train-loop-and-polyglot-benchmarks.md)
> **Status:** Complete
> **Lead Architect:** MD ALIM UL KARIM

---

## 1. Objectives

1. Develop a Python benchmark runner `03-ai-scripts/43-run-search-benchmarks.py` that executes reproducible micro-benchmarks comparing:
   - **File Search:** GitMap `find` / `ff` vs PowerShell `Get-ChildItem -Recurse` vs Python `os.walk` / `Path.rglob`.
   - **Project Search (Content Search):** GitMap `search` / `aum search` vs PowerShell `Get-ChildItem | Select-String` vs Python regex grep.
   - **Grid Search (Multi-filter Matrix):** GitMap scoped multi-extension search vs PowerShell multi-filter pipeline vs Python multi-process filter.
2. Measure:
   - Execution latency (milliseconds / seconds)
   - Matches found
   - Memory overhead / RSS
   - Relative speedup ratios (GitMap vs PowerShell, GitMap vs Python)
3. Emit structured JSON results and formatted markdown comparison tables.
