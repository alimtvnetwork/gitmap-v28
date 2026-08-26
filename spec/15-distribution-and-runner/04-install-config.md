# Install Config

**Version:** 1.0.0
**Updated:** 2026-04-19

---

## Purpose

`install-config.json` is the **single source of truth** for the default folder list pulled by `install.sh` / `install.ps1`. Both installers MUST read it on startup and use it as defaults for `--folders` / `-Folders`.

---

## Schema

```json
{
  "repo":    "<owner>/<repo>",
  "branch":  "<default branch>",
  "folders": ["<folder1>", "<folder2>", ...]
}
```

| Field | Type | Required | Default | Notes |
|-------|------|----------|---------|-------|
| `repo` | string | yes | `"alimtvnetwork/coding-guidelines-v24"` | GitHub `owner/repo` |
| `branch` | string | yes | `"main"` | Branch to install from when no `--version` is given |
| `folders` | string[] | yes | `["spec","linters","linter-scripts","linters-cicd"]` | Top-level folders to fetch |

Subpaths (e.g. `"spec/14-update"`) are allowed in `folders` for partial installs.

---

## Authoritative default

The committed `install-config.json` at the repo root MUST equal:

```json
{
  "repo": "alimtvnetwork/coding-guidelines-v24",
  "branch": "main",
  "folders": [
    "spec",
    "linters",
    "linter-scripts",
    "linters-cicd"
  ]
}
```

This list MUST stay in sync with the §"Default install layout" table in [`00-overview.md`](./00-overview.md) and the §"What gets installed (default)" table in [`01-install-contract.md`](./01-install-contract.md).

---

## Override precedence

Both installers MUST resolve `repo`, `branch`, and `folders` in this order (highest wins):

1. CLI flag (`--repo`, `--branch`, `--folders`, or `--version` for branch/tag)
2. Custom config file specified via `--config <path>`
3. The repo-root `install-config.json`
4. Hard-coded fallback in the installer source

---

## Custom config example

A team may ship its own config in its private repo:

```json
{
  "repo": "acme-corp/coding-guidelines-v24-fork",
  "branch": "internal",
  "folders": ["spec", "linters-cicd"]
}
```

Invocation:

```bash
./install.sh --config team-config.json
```

---

## Anti-requirements

- MUST NOT include any secret, token, or auth field. The installers operate on public repos only.
- MUST NOT include feature flags or runtime behavior — config is for source coordinates, period.
- MUST NOT add a fifth folder without updating §"Default install layout" and §"What gets installed (default)" in the same change.

---

## Cross-references

- [`./00-overview.md`](./00-overview.md)
- [`./01-install-contract.md`](./01-install-contract.md)
- [`./03-release-pipeline.md`](./03-release-pipeline.md)

---

*Install config — v1.0.0 — 2026-04-19*
