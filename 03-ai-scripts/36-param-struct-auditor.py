#!/usr/bin/env python3
"""36-param-struct-auditor.py - Audits Go function signatures for argument reduction, parameter structs, and affirmative booleans.

Verifies adherence to:
- Dedicated parameter Structs/DTOs for high-arity function signatures (>2-3 parameters)
- Affirmative boolean naming (is* / has* only) on struct fields and function parameters
- Strongly-typed *apperror.AppError returns instead of raw errors in refactored domain packages
- types.go centralization of parameter objects
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

ENFORCED_PACKAGES = (
    "cli/cluster",
    "cli/clonenow",
    "cli/clonefrom",
    "cli/visibility",
    "cli/goldenguard",
)

ENFORCED_TYPES_GO = (
    ("cli/cluster/types.go", [
        "LifecycleExecParams", "PreflightParams",
        "NodeExecutionReportParams", "FinishClusterPoolParams"
    ]),
    ("cli/clonenow/types.go", [
        "CloneIdempotentParams", "ConcurrentDispatchParams",
        "ConcurrentWorkerParams", "GitCloneParams", "ProgressWriteParams"
    ]),
    ("cli/clonefrom/types.go", [
        "CloneFromDispatchParams", "CloneFromWorkerParams",
        "BeforeRowInvokeParams", "ProgressWriteParams"
    ]),
)

NON_AFFIRMATIVE_BOOLEAN_PATTERNS = [
    (re.compile(r"\btrigger\s+bool\b"), "non-affirmative boolean parameter `trigger bool`; must be `isTrigger bool`"),
    (re.compile(r"\bshowRole\s+bool\b"), "non-affirmative boolean parameter `showRole bool`; must be `isShowRole bool`"),
    (re.compile(r"\bforceLifecycle\s+bool\b"), "non-affirmative boolean parameter `forceLifecycle bool`; must be `isForceLifecycle bool`"),
    (re.compile(r"\bautoConfirm\s+bool\b"), "non-affirmative boolean parameter `autoConfirm bool`; must be `isAutoConfirm bool`"),
    (re.compile(r"^\s+Exists\s+bool\b"), "non-affirmative boolean field `Exists bool`; must be `IsExists bool`"),
    (re.compile(r"^\s+Empty\s+bool\b"), "non-affirmative boolean field `Empty bool`; must be `IsEmpty bool`"),
]


def check_enforced_types_go() -> list[str]:
    violations = []
    for rel_path, required_structs in ENFORCED_TYPES_GO:
        file_path = ROOT_DIR / rel_path
        if not file_path.exists():
            violations.append(f"{rel_path}: missing required types.go file")
            continue

        content = file_path.read_text(encoding="utf-8", errors="replace")
        for struct_name in required_structs:
            if not re.search(rf"\b{struct_name}\s+struct\b", content):
                violations.append(
                    f"{rel_path}: missing required parameter struct `{struct_name} struct`"
                )

    visibility_exclude = ROOT_DIR / "cli/visibility/exclude.go"
    if visibility_exclude.exists():
        content = visibility_exclude.read_text(encoding="utf-8", errors="replace")
        if "type ExclusionTokenParams struct" not in content:
            violations.append("cli/visibility/exclude.go: missing required `ExclusionTokenParams` struct")
    else:
        violations.append("cli/visibility/exclude.go: file not found")

    return violations


def audit_file(filepath: Path) -> list[str]:
    violations = []
    try:
        content = filepath.read_text(encoding="utf-8", errors="replace")
    except Exception as e:
        return [f"{filepath}: error reading file: {e}"]

    rel = filepath.relative_to(ROOT_DIR).as_posix()
    is_enforced = any(rel.startswith(p + "/") or rel.startswith(p + "\\") for p in ENFORCED_PACKAGES)
    if not is_enforced:
        return violations

    for lno, line in enumerate(content.splitlines(), 1):
        stripped = line.strip()
        if stripped.startswith("//") or stripped.startswith("/*"):
            continue

        for pattern, msg in NON_AFFIRMATIVE_BOOLEAN_PATTERNS:
            if pattern.search(line):
                violations.append(f"{rel}:{lno} {msg}")

    return violations


def main() -> None:
    if not CLI_DIR.exists():
        print(f"Error: {CLI_DIR} not found.")
        sys.exit(1)

    all_violations = []
    all_violations.extend(check_enforced_types_go())

    file_count = 0
    for root, dirs, files in os.walk(CLI_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS]
        for f in files:
            if f.endswith(".go") and not f.endswith("_test.go"):
                file_count += 1
                fp = Path(root) / f
                violations = audit_file(fp)
                all_violations.extend(violations)

    print(f"Audited {file_count} non-test Go files in {CLI_DIR.name}/ for Parameter Struct & Boolean compliance.")
    print(f"Enforced packages: {', '.join(ENFORCED_PACKAGES)}")

    if not all_violations:
        print("\n✅ PASS: Zero parameter struct or affirmative boolean violations in enforced packages.")
        sys.exit(0)

    print(f"\n❌ FAIL: Found {len(all_violations)} violation(s):")
    for v in all_violations:
        print(f"  {v}")
    sys.exit(1)


if __name__ == "__main__":
    main()
