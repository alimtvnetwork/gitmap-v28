# Subtask 01: Install Context Entry & Antigravity Open

- **Plan**: `334-gitmap-installation-and-ui-enhancement-instructions`
- **Subtask**: 01
- **Status**: PENDING

---

## Action Items
1. In `cli/cmdinstall/installctxentries.go`, add `05_install_ctx` in `installChildren()` pointing to `["install", "context-menu"]`.
2. Update `cli/cmdinstall/installctxentries_argv_test.go` to register expected entry `60_install/05_install_ctx`.
3. In `cli/cmdopen/open_target.go`, add `isAntigravityTarget(arg string) bool` checking `agy`, `antigravity`, `ag`.
4. In `cli/cmdopen/open.go`, dispatch `isAntigravityTarget` to `cmdagy.RunAgyOpen(path)`.
5. In `cli/cmd/roottooling.go`, register `open-agy` and `open-antigravity` aliases.
6. Verify via `go test ./cli/cmdinstall` and `./cli/cmdopen`.
