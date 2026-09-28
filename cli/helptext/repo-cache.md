# `gitmap repo-cache` (`rc` / `repo-storage`) — Reusable Scripts & Scratch Storage Repository

## Synopsis

```bash
gitmap repo-cache <subcommand> [arguments] [flags]
gitmap repo-storage <subcommand> [arguments] [flags]
gitmap rc <file|folder|text|ls|init|help> [arguments] [flags]
gitmap cd rc
```

`gitmap repo-cache` (aliases `gitmap rc` and `gitmap repo-storage`) manages the dedicated **`repo-cache`** special repository for archiving reusable PowerShell (`.ps1`) scripts, diagnostic test harnesses, one-off automation utilities, and cross-project scratch assets. Instead of cluttering application repositories or losing helpful temporary scripts after a session, `gitmap rc` sequences them under `XX-<repo-name>/01-<slug>.ext` and auto-commits & pushes them for instant reuse across any project.

---

## Sequenced Hierarchy (`XX-<repo>/01-<slug>.ext`)

Every script, file, or folder stored with `gitmap rc` is placed into a 2-digit project directory (`XX-<repo-name>/`) and assigned a 2-digit sequential filename (`01-<slug>.ext`):

```text
repo-cache/
├── 01-gitmap-v28/
│   ├── 01-verify-sqlite-wal.ps1
│   ├── 02-benchmark-scanner.ps1
│   └── 03-test-fixtures/
└── 02-cluster-fleet/
    └── 01-check-ssh-latency.sh
```

- **Reusable PowerShell `.ps1` & Test Scripts**: Store `.ps1`, `.sh`, or `.py` diagnostic scripts once and execute or reference them from any workspace.
- **Auto-Commit & Push**: Automatically commits and pushes newly stored scripts to `repo-cache` via `commit-in`.

---

## Navigation Shortcut (`gitmap cd rc`)

Jump directly into the `repo-cache` repository directory:

```bash
gitmap cd rc
gitmap cd repo-cache
gitmap cd repo-storage
```

---

## Subcommands & Examples

### 1. Archive a Reusable Script or File (`file`)
Store a PowerShell `.ps1` script or test utility into `repo-cache/XX-<repo>/01-<slug>.ps1`:

```bash
gitmap rc file ./scratch/verify-cluster.ps1
gitmap rc file ./tools/reset-db.ps1 --slug reset-db --repo gitmap-v28
```

### 2. Archive a Scratch or Fixture Folder (`folder`)
Store a folder of test harnesses or diagnostic scripts into `repo-cache/XX-<repo>/01-<slug>/`:

```bash
gitmap rc folder ./scratch/ps1-harnesses --slug ps1-harnesses
```

### 3. Store Inline PowerShell or Script Content (`text`)
Save an inline PowerShell `.ps1` or shell script directly into `repo-cache`:

```bash
gitmap rc text "Get-Process | Where-Object { $_.CPU -gt 50 }" --slug high-cpu-check --ext .ps1
gitmap rc text "go version; git status -sb" --slug env-check --ext .ps1
```

### 4. List Cached Scripts & Assets (`ls`)
Display all sequenced project folders and cached scripts in `repo-cache`:

```bash
gitmap rc ls
gitmap rc ls --json
```

### 5. Initialize or Verify `repo-cache` (`init`)
Create or initialize the local `repo-cache` repository directory:

```bash
gitmap rc init
```

---

## Customizing the Special Cache Repository Name (`gitmap settings`)

By default, GitMap resolves the repository named `repo-cache`. Customize the repository folder name globally via `gitmap settings`:

```bash
gitmap settings set special_repos.cache_name repo-cache
gitmap settings set cache_repo repo-storage
gitmap settings
```
