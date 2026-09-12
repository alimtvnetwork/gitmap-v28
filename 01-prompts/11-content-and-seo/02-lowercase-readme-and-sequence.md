# Lowercase Filename Enforcement & Sequence Re-Ordering — Workflow (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

Enforce these naming rules across the entire repository:

1. All README files must be lowercase: rename every `README.md`, `Readme.md`, `ReadMe.md`, etc. to `readme.md`. Apply recursively at every depth (root, subfolders, packages, specs, prompts, scripts - everywhere). Update every internal link and import reference to match.

2. Sequence-prefixed markdown files must use `xx-lower-case.md` slug form: any markdown file that begins with a numeric sequence prefix must follow the pattern `NN-kebab-lower-case.md` where:

   - `NN` is a two-digit zero-padded number (`01`, `02`, ..., `99`)
   - The remainder is all lowercase, words separated by single hyphens (`-`)
   - No spaces, no underscores, no PascalCase or camelCase, no uppercase letters
   - The `.md` extension is lowercase

   Examples:

   - ✅ `01-overview.md`, `05-coding-guidelines.md`, `13-cicd-pipeline.md`
   - ❌ `1-Overview.md`, `01_Coding_Guidelines.md`, `13-CICD-Pipeline.MD`

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## Execution Steps

1. Scan the whole repo for non-conforming filenames (case-insensitive `readme` not equal to `readme.md`, and any `^\d+[-_ ]` markdown file not matching `^\d{2}-[a-z0-9]+(-[a-z0-9]+)*\.md$`).

2. Rename each offending file using `git mv` (preserve history).

3. Update every reference: markdown links, code imports, doc indexes, sidebars, and `.lovable/memory/01-index.md`.

4. Verify with a final scan - fail loudly if any non-conforming file remains.

5. Run the build and link checker; fix any broken references.

## Important

- Do not skip nested folders.
- Do not leave both `README.md` and `readme.md` (case-only renames on case-insensitive filesystems require a two-step `git mv`).
- Bump the minor version of the codebase after this change, per repo convention.

---

*This prompt is version 1.0.*

---

title: Lowercase Readme And Sequence Slugs
slug: lowercase-readme-and-sequence

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
