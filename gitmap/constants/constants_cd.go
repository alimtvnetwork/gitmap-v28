package constants

// gitmap:cmd top-level
// CD CLI commands.
const (
	CmdCD      = "cd"
	CmdCDAlias = "go"
)

// gitmap:cmd top-level
// CD subcommands.
const (
	CmdCDRepos        = "repos"         // gitmap:cmd skip
	CmdCDSetDefault   = "set-default"   // gitmap:cmd skip
	CmdCDClearDefault = "clear-default" // gitmap:cmd skip
)

// CD help text.
const HelpCD = "  cd (go) <name>      Navigate to a tracked repo directory"

// CD file.
const CDDefaultsFile = "cd-defaults.json"

// CD messages.
const (
	MsgCDMultipleHeader    = "Multiple locations found for \"%s\":\n"
	MsgCDMultipleRowFmt    = "  %d  %s\n"
	MsgCDPickPrompt        = "\nPick [1-%d]: "
	MsgCDReposHeader       = "TRACKED REPOS\n"
	MsgCDReposRowFmt       = "  %d  %s\n"
	MsgCDDefaultSet        = "Default set for %s: %s\n"
	MsgCDDefaultCleared    = "Default cleared for %s\n"
	ErrCDUsage             = "usage: gitmap cd <repo-name|repos> [--group <name>] [--pick] [<inner-command> ...]\n"
	ErrCDChdirFmt          = "cd: cannot enter %s: %v\n"
	ErrCDNotFound          = "no repo found matching '%s'\n"
	ErrCDInvalidPick       = "invalid selection\n"
	ErrCDSetDefaultUsage   = "usage: gitmap cd set-default <name> <path>\n"
	ErrCDClearDefaultUsage = "usage: gitmap cd clear-default <name>\n"
	ErrCDDefaultNotFound   = "no default set for '%s'\n"
)

// CD flag descriptions.
const (
	FlagDescCDGroup = "Filter repos list by group"
	FlagDescCDPick  = "Force interactive picker even if a default is set"
)

// CD shell wrapper functions — installed by setup/completion.
const CDFuncMarker = "# gitmap command wrapper v1"

// CDFuncMarkerEnd closes the managed command-wrapper block for rewrites.
const CDFuncMarkerEnd = "# gitmap command wrapper v1 end"

// CDFuncMarkerEndLegacy is the v5.7-style close marker accepted for rewrites.
const CDFuncMarkerEndLegacy = "# end gitmap command wrapper v1"

// Legacy markers identify older/broken command-wrapper blocks that must be
// replaced instead of causing setup to skip installation.
const CDFuncMarkerLegacy = "# gitmap shell wrapper v2"

// CD shell wrapper env var — set by wrappers so the binary can detect them.
const (
	EnvGitmapWrapper        = "GITMAP_WRAPPER"
	EnvGitmapCommandWrapper = "GITMAP_COMMAND_WRAPPER"
	EnvGitmapWrapperVal     = "1"
)

// Shell-handoff sentinel file env var. Set by the wrapper function to a
// writable temp file path; commands like clone-next, as, cd write the
// destination directory to that file, and the wrapper cds to it.
//
// Spec: spec/04-generic-cli/21-post-install-shell-activation/01-contract.md
const EnvGitmapHandoffFile = "GITMAP_HANDOFF_FILE"

// Shell-handoff error format.
const ErrShellHandoffWriteFmt = "  ⚠ Could not write shell-handoff file %s: %v\n"

// CD wrapper verification messages.
const (
	WarnWrapperInstallFmt = "  ⚠ Could not auto-install command wrapper: %v\n"
	MsgWrapperNotLoaded   = "  %s! Command wrapper not active%s — 'gitmap cd' printed the path but cannot change your directory yet.\n    Wrapper install/update was attempted; run: %s. $PROFILE%s (PowerShell) or %ssource ~/.bashrc%s / %ssource ~/.zshrc%s, then retry.\n"
	MsgWrapperVerifyOK    = "  %s✓%s Command wrapper is active (gitmap resolves as a function)\n"
	MsgWrapperVerifyTip   = "\n  %s→%s To activate: restart your terminal or reload your profile\n    PowerShell: %s. $PROFILE%s | Bash: %ssource ~/.bashrc%s | Zsh: %ssource ~/.zshrc%s\n"
)

// CDFuncBash installs gitmap and gcd wrappers for Bash.
const CDFuncBash = `gcd() {
  local dest status
  dest="$(GITMAP_COMMAND_WRAPPER=1 GITMAP_WRAPPER=1 command gitmap cd "$@")"
  status=$?
  if [ $status -ne 0 ]; then
    return $status
  fi
  if [ -n "$dest" ] && [ -d "$dest" ]; then
    builtin cd "$dest" || return $?
  fi
}

gitmap() {
  if [ "$1" = "cd" ] || [ "$1" = "go" ]; then
    local dest status
    dest="$(GITMAP_COMMAND_WRAPPER=1 GITMAP_WRAPPER=1 command gitmap "$@")"
    status=$?
    if [ $status -ne 0 ]; then
      return $status
    fi
    if [ -n "$dest" ] && [ -d "$dest" ]; then
      builtin cd "$dest" || return $?
    fi
    return 0
  fi
  local handoff status
  handoff="$(mktemp -t gitmap-handoff.XXXXXX 2>/dev/null)" || handoff=""
  if [ -n "$handoff" ]; then
    GITMAP_HANDOFF_FILE="$handoff" GITMAP_COMMAND_WRAPPER=1 GITMAP_WRAPPER=1 command gitmap "$@"
    status=$?
    if [ -s "$handoff" ]; then
      local target
      target="$(cat "$handoff")"
      if [ -n "$target" ] && [ -d "$target" ]; then
        builtin cd "$target" || true
      fi
    fi
    rm -f "$handoff"
    return $status
  fi
  command gitmap "$@"
}

gitm() { gitmap "$@"; }
# gitmap command wrapper v1 end`

// CDFuncZsh installs gitmap and gcd wrappers for Zsh.
const CDFuncZsh = `gcd() {
  local dest
  local status
  dest="$(GITMAP_COMMAND_WRAPPER=1 GITMAP_WRAPPER=1 command gitmap cd "$@")"
  status=$?
  if (( status != 0 )); then
    return $status
  fi
  if [[ -n "$dest" && -d "$dest" ]]; then
    builtin cd "$dest" || return $?
  fi
}

gitmap() {
  if [[ "$1" == "cd" || "$1" == "go" ]]; then
    local dest
    local status
    dest="$(GITMAP_COMMAND_WRAPPER=1 GITMAP_WRAPPER=1 command gitmap "$@")"
    status=$?
    if (( status != 0 )); then
      return $status
    fi
    if [[ -n "$dest" && -d "$dest" ]]; then
      builtin cd "$dest" || return $?
    fi
    return 0
  fi
  local handoff status
  handoff="$(mktemp -t gitmap-handoff.XXXXXX 2>/dev/null)" || handoff=""
  if [[ -n "$handoff" ]]; then
    GITMAP_HANDOFF_FILE="$handoff" GITMAP_COMMAND_WRAPPER=1 GITMAP_WRAPPER=1 command gitmap "$@"
    status=$?
    if [[ -s "$handoff" ]]; then
      local target
      target="$(cat "$handoff")"
      if [[ -n "$target" && -d "$target" ]]; then
        builtin cd "$target" || true
      fi
    fi
    rm -f "$handoff"
    return $status
  fi
  command gitmap "$@"
}

gitm() { gitmap "$@"; }
# gitmap command wrapper v1 end`

// CDFuncPowerShell installs gitmap and gcd wrappers for PowerShell.
const CDFuncPowerShell = `function gcd {
  $real = Get-GitmapCommand
  if (-not $real) {
    Write-Error "gitmap executable not found"
    return
  }
  $env:GITMAP_WRAPPER = "1"
  $env:GITMAP_COMMAND_WRAPPER = "1"
  $dest = [string](& $real cd @args | Out-String)
  if ($LASTEXITCODE -ne 0) {
    return
  }
  $dest = $dest.Trim()
  if ($dest -and (Test-Path -LiteralPath ([string]$dest))) {
    Set-Location -LiteralPath ([string]$dest)
  }
}

function Get-GitmapCommand {
  $cmd = Get-Command gitmap.exe -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($cmd) {
    return $cmd.Source
  }
  $cmd = Get-Command gitmap -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($cmd) {
    return $cmd.Source
  }
  return $null
}

function gitmap {
  $real = Get-GitmapCommand
  if (-not $real) {
    Write-Error "gitmap executable not found"
    return
  }
  if ($args.Count -gt 0 -and ($args[0] -eq 'cd' -or $args[0] -eq 'go')) {
    $env:GITMAP_WRAPPER = "1"
    $env:GITMAP_COMMAND_WRAPPER = "1"
    $dest = [string](& $real @args | Out-String)
    if ($LASTEXITCODE -ne 0) {
      return
    }
    $dest = $dest.Trim()
    if ($dest -and (Test-Path -LiteralPath ([string]$dest))) {
      Set-Location -LiteralPath ([string]$dest)
    }
    return
  }
  $handoff = [System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), "gitmap-handoff-$([System.Guid]::NewGuid().ToString('N')).txt")
  try {
    $env:GITMAP_HANDOFF_FILE = $handoff
    $env:GITMAP_WRAPPER = "1"
    $env:GITMAP_COMMAND_WRAPPER = "1"
    & $real @args
    if ((Test-Path -LiteralPath $handoff) -and ((Get-Item -LiteralPath $handoff).Length -gt 0)) {
      $target = [string](Get-Content -LiteralPath $handoff -Raw)
      $target = $target.Trim()
      if ($target -and (Test-Path -LiteralPath ([string]$target))) {
        Set-Location -LiteralPath ([string]$target)
      }
    }
  }
  finally {
    Remove-Item -LiteralPath $handoff -ErrorAction SilentlyContinue
    Remove-Item Env:\GITMAP_HANDOFF_FILE -ErrorAction SilentlyContinue
  }
}

function gitm { gitmap @args }
# gitmap command wrapper v1 end`

// CD function messages.
const (
	MsgCDFuncInstalled = "Installed 'gitmap'/'gcd' shell wrappers — restart your terminal or source your profile\n"
	MsgCDFuncAlready   = "Shell wrappers for 'gitmap'/'gcd' already installed\n"
)
