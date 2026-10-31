# gitmap autobot

Instantly configure the machine Git identity to the GitMap automation bot (`GitMap automation <info@rgitmap.com>`).

## Usage

```bash
gitmap autobot
gitmap autobot [--global|--project] [--email <email>]
gitmap user autobot [--global|--project] [--email <email>]
gitmap autobot --help
```

## Flags

| Flag | Description |
|---|---|
| `--global` | Apply the autobot identity globally |
| `--project` | Bind the autobot identity to the current repository |
| `--email <email>` | Use a custom email for the autobot identity |
| `--empty-email` | Set the autobot identity with an empty email |

## Examples

```bash
gitmap autobot
gitmap autobot --global
gitmap user info
```
