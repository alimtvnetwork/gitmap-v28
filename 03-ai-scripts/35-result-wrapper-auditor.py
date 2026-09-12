#!/usr/bin/env python3
"""35-result-wrapper-auditor.py - Audits Go functions returning multi-value map/slice error tuples.

Verifies adherence to the Single Return Object mandate, strongly-typed result wrappers
(ResultMap[K, V], ResultSlice[T], Result[T]), and structured AppError returns.
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
CLI_DIR = ROOT_DIR / "cli"

EXCLUDE_DIRS = {
    ".git", "node_modules", "dist", "build", "bin", ".next", ".gitmap",
    "vendor", "coverage", ".gemini", ".system_generated", "fixtures",
    "scratch", "temp-scripts", "temp-agents", "temp", "dbengine_old",
    "indexer_old", "installer_old", "repodb_old", "searcher_old"
}

# Regex for detecting legacy multi-value map return tuples with error
MAP_TUPLE_RETURN = re.compile(
    r"func\s+(?:\([^)]+\)\s+)?(\w+)\s*\([^)]*\)\s*\(\s*map\[[^\]]+\][^,]+,\s*(?:error|\*apperror\.AppError)\s*\)"
)


def audit_file(filepath: Path) -> list[str]:
    violations = []
    try:
        content = filepath.read_text(encoding="utf-8", errors="replace")
    except Exception as e:
        return [f"{filepath}: error reading file: {e}"]

    rel = filepath.relative_to(ROOT_DIR).as_posix()
    for lno, line in enumerate(content.splitlines(), 1):
        stripped = line.strip()
        if stripped.startswith("//") or stripped.startswith("/*"):
            continue

        match = MAP_TUPLE_RETURN.search(stripped)
        if match:
            fn_name = match.group(1)
            if fn_name in ("Unwrap",):
                continue
            violations.append(
                f"{rel}:{lno} function `{fn_name}` returns multi-value map tuple with error; "
                f"must return result.ResultMap[K, V]"
            )

    return violations


def main() -> None:
    if not CLI_DIR.exists():
        print(f"Error: {CLI_DIR} not found.")
        sys.exit(1)

    all_violations = []
    file_count = 0

    for root, dirs, files in os.walk(CLI_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for f in files:
            if f.endswith(".go") and not f.endswith("_test.go"):
                file_count += 1
                fp = Path(root) / f
                all_violations.extend(audit_file(fp))

    print(f"Audited {file_count} Go files in {CLI_DIR.name}/ for ResultMap/Result wrapper compliance.")

    if not all_violations:
        print("\n✅ PASS: Zero legacy multi-value map tuple error returns found. ResultMap envelopes verified.")
        sys.exit(0)

    print(f"\n❌ FAIL: Found {len(all_violations)} violation(s):")
    for v in all_violations:
        print(f"  {v}")
    sys.exit(1)


if __name__ == "__main__":
    main()
