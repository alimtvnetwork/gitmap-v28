# Component Spec — GitMap and Antigravity CI/CD and Code Refactoring

## Component 1: Antigravity-Manager release train

- Input: latest main tip, CI status for all workflows on that tip.
- Steps: confirm CI green on the tip; if red, diagnose from job logs, fix, commit via gitmap cpb, re-check until green. Then run the repo bump script for a minor bump (current tip version + 0.0.1 on the minor line — 4.266.0 → 4.267.0 unless the tip moved). Create the release tag using gitmap, then `gitmap release`. Verify: `git ls-remote origin "vX.Y.Z^{}"` shows the tag on the release commit; `gh release view vX.Y.Z` shows published, not draft, with the full asset set, marked Latest.
- Guard: AGM main ruleset 23609023 rejects force-push. Never force-push, never rewrite. If the tip moves under the worker, fetch and re-base the bump onto the new tip.
- Out of scope: any AGM-wide enum refactor; any change outside what CI-green and the release require.

## Component 2: gitmap-v28 CI-green

- Input: latest main tip (d886982d at plan time: the pending gosec fix eac563f2 plus the merge of origin 32d42a14 were just pushed; CI runs in flight on it).
- Steps: watch the runs on the tip. If green, done. If red, diagnose from the failing job logs, fix, commit via gitmap cpb, re-watch. Known history: v6.600.0 tag points at broken commit 208c46ef (unused strings import in cli/cmd/rootcore.go, fixed by 4bd9c937); race-detector failed on the other session's 4113c0c3 tip. Diagnose from the current tip, not from memory.
- Guard: never re-release v6.600.0. The recovery release is a new minor bump after the refactor waves (Component 4).

## Component 3: gitmap-v28 enum/constants refactor (the deep work)

Work group by group, package by package, incremental commits per group, CI check per group or batch. Audit first, then waves. Rules, each non-negotiable:

1. Enum-per-package: every enum gets its own package with the type convention from the coding guidelines. No enum shares a package with unrelated code. Where an enum already lives in its own package, record it as done; where it does not, move it.
2. No string concatenation with `+` for enum values or mixed strings: use the enum's String() method (automatic conversion) or combine values into a reusable constant. Exemplar to fix first and use as the reference: cli_cmd_cargo/cargo_exe.go.
3. No regular errors anywhere: always the app error type, enriched with the variables and paths involved so future debugging is easy. Named example: applyDefaultPragmas in the constants section returns a regular error — convert it to the app error with context.
4. Flag comparisons: cli/cmdpool/pool_cli_input.go around line 79 compares args with `==` chains (group flag, -p flag, many flags). Replace with a map/dictionary lookup plus explicit variables; flags that come from enums must be converted to a flag type first. Build one reusable dictionary and check membership against it. Apply the same anywhere `==` chains compare flags or args.
5. Regex: lazy regex (compiled once, package level) everywhere; some places already do — finish the rest.
6. sqlite_config.go: (a) naming — identifiers such as Database/Db used inconsistently must follow PascalCase per the naming convention, everywhere they appear; (b) wrapped returns — methods returning two things (app error + object) must be combined into one class type: a generic type that is also defined without generics so it looks flat. Before touching this file, read the wrapped-return / wrapped-results prompts in coding-guidelines-v24 and the recent gitmap changes that introduced this pattern, and follow that convention exactly.
7. File cap: keep every refactored file under 300 lines; split when needed.
8. Each wave: audit the group, refactor, gofmt + go build pre-flight, commit via gitmap cpf/cpb, CI check.

## Component 4: gitmap-v28 release

- Only after CI green and the refactor waves committed: minor bump via the bump script from the current tip version, tag via gitmap, release via `gitmap release`. Verify tag via ls-remote and the published release (assets, Latest).

## Component 5: coding-guidelines-v24 prompt update

- The enum CG execute prompt must strictly state the enum-per-package + type convention rule, with the cargo_exe.go example called out as the anti-pattern and the corrected form shown. Mirror to .agents/skills and .cursor/skills mirrors as applicable. Commit via the gitmap suite as the user.
- Depends on Component 3's exemplar fix so the prompt quotes the real corrected code.
