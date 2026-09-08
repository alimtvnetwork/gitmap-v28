# Subtask 02: JSON & ZIP Export / Import for Installers

## 1. Description
Implement flexible export and import supporting both single JSON file format and ZIP archive bundles.
- `gitmap install export <slug> [-o output] [--format json|zip] [--all]`
- `gitmap install import <file.json|file.zip>`
- Provide dual parity with `gitmap installer export` and `gitmap installer import`.
- Auto-detect format based on file extension (`.json` vs `.zip`).
- Support single installer JSON and JSON array of installers for batch exports.

## 2. Files to Modify / Create
- `gitmap/cmd/install_export.go`: Implement `runInstallExport(args []string) error` and `runInstallImport(args []string) error`.
- `gitmap/cmd/install.go`: Route `export` and `import` in `runInstall`.
- `gitmap/cmd/installer_export.go` / `gitmap/cmd/installer_import.go`: Share or reuse JSON export and import logic.
- `gitmap/cmd/install_export_test.go`: Unit tests for JSON and ZIP export/import round-trip.

## 3. Invariants
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested if statements.
- Strict relative paths.
