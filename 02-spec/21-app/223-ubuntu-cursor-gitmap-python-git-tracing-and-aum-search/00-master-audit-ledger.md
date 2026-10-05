# Master Audit Ledger: 223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search

> **Task ID:** 223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search  
> **Status:** Phase 1B Spec Authoring Active  
> **Created:** 2026-10-05  
> **Protocol:** execute-parent-task-with-n-steps-v6 (A = 2, H = 2)  
> **Precedence:** Top-Instruction Priority Mandate (User preamble outranks all)

---

## 1. User Request (Verbatim)

```text
have you improved and added cursor to the ubuntu system using gitmap or python or shell script for now and keep track in the repo-secrets folder

try to test and make sure by testing properly that these works, find issues if there is any fix it and make a minor bump and check gitmap pe -t until ci cd is green
python codes needs to run from gitmap please ensure that

make sure we have these git commands from gitmap because it can also keep a trace to split db for heat map understanding?? Clear???

make sure the agent calls are optimized using gitmap https://prnt.sc/Rk0BXJBXgkeD
https://prnt.sc/C0T16jJtXhE_
https://prnt.sc/2k7l4_Q0VFyM

Make sure we have the se AUM searches verified properly in other or current repo against other search tools to see if we have any mistakes and can be improved anything??

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Improve and add cursor to the Ubuntu system using GitMap, Python, or shell script and keep track in the repo-secrets folder.
5. Test thoroughly to ensure functionality, find and fix any issues, make a minor version bump, and check `gitmap pe -t` until CI/CD is green.
6. Ensure Python code runs from GitMap.
7. Ensure Git commands from GitMap can trace to split db for heat map understanding.
8. Optimize agent calls using GitMap.
```

---

## 2. Requirement Matrix & Actionable Scope

| ID | Requirement | Technical Scope & Target Files | Acceptance Criteria |
|:---|:---|:---|:---|
| **REQ-01** | Cursor IDE on Ubuntu Fleet (`u1`) | `repo-secrets/05-scripts/setup-cursor-ubuntu.py`, `.sh`, `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`, `cli/cmdcursor/cursor_install.go` | Autonomously fetch download URL from Cursor API, provision on `u1`, install `--no-sandbox` wrapper at `~/.local/bin/cursor` and `/usr/local/bin/cursor`, verify headless execution, and record JSON status in `repo-secrets/`. |
| **REQ-02** | Python Execution via GitMap (`gitmap py`) | `cli/cmdpy/py_cmd.go`, `cli/cmd/rootutility.go`, `cli/cmd/rootcore.go` | Route `gitmap py` and `gitmap python` to execute inline Python code (`-c "<code>"`) and script files (`gitmap py <script.py> [args]`), capturing execution into `store.RecordAiExecution` and `CommandHistory`. |
| **REQ-03** | Git Command Split-DB Heatmap Tracing | `cli/cmd/rootgit.go`, `cli/cmd/commit_push.go`, `cli/store/command_history_split_db.go` | Ensure all `gitmap git <cmd>` commands log execution duration, timestamp, and exit code to Split-DB `commands.db` (`CommandHistory`) to fuel developer activity heatmap visualization. |
| **REQ-04** | AUM Search Quote Stripping & Agent Optimization | `cli/cmdautomation/search.go`, `cli/cmdautomation/automation_cmd.go` | Automatically strip escaped outer quotes (`\"...\"`, `"...""`) from `searchOpts.Pattern` in `gitmap aum search` so agent queries never return false zero matches; benchmark and verify against pattern searches. |
| **REQ-05** | Remote Verification, Minor Bump & Green CI/CD | `03-ai-scripts/37-bump-version.py`, `cli/constants/constants.go`, `gitmap pe -t` | Bump minor version (`v6.483.0`), run all repository linters, commit with hyphen format (`cursor - ...`), and verify pipeline passes with `gitmap pe -t`. |

---

## 3. Disjoint Multi-Agent Subtask Ownership

```
Wave 1 (Phase 1B - Spec Authoring):
├── Spec Author 01 (Subagent 1)
│   ├── 02-spec/21-app/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/01-architecture-spec.md
│   ├── .ai-memory/plans/subtasks/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/01-cursor-ubuntu-fleet-setup-and-repo-secrets.md
│   └── .ai-memory/plans/subtasks/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/02-gitmap-py-runner-and-ai-telemetry.md
└── Spec Author 02 (Subagent 2)
    ├── 02-spec/21-app/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/02-component-and-search-spec.md
    ├── .ai-memory/plans/subtasks/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/03-git-split-db-heatmap-command-tracing.md
    ├── .ai-memory/plans/subtasks/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/04-aum-search-quote-unquoting-and-agent-optimization.md
    └── .ai-memory/plans/subtasks/223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/05-remote-u1-testing-and-release-ceremony.md
```
