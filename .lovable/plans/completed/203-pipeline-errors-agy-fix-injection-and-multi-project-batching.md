# 203: Pipeline Errors AGY Fix Injection, Full Logs Embedding & Multi-Project Batching

## User Request (Verbatim)

```text
PS D:\work\Antigravity-Manager> gitmap pipeline errors agy fix                                                              
                                                                                                                            
  ✔ Prepared CI/CD pipeline fix prompt for Antigravity IDE!                                                                 
                                                                                                                            
    • Target Repo:    alimtvnetwork/Antigravity-Manager                                                                     
    • Pipeline Logs:  12.9 KB (failing error logs)                                                                          
    • Fix Prompt:     47.3 KB (D:\work\Antigravity-Manager\01-prompts\07-bug-fix\01-fix-with-rca.md)                        
    • Total Payload:  47.3 KB                                                                                               
    • Saved Payload:  .ai-memory/temp/active-agy-pipeline-fix-prompt.txt                                                    
    • Clipboard:      Copied to OS clipboard ✅                                                                              
                                                                                                                            
  Ready! Paste into Antigravity IDE chat window (Ctrl+V / Cmd+V) to start the fix loop.                                     
                                                                                                                            
  ✓ Follow-up Verification Prompt Queued: .ai-memory/temp/queued-agy-followup-prompt.txt ("Is it fixed?")                   
    • Queue Ledger:   .ai-memory/temp/agy-prompt-queue.json                                                                 
    • Antigravity CLI: Detected at C:\Users\Administrator\AppData\Local\agy\bin\agy.exe (run: agy -c)                       
                                                                                                                            
PS D:\work\Antigravity-Manager>                                                                                             

File Prompt text also not good and doen't contain the error logs and it didn't run with AGY at all for that project, full bug not done??

FILE:  .ai-memory/temp/active-agy-pipeline-fix-prompt.txt   

# Verification Check: Is It Fixed?

> **Target Repository:** `alimtvnetwork/Antigravity-Manager`
> **Workflow Run ID:** `#35334190423`
> **Commit SHA:** `f6ca20601e2e5abca0150c5090098690c0a6a240`

/goal Verify whether the CI/CD pipeline failure and errors have been completely resolved and all quality gates pass.

Please perform the following verification steps:
1. Confirm all root causes identified in the 4-part RCA have been properly remediated.
2. Verify that modified files strictly satisfy coding guidelines (<= 8-15 line functions, affirmative booleans, AppError envelopes).
3. Check live pipeline status or local test/linter verification ('gitmap pipeline status').
4. Confirm zero regressions and all quality checks pass before concluding.

Here, the issues are very critical. So when we say errors AGY fix, it does not put the prompt, um, into the, into the AGY. That means AG, AGY or Antigravity is not injected with the prompt and error log to fix it. That's a very big issue. I think you should have fixed it. Uh, second is that error logs are not into the text file, which is very lame. Uh, that is another bug. The-- another thing is like the file path that you have shown, that should be a full file path. In the terminal, whenever you show the file path, it should be full file path. Remember that. Um, so lot of the issues I think, uh, is not even clear or done properly. Okay. Um, so I-I'm not sure, um, what you have done here. The task is not very clear, I believe. So, uh, first of all, you have to understand that do you understand the error? Do you understand why it is happening? Do you know the fix? Uh, these are things I think you should clear out because it should run on Antigravity when I say that with the error log and everything. And also it can keep a backup of the prompt. That, that is fine in the text log. That's absolutely fine to showcase in the terminal. But also at the same time, I should be able to see that, that is on the queue. But it didn't happen, so that means you didn't write the code properly. That needs to be fixed. And there are several other examples with the error fix with AGY. I, I, I think I have given you lots of alternatives. I'm not sure, uh, if you know those and those helps are there or not. So please make sure that this is done very, very clearly, and we should be able to run this like on multiple repositories at once in the future. Uh, if we run it on, uh, let's say on, on top of the root repository or let's say we, we can have a bin projects that we are working on, we can run the pipeline errors fix so that it run parallelly on all those projects to get all the pipeline information, collect these, and push it to the, uh, Antigravity to fix, enqueue the task to fix those. Do you understand? So there should be a limit. By default, it should only do for three projects. If I run again, then it should do it for the next two projects, things like that. So that should explain in the terminal help as well, UI help as well. That is very, very priority. Okay, before doing anything, update the help, also the doc CMD and everything. Is it clear?
```

---

## Extracted Actionable Task List

- [x] **Task 1: Direct Antigravity (AGY) Injection & Execution**:
  - Automatically dispatch and inject the prompt payload directly into Antigravity (`agy` CLI or IDE session via `agy prompt "<payload>"`, `agy -c`, or active project dispatch), instead of merely copying to the clipboard.
  - Detect installed `agy` CLI binary or running Antigravity IDE and launch the fix task in the background/foreground without requiring manual paste.

- [x] **Task 2: Fix Full Error Logs & Fix Prompt Assembly**:
  - Ensure `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt` ALWAYS contains the full failing CI/CD pipeline error logs embedded within the 4-part RCA fix prompt payload (not just the follow-up verification check!).
  - Guarantee that `active-agy-pipeline-fix-prompt.txt` contains the primary fix RCA payload while `queued-agy-followup-prompt.txt` contains the secondary follow-up verification check.
  - Fix metric reporting so `Fix Prompt: XX KB`, `Pipeline Logs: YY KB`, and `Total Payload: ZZ KB` accurately reflect the individual component and concatenated lengths.

- [x] **Task 3: Full File Paths in Terminal Output**:
  - In all terminal outputs (`Saved Payload`, `Followup Queued`, `Queue Ledger`, `Destination File`), display strictly the **full absolute path** on disk (e.g. `D:\work\Antigravity-Manager\.ai-memory\temp\active-agy-pipeline-fix-prompt.txt`), never relative paths (`.ai-memory/temp/...`).

- [x] **Task 4: Prompt Queue Visibility & Auto-Enqueuing**:
  - Persist and showcase the prompt queue status in the terminal output with clear active vs queued items.
  - Enable inspection via `gitmap agy queue` or `gitmap pipeline errors agy queue` and show the queued verification prompt ready to fire once the fix completes.

- [x] **Task 5: Multi-Project Parallel Batching with Configurable Limits**:
  - Support running `gitmap pipeline errors agy fix --all` / `gitmap pipeline errors agy fix --projects <N>` across tracked repositories (or pinned projects).
  - Default batch limit: **3 projects** per run.
  - Track processed projects so subsequent executions pick up the next batch (e.g. next 2 or 3 projects) sequentially without duplicate dispatch.
  - Provide `--limit <N>` (default: 3) and `--batch-size <N>` flags.

- [x] **Task 6: Documentation-First, Terminal Help & Web UI Parity**:
  - Update `cli/helptext/pipeline.md` and `cli/helptext/agy-fix-pipeline.md` with multi-project batching, automated AGY injection, and full path output.
  - Update `src/data/commands.ts`, `src/pages/Pipeline.tsx`, and `src/pages/AGYPrompts.tsx` to document direct injection, queue states, and batching limits.

- [x] **Task 7: Strict Coding Guidelines & Atomic Git Push**:
  - Enforce Go files $\le 100$ lines and functions $\le 15$ lines.
  - Zero nested if statements, affirmative booleans only, AppError envelopes.
  - Fast lint verification (`check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-relative-paths.py`).
  - Strict ban on routine `go test` and `go build`.
  - Single atomic commit at final step pushed to `origin/main`.
