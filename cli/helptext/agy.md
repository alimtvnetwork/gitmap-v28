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
| `prompt [slug/path] [text]` | Send a prompt to Antigravity IDE/CLI using Queue & Injection Protocol |
| `prompt read [conv-id]` | Read and display user prompts for a conversation from Antigravity brain |
| `prompt ls [limit]` | List conversation prompts across workspaces (alias: `list`) |
| `queue [status\|ls\|clear\|pop]` | Manage pending prompt queue entries |
| `status` | Show Antigravity system status (IDE process, IDLE vs RUNNING, prompt queue) |
| `rerun last [N]` | Replay last N prompts with optional prefix verification prompt template (`-p`) |
| `list-prompts [N]` | List historical prompts for current project or across all projects |
| `scan [path]` | Scan repositories, recent prompt activity counts (24h/7d), and prompt archives |
| `fix-pipeline [repo]` | Extract failing CI/CD pipeline errors and generate 4-part RCA prompt (alias: `aef`) |
| `open [path]` | Locate Antigravity IDE and open target workspace |
| `optimize-projects` | Remove stale or duplicate workspace configurations |
| `clean-cache` | Clear temporary Antigravity cache files and transcripts |
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
```

---

## 1. agy rerun

Replay the last N user prompts recorded in Antigravity conversation transcripts (`transcript.jsonl`). Supports prepending a prefix template to verify or guide execution.

### Usage

```bash
gitmap agy rerun last [N] [flags]
```

### Options & Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--prompt` | `-p` | string | `"is-done"` | Prompt template name, prefix, or ID |
| `--no-clipboard` | | boolean | `false` | Do not copy constructed prompt to system clipboard |
| `--dry-run` | `-d` | boolean | `false` | Preview constructed prompt without triggering execution |
| `--help` | `-h` | boolean | `false` | Show help for agy rerun |

### Default Verification Template (`is-done`)

When `-p is-done` is specified (or used by default):

```text
Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes.
```

### Prompt Template Resolution Order

1. **Exact Slug or ID Match**: Matches template ID (e.g., `-p is-done`, `-p code-review`).
2. **Prefix Match**: Matches any registered template starting with the given string (e.g., `-p is` resolves to `is-done`).
3. **Built-in Fallback**: Falls back to the canonical `is-done` verification prompt.

### Examples

```bash
# Rerun the very last prompt with the default is-done verification prefix
gitmap agy rerun last 1

# Rerun the last 5 prompts prefixed with a custom template
gitmap agy rerun last 5 -p code-review

# Preview constructed prompt without copying to clipboard or dispatching
gitmap agy rerun last 3 --dry-run
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
  IDE Executable: FOUND (C:\Users\Administrator\AppData\Local\Programs\antigravity\Antigravity.exe)
  IDE Process:    RUNNING (PID: 14280, antigravity.exe)
  Filesystem:     ACCESSIBLE
    • Brain Logs: C:\Users\Administrator\.gemini\antigravity\brain
    • Projects:   C:\Users\Administrator\.gemini\config\projects
  Workspace:      D:\work\gitmap
    • Status:     IDLE (ID: 6c46400d-e873-4539-9d52-a96fee2786b6, 42 steps)
  Prompt Queue:   1 active, 0 queued
  Overall Health: HEALTHY
```

---

## See Also

- [prompts-template](prompts-template.md) — Manage prompt prefix and verification templates
- [pipeline](pipeline.md) — CI/CD pipeline status, logs, and database management
- [storage](storage.md) — Inspect disk space, pipeline logs, and SQLite database storage

