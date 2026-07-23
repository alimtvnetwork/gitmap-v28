# 108 — `install-quick.sh` Auto-Source Wrapper

**Status:** Implemented (2026-04-22)
**Companion:** `install-quick.sh`, `gitmap/scripts/install.sh`, spec 95 (versioned repo discovery)

## Problem

When users installed gitmap via the README one-liner:

```
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install-quick.sh | bash
```

…the install succeeded but `gitmap` was not callable on the very next prompt.
The user had to **manually run `source ~/.zshrc` (or open a new terminal)**
before PATH picked up the new binary.

## Why this happened

This is a hard POSIX constraint, not a bug:

1. `curl ... | bash` spawns a **child** `bash` process.
2. The child can write `export PATH=...` to `~/.zshrc` (which `install.sh`
   does, via its `add_to_path` helper).
3. The child can even run `source ~/.zshrc` itself — but that source mutates
   only **the child's** environment.
4. The child exits. Its environment is destroyed. The user's interactive
   shell never observed the change.

A child process **cannot mutate its parent shell's environment**. Period.

## Solution: dual-mode execution

`install-quick.sh` now supports two installation modes:

### Eval-mode (recommended, auto-activating)

```
eval "$(curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install-quick.sh)"
```

Because `eval` runs the script's source text **inside the user's
interactive shell**, the trailing `. ~/.zshrc` mutates the *live* shell.
After the one-liner returns, `gitmap --help` works immediately on the next
prompt — no manual source, no new terminal.

The script detects eval-mode at startup (`__gitmap_detect_eval_mode`) by
inspecting `$0` and `BASH_SOURCE[0]`. When eval'd in zsh or bash, `$0` is
the user's shell (`zsh`, `-bash`, etc.) and `BASH_SOURCE[0]` is empty —
those signals trigger the auto-source path.

### Pipe-mode (legacy, prints manual hint)

```
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install-quick.sh | bash
```

Still works exactly as before. The script detects it's running in a child
bash and — instead of pretending it can source — prints a high-contrast
banner with the exact `source <profile>` command to run, plus a tip
suggesting eval-mode for next time.

## Implementation details

- **All heavy work runs inside a subshell** (`__gitmap_quick_install_main`)
  with `set -eu`. This isolates the installer's strict-mode options from
  the user's interactive shell, which would otherwise inherit them in
  eval-mode and break their session.
- **No top-level `set -euo pipefail`** for the same reason. Strict mode
  applies only inside the subshell.
- **PATH_RELOAD detection.** `install.sh` already prints a `source <file>`
  line as part of its post-install summary. The wrapper captures the
  installer's stderr via `tee`, greps the first `source ` or `. ` token,
  and writes the resolved profile to `<install_dir>/.gitmap-last-profile`
  so the outer driver can `.`-source it after the subshell exits.
- **Fallback profile.** If the hint can't be parsed, the wrapper picks
  `~/.zshrc`, `~/.bashrc`, or `~/.profile` based on `$SHELL`.
- **Verification.** After sourcing, the wrapper runs
  `command -v gitmap` and reports the resolved binary path. If it's still
  missing, it falls back to the manual-hint banner.

## Backwards compatibility

- The local-file invocations (`./install-quick.sh`, `./install-quick.sh
  --dir /opt/gitmap`, `--no-discovery`, `--probe-ceiling N`) all work
  unchanged.
- The pipe-mode invocation works unchanged — only the post-install
  messaging is different (loud banner instead of buried hint).
- All flag parsing and versioned-repo discovery (spec 95) is preserved
  bit-for-bit; the only changes are (a) the entry point dispatcher,
  (b) capturing installer stderr to extract PATH_RELOAD, (c) the
  end-of-run source-or-hint branch.

## README update (recommended, not yet applied)

The README install snippet should be updated to lead with eval-mode:

````markdown
**Linux / macOS** (auto-activates PATH in current shell):

```bash
eval "$(curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install-quick.sh)"
```

Legacy pipe-mode (prints a manual `source` hint at the end):

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install-quick.sh | bash
```
````

## See also

- `gitmap/scripts/install.sh::add_to_path` — writes PATH entries to rc files
- `gitmap/scripts/install.sh::detect_active_pwsh` — picks the right reload command per shell
- `spec/01-app/95-installer-script-find-latest-repo.md` — versioned repo discovery
