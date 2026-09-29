# Specification 184: Import-All-JSON, What-Configs (`wc`), and Fleet Deploy Hardening

## 1. Executive Summary & Architectural Mission

GitMap environments routinely handle heterogeneous configuration manifests spanning multiple subsystems: SSH cluster nodes, automation macros, commit-in/pull orchestration configurations, UI preferences, and split-DB templates.

To streamline bulk operations, eliminate manual file-by-file imports, diagnose arbitrary configuration files in a repository, and prevent Windows remote command execution overflows, this specification establishes:
1. **`gitmap import-all-json` Engine**:
   - Universal bulk importer supporting file patterns, directory scanning, and globs (`*`, `*.json`, `specific.json`).
   - Deep format inspection to infer envelope and subsystem types (`ssh-nodes`, `macro`, `commit-pull-config`, `ui-settings`, `templates`).
   - Execution dispatching to subsystem handlers with preview (`--dry-run`) and automatic confirmation (`-y`).
2. **`gitmap what-configs` (`wc`) Inspection Utility**:
   - Polyglot configuration inspection utility analyzing JSON, YAML, TOML, SQLite databases, INI, and ENV files.
   - Structured tabular ANSI report identifying format category, management impact, and exact recommended CLI commands.
3. **Fleet Deployment Hardening & Windows Command-Line Length Guard**:
   - Mitigation of Windows `cmd.exe` 8,191-character limit (`The command line is too long.`) during remote fleet deployments.
   - Direct file streaming (`streamDirectToRemote`) to `.gitmap/fleet_config_deploy.json` for payloads exceeding 1,500 characters, followed by decoupled CLI execution and deferred cleanup.
4. **SSH Join Parser Disambiguation & Bi-directional Communication Guard**:
   - Inverted parameter detection allowing both `gitmap ssh join <alias> <user@host>` and `gitmap ssh join <user@host> <alias>`.
   - Safe parsing of passwords containing special characters (such as `@`) without triggering erroneous multi-target enrollment.
   - PowerShell-compliant separator syntax for Windows remote verification commands.
5. **`gitmap deploy import` Command Discovery**:
   - Native integration of `deploy import` and `deploy import-json` within `gitmap deploy` router and rich help displays.

---

## 2. CLI Command: `gitmap import-all-json`

### 2.1 Invocation Syntax
- `gitmap import-all-json [file_or_pattern ...]`
- `gitmap import-all-json *` (inspects and imports all matching JSON files in directory)
- `gitmap import-all-json *.json` (targets only `.json` files)
- `gitmap import-all-json path/to/config.json` (imports a specific file)
- Flags:
  - `--dry-run`: Previews matched types and planned execution steps without mutating system state.
  - `-y`, `--yes`: Confirms all prompt bypasses for batch execution.
- Aliases: `import-all-json`, `importalljson`, `import-json-all`, `importall`.

### 2.2 Subsystem Importers Mapping

| Format Identifier | Target Subsystem Importer | State Mutation |
| :--- | :--- | :--- |
| `ssh-nodes` | `RunSSHNodesImportJSON(path, "-y")` | Persists remote nodes into `installation.db` and credential vault. |
| `macro` / `macro-bundle` | `ImportMacroFile(path, true)` | Inserts macro step definitions into `macros` SQLite table. |
| `commit-pull-config` | `ImportCommitPullConfig(path)` | Updates default commit-in and fast pull configurations. |
| `ui-settings` | `ImportSettingsFromFile(path)` | Updates terminal and web UI layout/palette preferences. |
| `templates` | `ImportTemplates(path)` | Imports declarative state templates into split-DB storage. |

---

## 3. CLI Command: `gitmap what-configs` (`wc`)

### 3.1 Invocation Syntax
- `gitmap what-configs [pattern_or_file ...]`
- `gitmap wc *` (scans all config files: `.json`, `.yaml`, `.yml`, `.toml`, `.db`, `.sqlite`, `.env*`, `.ini`)
- `gitmap wc *.json`
- `gitmap wc specific.json`
- Aliases: `what-configs`, `wc`, `what-config`, `whatconfigs`.

### 3.2 Detection & Classification Engine

Files are inspected and categorized into:
- **GitMap JSON Envelopes & Flat Configs**: Detects envelope `attributes.type` or structural signatures (`nodes`, `macros`, `pull_policy`).
- **YAML / YML**: CI/CD workflows, Docker Compose, or Kubernetes definitions.
- **SQLite / DB**: GitMap split-DBs (`gitmap.db`, `installation.db`, `repodb/*.db`) or application SQLite databases.
- **TOML**: Go/Rust package configurations (`Cargo.toml`, `Gopkg.toml`).
- **ENV**: Environment variable files (`.env`, `.env.local`, `.env.example`).
- **INI**: Desktop, gitconfig, and system configuration files.

### 3.3 Output Layout
Renders an ANSI formatted table detailing:
- **FILE**: Path to inspected file.
- **CATEGORY**: Detected configuration format.
- **RECOMMENDED COMMAND**: Recommended GitMap or CLI command for manipulation.
- **SYSTEM IMPACT**: Scope of changes if applied to the environment.

---

## 4. Remote Fleet Deployment Hardening

### 4.1 Windows Command-Line Length Limitation
Windows `cmd.exe` imposes an 8,191-character buffer limit on command execution strings. Encoding multi-node fleet manifests as inline base64 arguments within remote execution scripts causes `The command line is too long.` errors on Windows fleet nodes.

### 4.2 Streamed File Transfer Mitigation
- When the serialized base64 payload size is >= 1,500 bytes (or when target node OS is Windows), `DeployConfigSSHToNode` bypasses inline script execution.
- The configuration payload is streamed directly to `.gitmap/fleet_config_deploy.json` in the user's home directory via `streamDirectToRemote`.
- GitMap invokes `gitmap ssh nodes import-json .gitmap/fleet_config_deploy.json` on the remote node.
- A deferred cleanup command securely removes `.gitmap/fleet_config_deploy.json` upon completion.

---

## 5. SSH Join Disambiguation & OS Syntax

### 5.1 Symmetric Positional Argument Parsing
The SSH enrollment parser (`sshjoin_cmd.go`, `ssh_parser.go`, `sshjoin_add_pass_cmd.go`) resolves positional parameters symmetrically:
- If positional arg 0 matches `user@host:port` (or IP address), arg 1 is treated as the alias name.
- If positional arg 0 is an alias name, arg 1 is treated as the host target address.

### 5.2 Password Special Characters Protection
Passwords ending with `@` (e.g. `secret123@`) are protected against naive `@` splitting. The parser validates that the token following `@` resolves to a valid IP or hostname before treating the argument as a remote connection target.

### 5.3 Windows Command Separators
Bidirectional SSH connectivity confirmation (`sshjoin_enroll.go`) detects target OS. For Windows targets, chained commands use PowerShell-compatible syntax (`command1; if ($?) { command2 }`) rather than POSIX `command1 && command2` or `||`.

---

## 6. Verification and Regression Testing

- `cli/cmd/import_all_json_cmd_test.go`: Validates dry-run and bulk execution across test envelope fixtures.
- `cli/cmd/what_configs_cmd_test.go`: Validates format classification and table rendering for diverse file types.
- Remote deployment verified on live fleet nodes (`main`, `w3`).
- All `cli/cmdssh/...` and `cli/cmd/...` test suites pass.
