# `gitmap agy clear`

Clear missing or stale Antigravity workspace projects.

## Usage

```bash
gitmap agy clear [flags]
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--except`, `-e` | `""` | Comma-separated list or file path of IDs, names, slugs, or starts-with prefixes to preserve |
| `--missing`, `-m` | `true` | Only remove projects whose folders no longer exist on disk |
| `--dry-run`, `-d` | `false` | Preview clearance in a formatted table with IDs without deleting |
| `--yes`, `-y` | `false` | Skip confirmation prompt |

## Preview Table with IDs

When `gitmap agy clear` is run, it displays a structured preview table with IDs and directory slugs so you can easily identify what to preserve:

```text
  Targeting 4 project(s) to clear:

    ID           NAME                     SLUG                 PATH
    ────────────────────────────────────────────────────────────────────────────────────────
    1a9408cc     repo                     repo                 C:\Users\Alim\AppData\Local\Temp\...\repo
    6207fc01     repo                     repo                 C:\Users\Alim\AppData\Local\Temp\...\repo

    Tip: Exclude items using: --except "<id, name, slug, or starts-with text>"
```

## Flag Examples

### 1. Preview Clearance with Table (`--dry-run`, `-d`)
```bash
gitmap agy clear --dry-run
```

### 2. Exclude Based on ID, Slug, or Starts-With Text (`--except`, `-e`)
```bash
# Exclude by ID prefix, directory slug, or name starts-with prefix:
gitmap agy clear --except "1a9408cc, repo, wp-"
```

### 3. Clear Using Whitelist CSV File
```bash
gitmap agy clear --except "my-whitelist.csv"
```

### 4. Non-Interactive Headless Mode (`--yes`, `-y`)
```bash
gitmap agy clear --yes
```
