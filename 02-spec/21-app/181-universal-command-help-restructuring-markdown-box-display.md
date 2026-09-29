# Specification 181: Universal Command Help Restructuring, Markdown-to-Box Transformation, and Catalog Fallback

## 1. Executive Summary & Problem Statement

GitMap contains over 300 CLI commands and aliases. Historically, help was split across three distinct mechanisms:
1. **New Commands (v6.320+)**: Handcrafted modern double-line boxed terminal menus (`termhelp.RenderMenu` / `╔═╗`, `║ ║`, `╚═╝`) featuring clean sections: `Usage:`, `Description:`, `Subcommands:` / `Actions:`, `Flags:`, `Examples:`, and `Tips:`.
2. **Legacy Commands with Markdown**: Over 340 embedded markdown files in `cli/helptext/*.md`. When displayed, these previously rendered as raw markdown (`#`, `##`, markdown table pipes `|`, raw bash fences) or cyan block markers (`▌ `), lacking the unified look-and-feel of modern commands.
3. **Commands without Dedicated Help Files**: Certain aliases and subcommands lacked a markdown file in `cli/helptext/`, causing `gitmap <cmd> help` to bypass the help interceptor, fall through to `runDispatch`, and potentially fail with execution errors (e.g. "missing arguments") instead of showing documentation.

## 2. Requirements & Goals

### Requirement 1: Markdown-to-Modern-Box Restructuring Engine
Transform all markdown help rendered through `helptext.Print` / `render.RenderANSI` into the modern framed terminal display:
- **Title Banner**: The top Level 1 heading (`# gitmap <cmd>` or `# <cmd>`) is parsed and rendered inside the double-line box (`╔══════════════════════════════════════════════════╗`).
- **Clean Sections**: `## Subcommands`, `## Flags`, `## Examples`, `## Usage` are restructured into indented, colored section headers (`Usage:`, `Subcommands:`, `Flags:`, `Examples:`).
- **Table Transformation**: Markdown tables (`| Flag | Description |`) are parsed and rendered as clean, aligned two-column terminal listings without raw pipe characters.
- **Code Fences**: Backtick fences (````bash ... ````) are stripped and indented into syntax-highlighted terminal examples.

### Requirement 2: Dynamic Catalog Fallback for Every Command
- If a command is called with `gitmap <command> help` and lacks a dedicated `.md` file, the help interceptor in `cli/cmd/root.go` MUST NOT fall through to logic execution.
- Instead, it consults `cli/helptext/catalog.go` and `cli/constants/` to dynamically construct and display a complete modern framed box help screen with:
  - Framed header: `╔═╗ ... ╚═╝`
  - `Description:` from catalog registry
  - `Usage: gitmap <command> [arguments] [flags]`
  - Realistic `Examples:` demonstrating practical usage
  - Useful `Tips:`

### Requirement 3: Guarantee `gitmap <command> help` Never Fails or Mutates State
- For 100% of discovered commands in `completion.AllCommands()`, running `gitmap <cmd> help` or `gitmap <cmd> --help` MUST print the structured help screen and exit with code 0 without executing mutative operations.

### Requirement 4: Verification, Benchmarking, & CI Prompt
- Author comprehensive unit tests covering the markdown transformation engine and catalog fallback.
- Run `03-ai-scripts/42-clean-test-and-build-caches.py`.
- Create a dedicated verification prompt in `01-prompts/` to validate all commands with precision.
