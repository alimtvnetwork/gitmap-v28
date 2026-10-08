# Master Plan: Task 240 — GitMap MCP AI Analysis, Task-Based Safe Removal, Meta Muse Installer, and AGM Fleet Sync

## Executive Overview
This plan establishes GitMap as the primary AI-native Model Context Protocol (MCP) companion and execution engine across multi-repository workspaces. It implements comprehensive AI task analysis and reasoning tracking in Split-DB SQLite, task-based safe file removal with OS temp directory backup and undo capabilities, LLM decision-making summarization (`gitmap ai train`), cross-platform Meta Muse installation (`gitmap muse install`), smart repository creation with automatic clone fallback, high-speed cached Git logs and branch comparison, AGM fleet account deployment with secure auto-removal of temporary backup JSONs, and a clean, concise default help menu.

---

## User Request (Verbatim)
```text
https://prnt.sc/iBoFN5X62bs7
https://prnt.sc/Zc2_K0Fr3sWW

# High Priority Instruction

Now I wanted to learn what can we improve in the `gitmap`. What are the functionalities that we can create that `gitmap` can behave like a MCP server for the AI? And also in the future, AI does the analysis. AI does analysis, AI does searching, finding the files for a task. I want you to create in the `gitmap` these whole things, that analysis or anything that the AI does for the repo, and that would be AI analysis. That would be AI analysis part that would have for a task. It would start for a task, and then for that task, it will also put all these lines that it is reading, which files it is reading, which lines it is reading. So all these things, we should have AI analysis command that would use, again, the cache DB, a little bit of split DB data from the cache DB, and also itself. And again, it would use the split DB concept for the repo, so that would have a specific folder, AI analysis. So anytime that AI edits, modifies, or read based on the task, it would create that task, have a master table, and also the sub-table contains the file lines. So initially, it would have which file, related path, things like that in the master table. Very simple, straightforward information. With this, it would have child table where we have the lines we are modifying, and also we should have the reasoning in the table. So that means when AI reads a file, like the analysis, the reasoning for the task, it would also write into the reasoning section, like why it is reading or why it is modifying. Using `gitmap`, so `gitmap` have or needs to have these commands that would deal with the situation. And also, in cases like removing file path. So rather than removing an item directly, any AI should remove using `gitmap`. So `gitmap` should have similar command that actually accommodate the removal. And again, every removal is task-based, so that can be undo and removes can be kept as a backup in the OS temp directory for time being. But if the temporary is removed, then we cannot revert it. Remember that. There are a lot of AI scripts I do see, like the doc path linter. Also, reading multiple files, I believe, using Python files. So whatever the Python code that AI writes, that can also go through that analysis phase. The reason I'm saying this, that we could use that data to feed AI anytime, let's say, "Hey, learn this. What was the previous decision-making? Why it took the decision-making, and what was the reason?" We should be able to have a command like AI analysis LLM train, and that would basically give a summary with the next, next to any LLM, that it would learn each one of the task, why the decision was made and how it made. And also, it should have a clear method so that we can clear the decision-making or analysis, because over time it would become very big. These type of things. And that can be exported and imported from one system to another using JSON or using ZIP or also using the SQLite DB. Remember that. Now, sometimes we can do the filter, I believe, get child path. So this type of command, we need to have in `gitmap` with an example, not to use this type of get child, but use the `gitmap` functionality. Also, the git grep should be avoided and should provide the example with our example or our things. Remember that. Also, in the `gitmap`, try to install or have the installer feature for the Muse, Meta Muse. I'm giving you the screenshot so you can also search for the AI to understand how to install this in different platforms. So this is for Windows platform. There is a platform installation for Ubuntu, Linux and macOS, so try the different platform installation as well. Now, another important factor that `gitmap` can create a repository using the `gitmap` repo create. Now, if in case the repo is already created, then `gitmap` will say, "Do you want to clone it?" If the user says UI, then it's going to clone instead of create because it's already there. So we need to have this type of features. First, all the screenshots and things that I have given you, I want you to take the screenshots into the assets and try to define the spec what I have described. So my plan is to have all kinds of these functionality should be running from `gitmap`, and `gitmap` would be a prime MCT server for any AI to work with and very faster. Do you understand? So first, define the spec, show me the commands that you think should be created, and then you start working on it. Also, at the same time, please confirm from the `Antigravity Manager (`Antigravity Manager (AGM)`)` tool, Antigravity Manager, and from the map, we should be able to import, export, and deploy all the configuration across the nodes. That's the first thing. Second is that we should be able to deploy the accounts with the removal. That means once the accounts are deployed, it will remove those JSON files automatically. You need to confirm that this is how the implementation is. We can do it from `Antigravity Manager (`Antigravity Manager (AGM)`)` to other machines. `Antigravity Manager (`Antigravity Manager (AGM)`)` will use the `gitmap` to access to the machines and then use its functionality to deploy the JSON configurations or account information. Once it is deployed, the files will be removed, and during the transaction, make sure the files are encrypted so that the other parties cannot see it. Remember that. It needs to be kept because there are credentials that needs to be encrypted. And we can use JWT method if that is helpful during the transition. And we can have, let's say, REST endpoints from `gitmap` and the `Antigravity Manager (`Antigravity Manager (AGM)`)` both. So both should have the simultaneous method to deploy the settings from one to the another. So this is something that I want. Also, I want you to test that, so that you can confirm that it is a working copy. And also for the Git logs, I think time to time, AI uses the Git logs and other methods. I want those to be coming from the `gitmap` as well. So you should have methods which actually deals with the situation. Older version needs to be better, faster with caching, so that you can always show the things faster if we already have run this before. Remember that giving the data from the caching SQL DB concept, it will be so much faster than directly looking into the Git. But also we can reuse the Git if necessary. And also we should have two branch comparison function easily from `gitmap`. There should be examples for it. And also, if we do the `gitmap` now on, do not put all the commands in the help, just the version and the footer notes is fine. If the user gives help, then they will see the full help. Do you understand all these features and tasks? If you have any question, concern, let me know.
```

---

## Visual Assets & Screenshot Register
The 5 screenshots provided by the user have been captured and preserved into `assets/screenshots/`:
1. `assets/screenshots/240-ai-rm-cleanup-patch.png`: Agents running `Remove-Item patch_*.py -Force` instead of GitMap native safe removal.
2. `assets/screenshots/240-doc-path-linter-error.png`: AI scripts encountering dead backtick path references in documentation.
3. `assets/screenshots/240-get-childitem-filter.png`: Agents running `Get-ChildItem -Path 03-ai-scripts -Filter "*reconcile*"` instead of `gitmap find`.
4. `assets/screenshots/240-git-grep-func-runfix.png`: Agents running `git grep "func runFix" cli/` instead of `gitmap aum search`.
5. `assets/screenshots/240-muse-meta-installer.png`: Meta Muse agent installation on Windows (`irm https://dev.meta.ai/install.ps1 | iex`).

---

## Subtasks Decomposition & Work Wave Allocation

| Subtask ID | Title | Owned Files | Status | Description |
| :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | AI Analysis Subsystem & Split-DB Schema | `cli/store/ai_analysis_split_db.go`, `cli/cmdai/ai_analysis_*.go`, `cli/store/split_db_path.go` | COMPLETED | Track AI file operations (reads, edits, deletes) and reasoning in Split-DB SQLite |
| **Subtask-02** | Task-Based Safe File Removal & Undo | `cli/cmdrm/rm_*.go`, `cli/cmd/rm.go` | COMPLETED | Task-based safe deletions backed up to OS temp with undo capability (`gitmap rm --undo`) |
| **Subtask-03** | AI Analysis LLM Train, Export & Cleanup | `cli/cmdai/ai_train.go`, `cli/cmdai/ai_export.go`, `cli/cmdai/ai_clear.go` | COMPLETED | Summarize reasoning history for LLMs, export/import (JSON/ZIP/DB), and prune records |
| **Subtask-04** | Meta Muse Multi-Platform Native Installer | `cli/cmdinstall/install_muse.go`, `cli/cmd/muse_cmd.go`, `cli/constants/constants_install.go` | COMPLETED | Native installer for Meta Muse across Windows, Linux/Ubuntu, and macOS |
| **Subtask-05** | Smart Repo Creation & Clone Fallback | `cli/cmd/repo_create_smart.go`, `cli/cmd/create_cmd.go` | COMPLETED | Check local/remote existence and offer automatic clone fallback |
| **Subtask-06** | AGM Fleet Sync, Transit Encryption & Auto-Removal | `cli/cmdnodes/nodes_deploy_agm.go`, `../Antigravity-Manager` | COMPLETED | Direct transit streaming, encrypted payload, and auto-removal of unencrypted backup JSONs |
| **Subtask-07** | High-Speed Cached Git Logs & Branch Comparison | `cli/cmdlog/log_*.go`, `cli/cmdgit/branch_compare.go`, `cli/cmd/branch.go` | COMPLETED | Split-DB cached Git logs (`gitmap log`) and fast branch comparison (`gitmap diff-branch`) |
| **Subtask-08** | Concise Default Help Menu & Anti-Pattern Docs | `cli/cmd/rootusage.go`, `cli/cmd/root.go`, `cli/helptext/`, `02-spec/02-coding-guidelines/` | COMPLETED | Compact bare `gitmap` output, full menu on `gitmap help`, anti-pattern guidelines |

---

## Multi-Agent Execution Strategy (A = 2, H = 2)

- **Wave 1 (Subtasks 01, 02, 03, 04):**
  - **Worker 01:** Owns Subtask-01 (AI Analysis Subsystem) & Subtask-03 (LLM Train & Export).
  - **Worker 02:** Owns Subtask-02 (Safe File Removal & Undo) & Subtask-04 (Meta Muse Multi-Platform Installer).
- **Wave 2 (Subtasks 05, 06, 07, 08):**
  - **Worker 01:** Owns Subtask-05 (Smart Repo Create) & Subtask-07 (Cached Git Logs & Branch Comparison).
  - **Worker 02:** Owns Subtask-06 (AGM Fleet Auto-Removal & Transit Security) & Subtask-08 (Concise Help Menu & Anti-Pattern Docs).

---

## Non-Negotiable Governance
1. Search codebase exclusively via GitMap (`gitmap aum search`, `gitmap find`). TOTAL BAN on `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`.
2. Strictly relative Git paths only (`cli/...`, `02-spec/...`, `.ai-memory/...`).
3. Clean GitMap chore commit format: `gitmap cpc "<module> - <summary>"`. No colons inside message arguments.
4. All quality gates must pass with zero violations before release: `check-nested-ifs.py`, `check-enum-and-boolean.py`, `check-boolean-guidelines.py`, `check-error-management.py`, `check-relative-paths.py`, `check-legacy-refs.py`.
