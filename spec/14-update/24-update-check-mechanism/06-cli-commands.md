# CLI Commands — `update-check` & `do-update`

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## 1. `update-check`

### Synopsis

```
<cli> update-check [--async] [--force] [--json]
```

### Flags

| Flag | Default | Behavior |
|------|---------|----------|
| `--async` | off | Spawns a detached child process, returns immediately. The child runs `update-check --force` (no flag) and writes silently. |
| `--force` | off | Bypasses the `IsCheckDue` interval gate. Always runs the discovery. |
| `--json`  | off | Prints the combined JSON document instead of the formatted human output. |

> **No `--persist` flag.** Both sync and `--async` modes always persist
> to the store. (Resolved Decision #7.)

### Default (sync) behavior

1. Call `service.Initialize()`.
2. Call `service.RunAndPersist(ctx)` — this fetches, builds the
   combined JSON, and writes to the store inside try/catch.
3. Print the formatted output:

```
Current version: V1.5.0
Latest version:  V1.7.0
Status:          UpdateFound
Published at:    2026-04-15 10:00:00 UTC
Install (this OS):
  iwr -useb https://github.com/AlifKarim/repo-v17/releases/download/v1.7.0/Install.ps1 | iex

Run `<cli> do-update` to upgrade now.
Next automatic check: in 12 hours.
```

If `NewRepoUrl` is set, an extra line is printed:

```
⚠  This project has moved to: https://github.com/AlifKarim/repo-v20
```

### `--async` behavior

1. Spawn the **same CLI binary** with arguments `["update-check", "--force"]`.
2. Detach: redirect stdin/stdout/stderr to `/dev/null` (or
   `NUL` on Windows), set the platform-appropriate "no-console" flag
   (`syscall.SysProcAttr{HideWindow: true, CreationFlags: DETACHED_PROCESS}`
   on Windows; `setsid()` on POSIX).
3. **Do not** `Wait()` on the child. Return exit code 0 immediately.
4. The child writes its own output to the file-system log only.

### `--force` behavior

Skips `IsCheckDue()`. Used by:
* The `--async` child (so the interval gate doesn't short-circuit it).
* Power users who want to re-check immediately.

---

## 2. `do-update`

### Synopsis

```
<cli> do-update
```

No flags. **Unattended** — no interactive confirmation. The trailing
warning printed by prior commands serves as notice. (Resolved
Decision #2.)

### Behavior

1. Read the last persisted row via `service.GetLastResult()`.
2. If `HasUpdate == false`, print "Already up to date." and exit 0.
3. If `Selected == nil`, print "No update payload available — run
   `update-check --force` first." and exit 1.
4. Detect the current OS; pick `Install.Windows` or `Install.Unix`.
5. Execute the install command. The command is the **pinned installer
   one-liner** from [16-generic-release/08](../../16-generic-release/08-version-pinned-release-installers.md)
   — never the "latest" redirect.
6. On success: set `UpdateStatusId = UpdateApplied`, log to file, print
   success message.
7. On failure: set `UpdateStatusId = Failed`, write `ErrorMessage`,
   log to file, exit non-zero.

### Output

```
Updating <cli> from V1.5.0 to V1.7.0…
Downloading installer…
Verifying checksum (Sha256)…
Installing…
✔  Update complete. Restart any open shells to use the new version.
```

---

## 3. Exit Codes

| Code | Meaning |
|------|---------|
| 0    | Success (or `--async` dispatch succeeded) |
| 1    | No update payload available / pre-condition failed |
| 2    | Network error (all six probes failed) |
| 3    | Persistence error |
| 4    | `do-update` install command failed |

---

## 4. Conformance

* The commands MUST be visible in `<cli> --help`.
* They MUST NOT require any DB/network setup beyond what the rest of
  the CLI already needs — `update-check` works on a fresh install via
  the JSON fallback (see [09-json-fallback-store.md](./09-json-fallback-store.md)).

---

*CLI Commands — v1.0.0 — 2026-04-20*
