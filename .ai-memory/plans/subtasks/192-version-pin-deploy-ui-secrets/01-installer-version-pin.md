# Subtask 192.1: Installer Version Pinning & Status in Installer Subsystem

## Spec Reference
- `02-spec/21-app/182-version-pinning-macro-deploy-ui-settings-secret-flags.md` §2.1

## Deliverables
1. **Schema & Model**:
   - Update `model.InstallerScript` in `cli/model/installer.go` to add `PinnedVersion string` field.
   - In `cli/store/`, support storing and retrieving pinned version metadata in SQLite Split-DB.
2. **Commands**:
   - Implement `gitmap installer pin <target> <version>` and alias `gitmap pin <target> <version>` / `gitmap pin installer <target> <version>`.
   - Implement `gitmap installer unpin <target>` / `gitmap unpin <target>`.
3. **Display**:
   - Update `executeInstallerLs` and `printInstallerTableHeader` in `cli/cmdinstaller/installer_ls.go` to display `PINNED` column or badge.
   - Display a summary footer indicating total installers and pinned count.
