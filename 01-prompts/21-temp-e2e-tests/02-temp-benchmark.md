# Local Search Benchmarking Workflow — AUM Search vs Scanner

> [!IMPORTANT]
> Prompt Version: 1.0.0
> Execution Constraint: Strictly local only. NEVER run in CI/CD or routine pipeline checks.
> Output Target: Root `benchmark.md` (uncommitted / gitignored).

/goal Perform comparative performance benchmarking between GitMap AUM search and standard repository scanning, documenting latency, memory, and indexing throughput.

---

### Benchmark Execution Protocol

1. **Pre-flight Checks**:
   - Ensure local workspace is clean and database is accessible.
   - Confirm `benchmark.md` is present in `.gitignore` to prevent accidental committing.

2. **Benchmark Scenarios**:
   - **Scenario A (AUM Cached / Index Search)**: Query workspace records using GitMap AUM search indexing.
   - **Scenario B (Standard Directory Scanner)**: Perform un-indexed filesystem crawl and parse across repository paths.
   - **Metrics Captured**:
     - Execution latency (milliseconds / seconds).
     - Repositories / files processed per second.
     - Memory footprint and exit code.

3. **Output Format (`benchmark.md`)**:
   - Header with timestamp and execution environment details (OS, CPU, repo size).
   - Comparison table with metrics across both search strategies.
   - Analysis of search efficiency and acceleration factors.

4. **Safety & Git Hygiene**:
   - Do NOT commit `benchmark.md` to version control.
   - Do NOT include benchmarking scripts in CI/CD test matrices.
