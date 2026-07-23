# Self-Update Gold Standard (Generic CLI)

> Canonical, framework-agnostic blueprint for implementing safe, terminal-stable
> self-update in any CLI. Distilled from the gitmap-v27 reference implementation.

This spec is the **single instruction set** for AI assistants and engineers
adding `tool update` to a new CLI. If your CLI updates itself from a
PATH-managed executable, you MUST follow every rule below.

Cross-references (deep dives):

- [`spec/08-generic-update/01-self-update-overview.md`](../08-generic-update/01-self-update-overview.md)
- [`spec/08-generic-update/02-deploy-path-resolution.md`](../08-generic-update/02-deploy-path-resolution.md)
- [`spec/08-generic-update/03-rename-first-deploy.md`](../08-generic-update/03-rename-first-deploy.md)
- [`spec/08-generic-update/05-handoff-mechanism.md`](../08-generic-update/05-handoff-mechanism.md)
- [`spec/08-generic-update/06-cleanup.md`](../08-generic-update/06-cleanup.md)
- [`spec/08-generic-update/07-console-safe-handoff.md`](../08-generic-update/07-console-safe-handoff.md)
- [`spec/08-generic-update/08-repo-path-sync.md`](../08-generic-update/08-repo-path-sync.md)
- [`spec/03-general/02f-self-update-orchestration.md`](../03-general/02f-self-update-orchestration.md)

---

## 1. Why a dedicated update flow exists

A running executable is **locked by the OS** on Windows. You cannot
`Copy-Item -Force` over it. You cannot rebuild straight to its install path.
Naive solutions either:

- silently leave the user on the old binary (no actual upgrade), or
- corrupt PATH state when partial copies fail, or
- detach the terminal and lose the user's session.

The gold standard solves all three with a **two-phase handoff** + **rename-first
deploy** + **mandatory post-update verification** + **artifact cleanup**.

---

## 2. The Two-Phase Handoff (MANDATORY)

### Phase 1 — Handoff from active binary

1. `tool update` resolves repo path (flag → embedded → DB → prompt → updater fallback).
2. Copies the **currently running binary** to a sibling file in the same
   directory: `<tool>-update-<pid>.exe` (Windows) or `<tool>-update-<pid>` (Unix).
3. If that directory is not writable, falls back to `%TEMP%` / `$TMPDIR`.
4. Launches the handoff copy with a hidden subcommand
   (e.g. `update-runner`) using **`cmd.Run()` (foreground/blocking)**.
5. Parent waits. Terminal stays attached. Stdout/Stderr/Stdin are inherited.

### Phase 2 — Update pipeline from handoff copy

The handoff copy does the real work because it is a **different file** —
the install path is no longer locked.

1. Re-resolve repo path.
2. Run the build/deploy pipeline (`run.ps1 -Update` on Windows,
   `run.sh --update` on Unix).
3. Sync the active PATH binary using **rename-first** (see §3).
4. Print version BEFORE and AFTER from the actual executables (not constants).
5. Run `tool changelog --latest` from the **updated** binary.
6. Run `tool update-cleanup` to remove handoff and `.old` artifacts.

### Hard prohibitions

- ❌ **Never** use `cmd.Start()` + `os.Exit(0)` to detach. It kills the terminal session.
- ❌ **Never** copy-overwrite the active PATH binary directly — it is locked.
- ❌ **Never** generate update scripts that contain `Read-Host` or any prompt.
  Scripts run in non-interactive PowerShell.
- ❌ **Never** print the post-update version from a baked-in constant. Always
  invoke the deployed binary and capture its output.

---

## 3. Rename-First Deploy (MANDATORY in update mode)

Windows allows you to **rename a running executable** but not overwrite it.
Exploit that:

```
1. Rename  <deployPath>/tool.exe → <deployPath>/tool.exe.old
2. Copy    <buildOutput>/tool.exe → <deployPath>/tool.exe
3. On failure: rename .old back, abort, exit non-zero
4. Copy-retry loop (20 × 500 ms) is fallback ONLY if rename fails
```

This is the difference between "the user still runs the old version after
update" and "the upgrade actually took effect this session."

In **non-update mode** (fresh build/deploy where the binary isn't locked),
plain copy-with-backup is acceptable.

---

## 4. Repo Path Resolution Order

Resolve in this exact priority and stop at the first hit:

1. `--repo-path <path>` CLI flag
2. Embedded constant (set at build time via ldflags)
3. SQLite/JSON state DB previously saved
4. Interactive prompt (only if a TTY is attached)
5. Fallback to `<tool>-updater` binary (release-based update) if present
6. Fail with actionable error listing all 4 install options

Persist any newly resolved path back to the state DB so future runs are silent.

---

## 5. Post-Update Validation (acceptance criteria)

The update **MUST FAIL LOUDLY** if any of these don't hold:

- `which tool` resolves to a binary that reports the new version
- Active PATH binary version == deployed binary version
- Changelog output is produced by the **updated** binary
- Cleanup ran (no leftover `<tool>-update-<pid>.exe` or `.old` files)

Print a unified before → after summary:

```
Current version: v2.49.0
Latest version:  v2.49.1
v2.49.0 → v2.49.1
✓ Updated to v2.49.1
```

---

## 6. Cleanup Command (`tool update-cleanup`)

Required subcommand. Idempotent. Removes:

- `<installDir>/<tool>-update-*.exe` and `<installDir>/<tool>-update-*` (Unix)
- `<installDir>/<tool>.exe.old` and `<installDir>/<tool>.old`
- `%TEMP%/<tool>-update-*.*` fallback artifacts
- Any `<tool>.exe.bak` left by failed deploys

Runs automatically at the end of a successful update. Must also be safe to
invoke manually at any time. Print "Nothing to clean up" when empty.

---

## 7. Cross-Platform Implementation Notes

| Concern | Windows | Unix |
|---|---|---|
| Handoff copy name | `<tool>-update-<pid>.exe` | `<tool>-update-<pid>` |
| Executable bit | n/a | `os.Chmod(0o755)` after copy |
| Pipeline runner | `pwsh -ExecutionPolicy Bypass -File run.ps1 -Update` | `bash run.sh --update` |
| Sync strategy | rename-first, copy fallback | direct copy is safe (no exclusive lock) |
| Script encoding | UTF-8 **with BOM** | UTF-8 plain |
| Profile reload after install | dot-source `$PROFILE` in-process | print `source ~/.zshrc` one-liner |

---

## 8. Implementation Checklist (copy into PR description)

- [ ] `tool update` enters via `runUpdate()` — only handles handoff, never deploy
- [ ] `tool update-runner` is a **hidden** subcommand (not in `--help`)
- [ ] `tool update-cleanup` is a public subcommand and runs at end of update
- [ ] Handoff copy created in install dir, falls back to TEMP
- [ ] Parent uses `cmd.Run()` (foreground). Grep proves no `cmd.Start()` in update path
- [ ] Generated PS1 written with UTF-8 BOM, contains zero `Read-Host` calls
- [ ] Rename-first deploy gated behind `-Update` switch in `run.ps1`
- [ ] Repo path resolved via 6-step fallback chain, persisted to DB
- [ ] Updater fallback (`<tool>-updater`) attempted before failing with no repo
- [ ] Post-update version comparison printed from actual binaries
- [ ] `changelog --latest` invoked from the updated binary
- [ ] Help file `helptext/update.md` documents all 6 example scenarios
- [ ] Integration test: full update on Windows + Unix, no terminal disconnect
- [ ] Lint check / grep guard: no `cmd.Start(` inside `cmd/update*.go`

---

## 9. Reference Implementation Files (gitmap-v27)

| Concern | File |
|---|---|
| Update entrypoint + handoff | `gitmap-v27/cmd/update.go` |
| Script generation + execution | `gitmap-v27/cmd/updatescript.go` |
| PowerShell pipeline (rename-first) | `run.ps1` (search for `-Update` branch) |
| Cleanup | `gitmap-v27/cmd/updatecleanup.go` |
| Helptext example | `gitmap-v27/helptext/update.md` |

Use these as the literal template when bootstrapping a new CLI.
