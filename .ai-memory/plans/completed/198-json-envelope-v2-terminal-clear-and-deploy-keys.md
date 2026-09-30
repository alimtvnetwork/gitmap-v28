# Plan 198: JSON Envelope V2, Terminal Clear, Deploy-Keys Variants, and Friendly Import CLI

## 1. Objectives & Context

Implement seven core improvements across GitMap CLI:
1. **Terminal Clear Unification**: `gitmap terminal clear`, `gitmap clear terminal`, and `gitmap clear` must do the same thing: clear terminal screen buffer and command history/suggestions without asking for dev tools cache cleanup.
2. **Deploy-Keys Variants**: Wire all variants (`deploy all-keys`, `deploy all keys`, `deploy keys all`, `deploy key all`, `deploy keys`, `deploy-keys-all`, `deploy-all-keys`) to `cmdssh.RunSSHDeployKeysCLI`.
3. **Bare Import Guidance**: Fix `gitmap import` without args or `--confirm` to show interactive guidance instead of dumping `E9000` stack traces.
4. **Root Unknown Command Handling**: Stop dumping multi-page help on typos; show concise error message and suggestions.
5. **JSON Envelope V2 & Underscore Elimination**:
   - Eliminate underscores in `gitmap-ssh-nodes.json` and models (lowerCamelCase).
   - Add attributes: `importCommand`, `exportCommand`, `helpCommand`, `notes`, `version`, `gitmapVersion`, `workDirectory`, `defaultWorkDirectory`, `isWorkDirectoryApplied`, `isWorkDirectoryEnforced`.
   - Top-level `variables` support.
   - Enforce 1-based indexing for IDs in `gitmap-final.json`.
6. **Multi-JSON Merge Command (`gitmap merge-json` / `gitmap json merge`)**: Scan directory/pattern for JSONs of same type, deduplicate items, re-index with 1-based IDs, and emit a consolidated JSON document.
7. **Refactor `import-ssh-nodes.ps1`**: Remove hardcoded absolute paths, support relative/local execution, and display one-liner CLI command.
8. **Release Ceremony**: Version bump to `v6.415.0`, commit, tag, and push.

---

## 2. Implementation Steps

- [x] Task 1: Terminal Clear & Screen Buffer
  - Add screen buffer reset logic (`\033[H\033[2J\033[3J` + console clear) in `cli/osclean/` or `cli/cmdos/`.
  - Disambiguate `clear` vs `clean` in `cli/cmd/clean_dev_entry.go` and `cli/cmd/roottooling.go`.
  - Ensure `clear terminal`, `terminal clear`, and bare `clear` never prompt for devtools cache cleanup.
- [x] Task 2: Deploy-Keys Command Variants & Suggestions
  - Update `cli/cmdssh/ssh_deploy_router.go` to route `all-keys`, `all keys`, `keys all`, `key all`, `keys`.
  - Register `deploy-all-keys` and `deploy-keys-all` in root dispatch and completions (`rootsuggest.go`, `roottooling.go`).
- [x] Task 3: Bare `gitmap import` Friendly Guidance
  - Refactor `cli/cmd/importcmd.go` to render helpful suggestions menu when invoked with no args or without `--confirm`.
  - Suppress raw `E9000` apperror stack trace.
- [x] Task 4: Root Unknown Command Formatting
  - Update `cli/cmd/rootsuggest.go` `handleUnknownCommand` to remove `printUsage()` dump.
  - Render concise 2-line error with suggestions.
- [x] Task 5: JSON Envelope V2 & CamelCase Normalization
  - Update `cli/jsonenvelope/envelope.go`: add enriched attributes (`GitmapVersion`, `ImportCommand`, `ExportCommand`, `HelpCommand`, `Notes`, `WorkDirectory`, etc.) and top-level `Variables`.
  - In `cli/cmdssh/ssh_nodes_export_import.go`: convert all fields to lowerCamelCase with bidirectional unmarshaling for legacy snake_case.
  - Update fixtures and unit tests.
- [x] Task 6: Multi-JSON Merge Command (`gitmap merge-json`)
  - Create `cli/cmd/merge_json_cmd.go` with deduplication and 1-based ID re-indexing.
  - Wire command in `cli/cmd/roottooling.go` and `cli/cmd/root_cobra_completion.go`.
- [x] Task 7: Repo-Secrets Updates
  - Refactor `D:\work\repo-secrets\01-gitmap\import-ssh-nodes.ps1` to use `$PSScriptRoot` and relative paths.
  - Convert `D:\work\repo-secrets\01-gitmap\gitmap-ssh-nodes.json` to camelCase and Envelope V2.
  - Update `D:\work\repo-secrets\gitmap-final.json` and `07-final-network-machine/gitmap-final.json` to 1-based IDs and Envelope V2.
- [x] Task 8: Verification & Release
  - Complete release bump to `v6.415.0`.
  - Commit, tag, push, and consolidate plan.
