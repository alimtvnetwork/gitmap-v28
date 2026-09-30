# Specification: Fleet Update JSON Communication, Running Prompts Backup/Restore E2E & Deploy Polish

**Specification Version:** 1.0.0
**Status:** Active
**Spec Reference:** `02-spec/21-app/179-fleet-update-json-prompts-backup-and-deploy-polish.md`

---

## User Request (Verbatim)

```text
Okay. Currently, most of the things are working fine. I think it's good, but I think I do have some requests that you could modify. First of all, the gitmap update all. It could be hyphen all or space all or UA. This would do the same thing. But now coming to the point how it's going to do it, the update is a delicate situation. So if we are running from the external machine, we should actually use JSON as a communication so that the help text does not get included in it. And currently, the way that you are doing the formatting or displaying it looks very poor. Okay? So I think this is where we need to work on. We need to make sure the help and things, these are really high quality. Okay? So currently, the output of the part two, how you're doing it, it's quite poor. Okay? And this is where we need to work on. And also make sure that, at the end, you bump the minor version and make a release and check the Gitmap PE. Also now I wanted you to check, so you can even run a test command on, let's say, scripts picture I or something like this. And you check by the running prompts LS that this project is actually running and other projects, but this project needs to be there. Or you could do other stuff like the own project. Let's say the status, alien status sample, just put a hi there and check if this comes up in the running prompt when it is running or put some weight. Okay? Like sleep and then say hi. You can do that so you can trace back. So that's one thing. Another is backing up the running prompts. So each one of the prompts from all the projects, the running prompts would be backed up. Then you can restore the prompt, running prompt, that should inject and run the prompt that we have taken the backup. So make sure that you do the end-to-end testing so that you can confirm it, that it's working very well and there is no confusion. Okay? It's been several times that we are trying that. I hope there should be no issues, right? If possible, you can also test out with some file the Gitmap SC or deploy command, sync left, sync right command. Okay. And also update the help text and also UI help text. These are very important. Do not miss it. Is it clear? Can you please do that? At the end, do a minor bump and a release and do the Gitmap PE to check the errors. Okay.
```

---

## 1. Domain Architecture & Problem Statement

### 1.1 Fleet Update (`update all`, `update-all`, `ua`)
- Currently, when running fleet binary update across nodes (`gitmap update all` or `gitmap ua`), remote SSH execution can leak raw PowerShell/bash stdout (such as `WARNING: [Test-RepoExists] The remote server returned an error: (404) Not Found`, probing logs, and installer text) into the terminal.
- Remote nodes must communicate back structured JSON (`--json`) when invoked by fleet orchestrator.
- The terminal output must be formatted with clean modern styling (node alias, IP, duration, status, installed version) instead of dumping unparsed installer noise.
- Aliases supported: `gitmap update all`, `gitmap update-all`, `gitmap update -all`, `gitmap ua`.

### 1.2 Running Prompts End-to-End Verification
- `gitmap agy running-prompts ls`: Discovers active and queued prompts across projects in Antigravity workspaces.
- `gitmap agy running-prompts backup` (and `gitmap backup-running-prompts`): Backs up active and queued prompts from all projects into SQLite Split-DB (`backup-prompts/sql.db`).
- `gitmap agy running-prompts restore` (and `gitmap restore-running-prompts`): Restores and re-injects prompts back into project prompt queues (`agy-prompt-queue.json`).
- Must be verified end-to-end with real/mock projects and trace delays.

### 1.3 Deploy & Sync Modes Verification
- `gitmap deploy <node> <source> <dest>` supporting `--sync`, `--sync-right`, `--sync-left`, `--overwrite`, and `--json`.
- Commands `gitmap deploy-right`, `gitmap deploy-left`, and `gitmap sc` (smart copy) verified with clean JSON output and conflict detection.

### 1.4 High-Quality Help Text & UI Help
- Author comprehensive, beautifully formatted help text in `cli/helptext/update.md`, `cli/helptext/running-prompts.md`, `cli/helptext/deploy.md`.
- Ensure terminal renderers (`update help`, `running-prompts help`, `deploy help`) match the cyan/yellow boxed terminal UI standards.

---

## 2. Actionable Deliverables

- `Task-01`: Fleet Update (`update all`, `update-all`, `ua`) with JSON Communication & Clean Terminal Rendering.
- `Task-02`: Running Prompts End-to-End Discovery, Backup, and Restore Injection.
- `Task-03`: Deploy & Sync Modes (`deploy`, `deploy-right`, `deploy-left`, `sc`) Testing.
- `Task-04`: High-Quality Help Text & UI Documentation.
- `Task-05`: Minor Version Bump (`v6.361.0`) and Live `gitmap pe` Verification.
