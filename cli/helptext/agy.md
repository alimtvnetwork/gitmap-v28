# agy

Manage Google Antigravity (AGY) workspaces, rerun historical prompts with prefix verification templates, inspect and diff prompts across projects with non-admin VS Code, scan project activity, and orchestrate remote execution across SSH, Cluster, SC, and Local targets.

## Synopsis

```bash
gitmap agy <subcommand> [flags]
gitmap ag <subcommand> [flags]
gitmap antigravity <subcommand> [flags]
```

## Aliases & Shorthands

- `gitmap agy`
- `gitmap ag`
- `gitmap antigravity`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `ping` | Ping Antigravity IDE and inspect environment health (alias: `check`) |
| `prompt-project <prefix>` | Target Antigravity project by prefix with `-name`, `-txt`, `--pf/--sf` (alias: `p`) |
| `prompt [-n name] [-t txt]` | Dispatch prompt to current repo (auto-adds project, queues read-all first) |
| `prompt-with-name <name>` | Dispatch named template with optional `-txt` (alias: `pwn`; same as `prompt -name`) |
| `prompt-txt "<text>"` | Dispatch direct text (alias: `pt`; defaults to read-all prefix with 2 newlines) |
| `prompt ls` | List available prompt templates in formatted table (alias: `list prompts`) |
| `prompt read [conv-id]` | Read and display user prompts for a conversation from Antigravity brain |
| `queue [status\|ls\|clear\|pop]` | Manage pending prompt queue entries |
| `status` | Show Antigravity system status (IDE process, IDLE vs RUNNING, prompt queue) |
| `rerun last [N]` | Replay last N prompts with optional prefix verification prompt template (`-p`) |
| `list-prompts [N]` | List historical prompts for current project or across all projects |
| `scan [path]` | Scan repositories, recent prompt activity counts (24h/7d), and prompt archives |
| `fix-pipeline [repo]` | Extract failing CI/CD pipeline errors and generate 4-part RCA prompt (alias: `aef`) |
| `open [path]` | Locate Antigravity IDE and open target workspace |
| `optimize-projects` | Remove stale or duplicate workspace configurations |
| `clean-cache` | Clear temporary Antigravity cache files and transcripts |
| `recreate-project` | Purge cache/convs, re-register project in AGY, and start read-memory conv (alias: `recreate`, `rp`) |
| `pin-projects` | Manage pinned Antigravity projects |
| `plugins` | Inspect and configure Antigravity plugins |
| `reconcile` | Reconcile disk repositories with registered projects |
| `help` | Show this Antigravity command suite documentation |

---

## Examples

```bash
# Rerun the last prompt with the default verification template
gitmap agy rerun last 1

# List prompt activity and open multi-project changes in VS Code
gitmap agy list-prompts --projects 3

# Scan repositories and prompt archives
gitmap agy scan

# Fix pipeline errors with Antigravity
gitmap agy fix-pipeline

# Recreate current project (purge cache/convs, re-register, start read-memory conversation)
gitmap agy recreate-project

# Recreate multiple selected projects by sequence, alias, or path
gitmap agy recreate 1, 2, 3
gitmap agy recreate gitmap, test-gitmap
gitmap agy recreate D:\test-gitmap\test-gitmap
```

---

## 1. agy rerun

Rerun the last prompt for a specific project with Antigravity IDE restart and automatic media/picture replay. Supports targeting projects by index sequence (`1`, `2`, `3`, `4`), terminating running IDE processes, relaunching into the target workspace, and re-injecting the exact active prompt payload with all attached screenshots and pictures.

Also supports replaying the last N user prompts recorded in conversation transcripts into clipboard (`gitmap agy rerun last [N]`).

### Usage

```bash
# Rerun project by sequence index with IDE restart and prompt + picture replay
gitmap agy rerun [1|2|3|4|project-name] [flags]
gitmap agy rerun-restart [1|2|3|4] [flags]
gitmap agy rr [1|2|3|4] [flags]

# Legacy clipboard replay of last N prompts
gitmap agy rerun last [N] [flags]
```

### Options & Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--restart` | `-r` | boolean | `true` | Restart the Antigravity IDE process and replay prompt |
| `--no-restart` | | boolean | `false` | Replay prompt directly without restarting the IDE process |
| `--project` | `-P` | string | `""` | Target project by index (`1`, `2`, ...) or name/slug |
| `--conversation` | `-c` | string | `""` | Target conversation ID (auto-resolved if omitted) |
| `--prompt` | `-p` | string | `"is-done"` | Prefix prompt template name or verification template |
| `--no-clipboard` | | boolean | `false` | Do not copy constructed prompt to system clipboard |
| `--dry-run` | `-d` | boolean | `false` | Preview restart actions and extracted payload without execution |
| `--help` | `-h` | boolean | `false` | Show help for agy rerun |

### Default Verification Template (`is-done`)

When `-p is-done` is specified (or used by default):

```text
Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes.
```

### Examples

```bash
# Rerun the prompt of project #1 with IDE restart and full picture/media replay
gitmap agy rerun 1

# Rerun the prompt of project #2
gitmap agy rerun 2

# Rerun project #3 without restarting the IDE process
gitmap agy rerun 3 --no-restart

# Preview rerun plan for project #1 without executing
gitmap agy rerun 1 --dry-run

# Legacy: Replay the very last prompt to OS clipboard
gitmap agy rerun last 1
```

---

## 2. agy list-prompts

Inspect past user prompts per project or across all projects. When `--projects <M>` is specified, prompt diffs and recent changes across the last M commit projects are aggregated into a temporary workspace document and opened in VS Code.

### Usage

```bash
gitmap agy list-prompts [N] [flags]
```

### Options & Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--all-projects` | | boolean | `false` | List historical prompts across all discovered projects |
| `--projects` | | integer | `0` | Showcase prompt changes across last M commit projects in VS Code |
| `--project` | | string | `""` | Filter prompts by project name prefix |
| `--json` | `-j` | boolean | `false` | Output prompt history in structured JSON format |
| `--help` | `-h` | boolean | `false` | Show help for agy list-prompts |

### Non-Admin VS Code Opener Contract

When `--projects <M>` is executed:

1. Aggregates prompt history and recent commit changes for the M most recently active projects into `.ai-memory/temp/prompts-diff.md`.
2. Launches VS Code using:
   ```bash
   code -n <path-to-temp-file>
   ```
3. **New Window Isolation (`-n`)**: Ensures the inspection session never disrupts or overwrites currently active editor windows.
4. **Non-Admin Execution**: VS Code is never launched with administrative or elevated privileges, preventing accidental elevation risks.

### Examples

```bash
# List all historical prompts for the current active project
gitmap agy list-prompts

# List last 10 prompts for the current project
gitmap agy list-prompts 10

# List last 10 prompts across all registered projects
gitmap agy list-prompts 10 --all-projects

# Inspect prompt changes for last 10 commit projects in a new VS Code window
gitmap agy list-prompts 10 --projects 10

# Search project by name prefix and display its last 10 prompts
gitmap agy list-prompts 10 --project my-api-service
```

---

## 3. agy scan

Scan directory trees for Git repositories and Antigravity workspaces, reporting prompt activity, recent prompt counts (last 24h and 7d), and historical prompt archives.

### Usage

```bash
gitmap agy scan [path] [flags]
```

### Options & Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--json` | `-j` | boolean | `false` | Output scan statistics in structured JSON format |
| `--backup` | | boolean | `false` | Archive discovered project prompts to local backup storage |
| `--help` | `-h` | boolean | `false` | Show help for agy scan |

### Examples

```bash
# Scan current workspace directory
gitmap agy scan

# Scan an explicit directory tree
gitmap agy scan d:/work/projects

# Output scan results in JSON format
gitmap agy scan --json
```

---

## 4. Remote Delegation & Triad Parity

Antigravity commands can be executed across local repositories and remote nodes seamlessly:

| Transport | Command | Target |
|-----------|---------|--------|
| **SSH** | `gitmap ssh exec agy <command>` | Targeted or all online SSH machines |
| **Cluster** | `gitmap cluster exec agy <command>` | Targeted or all cluster worker nodes |
| **SC** | `gitmap sc exec agy <command>` | Distributed fan-out across server-client nodes |
| **Local** | `gitmap exec agy <command>` | Uniform execution across all tracked local repositories |

### Remote Delegation Examples

```bash
# Execute agy rerun on remote SSH host "devbox"
gitmap ssh exec devbox agy rerun last 1 -p is-done

# Execute agy rerun across all cluster worker nodes
gitmap cluster exec agy rerun last 1 -p is-done

# Scan workspaces remotely via Servers-Clients (SC)
gitmap sc exec agy scan

# Run agy status locally across all tracked repositories
gitmap exec agy status
```

---

## 5. agy fix-pipeline (aef)

Extract failing CI/CD pipeline error logs, recent git commits, embed complete error logs and 4-part RCA prompt into `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`, inject fix directly into Google Antigravity, and queue follow-up verification.

### Usage

```bash
gitmap agy fix-pipeline [repo] [flags]
gitmap agy fix [repo] [flags]
gitmap aef [repo] [flags]
```

### Options & Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--all` | | boolean | `false` | Scan all tracked repositories and fix failing pipelines in batches |
| `--projects` | | integer | `3` | Number of failing projects to process per batch |
| `--limit` | | integer | `3` | Limit number of failing projects per batch run (synonym for `--projects`) |
| `--reset-batch` | | boolean | `false` | Reset multi-project batch cursor to the first project |
| `--no-inject` | | boolean | `false` | Skip direct Antigravity CLI/IDE injection |
| `--clip` | | boolean | `false` | Copy generated prompt directly to system clipboard |
| `--tempfile` | | string | `""` | Write prompt to custom `.ai-memory/temp/<filename>` |
| `--help` | `-h` | boolean | `false` | Show help for fix-pipeline |

### Examples

```bash
# Fix pipeline errors for current repository with direct AGY injection
gitmap agy fix-pipeline

# Batch fix failing pipelines across all projects (first 3 projects)
gitmap agy fix-pipeline --all

# Run again to automatically process the next batch of failing projects
gitmap agy fix-pipeline --all

# Specify custom batch limit of 5 projects
gitmap agy fix-pipeline --all --limit 5

# Reset batch cursor to start from the beginning
gitmap agy fix-pipeline --reset-batch
```

---

## 6. agy prompt

Send an arbitrary prompt to Antigravity using the Queue & Injection Protocol. Automatically routes prompt delivery based on IDE process presence and active conversation execution state:

- **Active IDE is RUNNING (busy)**: Appends the prompt to `.ai-memory/temp/agy-prompt-queue.json` without interrupting the running agent.
- **Active IDE is IDLE**: Stages the prompt to `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt` and copies content to the OS system clipboard.
- **No IDE Active**: Falls back to the Antigravity CLI background runner if available.

### Usage

```bash
gitmap agy prompt [slug/id/path] [prompt-text]
gitmap agy prompt [prompt-text]
```

### Examples

```bash
# Send prompt to current active project
gitmap agy prompt "Refactor error handling to use AppError envelopes"

# Send prompt to specific project path
gitmap agy prompt d:/work/my-service "Implement health check endpoint"
```

### Inspection Subcommands

| Subcommand | Usage | Description |
|------------|-------|-------------|
| `read` | `gitmap agy prompt read [conv-id]` | Read all user prompts from a conversation transcript |
| `ls` / `list` | `gitmap agy prompt ls [limit]` | List recent conversation prompts across workspaces |

### Inspection Examples

```bash
# Read prompts for the active workspace conversation
gitmap agy prompt read

# Read prompts for a specific conversation ID
gitmap agy prompt read 6c46400d-e873-4539-9d52-a96fee2786b6

# List the last 20 conversation prompts across all workspaces
gitmap agy prompt ls 20
```

---

## 7. agy queue

Manage the pending prompt queue stored in `.ai-memory/temp/agy-prompt-queue.json`.

### Usage

```bash
gitmap agy queue [subcommand]
gitmap agy q [subcommand]
```

### Subcommands

| Subcommand | Description |
|------------|-------------|
| `status` | Show active prompt, queued count, and last update timestamp (default) |
| `ls` / `list` | List all queued prompts with IDs, types, titles, and creation timestamps |
| `pop` | Pop the next queued prompt, mark it active, stage to file, and copy to clipboard |
| `clear` | Remove all queued and active items from the prompt queue |

### Examples

```bash
# View queue status
gitmap agy queue status

# List all queued prompts
gitmap agy queue ls

# Pop next prompt into clipboard and active staging file
gitmap agy queue pop

# Clear prompt queue
gitmap agy queue clear
```

---

## 8. agy status

Display Antigravity system status including running IDE process, active conversation execution state (IDLE vs RUNNING), pending prompt queue count, followed by the project status table.

### Usage

```bash
gitmap agy status
gitmap agy st
```

### Output Example

```text
● Antigravity System Status
  IDE Process:   RUNNING (PID: 14280)
  Conversation:  IDLE (ID: 6b7f91b5-6b40-45f8-87d6-fbd200388af2)
  Prompt Queue:  2 pending prompt(s)
```

---

## 9. agy ping

Ping the Antigravity desktop IDE environment and inspect health across 5 core dimensions:
- **IDE Executable**: Locates `Antigravity.exe` or candidate paths across Windows, Linux, and macOS.
- **IDE Process**: Detects whether the Antigravity desktop process is currently running (PID) or stopped (offline filesystem mode).
- **Filesystem Health**: Checks accessibility of Antigravity brain logs (`~/.gemini/antigravity/brain/`), project configs (`~/.gemini/config/projects/`), and conversation databases (`~/.gemini/antigravity/conversations/`).
- **Workspace Execution State**: Resolves matching conversation for the target workspace, reporting execution state (`IDLE` vs `RUNNING`) and total steps.
- **Prompt Queue Status**: Reports pending and active verification prompts in `.ai-memory/temp/agy-prompt-queue.json`.

### Usage

```bash
gitmap agy ping [flags]
gitmap agy check [flags]
```

### Options & Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--json` | `-j` | boolean | `false` | Output ping report in structured JSON format |
| `--workspace` | `-w` | string | `""` | Target workspace directory to inspect |
| `--help` | `-h` | boolean | `false` | Show help for agy ping |

### Output Example

```text
● Antigravity IDE Ping & Health Report
  IDE Executable: FOUND (~/.local/bin/antigravity)
  IDE Process:    RUNNING (PID: 14280, antigravity)
  Filesystem:     ACCESSIBLE
    • Brain Logs: ~/.gemini/antigravity/brain
    • Projects:   ~/.gemini/config/projects
  Workspace:      .
    • Status:     IDLE (ID: 6c46400d-e873-4539-9d52-a96fee2786b6, 42 steps)
  Prompt Queue:   1 active, 0 queued
  Overall Health: HEALTHY
```

---

## See Also

- [prompts-template](prompts-template.md) — Manage prompt prefix and verification templates
- [pipeline](pipeline.md) — CI/CD pipeline status, logs, and database management
- [storage](storage.md) — Inspect disk space, pipeline logs, and SQLite database storage
