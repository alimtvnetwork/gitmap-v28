# 01 — Macro Step Execution Failure: `open chrome` (E9000:EXECUTION)

**Status:** PENDING
**Priority:** Normal
**Category:** CLI / Macro Execution Engine
**Reported:** 2026-09-05

---

## 1. Problem Statement & User Session Logs

During macro execution using `gitmap alim` (macro containing 9 steps created via `gitmap macro add "alim"`), steps 1 through 7 executed successfully, but step 8 (`open chrome`) failed with exit status 1:

```
PS M:\Work-files\export\export\chrome-ext> gitmap macro add "alim"

  ● Interactive Macro Builder: "alim"
  Enter commands one per line (empty line or 'done' to save, 'cancel' to abort):
  (Tip: type 'rec' or 'record' to launch live command recording session)

  Step 1> cd %temp%
  Step 2> echo PWD
  Step 3> gitmap cd macro
  Step 4> gitmap pull
  Step 5> cd d:\work
  Step 6> gitmap mkdir -p "sample"
  Step 7> gitmap rm "sample"
  Step 8> open chrome
  Step 9> open "linkedin.com"
  Step 10> done

✓ Created macro "alim" (9 step(s))
   1. cd %temp%
   2. echo PWD
   3. gitmap cd macro
   4. gitmap pull
   5. cd d:\work
   6. gitmap mkdir -p "sample"
   7. gitmap rm "sample"
   8. open chrome
   9. open "linkedin.com"

Run with: gitmap macro run alim (or gitmap alim)

PS M:\Work-files\export\export\chrome-ext> gitmap alim
  alim
  ├── cd %temp%
  ├── echo PWD
  ├── gitmap cd macro
  ├── gitmap pull
  ├── cd d:\work
  ├── gitmap mkdir -p "sample"
  ├── gitmap rm "sample"
  ├── open chrome
  └── open "linkedin.com"

  ▶ Executing Macro: "alim" (9 steps)

  [ 1/9] ➜ cd C:\Users\Alim\AppData\Local\Temp ...   ➜ 📁 Directory: C:\Users\Alim\AppData\Local\Temp
✔ ok (0.0s)
  [ 2/9] ➜ echo PWD ... ✔ ok (0.2s)
  [ 3/9] ➜ gitmap cd macro ... ✔ ok (0.3s)
  [ 4/9] ➜ gitmap pull ... ✔ ok (11.7s)
  [ 5/9] ➜ cd d:\work ...   ➜ 📁 Directory: d:\work
✔ ok (0.0s)
  [ 6/9] ➜ gitmap mkdir -p "sample" ... ✔ ok (0.3s)
  [ 7/9] ➜ gitmap rm "sample" ... ✔ ok (0.2s)
  [ 8/9] ➜ open chrome ... ✖ failed (0.7s)
  ✖ Step 8 failed: exit status 1
gitmap: [E9000:EXECUTION] macro.Execute:
PS M:\Work-files\export\export\chrome-ext>
```

---

## 2. Preliminary Technical Analysis

1. **Root Cause Hypothesis:**
   - On Windows, `open` is not a native command-line binary (unlike macOS `open`).
   - When a macro step specifies `open chrome` or `open "linkedin.com"`, the command dispatcher invokes the system shell (`cmd.exe /c open ...` or `powershell.exe -Command open ...`), which fails with exit status 1 (`'open' is not recognized as an internal or external command`).
   - GitMap already has built-in `gitmap chrome open` and `gitmap open` subcommands, as well as cross-platform browser opening utilities.
2. **Investigation Questions (to be clarified by user):**
   - Should `macro.Execute` provide a built-in cross-platform shim for `open <target>` (mapping to `explorer.exe` / `start` on Windows, `open` on macOS, and `xdg-open` on Linux)?
   - Or should `open chrome` automatically route to `gitmap chrome open` / default browser launcher?
   - Should macro execution support `--continue-on-error` or interactive step recovery?

---

## 3. Pending Implementation Roadmap

1. Audit `macro.Execute` command dispatching logic in `gitmap/cmd/` (check how steps are tokenized and executed).
2. Add cross-platform `open` built-in handler or alias expansion inside the macro runner.
3. Test macro execution on Windows with `open chrome` and `open "url"`.
4. Verify error management and exit code reporting.
