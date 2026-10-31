# gitmap thelp specific

Command-specific terminal help query system.

## Usage

```bash
gitmap thelp <command>           # Terminal help for specific command
gitmap help <command>            # Alias for specific command help
gitmap thelp merge               # Dedicated merge & transfer suite help
gitmap thelp merge-right         # One-way merge into RIGHT endpoint
gitmap thelp merge-left          # One-way merge into LEFT endpoint
gitmap thelp merge-in            # Multi-source chronological commit replay
gitmap thelp commit-in           # Replay commit history into target repository
gitmap thelp login               # Authentication, token, SSH & AI Auto setup
gitmap thelp settings            # Speed settings and AI Auto configuration
```

## Key Topics & Examples

| Topic | Description | Example Invocation |
|---|---|---|
| `merge` | Full merge & transfer suite | `gitmap thelp merge` |
| `merge-right` | One-way merge to RIGHT | `gitmap merge-right ./src https://github.com/org/repo -y` |
| `merge-left` | One-way merge to LEFT | `gitmap merge-left ./dst https://github.com/org/upstream -y` |
| `merge-in` | Chronological commit-in | `gitmap merge-in ./target all --default` |
| `commit-in` | Replay commit history | `gitmap cin ./target ./fork1,./fork2` |
| `login` | Provider auth & AI bot | `gitmap login set ai-auto gemini` |
| `settings` | Unified speed settings | `gitmap settings set ai_auto.default_agent antigravity` |

## Flags

| Flag | Description |
|---|---|
| `--pretty` | Force ANSI colors in pagers or redirects |
| `--no-pretty` | Disable ANSI formatting for plaintext files |
| `--json` | Return structured JSON metadata where supported |

## Examples

```bash
# Query specific help for merge suite
gitmap thelp merge

# Query specific help for commit replay
gitmap thelp commit-in

# Query specific help for AI automation login
gitmap thelp login
```
