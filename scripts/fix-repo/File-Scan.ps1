<#
.SYNOPSIS  File-traversal + binary-detection helpers for fix-repo.ps1.
#>

$ErrorActionPreference = 'Stop'

$Script:BinaryExtensions = @(
    '.png','.jpg','.jpeg','.gif','.webp','.ico','.pdf',
    '.zip','.tar','.gz','.tgz','.bz2','.xz','.7z','.rar',
    '.woff','.woff2','.ttf','.otf','.eot',
    '.mp3','.mp4','.mov','.wav','.ogg','.webm',
    '.class','.jar','.so','.dylib','.dll','.exe','.pyc'
)

$Script:MaxFileBytes = 5MB

function Get-TrackedFiles {
    param([string]$RepoRoot)
    Push-Location $RepoRoot
    try {
        $list = & git ls-files 2>$null
        if ($LASTEXITCODE -ne 0) { return @() }
        return $list
    } finally { Pop-Location }
}

function Test-IsBinaryExtension {
    param([string]$Path)
    $ext = [System.IO.Path]::GetExtension($Path).ToLowerInvariant()
    return $Script:BinaryExtensions -contains $ext
}

function Test-HasNullByte {
    param([string]$FullPath)
    $stream = [System.IO.File]::OpenRead($FullPath)
    try {
        $buf = New-Object byte[] 8192
        $read = $stream.Read($buf, 0, $buf.Length)
        for ($i = 0; $i -lt $read; $i++) {
            if ($buf[$i] -eq 0) { return $true }
        }
        return $false
    } finally { $stream.Dispose() }
}

function Test-IsSkippablePath {
    param([string]$FullPath)
    $info = Get-Item -LiteralPath $FullPath -Force -ErrorAction SilentlyContinue
    if (-not $info) { return $true }
    if ($info.Attributes -band [System.IO.FileAttributes]::ReparsePoint) { return $true }
    if ($info.Length -gt $Script:MaxFileBytes) { return $true }
    return $false
}
