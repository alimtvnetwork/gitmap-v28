# Goal

Generate a generic, tool-agnostic AI Instruction Specification (`ssh-commands.md`) at the root of the repository to guide any AI agent in implementing a comprehensive SSH and SSH Profile management system for a CLI tool (referred to as `<cli>`). 

## 50/50 Strategy Allocation

- **Phase 1 (Planning)**: We are currently writing this detailed execution plan.
- **Phase 2 (Execution)**: We will write the `ssh-commands.md` file to the root, commit it to the repository, trigger a version bump, update the changelog, and finally output the markdown text directly to the chat along with the End of Run Summary and Compliance Checklists.

## Execution Steps

1. **Draft `ssh-commands.md`**:
   - Write a strict AI-to-AI prompt instruction.
   - Section 1: Purpose & mental model (SQLite as source of truth, `~/.ssh/` for private keys, `<cli-dir>/` for repo bindings).
   - Section 2: Terminal output reference (Generic rendering contract with symbols, colors, alignment, fallback mechanisms).
   - Section 3: Read-and-adopt behavior (`<cli> ssh` logic, fallback chain for email).
   - Section 4: Clipboard requirement (OS-specific clipboard binaries).
   - Section 5: Base SSH subcommands (`create`, `ls`, `st`, `cp`, `view`, `rm`, `config`, `join`).
   - Section 6: Profiles command tree (`profiles create`, `set`, `set-repos`, `rm`, `github-desktop`, `export`, `import`).
   - Section 7: Help system contract (deep nesting with `>>`, UI/Terminal synchronization).
   - Section 8: Data model (SQLite schemas, JSON schemas).
   - Section 9: Implementation Checklist (Acceptance criteria).
2. **Commit & Release**:
   - Add `ssh-commands.md`.
   - Run python bump version script.
   - Update `readme.md`, `version.json`, `changelog.md`, `.lovable/what-to-read.md`.
3. **Output**:
   - Print the exact markdown contents to the chat.
   - Provide the End of Run Summary.
