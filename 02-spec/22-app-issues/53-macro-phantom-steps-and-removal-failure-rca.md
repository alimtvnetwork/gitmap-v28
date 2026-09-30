# RCA-53: Macro Execution Failure on Missing Target Deletion (`rm test`) and Interactive Edit Step Discard Regression

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
4. The user was surprised that `rm test` was originally recorded as step 1 while new commands typed in edit mode were ignored.

---

## 2. Root Cause

1. **Platform Non-Idempotent `rm` Behavior on Windows:**
   On Windows PowerShell, `rm` is an alias for the `Remove-Item` cmdlet. Unlike POSIX shells (`rm -f`), PowerShell throws a terminating `ItemNotFoundException` if the specified target does not exist. GitMap's macro step runner directly piped `rm <target>` into the shell without existence checking or suppression flags, causing immediate step failure.
2. **In-Builder Helper Hijack in Macro Editor:**
   In `cli/cmdmacro/macro_add_helpers.go`, commands starting with `mkdir` and `ls` are categorized as in-builder inspection helpers (`isExactHelper` and `hasBasicPrefix`). In `macro_edit.go`, `handleEditSpecialAction` calls `processInBuilderCommand` before `recordStepLine`. Consequently:
   - When the user types `mkdir -p test`, the helper intercepts it, executes the directory creation live on the host filesystem, and advises: `(created directory live. Enter command for Step, or '+add' to record 'mkdir -p test')`.
   - If the user does not explicitly type `+add`, the command is discarded from the recorded macro step list.
   - When the user exits the editor, only the pre-existing steps are written back to the SQLite store.
3. **Inconsistent UX Contract:**
   The macro editor header displays:
   `Type any shell command to append it (executes live in terminal):`
   This creates a direct contradiction: the prompt states that typing any shell command will append it, yet built-in commands like `mkdir` and `ls` are treated as non-recorded helpers unless followed by `+add`.

---

## 3. Resolution

1. **Idempotent Removal Shim in Macro Step Execution:**
   - In `cli/macro/execute.go`, when running on Windows PowerShell and encountering steps beginning with `rm `, `rmdir `, or `Remove-Item `, transform the command into an idempotent script:
     ```powershell
     if (Test-Path '<target>') { Remove-Item -Recurse -Force '<target>' }
     ```
   - Provide a native `gitmap rm` / `safe-rm` command in `cli/cmd/` that safely removes files/directories and returns exit 0 if they already do not exist.
2. **Interactive Macro Edit Recording Polish:**
   - In `cli/cmdmacro/macro_edit.go`, distinguish between navigation helpers (e.g. `cd`, `pwd`, `:ls`) and file modification commands (`mkdir`, `touch`).
   - In edit mode, when the user executes a command that modifies the workspace (like `mkdir -p <dir>`), automatically append it as a step unless the command is prefixed with a colon `:` indicating a transient helper.
3. **Macro Terminal Visual Polish:**
   - Render structured box borders, colored status badges (`[✔] ok`, `[✖] failed`), runtime timings, and clean error envelopes to replace unformatted error blocks.

---

## 4. Prevention & Learnings

- **Cross-Platform Shell Idempotence:** Deletion commands (`rm`, `del`) in automation scripts and macros must always be idempotent across Linux, macOS, and Windows.
- **Explicit Editor UX:** Never silently discard user commands in an interactive editor where the prompt promises to append typed commands.
- **Task & Audit Mode Integration:** All macro mutations (create, edit, run, delete) must be recorded into the persistent task audit table to preserve execution history and reproducibility.
