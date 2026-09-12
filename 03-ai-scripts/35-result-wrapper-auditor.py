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

# Regex for detecting legacy multi-value slice return tuples with error
# Excludes []byte and []rune which represent low-level stream/serialization primitives
SLICE_TUPLE_RETURN = re.compile(
    r"func\s+(?:\([^)]+\)\s+)?(\w+)\s*\([^)]*\)\s*\(\s*\[\](?!(?:byte|rune)\b)[^,]+,\s*(?:error|\*apperror\.AppError)\s*\)"
)

# Regex for detecting value receiver declarations on Result types
VALUE_RECEIVER_RESULT = re.compile(
    r"func\s+\((\w+)\s+(Result(?:Slice|Map)?\[[^\]]+\])\)"
)

# Regex for detecting clumsy compound cardinality checks
CLUMSY_CARDINALITY_CHECK = re.compile(
    r"(?:\.IsFailure\(\)\s*\|\|\s*[^.\s]+\.Count\(\)\s*!=\s*\d+|[^.\s]+\.Count\(\)\s*!=\s*\d+\s*\|\|\s*[^.\s]+\.IsFailure\(\))"
)

# Enforced subsystems that must define domain models & Result aliases in types.go
ENFORCED_TYPES_GO_PACKAGES = (
    "cli/result",
    "cli/cmdschedule",
    "cli/pipelinedb",
    "cli/macro",
    "cli/db",
    "cli/cluster",
    "cli/cmdpurge",
    "cli/clonenow",
    "cli/clonefrom",
    "cli/cmdprompt",
    "cli/downloaderconfig",
    "cli/movemerge",
    "cli/lazyregex",
    "cli/archive",
)

NON_AFFIRMATIVE_DEFINED = re.compile(r"\bdefined\s+bool\b")

# Enforced subsystems/files that must strictly use ResultSlice[T]
RESULT_SLICE_ENFORCED_PREFIXES = (
    "cli/macro/",
    "cli/pipelinedb/",
    "cli/cmdprompt/",
    "cli/cmdschedule/",
    "cli/cluster/pathalias.go",
    "cli/db/",
)

EXCLUDED_FUNCTIONS = {
    "Unwrap", "MarshalJSON", "UnmarshalJSON", "Read", "Write",
}


def audit_file(filepath: Path) -> tuple[list[str], int]:
    violations = []
    unmigrated_slices = 0
    try:
        content = filepath.read_text(encoding="utf-8", errors="replace")
    except Exception as e:
        return [f"{filepath}: error reading file: {e}"], 0

    rel = filepath.relative_to(ROOT_DIR).as_posix()
    is_slice_enforced = any(rel.startswith(p) or rel == p for p in RESULT_SLICE_ENFORCED_PREFIXES)

    for lno, line in enumerate(content.splitlines(), 1):
        stripped = line.strip()
        if stripped.startswith("//") or stripped.startswith("/*"):
            continue

        map_match = MAP_TUPLE_RETURN.search(stripped)
        if map_match:
            fn_name = map_match.group(1)
            if fn_name not in EXCLUDED_FUNCTIONS:
                violations.append(
                    f"{rel}:{lno} function `{fn_name}` returns multi-value map tuple with error; "
                    f"must return result.ResultMap[K, V]"
                )

        slice_match = SLICE_TUPLE_RETURN.search(stripped)
        if slice_match:
            fn_name = slice_match.group(1)
            if fn_name not in EXCLUDED_FUNCTIONS:
                if is_slice_enforced:
                    violations.append(
                        f"{rel}:{lno} function `{fn_name}` returns multi-value slice tuple with error; "
                        f"must return result.ResultSlice[T]"
                    )
                else:
                    unmigrated_slices += 1

        val_match = VALUE_RECEIVER_RESULT.search(stripped)
        if val_match:
            typ_name = val_match.group(2)
            violations.append(
                f"{rel}:{lno} method declared on value receiver `{typ_name}`; "
                f"must use pointer receiver `(r *{typ_name})` for null safety"
            )

        clumsy_match = CLUMSY_CARDINALITY_CHECK.search(stripped)
        if clumsy_match:
            violations.append(
                f"{rel}:{lno} clumsy compound cardinality check; "
                f"use `.IsCountOtherThan(N)` instead"
            )

    return violations, unmigrated_slices


def check_types_go_centralization() -> list[str]:
    violations = []
    for pkg_rel in ENFORCED_TYPES_GO_PACKAGES:
        pkg_dir = ROOT_DIR / pkg_rel
        types_file = pkg_dir / "types.go"
        if not types_file.exists():
            violations.append(f"{pkg_rel}/types.go is missing; package must define domain models & Result aliases in types.go")

    for rpath in ("cli/result/result.go", "cli/result/types.go"):
        rfile = ROOT_DIR / rpath
        if rfile.exists():
            content = rfile.read_text(encoding="utf-8", errors="replace")
            for lno, line in enumerate(content.splitlines(), 1):
                if NON_AFFIRMATIVE_DEFINED.search(line):
                    violations.append(f"{rpath}:{lno} non-affirmative boolean field `defined bool`; must be `isDefined bool`")

    sched_export = ROOT_DIR / "cli/cmdschedule/schedule_export.go"
    if sched_export.exists():
        content = sched_export.read_text(encoding="utf-8", errors="replace")
        for lno, line in enumerate(content.splitlines(), 1):
            if "type scheduleExportBundle struct" in line or "type scheduleExportOpts struct" in line:
                violations.append(f"cli/cmdschedule/schedule_export.go:{lno} unexported inline struct; must be exported in cli/cmdschedule/types.go")

    return violations


def main() -> None:
    if not CLI_DIR.exists():
        print(f"Error: {CLI_DIR} not found.")
        sys.exit(1)

    all_violations = []
    all_violations.extend(check_types_go_centralization())
    total_unmigrated_slices = 0
    file_count = 0

    for root, dirs, files in os.walk(CLI_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for f in files:
            if f.endswith(".go"):
                is_test = f.endswith("_test.go")
                file_count += 1
                fp = Path(root) / f
                violations, unmigrated = audit_file(fp)
                if is_test:
                    # In test files, only report clumsy cardinality checks
                    test_violations = [v for v in violations if "clumsy" in v]
                    all_violations.extend(test_violations)
                else:
                    all_violations.extend(violations)
                    total_unmigrated_slices += unmigrated

    print(f"Audited {file_count} Go files in {CLI_DIR.name}/ for ResultMap/ResultSlice compliance.")
    print(f"ResultSlice enforced subsystems: {', '.join(RESULT_SLICE_ENFORCED_PREFIXES)}")
    print(f"Backlog slice returns pending future migrations across other packages: {total_unmigrated_slices}")

    if not all_violations:
        print("\n✅ PASS: Zero legacy map/slice error return tuples in enforced subsystems. Result envelopes verified.")
        sys.exit(0)

    print(f"\n❌ FAIL: Found {len(all_violations)} violation(s):")
    for v in all_violations:
        print(f"  {v}")
    sys.exit(1)


if __name__ == "__main__":
    main()
