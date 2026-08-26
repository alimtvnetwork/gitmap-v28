# Fundamentals — Discovery Algorithm

> **Version:** 1.0.0
> **Parent:** [00-overview.md](./00-overview.md)

---

## 1. Naming Convention

The repository name encodes the major version as a `-v{N}` suffix:

```
repo-v15   →  major version 15
repo-v16   →  major version 16
```

The owner is either a GitHub user or organization — the discovery code
uses a single placeholder `{owner}` that resolves to whichever is
configured in `06-Seedable-Config`.

---

## 2. Status Script Location

| OS | URL template |
|----|--------------|
| Windows | `https://raw.githubusercontent.com/{owner}/{repo}/main/Status.ps1` |
| Unix    | `https://raw.githubusercontent.com/{owner}/{repo}/main/Status.sh` |

The script is **always** fetched from `raw.githubusercontent.com` and
**never** from a local copy. Local copies are not trusted and not used.

Standalone invocation (executed when an end-user pastes from the README):

```powershell

# Windows

irm https://raw.githubusercontent.com/{owner}/repo-v20/main/Status.ps1 | iex
```

```bash

# Unix

curl -fsSL https://raw.githubusercontent.com/{owner}/repo-v20/main/Status.sh | bash
```

---

## 3. Parallel V → V+5 Discovery

Given the installed CLI is at `repo-v{N}`, the discoverer dispatches
**six** probes simultaneously: the current repo (V) and five lookahead
repos (V+1 through V+5).

### Probe matrix (example: current = `repo-v15`)

| Probe | Repo | Windows URL | Unix URL |
|-------|------|-------------|----------|
| V (current) | `repo-v15` | `…/repo-v15/main/Status.ps1` | `…/repo-v15/main/Status.sh` |
| V+1 | `repo-v16` | `…/repo-v16/main/Status.ps1` | `…/repo-v16/main/Status.sh` |
| V+2 | `repo-v17` | `…/repo-v17/main/Status.ps1` | `…/repo-v17/main/Status.sh` |
| V+3 | `repo-v18` | `…/repo-v18/main/Status.ps1` | `…/repo-v18/main/Status.sh` |
| V+4 | `repo-v19` | `…/repo-v19/main/Status.ps1` | `…/repo-v19/main/Status.sh` |
| V+5 | `repo-v20` | `…/repo-v20/main/Status.ps1` | `…/repo-v20/main/Status.sh` |

### Rules

1. All six probes fire **at once** (Go routines, `Promise.all`,
   `asyncio.gather`, etc. — language-appropriate concurrency).
2. **HTTP 404** on either the repo or the file is treated as
   *not-found*. No retry.
3. **HTTP 200** with valid JSON is treated as *found*.
4. **HTTP 200 with malformed JSON** is treated as *not-found* and the
   parse error is logged (does not abort discovery).
5. The **highest version returning a valid response wins.**
6. Discovery **never** walks past V+5, even if all five lookahead
   probes succeed. A new V+5 window begins on the next scheduled check
   if the user has installed the new version.

---

## 4. Selecting OS

| Detected OS | Script |
|-------------|--------|
| `windows`   | `Status.ps1` |
| `linux`, `darwin`, anything POSIX | `Status.sh` |

Both URLs are still recorded in the combined JSON so cross-platform
consumers (e.g., a CI dashboard) can show install commands for either OS.

---

## 5. HTTP Client Requirements

| Setting | Value |
|---------|-------|
| Timeout (per probe) | 5 seconds |
| Total discovery deadline | 10 seconds |
| Retries | 0 (a missed probe is treated as not-found) |
| User-Agent | `<CliName>/<Version> UpdateChecker` |
| Proxy | Inherit from environment (`HTTP_PROXY`, `HTTPS_PROXY`) |

A discovery that times out is **not** an error — it is recorded as
"no update found" and the next check happens on the next interval.

---

*Fundamentals — v1.0.0 — 2026-04-20*
