# Master Execution Plan: GitMap Runner, Audit DB, Supabase Vault, PE-AI & Muse Skills

> **Plan Version:** 1.0.0  
> **Status:** Active  
> **Parent Task:** Task-239  
> **Associated Specs:**  
> - `02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/01-architecture-spec.md`  
> - `02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/02-component-spec.md`  

---

## 1. Goal & Objectives

1. **Universal File Runner (`gitmap run <file>`):** Automatically detect file extensions (`.py`, `.ps1`, `.sh`, `.js`, `.ts`, `.go`) and execute with the appropriate interpreter. Maintain an audit log in SQLite task DB, and record errors in a dedicated run errors DB. Expose `gitmap run errors` / `gitmap run history` to easily discover and fix failing files.
2. **Multi-Supabase Database Support with Encrypted Secrets:** Enable multiple Supabase databases in GitMap via one-liner CLI commands (`gitmap supabase add <alias> <url> <anon_key> <service_key> [db_url]`). Enforce AES-GCM / RSA vault encryption so no password or service key is ever stored in cleartext.
3. **Chore Commits (`cpc`) & Help Text Full Forms:** Add `gitmap cpc "<module> - <summary>"` (Commit & Push Chore) alongside `cpf`, `cpb`, and `cpr`. Expand CLI help text to document full forms (`cpf`: commit-push-feature, `cpb`: commit-push-bug, `cpc`: commit-push-chore, `cpr`: commit-push-release, `pas`: pull-all-ssh, etc.).
4. **Pipeline PE-T Enhancements, `--ai` Flag & Source Details:**
   - Add `gitmap pe -t`: Poll every 2 minutes; if errors are detected in the logs/stack trace, immediately showcase via `gitmap pe` and exit.
   - Add `gitmap pe -ud` (`--until-done`): Run continuously until full CI/CD completes, then display final status.
   - Add `gitmap pe --ai`: Suppress copying logs to clipboard/terminal.
   - Enhance PE details with log source (logger name, workflow name, URL, log path).
   - Immediate Git hash retrieval: If commit SHA is already in database and completed, pull immediately from DB without waiting.
5. **Muse Master Prompt & GitMap Skills Modernization:** Update Muse master prompt (`01-prompts/27-muse-prompts/01-muse-master-prompt.md`) and skills (`.agents/skills/muse-master-prompt/`, `.agents/skills/gitmap/`) to integrate native GitMap commands (`gitmap aum search`, `gitmap task`, `gitmap run`, `gitmap pe --ai`, `gitmap sync`). Update coding guidelines to mandate `gitmap pe --ai`.

---

## 2. User Request (Verbatim)

```text
Hi there. Can you please help me with the discover repositories. You have committed all the skills and updates to all the projects, I believe. I want you to update the Muse mastery prompt a little bit more and try to include in most parts the `gitmap` skills and also check that in the coding guideline, have we improved the `gitmap` skills? And can we add more features from `gitmap` in the skills section? You have to confirm me. For the `gitmap`, I think what you have is so far good. I have checked. Also `gitmap`, if in any case that it is trying to run the Python or let's say `gitmap`, could be a run command and a space, we could do Python file, PowerShell, Bash, any file. So if that is something that we are running, then `gitmap` will automatically detect the extension and then try to run with that. And also it will keep an audit log. That means it will enqueue the task to the task DB and then do it. When completed, it will mark that this is done. If there is an error, then it would save that error to the errors DB, and we can get the run errors. Then we can get all the files that has been run and has errors we can detect and fix. Also at the same time in the `gitmap`, we wanted to include the multiple, let's say, Supabase databases, and there should be one liner with help as well to add it and to communicate with multiple machines. We can use the Supabase database, but make sure that we do not keep any password or anything directly into the Supabase. If we have to keep it, it has to be encrypted. Remember that. And also in the `gitmap`, I think we need to have a little bit of chore commits which is missing. We have featured commits, we have bug fix commit, release commit, pull commits as well. So also write the pull forms of these methods in the help. And also we need to improve the PE-T. So PE-T, whenever it finds the error, it will try to showcase the errors. So `gitmap` space PE-T, it would have two versions. So hyphen T would wait for at least every two minutes to check if there is an error. If it finds an error that is there in the stack trace found in the logs, then it would immediately try to showcase as `gitmap` PE and exit out. And there would be another one, hyphen T, hyphen until done. So UD, hyphen UD, until done. That would actually run infinitely until the full CI/CD completes, and then it would run the `gitmap` PE and gives the error. Now, in the `gitmap` PE, it usually copies the logs to the memory, which is not helpful for the terminal and AI. So make sure that we have another thing called `gitmap` PE space hyphen hyphen AI. Now, in this case, it would not copy the text in the terminal, and all the AI should be using `gitmap` D hyphen hyphen AI. And also it needs to be updated in the skill. So even though the user says run `gitmap` PE as an AI, it should always run `gitmap` PE hyphen hyphen AI so that it does not copy to the terminal because we don't need to. AI will have the traces and then try to fix it. And try to enhance the PE information, like where the log is coming from, which stack that is the logger name or the pipeline name. I hope that these are there, but also if you can improve, do improve. And make sure that it tries to find out the last Git hash. So if the hash is already saved to the databases, it would immediately pull without waiting. So remember that. So all these things I asked you to do, try to make a list of it, write the spec first. You cannot write the spec for the `gitmap` in the coding guideline, but yes, I think you should write it in the `gitmap` first and then update the skills here. And then also improve the Muse skills. Do you understand?
```

---

## 3. Subtask Breakdown

- **Subtask 01:** `01-universal-file-runner-and-task-audit.md`
  - Implement `cli/cmdrun/` universal file runner supporting `.py`, `.ps1`, `.sh`, `.js`, `.ts`, `.go`.
  - Auto-extension resolution and interpreter discovery.
  - Enqueue task to SQLite task DB before running, mark complete upon exit 0.
  - On failure, record run error into SQLite errors DB.
  - Implement `gitmap run errors` / `gitmap run-errors` / `gitmap run history` to list failed runs.
  - Wire into `cli/cmd/roottooling.go` and `cli/cmd/macro_root_dispatch.go`.

- **Subtask 02:** `02-multi-supabase-vault-and-encryption.md`
  - Implement `cli/cmdsupabase/` package for multi-project Supabase management.
  - One-liner registration: `gitmap supabase add <alias> <url> <anon_key> <service_key> [db_url]`.
  - AES-GCM / RSA vault encryption for all secrets, service keys, and database passwords.
  - Storage in `.gitmap/data/installation/supabase/sql.db`.
  - Wire into `cli/cmd/rootdata.go` and `cli/cmd/roottooling.go`.

- **Subtask 03:** `03-chore-commits-and-full-forms-help.md`
  - Add `CmdCommitPushChore = "commit-push-chore"` and `CmdCommitPushChoreAlias = "cpc"`.
  - Implement `runCommitPushChore` in `cli/cmd/commit_push.go` prepending `Chore: `.
  - Wire into `cli/cmd/roottooling.go`, `cli/cmd/nodes_cmd.go`, and `cli/cmd/sends_cmd.go`.
  - Update CLI help menus (`cli/cmd/commit_help_menu.go`, `cli/helptext/`) documenting full forms of all commit methods (`cpf`, `cpb`, `cpc`, `cpr`, `pcp`, `pas`).

- **Subtask 04:** `04-pipeline-pe-ai-until-done-and-source-details.md`
  - Add `--ai` flag to `cli/cmdpipeline/pipeline_flags.go` to suppress clipboard copying and clipboard notices.
  - Implement 2-minute interval checking in `gitmap pe -t`, immediately showcasing errors and exiting when detected in logs/stack traces.
  - Add `-ud` / `--until-done` flag running continuously until the pipeline finishes.
  - Enhance failure diagnostics with log source details (logger name, workflow name, URL, log path).
  - Immediate Git hash retrieval when SHA is already completed in database.

- **Subtask 05:** `05-muse-master-prompt-and-gitmap-skills-modernization.md`
  - Update `01-prompts/27-muse-prompts/01-muse-master-prompt.md` with modern GitMap primacy, `gitmap task`, `gitmap run`, `gitmap pe --ai`, and `gitmap sync`.
  - Update `.agents/skills/muse-master-prompt/skill.md` and `.cursor/skills/muse-master-prompt/skill.md`.
  - Update `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md` with new features (`gitmap run`, `gitmap pe --ai`, `gitmap pe -ud`, `gitmap cpc`, `gitmap supabase`).
  - Update `02-spec/02-coding-guidelines/06-ai-optimization/` to mandate `gitmap pe --ai` and GitMap tool primacy.

---

## 4. Execution Tracking Ledger

| Subtask ID | File | Owner | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| Subtask-01 | `.ai-memory/plans/subtasks/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/01-universal-file-runner-and-task-audit.md` | Worker 01 | PENDING | Pending dispatch |
| Subtask-02 | `.ai-memory/plans/subtasks/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/02-multi-supabase-vault-and-encryption.md` | Worker 02 | PENDING | Pending dispatch |
| Subtask-03 | `.ai-memory/plans/subtasks/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/03-chore-commits-and-full-forms-help.md` | Worker 01 | PENDING | Pending dispatch |
| Subtask-04 | `.ai-memory/plans/subtasks/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/04-pipeline-pe-ai-until-done-and-source-details.md` | Worker 02 | PENDING | Pending dispatch |
| Subtask-05 | `.ai-memory/plans/subtasks/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/05-muse-master-prompt-and-gitmap-skills-modernization.md` | Lead | PENDING | Pending dispatch |
