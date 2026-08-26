# 07 — Console-Safe Handoff

## Purpose

Define the exact self-update handoff pattern that avoids breaking the
user's console session on Windows while still solving the running-binary
file-lock problem.

This document exists because many implementations try to fix the lock by
launching the updater asynchronously and exiting the parent immediately.
That "fix" often detaches the update from the active console and makes
the terminal output unreliable or invisible.

---

## The Failure Pattern

The broken pattern usually looks like this:

1. `<binary> update` copies itself to `<binary>-update-<pid>.exe`
2. The parent launches the copy asynchronously
3. The parent exits immediately
4. The worker continues in a detached or partially detached session

Common broken variants:

- `cmd.Start()` + `os.Exit(0)`
- PowerShell `Start-Process` without preserving the current console
- `cmd /c start ...`
- Launching the generated update script in the background

### Symptoms

- The prompt returns before the update is actually done
- Live update logs disappear or appear in a separate window
- Stdout/stderr ordering becomes confusing or incomplete
- Interactive script input fails or hangs
- The console session appears broken after update starts

---

## Root Cause

The console problem is different from the file-lock problem.

- **File-lock problem**: Windows does not let a running `.exe` be
  overwritten in place.
- **Console problem**: the parent process that owns the active terminal
  detaches from the worker too early, so the worker no longer runs as a
  normal foreground command in the same session.

The correct fix is **NOT** to detach more aggressively.

The correct fix is:

1. Run the updater from a **different binary file** (handoff copy)
2. Keep the worker in the **same foreground console session**
3. Use **rename-first** during deploy/PATH sync to deal with locks

---

## The Gitmap Fix Pattern

`gitmap` solves this in two blocking layers:

### Layer 1 — Parent → Handoff Worker

In `gitmap/cmd/update.go`:

- `runUpdate()` resolves repo path and executable path
- `createHandoffCopy(selfPath)` creates `gitmap-update-<pid>.exe`
- `launchHandoff(copyPath, repoPath)` starts the worker

The important part is that `launchHandoff()` uses:

```go
cmd := exec.Command(copyPath, copyArgs...)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Stdin = os.Stdin

if err := cmd.Run(); err != nil {
    handleHandoffError(err)
}
```

This is **foreground/blocking** execution. The parent waits.

### Layer 2 — Worker → Update Script

In `gitmap/cmd/updatescript.go`:

```go
cmd := exec.Command("powershell.exe",
    "-ExecutionPolicy", "Bypass",
    "-NoProfile", "-NoLogo", "-File", scriptPath)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Stdin = os.Stdin

err := cmd.Run()
```

So the generated PowerShell script is also run in the same attached,
foreground console session.

---

## Why Blocking Is Correct

This is the key architectural point that AI implementers often miss:

> The parent may still be alive while waiting, but the worker is running
> from a DIFFERENT file.

That means the parent's file lock on the original binary does **not**
conflict with the worker binary itself.

The remaining lock-sensitive step is when the updated binary is synced to
the active PATH location. That is handled by **rename-first**:

1. Rename active binary to `.old`
2. Copy the new binary into place
3. Restore `.old` if copy fails

So the right model is:

- **handoff copy** solves self-replacement
- **foreground blocking** keeps console stable
- **rename-first** solves locked-path sync

Do not trade away console stability just to make the parent exit early.

---

## Mandatory Rules for AI Implementation

1. **The parent MUST use `cmd.Run()`** when launching the handoff worker
2. **The worker MUST also use `cmd.Run()`** when launching the update script
3. **Always pipe `Stdout`, `Stderr`, and `Stdin`** from parent to worker
   and from worker to script
4. **Never use `cmd.Start()` + `os.Exit(0)`** for update handoff
5. **Never use `Start-Process` or `cmd /c start`** for the main update
   worker unless you can prove console attachment, waiting behavior, and
   stream forwarding are fully preserved
6. **Never add `Read-Host` or any interactive prompt** to generated
   update scripts executed via `exec.Command`
7. **Use a different executable file for the worker** (same directory
   first, temp directory fallback)
8. **Use rename-first for Windows sync/deploy**, not copy-overwrite as
   the primary strategy
9. **Propagate worker/script exit codes** back to the original command

---

## Reference Implementation Pattern

```go
func runUpdate() {
    repoPath := resolveRepoPath()

    selfPath, err := os.Executable()
    if err != nil {
        fail(err)
    }

    copyPath := createHandoffCopy(selfPath)

    cmd := exec.Command(copyPath, "update-runner", "--repo-path", repoPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        propagateExit(err)
    }
}
```

```go
func runUpdateScript(scriptPath string) {
    cmd := exec.Command("powershell.exe",
        "-ExecutionPolicy", "Bypass",
        "-NoProfile", "-NoLogo", "-File", scriptPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        propagateExit(err)
    }
}
```

---

## Explicit Anti-Pattern

Do **not** implement this:

```go
cmd := exec.Command(copyPath, "update-runner")
_ = cmd.Start()
os.Exit(0)
```

Why it is wrong:

- The original command returns too early
- The worker is no longer guaranteed to stay attached to the same console
- Output visibility becomes inconsistent
- Users think the update is frozen or broken
- Follow-up steps may run after the shell prompt already returned

---

## Validation Checklist

An implementation is only correct if all of these are true:

1. `tool update` prints progress in the **same console session** from
   start to finish
2. The shell prompt does **not** return until the worker and script both
   finish
3. No extra console window is opened
4. Failures in the script return a non-zero exit code to the caller
5. The generated script does not prompt for input
6. Windows replacement still succeeds because deploy/sync uses
   rename-first
7. Diagnostic output prints both:
   - active binary path
   - handoff copy path

---

## Implementation Checklist for AI

When implementing self-update, apply these steps in order:

1. Add a hidden worker command such as `update-runner`
2. In `update`, create a handoff copy of the current executable
3. Launch the handoff copy with `exec.Command(...)`
4. Attach `Stdout`, `Stderr`, `Stdin`
5. Call `cmd.Run()` and wait
6. In the worker, generate/run the platform update script in foreground
7. Remove interactive prompts from the generated script
8. Use rename-first for Windows deploy and PATH sync
9. Propagate exit codes back to the original command
10. Add validation logs so failures are diagnosable without guessing

If any step uses detached async launch, the implementation is wrong.

---

## Related

- [05-handoff-mechanism.md](05-handoff-mechanism.md) — Copy-and-handoff
  foundation for Windows self-replacement
- [03-rename-first-deploy.md](03-rename-first-deploy.md) — Rename-first
  strategy for locked binaries
- [22-update-command-workflow.md](22-update-command-workflow.md) — End-to-end update runner flow in this folder
- [../12-cicd-pipeline-workflows/06-self-update-mechanism.md](../12-cicd-pipeline-workflows/06-self-update-mechanism.md) — Pipeline duties and self-update contract
- [../17-consolidated-guidelines/17-self-update-app-update.md](../17-consolidated-guidelines/17-self-update-app-update.md) — Consolidated app-level orchestration and handoff guidance

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
