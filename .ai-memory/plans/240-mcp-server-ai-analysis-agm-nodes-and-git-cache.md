# Master Execution Plan: GitMap AI MCP Server, AI Analysis Engine, AGM Node Deployment & Git Cache

> **Plan Version:** 1.0.0  
> **Status:** Active  
> **Parent Task:** Task-240  
> **Associated Specs:**  
> - `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md`  
> - `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-component-and-cli-spec.md`  

---

## 1. Goal & Objectives

1. **AI MCP Server & AI Analysis Split DB Engine (`gitmap ai-analysis` / `gitmap ai`):** Transform GitMap into a primary Model Context Protocol (MCP) server for AI assistants. Provide task-based analysis tracking storing files inspected/modified, exact line ranges, and architectural reasoning in a dedicated Split SQLite database (`.gitmap/data/ai-analysis/ai-analysis.db`). Master table `AiTask`, file table `AiTaskFile`, and granular reasoning child table `AiTaskLine`.
2. **Task-Based File Removal with OS Temp Staging & Revert (`gitmap rm` / `gitmap file remove`):** Provide AI-safe file removal command that stages deleted files in an OS temp backup directory (`<temp>/gitmap/removed/<taskId>/`) allowing full task-level undo/revert before temp purge.
3. **AI Analysis LLM Training & Portability:** Summarization command (`gitmap ai-analysis llm-train`) to train/feed LLMs past architectural decisions with sequential next-step navigation, pruning/clearing (`gitmap ai-analysis clear`), and cross-system export/import (JSON, ZIP, SQLite).
4. **Directory Inventory & Grep Replacement (`gitmap child-path` / `gitmap ls-child`):** Replace slow PowerShell `Get-ChildItem` and raw `git grep` with high-speed GitMap native commands with clear documentation and migration examples.
5. **Meta Muse Cross-Platform Installer (`gitmap install muse` / `gitmap muse install`):** Automated installer supporting Windows (PowerShell/winget), Ubuntu/Linux (curl/apt), and macOS (brew/curl) based on `assets/screenshots/240-muse-meta-installer.png`.
6. **Smart Repository Creation & Clone Fallback (`gitmap repo create` / `gitmap cr`):** Detect if target repository already exists on remote/local and prompt user to clone instead of failing or blindly force-pushing.
7. **Antigravity Manager (AGM) Fleet Deployment & Account Purge:** Securely export, encrypt (AES-256-GCM), and deploy AGM configuration across cluster nodes, automatically purging sensitive JSON credential files after transmission.
8. **High-Speed Cached Git Logs & Two-Branch Comparison:** SQLite commit graph caching in `.gitmap/data/git-cache/<slug>/sql.db` for instant 0ms `gitmap log` viewing and structural two-branch comparison (`gitmap branch compare`).
9. **Compact Root Help Display:** Minimal default `gitmap` output showing version and footer notes only, reserving full command listings for explicit `gitmap help` / `-h`.

---

## 2. User Request (Verbatim)

```text
Now I wanted to learn what can we improve in the `gitmap`. What are the functionalities that we can create that `gitmap` can behave like a MCP server for the AI? And also in the future, AI does the analysis. AI does analysis, AI does searching, finding the files for a task. I want you to create in the `gitmap` these whole things, that analysis or anything that the AI does for the repo, and that would be AI analysis. That would be AI analysis part that would have for a task. It would start for a task, and then for that task, it will also put all these lines that it is reading, which files it is reading, which lines it is reading. So all these things, we should have AI analysis command that would use, again, the cache DB, a little bit of split DB data from the cache DB, and also itself. And again, it would use the split DB concept for the repo, so that would have a specific folder, AI analysis. So anytime that AI edits, modifies, or read based on the task, it would create that task, have a master table, and also the sub-table contains the file lines. So initially, it would have which file, related path, things like that in the master table. Very simple, straightforward information. With this, it would have child table where we have the lines we are modifying, and also we should have the reasoning in the table. So that means when AI reads a file, like the analysis, the reasoning for the task, it would also write into the reasoning section, like why it is reading or why it is modifying. Using `gitmap`, so `gitmap` have or needs to have these commands that would deal with the situation. And also, in cases like removing file path. So rather than removing an item directly, any AI should remove using `gitmap`. So `gitmap` should have similar command that actually accommodate the removal. And again, every removal is task-based, so that can be undo and removes can be kept as a backup in the OS temp directory for time being. But if the temporary is removed, then we cannot revert it. Remember that. There are a lot of AI scripts I do see, like the doc path linter. Also, reading multiple files, I believe, using Python files. So whatever the Python code that AI writes, that can also go through that analysis phase. The reason I'm saying this, that we could use that data to feed AI anytime, let's say, "Hey, learn this. What was the previous decision-making? Why it took the decision-making, and what was the reason?" We should be able to have a command like AI analysis LLM train, and that would basically give a summary with the next, next to any LLM, that it would learn each one of the task, why the decision was made and how it made. And also, it should have a clear method so that we can clear the decision-making or analysis, because over time it would become very big. These type of things. And that can be exported and imported from one system to another using JSON or using ZIP or also using the SQLite DB. Remember that. Now, sometimes we can do the filter, I believe, get child path. So this type of command, we need to have in `gitmap` with an example, not to use this type of get child, but use the `gitmap` functionality. Also, the git grep should be avoided and should provide the example with our example or our things. Remember that. Also, in the `gitmap`, try to install or have the installer feature for the Muse, Meta Muse. I'm giving you the screenshot so you can also search for the AI to understand how to install this in different platforms. So this is for Windows platform. There is a platform installation for Ubuntu, Linux and macOS, so try the different platform installation as well. Now, another important factor that `gitmap` can create a repository using the `gitmap` repo create. Now, if in case the repo is already created, then `gitmap` will say, "Do you want to clone it?" If the user says UI, then it's going to clone instead of create because it's already there. So we need to have this type of features. First, all the screenshots and things that I have given you, I want you to take the screenshots into the assets and try to define the spec what I have described. So my plan is to have all kinds of these functionality should be running from `gitmap`, and `gitmap` would be a prime MCT server for any AI to work with and very faster. Do you understand? So first, define the spec, show me the commands that you think should be created, and then you start working on it. Also, at the same time, please confirm from the `Antigravity Manager (`Antigravity Manager (AGM)`)` tool, Antigravity Manager, and from the map, we should be able to import, export, and deploy all the configuration across the nodes. That's the first thing. Second is that we should be able to deploy the accounts with the removal. That means once the accounts are deployed, it will remove those JSON files automatically. You need to confirm that this is how the implementation is. We can do it from `Antigravity Manager (`Antigravity Manager (AGM)`)` to other machines. `Antigravity Manager (`Antigravity Manager (AGM)`)` will use the `gitmap` to accelerate all the things and from the `gitmap`, we would be able to do this. And remember, the branch, we should be able to compare branches using the `gitmap`, `gitmap` branch compare, or we can say `gitmap` compare. So whatever the command that you could think of, but we should have this type of comparison feature from `gitmap` that will give you instant result using the SQLite. And that would be Git cached, remember, like we did for the Git log. So whenever we run the Git log, that usually takes longer time, but if it is Git log cache, then it will take instantly less than one second, less than a millisecond, to show that. And then it would have git log with a fast, or we could also do `gitmap` log. So whatever the command that you think of, you can do this. And another thing that `gitmap` root help should be very minimal. It would just say what is the version, and if there are any notes that we should provide, then it will provide notes in the footer. That's all. It should not show all the commands that it has. But if the user says `gitmap help`, only then it would show all the commands. Remember that.
```

---

## 3. Subtask Breakdown

- **Subtask 01 (`01-spec-and-plan.md`):** Canonical Architecture Spec, Component Spec, and Subtask Decomposition.
- **Subtask 02 (`02-compact-root-help.md`):** Compact Root Help & Footer Display (`gitmap` Default Output).
- **Subtask 03 (`03-ai-analysis-split-db.md`):** AI Analysis & Reasoning Split DB Engine (`gitmap ai-analysis`).
- **Subtask 04 (`04-task-file-removal.md`):** Task-Based Safe File Removal & Staging Undo (`gitmap rm`).
- **Subtask 05 (`05-llm-train-clear-transfer.md`):** AI Analysis LLM Training Summarizer & Data Portability.
- **Subtask 06 (`06-child-path-inventory.md`):** Directory Inventory & Search Replacement Commands (`gitmap child-path`).
- **Subtask 07 (`07-muse-installer.md`):** Meta Muse Multi-Platform Automated Installer (`gitmap install muse`).
- **Subtask 08 (`08-smart-repo-create.md`):** Smart Repository Creation & Interactive Clone Fallback (`gitmap repo create`).
- **Subtask 09 (`09-agm-deploy-purge.md`):** AGM Node Configuration Deployment & Sensitive Account Purge.
- **Subtask 10 (`10-git-cache-compare.md`):** High-Speed Cached Git Logs & Two-Branch Diff Comparison.
- **Subtask 11 (`11-verify-and-push.md`):** Quality Verification, File-Scoped Linting & Single Atomic Push.

---

## 4. Execution Tracking Ledger

| Subtask ID | File | Owner | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- |
| Subtask-01 | `01-spec-and-plan.md` | Lead Orchestrator | DONE | `01-architecture-spec.md`, `02-component-and-cli-spec.md`, master plan, and subtasks 01-11 authored and verified |
| Subtask-02 | `02-compact-root-help.md` | Worker 01 | PENDING | - |
| Subtask-03 | `03-ai-analysis-split-db.md` | Worker 02 | PENDING | - |
| Subtask-04 | `04-task-file-removal.md` | Worker 01 | PENDING | - |
| Subtask-05 | `05-llm-train-clear-transfer.md` | Worker 02 | PENDING | - |
| Subtask-06 | `06-child-path-inventory.md` | Worker 01 | PENDING | - |
| Subtask-07 | `07-muse-installer.md` | Worker 02 | PENDING | - |
| Subtask-08 | `08-smart-repo-create.md` | Worker 01 | PENDING | - |
| Subtask-09 | `09-agm-deploy-purge.md` | Worker 02 | PENDING | - |
| Subtask-10 | `10-git-cache-compare.md` | Worker 01 | PENDING | - |
| Subtask-11 | `11-verify-and-push.md` | Lead Orchestrator | PENDING | - |
