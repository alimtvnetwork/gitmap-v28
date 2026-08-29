#!/usr/bin/env python3
"""
check-error-management.py
Linter enforcing repository-wide error management rules:
1. Zero bare panic(...) in Go source code (outside cliexit/handle.go and tests).
2. Zero bare os.Exit(...) in Go source code (outside cliexit/handle.go and tests).
3. Zero empty catch {} or except: pass blocks.
4. Zero unhandled swallowed error assignments (_ = err) without explanation.
"""

import os
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

ROOT_DIR = Path(__file__).resolve().parent.parent

EXCLUDED_DIRS = {
    ".git",
    "node_modules",
    "dist",
    "build",
    "bin",
    ".gemini",
    "vendor",
    "tmp",
    "linter-scripts",
    "scripts",
    "gitmap-updater",
}

# Whitelisted files where os.Exit or panic are intentionally managed
WHITELISTED_GO_FILES = {
    "handle.go",
    "handle_test.go",
    "cliexit.go",
    "cliexit_test.go",
    "report.go",
    "report_test.go",
    "kind.go",
    "kind_test.go",
    "main.go",
}


def should_scan_dir(dirpath: Path) -> bool:
    for part in dirpath.parts:
        if part in EXCLUDED_DIRS or part.startswith("."):
            if part not in {".lovable", ".github", ".agents"}:
                return False
    return True


def check_go_file(filepath: Path) -> list[str]:
    violations = []
    if filepath.name in WHITELISTED_GO_FILES or filepath.name.endswith("_test.go"):
        return violations

    try:
        content = filepath.read_text(encoding="utf-8", errors="ignore")
    except Exception as e:
        return [f"{filepath}: Failed to read file: {e}"]

    lines = content.splitlines()
    for idx, line in enumerate(lines, 1):
        stripped = line.strip()
        if stripped.startswith("//") or stripped.startswith("/*"):
            continue

        # Check for bare panic("...") or raw panic(
        if re.search(r'\bpanic\s*\(', line):
            # allow debug comment waiver if explicitly marked
            if "// lint-allow: panic" not in line:
                violations.append(f"{filepath}:{idx}: Bare panic() call detected: '{stripped}'")

        # Check for bare os.Exit(
        if re.search(r'\bos\.Exit\s*\(', line):
            if "// lint-allow: os.Exit" not in line:
                violations.append(f"{filepath}:{idx}: Bare os.Exit() call detected (use cliexit.HandleError): '{stripped}'")

        # Check for swallowed error _ = err
        if re.search(r'_\s*=\s*err\b', line):
            if "// lint-allow: ignore-error" not in line:
                violations.append(f"{filepath}:{idx}: Swallowed error '_ = err' without waiver comment: '{stripped}'")

    return violations


def check_ts_js_file(filepath: Path) -> list[str]:
    violations = []
    if filepath.name.endswith(".test.ts") or filepath.name.endswith(".spec.ts"):
        return violations

    try:
        content = filepath.read_text(encoding="utf-8", errors="ignore")
    except Exception as e:
        return [f"{filepath}: Failed to read file: {e}"]

    lines = content.splitlines()
    for idx, line in enumerate(lines, 1):
        stripped = line.strip()
        if stripped.startswith("//") or stripped.startswith("/*"):
            continue

        # Check for empty catch block
        if re.search(r'catch\s*(\([^)]*\))?\s*\{\s*\}', line):
            violations.append(f"{filepath}:{idx}: Empty catch block detected: '{stripped}'")

    return violations


def check_py_file(filepath: Path) -> list[str]:
    violations = []
    if filepath.name.startswith("test_") or filepath.name.endswith("_test.py"):
        return violations

    try:
        content = filepath.read_text(encoding="utf-8", errors="ignore")
    except Exception as e:
        return [f"{filepath}: Failed to read file: {e}"]

    lines = content.splitlines()
    for idx, line in enumerate(lines, 1):
        stripped = line.strip()
        if stripped.startswith("#"):
            continue

        # Check for except: pass or except Exception: pass
        if re.search(r'except(\s+[^:]+)?:\s*pass\b', stripped):
            if "# lint-allow: empty-except" not in line:
                violations.append(f"{filepath}:{idx}: Silent 'except: pass' block detected: '{stripped}'")

    return violations


def main() -> int:
    all_violations = []
    scanned_count = 0

    for root, dirs, files in os.walk(ROOT_DIR):
        root_path = Path(root)
        dirs[:] = [d for d in dirs if d not in EXCLUDED_DIRS and not d.startswith(".")]

        for f in files:
            filepath = root_path / f
            ext = filepath.suffix.lower()

            if ext == ".go":
                scanned_count += 1
                all_violations.extend(check_go_file(filepath))
            elif ext in {".ts", ".tsx", ".js", ".jsx"}:
                scanned_count += 1
                all_violations.extend(check_ts_js_file(filepath))
            elif ext == ".py":
                scanned_count += 1
                all_violations.extend(check_py_file(filepath))

    print(f"Scanned {scanned_count} source files for error management compliance.")

    if all_violations:
        print(f"\n❌ FAILED: Found {len(all_violations)} error management violation(s):")
        for v in all_violations:
            print(f"  - {v}")
        return 1

    print("✅ PASS: All error management checks passed (zero bare panics, exits, or swallowed errors).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
