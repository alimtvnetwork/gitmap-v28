# `fix-repo.sh` / `fix-repo.ps1` — Normative Contract (AI Execution Prompt)

> **/goal** Define the normative detection algorithm, flag set, replacement rules, and exit code contract for `fix-repo.sh` and `fix-repo.ps1` cross-platform repository name rewriters.
> **/learn** Master numeric-overflow guards, regex URL parsing without altering host segments, lossless UTF-8 file traversal via git ls-files -z, and closed flag validation.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Parse remote URL into `{Host, Owner, RepoFull}` and validate `-v(\d+)$` version suffix with CurrentVersion >= 1.
- [ ] `/learn` Never modify non-target numeric suffixes (e.g. numeric-overflow guard: `coding-guidelines-v170` must never match `v17`).
- [ ] `/goal` Enforce closed mode flags (`-2`, `-3`, `-5`, `-all`) and propagate deterministic exit codes (`2` through `7`).
- [ ] `/learn` Verify dry-run safety and idempotency where a second consecutive run produces zero file modifications.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version:** 1.0.0 · **Updated:** 2026-04-28
> Implements: `fix-repo.sh`, `fix-repo.ps1`, `scripts/fix-repo/*`

---

## 1. Goal

Detect the current repo's base name and version from its git remote,
then rewrite every prior-version token (`{Base}-v{N}`) to the current
version (`{Base}-v{Current}`) across all tracked text files.

---

## 2. Detection algorithm

1. **Repo root** — `git rev-parse --show-toplevel`. Failure → exit `2`
   (`E_NOT_A_REPO`).
2. **Remote URL** — `git remote get-url origin`; on failure, first
   `(fetch)` line of `git remote -v`. Both empty → exit `3`
   (`E_NO_REMOTE`).
3. **Parse URL** into `{Host, Owner, RepoFull}`. Supported forms:

   | Form | Pattern |
   |------|---------|
   | HTTPS | `https?://<host>[:port]/<owner>/<repo>(.git)?[/...]` |
   | SSH (scp-like) | `git@<host>:<owner>/<repo>(.git)?` |
   | SSH (ssh://) | `ssh://git@<host>[:port]/<owner>/<repo>(.git)?` |

   Strip a single trailing `.git`. Anything past the repo segment is ignored.
4. **Split `RepoFull`** on suffix regex `-v(\d+)$` (case-sensitive `-v`).
   Miss → exit `4` (`E_NO_VERSION_SUFFIX`).
5. **Validate** `CurrentVersion >= 1`. `<=0` → exit `5` (`E_BAD_VERSION`).
   `Current == 1` AND mode != `--all` → exit `0` "nothing to replace".

---

## 3. Flag set (exhaustive, closed)

Exactly one mode flag may be passed. Default = `--2` / `-2`.

| PowerShell | Bash | Meaning |
|---|---|---|
| (none) / `-2` | (none) / `--2` | Replace last **2** versions |
| `-3` | `--3` | Replace last **3** versions |
| `-5` | `--5` | Replace last **5** versions |
| `-all` | `--all` | Replace every version `1..Current-1` |
| `-DryRun` | `--dry-run` | Report only; no writes |
| `-Verbose` | `--verbose` | Print every modified file |

Mode set is **closed**: `--4`, `--6`, etc. → exit `6` (`E_BAD_FLAG`).
Two mode flags → exit `6`. `--dry-run` and `--verbose` may combine
freely with any mode flag.

---

## 4. Target version computation

Let `M` = mode integer (`2`, `3`, `5`, or `Current-1` for `--all`):

```
TargetVersions = { v ∈ ℤ | max(1, Current - M) <= v <= Current - 1 }
```

Examples (Current = 18): `--2` → `{16,17}`; `--all` → `{1..17}`.
Examples (Current = 3): `--3` → `{1,2}` (clamped, never below 1).

---

## 5. Replacement rules

1. **Token:** literal `{RepoBase}-v{N}` for each `N ∈ TargetVersions`.
2. **Replacement:** literal `{RepoBase}-v{CurrentVersion}`.
3. **Match:** plain substring, case-sensitive, with a numeric-overflow
   guard — token MUST NOT be immediately followed by a digit.
   - Matches inside `https://github.com/x/coding-guidelines-v24`
   - Does NOT match inside `coding-guidelines-v170`
4. **URL handling:** plain substring everywhere; host is preserved
   automatically because it is not part of the token.
5. **Order:** ascending by `N`. Idempotent — second run changes 0 files.

---

## 6. File traversal

1. `git ls-files -z` from repo root (honors `.gitignore`, skips `.git/`,
   excludes submodule contents).
2. Skip if any of: symlink, size > 5 MiB, NUL byte in first 8192 bytes,
   or extension in the always-binary set (`.png .jpg .jpeg .gif .webp
   .ico .pdf .zip .tar .gz .tgz .bz2 .xz .7z .rar .woff .woff2 .ttf .otf
   .eot .mp3 .mp4 .mov .wav .ogg .webm .class .jar .so .dylib .dll .exe
   .pyc`).
3. Read as UTF-8 with **lossless** byte fallback.
4. Preserve original line endings and trailing-newline presence.

---

## 7. Output contract

Header (always):

```
fix-repo  base=<RepoBase>  current=v<N>  mode=<flag>
targets:  v<a>, v<b>, ...
host:     <Host>  owner=<Owner>
```

Per-file (only when `--verbose`):

```
modified: <path> (<int> replacements)
```

Summary (always):

```
scanned: <int> files
changed: <int> files (<int> replacements)
mode:    <write|dry-run>
```

---

## 8. Exit codes

| Code | Symbol | Meaning |
|---|---|---|
| `0` | OK | Success (incl. dry-run, "nothing to replace") |
| `2` | E_NOT_A_REPO | `git rev-parse` failed |
| `3` | E_NO_REMOTE | No remote URL found / unparseable |
| `4` | E_NO_VERSION_SUFFIX | Repo name has no `-vN` suffix |
| `5` | E_BAD_VERSION | `N <= 0` or non-integer |
| `6` | E_BAD_FLAG | Unknown / conflicting flags |
| `7` | E_WRITE_FAILED | At least one file failed to write |

---

## 9. CODE RED compliance (mandatory)

- Functions ≤ 8–15 effective lines.
- Files ≤ 300 lines (split helpers into `scripts/fix-repo/*`).
- Zero nested `if`. Guard-and-return only.
- Booleans: `Is*` / `Has*` (PS), `is_*` / `has_*` (Bash).
- No swallowed errors; every failure is logged AND reflected in exit code.
- All exit codes named (no magic numbers in branching code).

---

## 10. Non-goals

- Rewriting URLs across hosts.
- Acting on files outside the working tree.
- Touching files not tracked by git.
- Migrating away from `-vN` naming.
- `.bak` files or auto-staging — recovery is `git checkout -- .`.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-INSTALL-002: Fix-Repo Version Token Rewriter Specification Conformance

**Given** Installer and auto-fix repository management specifications.
**When** Audited against this installation specification.
**Then** Zero contract or visibility violations are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only
```
**Expected:** exit 0. Zero violations.
