# Execution Plan: 202 GitMap Pull Errors and SplitDB Error Storage

Spec Reference: [02-spec/21-app/202-gitmap-pull-errors-splitdb/01-overview.md](../../../02-spec/21-app/202-gitmap-pull-errors-splitdb/01-overview.md)

## Goal
Fix the `gitmap fix` untracked file stash collision, architect a SplitDB engine for storing repository errors, and expose visibility through CLI fleet commands.

## Architecture Context
The application utilizes an SQLite SplitDB architecture. The root error DB should maintain a high-level manifest of errors (where they occur), and repository-specific SQLite databases will hold the detailed stack traces and contextual information. Fleet nodes must be capable of emitting this error payload in JSON, and the local orchestrator node must parse, persist, and render these errors in the terminal.

## Tasks & Subtask Mapping

- **Task-02: Fix Stash Pop Untracked File Collision Root Cause**
  - Mapped to: `subtasks/202-gitmap-pull-errors-splitdb/01-fix-stash-pop-collision.md`
- **Task-03: SplitDB Error Storage Architecture**
  - Mapped to: `subtasks/202-gitmap-pull-errors-splitdb/02-splitdb-error-storage.md`
- **Task-04: GitMap CLI Error Visibility & Fleet Commands**
  - Mapped to: `subtasks/202-gitmap-pull-errors-splitdb/03-cli-error-fleet-commands.md`
