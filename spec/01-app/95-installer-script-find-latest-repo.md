# Versioned Repository Discovery for Installers

**Status:** Draft
**Audience:** Any AI or human implementing an installer / bootstrap script that downloads from a versioned GitHub repository.
**Scope:** Generic. Applies to any project whose repo name ends in a `-v<N>` suffix (e.g. `gitmap-v28`, `myapp-v7`, `cli-tool-v12`).

---

## 1. Purpose

When an installer is invoked against a versioned repo URL (e.g. `https://github.com/<owner>/<name>-v3`), it MUST:

1. Treat the requested repo as the **baseline**.
2. Probe a bounded sequence of higher-numbered sibling repos (`-v4`, `-v5`, …) under the same owner.
3. Pick the **highest existing** repo as the **effective install target**.
4. Delegate execution to **that** repo's installer script.
5. Log every step (probe, hit, miss, choice, delegation) clearly.

This guarantees a stale install URL (e.g. someone bookmarked `gitmap-v28` years ago) still pulls the user onto the latest major repo line without manual intervention.

---

## 2. Terminology

| Term              | Meaning                                                              |
|-------------------|----------------------------------------------------------------------|
| **Baseline repo** | The repo URL the user originally invoked (`<owner>/<name>-v<N>`).    |
| **Suffix**        | The trailing `-v<N>` segment. `N` MUST be a positive integer.        |
| **Probe ceiling** | Maximum `N` to try. Default `30`. Configurable.                      |
| **Effective repo**| Highest `-v<M>` repo (`M >= N`) that exists. Falls back to baseline. |
| **Installer path**| Path inside the repo to the installer (e.g. `gitmap-v28/scripts/install.ps1` or `install.sh`). MUST be identical across versions. |

---

## 3. Generic URL Pattern

```
https://github.com/<owner>/<name>-v<N>
                           ^^^^^^^^
                           suffix
```

**Splitting rule:**
- Find the **last** `-v<digits>` at the end of the repo name.
- Everything before it is the **stem**: `<name>`.
- `<digits>` is the baseline integer `N`.

**Examples:**

| Input URL                                          | stem        | N  |
|----------------------------------------------------|-------------|----|
| `https://github.com/alimtvnetwork/gitmap-v28`       | `gitmap-v28`    | 3  |
| `https://github.com/acme/widgets-v12`              | `widgets`   | 12 |
| `https://github.com/foo/bar-baz-v1`                | `bar-baz`   | 1  |
| `https://github.com/foo/no-suffix-here`            | (no match — skip discovery, install baseline as-is) |

---

## 4. Discovery Algorithm

```
INPUT:  baselineUrl, probeCeiling = 30
OUTPUT: effectiveUrl

(stem, N) = parseSuffix(baselineUrl)
if no suffix:
    log "No -v<N> suffix detected; installing baseline as-is."
    return baselineUrl

effectiveN = N
log "Baseline: <owner>/<stem>-v<N>. Probing -v<N+1>..-v<probeCeiling>."

for M in (N+1) .. probeCeiling:
    candidateUrl = "https://github.com/<owner>/<stem>-v<M>"
    exists = headCheck(candidateUrl)         // see §5
    if exists:
        log "  [HIT]  -v<M> exists"
        effectiveN = M
    else:
        log "  [MISS] -v<M>"
        // FAIL FAST: stop on first miss above the current effective.
        // Versioned repos are expected to be contiguous.
        break

if effectiveN == N:
    log "No higher version found. Using baseline -v<N>."
else:
    log "Latest available: -v<effectiveN>. Switching from -v<N>."

effectiveUrl = "https://github.com/<owner>/<stem>-v<effectiveN>"
return effectiveUrl
```

### Fail-fast rationale

Project versions are expected to be **contiguous** (`v3, v4, v5, …`). On the first MISS after the last HIT, stop probing. This keeps the probe to **at most 1 extra request** beyond the latest existing version in the common case.

### Non-contiguous fallback (optional, off by default)

If a project is known to skip versions, the implementation MAY accept a `--no-fail-fast` flag that probes the full range up to `probeCeiling` and picks the maximum hit. This MUST be opt-in. Default is fail-fast.

---

## 5. Existence Check

Use a cheap HTTP probe — **do not** clone, **do not** download release assets.

| Method          | URL                                                          | Success         |
|-----------------|--------------------------------------------------------------|-----------------|
| `HEAD` (preferred) | `https://github.com/<owner>/<stem>-v<M>`                  | HTTP `200`      |
| `GET` (fallback)   | same URL, drop the body                                   | HTTP `200`      |
| `GET` (API form)   | `https://api.github.com/repos/<owner>/<stem>-v<M>`        | HTTP `200` JSON |

**Treat as MISS:** `404`, `301`/`302` to a different owner/repo, network errors, timeouts.

**Timeout:** 5 seconds per probe. Fail = MISS, not abort.

**Auth:** Anonymous. If `GITHUB_TOKEN` is present in env, MAY use it to raise rate limits — never required.

---

## 6. Delegation

After resolving `effectiveUrl`:

1. Construct the new installer URL by replacing the stem in the original installer URL:
   ```
   https://raw.githubusercontent.com/<owner>/<stem>-v<effectiveN>/<branch>/<installerPath>
   ```
2. Re-invoke the installer **from the effective repo**, passing through all original flags (`--dir`, `--version`, etc.).
3. Exit with the delegated installer's exit code.

**Do NOT** continue executing the original installer after delegating. The newer installer is the source of truth.

**Loop guard:** the delegated installer MUST detect it was invoked with an env var (e.g. `INSTALLER_DELEGATED=1`) and skip discovery on its second run. Set this var before delegating.

---

## 7. Logging Format

Every step prints to stdout in this format:

```
  [discovery] baseline: <owner>/<stem>-v<N>
  [discovery] probe ceiling: <probeCeiling>
  [discovery] HEAD https://github.com/<owner>/<stem>-v4 ... HIT
  [discovery] HEAD https://github.com/<owner>/<stem>-v5 ... HIT
  [discovery] HEAD https://github.com/<owner>/<stem>-v6 ... MISS (fail-fast)
  [discovery] effective: <owner>/<stem>-v5
  [discovery] delegating to https://raw.githubusercontent.com/<owner>/<stem>-v5/main/<installerPath>
```

Use a consistent prefix (`[discovery]`) so it's grep-friendly.

---

## 8. Configuration Knobs

| Flag / env var              | Default | Purpose                                               |
|-----------------------------|---------|-------------------------------------------------------|
| `--probe-ceiling <N>`       | `30`    | Highest version to try.                                |
| `--no-discovery`            | off     | Skip §4 entirely; install baseline.                    |
| `--no-fail-fast`            | off     | Probe full range; max-hit wins.                        |
| `INSTALLER_DELEGATED=1`     | unset   | Set by parent; child skips discovery (loop guard).     |
| `GITHUB_TOKEN`              | unset   | Optional auth for higher rate limits.                  |

---

## 9. Reference Implementations (sketch)

### PowerShell

```powershell
function Resolve-EffectiveRepo {
    param([string]$Owner, [string]$Stem, [int]$BaselineN, [int]$Ceiling = 30)

    $effective = $BaselineN
    Write-Host "  [discovery] baseline: $Owner/$Stem-v$BaselineN"

    for ($m = $BaselineN + 1; $m -le $Ceiling; $m++) {
        $url = "https://github.com/$Owner/$Stem-v$m"
        try {
            $resp = Invoke-WebRequest -Uri $url -Method Head -TimeoutSec 5 `
                -UseBasicParsing -ErrorAction Stop
            if ($resp.StatusCode -eq 200) {
                Write-Host "  [discovery] HEAD $url ... HIT"
                $effective = $m
                continue
            }
        } catch {
            Write-Host "  [discovery] HEAD $url ... MISS (fail-fast)"
            break
        }
    }

    Write-Host "  [discovery] effective: $Owner/$Stem-v$effective"
    return $effective
}
```

### Bash

```bash
resolve_effective_repo() {
    local owner="$1" stem="$2" baseline="$3" ceiling="${4:-30}"
    local effective="$baseline" m url

    printf '  [discovery] baseline: %s/%s-v%s\n' "$owner" "$stem" "$baseline"

    for (( m = baseline + 1; m <= ceiling; m++ )); do
        url="https://github.com/${owner}/${stem}-v${m}"
        if curl -sfI --max-time 5 "$url" >/dev/null 2>&1; then
            printf '  [discovery] HEAD %s ... HIT\n' "$url"
            effective=$m
        else
            printf '  [discovery] HEAD %s ... MISS (fail-fast)\n' "$url"
            break
        fi
    done

    printf '  [discovery] effective: %s/%s-v%s\n' "$owner" "$stem" "$effective"
    echo "$effective"
}
```

---

## 10. Edge Cases

| Case                                        | Behaviour                                                  |
|---------------------------------------------|------------------------------------------------------------|
| Repo URL has no `-v<N>` suffix              | Skip discovery; install baseline.                          |
| Baseline repo itself doesn't exist (404)    | Probe anyway — user might be on an old name. If nothing found, exit 1 with clear error. |
| GitHub rate limit (HTTP 403)                | Treat as MISS, log a hint about `GITHUB_TOKEN`. Do not abort. |
| All probes time out                         | Use baseline; log a network-quality warning.               |
| `effectiveN > 100`                          | Cap at `probeCeiling`; never probe unbounded.              |
| Delegated installer also has discovery      | Loop guard env var prevents infinite recursion.            |

---

## 11. Out of Scope

- Choosing a specific **release tag** within a repo. That stays the existing installer's job (`--version v2.86.0`).
- Migrating user data between major versions. Discovery only changes the source repo, not on-disk state.
- Cross-owner discovery (forks). Discovery is restricted to the same `<owner>`.

---

## 12. Acceptance Checklist

An implementation conforms when:

- [x] Parses `-v<N>` suffix from the URL's last path segment.
- [x] Probes `N+1 .. ceiling` with HEAD requests, fail-fast on first MISS.
- [x] Logs every probe with `[discovery]` prefix.
- [x] Delegates to the highest existing `-v<M>` installer, passing through flags.
- [x] Sets `INSTALLER_DELEGATED=1` to prevent recursion.
- [x] Honours `--no-discovery`, `--probe-ceiling`. (`--no-fail-fast` not yet implemented — fail-fast is the only mode.)
- [x] Falls back gracefully on network errors (use baseline, never crash).
- [x] No more than `probeCeiling - N + 1` HTTP requests in worst case.

**Implementation status (v2.88.0):**

| Script                          | Discovery | Loop guard | Knobs                            |
|---------------------------------|-----------|------------|----------------------------------|
| `install-quick.ps1`             | ✅        | ✅         | `-NoDiscovery`, `-ProbeCeiling`  |
| `install-quick.sh`              | ✅        | ✅         | `--no-discovery`, `--probe-ceiling` |
| `gitmap-v28/scripts/install.ps1`    | ✅        | ✅         | `-NoDiscovery`, `-ProbeCeiling`  |
| `gitmap-v28/scripts/install.sh`     | ✅        | ✅         | `--no-discovery`, `--probe-ceiling` |

**Outstanding:**

- [ ] `--no-fail-fast` opt-in mode (probe full range, max-hit wins) — deferred until a non-contiguous version line appears.


---

## 13. Concrete Example (gitmap-v28)

User runs:

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.sh | bash
```

Discovery flow:

```
  [discovery] baseline: alimtvnetwork/gitmap-v28
  [discovery] probe ceiling: 30
  [discovery] HEAD https://github.com/alimtvnetwork/gitmap-v28 ... HIT
  [discovery] HEAD https://github.com/alimtvnetwork/gitmap-v28 ... MISS (fail-fast)
  [discovery] effective: alimtvnetwork/gitmap-v28
  [discovery] delegating to https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install-quick.sh
```

The user transparently lands on `gitmap-v28` even though they invoked `gitmap-v28`.

---

## See Also

- [83-install-bootstrap.md](83-install-bootstrap.md) — Bootstrap installer mechanics.
- [94-install-script.md](94-install-script.md) — Canonical installer contract.
- [spec/08-generic-update/](../08-generic-update/) — In-place self-update (different concern: same repo, new release tag).
