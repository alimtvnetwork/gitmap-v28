# Spec 363-01: Command Alias Disambiguation, Multi-Node File Copy (CP) & Move (MV) Engine

> **Specification Sequence:** `363-01`  
> **Canonical Path:** `02-spec/21-app/363-gitmap-command-enhancements/01-alias-and-copy-move-spec.md`  
> **Parent Specification:** `02-spec/21-app/363-gitmap-command-enhancements/00-overview.md`  
> **Parent Plan:** `.ai-memory/plans/pending/363-gitmap-command-enhancements.md`  
> **Subtask Plan:** `.ai-memory/plans/subtasks/363-gitmap-command-enhancements/01-alias-and-copy-move.md`  
> **Author:** GitMap Architecture & Swarm Team  
> **Date:** 31-Oct-26  
> **Status:** `approved`  

---

## 1. Executive Summary & Problem Statement

In standard developer terminal tooling, `cp` is universally recognized as the copy command, and `mv` is universally recognized as the move command. Historically within GitMap, `cp` was overloaded as an alias for `commit-push` (`constants.CmdCommitPushAlias = "cp"` in `cli/constants/constants_commands.go`), while `copy` was routed exclusively to local memory buffer and clipboard operations (`cli/cmdcopy/`). Concurrently, `mv` was routed exclusively to local repository and directory migration via `cmdmerge.RunMove` (`cli/cmdmerge/`).

This legacy layout caused significant user friction:
1. Developers attempting to copy files between machines or locally using `gitmap cp <source> <dest>` inadvertently triggered Git commit and push routines.
2. Developers lacked a unified, intuitive interface for transferring files across fleet nodes over SSH without typing verbose subcommands like `gitmap ssh cp` or `gitmap ssh mv`.
3. Clipboard and memory buffer copying lacked seamless disambiguation from file and multi-node transfers.

This specification establishes a clean architecture:
- **Universal Alias Reassignment:** `cp` is permanently reassigned to `copy` (`gitmap copy` / `gitmap cp`). Legacy `gitmap cp` for commit-push is retired.
- **Semantic Commit Verbs Preservation:** Dedicated semantic commit shortcuts (`cpf` for feature, `cbf` for bug fix, `cpc` for chore, `cpr` for release, and `cfr` for refactor) remain 100% active and protected.
- **Unified Multi-Node Copy Engine:** `gitmap copy` / `gitmap cp` supports 8 intuitive syntax variants covering local-to-node, local-to-fleet, remote-to-local, remote-to-remote, recursive directory copy, and local clipboard fallback.
- **Unified Move Engine & Disambiguation:** `gitmap move` / `gitmap mv` dynamically disambiguates between local repository/directory moves (`cmdmerge.RunMove`) and multi-node network moves with verified remote/local deletion (`cmdssh.RunSSHMvExtended`).

---

## 2. Alias Reassignment & Routing Architecture

### 2.1 Alias Mapping Matrix

| Invocation | Command Meaning | Target Dispatcher Function | Previous Behavior |
| :--- | :--- | :--- | :--- |
| `gitmap copy ...` | Unified Copy Engine | `cmdcopy.RunUnifiedCopy(args)` | Local clipboard buffer only |
| `gitmap cp ...` | Unified Copy Engine (Primary Alias) | `cmdcopy.RunUnifiedCopy(args)` | `cmdcommit.RunCommitPush` (Collision) |
| `gitmap move ...` | Unified Move Engine | `cmdcopy.RunUnifiedMove(args)` | `cmdmerge.RunMove` (Local only) |
| `gitmap mv ...` | Unified Move Engine (Primary Alias) | `cmdcopy.RunUnifiedMove(args)` | `cmdmerge.RunMove` (Local only) |
| `gitmap cpf ...` | Commit & Push Feature | `cmdcommit.RunCommitPushFeature(args)` | Unchanged (Protected) |
| `gitmap cbf ...` | Commit & Push Bug Fix | `cmdcommit.RunCommitPushBugFix(args)` | Unchanged (Protected) |
| `gitmap cpc ...` | Commit & Push Chore | `cmdcommit.RunCommitPushChore(args)` | Unchanged (Protected) |
| `gitmap cpr ...` | Commit & Push Release | `cmdcommit.RunCommitPushRelease(args)` | Unchanged (Protected) |
| `gitmap cfr ...` | Commit & Push Refactor | `cmdcommit.RunCommitPushRefactor(args)` | Unchanged (Protected) |
| `gitmap commit-push ...` | Explicit Commit & Push | `cmdcommit.RunCommitPush(args)` | Unchanged (Protected) |

### 2.2 Root Dispatcher Adjustments

In `cli/cmd/roottooling_util.go`:
- Remove `constants.CmdCommitPushAlias` (`"cp"`) from `constants.CmdCommitPush` registration.
- Retain `constants.CmdCommitPush` (`"commit-push"`, `"cmp"`).

In `cli/cmd/rootutility.go`:
- Register `constants.CmdCopy` (`"copy"`) and `constants.CmdCopyAlias` (`"cp"`) to call `cmdcopy.RunUnifiedCopy(argsTail())`.
- Retain existing memory-specific aliases (`"cp-mem"`, `"copy-to-memory"`) as explicit shortcuts.

In `cli/cmd/rootdata.go`:
- Re-route `constants.CmdMove` (`"move"`) and `constants.CmdMv` (`"mv"`) to call `cmdcopy.RunUnifiedMove(argsTail())`.

---

## 3. Unified Copy Engine (`copy` / `cp`) — 8 Syntax Variants

The unified copy engine parses source and destination descriptors to determine whether the operation is a fleet node transfer, a remote pull, a remote relay, a recursive directory copy, or a local clipboard operation.

### Variant 1: Local to Single Node
- **Syntax:** `gitmap cp <local-src> <node>:<dest>`
- **Behavior:** Reads `<local-src>` from the local filesystem and transfers bytes over SSH to destination path `<dest>` on node `<node>`.
- **Example:** `gitmap cp ./build/app.exe worker-1:/opt/bin/app.exe`

### Variant 2: Local to Single Node with Target Flag
- **Syntax:** `gitmap cp <local-src> <dest> -t <node>` (or `--target <node>`)
- **Behavior:** Explicit target flag specifies remote node. Resolves `<dest>` on target `<node>`.
- **Example:** `gitmap cp ./config.json /etc/gitmap/config.json -t storage-node`

### Variant 3: Local to All Nodes (Fleet Broadcast)
- **Syntax:** `gitmap cp <local-src> all:<dest>` OR `gitmap cp <local-src> <dest> --all`
- **Behavior:** Concurrently broadcasts `<local-src>` to all online fleet nodes registered in `gitmap.db`. Bypasses unreachable nodes gracefully with warning logs.
- **Example:** `gitmap cp ./scripts/healthcheck.sh all:/usr/local/bin/healthcheck.sh`
- **Example:** `gitmap cp ./scripts/healthcheck.sh /usr/local/bin/healthcheck.sh --all`

### Variant 4: Remote Node to Local (Fetch / Pull)
- **Syntax:** `gitmap cp <node>:<remote-src> <local-dest>`
- **Behavior:** Connects to `<node>` over SSH, streams `<remote-src>` bytes back, and writes to local file or folder `<local-dest>`.
- **Example:** `gitmap cp build-server:/var/log/pipeline.log ./logs/remote-pipeline.log`

### Variant 5: Remote Node to Remote Node (Inter-Node Relay)
- **Syntax:** `gitmap cp <node1>:<src> <node2>:<dest>`
- **Behavior:** Streams bytes from `<node1>` through a streaming memory pipeline and writes them directly to `<node2>` over SSH without intermediate disk caching.
- **Example:** `gitmap cp node-alpha:/data/backup.tar.gz node-beta:/backups/backup.tar.gz`

### Variant 6: Remote Node to All Nodes (Broadcast from Remote)
- **Syntax:** `gitmap cp <node1>:<src> all:<dest>`
- **Behavior:** Fetches `<src>` once from `<node1>` into an in-memory buffer, then broadcasts concurrently to all other online nodes in the fleet.
- **Example:** `gitmap cp primary-master:/etc/gitmap/cluster.token all:/etc/gitmap/cluster.token`

### Variant 7: Recursive Directory Copy
- **Syntax:** `gitmap cp -r <local-dir> <node>:<dest-dir>` OR `gitmap cp -r <node>:<remote-dir> <local-dir>`
- **Behavior:** Recursively traverses the directory tree, packs or streams files maintaining relative paths and permissions, and recreates directories on the target.
- **Example:** `gitmap cp -r ./dist/ u1-ubuntu:/home/user/dist/`

### Variant 8: Memory / Clipboard Disambiguation
- **Syntax:** `gitmap copy <text>` OR `gitmap copy --file <path>` OR `cat file | gitmap copy`
- **Disambiguation Rule:** When arguments contain NO node colon syntax (`:` where prefix is not a Windows drive letter `C:`), NO `-t`/`--target` flag, NO `--all` flag, and either:
  1. A single argument that does not point to an existing local destination pairing, OR
  2. Flag `--file` / `--clip`, OR
  3. Piped standard input.
- **Behavior:** Stores text directly into the OS clipboard and persistent GitMap SQLite memory buffer (`db.SaveMemoryEntry`).
- **Example:** `gitmap copy "npm run build && npm test"`
- **Example:** `gitmap copy --file ./notes.txt`

---

## 4. Unified Move Engine (`move` / `mv`) & Disambiguation

The move engine (`gitmap move` / `gitmap mv`) dynamically determines whether the user intended a local repository/directory reorganization or a cross-node fleet move.

```
+--------------------------------------------------------------------+
|                       gitmap mv <src> <dest>                       |
+--------------------------------------------------------------------+
                                  |
                                  v
              +---------------------------------------+
              | Has Node Colon or Fleet Target Flags? |
              +---------------------------------------+
                     /                         \
           No (Local Only)                Yes (Node / Fleet)
                   /                             \
                  v                               v
    +---------------------------+   +-------------------------------+
    | Call `cmdmerge.RunMove`   |   | Call `cmdssh.RunSSHMvExtended`|
    | (Local Repo / Folder Move)|   | (Multi-Node Move with Delete) |
    +---------------------------+   +-------------------------------+
```

### 4.1 Disambiguation Criteria

A command invocation is routed to **Multi-Node Move** (`cmdssh.RunSSHMvExtended`) if ANY of the following positive boolean predicates evaluate to true:
- `hasSrcNode`: Source string contains `<node>:<path>` (excluding Windows drive letter paths like `C:\...`).
- `hasDestNode`: Destination string contains `<node>:<path>` or `all:<path>`.
- `hasTargetFlag`: Flag `-t` or `--target` is present in arguments.
- `hasAllFlag`: Flag `--all` is present in arguments.

If NONE of the above conditions are met (`isLocalOnly == true`), the invocation is routed to **Local Move** (`cmdmerge.RunMove`).

### 4.2 Multi-Node Move Variants & Atomic Cleanup Guarantee

| Variant | Command Syntax | Transport Phase | Post-Verification Cleanup |
| :--- | :--- | :--- | :--- |
| **Local to Single Node** | `gitmap mv <local-src> <node>:<dest>` | SSH stream to `<node>:<dest>` | Delete local file `<local-src>` only after remote hash/byte verification succeeds. |
| **Local to Fleet** | `gitmap mv <local-src> all:<dest>` | Broadcast to all online nodes | Delete local file `<local-src>` only after ALL online target nodes report success. |
| **Remote to Local** | `gitmap mv <node>:<remote-src> <local-dest>` | SSH stream from `<node>` to local | Execute remote SSH `rm` on `<remote-src>` only after local write succeeds. |
| **Remote to Remote** | `gitmap mv <node1>:<src> <node2>:<dest>` | Stream `<node1>` -> memory -> `<node2>` | Execute remote SSH `rm` on `<node1>:<src>` only after `<node2>` write succeeds. |

> [!IMPORTANT]
> The source file MUST NEVER be deleted if the transport phase encounters any network interruption, checksum mismatch, or disk write error. Leaf error sovereignty ensures any transport error aborts cleanup immediately.

---

## 5. Semantic Commit Verbs Protection

GitMap enforces dedicated three-letter semantic commit verbs for conventional development workflows. These verbs remain fully protected and isolated from the `cp` alias reassignment:

| Semantic Shortcut | Canonical Command | Commit Subject Format | Operating Policy |
| :--- | :--- | :--- | :--- |
| `gitmap cpf "<msg>"` | `gitmap commit-push-feature` | `"MODULE - summary"` | 100% Protected |
| `gitmap cbf "<msg>"` | `gitmap commit-push-bugfix` | `"MODULE - summary"` | 100% Protected |
| `gitmap cpc "<msg>"` | `gitmap commit-push-chore` | `"MODULE - summary"` | 100% Protected |
| `gitmap cpr "<msg>"` | `gitmap commit-push-release` | `"v<version> summary"` | 100% Protected (No `VERSION -` prefix) |
| `gitmap cfr "<msg>"` | `gitmap commit-push-refactor` | `"MODULE - summary"` | 100% Protected |

### Disambiguation for `cpc` (Chore Commit vs. Chrome Profile Copy)
- When invoked in a Git repository with a single string argument, `gitmap cpc "chore: message"` routes to chore commit and push.
- When invoked with directory arguments or flags (e.g. `gitmap cpc "Profile 1" "Profile 2"`), sub-routing resolves to Chrome profile operations.

---

## 6. Go Architecture, Types, Interfaces & Signatures

All types and function signatures adhere strictly to GitMap coding guidelines:
- Return type `*apperror.AppError` or `error` for leaf error sovereignty.
- Strongly-typed enums ending in `Type`.
- Positive booleans only (`is*`, `has*`).
- Function length <= 15 lines max.

### 6.1 Package Constants & Enums

Located in `cli/constants/constants_commands.go` and `cli/cmdcopy/constants/constants.go`:

```go
package constants

// Command and alias constants
const (
	CmdCopy      = "copy"
	CmdCopyAlias = "cp"
	CmdMove      = "move"
	CmdMv        = "mv"
)

// Typed Enum for Copy Operating Mode
type CopyModeType string

const (
	CopyModeClipboard     CopyModeType = "clipboard"
	CopyModeLocalToNode   CopyModeType = "local_to_node"
	CopyModeLocalToFleet  CopyModeType = "local_to_fleet"
	CopyModeNodeToLocal   CopyModeType = "node_to_local"
	CopyModeNodeToNode    CopyModeType = "node_to_node"
	CopyModeNodeToFleet   CopyModeType = "node_to_fleet"
)

// Typed Enum for Node Destination Specifier
type NodeTargetType string

const (
	NodeTargetSingle NodeTargetType = "single"
	NodeTargetAll    NodeTargetType = "all"
	NodeTargetLocal  NodeTargetType = "local"
)
```

### 6.2 Descriptor & Parameter Structs

Located in `cli/cmdcopy/types.go`:

```go
package cmdcopy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CopyDescriptor holds parsed source and destination attributes.
type CopyDescriptor struct {
	SrcNode          string
	SrcPath          string
	DestNode         string
	DestPath         string
	HasSrcNode       bool
	HasDestNode      bool
	IsAllNodes       bool
	IsRecursive      bool
	IsClipboardMode  bool
	IsDeleteSource   bool
}

// CopyExecutionOptions encapsulates runtime flags.
type CopyExecutionOptions struct {
	TargetNode       string
	HasTargetNode    bool
	IsAll            bool
	IsRecursive      bool
	IsForceCopy      bool
	IsVerbose        bool
}
```

### 6.3 Core Dispatcher Signatures

```go
// RunUnifiedCopy parses arguments and dispatches to clipboard or multi-node transport.
func RunUnifiedCopy(args []string) error

// RunUnifiedMove parses arguments and dispatches to local move or multi-node move.
func RunUnifiedMove(args []string) error

// ParseCopyDescriptor inspects arguments and extracts node and path targets.
func ParseCopyDescriptor(args []string, isDeleteSource bool) (*CopyDescriptor, error)

// ExecuteNodeTransfer coordinates file bytes transport across SSH nodes.
func ExecuteNodeTransfer(desc *CopyDescriptor) error

// ExecuteRemoteDelete safely removes remote file via SSH session after verified copy.
func ExecuteRemoteDelete(node, remotePath string) error
```

### 6.4 Implementation Exemplar (<= 15 Lines per Function)

```go
// RunUnifiedCopy routes copy requests based on target descriptors.
func RunUnifiedCopy(args []string) error {
	desc, err := ParseCopyDescriptor(args, false)
	if err != nil {
		return err
	}

	if desc.IsClipboardMode {
		return RunCopyCmd(args)
	}

	return ExecuteNodeTransfer(desc)
}

// RunUnifiedMove routes move requests to local movemerge or remote node transport.
func RunUnifiedMove(args []string) error {
	desc, err := ParseCopyDescriptor(args, true)
	if err != nil {
		return err
	}

	if !desc.HasSrcNode && !desc.HasDestNode && !desc.IsAllNodes {
		return cmdmerge.RunMove(args)
	}

	return ExecuteNodeTransfer(desc)
}
```

---

## 7. Error Handling & Leaf Error Sovereignty

All leaf operations return structured errors using `apperror.AppError`:
- Network connection failures return `apperror.NewSSHConnectionError(...)`.
- Non-existent local files return `apperror.NewNotFoundError(...)`.
- Permission denials return `apperror.NewPermissionError(...)`.
- Destination write failures return `apperror.NewIOError(...)`.

No leaf function invokes `os.Exit` or `cliexit.HandleError`. All errors propagate to `cli/cmd/root.go` dispatch handlers for database audit capture in `finishCommandAudit`.

---

## 8. Documentation, Shell Completions & UI Updates

1. **Commands Catalog (`src/data/commands.ts`):**
   - Update `copy` entry: alias `cp`, descriptions for 8 multi-node variants.
   - Update `move` entry: alias `mv`, descriptions for local and multi-node move.
   - Remove `cp` from `cpar` and `commit-push` alias listings.
2. **Help Text Renderers (`cli/cmdcopy/` & `cli/cmdssh/`):**
   - Render multi-node syntax examples in Tokyo Night styled panels.
3. **Shell Completions (`cli/completion/`):**
   - Register `cp` pointing to `copy` completions.
   - Register `mv` pointing to `move` completions.

---

## 9. Acceptance Criteria & Invariants

1. **Alias Integrity:** `gitmap cp` executes file/multi-node copy, NOT commit-push.
2. **Semantic Verbs Protection:** `gitmap cpf`, `gitmap cbf`, `gitmap cpc`, `gitmap cpr`, and `gitmap cfr` execute conventional commits properly.
3. **Local to Single Node:** `gitmap cp src.txt node1:/tmp/src.txt` successfully transfers file over SSH.
4. **Target Flag Syntax:** `gitmap cp src.txt /tmp/src.txt -t node1` resolves and transfers to node1.
5. **Fleet Broadcast:** `gitmap cp src.txt all:/tmp/src.txt` and `gitmap cp src.txt /tmp/src.txt --all` transfer to all online nodes in `gitmap.db`.
6. **Remote to Local:** `gitmap cp node1:/tmp/src.txt ./src.txt` pulls the remote file locally.
7. **Remote to Remote:** `gitmap cp node1:/tmp/data node2:/tmp/data` transfers across nodes without local disk caching.
8. **Recursive Copy:** `gitmap cp -r ./dir node1:/tmp/dir` transfers directories recursively.
9. **Clipboard Fallback:** `gitmap copy "text"` or `gitmap cp "text"` (without destination or node syntax) persists text into memory/clipboard.
10. **Move Disambiguation:** `gitmap mv repo1 repo2` delegates to `cmdmerge.RunMove`. `gitmap mv file.txt node1:/tmp/file.txt` transfers and removes local file only upon verified success.
11. **Positive Booleans Only:** All code and structures use `is*` and `has*` boolean fields.
12. **Relative Paths Only:** All file references and links in documentation use repository-relative paths.
