# Install Contract

**Version:** 1.0.0
**Updated:** 2026-04-19

---

## Purpose

Define exactly what `install.sh` / `install.ps1` (at the repo root) MUST do, what they pull, and how they behave on errors and re-installs.

---

## One-liner invocations

These invocations MUST work without any prior setup other than `curl` (or `wget`) and `bash` (Linux/macOS) or PowerShell 7+ (Windows):

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.sh | bash
```

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.ps1 | iex
```

---

## What gets installed (default)

When invoked with no flags, both installers MUST pull the following four folders from `alimtvnetwork/coding-guidelines-v24@main` and unpack them into the **current working directory**:

| Folder | Mandatory | Purpose |
|--------|-----------|---------|
| `spec/` | yes | Full coding-guidelines spec tree (568+ files) |
| `linters/` | yes | Language-specific lint plugins, ESLint configs, tree-sitter queries |
| `linter-scripts/` | yes | Legacy orchestrator (validator wrappers, helper scripts) |
| `linters-cicd/` | yes | Python check suite, registry, `run-all.sh`, baseline |

This list is loaded from [`install-config.json`](./04-install-config.md). The list above MUST equal the `folders` array in that file.

> **Why all four?** The user wants every distributable component installed by default. If they want a subset, they override with `--folders` (Bash) or `-Folders` (PowerShell).

---

## Versioning

| Mode | Flag | Behavior |
|------|------|----------|
| Latest tag | (none) | Probes for newer `coding-guidelines-vN` repos (middle-out, parallel, 2s timeout). Hands off to newer installer if found. |
| Pinned version | `--version vX.Y.Z` (Bash) / `-Version vX.Y.Z` (PS) | Downloads the tagged tarball/zipball and extracts the four folders. |
| Pinned branch | `--branch <name>` / `-Branch <name>` | Downloads from the specified branch head. |
| Skip probe | `-n`, `--no-probe`, `--no-latest` (Bash) / `-NoProbe` (PS) | Use the running installer as-is, no version detection. |

---

## File-merge semantics

| Mode | Flag | Behavior |
|------|------|----------|
| Default | (none) | Overwrite existing files silently (legacy behavior) |
| Interactive | `--prompt` / `-Prompt` | Ask `[y]es / [n]o / [a]ll / [s]kip-all` for each existing file |
| Force | `--force` / `-Force` | Overwrite every existing file (mutually exclusive with `--prompt`) |
| Dry run | `--dry-run` / `-DryRun` | Print would-create/would-overwrite, write nothing |

---

## Listings (no install)

| Mode | Flag | Output |
|------|------|--------|
| List release tags | `--list-versions` / `-ListVersions` | Up to 50 tags from `api.github.com/repos/<repo>/releases` |
| List top-level folders | `--list-folders` / `-ListFolders` | Top-level dirs in the chosen ref (via `api.github.com/repos/<repo>/contents`) |

---

## Cleanup contract

Both installers MUST:

1. Create a unique temp directory under the OS temp root.
2. Download the archive and extract into the temp directory.
3. Copy/merge files into the destination.
4. **Always** clean up the temp directory on exit (success or failure), via `trap` (Bash) or `try/finally` (PS).
5. Verify cleanup happened and emit a warning if it didn't.

---

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Generic failure (network, parse error, missing tools) |
| 2 | Bad CLI flag combination (e.g. `--prompt` with `--force`) |

---

## Error messages

Errors MUST include:

- The remote URL that failed (so the user can `curl` it manually to debug).
- The local path that was being written (so the user can check permissions).
- One actionable next step ("retry with `--no-probe`", "delete X and re-run", etc.).

No silent failures. No bare stack traces.

---

## Anti-requirements

- MUST NOT require Node.js, Python, or Go on the host. (The installers themselves use only `curl`/`wget`/`tar`/`unzip` on Bash, or built-in PS cmdlets on Windows.)
- MUST NOT require GitHub authentication for public repos.
- MUST NOT install or modify anything outside `<dest>/{spec,linters,linter-scripts,linters-cicd}/`.

---

## Cross-references

- [`./00-overview.md`](./00-overview.md) — Distribution overview
- [`./03-release-pipeline.md`](./03-release-pipeline.md) — Where the artifacts come from
- [`./04-install-config.md`](./04-install-config.md) — Folder-list contract

---

*Install contract — v1.0.0 — 2026-04-19*
