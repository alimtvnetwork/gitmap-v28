# gitmap push-fix

Autonomous diagnosis, self-healing, and direct GitHub push recovery.

## Usage

```bash
gitmap push-fix [flags]
gitmap push fix [flags]
gitmap pfx [flags]
```

## Description

`gitmap push-fix` performs an end-to-end diagnosis and automated remediation of failed git push operations. It identifies and resolves root causes:

1. **Pre-flight Git State Inspection**: Verifies repository root, detects active branch, unpushed commits, and remote tracking state.
2. **Remote Authentication Probing**: Executes non-interactive SSH public key verification and transport detection under bounded 5-second timeouts with anti-hang safeguards (`GIT_TERMINAL_PROMPT=0`, `GCM_INTERACTIVE=never`).
3. **Self-Healing & Key Realignment**:
   - Detects authorized key candidates (e.g. `id_rsa_-y` or rotated keys) and promotes them to standard `~/.ssh/id_rsa`.
   - Purges transient or dashed keys from GitMap's SQLite registry.
   - Synchronizes `~/.ssh/config` and `%SystemDrive%\Users\%USERNAME%\.ssh\config` with canonical `Host github.com`.
   - Automatically converts blocked or unauthenticated HTTPS remotes to verified SSH (`git@github.com:...`).
4. **Push Execution & Self-Healing**:
   - Auto-binds new branches via `-u` (`--set-upstream`).
   - Automatically handles non-fast-forward diverged branches by running `git pull --rebase` and retrying.
   - Pushes directly to GitHub with clear terminal telemetry cards.

## Flags

| Flag | Short | Type | Description |
|---|---|---|---|
| `--dry-run` | `-n` | bool | Simulate diagnosis and push without altering remote |
| `--remote` | `-r` | string | Target remote name (default: `origin`) |
| `--branch` | `-b` | string | Target branch name (default: current branch) |
| `--force` | `-f` | bool | Push with safe `--force-with-lease` |
| `--yes` | `-y` | bool | Non-interactive mode (bypass prompts) |
| `--ssh` | | bool | Force conversion of remote to SSH |
| `--https` | | bool | Force conversion of remote to HTTPS |

## Examples

```bash
# Diagnose and fix push authentication in current repository
gitmap push-fix

# Dry-run test without pushing
gitmap push-fix --dry-run

# Push fix via the push subcommand
gitmap push fix

# Target specific remote and branch
gitmap push-fix --remote origin --branch feature/login

# Push with force-with-lease after conflict resolution
gitmap push-fix --force
```
