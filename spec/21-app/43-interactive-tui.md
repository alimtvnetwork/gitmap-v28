# Interactive TUI Mode

## Overview

An interactive terminal user interface for browsing, searching, and
managing repositories. Built with Bubble Tea (charmbracelet) for a
rich, keyboard-driven experience.

---

## Command

```
gitmap-v28 interactive
gitmap-v28 i
```

Launches the full-screen TUI with four views accessible via tab
navigation.

---

## Views

### 1. Repo Browser (Default)

Full-screen list of all tracked repositories with:
- Fuzzy search (type to filter)
- Arrow keys / j/k to navigate
- Space to toggle selection (multi-select)
- Enter to open detail panel
- `/` to focus search input

Columns: Slug, Branch, Path, Project Type

### 2. Batch Actions

After selecting repos in the browser:
- `p` — Pull selected repos
- `x` — Run a git command across selected
- `s` — Show status for selected
- `g` — Add selected to a group

An action bar at the bottom shows available keybindings based on
current selection count.

### 3. Group Management

- List all groups with member counts
- Create new group (inline prompt)
- Add/remove repos from group (uses browser selection)
- Delete group (with confirmation)
- Arrow keys to navigate groups, Enter to show members

### 4. Status Dashboard

Live-refreshing view showing:
- Dirty/clean indicator per repo
- Current branch with ahead/behind counts
- Stash count
- Auto-refresh interval configurable via `dashboardRefresh` in
  `config.json` (seconds, default 30)
- Manual refresh via `r` key

### 5. Releases

Browse release history from the database:
- Version, tag, branch, source, date columns
- Enter to toggle detail view (commit SHA, changelog, notes)
- Draft / pre-release / latest indicators
- `r` to refresh from database

### 6. Logs

Browse recent command history from the database:
- Command, alias, args, duration, exit code, date columns
- Enter to toggle detail view (flags, summary, repo count)
- `r` to refresh from database
- Shows duration in human-readable format (ms, s, m)

---

## Package Structure

```
gitmap-v28/tui/
├── tui.go           # Entry point, root model, tab switching
├── browser.go       # Repo list with fuzzy search
├── actions.go       # Batch action executor
├── groups.go        # Group management view
├── dashboard.go     # Live status dashboard
├── releases.go      # Release history browser
├── relformat.go     # Release formatting helpers
├── logs.go          # Command history log viewer
├── logformat.go     # Log formatting helpers
├── zipgroups.go     # Zip group browser
├── aliases.go       # Alias browser
├── keys.go          # Key bindings
└── styles.go        # Lipgloss style definitions
```

All files under 200 lines. View count: 8.

---

## Key Bindings

| Key       | Context        | Action                        |
|-----------|----------------|-------------------------------|
| Tab       | Global         | Switch between views          |
| q / Esc   | Global         | Quit TUI                      |
| /         | Browser        | Focus search input            |
| j / ↓     | Browser/Groups | Move cursor down              |
| k / ↑     | Browser/Groups | Move cursor up                |
| Space     | Browser        | Toggle repo selection         |
| Enter     | Browser        | Show repo detail              |
| a         | Browser        | Select all                    |
| p         | Batch          | Pull selected repos           |
| x         | Batch          | Execute command on selected   |
| s         | Batch          | Show status for selected      |
| g         | Batch          | Add selected to group         |
| c         | Groups         | Create new group              |
| d         | Groups         | Delete group (with confirm)   |
| r         | Dashboard      | Force refresh                 |

---

## Dependencies

```
github.com/charmbracelet/bubbletea
github.com/charmbracelet/bubbles
github.com/charmbracelet/lipgloss
github.com/sahilm/fuzzy
```

---

## Data Access

The TUI reads from the SQLite database using existing `store` package
methods:
- `ListRepos()` — all tracked repos
- `FindBySlug()` — single repo lookup
- `ListGroups()` — all groups
- `GroupRepos()` — repos in a group

Batch actions delegate to existing command logic (pull, exec, status)
running in the background with output captured and displayed in the TUI.

---

## Constraints

- All files under 200 lines
- No magic strings — all text in `constants/constants_tui.go`
- Styles use lipgloss with colors matching terminal accent palette
- Graceful degradation: if terminal doesn't support TUI, fall back to
  standard output with an error message
- No external process spawning for navigation — all in-process
