# gitmap vscode repair (alias: fix, doctor)

Diagnose and repair VS Code binary startup crashes and Project Manager JSON integrity.

## Usage

```bash
gitmap vscode repair [flags]
gitmap vscode fix [flags]
gitmap vscode doctor [flags]
```

## Flags

- `-k, --kill`: Terminate lingering/stuck `Code.exe` or `inno_updater` processes before running diagnostics.
- `-f, --force`: Attempt reinstall via `winget` if binary file mismatches cannot be repaired locally.
- `-d, --dry-run`: Preview diagnostic checks and repairs without modifying configuration files on disk.

## Actions Performed

1. **Process Cleanup**: Identifies and terminates orphaned VS Code processes that lock executables during background auto-updates.
2. **Project Manager JSON Validation**: Validates syntax in `projects.json`, creates automatic `.bak` backups on corruption, ensures required schema fields (`paths`, `tags`, `enabled`, `profile`), and synchronizes modern `globalStorage` and legacy `User` paths.
3. **Binary Integrity Check**: Validates `Code.exe` / `code.cmd` execution. Detects missing commit resource subfolders caused by `win32VersionedUpdate` race conditions (`icu_util.cc` data errors) and syncs missing commit folders between installation paths.
4. **Winget Fallback**: Runs official winget installer repair if binaries cannot be recovered automatically.
