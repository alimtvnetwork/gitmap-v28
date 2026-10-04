# Specification 183: Typed JSON Envelope Architecture, Format Inspection (`which-format`), and Repo-Secrets Normalization

## 1. Executive Summary & Architectural Mission

GitMap operates with diverse JSON data formats across multiple subsystems: SSH cluster nodes, automation macros, commit-in/pull orchestration configurations, UI preferences, test inventories, and pipeline states.

To prevent format ambiguity, enable autonomous format discovery, and provide clear operational guidance, this specification establishes:
1. **Typed JSON Envelope Standard**: A unified envelope featuring two top-level sections:
   - `attributes`: Schema type, origin source, generator command/method, version, and generation timestamp.
   - `data`: Core subsystem payload.
2. **Format Inspection & Recommendation CLI (`gitmap which-format`)**:
   - Inspects single or multiple JSON files (or scans all `*.json` in a directory if no arguments are provided).
   - Identifies matched format types, displays suggested import commands, and describes exact system state impacts.
   - Flags unmatched/unsupported files clearly without crashing.
   - Provides a single-line batch command with `-y` bypass recommendations.
3. **Dual-Format Ingestion**: All subsystem importers transparently accept both the new typed envelope and legacy flat JSONs.
4. **Repo-Secrets Normalization**: Standardizes all JSON manifests in `./repo-secrets` to the typed envelope while strictly guarding secret data.

---

## 2. Typed JSON Envelope Specification

### 2.1 JSON Schema Contract

```json
{
  "attributes": {
    "type": "ssh-nodes",
    "source": "repo-secrets/01-gitmap/gitmap-ssh-nodes.json",
    "how": "gitmap ssh export",
    "version": "1.0",
    "timestamp": "2026-09-29T18:50:00Z"
  },
  "data": {
    "total_nodes": 5,
    "nodes": [ ... ]
  }
}
```

### 2.2 Standard Type Registry

| Type Identifier | Description | Importer Target | System State Change |
| :--- | :--- | :--- | :--- |
| `ssh-nodes` | Remote SSH cluster machine definitions | `gitmap sj import <file> -y` / `gitmap ssh import <file> -y` | Enrolls SSH nodes into `installation.db` and credential vault. |
| `macro` / `macro-bundle` | Automation macro commands and step trees | `gitmap macro import <file> -y` | Adds macro commands to SQLite `macros` repository table. |
| `commit-pull-config` | Commit-in and pull orchestration config | `gitmap commitin --config <file>` | Sets commit-in flags, auto-tag, and sync rules. |
| `ui-settings` | Terminal & Web UI layout and graphics options | `gitmap cmdui import-settings <file>` | Updates `~/.gitmap/ui_settings.json` preferences. |
| `test-inventory` | Test execution times and historical inventory | `gitmap test-inventory import <file>` | Updates `.ai-memory/test-inventory.json` execution metrics. |
| `pipeline-config` | Pipeline AI thresholds and runner configurations | `gitmap pipeline import <file>` | Updates pipeline monitoring configurations. |

---

## 3. CLI Command: `gitmap which-format`

### 3.1 Invocation Syntax
- `gitmap which-format [file1.json file2.json ...]`
- `gitmap which-format` (scans current working directory for `*.json`)
- `gitmap which-format --dir <path>`
- Aliases: `gitmap which format`, `gitmap format which`, `gitmap format inspect`

### 3.2 Output Layout (ANSI Boxed)
For each file:
- **Filename**: relative or provided path.
- **Match Status**: `MATCHED` or `UNMATCHED`.
- **Identified Type**: `attributes.type` or inferred legacy schema.
- **Suggested Import Command**: exact command line string.
- **System Impact**: summary of what will change upon import.
- **Batch Command**: single copy-pasteable command to import all matched JSONs.
- **Prompt Bypass Advice**: recommendation to append `-y` to skip interactive prompts.

---

## 4. Repo-Secrets Normalization & Security Constraints
- All JSON manifests in `./repo-secrets\01-gitmap\` and machine folders (`04-w1-machine`, `05-w2-machine`, etc.) are transformed to the typed envelope.
- Sensitive credentials, passwords, and private keys remain strictly inside `./repo-secrets`.
- Automated test suites in `gitmap` use isolated, sanitized mock fixtures in `cli/jsonenvelope/fixtures/` with zero production secrets.
