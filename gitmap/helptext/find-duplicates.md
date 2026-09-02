# gitmap find-duplicates

Audit and remediate duplicate projects and repositories across Antigravity, VS Code, Chrome, and Git.

## Usage

```bash
gitmap find-duplicates [platform] [flags]
```

## Supported Platforms

| Platform | Scope |
|---|---|
| `agy` | Antigravity JSON project configurations in `~/.gemini/config/projects/` |
| `vscode` | VS Code workspace and project manager storage |
| `chrome` | Chrome browser profiles and offline exports |
| `git` | Tracked Git repositories and clone directories |
| `(none)` | Run audits across all four platforms simultaneously |

## Examples

### Scan All Platforms

```bash
gitmap find-duplicates
```

### Scan Specific Platform

```bash
gitmap agy find-duplicates
gitmap vscode find-duplicates
gitmap chrome find-duplicates
gitmap git find-duplicates
```

**Output:**

```text
  ── Antigravity (AGY) Duplicate Projects ──
  Found 1 duplicate project group(s) (2 duplicate entries total):

  Group 1: Path: D:\wp-work\riseup-asia\gitmap (3 entries)
    PROJECT ID                             NAME                   UPDATED
    ──────────────────────────────────────────────────────────────────────────
    0349c4d0-5a91-4f3e-800f-81fd53fc724f   gitmap-v28             2026-09-02T...

  Remediation & Fix Commands for Antigravity:
  ──────────────────────────────────────────────────────────────────────────
  ● Fix Single (Delete specific duplicate project ID):
    gitmap agy rm 2c666070-443d-489a-8cb7-72e67d3e2859

  ● Fix All Together (Deduplicate & keep newest per path):
    gitmap agy optimize-projects   (alias: gitmap agy --repeat-fix)
```

## See Also

- `gitmap agy optimize-projects`: Deduplicate Antigravity projects
- `gitmap vscode optimize-projects`: Deduplicate VS Code projects
- `gitmap clone --fix`: Auto-clean repeated clone directories
