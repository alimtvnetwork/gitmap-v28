# `gitmap create`

Create and initialize a local repository, automatically slugify human-readable names with spaces, and automatically provision and push to GitHub.

## Subcommands & Usage

```
gitmap create [repo] <name> [folder] [slug] [flags]
gitmap repo-create (repoc) <name> [folder] [slug] [flags]
gitmap create-repo (crepo) <name> [folder] [slug] [flags]
gitmap create-local-repo (clr) <name> [folder] [slug] [flags]
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

gitmap create web-client --org my-org --public

# Create local-only repository without pushing to remote

gitmap create scratch-tool --no-remote

# Structured JSON output for CI/CD automation

gitmap create pipeline-tester --json
```
