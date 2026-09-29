# Subtask 190.6: CI/AI Verification Prompt & Release Orchestration

> **Parent Plan:** [Plan 190](../../190-universal-command-help-modernization-llm-train-loop-and-polyglot-benchmarks.md)  
> **Status:** Complete  
> **Lead Architect:** MD ALIM UL KARIM  

---

## 1. Objectives

1. Author `01-prompts/21-verify-universal-help-llm-and-benchmarks.md` containing precision criteria that any AI or CI job can execute to verify:
   - All primary commands and subcommands respond to `<cmd> help` with modern framed formatting.
   - `gitmap llm-train --loop` runs through the 5 phases with 0 errors.
   - `gitmap llm --url` emits the valid GitHub raw URL.
   - `benchmark.md` and `readme.md` contain verified File Search, Project Search, and Grid Search benchmark tables.
2. Run test and cache cleanup with `03-ai-scripts/42-clean-test-and-build-caches.py`.
3. Orchestrate minor version release, tag, push, and confirm CI/CD is 100% green.
