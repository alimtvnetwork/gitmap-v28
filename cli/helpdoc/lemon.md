# gitmap lemon

Background lemon daemon and email command service.

## Alias

lm

## Usage

```bash
gitmap lemon run [--interval 8s] [--once] [--verbose]
gitmap lemon status [--json]
gitmap lemon poll
gitmap lemon service <install|uninstall|start|stop|restart|status>
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --interval, -i | 8s | Polling interval (clamped to minimum 5s guard) |
| --once | false | Run a single polling pass and exit |
| --verbose, -v | false | Enable verbose logging output |
| --json | false | Output status as structured JSON |

## Examples

### Example 1: Run lemon background daemon in foreground

```bash
gitmap lemon run --interval 8s
```

### Example 2: Check status of lemon service

```bash
gitmap lemon status
```

### Example 3: Single-shot poll pass

```bash
gitmap lemon poll
```
