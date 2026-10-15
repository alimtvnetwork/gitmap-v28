# 03 — Help Example Format (canonical)

## 1. Contract

Every command MUST expose `gitmap help example <command>` rendering three
parts in this order, every time:

1. **What it does** — one plain-English paragraph (≤3 sentences).
2. **Concrete invocation** — the literal command line, copy-pasteable, with
   real-looking arguments (no `<foo>` placeholders; use real values like
   `alimtvnetwork/story-telling-v1`).
3. **Expected output** — the representative output with unstable values
   normalized to placeholders: `<sha>`, `<date>`, `<machine>`, `<count>`,
   `<path>`. Never paste raw machine-specific output as the contract.

Render layout:

```text
gitmap scan — help example

WHAT IT DOES
  Scans a directory tree for git repositories and caches them in the
  Split-DB. Use it before any repo-wide command.

INVOCATION
  $ gitmap scan ~/workspace/repos

EXPECTED OUTPUT
  scanned: 14 repos in 0.83s (cache: gitmap.db)
  path                                branch  status
  repos/gitmap-v28                    main    clean
  repos/story-telling-v1              main    dirty
```

## 2. Single source of truth

One data registry (new file `cli/helpdoc/examples.go`): a Go map keyed by
command verb → `HelpExample{WhatItDoes string, Invocations []Example}` where
`Example{Invocation, ExpectedOutput string}`. Renderers consume the registry:

- `gitmap help example <cmd>` (new help subcommand)
- `gitmap help <cmd>` ("Examples" section appended)
- docs web UI help screens (`cli/helpdoc/*.md` generated — never hand-edited)

The existing `03-ai-scripts/51-helptext-generator.py` flow stays: generator
reads the registry, not hand-written markdown. Hand-maintaining the same
example in two formats is banned.

## 3. Coverage rule

Every dispatched command AND alias gets an entry. Audit rule: the set of
registry keys must equal the set of dispatcher verbs (minus the root `help`
meta-verbs, which get their own example page). The missing-command audit
(spec file 04) reuses this equality check.

## 4. Quality gates

- Invocation must be runnable as shown (verified in E2E: run it, compare
  output shape — values normalized per §1).
- Expected output must be ≤25 lines; full dumps go to `--json` instead.
- Examples must render correctly under `--theme mono` and when piped
  (no ANSI, no control chars in output).
- One example is the minimum; up to 3 where a command has genuinely
  distinct modes (e.g. `scan` bare vs `scan export`).
