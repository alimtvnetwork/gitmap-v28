# startup-remove (sr)

Remove a single user-scoped autostart entry that was created by
gitmap. Third-party entries are NEVER touched, even if you pass
their name.

## Synopsis

```
gitmap startup-remove <name>
gitmap startup-remove --dry-run <name>
gitmap startup-remove --output=json <name>
gitmap sr <name>
```

`<name>` is the entry name as printed by `gitmap startup-list` — the
basename without the platform extension. A trailing platform
extension is tolerated for convenience:

- Linux/Unix: `gitmap-foo` or `gitmap-foo.desktop`
- macOS: `gitmap.foo` or `gitmap.foo.plist`

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--dry-run` | `false` | Show what would be deleted (or refused/no-op) without touching the filesystem |
| `--backend` | (try all) | Windows only: `registry` (HKCU per-user), `registry-hklm` (HKLM machine-wide; **writes require admin**, dry-run does not), or `startup-folder`. Scopes the removal to one backend; without it, gitmap probes HKCU registry, then HKLM registry, then the Startup folder, and removes the first match |
| `--output` | `terminal` | Output mode: `terminal` (legacy human lines) or `json` (status object — see below) |
| `--json-indent` | `2` | Spaces per indent level for `--output=json` (`0` = minified single line). Range: 0..8. Ignored for `--output=terminal` |

`--dry-run` runs the full classification (existence + marker check)
but skips the actual unlink. The same four outcomes are reported,
each prefixed with `(dry-run)` so log scrapers can distinguish a
preview from a real action.

`--backend` is ignored on Linux/macOS — those OSes have one
canonical autostart backend each (`.desktop` and `.plist`
respectively). On Windows, passing `--backend=registry` will NOT
fall back to the Startup folder if the registry entry is missing
(returns `(no-op)` instead). This is deliberate: scoped removal
must never silently delete an entry the user did not target.

## `--output=json`

Emits a single-element JSON array containing one consistent status
object — the SAME shape `startup-add --output=json` produces — so a
single jq filter handles both commands:

```json
[
  {
    "command": "startup-remove",
    "action": "deleted",
    "name": "watch",
    "target": "/home/me/.config/autostart/gitmap-watch.desktop",
    "owner": "gitmap",
    "force_used": false,
    "dry_run": false
  }
]
```

Field meaning:

- **`command`** — `"startup-remove"`.
- **`action`** — one of `deleted`, `noop`, `refused`, `bad_name`.
- **`name`** — the positional name you passed.
- **`target`** — absolute file path (or `HKCU\...` registry path on
  Windows). Empty string for `noop` / `bad_name`.
- **`owner`** — `gitmap` (we deleted it), `none` (no entry by that
  name), `third-party` (refused), or `unknown` (bad name).
- **`force_used`** — always `false` for `startup-remove` (kept in the
  schema so add/remove records have rectangular fields).
- **`dry_run`** — `true` when `--dry-run` was passed.

Key order is byte-locked across Go versions (stablejson encoder).

## Outcomes

### Live (default)

| Status   | Message prefix | Meaning                                              | Exit |
|----------|----------------|------------------------------------------------------|------|
| Removed  | `✓ Removed`    | File existed, carried the gitmap marker, was deleted | 0    |
| No-op    | `(no-op)`      | No file by that name in the autostart dir            | 0    |
| Refused  | `(refused)`    | File exists but lacks the gitmap marker              | 0    |
| Bad name | `(refused)`    | Name is empty or contains a path separator           | 0    |

### Dry-run (`--dry-run`)

| Status   | Message prefix       | Meaning                                                       | Exit |
|----------|----------------------|---------------------------------------------------------------|------|
| Removed  | `(dry-run) would`    | File would be deleted on a live run                           | 0    |
| No-op    | `(dry-run) no...`    | No file by that name — nothing to remove                      | 0    |
| Refused  | `(dry-run) ... NOT`  | File exists but lacks the gitmap marker — would refuse        | 0    |
| Bad name | `(dry-run) name ...` | Name is empty or contains a path separator — would refuse     | 0    |

All eight outcomes exit 0 — the command is idempotent and safe to
script under both modes. A real I/O error (permission denied, etc.)
exits 1.

## Safety

- The marker is re-checked at remove time (not trusted from a stale
  `startup-list` snapshot), so a race between listing and removing
  cannot trick the command into deleting a third-party file that
  appeared after the listing. `--dry-run` runs the SAME re-check —
  preview accuracy is identical to a live run.
- The marker grammar is OS-specific:
  - Linux/Unix: `X-Gitmap-Managed=true` line in the `.desktop` body.
  - macOS: `<key>XGitmapManaged</key><true/>` in the `.plist` dict.
- Names containing `/`, `\`, or NUL bytes are rejected up-front to
  prevent path traversal — including under `--dry-run`.

## Examples

```sh
# Preview what a removal would do, without touching the file:
gitmap startup-remove --dry-run gitmap-sync-watcher
#   (dry-run) would remove gitmap-managed autostart entry: /home/me/.config/autostart/gitmap-sync-watcher.desktop

# Then commit:
gitmap startup-remove gitmap-sync-watcher
#   ✓ Removed gitmap-managed autostart entry: /home/me/.config/autostart/gitmap-sync-watcher.desktop
```

## Platform notes

Linux/Unix and macOS are supported. On Windows the command exits
with a clear "unsupported OS" message.

### macOS LaunchAgent caveat

`startup-remove` deletes the `.plist` file but does NOT call
`launchctl unload`. A removed agent takes effect at the next login
or after a manual `launchctl unload <path>` while the user's GUI
session is active. This keeps automated uninstall scripts working in
CI / SSH sessions where `launchctl` is unavailable. `--dry-run` does
not call `launchctl` either.

## Windows backends reference

For the full ownership-detection model and how the
`<hive>\Software\Gitmap\Startup*` tracking subkey makes
`startup-remove` idempotent (missing entry → `noop`, third-party
entry → `refused`, managed entry → `deleted` of both the Run-key
value / `.lnk` and the tracking subkey), see *Windows backends*
in [startup-add](startup-add.md). `--force` is intentionally NOT
accepted here — `startup-remove` will never delete a non-tracked
entry, regardless of any flag.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter startup-remove
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
