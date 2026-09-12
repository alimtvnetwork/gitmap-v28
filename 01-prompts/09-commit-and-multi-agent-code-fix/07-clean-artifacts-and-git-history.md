# Artifact Sanitization & Git History Preservation — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Ensure that NO assets, zip files from artifacts, test data, temporary scratch scripts, or extraneous generated code are accidentally committed to or retained in the Git repository.

If any temporary files, artifact dumps, or unwanted generated files are detected in the workspace or git index, the AI MUST NOT silently commit them or guess about their disposition. Instead, execute the following structured protocol:

---

### Step 1: Detect and Itemize Unwanted Artifacts & Temporary Files

Scan the entire repository, untracked files, and recent commit history for:

- Artifact zip files, archive bundles, or binary dumps.
- Temporary test data files, scratch scripts, mock output files, or debug logs.
- Extraneous auto-generated code, intermediate build artifacts, or temporary caches.

For EVERY detected candidate file, compile an explicit report detailing:

1. File Path & Name: Exact location in the workspace.
2. File Size: Size in bytes / KB / MB.
3. Creation Step & Origin: Exactly how and at which step/command it was generated.
4. Purpose & Context: Why the file was originally created during the task.

---

### Step 2: Present Positively Framed Questions to Developer

Present the itemized list to the developer. All questions MUST be framed with strictly positive phrasing (following the project's boolean and coding standards in `spec/02`, with zero negative words, no `!`, no double negatives, and no inverted questions).

Example positive review format:

```text
[Item 1]
- Path: /assets/artifacts/test-bundle.zip
- Size: 14.2 MB
- Generated At: Step 3 (Artifact packaging run)
- Purpose: Temporary zip archive created for test packaging.
Question: Would you like to keep this file in the project? (Yes / No)
```
or
```text
Question: Would you like to remove this file from both the filesystem and Git history? (Yes / No)
```

Never ask negative or inverted questions such as "Do you want to avoid not deleting this?" or "Is this not needed?". Keep all questions clear, direct, and positively stated.

---

### Step 3: Dual Removal (Filesystem + Git History Purge)

When the developer indicates that a file should be removed:

1. Filesystem Removal: Delete the file from the local working tree.
2. Git History Purge: Remove and purge the file from Git history and index (e.g., using `git rm --cached` or Git history filtering) so that large binary blobs or unwanted scripts are not permanently bloating the repository's git object database (`.git` size).
3. Verify Clean Repository State: Verify `git status` and repository size to confirm the files are completely eradicated.

---

### Step 4: Verification and Git Push

1. Ensure the working tree is clean and builds pass without errors.
2. Ensure no unintended artifacts remain in tracked files.
3. Commit clean state with a clear descriptive commit message.
4. Push all changes to the remote Git repository.

---

## Action Items — Must Follow (Non-Negotiable)

- [ ] Scan the entire repository and git working tree for artifact zip files, test data, temporary scratch scripts, and extraneous generated code.
- [ ] For each detected candidate file, itemize its full path, exact file size, generation step origin, and creation context.
- [ ] Present the candidate files to the developer with strictly POSITIVE question framing (no negative phrasing, no double negatives).
- [ ] If removal is chosen, execute dual removal: delete the file from the filesystem AND purge it from Git history to prevent repository bloat.
- [ ] Verify that `.git` repository size is minimized and no dangling binary blobs remain in tracked history.
- [ ] Confirm that all tests and builds pass following the removal.
- [ ] Commit all clean changes with a concise, descriptive commit message.
- [ ] Push the clean state to the Git repository before finalizing the task.

---

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

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

## No Automatic Releases (Strict Policy)

You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Before Writing Code

Read and follow spec folders `spec/02-coding-guidelines/`, `spec/03-error-manage/`, and `spec/04-database-conventions/` before writing any code. Error management must be followed. Code must be DRY.

---

## Metadata

- slug: clean-artifacts-and-git-history
- status: active
