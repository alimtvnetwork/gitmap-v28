# Subtask 04: Verification and Release

- **Plan**: `334-gitmap-installation-and-ui-enhancement-instructions`
- **Subtask**: 04
- **Status**: PENDING

---

## Action Items
1. Run FastGate pre-commit checks (`gitmap py 03-ai-scripts/50-fastgate.py`).
2. Run targeted unit tests (`go test ./cli/cmdinstall`, `./cli/cmdopen`, `./cli/cmdchrome`).
3. Execute minor version bump (`37-bump-version.py -t minor`).
4. Commit release via `gitmap cpr "v6.617.0 installation right-click command open antigravity app-window browser and ui animations"`.
5. Tag and push to upstream repository.
6. Recompile local binary and register context menu with `gitmap context-menu install`.
7. Move plan 334 to `.ai-memory/plans/completed/`.
