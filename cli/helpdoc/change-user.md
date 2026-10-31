# gitmap change-user

Set the Git author name and email directly, optionally scoped globally or to the current repository.

## Usage

```bash
gitmap change-user <name> [email]
gitmap user change-user [name] [email] [--global|--project]
gitmap change-user --help
```

## Flags

| Flag | Description |
|---|---|
| `--global` | Apply the identity globally |
| `--project` | Bind the identity to the current repository |
| `--name <string>` | Specify the author name |
| `--email <string>` | Specify the author email |

## Examples

```bash
gitmap change-user "Your Name" you@example.com
gitmap change-user --global "Your Name" you@example.com
gitmap user info
```
