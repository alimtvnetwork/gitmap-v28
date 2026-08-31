# sequence

Manages, lists, and re-sequences numbered files (e.g. `01-index.md`, `02-shared-engine.py`, `03-file-manipulator.py`) with pinning, shifting, zero-padding, and repo-scoped SQLite database persistence.

## Aliases

seq

## Usage

```bash
gitmap sequence [subcommand] [directory] [flags]
gitmap seq [subcommand] [directory] [flags]
```

## Subcommands

| Subcommand           | Description                                                          |
|----------------------|----------------------------------------------------------------------|
| `list`, `ls`         | Scan and display all sequenced/unsequenced files in human or JSON    |
| `fix`, `reorder`     | Re-sequence files with zero-padding, pin mapping, or sequence shifts |
| `get`                | Retrieve saved sequence metadata from the repo-scoped SQLite DB      |
| `history`, `hist`    | View sequence operation audit history from the repo-scoped SQLite DB |
| `help`               | Display this help documentation                                      |

## Flags

| Flag              | Type    | Default | Description                                                     |
|-------------------|---------|---------|-----------------------------------------------------------------|
| `--json`          | boolean | false   | Output structured JSON payload for AI agents and scripts        |
| `--dry-run`       | boolean | false   | Preview proposed renames without changing files                 |
| `--pin`           | string  | ""      | Pin specific file base names to exact sequence numbers          |
| `--start`         | integer | 1       | Starting sequence index (e.g. 1 for 01-, 0 for 00-)             |
| `--shift`         | integer | 0       | Increment/decrement all unpinned sequence numbers by an offset  |
| `--order-by-time` | boolean | false   | Order unpinned files by filesystem modification time            |
| `--order-by-az`   | boolean | false   | Order unpinned files alphabetically by text suffix              |

## Why use this instead of shell scripts?

Manually renaming dozens of sequential files in PowerShell or Bash is error-prone, risks overwriting collisions, and breaks git history.

Instead of writing custom loops:
```powershell
Get-ChildItem .lovable/ai-fix-scripts | ForEach-Object { ... }
```

Use GitMap's atomic sequence manager with Git history preservation and SQLite tracking:
```bash
gitmap seq fix .lovable/ai-fix-scripts --pin "index.md=01,shared-engine.py=02" --start 1
```

## Examples

### List Sequenced Files in Human-Readable Table

```bash
$ gitmap sequence list .lovable/ai-fix-scripts

Directory Sequence: .lovable/ai-fix-scripts (13 files, 11 sequenced)
SEQ     FILENAME                             BASE NAME
---------------------------------------------------------------------------
01      01-index.md                          index.md
02      02-shared-engine.py                  shared-engine.py
03      03-file-manipulator.py               file-manipulator.py
04      04-guideline-autofixer.py            guideline-autofixer.py
```

### Inspect Sequences as Machine-Readable JSON for AI Agents

```bash
$ gitmap seq list .lovable/ai-fix-scripts --json
{
  "directory": ".lovable/ai-fix-scripts",
  "totalFiles": 13,
  "sequencedFiles": 13,
  "files": [
    {
      "sequence": 1,
      "filename": "01-index.md",
      "baseName": "index.md",
      "extension": ".md",
      "path": ".lovable/ai-fix-scripts/01-index.md"
    },
    {
      "sequence": 2,
      "filename": "02-shared-engine.py",
      "baseName": "shared-engine.py",
      "extension": ".py",
      "path": ".lovable/ai-fix-scripts/02-shared-engine.py"
    }
  ]
}
```

### Pin Top Files and Sequence Subsequent Files

Pin `index.md` to `01` and `shared-engine.py` to `02`, with all other files sequentially numbered after them:

```bash
gitmap seq fix .lovable/ai-fix-scripts --pin "index=01,shared-engine=02" --start 1
```

### Preview Renames with Dry Run

```bash
gitmap seq fix .lovable/ai-fix-scripts --pin "index=01,shared-engine=02" --dry-run
```

### Retrieve Saved Sequences from Repo-Scoped SQLite Database

```bash
gitmap seq get .lovable/ai-fix-scripts
```

## See Also

- [folder](folder.md) — Export folder structure in text, markdown, JSON, or YAML
- [fix-seq-files](fix-seq-files.md) — Low-level file sequence utility
- [llm](llm.md) — Capability document for automated AI agent integrations
