# Component Technical Specification: Multi-Node Update Engine, Zip SCP Distribution, Prompt Search Primacy, and .gitignore Backup Policy

- **Spec Document:** `02-spec/21-app/65-gitmap-update-all-zip-and-fixes/02-component-spec.md`
- **Spec ID:** 65-02
- **Status:** Approved / Ready for Implementation
- **Author:** Spec Author 02
- **Parent Plan:** [.ai-memory/plans/pending/65-gitmap-update-all-zip-and-fixes.md](../../../.ai-memory/plans/pending/65-gitmap-update-all-zip-and-fixes.md)
- **Target Subtasks:**
  - [Subtask 65.2: Multi-Node Update Engine & Zip SCP Distribution](../../../.ai-memory/plans/subtasks/65-gitmap-update-all-zip-and-fixes/02-update-all-zip-distribution.md)
  - [Subtask 65.4: Canonical Prompts GitMap Search Primacy & Sub-Node Diagnostics](../../../.ai-memory/plans/subtasks/65-gitmap-update-all-zip-and-fixes/04-prompt-gitmap-search-primacy.md)
  - [Subtask 65.6: .gitignore Policy Engine: Retain .gitmap/ & Exclude Only .gitmap/backup/](../../../.ai-memory/plans/subtasks/65-gitmap-update-all-zip-and-fixes/06-gitignore-policy-gitmap-backup.md)

---

## 1. Executive Summary & Architectural Scope

This specification details the technical component designs for three interrelated capabilities within GitMap v28:
1. **Multi-Node Update Engine via Packaged Zip & SCP Distribution (Task-02):** Enables updating distributed cluster nodes and worker machines by packaging updated binaries locally (`gitmap` or Anti-Gravity Manager `agm`), streaming compressed archives across SSH sessions via `StreamFileToRemote`, and executing OS-tailored remote installation scripts (PowerShell for Windows, Shell/Bash for POSIX). This reduces external WAN traffic and circumvents remote network bandwidth bottlenecks.
2. **Canonical Prompts & Coding Guidelines: GitMap Search Primacy (Task-04):** Updates the high-discipline autonomous execution prompt (`[V6] Parent Task N-Step Continuous Loop`) and its associated Antigravity skill to enforce an explicit ban on raw unindexed search utilities (`rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`), replacing them with high-speed GitMap commands (`gitmap aum search`, `gitmap search`). Additionally establishes duplicate package path disambiguation rules and failure sub-node diagnostic branching.
3. **.gitignore Policy & Backup Exemption Engine (Task-06):** Refactors GitMap's ignore enforcement and automated pull sanitation to ensure that `.gitmap/` is preserved and version-controlled by default. Excludes only `.gitmap/backup/` and temporary resume files, eliminates legacy development repo conditionals, and enforces strict pattern normalization to remove duplicate ignore lines.

---

## 2. Component 1: Multi-Node Update Engine & Zip SCP Distribution

### 2.1 Problem Statement & Operational Rationale
Existing fleet update mechanisms (`gitmap update all`, `gitmap ua`) instruct remote nodes to independently fetch installers and binary assets over the public internet from GitHub releases (`irm https://raw.githubusercontent.com/...`). In distributed multi-node environments:
- Multiple remote worker machines simultaneously saturated external WAN connections downloading identical binary assets.
- Air-gapped, firewalled, or egress-restricted worker nodes failed to perform updates due to blocked GitHub network access.
- There was no unified command to update Anti-Gravity Manager (`agm`) across the entire fleet in parallel using pre-packaged or locally staged artifacts.

### 2.2 CLI Surface & Routing Table

The CLI router in `cli/cmd/rootutility.go` and dispatcher in `cli/cmdupdate/update_fleet.go` are extended to support the following invocations:

| Invocation | Command Token | Arguments | Action |
| :--- | :--- | :--- | :--- |
| `gitmap update all --include-others` | `update` | `all`, `--include-others` | Standard fleet update expanding targets to all known cluster nodes |
| `gitmap update all zip --include-others` | `update` | `all`, `zip`, `--include-others` | Zip-packaged local SCP distribution to SSH targets and cluster nodes |
| `gitmap update-all-zip --include-others` | `update-all-zip` | `--include-others` | Direct zip-packaged fleet distribution |
| `gitmap uaz --include-others` | `uaz` | `--include-others` | Shorthand alias for `update-all-zip` |
| `gitmap uaz agm --include-others` | `uaz` | `agm`, `--include-others` | Zip-packaged Anti-Gravity Manager (AGM) distribution across all nodes |

#### Router Registration (`cli/cmd/rootutility.go`)
```go
func utilityCoreEntries() []dispatchEntry {
    return []dispatchEntry{
        // ...
        {[]string{
            constants.CmdUpdate, 
            "ua", 
            "update-all", 
            "updateall", 
            "uaz", 
            "update-all-zip", 
            "updateallzip",
        }, runUpdateHelp},
        // ...
    }
}
```

### 2.3 Options Data Structure & Flag Parsing

In `cli/cmdupdate/update_fleet.go`:

```go
type FleetUpdateOptions struct {
    Target        string // Single target alias/IP filter
    Pkg           string // Target application: "gitmap", "agm", or "all"
    Except        string // Comma-separated exclusion list
    IsAll         bool   // Flag indicating fleet-wide update
    IsZip         bool   // Flag indicating zip packaging and SCP distribution
    IncludeOthers bool   // Expand beyond saved SSH connections to discovered cluster nodes
    IsDryRun      bool   // Simulate execution without performing modifications
    IsForce       bool   // Force update regardless of existing version match
}
```

#### Parsing Rules:
1. `isAllFleetToken(token)` recognizes `--all`, `-all`, `all`, `all-nodes`, `allnodes`, `-a`.
2. When parsing args, the presence of token `zip` (positional or `--zip`) sets `opts.IsZip = true`.
3. The presence of `--include-others` or `--includeothers` sets `opts.IncludeOthers = true`.
4. Commands `uaz`, `update-all-zip`, `updateallzip` automatically imply `opts.IsAll = true` and `opts.IsZip = true`.
5. Package targets `agm`, `ag-manager`, `antigravity-manager` resolve `opts.Pkg = "agm"`.

### 2.4 Target Resolution with `--include-others`

When `--include-others` is false, `loadDefaultFleetTargets()` loads only configured records from `cmdssh.FetchAllSSHConnections()`.
When `--include-others` is true:
1. Query `cmdssh.FetchAllSSHConnections()` for primary SSH hosts.
2. Query `db.ListClusterNodes()` or registered worker nodes from Split-DB.
3. Merge targets into a single slice of `FleetTarget`, deduplicating by normalized IP address:
   ```go
   func loadFleetTargetsWithCluster(includeOthers bool) ([]FleetTarget, error)
   ```

### 2.5 Local Packaging & Zip Asset Pipeline

When `opts.IsZip` is true:
1. **Asset Identification:**
   - For `gitmap`: Locate the active host executable via `os.Executable()`. Read binary bytes and metadata.
   - For `agm`: Locate `agm.exe` (Windows) or `agm` (Unix) via `exec.LookPath("agm")` or configured AGM path. If missing locally, download once to local temp directory.
2. **Zip Creation (`createUpdatePackageZip`):**
   - Package the target executable into an in-memory zip archive buffer (`bytes.Buffer` using `archive/zip`).
   - Include the executable file:
     - Windows: `gitmap.exe` or `agm.exe`
     - Linux/macOS: `gitmap` or `agm` (with `0755` permissions mode bits preserved)
   - Include an embedded launcher script `install_remote.ps1` (for Windows) and `install_remote.sh` (for Unix) to cleanly terminate existing processes, replace the target binary, and write telemetry.
3. **Payload Caching:**
   - Cache the generated zip buffer in memory for the duration of the command invocation so all concurrent worker threads distribute identical bytes without re-compressing per node.

### 2.6 SSH Streaming via `StreamFileToRemote` & Remote Execution

1. **Remote Destination Staging:**
   - Destination directory on remote target:
     - Windows: `%TEMP%\gitmap_update_<pkg>.zip`
     - POSIX: `/tmp/gitmap_update_<pkg>.zip`
2. **Direct SSH Streaming:**
   - Use `cmdssh.StreamFileToRemote(client, destZipPath, zipData, target.OS)`.
   - Streams raw zip archive directly across the SSH channel without relying on external SCP/SFTP binaries or secondary ports.
3. **Remote Unpack & Installation Execution:**
   - Construct OS-specific uninstallation and upgrade command:
     - **Windows (`PowerShell`):**
       ```powershell
       powershell -NoProfile -ExecutionPolicy Bypass -Command "& {
           $ErrorActionPreference = 'Stop'
           $zipPath = '%TEMP%\gitmap_update_gitmap.zip'
           $destDir = Split-Path (Get-Command gitmap -ErrorAction SilentlyContinue).Path
           if (-not $destDir) { $destDir = [System.IO.Path]::Combine($env:LOCALAPPDATA, 'Programs', 'gitmap') }
           Expand-Archive -Path $zipPath -DestinationPath $destDir -Force
           Remove-Item -Path $zipPath -Force -ErrorAction SilentlyContinue
           $curr = (gitmap version 2>$null | Out-String).Trim()
           @{ success = $true; current_version = $curr; details = ('Zip installed: ' + $curr) } | ConvertTo-Json -Compress
       }"
       ```
     - **POSIX (`sh` / `bash`):**
       ```bash
       sh -c '
           ZIP="/tmp/gitmap_update_gitmap.zip"
           DEST=$(dirname "$(command -v gitmap 2>/dev/null || echo /usr/local/bin/gitmap)")
           unzip -o "$ZIP" -d "$DEST" >/dev/null 2>&1
           chmod +x "$DEST/gitmap"
           rm -f "$ZIP"
           CURR=$(gitmap version 2>/dev/null | head -n1)
           printf "{\"success\":true,\"current_version\":\"%s\",\"details\":\"Zip installed: %s\"}" "$CURR" "$CURR"
       '
       ```
4. **Structured Telemetry Ingestion:**
   - Output from remote command is parsed by `ParseFleetUpdateTelemetry`.
   - Results are rendered in the standard cyan/green table via `renderFleetUpdateSummary`.

---

## 3. Component 2: Canonical Prompts & Coding Guidelines Update (Search Primacy)

### 3.1 Problem Statement & Rationale
Autonomous agent execution frequently degrades when subagents run unguided search tools such as PowerShell `Select-String`, recursive `Get-ChildItem`, or external unindexed tools (`rg`, `grep`, `findstr`). These commands:
- Bypass GitMap's multi-tier Split-DB indexed caching and hot memory cache.
- Trigger high CPU load and filesystem lock contention on large repositories.
- Fail to provide normalized symbol matching across polyglot project structures.
- Duplicate reporting when packages have identical names across nested directory hierarchies.

### 3.2 Target Artifacts
- **Primary Prompt:** `01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md`
- **Antigravity Skill:** `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`

### 3.3 Rule Specifications & Explicit Tool Bans

#### Rule 1: GitMap Search Primacy (Mandatory Command Injection)
All code searching, pattern scanning, and symbol discovery across repository files MUST use GitMap high-speed search tools.
- **Strictly Banned:**
  - `Select-String` (PowerShell)
  - `Get-ChildItem -Recurse` / `gci -r`
  - `grep` / `egrep` / `fgrep`
  - `git grep`
  - `findstr`
  - raw `rg` / `ripgrep`
- **Mandatory GitMap Commands:**
  - Text & Pattern Search:
    ```bash
    gitmap aum search "<pattern>" [directory] [-e <.extension>] [-r] [-i]
    ```
  - Fast Indexed Symbol Lookup:
    ```bash
    gitmap search "<symbol>"
    ```
  - Fast File Listing:
    ```bash
    gitmap lf [directory]
    ```
  - Fast File Inspection:
    ```bash
    gitmap cat <relative_filepath>
    ```

#### Rule 2: Package Path Disambiguation Rule
When a package or module appears multiple times in a repository (e.g., duplicated package names in different directory trees), agents must never refer to the package by bare name alone:
- **Requirement:** Always qualify repeated packages with their distinct relative repository paths (e.g., `cli/cmdpull/pull.go` vs `cli/cmdignore/fix_ignore.go`).
- **Enforcement:** Subtask files and agent reports must cite explicit file paths rather than abstract package names.

#### Rule 3: Sub-Node Diagnostic Branching for Failures
When an operation, command, or package test fails during execution, agents must format findings into structured sub-nodes with up to two concrete alternative remediation paths:
```markdown
- **Failure Condition:** [Concise description of error / conflict]
  - **Sub-node [Option A]:** [Primary remediation approach, e.g., rebase / re-stage / targeted fix]
  - **Sub-node [Option B]:** [Alternative remediation approach, e.g., fallback branch / config bypass]
```

---

## 4. Component 3: .gitignore Policy & Backup Exemption Engine

### 4.1 Problem Statement & Rationale
GitMap creates a local metadata and cache folder `.gitmap/` in repositories. Previously:
1. `cli/cmdpull/pull.go` contained logic that treated any repository without "gitmap" in its name as a foreign repo and automatically forced `.gitmap/` into the default ignore list.
2. `cli/cmdignore/ignore_groups.go` had `ImmutableDefaultRules` defining both `.gitmap/` and `.gitmap/backup/` as immutable entries that were forcefully injected into `.gitignore`.
3. `cli/cmdignore/fix_ignore.go` flagged repositories missing `.gitmap/` in `.gitignore` as an issue and automatically added `.gitmap/` during automated remediation (`gitmap fix-ignore-all`).

This behavior broke repositories that deliberately track `.gitmap/` metadata, configuration files, and project settings in git. The correct policy is:
- `.gitmap/` root folder is **NOT** ignored by default and is allowed to be version-controlled.
- Only `.gitmap/backup/` (containing database backups and temporary dumps) is ignored by default.
- If a user explicitly chooses to add `.gitmap/` to their `.gitignore`, GitMap respects it, but GitMap will never automatically enforce or inject `.gitmap/`.
- Duplicate patterns must be completely removed by normalizing ignore lines.

### 4.2 Code Changes & Architectural Refactoring

#### 1. `cli/cmdpull/pull.go`
Refactor `buildDefaultTrackedPathsToCheck`:
```go
func buildDefaultTrackedPathsToCheck(repoDir string) []string {
    return []string{
        ".gitmap/backup/",
        "antigravity-resume_task.json",
        ".antigravity_resume_task.json",
        "antigravity_resume_task.json",
        ".antigravity-resume_task.json",
    }
}
```
- Remove the `.gitmap/` prefix injection.
- Remove obsolete helper `isGitmapDevelopmentRepo(repoDir)`.

#### 2. `cli/cmdignore/ignore_groups.go`
Update immutable defaults and enforcement helper:
```go
// ImmutableDefaultRules defines the mandatory immutable rules for GitMap.
var ImmutableDefaultRules = []string{
    ".gitmap/backup/",
}

// EnsureImmutableRules ensures .gitmap/backup/ is present in patterns.
func EnsureImmutableRules(patterns []string) []string {
    hasBackup := hasExactPattern(patterns, ".gitmap/backup/")
    if !hasBackup {
        return append(patterns, ".gitmap/backup/")
    }
    return patterns
}
```
- Remove `.gitmap/` from `ImmutableDefaultRules`.
- Remove `.gitmap/` prepend from `EnsureImmutableRules`.

#### 3. `cli/cmdignore/fix_ignore.go`
Refactor analysis, deduplication, and assembly logic:
1. **Analyze Gitignore Data:**
   - Check `issue.MissingBackupDir`: verify presence of `.gitmap/backup/`.
   - Do **NOT** flag missing `.gitmap/` as an issue.
2. **Assemble Cleaned Gitignore:**
   ```go
   func assembleCleanedGitignore(cleaned []string, seen map[string]bool) string {
       hasBackup := seen[".gitmap/backup/"] || seen[".gitmap/backup"]
       if !hasBackup {
           cleaned = append(cleaned, "", "# GitMap Backup Persistence", ".gitmap/backup/")
           seen[".gitmap/backup/"] = true
       }
       text := strings.Join(cleaned, "\n")
       return strings.TrimRight(text, "\r\n") + "\n"
   }
   ```
3. **Pattern Normalization & Deduplication:**
   - Normalize patterns by trimming whitespace and stripping leading/trailing slashes:
     ```go
     norm := strings.Trim(trimmed, "/ \t\r\n")
     ```
   - Matches `dist`, `/dist`, `dist/`, and `/dist/` as identical patterns, eliminating duplicate ignore rules.

---

## 5. Security & Boundary Verification

1. **Path Traversal Guards:** All remote archive extraction scripts must extract within target application directories without `..` traversal vulnerabilities.
2. **Credential Safety:** Remote SSH dials use encrypted passwords or local SSH private keys; zero plaintext credentials logged to console or stored in telemetry.
3. **Immutability of User Customizations:** The `.gitignore` engine never deletes custom user rules; it only deduplicates identical lines and ensures the presence of `.gitmap/backup/`.

---

## 6. Implementation Checklist & Traceability

- [ ] Task-02: Register `uaz`, `update-all-zip`, `updateallzip` in `cli/cmd/rootutility.go`.
- [ ] Task-02: Add `IsZip` and `IncludeOthers` to `FleetUpdateOptions` in `cli/cmdupdate/update_fleet.go`.
- [ ] Task-02: Implement `createUpdatePackageZip` and streaming installation in `cli/cmdupdate/update_fleet.go`.
- [ ] Task-02: Add unit tests in `cli/cmdupdate/update_fleet_test.go` covering zip distribution mocks.
- [ ] Task-04: Update `01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md` with GitMap search primacy rules and examples.
- [ ] Task-04: Update `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` in lockstep.
- [ ] Task-06: Remove `.gitmap/` from default tracked check in `cli/cmdpull/pull.go`.
- [ ] Task-06: Update `ImmutableDefaultRules` and `EnsureImmutableRules` in `cli/cmdignore/ignore_groups.go`.
- [ ] Task-06: Fix `.gitignore` analysis, assembly, and deduplication in `cli/cmdignore/fix_ignore.go`.
