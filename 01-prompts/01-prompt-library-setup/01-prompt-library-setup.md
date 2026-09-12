# Library Setup & Prompt Architecture Scaffold — Setup Spec (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

Create a prompts folder in the root of the repo: `01-prompts`. Inside this you will have the general prompts. So inside this folder will be `01-general`, and afterwards you will have the project names, let's say `02-<project-name>`. Inside this you will always have the file as a sequence: `01`, `02` hyphen the prompt slug. So whatever prompt I give you, you create the slug or title for that and put it into that prompt folder specifically and make a commit.

Update your `.lovable` folder with this memory in an MD file. Make sure you also create a "what to read" MD file that would resolve this. Also in the root of the repo, create the root readme file — all readme files should be in lowercase. Remember: hyphen and lowercase, that's the format. In the root of the repo, explain the folder structure and how the files need to be. Empty folders you create with a `.gitkeep` file so that they exist.

Make sure that any AI who reads the root readme file will know what the memory structure is, which file it needs to read, and where the prompt is and how the prompt is.

Every time you put the prompt, you just do a proofread, nothing much. Mention this in the memory: what it needs to do if a prompt is given with "proofread prompt" or "next prompt". Usually the project name is given in the prompt. If no project name, then it is a general prompt.

Inside the general prompts, we can have other folders like bug fix or things like that. Every one of these folders can be `01`, `02` sequence. It could be bug fix, DRY code, coding guidelines, things like that. All kinds of folders can be there, and they can have a `.gitkeep` file. Based on the category, you put that prompt in there. If it is mentioned for a project, then you name the project name as the folder as mentioned — the sequence, hyphen, the name of the project in lowercase. Inside this you have `01`, `02` sequence and the prompts.

Whatever I'm saying, write it into the memory so that any AI who reads it understands everything about the project structure.

Any time I give you the prompt, you just remove the filler words. You don't do anything else. Keep the exact words exactly as they are. Then you put the action items and checked items that must be followed. Create a section called "Action Items — Must Follow (Non-Negotiable)" and put all the terms that have been discussed. Make sure it is written as a checklist that AI should follow. If there is any folder structure that is defined or discussed, that should also be in a section like "Folder Structure". If there is a database segment or database design discussed, put that at the end of the prompt. If nothing is discussed regarding databases, skip that. If no folder structure is discussed, skip that.

If any link or file is given, the AI should write that file into the specific folder. If in confusion, it can ask. Usually asset images should go into the assets folder. Specs should go into the spec folder. Inside the spec folder there should be folder `21`, and inside this the app spec should be there.

Always before writing the code, the AI should read the coding guidelines, understand the coding guidelines, and error management must be followed. This needs to be added with every prompt — that is, inside the spec folder, folder `02`, folder `03` must be followed, and `04` as well. Just mention these folder numbers. That would be all right. These are very important when writing the code. Make sure the code needs to be DRY. If it is dealing with code and spec, every time the spec needs to add that.

Inside the prompts folder, these prompts which I have shared here should go as a sample prompt inside the general prompts, in the prompt library setup. Inside the general folder there should be a prompt library setup, and it should have this MD file with this slug, like `01-prompt-library-setup.md`, that contains what we just discussed. Also refer to this from the root readme file and from inside the `.lovable` folder "what to read".

## Action Items — Must Follow (Non-Negotiable)

- [ ] Create `01-prompts/` at the repo root as the prompt archive.
- [ ] Store general prompts under `01-prompts/01-general/`; project prompts under `01-prompts/<NN>-<project-name>/`.
- [ ] Route by project name: project name mentioned in the prompt -> project folder; no project name -> general.
- [ ] Inside `01-general/`, group by category folder (`01-bug-fix`, `02-dry-code`, `03-coding-guidelines`, ...), each with a two-digit sequence prefix.
- [ ] Name every folder and file with a two-digit sequence prefix, hyphen, lowercase slug (`01-prompt-library-setup.md`).
- [ ] All readme files are lowercase `readme.md`.
- [ ] Keep every empty folder tracked with a `.gitkeep` file.
- [ ] On each new prompt: remove filler words only. Keep the exact wording otherwise. No rewriting, no summarising.
- [ ] Every prompt file has: the proofread prompt, then "Action Items — Must Follow (Non-Negotiable)" as a checklist.
- [ ] Add a "Folder Structure" section only if a folder structure was discussed.
- [ ] Add a "Database" section at the end only if database design was discussed.
- [ ] For any code-related prompt, add the standard footer: read spec folders `02`, `03`, `04` before writing code; error management must be followed; code must be DRY.
- [ ] Files/links supplied with a prompt: assets/images -> `assets/`, specs -> `spec/` (app spec in `spec/`). Ask when placement is unclear.
- [ ] Commit after storing each prompt.
- [ ] Keep the root `readme.md`, `.lovable/what-to-read.md`, and `.lovable/prompts.md` in sync so any AI can discover the structure.

## Folder Structure

```text
01-prompts/
  01-general/
    01-prompt-library-setup/
      01-prompt-library-setup.md
    02-bug-fix/            (example category, .gitkeep when empty)
    03-dry-code/           (example category, .gitkeep when empty)
    04-coding-guidelines/  (example category, .gitkeep when empty)
  02-<project-name>/
    01-<prompt-slug>.md
    02-<prompt-slug>.md
assets/                    images and other binary assets
spec/
  02/                      coding guidelines
  03/                      error management
  04/                      additional mandatory rules
  21/                      app spec
readme.md                  root, lowercase
.lovable/memory/
  what-to-read.md
  prompt-library.md
```

## Before Writing Code (applies to every code instruction)

Read and follow spec folders `02`, `03` and `04` before writing any code. Error management must be followed. Code must be DRY.

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
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

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

---

## Instructions Index

| # | File | Title | Purpose | Status |
|---|------|-------|---------|--------|
| 001 | [001-conversation-log-and-context-wrapper.md](../../prompts/001-conversation-log-and-context-wrapper.md) | Conversation Log & Context Wrapper | Staging prompt for conversation log persistence and context rewriting | active |
| 001 | [001-interactive-macro-builder-commands](../../prompts/001-interactive-macro-builder-commands.md) | Interactive Macro Builder: Dynamic PWD Header, In-Builder File Listing, and Helper Commands | Enhanced interactive macro builder UX with PWD display, directory file listing, and built-in search/find/replace commands | active |
| 002 | [002-gitmap-interactive-builder-and-clean-sync](../../prompts/002-gitmap-interactive-builder-and-clean-sync.md) | Interactive Macro Builder Verification, Git Hygiene & Atomic Release Sync | Final verification, clean artifact exclusion, single commit grouping, and remote synchronization | active |
