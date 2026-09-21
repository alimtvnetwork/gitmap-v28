# antigravity (agy, ag) — Antigravity Workspace Management

Manage Antigravity projects, project groups, prompt templates, prompt queues, undo/redo state, plugins, and settings.

## Subcommands

### Prompt Management & Dispatch

- `gitmap agy prompt-project <projectStartsWithName>` (aliases: `p`, `prompt-p`): Target an Antigravity project by name, ID, or path prefix and dispatch a prompt with optional template and text.
- `gitmap agy prompt`: Send or enqueue a prompt for the current git repository project (defaults to current repo root). Automatically registers the repository in Antigravity if missing, and dual-queues the Read Memory protocol before the user prompt.
- `gitmap agy prompt-with-name <templateName>` (alias: `pwn`): Dispatch a named prompt template to the current repository project with optional additional text.
- `gitmap agy prompt-txt "<text>"` (alias: `pt`): Dispatch direct text as a prompt to the current repository project.
- `gitmap agy prompt ls` (alias: `gitmap agy list prompts`): List all available prompt templates in a formatted table with Name, Description, Preview, and Invocation Example.

### Workspace & Project Management

- `gitmap agy group [add|rm|ls|show|export|import|prompt]`: Manage and prompt named groups of Antigravity projects.
- `gitmap agy undo`: Revert the last clear or project removal action from an automatic snapshot.
- `gitmap agy redo`: Reapply the undone Antigravity project state.
- `gitmap agy plugin [ls|install <slug>]`: List installed Antigravity plugins and install new plugins.
- `gitmap agy settings [export|import <file.json>]`: Export and import Antigravity configuration in JSON format.
- `gitmap agy ls`: List all discovered Antigravity projects with health status.
- `gitmap agy clear`: Remove stale or missing projects (snapshots state before deletion).
- `gitmap agy optimize-projects`: Deduplicate and optimize registered project entries.
- `gitmap agy export-projects <file.zip>`: Create a zip archive backup of all projects.
- `gitmap agy import-projects <file.zip>`: Restore projects from a zip archive backup.
- `gitmap agy install [ide|manager|cli|all]`: Install Google Antigravity Desktop IDE (default), Manager GUI, or CLI.
- `gitmap agm install`: Install Antigravity Manager GUI / Tools (or `gitmap install agm`).

## Prompt Command Flags

The prompt commands (`prompt`, `prompt-project`, `prompt-with-name`, `prompt-txt`) support the following flags:

- `-n`, `--name <name>`: Prompt template name to load from template store (e.g. `is-done`, `read-all`).
- `-t`, `--txt <text>`: Additional prompt text to include alongside template content.
- `--prefix`, `--pf`: Place template content before text (default mode, separated by two newlines: `<template>\n\n<text>`).
- `--suffix`, `--sf`: Place template content after text (separated by two newlines: `<text>\n\n<template>`).

## Auto-Registration & Dual-Queue Protocol

When executing `gitmap agy prompt`, `prompt-with-name`, or `prompt-txt` from a git repository:
1. **Repository Auto-Detection**: Project name is not required when running from a git repository.
2. **Auto-Registration**: GitMap inspects `~/.gemini/config/projects/` for an existing project matching the repository root. If not found, it automatically registers the repository in Antigravity workspaces.
3. **Dual-Queue Protocol**: GitMap enqueues two sequential prompts:
   - **Prompt 1**: Default Read Memory prompt (`Execute enhanced Read Memory protocol. Defensively load memory, specs, constraints, and pending plans before taking action.`) under type `read_memory`.
   - **Prompt 2**: The assembled user prompt under type `user_prompt`.
4. **Default Template**: If `-name` / `-n` is omitted, it defaults to `read-all` (same as `prompt-with-name read-all`).
5. **Formatting**: By default, template is placed before text separated by two newlines (`<template>\n\n<text>`). Use `--suffix` (`--sf`) for post-text placement (`<text>\n\n<template>`).
6. **Staging & Clipboard**: The assembled prompt is staged to `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt` and copied to the system clipboard.

## Installation Flags

The unified Antigravity installer engine supports the following multi-platform flags:

- `--force`, `-f`: Force reinstallation even if already present on the system.
- `--prefix`, `--dir <path>`: Custom installation directory prefix (overrides system defaults).
- `--dry-run`, `-n`: Simulate installation without downloading or modifying the host filesystem.
- `--version <ver>`: Specific release version to install (defaults to 2.13.0).
- `--yes`, `-y`: Automatic yes to interactive prompts.
- `--verbose`, `-v`: Enable detailed diagnostic output during installation.

## Examples

```bash
# 1. Target project by name/prefix with template name, text, and prefix mode (default, 2 newlines)
gitmap agy prompt-project my-app -name read-all -txt "what we want to add here" --prefix
gitmap agy p my-app -n read-all -t "what we want to add here" --pf

# 2. Prompt current repository (auto-adds project if missing, queues read-all first, then current prompt)
gitmap agy prompt -name read-all -txt "what we want to add here" --prefix
gitmap agy prompt -n read-all -t "what we want to add here" --pf

# 3. Prompt with template name as argument (same as prompt -name read-all)
gitmap agy prompt-with-name read-all -txt "what we want to add here" --prefix
gitmap agy pwn read-all -t "what we want to add here" --pf

# 4. Prompt direct text (defaults to read-all template prefix separated by 2 newlines)
gitmap agy prompt-txt "what we want to add here" --prefix
gitmap agy pt "what we want to add here" --pf

# 5. Suffix mode (text followed by 2 newlines and template)
gitmap agy prompt -name read-all -txt "what we want to add here" --suffix
gitmap agy prompt-txt "what we want to add here" --sf

# 6. List all available prompt templates in table
gitmap agy prompt ls
gitmap agy list prompts

# 7. Install Google Antigravity Desktop IDE and CLI
gitmap install antigravity
gitmap install agy

# 8. Manage project groups
gitmap agy group add frontend repo-a repo-b
gitmap agy group ls
gitmap agy group prompt frontend "Review styling and line gaps"

# 9. Undo a recent clear action
gitmap agy clear --all
gitmap agy undo

# 10. Inspect and install plugins
gitmap agy plugin ls
gitmap agy plugin install firebase-agent-plugin

# 11. Export and import settings
gitmap agy settings export my-settings.json
gitmap agy settings import my-settings.json
```
