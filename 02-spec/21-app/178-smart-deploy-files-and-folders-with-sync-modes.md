# Spec 178: Smart File and Folder Deployment with Bidirectional Sync Modes

## 1. Metadata
- **Spec ID:** 178
- **Title:** Smart File and Folder Deployment with Bidirectional Sync Modes
- **Category:** Core Tooling / Cluster SSH Automation
- **Status:** Active
- **Created:** 2026-09-27
- **Target CLI Version:** v6.357.0

---

## 2. User Request (Verbatim)

```text
these commands should work

 gitmap deploy <alias, ip, seq, id> <from or current releative file or folder> <abs path or path from default work dir into that machine> # clear???

 gitmap deploy <alias, ip, seq, id> <from or current releative file or folder> <abs path or path from default work dir into that machine> --json # clear???, will show summary response in json mode, if same file match or found then will ask for permission if we want to replace to skip

gitmap deploy <alias, ip, seq, id> <from or current releative file or folder> <abs path or path from default work dir into that machine> --overrite(o)/skip/sync/sync-right/sync-left  # don't seek for prompt, copy and replace , sync means it will only take if left one is latest by date or copy from right to left if the right has the latest date, for folder can do parallel replace ?? can???, sync right means don't pass any file to left, sync left means keep the left intact right on conflicts 

also create prompts for deploy-right, deply left
```

---

## 3. Architectural Overview

### 3.1 Command Syntax & Dispatch
GitMap introduces `gitmap deploy` (and `gitmap ssh deploy`) as an intelligent, unified file and folder synchronization engine across cluster nodes:

```bash
gitmap deploy <target> <source> <destination> [flags]
```

Where:
- `<target>`: Identifies the remote host using `<alias, ip, seq, id>`:
  - Alias: `w1`, `w2`, `alpha-win`
  - IP Address: `192.168.1.3`, `10.20.0.11`
  - Sequence Index: `1`, `4` (1-indexed matching `gitmap ssh ls`)
  - Host ID: Database ID from `ssh_hosts` table (e.g. `host-1`)
- `<source>`: Local file or folder path (relative to current working directory or absolute).
- `<destination>`: Remote path. If relative, resolved against the remote user's default working directory (`D:/work` or `~`). If absolute, written directly to that path.

### 3.2 Conflict Resolution Modes
When target files already exist on the remote host:
1. **Interactive Prompt (Default in interactive TTY):** When neither `--overwrite`, `--skip`, nor `--sync*` is specified and a conflict is detected, prompts the user:
   `File '<dest>' exists on remote node. Overwrite? [y/N/a(ll)/s(kip all)]`
2. **`--overwrite` / `-o`:** Unconditionally overwrites remote files with local files without prompt.
3. **`--skip`:** Skips any file that already exists on the remote node.
4. **`--sync`:** Full bidirectional synchronization based on modification timestamp (mtime):
   - If local file is newer than remote file: copy local -> remote.
   - If remote file is newer than local file: copy remote -> local.
   - If timestamps/sizes match: skip.
5. **`--sync-right`:** Unidirectional rightward sync (Local -> Remote):
   - Only copies files where local is newer than remote.
   - Never copies files from remote back to local.
   - Prevents stale local copies from overwriting newer remote work.
6. **`--sync-left`:** Unidirectional leftward sync (Remote -> Local) / Local Protection:
   - Keeps local files completely intact on conflicts.
   - Pulls newer remote files down to local workspace without overwriting local changes.

### 3.3 Folder Deployment & Parallel Replacement
When `<source>` is a directory:
- Recursively inventories all files in the source hierarchy.
- Probes remote directory structure and file metadata over SSH.
- Executes parallel transfers using a bounded worker pool (up to 8 concurrent streams).
- Creates remote directories on the fly.

### 3.4 JSON Telemetry
When `--json` is supplied, suppresses interactive terminal progress bars and renders a machine-readable JSON summary:
```json
{
  "target": "w1",
  "ip": "192.168.1.3",
  "source": "dist/",
  "destination": "D:/work/dist",
  "mode": "sync-right",
  "totalFiles": 142,
  "transferredFiles": 18,
  "skippedFiles": 124,
  "transferredBytes": 4821042,
  "conflicts": 0,
  "durationMs": 1420,
  "success": true
}
```

---

## 4. Acceptance Criteria
- [ ] `gitmap deploy <target> <src> <dest>` routes directly from CLI root and `gitmap ssh deploy`.
- [ ] Target resolution supports `<alias, ip, seq, id>` natively.
- [ ] Relative source and destination paths resolve accurately against local CWD and remote default workdir.
- [ ] Single file deployment succeeds with streaming transport.
- [ ] Directory folder deployment recurses subfolders and supports parallel worker replacement.
- [ ] Conflict flags (`--overwrite` / `-o`, `--skip`, `--sync`, `--sync-right`, `--sync-left`) execute deterministically.
- [ ] `--json` flag formats complete structured output and exit codes.
- [ ] Canonical prompts created in `01-prompts/deploy-right.md` and `01-prompts/deploy-left.md`.
- [ ] Documentation and help text updated in `cli/helptext/deploy.md` and terminal help renderer.
