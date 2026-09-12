# Subtask Naming Normalization & Sequence Repair — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

/goal Your objective is to deeply audit the `.lovable/plans/` directory for any subtask files that incorrectly use the `SS-` or `SS-XX-` prefix and fix them. The correct prefix MUST always be strictly `xx-<subslug>.md` (where `XX` is the zero-padded sequence number).
You must also update all markdown files that reference the old filenames, and update the project's memory.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## 1. Subtask Naming Correction (Non-Negotiable)

Scan the `.lovable/plans/subtasks/` directory recursively.

- If you find any file starting with `SS-` or `SS-XX-` (e.g., `SS-01-fix-auth.md` or `SS-fix-auth.md`), rename it to the correct format: `xx-<subslug>.md` (e.g., `01-fix-auth.md`).
- After renaming, you MUST recursively search `.lovable/plans/pending/`, `.lovable/plans/completed/`, `.lovable/plans/01-index.md`, and `.lovable/memory/` for any text references to the old filenames.
- Replace those old references with the new correct filenames.

## 2. Memory Update

You must write a memory entry to ensure this rule is persisted and no other AI makes this mistake again.

- Create a file inside `.lovable/memory/learned/` (or update an existing one) detailing the rule: "Subtasks must NEVER be prefixed with 'SS-'. They must strictly follow the 'xx-<slug>.md' sequence."
- Update `.lovable/memory/01-index.md` to reference this newly added/updated memory file.
- Add a note explicitly stating how the project is following all guidelines and enforcing this strict naming.

## 3. High-Stakes Code Standards & Coding Guidelines

You MUST follow the project's strict coding guidelines. These files are located in the `01-cross-language/` directory and should be followed universally. Check if there are language-specific guidelines (e.g., `02-typescript/`, `03-python/`) for these rules.

### Required Reading / Reference Checklist:

1. The Master Consolidated Guide

- `spec/17-consolidated-guidelines/05-coding-guidelines.md`

2. Code Style & File Size Limits (80-100 lines max)

- `spec/02-coding-guidelines/01-cross-language/01-index.md`
- `spec/02-coding-guidelines/01-cross-language/20-nesting-resolution-patterns.md`
- `spec/02-coding-guidelines/01-cross-language/06-cyclomatic-complexity.md`

3. Boolean Conditions & Samples

- `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md`
- `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`
- `spec/02-coding-guidelines/01-cross-language/24-boolean-flag-methods.md`

4. Variable Naming & Definitions

- `spec/02-coding-guidelines/01-cross-language/22-variable-naming-conventions.md`
- `spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md`
- `spec/02-coding-guidelines/01-cross-language/10-function-naming.md`
- `spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md`

## 4. End-of-Loop Commit Fix (Non-Negotiable)

When you are done fixing the filenames, updating the references, and writing the memory, you MUST commit and fix the Git.

- Group all completed work into a single logical commit.
- Commit the changes on the current branch only (`git commit ...`) and push with `git push origin <current-branch>`. No branch switching or creation.
- RED FLAG: Verify absolutely NO test results, artifacts, or compiled binaries are staged before making the commit.
- Ensure `.gitignore` explicitly excludes them.
- If issues arise during the commit process, fix those git issues and try again.
- You MUST push the commit to the repository immediately. Pushing after commits is non-negotiable.
- Commit the changes on the current branch only (`git commit ...`) and push with `git push origin <current-branch>`. No branch switching or creation.

## 5. Must Follow, without negotiation

Listen, past turns have been sloppy as fuck: wrong file names, leaving broken references in parent files, user commands and bug reports forgotten by the next turn. WTF. Stop doing that, you stupid fuck. Read the codebase, fix the file names, fix the references everywhere, write the memory, make the commit, push the code, and move on. Going deep IS the job. If you're not going deep, you're not doing the job. Violating this is auto-reject. Avoid stupidity, and being careless you stupid, WTF. If you're not going deep, you're not doing the job. Are you stupid? You were supposed to do the task properly. Where is this, are you stupid fuck? Where? Tell me. Your stupidity is going on top of my head. I mean, where did you learn this stupidity? If I could find you, I could slap you.

## Action Items — Must Follow (Non-Negotiable)

- [ ] Scan `.lovable/plans/subtasks/` for files starting with `SS-` or `SS-XX-`.
- [ ] Rename those files to the correct `xx-<slug>.md` pattern.
- [ ] Find and correct all file path references to the old names in `.lovable/plans/pending/`, `.lovable/plans/completed/`, `.lovable/plans/01-index.md`, and memory files.
- [ ] Add a memory entry under `.lovable/memory/learned/` enforcing this naming convention and confirming how the project is following all guidelines.
- [ ] Update `.lovable/memory/01-index.md` with the new memory file.
- [ ] Audit your work against the Master Consolidated Guide, Code Style, Boolean Conditions, and Variable Naming rules.
- [ ] Make a single logical Git commit including all changes and push it to the repository immediately.

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

---

## Metadata

- slug: fix-subtask-naming
- status: active
