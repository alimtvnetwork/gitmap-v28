# gitmap os

Manage operating system configurations, environment diagnostics, and symlink integrity.

## Usage

```bash
gitmap os [subcommand] [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| display [subcommand] | Inspect and configure OS display settings, desktop session & timeouts |
| fix-link [path] | Inspect and repair broken symlinks and shared directories |
| status | Display operating system summary, platform details, and link health |
| help | Show usage information for OS commands |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --target <path> | "" | Explicit target destination for broken link repair |
| --force (-f) | false | Overwrite broken link or file even if target does not yet exist |
| --recursive (-r) | false | Recursively inspect subdirectories for broken symlinks |
| --dry-run (-n) | false | Simulate inspection and repair without modifying filesystem |
| --json | false | Output status and link results as structured JSON |

## Examples

### Inspect Display Server and Screen Blanking Settings

```bash
gitmap os display
```

Output:

```text
▶ OS Display & Screen Status (linux)
  • Operating System: linux (amd64)
  • Display Server:   wayland
  • Desktop Session:  ubuntu:GNOME
  • Display Timeout:  Never (inhibited)
  • Sleep Timeout:    Never (inhibited)
  • Mode:             Never-Sleep (inhibited)
```

### Disable Display Sleep and Blanking Timeouts

```bash
gitmap os display never-sleep
```

Output:

```text
✓ Power setting updated: Never-Sleep mode enabled (display & standby timeouts disabled).
```

### Set Screen Timeout to 15 Minutes

```bash
gitmap os display set 15
```

Output:

```text
✓ Power setting updated: Display timeout=15m, Sleep timeout=15m.
```

### Check OS and Symlink Status

```bash
gitmap os status
```

Output:

```text
▶ OS Environment:
  • Operating System: linux (ubuntu 24.04)
  • Architecture:     amd64
  • User Home:        /home/developer
  • Desktop Links:    1 healthy, 0 broken
```

### Repair Desktop Shared Directories Symlink

```bash
gitmap os fix-link ~/Desktop/SharedDirectories --target /mnt/hgfs
```

Output:

```text
▶ Inspecting and repairing symlinks...
  ✓ Repaired: /home/developer/Desktop/SharedDirectories -> /mnt/hgfs (was dangling)
  • Summary: 1 checked, 1 repaired, 0 broken
```
