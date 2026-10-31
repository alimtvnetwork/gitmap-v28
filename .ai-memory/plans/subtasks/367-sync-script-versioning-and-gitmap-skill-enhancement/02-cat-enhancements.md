# Subtask 02: GitMap Cat Flag & Query Filtering Enhancements

- **Task Code:** `TASK-02`
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 02`
- **Status:** completed
- **Owner:** CLI & Macro Engineer
- **Specification:** [02-cat-enhancements-spec.md](../../../02-spec/21-app/367-sync-script-versioning-and-gitmap-skill-enhancement/02-cat-enhancements-spec.md)

## Actions
1. Add `-t` flag parsing in `cli/cmdmacro/cat_options.go` alongside `--tail` and `-tail`.
2. Update `cli/cmdmacro/cat_reader.go` so `opts.TailCount` and `opts.HeadCount` work seamlessly with `opts.Query`:
   - If query and tail are both specified (e.g. `gitmap cat file.jsonl -t 5 -q "media|USER_INPUT"`), extract matching lines and return the last `TailCount` matches.
3. Update `cli/cmdpy/py_cmd.go` to expand inline code warnings to flag `os.walk` and `time.sleep` anti-patterns.
4. Update `cli/helpdoc/cat.md` and `cli/helpdoc/py.md`.
