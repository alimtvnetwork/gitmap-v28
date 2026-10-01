# Subtask 190.3: LLM Train Self-Loop and Machine URL Corpus

> **Parent Plan:** [Plan 190](../../completed/190-universal-command-help-modernization-llm-train-loop-and-polyglot-benchmarks.md)
> **Status:** Complete
> **Lead Architect:** MD ALIM UL KARIM

---

## 1. Objectives

1. Enhance `cli/cmd/llm/` with autonomous self-looping capabilities:
   - Add `--loop` / `--self-loop [N]` flag to run iterative self-training simulation or execution across the 5 AI agent operational phases.
   - Add `--url` to output direct public raw URL for easy download by any AI/LLM model.
   - Add `--json` to output machine-readable command dictionary and parameter specifications.
2. Register `gitmap llm-train` and `gitmap train` aliases in root CLI dispatch.
3. Update `cli/helptext/llm.md` and embedded markdown with the complete catalog of GitMap tools (AUM, Split-DB, WPR, Folder-Tree, Deploy, SSH Fleet, Pipeline-AI).
