# Specification: Coding Guidelines Multi-Repo, Version Status & Status Dirty Filter

## 1. Overview

This specification covers:
1. **Coding Guidelines Status & Version**: `gitmap cg version` and `gitmap cg status` reading `"coding-guidelines"` from `version.json`.
2. **Multi-Repo Auto-Discovery**: When `gitmap cg install/update`, `gitmap pull`, or `gitmap push` is executed in a folder containing multiple Git repositories, auto-discover all child Git repositories.
3. **Flexible Target Specifiers**: Support repository folder paths, aliases, numeric IDs, and remote Git URLs for `gitmap cg repo install <targets>` and `gitmap cg repo update <targets>`.
4. **Selective Update**: `gitmap cg update` only updates repos with existing `"coding-guidelines"` configuration in `version.json`.
5. **Status Dirty Filter**: `gitmap status --dirty` / `gitmap status --only-dirty` displaying only repos with unstaged, staged, or unpushed/unpulled changes.
6. **Markdown Help Output**: Support `--markdown` on help commands to output formatted markdown examples.

## 2. version.json Schema for Coding Guidelines

```json
{
  "version": "6.104.0",
  "coding-guidelines": {
    "version": "v24.0.0",
    "installed_at": "2026-08-26T14:30:00Z",
    "status": "active"
  }
}
```

## 3. Command Routing & Options

- `gitmap cg version [repo...]`: Reads and prints the version of coding guidelines installed.
- `gitmap cg status [repo...]`: Displays status table across repositories showing whether coding guidelines are installed and their version.
- `gitmap cg repo install <path|alias|id|url>`: Installs to specified repos.
- `gitmap cg repo update <path|alias|id|url>`: Updates only if coding guidelines are currently installed in `version.json`.
- `gitmap status --dirty`: Filters status table to only dirty repositories.
- `gitmap pull` / `gitmap push` (in parent folder): Auto-discovers child repositories and performs batch operation.
