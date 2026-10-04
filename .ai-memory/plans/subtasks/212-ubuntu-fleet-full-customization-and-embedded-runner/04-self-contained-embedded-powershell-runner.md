# Subtask 212.4: Self-Contained Master Embedded PowerShell Runner

- **Parent Plan:** [212-ubuntu-fleet-full-customization-and-embedded-runner.md](../../pending/212-ubuntu-fleet-full-customization-and-embedded-runner.md)
- **Spec Reference:** [02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md](../../../../02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `scripts/master-embedded-ubuntu-runner.ps1`, `scripts/`

---

## 1. Context & Objective

Managing remote Linux workstations across an enterprise fleet often suffers from script distribution friction:
1. Multi-file deployment requires copying auxiliary `.sh` scripts via SCP or SFTP before execution.
2. Windows Git checkouts frequently introduce Windows CRLF line endings into bash scripts, triggering fatal syntax errors (`\r: command not found`, `syntax error near unexpected token $'\r'`).
3. Out-of-sync script versions between Windows orchestrators and remote nodes cause non-reproducible deployments.

The objective is to architect and implement `scripts/master-embedded-ubuntu-runner.ps1`, a completely self-contained PowerShell automation runner that encapsulates all required Linux provisioning and configuration routines as embedded multi-line PowerShell string templates (`@' ... '@`). The runner pipes each payload over SSH using `tr -d '\r' | bash -s`, completely eliminating CRLF issues and external script file dependencies.

---

## 2. Technical Architecture: Embedded Template Streaming

### 2.1 The CRLF Elimination Streaming Pattern

Rather than transferring physical files to `/tmp` on remote nodes, the PowerShell orchestrator directly pipes embedded bash templates via standard input:

```powershell
function Invoke-RemoteBashTemplate {
    param (
        [string]$Node = "ubuntu-fleet-01",
        [string]$User = "a",
        [string]$TemplateName,
        [string]$BashScriptContent,
        [string[]]$ScriptArgs = @()
    )

    Write-Host "[RUNNER] Executing stage: $TemplateName on $User@$Node..." -ForegroundColor Cyan
    
    # tr -d '\r' strips carriage returns on the fly before passing to bash
    $RemoteCommand = "tr -d '\r' | bash -s -- $($ScriptArgs -join ' ')"
    
    $ProcessInfo = New-Object System.Diagnostics.ProcessStartInfo
    $ProcessInfo.FileName = "ssh.exe"
    $ProcessInfo.Arguments = "-o BatchMode=yes -o StrictHostKeyChecking=accept-new $User@$Node `"$RemoteCommand`""
    $ProcessInfo.RedirectStandardInput = $true
    $ProcessInfo.RedirectStandardOutput = $true
    $ProcessInfo.RedirectStandardError = $true
    $ProcessInfo.UseShellExecute = $false
    $ProcessInfo.CreateNoWindow = $true

    $Process = [System.Diagnostics.Process]::Start($ProcessInfo)
    $Process.StandardInput.Write($BashScriptContent)
    $Process.StandardInput.Close()

    $StdOut = $Process.StandardOutput.ReadToEnd()
    $StdErr = $Process.StandardError.ReadToEnd()
    $Process.WaitForExit()

    if ($Process.ExitCode -ne 0) {
        Write-Error "[FAIL] Stage $TemplateName failed (Code: $($Process.ExitCode)): $StdErr"
        return $false
    }
    
    Write-Host $StdOut
    Write-Host "[SUCCESS] Stage $TemplateName completed successfully." -ForegroundColor Green
    return $true
}
```

---

## 3. Embedded Script Modules

`scripts/master-embedded-ubuntu-runner.ps1` embeds five core operational modules:

### Module 1: GNOME Wallpaper Synchronization (`$Bash_SetWallpaper`)
- Locates active user D-Bus session bus socket at `unix:path=/run/user/1000/bus`.
- Ensures default wallpaper directory `$HOME/Pictures/Wallpapers/` exists.
- Executes `gsettings set org.gnome.desktop.background picture-uri` and `picture-uri-dark`.
- Sets scaling style to `picture-options "zoom"`.

### Module 2: Remote GUI Application Launching (`$Bash_RemoteGuiLauncher`)
- Configures environment wrappers for Wayland/XWayland graphical sessions.
- Deploys or executes `systemd-run --user` wrappers to spawn graphical apps without display authorization errors:
  ```bash
  systemd-run --user /usr/local/bin/antigravity $HOME/git-work/gitmap
  ```
- Verifies process execution and reports PID via `systemctl --user status`.

### Module 3: VMware Shared Folders Automount & Permissions (`$Bash_VerifyVmwareShared`)
- Validates `open-vm-tools` and `open-vm-tools-desktop` package installation.
- Tests `/mnt/hgfs` mount status.
- Configures `/etc/systemd/system/mnt-hgfs.automount` and `/etc/systemd/system/mnt-hgfs.mount` with `allow_other,uid=1000,gid=1000`.
- Verifies write access for non-root user `a`.

### Module 4: Antigravity IDE Launcher & Account Switcher Check (`$Bash_CheckAntigravitySetup`)
- Verifies that `$HOME/.local/bin/antigravity` or `/usr/local/bin/antigravity` ELF binary wrapper resolves without recursion bugs.
- Checks `~/.gemini/oauth_creds.json` permissions (`600`).
- Validates `agm` binary availability and account synchronization.

### Module 5: Path Normalization & Integrity Check Receiver (`$Bash_NormalizeReceiver`)
- Validates SQLite database tables in `$HOME/.gemini/antigravity/conversation_summaries.db`.
- Executes SQL replacement for remaining Windows paths.
- Validates JSON formatting with `jq`.

---

## 4. CLI Routing & Parameter Dispatch

The runner provides flexible switches for targeted or end-to-end execution:

```powershell
param (
    [string]$TargetHost = "ubuntu-fleet-01",
    [string]$User = "a",
    [switch]$All,
    [switch]$Wallpaper,
    [switch]$GuiLaunch,
    [switch]$MountHgfs,
    [switch]$SyncBrain,
    [switch]$CheckHealth,
    [string]$WallpaperPath = "$HOME/Pictures/Wallpapers/default.png"
)
```

- `.\scripts\master-embedded-ubuntu-runner.ps1 -All` $\rightarrow$ Executes the complete 5-stage pipeline sequentially.
- `.\scripts\master-embedded-ubuntu-runner.ps1 -Wallpaper -WallpaperPath <path>` $\rightarrow$ Configures wallpaper only.
- `.\scripts\master-embedded-ubuntu-runner.ps1 -MountHgfs` $\rightarrow$ Checks and enables VMware shared folders.
- `.\scripts\master-embedded-ubuntu-runner.ps1 -CheckHealth` $\rightarrow$ Executes telemetry, mount, and D-Bus sanity checks.

---

## 5. Remediation Checklist

- [ ] Create `scripts/master-embedded-ubuntu-runner.ps1` with strict parameter validation and help documentation.
- [ ] Embed `$Bash_SetWallpaper` template with D-Bus socket export and dual light/dark mode `gsettings` keys.
- [ ] Embed `$Bash_RemoteGuiLauncher` template supporting `systemd-run --user` launching of Antigravity, text editors, and terminals.
- [ ] Embed `$Bash_VerifyVmwareShared` template verifying `/mnt/hgfs` mount and user `a` read/write permissions.
- [ ] Embed `$Bash_CheckAntigravitySetup` template validating non-recursive launcher and credentials.
- [ ] Implement `tr -d '\r' | bash -s` execution pipeline with error capture and colored console output.
- [ ] Add execution telemetry and timing benchmarks per stage.
- [ ] Document complete usage syntax in [02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md](../../../../02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md).
