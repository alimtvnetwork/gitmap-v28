# Error Handling Policy

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## Purpose

Defines the **end-to-end error contract** for the update-check
subsystem: where errors are caught, where they are logged, what is
persisted, and what the user sees.

This spec is bound by the project's 🔴 CODE RED principles — errors
are **never silently swallowed** except at one named boundary
(the pre-command hook, §4).

---

## 1. Error Categories

| Category | Trigger | Handled by |
|----------|---------|------------|
| `NetworkError`     | HTTP timeout, DNS failure, connection refused | `CheckForUpdate` |
| `HttpStatusError`  | Non-200, non-404 response | `CheckForUpdate` |
| `JsonParseError`   | 200 response with malformed JSON | `CheckForUpdate` |
| `ValidationError`  | JSON missing required fields (see [02 §Validation](./02-status-script-json.md#validation)) | `CheckForUpdate` |
| `PersistenceError` | DB write / file write failure | `PersistResult` |
| `InstallError`     | `do-update` install command non-zero exit | `ApplyUpdate` |

---

## 2. Per-Probe Failure (Discovery)

A single probe failure is **not** a subsystem error — it is a normal
data point.

```
For each probe in [V, V+1, V+2, V+3, V+4, V+5]:
    result = HttpGet(url, timeout=5s)
    if result.IsTimeout or result.Status == 404:
        Candidates.append({ Version, Found: false, Url: null, Reason: "NotFound" })
        continue
    if result.Status != 200:
        Logger.Warn("Probe non-200", url, result.Status)
        Candidates.append({ Version, Found: false, Url: null, Reason: "Http" })
        continue
    if not TryParseJson(result.Body):
        Logger.Warn("Probe malformed JSON", url, parseErr)
        Candidates.append({ Version, Found: false, Url: null, Reason: "Parse" })
        continue
    Candidates.append({ Version, Found: true, Url: url, Payload: parsed })
```

Discovery succeeds as long as **at least the V probe** completes (so
the CLI knows its own current state). If even V fails, the whole
`CheckForUpdate` returns a `NetworkError` and the row is marked
`Failed` (see §3).

---

## 3. Subsystem-Level Failure

When `CheckForUpdate` or `PersistResult` returns an error:

1. **File-system log.** Append a structured line to
   `~/.<CliName>/Logs/UpdateChecker.log`:

   ```
   2026-04-20T12:34:56Z LEVEL=Error Op=CheckForUpdate Code=NetworkError
     Owner=AlifKarim Repo=repo-v15 Message="dial tcp: i/o timeout"
     File=internal/update/probe.go Line=87
   ```

   File path is mandatory. PascalCase key constants apply
   (`LogKeyOp`, `LogKeyCode`, `LogKeyOwner`, …).

2. **Database row.** Update the existing `UpdateChecker` row (do not
   insert a new one):
   * `ErrorMessage = <human-readable reason>`
   * `ErrorAt = <Now UTC>`
   * `UpdateStatusId = 4` (`Failed`)
   * **Do NOT clear** `LatestVersion`, `HasUpdate`, `Selected`, or
     `RawJson` from the prior successful check. A failed re-check
     must not erase a known-good update notice.

3. **Caller propagation.** The error is returned to the caller
   verbatim (wrapped via the project's `apperror` package, with file
   and line metadata).

---

## 4. Pre-Command Hook Boundary (Single Permitted Swallow)

The pre-command hook (see [07-pre-command-hook.md](./07-pre-command-hook.md))
is the **only** place where an error from this subsystem is allowed to
be logged-and-swallowed. Rationale: the user typed `<cli> some-command`
— their command must run regardless of update-check health.

```
err := service.IsCheckDue()
if err != nil {
    Logger.Warn("PreHook IsCheckDue failed", err)   // log
    return                                          // swallow
}
```

This is the **only** exception. Every other code path obeys CODE RED
P1 (no swallowing).

---

## 5. `do-update` Failure

If the install command exits non-zero:

1. Capture stdout + stderr (truncate to 4 KiB) into `ErrorMessage`.
2. Set `UpdateStatusId = Failed`.
3. Log the full output to `~/.<CliName>/Logs/UpdateChecker.log`.
4. Exit code 4 (see [06 §3](./06-cli-commands.md#3-exit-codes)).
5. **Leave `HasUpdate = true`** so the user keeps seeing the warning
   and can retry.

---

## 6. Try/Catch Granularity

| Layer | Wrap? | Notes |
|-------|-------|-------|
| HTTP per-probe | yes | Captures `NetworkError`, `HttpStatusError`, `JsonParseError` per probe |
| Combined-JSON assembly | yes | Pure CPU; should never fail, but defends against panics |
| DB transaction (`PersistResult`) | yes | `BEGIN IMMEDIATE` / `COMMIT` or rollback on any error |
| JSON-fallback file write | yes | Atomic write via `tmp + rename`; rollback by deleting tmp |
| Install command exec | yes | Captures stdout/stderr/exitCode |

Try/catch blocks MUST NOT contain nested `if`s (CODE RED P6) — extract
helpers.

---

## 7. Log File Rotation

`UpdateChecker.log` is a single file, capped at 1 MiB, rotated to
`UpdateChecker.log.1` (overwrites). No multi-generation history. The
user can clear it at any time without affecting behavior.

---

*Error Handling Policy — v1.0.0 — 2026-04-20*
