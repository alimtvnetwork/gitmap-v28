# Subtask 03: Dynamic Custom Tools Section in Gitmap Install LS

## 1. Description
Integrate custom installer records from SQLite into `gitmap install ls` so that newly added installers appear dynamically in a dedicated `Custom Tools` section.
- Query `store.DB.ListInstallers()` during `printInstallListGrouped()`.
- If custom installers exist, render the category block `Custom Tools`.
- Display status dot (● if verify passes / present in split DB; ○ otherwise), tool slug/name, version, and description.
- Preserve table alignment with existing `gitmap install ls` columns.

## 2. Files to Modify / Create
- `gitmap/cmd/installlist.go`: Add `loadCustomInstallers()` and render `Custom Tools` block.
- `gitmap/cmd/install_unit_test.go`: Add test verifying custom installers appear in `install ls`.

## 3. Invariants
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested if statements.
- Strict relative paths.
