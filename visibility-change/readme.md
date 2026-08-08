# visibility-change scripts

Script-only fallback for the Go-native `gitmap make-public` /
`gitmap make-private` commands. Use these when you can't ship the
`gitmap` binary (CI runners, ad-hoc remotes, etc.).

## Layout

```
scripts/
  visibility-change.ps1          # PowerShell entry point
  visibility-change.sh           # Bash entry point
  visibility-change/
    Provider.ps1 / provider.sh   # origin URL parsing, gh/glab detection
    Apply.ps1    / apply.sh      # apply + verify + confirm prompt
```

## Usage

PowerShell:

```powershell
./scripts/visibility-change.ps1                  # toggle
./scripts/visibility-change.ps1 -Visible pub     # force public (prompts)
./scripts/visibility-change.ps1 -Visible pri     # force private
./scripts/visibility-change.ps1 -Yes -DryRun     # preview, no API call
```

Bash:

```bash
./scripts/visibility-change.sh                   # toggle
./scripts/visibility-change.sh --visible pub     # force public (prompts)
./scripts/visibility-change.sh --visible pri     # force private
./scripts/visibility-change.sh --yes --dry-run   # preview, no API call
```

## Requirements

- `gh` (GitHub) or `glab` (GitLab) on `PATH`, already authenticated.
- `git` with an `origin` remote pointing at github.com or gitlab.com.
- For self-hosted GitLab, set `VISIBILITY_GITLAB_HOSTS=host1,host2`.

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success (or already at target visibility) |
| 2 | Not inside a git repository |
| 3 | No `origin` remote |
| 4 | Unsupported provider host / unparseable owner/repo |
| 5 | Provider CLI missing, not authenticated, or apply failed |
| 6 | Bad flag / argument |
| 7 | Confirmation required (re-run with `--yes` / `-Yes`) |
| 8 | Verification failed (visibility did not change) |

These match the exit codes emitted by `gitmap make-public` /
`gitmap make-private` so wrappers can branch on the same numbers.
