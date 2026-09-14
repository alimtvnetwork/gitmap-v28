# gitmap github-desktop

Register git repositories with GitHub Desktop.

> **`github-desktop` (`gd`) and `desktop-sync` (`ds`) are the same command.**
> Use whichever you remember.

## Alias

gd (also: ds, desktop-sync)

## Usage

    gitmap gd                       # register CWD or every DB-tracked repo under CWD
    gitmap ds                       # same thing
    gitmap gd D:\path\to\repo       # register an explicit folder
    gitmap gd --all                 # register every repo in the gitmap DB
    gitmap gd --install             # auto-install GitHub Desktop if missing, then register CWD
    gitmap gd group ls              # list all GitHub Desktop repository groups
    gitmap gd group add <name> <p...> # add repositories to a group
    gitmap gd group rm <name> [p]   # remove a repository or delete a group

## Flags

    --all           Register every repo currently tracked in the gitmap database,
                    regardless of where you ran the command from.
    -i, --install   Automatically invoke the GitMap installer to install GitHub Desktop
                    if the CLI is missing before registering repositories.

## Prerequisites

- GitHub Desktop installed with the `github` CLI on PATH.
- A git repo at CWD, an explicit path, or under a registered scan root.
- **No prior `gitmap scan` is required.**

## Missing CLI Remediation

When GitHub Desktop is not installed or the `github` CLI cannot be found on PATH, GitMap provides several installation options:

### Option 1: GitMap Installer (Recommended)

```bash
gitmap install github-desktop
# Or use the short alias:
gitmap in gd
# Or pass --install directly to auto-install and register:
gitmap gd --install
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

### Example 1: Register the current repo

    cd D:\work\macro-ahk
    gitmap gd

**Output:**

    ✓ Registered: macro-ahk
    GitHub Desktop: 1 added · 0 skipped · 0 failed

### Example 2: Register an explicit path

    gitmap github-desktop D:\projects\billing-svc

**Output:**

    ✓ Registered: billing-svc
    GitHub Desktop: 1 added · 0 skipped · 0 failed

### Example 3: Bulk-register under a scan root

    cd D:\wp-work
    gitmap gd

**Output:**

    [1/14] my-api ............ ✓
    [2/14] web-app ........... ✓
    ...
    GitHub Desktop: 14 added · 0 skipped · 0 failed

### Example 4: Auto-install GitHub Desktop if missing

    gitmap gd --install

**Output:**

    ✓ Installed: github-desktop
    ✓ Registered with GitHub Desktop: D:\work\macro-ahk

## See Also

- [desktop-sync](desktop-sync.md) — same command, different name
- [scan](scan.md) — `--github-desktop` registers during scan
- [clone](clone.md) — `--github-desktop` registers during clone
- [scan-gd (spec 102)](../02-spec/01-app/102-scan-gd.md) — design doc for bulk mode

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter github-desktop
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
