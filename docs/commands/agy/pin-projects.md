# gitmap agy pin-projects

Manage pinned / starred Antigravity projects for quick filtering, inspection, and batch workflows.

```bash
gitmap agy pin-projects <subcommand> [flags]
```

Aliases: `pin-project`, `pinned-projects`, `pinned`, `pins`

---

## Subcommands

### 1. `ls` (List Pinned Projects)

List all pinned projects in a formatted terminal table or export as JSON.

```bash
# Formatted table view
gitmap agy pin-projects ls

# Short alias
gitmap agy pin-projects list

# JSON output for scripting
gitmap agy pin-projects ls --json
```

Output includes project name, short ID, default branch, status (`✔ active` or `✖ missing`), relative pinned timestamp, and filesystem path.

### 2. `add` (Pin Projects)

Pin one or more projects by ID, fuzzy project name, or filesystem directory path.

```bash
# Pin current repository/workspace
gitmap agy pin-projects add .

# Pin by name
gitmap agy pin-projects add gitmap-v28

# Pin by project ID
gitmap agy pin-projects add 46d05021-30e1-4036-acfb-3020489125eb

# Pin multiple projects simultaneously
gitmap agy pin-projects add proj-1 proj-2 /path/to/repo
```

### 3. `rm` (Unpin Projects)

Remove one or more projects from the pinned list.

```bash
# Unpin by name
gitmap agy pin-projects rm gitmap-v28

# Unpin by ID
gitmap agy pin-projects rm 46d05021-30e1-4036-acfb-3020489125eb

# Unpin all projects
gitmap agy pin-projects rm --all
```

---

## Filter Main Project List by Pinned

You can also filter the full project inventory table using the `--pinned` (`-p`) flag:

```bash
gitmap agy ls --pinned
gitmap agy ls -p
```
