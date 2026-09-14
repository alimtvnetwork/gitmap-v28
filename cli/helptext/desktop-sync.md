# gitmap desktop-sync

Register git repositories with GitHub Desktop.

> **`desktop-sync` (`ds`) is an alias of `github-desktop` (`gd`).**
> Both commands do the same thing. Use whichever you remember.

## Alias

ds (also: gd, github-desktop)

## Usage

    gitmap ds                       # register CWD or every DB-tracked repo under CWD
    gitmap gd                       # same thing
    gitmap ds D:\path\to\repo       # register an explicit folder
    gitmap ds --all                 # register every repo in the gitmap DB
    gitmap ds --install             # auto-install GitHub Desktop if missing, then sync

## Flags

    --all           Register every repo currently tracked in the gitmap database,
                    regardless of where you ran the command from.
    -i, --install   Automatically invoke the GitMap installer to install GitHub Desktop
                    if the CLI is missing before registering repositories.

## Prerequisites

- GitHub Desktop installed with the `github` CLI on PATH.
- A git repo at CWD, explicit path, or under a registered scan root.
- **No prior `gitmap scan` is required.**

## Missing CLI Remediation

When GitHub Desktop is not installed or the `github` CLI cannot be found on PATH, GitMap provides several installation options:

### Option 1: GitMap Installer (Recommended)

```bash
gitmap install github-desktop
# Or use the short alias:
gitmap in gd
# Or pass --install directly to auto-install and register:
gitmap ds --install
```

### Option 2: Package Manager Fallback

Install via your operating system's standard package manager:

- **Linux (Snap):**
  ```bash
  sudo snap install github-desktop --beta
  ```
- **Windows (Winget):**
  ```powershell
  winget install --id GitHub.GitHubDesktop
  ```
- **macOS (Homebrew):**
  ```bash
  brew install --cask github-desktop
  ```

### Option 3: Official GUI Installer

Download and install directly from the official website:
- [GitHub Desktop Official Website](https://desktop.github.com)

## Resolution order (no args)

1. Is GitHub Desktop installed? If not → exit with install hint.
2. Is CWD itself a git repo (`.git` dir OR `.git` file for worktrees)? Register it.
3. Is CWD inside a registered scan root? Bulk-register every tracked repo under it.
4. Otherwise → friendly hint + exit 3.

## Examples

### Example 1: Register the current folder (single repo)

    cd D:\work\macro-ahk
    gitmap ds

**Output:**

    ✓ Registered: macro-ahk
    GitHub Desktop: 1 added · 0 skipped · 0 failed

### Example 2: Bulk-register everything under a scan root

    cd D:\wp-work
    gitmap gd

**Output:**

    [1/14] my-api ............ ✓
    [2/14] web-app ........... ✓
    [3/14] billing-svc ....... already registered
    ...
    GitHub Desktop: 12 added · 2 skipped · 0 failed

### Example 3: Register every DB-tracked repo regardless of CWD

    gitmap ds --all

### Example 4: Auto-install GitHub Desktop if missing

    gitmap ds --install

**Output:**

    ✓ Installed: github-desktop
    GitHub Desktop: 14 added · 0 skipped · 0 failed

## See Also

- [github-desktop](github-desktop.md) — same command, different name
- [scan](scan.md) — populate the database first if you want bulk mode
- [clone](clone.md) — `--github-desktop` registers as it clones

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter desktop-sync
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
