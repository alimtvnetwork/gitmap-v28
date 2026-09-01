# `gitmap install gitmap-oneliner`

## Purpose

Print the canonical one-liners that bootstrap the gitmap installer on
Windows (PowerShell) and macOS/Linux (bash). This gives users an
instantly copy-pasteable snippet without leaving the terminal — no need
to open the README or `/install-gitmap` page.

## Synopsis

```
gitmap install gitmap-oneliner
```

No flags. The output is identical on every platform — both sections are
always printed so the user can copy whichever applies to their target
machine.

## Output

```
  📦 Install gitmap v<current-version>
  ─────────────────────

  🪟  Windows · PowerShell
     irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.ps1 | iex

  🐧  Linux / macOS
     curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.sh | sh
```

## Implementation contract

- **URLs are fixed.** They always point at the canonical
  `alimtvnetwork/gitmap-v28` repo on the `main` branch. The same URLs
  are already used by `MsgInstallHintWindows` / `MsgInstallHintUnix` in
  `gitmap/constants/constants_release.go` — this command reuses those
  constants verbatim so there is exactly one source of truth.
- **Rendering is dynamic.** The header, icons, and section ordering are
  produced by Go (`runInstallGitmapOneliner` in
  `gitmap/cmd/installgitmaponeliner.go`) using the current
  `constants.Version`, not a baked-in literal.
- **Dispatch** wires through `specialInstallHandler` in
  `gitmap/cmd/install.go` so it bypasses the generic
  detect → confirm → install pipeline.

## Related

- `spec/01-app/108-cross-platform-install-update.md` — full matrix.
- `gitmap/scripts/install.ps1` / `install.sh` — the scripts the
  one-liners actually fetch.
