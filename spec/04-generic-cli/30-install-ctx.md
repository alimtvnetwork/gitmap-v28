# `install ctx` — Cross-Platform Right-Click Context Menu

> Status: ships on **Windows** (HKCU registry cascade), **macOS**
> (Automator Quick Action `.workflow` bundles in `~/Library/Services`)
> and **Linux** (Nautilus scripts + Dolphin service menu + Thunar
> `uca.xml`). The flat menu used on macOS/Linux is generated from the
> same `[]ctxEntry` table that drives the Windows nested cascade —
> single source of truth.

## 1. Purpose

Add a single `gitmap-v28 ▸` entry to the Windows Explorer right-click menu
on **folder backgrounds** (clicking inside a folder) and on **folder
items** (right-clicking a folder). The entry expands into nested
category submenus that invoke `gitmap-v28` subcommands against the clicked
folder (`%V`).

Installed via:

```
gitmap-v28 install ctx              # add the menu (HKCU only — no admin)
gitmap-v28 install ctx --explain    # bake an "> gitmap-v28 <args>" announce into every
                                # generated entry; users see the resolved
                                # invocation before it executes. Process-local
                                # toggle: re-running without --explain
                                # regenerates plain entries.
gitmap-v28 uninstall ctx            # remove the menu
```

### 1.1 `--explain` rendering

| Mode     | Where the announce appears                                      |
| -------- | --------------------------------------------------------------- |
| Terminal | First line in the spawned terminal: `> gitmap-v28 pull --all`       |
| Silent   | Prepended to the OS notification body, ahead of stdout/stderr   |
| Prefill  | Not affected — the user is typing the command themselves        |

Implemented by `ctxExplainPrefixPwsh` / `ctxExplainPrefixSh` /
`ctxExplainAnnounce` in `gitmap-v28/cmd/installctxexplain.go`, gated by
the package-local `ctxExplainEnabled` flag set in `runInstallCtx`.

`ctx` is added to the existing install-tool table alongside `vscode-ctx`
and `pwsh-ctx`; this spec is **strictly additive** — neither of the
existing two commands is altered.

## 2. Menu Structure

One nested layout under a top-level `gitmap-v28` cascade. Categories mirror
the CLI command groups so users discover commands the same way they do
on the terminal.

```
gitmap-v28 ▸
├─ Scan ▸
│   ├─ Scan here                       (gitmap-v28 scan)
│   ├─ Rescan                          (gitmap-v28 rescan)
│   └─ Find next                       (gitmap-v28 find-next)
├─ Clone ▸
│   ├─ Clone-next here                 (gitmap-v28 clone-next)
│   ├─ Pull                            (gitmap-v28 pull)
│   └─ Pull all                        (gitmap-v28 pull-all)
├─ Release ▸
│   ├─ Release current                 (gitmap-v28 release)
│   ├─ Release next (bump minor)       (gitmap-v28 release --bump minor)
│   ├─ Release pull                    (gitmap-v28 release-pull)
│   ├─ Release pending          [N]    (gitmap-v28 release-pending)
│   ├─ List releases                   (gitmap-v28 list-releases)
│   └─ List versions                   (gitmap-v28 list-versions)
├─ Repos ▸
│   ├─ Go projects                     (gitmap-v28 go-repos)
│   ├─ Node projects                   (gitmap-v28 node-repos)
│   ├─ React projects                  (gitmap-v28 react-repos)
│   ├─ C++ projects                    (gitmap-v28 cpp-repos)
│   ├─ C# projects                     (gitmap-v28 csharp-repos)
│   ├─ Rust projects        [future]   (gitmap-v28 rust-repos)
│   └─ PHP projects         [future]   (gitmap-v28 php-repos)
├─ Visibility ▸
│   ├─ Make public                     (gitmap-v28 visibility public)
│   └─ Make private                    (gitmap-v28 visibility private)
├─ Tools ▸
│   ├─ Fix repo                        (gitmap-v28 fix-repo)
│   ├─ Diff                            (gitmap-v28 diff)
│   ├─ Logs                            (gitmap-v28 logs)
│   ├─ History                         (gitmap-v28 history)
│   └─ Update                          (gitmap-v28 update)
├─ ─────────────                       (separator)
├─ Open terminal here                  (open pwsh, prefill `gitmap-v28 `)
└─ Docs                                (gitmap-v28 docs)
```

Rust/PHP entries are stubbed in the menu only when the underlying
commands ship (gated by `constants.HasRustRepos` / `HasPhpRepos`
build-time flags). Until then the rows are omitted.

### 2.1 Windows registry layout

Use the legacy `MUIVerb` + `SubCommands` cascade (no COM handler). All
keys live under **HKCU** so install requires no elevation:

```
HKCU\Software\Classes\Directory\Background\shell\gitmap-v28
    (Default)        = (empty)
    MUIVerb          = "gitmap-v28"
    SubCommands      = ""               ; empty => use ExtendedSubCommandsKey
    Icon             = "<gitmap.exe path>,0"

HKCU\Software\Classes\Directory\Background\shell\gitmap-v28\shell\01_scan
    MUIVerb          = "Scan"
    SubCommands      = ""
    HKCU\...\01_scan\shell\01_scan_here
        MUIVerb      = "Scan here"
        \command (Default) = "<exec template>"
    ...
```

Mirror the same tree under `Directory\shell\gitmap-v28` so right-clicking
the folder **item** (not just background) also works. Generation is
table-driven from a single `[]ctxEntry` slice — see §4.

## 3. Menu → Command Mapping (authoritative)

`%V` (Windows) / `$1` (mac/Linux) is the **clicked folder**, used as the
working directory for every entry. No flag is added unless listed
below — defaults apply.

| KeyName              | Visible label                  | Exact `gitmap-v28` invocation              | Mode     | Notes                                                                              |
| -------------------- | ------------------------------ | -------------------------------------- | -------- | ---------------------------------------------------------------------------------- |
| `10_scan/10_scan_here`     | Scan here                | `gitmap-v28 scan`                          | Terminal | Walks current folder with default `--max-depth` and worker count from `git-setup.json`. No flags. |
| `10_scan/20_rescan`        | Rescan                   | `gitmap-v28 rescan`                        | Terminal | Re-runs the most recent scan against the same root.                               |
| `10_scan/30_find_next`     | Find next                | `gitmap-v28 find-next`                     | Silent   | Probes for the next available `<base>-vN+1` sibling. No `--scan-folder` (uses cwd); no `--json` (output goes to notification verbatim). |
| `20_clone/10_clone_next`   | Clone-next here          | `gitmap-v28 clone-next`                    | Terminal | Flattens to base-name folder by default (v2.75.0+). |
| `20_clone/20_pull`         | Pull                     | `gitmap-v28 pull`                          | Terminal | Fast-forward pull on current repo only. |
| `20_clone/30_pull_all`     | Pull all (every tracked repo) | `gitmap-v28 pull-all`                 | Terminal | **Power-user batch.** Windows: `Extended` REG_SZ → Shift+right-click only. macOS/Linux: visible but gated by `osascript display dialog` / `zenity`/`kdialog`/`xmessage` confirm. Forwards to `runPull` with `--all` injected. |
| `30_release/10_release`    | Release current          | `gitmap-v28 release`                       | Terminal | Re-tags `HEAD` at current `constants.Version`. Interactive prompts for missing notes. |
| `30_release/20_release_next` | Release next (bump minor) | `gitmap-v28 release --bump minor`        | Terminal | Uses `constants.FlagBumpDash` + `constants.BumpMinor` — no string literals in the entry. Patch / major variants intentionally omitted from the menu (rarely used; users type them). |
| `30_release/30_release_pull` | Release pull           | `gitmap-v28 release-pull`                  | Terminal | `git pull --ff-only` then `release`. Hard-fails on divergent history. |
| `30_release/40_release_pending` | Release pending     | `gitmap-v28 release-pending`               | Silent   | Prints commits since last release. |
| `30_release/50_list_releases` | List releases         | `gitmap-v28 list-releases`                 | Silent   | Single-repo view. (`--all-repos` deliberately omitted — would surprise the user clicking inside one folder.) |
| `30_release/60_list_versions` | List versions         | `gitmap-v28 list-versions`                 | Silent   | Single-repo `RepoVersionHistory` view. |
| `40_repos/10_go`           | Go projects              | `gitmap-v28 go-repos`                      | Silent   | Filters DB by `go.mod` detection. No `--json` (notification gets human text). |
| `40_repos/20_node`         | Node projects            | `gitmap-v28 node-repos`                    | Silent   | `package.json` detection. |
| `40_repos/30_react`        | React projects           | `gitmap-v28 react-repos`                   | Silent   | React dependency in `package.json`. |
| `40_repos/40_cpp`          | C++ projects             | `gitmap-v28 cpp-repos`                     | Silent   | CMakeLists / `.cpp` detection. |
| `40_repos/50_csharp`       | C# projects              | `gitmap-v28 csharp-repos`                  | Silent   | `.csproj` / `.sln` detection. |
| `50_visibility/10_public`  | Make public              | `gitmap-v28 make-public`                   | Terminal | Calls `gh repo edit --visibility public` (or `glab`). Interactive confirm. |
| `50_visibility/20_private` | Make private             | `gitmap-v28 make-private`                  | Terminal | Calls `gh repo edit --visibility private`. Interactive confirm. |
| `60_tools/10_fix_repo`     | Fix repo                 | `gitmap-v28 fix-repo`                      | Terminal | Rewrites stale `<base>-vN` tokens to current version. No `--strict` from the menu. |
| `60_tools/20_diff`         | Diff                     | `gitmap-v28 diff`                          | Terminal | Wraps `git diff` with gitmap-v28's pager. |
| `60_tools/30_history`      | History                  | `gitmap-v28 history`                       | Terminal | Local `CliInvocation` history; pages with `--limit 50` default. |
| `60_tools/40_update`       | Update                   | `gitmap-v28 update`                        | Terminal | Self-update — Terminal so the user sees the new-version banner. |
| `90_terminal`              | Open terminal here       | (no command — prefill prompt)          | Prefill  | Opens shell at folder, writes `gitmap-v28 ` literal so the user can finish typing. |
| `91_docs`                  | Docs                     | `gitmap-v28 docs`                          | Silent   | Prints help-dashboard URL. |

### 3.1 Execution-mode templates

| Mode       | Windows                                                                                                                | macOS                                                                                                                              | Linux                                                                                                       |
| ---------- | ---------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `Terminal` | `pwsh -NoExit -NoProfile -Command "Set-Location '%V'; & '<exe>' <args>"`                                               | `osascript -e 'tell application "Terminal" to do script "cd \"$1\" && <exe> <args>"' -e '… activate'`                              | `D="${1:-$PWD}"; cd "$D" && x-terminal-emulator -e sh -c "'<exe>' <args>; exec $SHELL" &`                   |
| `Silent`   | `pwsh -NoProfile -WindowStyle Hidden -Command "Set-Location '%V'; & '<exe>' <args> 2>&1 \| Out-String \| msg.exe * $_"` | `cd "$1" && OUT=$('<exe>' <args> 2>&1); osascript -e "display notification \"$(echo $OUT \| head -c 200)\" with title \"<label>\""` | `D="${1:-$PWD}"; cd "$D"; OUT=$('<exe>' <args> 2>&1); notify-send 'gitmap-v28' "$(echo $OUT \| head -c 200)"`   |
| `Prefill`  | `pwsh -NoExit -NoProfile -Command "Set-Location '%V'; Write-Host -NoNewline 'gitmap-v28 '"`                                | `osascript -e 'tell application "Terminal" to do script "cd \"$1\" && printf \"gitmap-v28 \""' -e '… activate'`                        | `x-terminal-emulator -e sh -c 'printf "gitmap-v28 "; exec $SHELL'`                                              |

### 3.2 Constant references (no magic strings)

Every value above resolves through a constant in
`gitmap-v28/constants/`:

| Used in entry                          | Constant                                                       |
| -------------------------------------- | -------------------------------------------------------------- |
| `release --bump minor`                 | `constants.CmdRelease` + `constants.FlagBumpDash` + `constants.BumpMinor` |
| `--scan-folder` (find-next, not used)  | `constants.FindNextFlagScanFolder`                             |
| `--json` (find-next, intentionally NOT used by menu) | `constants.FindNextFlagJSON`                       |
| `--all-repos` (list-releases, not used by menu) | `constants.FlagAllRepos`                              |
| Every `Cmd*`                            | `gitmap-v28/constants/constants_cli.go` + per-domain `constants_*.go` |

## 4. Implementation Layout

```
gitmap-v28/cmd/installctx.go            // entry point — runInstallCtx / runUninstallCtx
gitmap-v28/cmd/installctxentries.go     // []ctxEntry table (single source of truth)
gitmap-v28/cmd/installctxregistry.go    // reg add/delete helpers (table-driven)
gitmap-v28/cmd/installctxnotify.go      // probe BurntToast/msg.exe at install time
gitmap-v28/constants/constants_installctx.go  // all literals (tool name, key paths, MUIVerbs, flag names)
```

`ctxEntry` shape:

```go
type ctxEntry struct {
    KeyName  string   // "10_release_next" — numeric prefix preserves order
    MUIVerb  string   // "Release next (bump minor)"
    Args     []string // {"release", "--bump", "minor"}
    Mode     ctxMode  // Silent | Terminal | Prefill
    Category string   // "Release" — empty = top-level under gitmap-v28
}
```

The same slice drives:
- install (write keys),
- uninstall (delete the `gitmap-v28` subtree only, never neighbors),
- a unit test that asserts every entry references a real `Cmd*`
  constant from `constants_cli.go` (catches drift when commands are
  renamed).

### 4.1 Wire-up to the existing install dispatcher

`gitmap-v28/constants/constants_installctx.go`:

```go
const ToolCtx = "ctx"
```

`gitmap-v28/cmd/install.go::specialInstallHandler`:

```go
case constants.ToolCtx:
    return func(installOptions) { runInstallCtx() }
```

`gitmap-v28/cmd/uninstall.go` mirrors the `vscode-ctx` / `pwsh-ctx` branches.

### 4.2 Tool-table entry

Append to `constants_install.go`:

| Field           | Value                                                |
| --------------- | ---------------------------------------------------- |
| `ToolCtx`       | `"ctx"`                                              |
| description     | `"Add gitmap-v28 to Windows right-click context menu"`   |
| `allInstallable`| **omit** — `install all` should NOT install `ctx` (it changes Explorer chrome; users opt in explicitly). |

## 5. Acceptance Criteria

1. `gitmap-v28 install ctx` on Windows writes the full key tree under
   `HKCU\Software\Classes\Directory\{Background,}\shell\gitmap-v28` and
   prints `✓ gitmap-v28 context menu installed (X/X registry keys).`.
2. `gitmap-v28 uninstall ctx` deletes **only** the `gitmap-v28` subtree from
   both locations and prints a parallel summary. `vscode-ctx` /
   `pwsh-ctx` keys are untouched.
3. Right-clicking a folder background shows `gitmap-v28 ▸` with all five
   category submenus + the separator + Open-terminal + Docs entries.
4. Each `Terminal`-mode entry opens a non-closing `pwsh` window at the
   clicked folder and runs `gitmap-v28 <args>`.
5. Each `Silent`-mode entry surfaces output via the
   install-time-detected notifier (BurntToast → msg.exe → temp-log).
6. **Open terminal here** opens `pwsh` at the folder with a literal
   `gitmap-v28 ` prompt waiting for input (no command executed yet).
7. On non-Windows, both commands print the same OS-not-supported error
   the existing `vscode-ctx` handler prints, then exit 0 (parity).
8. A unit test (`installctxentries_test.go`) asserts every `ctxEntry.Args[0]`
   is one of the `Cmd*` constants in `constants_cli.go`.

## 6. Constraints

- All literals (registry paths, MUIVerbs, command templates, error
  strings) live in `constants_installctx.go` — no string literals in
  `installctx*.go`.
- Functions ≤15 lines; files ≤200 lines (split into the four files
  above).
- HKCU only — never write to `HKLM` (would require admin and affect
  other users).
- Uninstall must use `reg delete /f` scoped to the `gitmap-v28` key only —
  never wildcard the parent `shell` key.
- Use `cliexit.Reportf` for any error print (not bare `fmt.Fprintf`),
  per the `check-bare-stderr-err.sh` CI gate.

## 7. macOS / Linux — Implementation Notes

Because Finder Services and Linux file-manager menus do not support
arbitrary nested cascades, `flattenCtxMenu()` (in
`gitmap-v28/cmd/installctxflatten.go`) collapses each `Category ▸ Child`
into a single labelled `flatCtxEntry`:

```
gitmap-v28: Release — Release next (bump minor)
gitmap-v28: Tools — Fix repo
gitmap-v28: Open terminal here
```

### 7.1 macOS — `~/Library/Services/<slug>.workflow`

For every flat entry we generate one Automator Quick Action bundle
containing `Contents/Info.plist` (registers the service for
`public.folder`) and `Contents/document.wflow` (a single
`Run Shell Script` action). The shell script:

- **Terminal mode** → `osascript` to open Terminal.app at the folder
  and run `gitmap-v28 <args>`.
- **Silent mode** → run inline, surface output via
  `display notification`.
- **Prefill mode** → open Terminal.app and `printf "gitmap-v28 "` to leave
  a prompt.

After install, the user runs `pkill -KILL -u $USER cfprefsd` (or logs
out/in) to refresh Finder. No code-signing or notarization required —
`.workflow` bundles installed under the user's home are trusted.

### 7.2 Linux — Nautilus + Dolphin + Thunar

| Manager  | Path                                                  | Format                                   |
| -------- | ----------------------------------------------------- | ---------------------------------------- |
| Nautilus | `~/.local/share/nautilus/scripts/gitmap-v28/<label>`      | One executable shell script per entry; filename = menu label. |
| Dolphin  | `~/.local/share/kio/servicemenus/gitmap-ctx.desktop`  | Single `.desktop` with `Actions=` listing every entry under `X-KDE-Submenu=gitmap-v28` (KDE renders this as a real cascade). |
| Thunar   | `~/.config/Thunar/uca.xml`                            | Marker-delimited (`<!-- gitmap-ctx-begin --> … end -->`) `<action>` block; uninstall strips the block in place, leaving foreign actions intact. |

All three use `x-terminal-emulator` for Terminal/Prefill modes and
`notify-send` (with `echo` fallback) for Silent mode.

### 7.3 Future managers (out of scope)

Nemo, Caja and PCManFM use private formats with no shared schema; they
are not covered. Nautilus/Dolphin/Thunar already cover GNOME, KDE and
XFCE — roughly 95% of Linux desktop usage.

## 8. Cross-References

- Existing pattern: `gitmap-v28/cmd/installctxmenu.go` (`vscode-ctx`,
  `pwsh-ctx`) — copy the `runRegistryCommands` reporting style.
- Memory: `mem://features/install-ctx-menu`.
- CI gate: `.github/scripts/check-bare-stderr-err.sh` — must pass.
