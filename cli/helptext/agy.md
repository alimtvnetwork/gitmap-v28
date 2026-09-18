# gitmap agy

Manage Google Antigravity (AGY) workspaces, rerun historical prompts with prefix verification templates, inspect and diff prompts across projects with non-admin VS Code, scan project activity, and orchestrate remote execution across SSH, Cluster, SC, and Local targets.

## Synopsis

```bash
gitmap agy <subcommand> [flags]
```

## Aliases & Shorthands

- `gitmap agy`
- `gitmap ag`
- `gitmap antigravity`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `rerun last [N]` | Replay last N prompts with optional prefix prompt template |
| `list-prompts [N]` | List historical prompts for current project or across all projects |
| `scan` | Scan repositories, recent prompt activity counts, and prompt archives |
| `fix-pipeline [repo]` | Extract failing CI/CD pipeline errors and generate RCA prompt |
| `optimize-projects` | Remove stale or duplicate workspace configurations |
| `clean-cache` | Clear temporary Antigravity cache files |
| `open [path]` | Open target workspace in Antigravity |
| `pin-projects` | Manage pinned Antigravity projects |
| `plugins` | Inspect and configure Antigravity plugins |
| `reconcile` | Reconcile disk repositories with registered projects |

---

## Examples

```bash
# Rerun the last prompt with the default verification template
gitmap agy rerun last 1

# List prompts across all projects
gitmap agy list-prompts 10 --all-projects

# Inspect prompt changes for last 10 commit projects in VS Code
gitmap agy list-prompts 10 --projects 10

# Filter prompts for a specific project prefix
gitmap agy list-prompts 10 --project my-project

# Scan directory for repositories and prompt archives
gitmap agy scan
```

---

## 1. agy rerun

Replay the last N user prompts recorded in Antigravity conversation transcripts. Supports adding a prefix prompt template to verify or guide the execution.

### Usage

```bash
gitmap agy rerun last [N] [flags]
```

### Flags

```text
  -p, --prompt string     Prompt template name, prefix, or ID (default: "is-done")
      --no-clipboard      Do not copy the constructed prompt to clipboard
  -d, --dry-run           Preview constructed prompt without triggering execution
  -h, --help              Help for agy rerun
```

### Default Template (`is-done`)

When `-p is-done` is specified (or used by default):
```text
Is it done properly? Can we check properly the missing items from the task that is mentioned below? Please check it carefully. Do not make any mistakes.
```

### Examples

```bash
# Rerun the very last prompt with the default is-done verification prefix
gitmap agy rerun last 1

# Rerun the last 5 prompts prefixed with a custom template
gitmap agy rerun last 5 -p code-review

# Preview constructed prompt without clipboard copy
gitmap agy rerun last 3 --dry-run
```

---

## 2. agy list-prompts

Inspect past user prompts per project or across all projects. When `--projects <N>` is used, prompt diffs for the last N commit projects are formatted and opened in a new non-admin VS Code window.

### Usage

```bash
gitmap agy list-prompts [N] [flags]
```

### Flags

```text
      --all-projects      List historical prompts across all discovered projects
      --projects int      Showcase prompt changes across last N commit projects in VS Code
      --project string    Filter prompts by project name prefix
  -j, --json              Output prompt history in structured JSON format
  -h, --help              Help for agy list-prompts
```

### Non-Admin VS Code Inspection (`--projects <N>`)

Running `gitmap agy list-prompts 10 --projects 10` aggregates prompt history and commit changes for the 10 most recently committed projects into a structured temporary workspace document and launches:
```bash
code -n <temp-file>
```
- **New Window Isolation (`-n`)**: Never disrupts or overwrites currently open VS Code editor windows or workspaces.
- **Non-Admin Security**: Explicitly avoids administrative elevation for safe inspection.

### Examples

```bash
# List last 10 prompts for the current active project
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

### Examples

```bash
# Scan current workspace directory
gitmap agy scan

# Scan an explicit directory tree
gitmap agy scan d:/work/projects
```

---

## 4. Remote Execution & Triad Parity

Run Antigravity commands seamlessly across local and remote machines:

```bash
# Execute agy command on SSH nodes
gitmap ssh exec agy rerun last 1 -p is-done

# Execute agy command across all cluster nodes
gitmap cluster exec agy rerun last 1 -p is-done

# Execute agy command via Servers-Clients (SC) delegation
gitmap sc exec agy scan

# Execute agy command locally across all tracked repositories
gitmap exec agy status
```

---

## 5. agy fix-pipeline (aef)

Extract failing CI/CD pipeline error logs, recent git commits, embed complete error logs and 4-part RCA prompt into `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`, inject fix directly into Google Antigravity, and queue follow-up verification. Supports multi-project parallel batching across repositories with persistent batch cursor tracking.

### Usage

```bash
gitmap agy fix-pipeline [repo] [flags]
gitmap agy fix [repo] [flags]
gitmap aef [repo] [flags]
```

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

