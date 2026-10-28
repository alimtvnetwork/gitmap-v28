# Component and CLI Specification: GitMap Skills & Features Update

**Sequence ID:** 337  
**Document:** `02-component-and-cli-spec.md`  
**Target Packages:** `cli/cmd`, `cli/helpdoc`, `cli/cmdcg`  

---

## 1. CLI Signatures & Flags

### 1.1 PowerShell Runner (`gitmap pwsh` / `gitmap ps`)
```bash
gitmap pwsh [flags] <script-or-args...>
gitmap ps -c "<script>"
```
- Behavior:
  - If `-c` or `-Command` is provided, arguments are safely grouped as a single script execution block.
  - On Windows, script blocks are safely executed without command-line argument splitting splitting `=` operators.

### 1.2 Move/Merge Operations Help
```bash
gitmap help example merge-right
gitmap help example merge-left
```
- Outputs 3-part canonical layout (What It Does, Invocation, Expected Output).

### 1.3 Coding Guidelines & Prompt Install
```bash
gitmap help example cg
gitmap help example install-prompts
gitmap cg --help
```
- Displays complete usage lines with `--all`, `--dry-run`, `--exclude`, `--version`.

---

## 2. Testing and Validation Plan

1. Verify `gitmap pwsh -c '$a = 123; Write-Host $a'` outputs `123` cleanly on Windows.
2. Verify `gitmap help example merge-right` outputs curated examples matching `merge-right.md`.
3. Verify `gitmap help example merge-left` outputs curated examples matching `merge-left.md`.
4. Verify `gitmap help example cg` and `gitmap help example install-prompts` output curated examples.
5. Verify `python .github/scripts/tests/test_ci_scripts.py` passes 100%.
6. Run unit tests in `cli/cmd`, `cli/helpdoc`, `cli/cmdcg`.
