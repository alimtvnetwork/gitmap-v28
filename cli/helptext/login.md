# gitmap login

Log in to GitHub so private-repository clone, pull, and push work without extra flags.

## Usage

```bash
gitmap login                  # Interactive: pick token paste or browser login
gitmap login --token <PAT>    # Log in with a personal access token (non-interactive)
gitmap login --web            # Log in via the browser
gitmap login --status         # Show current login state
gitmap logout                 # Remove the stored GitHub credential
```

## Flags

| Flag | Description |
|------|-------------|
| `--token <PAT>` | Personal access token (classic or fine-grained). Validated against api.github.com before storing. |
| `--web`, `--browser` | Browser login. Uses `gh auth login --web` when the GitHub CLI is installed, otherwise opens the token creation page and prompts for a secure paste. |
| `--status` | Show login state with the token masked. |
| `--no-verify` | With `--token`: store without validating (offline use). |
| `-h`, `--help` | Show this help. |

## How it works

1. The token is validated with `GET https://api.github.com/user` (unless `--no-verify`).
2. On success it is stored in git global config (`github.token`).
3. `gitmap clone`, `pull`, and `push` pick it up automatically via token resolution — no `--token` flag needed afterwards.

## Examples

```bash
gitmap login --token ghp_xxxxxxxxxxxxxxxxxxxx
gitmap login --web
gitmap login --status
```

## See also

- `gitmap token list` — show the resolved token and its source
- `gitmap logout` — remove the stored credential
