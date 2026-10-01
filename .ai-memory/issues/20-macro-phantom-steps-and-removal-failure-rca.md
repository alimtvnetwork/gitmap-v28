# RCA-20: Macro Execution Failure on Missing Target Deletion (`rm test`) and Interactive Edit Step Discard Regression

Spec Reference: [02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md](../../02-spec/22-app-issues/53-macro-phantom-steps-and-removal-failure-rca.md)

## 1. Symptom

During interactive macro execution and editing in GitMap:
1. Executing a macro containing `rm test` (`gitmap alim1`) failed with `exit status 1` when the target directory `C:\Users\Alim\test` did not exist:
   ```text
   rm : Cannot find path 'C:\Users\Alim\test' because it does not exist.
   At line:1 char:1
   + rm test
   + ~~~~~~~
       + CategoryInfo          : ObjectNotFound: (C:\Users\Alim\test:String) [Remove-Item], ItemNotFoundException
       + FullyQualifiedErrorId : PathNotFound,Microsoft.PowerShell.Commands.RemoveItemCommand

   ✖ failed (0.2s)
   ✖ Step 1 failed: exit status 1
   ```
2. When the user entered interactive macro editing via `gitmap macro edit alim1` to resolve the issue, they ran `mkdir -p test` and `ls`:
   ```text
   [PWD: C:\Users\Alim]
   Edit [2]> mkdir -p test
   ✓ Created directory: C:\Users\Alim\test
   (created directory live. Enter command for Step, or '+add' to record 'mkdir -p test')

   [PWD: C:\Users\Alim]
   Edit [2]> ls
   ...
   (inspected directory. Enter command for Step, or '+add' to record 'ls')

   [PWD: C:\Users\Alim]
   Edit [2]> exit
   ✔ Macro "alim1" successfully updated (1 steps)
   ```
3. Upon exiting, the user discovered that neither `mkdir -p test` nor `ls` had been saved as macro steps. The macro retained only the original failing step `rm test`.

---

## 2. Root Cause

Windows PowerShell's `rm` (`Remove-Item`) fails terminatingly on non-existent targets rather than being idempotent, the installed GitMap binary at `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` was stale (built before safe removal and edit recording shims), and interactive edit mode treated workspace-modifying commands as transient unrecorded helpers.

1. **Platform Non-Idempotent `rm` Behavior on Windows:**
   On Windows PowerShell, `rm` is an alias for the `Remove-Item` cmdlet. Unlike POSIX shells (`rm -f`), PowerShell throws a terminating `ItemNotFoundException` if the specified target does not exist. GitMap's macro step runner directly piped `rm <target>` into the shell without existence checking or suppression flags, causing immediate step failure.
2. **In-Builder Helper Hijack in Macro Editor:**
   In `cli/cmdmacro/macro_add_helpers.go`, commands starting with `mkdir` and `ls` were categorized as in-builder inspection helpers (`isExactHelper` and `hasBasicPrefix`). In `macro_edit.go`, `handleEditSpecialAction` calls `processInBuilderCommand` before `recordStepLine`. Consequently:
   - When the user types `mkdir -p test`, the helper intercepted it, executed the directory creation live on the host filesystem, and advised: `(created directory live. Enter command for Step, or '+add' to record 'mkdir -p test')`.
   - If the user did not explicitly type `+add`, the command was discarded from the recorded macro step list.
   - When the user exited the editor, only the pre-existing steps were written back to the SQLite store.

---

## 3. Resolution

1. **Idempotent Removal Shim in Macro Step Execution:**
   - In `cli/macro/safe_rm.go` and `cli/macro/execute.go`, when running on Windows PowerShell and encountering steps beginning with `rm `, `rmdir `, or `Remove-Item `, transform the command into an idempotent script:
     ```powershell
     foreach ($__target in @(...)) { if (Test-Path -LiteralPath $__target) { Remove-Item -Recurse -Force -LiteralPath $__target } }
     ```
2. **Interactive Macro Edit Recording Polish:**
   - In `cli/cmdmacro/macro_edit.go`, allow explicit helpers only when prefixed with `:` (e.g. `:ls`, `:cat`), `cd`, `pwd`, or `help`. All standard workspace commands (`mkdir -p test`, `rm`, etc.) append directly as steps with live execution.
3. **Audit Logging Integration:**
   - Record macro mutations into `TaskHistory` via `cmdtask.RecordTaskAudit`.

---

## 4. Prevention & Learnings

- **Cross-Platform Shell Idempotence:** Deletion commands (`rm`, `del`) in automation scripts and macros must always be idempotent across Linux, macOS, and Windows.
- **Explicit Editor UX:** Never silently discard user commands in an interactive editor where the prompt promises to append typed commands.
- **Task & Audit Mode Integration:** All macro mutations (create, edit, run, delete) must be recorded into the persistent task audit table to preserve execution history and reproducibility.
