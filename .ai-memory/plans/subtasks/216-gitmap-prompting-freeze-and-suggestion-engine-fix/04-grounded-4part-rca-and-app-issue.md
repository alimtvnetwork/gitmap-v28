# Subtask 216.04: Grounded 4-Part RCA & App Issue Registry Indexing

## 1. Context & Objective

As part of task 216 (`216-gitmap-prompting-freeze-and-suggestion-engine-fix`), App Issue 68 (`02-spec/22-app-issues/68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md`) was authored to document the forensic root-cause analysis, surgical code resolutions, and prevention guidelines for:
1. Win32 console input code page corruption (`SetConsoleCP(65001)`).
2. PowerShell wrapper pipe redirection (`| Out-String`) causing TTY stdin detachment.
3. Suggestion candidate pool isolation and Cobra completion runnable leaf pruning.

The objective of this subtask is to thoroughly review the authored RCA document against repository standards, verify all four sections (Symptom, Root Cause, Resolution, Prevention), confirm exact code line anchors, and register App Issue 68 in the master index table of `02-spec/22-app-issues/readme.md`.

---

## 2. Target Files

- `02-spec/22-app-issues/68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md`
- `02-spec/22-app-issues/readme.md`

---

## 3. Scope & Verification Checklist

1. **Four-Part RCA Structure Verification:**
   - **Section 1: Symptom:** Verify complete description of user symptoms:
     - Terminal keystroke freeze on confirmation prompts (`[y/N]`, secrets, clean).
     - Disambiguation pick list deadlock/EOF during `gitmap cd` / `gcd`.
     - Typo suggestion absence for core commands (`instlal`, `seach`, `fnd`, `commti`, `ap`).
     - Truncated shell tab autocompletion (9 commands instead of 512).
   - **Section 2: Root Cause:** Verify exhaustive analysis with exact line references:
     - `cli/cmd/console_windows.go:35`: `SetConsoleCP(65001)` Win32 conhost input buffer deadlock/EOF.
     - `cli/constants/constants_cd.go:201, 237`: `| Out-String` stdout buffering and stdin detachment.
     - `cli/cmd/rootsuggest.go:14-30`: Omission of core commands from `primaryTopCommands`.
     - `cli/cmd/root_cobra_completion.go:214-222`: Non-runnable leaf commands (`Run == nil`) pruned by Cobra.
   - **Section 3: Resolution:** Verify concrete diff specifications:
     - Removal of `SetConsoleCP(65001)` in `console_windows.go`.
     - Migration from `| Out-String` to foreground execution with `GITMAP_HANDOFF_FILE` in PowerShell wrappers.
     - Merging `completion.AllCommands()` with deduplication in `collectTopCommandCandidates()`.
     - Provision of stub `Run` functions and `ValidArgsFunction` in Cobra completion tree.
   - **Section 4: Prevention & Learnings:** Verify actionable preventive measures:
     - Static linter gate banning `SetConsoleCP` across all Go files.
     - Shell wrapper IPC standard mandating file-based handoff.
     - Automated tests in `rootsuggest_test.go` and `root_cobra_completion_test.go`.

2. **App Issue Registry Indexing:**
   - In `02-spec/22-app-issues/readme.md`, append row 68 to the master Contents table:
     ```markdown
     | 68 | [68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md](68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md) | GitMap Interactive Prompt Input Freeze and Suggestion Engine Collapse: RCA & Fix | Resolved |
     ```
   - Verify alignment, links, and table formatting.

3. **Acceptance Criteria Gate (AC-AI-001):**
   - Confirm that the issue document contains all four sections (Reproduction/Symptom, Cause, Fix/Resolution, Prevention) and references at least one commit/PR (`cli - fix prompt input freeze and suggestion engine`).

---

## 4. Implementation Steps

### Step 1: Ingest and Validate Issue 68 RCA File
- Inspect `02-spec/22-app-issues/68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md`.
- Ensure all line anchors (`cli/cmd/console_windows.go:35`, `cli/constants/constants_cd.go:201, 237`, `cli/cmd/rootsuggest.go:14-30`, `cli/cmd/root_cobra_completion.go:214-222`) are intact and accurately described.

### Step 2: Index Issue 68 in `02-spec/22-app-issues/readme.md`
- Edit `02-spec/22-app-issues/readme.md` under Section `## Contents`.
- Add row 68 immediately following row 67.
- Ensure the relative link `[68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md](68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md)` resolves cleanly.

### Step 3: Run Acceptance Verification
- Verify link integrity and Markdown rendering.
- Ensure zero broken links or orphaned documents.

---

## 5. Acceptance Criteria

- [ ] `02-spec/22-app-issues/68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md` satisfies all 4 RCA sections with exact code references.
- [ ] Issue 68 is properly indexed in `02-spec/22-app-issues/readme.md`.
- [ ] Relative Markdown link in `readme.md` resolves directly to the RCA file.
- [ ] AC-AI-001 audit criteria are fully satisfied.
