# Plan 182: Antigravity Table Grouping, Tasks Suite, Author/Sponsor & Install Footer Deduplication

> **Task Origin & Objective**:
> - User prompt requested:
>   1. `gitmap agy ls` table redesign: reduce text with fixed length and middle ellipsis (`...`) for project names and paths; group projects by common parent root folder so they display folder by folder; output duplicate projects in a dedicated section outside the table.
>   2. Explicitly confirm and ensure that `gitmap agy rm` / `agy rm` deletes the project configuration ONLY from `~/.gemini/config/projects/*.json` and NEVER removes files from the filesystem.
>   3. Align `gitmap agy` help: reduce verbose description for `all-projects-read-memory-prompt`, make help colorful (ANSI colors), and ensure clean alignment.
>   4. Add shell auto-completion script support across shells (`powershell`, `bash`, `zsh`, `fish`) for all commands (`gitmap completion <shell>`).
>   5. `gitmap tasks` command suite: wire `gitmap tasks` (and alias `task`) with subcommands `list`, `history`, `undo`, `redo`, and ensure all gitmap commands register their operation into `pending_tasks` prior to execution and mark complete upon success.
>   6. Add top-level `author` and `sponsor` (and `credits`) commands revealing author info (MD ALIM UL KARIM) and sponsor info (RISE UP ASIA LLC).
>   7. Deduplicate installer/update binary identity footer between `install.ps1` and `cmdupdate`.

---

## 1. Problem Analysis & RCA
- `cli/cmdagy/agy_ls.go` and `agy_ls_table.go`: rendered a flat table without root folder grouping, and without length caps on names or paths. Duplicates were not reported unless `gitmap agy find-duplicates` was called explicitly.
- `cli/cmdagy/agy_projects.go`: `deleteProjectFile` only calls `os.Remove(filePath)` on the config JSON (preserving the filesystem), but printed no confirmation message, creating user uncertainty.
- `cli/cmdagy/agy_read_memory_prompt.go`: overly long `Short` string caused Cobra's default help formatter to stretch command column widths to 31+ characters without color.
- `cli/completion/`: only supported PowerShell, Bash, and Zsh. Fish shell completion was missing.
- `cli/cmd/roottooling.go`: registered `task` and `tk` but lacked `tasks`. No dedicated tasks controller with `list`, `history`, `undo`, `redo`. Commands were audited into `command_history` but did not automatically register in `pending_tasks`.
- `cli/cmd/author_sponsor.go`: missing top-level commands to recognize author MD ALIM UL KARIM and sponsor RISE UP ASIA LLC.
- `install.ps1` / `cmdupdate`: `install.ps1` ran `& $binPath binary`, and `finishRemoteUpdate` also ran `printPostUpdateIdentity()`, causing the binary identity box to print twice during updates.

---

## 2. Subtasks
- **Subtask 01**: `cli/cmdagy/agy_ls_table.go`, `cli/cmdagy/agy_ls.go`, `cli/cmdagy/agy_ls_table_test.go` — Root Directory Grouping, Middle Ellipsis Truncation & External Duplicate Reporting.
- **Subtask 02**: `cli/cmdagy/agy_projects.go`, `cli/cmdagy/agy_read_memory_prompt.go`, `cli/cmdagy/agy_cmd.go`, `cli/constants/constants_completion.go`, `cli/completion/fish.go`, `cli/completion/completion.go` — Deletion Safety Guarantee, Colorful Agy Help & Fish Shell Auto-Completion.
- **Subtask 03**: `cli/constants/constants_task.go`, `cli/store/tasktype.go`, `cli/cmd/tasks.go`, `cli/cmd/audit.go`, `cli/cmd/roottooling.go`, `cli/cmd/rootsuggest.go`, `cli/cmd/tasks_test.go` — Tasks Command Suite & Universal Pending Task Queueing.
- **Subtask 04**: `cli/cmd/author_sponsor.go`, `cli/helptext/catalog.go`, `cli/cmdupdate/updateremoteinstall.go`, `install.ps1`, `cli/scripts/install.ps1` — Author/Sponsor Terminal Cards & Installer Footer Deduplication.

---

## 3. Target Release
- Minor release: `v6.250.0`.
