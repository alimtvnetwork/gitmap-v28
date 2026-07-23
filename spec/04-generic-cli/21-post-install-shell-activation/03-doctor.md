# Post-Install Shell Activation — Doctor Check

> **Parent spec:** [../21-post-install-shell-activation.md](../21-post-install-shell-activation.md)
> **Sibling files:**
> - [01-contract.md](01-contract.md) — Required behaviours and activation flow
> - [02-snippets.md](02-snippets.md) — Per-shell profile snippet bodies
> - [04-idempotency.md](04-idempotency.md) — Rewrite, removal, and version-bump rules

## `doctor` Wrapper Check

`doctor` MUST emit one of these three outcomes:

| Status | Stdout | Exit |
|--------|--------|------|
| LOADED | `[OK] Shell wrapper active (TOOLNAME_WRAPPER=1)` | 0 |
| INSTALLED_BUT_NOT_LOADED | `[!!] Shell wrapper installed but not loaded — run: source ~/.zshrc` | 1 |
| NOT_INSTALLED | `[!!] Shell wrapper missing — run: toolname setup` | 1 |

---

## Detection Algorithm

1. Read `<TOOL>_WRAPPER` from environment → if `"1"` → **LOADED**.
2. Read profile file (`$PROFILE` / `~/.bashrc` / `~/.zshrc` /
   `~/.config/fish/config.fish`) and search for the marker line
   `# <tool> shell wrapper v<N>`. If found → **INSTALLED_BUT_NOT_LOADED**.
3. Otherwise → **NOT_INSTALLED**.

The reload one-liner printed for status 2 MUST match the detected
shell from [01-contract.md](01-contract.md#shell-detection-rules).

---

## Implementation Checklist For New CLIs

1. Add `<TOOL>_WRAPPER` constant to the constants package.
2. Define `ShellWrapperMarkerPrefix` and `ShellWrapperMarkerSuffix`
   constants (e.g. `# toolname shell wrapper v2`).
3. Implement `setup/wrapper.go` with:
   - `DetectShell()`
   - `ResolveProfilePath(shell)`
   - `InjectSnippet(profilePath, shell)` (idempotent via marker)
   - `RemoveSnippet(profilePath)` for upgrades
   - `TryInSessionActivate(shell)` (PowerShell only succeeds today)
   - `PrintReloadInstruction(shell)`
4. Add `doctor` check `checkShellWrapper()` returning one of the three
   statuses above.
5. Add stderr warnings to every shell-dependent subcommand using
   `os.Getenv("<TOOL>_WRAPPER")`.
6. Cover with tests (see [04-idempotency.md](04-idempotency.md#testing-requirements)).
