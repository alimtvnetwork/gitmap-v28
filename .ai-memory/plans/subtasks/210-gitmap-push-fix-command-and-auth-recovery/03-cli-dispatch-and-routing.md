# Subtask Plan 03: CLI Command Registration, Dual-Routing Dispatcher & Flag Parsing

> **Task Reference:** `210-gitmap-push-fix-command-and-auth-recovery`  
> **Subtask ID:** `03-cli-dispatch-and-routing`  
> **Status:** Pending Execution  
> **Assigned Agent Role:** Core CLI Dispatch & Infrastructure Engineer  

---

## 1. Overview & Objectives

Implement the CLI command interface and dual-routing entry points for the `gitmap push-fix` workflow. This subtask ensures that users can invoke push recovery either directly via `gitmap push-fix` (or aliases `pf`, `pushfix`, `fix-push`) or via the push subcommand `gitmap push fix` / `gitmap ph fix`.

### Key Deliverables:
1. Define constants `CmdPushFix`, `CmdPushFixAlias`, `CmdPushFixSolid`, and `CmdPushFixInvert` in `cli/constants/constants_cli.go`.
2. Migrate `CmdProfileAlias` from `"pf"` to `"prf"` in `cli/constants/constants_profile.go` to eliminate duplicate identifier collisions and ensure `TestTopLevelCmdConstantsAreUnique` passes.
3. Update `topLevelCmds()` in `cli/constants/cmd_constants_test.go` to maintain 100% AST parity for `TestTopLevelCmdRegistryMatchesAST`.
4. Register top-level dispatch routing in `cli/cmd/rootcore.go` and export bridge helper `runPushFix` in `cli/cmd/clihelpers.go`.
5. Intercept `fix` argument inside `runPush` in `cli/cmdpull/push.go` and route to `RunPushFix`.
6. Implement `PushFixFlags` parser supporting `-n` / `--dry-run`, `-r` / `--remote`, `-b` / `--branch`, `-f` / `--force`, `-y` / `--yes`, `--ssh`, and `--https` in `cli/cmdpull/push_fix_flags.go`.
7. Author comprehensive unit tests for flag parsing and dual-routing dispatch.

---

## 2. Step-by-Step Implementation Steps

### Step 3.1: Define CLI Constants & Resolve Namespace Collision
1. **Edit `cli/constants/constants_cli.go`**:
   - In the const block under comment `// gitmap:cmd top-level`, add:
     ```go
     // CmdPushFix diagnoses failed pushes, repairs authentication, and safely completes git push.
     CmdPushFix       = "push-fix"
     CmdPushFixAlias  = "pf"
     CmdPushFixSolid  = "pushfix"
     CmdPushFixInvert = "fix-push"
     ```
2. **Edit `cli/constants/constants_profile.go`**:
   - Change `CmdProfileAlias = "pf"` to `CmdProfileAlias = "prf"`.
   - Update `constants_profile_test.go` if any profile alias tests reference `"pf"`.
   - Rationale: The short alias `"pf"` is designated for `push-fix` as specified in system requirements. Modifying `CmdProfileAlias` avoids duplicate alias failure in `TestTopLevelCmdConstantsAreUnique`.

### Step 3.2: Register Constants in `topLevelCmds()` Registry
1. **Edit `cli/constants/cmd_constants_test.go`**:
   - Add the new constants to `topLevelCmds()`:
     ```go
     "CmdPushFix":       CmdPushFix,
     "CmdPushFixAlias":  CmdPushFixAlias,
     "CmdPushFixSolid":  CmdPushFixSolid,
     "CmdPushFixInvert": CmdPushFixInvert,
     ```
   - Verify `CmdProfileAlias` remains present in the map (now mapping to `"prf"`).

### Step 3.3: Register Dispatch Handlers in CLI Root
1. **Edit `cli/cmd/rootcore.go`**:
   - In `coreWorkflowEntries()` (or `coreEntries()`), add the dispatch entry:
     ```go
     {[]string{
         constants.CmdPushFix,
         constants.CmdPushFixAlias,
         constants.CmdPushFixSolid,
         constants.CmdPushFixInvert,
     }, func() error { return runPushFix(argsTail()) }},
     ```
2. **Edit `cli/cmd/clihelpers.go`**:
   - Add bridge wrapper function:
     ```go
     func runPushFix(args []string) error {
         return cmdpull.RunPushFix(args)
     }
     ```

### Step 3.4: Intercept `fix` Subcommand in `runPush`
1. **Edit `cli/cmdpull/push.go`**:
   - At the beginning of `runPush(args []string) error`:
     ```go
     if len(args) > 0 {
         first := strings.ToLower(strings.TrimSpace(args[0]))
         if first == "fix" || first == "repair" || first == "push-fix" {
             return runPushFix(args[1:])
         }
     }
     ```
2. **Edit `cli/cmdpull/exports.go`**:
   - Export entry point `RunPushFix`:
     ```go
     // RunPushFix handles the "push-fix" command or "push fix" subcommand.
     func RunPushFix(args []string) error {
         return runPushFix(args)
     }
     ```

### Step 3.5: Implement `PushFixFlags` Parser
1. **Create `cli/cmdpull/push_fix_flags.go`**:
   - Define struct:
     ```go
     package cmdpull

     import (
         "strings"
         "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
     )

     type PushFixFlags struct {
         DryRun         bool
         RemoteName     string
         BranchName     string
         ForceWithLease bool
         AssumeYes      bool
         ConvertSSH     bool
         ConvertHTTPS   bool
         ShowHelp       bool
     }

     func DefaultPushFixFlags() PushFixFlags {
         return PushFixFlags{
             DryRun:         false,
             RemoteName:     "origin",
             BranchName:     "",
             ForceWithLease: false,
             AssumeYes:      false,
             ConvertSSH:     false,
             ConvertHTTPS:   false,
             ShowHelp:       false,
         }
     }

     func ParsePushFixFlags(args []string) (PushFixFlags, error) {
         flags := DefaultPushFixFlags()
         for i := 0; i < len(args); i++ {
             arg := args[i]
             switch {
             case arg == "-h" || arg == "--help" || arg == "help":
                 flags.ShowHelp = true
                 return flags, nil
             case arg == "-n" || arg == "--dry-run" || arg == "-dry-run":
                 flags.DryRun = true
             case arg == "-f" || arg == "--force" || arg == "-force" || arg == "--force-with-lease":
                 flags.ForceWithLease = true
             case arg == "-y" || arg == "--yes" || arg == "-yes":
                 flags.AssumeYes = true
             case arg == "--ssh" || arg == "-ssh" || arg == "--sh":
                 flags.ConvertSSH = true
             case arg == "--https" || arg == "-https" || arg == "--ht":
                 flags.ConvertHTTPS = true
             case arg == "-r" || arg == "--remote":
                 if i+1 < len(args) {
                     i++
                     flags.RemoteName = args[i]
                 }
             case strings.HasPrefix(arg, "--remote="):
                 flags.RemoteName = strings.TrimPrefix(arg, "--remote=")
             case strings.HasPrefix(arg, "-r="):
                 flags.RemoteName = strings.TrimPrefix(arg, "-r=")
             case arg == "-b" || arg == "--branch":
                 if i+1 < len(args) {
                     i++
                     flags.BranchName = args[i]
                 }
             case strings.HasPrefix(arg, "--branch="):
                 flags.BranchName = strings.TrimPrefix(arg, "--branch=")
             case strings.HasPrefix(arg, "-b="):
                 flags.BranchName = strings.TrimPrefix(arg, "-b=")
             }
         }
         return flags, nil
     }
     ```

### Step 3.6: Unit Testing for Flags & Routing
1. **Create `cli/cmdpull/push_fix_flags_test.go`**:
   - Test default flag values.
   - Test short flags parsing (`-n`, `-f`, `-y`, `-r upstream`, `-b feature`).
   - Test long flags parsing (`--dry-run`, `--force`, `--yes`, `--remote=upstream`, `--branch=feature`).
   - Test mutual exclusivity resolution between `--ssh` and `--https`.
   - Test help flag triggering.
2. **Create `cli/cmdpull/push_fix_dispatch_test.go`**:
   - Test interception of `"fix"` in argument list.
   - Test delegation from `RunPushFix` to runner logic.

---

## 3. Verification & Quality Commands

```bash
# Verify AST and Constants parity
go test -v ./cli/constants -run "TestTopLevelCmd"

# Verify Flag parsing unit tests
go test -v ./cli/cmdpull -run "TestPushFixFlags"

# Verify coding guidelines and hygiene
python 03-ai-scripts/05-guideline-autofixer.py cli/constants cli/cmd cli/cmdpull --check-only
python linter-scripts/check-relative-paths.py
```
