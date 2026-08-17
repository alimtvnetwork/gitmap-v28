# Comprehensive Execution Plan

## 1. Sequence & Priority
We will execute the 5 pending tasks sequentially to avoid Git merge conflicts and database schema migration collisions. 

1. **`02-ssh-aware-clone.md`**: Smallest scope, adds transport persistence (`IdentifiedTransport`) to the SQLite DB and updates clone consumers. Must run first because it defines the foundational SSH logic.
2. **`03-reclone-transport-and-vscode-open.md`**: Relies on the `IdentifiedTransport` schema migration from `02`. Wires it into the `cfr` / `clone-now` flows. Adds `gitmap code`.
3. **`04-cfr-cg-os-aware-coding-guidelines.md`**: Modifies the `cfr` dispatch flow updated in `03`. Adds the `cg` modifier.
4. **`05-gitmap-improvements.md`**: The UI parallelization overhaul (Phases 4-6). Huge refactor of `gitmap clone gitmap.json`, push/pull. 
5. **`01-bulk-visibility-mapub-mapri.md`**: Massive 40-step feature adding wildcard commands and `GitMapRun` DB tracking. Best left for last to avoid blocking smaller wins.

## 2. Root Cause Protocols
Before any sub-task begins, we will analyze `.lovable/memory/last-failure.md` and explicitly identify the root cause of the bug (e.g. why `cfr` loses SSH transport) in the PR/commit message.

## 3. Autonomous Agent Flow
We will spawn specialized sub-agents (`self`) one at a time for each pending task file. Each agent will run end-to-end (code, test, commit) and then return control to the Orchestrator to begin the next task in the sequence.
