# Subtask 02: `cli/cloner` Constants & Enums Refactoring

> **Parent Plan:** [.ai-memory/plans/pending/319-package-constants-and-enums-dry-restructuring.md](../../pending/319-package-constants-and-enums-dry-restructuring.md)  
> **Target Package:** `cli/cloner`  
> **Status:** Ready for Execution (Next Turn)  

---

## 1. Scope & Technical Deliverables
1. **Enum: `CloneProtocolType`** (`cli/cloner/enums/clone_protocol_enum.go`):
   - Values: `CloneProtocolInvalid` (0), `CloneProtocolSSH`, `CloneProtocolHTTPS`, `CloneProtocolGit`, `CloneProtocolFile`.
   - Methods: `IsValid() bool`, `String() string`, `IsSSH() bool`, `IsHTTPS() bool`.
2. **Enum: `CloneDepthType`** (`cli/cloner/enums/clone_depth_enum.go`):
   - Values: `CloneDepthInvalid` (0), `CloneDepthFull`, `CloneDepthShallow`, `CloneDepthBlobless`, `CloneDepthTreeless`.
   - Methods: `IsValid() bool`, `String() string`, `IsPartial() bool`, `IsBlobless() bool`.
3. **Constants: `cli/cloner/constants/constants.go`**:
   - Git CLI args: `GitArgClone`, `GitArgDepth`, `GitArgFilterBlob`, `GitArgFilterTree`, `GitArgSingleBranch`.
   - Process environment variables: `EnvGitTerminalPrompt`, `EnvGitSSHCommand`.
   - Progress tokens: `EventCloning`, `EventResolvingDeltas`, `EventComplete`.
4. **Call-Site Refactoring:**
   - Update `cli/cloner/cloner.go`, `cli/cloner/transport.go`, `cli/cloner/options.go` to eliminate inline strings.

---

## 2. Coding Guidelines Invariants
- Files $\le 300$ lines.
- Functions $\le 15$ lines.
- Positive booleans only (`is*`, `has*`).
- Zero nested `if` statements.

---

## 3. Verification Evidence Contract
- Run `gofmt -w cli/cloner/`.
- AST pass via FastGate.
