# Subtask 2: Terminal UI & Settings (URL and Templates)

1. Modify gitmap settings struct and migration logic to add:
   - CommitReplayKeepUrl (bool, default alse)
   - CommitReplayTemplates (map[string]string) to map generic messages to custom ones.
2. In the replay engine (where commits are copied/applied), strip the original commit URL unless CommitReplayKeepUrl is true.
3. In the replay engine, parse the commit title. If it matches a generic title (e.g., "Changes", "Lovable update", "Work in progress"), replace it using the CommitReplayTemplates map (or hardcoded defaults if map is empty).
4. Update the terminal renderer (e.g., [commit-right] [1/292] SHA title) to add left padding (e.g.,   [commit-right]), color code the SHA (cyan), the index (dim), and the title (white or colored based on type).
