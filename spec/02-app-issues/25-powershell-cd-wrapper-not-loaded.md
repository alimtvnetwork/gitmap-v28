# `gitmap-v27 cd` still prints the path in PowerShell after setup

## Ticket

After the shell-wrapper work landed, some Windows users still see
`gitmap-v27 cd <repo>` print the resolved repository path instead of moving the
active PowerShell session into that directory.

## Symptoms

1. User runs `gitmap-v27 setup` and sees the completion / cd wrapper install step.
2. User opens PowerShell and runs `gitmap-v27 cd gitmap-v27`.
3. The command resolves the SQLite record correctly and prints the repo path.
4. The prompt stays in the original directory because no shell function takes
   over the result.

## Root Cause

The wrapper logic itself was correct after issue 24, but the PowerShell profile
installation path was not reliable:

- The installer read `os.Getenv("PROFILE")`, but `$PROFILE` is a PowerShell
  variable, not a normal exported environment variable, so Go usually received
  an empty string.
- The fallback targeted only one guessed profile path, which missed common
  Windows PowerShell / cross-host profile locations.
- Profile parent directories were not created before appending content, so new
  profile targets could silently fail to materialize.

That meant `gitmap-v27 setup` often wrote the managed `gitmap-v27` / `gcd` wrapper to a
file the active shell never loaded. When that happened, PowerShell executed the
real binary directly, and the binary could only print the destination path.

## Fix

Update PowerShell shell integration installation to:

1. Resolve all relevant PowerShell profile targets instead of assuming `PROFILE`.
2. Prefer engine-reported `CurrentUserAllHosts` profile paths from both
   `powershell` and `pwsh` when available.
3. Fall back to the standard all-hosts profile files for Windows PowerShell and
   PowerShell 7 when probing is unavailable.
4. Install the `gitmap-v27` / `gcd` wrapper and completion source line into every
   resolved PowerShell profile target.
5. Create profile parent directories before writing.

## Prevention

1. Do not treat shell-scoped variables like `$PROFILE` as process environment
   variables.
2. Test shell integrations at the profile-loading boundary, not only at the
   wrapper-function boundary.
3. Cover missing profile directories and multi-profile PowerShell installs in
   automated tests.
4. Keep the direct binary contract (`stdout` path) separate from the shell
   integration contract (`cd` / `Set-Location` in the parent shell).

## Related

- `spec/02-app-issues/24-cd-command-does-not-change-shell-directory.md`
- `gitmap-v27/completion/install.go`
- `gitmap-v27/completion/cdfunction.go`