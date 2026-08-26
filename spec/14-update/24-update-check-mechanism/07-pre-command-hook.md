# Pre-Command Hook — Interval Gate & Trailing Warning

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## Purpose

Every CLI invocation runs a tiny, **non-blocking** hook before the
user's command executes, and a tiny **post-execution** hook after it
finishes. Together they:

1. Spawn a background `update-check --async` when the interval has
   elapsed (fire-and-forget).
2. Append a one-line warning to the user's output if a pending update
   exists.

The hook is the only place CLI command code interacts with the update
subsystem — command handlers themselves remain ignorant of it.

---

## 1. Opt-Out

Controlled by `06-Seedable-Config`:

```jsonc
{
  "Update": {
    "BackgroundUpdateCheckEnabled": true,   // default
    "CheckIntervalHours": 12,               // default
    "PendingUpdateWarningEnabled": true     // default
  }
}
```

If `BackgroundUpdateCheckEnabled == false`, the pre-hook is a no-op.
If `PendingUpdateWarningEnabled == false`, the post-hook is a no-op.
The two flags are independent — a user MAY disable spawning while
keeping the warning, or vice versa.

---

## 2. Pre-Hook Algorithm

```
fn PreCommandHook(ctx, cmdName) {
    if cmdName in {"update-check", "do-update"} {
        return            // never recurse
    }
    if not Config.BackgroundUpdateCheckEnabled {
        return
    }
    due, err := service.IsCheckDue()
    if err != nil {
        Logger.Warn("PreHook IsCheckDue failed", err)
        return            // never block the user's command
    }
    if not due {
        return
    }
    SpawnDetached(<cli>, ["update-check", "--async"])
    return                // do not Wait()
}
```

### Hard rules

1. The hook MUST return in **< 50 ms** on the happy path. Spawning is
   the only side effect.
2. Any error from `IsCheckDue` is logged and swallowed at this
   boundary — the user's command must run regardless. (This is the
   single permitted "swallow at boundary" exception to Code Red P1,
   and it MUST be logged.)
3. The hook MUST NOT `Wait()` on the child. The child is fully
   detached (stdin/stdout/stderr → null device).
4. The hook MUST NOT print anything to the user's stdout/stderr.

---

## 3. Post-Hook Algorithm

```
fn PostCommandHook(ctx, cmdName) {
    if cmdName in {"update-check", "do-update"} {
        return
    }
    if not Config.PendingUpdateWarningEnabled {
        return
    }
    msg, has := service.GetPendingWarning()
    if not has {
        return
    }
    fmt.Fprintln(os.Stderr, msg)
}
```

### Warning line format

```
⚠  An update is available (V1.5.0 → V1.7.0). Run `<cli> do-update`. Next check in 12h.
```

If `Selected.NewRepoUrl` is non-null, a second line is appended:

```
↪  This project has moved to https://github.com/AlifKarim/repo-v20
```

The warning goes to **stderr** so it never contaminates piped stdout
(`<cli> list | jq …` keeps working).

---

## 4. Integration Points

| CLI Framework | Pre-hook | Post-hook |
|---------------|----------|-----------|
| `cobra` (Go) | `PersistentPreRunE` on the root command | `PersistentPostRunE` |
| `clap` (Rust) | Wrapper in `main()` before `App::run()` | After `App::run()` |
| `commander` (TS) | `program.hook("preAction", …)` | `program.hook("postAction", …)` |
| `symfony/console` (PHP) | `ConsoleEvents::COMMAND` listener | `ConsoleEvents::TERMINATE` listener |

The hook implementation lives next to `UpdateCheckerService` (same
package/folder); it is not duplicated per command.

---

## 5. Interaction With `--async` Child

The detached child runs `update-check --force` (so it skips its own
interval gate) and writes silently. The parent process has already
returned to the user's shell long before the child finishes. If the
child crashes, the next pre-hook invocation will see `IsCheckDue ==
true` again and re-spawn — there is no retry loop, no backoff, and no
state machine. The interval gate **is** the rate limit.

---

*Pre-Command Hook — v1.0.0 — 2026-04-20*
