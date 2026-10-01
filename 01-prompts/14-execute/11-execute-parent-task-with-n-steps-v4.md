# [V4] Parent Task N-Step Loop: Antigravity-Native Orchestrator (must follow)

```text
N = 300         Step ceiling for the whole run (edit before running)
A = 2           Worker subagents per wave (invoke_subagent)
H = 2           Subtasks per worker per wave
COMMIT = push   push | local | none

PHASE_1_BUDGET = N / 2   Steps 1 .. 150: preflight, capture, spec, plan, subtasks
PHASE_2_BUDGET = N / 2   Steps 151 .. 300: worker waves, acceptance, consolidation, commit
```

> [!IMPORTANT]
> Prompt Version: 4.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI), checked against the Antigravity docs on 2026-09-30.
> Invoke: paste this prompt below your task, or run the skill `/execute-parent-task-with-n-steps-v4 <task>`.
> Parameters are read-only after Step 0. N is a ceiling, not a quota: finishing early is success, and padding steps is a failure.
> One step is one round of tool calls by the lead agent. The lead records the step count in the ledger after every round.

[/goal](slashCommand;goal) Execute the parent task end to end: capture it verbatim, plan it in the repo, run it through `A` worker subagents in disjoint file boxes, prove every claim with evidence, and finish with one commit that holds only this task's files.

[/learn](slashCommand;learn) Top-Instruction Priority Mandate: the instructions above this prompt and in the invoking message outrank everything below. Each rule is stated once (R1 to R15) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

## 1. Precedence (highest first)

1. The user's instructions above this prompt or in the invoking message, including edited parameters.
2. Platform limits: the tools in your tool list, the Artifact Review Policy, permission prompts, and hooks. Never claim to override them.
3. Repo rules: `.ai-memory/strictly-avoid.md` and `.ai-memory/coding-guidelines.md` (canonical source: `02-spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md`).
4. This prompt.

If two sources at the same level conflict, follow the stricter one and record the conflict in the ledger.

## 2. Rules (cite by ID)

- **R1 No builds or test suites.** Never run `go build`, `npm run build`, `go test`, `pytest`, `npm test`, or `03-ai-scripts/06-cicd-local-runner.py`. CI verifies them. Only an explicit user command lifts this.
- **R2 Targeted checks only.** Run the checks in section 9 on the files or folders this task changed. A check that scans 0 files is a FAIL, not a pass.
- **R3 Evidence or it did not happen.** Every DONE, PASS, or "verified" names a path, a diff, or a command with its exit code.
- **R4 Never invent.** Before first use, confirm that a tool is in your tool list and that a script, flag, or path exists (`--help`, a directory listing). If it is missing, use the fallback and log it.
- **R5 Subagents are mandatory when `invoke_subagent` exists.** Use `research` subagents for read-only discovery and `self` subagents for edits. If the tool is absent, run the same waves yourself, one subtask at a time, and log `SOLO_FALLBACK: invoke_subagent absent`. That is the only allowed solo path.
- **R6 One owner per file.** Within a wave, each file has exactly one owner. Shared indexes (`.ai-memory/plans/readme.md`, `.ai-memory/prompts.md`, `.ai-memory/what-to-read.md`, `02-spec/21-app/readme.md`, and any folder `readme.md`) belong to the lead.
- **R7 Git safety.** Workers never run `git add`, `git commit`, or `git push`. Nobody runs `git reset --hard`, `git checkout --`, `git clean`, `git stash`, or a force push.
- **R8 Stage by explicit path.** Use `git add -- <paths from the ledger stage list>`. Never stage everything at once with the `-A` flag or `.` as the path.
- **R9 One atomic commit** at the end of the run, never one per file. Push according to `COMMIT`.
- **R10 No releases.** No version bumps, `version.json` edits, changelog dates, or release commands unless the user explicitly asks for a release.
- **R11 Paths.** No absolute paths or `file:///` URIs in repo files. Paths are relative to the repo root, and new file names are lowercase.
- **R12 No polling.** After dispatching subagents, yield and wait for their result messages. For remote CI, use `gitmap pipeline-ai status --json` and wait on its `etaSeconds`; never loop `gh run view`.
- **R13 Retry cap.** The same failing call twice means change the approach. A subtask that fails two remediation rounds is marked FAILED with an RCA (section 11), and the run continues.
- **R14 Ambiguity.** Non-blocking: take the most conservative option consistent with the spec, record it under Assumptions, and continue. Blocking (destructive, irreversible, needs credentials, or contradicts a spec): ask once with `ask_question` if it exists, log it in `.ai-memory/ambiguous-questions/01-new-ambiguity/`, and keep working on every unblocked Task-ID.
- **R15 No generated artifacts.** Never commit caches, test reports, binaries, or newly generated code (Hard Rule 1). Run a code generator only when the subtask's spec requires it, and commit its output only if `git ls-files` shows the repo already tracks those files, so CI sees no drift.

## 3. Step 0: Preflight (counts toward Phase 1)

1. **Tool map.** Record which of these are in your tool list, and use only those: `invoke_subagent`, `send_message`, `manage_subagents`, `ask_question`, artifact writes (`write_to_file` with artifact metadata), `run_command`. Older builds offer `task_boundary` and `notify_user` instead; if `task_boundary` exists, set PLANNING for Phase 1, EXECUTION for Phase 2, and VERIFICATION for acceptance and the final checklist.
2. **Script map.** Record whether `gitmap --version` and `python --version` succeed, and whether each script you plan to call exists.
3. **Tree state.** Record the branch and `git status --porcelain`. Every file that is already modified is not yours: never edit or stage it unless a Task-ID names it. Confirm the root `readme.md` is lowercase.
4. **Context.** Read the "Before any task" list in `.ai-memory/what-to-read.md`, then `.ai-memory/strictly-avoid.md` and `.ai-memory/coding-guidelines.md`.
5. **Resume.** If a ledger at `.ai-memory/temp-agents/NN-<slug>/ledger.md` already records this same request, continue from its last recorded step instead of starting over.
6. **Review policy.** With "Request Review", the platform pauses after plan artifacts; after approval, continue from the ledger. Unattended runs need "Always Proceed", which only the user can set.

Phase 1A creates this ledger (gitignored) as the run's first file change; update it after every round:

```markdown
# Ledger: NN-<slug>
Request: <first line of the verbatim request>
Step: 12 / 300 (Phase 1: 12 / 150, Phase 2: 0 / 150)
Branch: <branch> | Tree at start: clean, or dirty with <paths>
Tools: invoke_subagent=yes send_message=yes ask_question=yes artifacts=yes gitmap=no
| Task-ID | Subtask | Owner | Owned files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | 01-<name> | Worker 01 | <paths> | DONE | <check> exit 0, 3 files |
Assumptions: <list or none>
Conflicts: <list or none>
Stage list: <every path this task created or changed>
```

## 4. Phase 1A: Capture and Confirm

1. Capture the parent task verbatim, including every instruction above this prompt. It goes, unedited, under `## User Request (Verbatim)` in both the spec and the plan.
2. Save each screenshot or base64 image to `assets/screenshots/<slug>-<NN>.png` and reference it by relative link. Never put base64 data or temporary URLs in repo files.
3. Split the request into ordered Task-IDs (`Task-01`, `Task-02`, ...). Nothing the user asked for may be dropped or merged away.
4. Pick NN, the smallest two-digit number not used in `.ai-memory/plans/pending/`, `.ai-memory/plans/completed/`, or `02-spec/21-app/`, and a lowercase, hyphenated slug.
5. Print the breakdown below, and in the same response write the ledger, your first file change. Do not ask "Should I proceed?".

```markdown
### Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS]`
   - **Understood:** `[YES]` — [one sentence proving you understood the intent, scope, and constraints]
   - **Actionable Scope:** [precise deliverable]
   - **Target Files / Area:** `[relative/path/or/module]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED]`
   - **Understood:** `[YES]` — [one sentence]
   - **Actionable Scope:** [precise deliverable]
   - **Target Files / Area:** `[relative/path/or/module]`

Proceeding to Phase 1B (tool call below).
```

Each task is its own numbered item, separated by a blank line. If a task is not fully understood, write `[PARTIAL]` and the open question (R14) instead of `[YES]`.

## 5. Phase 1B: Spec, Plan, and Subtasks (Phase 1 budget)

1. **Blueprint first.** The lead alone writes the spec overview and the plan outline. Workers never write competing plans.
2. **Discovery.** Dispatch up to `A` `research` subagents on disjoint areas, then yield as in section 6. Together they cover every folder of `02-spec/` and `.ai-memory/` that the Task-IDs touch, and return paths, symbols, and call sites with line numbers.
3. **Spec.** Write `02-spec/21-app/NN-<slug>.md` (150 lines or fewer), or the folder `02-spec/21-app/NN-<slug>/` with `01-overview.md`, `02-data-contracts.md`, `03-visual-and-ux.md`, and `04-verification-gates.md`. Register it in `02-spec/21-app/readme.md`.
4. **Plan.** Write `.ai-memory/plans/pending/NN-<slug>.md` with the spec link, the verbatim request, the blast radius, and the Task-ID to subtask map. Register it in `.ai-memory/plans/readme.md`.
5. **Subtasks.** Write one file per subtask in `.ai-memory/plans/subtasks/NN-<slug>/`, holding only what is unique to it:

```markdown
# Subtask [01]: [Descriptive Subtask Name]
Traceability ID: Task-01
Spec Reference: [02-spec/21-app/NN-<slug>.md](../../../02-spec/21-app/NN-<slug>.md)
Owned Files: [relative paths; mark files to create as NEW]
Action: [exact changes: functions, types, logic]
Acceptance Criteria: [2 to 4 testable conditions]
Targeted Verification: [a check from section 9 and the scope it must scan]
```

6. **Artifacts.** If artifact writes exist, publish the Task-ID checklist as the task artifact and a short implementation plan artifact that links to the repo plan. The repo files remain the source of truth.
7. **Readiness gate.** Write the evidence for each line into the ledger; on any failure, fix it and check again:
   - Every Task-ID maps to at least one subtask, and every subtask to a Task-ID.
   - Owned-file sets are pairwise disjoint within each wave, and shared indexes are lead-owned (R6).
   - Every existing path in a subtask exists; new files are marked NEW.
8. Go straight to Phase 2. A run that stops after planning is incomplete.

## 6. Phase 2: Worker Waves (Phase 2 budget)

1. **Build a wave.** Take up to `A × H` ready subtasks whose owned files do not overlap, at most `H` per worker.
2. **Dispatch the wave in one `invoke_subagent` call.** The only fields are `TypeName`, `Role`, `Prompt`, and the optional `Workspace`. There is no `Model` field. Keep `Workspace` at `inherit`: `branch` creates an isolated git worktree that this prompt has no merge step for, so use it only when the user asks.

```json
{
  "Subagents": [
    { "TypeName": "self", "Role": "Worker 01 Auth Module", "Prompt": "<worker brief for Subtasks 01 and 02>", "Workspace": "inherit" },
    { "TypeName": "self", "Role": "Worker 02 Billing Module", "Prompt": "<worker brief for Subtasks 03 and 04>", "Workspace": "inherit" }
  ]
}
```

Discovery uses the same shape with `"TypeName": "research"`. Use a custom TypeName only if your tool description lists it; otherwise use `self` with a specific Role.

3. **Worker brief.** Subagents start with no chat history, so each Prompt must stand alone:

```text
You are Worker 01 for task NN-<slug>. You have no chat history; this brief is everything.
Read first: .ai-memory/coding-guidelines.md, then your subtask files:
  .ai-memory/plans/subtasks/NN-<slug>/01-<name>.md
  .ai-memory/plans/subtasks/NN-<slug>/02-<name>.md
Owned files (edit only these): <paths>. Everything else is read-only.
Never: run builds or test suites; run git add, commit, push, reset, checkout, clean, or stash;
edit shared index files; spawn subagents; invent scripts, flags, or paths (check first).
Verify with: <exact check commands from the subtask>. A check that scans 0 files is a FAIL.
Reply with exactly this block, once per subtask, then stop:
TASK: Task-01 | STATUS: DONE | FAILED | BLOCKED
FILES_CHANGED: <paths>
CHECKS: <command> -> exit <code>, <files scanned>
ACCEPTANCE: AC1 PASS <evidence>; AC2 PASS <evidence>
ASSUMPTIONS: <list or none>
BLOCKERS: <list or none>
```

4. **Yield.** After any dispatch, print one line (for example `Dispatched Worker 01 and Worker 02 for Subtasks 01 to 04; waiting for their results.`) and stop calling tools. Results arrive as messages (R12). Apart from blocking questions (R14) and the final report, this is the only time the lead ends a response without a tool call. If a worker is waiting on a permission prompt, that is the user's action; do not dispatch it again. Use `manage_subagents` only to recover after an interruption.
5. **Accept a report only after checking it yourself:**
   - `git diff --stat -- <owned files>` matches FILES_CHANGED, and `git status --porcelain` lists no file outside the stage list, the preflight list, and this wave's owned files.
   - Re-run the subtask's check and confirm exit 0 with a non-zero file count.
   - Every acceptance criterion has concrete evidence.
6. **Reject.** Send the specific failures back to the same worker with `send_message`; it keeps its context. A worker that touched files outside its box must revert its own edits (R7). After two failed rounds, mark the subtask FAILED, write the RCA, and continue (R13).
7. **Record.** Update the ledger (status, evidence, stage list), then run `python 03-ai-scripts/33-test-inventory-generator.py --record <changed files>`.
8. **Loop.** Dispatch the next wave until every subtask is DONE or FAILED, or the budget runs out. If it runs out, write the resume point into the ledger and the plan, and report what is left.

## 7. Phase 3: Consolidate and Commit

1. **Consolidate.** Merge the finished subtasks into `.ai-memory/plans/completed/NN-<slug>.md` with a header that says how the task started, links the spec, and states the real step count from the ledger. Delete the subtask folder and the pending plan, and keep the spec. Move the entry in `.ai-memory/plans/readme.md` and add it to the Recent Completed Tasks Register there.
2. **Indexes.** Register every new spec, standard, or prompt in its folder index during this run. Add to `.ai-memory/what-to-read.md` only material that future agents must read before a task, never temporary files or subtasks.
3. **Stage (R8).** Drop from the stage list any file that was created and deleted during this run. Run `git add -- <every path in the stage list>`, then confirm that `git diff --cached --name-only` matches the stage list exactly. Unstage anything extra with `git restore --staged -- <path>`.
4. **Commit once (R9)** with a semantic prefix (`Feature:`, `Fix:`, `Docs:`). `gitmap cpf "<summary>"` and `gitmap cpb "<summary>"` stage everything, so use them only when preflight recorded a clean tree.
5. **Push according to `COMMIT`.** `push`: `git push origin <branch>`. `local`: commit without pushing. `none`: leave the changes staged and list them. Never force push, and never start a release (R10).
6. **Walkthrough.** If artifact writes exist, publish a walkthrough artifact with the summary and the check evidence.

## 8. Final Report (in chat, in this order)

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]` — [check or diff evidence]
- ❌ **Task-02: [Descriptive Task Title]** — `[Failed]` — [reason and RCA link]
- ⏳ **Task-03: [Descriptive Task Title]** — `[Deferred]` — [what blocks it]

### Modified Files Summary

- [relative/path/one.ext]

### Steps Used

- Phase 1: [x] of [PHASE_1_BUDGET], Phase 2: [y] of [PHASE_2_BUDGET] (from the ledger)

### Implementation Confidence Score

- Confidence: [checklist items with evidence] / [all checklist items] = [percent]
- Gaps: [each checklist item without evidence, or none]
```

Every task gets its own `- ` bullet on its own line; never run tasks together on one line. Take the modified files from `git show --name-only HEAD`, or from the staged list, never from memory. Then emit this audit prompt with the links filled in:

```markdown
### Independent AI Audit & Verification Instructions

You are an independent auditor. Verify the implementation against the spec and the verbatim request, and fix the gaps directly.
- Spec and verbatim request: [02-spec/21-app/NN-<slug>.md](02-spec/21-app/NN-<slug>.md)
- Consolidated plan: [.ai-memory/plans/completed/NN-<slug>.md](.ai-memory/plans/completed/NN-<slug>.md)
- Changed files: [relative/path/one.ext](relative/path/one.ext)
1. Read the spec in full, especially the verbatim request and the acceptance criteria.
2. Compare each changed file line by line with each requirement, and list every missing, partial, or non-compliant item with its file and line.
3. Fix the gaps without asking, following `.ai-memory/coding-guidelines.md`, then re-run the targeted checks.
4. Score: Verbatim Adherence [X/100], Completeness [Y/100], Guideline Compliance [Z/100], Overall [(X+Y+Z)/3].
5. End with PASS or FAIL, your confidence, and the evidence behind it.
```

## 9. Targeted Checks (R2)

Confirm that each script exists before using it (R4), and scope each check to what this task changed.

- **Newlines and boolean names (no edits):** `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only --ext <.ext>`. Pass a folder, not a file: a single file is only scanned if it is already in the script's cache. Confirm the summary counts your files.
- **Prompt, skill, and plan markdown:** `python linter-scripts/check-prompts-loaded.py`, `python linter-scripts/check-prompt-and-spec-paths.py`, and `python linter-scripts/check-sequence-integrity.py`.
- **Forbidden strings:** `python linter-scripts/check-forbidden-strings.py`.
- **Language linters** that already exist in the repo (for example in `linter-scripts/` or `linters-cicd/`), run on the changed files only.

## 10. Discovery Tools

- **GitMap (when preflight found it):** `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r]` (streaming live search; replaces `Select-String` / `git grep`), `gitmap search "<term>"`, `gitmap find "<pattern>" -ext <ext>`, `gitmap lf <pattern>`, `gitmap cat <file>`. TOTAL BAN on PowerShell `Select-String`, `Get-ChildItem -Recurse`, `git grep`, `grep`, or `findstr`.
- **Python fallback:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`, `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<text>" --limit 50`, `python 03-ai-scripts/17-fast-file-reader.py --read-file <file> --max-bytes 100000`, `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`.
- Before writing a helper script, check `03-ai-scripts/readme.md` for an existing one. Put independent reads in the same response. Never run `gitmap pa` or `gitmap pae` unless the user asks.

## 11. RCA Routing (bugs, failures, and "fix with RCA")

- **CI/CD, workflow, lint-gate, or runner failures:** `.ai-memory/cicd-issues/NN-<slug>.md`, indexed in `.ai-memory/cicd-index.md`.
- **Application bugs (logic, UI, CLI, API):** `02-spec/22-app-issues/NN-<slug>.md` with Reproduction, Cause, Fix, and Prevention, indexed in `02-spec/22-app-issues/readme.md`.
- **Failed subtasks (R13):** the same routing, linked from the ledger and the plan.

## 12. Final Checklist

Each item needs evidence in the ledger; an item without evidence counts as failed.

- [ ] Precedence: the instructions above the prompt were captured verbatim and followed.
- [ ] R5: `invoke_subagent` was used for discovery and edits, or `SOLO_FALLBACK` was logged with its reason.
- [ ] R6: owned-file sets were disjoint in every wave, and only the lead edited shared indexes.
- [ ] R1 and R2: no builds or test suites ran, and every check exited 0 with a non-zero file count.
- [ ] R3: every Task-ID has evidence, and no template tokens such as `[Descriptive Task Title]`, `<slug>`, or `TODO` remain in changed files.
- [ ] The spec, plan, subtasks, and indexes were written and registered during this run.
- [ ] Changed code follows `.ai-memory/coding-guidelines.md`, and every worker violation was rejected.
- [ ] Blast radius: callers of every changed symbol were searched via `gitmap aum search "<symbol>" [dir]` and updated (never `Select-String` or `git grep`).
- [ ] R8 and R9: the staged list equals the stage list, there is one commit, and it was pushed according to `COMMIT`.
- [ ] R10 and R15: no release, no version edit, and no generated artifacts committed.
- [ ] The final report lists tasks line by line, modified files from git, step counts from the ledger, and the audit prompt.

## MUST FOLLOW NON-NEGOTIABLE

This block is the user's own wording, kept verbatim except for one clause changed to match R1.

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run the targeted checks (builds and full unit tests stay in CI per R1), group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
