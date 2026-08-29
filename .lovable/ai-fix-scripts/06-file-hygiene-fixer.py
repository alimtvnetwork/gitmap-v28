#!/usr/bin/env python3
"""06-file-hygiene-fixer.py - Enforces Unix LF line endings, UTF-8 without BOM, and single trailing newline at EOF."""
import os
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent.parent
EXCLUDE_DIRS = {'.git', 'node_modules', 'dist', 'build', '.next', '.gitmap', 'vendor', 'brain', '.gemini'}
TARGET_EXTS = {
    '.go', '.ts', '.tsx', '.js', '.jsx', '.json', '.md', '.py', '.sh', '.ps1',
    '.css', '.html', '.yaml', '.yml', '.toml', '.sql'
}


def sanitize_file(filepath: Path) -> bool:
    try:
        raw_bytes = filepath.read_bytes()
    except Exception:
        return False

    has_bom = raw_bytes.startswith(b'\xef\xbb\xbf')
    if has_bom:
        raw_bytes = raw_bytes[3:]

    try:
        text = raw_bytes.decode('utf-8')
    except Exception:
        return False

    # Check CRLF
    has_crlf = '\r' in text
    if has_crlf:
        text = text.replace('\r\n', '\n').replace('\r', '\n')

    # Strip multiple trailing newlines and ensure exactly one trailing newline
    trimmed = text.rstrip(' \t\r\n')
    new_text = trimmed + '\n' if trimmed else ''

    new_bytes = new_text.encode('utf-8')
    if new_bytes == raw_bytes:
        return False

    filepath.write_bytes(new_bytes)
    return True


def process_single_file(fp: Path) -> tuple[int, int]:
    if fp.suffix.lower() not in TARGET_EXTS:
        return 0, 0
    is_modified = sanitize_file(fp)
    return 1, (1 if is_modified else 0)


def process_target_dir(td: Path) -> tuple[int, int]:
    if td.is_file():
        return process_single_file(td)

    scanned = 0
    modified = 0
    for root, dirs, files in os.walk(td):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for file in files:
            fp = Path(root) / file
            s_inc, m_inc = process_single_file(fp)
            scanned += s_inc
            modified += m_inc
    return scanned, modified


def main():
    target_dirs = [ROOT_DIR]
    if len(sys.argv) > 1:
        target_dirs = [Path(p) for p in sys.argv[1:]]

    total_scanned = 0
    total_modified = 0

    for td in target_dirs:
        s, m = process_target_dir(td)
        total_scanned += s
        total_modified += m

    print(f"File Hygiene Fixer: Scanned {total_scanned} files, normalized {total_modified} file(s) to LF/UTF-8.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
