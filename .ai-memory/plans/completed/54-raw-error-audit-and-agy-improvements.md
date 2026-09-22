# Master Spec: Raw Error Audit & Elimination of Swallowed Errors with AGY Workspace Enhancements

## Outcome Summary
- **Status**: COMPLETED
- **Completed Date**: 2026-09-22
- **Master Plan ID**: 54

## User Request (Verbatim)

```text
Please go through all the code base wherever you are managing or handling error. Make sure the error also contains the raw error. So if there is a raw error, there should be no swallow. Okay? And you should make sure of that everywhere in the code. And I want you to go through it deep and you find those. Okay? Do you understand my concept? Can you please follow through?

https://prnt.sc/zAm6uanENZGF

gitmap agy rm-rejoin-read(rrr) <id>/<seq>/project-slug-starts-with,<id> # it will rename the conversation with project name or slug , clear??
gitmap agy rm  <id>/<seq>/project-slug-starts-with,<id>
gimay agy rm help
gimay agy rm-rejoin-read help
gitmap agy pins # shouw show ls 
gitmap agy pins help/ls/add/rm/ <id>/<seq>/project-slug-starts-with,<id>
gitmap agy rm-rejoin-pin-read(rrpr) <id>/<seq>/project-slug-starts-with,<id> # additionall running after rm-rejoin-read it pins if not 
already

So, git map AGY pins should show the AGY pin items first. So it should have add, remove, edit, help on this command, and the help should display nicely what it does, how the pins work. So pins LS or just pins will show the list of pin projects. These are there. And also there is a concern in the remove section. So in the remove section, it shows remove the ID. Rather than giving the ID, we should be giving also... ID should work. Also what would work is the sequence number. So when the table is created, it should show a sequence number in the left-hand side. Okay, let's say three digits sequence. So if the sequence is also given the number, three-digit format, then also the project would be removed. So we could do remove using sequence. We could do the start of the name. Okay. Along with that, we can actually remove also, let's say, folder by folder projects as well. We can do RM folder remove, and we can just select a folder root, and that would remove all these, let's say, nested projects along with that. So that is another command, like RM or folder RM. Okay. Actually, the AGY commands remove, it should not remove the actual project from the files and disk. So that needs to be mentioned clearly. It says, "Files on disk preserved." Okay, that's correct, but it needs to be a bit more detailed if I do a RM.help. Okay. Or remove as well, remove or delete. So all these commands needs to be there. Same behavior, just aliasing. Also, one more command I want to have, like git map AGY RM rejoin read. RM-rejoin-read. What it would do is that it would remove that project and conversation fully from the Gemini folder as well, very clearly, but also it keeps the project file. Do not remove the git. That is done. Also, every one of these steps needs to be done with the task entering so that we can do undo as well. So after removal, it will also say how user can undo or redo that. Remember that. Also, the history command needs to be shown on every removal, not only AGY, but all the AGY command. Remember that. Now, with this command, what I'm saying, like RM rejoin read. Okay? So that would remove the project from the Antigravity IDE, remove its conversation from the Gemini folder in the actual file system, conversation and memory. So it would do some cleanup for that project, and then it would add the project one more time, and then it would run the read prompt. So we can find the read prompt and feed that read prompt to that project and run it. So that's what it does. Also along with that, we should have another command that actually does the git map RM rejoin read. So that's like the... We could do ID sequence. We could do project slug. So all of these should work. It's not only for this case, but also RM and so all these cases. Okay, remember to do that. So it should also give the suggestions. If I do a tab, it would give me suggestion. Also in RM, if I do a tab, it would give me the suggestions. Remember that. That's very important. And I could do for multiple projects, actually. So another ID and things like that. So the same thing I could do for the other cases, AGY RM, the same way that the other cases we could do. Also, we could do git map AGY RM help. The same way we could do for the other commands help. Also the pins help. Okay. We have AGY pin. So, okay. We have AGY. We have help, LS, add, remove, add, remove, ETC. Also, it can be removed with all of these things, like the ID sequence or the other stuff. Okay, so these are on top of my head. So rejoin read should actually run the read prompt on top of after adding it. So we should have also another command that would pin the stuff. Rejoin. That is like rejoin. And also this one also, every read conversation, if we do it for the first time, it will rename the conversation with the project name or slug here. So it's just the additional stuff that would have a short version, like RRR. This will have RRBR version. So I've added all these things. You should first read, understand, and then-

gitmap agy ls
should show these 3 things, conv name, id, project name , please in the table display
```

## Visual Specification Reference

User reference screenshot saved to local persistent storage:
- `assets/screenshots/agy-rm-pins-raw-error-01.png`

## Key Accomplishments

1. **Elimination of Swallowed Raw Errors (Task-01)**:
   - Audited repository error wrapping sites across `cli/cmdagy`, `cli/cmd`, `cli/cmdssh`, `cli/cmdvhost`.
   - Replaced dropped `err` returns with `apperror.WrapSimple(err, op)`, `apperror.Wrap(err, op, ctx)`, `apperror.WrapNotFound`, `apperror.WrapValidation`, and `apperror.WrapExecution`.
   - Guaranteed every `AppError` preserves `err` in `Cause`.
   - Verified 3,601 source files with `python linter-scripts/check-error-management.py` (0 errors, 0 panics, 0 swallowed errors).

2. **`gitmap agy ls` Table Enhancement (Conversation Name, ID, Project Name)**:
   - Added `AgyLatestConv` loading from `~/.gemini/antigravity/conversation_summaries.db` in `cli/cmdagy/agy_conv_latest.go`.
   - Table columns rendered: `SEQ`, `PROJECT`, `ID`, `CONV NAME`, `CONV ID`, `BRANCH`, `STATUS`, `PATH`.
   - Displayed both conversation name (title) and conversation ID alongside project name and project ID with proper column truncation and alignment.

3. **Extended `gitmap agy rm` Multi-Target & Folder Resolution (Task-03)**:
   - Implemented 3-digit sequence (`001`), short/full project ID, slug prefix, comma-separated lists (`001,gitmap,418bd745`), and folder batch removal (`--folder <dir>` / `rm folder <dir>`).
   - Created dedicated, styled help guide in `cli/cmdagy/agy_rm_help.go` explicitly declaring: "IMPORTANT: Project files, source repositories, and git history on disk are strictly PRESERVED."

4. **Complete `gitmap agy pins` Suite (Task-04)**:
   - Defaulted `gitmap agy pins` to `ls`, displaying pinned projects first.
   - Subcommands supported: `add`, `rm`, `edit`, `help`, `ls`.
   - Added rich help rendering in `agy_pins_help.go` and empty state command hints in `agy_pin_projects_table.go`.
   - Wired root dispatching in `cli/cmd/root.go` for `gitmap pins`.

5. **`gitmap agy rm-rejoin-read` (`rrr`) and `rm-rejoin-pin-read` (`rrpr`) (Task-05, Task-06)**:
   - Purges stale conversation records from `conversation_summaries.db` and brain artifacts.
   - Preserves source repositories and git history on disk.
   - Re-enrolls the project and dispatches the enhanced Read Memory protocol.
   - Automatically updates initial conversation title to the project name or slug.
   - For `rrpr`: automatically pins the project into bookmarks upon rejoining.
   - Wired root dispatching for `gitmap rrr`, `gitmap rrpr`, `gitmap rrbr`.

6. **Task History Tracking & Undo Guidance (Task-07)**:
   - Logged snapshots before mutations to `TaskHistory`.
   - Displayed clear undo guidance (`gitmap agy undo`, `gitmap agy history`) after every project removal and rejoin.
   - Provided shell tab completion (`completeAgyRmArgs`) for sequence numbers, IDs, and slugs.

## Quality Linters Verified

- `python linter-scripts/check-file-sizes.py`: All 13 modified/created files in `cli/cmdagy` strictly $\le 100$ lines.
- `python linter-scripts/check-function-lengths.py --root cli/cmdagy`: All functions strictly $\le 15$ lines.
- `python linter-scripts/check-nested-ifs.py`: 0 violations across all 23 changed files (PASS).
- `python linter-scripts/check-boolean-guidelines.py`: 0 violations across all 23 changed files (PASS).
- `python linter-scripts/check-error-management.py`: 0 violations across 3,601 files (PASS).
