# RCA: Go Channel Deadlock in `gitmap aum search` and Scoped Learning Protocol

- **Date:** 2026-09-21
- **Affected Commands:** `gitmap aum search "train"`, `gitmap llm train --help`, `gitmap llm --help`
- **Scope:** CLI Automation (`cmdautomation/search.go`), LLM Chained Curriculum (`cmd/llm/`)

---

## 1. Symptom

1. **`gitmap aum search "train"` Hangs Indefinitely:**
   Executing `gitmap aum search "train"` hung indefinitely without emitting search output or completing, creating a hanging background process that consumed resources.
2. **`gitmap llm train --help` Crash:**
   Running `gitmap llm train --help` failed with:
   ```text
   gitmap llm: execute failed: [E9000:EXECUTION] parse train flags: flag: help requested (at=llm/llm_train.go:38)
   ```
3. **Bad Instruction & Anti-Pattern Loop:**
   AI agents and users attempting to "learn" or find training materials resorted to executing naive, broad keyword searches like `gitmap aum search "train"` across the entire repository because:
   - `gitmap llm --help` did not document the `train` or `chain` subcommands.
   - The chained curriculum listed `gitmap aum search --help` as Step 1 without enforcing directory or extension scoping.
   - The Antigravity skill omitted guidance on how models should properly acquire GitMap capabilities without scanning disk.

---

## 2. Root Cause

1. **Go Channel Deadlock in `dispatchSearchWorkers`:**
   In `cli/cmdautomation/search.go`, `results := make(chan []SearchMatch, opts.Workers)` allocated a channel with buffer capacity equal to worker count (e.g. 16).
   When searching a common word like "train" (49 hits across 3,566 files), each worker discovering matches attempted `results <- matches`. Once 16 slices were sent, the buffer filled, and worker goroutines blocked on channel send.
   Meanwhile, the main goroutine blocked on `wg.Wait()` before calling `gatherSearchResults(results)`.
   Because the consumer never read from `results` until workers finished, and workers could not finish while blocked on `results`, a circular channel deadlock occurred.
2. **Help Flag Error Wrapping in `parseTrainFlags`:**
   In `cli/cmd/llm/llm_train.go`, `fs.Parse(args)` returned `flag.ErrHelp` when `--help` was requested. Wrapping `flag.ErrHelp` into `apperror.WrapSimple` caused a fatal command failure (`E9000`) instead of a clean exit code 0.
3. **Missing Scoped Search Guardrails:**
   Documentation and curriculum lacked explicit guardrails forbidding unscoped keyword searches across thousands of repository files.

---

## 3. Resolution

1. **Non-Blocking Channel Streaming:**
   Refactored `dispatchSearchWorkers` in `cli/cmdautomation/search.go` to close `results` asynchronously in a background goroutine (`go waitAndCloseResults(&wg, results)`). The main goroutine immediately drains `results` via `gatherSearchResults(results)` in real time, preventing channel blockage regardless of hit volume.
2. **Clean Help Flag Handling:**
   Updated `parseTrainFlags` and `parseDefaultLlmFlags` to detect `flag.ErrHelp`, print dedicated usage instructions, and return clean `nil` (`exit 0`).
3. **Structured Learning Protocol & Scoped Search Directives:**
   - Updated `ChainedDiscoveryHeader` and `OperationalDirectivesText` in `cli/cmd/llm/llm_types.go` and `cli/cmd/llm/llm_train.go`:
     - Explicitly instructed agents: *To learn GitMap, run `gitmap llm train` or inspect `.agents/skills/gitmap/SKILL.md`. NEVER run unconstrained searches like `gitmap aum search "train"` to discover capabilities.*
     - Updated Step 1 in curriculum to show scoped search: `gitmap aum search "func Run" cli --ext .go`.
   - Updated `.agents/skills/gitmap/SKILL.md` template with scoped search guidelines and the LLM onboarding suite.
   - Updated Spec 127 (`02-spec/21-app/127-llm-train-and-chained-agent-curriculum.md`).
4. **Unit Test Verification:**
   - Created `cli/cmdautomation/search_test.go` (`TestRunSearch_ConcurrencyNoDeadlock` asserting 25 matches with 4 workers).
   - Added `TestRunTrainHelp` and `TestRunLlmHelp` in `cli/cmd/llm/llm_train_test.go`.

---

## 4. Prevention & Learnings

- **Always Drain Channels Concurrently:** In worker-pool architectures, never invoke `wg.Wait()` on the main goroutine before reading from the output channel unless the channel is buffered to hold the *entire maximum possible* items. Always drain on the main thread while closing from a background waiter goroutine.
- **`flag.ErrHelp` Is Never a Failure:** Command-line parsers must always intercept `flag.ErrHelp` and exit cleanly with code 0.
- **Agent Instruction Scoping:** AI discovery commands must always provide explicit scoping flags (`[dir]`, `--ext`) to prevent LLMs from emitting runaway whole-repository text scans.
