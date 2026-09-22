# Plan 79: AGY Prompts Injection Error Handling, Conversation Rename, Table Polish, Unique Sequence Persistence, and AUM Search Optimization

## User Request (Verbatim)
```text
https://prnt.sc/9gPwXSiL4y2m

https://prnt.sc/Bl-7E2INIrex

agy prompts inject still nto fixed why, not stacktrace why, didn't you uise proper apprError, where is the stack trace are you stupid???
make sure all the agy prompts has help with examples and read-all should n't be fix but one single prompt name, we should be able to see list of prompts and inject as we wish , do you understamd??

if you cannot find antigravtiy ide show the path what you found and also clearly check if agy cli helps to connect or not and also 

the sequence needs to be uniqye not per group fix it

and alsothe sequence needs to be saved into sqlite db so that remains same and can be used later on

if we do agy prompt help it should give all prompts example and prompts details how it works

clear???

also add command to rename conversation

gitmap agy conv-rename (cr) <id>/<seq>/startswith/<slug>/<path> "new name" # show status when done, clear???

https://prnt.sc/Bl-7E2INIrex

from these table remove the column , branch and status ( if not active don't display it)

also make sure that aum search are faster do benchmarking and correctness testing please.

FIx it with RCA and release minor bump
```

## 4-Part Root Cause Analysis (RCA)

### 1. Root Cause
- **AGY Prompt Injection Error Swallowing**: `tryDispatchExisting` and `tryInjectNewConversation` in `agy_agentapi_dispatch.go` dropped the `*apperror.AppError` returned by `AgentAPISendMessage` and `AgentAPINewConversation`. On failure, a generic string message was returned with zero stack trace, zero candidate path diagnostics, and zero `agy` CLI availability verification.
- **Forced Dual-Queue Prompting**: `EnqueueWithDualQueuePolicy` in `prompt_queue.go` forcibly queued `[1] Read all first prompt, [2] Current prompt` regardless of user input, preventing single prompt injection.
- **Per-Group Sequence Reset**: `printAgyTableRow` in `agy_ls_group.go` used the group-relative loop index `index + 1`, causing `SEQ` to reset to `001` for every folder group and lacking SQLite persistence.
- **Table Line Wrapping**: `BRANCH` and `STATUS` columns occupied excessive horizontal terminal width, and inactive/missing projects were rendered, causing visual line wrapping.
- **Missing Conversation Rename Command**: No CLI command existed to rename conversations by ID, SEQ, prefix, slug, or path.
- **AUM Search Inefficiencies**: `searchFileContent` split entire file buffers into lines with `bytes.Split` for every file (even files containing zero pattern matches), allocating memory unnecessarily.

### 2. Immediate Fix
- **Error Capture & Diagnostics**: Updated `tryInjectExistingConversation` and `tryInjectNewConversation` to return `*apperror.AppError`. Updated `makeOfflineFallbackResult` to attach `appErr`, print bounded stack traces, display searched IDE candidate paths, and check `agy` CLI connectivity.
- **Dynamic Single Prompts & Rich Help**: Replaced forced dual-queue policy with `EnqueueSinglePrompt`. Updated `PrintAgyPromptHelp` to load all templates from `cmdprompttemplate.LoadTemplates()` and render comprehensive examples.
- **Persistent Global Unique Sequences**: Created SQLite table `AgyProjectSequence` in `data/gitmap.db` via `agy_conv_seq.go`. Monotonically assigned unique sequence numbers across all groups.
- **Table Cleanup**: Removed `BRANCH` and `STATUS` columns from `agy_ls_render.go` and filtered out inactive projects in `agy_ls_group.go`.
- **Conversation Rename (`conv-rename` / `cr`)**: Implemented `ConvRenameCmd` and resolution engine in `agy_conv_rename.go` and `agy_conv_rename_resolve.go`.
- **AUM Search Optimization**: Added `hasFileMatch` pre-filtering before line splitting, replaced `bytes.Split` with zero-alloc `bytes.IndexByte` line scanning, and cached pre-lowercased pattern bytes.

### 3. Verification
- `go vet ./...` in `cli/`: PASS.
- AUM Search Benchmarks:
  - `BenchmarkRunSearch_Literal-20`: 3.42 ms/op (100 files).
  - `BenchmarkRunSearch_Regex-20`: 3.47 ms/op (100 files).
  - `TestRunSearch_Correctness`: PASS.
- `gitmap agy prompt help`: Verified rich help output with template catalog and examples.
- `gitmap agy ls`: Verified single-line table without `BRANCH` or `STATUS`, unique persistent `SEQ` numbers, and inactive project omission.
- `gitmap agy cr 15 "Gitmap Development Session"`: Successfully resolved and renamed conversation.
- `gitmap agy prompt -n is-done -t "..."`: Verified single prompt injection into active Antigravity session.

### 4. Prevention
- Enforce `AppError` return types on all external IPC and CLI dispatch functions.
- Centralize sequence numbering in SQLite database co-located with binary executable.
- Enforce strict width budgets on terminal table formatters to prevent multi-line row wrapping.

## Extracted Actionable Task List & Outcomes
- **Task-01:** AGY Prompts Injection Error Handling & Full Stack Trace with AppError — **COMPLETED**
- **Task-02:** AGY Prompts Help with Examples, Single Prompt Injection & Dynamic List Selection — **COMPLETED**
- **Task-03:** Antigravity IDE Path Diagnostics & AGY CLI Connection Verification — **COMPLETED**
- **Task-04:** Global Unique Sequence Number & SQLite DB Persistence for Conversations — **COMPLETED**
- **Task-05:** Conversation Rename Command (`gitmap agy conv-rename` / `cr`) — **COMPLETED**
- **Task-06:** AGY Conversation Table Formatting: Remove Branch & Conditional Status Column — **COMPLETED**
- **Task-07:** AUM Search Acceleration, Benchmarking & Correctness Testing — **COMPLETED**
- **Task-08:** Release Orchestration & Minor Version Bump (`v6.307.0` -> `v6.308.0`) — **READY**
