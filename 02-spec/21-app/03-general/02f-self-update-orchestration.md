# Self-Update Orchestration (Windows-Safe)

Part of [PowerShell Build & Deploy Patterns](02-powershell-build-deploy.md).

When a CLI updates itself from a PATH-managed executable, use a **three-phase handoff** so the active binary lock is released before deploy AND so cleanup of the locked artifacts can complete after deploy.

## Phase 1: Handoff from active binary

1. `tool update` creates a handoff copy in the same active binary directory (for example `toolname-update-<pid>.exe`, fallback to `%TEMP%` if locked).
2. It launches the handoff copy with a hidden worker command (e.g. `update-runner`).
3. The parent **waits for the worker** using foreground/blocking execution (`cmd.Run()`). This keeps the terminal session stable. The handoff copy is a different file so there is no lock conflict with deploy.

## Phase 2: Execute update from handoff copy

1. Resolve repo root from embedded/configured repo path.
2. Run `run.ps1 -Update` (full pipeline: pull, build, deploy).
3. Sync active PATH binary using **rename-first** strategy in update mode:
   - Rename active binary to `.old` (Windows allows renaming a running exe).
   - Copy deployed binary to the active path.
   - Fall back to copy-retry loop (20 x 500ms) only if rename fails.
4. Read and print versions from the binaries (before update and after update) using `tool version`.
5. Show latest notes using the updated binary (`tool changelog --latest`).

## Phase 3: Cleanup handoff to the freshly deployed binary

After Phase 2 completes, the still-running handoff copy CANNOT clean up its own files. On Windows, both `gitmap-update-<pid>.exe` (held by the live handoff process) and `gitmap.exe.old` (briefly retained by the OS or AV) are locked when invoked from inside the handoff process tree. Running `update-cleanup` from there produces:

```
Error: could not remove cleanup artifact at ...\gitmap-update-<pid>.exe: Access is denied.
Error: could not remove cleanup artifact at ...\gitmap.exe.old: Access is denied.
```

The fix is a **third handoff** to the freshly deployed binary, which is a different file with no shared lock:

1. The handoff copy resolves the path to the deployed binary (`exec.LookPath("gitmap-v28")`, falling back to the active binary's deploy directory).
2. It spawns the deployed binary detached with `update-cleanup`:
   - **Windows**: `cmd.exe /C ping 127.0.0.1 -n 3 >nul & start "" /B "<deployed>" update-cleanup`. The ping delays ~2s so the handoff process can exit and release its file lock; `start "" /B` detaches without opening a window.
   - **Unix**: invoke `<deployed> update-cleanup` directly (no lock conflicts).
3. The handoff process exits immediately after spawning. Its file becomes unlocked, and the deployed binary's cleanup pass removes both the handoff copy and the `.old` backup.

If the handoff copy *is* the deployed binary (Unix in-place update), Phase 3 just calls `runUpdateCleanup` inline — no spawn needed.

## Critical Rules

- The Phase 1 parent MUST use `cmd.Run()` (foreground/blocking). Using `cmd.Start()` + `os.Exit(0)` (async) breaks the terminal session.
- PATH sync MUST use rename-first in update mode. Copy-overwrite fails on Windows when any process holds the binary.
- Generated update scripts MUST NOT contain `Read-Host` or any interactive prompts — they run in non-interactive PowerShell sessions.
- Phase 3 MUST spawn the deployed binary detached on Windows. Running cleanup inline from the handoff copy is the bug we are solving.
- Phase 3 is best-effort. If the spawn fails, the user can re-run `<tool> update-cleanup` manually. Never fail the update because cleanup failed.

## Required Validation

- Fail the update if active version still does not match deployed version after sync.
- Version/changelog output must come from the updated executable, not static constants.
- After Phase 3 completes, neither `<tool>-update-<pid>.exe` nor `<tool>.exe.old` should remain in the deploy directory.

## Cross-References

- Generic spec: [05-handoff-mechanism.md](../08-generic-update/05-handoff-mechanism.md) §Solution: Copy-and-Handoff
- Generic spec: [06-cleanup.md](../08-generic-update/06-cleanup.md) §Phase 3 Cleanup Handoff
- Generic spec: [08-repo-path-sync.md](../08-generic-update/08-repo-path-sync.md) §Post-Deploy Repo Path Sync
- Self-update mechanism: [03-self-update-mechanism.md](03-self-update-mechanism.md)
