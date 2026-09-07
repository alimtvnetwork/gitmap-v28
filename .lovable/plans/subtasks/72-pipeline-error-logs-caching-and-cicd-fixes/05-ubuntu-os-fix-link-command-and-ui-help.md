# Subtask 05: Ubuntu OS Fix-Link Command, Help Text & UI Help

## Scope
- Implement `gitmap os fix-link [path]` (and `gitmap os`, `gitmap fix-link`):
  - In `gitmap/constants/constants_cli.go`:
    - `CmdOS = "os"`
    - `CmdFixLink = "fix-link"`
    - `CmdFixLinkAlias = "fixlink"`
    - `SubCmdFixLink = "fix-link"`
  - In `gitmap/cmd/os.go` & `gitmap/cmd/os_fixlink.go`:
    - Parse flags: `--target`, `--force`, `--dry-run`, `--recursive`, `--json`.
    - Detect and repair broken symlinks on Ubuntu/Linux (and other Unix systems).
    - Handle standard links: `~/Desktop/SharedDirectories` -> `/mnt/hgfs`, `/usr/local/bin/gitmap` / `~/.local/bin/gitmap` binary link, and arbitrary `$path`.
    - Directory recursive scanning for broken symlinks.
  - In `gitmap/cmd/rootutility.go`:
    - Wire `constants.CmdOS` and `constants.CmdFixLink` into dispatch table.
  - In `gitmap/helptext/os.md` and `gitmap/helptext/fix-link.md`:
    - Full markdown documentation conforming to catalog requirements.
  - In `gitmap/helptext/catalog.go`:
    - Update catalog summary for `os` and `fix-link`.
  - In `gitmap/constants/constants_help.go` & `gitmap/cmd/rootusage_groups.go`:
    - Add UI help lines and display under `GET STARTED` / `ENVIRONMENT & TOOLS` or `INSTALLERS & TOOLS`.

## Files Touched
- `gitmap/constants/constants_cli.go`
- `gitmap/constants/constants_help.go`
- `gitmap/constants/cmd_constants_test.go`
- `gitmap/cmd/os.go`
- `gitmap/cmd/os_fixlink.go`
- `gitmap/cmd/os_fixlink_test.go`
- `gitmap/cmd/rootutility.go`
- `gitmap/cmd/rootusage_groups.go`
- `gitmap/helptext/os.md`
- `gitmap/helptext/fix-link.md`
- `gitmap/helptext/catalog.go`
- `gitmap/helptext/coverage_test.go`
