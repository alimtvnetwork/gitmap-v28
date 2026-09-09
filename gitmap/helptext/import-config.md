# gitmap import-config

Import developer tool and application configuration files from a JSON bundle into native platform directories.

## Simulation

```
$ gitmap import-config uttorrent
  ✓ Imported utorrent configuration from uttorrent.json (1 file(s) restored)
$ gitmap import-config all ./my-configs/
  ✓ Batch import complete: 3 tool configuration(s) processed
```

## Usage

```bash
gitmap import-config <tool|all> [path]
gitmap config-import <tool|all> [path]
gitmap improt-config <tool|all> [path]
```

## Supported Tools

| Tool | Aliases | Default Source | Description |
|------|---------|----------------|-------------|
| vscode | code, vcode | vscode.json | Restores settings.json, keybindings.json, extensions.txt |
| qtorrent | qbittorrent, qbit | qtorrent.json | Restores qBittorrent.ini (Windows) or qBittorrent.conf (Linux/macOS) |
| utorrent | uttorrent, u-torrent | uttorrent.json / utorrent.json | Restores settings.dat |
| all | * | <folder>/*.json | Scans folder and imports all discovered configurations |

## Examples

### Import uTorrent settings from uttorrent.json

```bash
gitmap import-config uttorrent
```

### Import qBittorrent settings from custom path

```bash
gitmap import-config qtorrent ./my-configs/qtorrent.json
```

### Import VS Code settings

```bash
gitmap import-config vscode
```

### Batch import all configurations from a folder

```bash
gitmap import-config all ./my-configs/
```
