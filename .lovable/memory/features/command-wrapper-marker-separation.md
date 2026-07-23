---
name: Command Wrapper Marker Separation
description: `gitmap`/`gcd` command-wrapper marker and runtime sentinel are separate from PATH snippet marker/env; Windows installers and release ZIPs must ship/load the PowerShell command wrapper plus `gitmap.ps1` shim. v5.9.0+.
type: feature
---

# Command Wrapper Marker Separation (v5.3.0+; installer hardening v5.4.0+; shim v5.5.0+; release asset fix v5.7.0+; profile hardening v5.9.0+)

## Rule

Never use the PATH snippet marker (`# gitmap shell wrapper v2 - managed by ...`)
or `GITMAP_WRAPPER` as proof that the interactive `gitmap` shell function is
installed/active.

The actual command wrapper must use:

- Profile marker: `# gitmap command wrapper v1`
- Runtime sentinel: `GITMAP_COMMAND_WRAPPER=1`

## Root cause fixed

The PATH snippet and command wrapper both used `# gitmap shell wrapper v2`, and
the PATH snippet exported `GITMAP_WRAPPER=1`. That made `gitmap setup` skip
installing the real `function gitmap { ... }` / `gcd` block and made doctor/setup
verification report success even though PowerShell still resolved `gitmap` as the
exe. Result: `gitmap cd <repo>` printed a path but could not change directory.

## Prevention

- `completion.appendCDFunction` must only skip when `CDFuncMarker` is present.
- `isWrapperActive` must check `EnvGitmapCommandWrapper`, not `EnvGitmapWrapper`.
- Keep `GITMAP_WRAPPER` for legacy compatibility only; do not use it for active
  command-wrapper detection.
- Windows install/profile snippets must install and load the PowerShell
  `function gitmap` / `function gcd` wrapper. Adding PATH alone is not enough,
  because an executable can never `Set-Location` in the parent PowerShell session.
- Windows installs must also write `gitmap.ps1` beside `gitmap.exe`. PowerShell
  prefers the `.ps1` script over the `.exe` on PATH, and scripts run in-process,
  so this shim can call the exe then `Set-Location` even when the user's profile
  function has not been reloaded yet.
- Installer/setup must replace stale `# gitmap command wrapper v1` blocks when
  the wrapper body changes; marker presence alone is not proof the body is current.
- Release ZIPs and release-specific `install.ps1` must ship/install `gitmap.ps1`
  beside `gitmap.exe`. The v5.5.0 fix failed for release installs because the
  pipeline zipped only the exe, so users never received the shim.
- Windows installers must write the command wrapper to all standard current-user
  PowerShell profile files and load it into the installer session; writing only
  `$PROFILE` can leave new `powershell.exe` / `pwsh.exe` windows resolving the raw exe.