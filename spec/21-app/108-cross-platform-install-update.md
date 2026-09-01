# Cross-Platform Install / Update / Uninstall

This spec is the **single source of truth** for installing, updating, and
removing `gitmap` on Windows, macOS, and Linux. The README, the in-app
docs page (`/install-gitmap`), `gitmap self-install --help`, and the
`gitmap update --help` text all reference (and must stay aligned with)
the matrices below.

Companion specs:

- `spec/01-app/94-install-script.md` — internals of `install.ps1` / `install.sh`.
- `spec/07-generic-release/03-install-scripts.md` — release-side packaging.
- `spec/07-generic-release/09-generic-install-script-behavior.md` — version-resolution contract.
- `spec/01-app/90-self-install-uninstall.md` — `self-install` / `self-uninstall` semantics.

## Why a unified reference

Users land on three different surfaces (README, web docs, `--help`) and
have to combine snippets across them to get a working install. This spec
collapses all three flows into one matrix per platform, with the same
copy-paste-safe one-liners everywhere.

## Platform matrix

| Action            | Windows (PowerShell)                                                                                                                          | macOS / Linux (bash / zsh)                                                                                                                |
|-------------------|------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| Install (default) | `irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.ps1 \| iex`                                         | `curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.sh \| sh`                                |
| Install (prompt)  | `irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.ps1 \| iex`                                                  | `curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.sh \| bash`                                       |
| Install (pinned)  | `$ver='vX.Y.Z'; & ([scriptblock]::Create((irm .../install.ps1))) -Version $ver -NoDiscovery`                                                   | `curl -fsSL .../install.sh \| bash -s -- --version vX.Y.Z --no-discovery`                                                                  |
| Update (in-place) | `gitmap update`                                                                                                                                | `gitmap update`                                                                                                                            |
| Update (pinned)   | `gitmap self-install --version vX.Y.Z --yes`                                                                                                   | `gitmap self-install --version vX.Y.Z --yes`                                                                                               |
| Uninstall         | `irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/uninstall-quick.ps1 \| iex`                                                | `curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/uninstall-quick.sh \| bash`                                     |
| Doctor / verify   | `gitmap doctor`                                                                                                                                | `gitmap doctor`                                                                                                                            |

The pinned commands assume the contract from
`spec/07-generic-release/09-generic-install-script-behavior.md`: with
`--version vX.Y.Z` (or `-Version`) the installer **never** falls back —
missing tag exits 1.

## Install resolution flow

Both installers follow the same algorithm (kept in lockstep by
`gitmap/clonenext/remoteupdate.go` and the install scripts):

1. **Strict mode** — if `--version <tag>` is supplied, fetch that tag
   verbatim. Missing → `exit 1`. No `latest`, no sibling probe.
2. **Discovery mode** — when no tag is supplied, run a 20-parallel
   `-v<N+i>` HEAD probe against sibling repos. Highest hit wins.
3. **Releases fallback** — if no sibling probe hits, request
   `releases/latest` from the canonical repo.
4. **HEAD fallback** — if no release exists yet, build from the default
   branch HEAD as a last resort.

Discovery mode is what `install-quick.{ps1,sh}` and the bare `install.*`
one-liners use. Strict mode is for CI and reproducible installs.

## Shell PATH activation

After installing the binary, the script writes a marker-block PATH
snippet into the user's shell profile(s). The exact targets depend on
`--shell-mode`:

| Mode value          | Profiles touched                                                          |
|---------------------|---------------------------------------------------------------------------|
| `auto` (default)    | Detected current shell + `~/.profile` on Unix; PowerShell profile on Win  |
| `both`              | zsh + bash + `~/.profile` + fish (if present) + pwsh                      |
| `zsh` / `bash` / `pwsh` / `fish` | Only that shell family                                          |
| `<a>+<b>` combos    | Strict union of listed families (no `~/.profile`, no auto-detect)         |

Snippet templates live in `gitmap/constants/constants_pathsnippet.go`
so `install.sh`, `install.ps1`, and `gitmap setup print-path-snippet`
emit byte-identical bytes.

## Update lifecycle

`gitmap update` handles three scenarios in this order:

1. **Linked source repo present** — `git pull` + `go build` + redeploy.
2. **No source repo, `gitmap-updater` installed** — delegate to the
   updater (downloads the matching release asset).
3. **No source repo, no updater** — print the four-option fallback
   panel documented in `gitmap/helptext/update.md` (re-install
   one-liner, clone+build, manual download, `--repo-path`).

The Phase 3 cleanup handoff always writes a structured durable log to
`<TMP>/gitmap-update-handoff-YYYYMMDD.log` so failed updates remain
recoverable on Windows where stdout can be swallowed by the launcher.

## Verification (post-install)

After any install or update path, the canonical verification flow is:

```
gitmap version          # confirms binary is on PATH and reports build
gitmap doctor           # checks PATH, profile snippets, deploy folder, DB
gitmap setup print-path-snippet   # prints the exact bytes the installer wrote
```

`gitmap doctor` is platform-aware and surfaces the exact remediation
command for any failed check (e.g. "PATH missing — re-run
`gitmap self-install --shell-mode <detected>`").

## Uninstall lifecycle

`uninstall-quick.{ps1,sh}` first attempts the canonical
`gitmap self-uninstall`. If `gitmap` is no longer on PATH, the script
falls back to a manual sweep that removes:

- The deploy folder (`D:\gitmap` on Windows; `~/.local/bin/gitmap` on Unix).
- PATH marker blocks from every profile under `--shell-mode`.
- (Optional) `%APPDATA%\gitmap` / `~/.config/gitmap` user data.

Useful flags (both scripts honor the same semantics):

| Flag             | Effect                                                            |
|------------------|-------------------------------------------------------------------|
| `--yes` / `-Yes` | Skip the "delete user data?" prompt and assume yes                |
| `--keep-data`    | Always keep user data even when `--yes` is set                    |
| `--dir <path>`   | Override the auto-detected deploy root                            |

## Out of scope

- Container images / Homebrew tap / Scoop bucket — tracked separately
  in `spec/07-generic-release/05-release-assets.md`.
- Network-restricted installs (offline tarball) — covered by the
  pinned-version snippet in
  `spec/07-generic-release/08-pinned-version-install-snippet.md`.
