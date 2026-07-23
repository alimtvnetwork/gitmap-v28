# 26 — Startup Management

> Cross-platform "run on login / boot" registration for arbitrary
> commands and persistent environment variables. Available on Windows,
> Linux, and macOS.
>
> Introduced in **v3.125.0**.

## 1. Goals

| Goal | Rationale |
|------|-----------|
| Single CLI surface for registering startup entries on any OS | Users shouldn't have to learn registry edits, `systemctl --user enable`, or `launchctl load` |
| Multiple backend mechanisms per OS, user-selectable | Different real-world constraints (no admin, headless server, GUI app, delayed start) need different mechanisms |
| Every change logged in `StartupEntries` table | Audit trail + reliable `startup list` regardless of which backend created the entry |
| Symmetric `add` / `remove` / `list` / `enable` / `disable` | Consistent verbs across backends |
| User-scope by default; `--system` opt-in for machine-wide | Matches `gitmap-v27 self-install` philosophy — no admin required for the common case |

## 2. CLI Surface

Both forms are first-class. The umbrella form is preferred for humans;
flat aliases keep parity with existing flat verbs (`self-install`,
`clone-next`, etc.).

### Umbrella

```
gitmap-v27 startup add      <name> <command> [--method <m>] [--scope user|system] [--delay <duration>]
gitmap-v27 startup add      --interactive
gitmap-v27 startup list     [--json]
gitmap-v27 startup remove   <name> [--method <m>]
gitmap-v27 startup enable   <name>
gitmap-v27 startup disable  <name>
gitmap-v27 startup env-add  <KEY=VALUE> [--scope user|system]
gitmap-v27 startup env-list [--json]
gitmap-v27 startup env-rm   <KEY>      [--scope user|system]
```

### Flat aliases

```
gitmap-v27 startup-add        ⇄ gitmap-v27 startup add
gitmap-v27 startup-add-i      ⇄ gitmap-v27 startup add --interactive
gitmap-v27 startup-list       ⇄ gitmap-v27 startup list
gitmap-v27 startup-remove     ⇄ gitmap-v27 startup remove
gitmap-v27 startup-enable     ⇄ gitmap-v27 startup enable
gitmap-v27 startup-disable    ⇄ gitmap-v27 startup disable
gitmap-v27 env-add            ⇄ gitmap-v27 startup env-add
gitmap-v27 env-list           ⇄ gitmap-v27 startup env-list
gitmap-v27 env-rm             ⇄ gitmap-v27 startup env-rm
```

## 3. Backends

### Windows (`runtime.GOOS == "windows"`)

| Method ID | Mechanism | Scope | Notes |
|-----------|-----------|-------|-------|
| `registry-run` | `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` (or `HKLM` for `--system`) | user / system | Default; no admin for HKCU |
| `startup-folder-lnk` | `.lnk` shortcut in `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup` | user only | Generated via PowerShell `WScript.Shell` |
| `startup-folder-cmd` | Plain `.cmd` file in same Startup folder | user only | Easiest to inspect, no PowerShell dependency |
| `task-scheduler` | `schtasks.exe /create /sc ONLOGON` (or `ONSTART` for system) | user / system | Only backend supporting `--delay` and `/RL HIGHEST` |

### Linux (`runtime.GOOS == "linux"`)

| Method ID | Mechanism | Scope | Notes |
|-----------|-----------|-------|-------|
| `systemd-user` | `~/.config/systemd/user/<name>.service` + `systemctl --user enable` | user | Preferred when `systemctl` is on PATH |
| `desktop-autostart` | `~/.config/autostart/<name>.desktop` | user | XDG-compliant; works for GNOME/KDE/XFCE GUI sessions |
| `cron-reboot` | `crontab -l` append `@reboot <cmd>` | user / system | Always-works fallback even on minimal systems |

### macOS (`runtime.GOOS == "darwin"`)

| Method ID | Mechanism | Scope | Notes |
|-----------|-----------|-------|-------|
| `launchagent` | `~/Library/LaunchAgents/<name>.plist` + `launchctl load` | user | Default; supports `RunAtLoad` and `KeepAlive` |
| `shell-rc` | Append `# gitmap-startup:<name>` block to `~/.zshrc` / `~/.bashrc` | user | Terminal-launch-only; for env-var-style entries |

## 4. Method Resolution

When `--method` is omitted:

| OS | Default |
|----|---------|
| Windows | `registry-run` |
| Linux | `systemd-user` if `systemctl` on PATH, else `cron-reboot` |
| macOS | `launchagent` |

When `--interactive` is passed, `gitmap-v27 startup add` prompts the user
to pick from the list of methods supported on the current OS, then
prompts for any method-specific options (delay, elevated, etc.).

## 5. The `StartupEntries` Table

Every add/remove operation, regardless of backend, writes a row so
`gitmap-v27 startup list` can report a unified view without scraping
each backend.

```sql
CREATE TABLE IF NOT EXISTS StartupEntries (
    ID            INTEGER PRIMARY KEY AUTOINCREMENT,
    Name          TEXT NOT NULL UNIQUE,
    Command       TEXT NOT NULL,
    Method        TEXT NOT NULL,            -- one of the Method IDs above
    Scope         TEXT NOT NULL,            -- 'user' or 'system'
    Enabled       INTEGER NOT NULL DEFAULT 1,
    DelaySeconds  INTEGER NOT NULL DEFAULT 0,
    BackendPath   TEXT,                     -- absolute path or registry key
    CreatedAt     TEXT NOT NULL,
    UpdatedAt     TEXT NOT NULL
);
```

A second table tracks env vars:

```sql
CREATE TABLE IF NOT EXISTS StartupEnvVars (
    ID         INTEGER PRIMARY KEY AUTOINCREMENT,
    KeyName    TEXT NOT NULL,
    KeyValue   TEXT NOT NULL,
    Scope      TEXT NOT NULL,                -- 'user' or 'system'
    Backend    TEXT NOT NULL,                -- registry / etc-environment / shell-rc
    UNIQUE (KeyName, Scope)
);
```

## 6. Error Handling

All errors follow the standard format documented in
`spec/04-generic-cli/07-error-handling.md`:

```
Error: [message] at [path]: [reason] (operation: [op], reason: [why])
```

Constants live in `gitmap-v27/constants/constants_startup.go` and are
named `ErrStartup*` for return-wrapping (with `%w`) and
`ErrStartupFmt*` for direct stderr formatting (with trailing `\n`).

## 7. Cross-References

- `spec/04-generic-cli/07-error-handling.md` — error format
- `spec/04-generic-cli/15-constants-reference.md` — `Err*` naming
- `spec/04-generic-cli/22-self-update-gold-standard.md` — `gitmap-v27 self-install` precedent for OS-specific install logic
- `spec/04-generic-cli/21-post-install-shell-activation.md` — shell-rc append/remove pattern reused for `shell-rc` backend
