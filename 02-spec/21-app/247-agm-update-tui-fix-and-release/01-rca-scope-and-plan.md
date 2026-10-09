# 247 — AGM update TUI fix and release: RCA, scope, and plan

## 1. Problem statement

`gitmap agm update` has two terminal-UI failures. On Linux/macOS the command goes
completely silent for the minutes the installer runs (`CombinedOutput` swallows
everything), then prints either one success line or a full raw output dump plus a
stack trace on failure — nothing in between. On Windows (and the Unix dispatch
path) the installer inherits the terminal's stdio directly, so its raw escape
sequences collide with anything gitmap is rendering and corrupt the screen. The
user asked for a GitLab-style professional progress UI and for the screen crash
to be root-caused and fixed, then released.

## 2. Root cause

- **Linux/macOS silence:** `runAgmUpdateLinuxQuiet` (`cli/cmdinstall/agm_update.go:103`)
  runs the installer via `cmd.CombinedOutput()` (`agm_update.go:114`) — all
  installer output is buffered until the process exits. On failure,
  `printUpdateFailureDetails(err, out)` (`agm_update.go:126`, called at :116)
  dumps the entire buffer plus a stack trace. No streaming, no stages, no UI.
- **Screen corruption:** `dispatchAgManagerWindowsWithVersion`
  (`cli/cmdinstall/installagmanager.go:156-166`) sets
  `cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin` (:163);
  the Unix sibling `dispatchAgManagerUnixWithVersion` (:168-176) does the same
  (:173-175). The child installer writes escape sequences straight to the
  terminal with zero coordination with gitmap's own output — two writers, one
  screen, no renderer owning it. That is the crash.
- **Common thread:** in both paths no single renderer owns the screen. Output is
  either fully swallowed or fully uncoordinated.
- Note: there is NO `space update` subcommand (`cli/cmdspace/` has zero hits) —
  the command in question is `gitmap agm update` (`runAgmUpdateCmd`,
  `agm_update.go:60`).

## 3. Fix design

ONE coordinated, mutex-guarded renderer fed by CAPTURED (never inherited) child
output. Reusable precedents already in the repo: `cli/cloner/batchprogress.go:18`
(BatchProgress), `cli/cluster/pool.go` (pterm spinners), `cli/clonepick/picker.go`
(bubbletea). pterm v0.12.83 is already in go.mod — no new dependencies.

- **(a) Linux path:** replace the `CombinedOutput` silence with streaming: run the
  installer with stdout/stderr pipes, feed lines into a pterm-based GitLab-style
  UI — spinner + live log tail + per-stage checkmarks as installer stages
  complete. The full log is still retained in memory for the failure panel.
- **(b) Windows/Unix dispatch:** stop inheriting the terminal. Capture stdio via
  pipes in `dispatchAgManagerWindowsWithVersion` and
  `dispatchAgManagerUnixWithVersion` and feed the SAME renderer as (a); stdin
  pass-through only when the installer genuinely prompts (default: no PTY
  inheritance).
- **(c) Failure path:** replace the raw dump + stack trace with a structured
  pterm error panel (what failed, which stage, last N log lines, next step).
  The retained full log goes to the error-log pipeline, not the screen.
- **(d) Non-TTY fallback:** when stdout is not a TTY (piped, CI, `TERM=dumb`),
  emit plain lines, zero escape sequences.
- **glyphs.Install decision (explicit):** the update UI renders through the
  EXISTING glyphs pipeline — NO `byteFaithfulCommands` exemption. The exemption
  (task 244) exists for byte-faithful *content* (`cat`/`view`/`type`); the update
  UI is *presentation*, so the legacy-terminal transliteration is desirable, and
  pterm already degrades on dumb terminals. `cli/glyphs/` is not touched.

## 4. File box

IN SCOPE — `cli/cmdinstall/agm_update.go`, `cli/cmdinstall/installagmanager.go`,
plus ONE new helper file/package under `cli/cmdinstall/` if a shared renderer is
needed (keep it inside `cmdinstall`; do not create cross-cutting TUI packages).

EXPLICITLY OUT OF SCOPE — `cli/cmdupdate/`, `cli/cmddownload/` (the `gitmap update`
self-updater shares the swallow/inherit pattern but is a follow-up task, not this
one), `cli/glyphs/`, anything task 246 owns. Task-246 collision assessed LOW:
246 never touches `cli/cmdinstall/`.

## 5. Release plan

Version **6.517.0** (next minor: `version.json` is stale at `6.515.0`; tag
`v6.516.0` was cut directly by 243b).

1. Bump via `03-ai-scripts/37-bump-version.py` (or manual equivalent if the
   script's contract changed — verify first).
2. Changelog entry for the TUI fix.
3. Annotated tag `v6.517.0` with the tag message exactly `Release v6.517.0`.
4. Push `main` + the tag; verify the tag exists on the remote.

DOCUMENTED DEVIATION: the release orchestrator's quality gates run the test
suite — SKIPPED per the standing never-run-tests rule (owner's explicit command
required). E2E verification below is the quality evidence for this release.

## 6. Test plan

- `go build ./...` from `cli/` only. NEVER `go test` (standing rule).
- E2E with a MOCK installer script (emits staged stdout lines, sleeps between
  stages, optional `--fail` mode): run the fixed `gitmap agm update` against the
  mock and verify (i) the progress UI renders (spinner/stages, no minutes of
  silence), (ii) success path completes, (iii) failure path shows the structured
  panel and no raw dump/stack trace, (iv) non-TTY run emits plain lines only.
- All e2e artifacts (mock script, logs) live OUTSIDE the repo and are never
  committed.

## 7. Non-goals

- No `gitmap update` self-updater changes (follow-up task).
- No `cli/cmd/root.go` or `cli/glyphs/` changes.
- No new Go dependencies (pterm is already vendored in go.mod).
- No refactoring of the installer scripts themselves (`install.ps1`/`install.sh`
  content is unchanged; only how their output is captured and rendered changes).
