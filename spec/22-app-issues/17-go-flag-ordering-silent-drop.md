# Post-Mortem #17: Go Flag Package Ordering — Silent Flag Drop

## Problem

Running `gitmap-v28 release v2.55.0 -y` silently ignored the `-y` flag,
causing the release to pause for an interactive confirmation prompt
even though the user explicitly passed `--yes`.

### Impact

1. **Silent flag drop** — no error was reported; the flag simply had
   no effect, making it extremely difficult to diagnose.
2. **CI/automation breakage** — non-interactive pipelines using `-y`
   would hang waiting for user input.
3. **Inconsistent UX** — `gitmap-v28 release -y v2.55.0` worked correctly,
   but `gitmap-v28 release v2.55.0 -y` did not, violating the principle
   of least surprise.

### Root Cause

Go's standard `flag` package **stops parsing at the first non-flag
argument**. When the user writes:

```
gitmap-v28 release v2.55.0 -y
```

The `flag.Parse()` call sees `v2.55.0` as the first positional argument
and stops. Everything after it — including `-y` — is placed in
`flag.Args()` as unparsed trailing arguments.

```
Parsed flags: (none)
Positional:   ["v2.55.0", "-y"]
```

This is documented Go behavior (`flag.Parse` uses `flag.ContinueOnError`
or `flag.ExitOnError`, but both stop at the first non-flag token), but
it differs from GNU-style flag parsing where flags can appear anywhere
in the argument list.

## Fix

Added `reorderFlagsBeforeArgs()` in `cmd/releaseargs.go`. This function
runs **before** `flagSet.Parse()` and reorders the raw argument slice so
that all flag-like tokens (starting with `-`) appear before positional
arguments.

### Before (broken)

```
args = ["v2.55.0", "-y"]
flagSet.Parse(args)
// -y is never parsed
```

### After (fixed)

```
args = reorderFlagsBeforeArgs(["v2.55.0", "-y"])
// args = ["-y", "v2.55.0"]
flagSet.Parse(args)
// -y is correctly parsed
```

### Value-Flag Awareness

Flags that consume the next argument as a value (e.g., `--bump patch`,
`-N "release notes"`) are kept together with their value during
reordering. A lookup table of known value flags ensures the reorder
logic doesn't split a flag from its argument:

```go
valueFlags := map[string]bool{
    "--assets": true, "--commit": true, "--branch": true,
    "--bump": true, "--notes": true, "--targets": true,
    "--bundle": true, "--zip-group": true,
    "-N": true, "-Z": true,
}
```

## Key Files

| File | Role |
|------|------|
| `cmd/releaseargs.go` | `reorderFlagsBeforeArgs()` — pre-parse argument reordering |
| `cmd/release.go` | `parseReleaseFlags()` — calls reorder before `flagSet.Parse()` |

## Lessons Learned

1. **Go's `flag` package is not GNU-style** — flags after the first
   positional argument are silently ignored. Any CLI that accepts both
   flags and positional args must reorder or use a third-party parser.
2. **Silent failures are worse than crashes** — the flag was simply
   dropped with no diagnostic output, making the bug invisible until
   someone noticed the prompt appearing in CI.
3. **Argument reordering is a one-time, centralized fix** — by
   reordering before parsing, all current and future flags benefit
   without per-flag workarounds.

## Version

Introduced: **v2.58.0**
