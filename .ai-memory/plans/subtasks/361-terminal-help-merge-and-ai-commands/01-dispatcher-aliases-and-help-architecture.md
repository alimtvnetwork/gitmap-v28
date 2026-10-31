# Subtask 01: Dispatcher Aliases & Help Architecture

- **Plan:** `361-terminal-help-merge-and-ai-commands`
- **Slug:** `gitmap-feature-implementation-and-update - Task 01`
- **Status:** `pending`

## Steps
1. Add `"thelp"` and `"th"` to `cli/cmd/rootutility.go` under `constants.CmdHelp`.
2. Add `"merge-in"` and `"mergein"` to `cli/cmd/rootcore.go` under `constants.CmdCommitIn`.
3. In `cli/cmd/rootutility_help.go`, update `normalizeHelpTopic` to map `"merge-in"`, `"mergein"`, `"thelp"`, `"th"`, and `"specific"` cleanly.
4. Support `gitmap thelp specific` and `gitmap help specific` with dedicated guidance.
