# gitmap rest-enable

Enables and starts the background REST OS daemon service (`gitmap-daemon`)
to enable cross-machine automation, remote prompt injection, and cluster coordination.

## Usage

    gitmap rest-enable
    gitmap agy rest-enable

## Supported Platforms

- Windows: Windows Service Controller (`winsvc`)
- Linux: Systemd unit
- macOS: Launchd daemon plist

## Examples

```bash
# Enable background REST service
gitmap rest-enable

# Enable via AGY subsystem
gitmap agy rest-enable
```

