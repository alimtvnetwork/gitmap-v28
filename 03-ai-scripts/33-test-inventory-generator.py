#!/usr/bin/env python3
"""
33-test-inventory-generator.py
==============================
Generates and maintains the centralized test inventory manifest at `.ai-memory/test-inventory.json`,
test heatmap at `.ai-memory/test-heatmap.json` (and `.ai-memory/cicd/test-heatmap.json`),
and provides safe, atomic file change recording into `.ai-memory/temp/recent-file-changes.json`
with file locking to ensure concurrency safety across multi-agent turns.

Usage:
  # Scan and generate / update test inventory and heatmap:
  python 03-ai-scripts/33-test-inventory-generator.py

  # Display heatmap report:
  python 03-ai-scripts/33-test-inventory-generator.py --heatmap --top 25

  # Record modified files safely under lock:
  python 03-ai-scripts/33-test-inventory-generator.py --record "cli/cmd/root.go"

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
import subprocess
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

REPO_ROOT = Path(__file__).resolve().parent.parent
LOVABLE_DIR = REPO_ROOT / ".ai-memory"
TEMP_DIR = LOVABLE_DIR / "temp"
CICD_DIR = LOVABLE_DIR / "cicd"
TEST_INVENTORY_PATH = LOVABLE_DIR / "test-inventory.json"
TEST_HEATMAP_PATH = LOVABLE_DIR / "test-heatmap.json"
CICD_HEATMAP_PATH = CICD_DIR / "test-heatmap.json"
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


def try_open_lock(lock_path: Path, pid_bytes: bytes) -> bool:
    """Tries creating and writing lock file."""
    try:
        fd = os.open(str(lock_path), os.O_CREAT | os.O_EXCL | os.O_RDWR)
        os.write(fd, pid_bytes)
        os.close(fd)
        return True
    except FileExistsError:
        return False


def acquire_lock_file(lock_path: Path, timeout_sec: float) -> bool:
    """Spins until acquiring the lock file or timing out."""
    start_time = time.time()
    pid_bytes = f"{os.getpid()}\n".encode("utf-8")
    while not try_open_lock(lock_path, pid_bytes):
        if time.time() - start_time > timeout_sec:
            try:
                lock_path.unlink()
            except OSError:
                pass
        time.sleep(0.05)
    return True


def release_lock_file(lock_path: Path) -> None:
    """Releases the lock file."""
    try:
        lock_path.unlink()
    except OSError:
        pass


@contextlib.contextmanager
def file_lock(lock_path: Path, timeout_sec: float = 10.0):
    """Acquires a cross-platform cooperative lock file with timeout."""
    lock_path.parent.mkdir(parents=True, exist_ok=True)
    acquire_lock_file(lock_path, timeout_sec)
    try:
        yield
    finally:
        release_lock_file(lock_path)


def remove_temp_file(temp_file: Path) -> None:
    """Removes a temporary file safely."""
    if temp_file.exists():
        try:
            temp_file.unlink()
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
        remove_temp_file(temp_file)
        raise


def flush_go_func(funcs: dict[str, str], cur_func: str | None, cur_lines: list[str]) -> None:
    """Hashes and stores the previous Go function chunk."""
    if cur_func:
        chunk = "\n".join(cur_lines).encode("utf-8")
        funcs[cur_func] = hashlib.sha256(chunk).hexdigest()[:16]


def parse_go_declarations(lines: list[str]) -> dict[str, str]:
    """Parses Go function declarations from lines."""
    funcs: dict[str, str] = {}
    cur_func, cur_lines = None, []
    for line in lines:
        m = FUNC_START_RE.match(line)
        if m:
            flush_go_func(funcs, cur_func, cur_lines)
            cur_func, cur_lines = m.group(1), [line]
        elif cur_func:
            cur_lines.append(line)
    flush_go_func(funcs, cur_func, cur_lines)
    return funcs


def extract_go_declarations(filepath: Path) -> dict[str, str]:
    """Extracts function/method declarations and their line spans."""
    if not filepath.is_file():
        return {}
    try:
        lines = filepath.read_text(encoding="utf-8", errors="replace").splitlines()
        return parse_go_declarations(lines)
    except Exception:
        return {}


def flush_go_test(tests: dict[str, tuple[str, str]], cur_test: str | None, cur_lines: list[str]) -> None:
    """Hashes and stores previous Go test body."""
    if cur_test:
        body = "\n".join(cur_lines)
        chunk = body.encode("utf-8")
        tests[cur_test] = (hashlib.sha256(chunk).hexdigest()[:16], body)


def parse_go_test_lines(lines: list[str]) -> dict[str, tuple[str, str]]:
    """Parses Go test function definitions from lines."""
    tests: dict[str, tuple[str, str]] = {}
    cur_test, cur_lines = None, []
    for line in lines:
        m = TEST_START_RE.match(line)
        if m:
            flush_go_test(tests, cur_test, cur_lines)
            cur_test, cur_lines = m.group(1), [line]
        elif cur_test:
            cur_lines.append(line)
    flush_go_test(tests, cur_test, cur_lines)
    return tests


def extract_go_tests(filepath: Path) -> dict[str, tuple[str, str]]:
    """Extracts Go test functions (Test*) from a test file: tname -> (hash, body)."""
    if not filepath.is_file():
        return {}
    try:
        lines = filepath.read_text(encoding="utf-8", errors="replace").splitlines()
        return parse_go_test_lines(lines)
    except Exception:
        return {}


def flush_py_test(tests: dict[str, str], cur_test: str | None, cur_lines: list[str]) -> None:
    """Hashes and stores previous Python test chunk."""
    if cur_test:
        chunk = "\n".join(cur_lines).encode("utf-8")
        tests[cur_test] = hashlib.sha256(chunk).hexdigest()[:16]


def parse_python_test_lines(lines: list[str]) -> dict[str, str]:
    """Parses Python test functions from lines."""
    tests: dict[str, str] = {}
    cur_test, cur_lines = None, []
    for line in lines:
        m = PY_TEST_START_RE.match(line.strip())
        if m:
            flush_py_test(tests, cur_test, cur_lines)
            cur_test, cur_lines = m.group(1), [line]
        elif cur_test:
            cur_lines.append(line)
    flush_py_test(tests, cur_test, cur_lines)
    return tests


def extract_python_tests(filepath: Path) -> dict[str, str]:
    """Extracts Python test functions (test_*) from a test file."""
    if not filepath.is_file():
        return {}
    try:
        lines = filepath.read_text(encoding="utf-8", errors="replace").splitlines()
        return parse_python_test_lines(lines)
    except Exception:
        return {}


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


def check_duration_keywords(clean_code: str, test_name: str) -> float:
    """Adjusts base duration based on keywords in test code."""
    duration = 0.005
    if "exec.Command" in clean_code:
        duration += 3.0
    if "time.Sleep" in clean_code:
        duration += 1.5
    if "net.Listen" in clean_code or "http.Get" in clean_code or "http.Post" in clean_code:
        duration += 0.5
    if "git" in test_name.lower() and ("subprocess" in clean_code.lower() or "exec" in clean_code.lower()):
        duration += 2.0
    return duration


def estimate_test_duration(
    filepath: Path, test_name: str, content: str, slow_threshold: float = 4.0
) -> tuple[float, str, bool]:
    """Estimates test duration in seconds and categorizes tier (slow vs fast)."""
    rel = str(filepath).replace("\\", "/")
    if "tests/heavy_test" in rel:
        return 5.0, "slow", True
    code_lines = [line for line in content.splitlines() if not line.strip().startswith("//")]
    duration = check_duration_keywords("\n".join(code_lines), test_name)
    is_slow = bool(duration >= slow_threshold)
    tier = "slow" if is_slow else "fast"
    return round(duration, 3), tier, is_slow


def resolve_heavy_test_target(repo_root: Path, rel_test_file: str, tf: Path) -> tuple[str, str]:
    """Resolves heavy test target file and hash."""
    stem = tf.name.replace("_e2e_test.go", "").replace("_test.go", "")
    prefix = stem.split("_")[0]
    cand_dir = repo_root / "cli" / prefix
    if cand_dir.is_dir():
        go_files = sorted([f for f in cand_dir.glob("*.go") if not f.name.endswith("_test.go")])
        if go_files:
            return normalize_repo_rel(go_files[0]), compute_file_hash(go_files[0])
    return "", ""


def resolve_exact_or_suffix_target(tf: Path, pkg_dir: Path) -> tuple[str, str]:
    """Resolves standard Go source targets by exact name or unit/e2e suffix."""
    exact_path = pkg_dir / tf.name.replace("_test.go", ".go")
    if exact_path.is_file():
        return normalize_repo_rel(exact_path), compute_file_hash(exact_path)
    suffixes = ("_unit_test.go", "_e2e_test.go", "_integration_test.go", "_helpers_test.go")
    for suffix in suffixes:
        if tf.name.endswith(suffix):
            cand_p = pkg_dir / tf.name.replace(suffix, ".go")
            if cand_p.is_file():
                return normalize_repo_rel(cand_p), compute_file_hash(cand_p)
    return "", ""


def resolve_target_file(tf: Path, pkg_dir: Path, repo_root: Path, rel_test_file: str) -> tuple[str, str]:
    """Intelligently resolves relative target source file and hash for a test file."""
    tgt, h = resolve_exact_or_suffix_target(tf, pkg_dir)
    if tgt:
        return tgt, h
    if "tests/heavy_test" in rel_test_file:
        tgt_heavy, h_heavy = resolve_heavy_test_target(repo_root, rel_test_file, tf)
        if tgt_heavy:
            return tgt_heavy, h_heavy
    pkg_go_files = sorted([f for f in pkg_dir.glob("*.go") if not f.name.endswith("_test.go")])
    if pkg_go_files:
        return normalize_repo_rel(pkg_go_files[0]), compute_file_hash(pkg_go_files[0])
    return rel_test_file, compute_file_hash(tf)


def get_decay_weight(commit_idx: int) -> float:
    """Calculates recency decay weight based on commit position."""
    if commit_idx <= 10:
        return 1.0
    if commit_idx <= 25:
        return 0.6
    return 0.3


def is_git_log_path(line: str) -> bool:
    """Verifies that a line is a file path rather than git metadata."""
    has_content = bool(line and not line.isspace())
    is_header = bool(line.startswith(("commit ", "Author:", "Date:", "Merge:")))
    return has_content and not is_header


def parse_git_log_output(lines: list[str]) -> dict[str, float]:
    """Parses git log lines to aggregate file churn with recency decay."""
    churn_map: dict[str, float] = {}
    commit_idx = 0
    for line in lines:
        sline = line.strip()
        if sline.startswith("commit "):
            commit_idx += 1
        elif is_git_log_path(sline):
            w = get_decay_weight(commit_idx)
            frel = normalize_repo_rel(sline)
            churn_map[frel] = round(churn_map.get(frel, 0.0) + w, 2)
    return churn_map


def extract_git_commit_churn(repo_root: Path, commits: int = 50) -> dict[str, float]:
    """Tally commit churn per file with recency decay over the last N commits."""
    cmd = ["git", "log", f"-n{commits}", "--name-only", "--format=commit %H"]
    try:
        proc = subprocess.run(cmd, cwd=str(repo_root), capture_output=True, text=True, check=False)
        return parse_git_log_output(proc.stdout.splitlines())
    except Exception:
        return {}


def extract_test_id_from_failure_line(line: str) -> str:
    """Extracts a test identifier from a log error line."""
    m = re.search(r"Test\s+([A-Za-z0-9_./\\-]+(?:\.[A-Za-z0-9_]+)?)\s+failed", line)
    if m:
        return m.group(1).replace("\\", "/")
    m2 = re.search(r"\[([A-Za-z0-9_./\\-]+\.Test[A-Za-z0-9_]+)\]", line)
    if m2:
        return m2.group(1).replace("\\", "/")
    return ""


def extract_temp_failures(temp_fail_dir: Path, failure_map: dict[str, int]) -> None:
    """Scans .ai-memory/temp/failures/ for historical failure records."""
    if not temp_fail_dir.is_dir():
        return
    for log_file in temp_fail_dir.glob("*.log"):
        tid = ""
        try:
            first_line = log_file.read_text(encoding="utf-8", errors="replace").splitlines()[0]
            tid = extract_test_id_from_failure_line(first_line)
        except Exception:
            pass
        cand_tid = tid if tid else log_file.stem
        failure_map[cand_tid] = failure_map.get(cand_tid, 0) + 1
        tfunc = cand_tid.split(".")[-1]
        failure_map[tfunc] = failure_map.get(tfunc, 0) + 1


def parse_run_error_text(err_text: str, run_fails: set[str]) -> None:
    """Extracts test IDs from error message text."""
    for line in err_text.splitlines():
        tid = extract_test_id_from_failure_line(line)
        if tid:
            run_fails.add(tid)


def parse_run_errors_json(errors_json_path: Path) -> set[str]:
    """Extracts distinct test failure IDs from a single run errors.json."""
    run_fails: set[str] = set()
    if not errors_json_path.is_file():
        return run_fails
    try:
        data = json.loads(errors_json_path.read_text(encoding="utf-8"))
        for entry in data:
            parse_run_error_text(entry.get("error", ""), run_fails)
    except Exception:
        pass
    return run_fails


def extract_cicd_run_failures(runs_dir: Path, failure_map: dict[str, int]) -> None:
    """Scans .ai-memory/cicd/runs/ subdirectories for historical run failures."""
    if not runs_dir.is_dir():
        return
    for run_folder in runs_dir.iterdir():
        if not run_folder.is_dir():
            continue
        run_fails = parse_run_errors_json(run_folder / "errors.json")
        for tid in run_fails:
            failure_map[tid] = failure_map.get(tid, 0) + 1
            tfunc = tid.split(".")[-1]
            failure_map[tfunc] = failure_map.get(tfunc, 0) + 1


def extract_historical_failures(repo_root: Path) -> dict[str, int]:
    """Queries .ai-memory/temp/failures/ and .ai-memory/cicd/runs/ for failure counts."""
    failure_map: dict[str, int] = {}
    rpath = Path(repo_root)
    extract_temp_failures(rpath / ".ai-memory" / "temp" / "failures", failure_map)
    extract_cicd_run_failures(rpath / ".ai-memory" / "cicd" / "runs", failure_map)
    return failure_map


def extract_git_status_files(repo_root: Path) -> set[str]:
    """Extracts currently modified, added, or untracked files from git status."""
    files: set[str] = set()
    try:
        proc = subprocess.run(["git", "status", "--porcelain"], cwd=str(repo_root), capture_output=True, text=True, check=False)
        for line in proc.stdout.splitlines():
            if len(line) >= 4:
                files.add(normalize_repo_rel(line[3:].strip()))
    except Exception:
        pass
    return files


def read_recent_changed_files(repo_root: Path) -> set[str]:
    """Reads recorded file changes from recent-file-changes.json."""
    recent_path = Path(repo_root) / ".ai-memory" / "temp" / "recent-file-changes.json"
    if not recent_path.is_file():
        return set()
    try:
        data = json.loads(recent_path.read_text(encoding="utf-8"))
        return set(normalize_repo_rel(f) for f in data.get("files", []))
    except Exception:
        return set()


def extract_dirty_files(repo_root: Path) -> set[str]:
    """Combines working-tree git dirty files and recent recorded changes."""
    git_files = extract_git_status_files(repo_root)
    recent_files = read_recent_changed_files(repo_root)
    return git_files.union(recent_files)


def compute_failure_score(test_meta: dict[str, Any], failure_map: dict[str, int]) -> tuple[float, int, bool]:
    """Calculates failure score up to 40 points."""
    tid = test_meta.get("id", "")
    tfunc = test_meta.get("test_func", "")
    fail_count = failure_map.get(tid, 0) or failure_map.get(tfunc, 0)
    is_failing = bool(test_meta.get("last_status") == "failed" or fail_count > 0 and test_meta.get("last_status") != "passed")
    if test_meta.get("last_status") == "failed":
        return 40.0, fail_count, True
    return min(40.0, float(fail_count * 10.0)), fail_count, is_failing


def compute_churn_score(test_meta: dict[str, Any], churn_map: dict[str, float]) -> tuple[float, float]:
    """Calculates churn score up to 35 points based on target and test churn."""
    target = test_meta.get("target_file", "")
    tfile = test_meta.get("test_file", "")
    pkg = test_meta.get("package", "")
    churn = churn_map.get(target, 0.0) + churn_map.get(tfile, 0.0)
    if churn == 0.0 and pkg:
        pkg_churn = sum(v for k, v in churn_map.items() if k.startswith(pkg + "/"))
        churn = min(3.0, round(pkg_churn * 0.3, 2))
    score = min(35.0, round(churn * 7.0, 1))
    return score, round(churn, 2)


def compute_dirty_score(test_meta: dict[str, Any], dirty_files: set[str]) -> tuple[float, bool]:
    """Calculates dirty file boost up to 15 points."""
    target = test_meta.get("target_file", "")
    tfile = test_meta.get("test_file", "")
    pkg = test_meta.get("package", "")
    is_direct_dirty = bool(target in dirty_files or tfile in dirty_files)
    if is_direct_dirty:
        return 15.0, True
    has_pkg_dirty = bool(pkg and any(df.startswith(pkg + "/") for df in dirty_files))
    if has_pkg_dirty:
        return 10.0, False
    return 0.0, False


def compute_speed_weight(duration_sec: float) -> float:
    """Calculates speed weight up to 10 points for faster tests."""
    if duration_sec <= 0.05:
        return 10.0
    if duration_sec <= 0.5:
        return 8.0
    if duration_sec <= 2.0:
        return 5.0
    if duration_sec < 4.0:
        return 3.0
    return 0.0


def assign_heat_tier(score: float, is_failing: bool, is_dirty: bool) -> str:
    """Categorizes heat score into hot, warm, or cold tier."""
    if score >= 60.0 or is_failing or is_dirty:
        return "hot"
    if score >= 25.0:
        return "warm"
    return "cold"


def calculate_composite_score(s_fail: float, s_churn: float, s_dirty: float, dur: float) -> float:
    """Combines sub-scores into final bounded heat score."""
    w_speed = compute_speed_weight(dur)
    raw_score = s_fail + s_churn + s_dirty + w_speed
    return max(0.0, min(100.0, round(raw_score, 1)))


def compute_test_heat_score(
    test_meta: dict[str, Any], churn_map: dict[str, float], failure_map: dict[str, int], dirty_files: set[str],
) -> tuple[float, str, float, int]:
    """Computes composite Heat Score H(T) in [0, 100] and updates test metadata."""
    s_fail, fcount, is_failing = compute_failure_score(test_meta, failure_map)
    s_churn, churn_val = compute_churn_score(test_meta, churn_map)
    s_dirty, is_dirty = compute_dirty_score(test_meta, dirty_files)
    score = calculate_composite_score(s_fail, s_churn, s_dirty, test_meta.get("duration_sec", 0.0))
    tier = assign_heat_tier(score, is_failing, is_dirty)
    test_meta.update({"heat_score": score, "heat_tier": tier, "churn_count": churn_val, "failure_count": fcount})
    return score, tier, churn_val, fcount


def create_go_test_entry(
    tid: str, rel_pkg: str, rel_test_file: str, tname: str, thash: str,
    rel_target: str, target_func: str, code_hash: str, dur: float, tier: str, is_slow: bool,
) -> dict[str, Any]:
    """Creates a base Go test metadata entry dictionary."""
    return {
        "id": tid, "package": rel_pkg, "test_file": rel_test_file,
        "test_func": tname, "test_hash": thash, "target_file": rel_target,
        "target_func": target_func, "code_hash": code_hash, "duration_sec": dur,
        "tier": tier, "is_slow": is_slow, "last_status": "never_run",
        "last_run_at": "", "needs_run": True,
    }


def create_scored_go_test(
    pkg: str, rtest: str, tname: str, thash: str, tbody: str, rtgt: str, chash: str, th: float,
    maps: tuple[dict[str, float], dict[str, int], set[str]],
) -> tuple[str, dict[str, Any]]:
    """Creates and heat-scores a single Go test entry."""
    tid = f"{pkg}.{tname}"
    tfunc = tname[4:].split("_")[0] if tname.startswith("Test") else ""
    dur, tier, is_slow = estimate_test_duration(Path(rtest), tname, tbody, th)
    entry = create_go_test_entry(tid, pkg, rtest, tname, thash, rtgt, tfunc, chash, dur, tier, is_slow)
    compute_test_heat_score(entry, maps[0], maps[1], maps[2])
    return tid, entry


def index_go_test_file(
    tf: Path, pkg_dir: Path, repo_root: Path, th: float,
    maps: tuple[dict[str, float], dict[str, int], set[str]],
) -> dict[str, Any]:
    """Indexes all test functions in a single Go test file."""
    rtest, rpkg = normalize_repo_rel(tf), normalize_repo_rel(pkg_dir)
    rtgt, chash = resolve_target_file(tf, pkg_dir, repo_root, rtest)
    file_tests: dict[str, Any] = {}
    for tname, (thash, tbody) in extract_go_tests(tf).items():
        tid, entry = create_scored_go_test(rpkg, rtest, tname, thash, tbody, rtgt, chash, th, maps)
        file_tests[tid] = entry
    return file_tests


def is_skippable_dir(rel_path: str) -> bool:
    """Identifies directories to ignore during indexing."""
    return any(p in rel_path for p in (".git", "node_modules", ".ai-memory", "dist"))


def scan_go_tests(
    repo_root: Path, slow_threshold: float = 4.0, force_run_all: bool = False,
    churn_map: dict[str, float] | None = None, failure_map: dict[str, int] | None = None,
    dirty_files: set[str] | None = None,
) -> tuple[dict[str, Any], int, int]:
    """Scans and indexes all Go test files with heat scores and tiers."""
    tests: dict[str, Any] = {}
    maps = (churn_map or {}, failure_map or {}, dirty_files or set())
    for tf in repo_root.rglob("*_test.go"):
        if not is_skippable_dir(normalize_repo_rel(tf)):
            tests.update(index_go_test_file(tf, tf.parent, repo_root, slow_threshold, maps))
    return tests, len(tests), len(set(t["package"] for t in tests.values()))


def create_scored_py_test(
    rdir: str, rfile: str, tname: str, thash: str, maps: tuple[dict[str, float], dict[str, int], set[str]],
) -> dict[str, Any]:
    """Creates and heat-scores a Python test entry."""
    tgt = rfile.replace("test_", "").replace("_test.py", ".py")
    meta = {
        "id": f"{rdir}.{tname}", "package": rdir, "test_file": rfile, "test_func": tname,
        "test_hash": thash, "target_file": tgt, "target_func": "", "code_hash": thash,
        "duration_sec": 0.05, "tier": "fast", "is_slow": False, "last_status": "never_run",
        "last_run_at": "", "needs_run": True,
    }
    compute_test_heat_score(meta, maps[0], maps[1], maps[2])
    return meta


def index_py_tests_in_file(
    p: Path, rel_dir: str, rel_file: str, maps: tuple[dict[str, float], dict[str, int], set[str]],
) -> dict[str, Any]:
    """Extracts and scores Python tests in a single file."""
    tests: dict[str, Any] = {}
    for tname, thash in extract_python_tests(p).items():
        meta = create_scored_py_test(rel_dir, rel_file, tname, thash, maps)
        tests[meta["id"]] = meta
    return tests


def create_scored_ts_test(
    rdir: str, rfile: str, tname: str, thash: str, maps: tuple[dict[str, float], dict[str, int], set[str]],
) -> dict[str, Any]:
    """Creates and heat-scores a TypeScript test entry."""
    tgt = rfile.replace(".test.ts", ".ts").replace(".test.tsx", ".tsx").replace(".spec.ts", ".ts")
    meta = {
        "id": f"{rdir}.{tname}", "package": rdir, "test_file": rfile, "test_func": tname,
        "test_hash": thash, "target_file": tgt, "target_func": "", "code_hash": thash,
        "duration_sec": 0.005, "tier": "fast", "is_slow": False, "last_status": "never_run",
        "last_run_at": "", "needs_run": True,
    }
    compute_test_heat_score(meta, maps[0], maps[1], maps[2])
    return meta


def index_ts_tests_in_file(
    p: Path, rel_dir: str, rel_file: str, maps: tuple[dict[str, float], dict[str, int], set[str]],
) -> dict[str, Any]:
    """Extracts and scores TypeScript tests in a single file."""
    tests: dict[str, Any] = {}
    for tname, thash in extract_ts_tests(p).items():
        meta = create_scored_ts_test(rel_dir, rel_file, tname, thash, maps)
        tests[meta["id"]] = meta
    return tests


def index_walk_files(
    root: str, files: list[str], maps: tuple[dict[str, float], dict[str, int], set[str]],
) -> dict[str, Any]:
    """Processes Python and TS test files found in a directory."""
    tests: dict[str, Any] = {}
    rdir = normalize_repo_rel(root)
    for f in files:
        p, rfile = Path(root) / f, normalize_repo_rel(Path(root) / f)
        if (f.startswith("test_") or f.endswith("_test.py")) and f.endswith(".py"):
            tests.update(index_py_tests_in_file(p, rdir, rfile, maps))
        elif f.endswith((".test.ts", ".test.tsx", ".spec.ts")):
            tests.update(index_ts_tests_in_file(p, rdir, rfile, maps))
    return tests


def scan_python_and_ts_tests(
    repo_root: Path, slow_threshold: float = 4.0,
    churn_map: dict[str, float] | None = None, failure_map: dict[str, int] | None = None,
    dirty_files: set[str] | None = None,
) -> dict[str, Any]:
    """Scans and indexes Python and TypeScript test files with heat scores."""
    tests_dict: dict[str, Any] = {}
    maps = (churn_map or {}, failure_map or {}, dirty_files or set())
    for root, _, files in os.walk(repo_root):
        if not is_skippable_dir(normalize_repo_rel(root)):
            tests_dict.update(index_walk_files(root, files, maps))
    return tests_dict


def count_heat_tiers(tests: dict[str, Any]) -> tuple[int, int, int]:
    """Counts tests in each heat tier."""
    hot = sum(1 for t in tests.values() if t.get("heat_tier") == "hot")
    warm = sum(1 for t in tests.values() if t.get("heat_tier") == "warm")
    cold = sum(1 for t in tests.values() if t.get("heat_tier") == "cold")
    return hot, warm, cold


def calculate_inventory_summary(tests: dict[str, Any], pkgs: int, th: float) -> dict[str, Any]:
    """Aggregates inventory counts, durations, and heat tier distribution."""
    slow = [t for t in tests.values() if t.get("is_slow", False) or t.get("tier") in ("slow", "heavy")]
    fast = [t for t in tests.values() if not (t.get("is_slow", False) or t.get("tier") in ("slow", "heavy"))]
    sdur = round(sum(t.get("duration_sec", 0.0) for t in slow), 2)
    fdur = round(sum(t.get("duration_sec", 0.0) for t in fast), 2)
    hot, warm, cold = count_heat_tiers(tests)
    return {
        "total": len(tests), "cached": 0, "dirty": len(tests), "packages": pkgs,
        "slow_tests": len(slow), "fast_tests": len(fast), "heavy_tests": len(slow), "unit_tests": len(fast),
        "hot_tests": hot, "warm_tests": warm, "cold_tests": cold, "slow_threshold_sec": th,
        "estimated_slow_sec": sdur, "estimated_fast_sec": fdur, "estimated_heavy_sec": sdur, "estimated_unit_sec": fdur,
    }


def persist_inventory_manifests(inventory: dict[str, Any]) -> None:
    """Writes inventory and heatmap JSON files atomically."""
    atomic_write_json(TEST_INVENTORY_PATH, inventory)
    atomic_write_json(TEST_HEATMAP_PATH, inventory)
    atomic_write_json(CICD_HEATMAP_PATH, inventory)


def build_test_inventory(
    repo_root: Path, slow_threshold: float = 4.0, force_run_all: bool = False, commits: int = 50,
) -> dict[str, Any]:
    """Compiles complete test inventory and persists test-inventory.json and heatmap manifests."""
    r = Path(repo_root)
    c_map, f_map, d_files = extract_git_commit_churn(r, commits), extract_historical_failures(r), extract_dirty_files(r)
    g_tests, _, _ = scan_go_tests(r, slow_threshold, force_run_all, c_map, f_map, d_files)
    p_tests = scan_python_and_ts_tests(r, slow_threshold, c_map, f_map, d_files)
    all_t = {**g_tests, **p_tests}
    sorted_t = dict(sorted(all_t.items(), key=lambda it: it[1].get("heat_score", 0.0), reverse=True))
    summary = calculate_inventory_summary(sorted_t, len(set(t["package"] for t in sorted_t.values())), slow_threshold)
    now_str = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    inv = {"version": 1, "updated_at": now_str, "total_tests": len(sorted_t), "summary": summary, "tests": sorted_t}
    persist_inventory_manifests(inv)
    return inv


def find_associated_tests(file_set: set[str], inv_tests: dict[str, Any]) -> set[str]:
    """Identifies test files affected by a set of modified paths."""
    associated_tests: set[str] = set()
    for fpath in file_set:
        stem = Path(fpath).stem
        fdir = str(Path(fpath).parent).replace("\\", "/")
        for tid, tmeta in inv_tests.items():
            tgt = tmeta.get("target_file", "")
            tfile = tmeta.get("test_file", "")
            pkg = tmeta.get("package", "")
            is_hit = bool(tgt == fpath or tfile == fpath or (tgt and Path(tgt).stem == stem) or (pkg and (pkg == fdir or fpath.startswith(pkg + "/"))))
            if is_hit:
                associated_tests.add(tfile)
    return associated_tests


def load_existing_change_files() -> set[str]:
    """Loads existing tracked changes from disk."""
    if not RECENT_CHANGES_PATH.is_file():
        return set()
    try:
        return set(json.loads(RECENT_CHANGES_PATH.read_text(encoding="utf-8")).get("files", []))
    except Exception:
        return set()


def load_inventory_test_map() -> dict[str, Any]:
    """Loads existing test inventory dictionary."""
    if not TEST_INVENTORY_PATH.is_file():
        return {}
    try:
        return json.loads(TEST_INVENTORY_PATH.read_text(encoding="utf-8")).get("tests", {})
    except Exception:
        return {}


def record_recent_changes(changed_files: list[str]) -> dict[str, Any]:
    """Safely appends distinct modified relative paths to recent-file-changes.json under lock."""
    with file_lock(LOCK_FILE_PATH):
        file_set = load_existing_change_files()
        file_set.update(normalize_repo_rel(f) for f in changed_files if normalize_repo_rel(f))
        associated = find_associated_tests(file_set, load_inventory_test_map())
        now_str = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        payload = {"version": 1, "updated_at": now_str, "files": sorted(list(file_set)), "associated_tests": sorted(list(associated))}
        atomic_write_json(RECENT_CHANGES_PATH, payload)
        return payload


def get_inventory_updated_dt(inv: dict[str, Any], path: Path) -> datetime.datetime:
    """Extracts timestamp from manifest or fallback file mtime."""
    raw_updated = inv.get("updated_at", "")
    if raw_updated:
        return datetime.datetime.fromisoformat(raw_updated.replace("Z", "+00:00"))
    mtime = os.path.getmtime(path)
    return datetime.datetime.fromtimestamp(mtime, datetime.timezone.utc)


def audit_inventory_freshness(inv: dict[str, Any], path: Path, max_age_days: float) -> tuple[bool, dict[str, Any]]:
    """Evaluates age and profiling data for an inventory manifest."""
    updated_dt = get_inventory_updated_dt(inv, path)
    now_dt = datetime.datetime.now(datetime.timezone.utc)
    age_days = round((now_dt - updated_dt).total_seconds() / 86400.0, 2)
    summary = inv.get("summary", {})
    tests = inv.get("tests", {})
    has_profiled = any(t.get("duration_sec", 0.0) > 0.0 for t in tests.values())
    is_fresh = bool(age_days <= max_age_days and has_profiled)
    return is_fresh, {
        "exists": True, "is_fresh": is_fresh, "age_days": age_days, "max_age_days": max_age_days,
        "updated_at": inv.get("updated_at", updated_dt.isoformat()), "total_tests": len(tests),
        "slow_tests": summary.get("slow_tests", 0), "fast_tests": summary.get("fast_tests", 0),
        "has_profiled": has_profiled,
    }


def check_inventory_age(max_age_days: float = 5.0) -> tuple[bool, dict[str, Any]]:
    """Audits the freshness of .ai-memory/test-inventory.json against max_age_days."""
    if not TEST_INVENTORY_PATH.is_file():
        return False, {"exists": False, "is_fresh": False, "age_days": None, "message": "Missing manifest"}
    try:
        inv = json.loads(TEST_INVENTORY_PATH.read_text(encoding="utf-8"))
        return audit_inventory_freshness(inv, TEST_INVENTORY_PATH, max_age_days)
    except Exception as err:
        return False, {"exists": True, "is_fresh": False, "age_days": None, "error": str(err)}


def print_age_report_details(audit: dict[str, Any]) -> None:
    """Prints detailed status fields for inventory age audit."""
    is_fresh = bool(audit.get("is_fresh", False))
    status_str = "✅ FRESH (<= threshold)" if is_fresh else "⚠️  STALE (> threshold or unprofiled)"
    has_prof = bool(audit.get("has_profiled", False))
    print(f"📁 Manifest Path : {normalize_repo_rel(TEST_INVENTORY_PATH)}")
    print(f"🕒 Last Updated  : {audit.get('updated_at')} ({audit.get('age_days')} days ago)")
    print(f"🏷️  Status        : {status_str}")
    print(f"📊 Tests Catalog : {audit.get('total_tests')} (Slow: {audit.get('slow_tests')}, Fast: {audit.get('fast_tests')})")
    print(f"🎯 Profiled Data : {'Yes' if has_prof else 'No'}")
    print("=" * 80)


def display_age_report(audit: dict[str, Any]) -> None:
    """Formats and prints test inventory age audit report to console."""
    print("=" * 80)
    print("📦 TEST INVENTORY FRESHNESS AUDIT")
    print("=" * 80)
    if not audit.get("exists", False):
        print(f"❌ Manifest Missing: {TEST_INVENTORY_PATH}")
        print("=" * 80)
        return
    print_age_report_details(audit)


def format_heatmap_row(rank: int, test_meta: dict[str, Any]) -> str:
    """Formats a single test row for terminal heatmap display."""
    score = test_meta.get("heat_score", 0.0)
    tier = str(test_meta.get("heat_tier", "cold")).upper()
    dur = f"{test_meta.get('duration_sec', 0.0):.2f}s"
    churn = f"{test_meta.get('churn_count', 0.0):.1f}"
    fail = str(test_meta.get("failure_count", 0))
    tid = test_meta.get("id", "")
    if len(tid) > 42:
        tid = tid[:39] + "..."
    return f"{rank:>3}  {score:>5.1f}  {tier:<4}  {dur:>6}  {churn:>5}  {fail:>4}  {tid}"


def print_heatmap_table(tests: list[dict[str, Any]], top_n: int) -> None:
    """Prints tabular heatmap tests up to top_n limit."""
    print("-" * 80)
    print(f"{'Rank':>3}  {'Score':>5}  {'Tier':<4}  {'Dur':>6}  {'Churn':>5}  {'Fail':>4}  {'Test Target':<42}")
    print("-" * 80)
    for idx, t in enumerate(tests[:top_n], start=1):
        print(format_heatmap_row(idx, t))
    print("=" * 80)


def display_heatmap_report(inventory: dict[str, Any], top_n: int = 25) -> None:
    """Prints formatted heatmap summary and table to stdout."""
    tests = list(inventory.get("tests", {}).values())
    tests.sort(key=lambda t: t.get("heat_score", 0.0), reverse=True)
    summary = inventory.get("summary", {})
    limit = min(top_n, len(tests))
    print("=" * 80)
    print(f"🔥 TEST HEATMAP REPORT (Top {limit})")
    print("=" * 80)
    s_total = summary.get("total", len(tests))
    s_hot = summary.get("hot_tests", 0)
    s_warm = summary.get("warm_tests", 0)
    s_cold = summary.get("cold_tests", 0)
    print(f"Summary: Total: {s_total} | Hot: {s_hot} | Warm: {s_warm} | Cold: {s_cold}")
    print_heatmap_table(tests, top_n)


def add_inventory_cli_args(p: argparse.ArgumentParser, default_th: float) -> None:
    """Registers command line options for inventory generator."""
    p.add_argument("--record", nargs="+", help="Record modified file paths.")
    p.add_argument("--query-recent", action="store_true", help="Display recent changes.")
    p.add_argument("--clear", action="store_true", help="Clear recent changes log.")
    p.add_argument("--slow-threshold", type=float, default=default_th, help="Slow threshold.")
    p.add_argument("--force-run-all", action="store_true", help="Force run all tests.")
    p.add_argument("--check-age", "--age", action="store_true", help="Check manifest age.")
    p.add_argument("--max-age-days", type=float, default=5.0, help="Max age in days.")
    p.add_argument("--json", action="store_true", help="Output JSON format.")
    p.add_argument("--heatmap", action="store_true", help="Display heatmap report.")
    p.add_argument("--top", type=int, default=25, help="Top N tests in report.")
    p.add_argument("--commits", type=int, default=50, help="Commits to scan for churn.")


def build_arg_parser() -> argparse.ArgumentParser:
    """Constructs command line argument parser."""
    default_th = float(os.environ.get("GITMAP_SLOW_TEST_THRESHOLD", "4.0"))
    p = argparse.ArgumentParser(description="Test inventory generator & atomic change recorder.")
    add_inventory_cli_args(p, default_th)
    return p


def handle_check_age(args: argparse.Namespace) -> None:
    """Handles --check-age CLI option."""
    is_fresh, audit = check_inventory_age(max_age_days=args.max_age_days)
    if args.json:
        print(json.dumps(audit, indent=2))
    else:
        display_age_report(audit)
    sys.exit(0 if is_fresh else 1)


def handle_record_changes(files: list[str]) -> None:
    """Handles --record CLI option."""
    result = record_recent_changes(files)
    print(f"Recorded {len(files)} modified file(s). Total tracked: {len(result.get('files', []))}")
    print(f"Associated test files to run on release: {len(result.get('associated_tests', []))}")


def handle_query_recent() -> None:
    """Handles --query-recent CLI option."""
    if RECENT_CHANGES_PATH.is_file():
        print(RECENT_CHANGES_PATH.read_text(encoding="utf-8"))
    else:
        print("No recent changes recorded.")


def handle_clear_recent_log() -> None:
    """Handles --clear CLI option."""
    with file_lock(LOCK_FILE_PATH):
        if RECENT_CHANGES_PATH.is_file():
            RECENT_CHANGES_PATH.unlink()
    print("Cleared recent changes log.")


def display_inventory_summary(inv: dict[str, Any], slow_threshold: float) -> None:
    """Prints standard inventory generation summary."""
    print(f"Generated test inventory at {normalize_repo_rel(TEST_INVENTORY_PATH)}")
    print(f"Total Tests Indexed : {inv['total_tests']}")
    print(f"Total Packages      : {inv['summary']['packages']}")
    print(f"Slow Tests (> {slow_threshold}s) : {inv['summary']['slow_tests']}")
    print(f"Fast Tests (<= {slow_threshold}s): {inv['summary']['fast_tests']}")
    s = inv["summary"]
    print(f"Heatmap Distribution : Hot={s['hot_tests']}, Warm={s['warm_tests']}, Cold={s['cold_tests']}")


def handle_inventory_run(args: argparse.Namespace) -> None:
    """Handles full test inventory generation and heatmap presentation."""
    inv = build_test_inventory(REPO_ROOT, args.slow_threshold, args.force_run_all, args.commits)
    if args.json:
        print(json.dumps(inv, indent=2))
        return
    if args.heatmap:
        display_heatmap_report(inv, top_n=args.top)
        return
    display_inventory_summary(inv, args.slow_threshold)


def main() -> None:
    """Main CLI entrypoint."""
    parser = build_arg_parser()
    args = parser.parse_args()
    if args.check_age:
        handle_check_age(args)
    elif args.clear:
        handle_clear_recent_log()
    elif args.record:
        handle_record_changes(args.record)
    elif args.query_recent:
        handle_query_recent()
    else:
        handle_inventory_run(args)


if __name__ == "__main__":
    main()
