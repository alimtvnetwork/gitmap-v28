# Learned Memory 22: Nuclear Package Modularization (Phase 8) - cmdclone, cmdpull, cmdupdate

## Key Insights & Architectural Patterns

1. **Massive Domain Command Extraction (`cmdclone`, `cmdpull`, `cmdupdate`)**:
   - `cmdclone`: Extracted 83 files (`clone*.go`, `reclone*.go`, `directclone*.go`) out of `cli/cmd`. Reclone and directclone commands naturally share internal transport resolvers, destination validation, and clone history records with clone; grouping them into `cli/cmdclone` avoided unnecessary bridge forwarders.
   - `cmdupdate`: Extracted 23 files (`update*.go`) out of `cli/cmd`. Isolated Windows-specific process attributes (`processattr_windows.go` vs `processattr_other.go`) directly inside `cli/cmdupdate`.
   - `cmdpull`: Extracted 21 files (`pull*.go`, `push*.go`) out of `cli/cmd`. Push commands share identical parallel workers and destination targets with pull, making `cli/cmdpull` a cohesive cluster.

2. **Acyclic Domain Subpackage Hooks & Forwarders**:
   - Upward orchestration (such as `RunCodingGuidelinesInstall`, `CommitCodingGuidelines`, `RunCFRPPriorVersionPrivatize`, `EscapeCwdIfInside`, `RunStatus`, `RequireOnline`, `CreatePendingTask`) are decoupled via function pointers in `exports.go` (e.g., `cmdclone.RunCodingGuidelinesInstallFn`).
   - In `cli/cmd/clihelpers.go:init()`, all hooks are wired to existing `cmd` logic, preventing any import cycles (`cli/cmd` -> `cmdclone` -> `cli/cmd` is forbidden).
   - In `cli/cmd/rootflags.go`, `type CloneFlags = cmdclone.CloneFlags` provides transparent flag struct compatibility.

3. **Dead Code Elimination on Extraction**:
   - Moving internal helpers across packages frequently reveals dead or redundant helper wrappers (e.g. `parseCloneFlags`, `syncClonedReposToVSCodePM`, `buildClonePMPair`, `fileExists`, `checkHelp`).
   - Cleanly removing unused symbols keeps `golangci-lint` (strict) and `unused` guards completely green.

4. **Heavy Test Isolation**:
   - Subprocess tests that spawn long-running background loops (like `cli/probe/background_test.go`) were moved to `cli/tests/heavy_test/probe_background_test.go`.
   - Unit tests inside domain packages remain 100% pure in-memory, keeping lint and test times minimal.
