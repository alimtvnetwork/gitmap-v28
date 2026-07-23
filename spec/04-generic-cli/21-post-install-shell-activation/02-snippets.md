# Post-Install Shell Activation — Snippets

> **Parent spec:** [../21-post-install-shell-activation.md](../21-post-install-shell-activation.md)
> **Sibling files:**
> - [01-contract.md](01-contract.md) — Required behaviours and activation flow
> - [03-doctor.md](03-doctor.md) — `doctor` wrapper status detection
> - [04-idempotency.md](04-idempotency.md) — Rewrite, removal, and version-bump rules

## Profile Snippet Contract

Every snippet MUST start with the marker line and MUST end with a
matching closing marker so the CLI can rewrite or remove it without
disturbing surrounding content:

```
# <tool> shell wrapper v2 — managed by `<tool> setup`. Do not edit manually.
...snippet body...
# <tool> shell wrapper v2 end
```

The detection variable name MUST follow `<TOOL>_WRAPPER` (uppercased,
underscores) so multiple CLIs can coexist in one profile.

---

## PowerShell (`$PROFILE`)

```powershell
# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.
$env:TOOLNAME_WRAPPER = "1"
function gcd { Set-Location (toolname cd @args) }
# toolname shell wrapper v2 end
```

Profile path resolution: `$PROFILE` (typically
`Documents\PowerShell\Microsoft.PowerShell_profile.ps1` on PS 7+ or
`Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1` on PS 5).

---

## Bash (`~/.bashrc`)

```bash
# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.
export TOOLNAME_WRAPPER=1
gcd() { cd "$(toolname cd "$@")" ; }
# toolname shell wrapper v2 end
```

---

## Zsh (`~/.zshrc`)

```bash
# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.
export TOOLNAME_WRAPPER=1
gcd() { cd "$(toolname cd "$@")" ; }
# toolname shell wrapper v2 end
```

---

## Fish (`~/.config/fish/config.fish`)

```fish
# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.
set -gx TOOLNAME_WRAPPER 1
function gcd; cd (toolname cd $argv); end
# toolname shell wrapper v2 end
```

---

## Cross-Platform Parity Table

| Capability | PowerShell | Bash | Zsh | Fish |
|------------|-----------|------|-----|------|
| Profile detection | ✅ `$PROFILE` | ✅ `~/.bashrc` | ✅ `~/.zshrc` | ✅ `~/.config/fish/config.fish` |
| Marker-based snippet | ✅ | ✅ | ✅ | ✅ |
| Wrapper detection env var | ✅ `$env:` | ✅ `export` | ✅ `export` | ✅ `set -gx` |
| In-session activation | ✅ dot-source | ❌ (print one-liner) | ❌ (print one-liner) | ❌ (print one-liner) |
| `doctor` LOADED check | ✅ | ✅ | ✅ | ✅ |
| `doctor` INSTALLED_BUT_NOT_LOADED | ✅ | ✅ | ✅ | ✅ |
| Reload one-liner printed | ✅ `. $PROFILE` | ✅ `source ~/.bashrc` | ✅ `source ~/.zshrc` | ✅ `source ~/.config/fish/config.fish` |
