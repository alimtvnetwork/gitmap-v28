# Run Pattern (`-R` Flag)

Part of [PowerShell Build & Deploy Patterns](02-powershell-build-deploy.md).

## Forwarding Arguments

Use `ValueFromRemainingArguments` to capture all trailing arguments
after the `-R` switch and forward them to the built binary:

```powershell
[switch]$R,
[Parameter(ValueFromRemainingArguments=$true)]
[string[]]$RunArgs
```

## Path Resolution

Resolve relative paths to absolute before passing to the binary,
since `Start-Process` may run from a different working directory:

```powershell
foreach ($arg in $CliArgs) {
    if ($arg -match '^(\.\.[\\/]|\.[\\/]|\.\.?$)') {
        $path = Resolve-Path -LiteralPath $arg -ErrorAction SilentlyContinue
        if ($path) { $resolved += $path.Path }
        else { $resolved += [System.IO.Path]::GetFullPath((Join-Path $baseDir $arg)) }
    } else {
        $resolved += $arg
    }
}
```

## Default Behavior

If `-R` is used with no arguments, default to a sensible action
(e.g., process the parent folder of the repo).

## Context Logging

Before executing, print diagnostic info:

```
[RUN] Executing toolname
→ Runner CWD: D:\projects\my-tool
→ Command: toolname scan D:\projects
→ Scan target: D:\projects
```
