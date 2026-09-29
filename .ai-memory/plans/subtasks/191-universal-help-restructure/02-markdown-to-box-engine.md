# Subtask 191.2: Markdown-to-Modern-Box Transformation Engine

## Objective
Upgrade `cli/render/pretty.go` and `cli/render/prettypost.go` to reformat all embedded markdown help files into the modern double-line framed box display format.

## Requirements
1. Detect top-level `# <command>` or `# gitmap <command>` heading and format it as a double-line framed box banner.
2. Transform `## Section` into styled terminal headers (`Usage:`, `Subcommands:`, `Flags:`, `Examples:`).
3. Transform markdown tables (`| Flag | Description |`) into clean, aligned two-column lines without raw pipes.
4. Clean code fence syntax blocks and provide syntax highlighting.
