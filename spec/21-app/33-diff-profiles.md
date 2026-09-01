# Diff Profiles — Compare Repos Across Profiles

## Overview

The `diff-profiles` command compares tracked repositories between two
database profiles and highlights additions, removals, and differences.

---

## How It Works

1. Open both profile databases using `store.OpenProfile()`.
2. Load the full repo list from each.
3. Match repos by `RepoName` (the unique repository identifier).
4. Categorize each repo as:
   - **Only in A** — exists in the first profile only
   - **Only in B** — exists in the second profile only
   - **Different** — exists in both but with different paths or URLs
   - **Same** — identical in both (hidden by default)

---

## Commands

### `gitmap diff-profiles <profileA> <profileB>` (alias: `dp`)

Compare two profiles side by side.

```bash
gitmap diff-profiles default work
gitmap dp work personal
```

**Output:**
```
Comparing profiles: default ↔ work

ONLY IN default:
  docs-site          C:\repos\docs-site

ONLY IN work:
  client-portal      D:\work\client-portal

DIFFERENT:
  shared-lib
    default: C:\repos\github\shared-lib  (https)
    work:    D:\work\shared-lib           (ssh)

SAME: 12 repos (use --all to show)

Summary: 1 only-left | 1 only-right | 1 different | 12 same
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--all` | Include identical repos in the output |
| `--json` | Output as structured JSON |

### JSON Output

```json
{
  "profileA": "default",
  "profileB": "work",
  "onlyInA": [{"name": "docs-site", "path": "C:\\repos\\docs-site"}],
  "onlyInB": [{"name": "client-portal", "path": "D:\\work\\client-portal"}],
  "different": [
    {
      "name": "shared-lib",
      "a": {"path": "C:\\repos\\github\\shared-lib", "mode": "https"},
      "b": {"path": "D:\\work\\shared-lib", "mode": "ssh"}
    }
  ],
  "same": 12
}
```

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_diffprofile.go` | Command names, headers, messages |
| `cmd/diffprofiles.go` | Command entry and flag parsing |
| `cmd/diffprofilesops.go` | Comparison logic and output formatting |

---

## Constraints

- Both profiles must exist; exit 1 with error if either is missing.
- Matching is by `RepoName`, not by path or URL.
- All files under 200 lines, all functions 8–15 lines.
