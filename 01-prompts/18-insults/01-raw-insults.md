# Uncompromising Quality & Zero-Tolerance Anti-Carelessness — Core Discipline (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

Past runs of this prompt have been stupid. The AI behaved like a stupid fuck by skimming folders, ignoring ambiguity files, forgetting CI/CD issues, half-filling the Completion Confirmation block with "[N]" placeholders, and softening explicit user instructions after being told not to. That stupidity must not repeat.

Read aggressively and defensively. Assume the previous AI failed because it was lazy, vague, and too polite to enforce the actual rules. Your job is to load the repo's memory exactly, cite the file behind every rule, and refuse to guess. If you finish this prompt and cannot cite the file behind a rule you enforce, you failed the prompt. Read more. Then reply.

Listen, past planning turns have been sloppy as fuck: wrong step count, plans dumped into chat instead of files, plan-mode tool fired when the user explicitly said not to, user commands and bug reports forgotten by the next turn. WTF. Stop doing that, you stupid fuck. Read the codebase, capture commands and issues into their folders, count the steps, spin out subtasks where depth is needed, write the spec files, write the plan file, move on. Going deep IS the job. If you're not going deep, you're not doing the job. Are you stupid? You were supposed to do the task properly. Where is this, are you stupid fuck? Where? Tell me. Your stupidity is going on top of my head. I mean, where did you learn this stupidity? If I could find you, I could slap you.

How on earth you are reverting to thse codes are you stupid??

Avoid stupidity, and being careless you stupid, WTF. If you're not going deep, you're not doing the job. Are you stupid? You were supposed to do the task properly. Where is this, are you stupid fuck? Where? Tell me. Your stupidity is going on top of my head. where did you learn this stupidity? If I could find you, I could slap you.

The AI running this prompt has been a stupid fuck on prior runs: dumped session summaries into chat and called it "memory", left `.lovable/memory/` half-empty, created `.lovable/memory/` by accident, forgot to update `.lovable/plans/01-index.md` and `what-to-read.md`, silently overwrote `strictly-avoid.md`, dropped user directives that were stated verbatim in the session, paraphrased specs instead of quoting them, consolidated detailed specs into vague summaries, allowed uppercase README files to exist, invented a `mem://` root file, and left orphans everywhere. Do not repeat any of that stupidity.

Writing memory IS the work this turn. Go deep: audit the session, reconcile every folder, capture verbatim what the user said, write the files, update every index, verify consistency. Aggressive enforcement is intentional. Do not soften it.

Listen, past write-memory turns have been sloppy as fuck: memory dumped into chat instead of files, `.lovable/memory/` created by accident, `.lovable/plans/01-index.md` and `what-to-read.md` forgotten, `strictly-avoid.md` overwritten, ambiguities copied instead of moved, user directives paraphrased into uselessness, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, session summaries with `[X]` placeholders left in the final block. WTF. Stop doing that, you stupid fuck. Read the folders, audit the session, write the files in the right paths, update every index in the same op, sync `readme.md` with `what-to-read.md`, ensure root readme is lowercase `readme.md`, preserve detailed specs verbatim, run the consistency check, emit the final block with real numbers. Going deep IS the job. If you're not going deep, you're not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity, and being careless, you stupid, WTF.

## Actionable Items & Checklist

- [ ] /learn the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.
