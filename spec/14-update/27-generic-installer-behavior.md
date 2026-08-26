# Generic Installer Behavior Specification

> **Version:** 1.0.0
> **Status:** Normative — applies to ALL installer scripts in ALL repositories
> **Audience:** Any AI or human implementing an installer for any repo
> **Companion specs:**
> - `25-release-pinned-installer.md` (release-pinned variant details)
> - `24-update-check-mechanism/01-fundamentals.md` (V → V+N parallel discovery)
> - `26-repo-major-version-migrator.md` (cross-major repo migration)

---

## 0. Why This Spec Exists

Every repository tends to ship multiple installer entry points
(`quick-install`, `release-install`, feature-specific installers like
`error-manage-install`, etc.). Without a single shared contract these
drift apart and produce confusing results — some auto-upgrade, some
silently fall back to `main`, some cross major-version repo boundaries.

**This spec is the single source of truth for installer behavior.** Any
AI handed this document — with no other context about the repo — must be
able to implement a correct installer for that repo.

The spec is intentionally **repository-agnostic**. It uses placeholders
(`{owner}`, `{repo-base}`, `{N}`, `{tag}`) and never names a specific
project.

---

## 1. Scope

### 1.1 Installer scripts in scope

Any script whose primary purpose is to fetch + extract + install a
release of a repository. Examples (names are illustrative, not required):

| Family | Typical names |
|--------|---------------|
| Quick / generic | `install.sh`, `install.ps1`, `quick-install.*` |
| Release-pinned | `release-install.sh`, `release-install.ps1` |
| Feature/bundle  | `error-manage-install.*`, `linters-install.*`, `slides-install.*`, `wp-install.*`, `cli-install.*`, `splitdb-install.*`, `consolidated-install.*` |
| Repo-specific helpers | anything else that downloads and installs |

### 1.2 Not in scope

- Build scripts that operate on an already-checked-out tree.
- Package-manager wrappers (`npm i`, `cargo install`) that do not
  themselves resolve a version.
- CI workflows (they MAY call installers, but are not installers).

---

## 2. Terminology

| Term | Meaning |
|------|---------|
| **Strict version** | A version explicitly supplied by the user via CLI flag (`--version`, `-Version`), env var, or baked-in placeholder in a release-asset copy of the script. |
| **Implicit mode** | No strict version supplied — installer may pick a source. |
| **Pinned mode** | Strict version supplied — installer MUST install exactly that. |
| **Main branch fallback** | Downloading the tip of the default branch (usually `main`) as a tarball/zipball. |
| **Versioned repo** | A repo whose name encodes a major version, e.g. `repo-v15`, `repo-v16`. |
| **V → V+N discovery** | Probing the current repo and the next N versioned repos in parallel to find the highest existing one. |

---

## 3. The Two Modes (Normative)

Every installer MUST operate in exactly one of two modes per invocation,
selected at the very start:

```
┌──────────────────────────────────────────────────────────────┐
│  Strict version supplied?                                    │
│   ├─ YES → PINNED MODE  (§4)                                 │
│   └─ NO  → IMPLICIT MODE (§5)                                │
└──────────────────────────────────────────────────────────────┘
```

No installer may switch modes mid-run.

---

## 4. PINNED MODE (Strict — Hard Rules)

Triggered when ANY of the following resolve to a non-empty value:

1. CLI flag — Bash `--version <tag>`, PowerShell `-Version <tag>`.
2. Environment variable — `INSTALLER_VERSION` (or repo-defined equivalent).
3. Baked-in `__VERSION_PLACEHOLDER__` substituted at release-asset build
   time (see `25-release-pinned-installer.md`).
4. The script was downloaded from a URL containing `/releases/download/<tag>/`.

### 4.1 MUST

- Install **exactly** the resolved tag. Nothing else.
- Download from `…/releases/download/<tag>/…` first, then fall back to the
  tag tarball (`…/archive/refs/tags/<tag>.tar.gz` or
  `codeload.github.com/{owner}/{repo}/zip/refs/tags/<tag>`).
- Fail loudly with a non-zero exit if the tag does not exist.
- Print the resolved tag in the startup banner.

### 4.2 MUST NOT

- ❌ Query `…/releases/latest` or any "what's newest?" endpoint.
- ❌ Fall back to the `main` branch.
- ❌ Cross repository boundaries (no V → V+N discovery, no
  `repo-v15 → repo-v16` jump). The installer stays on the SAME repo
  the tag belongs to.
- ❌ Pick a "compatible" or "nearest" version.
- ❌ Print upgrade nags or "newer version available" notices.
- ❌ Silently downgrade to implicit mode on 404 — fail with exit code
  `3` (release/asset not found) per §8.

### 4.3 Override precedence

CLI flag > env var > baked-in placeholder > URL-encoded tag.
If two sources disagree, emit a WARNING line (not an error) and use the
higher-precedence value.

---

## 5. IMPLICIT MODE (Default — Discovery Allowed)

Triggered when no strict version is resolvable.

### 5.1 Source order (MUST be tried in this order)

1. **Latest release** of the current repo
   (`https://api.github.com/repos/{owner}/{repo}/releases/latest`).
   Use the resolved tag exactly as if PINNED MODE had been requested
   for that tag — except that 404 here is NOT fatal, it falls through.
2. **V → V+N parallel discovery** (§6) — only if the current repo name
   matches the versioned-repo pattern.
3. **Main branch fallback** — download the default branch tarball
   (`…/archive/refs/heads/main.tar.gz` or codeload equivalent).

At each step, log which source was tried and the outcome
(`found` / `not-found` / `error`). Stop at the first success.

### 5.2 Main branch fallback rules

- Used only after both (1) and (2) fail.
- MUST print a clearly visible warning:
  `⚠️  No release found — installing from main branch (unstable).`
- MUST record `source: main-branch` in any post-install metadata
  (e.g. `installed-version.json`).

### 5.3 Opt-outs

Implementations MUST offer the following flags to disable individual sources (exact names; no synonyms):

| Flag (Bash / PowerShell) | Effect |
|--------------------------|--------|
| `--no-discovery` / `-NoDiscovery` | Skip §5.1 step 2 (V → V+N). |
| `--no-main-fallback` / `-NoMainFallback` | Skip §5.1 step 3. |
| `--offline` / `-Offline` | Skip all network operations; use local archive only. |

---

## 6. V → V+N Parallel Repo Discovery

### 6.1 Naming convention (REQUIRED for discovery)

The repo name MUST end in `-v{N}` where `{N}` is a positive integer:

```
{repo-base}-v{N}        e.g.  coding-guidelines-v24
                              movie-cli-v2
                              repo-v20
```

If the repo does not match this pattern, discovery is skipped and the
installer proceeds to main-branch fallback.

### 6.2 Algorithm

1. Parse current `{N}` from the repo name.
2. Build candidate names for `{N+1}` through `{N+LOOKAHEAD}`.
   - **Default `LOOKAHEAD = 20`** (per user requirement; was 5 in the
     older `24-update-check-mechanism/01-fundamentals.md`). New
     installers MUST adopt `LOOKAHEAD = 20`. Installers explicitly
     tagged `legacy` in their manifest MAY retain `LOOKAHEAD = 5`;
     all other implementations MUST use 20.
3. Fire all `LOOKAHEAD + 1` probes (V plus V+1..V+LOOKAHEAD)
   **in parallel**:
   - Bash: `&` + `wait`, or GNU `parallel`.
   - PowerShell: `Start-Job` / `ForEach-Object -Parallel`.
   - Go: goroutines + `errgroup`.
   - Node: `Promise.all`.
   - Python: `asyncio.gather`.
4. Each probe is a `HEAD` (or lightweight `GET`) against
   `https://github.com/{owner}/{repo-base}-v{K}` and/or its
   `releases/latest` endpoint.
5. **HTTP 200** = candidate exists. **HTTP 404** = does not exist
   (no retry). Any other status = treat as not-found and log.
6. The **highest `K` returning 200 wins.** Switch the active repo to
   that candidate and re-enter §5.1 step 1 against it.
7. Discovery NEVER walks past `V+LOOKAHEAD` in a single run. A new
   window opens on the next invocation if the user has installed the
   newer version.

### 6.3 HTTP client requirements

| Setting | Value |
|---------|-------|
| Per-probe timeout | 5 seconds |
| Total discovery deadline | 10 seconds |
| Retries | 0 |
| User-Agent | `<InstallerName>/<Version> Discovery` |
| Proxy | Inherit `HTTP_PROXY` / `HTTPS_PROXY` |

A timed-out discovery is NOT an error — it falls through to the next
source in §5.1.

### 6.4 Disabled in PINNED MODE

Discovery is **forbidden** in PINNED MODE. A pinned tag belongs to
exactly one repo; jumping to `repo-v{N+1}` would violate §4.2.

---

## 7. Logging Expectations

Every installer MUST print, at startup, a banner showing:

```
    📦 {InstallerName} {InstallerVersion}
       mode:    pinned | implicit
       repo:    {owner}/{repo}
       version: {tag | "discovering…" | "main"}
       source: release-asset | tag-tarball | main-branch | local-archive
```

At completion, print:

```
    ✅ installed {tag} from {source} → {dest}
```

On failure, print the exit code and a one-line cause.

---

## 8. Exit Codes (Normative)

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Generic failure (missing tool, unknown flag, network error in implicit mode after all sources failed) |
| `2` | Offline mode required a network operation |
| `3` | Pinned release / asset not found (PINNED MODE only) |
| `4` | Verification failed (checksum / required-paths check) |
| `5` | Inner installer / handoff rejected |

Implementations MAY add codes ≥ `10` for repo-specific conditions but
MUST NOT redefine `0–5`.

---

## 9. Security Considerations

- Always use HTTPS. Reject `http://` URLs.
- When checksums are published alongside a release, verify them by
  default; allow `--no-verify` only with a loud warning.
- Never `eval` downloaded content without first writing it to disk and
  verifying integrity (where checksums exist).
- Never echo secrets (tokens, env vars matching `*TOKEN*`, `*KEY*`,
  `*SECRET*`) in the banner or logs.

---

## 10. Acceptance Criteria

An installer conforms to this spec iff:

1. ✅ Without any version flag, it installs the latest release of the
   current repo, falling through V → V+N discovery, then main, in that
   order, with visible logging at each step.
2. ✅ With a `--version <tag>` flag, it installs exactly that tag and
   refuses every fallback in §4.2.
3. ✅ A non-existent pinned tag exits with code `3` and a clear error.
4. ✅ Main-branch fallback prints the "unstable" warning and tags
   metadata `source: main-branch`.
5. ✅ V → V+N discovery fires probes in parallel and respects the
   `LOOKAHEAD = 20` default.
6. ✅ Discovery is skipped entirely in PINNED MODE.
7. ✅ The banner shows `mode`, `repo`, `version`, `source`.
8. ✅ Exit codes match §8.
9. ✅ The spec file path (this document) appears in the repo's installer
   readme so future maintainers can find it.

---

## 11. Reference: How To Apply This To Any Repo

An AI given only this document and a target repo should:

1. Identify all installer scripts (see §1.1).
2. For each: classify as PINNED-only, IMPLICIT-only, or hybrid.
3. Implement §3 mode dispatch at the top of the script.
4. Wire PINNED MODE to §4 rules.
5. Wire IMPLICIT MODE to §5 source order, §6 discovery if the repo name
   matches the `-v{N}` pattern, and §5.2 fallback otherwise.
6. Add the §7 banner and §8 exit codes.
7. Add a paragraph in the repo's `readme.md` linking back to this spec.

---

## 12. Where This Spec Lives

**Canonical path:** `spec/14-update/27-generic-installer-behavior.md`

Share that exact path with any AI to bootstrap a conformant installer in
any repository.

---

*Generic Installer Behavior — v1.0.0 — 2026-04-22*