# `gitmap-v27 cd` prints the path but does not change the current shell directory

## Ticket

Users run `gitmap-v27 cd <repo>` in an interactive shell and expect the terminal
itself to move into that repository, but the command only prints the absolute
path.

## Symptoms

1. User runs `gitmap-v27 cd gitmap-v27`.
2. The command resolves the repo from the SQLite database correctly.
3. The terminal prints the absolute path.
4. The prompt stays in the original working directory.

## Root Cause

`gitmap-v27 cd` is implemented as a normal CLI subcommand. A child process can
print the destination path, but it cannot mutate the working directory of the
parent shell process that launched it.

The original implementation documented this limitation and installed only a
separate `gcd` shell helper during `gitmap-v27 setup`. That left a behavioral gap:

- `gitmap-v27 cd <repo>` still executed the real binary directly.
- The binary correctly printed the destination path.
- No shell-level wrapper intercepted that output and called `cd` /
  `Set-Location` in the parent shell.

So the lookup logic was working; the shell integration contract was incomplete.

## Fix

Upgrade the setup-installed shell integration so supported shells receive a
managed wrapper for **both** `gitmap-v27` and `gcd`:

1. Install a shell function named `gitmap-v27`.
2. Intercept `gitmap-v27 cd ...` and `gitmap-v27 go ...`.
3. Delegate to the real executable (`command gitmap-v27` on Bash/Zsh,
   `gitmap.exe`/application lookup on PowerShell).
4. Capture the emitted path.
5. Change the parent shell directory with `cd` / `Set-Location`.
6. Pass all non-`cd` commands through unchanged.

The wrapper uses a new managed marker so rerunning `gitmap-v27 setup` appends the
new shell integration even when older `gcd`-only profile entries already exist.

## Prevention

1. Any CLI feature that must affect the current shell session must ship with a
   shell wrapper, not only a child-process subcommand.
2. The direct binary behavior and the shell-integrated behavior must be treated
   as separate contracts and tested separately.
3. Profile-managed shell snippets must use versioned markers so upgrades can be
   rolled out safely without requiring manual profile cleanup.

## Related

- `spec/01-app/31-cd.md`
- `gitmap-v27/completion/cdfunction.go`
- `gitmap-v27/constants/constants_cd.go`