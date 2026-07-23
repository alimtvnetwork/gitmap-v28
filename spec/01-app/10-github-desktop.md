# GitHub Desktop Integration

## Overview

GitMap can register discovered or cloned repositories with GitHub Desktop
using the `--github-desktop` flag. This works on both the `scan` and
`clone` commands.

## Prerequisites

GitHub Desktop must be installed and its CLI tool (`github`) must be
available on the system `PATH`. On Windows, this is typically set up
automatically when GitHub Desktop is installed.

## Usage

### With Scan

```bash
# Scan and add all discovered repos to GitHub Desktop
gitmap scan ./projects --github-desktop
```

After scanning, each discovered repository is opened in GitHub Desktop
via `github <absolute-path>`. This registers the repo so it appears
in the GitHub Desktop repository list.

### With Clone

```bash
# Clone and add all successfully cloned repos to GitHub Desktop
gitmap clone ./.gitmap/output/gitmap.json --target-dir ./restored --github-desktop
```

After cloning, only **successfully cloned** repositories are registered.
Failed clones are skipped.

## Behavior

| Scenario | Result |
|----------|--------|
| GitHub Desktop installed | Repos are added one by one |
| GitHub Desktop not installed | Prints message, skips gracefully |
| Repo fails to register | Logged as failure, continues with next |
| Mixed success/failure | Summary shows counts for both |

## Output

The flag produces per-repo feedback and a summary:

```
  ✓ Added to GitHub Desktop: my-app
  ✓ Added to GitHub Desktop: core-lib
  ✗ Failed to add docs: exit status 1
GitHub Desktop: 2 added, 1 failed
```

## Implementation

| Component | File | Responsibility |
|-----------|------|----------------|
| `desktop` package | `gitmap/desktop/desktop.go` | Core integration logic |
| `constants` | `gitmap/constants/constants.go` | CLI binary name, messages |
| `cmd/scan.go` | `addToDesktop()` | Wires flag to scan workflow |
| `cmd/clone.go` | `registerCloned()` | Wires flag to clone workflow |

### Detection

The `desktop.isInstalled()` function uses `exec.LookPath` to check
if the `github` CLI is on the system `PATH`.

### Registration

Each repo is registered by calling `github <absolute-path>`, which
opens the repo in GitHub Desktop and adds it to the repository list.

## CLI Reference

| Command | Flag | Description |
|---------|------|-------------|
| `gitmap scan` | `--github-desktop` | Add discovered repos to GitHub Desktop |
| `gitmap clone` | `--github-desktop` | Add cloned repos to GitHub Desktop |

## Limitations

- Only works on systems where GitHub Desktop is installed.
- The `github` CLI must be on the system `PATH`.
- Each repo is opened individually, which may briefly flash the
  GitHub Desktop window per repo.
