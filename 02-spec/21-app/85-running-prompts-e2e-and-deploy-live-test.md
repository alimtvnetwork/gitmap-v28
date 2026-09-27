# Canonical Spec: Antigravity Running Prompts E2E Verification & Smart Deploy Polish

**Spec ID:** SPEC-APP-85
**Status:** Active
**Author:** Antigravity Master Orchestrator
**Date:** 2026-09-28
**Version:** v6.364.0

---

## 1. User Request (Verbatim)

```text
is it done properly???

Okay. Currently, most of the things are working fine. I think it's good, but I think I do have some requests that you could modify. First of all, the gitmap update all. It could be hyphen all or space all or UA. This would do the same thing. But now coming to the point how it's going to do it, the update is a delicate situation. So if we are running from the external machine, we should actually use JSON as a communication so that the help text does not get included in it. And currently, the way that you are doing the formatting or displaying it looks very poor. Okay? So I think this is where we need to work on. We need to make sure the help and things, these are really high quality. Okay? So currently, the output of the part two, how you're doing it, it's quite poor. Okay? And this is where we need to work on. And also make sure that, at the end, you bump the minor version and make a release and check the Gitmap PE. Also now I wanted you to check, so you can even run a test command on, let's say, scripts picture I or something like this. And you check by the running prompts LS that this project is actually running and other projects, but this project needs to be there. Or you could do other stuff like the own project. Let's say the status, alien status sample, just put a hi there and check if this comes up in the running prompt when it is running or put some weight. Okay? Like sleep and then say hi. You can do that so you can trace back. So that's one thing. Another is backing up the running prompts. So each one of the prompts from all the projects, the running prompts would be backed up. Then you can restore the prompt, running prompt, that should inject and run the prompt that we have taken the backup. So make sure that you do the end-to-end testing so that you can confirm it, that it's working very well and there is no confusion. Okay? It's been several times that we are trying that. I hope there should be no issues, right? If possible, you can also test out with some file the Gitmap SC or deploy command, sync left, sync right command. Okay. And also update the help text and also UI help text. These are very important. Do not miss it. Is it clear? Can you please do that? At the end, do a minor bump and a release and do the Gitmap PE to check the errors. Okay.
```

---

## 2. Technical Architecture & Invariants

1. **Running Prompts Traceability & Discovery:**
   - Antigravity queue files (`agy-prompt-queue.json`) and active fix prompt files (`active-agy-pipeline-fix-prompt.txt`) across registered projects are read by `cli/cmdagy/agy_running_prompts_io.go`.
   - `gitmap agy running-prompts ls` aggregates and renders items in a clean table showing Project, Status (`RUNNING` / `ENQUEUED`), Word Count, and Truncated Prompt Preview.
2. **Split-DB Backup & Restoration Invariant:**
   - Backup writes batches to SQLite database `data/backup-prompts/sql.db` using schema `PromptBackupBatch` and `PromptBackupItem`.
   - Restore reads the latest active batch, resolves each item's project queue file via `resolveQueueFilePath(projectPath)`, injects un-enqueued prompts, and marks the batch as restored in `PromptRestoreLedger`.
3. **Smart Deploy Command Matrix:**
   - `gitmap deploy <target> <src> <dest>`: bidirectional or default push.
   - `gitmap deploy-right <target> <src> <dest>`: push newer from left to right only (`--sync-right`).
   - `gitmap deploy-left <target> <src> <dest>`: pull newer from right to left only (`--sync-left`).
   - `gitmap sc deploy <subcommand>`: fleet broadcast wrapper.
4. **Help Menu Visual Standard:**
   - Cyan/yellow boxed two-column terminal menus rendered via `termhelp.HelpMenu`.
