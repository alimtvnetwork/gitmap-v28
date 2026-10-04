# Subtask 04: Intelligent Work Directory & Relative Path Resolution

> **Parent Plan:** [67-nodes-cfr-remote-fleet-clone-enhancement](../../pending/67-nodes-cfr-remote-fleet-clone-enhancement.md)
> **Component Spec:** [02-component-spec.md](../../../../02-spec/21-app/67-nodes-cfr-remote-fleet-clone-enhancement/02-component-spec.md)
> **Primary Targets:**
> - `cli/cmdnodes/nodes_clone_file.go`
> - `cli/cmdnodes/nodes_clone.go`
> - `cli/cmdnodes/nodes_clone_types.go`
> - `cli/cmdnodes/nodes_clone_remote.go`
> - `cli/cmdnodes/nodes_clone_test.go`

---

## 1. Objective

Enable intelligent current working directory (CWD) relative path preservation and custom destination path overrides for `gitmap nodes cfr`, `nodes clone`, and `nodes cfrp`:
1. If the user invokes the command from a subfolder inside their local work directory (e.g. `./presentations-repos`), automatically mirror and clone into `<remoteWorkDir>/presentations-repos` on all remote fleet nodes (`./presentations-repos` on Windows, `~/work/presentations-repos` on Linux).
2. If the user invokes the command from outside their work directory (e.g. `C:\Users\Admin` or `/tmp`), default to the remote work directory roots (`D:\work` / `~/work`).
3. Support explicit custom destination paths via flags (`-d`, `--dir`, `--dest`, `--target-dir`) or an optional positional destination argument (`gitmap nodes cfr <target> [dest]`).
4. Ensure target directories are created safely on remote machines (`mkdir -p` on Linux, `New-Item -ItemType Directory -Force` on Windows) prior to executing `gitmap <kind>`.

---

## 2. Detailed Technical Scope

### 2.1 Struct Enhancements (`cli/cmdnodes/nodes_clone_types.go`)
Extend `NodesCloneOptions` to store relative path telemetry and custom destination state:
- `TargetDir string`: Explicit or resolved destination path.
- `RelativeSubDir string`: Clean relative subfolder path from work directory root (e.g. `presentations-repos` or `clients/subproject`).
- `IsInsideWorkDir bool`: True when local CWD is located within the configured work directory.
- `HasCustomDest bool`: True when an explicit destination flag or positional argument was supplied.

### 2.2 Work Directory Resolution Helpers (`cli/cmdnodes/nodes_clone_file.go`)
Implement concise, modular functions (each $\le 15$ lines):
- `ResolveLocalWorkDir() string`:
  - Query default work directory from `store.OpenDefault()`.
  - Fall back to `D:\work` on Windows or `~/work` on Unix.
- `ResolveCwdRelativeSubDir(workDir string) (string, bool)`:
  - Fetch `os.Getwd()`.
  - Check `fsutil.IsInsideWorkDir(cleanCwd, cleanWork)`.
  - Compute `filepath.Rel(cleanWork, cleanCwd)`.
  - Return relative subpath normalized with forward slashes (`filepath.ToSlash`), or empty string if at work directory root.
- `ResolveRemoteTargetDir(osType, relativeSubDir, customDest string) string`:
  - If `customDest != ""`, format and return `customDest`.
  - If `relativeSubDir != ""`, append to base remote work root (`./<rel>` on Windows, `~/work/<rel>` on Linux).
  - Otherwise, return base remote work root (`D:\work` / `~/work`).

### 2.3 Argument & Flag Parsing (`cli/cmdnodes/nodes_clone.go`)
Enhance `parseNodesCloneOptions(kind NodesCloneKind, raw []string)`:
- Support destination flags:
  - `-d <dir>`, `--dir <dir>`, `--dest <dir>`, `--target-dir <dir>`.
  - Inline forms: `--dir=<dir>`, `--dest=<dir>`, `--target-dir=<dir>`.
- Extract optional positional destination argument:
  - In `extractTargetDirFromPassArgs`, correctly detect if the second positional parameter is a destination path rather than a URL or second target.
- Automatically populate `opts.IsInsideWorkDir`, `opts.RelativeSubDir`, and `opts.TargetDir` if no explicit destination was provided.

### 2.4 Remote Pre-Creation Commands (`cli/cmdnodes/nodes_clone_remote.go`)
Update remote shell command generation so target directories are guaranteed to exist prior to navigation:
- **Windows PowerShell:**
  ```powershell
  if (!(Test-Path -Path "<targetDir>")) { New-Item -ItemType Directory -Path "<targetDir>" -Force | Out-Null }; Set-Location "<targetDir>"; gitmap <kind> <args>
  ```
- **Linux/Unix Bash:**
  ```bash
  mkdir -p <targetDir> && cd <targetDir> && gitmap <kind> <args>
  ```

---

## 3. Implementation Steps

1. **Step 1:** Update `NodesCloneOptions` in `cli/cmdnodes/nodes_clone_types.go` with `RelativeSubDir`, `IsInsideWorkDir`, and `HasCustomDest`.
2. **Step 2:** Author `ResolveLocalWorkDir`, `ResolveCwdRelativeSubDir`, and `ResolveRemoteTargetDir` in `cli/cmdnodes/nodes_clone_file.go`.
3. **Step 3:** Update `parseNodesCloneOptions` and helper functions in `cli/cmdnodes/nodes_clone.go` to parse destination flags and compute work directory relative paths.
4. **Step 4:** Modify `buildWindowsWorkDirExecString` and `buildUnixWorkDirExecString` in `cli/cmdnodes/nodes_clone_remote.go` to inject pre-creation commands.
5. **Step 5:** Add comprehensive unit tests in `cli/cmdnodes/nodes_clone_test.go` verifying:
   - Subfolder inside work directory -> remote path has subfolder appended.
   - CWD outside work directory -> remote path defaults to `D:\work` / `~/work`.
   - Explicit flag override (`--dest`) -> takes absolute precedence.
   - Positional destination override -> takes precedence.

---

## 4. Verification & Quality Gates

- Run unit tests:
  ```powershell
  go test -v ./cli/cmdnodes/ -run TestWorkDirAndRelativePath
  ```
- Verify zero compiler errors across `cli/`:
  ```powershell
  go build ./cli/...
  ```
- Verify compliance with coding guidelines (all functions $\le 15$ lines, positive booleans, no raw git commands).
