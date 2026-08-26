# 10 — Last-Release Detection

**Version:** 1.0.0  
**Updated:** 2026-04-17

---

## Purpose

Define a standalone, reusable script that resolves and prints the
latest released version of the CLI tool. Used by `run.ps1`, the
`update` worker, CI pipelines, and ad-hoc manual checks.

---

## Why a Standalone Script

| Concern | Inline in `run.ps1` | Standalone `Get-LastRelease.ps1` |
|---------|---------------------|----------------------------------|
| Reuse from CI | ❌ Duplicates code | ✅ Single source of truth |
| Reuse from updater binary | ❌ Bash/PS bifurcation | ✅ Both call same script |
| `run.ps1` size | ❌ Bloats main script | ✅ Stays focused on build/deploy |
| Manual invocation | ❌ Awkward (must invoke main script) | ✅ Direct: `pwsh Get-LastRelease.ps1` |

---

## File Layout

```
<repo-root>/
├── run.ps1
├── run.sh
└── scripts/
    ├── Get-LastRelease.ps1     ← PowerShell version
    └── get-last-release.sh     ← Bash version (mirror)
```

Both scripts MUST produce identical output for the same inputs.

---

## Resolution Strategy

The script resolves "latest release" through a tiered fallback:

| Tier | Source | Used when |
|------|--------|-----------|
| 1 | `git tag --sort=-v:refname | head -1` (local) | Repo has tags locally |
| 2 | GitHub Releases API: `/repos/<owner>/<repo>/releases/latest` | Online check |
| 3 | `latest.json` from release assets | Tier 2 unavailable / rate-limited |
| 4 | Embedded version constant | Fully offline |

Each tier should print *which* tier produced the answer so the user
can interpret a stale result.

---

## PowerShell — `Get-LastRelease.ps1`

```powershell
[CmdletBinding()]
param(
    [string]$Owner = "<owner>",
    [string]$Repo  = "<repo>",
    [switch]$Json
)

$ErrorActionPreference = 'Stop'

function Get-FromGitTags {
    $tag = git tag --sort=-v:refname 2>$null | Select-Object -First 1
    if ($tag) { return @{ version = $tag; source = 'git-tag' } }
    return $null
}

function Get-FromGithubApi {
    $url = "https://api.github.com/repos/$Owner/$Repo/releases/latest"
    try {
        $r = Invoke-RestMethod -Uri $url -TimeoutSec 5 -Headers @{
            'User-Agent' = 'get-last-release'
        }
        return @{ version = $r.tag_name; source = 'github-api' }
    } catch { return $null }
}

$result = (Get-FromGitTags) ?? (Get-FromGithubApi)
if (-not $result) { throw "Could not resolve latest release" }

if ($Json) {
    $result | ConvertTo-Json -Compress
} else {
    Write-Host "$($result.version)  [$($result.source)]"
}
```

---

## Bash — `get-last-release.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

OWNER="${OWNER:-<owner>}"
REPO="${REPO:-<repo>}"
JSON="${1:-}"

from_git_tags() {
    local tag
    tag=$(git tag --sort=-v:refname 2>/dev/null | head -1 || true)
    [[ -n "$tag" ]] && { echo "{\"version\":\"$tag\",\"source\":\"git-tag\"}"; return 0; }
    return 1
}

from_github_api() {
    local url="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
    local v
    v=$(curl -fsSL --max-time 5 -H 'User-Agent: get-last-release' "$url" \
        | python3 -c "import json,sys;print(json.load(sys.stdin).get('tag_name',''))" 2>/dev/null || true)
    [[ -n "$v" ]] && { echo "{\"version\":\"$v\",\"source\":\"github-api\"}"; return 0; }
    return 1
}

result=$(from_git_tags || from_github_api || true)
[[ -z "$result" ]] && { echo "ERROR: could not resolve latest release" >&2; exit 1; }

if [[ "$JSON" == "--json" ]]; then
    echo "$result"
else
    echo "$result" | python3 -c "import json,sys;d=json.load(sys.stdin);print(f\"{d['version']}  [{d['source']}]\")"
fi
```

---

## Output Contract

Both scripts MUST honor:

| Mode | Output |
|------|--------|
| Default | `v1.2.3  [git-tag]` (human-readable, one line) |
| `--json` (Bash) / `-Json` (PS) | `{"version":"v1.2.3","source":"git-tag"}` (single line, parseable) |

Never print colored output by default — these scripts are commonly
piped into other commands.

---

## Constraints

- Network calls MUST have a 5-second timeout. The script is often
  used in tight loops (CI matrix, install scripts).
- The script MUST exit non-zero on total failure so callers can
  detect "no version available" without parsing.
- The `source` field is informational, not authoritative — callers
  that need a specific source should override the resolution order
  via env vars (e.g., `LASTRELEASE_SOURCE=github-api`).
- The script MUST work standalone — no dependencies on other repo
  scripts or build artifacts.

---

## Cross-References

- [04-build-scripts.md](04-build-scripts.md) §Version Verification — uses this script
- [01-self-update-overview.md](01-self-update-overview.md) §Skip-if-Current — calls this for the comparison
- [`../14-update/23-install-script-version-probe.md`](../14-update/23-install-script-version-probe.md) — install probe uses similar tiered fallback

---

*Last-release detection — v1.0.0 — 2026-04-17*
