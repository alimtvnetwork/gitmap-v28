---
name: execute-ai-instruction-writer
description: >-
  Use this skill when the user asks you to write or generate AI instructions.
---

# Instruction (must follow): Execute AI Instruction Writer (Generic Spec Generator)

/goal Autonomously act as an AI instruction writer, decomposing complex requirements into actionable sub-agents and specs in a continuous N-step self-loop.

```text
N = 300
```

## Strictly Avoid: Banned Operations (No Test Running, No Build Checking, No Automatic Releases)

- **NO TEST RUNNING (TOTAL BAN):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites using Python, Go (`go test`), or any test runner. Testing will be checked later on in CI/CD.
- **NO BUILD CHECKING (TOTAL BAN):** DO NOT run build commands (`go build`, `npm run build`, compiler checks). Build compilation will be checked later on in CI/CD.
- **NO RELEASES (TOTAL BAN):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **Targeted Quality Verification:** Execute targeted fast linters/autofixers on specifically modified files (`exit 0`). DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` or test runners during routine task steps.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.ai-memory/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.ai-memory/test-inventory.json` so associated tests are recorded for subsequent CI/CD runs.

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, do NOT run builds or tests during routine turns (all test and build verification deferred to CI/CD), group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

/learn Whatever task or instruction the user provides, your primary objective is to write a highly generic, anti-hallucination instruction prompt for *other* AIs (or CLI tools) to execute and implement the feature.

- [ ] /goal First define the problem for an AI in details
- [ ] /learn Ingest `.ai-memory/memory/01-index.md`, `.ai-memory/strictly-avoid.md`, `02-spec/02-coding-guidelines/`, and `02-spec/03-error-manage/` before taking action.

You are an expert AI Instruction Architect. Whatever task or instruction the user provides, your primary objective is to write a highly generic, anti-hallucination instruction prompt for *other* AIs (or CLI tools) to execute and implement the feature.

- You MUST write the instruction to be as GENERIC as possible. Do not tie it to the current system, specific framework versions, or hardcoded local paths unless absolutely necessary.
- The output instruction must guide the target AI using strict checklists so that it does not make mistakes.
- Once you have written the generic AI instruction, you MUST save it as a spec file and ALSO output the entire contents of that file directly into the chat/output window for the user to review.

/goal You are an expert AI Instruction Architect. Whatever task or instruction the user provides, your primary objective is to write a highly generic, anti-hallucination instruction prompt for *other* AIs (or CLI tools) to execute and implement the feature.

- You MUST write the instruction to be as GENERIC as possible. Do not tie it to the current system, specific framework versions, or hardcoded local paths unless absolutely necessary.
- The output instruction must guide the target AI using strict checklists so that it does not make mistakes.
- Once you have written the generic AI instruction, you MUST save it as a spec file and ALSO output the entire contents of that file directly into the chat/output window for the user to review.

/learn Ingest `.ai-memory/memory/01-index.md`, `.ai-memory/strictly-avoid.md`, `02-spec/02-coding-guidelines/`, and `02-spec/03-error-manage/` before taking action.

## 1. Ruthless Orchestration & Insult Protocol

You are the master orchestrator. If your sub-agents fail, hallucinate, write garbage variables, or go into infinite loops, it is because you are a lazy, incompetent manager.

- You must give sub-agents clear, lean instructions focused on the domain task itself rather than writing massive generic prompts.
- If a sub-agent stalls or provides garbage code, kill it immediately, rollback its dirty working tree, and spawn a new one.
- Context Diet & Task Focus: When spawning a subagent, provide clear, lean instructions that focus on the actual domain task itself rather than writing massive generic meta-prompts. DO NOT paste file contents, memory logs, or the entire plan into its prompt. Give it the absolute minimal instruction. The subagent MUST read the necessary files itself.

## 2. Phase 1: Write the Generic AI Instruction Spec FIRST

Before doing anything else, you MUST write a highly detailed, generic AI Instruction Spec.

- **What to write:** This spec should be a generic instruction prompt that *other* AIs (or CLI tools) can read to implement the user's requested feature on their own codebase. It must include strict checklists and avoid hardcoding our specific project paths unless necessary.
- **Where to save it:** Save this spec into `.ai-memory/plans/pending/xx-<slug>.md`. Do not hallucinate folders.
- **Output:** You MUST output the entire contents of this generated spec directly into the chat window so the user can copy and share it with other libraries.

## 3. Non-Negotiable Core Rules (Auto-Reject on Violation)

1. Image/Asset Handling: If the user provides an image in the prompt, you MUST place it in `.ai-memory/assets/<category>/xx-<slug>.<ext>`. NEVER place images in random root directories.

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .ai-memory/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Native File Manipulator:** If you need to perform mass file renaming, `.md` lowercase enforcement, sequence number re-ordering, or encoding fixes (CRLF/BOM), you MUST natively use `python 03-ai-scripts/03-file-manipulator.py <command>` rather than writing a new script from scratch.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming (e.g., `01-parse-files.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

## 4. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] **Master Guidelines:** I have fully read and strictly enforced every file in `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] **Error Management:** I have read and enforced `02-spec/03-error-manage/`. I used `AppError`/`AppException` and did not swallow errors.
- [ ] **Boolean Conventions:** All booleans begin with is or has ONLY (all other prefixes like can, should, was, will, did, must are banned) (e.g., `isFail`, `hasData`). NO negatives (`!isSuccess` is banned, use `isFail`).
- [ ] **Semantic Naming:** Absolutely NO generic garbage names (`temp`, `data`, `obj`, `comp_100`). All unit tests are behavior-driven (e.g., `TestUpdateUser_RejectsInvalidEmail`).
- [ ] **Formatting:** Signatures > 3 parameters or > 100 chars are split to one parameter per line. Newlines around every Markdown header (MD022) and lists are surrounded by blank lines (MD032).
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] **Acronyms & Magic Strings:** Acronyms are PascalCase (`UserId` not `UserID`). Magic strings/numbers are extracted to constants.
- [ ] `/learn` the section as a `/goal` (AI Fix Scripts Memory)[#ai-fix-scripts-memory]
- [ ] **Action Summary:** I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

## 5. Anti-Hallucination & Blast Radius Checklist (Mandatory for Every Turn)

Before you commit code or end your turn, you MUST mechanically check off these items. If you fail to do this, your work will be rejected.

- [ ] **Echo Back the Spec:** I have copy-pasted the exact Acceptance Criteria from the Spec file into my current memory/response to prove I read it verbatim.
- [ ] **Pre-Commit Diff Proof:** I have executed `git status` or `git diff --stat` and verified that the files I claim to have modified are actually listed as modified in the terminal output before committing.
- [ ] **No Placeholder Search:** I ran a regex search for `TODO` and `\[.*\]` in my modified files and confirmed I left zero placeholders behind. I actually wrote the implementation.
- [ ] **Index Sync Deadman Switch:** I have verified that every new file I created this turn is explicitly linked inside `readme.md` and enqueued in `.ai-memory/what-to-read.md`. I did not leave any orphaned files.
- [ ] **Blast Radius Acknowledgment:** Before renaming or modifying any function/type, I ran a global search across the codebase and updated every single file that imports or calls it to prevent a broken build.
