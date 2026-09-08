# Subtask 01: `gitmap install ls` Version Extraction, Gitmap Release Info & New Tools

## 1. Description
Enhance `gitmap install ls` to display real installed versions (e.g. `v22.14.0`, `1.24.1`, `2.47.0`) instead of generic `"found"`, register `gitmap` CLI at the top of Core Tools with `v6.199.0` and recent release tag summary, and register newly added developer tools: `composer`, `wp-cli`, `open-vm-tools`.

## 2. Files to Modify
- `gitmap/constants/constants_install.go`:
  - Add `ToolGitmap = "gitmap"`
  - Add `ToolComposer = "composer"`
  - Add `ToolWpCli = "wp-cli"`
  - Add `ToolOpenVmTools = "open-vm-tools"`
  - Add entries to `InstallToolDescriptions`
  - Add entries to `InstallToolCategories` (`ToolGitmap` at top of Core Tools, `ToolComposer` in Languages, `ToolWpCli` in Core, `ToolOpenVmTools` in DevOps)
- `gitmap/cmd/installlist.go`:
  - Implement fast version probing `detectToolVersion(tool string) string` with <= 300ms context timeout.
  - In `resolveToolStatus`: if tool is `constants.ToolGitmap`, return `constants.Version`. If binary is in PATH, probe `detectToolVersion(tool)`.
  - In `printInstallListGrouped`: render header banner with current Gitmap version (`v6.199.0`) and recent release tags.
- `gitmap/cmd/install_unit_test.go`:
  - Add tests verifying `gitmap` resolves to `constants.Version` and real tool versions are detected.

## 3. Invariants
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested ifs.
- Fast execution (<= 300ms context timeout per probe so `install ls` completes in under 1 second).
