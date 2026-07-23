---
name: ssh-view-copy-create
description: `gitmap ssh` adds `view`/`v`, `copy`/`cp`, `create` subcommands (v5.21.0+); copy uses clip/pbcopy/wl-copy/xclip/xsel; soft-fails when no tool present.
type: feature
---

# `gitmap ssh` view / copy / create (v5.21.0+)

- `gitmap ssh view <key>` (aliases `v`, existing `cat`) — prints public key.
- `gitmap ssh copy <key>` (alias `cp`) — prints public key + pushes to OS
  clipboard via `resolveClipboardTool()`:
  - Windows → `clip`
  - macOS → `pbcopy`
  - Linux → `wl-copy` → `xclip -selection clipboard` → `xsel --clipboard --input`
- `gitmap ssh create` — explicit alias for default `gitmap ssh` (generate).

Soft-fail contract: when no clipboard tool is on PATH, the key is still
printed to stdout and `MsgSSHCopyFallback` is emitted to stderr. The
command exits 0. Tool errors (e.g. clip.exe failed) print
`ErrSSHClipboard` to stderr but the command still exits 0 because the
key already reached stdout.

Files: `gitmap/cmd/sshcopy.go`, `gitmap/cmd/ssh.go`,
`gitmap/constants/constants_ssh.go`, `gitmap/helptext/ssh.md`.
