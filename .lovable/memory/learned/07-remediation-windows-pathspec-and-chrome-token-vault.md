# Remediation Windows Pathspec Resolution & Chrome Token Vault Architecture

## 1. Context & Problem Statement
During interactive git remediation (`gitmap fix` / `gitmap reconcile` option `2=wip`), Git commands failed on Windows with:
```text
error: pathspec 'local' did not match any file(s) known to git
error: pathspec 'changes"' did not match any file(s) known to git
✗ Fix failed: exit status 1
```

## 2. Root Cause Analysis
- Go's `exec.Command("cmd", "/c", commandString)` escapes internal quotes as `\"`.
- Windows `cmd.exe` does not unescape `\"` in compound command lines (`&&`), stripping quotes and passing `-m "wip: local changes"` as distinct argv tokens: `-m`, `"wip:`, `local`, `changes"`.
- Git interpreted `"wip:` as the commit message and treated subsequent unquoted tokens `local` and `changes"` as pathspecs.

## 3. Architecture & Resolution
1. **Discrete Argument Isolation (`RemediationStep`)**:
   - Replaced shell string concatenation with structured execution steps:
     ```go
     type RemediationStep struct {
         Name string   `json:"name"`
         Args []string `json:"args"`
     }
     ```
   - Each command step is executed directly with `exec.Command(step.Name, step.Args...)`, preserving exact argv boundaries without shell escaping bugs.
2. **Dynamic Step Synthesis Fallback (`synthesizeRecipeSteps`)**:
   - If a recipe is loaded from legacy state or missing `Steps`, `synthesizeRecipeSteps` dynamically reconstructs structured steps from `repoPath` based on recipe action/title (`Option 1/Stash`, `Option 2/Commit`, `Option 3/Discard`), preventing any fallback to `cmd /c`.
3. **Multi-Target Binary Deployment**:
   - On Windows, PowerShell executes wrappers (`gitmap.ps1`) targeting `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe`.
   - All builds must synchronize:
     - `bin/gitmap.exe`
     - Repository root `gitmap.exe`
     - `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe`
     - `C:\Users\Alim\AppData\Local\gitmap\gitmap.exe`
4. **Chrome Refresh Token Vault & Reversible Ciphers**:
   - Extracted from `Web Data` table `token_service` using lock-free WAL shadow copying.
   - Dual-layer reversible encoding:
     - 2-time Base64 (`doubleBase64`).
     - Alphanumeric Caesar Cipher (`caesarCipher`) with separate modular inverses: `mod 26` for letters, `mod 10` for digits.
     - Raw byte shift (`caesarByteShift`) wrapping `mod 256`.
   - Auto-restoration on `gitmap chrome import`.
