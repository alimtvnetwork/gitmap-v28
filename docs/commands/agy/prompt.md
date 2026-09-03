# `gitmap agy prompt` & `prompt-all-project`

Send user prompts into Antigravity workspace agent sessions.

## Subcommands

### 1. `gitmap agy prompt <project> <prompt>`
Send a prompt to a specific Antigravity project session.
```bash
gitmap agy prompt gitmap-v28 "Run linting and report status"
```

### 2. `gitmap agy prompt-all-project <title> <prompt>`
* **Alias:** `gitmap agy pap`
* Broadcast a prompt across all active Antigravity project workspaces.
```bash
gitmap agy pap "Repository Audit" "Check if repo has uncommitted changes"
```
