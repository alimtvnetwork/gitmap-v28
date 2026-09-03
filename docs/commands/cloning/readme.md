# Cloning & Workspace Synchronization

Commands for batch cloning, version increment cloning, and IDE registration.

## Commands

### `gitmap clone <source|json|csv|text>`
* **Alias:** `c`
* Clones repositories listed in a scan output file.
* Flags: `--target-dir <dir>`, `--safe-pull`, `--verbose`, `--fix`.

### `gitmap clone-sync <url> [urls...]`
* **Alias:** `cs`
* Clones one or more repositories and immediately registers them with VS Code Workspaces and Antigravity.

### `gitmap clone-next [v++|vN]`
* **Alias:** `cn`
* Clones the next versioned iteration of the current repo (e.g. `v27` $\rightarrow$ `v28`).
* Flags: `--delete`, `--keep`, `--no-desktop`, `--ssh-key <name>`, `--create-remote`.

### `gitmap desktop-sync`
* **Alias:** `ds`
* Syncs tracked repositories to GitHub Desktop from output.

### `gitmap github-desktop`
* **Alias:** `gd`
* Registers the current repository with GitHub Desktop directly without requiring a scan.
