# Subtask 89.03: Help Text, UI Help, and Unit Tests

**Plan:** 89-qtorrent-utorrent-installers-and-config-options.md  
**Status:** Completed  

## 1. Description & Requirements
- Register top-level command constants in `gitmap/constants/constants_cli.go`:
  - `CmdExportConfig = "export-config"`
  - `CmdExportConfigAlias = "config-export"`
  - `CmdImportConfig = "import-config"`
  - `CmdImportConfigAlias = "config-import"`
- Update `topLevelCmds()` in `gitmap/constants/cmd_constants_test.go` and verify AST parity with `TestTopLevelCmdRegistryMatchesAST`.
- Run `go generate ./...` in `gitmap/` to generate `completion/allcommands_generated.go`.
- Author command help text:
  - `gitmap/helptext/export-config.md` (fenced code blocks, <= 120 lines, verified by golden tests).
  - `gitmap/helptext/import-config.md` (fenced code blocks, <= 120 lines, verified by golden tests).
  - `gitmap/helptext/install.md` updated with qBittorrent and uTorrent in Supported Tools.
- Update web UI command database in `src/data/commands.ts`:
  - `export-config` with flags and examples (`gitmap export-config uttorrent`, `gitmap export-config all ./configs/`).
  - `import-config` with flags and examples (`gitmap import-config uttorrent`, `gitmap import-config all ./configs/`).
- Implement unit tests:
  - `gitmap/cmd/config_export_import_test.go` (normalization, default filename resolution, round-trip write and restore, batch export, CLI help flag dispatch).
  - `gitmap/cmd/install_packages_test.go` (choco, winget, apt, brew packages, aliases, categories, and descriptions).
- Verify all quality gates:
  - `python linter-scripts/check-nested-ifs.py --changed-only`
  - `python linter-scripts/check-enum-and-boolean.py --changed-only`
  - `python linter-scripts/check-error-management.py --changed-only`
