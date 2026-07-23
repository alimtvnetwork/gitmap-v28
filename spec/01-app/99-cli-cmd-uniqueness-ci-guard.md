# 99 — CLI `Cmd*` Uniqueness CI Guard + `topLevelCmds` Registry Contract

> **Audience:** NEA and any future AI/maintainer touching `gitmap-v27/constants/constants_cli.go` or its sibling files.
> **Status:** Active since gitmap-v27 **v3.11.0** (full-name uniqueness) and **v3.11.1** (alias uniqueness).
> **Related:** `mem://tech/go-namespace-constraints`, `mem://tech/constants-structure`, `mem://features/marker-comments`, `spec/04-generic-cli/15-constants-reference.md`, `spec/01-app/02-cli-interface.md`, `spec/01-app/38-command-help.md`, `spec/01-app/39-shell-completion.md`.

## 1. Goal

Two distinct `Cmd*` identifiers in the `constants` package must never share the same string value. A collision is one of two failure modes:

1. **Long-form name collision** — e.g. two constants both equal to `"release-alias"`. The Go compiler will catch this only if the names also collide; otherwise the runtime dispatcher silently routes one command and never the other.
2. **Short alias collision** — e.g. a hypothetical `CmdFooAlias = "ls"` shadowing the existing `CmdListAlias`. The compiler is happy, the dispatcher picks whichever case-arm runs first, and the user's `gitmap-v27 ls` jumps to the wrong feature without warning.

Both modes have shipped to users in past versions (`CmdReleaseAlias` bound twice in v3.10.x; `cd` / `go` shadowing in v3.10.x). The CI guard added in v3.11.x rejects them at `go test` time, **before** the build phase, so no further regressions can reach a release tag.

## 2. The two tests

Both live in `gitmap-v27/constants/cmd_constants_test.go`:

### 2.1 `TestTopLevelCmdConstantsAreUnique` — full-name uniqueness

Iterates every entry in `topLevelCmds()` and fails if any two entries share the same string value, regardless of length. Catches the v3.10.x-class redeclaration bugs.

### 2.2 `TestTopLevelCmdAliasesAreUnique` — short-alias uniqueness

Iterates the same registry but only compares entries whose value length is `<= maxAliasLen` (currently `2`). Skips empty strings and long-form names. Catches future regressions like `CmdFooAlias = "ls"` colliding with `CmdListAlias`.

The split is deliberate: the alias test surfaces a tighter, more actionable error message ("duplicate short alias `ls`") and remains useful even if a future refactor relaxes the long-form test (e.g. to permit intentional aliasing).

## 3. The `topLevelCmds()` registry — manual source of truth

```go
func topLevelCmds() map[string]string {
    return map[string]string{
        "CmdScan":      CmdScan,
        "CmdScanAlias": CmdScanAlias,
        // … one row per top-level Cmd* identifier …
    }
}
```

### 3.1 What goes in

Every constant whose name starts with `Cmd` AND that the root dispatcher (`gitmap-v27/cmd/rootdata.go`) routes as a **top-level** command or its alias. This is the same set the shell-completion generator considers — see `mem://features/marker-comments` for the marker-comment opt-in convention.

### 3.2 What stays out

* Subcommand verbs (`create`, `add`, `list`, `remove`, etc.) reused inside subcommand groups like `gitmap-v27 group create` / `gitmap-v27 alias create`. These are intentionally duplicated, marked with `// gitmap-v27:cmd skip` in `constants_cli.go`, and excluded from `topLevelCmds()`.
* Runner helpers like `CmdSelfUninstallRunner` (also marked `// gitmap-v27:cmd skip`).
* Anything in a non-CLI namespace (e.g. internal flag names, message strings).

### 3.3 Drift contract

The registry is **manual**. When you add or remove a top-level `Cmd*` constant:

1. Add or remove the matching row in `topLevelCmds()`.
2. Run `go test ./gitmap-v27/constants/ -run TestTopLevelCmd -v` locally.
3. Confirm both tests pass.

If the registry drifts (a new top-level constant exists in `constants_cli.go` but is missing from the map), the test will not catch the new constant's collisions — but `mem://features/marker-comments`'s `generate-check` CI step will catch the drift indirectly because the completion generator emits a different surface than the registry covers. A future hardening (see §5) closes this gap by deriving the map from the AST.

## 4. NEA / AI handoff checklist

When adding a new top-level CLI command:

1. Define `CmdFoo` (and optionally `CmdFooAlias`) in `gitmap-v27/constants/constants_cli.go` under the `// gitmap-v27:cmd top-level` block.
2. Add the dispatch arm in `gitmap-v27/cmd/rootdata.go`.
3. Append the matching rows to `topLevelCmds()` in `gitmap-v27/constants/cmd_constants_test.go`.
4. Run `go test ./gitmap-v27/constants/ -run TestTopLevelCmd -v` — must pass.
5. Run the completion generator (`go generate ./completion/...` or whatever the `Makefile` target is) — must produce no diff.
6. If you intentionally introduced a duplicate that the dispatcher handles via a subcommand group, mark the constant with `// gitmap-v27:cmd skip` and DO NOT add it to `topLevelCmds()`.

## 5. AST-derived parity test (registry drift guard)

Implemented in `gitmap-v27/constants/cmd_constants_parity_test.go` as `TestTopLevelCmdRegistryMatchesAST`. The test uses `go/parser` to walk every `constants_*.go` file in the package, collects every `Cmd*` string constant declared inside a const block marked `// gitmap-v27:cmd top-level` (excluding per-spec `// gitmap-v27:cmd skip` lines), and asserts the resulting set of identifier names is **exactly equal** to the keys of the manual `topLevelCmds()` map.

Two failure modes are reported with actionable messages:

* **Missing from registry** — A new top-level `Cmd*` constant exists in the AST but was not appended to `topLevelCmds()`. Fix: add the row, or mark the constant `// gitmap-v27:cmd skip` if it should be excluded.
* **Extra in registry** — A row in `topLevelCmds()` references a constant that no longer exists under a `// gitmap-v27:cmd top-level` block. Fix: remove the row, or restore the constant under an opted-in block.

Combined with §2.1 / §2.2, this closes the drift gap: the registry can no longer fall out of sync with the source of truth, which means the value-uniqueness and alias-uniqueness tests can never be silently bypassed by an unregistered constant.

## 6. History

| Version | Change |
|---|---|
| v3.11.0 | Initial `TestTopLevelCmdConstantsAreUnique` test + `topLevelCmds()` registry. Caught the `CmdReleaseAlias` and `cd`/`go` redeclaration bugs that motivated this guard. |
| v3.11.1 | Added `TestTopLevelCmdAliasesAreUnique` to specifically guard the short-alias namespace (length `<= 2`). |
| v3.12.0 | Spec doc written so NEA can extend the guard without re-deriving its rationale. Added `TestTopLevelCmdRegistryMatchesAST` (AST-derived parity test) — registry drift is now impossible. |
