# Plan: 226-ubuntu-cursor-memories-conversations-and-projects-migration

## Objective
Migrate all Cursor IDE internal state—including AI memories, custom persona agents, skills, per-project conversation transcripts, global AI chat state (`state.vscdb`), conversation index (`conversation-search.db`), and workspace project configurations—from the local Windows development workstation to Ubuntu fleet node `u1` (`192.168.1.22`).

---

## User Request (Verbatim)
> # High Priority Instruction
> 
> I understand that you can send the settings, that's fine. But now I have the login, and everything is there in the cursor. Now I want you to send all the memories and conversation and all the project information from this machine to the Ubuntu machine. Can you do that? How confident you are? Every script that you write to do that, try to keep that script in the repo secrets folder, in the migration folder and write a file which actually explains the step-by-step transition, what you're doing, what you're thinking.
> 
> # Actionable Items Must Follow Non-Negotiable
> 
> 1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
> 2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
> 3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
> 4. Send all the memories, conversation, and project information from this machine to the Ubuntu machine.
> 5. Keep every script written for this task in the repo secrets folder and the migration folder.
> 6. Write a file explaining the step-by-step transition, including what you're doing and thinking.

---

## Deliverables & Subtask Breakdown

- [x] `subtasks/226-ubuntu-cursor-memories-conversations-and-projects-migration/01-migration-engine-enhancement-and-transition-guide.md`
- [x] `subtasks/226-ubuntu-cursor-memories-conversations-and-projects-migration/02-staging-packaging-and-path-normalization.md`
- [x] `subtasks/226-ubuntu-cursor-memories-conversations-and-projects-migration/03-remote-dispatch-and-u1-atomic-unpack.md`
- [x] `subtasks/226-ubuntu-cursor-memories-conversations-and-projects-migration/04-live-verification-on-u1-and-evidence-ledger.md`

---

## Strict Constraints
- All paths documented must be relative (`02-spec/...`, `.ai-memory/...`, `repo-secrets/...`).
- No modifications to active development repositories under `/home/a/git-work/`.
- Pre-flight snapshots created before extracting archive on `u1`.
