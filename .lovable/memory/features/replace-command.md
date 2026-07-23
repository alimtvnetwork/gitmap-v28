---
name: Replace Command
description: gitmap replace — repo-wide literal text swap and version-suffix bump driven by remote URL, with audit and all modes (v3.96.0)
type: feature
---

`gitmap replace` performs deterministic find/replace across every text
file in the current repo. Default mode is interactive confirm before
write.

## Modes

- `gitmap replace "<old>" "<new>"` — literal text swap.
- `gitmap replace -N` — bump v(K-N)..v(K-1) → vK where K is the
  current version parsed from `git remote get-url origin`'s
  `<base>-vN` suffix.
- `gitmap replace --audit` — report-only scan, never writes.
- `gitmap replace all` — equivalent to `-N` with N = K-1.

## Flags

`--yes`/`-y`, `--dry-run`, `--quiet`/`-q`, `--ext` (comma-separated
extension allow-list, leading dot optional, deduplicated; empty = no
filter), `--ext-case sensitive|insensitive` (default `insensitive`).
Sensitive mode preserves user-supplied casing in `--ext` and matches
filenames byte-exact. v3.98.0+.

## Exclusions

Dirs: `.git`, `.gitmap`, `.release`, `node_modules`, `vendor`.
Path prefixes: `.gitmap/release`, `.gitmap/release-assets`.
Binary files (null-byte sniff in first 8 KiB).

## Files

- `gitmap/cmd/replace.go` — entrypoint + mode classifier.
- `gitmap/cmd/replaceflags.go` — flag parsing (`--audit` stripped pre-Parse).
- `gitmap/cmd/replacewalk.go` — repo walk, exclusions, binary sniff.
- `gitmap/cmd/replaceapply.go` — scan + atomic temp+rename writer.
- `gitmap/cmd/replaceversion.go` — remote slug parser → `(base, K)`.
- `gitmap/cmd/replaceversionrun.go` — `-N` and `all` runners.
- `gitmap/cmd/replaceaudit.go` — `--audit` line-level reporter.
- `gitmap/constants/constants_replace.go` — all messages/flags.
- `gitmap/helptext/replace.md` — embedded help.
- `spec/04-generic-cli/15-replace-command.md` — canonical spec.
