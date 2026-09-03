# `gitmap create`

Create and initialize a local repository and automatically provision and push to GitHub or GitLab.

## Simulation

```
$ gitmap create my-new-service --description "Microservice backend"
  ✓ Repository created successfully!
  ● Name:      my-new-service
  ● Path:      D:/wp-work/my-new-service
  ● Profile:   alimtvnetwork (github)
  ● Remote:    https://github.com/alimtvnetwork/my-new-service
```

## Subcommands & Usage

```
gitmap create [repo] <name> [flags]
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--profile <name\|1-N>` | Target GitHub/GitLab user or organization profile | Active default |
| `--org <name>` | Explicit organization owner | Active default |
| `--private` | Create repository as private | `true` |
| `--public` | Create repository as public | `false` |
| `--dir <path>` | Destination local directory | `./<name>` |
| `-d`, `--description <text>` | Repository description | `""` |
| `--no-remote` | Local git repository initialization only | `false` |
| `--json` | Output creation metadata as structured JSON | `false` |

## Examples

```bash

# Create private repository using default Git profile

gitmap create my-api

# Create repository explicitly under an organization

gitmap create web-client --org riseup-asia --public

# Create local-only repository without pushing to remote

gitmap create scratch-tool --no-remote

# Structured JSON output for CI/CD automation

gitmap create pipeline-tester --json
```
