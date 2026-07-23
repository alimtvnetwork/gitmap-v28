<div align="center">

# GitMap

**Git repository scanner, manager, and navigator CLI**

**Pinned version: v6.82.0**

[![CI](https://img.shields.io/github/actions/workflow/status/alimtvnetwork/gitmap-v28/ci.yml?style=flat-square&label=CI)](https://github.com/alimtvnetwork/gitmap-v28/actions/workflows/ci.yml)
[![Vulnerability Scan](https://img.shields.io/github/actions/workflow/status/alimtvnetwork/gitmap-v28/vulncheck.yml?style=flat-square&label=vulnerability%20scan)](https://github.com/alimtvnetwork/gitmap-v28/actions/workflows/vulncheck.yml)
[![Cross-Platform Build](https://img.shields.io/github/actions/workflow/status/alimtvnetwork/gitmap-v28/cross-platform.yml?style=flat-square&label=cross-platform%20build)](https://github.com/alimtvnetwork/gitmap-v28/actions/workflows/cross-platform.yml)
[![Startup Build-Tag Guard](https://img.shields.io/github/actions/workflow/status/alimtvnetwork/gitmap-v28/startup-build-tags.yml?style=flat-square&label=startup%20build-tag%20guard)](https://github.com/alimtvnetwork/gitmap-v28/actions/workflows/startup-build-tags.yml)
[![GitHub Release](https://img.shields.io/github/v/release/alimtvnetwork/gitmap-v28?style=flat-square&label=version)](https://github.com/alimtvnetwork/gitmap-v28/releases)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey?style=flat-square)](https://github.com/alimtvnetwork/gitmap-v28)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](./LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/alimtvnetwork/gitmap-v28/gitmap?style=flat-square)](https://goreportcard.com/report/github.com/alimtvnetwork/gitmap-v28/gitmap)

_Scan, catalog, clone, and manage all your Git repositories from a single CLI._

<br>

<img src="docs/demo.gif" alt="Animated terminal demo of GitMap commands: clone, scan, history, release, ssh, make-public/make-private, merge-left/right/both, and the interactive TUI" width="900">

<sub><i>30-second tour: clone · scan · history · release · ssh · visibility · merge · interactive TUI</i></sub>

</div>

---

## Why GitMap?

GitMap is extremely powerful — we built it to solve our own AI-coding workflow problems. Think of it as Git on steroids. It started in 2024 as a quick way to move every repo from my desktop to my laptop, and today our team can't get through a single day without it.

It manages Git repos (public and private), makes releases smooth, and every repo we ship is released through GitMap itself.

---

## 🚀 Install

GitMap is a **Windows-first** project. The commands below install the latest release with sensible defaults — no prompts, no drive picker. Use the **Quick** block if you want to pick a custom install drive.

### 🪟 Windows · PowerShell (latest)

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex
```

### 🐧 macOS · Linux · Bash (latest)

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | sh
```

### 📌 Pinned version (`v6.82.0`)

```powershell
# Windows · PowerShell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/release-version-v6.82.0.ps1 | iex
```

```bash
# macOS · Linux · Bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/release-version-v6.82.0.sh | bash
```

#### 🧭 Version matrix — `v6.82.0` install-script URLs

One-line installers and release assets for the pinned **v6.82.0** build across every supported platform. All URLs resolve to the exact release tag — no fallback, no discovery.

| Platform | Shell | Install-script URL | Release binary asset |
| --- | --- | --- | --- |
| **Windows** (amd64) | PowerShell | [`release-version-v6.82.0.ps1`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/release-version-v6.82.0.ps1) | [`gitmap-v28.82.0-windows-amd64.zip`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/gitmap-v28.82.0-windows-amd64.zip) |
| **macOS** (arm64) | Bash | [`release-version-v6.82.0.sh`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/release-version-v6.82.0.sh) | [`gitmap-v28.82.0-darwin-arm64.tar.gz`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/gitmap-v28.82.0-darwin-arm64.tar.gz) |
| **macOS** (amd64) | Bash | [`release-version-v6.82.0.sh`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/release-version-v6.82.0.sh) | [`gitmap-v28.82.0-darwin-amd64.tar.gz`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/gitmap-v28.82.0-darwin-amd64.tar.gz) |
| **Linux** (amd64) | Bash | [`release-version-v6.82.0.sh`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/release-version-v6.82.0.sh) | [`gitmap-v28.82.0-linux-amd64.tar.gz`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/gitmap-v28.82.0-linux-amd64.tar.gz) |
| **Linux** (arm64) | Bash | [`release-version-v6.82.0.sh`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/release-version-v6.82.0.sh) | [`gitmap-v28.82.0-linux-arm64.tar.gz`](https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.82.0/gitmap-v28.82.0-linux-arm64.tar.gz) |

> **Release page:** [github.com/alimtvnetwork/gitmap-v28/releases/tag/v6.82.0](https://github.com/alimtvnetwork/gitmap-v28/releases/tag/v6.82.0) · **Asset naming contract:** `gitmap-<version>-<os>-<arch>.<ext>` (`.zip` on Windows, `.tar.gz` elsewhere) — verified by the installer pre-flight HEAD probe.

### 🎯 Install — Quick (pick your install drive)

Use this **only** when you want to choose a specific drive or folder (e.g. install to `D:\` instead of the default location). It prompts for the install drive/folder, then delegates to the canonical installer above.

```powershell
# Windows · PowerShell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.ps1 | iex
```

```bash
# macOS · Linux · Bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.sh | bash
```

> **How install resolves a version:** every installer follows the generic contract in [`spec/07-generic-release/09-generic-install-script-behavior.md`](spec/07-generic-release/09-generic-install-script-behavior.md). In short — **strict tag mode** (`--version <tag>` / `-Version <tag>`) installs that exact release with **no fallback whatsoever** (no `latest`, no sibling probe, no main-branch HEAD; missing tag → exit 1). **Discovery mode** (no tag supplied) probes the next 20 `-v<N+i>` sibling repos in parallel, then falls back to `releases/latest`, and finally to the default branch HEAD as a last resort.

### 📁 Install to a custom directory

Pass an explicit install path on the one-liner — no prompts, no auto-detection. Useful for shared drives, portable installs, or pinning gitmap next to other dev tools. Replace `D:\tools\gitmap` / `/opt/gitmap` with your preferred path.

```bash
# Windows · PowerShell — pass -InstallDir through irm | iex
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1))) -InstallDir 'D:\tools\gitmap'

# macOS · Linux · Bash — pass --dir through curl | sh
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | sh -s -- --dir /opt/gitmap
```

> Add `-Version v5.63.0` (PowerShell) or `--version v5.63.0` (Bash) before the closing quote / end of line to pin a specific release into your custom directory.

---


## About GitMap

### What is GitMap?

**GitMap** ([this repository](https://github.com/alimtvnetwork/gitmap-v28))
is an all-in-one **Git workspace CLI**. It scans any folder tree,
catalogs every Git repository it finds, and lets you clone, mirror,
navigate, release, and manage all of them as a single map — instead
of one repo at a time. One binary, one `.gitmap/` directory, every
repo on your machine in one view.

### Why GitMap?

Because juggling dozens of repos with raw `git` commands does not
scale. Without GitMap you end up writing one-off shell scripts every
time you switch machines, audit branches across projects, or try to
remember which folder maps to which remote. GitMap replaces all of
that with a single, deterministic, cross-platform workflow:

- **One scan** to discover every repo on disk (however deeply nested).
- **One manifest** (`.gitmap/output/gitmap.{json,csv,txt}`) that is
  both a catalog AND a runnable migration script.
- **One reclone** to reproduce the exact folder layout on any new
  machine — no IDE junk, no build artifacts, just canonical source.
- **One CLI** for the everyday Git chores that normally need 5
  different tools: bulk visibility flips, SSH key sync, multi-repo
  merges, release packaging, interactive TUI, and shell completion.

If you have ever opened more than ~10 repos on the same machine,
GitMap pays for itself the first time you switch laptops.

### Why it exists — the two-hour origin story

GitMap started as a one-evening fix for a very ordinary problem.
The author needed to **migrate every single Git repository from one
laptop to a brand-new machine** — dozens of folders, scattered across
nested directories, each with its own remote, branch, and personal
quirks. Cloning them by hand would have taken a weekend; copying the
working trees would have dragged along build artifacts, half-finished
branches, and IDE junk that did not belong on a fresh box.

So in **about two hours, with the help of AI coding tools**, the
first version of `gitmap` was built: walk a folder, find every Git
repo, write a list of `git clone` commands, run them on the new
machine, and end up with the **exact same folder layout** — no
artifacts, no garbage, just the canonical source of every project.

That tiny utility worked. Then it kept growing. After **months of
daily use, refactors, and feature additions** it has turned into the
all-in-one Git workspace CLI you see today — a tool the author now
reaches for before almost any other Git operation.

### What GitMap actually does

At its heart `gitmap` does one thing extremely well: it treats your
disk as a **map of Git repositories** and lets you operate on that
map as a single object. Every command flows from that idea.

#### 🗺️ Scan & catalog
- Recursively walks any folder tree and discovers every Git repo
  underneath it (no matter how deeply nested).
- Records each repo's remote URLs (HTTPS **and** SSH), branch,
  relative path, and discovery URL into a deterministic
  `.gitmap/output/gitmap.{json,csv,txt}` manifest.
- Emits the same data as ready-to-run `git clone` instructions, so
  the catalog is **also** the migration script.

#### 🚚 Round-trip clone (`reclone`)
- Re-creates the **exact folder layout** of a previous scan on any
  new machine, using the recorded relative paths verbatim.
- Pre-flight safety prompt + dry-run summary + row-level manifest
  validation — you always see what's about to happen before any
  side effect touches your disk.
- Concurrent workers (`--max-concurrency`), `--on-exists` policy
  (skip / update / force), and HTTPS ⇄ SSH mode switching.

#### 🔢 Version tracking & history
- `clone-next` flattens versioned URLs (`…-v7`, `…-v8`, …) into a
  single base folder and records every cloned version in a local
  SQLite **`RepoVersionHistory`** table.
- `history`, `stats`, and `release` commands give you a per-repo
  timeline — when you cloned what, from where, on which machine.

#### 🤝 AI-tool friendly
- First-class helpers for working with **Lovable, Claude, Cursor,
  GitHub Copilot, and other AI coding agents**: structured manifests
  the agent can read, deterministic outputs that survive being
  diffed across runs, and command output formats designed to be
  pasted straight into a chat.
- Built-in `LLM.md` + `spec/` directory makes the codebase itself
  legible to AI — an explicit design choice, not an accident.

#### 🛠️ Self-managing installation
- `gitmap self-install` / `self-uninstall` manage the binary itself
  on every supported platform.
- Canonical installers (`gitmap/scripts/install.ps1` /
  `install.sh`) are the **default** one-liners — no prompts, sensible
  defaults, full PATH + data-folder setup. Quick installers
  (`install-quick.ps1` / `install-quick.sh`) layer a drive-picker
  prompt on top for users who want to install on a specific drive.
- `gitmap-updater` keeps the binary fresh; `self-uninstall` cleans
  up the PATH marker block and (optionally) the user data folder.

#### 🔀 Workspace operations
- `mv`, `merge-both`, `merge-left`, `merge-right` — move or merge
  two working trees with an interactive **L / R / S / A / B / Q**
  prompt and `--prefer-*` flags for non-interactive runs.
- `as` / `release-alias` (`ra`) / `release-alias-pull` (`rap`) —
  create labelled aliases of a release with concurrency-safe
  auto-stash/pop.
- `regoldens` (`rg`) — automated two-pass golden-fixture
  regeneration with built-in determinism verification.

#### 🖥️ Web docs UI
The repository ships with an **interactive documentation site**
(shown above) that mirrors every CLI command, every flag, and every
exit code — searchable, copy-paste-able, and synchronised with the
release metadata so the docs can never drift from the binary.

### TL;DR

> **A single CLI that maps, migrates, versions, and manages every
> Git repository on your machine — born from a two-hour migration
> hack, hardened into the author's daily driver.**

> **One-stop install/update reference**: [`spec/01-app/108-cross-platform-install-update.md`](spec/01-app/108-cross-platform-install-update.md) — the full Windows / macOS / Linux install · update · uninstall · verify matrix is also rendered in the docs at `/install-gitmap`.

## Quick Start

> Looking for install one-liners? See **[Install](#install)** at the top of this README.


### Uninstall — Quick (one-liner)

Removes the gitmap binary, deploy folder, PATH entries, and (optionally) the user data folder. First tries the canonical `gitmap self-uninstall`; falls back to a manual sweep if gitmap is no longer on PATH.

#### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/uninstall-quick.ps1 | iex
```

#### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/uninstall-quick.sh | bash
```

Useful flags (both scripts):

| Flag | Effect |
|---|---|
| `-Yes` / `-y` `--yes` | Skip the "delete user data?" prompt and assume yes |
| `-KeepData` / `--keep-data` | Always keep `%APPDATA%\gitmap` (Windows) or `~/.config/gitmap` (Unix) |
| `-InstallDir` / `--dir` | Override the auto-detected deploy root |

### Scan repos and see results

```bash
gitmap scan ~/projects
gitmap ls
```

The `[dir]` argument accepts **relative paths** and resolves them against
your current working directory. Common shortcuts:

```bash
gitmap scan .          # scan the current directory
gitmap scan ..         # scan the parent directory
gitmap scan ../..      # scan two folders up
gitmap scan ../../x    # scan the "x" folder two levels up
gitmap scan ~/work     # "~" expands to your home directory
```

When the resolved path differs from what you typed, gitmap prints a
one-line `↳ Resolved "<input>" → <abs>` hint so the target is
unambiguous. Non-existent paths exit early with a clear error instead of
silently falling back to the current directory.

### Navigate and pull

```bash
gitmap cd my-api
gitmap pull --all
```

Every command supports `--help` or `-h` for detailed usage with examples.

---

## Installation

### One-Liner Install (recommended)

The canonical installer (`install.ps1` / `install.sh`) is the **default**: no prompts, sensible install location, full PATH + data-folder setup. Use **install-quick** only when you want to choose the install drive.

#### Windows (PowerShell) — Default

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex
```

#### Linux / macOS (Bash) — Default

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | sh
```

#### Windows (PowerShell) — Full bootstrap (locked-down machines)

Use when execution policy / TLS settings block the short form above.

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1'))
```

#### Windows (PowerShell) — Quick (drive picker)

Prompts for the install drive/folder before delegating to the canonical installer.

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.ps1 | iex
```

#### Linux / macOS (Bash) — Quick (drive picker)

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.sh | bash
```

### Installer Options

**Windows (PowerShell):**

| Flag | Description | Example |
|------|-------------|---------|
| `-Version` | Pin a specific release | `-Version v2.51.0` |
| `-InstallDir` | Custom install directory | `-InstallDir C:\tools\gitmap` |
| `-Arch` | Force architecture (`amd64`, `arm64`) | `-Arch arm64` |
| `-NoPath` | Skip adding to user PATH | `-NoPath` |
| `-AllowFallback` | Use newest patch in same vMAJOR.MINOR if version missing | `-AllowFallback` |

**Linux / macOS (Bash):**

| Flag | Description | Example |
|------|-------------|---------|
| `--version` | Pin a specific release | `--version v2.55.0` |
| `--dir` | Custom install directory | `--dir /opt/gitmap` |
| `--arch` | Force architecture (`amd64`, `arm64`) | `--arch arm64` |
| `--no-path` | Skip adding to PATH | `--no-path` |
| `--allow-fallback` | Use newest patch in same vMAJOR.MINOR if version missing | `--allow-fallback` |

#### Non-Interactive / CI Installations

When installing via pipe (`irm ... | iex` or `curl ... | bash`), the terminal
is **non-interactive**. If the requested version is missing, the installer
**exits with code 1** without prompting.

To handle missing versions in automated environments:

1. **Use `--allow-fallback`** — Automatically picks the newest patch in the same
   minor series (e.g., `v3.38.0` requested but missing → uses `v3.38.5`):
   ```powershell
   irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v3.38.0/release-version-v3.38.0.ps1 | iex
   # Or with generic script:
   irm https://gitmap.dev/scripts/release-version.ps1 | iex; Install-Gitmap -Version "v3.38.0" -AllowFallback
   ```

2. **Pre-validate the version** — Use `gitmap list-versions` to confirm existence
   before installing.

#### Version-Specific Install Scripts

For reproducible installs, use the **per-version snapshot scripts** that are
 baked with the version at release time:

| Script | URL Pattern |
|--------|-------------|
| Pinned PowerShell | `https://github.com/alimtvnetwork/gitmap-v28/releases/download/{version}/release-version-{version}.ps1` |
| Pinned Bash | `https://github.com/alimtvnetwork/gitmap-v28/releases/download/{version}/release-version-{version}.sh` |
| Generic PowerShell | `https://gitmap.dev/scripts/release-version.ps1` (requires `-Version` param) |
| Generic Bash | `https://gitmap.dev/scripts/release-version.sh` (requires `--version` flag) |

**Specific version install (one-liner with fallback):**

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex; Install-Gitmap -Version "v2.51.0"
```

**Specific version + custom directory (one-liner):**

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex; Install-Gitmap -Version "v2.51.0" -InstallDir "D:\DevTools\gitmap"
```

**Custom directory install (downloaded script):**

```powershell
.\install.ps1 -InstallDir "D:\DevTools\gitmap"
```

**Pinned version + custom directory (downloaded script):**

```powershell
.\install.ps1 -Version v2.51.0 -InstallDir "C:\tools\gitmap"
```

**Linux / macOS — specific version:**

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | sh -s -- --version v2.51.0
```

> **Tip:** Use `gitmap list-versions` to see all available release versions before pinning.

### Clone & Setup (Development)

```bash
git clone https://github.com/alimtvnetwork/gitmap-v28.git gitmap
```

```bash
cd gitmap
./setup.sh
```

The setup script installs the pre-commit hook (golangci-lint), verifies your Go toolchain, and downloads dependencies. See [CONTRIBUTING.md](CONTRIBUTING.md) for the full development workflow.

### Update Source Before Building (avoid the `fileExists redeclared` regression)

If you have an existing local checkout, **always pull the latest source before
building**. Three releases (v3.92.0, v3.113.0, v3.114.0) eliminated a
`fileExists` symbol collision in `gitmap/cmd/`. A pre-v3.92.0 checkout
will fail to compile with:

```
cmd/updaterepo.go:118:6: fileExists redeclared in this block
        cmd/updatedebugwindows.go:148:6: other declaration of fileExists
```

This error is **always** a stale-checkout symptom — the current source
cannot produce it. Run the canonical update sequence:

```bash
cd /path/to/gitmap
git fetch origin
git checkout main
git pull --ff-only origin main
git status                          # must report "working tree clean"
git log -1 --format='%H %s'         # capture the SHA you're about to build
```

#### Verify the v3.92.0+ rename fix is present

Three quick checks confirm the redeclaration fix is in your tree. All
three must pass before you run `./run.sh` / `./run.ps1`:

**1. The declared version is v3.92.0 or newer:**

```bash
grep '^const Version = ' gitmap/constants/constants.go
# expected: const Version = "3.115.0"   (or higher)
```

**2. `gitmap/cmd/updatedebugwindows.go` does NOT declare a local helper:**

```bash
grep -nE '^func (fileExists|fileExistsLoose)\(' gitmap/cmd/updatedebugwindows.go
# expected: (no output — the helper moved to gitmap/fsutil in v3.113.0)
```

**3. The shared `fsutil` package exists and is imported by `cmd/`:**

```bash
test -f gitmap/fsutil/exists.go && echo "fsutil package present"
grep -l 'gitmap/fsutil' gitmap/cmd/updaterepo.go gitmap/cmd/updatedebugwindows.go
# expected: both file paths printed
```

If any check fails, your checkout is older than v3.113.0. Re-run the
update sequence above; if it still fails, your branch diverged before
the fix landed and needs `git rebase origin/main`.

#### Run the pre-build provenance stamp (recommended)

Since v3.115.0, `./run.sh` and `./run.ps1` automatically print a
provenance stamp before invoking `go build`:

```bash
bash scripts/build-stamp.sh           # prints SHA, version, file fingerprints
bash scripts/build-stamp.sh --strict  # exits 1 if a redeclaration risk is detected
```

A healthy stamp ends with:

```
guards
  redecl-risk-check       ok (no local fileExists* in cmd/ — fsutil migration present)
```

If you see `FAIL — fileExists/fileExistsLoose declared in both files`,
**stop and re-pull** — `go build` will fail with the redeclaration
error. The Windows equivalent is `pwsh scripts/build-stamp.ps1 -Strict`.

### Install-script behavior spec (shareable with any AI)

The canonical, repository-agnostic contract that every installer in this
project (and any sibling repo) MUST follow lives at:

> **[`spec/07-generic-release/09-generic-install-script-behavior.md`](spec/07-generic-release/09-generic-install-script-behavior.md)**

It defines the two install modes in one place:

- **Strict tag mode** — explicit `--version <tag>` installs that exact
  release with **no** fallback to `latest`, no `-v<N+i>` sibling probe,
  and no main-branch fallback. Missing tag → exit 1 with a canonical
  message.
- **Discovery mode** — no tag supplied → probe the next 20 `-v<N+i>`
  sibling repos in **parallel** (max-hit wins) → fall back to
  `releases/latest` → fall back to the default branch HEAD as a last
  resort.

The spec is intentionally generic (placeholders for `<owner>`, `<stem>`,
`<binary>`, `<installerPath>`) so you can hand it to any AI working on
any repository's installer and they will implement the same contract.

---

## What It Does

A portable CLI that scans directory trees for Git repositories, extracts clone URLs and branch info, and outputs structured data. Every scan produces **all outputs** automatically:

- **Terminal** — formatted table to stdout
- **CSV** — `gitmap.csv`
- **JSON** — `gitmap.json`
- **Folder Structure** — `folder-structure.md` (tree view of discovered repos)

All files are written to `.gitmap/output/` at the root of the scanned directory.

---

## Help UX — discover commands fast

`gitmap help` is bucketed into 5 intent super-categories with search, compact mode, and machine-readable JSON for IDE/script integrations.

```bash
gitmap help                       # full grouped help
gitmap help --compact             # one-line-per-command dense layout
gitmap help --groups              # only the group banners
gitmap help --filter clone        # substring search (alias: -f)
gitmap help --json                # machine-readable JSON (ANSI stripped)
gitmap help --json --filter ssh   # compose: scope JSON to matching rows
```

> **JSON Schema:** the `--json` payload conforms to [`spec/08-json-schemas/help-json.schema.json`](spec/08-json-schemas/help-json.schema.json) (JSON Schema draft 2020-12). A contract test validates runtime output against the schema on every build to prevent drift. Glyph rendering follows `--glyphs auto|rich|safe` (auto-falls-back to ASCII on legacy PowerShell 5.1).

---

## Command Reference


<div align="center">

### Scanning & Discovery

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `scan` | `s` | Scan directory for Git repos |
| `rescan` | `rsc` | Re-scan previously scanned directories |
| `list` | `ls` | Show all tracked repos with slugs |

```bash
gitmap scan ~/projects --output json --mode ssh
gitmap ls go                    # list Go projects
gitmap rescan                   # re-scan all known directories
```

→ [scan](gitmap/helptext/scan.md) · [rescan](gitmap/helptext/rescan.md) · [list](gitmap/helptext/list.md)

#### Scan rules — what counts as a repo, and how deep we walk

The scanner is intentionally strict so the catalog stays trustworthy.
These rules are stable across releases and are enforced by
[`gitmap/scanner/scanner.go`](gitmap/scanner/scanner.go) (see also
[`spec/01-app/03-scanner.md`](spec/01-app/03-scanner.md)).

**1. Repo markers — what makes a directory a "repo".** A directory is
recorded as a repo when it contains a child entry literally named
`.git` matching either of these forms:

- **`.git/` is a directory** → ✅ counted as a repo, **always**. This is the standard `git init` / `git clone` layout. The directory's own contents are not inspected — its mere presence is the marker.
- **`.git` is a regular file whose first bytes are exactly `gitdir:`** → ✅ counted as a repo. This is the `git worktree add` linked-checkout layout, and the layout submodules use when their `.git` was absorbed into the superproject. Only the first **256 bytes** are read; the prefix match is literal (lowercase `gitdir:`, no leading whitespace tolerated).
- **`.git` is a regular file but its contents do *not* start with `gitdir:`** → ❌ **ignored**. A stray `.git` text file (committed by accident, dropped by an editor, left over from a failed `git init`) does not create a false positive.
- **`.git` is a symlink** → ✅/❌ resolved as whichever target form above it points to (directory → counted; file → counted only if `gitdir:`-prefixed). A broken symlink is ignored.
- **`.git` is missing or unreadable** → ❌ ignored. Permission errors are treated as "not a marker" so one unreadable subtree does not abort the whole scan.
- **A directory ends in `.git` (e.g. `myrepo.git/`)** → ❌ ignored. Bare repos and `*.git` mirror folders are not catalogued by `gitmap scan`. Only a child entry named **exactly** `.git` qualifies.
- **Any other hand-rolled hint** (loose `HEAD` files, a `refs/` folder, a `config` with `[core]`) → ❌ ignored. The two forms above are the only positive signals.

The same rules in table form (kept for backwards-compatible cross-references):

| Marker form | Matches | Notes |
|---|---|---|
| `.git/` directory | Standard `git init` / `git clone` checkout | Counted unconditionally. |
| `.git` regular file | `git worktree add` linked checkouts; submodules whose `.git` was absorbed into the superproject | Only counted when the file's contents start with the literal prefix `gitdir:` (read budget: 256 bytes). A stray `.git` text file without that prefix is ignored to prevent false positives. |

Anything else under the directory — bare repos, `*.git` mirror folders,
hand-rolled `HEAD` files — is **not** treated as a repo by `gitmap scan`.

**2. Gitdir / worktree handling — no descent into discovered repos.**
Once a directory is recorded as a repo (by either marker form above),
the scanner does **not** descend into its subtree. This means:

- A worktree's `.git` file (`gitdir: /path/to/main/.git/worktrees/<name>`)
  registers the worktree directory itself as a repo. The linked main
  repo is recorded separately when its own `.git/` directory is reached
  via the normal walk.
- Submodules with absorbed `.git` files are recorded as repos at their
  own location, independent of the superproject.
- Nested repos hidden under a discovered repo (e.g. a `vendor/` checkout
  under a project that itself is a repo) are **not** discovered. Move
  them outside, scan them separately, or scan their parent directly.

**3. Default 4-level nesting cap.** The scanner refuses to descend more
than `DefaultMaxDepth = 4` directory levels below the scan root, even
when no `.git` marker has been found on the path. Depth is counted
from the root:

| Depth | Example path under `gitmap scan ~/code` |
|---|---|
| `0` | `~/code` (the scan root itself) |
| `1` | `~/code/<org>` |
| `2` | `~/code/<org>/<project>` |
| `3` | `~/code/<org>/<project>/<service>` |
| `4` | `~/code/<org>/<project>/<service>/<module>` |
| `5+` | **Not walked** under the default cap |

The cap exists to prevent runaway walks into dependency trees that
slipped past the exclude list (e.g. a forgotten `node_modules/` deep
inside a project). Repos discovered at any depth still stop their own
subtree from descending — the cap only matters for paths that have
**not** hit a `.git` marker yet.

Override via `ScanOptions.MaxDepth` when calling the library directly:
a positive value sets a custom cap, a negative value disables the cap
entirely (legacy unbounded behavior), and zero (the field's zero
value) keeps `DefaultMaxDepth = 4`.

**Excluded directory names** are skipped before the depth check fires
— see the per-scan `--config` exclude list and the project defaults
documented in [`gitmap/helptext/scan.md`](gitmap/helptext/scan.md).

#### Scan examples — markers, worktrees, and the depth cap in action

The three scenarios below show the rules above as concrete commands
and the rows you would (or would not) see. CSV output is shown
because v3.150.0 added a trailing `depth` column that makes the
4-level cap auditable at a glance.

**Example A — marker detection (`.git/` dir vs `.git` worktree file vs stray text file).**

Layout:

```text
~/code/
├── alpha/         .git/                     ← standard checkout
├── beta/          .git  (file: "gitdir: …") ← worktree-style marker
└── gamma/         .git  (file: "hello")     ← stray text, NOT a repo
```

```bash
gitmap scan ~/code --output csv
```

```csv
repoName,httpsUrl,sshUrl,branch,branchSource,relativePath,absolutePath,cloneInstruction,notes,depth
alpha,…,…,main,remote-head,alpha,/home/u/code/alpha,git clone …,,1
beta,…,…,main,remote-head,beta,/home/u/code/beta,git clone …,,1
```

`gamma/` is silently dropped: its `.git` file lacks the `gitdir:`
prefix, so rule 1 rejects it. Only the two valid markers produce
rows, both at depth 1 (immediate children of the scan root).

**Example B — worktrees and absorbed submodules don't double-count and don't hide siblings.**

Layout:

```text
~/work/
├── main-repo/                       .git/      ← superproject
│   ├── vendor/lib/                  .git/      ← nested repo: HIDDEN by rule 2
│   └── modules/auth/                .git file  ← absorbed submodule
└── main-repo-feature-x/             .git file  ← `git worktree add` checkout
```

```bash
gitmap scan ~/work --output csv
```

```csv
repoName,httpsUrl,sshUrl,branch,branchSource,relativePath,absolutePath,cloneInstruction,notes,depth
main-repo,…,…,main,remote-head,main-repo,/home/u/work/main-repo,git clone …,,1
main-repo-feature-x,…,…,feature-x,head,main-repo-feature-x,/home/u/work/main-repo-feature-x,git clone …,,1
```

The worktree checkout is its own row (rule 1, gitdir-prefixed file).
The absorbed submodule under `main-repo/modules/auth` and the
`vendor/lib` checkout are **both hidden** because rule 2 stops descent
the moment `main-repo/.git/` is recorded. To catalog them, scan
`~/work/main-repo/modules/` or `~/work/main-repo/vendor/` directly.

**Example C — what the default 4-level cap skips.**

Layout (depths annotated):

```text
~/mono/                                                       depth 0
├── team-a/                                                   depth 1
│   └── service-x/                          .git/             depth 2  ← FOUND
├── team-b/proj/svc/mod/                    .git/             depth 4  ← FOUND (at cap)
└── team-c/area/group/proj/svc/             .git/             depth 5  ← SKIPPED
```

```bash
gitmap scan ~/mono --output csv
```

```csv
repoName,httpsUrl,sshUrl,branch,branchSource,relativePath,absolutePath,cloneInstruction,notes,depth
service-x,…,…,main,remote-head,team-a/service-x,/home/u/mono/team-a/service-x,git clone …,,2
mod,…,…,main,remote-head,team-b/proj/svc/mod,/home/u/mono/team-b/proj/svc/mod,git clone …,,4
```

The depth-5 repo under `team-c/` is **not** in the output: the walker
read `team-c/area/group/proj/svc/` (depth 4) and saw no `.git` marker
on the path, so it refused to enqueue depth-5 children. To catch it,
either scan deeper into the subtree:

```bash
gitmap scan ~/mono/team-c --output csv      # new root → depth resets to 0
```

…or override the cap globally (positive = custom, negative =
unbounded):

```bash
# library callers only — set ScanOptions.MaxDepth = -1 for legacy
# unbounded walks; this is intentionally not a CLI flag so casual
# `gitmap scan` stays fast and bounded by default.
```

#### Copy-paste scan commands per scenario

Every example above used the minimal `gitmap scan <root> --output csv`
form so the marker / depth / rule-2 logic stayed in focus. In real
projects you almost always want the full triple — `--config` to pin
your exclude list, `--mode` to fix the URL column, and `--output csv`
to land on a spreadsheet-friendly artifact at a known path. The
blocks below reproduce each scenario above with that full triple,
copy-paste ready.

All blocks assume a project-local `gitmap.config.json` similar to
the one in [`data/config.json` exclude list](#dataconfigjson-exclude-list--sample-and-interaction-with-the-depth-cap)
below; substitute your own `--config` path freely. Output lands in
`./.gitmap/output/gitmap.csv` by default — pass `--output-path
<dir>` to redirect.

**Reproduce Example A — marker detection (`.git/` dir vs `.git` worktree file vs stray text file):**

```bash
# HTTPS clone URLs, project-local config, CSV to ./.gitmap/output/gitmap.csv
gitmap scan ~/code \
  --config ./gitmap.config.json \
  --mode https \
  --output csv

# SSH variant — same repos, sshUrl column populated, httpsUrl empty/secondary
gitmap scan ~/code \
  --config ./gitmap.config.json \
  --mode ssh \
  --output csv \
  --output-path ./reports/markers
```

The CSV header and row contract are identical between `--mode https`
and `--mode ssh`; only the `httpsUrl` / `sshUrl` column emphasis
shifts (both columns are always emitted; `--mode` selects which one
is used to build the `cloneInstruction` column).

**Reproduce Example B — worktrees and absorbed submodules:**

```bash
# Standard run — main-repo and main-repo-feature-x get rows;
# vendor/lib and modules/auth are hidden by rule 2.
gitmap scan ~/work \
  --config ./gitmap.config.json \
  --mode https \
  --output csv

# To also catalog the absorbed submodules / nested vendor checkouts,
# re-aim at their parent (rule 2 only stops descent under a recorded
# repo — these subdirs become depth-1 in their own scan):
gitmap scan ~/work/main-repo/modules \
  --config ./gitmap.config.json \
  --mode https \
  --output csv \
  --output-path ./reports/submodules

gitmap scan ~/work/main-repo/vendor \
  --config ./gitmap.config.json \
  --mode https \
  --output csv \
  --output-path ./reports/vendor
```

**Reproduce Example C — the 4-level depth cap:**

```bash
# Standard run — service-x (depth 2) and team-b/proj/svc/mod (depth 4)
# are emitted; team-c/area/group/proj/svc (depth 5) is skipped.
gitmap scan ~/mono \
  --config ./gitmap.config.json \
  --mode https \
  --output csv

# To catch the depth-5 repo, re-root at the at-cap directory:
gitmap scan ~/mono/team-c \
  --config ./gitmap.config.json \
  --mode https \
  --output csv \
  --output-path ./reports/team-c
```

The two scans compose additively in the database (upsert by
`AbsolutePath`) — running both produces one row per repo, never
duplicates. See [Reading at-cap CSV rows and rescanning a deeper
subfolder](#reading-at-cap-csv-rows-and-rescanning-a-deeper-subfolder)
for the full recipe and rationale.

**Reproduce the edge-case layout — skipped non-repos and marker-like cases:**

```bash
# All 7 negative cases are silently dropped; only real-repo,
# nested-under-real/inner, and worktree-link appear in the CSV.
gitmap scan ~/edge \
  --config ./gitmap.config.json \
  --mode https \
  --output csv \
  --output-path ./reports/edge-cases
```

To verify the silence, list every directory under `~/edge` that has
a `.git` child of any kind, then diff against the CSV's
`absolutePath` column:

```bash
# All candidate directories (anything with a .git child).
find ~/edge -maxdepth 2 -name .git -printf '%h\n' | sort > /tmp/candidates.txt

# What the scanner actually catalogued.
awk -F, 'NR>1 {print $7}' ./reports/edge-cases/gitmap.csv | sort > /tmp/found.txt

# The lines unique to /tmp/candidates.txt are the silent skips.
comm -23 /tmp/candidates.txt /tmp/found.txt
```

This is the canonical way to audit scan coverage in CI: the
`comm -23` output should be empty for healthy trees and equal to
the expected skip set for trees that intentionally include
edge-case fixtures (e.g. test repos for gitmap itself).

#### Reading at-cap CSV rows and rescanning a deeper subfolder

A `depth` value equal to the cap (`4` under defaults) is the
diagnostic signal you should learn to read. It does not mean "this
repo is 4 levels deep and that's all there is to know" — it means
**this row sits on the boundary, and any repos hidden in its
subtree were silently skipped on this scan**. Rows with `depth < 4`
are unambiguous: the walker reached them and recorded them, full
stop. Rows with `depth == 4` are the "investigate this" pile.

How to interpret a single row at a glance:

| `depth` value in the row | What it tells you | Action |
|---|---|---|
| `0`–`3` | Discovered well inside the cap. Nothing under it could have been skipped *for cap reasons* (rule 2 / exclude list still apply). | None — trust the row. |
| `4` (= `DefaultMaxDepth`) | Discovered exactly at the cap. The walker did NOT enqueue its depth-5+ children. If you expected nested repos under this row, they were silently dropped. | Re-scan that one subtree (recipe below). |
| (negative or absent) | You're not on a current build, or the column was stripped by post-processing. | Re-run with the latest `gitmap` binary; the column has been mandatory since v3.150.0. |

**Spot-check recipe — find every at-cap row in one shell pipe:**

```bash
# CSV — assumes the default header order (depth is column 10)
gitmap scan ~/code --output csv | awk -F, 'NR>1 && $10==4 {print $7}'

# JSON — same idea, jq filter
gitmap scan ~/code --output json | jq -r '.[] | select(.depth==4) | .absolutePath'
```

Each printed `absolutePath` is a candidate for the deeper-subfolder
rescan below.

**Rescan recipe — point the CLI at the deeper subfolder:**

When an at-cap row hides nested repos you actually want catalogued,
the simplest, most predictable fix is to re-run `gitmap scan` with
the at-cap directory itself as the new root. Depth resets to `0` at
the new root, so the cap effectively shifts 4 levels deeper into
the original tree:

```bash
# Original scan caps out — say team-c/area/group/proj/svc/ shows depth=4
gitmap scan ~/mono --output csv

# Re-aim the CLI at the at-cap directory; depth resets to 0 there,
# so its children (which were depth-5 in the original scan) are now
# depth-1 in the new scan and get walked normally.
gitmap scan ~/mono/team-c/area/group/proj/svc --output csv
```

Two important properties of this recipe:

- **It does not modify the original `last-scan.json`'s root**, so
  `gitmap rescan` from the parent shell will still replay the
  original `~/mono` scan. The deeper-subfolder run is its own
  independent scan-cycle (its own root, its own cached parameters,
  its own database upserts). If you want the deeper root to become
  the *default* for future `gitmap rescan`, run the deeper command
  from the same shell and it will overwrite `last-scan.json`.
- **It composes additively in the database.** Repos discovered by
  the deeper scan are upserted by `absolutePath`, so a repo that
  appears in both the shallow and deep scans is one row in the DB,
  not two. There's no risk of duplicate entries from running both.

When you genuinely need a single command that crosses the cap (e.g.
in CI where re-rooting per subtree is awkward), the library-level
override is `ScanOptions.MaxDepth = -1` for unbounded walks; this
is intentionally not exposed as a CLI flag so casual `gitmap scan`
invocations stay fast and bounded by default. See the
[scanner package docs](gitmap/scanner/scanner.go) for the call
shape.

#### CSV column reference — the 10 columns, in order

The CSV header is **stable across releases** (locked by
[`gitmap/formatter/csv_header_contract_test.go`](gitmap/formatter/csv_header_contract_test.go))
and produced by [`ScanRecord` in `gitmap/model/record.go`](gitmap/model/record.go).
Line endings are always `\r\n` (RFC 4180), the separator is a comma,
and fields containing commas, quotes, or newlines are double-quoted
per RFC 4180 — never escaped with backslashes.

```csv
repoName,httpsUrl,sshUrl,branch,branchSource,relativePath,absolutePath,cloneInstruction,notes,depth
```

| # | Column | Type | Source / meaning |
|---|---|---|---|
| 1 | `repoName` | string | Basename of the repo's working tree directory. For `~/code/alpha/.git/` this is `alpha`. |
| 2 | `httpsUrl` | string | `https://…` form of `origin`. Empty if the repo has no `origin` remote. |
| 3 | `sshUrl` | string | `git@host:owner/repo.git` form of `origin`. Empty if no `origin` remote. |
| 4 | `branch` | string | The branch we'd clone / check out: the remote `HEAD` target if known, otherwise the local `HEAD` branch, otherwise empty. |
| 5 | `branchSource` | enum | How column 4 was determined: `remote-head` (preferred), `head` (fallback to local `HEAD`), or empty when neither resolved. |
| 6 | `relativePath` | string | Path from the scan root to the repo. `team-a/service-x` for a repo at depth 2. |
| 7 | `absolutePath` | string | OS-absolute path. On Windows, drive-letter form (`C:\…`); separators are the host's native form. |
| 8 | `cloneInstruction` | string | Ready-to-paste `git clone -b <branch> <url> <relativePath>` command. The `-b <branch>` segment is included when columns 4–5 resolved a branch; the URL form follows the scan's `--mode` (`https` by default, `ssh` with `--mode ssh`). |
| 9 | `notes` | string | Free-text diagnostics. Empty for clean rows. May contain commas → will be quoted. |
| 10 | `depth` | integer | Directory level relative to the scan root: `0` = root, `1` = immediate child, capped at `DefaultMaxDepth = 4` unless `ScanOptions.MaxDepth` was overridden. |

#### Skipped non-repos and marker-like edge cases — what does NOT appear in CSV

`gitmap scan` is silent about everything it skips. The table below
makes that silence explicit: each row is a directory layout that
might *look* like a repo, the reason it was rejected, and a worked
example of the CSV that the scan produced (showing only the
neighbours that DID match, with the `depth` column called out).

Layout:

```text
~/edge/                                                  depth 0
├── real-repo/             .git/                         depth 1  ← FOUND
├── stray-text/            .git  (file: "hello world")   depth 1  ← skipped (no gitdir: prefix)
├── empty-dotgit/          .git  (file: 0 bytes)         depth 1  ← skipped (empty file != gitdir:)
├── uppercase/             .Git/                         depth 1  ← skipped (case-sensitive name match)
├── trailing-dotgit/       myrepo.git/                   depth 1  ← skipped (only literal ".git" qualifies)
├── bare-mirror.git/       HEAD, refs/, config           depth 1  ← skipped (bare repos are not catalogued)
├── broken-symlink/        .git -> /missing/path         depth 1  ← skipped (unresolvable symlink)
├── unreadable/            .git/  (chmod 000)            depth 1  ← skipped (read error treated as "no marker")
├── nested-under-real/                                   depth 1
│   └── inner/             .git/                         depth 2  ← FOUND (parent has no .git of its own)
├── worktree-link/         .git  (file: "gitdir: ...")   depth 1  ← FOUND (worktree marker)
└── deep/a/b/c/d/          .git/                         depth 5  ← skipped (past 4-level cap)
```

```bash
gitmap scan ~/edge --output csv
```

```csv
repoName,httpsUrl,sshUrl,branch,branchSource,relativePath,absolutePath,cloneInstruction,notes,depth
real-repo,https://github.com/u/real-repo.git,git@github.com:u/real-repo.git,main,head,real-repo,/home/u/edge/real-repo,git clone -b main https://github.com/u/real-repo.git real-repo,,1
inner,https://github.com/u/inner.git,git@github.com:u/inner.git,main,head,nested-under-real/inner,/home/u/edge/nested-under-real/inner,git clone -b main https://github.com/u/inner.git nested-under-real/inner,,2
worktree-link,https://github.com/u/main-repo.git,git@github.com:u/main-repo.git,feature-x,head,worktree-link,/home/u/edge/worktree-link,git clone -b feature-x https://github.com/u/main-repo.git worktree-link,,1
```

Why each skipped row was rejected — match the row to the rules in
the §1 bullet list above:

| Layout | Reason it's NOT in the CSV | Rule reference |
|---|---|---|
| `stray-text/.git` | File contents do not start with `gitdir:` | §1 bullet 3 |
| `empty-dotgit/.git` | Empty file — no `gitdir:` prefix to match | §1 bullet 3 |
| `uppercase/.Git/` | The marker name is case-sensitive (`.git`, lowercase) | §1 bullet 6 |
| `trailing-dotgit/myrepo.git/` | Only an entry named **exactly** `.git` qualifies | §1 bullet 6 |
| `bare-mirror.git/` | Bare repos and `*.git` mirror folders are out of scope | §1 bullet 6 |
| `broken-symlink/.git` | Symlink target does not exist → ignored | §1 bullet 4 |
| `unreadable/.git/` | Permission denied → treated as "no marker", not as an error | §1 bullet 5 |
| `deep/a/b/c/d/.git/` | Depth 5 exceeds `DefaultMaxDepth = 4` | §3 (depth cap) |

**Note on `nested-under-real/inner`**: it *is* discovered (depth 2)
because `nested-under-real/` itself is **not** a repo — it has no
`.git` of its own — so rule 2 ("no descent into a discovered repo")
does not fire. Rule 2 only stops descent under a directory that is
itself recorded as a repo. This is the most common source of
"why did gitmap find / miss this checkout?" confusion: check the
`depth` column and the parent's status, in that order.

**`notes` column for skipped rows**: the `notes` field is empty for
clean rows in the example above. A future change that surfaces *why*
a directory was skipped (e.g. emitting a row with
`notes="skipped: gitdir prefix missing"` and an otherwise empty
record) is tracked but **not yet implemented**; today, silence is
the contract and the CSV contains only successful discoveries.

#### `rescan` and the depth cap — what survives, what disappears, what's newly found

`gitmap rescan` is **not** an incremental diff against the database.
It reads the cached `last-scan.json` and replays the original
`gitmap scan` command verbatim — same root, same `--config`, same
`--mode`, and the same `DefaultMaxDepth = 4`. The output is whatever
a fresh walk produces today, period. Concretely:

| Repo state at rescan time | Was previously discovered? | Result |
|---|---|---|
| Still at depth ≤ 4 with a valid `.git` marker | yes | Re-emitted, same row, possibly with refreshed branch / clone-instruction. |
| Still at depth ≤ 4, but its `.git` marker is gone (deleted, broken worktree) | yes | **Dropped from the new output**. The DB row is reconciled away on the next scan-cycle commit — `rescan` does not "remember" it just because it was there last time. |
| Moved deeper than depth 4 since the last scan | yes | **Silently disappears** from the rescan output. The cap fires before the marker is reached; the previous discovery grants no special exemption. |
| Newly added at depth ≤ 4 (e.g. a fresh `git worktree add` checkout placed beside the superproject) | no | **Picked up** as a new row, depth filled in from the walker. Worktree-style `.git` files are recognized via the `gitdir:` prefix exactly as on the first scan (rule 1). |
| Newly added worktree **inside** a discovered repo's subtree | no | **Not picked up.** Rule 2 (no descent into a discovered repo) still wins — the worktree lives under the superproject's stopped subtree. Move it outside, or scan its parent directly. |
| Newly added at depth ≥ 5 | no | **Not picked up** for the same reason fresh `scan` would miss it: the cap fires before depth 5 is enqueued. |

The takeaway: `rescan` is a *replay*, not a *delta*. To widen what it
sees, edit the cached scan parameters (re-run `gitmap scan <root>`
with a different `--config` or a shallower starting root, which
overwrites `last-scan.json`). The library-level `ScanOptions.MaxDepth`
override applies to both `scan` and `rescan` since they share the
same walker — set it negative for unbounded walks when you really
need to catch a depth-5+ repo without restructuring directories.

#### `data/config.json` exclude list — sample and interaction with the depth cap

The `excludeDirs` field in `data/config.json` is a list of directory
**base names** (not paths, not globs) that the walker drops *before*
enqueueing — which means the exclude check runs strictly **before**
the depth check, so an excluded directory costs zero of your 4-level
budget. This is the difference that makes deep monorepos scan
quickly without raising the cap.

A representative `data/config.json`:

```json
{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": ".gitmap/output",
  "excludeDirs": [
    "node_modules",
    "vendor",
    ".venv",
    "venv",
    "__pycache__",
    "target",
    "dist",
    "build",
    ".next",
    ".cache",
    ".terraform"
  ],
  "notes": "",
  "dashboardRefresh": 0
}
```

Matching rules — short and strict:

- **Exact basename, case-sensitive.** `node_modules` excludes
  `~/code/app/node_modules/` but **not** `Node_Modules` or
  `node_modules.bak`.
- **No path patterns, no globs.** `vendor/protos` is not a valid
  entry — use `vendor` and accept that every `vendor/` directory in
  the tree is skipped.
- **Applied at every depth from 1 upward.** The check fires inside
  `handleSubdir` before the child is enqueued, so a `node_modules/`
  at depth 2 and one at depth 4 are both dropped equally.
- **Does not affect already-discovered repos.** A repo whose own
  basename appears in `excludeDirs` (rare, but possible) will still
  be skipped — exclusion wins over discovery for that directory.

How it interacts with `DefaultMaxDepth = 4`. Consider this layout:

```text
~/mono/                                                       depth 0
├── team-a/                                                   depth 1
│   └── service-x/                                            depth 2
│       └── node_modules/...500 dirs.../leaf/  .git/          depth 3+ ← never walked
└── team-b/proj/svc/                                          depth 4
    └── mod/                                    .git/         depth 5  ← skipped by cap
```

With `excludeDirs: ["node_modules"]`:

| Path | What happens | Why |
|---|---|---|
| `team-a/service-x/node_modules/...` | Pruned at depth 3 | `node_modules` matches the exclude basename → never enqueued, depth check never fires. The hundreds of nested directories inside cost zero budget. |
| `team-b/proj/svc/mod/` (depth 5) | Still skipped | Exclude list does **not** raise the cap. Depth check fires at depth-5 enqueue and refuses. |
| `team-b/proj/svc/` (depth 4) | Walked, no `.git` found | Inside the cap; its depth-5 children are not enqueued. |

In other words: **the exclude list buys you walk speed, not walk
depth.** A bloated `node_modules/` deep inside `team-a/` no longer
slows the scan or eats the budget — but a legitimate repo that
genuinely lives at depth 5 still requires either restructuring or
the `ScanOptions.MaxDepth` library override.

Edit the file, then re-run `gitmap scan <root>` (which overwrites
`last-scan.json` and seeds future `gitmap rescan` calls with the new
exclude list). The two-stage `data/config.json` validation added in
v3.149.0 will reject malformed enums or missing required keys before
the walker starts, so a typo here fails fast rather than silently
expanding the walk.

#### Worktree markers under excluded directories — confirmed skipped

The exclude check runs **before** any repo-marker check, so any kind
of `.git` marker — directory, `gitdir:` worktree file, or stray text
file — buried inside an excluded basename is **never inspected and
never reported.** The walker prunes the parent and moves on; nothing
inside is enqueued, opened, or stat'd beyond the `os.ReadDir` entry
that proved the basename matches.

This matters in practice because tools love to drop real worktrees
into directories you almost certainly want to skip:

- `pnpm` / `yarn` workspaces sometimes link sibling packages into
  `node_modules/<scope>/<pkg>` via a `.git` file pointing at the
  source repo's `.git/worktrees/<name>`.
- Vendored dependency mirrors (`vendor/<dep>/.git/`) — full clones
  kept for offline builds, especially in Go monorepos pre-modules.
- Build outputs (`dist/`, `target/`, `.next/`) that a release
  pipeline accidentally `git init`'d while debugging.
- IDE caches (`.cache/`, `.terraform/`) where a plugin checked out
  a helper repo.

Consider this layout with the default `excludeDirs` from above:

```text
~/mono/                                              depth 0
├── app/                                             depth 1
│   ├── .git/                                        ← REPO (discovered)
│   ├── node_modules/                                depth 2  ← excluded
│   │   ├── @scope/linked-pkg/
│   │   │   └── .git              ← gitdir: …/main/.git/worktrees/linked
│   │   └── legit-dep/
│   │       └── .git/             ← full nested clone
│   ├── vendor/                                      depth 2  ← excluded
│   │   └── upstream-fork/
│   │       └── .git/             ← real repo, mirrored offline
│   └── dist/                                        depth 2  ← excluded
│       └── .git/                 ← stray init from a CI experiment
└── tools/helper/                                    depth 2
    └── .git                      ← gitdir: …/tools/.git/worktrees/helper
```

What `gitmap scan ~/mono` reports:

| Path | Marker kind | In CSV? | Reason |
|---|---|---|---|
| `app/` | `.git/` directory | ✅ yes | Discovered at depth 1, before any exclude check applies to its children. |
| `app/node_modules/@scope/linked-pkg/.git` | `gitdir:` worktree file | ❌ no | `node_modules` pruned at depth 2 — the worktree file is never opened. |
| `app/node_modules/legit-dep/.git/` | `.git/` directory | ❌ no | Same prune; the directory entry is never read. |
| `app/vendor/upstream-fork/.git/` | `.git/` directory | ❌ no | `vendor` is excluded; the mirror is invisible to the scan. |
| `app/dist/.git/` | `.git/` directory | ❌ no | `dist` is excluded; the stray init is invisible. |
| `tools/helper/.git` | `gitdir:` worktree file | ✅ yes | `tools` is **not** excluded; the worktree file is read, the `gitdir:` prefix matches, and it is treated as a repo root. |

Two consequences worth internalizing:

1. **No silent diagnostic for excluded hits.** Because the prune
   happens before marker inspection, gitmap cannot tell you "I saw a
   `.git` inside `node_modules/` and skipped it" — the file was
   never even classified. If you suspect a real repo lives under an
   excluded basename, the only way to confirm is to scan it
   directly: `gitmap scan ~/mono/app/node_modules` (which still
   honors the exclude list at *its* depth-1 children, so for the
   most direct check, point the CLI deeper, e.g. `gitmap scan
   ~/mono/app/node_modules/@scope`).
2. **Worktrees vs nested clones are treated identically by the
   exclude rule.** The `gitdir:` prefix detection only happens
   *after* a directory is enqueued. Excluding a basename short-
   circuits both kinds of markers in the same step — there is no
   "prefer worktree" or "prefer real .git" preference to configure.

Quick verification recipe — confirms the survivors-only contract:

```bash
gitmap scan ~/mono --output csv --output-path ./reports/mono
awk -F, 'NR>1 {print $7}' ./reports/mono.csv | sort
# expected output (relative paths):
#   app
#   tools/helper
```

The two excluded-directory worktrees and the two excluded-directory
nested clones never appear — neither in the CSV, nor in the JSON,
nor in any error or warning channel. Silence is the contract.

---

<div align="center">

### Cloning & Sync

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `clone` | `c` | Clone from a structured file OR a direct URL |
| `clone-next` | `cn` | Clone next versioned iteration of current repo |
| `clone-fix-repo` | `cfr` | One-shot: `clone` → `cd` → `fix-repo --all` |
| `clone-fix-repo-pub` | `cfrp` | `cfr` + `make-public --yes` at the end |
| `push` | `ph` | `git push` in cwd, with `--ssh` / `--https` transport coercion + auto `pull --rebase` retry on non-fast-forward rejection |
| `pull` | — | `git pull` in cwd, with `--ssh` / `--https` transport coercion |
| `pull-release-cd` | `prc` | Multi-repo one-shot: comma-separated `<repo> <version>` list; URL tokens auto-cloned first, then `pull-release <ver> -y` per entry |
| `ssh view` / `copy` / `create` | `v` / `cp` / — | Inspect, copy, or create the SSH public key (clipboard-aware: `clip` / `pbcopy` / `wl-copy`) |
| `desktop-sync` | `ds` | Sync tracked repos with GitHub Desktop |
| `install gitmap-oneliner` | — | Print the Win/Mac install one-liners for the current repo |

```bash
# clone from a structured file
gitmap clone json --target-dir ./restored
gitmap clone csv                                # auto-resolves to ./gitmap-output/gitmap.csv
gitmap clone ./gitmap-output/gitmap.json --safe-pull
gitmap clone ./gitmap-output/gitmap.json --github-desktop

# clone a single repo by URL (auto-flattens versioned URLs)
gitmap clone https://github.com/alimtvnetwork/gitmap-v28
gitmap clone https://github.com/alimtvnetwork/gitmap-v28 my-folder
gitmap clone git@github.com:alimtvnetwork/gitmap-v28.git my-folder
gitmap clone https://github.com/alimtvnetwork/gitmap-v28 --replace   # see spec 96

# force the transport — auto-converts HTTPS <-> SSH-shorthand before git runs
gitmap clone https://github.com/alimtvnetwork/gitmap-v28.git --ssh    # -> git@github.com:alimtvnetwork/gitmap-v28.git
gitmap clone git@github.com:alimtvnetwork/gitmap-v28.git --https      # -> https://github.com/alimtvnetwork/gitmap-v28.git
gitmap clone "https://github.com/a/x,https://github.com/a/y" --ssh    # batch — every URL converted

# clone-next: jump to the next (or specific) versioned sibling
gitmap cn v++                                   # my-app-v3 -> my-app-v4
gitmap cn v15 --delete                          # jump to v15, delete current
gitmap cn v++ --create-remote                   # create GitHub repo if missing
gitmap cn v++ --no-flatten                      # keep nested folder layout

# clone-fix-repo / clone-fix-repo-pub — one-shot pipelines (full clone-flag parity)
gitmap cfr  https://github.com/acme/myrepo-v13.git                # clone + fix-repo --all
gitmap cfr  https://github.com/acme/myrepo-v13.git --ssh          # coerce to SSH, then run pipeline
gitmap cfrp git@github.com:acme/myrepo-v13.git --https            # coerce to HTTPS, fix, then make-public --yes
gitmap cfrp https://github.com/acme/myrepo-v13.git myrepo-fresh   # explicit destination folder

# push / pull — transport-aware (rewrites + persists remote.origin.url)
gitmap push                                     # plain git push in cwd
gitmap ph --ssh                                 # rewrite origin to SSH, push, persist
gitmap pull --https                             # rewrite origin to HTTPS, pull, persist
gitmap push --ssh origin main                   # extra args forward verbatim

# pull-release-cd (prc) — release N repos in one shot
gitmap prc gitmap v5.51.0, marco v2.5.0, https://github.com/me/other v3.1.0

# ssh — inspect / copy / create the public key (auto-copies to clipboard)
gitmap ssh view                                 # print + copy ~/.ssh/id_*.pub
gitmap ssh copy                                 # copy without printing
gitmap ssh create my-alias                      # generate + copy a new keypair

# install gitmap-oneliner — print the Win/Mac install one-liners for the current repo
gitmap install gitmap-oneliner
```

→ [clone](gitmap/helptext/clone.md) · [clone-next](gitmap/helptext/clone-next.md) · [clone-fix-repo](gitmap/helptext/clone-fix-repo.md) · [clone-fix-repo-pub](gitmap/helptext/clone-fix-repo-pub.md) · [push](gitmap/helptext/push.md) · [pull](gitmap/helptext/pull.md) · [pull-release-cd](gitmap/helptext/pull-release-cd.md) · [ssh](gitmap/helptext/ssh.md) · [desktop-sync](gitmap/helptext/desktop-sync.md)

---

<div align="center">

### Scan & Clone — `--config`, `--mode`, `--output` recipes

</div>

Concrete, copy-pasteable examples for the three flags you'll reach for most.
Defaults are `--config ./data/config.json`, `--mode https`, and
`--output terminal`. Source of truth: [`gitmap/helptext/scan.md`](gitmap/helptext/scan.md)
and [`gitmap/helptext/clone.md`](gitmap/helptext/clone.md).

#### `--config <path>` — point scan at a non-default config

```bash
# default: reads ./data/config.json relative to the binary
gitmap scan ~/projects

# point at a project-local config (commit it alongside your repo list)
gitmap scan ~/projects --config ./gitmap.config.json

# CI: point at a profile that excludes vendored & node_modules trees
gitmap scan /workspace --config /etc/gitmap/ci-profile.json --quiet

# different config for a different drive on Windows
gitmap scan D:\wp-work --config D:\gitmap\configs\wp.json
```

The `--config` path is recorded in the scan cache, so a follow-up
`gitmap rescan` replays the exact same config without re-typing it.

> 📖 **Full key reference:** every JSON key gitmap reads from `data/config.json`
> — defaults, allowed values, and the nested `release` shape — is documented
> in [`docs/config-schema.md`](docs/config-schema.md).

#### `--mode ssh|https` — pick the clone-URL flavor recorded in output

```bash
# HTTPS (default) — works without keys, prompts for token on private repos
gitmap scan ~/projects --mode https
# → records: https://github.com/<owner>/<repo>.git

# SSH — uses your ~/.ssh keys, no token prompt, works for private repos
gitmap scan ~/projects --mode ssh
# → records: git@github.com:<owner>/<repo>.git

# scan once in HTTPS, then re-scan in SSH for a CI machine that has keys
gitmap scan ~/projects --mode https --output json --output-path ./out/https
gitmap scan ~/projects --mode ssh   --output json --output-path ./out/ssh
```

The mode only affects the **URL string written to the output files** —
your local working copies are not touched. Downstream `gitmap clone`
honours whichever URL the file contains, so the choice flows end-to-end.

#### `--output csv|json|terminal` — pick the artifact format

```bash
# terminal (default) — pretty-print to stdout, no files written
gitmap scan ~/projects

# csv — machine-readable spreadsheet (one row per repo)
gitmap scan ~/projects --output csv
# → writes ./.gitmap/output/gitmap.csv

# json — full structured payload (branch, remote, tags, last-commit SHA, ...)
gitmap scan ~/projects --output json
# → writes ./.gitmap/output/gitmap.json

# custom output directory (handy in CI artifact uploads)
gitmap scan ~/projects --output json --output-path ./build/scan-results
```

#### Combined recipes (the patterns you'll actually use)

```bash
# 1. Daily local sync: HTTPS + JSON, cached config
gitmap scan ~/projects --config ~/.gitmap/personal.json --mode https --output json

# 2. CI snapshot for SSH-keyed runners: SSH + CSV
gitmap scan /workspace --config /etc/gitmap/ci.json --mode ssh --output csv

# 3. Quick one-off sanity check (no files, no config tweaks)
gitmap scan . --output terminal

# 4. Scan once, then bulk-clone elsewhere using the SSH JSON manifest
gitmap scan ~/projects --mode ssh --output json --output-path ./manifest
gitmap clone ./manifest/gitmap.json --target-dir /opt/restored --safe-pull

# 5. CSV → another machine → clone everything via HTTPS
gitmap scan ~/projects --mode https --output csv --output-path ./share
# (copy ./share/gitmap.csv to the other host, then:)
gitmap clone ./gitmap.csv --target-dir ~/work --github-desktop
```

`gitmap clone` automatically picks the right input parser from the file
extension (`.json` / `.csv` / `.txt`) or the shorthand keywords `json` /
`csv` / `text`, so the `--output` format you chose at scan time is the
format `clone` will read on the other side.

#### `--no-vscode-sync` — opt out of VS Code Project Manager sync

By default, every `gitmap` clone variant upserts each successfully-cloned
repo into the alefragnani.project-manager `projects.json` so the new
folder appears in VS Code's Project Manager sidebar immediately. Pass
`--no-vscode-sync` on any clone command to skip that step (useful in CI,
headless servers, or environments without VS Code installed):

| Command | Sync target when default |
|---------|--------------------------|
| `gitmap clone <url>` | the single resolved folder |
| `gitmap clone <json\|csv\|text>` | every repo that landed on disk (one batched pass) |
| `gitmap clone-next` (`cn`) | the freshly-flattened folder; batch mode covers every repo in the batch |
| `gitmap clone-fix-repo` (`cfr`) / `clone-fix-repo-pub` (`cfrp`) | forwarded to the underlying clone step |
| `gitmap clone-from` (`cf`) | every successfully-cloned row (no-op during dry-run) |
| `gitmap clone-pick` (`cpk`) | the sparse-checkout destination |
| `gitmap reclone` / `clone-now` (`cnow`) | every successfully re-cloned repo (no-op without `--execute`) |

```bash
gitmap clone https://github.com/acme/widget.git --no-vscode-sync
gitmap cn v++ --no-vscode-sync
gitmap clone-from repos.csv --execute --no-vscode-sync
```

When sync is skipped gitmap prints
`• VS Code Project Manager sync skipped (--no-vscode-sync)`. When the
VS Code extension directory is missing the sync is also a soft no-op (a
warning is logged; exit code is unchanged).

##### Windows path canonicalization & EvalSymlinks soft-fail

Every `rootPath` written to `projects.json` is run through a single
canonicalization helper (`canonicalizePMPath` in `gitmap/cmd/clonepmsync.go`)
before it reaches disk. This is what stops the same physical clone target
from producing two distinct VS Code sidebar entries when reached through
different shells, drive mappings, or path spellings on Windows.

**Canonicalization steps** (in order):

1. **`filepath.Clean`** — collapses mixed `/` and `\` separators,
   removes redundant `.` segments and trailing separators. Without
   this, a manifest authored on Linux with `RelativePath: "acme/widget"`
   joined onto a Windows abs path leaks the forward-slash form into
   `projects.json`.
2. **`filepath.EvalSymlinks`** — resolves symlinks AND Windows 8.3
   short names (`C:\PROGRA~1\…`) to their canonical long form.
   Without this, a `gitmap clone` invoked from a `cmd.exe` that
   resolved a `Program Files` ancestor to its short name produces a
   second `projects.json` row distinct from the long-form row a
   PowerShell-launched run would produce.
3. **Case-insensitive dedup on Windows** (handled inside
   `vscodepm.normalizePath`) — `C:\Foo\Repo` and `c:\foo\repo`
   resolve to a single entry. POSIX systems remain case-sensitive.

**Soft-fail policy** for `EvalSymlinks`: when the resolver errors
(path not yet on disk, permission denied, broken symlink), the
helper falls back to the **cleaned absolute path** rather than
failing the clone. The reasoning: a `projects.json` entry with the
pre-symlink path is always preferable to a swallowed clone — VS
Code will simply open the path verbatim, and the next `gitmap scan`
pass will re-canonicalize once the symlink is resolvable.

**Manifest-mode `RelativePath` joins** apply the same defensive
normalization at the source: every join site routes through
`model.CleanRelativePath` (`gitmap/model/relativepath.go`), which
runs `filepath.Clean(filepath.FromSlash(rel))` so cross-platform
manifests produce identical `AbsolutePath` strings on Windows and
POSIX.

**Tracing**: pass `--debug-paths` to `gitmap clone` to emit one
`[debug-paths] in=… clean=… resolved=…` line on **stderr** for every
canonicalize call (works on both the resolved-OK and soft-fail
branches). Equivalent to setting `GITMAP_DEBUG_PATHS=1` directly,
which CI environments can use without a CLI flag. The trace surfaces
the exact rewrite at the boundary, making 8.3 short-name and
symlink-ancestor dedup issues self-diagnosing.

---

<div align="center">

### Move & Merge

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `mv` | `move` | Move LEFT into RIGHT, then delete LEFT |
| `merge-both` | — | Fill missing files on BOTH sides; prompt on conflicts |
| `merge-left` | — | Copy from RIGHT into LEFT; prompt on conflicts |
| `merge-right` | — | Copy from LEFT into RIGHT; prompt on conflicts |

Each side (LEFT / RIGHT) can be a local folder OR a remote URL.
URL endpoints are auto-cloned (or pulled if already cloned), and
the result is committed + pushed back when the operation completes.

```bash
# move: classic file copy + delete source
gitmap mv ./gitmap-v28 ./gitmap-v28
gitmap mv ./gitmap-v28 https://github.com/alimtvnetwork/gitmap-v28
gitmap mv https://github.com/alimtvnetwork/gitmap-v28 ./another-folder
gitmap mv https://github.com/alimtvnetwork/gitmap-v28 \
         https://github.com/alimtvnetwork/gitmap-v28

# merge-both: bidirectional fill (each side gains what the other has)
gitmap merge-both ./gitmap-v28 ./gitmap-v28
gitmap merge-both ./gitmap-v28 https://github.com/alimtvnetwork/gitmap-v28
gitmap merge-both https://github.com/alimtvnetwork/gitmap-v28 \
                  https://github.com/alimtvnetwork/gitmap-v28

# merge-left: take RIGHT into LEFT
gitmap merge-left ./gitmap-v28 ./gitmap-v28
gitmap merge-left ./local https://github.com/alimtvnetwork/gitmap-v28

# merge-right: take LEFT into RIGHT
gitmap merge-right ./gitmap-v28 ./gitmap-v28
gitmap merge-right ./local https://github.com/alimtvnetwork/gitmap-v28

# bypass conflict prompts: source-side wins by default
gitmap merge-right ./gitmap-v28 ./gitmap-v28 -y
gitmap merge-both  ./gitmap-v28 ./gitmap-v28 -y --prefer-newer

# pin remote branch + preview
gitmap merge-right ./local https://github.com/owner/repo:develop
gitmap mv ./gitmap-v28 ./gitmap-v28 --dry-run
```

Conflict prompt keys: **L**eft / **R**ight / **S**kip /
**A**ll-left / **B**all-right / **Q**uit. Pass `-y` (or `-a`) to
bypass; combine with `--prefer-left` / `--prefer-right` /
`--prefer-newer` / `--prefer-skip` to override the default policy.

→ [spec/01-app/97-move-and-merge.md](spec/01-app/97-move-and-merge.md)

---

<div align="center">

### Repo-wide Rewrite, One-shot Publish & Chronological Replay

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `replace` | `rpl` | Repo-wide find/replace across every text file. Literal swap or version-suffix bump driven by the `-vK` git remote URL. `--audit` reports without writing. |
| `fix-repo` | `fr` | Rewrite prior `{base}-vN` versioned-repo-name tokens to the current version. Also rewrites bare `{base}` tokens (no `-vN` suffix) when v1 is in the target span, for repos that shipped pre-versioned. Negative-lookahead guards `-v1` from matching `-v18`; bare-base sweep is word-boundary guarded so `{base}.js`, `{base}_alt`, `{base}-v2`, and `myimg-pdf` are left untouched. `--strict` runs `go test` on touched packages. **v5.51.0+:** every write is backed up to `.gitmap/backup/<repo>/v<current>/fix-repo/<UTC-ts>/` and is restorable with `gitmap undo`. |
| `undo` | `ud` | Restore the latest `gitmap fix-repo` snapshot for the current repo + version. `--list` lists snapshots, `--snapshot <ts>` selects one, `--dry-run` previews. |
| `clone-fix-repo` | `cfr` | One-shot: `clone <url>` then `fix-repo --all` inside the new folder. Versioned URLs auto-flatten. |
| `clone-fix-repo-pub` | `cfrp` | Same as `cfr`, plus `make-public --yes` at the end. |
| `make-public` | — | Make the current repo **public** on GitHub or GitLab via the matching CLI (`gh` / `glab`). Verifies visibility post-edit. |
| `commit-in` | `cin` | Walk one or more SOURCE repos chronologically and APPEND each commit into a TARGET repo, preserving both AuthorDate AND CommitterDate. Idempotent via `ShaMap`. |

```bash
# Generic find/replace (preview first)
gitmap replace "github.com/old-org" "github.com/new-org" --dry-run
gitmap rpl     "github.com/old-org" "github.com/new-org" -y

# Bump prior 3 versions of {base}-vN -> current vK
gitmap replace -3 -y
gitmap replace all --audit                # report-only

# Just the version bump (Go-native fix-repo, --strict runs `go test`)
gitmap fix-repo --all --strict
gitmap fr -5 --dry-run --verbose

# Before/after — scope depends on current version
# Inside gitmap-v28 (bare sweep ACTIVE on v1->v2)
#   gitmap       -> gitmap-v28
#   gitmap-v28    -> gitmap-v28
#   gitmap.js    -> gitmap.js     (word-boundary)
# Inside gitmap-v28+ (bare sweep SKIPPED)
#   gitmap                 -> gitmap                 (binary/brand preserved)
#   https://.../gitmap     -> https://.../gitmap     (URL preserved)
#   gitmap-v28, gitmap-v28   -> gitmap-v<current>
#   gitmap-v28             -> gitmap-v28             (negative-lookahead)
# Force-suppress the bare sweep even on v1->v2:
gitmap fix-repo -2 --restrict no-version
gitmap fr       -2 -r nv

# Undo the most recent fix-repo write (v5.51.0+)
gitmap undo --list                # show snapshots, newest first
gitmap undo                       # restore the latest snapshot
gitmap undo --snapshot 20260519T134210Z --dry-run

# Clone + fix in one shot (versioned URLs auto-flatten)
gitmap clone-fix-repo https://github.com/acme/myrepo-v13.git
gitmap cfr            git@github.com:acme/myrepo-v13.git myrepo-fresh

# Clone + fix + publish public in one shot
gitmap cfrp https://github.com/acme/myrepo-v13.git

# Flip current repo public on GitHub or GitLab
gitmap make-public                          # interactive
gitmap make-public --yes                    # CI / scripts
gitmap make-public --dry-run --verbose      # preview the gh/glab call

# v5.63.0+ — bulk visibility (N most recent versions of a base repo)
gitmap make-public 3                                # flip current + 2 prior versions
gitmap make-public macro-ahk-v40 5                  # flip v40..v36 of macro-ahk
gitmap make-private https://github.com/me/foo-v9 2  # privatize v9 + v8

# v5.63.0+ — cfrp now offers to privatize prior public versions
gitmap cfrp https://github.com/me/foo-v40           # prompts after publishing v40
gitmap cfrp https://github.com/me/foo-v40 -y        # auto-privatize, no prompts

# Stitch every versioned sibling's history into one canonical repo
gitmap commit-in ./canonical all --save-profile Default --set-default
gitmap cin       ./canonical -3 --dry-run --function-intel on --languages Go,TypeScript
```

→ Specs: [`spec/04-generic-cli/15-replace-command.md`](spec/04-generic-cli/15-replace-command.md) ·
[`spec/04-generic-cli/27-fix-repo-command.md`](spec/04-generic-cli/27-fix-repo-command.md) ·
[`spec/03-commit-in/`](spec/03-commit-in/)

---

<div align="center">

### Git Operations

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `pull` | `p` | Pull a specific repo by name |
| `exec` | `x` | Run git command across all repos |
| `status` | `st` | Show repo status dashboard |
| `watch` | `w` | Live-refresh repo status dashboard |
| `has-any-updates` | `hau` | Check if remote has new commits |
| `latest-branch` | `lb` | Find most recently updated remote branch |

```bash
gitmap pull --group work --all
gitmap exec fetch --prune
gitmap watch --interval 10 --group work
gitmap lb 5 --format csv
```

→ [pull](gitmap/helptext/pull.md) · [exec](gitmap/helptext/exec.md) · [status](gitmap/helptext/status.md) · [watch](gitmap/helptext/watch.md) · [latest-branch](gitmap/helptext/latest-branch.md)

---

<div align="center">

### Navigation & Organization

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `cd` | `go` | Navigate to a tracked repo directory |
| `group` | `g` | Manage repo groups / activate for batch ops |
| `multi-group` | `mg` | Select multiple groups for batch operations |
| `alias` | `a` | Assign short names to repos |
| `as` | `s-alias` | Register the current Git repo + name in one shot (run from inside the repo) |
| `diff-profiles` | `dp` | Compare repos across two profiles |

```bash
gitmap cd my-api
gitmap g work && gitmap g pull
gitmap mg backend,frontend && gitmap mg status
gitmap alias set api github/user/api-gateway
gitmap as backend           # registers the current repo as 'backend' + adds it to the DB
gitmap as                   # uses the folder basename as the alias
gitmap alias suggest --apply
```

→ [cd](gitmap/helptext/cd.md) · [group](gitmap/helptext/group.md) · [multi-group](gitmap/helptext/multi-group.md) · [alias](gitmap/helptext/alias.md) · [as](gitmap/helptext/as.md) · [diff-profiles](gitmap/helptext/diff-profiles.md)

---

### Right-click context menu (`gitmap install ctx`)

One command wires the most-used gitmap actions into your file manager's
right-click menu — Windows Explorer (HKCU registry cascade), macOS
Finder (Automator Quick Actions in `~/Library/Services/`), and Linux
Nautilus / Dolphin / Thunar. Right-click any folder, pick the action,
and gitmap runs against that folder as cwd.

```bash
gitmap install ctx              # add the cascade
gitmap install ctx --explain    # bake an `> gitmap <args>` echo into every entry,
                                # so each click prints its exact resolved invocation
                                # before executing (terminal echo for Terminal mode,
                                # prepended to the notification text for Silent mode)
gitmap uninstall ctx            # remove every entry written by ctx (safe, idempotent)
```

The table below is the **authoritative** mapping between menu item and
the exact command + flags gitmap runs. It mirrors
[`spec/04-generic-cli/30-install-ctx.md`](spec/04-generic-cli/30-install-ctx.md)
§3 and the Go source of truth in
[`gitmap/cmd/installctxentries.go`](gitmap/cmd/installctxentries.go).

| Menu path                       | Command run                          | Mode     | Notes |
|---------------------------------|--------------------------------------|----------|-------|
| Scan ▸ Scan here                | `gitmap scan`                        | Terminal | Walks the clicked folder; honours all gitmap scan rules. |
| Scan ▸ Rescan                   | `gitmap rescan`                      | Terminal | Re-runs the most recent scan against the same root. |
| Scan ▸ Find next                | `gitmap find-next`                   | Silent   | Probes for the next `<base>-vN+1` sibling. **No `--scan-folder`** (uses cwd); **no `--json`** (output goes to the OS notification verbatim). |
| Clone ▸ Clone-next here         | `gitmap clone-next`                  | Terminal | Flattens to the base-name folder (v2.75.0+ default); version history tracked in `RepoVersionHistory`. |
| Clone ▸ Pull                    | `gitmap pull`                        | Terminal | Fast-forward pull on the current repo only. |
| Clone ▸ Pull all (every tracked repo) | `gitmap pull-all` (= `pull --all`) | Terminal | **Power-user batch.** On Windows it's hidden behind **Shift+right-click** (HKCU `Extended` REG_SZ). On macOS/Linux there's no Shift filter, so the entry is visible but gated by an OS confirm dialog (`osascript display dialog` / `zenity`/`kdialog`/`xmessage`) before it fans out. |
| Release ▸ Release current       | `gitmap release`                     | Terminal | Auto-bump-minor prompt; reads `latest.json` first, falls back to git tags. |
| Release ▸ Release next          | `gitmap release --bump minor`        | Terminal | Hard-pinned to `--bump minor`. `--bump major` and `--bump patch` are deliberately not surfaced (one-click should never ship a breaking change). |
| Release ▸ Pull release          | `gitmap pull-release`                | Terminal | `git pull` then `release`. Accepts `pr`, legacy `release-pull` / `relp` / `rlp`. |
| Release ▸ Release pending       | `gitmap release-pending`             | Silent   | Reports unpublished commits since the last pinned tag. |
| Release ▸ List releases         | `gitmap list-releases`               | Silent   | Single-repo view. **No `--all-repos`** — clicking inside one folder shouldn't surprise the user with a catalog-wide list. |
| Release ▸ List versions         | `gitmap list-versions`               | Silent   | Local + remote version siblings. |
| Repos ▸ Go / Node / React / C++ / C# | `gitmap go-repos` etc.          | Silent   | Filtered catalog views, surfaced via notification. |
| Visibility ▸ Make public/private | `gitmap make-public` / `make-private` | Terminal | Confirmation prompts run in the visible terminal. |
| Tools ▸ Fix repo                | `gitmap fix-repo`                    | Terminal | **No `--strict`** — the menu is for the fast one-shot fix; CI / scripted use is where `--strict` belongs. |
| Tools ▸ Diff / History / Update | `gitmap diff` / `history` / `update` | Terminal | gitmap-managed wrappers. |
| Git ▸ History (git log graph)   | `git log --oneline --graph --decorate --all -n 100` | Terminal | Raw git, bypasses gitmap. |
| Git ▸ Diff (git diff)           | `git diff --stat HEAD`               | Terminal | Raw git. |
| Git ▸ Log (git log)             | `git log -n 30 --pretty=format:%h %ad %s --date=short` | Terminal | Raw git. |
| Git ▸ Status (git status)       | `git status`                         | Silent   | Raw git, into notification. |
| Open terminal here              | (prefilled `gitmap ` prompt)         | Prefill  | Opens a terminal at the folder with `gitmap ` typed but not run. |
| Docs                            | `gitmap docs`                        | Silent   | Opens the docs URL via OS notification handler. |

**Execution modes**

- **Terminal** — opens a new terminal window that stays open after the
  command exits, so you can read multi-line output (`pwsh -NoExit` on
  Windows, Terminal.app via `osascript` on macOS, `x-terminal-emulator`
  on Linux).
- **Silent** — runs hidden; first ~200 chars of stdout/stderr surface
  via `msg.exe` (Windows), `display notification` (macOS), or
  `notify-send` with `echo` fallback (Linux).
- **Prefill** — opens a terminal and writes `gitmap ` at the prompt
  without pressing Enter — for ad-hoc commands.

---

<div align="center">

## 🚀 gitmap Release

**Current version:** `v3.50.0` · Cross-platform (Windows · Linux · macOS) · Single static binary

</div>

The `gitmap release` command turns a clean working tree into a versioned,
tagged, multi-target GitHub release in one step — branch + tag + push +
binary build + checksum + changelog body + GitHub Release page.

### Install or upgrade gitmap

#### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex
```

#### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | sh
```

#### Pin to an exact version

```powershell
# Windows — install v3.50.0 exactly, skip the "latest" lookup
$ver = 'v3.50.0'
$installer = irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1
& ([scriptblock]::Create($installer)) -Version $ver -NoDiscovery
```

```bash
# Linux / macOS — install v3.50.0 exactly
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh \
  | bash -s -- --version v3.50.0 --no-discovery
```

Verify:

```bash
gitmap --version
# gitmap v3.50.0
```

### Release CLI examples

| Goal | Command |
|------|---------|
| Auto-bump minor (`3.50.0 → 3.51.0`) and release | `gitmap r` |
| Bump patch and release | `gitmap release --bump patch` |
| Bump minor with binary + zip + checksums | `gitmap release --bump minor --bin --compress --checksums` |
| Release a specific version with notes | `gitmap release v3.51.0 -N "Performance pass"` |
| Release any registered repo from anywhere | `gitmap ra my-api v1.4.0` |
| Pull, then release a registered repo | `gitmap rap my-api v1.4.0` |
| Release gitmap itself from any directory | `gitmap release-self --bump patch` |
| Multi-repo: bump every repo under cwd | `gitmap r -y` (run from a folder of repos) |
| Preview without pushing | `gitmap release --bump patch --dry-run` |
| List build targets that will be produced | `gitmap release --list-targets` |

```bash
# End-to-end: register once, release from anywhere
cd ~/code/my-api
gitmap as my-api                  # one-time alias
cd ~                              # go anywhere
gitmap ra my-api v1.4.0 --pull    # pull --ff-only, release, push, build, attach
```

> **Auto-stash safety:** dirty trees are stashed before `release-alias`
> with a label like `my-api-1.4.0-1715000000` and popped on exit. Pass
> `--no-stash` to abort instead, or `--dry-run` to preview every step.

### What you get — output formats

Every successful release produces all of the following:

| Artifact | Where | Format |
|----------|-------|--------|
| Release branch | local + `origin` | `release/v3.51.0` |
| Annotated tag | local + `origin` | `v3.51.0` |
| Local manifest | `.gitmap/release/latest.json` | JSON (version, tag, branch) |
| GitHub Release page | github.com/&lt;owner&gt;/&lt;repo&gt;/releases | Title + body + assets |
| Release body | GitHub Release page | Markdown — CHANGELOG section + pinned-install snippet |
| Cross-compiled binaries (with `--bin`) | uploaded as release assets | `gitmap_v3.51.0_<os>_<arch>[.exe]` for 6 targets |
| Compressed archives (with `--compress`) | uploaded as release assets | `.zip` (Windows) / `.tar.gz` (Linux/macOS) |
| Checksums (with `--checksums`) | uploaded as release asset | `SHA256SUMS.txt` |
| Custom zip groups (with `--zip-group`) | uploaded as release assets | `.zip` bundles per group |

**Default build targets** (override with `--targets` or `release.targets` in `config.json`):

```
linux/amd64    linux/arm64
darwin/amd64   darwin/arm64
windows/amd64  windows/arm64
```

**Sample terminal output:**

```
  Creating release v3.51.0...
  ✓ Created branch release/v3.51.0
  ✓ Created tag v3.51.0
  ✓ Pushed branch and tag to origin
  ✓ Release metadata written to .gitmap/release/latest.json
  ✓ Committed release metadata on release/v3.51.0
  ✓ Marked v3.51.0 as latest release
  ✓ Using CHANGELOG.md as release body
  ✓ Attached gitmap_v3.51.0_windows_amd64.exe
  ✓ Attached gitmap_v3.51.0_linux_amd64
  ✓ Attached SHA256SUMS.txt

  ── Release v3.51.0 complete ──
```

### Configuration file

| Property | Value |
|----------|-------|
| **Default path** | `./data/config.json` (resolved relative to the gitmap binary) |
| **Override** | Edit the file directly — there is no `--config` flag; CLI flags always win over config values |
| **Format** | JSON, loaded once per command via `config.LoadFromFile` |
| **Missing file** | Silently falls back to built-in defaults (no error) |

**Full schema** (every field is optional; omit a field to use its default):

```json
{
  "defaultMode":      "https",
  "defaultOutput":    "terminal",
  "outputDir":        "./output",
  "excludeDirs":      ["node_modules", ".git"],
  "notes":            "",
  "dashboardRefresh": 30,
  "release": {
    "targets": [
      { "goos": "windows", "goarch": "amd64" },
      { "goos": "linux",   "goarch": "amd64" },
      { "goos": "darwin",  "goarch": "arm64" }
    ],
    "checksums": true,
    "compress":  true
  }
}
```

**Field meanings (release-relevant only):**

| Field | Type | Default | Effect |
|-------|------|---------|--------|
| `release.targets[]` | array of `{goos, goarch}` | built-in 6-target matrix | Cross-compile matrix used when `--bin` is set. `--targets` flag overrides this entirely. |
| `release.targets[].goos` | string | — | Go `GOOS` value: `windows`, `linux`, `darwin`, `freebsd`, … |
| `release.targets[].goarch` | string | — | Go `GOARCH` value: `amd64`, `arm64`, `386`, … |
| `release.checksums` | bool | `false` | Always emit `SHA256SUMS.txt`. Equivalent to passing `--checksums` on every release. |
| `release.compress` | bool | `false` | Always wrap assets in `.zip`/`.tar.gz`. Equivalent to passing `--compress` on every release. |
| `outputDir` | string | `./output` | Where non-release CLI exports land (scan reports, etc). Not used by `release` directly. |
| `excludeDirs` | array | `[]` | Folders to skip during scanning. Not used by `release` directly. |

**Resolution order** (last writer wins):

```
built-in defaults  <  data/config.json  <  CLI flags
```

So `release.compress: false` in config + `--compress` on the CLI → compression ON for that one run.

### CLI flag reference — `gitmap release` / `r`

```
gitmap release [version] [flags]
gitmap r       [version] [flags]
```

`version` is positional and optional. Forms accepted:

- `v3.51.0` — release exactly this tag
- *(omitted)* — auto-bump **minor** from the last release in `.gitmap/release/latest.json`, with a `[y/N]` prompt (skip with `-y`)
- combined with `--bump` is an error (mutually exclusive)

| Flag | Type | Default | Meaning |
|------|------|---------|---------|
| `--bump <level>` | string | (none) | Auto-increment version segment. Accepts `major`, `minor`, or `patch`. Mutually exclusive with positional `version`. |
| `-N`, `--notes <text>` | string | git commit subject | Release title / notes used as the GitHub Release body header. |
| `--commit <sha>` | string | `HEAD` | Create the release from a specific commit instead of `HEAD`. Mutually exclusive with `--branch`. |
| `--branch <name>` | string | current branch | Create the release from the latest commit of `<name>`. Mutually exclusive with `--commit`. |
| `--assets <path>` | string | (none) | Single file or directory to attach as release assets (in addition to `--bin`, `-Z`, `--zip-group`). |
| `-b`, `--bin` | bool | `false` | Cross-compile Go binaries for every target in the matrix and attach them as release assets. |
| `--targets <list>` | string | from config / built-ins | Comma-separated `goos/goarch` pairs (e.g. `windows/amd64,linux/arm64`). Overrides `release.targets` in config. |
| `--list-targets` | bool | `false` | Resolve the target matrix, print it, exit 0. No release is created. Useful for verifying config. |
| `--compress` | bool | `false` (or config) | Wrap each binary in a per-target archive: `.zip` for Windows, `.tar.gz` for Linux/macOS. |
| `--checksums` | bool | `false` (or config) | Emit `SHA256SUMS.txt` covering every uploaded asset and attach it to the release. |
| `-Z <path>` | repeatable | (none) | Ad-hoc zip: include a single file or folder as a release asset. May be passed multiple times. |
| `--zip-group <name>` | repeatable | (none) | Attach a persistent named **zip group** (defined via `gitmap zg add`) as a release asset. May be passed multiple times. |
| `--bundle <name>` | string | (none) | Combine all `-Z` items into one archive named `<name>.zip` instead of one archive per item. |
| `--draft` | bool | `false` | Create the GitHub Release as an unpublished draft. Branch/tag are still pushed. |
| `--dry-run` | bool | `false` | Print every step that would run (branch, tag, push, upload) without touching git or GitHub. |
| `--no-commit` | bool | `false` | Skip the post-release auto-commit + push of `.gitmap/release/latest.json`. |
| `-y`, `--yes` | bool | `false` | Auto-confirm every prompt: bare-release auto-bump, multi-repo scan, orphaned-metadata cleanup. |
| `--verbose` | bool | `false` | Write detailed stdout/stderr trace to a timestamped log file under `data/logs/`. |

**Mutually exclusive combinations** (gitmap will exit with an error):

| Combination | Why |
|-------------|-----|
| `<version>` **and** `--bump` | Either explicit or auto-bump — pick one. |
| `--commit` **and** `--branch` | A release has exactly one source commit. |

→ Detailed help: [release](gitmap/helptext/release.md) · [release-alias](gitmap/helptext/release-alias.md) · [release-self](gitmap/helptext/release-self.md) · [release-pending](gitmap/helptext/release-pending.md) · [changelog](gitmap/helptext/changelog.md)

### Changelog (recent versions)

Concise, grouped per version. Each entry calls out **💥 Breaking**, **✨ Enhancements**, and **🐛 Fixes**. Versions with nothing in a category omit it. Full history lives in [`CHANGELOG.md`](CHANGELOG.md); query it from the CLI with `gitmap changelog vX.Y.Z` or `gitmap cl --limit 5`.

#### v5.10.0 — 2026-05-16 — PowerShell wrapper loads last

- 🐛 **Fixes:** `gitmap setup` / installer rewrites now move the managed `gitmap`/`gcd` PowerShell command wrapper to the end of `$PROFILE`. This stops older `# gitmap shell wrapper v2` blocks or hand-edited snippets appearing later in the profile from overriding the fresh wrapper and leaving `gitmap cd` stuck on the raw-exe warning.

#### v5.9.0 — 2026-05-16 — Harden PowerShell `gitmap cd` activation

- 🐛 **Fixes:**
  - Windows installers now write the `gitmap`/`gcd` command wrapper to **all** standard current-user PowerShell profile files (Windows PowerShell + PowerShell 7+), not only the one profile visible to the installer host.
  - Installer also loads the wrapper into the current session even when PATH was already present; release ZIP shim supports the same `GITMAP_HANDOFF_FILE` handoff as the generated shim.

#### v5.8.0 — 2026-05-16 — `fix-repo` bare-base rewrite for pre-versioned v1 repos

- ✨ **Enhancements:** `gitmap fix-repo` now rewrites bare `{base}` occurrences (not just `{base}-v1`) when v1 is in the target span. Pre-versioned remotes shipped without a `-v1` suffix, so downstream references read e.g. `img-pdf` rather than `img-pdf-v1`; the previous rewriter skipped them entirely and reported `changed: 0`. The bare-base sweep is guarded by strict word-boundary checks (alnum / `_` / `-` / `.`) so `{base}-v2`, `{base}.js`, `{base}_alt`, and `myimg-pdf` are left untouched. Skipped automatically when v1 is not in the target set.

#### v5.7.0 — 2026-05-16 — Ship PowerShell shim in release installs

- 🐛 **Fixes:** Windows release ZIPs now contain both `gitmap.exe` **and** `gitmap.ps1`, and the release-specific `install.ps1` moves the shim beside the exe. Users installing from release assets previously never received `gitmap.ps1`, so PowerShell resolved the raw exe and printed the wrapper-not-active warning even after upgrading.

#### v3.52.0 — 2026-04-21 — CI lint baseline cache controls (docs)

- ✨ **Enhancements:** `spec/09-pipeline/01-ci-pipeline.md` now documents the two `workflow_dispatch` inputs (`lint_baseline_cache_version`, `lint_baseline_disable`) that let operators rotate or bypass the golangci-lint baseline cache without editing the workflow. Includes copy-paste `gh workflow run` examples and a new "Job: Lint Baseline Diff" section covering cache keys, seeding mode, and sticky PR comment behavior.
- 🐛 **Fixes:** none — documentation-only sync; no CI behavior change.

#### v3.51.0 — 2026-04-21 — `cn v+1 -f` flag parsing + cleaner release trailer

- ✨ **Enhancements:**
  - Three-stage progress layout (Prepare → Clone → Finalize) now shown when `gitmap cn -f` is used; default `cn` output is unchanged.
  - `MsgForceReleasing` rewritten to plainly describe the Windows file-lock release ("Stepping out of … to release the file lock").
- 🐛 **Fixes:**
  - `gitmap cn v+1 -f` no longer silently dropped the `-f` flag when it followed a positional version arg. Fixed via `reorderFlagsBeforeArgs(args)` in `gitmap/cmd/clonenextflags.go` and an updated value-flag map in `gitmap/cmd/releaseargs.go` (covers `--csv`, `--ssh-key`, `-K`, `--target-dir`).
  - `Force` now implies `Keep` for the prior-folder cleanup, suppressing the redundant "Remove current folder?" prompt.
  - `MsgInstallHintUnix` gained a trailing blank line so the post-release shell prompt no longer sits flush against the `curl … | sh` line.

#### v3.50.0 — 2026-04-21 — Force-flatten for `clone-next`

- 💥 **Breaking:** none. New flag is opt-in and defaults to existing behavior.
- ✨ **Enhancements:**
  - `gitmap cn -f` / `--force`: force a flat clone even when cwd IS the target folder. Chdirs to the parent before remove (releases the Windows file lock), then re-clones into `<base>/`.
  - Refuses the silent versioned-folder fallback under `-f` — you get a flat layout or a clear error, never a surprise rename.
  - `--force` / `-f` added to zsh + PowerShell completions and to `clone-next` help text.
- 🐛 **Fixes:** `MsgFlattenLockedHint` now mentions `-f` so users discover the escape hatch on the first lock warning.

#### v3.32.1 — 2026-04-20 — `gitmap status` legacy path fix

- 🐛 **Fixes:** `status` no longer fails with `could not load gitmap.json at output\gitmap.json`. It now reads from the unified `.gitmap/output/` path and transparently falls through to the SQLite database when the JSON file is missing.

#### v3.32.0 — 2026-04-20 — Scan output: hoisted base path

- ✨ **Enhancements:** `gitmap scan` post-scan summary prints `📂 Base: <path>` once and lists each artifact by filename only. All icon/label columns aligned to a 12-char gutter.

#### v3.31.0 — 2026-04-20 — Cross-dir release/clone-next + `has-change`

- ✨ **Enhancements:**
  - `gitmap r <repo> <vX.Y.Z>` and `gitmap cn <repo> <vX.Y.Z>` — run from anywhere; gitmap chdirs in, fetches/pulls (rebase), auto-stashes, releases, then chdirs back and pops. Single-arg forms unchanged.
  - New `gitmap has-change` (`hc`) command: prints `true`/`false` per dimension (`dirty`/`ahead`/`behind`) or `--all` for all three; `--fetch=false` for offline use.
- 🐛 **Fixes:** `gitmap ssh` no longer exits 1 when `~/.ssh/id_rsa` already exists outside gitmap's DB — disk check moved before DB check; `--force` backs up to `id_rsa.bak.<unix-ts>` first.

#### v3.30.0 — 2026-04-20 — Go Report Card badge URL

- 🐛 **Fixes:** README badge now points at `goreportcard.com/badge/github.com/alimtvnetwork/gitmap-v28/gitmap` (real module path) instead of the repo root, which 404'd because there is no `go.mod` at the root.

#### v3.28.0 — 2026-04-20 — Lucrative scan summary

- ✨ **Enhancements:** `gitmap scan` summary regrouped into three labeled sections (`📦 Output Artifacts`, `🗄️ Database`, post-scan log) with category icons per row.

#### v3.27.0 — 2026-04-20 — Real Go module path

- 💥 **Breaking:** `go.mod` module path renamed from a placeholder to `github.com/alimtvnetwork/gitmap-v28/gitmap`. Anyone importing the module by the old path must update their import lines. CLI users are unaffected.

#### v3.26.0 — 2026-04-20 — Constants collision audit + CI guard

- ✨ **Enhancements:** new CI guard rejects PRs that introduce duplicate `Cmd*` / `Msg*` / `Err*` identifiers across `gitmap/constants/`. Backfilled audit caught 0 collisions on `main`.

#### v3.25.0 — 2026-04-20 — `github-desktop` (`gd`) command

- ✨ **Enhancements:** `gitmap gd` registers cwd repo with GitHub Desktop without running a full `scan` first. Pairs well with `clone-next`.

#### v3.24.0 — 2026-04-20 — Quiet release output

- 🐛 **Fixes:** suppressed cosmetic `LF will be replaced by CRLF` warnings during the release pipeline (kept underlying behavior identical; only stderr noise is muted).

#### v3.22.0 — 2026-04-20 — Auto-register on release

- ✨ **Enhancements:** `gitmap r` auto-registers the cwd repo in the database if it isn't tracked yet, instead of failing with "repo not found". The new repo is tagged with the current scan folder.

> Versions older than v3.22 are summarized in [`CHANGELOG.md`](CHANGELOG.md). Notable jumps: **v3.21** (schema-version fast path + `db-migrate --force`), **v3.19** (bare release auto-bumps **minor** + multi-repo scan-dir release), **v3.17** (`Release.RepoId` foreign key + doctor duplicate-binary check), **v3.16** (repo renamed to `gitmap-v28`).

### Copy-paste workflows — scan, output, re-clone

End-to-end recipes for the data pipeline that feeds a release: discover repos with `scan`, capture the result as CSV/JSON, and rebuild the same set on another machine with `clone`. Every block below is a single copy-paste — no placeholders to edit unless wrapped in `<…>`.

#### Scan with HTTPS clone URLs (default)

```bash
# Scan the current directory tree, terminal output only
gitmap scan

# Scan a specific folder
gitmap scan D:\wp-work

# Quiet mode (skip the post-scan clone-help section)
gitmap scan D:\wp-work --quiet
```

#### Scan with SSH clone URLs

```bash
# SSH-style URLs (git@github.com:owner/repo.git) instead of HTTPS
gitmap scan D:\wp-work --mode ssh

# SSH + auto-register every repo with GitHub Desktop
gitmap scan D:\wp-work --mode ssh --github-desktop
```

#### Output to CSV

```bash
# Write CSV to the default location: ./.gitmap/output/gitmap.csv
gitmap scan D:\wp-work --output csv

# CSV + custom output directory
gitmap scan D:\wp-work --output csv --output-path ./reports

# CSV + SSH URLs in one go
gitmap scan D:\wp-work --output csv --mode ssh --output-path ./reports
```

#### Output to JSON

```bash
# Write JSON to the default location: ./.gitmap/output/gitmap.json
gitmap scan D:\wp-work --output json

# JSON + custom directory + open the folder when done
gitmap scan D:\wp-work --output json --output-path ./reports --open

# JSON + SSH URLs (handy for piping into another tool)
gitmap scan D:\wp-work --output json --mode ssh
```

> Every `gitmap scan` run **also** writes the standard artifact bundle to `./.gitmap/output/` regardless of `--output`: `gitmap.csv`, `gitmap.json`, `gitmap.txt`, `folder-structure.md`, `clone.ps1`, `direct-clone.ps1`, `direct-clone-ssh.ps1`, `register-desktop.ps1`, `last-scan.json`. The `--output` flag controls the **terminal** representation only.

#### Re-clone everything from a scan file

```bash
# Re-clone from the JSON file produced by a previous scan
gitmap clone ./.gitmap/output/gitmap.json

# Re-clone from CSV
gitmap clone ./.gitmap/output/gitmap.csv

# Re-clone from a plain text list
gitmap clone ./.gitmap/output/gitmap.txt

# Re-clone into a specific base directory
gitmap clone ./.gitmap/output/gitmap.json --target-dir D:\restored

# Re-clone + safe pull on existing repos (retries + diagnostics)
gitmap clone ./.gitmap/output/gitmap.json --target-dir D:\restored --safe-pull

# Re-clone + auto-register everything with GitHub Desktop (no prompt)
gitmap clone ./.gitmap/output/gitmap.json --target-dir D:\restored --github-desktop
```

##### Accepted input schema

`gitmap clone <file>` reads three formats. The dispatcher picks one
from the file extension (`.json`, `.csv`, `.txt`); pass `--format
json|csv|text` to override. JSON and CSV inputs are validated
**before** cloning so a typo in a field name fails loudly with the
exact row number instead of silently dropping repos.

**JSON** — top-level array of objects. Every object key must be one
of the fields below; unknown fields are rejected with the offending
row number (1-based) and the full list of accepted fields. Each row
must carry at least one of `httpsUrl` or `sshUrl`.

| Field              | Type    | Required           | Notes |
|--------------------|---------|--------------------|-------|
| `id`               | number  | no                 | Scan-internal row id; ignored by clone. |
| `slug`             | string  | no                 | Stable repo slug (e.g. `owner-repo`). |
| `repoId`           | string  | no                 | Provider-side repo id. |
| `repoName`         | string  | no                 | Display name; falls back to URL basename. |
| `httpsUrl`         | string  | one of these two   | HTTPS clone URL. |
| `sshUrl`           | string  | one of these two   | SSH clone URL (`ssh://` or `git@host:path`). |
| `discoveredUrl`    | string  | no                 | URL as originally seen on disk. |
| `branch`           | string  | no                 | Branch to check out after clone. |
| `branchSource`     | string  | no                 | How the branch was determined (`HEAD`, `config`, ...). |
| `relativePath`     | string  | recommended        | Destination folder, relative to `--target-dir`. |
| `absolutePath`     | string  | no                 | Original absolute path on the source machine. |
| `cloneInstruction` | string  | no                 | Pre-rendered `git clone ...` command (informational). |
| `notes`            | string  | no                 | Free-form notes carried through from scan. |
| `depth`            | number  | no                 | Shallow-clone depth; `0` means full history. |
| `transport`        | string  | no                 | `ssh` / `https` / `other`; emitted by scan, ignored by clone. |

**CSV** — first row is the header. Header column names use the
exact same identifiers as the JSON fields above (case-sensitive).
The header normalizer tolerates a UTF-8 BOM, surrounding
double-quotes, and leading/trailing whitespace, so files saved by
Excel or PowerShell `Out-File -Encoding utf8` work without manual
cleanup. The header must include at least one of `httpsUrl` or
`sshUrl`. Per-data-row errors report the **1-based data row number**
(header is row 0; first data row is row 1).

**Text** — one URL per line, blank lines and `#` comments allowed.
No schema; the URL basename becomes the destination folder.

##### HTTPS example (JSON)

```json
[
  {
    "repoName": "wp-onboarding",
    "httpsUrl": "https://github.com/alimtvnetwork/wp-onboarding.git",
    "branch": "main",
    "relativePath": "wp-onboarding"
  },
  {
    "repoName": "gitmap-v28",
    "httpsUrl": "https://github.com/alimtvnetwork/gitmap-v28.git",
    "relativePath": "gitmap-v28",
    "depth": 1
  }
]
```

```bash
gitmap clone repos.json --target-dir D:\restored
```

##### SSH example (CSV)

```csv
repoName,sshUrl,branch,relativePath
wp-onboarding,git@github.com:alimtvnetwork/wp-onboarding.git,main,wp-onboarding
gitmap-v28,git@github.com:alimtvnetwork/gitmap-v28.git,,gitmap-v28
```

```bash
gitmap clone repos.csv --target-dir D:\restored
```

##### Mixed HTTPS + SSH (JSON)

`httpsUrl` and `sshUrl` may both appear on the same row; the
`--mode https|ssh` flag (default `https`) decides which one is used
to clone. Rows that lack the chosen mode's URL fall back to the
other URL automatically, so a mixed manifest "just works" under
either mode.

```json
[
  {
    "repoName": "public-repo",
    "httpsUrl": "https://github.com/acme/public.git",
    "sshUrl":   "git@github.com:acme/public.git",
    "relativePath": "public-repo"
  },
  {
    "repoName": "ssh-only",
    "sshUrl":   "git@gitlab.example.com:team/private.git",
    "relativePath": "ssh-only"
  }
]
```

#### One-shot single repo (no scan file needed)

```bash
# Clone a single URL — versioned URLs (e.g. -v13) auto-flatten to <base>/
gitmap clone https://github.com/alimtvnetwork/wp-onboarding-v13.git

# Clone into a custom folder name (skips auto-flatten)
gitmap clone https://github.com/alimtvnetwork/wp-onboarding-v13.git my-onboarding

# Clone via SSH
gitmap clone git@github.com:alimtvnetwork/wp-onboarding-v13.git
```

#### Full backup-and-restore round-trip

```bash
# === On the source machine ===
gitmap scan D:\wp-work --mode ssh --output json
# → produces D:\wp-work\.gitmap\output\gitmap.json (+ all sibling artifacts)

# Copy the file to the new machine, then:

# === On the target machine ===
gitmap clone gitmap.json --target-dir D:\wp-work --github-desktop --safe-pull
# → re-clones every repo via SSH, registers each with GitHub Desktop,
#   and pulls if any of them already exist on disk.
```

##### Clone-from report schema (JSON & CSV output)

After every `gitmap clone` run on a structured input file, a per-row
report is written to `./.gitmap/clone-from-report-<unixts>.<ext>`.
Both formats carry the **same field set**, sourced 1:1 from
`clonefrom.Result` (see `gitmap/clonefrom/execute.go` and
`gitmap/clonefrom/summary.go`). The schema is pinned by
`constants.CloneFromReportSchemaVersion` and guarded by
`TestCloneFromReportJSON_SchemaVersion_Pinned`.

**JSON envelope** (`clone-from-report-<unixts>.json`):

```json
{
  "schemaVersion": 1,
  "transport": { "ssh": 0, "https": 0, "other": 0 },
  "rows": [ { "url": "...", "dest": "...", "...": "..." } ]
}
```

`rows` is always a JSON array (never `null`) and `transport` is
always emitted, even when all counters are zero, so downstream
parsers can treat the envelope shape as unconditional.

**CSV report** (`clone-from-report-<unixts>.csv`): UTF-8, CRLF line
endings (matches every other gitmap CSV — see
`csvcrlf_contract_test.go`), one header row + one row per result.
Column order matches the JSON field order below.

| Field (JSON / CSV column) | Type    | Source on `clonefrom.Result`                       | Populated by                                        |
|---------------------------|---------|----------------------------------------------------|-----------------------------------------------------|
| `url`                     | string  | `Result.Row.URL`                                   | `parse.go` / `parsecsv.go` (verbatim from input).   |
| `dest`                    | string  | `Result.Dest`                                      | `executeRow` → `resolveDest` (after `DeriveDest` fallback when `Row.Dest` is empty). |
| `branch`                  | string  | `Result.Row.Branch`                                | Input row; empty means "use remote HEAD".           |
| `depth`                   | number  | `Result.Row.Depth`                                 | Input row; `0` means full history.                  |
| `status`                  | string  | `Result.Status`                                    | `executeRow` — one of `ok` / `skipped` / `failed` (`constants.CloneFromStatus*`). |
| `detail`                  | string  | `Result.Detail`                                    | `executeRow`: empty for `ok`, `dest exists` for `skipped`, trimmed git stderr (capped at `GitErrorTrimLimit`) for `failed`. |
| `duration_seconds`        | number  | `Result.Duration.Seconds()` (3-decimal CSV format) | `executeRow` (`time.Since(start)` around the whole row, including skip/parent-dir/checkout phases). |

Envelope-only fields (JSON only — not present in CSV):

| Field             | Type   | Source                                                                                      |
|-------------------|--------|---------------------------------------------------------------------------------------------|
| `schemaVersion`   | number | `constants.CloneFromReportSchemaVersion`.                                                   |
| `transport.ssh`   | number | `clonefrom.TransportTally(results)` — count of rows whose `Row.URL` classifies as SSH.      |
| `transport.https` | number | `clonefrom.TransportTally(results)` — count of HTTPS rows.                                  |
| `transport.other` | number | `clonefrom.TransportTally(results)` — count of rows that are neither SSH nor HTTPS.         |

The `transport` tally matches the terminal `transport: N ssh, N https, N other`
line emitted by `RenderSummary` byte-for-byte, so JSON consumers never
have to re-derive it from the row URLs.

→ Detailed help: [scan](gitmap/helptext/scan.md) · [rescan](gitmap/helptext/rescan.md) · [clone](gitmap/helptext/clone.md) · [clone-next](gitmap/helptext/clone-next.md)

<div align="center">

### Release & Versioning

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `release` | `r` | Create release branch, tag, and push |
| `pull-release` | `pr` *(also `release-pull`, `relp`, `rlp`)* | `git pull` (ff/rebase/merge), then release the current repo |
| `release-alias` | `ra` | Release a repo by its registered alias from anywhere |
| `release-alias-pull` | `rap` | `release-alias` with implicit `--pull` (pull-then-release) |
| `release-self` | `rs` | Release gitmap itself from any directory |
| `release-branch` | `rb` | Create release branch without tagging |
| `temp-release` | `tr` | Create lightweight temp release branches |

#### Alias mapping — `pull-release` family (v5.6.0+)

`pull-release` was renamed from `release-pull` in v5.6.0 to read verb-first.
Every legacy spelling still routes to the exact same handler — no behavior
change, no deprecation warning. Pick whichever you have muscle memory for.

| Token            | Status                        | Routes to      | Notes                                              |
|------------------|-------------------------------|----------------|----------------------------------------------------|
| `pull-release`   | **Canonical** (v5.6.0+)       | `pull-release` | Preferred spelling; shown in help and completions. |
| `pr`             | **Canonical short alias**     | `pull-release` | Freed up by moving `prune` to `prn`.               |
| `release-pull`   | Legacy long form (pre-v5.6.0) | `pull-release` | Wired as alias; safe in scripts and CI.            |
| `relp`           | Legacy short alias            | `pull-release` | Wired as alias.                                    |
| `rlp`            | Legacy short alias            | `pull-release` | Wired as alias.                                    |
| `prn`            | **Canonical short alias** for `prune` (v5.6.0+) | `prune` | Previously `pr`; renamed to free `pr` for `pull-release`. `gitmap prune` (long form) is unchanged. |

```bash
gitmap release --bump patch
gitmap release --bump minor --bin --compress --checksums
gitmap release v3.0.0 -N "Major redesign"

# pull-release: every spelling below does the same thing
gitmap pull-release --rebase v1.4.0      # canonical
gitmap pr --ff-only                      # canonical short alias
gitmap release-pull --merge              # legacy long form
gitmap relp v1.4.0                       # legacy short alias
gitmap rlp --dry-run                     # legacy short alias

# Release any aliased repo from anywhere — no `cd` required
gitmap as my-api                                   # one-time, run from inside the repo
gitmap release-alias my-api v1.4.0
gitmap ra my-api v1.4.0 --pull                     # pull --ff-only, then release
gitmap release-alias-pull my-api v1.4.0            # equivalent thin verb
gitmap rap my-api v1.4.0 --dry-run

gitmap release-self --bump patch
gitmap tr 10 v1.$$ -s 5
```

> Dirty trees are auto-stashed before `release-alias` runs and restored on
> exit. Pass `--no-stash` to abort instead, or `--dry-run` to preview.

→ [release](gitmap/helptext/release.md) · [pull-release](gitmap/helptext/pull-release.md) · [release-alias](gitmap/helptext/release-alias.md) · [release-alias-pull](gitmap/helptext/release-alias-pull.md) · [release-self](gitmap/helptext/release-self.md) · [release-branch](gitmap/helptext/release-branch.md) · [temp-release](gitmap/helptext/temp-release.md)

---

<div align="center">

### Release History & Info

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `changelog` | `cl` | Show release notes |
| `changelog-generate` | `cg` | Auto-generate changelog from commits |
| `list-versions` | `lv` | List all available Git release tags |
| `list-releases` | `lr` | List release metadata from database |
| `release-pending` | `rp` | Show unreleased commits since last tag |
| `revert` | — | Revert to a specific release version |
| `clear-release-json` | `crj` | Remove orphaned release metadata files |
| `prune` | `prn` | Delete stale release branches (alias was `pr` before v5.6.0; `pr` now belongs to `pull-release`) |

```bash
gitmap changelog v2.49.0
gitmap release-pending
gitmap list-versions --json --limit 5
gitmap cg --from v2.22.0 --to v2.24.0 --write
gitmap revert v2.48.0
```

→ [changelog](gitmap/helptext/changelog.md) · [list-versions](gitmap/helptext/list-versions.md) · [list-releases](gitmap/helptext/list-releases.md) · [release-pending](gitmap/helptext/release-pending.md) · [revert](gitmap/helptext/revert.md) · [clear-release-json](gitmap/helptext/clear-release-json.md) · [prune](gitmap/helptext/prune.md)

> **CI Pipeline:** Pushing a `release/*` branch or `v*` tag triggers GitHub Actions to cross-compile 6 targets, generate checksums, and create a GitHub release with changelog and install instructions.

---

<div align="center">

### Data, Profiles & Bookmarks

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `export` | `ex` | Export database to file |
| `import` | `im` | Import repos from file |
| `profile` | `pf` | Manage database profiles |
| `bookmark` | `bk` | Save and run bookmarked commands |
| `rm` | `remove`, `del` | Untrack repos from the database without deleting files |
| `db-reset` | — | Reset the local SQLite database |

```bash
gitmap export && gitmap import gitmap-export.json
gitmap profile create work && gitmap profile switch work
gitmap bookmark save daily scan ~/projects
gitmap bookmark run daily
gitmap del old-repo
```

→ [export](gitmap/helptext/export.md) · [import](gitmap/helptext/import.md) · [profile](gitmap/helptext/profile.md) · [bookmark](gitmap/helptext/bookmark.md) · [rm](gitmap/helptext/rm.md) · [db-reset](gitmap/helptext/db-reset.md)

#### Chrome bookmarks export (`--root` / `--folder`, md / html / json)

```bash
# Whole tree → Markdown file
gitmap chrome export-bookmarks Lovable --out bm.md

# Only the bookmarks bar → HTML
gitmap chrome export-bookmarks Lovable --root bookmark_bar --format html --out bar.html

# Nested subtree → JSON
gitmap chrome export-bookmarks "Profile 1" --root bookmark_bar --folder "Work/Docs" --format json --out work-docs.json

# Other / synced roots
gitmap chrome export-bookmarks Default --root other  --format md
gitmap chrome export-bookmarks Default --root synced --format html --out synced.html
```

`--root` ∈ `bookmark_bar | other | synced` (case-insensitive). `--folder` is slash-delimited and case-insensitive. See [chrome](gitmap/helptext/chrome.md) for `--match` / `--title` filters.


---

<div align="center">

### History, Stats & Author Amendment

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `history` | `hi` | Show CLI command execution history |
| `history-reset` | `hr` | Clear command execution history |
| `stats` | `ss` | Show aggregated usage and performance metrics |
| `amend` | `am` | Rewrite commit author info |
| `amend-list` | `al` | List previous author amendments |

```bash
gitmap history --limit 10
gitmap stats --json
gitmap amend --name "John Doe" --email "john@example.com" --dry-run
```

→ [history](gitmap/helptext/history.md) · [stats](gitmap/helptext/stats.md) · [amend](gitmap/helptext/amend.md) · [amend-list](gitmap/helptext/amend-list.md)

---

<div align="center">

### Project Detection

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `go-repos` | `gr` | List detected Go projects |
| `node-repos` | `nr` | List detected Node.js projects |
| `react-repos` | `rr` | List detected React projects |
| `cpp-repos` | `cr` | List detected C++ projects |
| `csharp-repos` | `csr` | List detected C# projects |

```bash
gitmap go-repos
gitmap csharp-repos --json
```

→ [go-repos](gitmap/helptext/go-repos.md) · [node-repos](gitmap/helptext/node-repos.md) · [react-repos](gitmap/helptext/react-repos.md) · [cpp-repos](gitmap/helptext/cpp-repos.md) · [csharp-repos](gitmap/helptext/csharp-repos.md)

---

<div align="center">

### Tool Installation

</div>

Install developer tools and databases via platform package managers directly from the CLI.

#### Core Tools

| Tool | Keyword | Description |
|------|---------|-------------|
| Visual Studio Code | `vscode` | Code editor |
| Node.js | `node` | JavaScript runtime (includes Yarn, Bun) |
| pnpm | `pnpm` | Fast package manager |
| Python | `python` | Programming language |
| Go | `go` | Programming language |
| Git + LFS + gh | `git`, `git-lfs`, `gh` | Version control ecosystem |
| GitHub Desktop | `github-desktop` | Git GUI |
| C++ (MinGW) | `cpp` | C++ compiler |
| PHP | `php` | Programming language |
| PowerShell | `powershell` | Shell |

#### Databases

| Tool | Keyword | Description |
|------|---------|-------------|
| MySQL | `mysql` | Open-source relational database |
| MariaDB | `mariadb` | MySQL-compatible fork |
| PostgreSQL | `postgresql` | Advanced relational database |
| SQLite | `sqlite` | Embedded file-based database |
| MongoDB | `mongodb` | Document-oriented NoSQL |
| CouchDB | `couchdb` | Document database with REST API |
| Redis | `redis` | In-memory key-value store |
| Cassandra | `cassandra` | Wide-column distributed NoSQL |
| Neo4j | `neo4j` | Graph database |
| Elasticsearch | `elasticsearch` | Full-text search and analytics |
| DuckDB | `duckdb` | Analytical columnar database |

```bash
# Install a tool
gitmap install node
gitmap install postgresql

# Pin a specific version
gitmap install node --version 20.11.1

# Check if installed (no install)
gitmap install go --check

# Preview install command
gitmap install redis --dry-run

# Force a specific package manager
gitmap install vscode --manager winget

# List all supported tools
gitmap install --list

# Uninstall a tool
gitmap uninstall redis
```

**Default package managers by platform:**

| Platform | Default | Fallback |
|----------|---------|----------|
| Windows | Chocolatey | Winget |
| macOS | Homebrew | — |
| Linux | apt | snap |

Override in `config.json` → `install.defaultManager` or per-command with `--manager`.

→ [install](gitmap/helptext/install.md)

---

<div align="center">

### SSH Key Management

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `ssh` | — | Generate and manage SSH keys |

```bash
gitmap ssh --name work --path ~/.ssh/id_rsa_work
gitmap ssh cat --name work
gitmap ssh list
gitmap ssh config
```

→ [ssh](gitmap/helptext/ssh.md)

---

<div align="center">

### Zip Groups (Release Archives)

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `zip-group` | `z` | Manage named file collections for release archives |

```bash
gitmap z create docs-bundle
gitmap z add docs-bundle ./README.md ./CHANGELOG.md ./docs/
gitmap z show docs-bundle
gitmap release v3.0.0 --zip-group docs-bundle
```

→ [zip-group](gitmap/helptext/zip-group.md)

---

<div align="center">

### Environment & File-Sync

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `env` | `ev` | Manage persistent environment variables and PATH |
| `task` | `tk` | Manage file-sync watch tasks |

```bash
gitmap env set GOPATH "/home/user/go"
gitmap env path add /usr/local/go/bin
gitmap env list
gitmap task create my-sync --src ./src --dest ./backup
gitmap tk run my-sync --interval 10
```

→ [env](gitmap/helptext/env.md) · [task](gitmap/helptext/task.md)

---

<div align="center">

### Utilities

</div>

| Command | Alias | Description |
|---------|-------|-------------|
| `setup` | — | Interactive first-time configuration wizard |
| `doctor` | — | Diagnose PATH, deploy, and version issues |
| `update` | — | Self-update from source repo or gitmap-updater |
| `version` | `v` | Show version number |
| `completion` | `cmp` | Generate shell tab-completion scripts |
| `interactive` | `i` | Launch full-screen interactive TUI |
| `docs` | `d` | Open documentation website in browser |
| `seo-write` | `sw` | Auto-commit SEO messages |
| `gomod` | `gm` | Rename Go module path across repo |
| `dashboard` | `db` | Generate interactive HTML dashboard |

```bash
gitmap doctor --fix-path
gitmap update
gitmap completion powershell
gitmap interactive --refresh 10
gitmap dashboard --limit 100 --open
```

→ [setup](gitmap/helptext/setup.md) · [doctor](gitmap/helptext/doctor.md) · [update](gitmap/helptext/update.md) · [completion](gitmap/helptext/completion.md) · [interactive](gitmap/helptext/interactive.md) · [dashboard](gitmap/helptext/dashboard.md)

---

## Build & Deploy

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make all` | Lint → Test → Build (default) |
| `make setup` | Install hooks and dev tools |
| `make lint` | Run golangci-lint |
| `make test` | Run all tests |
| `make build` | Compile for current platform |
| `make vulncheck` | Scan dependencies for CVEs |
| `make release BUMP=patch` | Lint, test, then release |
| `make release-dry` | Preview release without executing |
| `make clean` | Remove build artifacts |
| `make fixtures-bump RUN='<test-regex>' PKG='<pkg>'` | Auto-refresh stale `// fixture-stamp:` markers, then re-verify |
| `make fixtures-bump-verify RUN='<test-regex>' PKG='<pkg>'` | Re-run the same selection without the autobump gate to confirm the rewrite stuck |

### Fixture Stamps & Auto-Bump

Test fixtures across `gitmap/` carry a `// fixture-stamp:` marker that
pins both a `gen=N` generation counter and a `sha=<hex>` hash of the
fixture body (excluding the marker line). On every test run,
`MustValidateBody` (in `gitmap/fixtureversion`) verifies that:

1. The recorded `gen=` matches the current generation.
2. The recorded `sha=` matches `BodyHashExcludingMarker(body)`.

If either drifts — typically because you intentionally regenerated a
golden fixture or bumped the generation number — the test fails with
a clear regenerate recipe instead of producing a confusing diff.

**`fixtures-bump` is the one-shot fix for that failure.** It re-runs
the selected tests with `GITMAP_FIXTURE_AUTOBUMP=1`, which lets
`MaybeAutoBumpFile` rewrite the marker in-place (new `gen` + freshly
computed `sha`), and then re-runs the same selection *without* the
gate to prove the file is now stable. The two-pass design is
deliberate: the verify pass is what actually convinces CI that the
autobump produced deterministic output, not just papered over a real
regression.

#### When to use it

- You changed a writer (formatter, scanner, render layer) on purpose
  and the fixture body legitimately needs to update — `sha=` will
  mismatch on the next run.
- You bumped the fixture generation counter to invalidate a cached
  expectation across the whole suite — `gen=` will mismatch.
- A `MustValidateBody` failure points you at a single test file and
  the new body is what you want to lock in.

#### When NOT to use it

- The fixture failure is a real regression. If you don't understand
  *why* the body changed, do not bump — investigate first. Autobump
  will silently bless a broken writer.
- The determinism pre-check (`AssertWriterDeterministic`) is failing.
  That gate runs writers 3× and refuses to rewrite fixtures on any
  byte divergence — fix the non-determinism first, then bump.

#### Example

```bash
# Fixture for the v9->v12 fix-repo rewriter test went stale after a
# legitimate writer change. Bump just that one test, then verify.
make fixtures-bump \
  RUN='TestFixRepoRewriteV9ToV12Fixture' \
  PKG='./gitmap/cmd/...'
```

The first pass rewrites the `// fixture-stamp:` line in the test
source (or the embedded fixture file the test points at). The second
pass runs the same selection with the env gate cleared — if it
passes, the rewrite is durable; if it fails, the writer is
non-deterministic and you have a real bug to fix.

### Build from Source

```bash
cd gitmap && go build -o ../gitmap .
```

### Build via run.ps1 (Windows)

```powershell
.\run.ps1                        # Full pipeline: pull, build, deploy, setup
.\run.ps1 -R scan                # Build + scan parent folder
.\run.ps1 -R scan D:\repos --mode ssh
.\run.ps1 -uninstall             # Run uninstall-quick.ps1 -Yes and exit
.\run.ps1 -reinstall             # Uninstall, then re-run run.ps1 with no args
.\run.ps1 -NoSetup               # Skip the auto `gitmap setup` after deploy
```

| Flag | Description |
|------|-------------|
| `-NoPull` | Skip `git pull` |
| `-NoDeploy` | Skip deploy step |
| `-NoSetup` | Skip auto-running `gitmap setup` after deploy |
| `-Update` | Update mode with post-update validation |
| `-uninstall` | Run `uninstall-quick.ps1 -Yes` and exit (alias: `-u`) |
| `-reinstall` | Uninstall, then re-invoke `run.ps1` with no args (alias: `-ri`) |
| `-R` | Run gitmap after build (trailing args forwarded) |

PowerShell flags are case-insensitive, so `-uninstall`, `-Uninstall`, and
`-UNINSTALL` are all equivalent — lowercase is preferred for typing speed.

---

## Project Structure

```
/                                repo root
├── README.md                    this file (pinned version + install)
├── CHANGELOG.md                 version-by-version history
├── gitmap/                      Go CLI source (the product)
│   ├── cmd/                     command handlers + JSON contract tests
│   ├── constants/               ALL string constants (no magic strings)
│   ├── model/                   shared Go structs (mirror JSON schemas)
│   ├── stablejson/              deterministic JSON writer (key-by-key)
│   ├── store/                   SQLite layer + migrations
│   ├── completion/              shell completion generators
│   ├── release/                 release workflow + semver
│   ├── formatter/               human / JSON / CSV output formatters
│   ├── helptext/                embedded markdown help files
│   └── scripts/                 embedded install / uninstall scripts
├── gitmap-updater/              standalone updater binary
├── spec/                        specifications (source of truth)
│   ├── 03-general/              cross-cutting rules (logging, build, prohibited)
│   ├── 04-generic-cli/          per-command specs
│   ├── 05-coding-guidelines/    Go style, naming, errors
│   ├── 08-json-schemas/         JSON Schema (draft-07) for every JSON output
│   └── 12-consolidated-guidelines/  merged style guide
├── scripts/                     helper shell / PowerShell scripts
├── .github/workflows/           CI pipelines
├── .gitmap/                     runtime artifacts — do NOT hand-edit
├── .lovable/                    AI memory + plans + prompts
└── src/                         React docs site (Vite + shadcn)
```

### What to Read (AI onboarding map)

A fresh AI session should open these files **in order** to understand the
project, the JSON output contracts, and how to add a new JSON surface:

1. This README — product overview + folder structure (you are here).
2. **[`.lovable/memory/project/what-to-read.md`](./.lovable/memory/project/what-to-read.md)**
   — full onboarding map, JSON contract triangle
   (schema ↔ model ↔ encoder ↔ test ↔ fixture), and the step-by-step
   recipe for adding a new JSON output.
3. `.lovable/memory/index.md` — master memory index (Core rules + topic links).
4. `.lovable/overview.md` and `.lovable/strictly-avoid.md` — invariants & hard NOs.
5. `spec/08-json-schemas/` — every JSON output's schema.
6. `gitmap/cmd/amendauditrender.go` + `amendaudit_jsonschema_contract_test.go`
   — canonical example of the encoder + contract-test pattern to copy.

### How JSON Outputs Are Structured

Every JSON-emitting command follows the same four-layer contract:

```
spec/08-json-schemas/<name>.schema.json          ← contract (draft-07)
gitmap/model/<name>.go                           ← Go struct
gitmap/cmd/<name>render.go                       ← stablejson encoder (key order = wire contract)
gitmap/cmd/<name>_jsonschema_contract_test.go    ← drift guard
```

See `.lovable/memory/project/what-to-read.md` §3–4 for the full recipe.

---

## Web UI Dashboard

GitMap includes a React-based documentation and dashboard UI:

```bash
npm install && npm run dev     # opens at http://localhost:5173
```

**Tech Stack:** Vite · TypeScript · React · shadcn/ui · Tailwind CSS

---

## Author

<div align="center">

### [Md. Alim Ul Karim](https://www.google.com/search?q=alim+ul+karim)

**[Creator & Lead Architect](https://alimkarim.com)** | [Chief Software Engineer](https://www.google.com/search?q=alim+ul+karim), [Riseup Asia LLC](https://riseup-asia.com)

</div>

A system architect with **20+ years** of professional software engineering experience across enterprise, fintech, and distributed systems. His technology stack spans **.NET/C# (18+ years)**, **JavaScript (10+ years)**, **TypeScript (6+ years)**, and **Golang (4+ years)**.

Recognized as a **top 1% talent at Crossover** and one of the top software architects globally. He is also the **Chief Software Engineer of [Riseup Asia LLC](https://riseup-asia.com)** and maintains an active presence on **[Stack Overflow](https://stackoverflow.com/users/361646/alim-ul-karim)** (2,452+ reputation, member since 2010) and **LinkedIn** (12,500+ followers).

|  |  |
|---|---|
| **Website** | [alimkarim.com](https://alimkarim.com/) · [my.alimkarim.com](https://my.alimkarim.com/) |
| **LinkedIn** | [linkedin.com/in/alimkarim](https://linkedin.com/in/alimkarim) |
| **Stack Overflow** | [stackoverflow.com/users/361646/alim-ul-karim](https://stackoverflow.com/users/361646/alim-ul-karim) |
| **Google** | [Alim Ul Karim](https://www.google.com/search?q=Alim+Ul+Karim) |
| **Role** | Chief Software Engineer, [Riseup Asia LLC](https://riseup-asia.com) |

### Riseup Asia LLC

[Top Leading Software Company in WY (2026)](https://riseup-asia.com)

| | |
|---|---|
| **Website** | [riseup-asia.com](https://riseup-asia.com) |
| **Facebook** | [riseupasia.talent](https://www.facebook.com/riseupasia.talent/) |
| **LinkedIn** | [Riseup Asia](https://www.linkedin.com/company/105304484/) |
| **YouTube** | [@riseup-asia](https://www.youtube.com/@riseup-asia) |

## License

This project is licensed under the [MIT License](./LICENSE).
