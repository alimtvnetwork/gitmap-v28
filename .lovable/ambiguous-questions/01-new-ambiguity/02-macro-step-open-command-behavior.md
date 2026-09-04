# Macro Step Execution Behavior for `open <target>`

Slug: macro-step-open-command-behavior
Status: open
Raised: 2026-09-05
Blocking: .lovable/plans/pending/01-macro-step-open-chrome-failure.md

## Question

During interactive macro execution (`gitmap alim`, where step 8 was `open chrome` and step 9 was `open "linkedin.com"`), step 8 failed on Windows with `exit status 1` (`[E9000:EXECUTION] macro.Execute:`) because `open` is not a native command-line binary on Windows.
The user instructed: *"Write this as pending task I will explain later"*.

What is the expected resolution behavior for `open <target>` within the macro runner?

1. Should `macro.Execute` provide an internal cross-platform shim for `open` (mapping to `explorer.exe` / `start` on Windows, `/usr/bin/open` on macOS, and `xdg-open` on Linux)?
2. Should `open chrome` specifically delegate to `gitmap chrome open` or the registered browser launcher?
3. Should macros support error tolerances (e.g. `--continue-on-error` or interactive retry)?

## Options considered

1. **Option A: Internal macro shim for `open`**: Detect `open` as the leading command word. On Windows, translate `open <target>` into `rundll32.exe url.dll,FileProtocolHandler <target>` or `cmd.exe /c start "" <target>`. On Linux, translate to `xdg-open <target>`. On macOS, preserve `open <target>`.
2. **Option B: Delegate to `gitmap open`**: Expand `open <target>` to invoke GitMap's built-in `gitmap open` command.
3. **Option C: Await user explanation**: Maintain plan `01-macro-step-open-chrome-failure.md` as pending until the user clarifies the desired architectural design.

## Impact if guessed wrong

Implementing an incorrect opening strategy could break macros expecting custom shell commands or cause unexpected detached process behavior on Windows.
