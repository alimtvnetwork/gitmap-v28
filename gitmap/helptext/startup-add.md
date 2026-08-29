# gitmap startup-add

Create a Linux/Unix XDG autostart entry that runs gitmap (or any
command) at login. The created `.desktop` file is tagged with
`X-Gitmap-Managed=true` so `startup-list` and `startup-remove` can
safely manage it without touching third-party autostart files.

## Alias

sa

## Usage

    gitmap startup-add --name <id> [--exec <path>] [--display-name <s>]
                       [--comment <s>] [--working-dir <path>]
                       [--backend registry|registry-hklm|startup-folder]
                       [--no-display] [--force]

## Flags

| Flag | Required | Description |
|------|----------|-------------|
| --name           | yes | Logical name; filename becomes `gitmap-<name>.desktop` |
| --exec           | no  | Command to run at login (default: path to running gitmap binary) |
| --display-name   | no  | Override the `Name=` field shown in session managers |
| --comment        | no  | Optional `Comment=` text |
| --working-dir    | no  | Working directory the entry runs in (see *Working directory* below) |
| --no-display     | no  | Set `NoDisplay=true` (hide from app menus, still autostarts) |
| --force          | no  | Overwrite an existing **gitmap-managed** entry (never overwrites third-party files) |
| --backend        | no  | Windows only: `registry` (default, HKCU per-user), `registry-hklm` (HKLM machine-wide; **requires admin**), or `startup-folder` |
| --output         | no  | Output mode: `terminal` (default human lines) or `json` (status object — see below) |
| --json-indent    | no  | Spaces per indent level for `--output=json` (`0` = minified). Range 0..8. Ignored for terminal |

## `--output=json`

Emits a single-element JSON array containing one consistent status
object — the SAME shape `startup-remove --output=json` produces, so
a single jq filter handles both:

```json
[
  {
    "command": "startup-add",
    "action": "created",
    "name": "watch",
    "target": "/home/me/.config/autostart/gitmap-watch.desktop",
    "owner": "gitmap",
    "force_used": false,
    "dry_run": false
  }
]
```

- **`action`** — one of `created`, `overwritten`, `exists`, `refused`, `bad_name`.
- **`owner`** — `gitmap` (we own the entry), `third-party` (refused), or `unknown` (bad name).
- **`target`** — absolute file path or `HKCU\...` registry path; empty for `bad_name`.
- **`force_used`** — reflects whether `--force` was passed.
- **`dry_run`** — always `false` for `startup-add` (kept so add/remove records are rectangular).

Key order is byte-locked across Go versions (stablejson encoder).

## Working directory

`--working-dir <path>` records a directory the entry should run in.
The value is rendered differently per OS but is always read back by
`startup-list`:

- **Linux/Unix**: written as `Path=<dir>` in the `.desktop` file
  (XDG-spec field). The session manager `chdir`s here before
  invoking `Exec=`.
- **macOS**: written as `<key>WorkingDirectory</key>` in the
  LaunchAgent plist. `launchd` `chdir`s here before exec'ing
  `ProgramArguments`.
- **Windows**: stored as a `WorkingDir` REG_SZ value in the gitmap
  tracking subkey at `HKCU\Software\Gitmap\StartupRegistry\<name>`
  (registry backend) or `HKCU\Software\Gitmap\StartupFolder\<name>`
  (startup-folder backend). The autostart command itself (Run-key
  value or `.lnk` target) is unchanged — Windows reads cwd from the
  `.lnk` `WorkingDirectory` field, which the current minimal Shell
  Link writer does not yet emit; the tracking-subkey value is the
  source of truth for tooling.

Pass an absolute path. Relative paths are accepted as-is and
interpreted by the OS at login time. Omit the flag (or pass `""`)
to inherit whatever directory the login session provides.

## Windows backends

| `--backend`         | Hive | Scope             | Admin? | Path |
|---------------------|------|-------------------|--------|------|
| `registry` (default)| HKCU | Current user      | no     | `HKCU\Software\Microsoft\Windows\CurrentVersion\Run\gitmap-<name>` |
| `registry-hklm`     | HKLM | Every user (machine-wide) | **yes** | `HKLM\Software\Microsoft\Windows\CurrentVersion\Run\gitmap-<name>` |
| `startup-folder`    | —    | Current user      | no     | `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup\gitmap-<name>.lnk` |

The `registry-hklm` backend writes the autostart value under
`HKEY_LOCAL_MACHINE` so it fires for **every interactive user** on
the machine — useful for shared workstations, kiosks, and lab
images. Tracking metadata is mirrored at
`HKLM\Software\Gitmap\StartupRegistry\<name>` so `startup-list
--backend=registry-hklm` can attribute each entry to gitmap
without touching `HKCU`.

Writes (add / remove) require administrator privileges. The
process token is checked **before** any registry mutation; if you
are not elevated the command exits with a friendly:

    startup-add: --backend=registry-hklm requires administrator privileges
    (re-run from an elevated shell, e.g. `Run as administrator` from the
    Start menu, or use the per-user `--backend=registry` default)

`startup-list --backend=registry-hklm` and `startup-remove
--backend=registry-hklm --dry-run` are read-only and work for any
user (no elevation needed).

### Ownership detection (Windows)

Every Windows backend uses a **two-gate** check to decide whether
an existing entry belongs to gitmap. Both gates must pass for the
entry to be considered "managed by gitmap"; otherwise it is
treated as **third-party** and refused (even with `--force`).

| Backend | Gate 1 (autostart surface) | Gate 2 (ownership marker) |
|---------|----------------------------|---------------------------|
| `registry`       | Run-key value `HKCU\…\Run\gitmap-<name>` exists | Tracking subkey `HKCU\Software\Gitmap\StartupRegistry\<name>` exists |
| `registry-hklm`  | Run-key value `HKLM\…\Run\gitmap-<name>` exists | Tracking subkey `HKLM\Software\Gitmap\StartupRegistry\<name>` exists |
| `startup-folder` | `gitmap-<name>.lnk` file exists in the user's Startup folder | Tracking subkey `HKCU\Software\Gitmap\StartupFolder\<name>` exists |

The Run-key value itself **never** carries an inline marker.
Windows treats every value under `…\Run` as an autostart command
and feeds it to the shell at login, so a sibling
`gitmap-<name>.gitmap-managed = "true"` value would surface in
Task Manager's Startup tab and be dispatched as the literal
command `true`. Keeping the marker in a separate scope under
`<hive>\Software\Gitmap` lets the Run key contain only real
autostart commands — exactly what a hand-edited entry would look
like.

The classifier returns one of three states for each candidate:

| State        | Gate 1 | Gate 2 | `startup-add` outcome |
|--------------|:-:|:-:|----------------------|
| **none**     | ✗  | ✗  | `created` (fresh write) |
| **gitmap**   | ✓  | ✓  | `exists` (no `--force`) / `overwritten` (with `--force`) |
| **third-party** | ✓ | ✗ | `refused` — never touched, even with `--force` |

A read error opening either key is treated as "third-party"
(refused) so an unreadable entry is never silently overwritten.

### `--force` overwrite behavior (Windows)

`--force` lifts **only one** check: the "already exists AND is
ours" guard. It does **not** weaken the ownership gate. Concretely:

- **gitmap-managed entry, no `--force`** → `exists` no-op, exit 0.
- **gitmap-managed entry, `--force`** → tracking subkey rewritten
  (with a fresh `CreatedAt`) and the Run-key value / `.lnk`
  replaced. Reported as `overwritten`.
- **Third-party entry, no `--force`** → `refused`, nothing
  written. Exit 0 so a provisioning script can keep going.
- **Third-party entry, `--force`** → still **`refused`**.
  `--force` cannot promote a non-tracked entry to managed status.

If a previous interactive user manually deleted
`<hive>\Software\Gitmap`, every existing Run-key value gitmap
created looks third-party from then on — Add will refuse to
overwrite them and `startup-list` will not show them. Re-create
the tracking subkeys (or just `startup-add … --force` on a
different `--name`) to recover; gitmap deliberately never
auto-claims a Run-key value it cannot prove ownership of.

### Idempotency via the tracking subkey

Because ownership lives in the tracking subkey, re-running the
same `startup-add` invocation is **always safe**:

```
gitmap startup-add --name watch --exec "C:\gitmap.exe watch"

# → created

gitmap startup-add --name watch --exec "C:\gitmap.exe watch"

# → exists  (no-op, exit 0 — safe to put in a provisioning script)

gitmap startup-add --name watch --exec "C:\gitmap.exe watch --quiet" --force

# → overwritten  (Run-key value updated, tracking subkey rewritten)

```

Crash safety: the tracking subkey is written **before** the
Run-key value / `.lnk`. If the process is killed between the two
writes, the next `startup-add` re-run sees a managed-but-inactive
record and safely overwrites it (gate 2 passes, gate 1 fails →
treated as `created` for Run-key purposes, the existing tracking
subkey's metadata gets refreshed). At no point can a partial
write leave a third-party-looking value behind.

`startup-remove` follows the same idempotent contract: a missing
entry is `noop` (exit 0), a third-party entry is `refused`, and
deleting a managed entry removes both the Run-key value and the
tracking subkey. A crash between those two deletes leaves an
orphaned tracking subkey that the next `startup-remove` will
clean up.

## Prerequisites

- Linux or other Unix with `~/.config/autostart` (XDG-compliant).
- macOS uses LaunchAgents — not handled here.
- Windows: any backend works for `--output=json`/terminal alike;
  `registry-hklm` additionally requires UAC elevation for writes.

## Safety

- Refuses to overwrite a `.desktop` file that does NOT carry the
  `X-Gitmap-Managed=true` marker, even with `--force`.
- On Windows, refuses to overwrite a Run-key value or `.lnk`
  whose tracking subkey under `<hive>\Software\Gitmap` is
  missing, even with `--force` — see *Ownership detection
  (Windows)* above.
- Names containing path separators (`/`, `\`) or NUL are rejected
  before any I/O.
- Atomic write (temp file + rename) so a crash mid-write cannot
  leave a half-written file the next login session would execute.
  On Windows the tracking subkey is written before the Run-key
  value / `.lnk` so a crash leaves the entry recoverable on the
  next `startup-add`.

## Examples

### Example 1: Add gitmap itself with default args

    gitmap startup-add --name watch --exec "$(command -v gitmap) watch"

**Output:**

    ✓ Created gitmap-managed autostart entry: /home/me/.config/autostart/gitmap-watch.desktop

### Example 2: Re-run is idempotent

    gitmap startup-add --name watch --exec "$(command -v gitmap) watch"

**Output:**

      (exists) gitmap-managed entry already at /home/me/.config/autostart/gitmap-watch.desktop — pass --force to overwrite

### Example 3: Update an existing entry

    gitmap startup-add --name watch \
      --exec "$(command -v gitmap) watch --quiet" --force

**Output:**

    ✓ Overwrote gitmap-managed autostart entry: /home/me/.config/autostart/gitmap-watch.desktop

### Example 4: Windows — `--force` cannot overwrite a third-party Run-key value

    gitmap startup-add --name watch --exec "C:\gitmap.exe watch" --force

**Output (third-party `gitmap-watch` value already in HKCU\…\Run, no tracking subkey):**

      (refused) "HKCU\Software\Microsoft\Windows\CurrentVersion\Run\gitmap-watch" exists but was NOT created by gitmap — refusing to overwrite

`--force` only lifts the "exists AND is ours" guard; the
ownership gate is unconditional. Rename your entry
(`--name watch2`) or delete the third-party value manually in
regedit before re-running.

### Example 5: Machine-wide install via `--backend=registry-hklm`

    gitmap startup-add --name kiosk --backend=registry-hklm \
      --exec "C:\gitmap.exe watch --profile kiosk"

**Output (run from an elevated shell):**

    ✓ Created gitmap-managed autostart entry: HKLM\Software\Microsoft\Windows\CurrentVersion\Run\gitmap-kiosk

**Output (run from a non-elevated shell):**

    startup-add: --backend=registry-hklm requires administrator privileges (re-run from an elevated shell, e.g. `Run as administrator` from the Start menu, or use the per-user `--backend=registry` default)

## See Also

- [startup-list](startup-list.md) — List entries gitmap created
- [startup-remove](startup-remove.md) — Delete a gitmap-managed entry

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter startup-add
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
