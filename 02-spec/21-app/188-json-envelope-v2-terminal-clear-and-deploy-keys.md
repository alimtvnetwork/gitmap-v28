# Spec 188: JSON Envelope V2, Terminal Clear, Deploy-Keys Variants, and Friendly Import CLI

## 1. Overview & Context

This specification formalizes four tightly-integrated developer ergonomics and data serialization enhancements across GitMap:
1. **Terminal Clear Unification**: `gitmap terminal clear` and `gitmap clear terminal` (along with bare `gitmap clear`) provide identical behavior: clearing terminal screen buffer (ANSI + OS console) and clearing/reseeding shell command history and suggestions without erroneously prompting for developer tools/devtools cache cleanup.
2. **Deploy-Keys Variants**: Comprehensive routing for all natural user phrasing: `gitmap deploy all-keys`, `gitmap deploy all keys`, `gitmap deploy keys all`, `gitmap deploy key all`, and `gitmap deploy keys`. All dispatch to mesh public key collection and distribution, with full terminal tab completion.
3. **Bare Import Guidance & Suppression of Unfriendly Stack Traces**: Executing bare `gitmap import` without arguments or `--confirm` no longer dumps internal stack traces (`E9000`). Instead, it displays a friendly interactive guidance catalog showing available envelope types, format inspection commands (`gitmap which-format`), bulk import commands (`gitmap import-all-json`), and explicit confirmation requirements.
4. **Root Unknown Command Handling**: Typos at the root CLI level no longer dump multi-page help catalogs. GitMap displays a concise suggestion indicating the exact typo, nearest "Did you mean" recommendations, and instructions on how to view help.
5. **Typed JSON Envelope V2 & Underscore Elimination**:
   - Complete removal of underscores from JSON models (e.g., in `gitmap-ssh-nodes.json`, snake_case fields like `auth_type`, `ip_address`, `worker_id` transition to lowerCamelCase `authType`, `ipAddress`, `workerId`).
   - JSON envelope attributes section expanded to include explicit `importCommand`, `exportCommand`, `helpCommand`, `notes`, `version`, `gitmapVersion`, and working directory settings (`workDirectory`, `defaultWorkDirectory`, `isWorkDirectoryApplied`, `isWorkDirectoryEnforced`).
   - Top-level `variables` dictionary support.
   - Enforce 1-based indexing for machine IDs in `gitmap-final.json`.
6. **Multi-JSON Merge CLI (`gitmap merge-json` / `gitmap json merge`)**: Scan directories/patterns for JSON files of identical format type, deduplicate items, re-index with 1-based IDs, and emit a consolidated JSON document with fresh envelope metadata.
7. **Generic SSH Import Script & One-Liner**: Refactor `import-ssh-nodes.ps1` to eliminate hardcoded absolute paths, support portable local execution, and surface one-liner GitMap CLI commands.

---

## 2. Architecture & Design Principles

### 2.1 Terminal Clear vs Dev Clean Disambiguation

Previously, `gitmap clear` routed to `runCleanTopLevel`, which defaulted to `cmdos.RunOSDevClean` when no subcommand matched, prompting:
`Proceed with dev tools cache cleanup? Type 'yes' to continue:`
This conflicted with the universal convention where operators type `clear` to clear their terminal screen and history.

**New Resolution Matrix**:
| CLI Invocation | Target Action | Screen Buffer Clear | History / Suggestion Clean | Devtools Cache Clean | Prompt |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `gitmap clear` | Terminal Clean | Yes | Yes | No | Yes (unless `-y`) |
| `gitmap clear terminal` | Terminal Clean | Yes | Yes | No | Yes (unless `-y`) |
| `gitmap terminal clear` | Terminal Clean | Yes | Yes | No | Yes (unless `-y`) |
| `gitmap clean` | Dev Tools Cache Clean | No | No | Yes | Yes (unless `-y`) |
| `gitmap clean dev` | Dev Tools Cache Clean | No | No | Yes | Yes (unless `-y`) |
| `gitmap devtool clean` | Dev Tools Cache Clean | No | No | Yes | Yes (unless `-y`) |

Screen buffer clearing emits ANSI escape sequences `\033[H\033[2J\033[3J` (reset cursor, clear viewable screen, clear scrollback buffer) accompanied by platform console reset on Windows.

### 2.2 Deploy-Keys Phrasing Matrix

All variants below route seamlessly to `cmdssh.RunSSHDeployKeysCLI`:
- `gitmap deploy all-keys`
- `gitmap deploy all keys`
- `gitmap deploy keys all`
- `gitmap deploy key all`
- `gitmap deploy keys`
- `gitmap deploy-all-keys`
- `gitmap deploy-keys-all`

### 2.3 Bare `gitmap import` Friendly Guidance

When `gitmap import` is invoked without `--confirm`:
- If `len(args) == 0`: Renders a formatted guidance panel explaining the database import requirement and cross-referencing typed envelope commands (`gitmap which-format`, `gitmap import-all-json`, `gitmap sj import`). Returns `nil` without dumping stack trace.
- If target file is provided without `--confirm`: Renders a safety warning indicating `--confirm` is required, with exact copy-paste command. Returns `nil` without dumping stack trace.

### 2.4 Concise Root Unknown Command Handling

In `cmd/rootsuggest.go`, `handleUnknownCommand`:
- Eliminates the call to `printUsage()`.
- Outputs a 2-line box with nearest suggestions:
  `Unknown command 'foo'. Run 'gitmap help' or 'gitmap <command> --help' to view available commands.`

### 2.5 JSON Envelope V2 Schema

```json
{
  "attributes": {
    "type": "ssh-nodes",
    "source": "gitmap-ssh-nodes.json",
    "how": "gitmap ssh export",
    "version": "2.0",
    "gitmapVersion": "6.415.0",
    "importCommand": "gitmap sj import gitmap-ssh-nodes.json -y",
    "exportCommand": "gitmap ssh export-json gitmap-ssh-nodes.json",
    "helpCommand": "gitmap which-format gitmap-ssh-nodes.json",
    "notes": "GitMap SSH Fleet Cluster Nodes and Credentials",
    "timestamp": "2026-09-30T00:00:00Z",
    "workDirectory": "D:\\work",
    "defaultWorkDirectory": "D:\\work",
    "isWorkDirectoryApplied": true,
    "isWorkDirectoryEnforced": false
  },
  "variables": {},
  "data": {
    "schemaVersion": "2.0",
    "exportedAt": "2026-09-30T00:00:00Z",
    "totalNodes": 6,
    "nodes": [
      {
        "workerId": "worker-1",
        "id": 1,
        "alias": "w3",
        "ipAddress": "node-w3",
        "username": "Administrator",
        "port": 22,
        "os": "windows",
        "authMethod": "key",
        "keyPath": "C:\\Users\\Administrator\\.ssh\\id_rsa"
      }
    ]
  }
}
```

### 2.6 Multi-JSON Merge Command (`gitmap merge-json`)

- Invocation: `gitmap merge-json [pattern/dir] [-o <file>] [--type <type>] [-y] [-n]`
- Invocation alias: `gitmap json merge [pattern/dir] [...]`
- Features:
  - Scans files, auto-detects schema format via `jsonenvelope.DetectFormat`.
  - Filters by type if `--type` is specified.
  - Merges items, deduplicating on primary keys (e.g. `alias`/`ipAddress`/`slug`/`name`).
  - Re-indexes numeric IDs starting from 1 (1-based indexing).
  - Emits Envelope V2 with updated `attributes`, `source`, `totalCount`, and `timestamp`.
