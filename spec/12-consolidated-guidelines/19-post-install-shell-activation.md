# Post-Install Shell Activation

> **Source:** Consolidated from the split spec under [spec/04-generic-cli/21-post-install-shell-activation/](../04-generic-cli/21-post-install-shell-activation/):
> - [01-contract.md](../04-generic-cli/21-post-install-shell-activation/01-contract.md) — Required behaviours, activation flow, shell detection, stderr warnings.
> - [02-snippets.md](../04-generic-cli/21-post-install-shell-activation/02-snippets.md) — Per-shell profile snippet bodies and parity table.
> - [03-doctor.md](../04-generic-cli/21-post-install-shell-activation/03-doctor.md) — `doctor` three-state detection and implementation checklist.
> - [04-idempotency.md](../04-generic-cli/21-post-install-shell-activation/04-idempotency.md) — Rewrite, removal, version-bump, and testing rules.
> - Index: [spec/04-generic-cli/21-post-install-shell-activation.md](../04-generic-cli/21-post-install-shell-activation.md)
>
> **Applies to:** Every CLI in this framework that ships shell-integrated subcommands (e.g., `cd`, `clone-next`, `go`).

## Purpose

After `setup` (or the bootstrap installer) runs, the user MUST be able to invoke the CLI **and any of its shell-integrated subcommands** in the **current terminal session** without restarting it. This guideline ensures every new CLI inherits a deterministic, AI-implementable activation flow on day one.

---

## Required Behaviours

| ID | Behaviour |
|----|-----------|
| PIA-1 | `setup` writes a profile snippet, idempotent via marker comment. |
| PIA-2 | `setup` exports a `<TOOL>_WRAPPER=1` env var so the binary can detect wrapper state. |
| PIA-3 | `setup` attempts in-process activation (PowerShell dot-source `$PROFILE`). |
| PIA-4 | When in-session activation cannot run, `setup` prints the **exact** reload one-liner. |
| PIA-5 | `doctor` reports one of `LOADED`, `INSTALLED_BUT_NOT_LOADED`, `NOT_INSTALLED`. |
| PIA-6 | Shell-dependent subcommands print a stderr warning when `<TOOL>_WRAPPER` is unset. |
| PIA-7 | Profile snippet uses versioned markers: `# <tool> shell wrapper v<N>` … `# <tool> shell wrapper v<N> end`. |

---

## Profile Snippet Contract

Every snippet MUST start and end with matching marker lines so the CLI can rewrite or remove it without disturbing surrounding content.

### PowerShell (`$PROFILE`)
```powershell
# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.
$env:TOOLNAME_WRAPPER = "1"
function gcd { Set-Location (toolname cd @args) }
# toolname shell wrapper v2 end
```

### Bash / Zsh
```bash
# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.
export TOOLNAME_WRAPPER=1
gcd() { cd "$(toolname cd "$@")" ; }
# toolname shell wrapper v2 end
```

### Fish
```fish
# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.
set -gx TOOLNAME_WRAPPER 1
function gcd; cd (toolname cd $argv); end
# toolname shell wrapper v2 end
```

The detection variable name MUST follow `<TOOL>_WRAPPER` (uppercased, underscores) so multiple CLIs can coexist in one profile.

---

## In-Session Activation

`setup` MUST attempt to activate the wrapper in the **current** shell before falling back to a printed instruction.

- **PowerShell**: dot-source `$PROFILE` in-process; verify `$env:<TOOL>_WRAPPER -eq "1"`.
- **Bash / Zsh / Fish**: a child process cannot source the parent's profile. Detect the active shell via `$SHELL` and print the exact one-liner (`source ~/.bashrc`, `source ~/.zshrc`, `source ~/.config/fish/config.fish`).

Always end with: *"Or open a new terminal window."*

---

## `doctor` Wrapper Check

| Status | Stdout | Exit |
|--------|--------|------|
| LOADED | `[OK] Shell wrapper active (TOOLNAME_WRAPPER=1)` | 0 |
| INSTALLED_BUT_NOT_LOADED | `[!!] Shell wrapper installed but not loaded — run: source ~/.zshrc` | 1 |
| NOT_INSTALLED | `[!!] Shell wrapper missing — run: toolname setup` | 1 |

Detection algorithm:
1. `<TOOL>_WRAPPER == "1"` → **LOADED**.
2. Profile contains `# <tool> shell wrapper v<N>` marker → **INSTALLED_BUT_NOT_LOADED**.
3. Otherwise → **NOT_INSTALLED**.

---

## Stderr Warning From Shell-Dependent Subcommands

Any subcommand that requires the wrapper MUST detect missing state and print:

```
  ⚠ Shell wrapper not active. The current command will print the path
    instead of changing directory. Run `toolname setup` (and reload
    your shell) to enable shell-integrated behaviour.
```

The warning MUST include both the setup action and the reload step.

---

## Idempotency Rules

- Rewrites MUST be safe to run repeatedly.
- The marker comment is the **only** legal anchor for rewrites — never line counts or content matching.
- Bumping the version (`v2` → `v3`) MUST remove every previous marker block before injecting the new one.
- Uninstall MUST delete the marker block in full and leave surrounding content byte-identical.

---

## Implementation Checklist (per CLI)

1. Add `<TOOL>_WRAPPER` constant to the constants package.
2. Define `ShellWrapperMarkerPrefix` / `ShellWrapperMarkerSuffix` constants.
3. Implement `setup/wrapper.go`:
   - `DetectShell()` / `ResolveProfilePath(shell)`
   - `InjectSnippet()` / `RemoveSnippet()` (marker-anchored, idempotent)
   - `TryInSessionActivate(shell)` (PowerShell only succeeds today)
   - `PrintReloadInstruction(shell)`
4. Add `doctor` check `checkShellWrapper()` returning the three statuses.
5. Add stderr warnings in every shell-dependent subcommand.
6. Cover with tests: fresh injection, re-injection (no duplicates), marker-based removal, all three `doctor` states.

---

## Constraints

- Snippets MUST be **ASCII only** — no em-dashes or Unicode arrows (PowerShell parser fails on UTF-8 in some hosts).
- The CLI MUST NOT modify any line outside its marker block.
- The CLI MUST NOT depend on the user editing their profile manually.
- Snippet body MUST be under 10 lines so users can audit it.
- Reload instructions MUST be a single copy-pasteable line.

---

## Cross-Platform Parity

| Capability | PowerShell | Bash | Zsh | Fish |
|------------|-----------|------|-----|------|
| Profile detection | `$PROFILE` | `~/.bashrc` | `~/.zshrc` | `~/.config/fish/config.fish` |
| In-session activation | ✅ dot-source | ❌ (print one-liner) | ❌ (print one-liner) | ❌ (print one-liner) |
| Reload one-liner | `. $PROFILE` | `source ~/.bashrc` | `source ~/.zshrc` | `source ~/.config/fish/config.fish` |

See the split sub-files linked at the top — [01-contract](../04-generic-cli/21-post-install-shell-activation/01-contract.md), [02-snippets](../04-generic-cli/21-post-install-shell-activation/02-snippets.md), [03-doctor](../04-generic-cli/21-post-install-shell-activation/03-doctor.md), [04-idempotency](../04-generic-cli/21-post-install-shell-activation/04-idempotency.md) — for full rationale, activation flow diagram, and historical bug references.
