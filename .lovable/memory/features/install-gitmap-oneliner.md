---
name: install-gitmap-oneliner
description: `gitmap install gitmap-oneliner` prints the canonical Win + macOS install one-liners with icons; URLs fixed, rendering dynamic via constants.Version + MsgInstallHint* reuse.
type: feature
---

# install gitmap-oneliner (v5.15.0+)

`gitmap install gitmap-oneliner` is a print-only special handler that
emits the canonical bootstrap one-liners for both platforms:

- 🪟 Windows · PowerShell — `irm .../install.ps1 | iex`
- 🐧 Linux / macOS — `curl -fsSL .../install.sh | sh`

## Wiring

- Tool name: `constants.ToolGitmapOneliner = "gitmap-oneliner"`
- Description registered in `InstallToolDescriptions` and listed in the
  Core `InstallToolCategories` group so `gitmap install --list` shows it.
- Dispatch: `specialInstallHandler` in `gitmap/cmd/install.go` returns
  `runInstallGitmapOneliner`, bypassing the detect/confirm/install pipe.
- Handler: `gitmap/cmd/installgitmaponeliner.go` reuses the existing
  `MsgInstallHintHeader/Windows/Unix` constants from
  `gitmap/constants/constants_release.go` — single source of truth for
  the install URLs.

## Contract

- **URLs are fixed** (canonical `alimtvnetwork/gitmap-v27` main branch).
  Never accept flags to override — that's what `self-install --version`
  is for.
- **Rendering is dynamic** — header uses `constants.Version`, icons +
  ordering are produced in Go, not stored as a baked literal block.
- **Both sections always print** regardless of host OS so the user can
  copy whichever applies.

Spec: `spec/01-app/109-install-gitmap-oneliner.md`.
