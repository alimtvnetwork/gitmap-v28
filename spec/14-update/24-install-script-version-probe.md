# 24 — Install Script Version Probe (Latest-Repo Hand-Off)

**Version:** 1.0.0
**Updated:** 2026-04-17
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Purpose

Define a generic, language-agnostic prologue that every CLI's
`install.ps1` and `install.sh` MUST run before doing any local install
work. The prologue probes for a newer numbered fork of the same
repository (e.g. `coding-guidelines-v24` → `coding-guidelines-v24`) and,
if a higher version exists, **hands off** to that repo's installer
instead of installing the older copy.

This guarantees that a user who runs an old `irm | iex` one-liner
**always ends up with the latest spec/CLI**, with no maintainer action.

---

## When This Runs

```
$ irm https://raw.githubusercontent.com/<owner>/<base>-v<N>/main/install.ps1 | iex
       │
       └── PROBE PROLOGUE (this spec)
              │
              ├── newer Vk found  →  irm https://raw.../<base>-v<k>/main/install.ps1 | iex
              │                        └── exit 0   (newer installer takes over)
              │
              └── nothing newer   →  fall through to local installer (existing logic)
```

**Mandatory ordering:** the probe is the FIRST thing the script does
after `$ErrorActionPreference = "Stop"` / `set -euo pipefail` and the
logging helpers. It must run **before** the banner, before reading
`install-config.json`, before any download.

---

## Self-Identification (Owner / Base / Current Version)

The script must detect `(owner, base, currentVersion)` from its own
invocation context. Detection is a 4-step ladder; the first match wins.

| # | Source | Example produces |
|---|--------|------------------|
| 1 | URL captured from `$MyInvocation.MyCommand.Definition` (PowerShell) or the script's own location resolved via `BASH_SOURCE` / `$0` (sh). When invoked via `irm \| iex`, PowerShell exposes the source URL inside `$MyInvocation`. | `https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.ps1` |
| 2 | Environment variables — `INSTALL_PROBE_OWNER`, `INSTALL_PROBE_BASE`, `INSTALL_PROBE_VERSION`. Used for local dev or CI. | — |
| 3 | Constants embedded at the top of the script: `$ProbeOwner`, `$ProbeBase`, `$ProbeVersion`. Maintainer fills these when forking. | — |
| 4 | Hard fail-open — if nothing resolves, log `"version probe disabled (no self-identity)"` and skip the probe. **Never block the local install.** | — |

### Required regex

```
^https?://[^/]+/(?<owner>[^/]+)/(?<base>[A-Za-z0-9._-]+?)-v(?<ver>\d+)/[^/]+/install\.(ps1|sh)(\?.*)?$
```

If the URL does not match (e.g. a fork without the `-v<N>` suffix),
treat as "no self-identity" and skip.

---

## Probe Algorithm

```
range   = currentVersion+1 .. currentVersion+20
url(N)  = https://raw.githubusercontent.com/<owner>/<base>-v<N>/main/install.<ext>
timeout = 2 seconds per request
method  = HEAD (fall back to GET with Range: bytes=0-0 if HEAD is blocked)
parallel = ALL 20 requests issued concurrently
expected success codes = 200, 301, 302
```

### Decision

1. Wait up to **`timeout × 2 = 4 s`** total for results to settle.
2. From the set of responding versions, pick the **highest** `N` (sort descending).
3. If that `N > currentVersion`, **hand off** (next section).
4. If no responses, log and continue to local install.

### Probe ordering optimization (middle-out + descending result scan)

When all 20 requests are dispatched in parallel, ordering is functionally
equivalent — every probe is in flight simultaneously and the timeout is per
request. So why bother ordering at all? Two reasons:

1. **Result iteration order matters for short-circuit logic.** When you walk
   the completed task set, do it **highest → lowest** so the first hit you
   accept is already the winner. No second pass, no `Sort-Object` on every
   iteration.
2. **Middle-out dispatch helps degraded environments.** On networks that
   serialise outbound HTTPS connections (corporate proxies, low-fd shells,
   throttled CI runners), 20 "parallel" requests effectively become a queue.
   Most active forks land in the middle of the +1..+20 window — start there
   and expand outward so the queue hits the likely winner sooner.

```
candidates = [mid, mid+1, mid-1, mid+2, mid-2, …]   # middle-out
result_scan_order = candidates sorted descending     # winner first
```

Reference implementations:

- **PowerShell:** `install.ps1` builds the middle-out array, dispatches all
  `HttpClient.SendAsync` calls, then iterates `$candidates | Sort-Object -Descending`
  and accepts the first `IsSuccessStatusCode`.
- **Bash:** `install.sh` builds the same `order=()` array, forks one
  `curl --max-time 2 -I` per candidate, then `ls "$tmp" | sort -n | tail -1`
  picks the highest hit.

This is a **portable trick any installer can adopt** — it costs ~10 lines and
makes the probe robust against degraded parallelism without changing the
public contract.

### Short-circuit / cancellation

After a result settles for the highest queued `N`, in-flight lower-N
probes MAY be cancelled to free sockets. This is an optimisation, not
required for correctness.

---

## Hand-Off

```powershell

# PowerShell hand-off

$newerUrl = "https://raw.githubusercontent.com/$Owner/$Base-v$Latest/main/install.ps1"
Write-OK "Newer version found: v$Latest (was v$Current). Handing off..."
Invoke-RestMethod -Uri $newerUrl | Invoke-Expression
exit $LASTEXITCODE   # propagate child exit
```

```bash

# Bash hand-off

newer_url="https://raw.githubusercontent.com/$OWNER/$BASE-v$LATEST/main/install.sh"
ok "Newer version found: v$LATEST (was v$CURRENT). Handing off..."
curl -fsSL "$newer_url" | bash
exit $?
```

### Hand-off invariants

| Rule | Reason |
|------|--------|
| The child install MUST inherit stdout/stderr (no detach). | User sees one continuous log. |
| The parent MUST exit with the child's exit code. | Pipelines and CI see the real outcome. |
| The probe is **not** repeated by the child. | Each install.* runs the prologue exactly once for its own `currentVersion`. The child will probe for its own +1..+20 and (correctly) find nothing higher, so it falls through to install. |
| Loop guard — env var `INSTALL_PROBE_HANDOFF_DEPTH` (int, default 0) is incremented on hand-off. If it reaches `3`, abort with a clear error to prevent runaway. | Defence in depth. |

---

## Logging Contract

User-visible output MUST follow this exact pattern so any AI / user
can recognise the probe phase in support logs:

```
▸ Detecting installer identity...
✓ Identity: alimtvnetwork/coding-guidelines-v24  (probing v6..v25)
▸ Probing 20 candidate versions in parallel (timeout 2s)...
✓ Newer version found: v19 (was v5). Handing off to v19 installer...
─────────────────────────────────────────
[child installer output begins here]
```

Or, when nothing newer is found:

```
▸ Detecting installer identity...
✓ Identity: alimtvnetwork/coding-guidelines-v24  (probing v20..v39)
▸ Probing 20 candidate versions in parallel (timeout 2s)...
✓ Already on latest (v19). Continuing local install...
```

Or, when the probe is skipped:

```
▸ Detecting installer identity...
⚠ Could not derive (owner/base/version) — skipping version probe.
```

---

## Failure Modes

| Condition | Behavior |
|-----------|-----------|
| No internet | Log `⚠ Probe skipped: network unreachable`. Continue local install. |
| GitHub rate-limit (HTTP 403/429) | Log warning with `X-RateLimit-Reset` if present. Continue local install. |
| TLS / certificate failure | Log warning. Continue local install. |
| Probe timeout exceeded | Use whatever results are in. If none, continue local install. |
| Hand-off child exits non-zero | Parent exits with same code. **Do not** retry locally — the child explicitly chose to fail. |
| `INSTALL_PROBE_HANDOFF_DEPTH ≥ 3` | Hard error, exit 1. Print: `Probe loop guard triggered — aborting.` |

**Golden rule:** the probe is a courtesy, not a gate. Any error in the
probe path falls back to "install whatever I am right now".

---

## PowerShell Reference Implementation

```powershell

# ── Version probe ─────────────────────────────────────────────────

function Invoke-LatestVersionProbe {
    # 1. Self-identify
    $sourceUrl = $null
    if ($MyInvocation.MyCommand.Definition) {
        $sourceUrl = $MyInvocation.MyCommand.Definition
    }
    if (-not $sourceUrl -and $env:INSTALL_PROBE_SOURCE_URL) {
        $sourceUrl = $env:INSTALL_PROBE_SOURCE_URL
    }

    $owner = $env:INSTALL_PROBE_OWNER
    $base  = $env:INSTALL_PROBE_BASE
    $cur   = $env:INSTALL_PROBE_VERSION

    if ($sourceUrl -match
        '^https?://[^/]+/(?<o>[^/]+)/(?<b>[A-Za-z0-9._-]+?)-v(?<v>\d+)/[^/]+/install\.ps1') {
        if (-not $owner) { $owner = $Matches.o }
        if (-not $base)  { $base  = $Matches.b }
        if (-not $cur)   { $cur   = $Matches.v }
    }

    # Fallback to embedded constants (filled by maintainer)
    if (-not $owner) { $owner = $script:ProbeOwner }
    if (-not $base)  { $base  = $script:ProbeBase }
    if (-not $cur)   { $cur   = $script:ProbeVersion }

    if (-not $owner -or -not $base -or -not $cur) {
        Write-Warn "Could not derive (owner/base/version) — skipping version probe."
        return
    }

    [int]$current = [int]$cur
    Write-OK "Identity: $owner/$base-v$current  (probing v$($current+1)..v$($current+20))"

    # 2. Loop guard
    [int]$depth = 0
    if ($env:INSTALL_PROBE_HANDOFF_DEPTH) {
        [int]::TryParse($env:INSTALL_PROBE_HANDOFF_DEPTH, [ref]$depth) | Out-Null
    }
    if ($depth -ge 3) {
        Write-Err "Probe loop guard triggered (depth=$depth) — aborting."
        exit 1
    }

    # 3. Parallel probe — PowerShell 7+ uses ForEach-Object -Parallel,
    #    Windows PowerShell 5.1 uses runspaces. Below: 5.1-compatible jobs.
    Write-Step "Probing 20 candidate versions in parallel (timeout 2s)..."
    $candidates = ($current + 1)..($current + 20)
    $jobs = foreach ($n in $candidates) {
        Start-Job -ScriptBlock {
            param($url, $n)
            try {
                $r = Invoke-WebRequest -Uri $url -Method Head `
                    -TimeoutSec 2 -UseBasicParsing -ErrorAction Stop
                if ($r.StatusCode -ge 200 -and $r.StatusCode -lt 400) {
                    return $n
                }
            } catch { }
            return $null
        } -ArgumentList "https://raw.githubusercontent.com/$owner/$base-v$n/main/install.ps1", $n
    }

    $results = @($jobs | Wait-Job -Timeout 4 | Receive-Job)
    $jobs    | Remove-Job -Force -ErrorAction SilentlyContinue
    $hits    = $results | Where-Object { $_ -is [int] } | Sort-Object -Descending
    $latest  = if ($hits.Count -gt 0) { $hits[0] } else { $current }

    if ($latest -le $current) {
        Write-OK "Already on latest (v$current). Continuing local install..."
        return
    }

    # 4. Hand-off
    $newerUrl = "https://raw.githubusercontent.com/$owner/$base-v$latest/main/install.ps1"
    Write-OK "Newer version found: v$latest (was v$current). Handing off to v$latest installer..."
    Write-Host ("-" * 50) -ForegroundColor DarkGray
    $env:INSTALL_PROBE_HANDOFF_DEPTH = ($depth + 1).ToString()
    Invoke-RestMethod -Uri $newerUrl | Invoke-Expression
    exit $LASTEXITCODE
}

# Maintainer-supplied fallback constants (used when URL parse fails)

$script:ProbeOwner   = "alimtvnetwork"
$script:ProbeBase    = "coding-guidelines"
$script:ProbeVersion = 19

Invoke-LatestVersionProbe

# ── End probe ─────────────────────────────────────────────────────

```

---

## Bash Reference Implementation

```bash

# ── Version probe ─────────────────────────────────────────────────

PROBE_OWNER_FALLBACK="alimtvnetwork"
PROBE_BASE_FALLBACK="coding-guidelines"
PROBE_VERSION_FALLBACK=19

invoke_latest_version_probe() {
    # 1. Self-identify
    local src_url="${INSTALL_PROBE_SOURCE_URL:-${BASH_SOURCE[0]:-$0}}"
    local owner="${INSTALL_PROBE_OWNER:-}"
    local base="${INSTALL_PROBE_BASE:-}"
    local cur="${INSTALL_PROBE_VERSION:-}"

    local re='^https?://[^/]+/([^/]+)/([A-Za-z0-9._-]+?)-v([0-9]+)/[^/]+/install\.sh'
    if [[ "$src_url" =~ $re ]]; then
        : "${owner:=${BASH_REMATCH[1]}}"
        : "${base:=${BASH_REMATCH[2]}}"
        : "${cur:=${BASH_REMATCH[3]}}"
    fi

    : "${owner:=$PROBE_OWNER_FALLBACK}"
    : "${base:=$PROBE_BASE_FALLBACK}"
    : "${cur:=$PROBE_VERSION_FALLBACK}"

    if [[ -z "$owner" || -z "$base" || -z "$cur" ]]; then
        warn "Could not derive (owner/base/version) — skipping version probe."
        return 0
    fi

    local current=$cur
    ok "Identity: $owner/$base-v$current  (probing v$((current+1))..v$((current+20)))"

    # 2. Loop guard
    local depth=${INSTALL_PROBE_HANDOFF_DEPTH:-0}
    if [[ $depth -ge 3 ]]; then
        err "Probe loop guard triggered (depth=$depth) — aborting."
        exit 1
    fi

    # 3. Parallel probe via background curl
    step "Probing 20 candidate versions in parallel (timeout 2s)..."
    local tmp; tmp=$(mktemp -d)
    local n
    for n in $(seq $((current+1)) $((current+20))); do
        (
            local url="https://raw.githubusercontent.com/$owner/$base-v$n/main/install.sh"
            local code
            code=$(curl -s -o /dev/null -w '%{http_code}' \
                        --max-time 2 -I "$url" 2>/dev/null || echo 000)
            if [[ "$code" =~ ^(200|301|302)$ ]]; then
                echo "$n" > "$tmp/$n"
            fi
        ) &
    done

    # Wait up to 4s for the background fan-out to settle
    local waited=0
    while [[ $waited -lt 4 ]]; do
        sleep 1
        waited=$((waited + 1))
        # Optimistic early exit when no jobs remain
        if [[ -z "$(jobs -r)" ]]; then break; fi
    done
    # Reap any stragglers
    wait 2>/dev/null || true

    local latest=$current
    if compgen -G "$tmp/*" >/dev/null; then
        latest=$(basename "$(ls "$tmp" | sort -n | tail -1)")
    fi
    rm -rf "$tmp"

    if [[ $latest -le $current ]]; then
        ok "Already on latest (v$current). Continuing local install..."
        return 0
    fi

    # 4. Hand-off
    local newer_url="https://raw.githubusercontent.com/$owner/$base-v$latest/main/install.sh"
    ok "Newer version found: v$latest (was v$current). Handing off to v$latest installer..."
    echo "──────────────────────────────────────────"
    export INSTALL_PROBE_HANDOFF_DEPTH=$((depth + 1))
    curl -fsSL "$newer_url" | bash
    exit $?
}

invoke_latest_version_probe

# ── End probe ─────────────────────────────────────────────────────

```

---

## Constraints (MUST / MUST NOT)

- MUST run before any local install work, banner, config read, or download.
- MUST be idempotent — calling twice on the same script is a no-op.
- MUST NOT prompt the user (one-liner installs are unattended).
- MUST NOT fail-close. Any probe error → continue local install.
- MUST honour `INSTALL_PROBE_HANDOFF_DEPTH` to prevent loops.
- MUST inherit stdout/stderr on hand-off (no `Start-Process -NoNewWindow` quirks).
- MUST cap probe range at +20 from current version.
- MUST cap per-request timeout at 2 seconds and total wall-time at 4 seconds.
- MUST log identity + range + outcome on stdout for traceability (sole debuggability hook for headless installs).

---

## Test Matrix

| Scenario | Expected |
|----------|----------|
| `irm .../coding-guidelines-v24/main/install.ps1 \| iex`, v19 exists | Hands off to v19, exits 0. |
| `irm .../coding-guidelines-v24/main/install.ps1 \| iex`, no v20..v39 | Falls through to local install. |
| Local dev: `.\install.ps1` (no URL context), constants point to v19 | If running script claims v19 and no v20..v39 exists, falls through. |
| Network down | Logs `⚠ network unreachable`, falls through. |
| GitHub returns 429 for all probes | Falls through. |
| Hand-off child exits 1 | Parent exits 1, no local fallback. |
| `INSTALL_PROBE_HANDOFF_DEPTH=3` already set | Hard error, exit 1. |
| URL has no `-v<N>` suffix and no fallback constants | Skip probe, fall through. |

---

## Cross-References

- [18 — Install Scripts](./19-install-scripts.md) — Local install pipeline that runs after the probe falls through.
- [05 — Handoff Mechanism](./06-handoff-mechanism.md) — Related but distinct: in-process binary self-replacement, not installer-to-installer.
- [20 — Network Requirements](./21-network-requirements.md) — Timeout, retry, and TLS conventions referenced above.
- [22 — Update Command Workflow](./23-update-command-workflow.md) — `<binary> update` uses similar probe semantics for binary-based updates.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect.
- [Riseup Asia LLC](https://riseup-asia.com) (2026)
