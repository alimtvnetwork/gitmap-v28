# Environment Variables

Manage persistent environment variables and PATH entries across platforms.

## Alias

ev

## Usage

    gitmap env <subcommand> [flags]

## Subcommands

| Subcommand  | Description                              |
|-------------|------------------------------------------|
| set         | Set a persistent environment variable    |
| get         | Get a managed variable's value           |
| delete      | Remove a managed variable                |
| list        | List all managed variables               |
| path add    | Add a directory to PATH                  |
| path remove | Remove a directory from PATH             |
| path list   | List managed PATH entries                |

## Flags

| Flag      | Default | Description                                          |
|-----------|---------|------------------------------------------------------|
| --system  | false   | Target system-level variables (Windows, requires admin) |
| --shell   | (auto)  | Target shell profile: bash, zsh (Unix only)          |
| --verbose | false   | Show detailed operation output                       |
| --dry-run | false   | Preview changes without applying                     |

## Prerequisites

- Windows: setx available (built-in)
- Unix: shell profile (~/.bashrc or ~/.zshrc) writable

## Examples

### Set and retrieve a variable

    $ gitmap env set GOPATH "/home/user/go"
      Set GOPATH=/home/user/go

    $ gitmap env get GOPATH
      GOPATH=/home/user/go

    $ gitmap env list
      Managed variables:
        GOPATH = /home/user/go

### Add a directory to PATH

    $ gitmap ev path add /usr/local/go/bin
      Added to PATH: /usr/local/go/bin

    $ gitmap env path list
      Managed PATH entries:
        /usr/local/go/bin

### Preview changes with dry-run

    $ gitmap env set NODE_ENV "production" --dry-run
      [dry-run] Would set NODE_ENV=production

    $ gitmap env path add /opt/tools/bin --dry-run
      [dry-run] Would add to PATH: /opt/tools/bin

### Target a specific shell profile with --shell

On Unix, `--shell` selects which profile file to write to. Without it,
gitmap auto-detects from `$SHELL`. Use it when you need to target a
profile other than your current login shell (e.g. seeding a CI image
or scripting from one shell while configuring another).

    $ gitmap env set EDITOR "nvim" --shell zsh
      Set EDITOR=nvim in ~/.zshrc

    $ gitmap env set EDITOR "nvim" --shell bash
      Set EDITOR=nvim in ~/.bashrc

    $ gitmap env path add /opt/go/bin --shell zsh
      Added to PATH (~/.zshrc): /opt/go/bin

    $ gitmap env delete EDITOR --shell bash --dry-run
      [dry-run] Would remove EDITOR from ~/.bashrc

Notes:
- `--shell` is Unix-only; on Windows it is silently ignored (Windows uses `setx`).
- Accepted values: `bash`, `zsh`. Other shells fall back to auto-detect.
- Combine with `--dry-run` to preview which profile file would change.

## See Also

- [install](install.md) — Install developer tools
- [doctor](doctor.md) — Diagnose PATH and version issues
- [setup](setup.md) — Configure Git global settings

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter env
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
