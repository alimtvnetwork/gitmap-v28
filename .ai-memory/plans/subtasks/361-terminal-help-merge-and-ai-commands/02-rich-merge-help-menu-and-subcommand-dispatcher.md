# Subtask 02: Rich Merge Help Menu & Subcommand Dispatcher

- **Plan:** `361-terminal-help-merge-and-ai-commands`
- **Slug:** `gitmap-feature-implementation-and-update - Task 02`
- **Status:** `pending`

## Steps
1. Create `cli/cmdmerge/merge_help_menu.go` with `RenderMergeHelp()`.
2. Format as a 2-column menu covering `merge-right`, `merge-left`, `merge-both`, `merge-in`, `commit-right`, `commit-left`, `commit-both`, and `merge write`.
3. Update `cli/cmdmerge/dispatchmovemerge.go` so `command == "merge"` routes subcommands (`right`, `left`, `both`, `in`) and displays `RenderMergeHelp()` on help or bare invocation.
4. Update `cli/cmd/rich_help_dispatcher.go` to route `"merge"`, `"merge-right"`, `"merge-left"`, `"merge-both"`, `"merge-in"` to `cmdmerge.RenderMergeHelp()`.
