# Subtask 01: Command Alias Disambiguation, Multi-Node File Copy (CP) & Move (MV)

- **Plan:** `364-gitmap-command-enhancements`
- **Subtask ID:** `TASK-01`
- **Slug:** `01-alias-and-copy-move`
- **Status:** `completed`
- **Owner:** Subagent 1 (Alias and Copy-Move Spec Author)
- **Specification:** [01-alias-and-copy-move-spec.md](../../../../02-spec/21-app/364-gitmap-command-enhancements/01-alias-and-copy-move-spec.md)
- **Parent Plan:** [.ai-memory/plans/completed/364-gitmap-command-enhancements.md](../../completed/364-gitmap-command-enhancements.md)

---

## 1. Execution Objective

Reassign CLI alias `cp` exclusively to `gitmap copy` and `mv` exclusively to `gitmap move`, retiring legacy `cp` commit-push aliasing while strictly protecting semantic commit verbs (`cpf`, `cbf`, `cpc`, `cpr`, `cfr`). Implement a unified multi-node file copy engine with 8 syntax variants and an extended move engine with local vs. fleet node disambiguation and atomic remote cleanup.

---

## 2. Detailed Implementation Steps

### Step 1: Command Constants & Alias Cleansing
- In `cli/constants/constants_commands.go`:
  - Deprecate / remove `CmdCommitPushAlias = "cp"`.
  - Add `CmdCopyAlias = "cp"`.
  - Add strongly-typed enum `CopyModeType` ending in `Type` (`CopyModeClipboard`, `CopyModeLocalToNode`, `CopyModeLocalToFleet`, `CopyModeNodeToLocal`, `CopyModeNodeToNode`, `CopyModeNodeToFleet`).
- In `cli/cmdcopy/constants/constants.go`:
  - Define `CmdCopyAlias = "cp"`, `CmdMove = "move"`, `CmdMv = "mv"`.
  - Define Tokyo Night UI color format strings and usage constants.

### Step 2: Root Dispatcher Alignment
- In `cli/cmd/roottooling_util.go`:
  - Remove `constants.CmdCommitPushAlias` (`"cp"`) from `constants.CmdCommitPush` registration block.
  - Verify `constants.CmdCommitPush` (`"commit-push"`, `"cmp"`) remains registered.
- In `cli/cmd/rootutility.go`:
  - Map `constants.CmdCopy` (`"copy"`) and `constants.CmdCopyAlias` (`"cp"`) to call `cmdcopy.RunUnifiedCopy(argsTail())`.
- In `cli/cmd/rootdata.go`:
  - Re-route `constants.CmdMove` (`"move"`) and `constants.CmdMv` (`"mv"`) to call `cmdcopy.RunUnifiedMove(argsTail())`.

### Step 3: Implement Unified Copy Engine (`cmdcopy`) — 8 Syntax Variants
Create `cli/cmdcopy/unified_copy.go` and `cli/cmdcopy/types.go` supporting:
1. **Local to Single Node:** `gitmap cp <local-src> <node>:<dest>`
2. **Local to Single Node with Flag:** `gitmap cp <local-src> <dest> -t <node>` (or `--target <node>`)
3. **Local to All Nodes (Fleet Broadcast):** `gitmap cp <local-src> all:<dest>` OR `gitmap cp <local-src> <dest> --all`
4. **Remote Node to Local (Pull):** `gitmap cp <node>:<remote-src> <local-dest>`
5. **Remote Node to Remote Node (Inter-Node Relay):** `gitmap cp <node1>:<src> <node2>:<dest>`
6. **Remote Node to All Nodes (Broadcast from Remote):** `gitmap cp <node1>:<src> all:<dest>`
7. **Recursive Directory Copy:** `gitmap cp -r <local-dir> <node>:<dest-dir>` / `gitmap cp -r <node>:<remote-dir> <local-dir>`
8. **Memory/Clipboard Disambiguation:** `gitmap copy <text>` or `gitmap copy --file <path>` (when no node colon, single arg / `--file` / piped stdin)

### Step 4: Implement Unified Move Engine (`cmdcopy` & `cmdssh`)
Create `cli/cmdcopy/unified_move.go` and `cli/cmdssh/ssh_mv_extended.go`:
- **Disambiguation Logic:**
  - Inspect args for node colon (`<node>:<path>`), `-t`/`--target`, or `--all`.
  - If purely local paths (`isLocalOnly == true`): delegate to `cmdmerge.RunMove(args)`.
  - If node colon or fleet flags present: delegate to `cmdssh.RunSSHMvExtended(args)`:
    - **Local to Single Node:** Stream to remote node; delete local source only upon verified completion.
    - **Local to Fleet:** Broadcast to all online nodes; delete local source only after all target nodes confirm success.
    - **Remote to Local:** Pull remote file to local; issue remote SSH `rm` command only upon local write success.
    - **Remote to Remote:** Relay bytes from source node to target node; issue remote SSH `rm` on source node only upon target write success.

### Step 5: Semantic Commit Verbs Protection Audit
- Verify `cpf` (Feature), `cbf`/`cpb` (Bugfix), `cpc` (Chore), `cpr` (Release), and `cfr` (Refactor) function as expected.
- In `cli/cmd/cmd_sends.go` and `cli/cmdcommit/commit_push.go`, ensure semantic shortcuts remain intact.
- Add regression tests to ensure typing `gitmap cp` does not stage or commit git changes.

### Step 6: UI, Help Text & Shell Completion Synchronization
- In `src/data/commands.ts`:
  - Update `copy` entry: alias `cp`, document the 8 copy syntax variants.
  - Update `move` entry: alias `mv`, document local repo move vs. multi-node file move.
  - Remove `cp` from `cpar` and `commit-push` alias documentation.
- In `cli/completion/powershell.go`, `cli/completion/allcommands_generated_part1.go`, and `cli/completion/allcommands_generated_part2.go`:
  - Ensure `cp` resolves to `copy` subcommands and flags.
  - Ensure `mv` resolves to `move` subcommands and flags.

---

## 3. Go Interface & Types Contract

```go
package cmdcopy

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CopyDescriptor encapsulates parsed source, destination, and mode flags.
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

// Function signatures adhering to leaf error sovereignty (<= 15 lines each)
func RunUnifiedCopy(args []string) error
func RunUnifiedMove(args []string) error
func ParseCopyDescriptor(args []string, isDeleteSource bool) (*CopyDescriptor, error)
func ExecuteNodeTransfer(desc *CopyDescriptor) error
```

---

## 4. File Modification Ledger

| File Path (Relative) | Action | Scope |
| :--- | :--- | :--- |
| `02-spec/21-app/363-gitmap-command-enhancements/01-alias-and-copy-move-spec.md` | Created | Architectural specification |
| `.ai-memory/plans/subtasks/363-gitmap-command-enhancements/01-alias-and-copy-move.md` | Created | Subtask execution plan |
| `cli/constants/constants_commands.go` | Modify | Reassign `CmdCopyAlias = "cp"`, remove `CmdCommitPushAlias = "cp"` |
| `cli/cmd/roottooling_util.go` | Modify | Detach `cp` from `RunCommitPush` |
| `cli/cmd/rootutility.go` | Modify | Route `copy` and `cp` to `cmdcopy.RunUnifiedCopy` |
| `cli/cmd/rootdata.go` | Modify | Route `move` and `mv` to `cmdcopy.RunUnifiedMove` |
| `cli/cmdcopy/unified_copy.go` | Create | 8-variant unified copy dispatcher and coordinator |
| `cli/cmdcopy/unified_move.go` | Create | Local vs. fleet move disambiguation |
| `cli/cmdssh/ssh_mv_extended.go` | Create | Multi-node move with verified deletion |
| `cli/cmdcopy/types.go` | Create | Typed structs and descriptors with positive booleans |
| `src/data/commands.ts` | Modify | Documentation and UI help updates |

---

## 5. Architectural Guardrails & Invariants

1. **Leaf Error Sovereignty:** No leaf functions may call `cliexit.HandleError` or `os.Exit`. All errors wrap and return `*apperror.AppError` or `error`.
2. **Positive Polarity Only:** Every boolean must begin with `is*` or `has*` (`isRecursive`, `hasSrcNode`, `hasDestNode`, `isAllNodes`, `isClipboardMode`, `isDeleteSource`). Zero negative booleans.
3. **15-Line Function Limit:** No function may exceed 15 lines. Preferred length is <= 8 lines.
4. **Zero Nested Ifs:** Flatten all branching logic using early return guard clauses.
5. **Atomic Move Safety:** The source file or directory must NEVER be deleted if destination transfer fails.
6. **Semantic Commit Preservation:** `cpf`, `cbf`, `cpc`, `cpr`, and `cfr` must continue executing semantic commit workflows without any behavioral regression.
