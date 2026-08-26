# UpdateCheckerService — Reusable Module Contract

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## Purpose

The `UpdateCheckerService` is the single, reusable class that owns the
entire detection pipeline. CLI commands, the pre-command hook, and any
future integrations all go through this service — no command handler
talks to HTTP, the database, or the JSON fallback directly.

---

## 1. Location

| Language | File path |
|----------|-----------|
| Go       | `internal/update/update_checker_service.go` |
| TypeScript | `src/update/UpdateCheckerService.ts` |
| PHP      | `src/Update/UpdateCheckerService.php` |
| Rust     | `src/update/update_checker_service.rs` (struct named in PascalCase) |

The service MUST be free of CLI-framework imports (no `cobra`, `clap`,
`commander`, etc.). It receives configuration via its constructor.

---

## 2. Constructor / Dependencies

```go
type UpdateCheckerServiceDeps struct {
    Owner            string         // GitHub owner (user or org)
    OwnerKind        string         // "User" | "Organization"
    CurrentRepo      string         // e.g., "repo-v15"
    CurrentVersion   string         // e.g., "V1.5.0"
    CliName          string         // for ~/.<CliName>/...
    HttpClient       HttpDoer       // injectable for tests
    Clock            Clock          // injectable for tests
    Logger           Logger         // file-system logger
    Store            UpdateStore    // SQLite or JSON fallback
    Config           SeedableConfig // 06-Seedable-Config reader
}
```

`UpdateStore` is an interface satisfied by both the SQLite ORM
implementation and the JSON-file implementation.

---

## 3. Public Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `Initialize()` | `error` | Ensures DDL, seeds `UpdateStatus`, validates config. Idempotent. |
| `IsCheckDue()` | `(bool, error)` | True iff `Now ≥ NextCheckDueAt` OR no row exists. |
| `CheckForUpdate(ctx)` | `(CombinedJson, error)` | Runs the six parallel probes, returns the combined JSON. Does **not** persist. |
| `PersistResult(combined)` | `error` | Writes to the store inside try/catch. Updates `LastCheckedAt`, `NextCheckDueAt`, `HasUpdate`, `UpdateStatusId`. |
| `RunAndPersist(ctx)` | `(CombinedJson, error)` | `CheckForUpdate` + `PersistResult` in one call. Used by both sync and `--async`. |
| `GetLastResult()` | `(CombinedJson, error)` | Reads the last persisted combined JSON for hooks/CLI output. |
| `GetPendingWarning()` | `(string, bool)` | If `HasUpdate`, returns the trailing warning line. |
| `ApplyUpdate(ctx)` | `error` | Implements `do-update`: loads OS-specific install command, executes it, records `UpdateApplied` or `Failed`. |

> No method blocks indefinitely. Every HTTP and DB call is bounded by
> the timeouts declared in [01-fundamentals.md](./01-fundamentals.md) §5.

---

## 4. Error Handling Contract

1. Every public method is wrapped in try/catch (or `defer recover` /
   `Result<>`); panics never escape.
2. On failure, the service:
   a. Writes a structured line to `~/.<CliName>/Logs/UpdateChecker.log`.
   b. Updates `UpdateChecker.ErrorMessage` and `ErrorAt`.
   c. Sets `UpdateStatusId = Failed` (4).
   d. Returns the original error to the caller (no swallowing — see
      Code Red P1).
3. A failed check **does not** clear `LatestVersion` / `HasUpdate` from
   the previous successful check. The trailing warning continues to
   surface until a successful check supersedes it.

See [08-error-handling.md](./08-error-handling.md) for the full policy.

---

## 5. Concurrency

* `CheckForUpdate` dispatches the six probes via the language's native
  concurrency primitive (`errgroup.Group`, `Promise.all`,
  `tokio::join!`, etc.).
* `PersistResult` takes a single store-level lock (SQLite `BEGIN
  IMMEDIATE` or file lock) so two concurrent writers cannot interleave.
* The service is safe to instantiate twice in the same process; the
  store lock prevents double-writes if the user runs `update-check`
  while a background check is also in flight.

---

## 6. Test Surface

The service exposes the `HttpDoer` and `Clock` seams specifically so
unit tests can:

1. Simulate the V+0..V+5 probe matrix with arbitrary 200/404 mixes.
2. Assert that the highest valid version wins.
3. Assert `NextCheckDueAt = LastCheckedAt + CheckIntervalHours`.
4. Assert that errors land in both the log file and the DB column.
5. Assert that `Failed` status does **not** clobber a prior
   `UpdateFound` warning.

---

*UpdateCheckerService Contract — v1.0.0 — 2026-04-20*
