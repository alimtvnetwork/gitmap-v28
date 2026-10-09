# Deprecation Policy

How gitmap retires commands, subcommands, flags, and aliases without breaking users.

## Rules

1. **Two-release warning window.** Nothing is removed or repurposed until it has printed a deprecation warning for at least two minor releases.
2. **Aliases never break.** An alias, once shipped, keeps working. If its target is renamed, the alias follows the new target and prints the warning.
3. **Warning format.** One line on stderr, prefixed with two-space + middle-dot, naming the replacement:
   `  · <old> is deprecated; use <new> instead`
   (mirrors `MsgScanWorkersConcurrencyAlias` in `cli/constants/constants_cli.go`).
4. **Warning at use time.** The warning prints when the deprecated thing is invoked — not in docs alone. Docs are not enforcement.
5. **Removal checklist.** Before removing: (a) warning shipped for 2+ releases, (b) `docs/deprecation-policy.md` removal log updated, (c) major-version or explicit owner sign-off for commands (flags may go in a minor).

## Removal log

| Removed | Replaced by | Warning since | Removed in |
|---|---|---|---|
| `gitmap fix <git-state>` (stash/wip/discard surface) | `gitmap stash` / `gitmap wip` / `gitmap discard` | repurposed in v6.520.0 without warning; warning added after v6.522.0 | not yet removed — warning phase |

## Notes

- The `fix` command was repurposed in v6.522.0 (content fixer). The old git-state surface now warns instead of silently failing.
- Flag deprecations follow the same pattern (see `cli/cmdscan/flags.go:152`).
