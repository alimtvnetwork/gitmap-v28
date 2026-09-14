# gitmap service

Cross-platform OS service management suite for Linux (systemd), Windows (sc.exe / Windows Services), and macOS (launchd).

## Alias

`srv`

## Usage

    gitmap service ls                          # list all managed and system services
    gitmap service status <name>               # show detailed status and state of a service
    gitmap service on <name>                   # start and enable a service
    gitmap service off <name>                  # stop and disable a service
    gitmap service create <name> <cmd...>      # create/register a new OS service
    gitmap service rm <name>                   # remove/unregister an OS service
    gitmap service export <name> [-f <file>]   # export service definition to JSON/YAML
    gitmap service import <file>               # import and register service definition

## Subcommands

| Subcommand | Description |
|---|---|
| `ls`, `list` | List installed OS services with their current status |
| `status`, `info` `<name>` | Display detailed status, PID, and configuration of a service |
| `on`, `start`, `enable` `<name>` | Start service execution and enable automatic startup |
| `off`, `stop`, `disable` `<name>` | Stop running service and disable automatic startup |
| `create`, `new`, `add` `<name> <cmd>` | Create and register a new system service |
| `rm`, `delete`, `del` `<name>` | Remove and unregister a system service |
| `export` `[name]` | Export service configuration as JSON or YAML |
| `import` `<file>` | Import service configuration from JSON or YAML |

## Examples

### List all services

```bash
gitmap service ls
```

### Check status of a service

```bash
gitmap service status nginx
```

### Start or enable a service

```bash
gitmap service on nginx
```

### Stop or disable a service

```bash
gitmap service off nginx
```
