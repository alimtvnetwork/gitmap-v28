# Verification Prompt: Universal Command Help, LLM Training Self-Loop & Search Benchmarks

> **Task ID:** Plan 190 / Subtask 190.6  
> **Target Release:** `v6.397.0`  
> **Specification:** `02-spec/21-app/180-universal-command-help-modernization-llm-train-loop-and-polyglot-benchmarks.md`  

---

## 1. Objective

Execute autonomous end-to-end verification of:
1. **Universal Command Help Interception (`gitmap <command> help`)**:
   Verify that all registered CLI commands display the standardized framed double-line box banner (`╔═╗`, `║ ║`, `╚═╝`), `Usage:`, aligned two-column sections, flags, examples, and tips without throwing execution errors or triggering unwanted live logic.
2. **LLM Training & Self-Looping Suite (`gitmap train`)**:
   Verify autonomous 5-phase self-looping execution (`--loop`, `--self-loop [N]`), raw public spec URL output (`--url`), and machine-readable JSON dictionary (`--json`).
3. **Polyglot Search Benchmarks & Documentation**:
   Verify that `benchmark.md` and root `readme.md` publish the 3 side-by-side comparison tables across File Search, Project Search, and Grid Search workloads.
4. **Cleanliness & CI/CD Integrity**:
   Verify that no temporary build binaries or test outputs are left in the tree and that all GitHub Actions workflows are 100% green.

---

## 2. Automated Verification Checklist

Run each command in sequence and verify the expected output:

### Step 1: Universal Help Command Verification

```bash
# Verify previously failing commands now intercept help cleanly and exit 0:
gitmap commit help
gitmap stash help
gitmap branch help
gitmap aum help
gitmap templates help
gitmap ui help
gitmap revert help
gitmap user help

# Verify modernized search and locator help:
gitmap search help
gitmap find help
gitmap status help
gitmap diff help
```

*Criteria:*
- Every command must exit with code `0`.
- Output MUST feature the styled box banner:
  ```text
    ╔══════════════════════════════════════════════════╗
    ║   Command Title (gitmap <command>)               ║
    ╚══════════════════════════════════════════════════╝
  ```
- Output MUST include `Usage:`, aligned command descriptions, and `[tip] Tip:`.

---

### Step 2: LLM Training & Autonomous Self-Loop Verification

```bash
# 1. Verify help display:
gitmap train help

# 2. Verify public instruction URL output:
gitmap train --url
# Expected: https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md

# 3. Verify machine-readable JSON dictionary:
gitmap train --json

# 4. Verify 5-phase self-looping execution:
gitmap train --loop
gitmap train --self-loop 2
```

*Criteria:*
- `gitmap train --url` outputs exact raw GitHub URL.
- `gitmap train --json` returns valid JSON with `phases` and `efficiency_factors`.
- `gitmap train --loop` prints all 5 lifecycle phases (`Discovery`, `Refactoring`, `Verification`, `Semantic Commit`, `CI Telemetry`) and reports 100% green.

---

### Step 3: Polyglot Search Benchmark Verification

```bash
# Run benchmark runner script:
python 03-ai-scripts/43-run-search-benchmarks.py --runs 2

# Verify output file generated:
cat tmp/benchmarks/search_benchmark_results.json
```

*Criteria:*
- Script runs without errors and produces all 3 tables:
  - Table 1: Wildcard File Search Benchmarks (`*config*.json`)
  - Table 2: Project Content Search Benchmarks (`SSHConnection`)
  - Table 3: Grid Search Benchmarks (Multi-Filter: `func Run` in `cli/*.go`)
- `benchmark.md` and `readme.md` both contain the verified comparison tables.

---

### Step 4: Cleanup & CI/CD Release Ceremony

```bash
# Run mandatory cache and test cleanup script:
python 03-ai-scripts/42-clean-test-and-build-caches.py

# Verify working directory is pristine:
git status

# Execute minor version release ceremony:
python 03-ai-scripts/18-release-orchestrator.py --type minor
```

*Criteria:*
- Working directory clean before release.
- Release branch and tag pushed to remote.
- GitHub Actions CI/CD workflows complete 100% green.
