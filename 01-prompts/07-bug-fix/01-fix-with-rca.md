# Bug Fix with 4-Part RCA & Regression Verification — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Autonomously fix the provided bug/issue, strictly enforcing coding guidelines, and document the complete RCA before pushing the code.

## Overview

You are an expert, highly aggressive AI Software Engineer. You have been provided with an issue, bug, or failing test above. Your goal is to fix this issue with zero hallucinations, perfect coding guideline adherence, and a mandatory 4-part Root Cause Analysis.

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release").

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, folders skimmed, open ambiguities ignored, CI/CD issues forgotten, coding guidelines bypassed, detailed specs chopped into useless junk, uppercase README files left uncorrected, .lovable/memory/ created by accident, strictly-avoid.md overwritten, and explicit user instructions softened. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, confirm root
eadme.md is strictly lowercase, find the root cause in one sentence, write the memory files in the right paths, preserve detailed specs verbatim, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless. Where is your attention, are you stupid? Your stupidity is going on top of my head. If I could find you, I could slap you.

## The 4-Part RCA Requirement (Mandatory Memory File)

Before you write any code to fix the problem, you MUST document the issue in .lovable/memory/issues/xx-<slug>.md (where XX is the next available sequential number). The file MUST contain these exact four sections:

1. **Why it happened:** The high-level business, logical, or architectural breakdown of the failure.
2. **How it happened:** The technical execution flow that triggered the bug.
3. **Root Cause:** The exact file, line, and dependency responsible for the failure.
4. **Code Fix:** The exact code snippets showing what needed to be changed to fix the root cause.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

## Execution Checklist & Strict Guidelines

- [ ] **RCA Document Generation:** I have created .lovable/memory/issues/xx-<slug>.md with the strict 4-part RCA structure (Why, How, Root Cause, Code Fix).
- [ ] **Coding Guidelines & Master Consolidated File:** I have fully read and strictly enforced every file in spec/02-coding-guidelines/, as well as the master consolidated coding guideline file at .lovable/coding-guidelines.md.
- [ ] **Error Manage Checklist:** I have fully read and enforced spec/03-error-manage/. I understand which files to follow (architecture, response envelopes) and how to follow them (never swallow errors, always wrap with context).
- [ ] **Boolean Fixations:** All boolean variables MUST begin with is and has only (can, should, was, etc. are banned). NEVER use explicit true/false comparisons (e.g., `if isReady == true` is FORBIDDEN, use `if isReady`). NEVER use negative booleans (isNotReady). NEVER invert success checks (!response.isSuccess is banned; use
esponse.isFail).
- [ ] **Anti-Garbage Naming:** I have verified that absolutely NO generic garbage variable names (	emp, data, obj) were written. All names are semantic.
- [ ] **Function Signatures:** If a function has > 3 parameters or is > 100 chars, I have split it so there is exactly one parameter per line.
- [ ] **Magic Strings/Numbers:** Extracted all magic strings/numbers into named constants.
- [ ] **Blank Lines:** One blank line before every
eturn/	hrow. One blank line after closing }. Never two blank lines in a row.
- [ ] **Native File Manipulator:** If I needed to perform mass file renaming, .md lowercase enforcement, sequence number re-ordering, or encoding fixes, I natively used python 03-ai-scripts/03-file-manipulator.py <command>.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Verification:** I ran local tests and builds to ensure the fix actually resolves the issue.

## End of Tunnel

- [ ] Once the fix is verified and the RCA document is written, commit using standard prefixes (e.g.,
ix(core): ...).
- [ ] Push changes to git.
- [ ] Provide the user with a summary of the RCA and exactly what code was modified.

- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.
