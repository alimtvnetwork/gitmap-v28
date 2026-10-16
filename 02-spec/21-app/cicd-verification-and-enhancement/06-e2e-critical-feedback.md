# 06 — E2E Critical Feedback (CLI-level, zero live side effects)

**Wave:** 2, WS4
**Date:** 2026-10-17
**Binary under test:** `gitmap v6.553.0` (commit `7956d80eb71a0fa59f0ff5d0778832de98bf6caf`)
**Repo HEAD at test time:** past v6.556.2 (per `/tmp/274-last100.txt`)
**Method:** `--help` / `--dry-run` / `--json` probes only. No PRs opened, no pushes, no state mutated.

---

## 0. The headline finding: the installed binary predates every feature under test

This has to come first because it conditions everything below.

The `~/.local/bin/gitmap` binary reports v6.553.0 / commit `7956d80eb`.
That commit is **not** in the last-100 list, and the 273 helpdoc feature
commit (`19b59869c` — "273 help example registry and 4 design system themes")
landed **after** the v6.553.0 release bump (`d80927b23`). Consequences:

- `gitmap help example <cmd>` — does not exist in this binary.
- `--theme navy-purple|tokyo-night|vscode-dark-plus|monokai-pro` — all rejected.
- `gitmap merge write` — "not recognized".
- `gitmap mr` — still the old `merge-right` alias.

All of these exist in HEAD **source** (verified by reading the code:
`cli/cmdmerge/merge_write*.go`, `cli/cmd/rootutility.go:200,263`,
`cli/cmd/themeflag.go` + `termout.ThemeSlugs()`), but **none of them could be
E2E-verified** against a running binary. So for items 1–3 the honest verdict
is: **not E2E-testable from the shipped binary; source-level review only.**
Everything in §1–§3 marked `[source]` was read, not run.

**Suggestion:** gate this wave on `gitmap version` == HEAD `version.json`
before probing, and rebuild the local binary from HEAD (the standing rule is
"update by rebuilding from source"). E2E against a binary that predates the
features is theater.

---

## 1. `mr` / `merge write` — [source] collision handled, but risky and un-dry-runnable

**What works (source):** The implementation exists at
`cli/cmdmerge/merge_write.go` + `merge_write_gh.go` + `merge_write_help.go`.
Flag set matches the spec: `--title` (required), `--body`, `--source`,
`--target`, `--draft`, `--yes`/`-y`. Shells out to `gh pr create`.
Confirmation prompt before creating; aborts on non-`y` answer, and on EOF
(`fmt.Scanln` error → `false` → abort), so piped/non-TTY usage fails closed.
Good.

**What's broken / awkward:**

1. **Silent alias reassignment — the biggest concern in this whole review.**
   `mr` was `merge-right`'s alias for years ("One-way file-level merge that
   writes only into RIGHT"). The new code sets `CmdMergeWriteA = "mr"` and
   demotes merge-right to `mrg` (`cli/constants/constants_movemerge.go:15,21`).
   So `gitmap mr LEFT RIGHT` — which a user may have in scripts, muscle
   memory, macros — will now open a **GitHub pull request** instead of merging
   files. That is a silent behavioral break of a destructive-ish command, and
   it landed in the same program that removed the deprecation policy
   (program 272 item 1). The removed policy is exactly what would have caught
   this. Concrete suggestion: either keep `mr` on merge-right and give
   merge-write a fresh alias (`mwr`), or detect the old arity (`mr LEFT RIGHT`
   with two positional args and no `--title`) and print an explicit migration
   error instead of doing something surprising.
2. **No `--dry-run`.** A command whose entire job is creating a real GitHub
   PR offers no simulation mode. The E2E plan even probed for it
   (`mr --dry-run` "if supported") — it isn't. For the highest-side-effect new
   command in the CLI, that's backwards. Add `--dry-run` that prints the exact
   `gh pr create` argv and exits.
3. **Redundant strip:** `stripMergeWriteSubcommand` is called in
   `DispatchMoveMerge` and again inside `runMergeWrite`→`parseMergeWriteFlags`.
   Second call is a no-op. Harmless, but it signals the flow wasn't re-read
   after writing.
4. **Mixed arg conventions:** `isMergeWriteInvocation(command, os.Args)`
   reads the global `os.Args` while everything around it passes `os.Args[2:]`.
   Works today; fragile for any future refactor that calls the dispatcher with
   synthetic argv (tests, TUI).

**Binary behavior (for the record, old `mr` = merge-right):**
`gitmap mr --dry-run` → `Usage: gitmap merge-right LEFT RIGHT [flags]` (exit 0).
`gitmap mr --help` → the generic boilerplate described in §4, **not**
merge-right's real help. So the alias already has a help-resolution bug in
the current binary, independent of the 273 change.

**E2E verdict:** not testable (binary too old). Source review: works, but §1.1
and §1.2 should block any release that ships this binary.

---

## 2. `gitmap help example <cmd>` — not E2E-testable; source wiring looks sane

**Binary:** all five probes (`mr`, `scan`, `cpf`, `nc`, `pe`, `spec next`)
return the same output:

```
No matches for "example"

Did you mean:
 • execute (exec) <name>      Replay recorded macro steps ...
 • chrome-profile-export (cpe) <name>  ...
 • clone-fix-repo-pub (cfrp) <url> [folder]  ...
```

(exit 0). The subcommand doesn't exist in this binary, so consistency of the
3-part format **cannot be verified**. Nothing here is a product bug — it's the
stale binary again.

**Source:** `cli/cmd/rootutility.go:200` (`isHelpExample`, case-insensitive
match on `os.Args[2]`) → `dispatchHelpExample` (line 263) →
`helpdoc.PrintHelpExample`. Wiring is straightforward; `help example` with no
command prints an index. One nit: the `EqualFold` check runs before topic
dispatch, so a hypothetical future command literally named `example` would be
unreachable — not a problem today.

**E2E verdict:** SKIPPED — unsafe to claim (binary predates feature).
**Suggestion:** re-run the five probes after a HEAD rebuild; spot-check that
the "expected output" part of the 3-part blocks is actually byte-accurate
against real command output, since these registries rot.

---

## 3. Four new themes — not E2E-testable; one genuine UX positive found

**Binary:** all four slugs rejected identically:

```
gitmap: [E1152:VALIDATION] cmd.themeflag.apply: invalid --theme value "navy-purple" (want: bright | standard | monochrome)
  origin: cmd/themeflag.go:54
```

(exit 0). Same for `tokyo-night`, `vscode-dark-plus`, `monokai-pro`, and a
control `bogus-theme`. So: no theme rendering could be compared.

**What works:** the invalid-theme error itself is genuinely good UX — it names
the rejected value and lists the accepted ones. HEAD source improves it
further (lists new slugs plus legacy aliases: bright/standard/monochrome/mono).
Keep this pattern; it's the best error message I hit in the whole session.

**Source (HEAD):** `stripThemeFlag` accepts `--theme <v>`, `--theme=<v>`,
`-theme`; invalid values exit via `cliexit.Exit(0)`... wait — the comment in
HEAD source says "exits with status 2" but the binary shows EXIT: 0 for the
same error class. If HEAD still exits 0 on invalid theme, that's §5's issue
persisting. Worth a check on rebuild: `gitmap --theme bogus; echo $?` must be
non-zero.

**E2E verdict:** SKIPPED (binary predates feature).
**Suggestion:** when re-testing, render the *same* command (e.g. `nc`) under
all four themes plus `--pretty=off` and diff — the 273 claim is "dumb-terminal
safe"; verify ANSI degrades cleanly when piped.

---

## 4. Generic fallback help fabricates help for commands that don't exist — REAL BUG

This is the most concrete defect found, and it's in the **current binary**:

```
$ gitmap cpu --help
  ╔══════════════════════════════════════════════════╗
  ║                    gitmap cpu                    ║
  ╚══════════════════════════════════════════════════╝
  Usage:
    gitmap cpu [arguments] [flags]
    gitmap cpu help
  Description:
    cpu    Execute, manage, and automate GitMap cpu operations.
  ...
```

(exit 0) — while:

```
$ gitmap cpu
╭─ Did you mean? ─╮
│ Command 'cpu' is not recognized. Closest matches:
│  • gitmap cd   (Confidence: 33%)
...
gitmap: [E1001:VALIDATION] cmd.dispatch: unrecognized command: cpu
```

`cpu` and `nro` don't exist in this binary (they shipped in v6.554.0, after the
binary's build commit). With `--help` they get a **plausible-looking but
entirely fabricated** help screen; without it they correctly error. So the
`--help` path and the dispatch path disagree about what exists, and `--help`
always wins with fiction. Same root cause explains `gitmap mr --help`
rendering boilerplate ("Execute, manage, and automate GitMap mr operations.")
while `gitmap mr --dry-run` correctly reaches merge-right
("Usage: gitmap merge-right LEFT RIGHT [flags]") and `gitmap merge-right
--help` renders the real help page. Aliases are not resolved on the help
path — inconsistently, since `gitmap rrb --help` and `gitmap pnr --help`
render correct alias help ("Alias for gitmap release remove-between").

**Why it matters:** a user typing `gitmap <typo> --help` gets confident,
official-looking documentation for a command that doesn't exist, with exit
code 0. Scripts and agents (like this one) will trust it. This is worse than
no help.

**Suggestions:**
1. Make the fallback renderer verify the command against the dispatcher
   registry first; unknown → the same "not recognized" + did-you-mean error
   as the no-flag path.
2. Resolve aliases on the help path (rrb/pnr prove it's possible; mr proves
   it's inconsistent).
3. Scrub the template: the fallback tip reads "Combine with '--json' for
   pipeline and AI automation **fspath**." — "fspath" is garbled filler.

---

## 5. Errors exit 0 — REAL BUG (scripting hazard)

Every error probed — E1001 unrecognized command, E1152 invalid theme —
printed an error and **exited 0**. Verified via `echo $?` after each. Any
script doing `gitmap <cmd> || handle_failure` will silently take the success
branch. HEAD's themeflag comment claims "exits with status 2" for invalid
theme, but the observed binary exits 0 — either the comment is aspirational
or the binary predates the fix; either way, re-verify on a HEAD build and
make E1001/E1152 (at minimum) non-zero. This is a one-line class of fix with
outsized CI/scripting value.

---

## 6. `gitmap nc` (codegen catalog) — WORKS

- `gitmap nc` → `Total: 105 | Returned: 100 | Limit: 100`, grouped listing,
  exit 0. Matches the "105 entries" claim from the codegen commit.
- `gitmap nc --json` → parses cleanly as JSON (`totalCommands: 105`,
  `filteredCommands: 100`, `limit: 100`, `commands: [...]`).
- `spec next` present and correct: "Issue next free spec number atomically
  via DB-backed counter", example `gitmap spec next`. `te all` present.
- `gitmap nc --help` renders real help (alias `nc` documented, usage lines).

Minor nits: the header says "GITMAP RECENT COMMANDS DISCOVERY ENGINE" while
`--help` titles it "gitmap new-commands" — two names for one command is a
small consistency wart. The 100-result limit vs 105 total is disclosed
honestly in the header, which is good. (Whether the catalog is truly
build-time generated rather than a checked-in snapshot couldn't be verified
without building — noted, not claimed.)

---

## 7. Other recent features (spot checks from `/tmp/274-last100.txt`)

- `gitmap rrb --help` / `gitmap pnr --help` → real help, correct alias
  attribution ("Alias for gitmap release remove-between /
  pin-not-to-remove"), synopsis lines. PASS.
- `gitmap cpu --help` / `gitmap nro --help` → the fabricated-help bug (§4).
  The commands themselves (os-hardware telemetry, nodes-running-only) are
  v6.554.0+, absent from this binary — not E2E-testable here.
- `gitmap merge write --help` → "Command 'merge' is not recognized" + good
  did-you-mean box (reconcile 80%, rerun 40%). Correct behavior *for this
  binary*; will change once HEAD is built.

Not E2E-tested (would cause side effects): actual `mr` PR creation,
`release remove-between`, anything under `pa`/`pull`. Marked per the rules.

---

## 8. Consolidated recommendations (priority order)

1. **Rebuild `~/.local/bin/gitmap` from HEAD before any further E2E.**
   The v6.553.0 binary contains none of program 273's features; wave 2 cannot
   validate what isn't in the binary. Then re-run §1–§3 probes.
2. **Reconsider the `mr` alias reassignment** (file-merge → PR-creation is a
   silent break; no deprecation policy left to catch it). Prefer a fresh alias
   (`mwr`) or an explicit migration error on the old arity.
3. **Add `--dry-run` to `merge write`.** The one command that creates real
   GitHub PRs is the one without a simulation mode.
4. **Fix the fabricated-help fallback** (§4): unknown command + `--help`
   must error like the dispatch path, and aliases must resolve consistently
   on the help path.
5. **Non-zero exit codes for E1001/E1152** (§5).
6. **Scrub "fspath"** from the fallback help template.
7. On the rebuild, verify: `help example` 3-part accuracy against real output
   (registries rot), theme rendering diffed across all four themes piped
   through `cat` (dumb-terminal safety claim), and `gitmap --theme bogus;
   echo $?` ≠ 0.

## 9. Scorecard

| Feature | Verdict |
|---|---|
| `mr` / `merge write` | NOT E2E-TESTED (binary predates it); source OK but §§1.1–1.2 are release blockers |
| `help example <cmd>` (5 probes) | NOT E2E-TESTED (not in binary); source wiring sane |
| 4 new themes | NOT E2E-TESTED (not in binary); invalid-value error message itself is good |
| `nc` catalog (105 entries, `--json`, `spec next`, `te all`) | PASS |
| `rrb` / `pnr` help | PASS |
| `cpu` / `nro` help | FAIL — fabricated help for nonexistent commands (§4) |
| Error exit codes | FAIL — errors exit 0 (§5) |

**No side effects were caused.** Every probe was `--help`, `--dry-run`, or
`--json`. No PRs opened, nothing pushed, no state mutated.
