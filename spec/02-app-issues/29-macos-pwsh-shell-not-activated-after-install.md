# macOS / Unix Install Does Not Activate the User's Active Shell (PowerShell-on-Unix)

> **Status:** Open — fix shipped in **v3.43.1** (2026-04-21)
> **Reporter:** end-user screenshot, macOS, pwsh inside Terminal/iTerm
> **Related specs:**
> - `spec/02-app-issues/22-installer-path-not-active-after-install.md` — multi-profile contract (this issue extends it)
> - `spec/04-generic-cli/21-post-install-shell-activation/02-snippets.md` — per-shell snippet bodies
> - `gitmap-v27/scripts/install.sh`, `install-quick.sh`

---

## 1. Symptoms

```text
PS /Users/ab_mahin/Downloads/core-v8> curl -fsSL https://raw.githubusercontent.com/.../install-quick.sh | bash
Install path: /Users/ab_mahin/.local/bin/gitmap-cli
  Shell: zsh
  PATH target: /Users/ab_mahin/.zshrc (added)
  Reload: . /Users/ab_mahin/.zshrc

OK To start using gitmap-v27 right now, run:

    . /Users/ab_mahin/.zshrc

  Or open a new terminal window.

Installed to:
gitmap-v27 quick installer
---------------------
Choose install folder. Press Enter to accept the default.
Default: /Users/ab_mahin/.local/bin
Install path: /Users/ab_mahin/.local/bin/gitmap-cli/gitmap-v27   ← prompt fires AGAIN
App folder on PATH: gitmap-v27 quick installer                   ← garbage label
---------------------
Choose install folder. Press Enter to accept the default.
Default: /Users/ab_mahin/.local/bin
Install path: /Users/ab_mahin/.local/bin/gitmap-cli           ← prompt fires a THIRD time

Done! Run 'gitmap-v27 --help' to get started.

PS /Users/ab_mahin/Downloads/core-v8> gitmap-v27 r v1.5.6
gitmap-v27: The term 'gitmap-v27' is not recognized as a name of a cmdlet, function, ...
PS /Users/ab_mahin/Downloads/core-v8> gitmap-v27
gitmap-v27: The term 'gitmap-v27' is not recognized as a name of a cmdlet, function, ...
PS /Users/ab_mahin/Downloads/core-v8> gitmap-v27 --help
gitmap-v27: The term 'gitmap-v27' is not recognized as a name of a cmdlet, function, ...
```

Two distinct failures:

1. **Wrong-shell PATH wiring** — install.sh wrote PATH to `.zshrc` /
   `.zprofile`, but the user's **interactive shell is PowerShell 7+
   (`pwsh`) on macOS**. Neither file is sourced by pwsh. Result:
   `gitmap-v27` is not on PATH in the active session OR in any future pwsh
   session.
2. **Prompt fires three times** — the install-quick.sh discovery
   delegation re-runs the script in a child bash, which prompts; the
   delegated install.sh prompts; and the outer install-quick.sh also
   falls through and prompts. The screenshot shows the user seeing
   "Choose install folder" three separate times.

---

## 2. Root Cause

### 2.1 Primary — pwsh on Unix is invisible to install.sh

`gitmap-v27/scripts/install.sh::add_to_path` detects the active shell with:

```bash
shell_name="$(basename "${SHELL:-/bin/bash}")"
```

On macOS, `$SHELL` reflects the **account's default login shell**
(zsh by default since Catalina) — it does NOT reflect the shell the
user is currently typing into. When the user opens pwsh from inside
Terminal.app, `$SHELL` stays `/bin/zsh` while the active interpreter
is `/usr/local/bin/pwsh`.

Consequences:

- The script branches into the zsh arm only, writing to `.zshrc` /
  `.zprofile` / `.profile`. None of these are sourced by pwsh.
- pwsh's profile lives at `~/.config/powershell/Microsoft.PowerShell_profile.ps1`
  on Unix (analogous to `$PROFILE` on Windows). Install.sh has zero
  branches that touch this path.
- The post-install "Reload" hint prints `. /Users/.../.zshrc`. If a
  pwsh user copy-pastes that, pwsh dot-sources a zsh script and
  errors out (or worse, partially succeeds with garbage).

The `print-path-snippet` subcommand DOES support `--shell pwsh` (see
`gitmap-v27/constants/constants_pathsnippet.go::PathSnippetPwshFmt`), but
install.sh never invokes that branch — it hardcodes
`snippet_shell="bash"` for everything that isn't fish.

### 2.2 Secondary — discovery delegation leaks into baseline path

`install-quick.sh` flow:

```bash
EFFECTIVE_REPO=$(resolve_effective_repo ...)
if [ "$EFFECTIVE_REPO" != "$REPO" ]; then
    invoke_delegated_installer "$EFFECTIVE_REPO" || true   # ← guarded
fi
# ... falls through to baseline prompt + install ...
```

`invoke_delegated_installer` ends with `exit $?`, which only runs if
the inline `bash -c "$script" _ ...` returns. Under `curl | bash`
piping with `read < /dev/tty` inside the child, the child's `read`
sometimes returns success with empty input (depending on tty
ownership), causing the delegated script to **complete the install**
(printing its own "Installed to:" + "Done!"), then `exit 0`. But
because `set -e` is active and the shell sees the delegated child
exit cleanly, control returns to the outer `|| true` arm and the
outer script ALSO runs its prompt + install pass. The user sees the
prompt 2-3 times and a duplicate install.

The "garbage label" `App folder on PATH: gitmap-v27 quick installer` is
the same issue rendered visually: the second installer's banner
overlapped with the first installer's `App folder on PATH:` line.

### 2.3 Why issue 22's fix didn't catch this

Issue 22 added multi-profile writes for **bash + zsh + fish + POSIX
sh**. PowerShell-on-Unix was never enumerated in the matrix because
the assumption at the time was "PowerShell users are on Windows and
get install.ps1". macOS Homebrew adoption of `pwsh` (and Linux
Microsoft repos) made that assumption wrong.

---

## 3. Fix

### 3.1 install.sh — detect and write pwsh profile on Unix

1. New helper `detect_active_pwsh()`:
   - Returns 0 if `$PSModulePath` is set in the current env (always
     set inside a pwsh session, even when invoked through a sh
     wrapper), OR if `command -v pwsh` succeeds AND the user opted
     in via `--with-pwsh`.
2. New helper `pwsh_profile_path()`:
   - Linux/macOS: `${HOME}/.config/powershell/Microsoft.PowerShell_profile.ps1`
   - Creates the parent directory if missing.
3. Extend `add_to_path` with a pwsh arm that calls
   `add_path_to_profile_pwsh "${dir}" "${profile}"` — a new sibling
   to `add_path_to_profile` that:
   - Renders the snippet via `gitmap-v27 setup print-path-snippet
     --shell pwsh --dir <dir> --manager installer`.
   - Uses the same `# gitmap-v27 shell wrapper v2 ...` marker block, so
     idempotent rewrite via the existing awk pass works unchanged.
4. When pwsh is the active shell:
   - `PATH_TARGET` becomes the pwsh profile path.
   - `PATH_RELOAD` becomes `. $PROFILE` (pwsh syntax).
   - Post-install banner prints the pwsh-flavored reload line.

### 3.2 install-quick.sh — make delegation terminal

Replace the silent `|| true` with an explicit success/failure split:

```bash
if [ "$EFFECTIVE_REPO" != "$REPO" ]; then
    if invoke_delegated_installer "$EFFECTIVE_REPO"; then
        exit 0   # delegated child already installed; do not fall through
    fi
    printf '  [discovery] [WARN] delegation failed; falling back to baseline\n' >&2
fi
```

`invoke_delegated_installer` already calls `exit $?`, so the only way
control returns to the parent is if `curl` itself fails. The new
explicit `exit 0` gate prevents the double-install.

### 3.3 Banner accuracy

When `PATH_SHELL == pwsh`, the post-install summary prints:

```text
  Shell: pwsh
  PATH target: /Users/<u>/.config/powershell/Microsoft.PowerShell_profile.ps1 (added)
  Reload: . $PROFILE
```

instead of the bash/zsh-flavored reload command.

---

## 4. Verification

Manual (macOS, pwsh open in Terminal.app):

```pwsh
PS> curl -fsSL https://raw.githubusercontent.com/.../install-quick.sh | bash
# expect: prompt fires ONCE, summary shows "Shell: pwsh", reload line is `. $PROFILE`
PS> . $PROFILE
PS> gitmap-v27 --version
# expect: prints v3.43.1
```

Automated regression (CI, future work):

- New test in `gitmap-v27/scripts/install_test.sh` (does not exist yet)
  that runs `bash install.sh --dir /tmp/x --no-discovery` with
  `PSModulePath=fake` exported and asserts the pwsh profile got the
  marker block.

---

## 5. History

| Version | Change |
|---|---|
| v3.43.1 | Issue filed + fix implemented. install.sh detects pwsh-on-Unix via `$PSModulePath` and writes the pwsh profile. install-quick.sh exits cleanly after successful delegation. Post-install banner emits pwsh-flavored reload line when applicable. |
