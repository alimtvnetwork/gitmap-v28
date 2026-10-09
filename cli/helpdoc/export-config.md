# gitmap export-config

Export developer tool and application configuration files into a portable JSON bundle.

## Simulation

```
$ gitmap export-config qtorrent
  ✓ Exported qtorrent configuration to qtorrent.json (1 file(s))
$ gitmap export-config vscode ./backup/
  ✓ Exported vscode configuration to backup/vscode.json (2 file(s))
```

## Usage

```bash
gitmap export-config <tool|all> [path]
gitmap config-export <tool|all> [path]
```

## Supported Tools

| Tool | Aliases | Default File | Description |
|------|---------|--------------|-------------|
| vscode | code, vcode | vscode.json | VS Code settings.json, keybindings.json, extensions list |
| qtorrent | qbittorrent, qbit | qtorrent.json | qBittorrent settings (qBittorrent.ini / qBittorrent.conf) |
| utorrent | uttorrent, u-torrent | uttorrent.json / utorrent.json | uTorrent settings (settings.dat) |
| all | * | <tool>.json | Export configurations for all supported tools to a folder |

## Examples

### Export qBittorrent settings to default qtorrent.json

```bash
gitmap export-config qtorrent
```

### Export VS Code settings to default vscode.json

```bash
gitmap export-config vscode
```

### Export uTorrent settings to uttorrent.json

```bash
gitmap export-config uttorrent
```

### Export all tool configurations into a backup folder

```bash
gitmap export-config all ./my-configs/
```
