#!/usr/bin/env python3
"""
33-test-inventory-generator.py
==============================
Generates and maintains the centralized test inventory manifest at `.lovable/test-inventory.json`
and provides safe, atomic file change recording into `.lovable/temp/recent-file-changes.json`
with file locking to ensure concurrency safety across multi-agent turns.

Usage:
  # Scan and generate / update test inventory:
  python 03-ai-scripts/33-test-inventory-generator.py

  # Record modified files safely under lock:
  python 03-ai-scripts/33-test-inventory-generator.py --record "gitmap/cmd/root.go"

  # Query tests associated with recent changes:
  python 03-ai-scripts/33-test-inventory-generator.py --query-recent

  # Clear recent changes log:
  python 03-ai-scripts/33-test-inventory-generator.py --clear
"""

import argparse
import contextlib
import datetime
import hashlib
import json
import os
from pathlib import Path
import re
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

REPO_ROOT = Path(__file__).resolve().parent.parent
LOVABLE_DIR = REPO_ROOT / ".lovable"
TEMP_DIR = LOVABLE_DIR / "temp"
TEST_INVENTORY_PATH = LOVABLE_DIR / "test-inventory.json"
RECENT_CHANGES_PATH = TEMP_DIR / "recent-file-changes.json"
LOCK_FILE_PATH = TEMP_DIR / "recent-file-changes.lock"

FUNC_START_RE = re.compile(r"^func\s+(?:\([^)]+\)\s+)?([A-Za-z0-9_]+)\s*\(")
TEST_START_RE = re.compile(r"^func\s+(Test[A-Za-z0-9_]*)\s*\(")
PY_TEST_START_RE = re.compile(r"^def\s+(test_[A-Za-z0-9_]*)\s*\(")
TS_TEST_START_RE = re.compile(r"""(?:it|test)\s*\(\s*["'`]([^"'`]+)["'`]""")


def compute_file_hash(filepath: Path) -> str:
    """Computes a 16-character SHA-256 hash for a file."""
    if not filepath.is_file():
        return ""
    try:
        data = filepath.read_bytes()
        return hashlib.sha256(data).hexdigest()[:16]
    except OSError:
        return ""


def normalize_repo_rel(path_input: str | Path) -> str:
    """Normalizes any path to a forward-slash relative path from the repository root."""
    raw_str = str(path_input).strip()
    try:
        p = Path(raw_str)
        if p.is_absolute():
            rel = p.resolve().relative_to(REPO_ROOT.resolve())
            raw_str = str(rel)
    except Exception:
        pass
    norm = raw_str.replace("\\", "/").strip()
    if norm.startswith("./"):
        norm = norm[2:]
    return norm


@contextlib.contextmanager
def file_lock(lock_path: Path, timeout_sec: float = 10.0):
    """Acquires a cross-platform cooperative lock file with timeout."""
    lock_path.parent.mkdir(parents=True, exist_ok=True)
    start_time = time.time()
    pid = os.getpid()
    acquired = False

    while not acquired:
        try:
            fd = os.open(str(lock_path), os.O_CREAT | os.O_EXCL | os.O_RDWR)
            os.write(fd, f"{pid}\n".encode("utf-8"))
            os.close(fd)
            acquired = True
        except FileExistsError:
            if time.time() - start_time > timeout_sec:
                try:
                    lock_path.unlink()
                except OSError:
                    pass
            time.sleep(0.05)

    try:
        yield
    finally:
        try:
            lock_path.unlink()
        except OSError:
            pass


def atomic_write_json(filepath: Path, data: Any) -> None:
    """Safely writes JSON data using a temporary file rename."""
    filepath.parent.mkdir(parents=True, exist_ok=True)
    temp_file = filepath.with_suffix(f".tmp.{os.getpid()}")
    try:
        with open(temp_file, "w", encoding="utf-8", newline="\n") as f:
            json.dump(data, f, indent=2)
            f.write("\n")
        temp_file.replace(filepath)
    except Exception:
        if temp_file.exists():
            try:
                temp_file.unlink()
            except OSError:
                pass
        raise


def extract_go_declarations(filepath: Path) -> dict[str, str]:
    """Extracts function/method declarations and their line spans."""
    funcs: dict[str, str] = {}
    if not filepath.is_file():
        return funcs

    try:
        lines = filepath.read_text(encoding="utf-8", errors="replace").splitlines()
    except Exception:
        return funcs

    current_func = None
    current_lines = []

    for line in lines:
        m = FUNC_START_RE.match(line)
        if m:
            if current_func:
                chunk = "\n".join(current_lines).encode("utf-8")
                funcs[current_func] = hashlib.sha256(chunk).hexdigest()[:16]
            current_func = m.group(1)
            current_lines = [line]
        elif current_func:
            current_lines.append(line)

    if current_func:
        chunk = "\n".join(current_lines).encode("utf-8")
        funcs[current_func] = hashlib.sha256(chunk).hexdigest()[:16]

    return funcs


def extract_go_tests(filepath: Path) -> dict[str, str]:
    """Extracts Go test functions (Test*) from a test file."""
    tests: dict[str, str] = {}
    if not filepath.is_file():
        return tests

    try:
        lines = filepath.read_text(encoding="utf-8", errors="replace").splitlines()
    except Exception:
        return tests

    current_test = None
    current_lines = []

    for line in lines:
        m = TEST_START_RE.match(line)
        if m:
            if current_test:
                chunk = "\n".join(current_lines).encode("utf-8")
                tests[current_test] = hashlib.sha256(chunk).hexdigest()[:16]
            current_test = m.group(1)
            current_lines = [line]
        elif current_test:
            current_lines.append(line)

    if current_test:
        chunk = "\n".join(current_lines).encode("utf-8")
        tests[current_test] = hashlib.sha256(chunk).hexdigest()[:16]

    return tests


def extract_python_tests(filepath: Path) -> dict[str, str]:
    """Extracts Python test functions (test_*) from a test file."""
    tests: dict[str, str] = {}
    if not filepath.is_file():
        return tests
    try:
        lines = filepath.read_text(encoding="utf-8", errors="replace").splitlines()
    except Exception:
        return tests

    current_test = None
    current_lines = []
    for line in lines:
        m = PY_TEST_START_RE.match(line.strip())
        if m:
            if current_test:
                chunk = "\n".join(current_lines).encode("utf-8")
                tests[current_test] = hashlib.sha256(chunk).hexdigest()[:16]
            current_test = m.group(1)
            current_lines = [line]
        elif current_test:
            current_lines.append(line)

    if current_test:
        chunk = "\n".join(current_lines).encode("utf-8")
        tests[current_test] = hashlib.sha256(chunk).hexdigest()[:16]

    return tests


def extract_ts_tests(filepath: Path) -> dict[str, str]:
    """Extracts TypeScript test descriptions (it/test) from a test file."""
    tests: dict[str, str] = {}
    if not filepath.is_file():
        return tests
    try:
        content = filepath.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return tests

    for m in TS_TEST_START_RE.finditer(content):
        tname = m.group(1).strip()
        tests[tname] = hashlib.sha256(tname.encode("utf-8")).hexdigest()[:16]

    return tests


def resolve_go_package_rel(pkg_dir: Path, repo_root: Path) -> str:
    """Resolves relative package path from repo root."""
    rel = pkg_dir.resolve().relative_to(repo_root.resolve()).as_posix()
    return rel


def estimate_test_duration(filepath: Path, test_name: str, content: str) -> tuple[float, str]:
    """Estimates test duration in seconds and categorizes tier (unit vs heavy)."""
    rel = str(filepath).replace("\\", "/")
    if "tests/heavy_test" in rel:
        return 3.0, "heavy"

    is_heavy = False
    duration = 0.005
    if "exec.Command" in content:
        is_heavy = True
        duration += 2.0
    if "time.Sleep" in content:
        is_heavy = True
        duration += 1.0
    if "net.Listen" in content or "http.Get" in content or "http.Post" in content:
        is_heavy = True
        duration += 0.5
    if "git" in test_name.lower() and ("subprocess" in content.lower() or "exec" in content.lower()):
        is_heavy = True
        duration += 1.5

    tier = "heavy" if is_heavy else "unit"
    return round(duration, 3), tier


def scan_go_tests(repo_root: Path) -> tuple[dict[str, Any], int, int]:
    """Scans and indexes all Go test files."""
    tests_dict: dict[str, Any] = {}
    go_test_files = list(repo_root.rglob("*_test.go"))

    for tf in go_test_files:
        rel_test_file = normalize_repo_rel(tf)
        if ".git" in rel_test_file or "node_modules" in rel_test_file or ".lovable" in rel_test_file:
            continue

        pkg_dir = tf.parent
        rel_pkg = normalize_repo_rel(pkg_dir)
        test_funcs = extract_go_tests(tf)
        test_file_hash = compute_file_hash(tf)

        file_content = ""
        try:
            file_content = tf.read_text(encoding="utf-8", errors="replace")
        except Exception:
            pass

        target_file_name = tf.name.replace("_test.go", ".go")
        target_path = pkg_dir / target_file_name
        rel_target = normalize_repo_rel(target_path) if target_path.exists() else ""
        code_hash = compute_file_hash(target_path) if target_path.exists() else test_file_hash

        for tname, thash in test_funcs.items():
            tid = f"{rel_pkg}.{tname}"
            target_func = ""
            if tname.startswith("Test"):
                cand = tname[4:].split("_")[0]
                target_func = cand

            dur, tier = estimate_test_duration(tf, tname, file_content)

            tests_dict[tid] = {
                "id": tid,
                "package": rel_pkg,
                "test_file": rel_test_file,
                "test_func": tname,
                "test_hash": thash,
                "target_file": rel_target,
                "target_func": target_func,
                "code_hash": code_hash,
                "duration_sec": dur,
                "tier": tier,
                "last_status": "never_run",
                "last_run_at": "",
                "needs_run": True,
            }

    return tests_dict, len(tests_dict), len(set(t["package"] for t in tests_dict.values()))


def scan_python_and_ts_tests(repo_root: Path) -> dict[str, Any]:
    """Scans and indexes Python and TypeScript test files."""
    tests_dict: dict[str, Any] = {}
    for root, _, files in os.walk(repo_root):
        rel_dir = normalize_repo_rel(root)
        if ".git" in rel_dir or "node_modules" in rel_dir or ".lovable" in rel_dir or "dist" in rel_dir:
            continue
        for f in files:
            p = Path(root) / f
            rel_file = normalize_repo_rel(p)
            if (f.startswith("test_") or f.endswith("_test.py")) and f.endswith(".py"):
                py_tests = extract_python_tests(p)
                for tname, thash in py_tests.items():
                    tid = f"{rel_dir}.{tname}"
                    tests_dict[tid] = {
                        "id": tid,
                        "package": rel_dir,
                        "test_file": rel_file,
                        "test_func": tname,
                        "test_hash": thash,
                        "target_file": rel_file.replace("test_", "").replace("_test.py", ".py"),
                        "target_func": "",
                        "code_hash": thash,
                        "duration_sec": 0.0,
                        "last_status": "never_run",
                        "last_run_at": "",
                        "needs_run": True,
                    }
            elif f.endswith(".test.ts") or f.endswith(".test.tsx") or f.endswith(".spec.ts"):
                ts_tests = extract_ts_tests(p)
                for tname, thash in ts_tests.items():
                    tid = f"{rel_dir}.{tname}"
                    target_file = rel_file.replace(".test.ts", ".ts").replace(".test.tsx", ".tsx").replace(".spec.ts", ".ts")
                    tests_dict[tid] = {
                        "id": tid,
                        "package": rel_dir,
                        "test_file": rel_file,
                        "test_func": tname,
                        "test_hash": thash,
                        "target_file": target_file,
                        "target_func": "",
                        "code_hash": thash,
                        "duration_sec": 0.005,
                        "tier": "unit",
                        "last_status": "never_run",
                        "last_run_at": "",
                        "needs_run": True,
                    }
    return tests_dict


def build_test_inventory(repo_root: Path) -> dict[str, Any]:
    """Compiles the complete test inventory across all supported languages."""
    go_tests, go_count, _ = scan_go_tests(repo_root)
    py_ts_tests = scan_python_and_ts_tests(repo_root)

    combined_tests = {**go_tests, **py_ts_tests}
    packages = set(t["package"] for t in combined_tests.values())

    heavy_count = len([t for t in combined_tests.values() if t.get("tier") == "heavy"])
    unit_count = len([t for t in combined_tests.values() if t.get("tier") == "unit"])
    heavy_dur = round(sum(t.get("duration_sec", 0.0) for t in combined_tests.values() if t.get("tier") == "heavy"), 2)
    unit_dur = round(sum(t.get("duration_sec", 0.0) for t in combined_tests.values() if t.get("tier") == "unit"), 2)

    inventory = {
        "version": 1,
        "updated_at": datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "total_tests": len(combined_tests),
        "summary": {
            "total": len(combined_tests),
            "cached": 0,
            "dirty": len(combined_tests),
            "packages": len(packages),
            "heavy_tests": heavy_count,
            "unit_tests": unit_count,
            "estimated_heavy_sec": heavy_dur,
            "estimated_unit_sec": unit_dur,
        },
        "tests": combined_tests,
    }

    atomic_write_json(TEST_INVENTORY_PATH, inventory)
    return inventory


def record_recent_changes(changed_files: list[str]) -> dict[str, Any]:
    """Safely appends distinct modified relative paths to recent-file-changes.json under lock."""
    with file_lock(LOCK_FILE_PATH):
        existing_data: dict[str, Any] = {}
        if RECENT_CHANGES_PATH.is_file():
            try:
                existing_data = json.loads(RECENT_CHANGES_PATH.read_text(encoding="utf-8"))
            except Exception:
                existing_data = {}

        file_set = set(existing_data.get("files", []))
        for f in changed_files:
            rel = normalize_repo_rel(f)
            if rel:
                file_set.add(rel)

        inventory_tests: dict[str, Any] = {}
        if TEST_INVENTORY_PATH.is_file():
            try:
                inv = json.loads(TEST_INVENTORY_PATH.read_text(encoding="utf-8"))
                inventory_tests = inv.get("tests", {})
            except Exception:
                pass

        associated_tests: set[str] = set()
        for fpath in file_set:
            stem = Path(fpath).stem
            fdir = str(Path(fpath).parent).replace("\\", "/")
            for tid, tmeta in inventory_tests.items():
                target = tmeta.get("target_file", "")
                test_file = tmeta.get("test_file", "")
                pkg = tmeta.get("package", "")
                if (
                    target == fpath
                    or test_file == fpath
                    or (target and Path(target).stem == stem)
                    or (pkg and (pkg == fdir or fpath.startswith(pkg + "/")))
                ):
                    associated_tests.add(test_file)

        recent_payload = {
            "version": 1,
            "updated_at": datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            "files": sorted(list(file_set)),
            "associated_tests": sorted(list(associated_tests)),
        }

        atomic_write_json(RECENT_CHANGES_PATH, recent_payload)
        return recent_payload


def main():
    parser = argparse.ArgumentParser(description="Test inventory generator & atomic change recorder.")
    parser.add_argument("--record", nargs="+", help="Record modified file paths to recent-file-changes.json under lock.")
    parser.add_argument("--query-recent", action="store_true", help="Display recently modified files and associated tests.")
    parser.add_argument("--clear", action="store_true", help="Clear recent changes log.")
    args = parser.parse_args()

    if args.clear:
        with file_lock(LOCK_FILE_PATH):
            if RECENT_CHANGES_PATH.is_file():
                RECENT_CHANGES_PATH.unlink()
        print("Cleared recent changes log.")
        return

    if args.record:
        result = record_recent_changes(args.record)
        print(f"Recorded {len(args.record)} modified file(s). Total tracked: {len(result.get('files', []))}")
        print(f"Associated test files to run on release: {len(result.get('associated_tests', []))}")
        return

    if args.query_recent:
        if RECENT_CHANGES_PATH.is_file():
            print(RECENT_CHANGES_PATH.read_text(encoding="utf-8"))
        else:
            print("No recent changes recorded.")
        return

    print("Scanning codebase to generate test inventory...")
    inv = build_test_inventory(REPO_ROOT)
    print(f"Generated test inventory at {normalize_repo_rel(TEST_INVENTORY_PATH)}")
    print(f"Total Tests Indexed : {inv['total_tests']}")
    print(f"Total Packages      : {inv['summary']['packages']}")


if __name__ == "__main__":
    main()
