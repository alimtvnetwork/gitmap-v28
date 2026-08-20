# gitmap sc

Short alias for `servers-clients`. Broadcasts commands across all server and client nodes in the cluster.

## Usage

```bash
gitmap sc <subcommand> [args] [flags]
```

## Examples

```bash
gitmap sc ps "Get-Date"
gitmap sc cmd "whoami" --except 2
gitmap sc pull --all
```

See also: `gitmap servers-clients`, `gitmap clients`, `gitmap cluster`
