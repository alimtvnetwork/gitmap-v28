# 26 — `cfrp` on a non-versioned repo: hard error with `--require-version` opt-in

> **Status:** specified — implemented in v4.42.x via the `--require-version`
> flag in `gitmap-v28/cmd/clonefixrepo.go`.

## Problem

Running `gitmap-v28 clone-fix-repo-pub <url>` (alias `cfrp`) on a URL whose
repository name does **not** end in a `-v<N>` suffix used to silently
proceed: it cloned, ran `fix-repo --all` (which had nothing to rewrite
because there was no version token), then flipped visibility public.

Two failure modes followed from this:

1. **User confusion.** `cfrp` was conceived for the versioned-repo
   workflow (`gitmap-v28`, `coding-guidelines-v23`). Running it on a
   plain repo silently skipped the rewrite step, so users believed
   the rewrite ran when it hadn't.
2. **Accidental publication.** `cfrp` always flips visibility to
   public. Pairing that with a no-op rewrite is a foot-gun on private
   repos that were never intended to be public-by-default.

## Decision

`cfrp` MUST refuse to run on a non-versioned repo unless the caller
opts in. The opt-in is the existing `--require-version` flag, but its
default is now **strict** (`true`): the flag now means "I confirm I
am working on a versioned repo and want the rewrite to be required."

Concretely:

| Scenario | Behavior |
|---|---|
| URL ends in `-v<N>` (e.g. `gitmap-v28`) | Run the full pipeline. |
| URL has no `-v<N>` suffix, no flag | **Hard error**, non-zero exit, no clone. |
| URL has no `-v<N>` suffix, `--require-version=false` | Skip + warn, continue. |

The skip path still flips visibility public, so `--require-version=false`
is the explicit "I know what I'm doing" escape hatch.

## Error message

```
✗ cfrp: <url> is not a versioned repository (no -v<N> suffix)
  cfrp expects a sibling-versioned repo (gitmap-v28, coding-guidelines-v23, ...).
  To proceed anyway (skip the rewrite, still flip visibility public):
      gitmap-v28 cfrp <url> --require-version=false
```

Exit code: `ExitCloneFixRepoBadFlag` — same family as a malformed
flag, so CI scripts that already key on this code keep working.

## Why not auto-skip silently?

The original behavior was auto-skip silently. Two CLI runs in
production accidentally published private prototype repos that
happened to lack a `-v<N>` suffix. The cost of the false-positive
"this rewrite didn't run" is much higher than the cost of one extra
flag for the rare bypass case.

## Cross-refs

- `spec/04-generic-cli/27-fix-repo-command.md` — the wrapped command.
- `mem://constraints/strictly-prohibited` — registry of "never silently
  skip a destructive step" rules.
