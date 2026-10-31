# GitMap v6.639.0

## What's Changed in v6.639.0

- **Multi-Node Deploy, Copy (CP), Move (MV) & Fleet Software Updates**

### Key Features & Architectural Enhancements

#### 1. Universal Alias Reassignment: `cp` for Copy, `mv` for Move
- **Retired Legacy Collision:** Reassigned `CmdCommitPushAlias` from `cp` to `cmp` so `cp` is permanently dedicated to file copying (`gitmap copy` / `gitmap cp`), respecting universal developer muscle memory.
- **Dedicated Move Alias:** Formally aligned `mv` with `gitmap move` (`gitmap move` / `gitmap mv`).
- **Preserved Semantic Verbs:** All discrete 3-letter semantic commit verbs (`cpf`, `cpb`/`cbf`, `cpr`, `cpc`, `cfr`) remain 100% intact.
- **Clean Disambiguation:** `gitmap copy` automatically routes to remote fleet copy when target nodes or remote delimiters (`node:path`) are passed, or falls back to OS clipboard and memory copying when plain text or local file arguments are supplied.

#### 2. Multi-Node File Transport (8 Copy Syntax Variants)
Implemented in [cli/cmdssh/ssh_cp.go](cli/cmdssh/ssh_cp.go):
1. **Local to Single Node:** `gitmap cp <local-src> <node>:<dest>`
2. **Local to Single Node with Flag:** `gitmap cp <local-src> <dest> -t <node>` (or `--target <node>`)
3. **Local to Fleet Broadcast:** `gitmap cp <local-src> all:<dest>` (or `--all`)
4. **Remote Node Pull to Local:** `gitmap cp <node>:<remote-src> <local-dest>`
5. **Direct Node-to-Node Transfer:** `gitmap cp <node1>:<src> <node2>:<dest>`
6. **Remote Node to Fleet Replicate:** `gitmap cp <node1>:<src> all:<dest>`
7. **Recursive Directory Copy:** `gitmap cp -r <local-dir> <node>:<dest-dir>` via streaming tar sessions
8. **Memory & Clipboard Disambiguation:** `gitmap copy <text>` or `gitmap copy --file <path>`

#### 3. Remote Node Move Engine (`gitmap mv` / `gitmap move`)
Implemented in [cli/cmdssh/ssh_mv.go](cli/cmdssh/ssh_mv.go):
- **Local to Remote Move:** Uploads file to remote node, verifies transfer, then removes local source file.
- **Remote to Local Move:** Pulls file to local filesystem, verifies write, then executes platform-native removal (`rm -f` on Linux, `Remove-Item` on Windows).
- **Inter-Node Move:** Relays file between remote cluster nodes and cleans source upon verified transfer.
- **Local to Fleet Move:** Broadcasts file to all reachable nodes and cleans local source.
- **Local Folder Disambiguation:** When executed without remote node prefixes, seamlessly preserves local repository/directory relocation (`cmdmerge.RunMove`).

#### 4. OS-Aware Binary Deployment (`gitmap deploy`)
Implemented in [cli/cmdssh/ssh_deploy_bin.go](cli/cmdssh/ssh_deploy_bin.go), [cli/cmdssh/ssh_deploy_compat.go](cli/cmdssh/ssh_deploy_compat.go), and [cli/cmdssh/ssh_deploy_feedback.go](cli/cmdssh/ssh_deploy_feedback.go):
- **Binary Magic Byte Inspection:** Inspects executable headers (`MZ` for PE Windows `.exe`, `\x7fELF` for Linux) to classify binary formats dynamically.
- **Automatic OS Incompatibility Filter:** By default, Windows `.exe` binaries are automatically skipped on Linux/Ubuntu targets with `SKIPPED_INCOMPATIBLE_OS` without failing the batch.
- **Force Override:** Pass `--force` or `--force-copy` to override OS filtering when using compatibility layers (Wine, emulation).
- **Concurrent Offline Bypass:** Non-blocking TCP liveness probes (1200ms timeout) skip unreachable nodes immediately without stalling.
- **Remote Installation Hooks:** Automatically applies permissions (`chmod 755`), moves binary to destination (`/usr/local/bin` or custom `--dest`), and executes self-test verification.
- **Diagnostic Feedback Docket:** Displays post-execution summary table with exact success/skip breakdowns and copy-pasteable next steps (cross-compilation commands, ping diagnostics).

#### 5. Fleet-Wide Software Updates
Implemented in [cli/cmdnodes/nodes_update_software.go](cli/cmdnodes/nodes_update_software.go):
- Centralized commands to upgrade developer tools across single nodes or the entire fleet:
  - `gitmap update cursor [node|all]`
  - `gitmap update vscode [node|all]`
  - `gitmap update agm [node|all]`
  - `gitmap update agy [node|all]` (or `update agy-tools`)
  - Subcommand alias: `gitmap node update <app> [node|all]`
- Cross-platform package manager execution: Winget silent updates on Windows, APT / script automation on Ubuntu/Linux.

#### 6. Terminal Installation Suggestions Remediation & Offline Parity
Implemented in [install.ps1](install.ps1), [cli/scripts/install.ps1](cli/scripts/install.ps1), [cli/fspath/install.ps1](cli/fspath/install.ps1), [local-install.ps1](local-install.ps1), [cli/cmd/root_cobra_completion.go](cli/cmd/root_cobra_completion.go), [cli/completion/powershell.go](cli/completion/powershell.go), and [cli/cmdinstall/install_prompt.go](cli/cmdinstall/install_prompt.go):
- **Persistent Profile Suggestions Reporting:** Actively tracks `$activeProfiles` alongside `$configuredProfiles` so suggestion status never disappears on re-runs or updates.
- **Offline Local Installer Parity:** Ported `Configure-PowerShellProfileSuggestions` and profile targets resolution into `local-install.ps1` Stage 5.
- **Cobra Completer for Short Alias (`gm`):** Registered argument completion for both `gitmap` and `gm`.
- **Release Install Hints:** Recognized unversioned canonical `gitmap` repository URLs.
- **Actionable `gitmap install gitmap` Suggestions:** Emits an interactive suggestions box pointing to `completion install`, `status`, `update`, and fleet operations.

#### 7. Documentation & UI Synchronization
- Updated [src/data/commands.ts](src/data/commands.ts) with `deploy`, `copy`, `move`, and `update-software` entries.
- Added and updated help documentation in [cli/helpdoc/copy.md](cli/helpdoc/copy.md), [cli/helpdoc/cp.md](cli/helpdoc/cp.md), [cli/helpdoc/mv.md](cli/helpdoc/mv.md), [cli/helpdoc/deploy-bin.md](cli/helpdoc/deploy-bin.md), and [cli/helpdoc/update.md](cli/helpdoc/update.md).
- Updated canonical specification in [02-spec/21-app/364-gitmap-command-enhancements/](02-spec/21-app/364-gitmap-command-enhancements/).

---

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.639.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.639.0/install.sh | sh
```
