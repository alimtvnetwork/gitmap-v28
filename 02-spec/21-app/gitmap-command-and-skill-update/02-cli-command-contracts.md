# CLI Command Contracts: `gitmap prnt` Suite

## 1. Syntax & Subcommand Signatures

```bash
# Single URL download with automatic brain sync
gitmap prnt <url> [destination-path]

# Multi-URL batch download (parallel fetching)
gitmap prnt <url1> <url2> <url3>...

# URL extraction only (prints resolved direct image links without downloading)
gitmap prnt <url1> <url2>... --extract-only
gitmap prnt <url1> <url2>... -e

# Structured JSON envelope output for AI automation
gitmap prnt <url1> <url2>... --json
gitmap prnt <url1> <url2>... -j

# Custom output directory for batch downloads
gitmap prnt <url1> <url2>... -o <directory>
gitmap prnt <url1> <url2>... --output-dir <directory>

# Aliases
gitmap download-prnt ...
gitmap prnt-download ...
gitmap asset download-prnt ...
```

---

## 2. Command Flags & Parameters

| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--extract-only` | `-e`, `-p` | Boolean | `false` | Resolve and print direct image URLs without writing files to disk |
| `--json` | `-j` | Boolean | `false` | Output results in standard Universal JSON Response Envelope |
| `--output-dir` | `-o` | String | `assets/screenshots` | Target directory for saving downloaded screenshot images |
| `--concurrency` | `-c` | Integer | `4` | Maximum parallel worker routines for batch fetching |
| `--help` | `-h` | Boolean | `false` | Display command help and usage examples |

---

## 3. Terminal Output Contracts

### Standard Batch Download Mode:
```text
✔ Successfully processed 3 screenshot asset(s):
  [1/3] https://prnt.sc/FnF7IMTncnd5
        Direct URL: https://img.lightshot.app/3_iJS9o8Sjm1aTq2i1lowQ.png
        Saved:      assets/screenshots/FnF7IMTncnd5_20261020_110409.png (48,991 B)
        Brain Sync: .gemini/antigravity/brain/<id>/screenshot_FnF7IMTncnd5.png
  [2/3] https://prnt.sc/oxgHWyUQkiWW
        Direct URL: https://i.imgur.com/example2.png
        Saved:      assets/screenshots/oxgHWyUQkiWW_20261020_110410.png (54,120 B)
        Brain Sync: .gemini/antigravity/brain/<id>/screenshot_oxgHWyUQkiWW.png
Completed 3 asset(s) in 320ms.
```

### Extraction Only Mode (`--extract-only`):
```text
https://prnt.sc/FnF7IMTncnd5 -> https://img.lightshot.app/3_iJS9o8Sjm1aTq2i1lowQ.png
https://prnt.sc/oxgHWyUQkiWW -> https://i.imgur.com/example2.png
Completed in 140ms.
```

### Universal JSON Envelope (`--json`):
```json
{
  "Status": {
    "IsSuccess": true,
    "IsFailed": false,
    "Code": 200,
    "Message": "OK",
    "Timestamp": "2026-10-20T03:00:00Z"
  },
  "Attributes": {
    "TotalRequested": 2,
    "TotalSucceeded": 2,
    "TotalFailed": 0,
    "ExtractOnly": false
  },
  "Results": [
    {
      "SourceURL": "https://prnt.sc/FnF7IMTncnd5",
      "DirectURL": "https://img.lightshot.app/3_iJS9o8Sjm1aTq2i1lowQ.png",
      "Slug": "FnF7IMTncnd5",
      "SavedPath": "assets/screenshots/FnF7IMTncnd5_20261020_110409.png",
      "BrainSyncPath": "C:\\Users\\...\\screenshot_FnF7IMTncnd5.png",
      "FileSizeBytes": 48991,
      "IsSuccess": true,
      "Error": ""
    }
  ],
  "Errors": null
}
```
