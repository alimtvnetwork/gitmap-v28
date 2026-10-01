---
name: python-dry-caching
description: Design, implement, and maintain high-performance, cached Python automation scripts leveraging 03-ai-scripts/02-shared-engine.py.
---

# Python DRY Architecture & Caching Specification

This skill guides the design, implementation, and maintenance of high-performance Python scripts in `03-ai-scripts/`.

## Architecture & Shared Engine

All standalone automation and inspection scripts in `03-ai-scripts/` must adhere to DRY principles:

1. **Centralized Engine (`03-ai-scripts/02-shared-engine.py`):**
   - Reuse shared constants, pre-compiled regexes (`RegexPatternType`), and directory resolution.
   - Never duplicate filesystem walking or cache management across scripts.

2. **Sub-Millisecond Caching (`tmp/cache/`):**
   - Utilize 2-tier caching (in-memory lookup followed by persistent disk serialization).
   - Scripts like `03-ai-scripts/17-fast-file-reader.py` and `03-ai-scripts/12-fast-cached-grep.py` provide <15ms exploration over large directory trees.

3. **Strict Python Standards:**
   - Standard library only (`os`, `sys`, `pathlib`, `re`, `argparse`, `json`).
   - UTF-8 standard output enforcement: `sys.stdout.reconfigure(encoding="utf-8")`.
   - Explicit return codes: `sys.exit(0)` on success, `sys.exit(1)` on violation.

## Key CLI Tools

- `python 03-ai-scripts/17-fast-file-reader.py --list-folder <path>`: Rapid cached folder discovery.
- `python 03-ai-scripts/17-fast-file-reader.py --read-file <path> --max-bytes <N>`: Fast bounded file reader.
- `python 03-ai-scripts/06-cicd-local-runner.py`: Executes all 36 local quality gates in parallel.
