# gitmap os dm

Inspect and configure Linux Display Managers (GDM3, LightDM, SDDM), session types, and Wayland compatibility.

## Usage

```bash
gitmap os dm [subcommand] [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| status (st) | Inspect current Display Manager, active session, and Wayland status |
| wayland <on\|off> | Enable or disable Wayland (forces X11 on GDM3 for headless / VM stability) |
| restart | Restart the display-manager service via systemctl |
| help | Display help and usage information |

## Examples

### Inspect Display Manager and Session Status

```bash
gitmap os dm status
```

### Disable Wayland (Force X11 Session on Ubuntu GDM3)

```bash
gitmap os dm wayland off
```

### Enable Wayland Session

```bash
gitmap os dm wayland on
```

### Restart Display Manager Service

```bash
gitmap os dm restart
```
