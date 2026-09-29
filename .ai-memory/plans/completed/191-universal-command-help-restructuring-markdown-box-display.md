# Master Plan 191: Universal Command Help Restructuring, Markdown-to-Box Transformation, and Catalog Fallback

## 1. Overview & Objectives
Restructure all legacy and markdown-rendered command help in GitMap to display with the modern terminal layout (double-line framed box, clean aligned two-column flags and subcommands, examples, and tips), and ensure 100% of commands in `completion.AllCommands()` cleanly respond to `gitmap <command> help` without side effects or execution errors.

## 2. Subtasks Breakdown

- [x] **191.1: Catalog Audit & Dynamic Fallback Help Generator**
  - Extended `cli/cmd/help_dynamic.go` and `cli/cmd/helpcheck.go` with dynamic fallback menu generator `RenderDynamicCommandHelp(command string) bool`.
  - When `gitmap <command> help` is invoked for any command without a `.md` file, dynamically render a modern boxed help screen.

- [x] **191.2: Markdown-to-Modern-Box Transformation Engine**
  - Upgraded `cli/render/prettypost.go` to detect Markdown Level 1 headings and render them inside double-line framed box banners (`termhelp` style).
  - Converted markdown table pipes (`| Flag | Description |`) into clean, padded two-column output.
  - Transformed `## Sections` into highlighted terminal section titles (`Usage:`, `Subcommands:`, `Flags:`, `Examples:`).

- [x] **191.3: Update Universal Interceptor in `cli/cmd/root.go`**
  - In `tryInterceptCommandHelp`, ensured that after checking core rich topics and `.md` files, it checks `RenderDynamicCommandHelp`.
  - Guaranteed 100% coverage: No command shall fall through to live execution when the second argument is `help` or a help flag.

- [x] **191.4: Unit Tests & Verification**
  - Added `TestRenderANSIBoxBannerAndTable` in `cli/render/pretty_test.go` and `TestRenderDynamicCommandHelp` in `cli/cmd/help_dynamic_test.go`.
  - Verified across representative commands (`backup help`, `clean help`, `clone help`, `add help`).

- [x] **191.5: Workstation Cache Clean & Automated Verification Prompt**
  - Executed `03-ai-scripts/42-clean-test-and-build-caches.py`.
  - Authored `01-prompts/22-verify-markdown-restructure-and-universal-help.md`.

- [x] **191.6: Minor Release Orchestration & CI Verification**
  - Release `v6.400.0` using `03-ai-scripts/29-release-orchestrator.py`.
  - Monitored and verified all GitHub Actions workflows are 100% green.
