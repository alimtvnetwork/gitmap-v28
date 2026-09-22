#!/usr/bin/env python3
"""Fast Parallel Multi-Worker Local CI/CD Runner with Incremental Caching & Real-Time Telemetry.

====================================================================================================
                        AI AGENT OPERATIONAL INSTRUCTION MANUAL
====================================================================================================

1. PARALLEL EXECUTION MODEL & BATCH BARRIER PIPELINE:
   The runner organizes all 33 quality gates into 7 sequential batch barriers:
   • Batch 1 — Linters, AST Checks & Static Analyzers (Gates 1-22):
     - Parallelism: CPU-bound, executes across min(8, CPU count) worker threads.
     - Scope: Spell check, nested if, boolean/enum conventions, error codes, AST helptext, govulncheck.
   • Batch 2 — Compile Gates (Gates 23-24):
     - Parallelism: I/O-bound, restricted to max 2 workers to prevent disk lock contention.
     - Scope: Go Compile Gate (bin/gitmap.exe), Web App Build.
   • Batch 3 — Packaging Gates (Gate 25):
     - Parallelism: Single dedicated worker to isolate GoReleaser snapshot release and avoid dist/ collisions.
     - Scope: GoReleaser Snapshot Build.
   • Batch 4 — E2E Smoke & Integration Suites (Gates 26-28):
     - Parallelism: High concurrency parallel worker pool across isolated temporary environments.
     - Scope: E2E Smoke Suite, History Purge/Pin.
   • Batch 5 — Coverage Generation (Gate 31):
     - Parallelism: Dedicated worker with 20m timeout for deep codebase coverage generation.
     - Scope: Go Test Coverage Profile (coverage.out).
   • Batch 6 — Coverage Verification (Gate 32):
     - Parallelism: Fast assertion runner.
     - Scope: Coverage Floor Guard.
   • Batch 7 — Race Detection (Gate 33):
     - Parallelism: Dedicated worker for race condition detection across hot packages.
     - Scope: Go Test Race (Hot Packages).

2. REAL-TIME TELEMETRY & ARTIFACT STREAMING (.ai-memory/cicd/):
   All telemetry is written immediately to disk with unbuffered os.fsync flushing:
   • .ai-memory/cicd/errors.log          : Real-time append-only Markdown failure stream containing
                                         failing commands, cwd, env, suspect files, and stack traces.
   • .ai-memory/cicd/errors.json         : Structured JSON array of all active failures:
                                         [{"name": ..., "cmd": ..., "code": ..., "suspect_files": [...]}]
   • .ai-memory/cicd/events.jsonl        : Real-time append-only NDJSON event stream:
                                         {"timestamp": ..., "event": "gate_started"|"gate_failed"|...}
   • .ai-memory/cicd/changelog.log       : Chronological human-readable summary of pipeline events.
   • .ai-memory/cicd/summary.json        : Live status metadata ("status": "running"|"completed"|"failed",
                                         total, passed, failed, remaining, cached counts, active_failures).
   • .ai-memory/cicd/state.json          : Persistent incremental cache fingerprinting Git HEAD SHA,
                                         uncommitted dirty file hashes, tool script mtimes, and gate status.
   • .ai-memory/cicd/run.log             : Chronological record of all executed, cached, and failed gates.
   • .ai-memory/cicd/latest/             : Symlink / junction pointing directly to current session folder
                                         under .ai-memory/cicd/runs/<timestamp>/.

3. AI AGENT PARALLEL REMEDIATION PLAYBOOK:
   When executing tasks or autonomous repair loops, follow this protocol:
   • Step 1 (Early Interception): Do NOT wait for the entire 33-gate suite to finish. As soon as an
     immediate failure banner appears or .ai-memory/cicd/errors.json contains an entry, begin remediation.
   • Step 2 (Inspect Suspect Files): View .ai-memory/cicd/errors.log using view_file to examine
     extracted suspect files and the root-cause stack trace.
   • Step 3 (Surgical Code Fix): Edit offending files following 02-spec/02-coding-guidelines/ (functions <= 15
     lines, blank line before returns, affirmative booleans, zero swallowed exceptions).
   • Step 4 (Targeted Verification): Validate the fix in isolation before running the full suite:
     python 03-ai-scripts/06-cicd-local-runner.py --filter "<Gate Name>"
     Or execute the underlying direct tool command as indicated in errors.log.
   • Step 5 (Fast Incremental Pass): Re-run the runner:
     python 03-ai-scripts/06-cicd-local-runner.py
     The incremental cache will instantly skip passing gates (~0.5ms each), validating only your fix.

4. INCREMENTAL CACHING MECHANICS & SKIP RULES:
   • Two-Tier Fingerprinting:
     - Tier 1: Computes repo delta Δ_repo from git rev-parse HEAD and git status --porcelain=v1 -uall.
     - Tier 2: Evaluates GateSpec matching rules against modified file paths, tool scripts, and configs.
   • Skip Evaluation (O(1) in-memory, ~0.5ms per gate):
     - Unchanged gates that passed in the previous run are skipped with label "[cached]".
     - Previously failed gates are NEVER skipped and always re-execute until green.
     - Upstream Invalidation: If "Go Compile Gate" rebuilds bin/gitmap.exe, downstream gates
       ("E2E Smoke Suite", "History Pin/Purge") are automatically invalidated.
   • Cache Bypass:
     - Use --force (or --fresh / --clean) to clear cache and force all 33 gates to run from scratch.

5. USAGE & CLI CHEAT SHEET:
   python 03-ai-scripts/06-cicd-local-runner.py                    # Fast incremental run (quiet on pass)
   python 03-ai-scripts/06-cicd-local-runner.py --filter "<name>"  # Test specific gate (e.g. --filter "Nested")
   python 03-ai-scripts/06-cicd-local-runner.py --all-paths       # Show all 33 gates (passed and cached)
   python 03-ai-scripts/06-cicd-local-runner.py --force           # Purge cache and re-execute all gates
   python 03-ai-scripts/06-cicd-local-runner.py --json            # Emit machine-readable execution JSON
   python 03-ai-scripts/06-cicd-local-runner.py --sync            # Run sequentially (1 worker thread)
"""
from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import asdict, dataclass
import fnmatch
import hashlib
import importlib
import json
import os
from pathlib import Path
import re
import shutil
import stat
import subprocess
import sys
import tempfile
import threading
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

CPU_CORES = os.cpu_count() or 16


def resolve_cpu_scaling_multiplier(free_pct: float, mem_free_pct: float, aggressive: bool) -> float:
    """Calculates thread pool multiplier based on CPU idle headroom and available memory."""
    if free_pct >= 70.0 and mem_free_pct >= 20.0:
        return 2.0 if aggressive else 1.5
    if free_pct >= 45.0:
        return 1.5 if aggressive else 1.25
    if free_pct >= 25.0:
        return 1.0

    return 0.75


def detect_cpu_freeness_and_workers(
    min_workers: int = 8,
    max_cap: int | None = None,
    aggressive: bool = True,
) -> tuple[float, int]:
    """Inspects system CPU freeness and dynamically calculates optimal worker threads.

    Evaluates idle capacity, physical vs logical core headroom, and memory availability
    to maximize hardware saturation up to 85-95% CPU without system starvation.
    """
    logical_cores = os.cpu_count() or 16
    free_pct = 75.0
    mem_free_pct = 60.0
    try:
        import psutil
        busy_pct = psutil.cpu_percent(interval=0.15)
        free_pct = max(1.0, 100.0 - busy_pct)
        mem = psutil.virtual_memory()
        mem_free_pct = (mem.available / mem.total) * 100.0
    except Exception:
        pass

    multiplier = resolve_cpu_scaling_multiplier(free_pct, mem_free_pct, aggressive)
    scaled_workers = int(logical_cores * multiplier)
    optimal = max(min_workers, scaled_workers)
    if max_cap is not None:
        optimal = min(optimal, max_cap)
    else:
        optimal = min(optimal, logical_cores * 3)

    # Memory Safety Guard: On machines with limited memory (<=16GB RAM or <8GB available),
    # cap worker threads to prevent Windows commit limit exhaustion (errno 1455)
    # which crashes IDE language servers and background tasks.
    try:
        import psutil
        vm = psutil.virtual_memory()
        if vm.total <= 16 * 1024**3 or vm.available < 8 * 1024**3:
            optimal = min(optimal, 4)
    except Exception:
        if os.name == "nt":
            optimal = min(optimal, 4)

    return round(free_pct, 1), optimal


INITIAL_FREE_CPU_PCT, DYNAMIC_WORKERS = detect_cpu_freeness_and_workers(min_workers=4)
DEFAULT_WORKERS = int(os.environ.get("CI_MAX_WORKERS", DYNAMIC_WORKERS))
DEFAULT_IO_WORKERS = int(os.environ.get("CI_MAX_IO_WORKERS", min(8, CPU_CORES)))
DEFAULT_TIMEOUT_SEC = int(os.environ.get("CI_TIMEOUT_SEC", 1200))
DEFAULT_ENCODING = "utf-8"
DEFAULT_JOB_ESTIMATE_SEC = 5.0
DEFAULT_HEARTBEAT_INTERVAL = float(os.environ.get("RUNNER_HEARTBEAT_INTERVAL", 25.0))
REPO_ROOT = Path(__file__).resolve().parent.parent

TMP_CACHE_DIR = REPO_ROOT / ".ai-memory" / "temp"
TMP_CACHE_DIR.mkdir(parents=True, exist_ok=True)
FAILURES_DIR = TMP_CACHE_DIR / "failures"
FAILURES_DIR.mkdir(parents=True, exist_ok=True)
RUNNER_ETA_FILE = TMP_CACHE_DIR / "runner-eta.json"


def get_repo_os_temp_dir(*subdirs: str) -> Path:
    """Returns a directory under the OS temp directory scoped by repository name."""
    target = Path(tempfile.gettempdir()) / "gitmap"
    if subdirs:
        target = target.joinpath(*subdirs)
    target.mkdir(parents=True, exist_ok=True)

    return target


def remove_readonly(func: Any, path: str, exc_info: Any) -> None:
    """Clears the Windows read-only attribute on files (such as git objects/packs) and retries removal."""
    try:
        os.chmod(path, stat.S_IWRITE)
        func(path)
    except Exception:
        pass


def robust_rmtree(path: Path | str) -> None:
    """Removes a directory tree robustly, handling Windows read-only git pack files."""
    p = Path(path)
    if p.exists():
        try:
            shutil.rmtree(p, onerror=remove_readonly)
        except OSError:
            pass


def clear_repo_build_temp() -> None:
    """Purges previous build artifacts and sweeps repo build temp to prevent disk waste."""
    build_dir = Path(tempfile.gettempdir()) / "gitmap" / "build"
    robust_rmtree(build_dir)
    build_dir.mkdir(parents=True, exist_ok=True)


def clear_repo_test_temp() -> None:
    """Sweeps repo test, sandbox, and auxiliary temporary directories to ensure test hygiene and free disk space."""
    temp_base = Path(tempfile.gettempdir()) / "gitmap"
    for category in ("test", "sandbox", "purge", "downloads", "handoff"):
        cat_dir = temp_base / category
        robust_rmtree(cat_dir)
        cat_dir.mkdir(parents=True, exist_ok=True)


def clear_stale_failures_log() -> None:
    """Clears stale test failure logs before running tests to prevent storage accumulation."""
    failures_dir = REPO_ROOT / ".ai-memory" / "temp" / "failures"
    robust_rmtree(failures_dir)
    failures_dir.mkdir(parents=True, exist_ok=True)


def prune_old_cicd_runs(keep_count: int = 5) -> None:
    """Prunes older session run directories in .ai-memory/cicd/runs to prevent disk waste."""
    runs_dir = REPO_ROOT / ".ai-memory" / "cicd" / "runs"
    if not runs_dir.exists():
        return
    try:
        run_dirs = [p for p in runs_dir.iterdir() if p.is_dir()]
        if len(run_dirs) <= keep_count:
            return
        run_dirs.sort(key=lambda p: p.stat().st_mtime, reverse=True)
        for old_dir in run_dirs[keep_count:]:
            robust_rmtree(old_dir)
    except Exception:
        pass


def clean_stale_temp_artifacts() -> None:
    """Removes orphaned build caches, old e2e sandboxes, and stale binaries from .ai-memory/temp to prevent bloat."""
    lovable_temp = REPO_ROOT / ".ai-memory" / "temp"
    if not lovable_temp.exists():
        return
    for item in lovable_temp.iterdir():
        if item.is_dir() and (item.name.startswith("go-build") or item.name in ("node-compile-cache", "gitmap", "cicd")):
            robust_rmtree(item)
        elif item.is_file() and item.name in ("gitmap", "gitmap.exe"):
            try:
                item.unlink(missing_ok=True)
            except OSError:
                pass


def clean_coverage_artifacts() -> None:
    """Removes loose intermediate coverage profile files after verification."""
    for cov_file in (REPO_ROOT / "coverage.out", REPO_ROOT / "cli" / "coverage.out"):
        if cov_file.exists():
            try:
                cov_file.unlink(missing_ok=True)
            except OSError:
                pass


def clear_repo_go_cache() -> None:
    """Purges Go build and test cache from repo-scoped temp and system to prevent disk accumulation."""
    gocache_dir = Path(tempfile.gettempdir()) / "gitmap" / "gocache"
    robust_rmtree(gocache_dir)
    gocache_dir.mkdir(parents=True, exist_ok=True)
    try:
        subprocess.run(
            ["go", "clean", "-cache", "-testcache"],
            cwd=str(REPO_ROOT / "cli"),
            capture_output=True,
            timeout=30,
            check=False,
        )
    except Exception:
        pass


def clear_worker_go_cache(worker_cache_dir: Path | str) -> None:
    """Purges an isolated worker's Go build cache immediately after package test execution."""
    robust_rmtree(worker_cache_dir)


REPO_OS_TEMP = get_repo_os_temp_dir()
REPO_BUILD_TEMP = get_repo_os_temp_dir("build")
REPO_TEST_TEMP = get_repo_os_temp_dir("test")
REPO_GOCACHE_TEMP = get_repo_os_temp_dir("gocache")

os.environ["GOTMPDIR"] = str(REPO_BUILD_TEMP)
os.environ["TMPDIR"] = str(REPO_TEST_TEMP)
os.environ["GOCACHE"] = str(REPO_GOCACHE_TEMP)
if os.name != "nt":
    os.environ["TEMP"] = str(REPO_TEST_TEMP)
    os.environ["TMP"] = str(REPO_TEST_TEMP)

CICD_DIR = REPO_ROOT / ".ai-memory" / "cicd"
CICD_DIR.mkdir(parents=True, exist_ok=True)
CACHE_DIR = CICD_DIR / "cache"
CACHE_DIR.mkdir(parents=True, exist_ok=True)

CACHE_STATE_FILE = CACHE_DIR / "state.json"
CACHE_INVENTORY_FILE = CACHE_DIR / "inventory.json"
CACHE_TIMINGS_FILE = CACHE_DIR / "timings.json"
CACHE_DEBOUNCE_FILE = CACHE_DIR / "debounce.json"

CICD_TEMP_DIR = CICD_DIR
CICD_RUNS_DIR = CICD_DIR / "runs"
CICD_LATEST_DIR = CICD_DIR / "latest"
CICD_ERRORS_LOG = CICD_DIR / "errors.log"
CICD_ERRORS_JSON = CICD_DIR / "errors.json"
CICD_RUN_LOG = CICD_DIR / "run.log"
CICD_EVENTS_JSONL = CICD_DIR / "events.jsonl"
CICD_CHANGELOG_LOG = CICD_DIR / "changelog.log"
CICD_SUMMARY_JSON = CICD_DIR / "summary.json"
CICD_STATE_JSON = CACHE_STATE_FILE
CICD_POINTER_FILE = CICD_DIR / "latest_run.txt"
CICD_LAST_RUN_CACHE = CACHE_DEBOUNCE_FILE

TIMING_FILE_PATH = CACHE_TIMINGS_FILE
TEST_INVENTORY_PATH = REPO_ROOT / ".ai-memory" / "test-inventory.json"
TEST_INVENTORY_CACHE_PATH = CACHE_INVENTORY_FILE

DISK_WRITE_LOCK = threading.RLock()

JOB_BATCHES: list[dict[str, Any]] = [
    {
        "name": "Linters & AST Checks",
        "max_workers": min(4, DEFAULT_WORKERS),
        "jobs": {

            "Go Format Check": [sys.executable, ".github/scripts/go-format-check.py", "--no-commit"],
            "Spell Check (misspell)": [sys.executable, ".github/scripts/misspell-changed.py"],
            "Nested If Linter": [sys.executable, "linter-scripts/check-nested-ifs.py"],

            "Boolean & Enum Linter": [sys.executable, "linter-scripts/check-enum-and-boolean.py"],
            "Boolean Guidelines Linter": [sys.executable, "linter-scripts/check-boolean-guidelines.py"],
            "Enum Guidelines Linter": [sys.executable, "linter-scripts/check-enum-guidelines.py"],
            "Schema Guidelines Linter": [sys.executable, "linter-scripts/check-schema-guidelines.py"],
            "Error Management Check": [sys.executable, "linter-scripts/check-error-management.py"],
            "Relative Path Check": [sys.executable, "linter-scripts/check-relative-paths.py"],
            "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
            "MWS Error Codes Check": [sys.executable, "linter-scripts/check-mws-error-codes.py"],
            "Interface Naming Check": [sys.executable, "linter-scripts/check-interface-naming.py"],
            "CLI Help Parity Check": [sys.executable, "03-ai-scripts/09-cli-help-auditor.py"],
            "Constants Registry AST Check": ["go", "test", "-C", "cli", "./constants/...", "-run", "TestTopLevelCmdRegistryMatchesAST", "-count=1"],
            "Constants Collision Check": ["go", "test", "-C", "cli", "./constants/...", "-run", "TestTopLevelCmdConstantsAreUnique", "-count=1"],
            "Helptext Parity Check": ["go", "test", "-C", "cli", "./helptext/...", "-count=1"],
            "govulncheck": [sys.executable, ".github/scripts/check-vulncheck.py"],
            "Startup Build-Tags (linux)": {"cmd": ["go", "build", "./startup/..."], "cwd": "cli", "env": {"GOOS": "linux", "GOARCH": "amd64", "CGO_ENABLED": "0"}},
            "Startup Build-Tags (darwin)": {"cmd": ["go", "build", "./startup/..."], "cwd": "cli", "env": {"GOOS": "darwin", "GOARCH": "amd64", "CGO_ENABLED": "0"}},
            "Startup Build-Tags (windows)": {"cmd": ["go", "build", "./startup/..."], "cwd": "cli", "env": {"GOOS": "windows", "GOARCH": "amd64", "CGO_ENABLED": "0"}},
            "golangci-lint (strict)": {"cmd": ["golangci-lint", "run", "--allow-serial-runners", "--issues-exit-code=1", "--timeout=10m", "-c", ".golangci.yml", "--path-prefix", "cli", "./..."], "cwd": "cli"},
            "Unused Code Guard (unused)": [sys.executable, ".github/scripts/check-unused-diff.py"],
            "Gosec G115 Guard (overflow)": [sys.executable, ".github/scripts/check-gosec-diff.py"],
            "GoCritic Guard (style)": [sys.executable, ".github/scripts/check-gocritic-diff.py"],
            "Cross-OS Vet (Windows)": {"cmd": ["go", "vet", "-C", "cli", "./..."], "env": {"GOOS": "windows", "GOARCH": "amd64"}},
            "Cross-OS Vet (Darwin)": {"cmd": ["go", "vet", "-C", "cli", "./..."], "env": {"GOOS": "darwin", "GOARCH": "amd64"}},
            "cmd/ Naming Check": [sys.executable, ".github/scripts/check-cmd-naming.py", "cli/cmd"],
            "Legacy Refs Check": [sys.executable, ".github/scripts/check-legacy-refs.py", "."],
            "Deploy Layout Check": [sys.executable, ".github/scripts/check-deploy-layout.py", "."],
            "constants/ Naming Check": [sys.executable, ".github/scripts/check-constants-naming.py"],
            "Golden Allow Leak Check": [sys.executable, ".github/scripts/check-no-golden-allow-leak.py"],
            "Bare Stderr Check": [sys.executable, ".github/scripts/check-bare-stderr-err.py"],
            "Changelog Version Sync": [sys.executable, ".github/scripts/check-changelog-version-sync.py"],
            "File Size Check": [sys.executable, ".github/scripts/file-size-check.py", "200"],
            "JSON Snapshot Fast Check": ["go", "test", "-C", "cli", "./cmd/...", "./formatter/...", "-failfast", "-count=1", "-run", "^(TestStartupListJSON|TestFindNextJSONContract|TestLatestBranchJSONContract|TestSchemaRegistry|TestAssertGoldenBytesDeterministic|TestExpectDelim|TestScanEveryObjectKeysPure|TestWriteJSON|TestWriteCSV)"],
        },
    },
    {
        "name": "Compile Gates",
        "max_workers": DEFAULT_IO_WORKERS,
        "jobs": {
            "Go Compile Gate": ["go", "build", "-C", "cli", "-o", "../bin/gitmap.exe", "."],
            "Web App Build": ["npm", "run", "build"],
        },
    },
    {
        "name": "Packaging Gates",
        "max_workers": 1,
        "jobs": {
            "GoReleaser Snapshot Build": {"cmd": ["goreleaser", "build", "--snapshot", "--clean", "--single-target"], "cwd": "cli"},
        },
    },
    {
        "name": "E2E Smoke Tests",
        "max_workers": None,
        "jobs": {
            "E2E Smoke Suite": [sys.executable, ".github/scripts/e2e-cli-smoke.py", "bin/gitmap.exe"],
            "History Purge Smoke": [sys.executable, ".github/scripts/smoke-history-purge.py", "bin/gitmap.exe"],
            "History Pin Smoke": [sys.executable, ".github/scripts/smoke-history-pin.py", "bin/gitmap.exe"],
        },
    },
    {
        "name": "Smart Unit Tests & Coverage",
        "max_workers": DEFAULT_WORKERS,
        "jobs": {
            "Go Smart Incremental Tests": {"type": "smart_go_tests", "cmd": ["go", "test", "smart-incremental"], "cwd": "cli"},
            "Go Test Coverage Profile": {
                "cmd": [
                    sys.executable, "-c",
                    "import subprocess, shutil; subprocess.run(['go', 'test', '-count=1', '-coverprofile=coverage.out', './visibility/...', './cmd/commitin/checkpoint/...', './jsonenv/...', './logging/...', './transport/...'], cwd='cli', check=True); shutil.copyfile('cli/coverage.out', 'coverage.out')"
                ]
            },
        },
    },
    {
        "name": "Python CI Script Unit Tests",
        "max_workers": 1,
        "jobs": {
            "Lint Script Unit Tests": [sys.executable, ".github/scripts/tests/test_ci_scripts.py"],
        },
    },
    {
        "name": "Coverage Verification",
        "max_workers": 1,
        "jobs": {
            "Coverage Floor Guard": [sys.executable, ".github/scripts/coverage-floor.py", "coverage.out"],
        },
    },
    {
        "name": "Race Detection",
        "max_workers": 1,
        "jobs": {
            "Go Test Race (Hot Packages)": {"cmd": ["go", "test", "-p", str(min(4, DEFAULT_WORKERS)), "-parallel", str(min(4, DEFAULT_WORKERS)), "-count=1", "-timeout=15m", "./cmd/...", "./cmdagy/...", "./cmdchromeprofile/...", "./cloneconcurrency/...", "./visibility/...", "./store/...", "./uipref/..."], "cwd": "cli", "env": {"GITMAP_IN_MEMORY_DB": "1"}},
        },
    },
]

ANSI_ESCAPE_REGEX = re.compile(r"\x1B(?:\[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])")
GLOBAL_TIMINGS: dict[str, float] = {}
EXCLUDE_DEFAULTS = [
    ".git/**", "node_modules/**", "dist/**", "bin/**", "vendor/**", ".ai-memory/**", ".tmp/**",
]

CLUSTER_GO_ALL = ["cli/**", "cli-updater/**/*.go", "cli/go.mod", "cli/go.sum"]
CLUSTER_GO_CMD = ["cli/cmd/**/*.go", "cli/cmdagy/**/*.go", "cli/cmdchromeprofile/**/*.go", "cli/constants/**/*.go", "cli/go.mod", "cli/go.sum"]
CLUSTER_GO_CONSTANTS = ["cli/constants/**/*.go", "cli/go.mod"]
CLUSTER_GO_HELPTEXT = ["cli/helptext/**/*.go", "cli/cmd/**/*.go", "cli/cmdagy/**/*.go", "cli/cmdchromeprofile/**/*.go", "cli/constants/**/*.go", "cli/go.mod"]
CLUSTER_GO_STARTUP = ["cli/startup/**/*.go", "cli/go.mod"]
CLUSTER_GO_RACE = ["cli/cmd/**/*.go", "cli/cmdagy/**/*.go", "cli/cmdchromeprofile/**/*.go", "cli/cloneconcurrency/**/*.go", "cli/visibility/**/*.go", "cli/store/**/*.go", "cli/uipref/**/*.go", "cli/go.mod"]
CLUSTER_WEB_APP = ["src/**/*", "public/**/*", "index.html", "package.json", "package-lock.json", "vite.config.ts", "tsconfig*.json", "tailwind.config.ts", "postcss.config.js"]
CLUSTER_LINTER_SCRIPTS = ["linter-scripts/**/*.py", ".github/scripts/**/*.py"]
CLUSTER_REPO_TEXT = ["cli/**", "src/**", "spec/**", "docs/**", "03-ai-scripts/**", "linter-scripts/**", ".github/**", "*.md", "*.json", "*.yml", "*.yaml"]
CLUSTER_MWS = ["02-spec/19-main-worker-service/**", "02-spec/14-update/**", "02-spec/03-error-manage/03-error-code-registry/**", "src/**/*.{ts,tsx}", "linter-scripts/check-mws-error-codes.*"]


def strip_ansi(text: str) -> str:
    """Removes terminal ANSI color escape codes from text."""
    clean_text = ANSI_ESCAPE_REGEX.sub("", text)

    return clean_text


@dataclass
class JobResult:
    """Encapsulates the execution outcome of an individual quality gate."""
    name: str
    cmd: list[str] | dict[str, Any]
    code: int | str
    out: str
    err: str
    elapsed: float
    cwd: str | None = None
    env_overrides: dict[str, str] | None = None
    is_cached: bool = False

    @property
    def is_success(self) -> bool:
        has_passed = (self.code == 0)

        return has_passed

    @property
    def is_timeout(self) -> bool:
        has_timed_out = (self.code == "timeout")

        return has_timed_out


def record_job_timing(job_name: str, elapsed_sec: float) -> None:
    """Stores the execution duration of a job for persistence."""
    GLOBAL_TIMINGS[job_name] = elapsed_sec


def format_duration(seconds: float | int) -> str:
    """Formats duration in seconds; if >= 60s, displays both seconds and minutes/seconds."""
    sec_float = float(seconds)
    is_under_minute = sec_float < 60.0
    if is_under_minute:
        return f"{sec_float:.2f}s"
    mins = int(sec_float // 60)
    rem_sec = sec_float % 60

    return f"{sec_float:.2f}s ({mins}m {rem_sec:.1f}s)"


def normalize_repo_rel(path_input: Any) -> str:
    """Normalizes path to forward-slash relative path within the repo."""
    if path_input is None:
        return ""
    raw_str = str(path_input).strip()
    try:
        p = Path(raw_str)
        if p.is_absolute():
            resolved = p.resolve()
            repo_resolved = REPO_ROOT.resolve()
            try:
                rel = resolved.relative_to(repo_resolved)
                raw_str = str(rel)
            except ValueError:
                rel = os.path.relpath(str(resolved), str(repo_resolved))
                if not rel.startswith(".."):
                    raw_str = rel
    except Exception:
        pass
    norm = raw_str.replace("\\", "/").strip()
    if norm.startswith("./"):
        norm = norm[2:]

    return norm


def expand_brace_patterns(pattern: str) -> list[str]:
    """Expands comma-separated brace sets into multiple patterns."""
    match = re.search(r"\{([^{}]+)\}", pattern)
    if not match:
        return [pattern]
    prefix = pattern[:match.start()]
    suffix = pattern[match.end():]
    expanded: list[str] = []
    for opt in match.group(1).split(","):
        expanded.extend(expand_brace_patterns(f"{prefix}{opt.strip()}{suffix}"))

    return expanded


def glob_to_regex(pattern: str) -> re.Pattern:
    """Translates standard recursive glob patterns into regex."""
    norm = normalize_repo_rel(pattern)
    escaped = re.escape(norm)
    res = escaped.replace(r"\*\*/", r"(?:.*/)?")
    res = res.replace(r"\*\*", r".*")
    res = res.replace(r"\*", r"[^/]*")
    res = res.replace(r"\?", r"[^/]")

    return re.compile(f"^{res}$", re.IGNORECASE)


class GateSpec:
    """Declarative specification of gate inputs, tool scripts, and artifact dependencies."""

    def __init__(
        self, name: str, tool_scripts: list[str] | None = None, configs: list[str] | None = None,
        relevant_patterns: list[str] | None = None, exclude_patterns: list[str] | None = None,
        artifact_inputs: list[str] | None = None, artifact_outputs: list[str] | None = None, upstream_gates: list[str] | None = None,
    ):
        self.name = name
        self._init_paths(tool_scripts, configs, artifact_inputs, artifact_outputs, upstream_gates)
        self._rel_res = self._compile_patterns(relevant_patterns or [])
        self._exc_res = self._compile_patterns((exclude_patterns or []) + EXCLUDE_DEFAULTS)

    def _init_paths(self, tools: list[str] | None, cfgs: list[str] | None, in_arts: list[str] | None, out_arts: list[str] | None, ups: list[str] | None) -> None:
        self.tool_scripts = [normalize_repo_rel(p) for p in (tools or [])]
        self.configs = [normalize_repo_rel(p) for p in (cfgs or [])]
        self.artifact_inputs = [normalize_repo_rel(p) for p in (in_arts or [])]
        self.artifact_outputs = [normalize_repo_rel(p) for p in (out_arts or [])]
        self.upstream_gates = ups or []

    def _compile_patterns(self, raw_patterns: list[str]) -> list[re.Pattern]:
        compiled: list[re.Pattern] = []
        for p in raw_patterns:
            for exp in expand_brace_patterns(p):
                compiled.append(glob_to_regex(exp))

        return compiled

    def matches_path(self, path: str) -> bool:
        """Determines whether a repository file change triggers this gate."""
        norm_path = normalize_repo_rel(path)
        if norm_path in self.tool_scripts or norm_path in self.configs:
            return True
        if any(rx.match(norm_path) for rx in self._exc_res):
            return False
        has_match = any(rx.match(norm_path) for rx in self._rel_res)

        return has_match


GATE_SPECS: dict[str, GateSpec] = {
    "Go Format Check": GateSpec("Go Format Check", tool_scripts=[".github/scripts/go-format-check.py"], relevant_patterns=["cli/**/*.go"]),
    "Spell Check (misspell)": GateSpec("Spell Check (misspell)", tool_scripts=[".github/scripts/misspell-changed.py"], configs=[".misspell-ignore"], relevant_patterns=CLUSTER_REPO_TEXT, exclude_patterns=["cli/completion/allcommands_generated.go"]),
    "Nested If Linter": GateSpec("Nested If Linter", tool_scripts=["linter-scripts/check-nested-ifs.py"], relevant_patterns=["cli/**/*.go", "src/**/*.{ts,tsx,js,jsx}"]),

    "Boolean & Enum Linter": GateSpec("Boolean & Enum Linter", tool_scripts=["linter-scripts/check-enum-and-boolean.py"], relevant_patterns=["cli/**/*.go", "src/**/*.{ts,tsx}"]),
    "Boolean Guidelines Linter": GateSpec("Boolean Guidelines Linter", tool_scripts=["linter-scripts/check-boolean-guidelines.py"], configs=["02-spec/02-coding-guidelines/**"], relevant_patterns=["cli/**/*.go", "src/**/*.{ts,tsx,js,jsx}"]),
    "Enum Guidelines Linter": GateSpec("Enum Guidelines Linter", tool_scripts=["linter-scripts/check-enum-guidelines.py"], configs=["02-spec/02-coding-guidelines/**"], relevant_patterns=["cli/**/*.go", "src/**/*.{ts,tsx}"]),
    "Error Management Check": GateSpec("Error Management Check", tool_scripts=["linter-scripts/check-error-management.py"], configs=["02-spec/03-error-manage/**"], relevant_patterns=["cli/**/*.go", "src/**/*.{ts,tsx}"]),
    "Relative Path Check": GateSpec("Relative Path Check", tool_scripts=["linter-scripts/check-relative-paths.py"], relevant_patterns=CLUSTER_REPO_TEXT, exclude_patterns=["*.png", "*.jpg", "*.exe", "*.zip", "*.sqlite"]),
    "Newline Styling Check": GateSpec("Newline Styling Check", tool_scripts=["linter-scripts/check-newline-styling.py"], relevant_patterns=["src/**/*.{ts,tsx,js}", "cli/**/*.go"]),
    "MWS Error Codes Check": GateSpec("MWS Error Codes Check", tool_scripts=["linter-scripts/check-mws-error-codes.py"], configs=["02-spec/19-main-worker-service/13-error-codes.md", "02-spec/19-main-worker-service/error-codes.json", "02-spec/03-error-manage/03-error-code-registry/error-codes-master.json", "linter-scripts/check-mws-error-codes.waivers.txt", "linter-scripts/check-mws-error-codes.unallocated.txt"], relevant_patterns=CLUSTER_MWS),
    "Interface Naming Check": GateSpec("Interface Naming Check", tool_scripts=["linter-scripts/check-interface-naming.py"], relevant_patterns=["cli/**/*.go"]),
    "CLI Help Parity Check": GateSpec("CLI Help Parity Check", tool_scripts=["03-ai-scripts/09-cli-help-auditor.py"], configs=["03-ai-scripts/02-shared-engine.py"], relevant_patterns=["cli/cmd/**/*.go", "cli/helptext/**/*.go", "03-ai-scripts/**/*.py"]),
    "Constants Registry AST Check": GateSpec("Constants Registry AST Check", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_CMD),
    "Constants Collision Check": GateSpec("Constants Collision Check", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_CONSTANTS),
    "Helptext Parity Check": GateSpec("Helptext Parity Check", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_HELPTEXT),
    "Lint Script Unit Tests": GateSpec("Lint Script Unit Tests", tool_scripts=[".github/scripts/tests/test_ci_scripts.py"], relevant_patterns=CLUSTER_LINTER_SCRIPTS + [".github/scripts/tests/**"], exclude_patterns=["__pycache__/**"]),
    "govulncheck": GateSpec("govulncheck", tool_scripts=[".github/scripts/check-vulncheck.py"], configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL),
    "Startup Build-Tags (linux)": GateSpec("Startup Build-Tags (linux)", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_STARTUP),
    "Startup Build-Tags (darwin)": GateSpec("Startup Build-Tags (darwin)", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_STARTUP),
    "Startup Build-Tags (windows)": GateSpec("Startup Build-Tags (windows)", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_STARTUP),
    "golangci-lint (strict)": GateSpec("golangci-lint (strict)", configs=[".golangci.yml", "cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL),
    "Unused Code Guard (unused)": GateSpec("Unused Code Guard (unused)", tool_scripts=[".github/scripts/check-unused-diff.py", ".github/scripts/check-single-linter-diff.py"], configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL),
    "Gosec G115 Guard (overflow)": GateSpec("Gosec G115 Guard (overflow)", tool_scripts=[".github/scripts/check-gosec-diff.py", ".github/scripts/check-single-linter-diff.py"], configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL),
    "GoCritic Guard (style)": GateSpec("GoCritic Guard (style)", tool_scripts=[".github/scripts/check-gocritic-diff.py", ".github/scripts/check-single-linter-diff.py"], configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL),
    "Cross-OS Vet (Windows)": GateSpec("Cross-OS Vet (Windows)", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL),
    "Cross-OS Vet (Darwin)": GateSpec("Cross-OS Vet (Darwin)", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL),
    "Go Compile Gate": GateSpec("Go Compile Gate", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL, artifact_outputs=["bin/gitmap.exe"]),
    "Web App Build": GateSpec("Web App Build", configs=["package.json", "package-lock.json", "vite.config.ts", "tsconfig*.json"], relevant_patterns=CLUSTER_WEB_APP, artifact_outputs=["dist/**"]),
    "GoReleaser Snapshot Build": GateSpec("GoReleaser Snapshot Build", configs=[".goreleaser.yaml", ".goreleaser.yml", "cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL, artifact_outputs=["dist/**"]),
    "E2E Smoke Suite": GateSpec("E2E Smoke Suite", tool_scripts=[".github/scripts/e2e-cli-smoke.py"], relevant_patterns=CLUSTER_GO_ALL + [".github/scripts/e2e-cli-smoke.py"], artifact_inputs=["bin/gitmap.exe"], upstream_gates=["Go Compile Gate"]),
    "History Purge Smoke": GateSpec("History Purge Smoke", tool_scripts=[".github/scripts/smoke-history-purge.py"], relevant_patterns=[".github/scripts/smoke-history-purge.py", "cli/cmd/**/*.go", "cli/store/**/*.go"], artifact_inputs=["bin/gitmap.exe"], upstream_gates=["Go Compile Gate"]),
    "History Pin Smoke": GateSpec("History Pin Smoke", tool_scripts=[".github/scripts/smoke-history-pin.py"], relevant_patterns=[".github/scripts/smoke-history-pin.py", "cli/cmd/**/*.go", "cli/store/**/*.go"], artifact_inputs=["bin/gitmap.exe"], upstream_gates=["Go Compile Gate"]),
    "Go Smart Incremental Tests": GateSpec("Go Smart Incremental Tests", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL),
    "Go Test Coverage Profile": GateSpec("Go Test Coverage Profile", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_ALL, artifact_outputs=["coverage.out"]),
    "Coverage Floor Guard": GateSpec("Coverage Floor Guard", tool_scripts=[".github/scripts/coverage-floor.py"], artifact_inputs=["coverage.out"], upstream_gates=["Go Test Coverage Profile"]),
    "Go Test Race (Hot Packages)": GateSpec("Go Test Race (Hot Packages)", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=CLUSTER_GO_RACE),
    "cmd/ Naming Check": GateSpec("cmd/ Naming Check", tool_scripts=[".github/scripts/check-cmd-naming.py"], relevant_patterns=["cli/cmd/**/*.go", "cli/cmdagy/**/*.go", "cli/cmdchromeprofile/**/*.go"]),
    "Legacy Refs Check": GateSpec("Legacy Refs Check", tool_scripts=[".github/scripts/check-legacy-refs.py"], relevant_patterns=CLUSTER_REPO_TEXT),
    "Deploy Layout Check": GateSpec("Deploy Layout Check", tool_scripts=[".github/scripts/check-deploy-layout.py"], relevant_patterns=CLUSTER_REPO_TEXT),
    "constants/ Naming Check": GateSpec("constants/ Naming Check", tool_scripts=[".github/scripts/check-constants-naming.py"], relevant_patterns=CLUSTER_GO_CONSTANTS),
    "Golden Allow Leak Check": GateSpec("Golden Allow Leak Check", tool_scripts=[".github/scripts/check-no-golden-allow-leak.py"], relevant_patterns=CLUSTER_REPO_TEXT),
    "Bare Stderr Check": GateSpec("Bare Stderr Check", tool_scripts=[".github/scripts/check-bare-stderr-err.py"], relevant_patterns=["cli/cmd/**/*.go", "cli/cmdagy/**/*.go", "cli/cmdchromeprofile/**/*.go"]),
    "Changelog Version Sync": GateSpec("Changelog Version Sync", tool_scripts=[".github/scripts/check-changelog-version-sync.py"], configs=["cli/constants/constants.go", "changelog.md"], relevant_patterns=["cli/constants/constants.go", "changelog.md"]),
    "File Size Check": GateSpec("File Size Check", tool_scripts=[".github/scripts/file-size-check.py"], relevant_patterns=CLUSTER_GO_ALL),
    "JSON Snapshot Fast Check": GateSpec("JSON Snapshot Fast Check", configs=["cli/go.mod", "cli/go.sum"], relevant_patterns=["cli/cmd/**/*.go", "cli/formatter/**/*.go"]),
}


def get_head_commit(repo_root: Path) -> str:
    """Retrieves current git HEAD commit SHA."""
    cmd = ["git", "rev-parse", "HEAD"]
    res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, check=False)
    if res.returncode == 0:
        return res.stdout.strip()

    return ""


def parse_porcelain_line(line: str, repo_root: Path) -> tuple[str, dict[str, Any]] | None:
    """Parses a single porcelain line into path and file stat metadata."""
    if len(line) < 4:
        return None
    raw_path = line[3:].strip().strip('"')
    if " -> " in raw_path:
        _, raw_path = raw_path.split(" -> ", 1)
    rel_path = normalize_repo_rel(raw_path)
    full_path = repo_root / rel_path
    mtime, size = 0.0, 0
    if full_path.is_file():
        stat = full_path.stat()
        mtime, size = stat.st_mtime, stat.st_size

    return rel_path, {"status": line[:2], "mtime": mtime, "size": size}


def get_dirty_files_map(repo_root: Path) -> dict[str, dict[str, Any]]:
    """Queries git status --porcelain=v1 -uall to get all dirty files."""
    cmd = ["git", "status", "--porcelain=v1", "-uall"]
    res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, check=False)
    if res.returncode != 0:
        return {}
    dirty_map: dict[str, dict[str, Any]] = {}
    for line in res.stdout.splitlines():
        parsed = parse_porcelain_line(line, repo_root)
        if parsed is not None:
            dirty_map[parsed[0]] = parsed[1]

    return dirty_map


def get_commit_diff_files(repo_root: Path, old_head: str, new_head: str) -> set[str]:
    """Gets committed file path differences between two commits."""
    if not old_head or not new_head or old_head == new_head:
        return set()
    cmd = ["git", "diff", "--name-only", f"{old_head}..{new_head}"]
    res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, check=False)
    if res.returncode == 0:
        return {normalize_repo_rel(p) for p in res.stdout.splitlines() if p.strip()}

    return set()


def is_ignored_repo_path(path_str: str) -> bool:
    """Checks if path falls inside internal directories to ignore."""
    ignored = (".git/", ".ai-memory/", ".tmp/", "node_modules/", "dist/", "bin/", "vendor/")
    has_match = any(path_str.startswith(prefix) for prefix in ignored)

    return has_match


def collect_dirty_delta(last_dirty: dict[str, Any], current_dirty: dict[str, Any]) -> set[str]:
    """Finds paths modified, added, or deleted between dirty snapshots."""
    delta: set[str] = set()
    for p, stat in current_dirty.items():
        if not is_ignored_repo_path(p):
            old = last_dirty.get(p)
            if not old or old.get("mtime") != stat["mtime"] or old.get("size") != stat["size"]:
                delta.add(p)
    for p in last_dirty:
        if not is_ignored_repo_path(p) and p not in current_dirty:
            delta.add(p)

    return delta


def compute_repo_delta(
    repo_root: Path, last_head: str, curr_head: str, last_dirty: dict, curr_dirty: dict
) -> set[str]:
    """Computes the set of changed repository files since the last recorded run."""
    delta = collect_dirty_delta(last_dirty, curr_dirty)
    committed = get_commit_diff_files(repo_root, last_head, curr_head)
    for p in committed:
        if not is_ignored_repo_path(p):
            delta.add(p)

    return delta


def compute_cmd_hash(raw_cmd: Any, env: dict[str, str] | None, cwd: str | None) -> str:
    """Generates MD5 hash for command invocation arguments and environment."""
    data = {"cmd": raw_cmd, "env": env or {}, "cwd": cwd or ""}
    serialized = json.dumps(data, sort_keys=True)

    return hashlib.md5(serialized.encode("utf-8")).hexdigest()


def check_file_stat(repo_root: Path, rel_path: str) -> tuple[float, int] | None:
    """Returns current mtime and size for file or None if missing."""
    target = repo_root / rel_path
    if not target.is_file():
        return None
    st = target.stat()

    return st.st_mtime, st.st_size


def check_scripts_modified(repo_root: Path, scripts: list[str], recorded: dict[str, Any]) -> str | None:
    """Checks if any tool script was modified since last run."""
    for script in scripts:
        current = check_file_stat(repo_root, script)
        prev = recorded.get(script)
        if current is None or prev != list(current):
            return f"Tool script modified: {script}"

    return None


def check_artifacts_valid(repo_root: Path, inputs: list[str], recorded: dict[str, Any]) -> str | None:
    """Verifies that all required upstream input artifacts exist and match stats."""
    for art in inputs:
        current = check_file_stat(repo_root, art)
        prev = recorded.get(art)
        if current is None or prev != list(current):
            return f"Required artifact missing or modified: {art}"

    return None


def check_prev_gate_status(record: dict[str, Any] | None, cmd_hash: str) -> str | None:
    """Checks previous pass status and command signature."""
    if not record or record.get("status") != "PASSED":
        return "Previous run did not pass"
    if record.get("cmd_hash") != cmd_hash:
        return "Command arguments or environment changed"

    return None


def check_upstream_and_scripts(spec: GateSpec, record: dict[str, Any], executed: set[str], root: Path) -> str | None:
    """Verifies upstream execution status and script timestamps."""
    for up in spec.upstream_gates:
        if up in executed:
            return f"Upstream gate {up!r} was re-executed"
    err = check_scripts_modified(root, spec.tool_scripts, record.get("tool_stats", {}))
    if err is not None:
        return err

    return check_artifacts_valid(root, spec.artifact_inputs, record.get("artifact_stats", {}))


def check_delta_match(spec: GateSpec, repo_delta: set[str]) -> str | None:
    """Checks if any changed file in repo delta matches gate spec patterns."""
    for changed in repo_delta:
        if spec.matches_path(changed):
            return f"File changed: {changed}"

    return None


def evaluate_gate_skip(
    spec: GateSpec, record: dict[str, Any] | None, cmd_hash: str, repo_delta: set[str], root: Path, executed: set[str]
) -> tuple[bool, str]:
    """Determines whether a gate can be skipped based on inputs and delta."""
    status_err = check_prev_gate_status(record, cmd_hash)
    if status_err is not None:
        return False, status_err
    dep_err = check_upstream_and_scripts(spec, record or {}, executed, root)
    if dep_err is not None:
        return False, dep_err
    delta_err = check_delta_match(spec, repo_delta)
    if delta_err is not None:
        return False, delta_err

    return True, "No input changes detected"


def clean_temp_file(tmp_path: Path) -> None:
    """Safely unlinks temporary file if it exists."""
    if tmp_path.exists():
        try:
            tmp_path.unlink()
        except OSError as err:
            sys.stderr.write(f"[WARN] Failed to unlink temp file {tmp_path}: {err}\n")


def retry_replace_or_overwrite(tmp_path: Path, file_path: Path, text: str, has_fsync: bool = False) -> None:
    """Attempts atomic replace with retry, falling back to direct write."""
    for attempt in range(5):
        try:
            os.replace(tmp_path, file_path)
            return
        except PermissionError:
            time.sleep(0.01 * (2 ** attempt))
    with open(file_path, "w", encoding=DEFAULT_ENCODING) as fh:
        fh.write(text)
        fh.flush()
        if has_fsync:
            os.fsync(fh.fileno())


def atomic_write_text(file_path: Path, text: str, has_fsync: bool = False) -> None:
    """Atomically writes text using tmp file with retry backoff for Windows file locks."""
    file_path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = file_path.with_suffix(f".tmp.{os.getpid()}.{threading.get_ident()}")
    try:
        with open(tmp_path, "w", encoding=DEFAULT_ENCODING) as fh:
            fh.write(text)
            fh.flush()
            if has_fsync:
                os.fsync(fh.fileno())
        retry_replace_or_overwrite(tmp_path, file_path, text, has_fsync=has_fsync)
    except OSError as err:
        sys.stderr.write(f"[WARN] Failed atomic write: {err}\n")
    finally:
        clean_temp_file(tmp_path)


def atomic_write_json(file_path: Path, data: Any, has_fsync: bool = False) -> None:
    """Atomically writes JSON structure to file."""
    payload = json.dumps(data, indent=2)
    atomic_write_text(file_path, payload, has_fsync=has_fsync)


def clear_existing_link(latest_dir: Path) -> None:
    """Removes existing junction or directory link cleanly."""
    if latest_dir.exists():
        try:
            os.rmdir(str(latest_dir))
        except OSError:
            shutil.rmtree(latest_dir, ignore_errors=True)


def create_os_junction_or_symlink(target_dir: Path, latest_dir: Path) -> None:
    """Creates NTFS junction on Windows or symlink on POSIX."""
    try:
        if os.name == "nt":
            import _winapi
            _winapi.CreateJunction(str(target_dir), str(latest_dir))
        else:
            latest_dir.symlink_to(target_dir, target_is_directory=True)
    except OSError as err:
        sys.stderr.write(f"[WARN] Failed to create junction or symlink: {err}\n")


def link_latest_session(target_dir: Path, latest_dir: Path) -> None:
    """Links or writes pointer for latest session directory."""
    pointer_file = target_dir.parent.parent / "latest_run.txt"
    try:
        pointer_file.write_text(normalize_repo_rel(target_dir), encoding=DEFAULT_ENCODING)
        clear_existing_link(latest_dir)
        create_os_junction_or_symlink(target_dir, latest_dir)
    except OSError as err:
        sys.stderr.write(f"[WARN] Failed to link latest session: {err}\n")


# ====================================================================================================
#              SMART INCREMENTAL GO TEST ENGINE & CODE-TO-TEST INVENTORY GENERATOR
# ====================================================================================================

FUNC_START_RE = re.compile(r"^func\s+(?:\([^)]+\)\s+)?([A-Za-z0-9_]+)\s*\(")
TEST_START_RE = re.compile(r"^func\s+(Test[A-Za-z0-9_]+)\s*\(")


def extract_go_function_hashes(filepath: Path) -> dict[str, str]:
    """Extracts function declarations and computes SHA256 body hashes from a Go source file."""
    funcs: dict[str, str] = {}
    try:
        with open(filepath, "r", encoding="utf-8", errors="ignore") as fp:
            lines = fp.readlines()
    except OSError:
        return funcs

    curr_name: str | None = None
    curr_lines: list[str] = []
    brace_depth = 0
    in_func = False

    for line in lines:
        if not in_func:
            m = FUNC_START_RE.match(line)
            if m and "{" in line:
                curr_name = m.group(1)
                curr_lines = [line]
                brace_depth = line.count("{") - line.count("}")
                if brace_depth > 0:
                    in_func = True
                else:
                    body = "".join(curr_lines)
                    funcs[curr_name] = hashlib.sha256(body.encode("utf-8")).hexdigest()[:16]
        else:
            curr_lines.append(line)
            brace_depth += line.count("{") - line.count("}")
            if brace_depth <= 0:
                in_func = False
                if curr_name:
                    body = "".join(curr_lines)
                    funcs[curr_name] = hashlib.sha256(body.encode("utf-8")).hexdigest()[:16]

    return funcs


def extract_go_test_functions(filepath: Path) -> dict[str, str]:
    """Extracts Test* functions and computes SHA256 body hashes from a Go test file."""
    tests: dict[str, str] = {}
    try:
        with open(filepath, "r", encoding="utf-8", errors="ignore") as fp:
            lines = fp.readlines()
    except OSError:
        return tests

    curr_name: str | None = None
    curr_lines: list[str] = []
    brace_depth = 0
    in_func = False

    for line in lines:
        if not in_func:
            m = TEST_START_RE.match(line)
            if m and "{" in line:
                curr_name = m.group(1)
                curr_lines = [line]
                brace_depth = line.count("{") - line.count("}")
                if brace_depth > 0:
                    in_func = True
                else:
                    body = "".join(curr_lines)
                    tests[curr_name] = hashlib.sha256(body.encode("utf-8")).hexdigest()[:16]
        else:
            curr_lines.append(line)
            brace_depth += line.count("{") - line.count("}")
            if brace_depth <= 0:
                in_func = False
                if curr_name:
                    body = "".join(curr_lines)
                    tests[curr_name] = hashlib.sha256(body.encode("utf-8")).hexdigest()[:16]

    return tests


FILE_HASH_MEMO_CACHE: dict[str, str] = {}


def compute_file_hash(filepath: Path) -> str:
    """Computes short SHA256 hex digest for an entire file with in-memory memoization."""
    key = str(filepath)
    if key in FILE_HASH_MEMO_CACHE:
        return FILE_HASH_MEMO_CACHE[key]
    try:
        digest = hashlib.sha256(filepath.read_bytes()).hexdigest()[:16]
        FILE_HASH_MEMO_CACHE[key] = digest
        return digest
    except OSError:
        return ""


def load_raw_test_inventory(path: Path) -> dict[str, Any]:
    """Loads existing test inventory JSON from disk if present."""
    for candidate in (path, CACHE_INVENTORY_FILE):
        if candidate.is_file():
            try:
                data = json.loads(candidate.read_text(encoding=DEFAULT_ENCODING))
                if isinstance(data, dict):
                    return data
            except Exception:
                pass
    return {}


def resolve_inventory_delta(repo_root: Path, delta: set[str] | None) -> set[str]:
    """Resolves repo delta set for test inventory short-circuiting."""
    if delta is not None:
        return delta
    prev = load_previous_state(False)
    curr_head = get_head_commit(repo_root)
    curr_dirty = get_dirty_files_map(repo_root)
    return compute_repo_delta(repo_root, prev.get("head", ""), curr_head, prev.get("dirty", {}), curr_dirty)


def is_pkg_file_in_delta(pkg: str, delta: set[str]) -> bool:
    """Checks if any file in the git delta belongs to the package directory."""
    prefix = pkg.rstrip("/") + "/"

    return any(f == pkg or f.startswith(prefix) for f in delta)


def is_test_cached_by_delta(t: dict[str, Any], delta: set[str], force: bool) -> bool:
    """Checks if test can be bypassed directly from git delta and prior success."""
    tgt, tf = t.get("target_file", ""), t.get("test_file", "")
    pkg = t.get("package", "")
    has_passed = t.get("last_status") == "passed"
    is_affected = bool((tgt and tgt in delta) or (tf and tf in delta) or (pkg and is_pkg_file_in_delta(pkg, delta)))

    return bool(not force and has_passed and not is_affected and t.get("code_hash"))


def evaluate_inventory_test_item(t: dict[str, Any], delta: set[str], root: Path, force: bool) -> tuple[bool, bool]:
    """Evaluates if test needs run and whether its dirty status changed."""
    if is_test_cached_by_delta(t, delta, force):
        has_changed = bool(t.get("needs_run") is not False)
        t["needs_run"] = False
        return False, has_changed
    pkg = t.get("package", "")
    is_pkg_affected = bool(pkg and is_pkg_file_in_delta(pkg, delta))
    tgt, tf = t.get("target_file", ""), t.get("test_file", "")
    ch = compute_file_hash(root / tgt) if tgt else ""
    th = compute_file_hash(root / tf) if tf else ""
    is_match = not force and not is_pkg_affected and t.get("last_status") == "passed" and ch == t.get("code_hash") and th == t.get("test_hash") and ch != ""
    has_changed = bool(t.get("needs_run") == is_match)
    t["needs_run"] = not is_match

    return not is_match, has_changed


def process_inventory_tests(tests: dict, delta: set[str], root: Path, force: bool) -> tuple[int, int, bool]:
    """Scans all inventory tests and computes dirty count and status transitions."""
    dirty_count, cached_count = 0, 0
    is_dirty_state = False
    for t in tests.values():
        needs_run, has_changed = evaluate_inventory_test_item(t, delta, root, force)
        if has_changed:
            is_dirty_state = True
        if needs_run:
            dirty_count += 1
        else:
            cached_count += 1
    return dirty_count, cached_count, is_dirty_state


def persist_inventory_if_dirty(inv: dict, dirty: int, cached: int, is_dirty: bool, force: bool) -> None:
    """Writes inventory.json only if dirty or status changed or missing from disk."""
    has_disk_file = TEST_INVENTORY_PATH.is_file() and CACHE_INVENTORY_FILE.is_file()
    if not is_dirty and dirty == 0 and not force and has_disk_file:
        return
    summary = inv.setdefault("summary", {})
    summary["total"] = len(inv.get("tests", {}))
    summary["cached"] = cached
    summary["dirty"] = dirty
    inv["updated_at"] = time.strftime("%Y-%m-%dT%H:%M:%SZ")
    atomic_write_json(TEST_INVENTORY_PATH, inv)
    atomic_write_json(CACHE_INVENTORY_FILE, inv)


def invoke_inventory_generator_fallback(repo_root: Path, slow_threshold: float, force: bool) -> dict[str, Any]:
    """Invokes external test inventory generator script."""
    try:
        sys.path.insert(0, str(Path(__file__).parent))
        inv_gen = importlib.import_module("33-test-inventory-generator")
        return inv_gen.build_test_inventory(repo_root, slow_threshold=slow_threshold, force_run_all=force)
    except Exception as exc:
        sys.stderr.write(f"[WARN] Failed to invoke 33-test-inventory-generator: {exc}\n")
        return {}


def build_or_update_test_inventory(
    repo_root: Path, force: bool = False, repo_delta: set[str] | None = None
) -> dict[str, Any]:
    """Discovers all repository tests, indexes source code functions, and caches hashes & timings."""
    existing_inv = load_raw_test_inventory(TEST_INVENTORY_PATH)
    if not existing_inv or not existing_inv.get("tests") or force:
        slow_threshold = float(os.environ.get("GITMAP_SLOW_TEST_THRESHOLD", "4.0"))
        gen_inv = invoke_inventory_generator_fallback(repo_root, slow_threshold, force)
        if gen_inv:
            return gen_inv
    delta = resolve_inventory_delta(repo_root, repo_delta)
    tests = existing_inv.get("tests", {})
    dirty, cached, is_dirty = process_inventory_tests(tests, delta, repo_root, force)
    persist_inventory_if_dirty(existing_inv, dirty, cached, is_dirty, force)
    return existing_inv


def print_inventory_summary(inventory: dict[str, Any]) -> None:
    """Displays formatted terminal summary of the test inventory manifest."""
    summ = inventory.get("summary", {})
    total = summ.get("total", len(inventory.get("tests", {})))
    cached = summ.get("cached", 0)
    dirty = summ.get("dirty", 0)
    pkgs = summ.get("packages", 0)
    total_time = sum(t.get("duration_sec", 0.0) for t in inventory.get("tests", {}).values())

    print("================================================================")
    print("           GITMAP CI/CD TEST INVENTORY & TIMINGS MANIFEST       ")
    print("================================================================")
    print(f"📁 Test Inventory File     : {TEST_INVENTORY_PATH}")
    print(f"📋 Total Tests Cataloged   : {total}")
    print(f"📦 Unique Test Packages    : {pkgs}")
    print(f"⏱️  Historical Total Time   : {total_time:.2f}s")
    print(f"✅ Cached (Unchanged Code) : {cached}")
    print(f"🔄 Dirty (Needs Execution) : {dirty}")
    print("================================================================", flush=True)


def resolve_package_test_target(pkg: str, repo_root: Path) -> tuple[Path, str]:
    """Resolves module working directory and package argument for go test."""
    if pkg.startswith("scripts/changelog"):
        mod_dir = repo_root / "scripts" / "changelog"
        sub = pkg[len("scripts/changelog"):].lstrip("/")
        return mod_dir, f"./{sub}" if sub else "."
    if pkg.startswith("04-code/golang"):
        mod_dir = repo_root / "04-code" / "golang"
        sub = pkg[len("04-code/golang"):].lstrip("/")
        return mod_dir, f"./{sub}" if sub else "."
    if pkg.startswith("cli-updater"):
        mod_dir = repo_root / "cli-updater"
        sub = pkg[len("cli-updater"):].lstrip("/")
        return mod_dir, f"./{sub}" if sub else "."
    if pkg.startswith("cli/"):
        return repo_root / "cli", "./" + pkg[len("cli/"):]
    if pkg == "cli":
        return repo_root / "cli", "."

    return repo_root / "cli", f"./{pkg}"


def build_standard_test_env() -> dict[str, str]:
    """Constructs hermetic test environment variables for Go test execution."""
    env = dict(os.environ)
    env["GOMAXPROCS"] = str(CPU_CORES)
    env["GITMAP_MOCK_GH"] = "1"
    env["GITMAP_FAST_PROBE"] = "1"
    env["GITMAP_IN_MEMORY_DB"] = "1"
    env["GITMAP_TEST"] = "1"
    env["GOTMPDIR"] = str(REPO_BUILD_TEMP)
    env["TMPDIR"] = str(REPO_TEST_TEMP)
    env["TEMP"] = str(REPO_TEST_TEMP)
    env["TMP"] = str(REPO_TEST_TEMP)
    env["GOCACHE"] = str(REPO_GOCACHE_TEMP)

    return env


def get_test_binaries_dir() -> Path:
    """Returns the OS temp directory scoped by repository name for precompiled test binaries."""
    target = get_repo_os_temp_dir("test_binaries")
    target.mkdir(parents=True, exist_ok=True)

    return target


def get_test_manifest_path() -> Path:
    """Returns the path to manifest.json for tracking precompiled test binaries."""
    return get_test_binaries_dir() / "manifest.json"


def load_test_manifest() -> dict[str, Any]:
    """Loads existing precompiled binary manifest from the OS temp directory."""
    path = get_test_manifest_path()
    if path.exists():
        try:
            return json.loads(path.read_text(encoding="utf-8"))
        except Exception:
            pass

    return {"version": "1.0.0", "packages": {}}


def save_test_manifest(manifest: dict[str, Any]) -> None:
    """Persists precompiled binary manifest tracking to the OS temp directory."""
    manifest["updated_at"] = time.strftime("%Y-%m-%dT%H:%M:%SZ")
    path = get_test_manifest_path()
    try:
        path.write_text(json.dumps(manifest, indent=2), encoding="utf-8")
    except Exception:
        pass


def resolve_package_dir(pkg: str, repo_root: Path) -> Path:
    """Resolves filesystem directory for a Go package within the repository."""
    cwd, rel = resolve_package_test_target(pkg, repo_root)
    if rel == ".":
        return cwd

    return cwd / rel.lstrip("./")


def compute_package_hash(pkg_dir: Path) -> str:
    """Computes composite hash of all Go source and test files in a package directory."""
    hasher = hashlib.sha256()
    if not pkg_dir.exists():
        return ""
    for f in sorted(pkg_dir.glob("*.go")):
        try:
            hasher.update(f.read_bytes())
        except Exception:
            pass

    return hasher.hexdigest()[:16]


def is_precompiled_binary_fresh(
    pkg: str, bin_path_str: str, pkg_dir: Path, manifest: dict[str, Any]
) -> bool:
    """Verifies that precompiled test binary exists and matches the current package hash."""
    bin_path = Path(bin_path_str)
    if not bin_path.exists():
        return False
    entry = manifest.get("packages", {}).get(pkg, {})
    saved_hash = entry.get("code_hash", "")
    current_hash = compute_package_hash(pkg_dir)

    return bool(saved_hash and saved_hash == current_hash)


def partition_safe_batches(pkgs: list[str], max_size: int = 15) -> list[list[str]]:
    """Partitions packages into batches ensuring zero duplicate base names per batch."""
    batches: list[list[str]] = []
    for p in pkgs:
        base = p.split("/")[-1]
        placed = False
        for b in batches:
            b_bases = {x.split("/")[-1] for x in b}
            if len(b) < max_size and base not in b_bases:
                b.append(p)
                placed = True
                break
        if not placed:
            batches.append([p])

    return batches


def build_warmup_batch_targets(batch: list[str], repo_root: Path, batch_dir: Path) -> tuple[Path, list[str]]:
    """Resolves working directory and command arguments for multi-package compilation."""
    cwd = resolve_package_test_target(batch[0], repo_root)[0]
    out_dir_arg = str(batch_dir.resolve()) + os.sep
    rel_targets = [resolve_package_test_target(p, repo_root)[1] for p in batch]

    return cwd, ["go", "test", "-c", "-o", out_dir_arg] + rel_targets


def compile_single_warmup_batch(
    batch: list[str], batch_idx: int, repo_root: Path, test_env: dict[str, str]
) -> tuple[dict[str, str], str]:
    """Compiles a batch of test packages into a dedicated batch directory in OS temp."""
    batch_dir = get_test_binaries_dir() / f"batch_{batch_idx}"
    batch_dir.mkdir(parents=True, exist_ok=True)
    cwd, cmd = build_warmup_batch_targets(batch, repo_root, batch_dir)
    res = subprocess.run(cmd, cwd=cwd, capture_output=True, text=True, env=test_env)
    if res.returncode != 0:
        return {}, res.stderr or res.stdout
    bin_ext = ".test.exe" if os.name == "nt" else ".test"
    found = {p: str(batch_dir / f"{p.split('/')[-1]}{bin_ext}") for p in batch}

    return found, ""


def filter_packages_needing_compile(
    pkgs: list[str], manifest: dict[str, Any], repo_root: Path, force: bool
) -> tuple[dict[str, str], list[str]]:
    """Separates warm cached packages from packages that require binary compilation."""
    bin_map: dict[str, str] = {}
    needs_compile: list[str] = []
    for p in pkgs:
        p_dir = resolve_package_dir(p, repo_root)
        entry = manifest.get("packages", {}).get(p, {})
        bin_path_str = entry.get("binary_path", "")
        if not force and is_precompiled_binary_fresh(p, bin_path_str, p_dir, manifest):
            bin_map[p] = bin_path_str
        else:
            needs_compile.append(p)

    return bin_map, needs_compile


def update_manifest_with_compiled_batch(
    batch_map: dict[str, str], manifest: dict[str, Any], repo_root: Path
) -> None:
    """Updates manifest metadata entries with newly compiled test binaries."""
    for p, bin_path in batch_map.items():
        p_dir = resolve_package_dir(p, repo_root)
        manifest.setdefault("packages", {})[p] = {
            "binary_path": bin_path,
            "code_hash": compute_package_hash(p_dir),
            "built_at": time.strftime("%Y-%m-%dT%H:%M:%SZ"),
        }


def cluster_packages_by_speed(pkgs: list[str], manifest: dict[str, Any]) -> list[str]:
    """Sorts packages clustering slow packages together and fast packages together."""
    pkg_entries = manifest.get("packages", {})

    return sorted(pkgs, key=lambda p: float(pkg_entries.get(p, {}).get("duration_estimate_sec", 0.0)), reverse=True)


def resolve_module_key(pkg: str) -> str:
    """Returns module grouping key for a Go package."""
    for prefix in ("scripts/changelog", "04-code/golang", "cli-updater"):
        if pkg.startswith(prefix):
            return prefix

    return "cli"


def group_packages_by_module(pkgs: list[str]) -> dict[str, list[str]]:
    """Groups packages by root Go module to prevent cross-module compilation."""
    modules: dict[str, list[str]] = {}
    for p in pkgs:
        modules.setdefault(resolve_module_key(p), []).append(p)

    return modules


def build_module_safe_batches(needs_compile: list[str], manifest: dict[str, Any]) -> list[list[str]]:
    """Partitions sorted needs into module-isolated safe batches."""
    sorted_needs = cluster_packages_by_speed(needs_compile, manifest)
    batches: list[list[str]] = []
    for mod_pkgs in group_packages_by_module(sorted_needs).values():
        batches.extend(partition_safe_batches(mod_pkgs, max_size=15))

    return batches


def process_compiled_warmup_batch(
    batch_map: dict[str, str], bin_map: dict[str, str], manifest: dict[str, Any], repo_root: Path
) -> None:
    """Updates binary mapping and registers package entries into manifest."""
    bin_map.update(batch_map)
    update_manifest_with_compiled_batch(batch_map, manifest, repo_root)


def run_warmup_batches_compilation(
    needs_compile: list[str], bin_map: dict[str, str],
    manifest: dict[str, Any], repo_root: Path
) -> str:
    """Compiles batches and registers generated binaries in manifest."""
    test_env = build_standard_test_env()
    batches = build_module_safe_batches(needs_compile, manifest)
    for idx, batch in enumerate(batches):
        batch_map, err = compile_single_warmup_batch(batch, idx, repo_root, test_env)
        if err:
            return f"Warmup compilation failed for batch {idx}: {err}"
        process_compiled_warmup_batch(batch_map, bin_map, manifest, repo_root)

    return ""





def precompile_test_binaries_warmup(
    pkgs: list[str], repo_root: Path, force: bool
) -> tuple[dict[str, str], str]:
    """Pre-compiles test packages in batched sweeps into OS temp with caching."""
    manifest = load_test_manifest()
    bin_map, needs_compile = filter_packages_needing_compile(pkgs, manifest, repo_root, force)
    if not needs_compile:
        return bin_map, ""
    err = run_warmup_batches_compilation(needs_compile, bin_map, manifest, repo_root)
    if err:
        return bin_map, err
    save_test_manifest(manifest)

    return bin_map, ""


def record_package_test_failure(pkg: str, test_id: str, content: str) -> None:
    """Writes surgical failure log for a failed test case into failures directory."""
    clean_tid = re.sub(r'[<>:"/\\|?*]', "_", test_id)
    fail_log = FAILURES_DIR / f"{clean_tid}.log"
    try:
        fail_log.write_text(content, encoding="utf-8")
    except Exception:
        pass


def extract_failing_test_names(output: str) -> list[str]:
    """Extracts failing test function names from native Go test binary output."""
    matches = re.findall(r"--- FAIL: ([^\s\(]+)", output)

    return list(dict.fromkeys(matches))


def process_precompiled_worker_result(
    pkg: str, proc: subprocess.CompletedProcess[str], pkg_tests: list[dict[str, Any]]
) -> tuple[int, int, str, dict[str, dict[str, Any]]]:
    """Evaluates process exit code; silent on pass, surgical error extraction on failure."""
    if proc.returncode == 0:
        test_results = {t["id"]: {"status": "passed", "elapsed": 0.01} for t in pkg_tests}
        return len(pkg_tests), 0, "", test_results
    out_text = (proc.stdout or "") + "\n" + (proc.stderr or "")
    failing_names = extract_failing_test_names(out_text)
    failed_count = len(failing_names) if failing_names else len(pkg_tests)
    passed_count = max(0, len(pkg_tests) - failed_count)
    record_package_test_failure(pkg, f"{pkg}.failure", out_text)
    test_results = {}
    for t in pkg_tests:
        has_failed = t["test_func"] in failing_names or not failing_names
        test_results[t["id"]] = {"status": "failed" if has_failed else "passed", "elapsed": 0.01}

    return passed_count, failed_count, out_text.strip(), test_results


def execute_test_binary_subprocess(cmd: list[str], pkg_dir: Path, timeout_sec: int) -> subprocess.CompletedProcess[str]:
    """Executes a precompiled test binary with DEVNULL stdin and hermetic env."""
    effective_timeout = max(180, timeout_sec)

    return subprocess.run(
        cmd, cwd=pkg_dir, capture_output=True, text=True,
        encoding="utf-8", errors="replace", timeout=effective_timeout,
        env=build_standard_test_env(), stdin=subprocess.DEVNULL
    )


def run_precompiled_package_worker(
    pkg: str, bin_path_str: str, pkg_tests: list[dict[str, Any]],
    repo_root: Path, timeout_sec: int
) -> tuple[int, int, str, dict[str, dict[str, Any]]]:
    """Worker executing precompiled package test binary with 8-way parallel goroutines."""
    bin_path = Path(bin_path_str)
    if not bin_path.exists():
        return 0, len(pkg_tests), f"Binary not found: {bin_path}", {}
    try:
        proc = execute_test_binary_subprocess([str(bin_path), "-test.parallel=8"], resolve_package_dir(pkg, repo_root), timeout_sec)
        return process_precompiled_worker_result(pkg, proc, pkg_tests)
    except Exception as ex:
        return 0, len(pkg_tests), f"Execution error in {pkg}: {ex}", {}



def run_git_stdout_lines(args: list[str], repo_root: Path) -> list[str]:
    """Runs a git command and returns non-empty stripped stdout lines."""
    git_bin = shutil.which("git") or "git"
    try:
        res = subprocess.run([git_bin] + args, cwd=str(repo_root), capture_output=True, text=True)
        if res.returncode == 0:
            return [line.strip() for line in res.stdout.splitlines() if line.strip()]
    except Exception:
        pass

    return []


def get_current_git_head_sha(repo_root: Path) -> str:
    """Returns current HEAD commit hash from git."""
    lines = run_git_stdout_lines(["rev-parse", "HEAD"], repo_root)

    return lines[0].strip() if lines else ""


def parse_porcelain_path(line: str) -> str:
    """Extracts normalized relative path from git porcelain status line."""
    p = line[2:].strip().strip('"')
    if " -> " in p:
        p = p.split(" -> ")[-1].strip()

    return p.replace("\\", "/")


def get_repo_changed_files_since_sha(last_sha: str, current_head: str, repo_root: Path) -> set[str]:
    """Returns set of relative paths changed uncommitted or committed since last_sha."""
    changed: set[str] = set()
    st_lines = run_git_stdout_lines(["status", "--porcelain", "-uall"], repo_root)
    changed.update(parse_porcelain_path(line) for line in st_lines if parse_porcelain_path(line))
    if last_sha and last_sha != current_head:
        df_lines = run_git_stdout_lines(["diff", "--name-only", f"{last_sha}..{current_head}"], repo_root)
        changed.update(line.replace("\\", "/") for line in df_lines)

    return changed


def is_pkg_changed_in_git(rel_pkg: str, changed_files: set[str]) -> bool:
    """Checks if any changed file resides inside the package directory."""
    prefix = rel_pkg.rstrip("/") + "/"

    return any(f == rel_pkg or f.startswith(prefix) for f in changed_files)


def is_package_execution_fresh(
    pkg: str, bin_path_str: str, pkg_dir: Path, manifest: dict[str, Any],
    changed_files: set[str], repo_root: Path
) -> bool:
    """Checks if test binary exists, last status passed, and no package files changed."""
    entry = manifest.get("packages", {}).get(pkg, {})
    if entry.get("last_status") != "passed":
        return False
    if not is_precompiled_binary_fresh(pkg, bin_path_str, pkg_dir, manifest):
        return False
    rel_pkg = pkg_dir.relative_to(repo_root).as_posix()

    return not is_pkg_changed_in_git(rel_pkg, changed_files)


def record_package_test_success(
    pkg: str, bin_path: str, manifest: dict[str, Any],
    repo_root: Path, current_head: str, test_count: int
) -> None:
    """Updates manifest entry with successful test execution metadata."""
    p_dir = resolve_package_dir(pkg, repo_root)
    manifest.setdefault("packages", {})[pkg] = {
        "binary_path": bin_path,
        "code_hash": compute_package_hash(p_dir),
        "last_git_sha": current_head,
        "last_status": "passed",
        "test_count": test_count,
        "built_at": time.strftime("%Y-%m-%dT%H:%M:%SZ"),
        "last_run_at": time.strftime("%Y-%m-%dT%H:%M:%SZ"),
    }


def record_package_test_failure_status(
    pkg: str, bin_path: str, manifest: dict[str, Any],
    repo_root: Path, current_head: str, test_count: int
) -> None:
    """Updates manifest entry with failed test execution metadata."""
    p_dir = resolve_package_dir(pkg, repo_root)
    manifest.setdefault("packages", {})[pkg] = {
        "binary_path": bin_path,
        "code_hash": compute_package_hash(p_dir),
        "last_git_sha": current_head,
        "last_status": "failed",
        "test_count": test_count,
        "built_at": time.strftime("%Y-%m-%dT%H:%M:%SZ"),
        "last_run_at": time.strftime("%Y-%m-%dT%H:%M:%SZ"),
    }


def should_skip_package_run(
    pkg: str, bin_path: str, manifest: dict[str, Any],
    changed_files: set[str], repo_root: Path, force: bool
) -> bool:
    """Evaluates if package is fresh and eligible to skip execution."""
    p_dir = resolve_package_dir(pkg, repo_root)

    return bool(not force and is_package_execution_fresh(pkg, bin_path, p_dir, manifest, changed_files, repo_root))


def separate_execution_packages(
    dirty_pkg_map: dict[str, list[dict[str, Any]]], bin_map: dict[str, str],
    manifest: dict[str, Any], changed_files: set[str], repo_root: Path, force: bool
) -> tuple[dict[str, list[dict[str, Any]]], dict[str, list[dict[str, Any]]]]:
    """Separates packages into cached done-checking packages and packages needing execution."""
    cached_pkgs: dict[str, list[dict[str, Any]]] = {}
    exec_pkgs: dict[str, list[dict[str, Any]]] = {}
    for pkg, b_tests in dirty_pkg_map.items():
        bin_path = bin_map.get(pkg, "") or manifest.get("packages", {}).get(pkg, {}).get("binary_path", "")
        if should_skip_package_run(pkg, bin_path, manifest, changed_files, repo_root, force):
            cached_pkgs[pkg] = b_tests
        else:
            exec_pkgs[pkg] = b_tests

    return cached_pkgs, exec_pkgs



def apply_cached_package_results(
    cached_pkgs: dict[str, list[dict[str, Any]]], tests: dict[str, Any]
) -> int:
    """Marks tests from cached packages as passed in inventory and returns count."""
    passed = 0
    for b_tests in cached_pkgs.values():
        passed += len(b_tests)
        for t in b_tests:
            tid = t.get("id", "")
            if tid in tests:
                tests[tid]["last_status"] = "passed"
                tests[tid]["needs_run"] = False

    return passed


def filter_tests_by_package_or_file(tests: dict[str, Any], queries: list[str], repo_root: Path) -> list[dict[str, Any]]:

    """Filters inventory tests based on code file paths, code file names, or Go package names."""
    matched_ids: set[str] = set()
    for q_raw in queries:
        q = q_raw.strip().replace("\\", "/").rstrip("/")
        if not q:
            continue
        q_base = os.path.basename(q)
        initial_count = len(matched_ids)
        for tid, t in tests.items():
            pkg = t.get("package", "").replace("\\", "/")
            tf = t.get("target_file", "").replace("\\", "/")
            test_f = t.get("test_file", "").replace("\\", "/")

            # 1. Package match (exact, relative, or Go module suffix)
            if pkg == q or pkg == f"cli/{q}" or pkg.endswith("/" + q) or pkg.split("/")[-1] == q:
                matched_ids.add(tid)
                continue
            if q.startswith("github.com/") and q.endswith(pkg):
                matched_ids.add(tid)
                continue

            # 2. File path match
            if tf == q or test_f == q or tf.endswith("/" + q) or test_f.endswith("/" + q):
                matched_ids.add(tid)
                continue

            # 3. File name match (basename)
            if os.path.basename(tf) == q_base or os.path.basename(test_f) == q_base:
                matched_ids.add(tid)
                continue

        # 4. Fallback: if no specific test targeted this .go file, find its package
        if len(matched_ids) == initial_count and q_base.endswith(".go"):
            found_pkg = None
            for t in tests.values():
                tf = t.get("target_file", "").replace("\\", "/")
                if os.path.basename(tf) == q_base:
                    found_pkg = t.get("package")
                    break
            if not found_pkg:
                for root, _, files in os.walk(repo_root / "cli"):
                    if q_base in files:
                        found_pkg = os.path.relpath(root, repo_root).replace("\\", "/")
                        break
            if found_pkg:
                for tid, t in tests.items():
                    if t.get("package") == found_pkg:
                        matched_ids.add(tid)

    return [tests[tid] for tid in matched_ids if tid in tests]


def filter_fast_mode_tests(tests: list[dict[str, Any]], is_fast: bool) -> list[dict[str, Any]]:
    """Filters dirty tests to hot and warm tiers when fast mode is active."""
    if not is_fast:
        return tests
    return [
        t for t in tests
        if t.get("heat_tier") in ("hot", "warm") or t.get("last_status") == "failed"
    ]


def sort_tests_by_heat(tests: list[dict[str, Any]]) -> None:
    """Sorts test list in-place by descending heat score for fail-fast execution."""
    tests.sort(key=lambda t: float(t.get("heat_score", 0.0)), reverse=True)


def is_pkg_code_changed(pkg: str, pkg_dir: Path, manifest: dict[str, Any]) -> bool:
    """Checks if package source files have changed since last recorded test run."""
    entry = manifest.get("packages", {}).get(pkg, {})
    if not entry or entry.get("last_status") != "passed":
        return True
    bin_path = entry.get("binary_path", "")
    if not bin_path or not Path(bin_path).exists():
        return True
    saved_hash = entry.get("code_hash", "")
    current_hash = compute_package_hash(pkg_dir)

    return bool(not saved_hash or saved_hash != current_hash)


def resolve_dirty_go_tests(
    tests: dict[str, Any], changed_files: set[str], manifest: dict[str, Any],
    repo_root: Path, include_heavy: bool, force: bool
) -> list[dict[str, Any]]:
    """Filters inventory tests to only packages changed in git or failing in manifest."""
    dirty = []
    checked_pkgs: dict[str, bool] = {}
    for t in tests.values():
        pkg = t.get("package", "")
        if pkg not in checked_pkgs:
            p_dir = resolve_package_dir(pkg, repo_root)
            is_git_diff = is_pkg_changed_in_git(pkg, changed_files)
            is_code_diff = is_pkg_code_changed(pkg, p_dir, manifest)
            checked_pkgs[pkg] = bool(force or (is_git_diff and is_code_diff))
        is_light = bool(include_heavy or "tests/heavy_test" not in t.get("test_file", ""))
        if checked_pkgs[pkg] and is_light:
            dirty.append(t)

    return dirty


def build_cached_smart_test_result(
    name: str, total_tests: int, current_head: str, elapsed: float
) -> JobResult:
    """Builds clean JobResult for zero-change cached pass."""
    out_msg = f"Passed {total_tests} tests in {format_duration(elapsed)} (all packages cached on git {current_head[:8]})"
    print(f"\n⚡ [Quad Runner] All Go test packages verified & cached (done checking on git {current_head[:8]}).\n", flush=True)

    return JobResult(
        name=name, cmd=["go", "test", "smart-incremental"], code=0,
        out=out_msg, err="", elapsed=elapsed, is_cached=True
    )


def run_smart_go_tests(
    name: str, timeout_sec: int, max_workers: int, force: bool, repo_root: Path,
    tel: TelemetryTracker | None = None, package_filter: list[str] | str | None = None
) -> JobResult:
    """Executes changed Go tests with unified priority worker pool."""
    start_time = time.monotonic()
    clear_repo_test_temp()
    clear_stale_failures_log()
    manifest = load_test_manifest()
    current_head = get_current_git_head_sha(repo_root)
    last_head = manifest.get("last_git_hash", "")
    changed_files = get_repo_changed_files_since_sha(last_head, current_head, repo_root)
    inventory = build_or_update_test_inventory(repo_root, force=force)
    all_tests = inventory.get("tests", {})
    tests = {tid: t for tid, t in all_tests.items() if t.get("test_file", "").endswith(".go")}
    if package_filter:
        queries = [package_filter] if isinstance(package_filter, str) else list(package_filter)
        target_tests = filter_tests_by_package_or_file(tests, queries, repo_root)
        if not target_tests:
            elapsed = round(time.monotonic() - start_time, 2)
            out_msg = f"[WARN] No unit tests found matching package/file query: {', '.join(queries)}"
            return JobResult(
                name=name, cmd=["go", "test", f"--pkg={queries}"], code=0,
                out=out_msg, err="", elapsed=elapsed, is_cached=False
            )
        dirty_tests = target_tests
    else:
        include_heavy = os.environ.get("GITMAP_RUN_HEAVY_TESTS") == "1"
        dirty_tests = resolve_dirty_go_tests(tests, changed_files, manifest, repo_root, include_heavy, force)

    is_fast_mode = os.environ.get("GITMAP_FAST_TESTS") == "1"
    dirty_tests = filter_fast_mode_tests(dirty_tests, is_fast_mode)

    if not dirty_tests:
        elapsed = round(time.monotonic() - start_time, 2)
        manifest["last_git_hash"] = current_head
        save_test_manifest(manifest)
        return build_cached_smart_test_result(name, len(tests), current_head, elapsed)

    slow_threshold = float(os.environ.get("GITMAP_SLOW_TEST_THRESHOLD", "4.0"))
    slow_tests = [
        t for t in dirty_tests
        if t.get("tier") in ("slow", "heavy")
        or t.get("is_slow", False)
        or float(t.get("duration_sec", 0.0)) >= slow_threshold
    ]
    fast_tests = [t for t in dirty_tests if t not in slow_tests]
    sort_tests_by_heat(fast_tests)
    sort_tests_by_heat(slow_tests)

    total_dirty = len(dirty_tests)
    passed_count = 0
    failed_count = 0
    error_outputs: list[str] = []

    # Calculate exact ETA from test inventory durations
    slow_dur = sum(float(t.get("duration_sec", 4.0)) for t in slow_tests)
    fast_dur = sum(float(t.get("duration_sec", 0.005)) for t in fast_tests)
    slow_eta = slow_dur / 8.0   # 4 workers * 2 tests
    fast_eta = fast_dur / 16.0  # 4 workers * 4 tests
    total_test_eta = max(1.0, round(slow_eta + fast_eta, 1))

    # Initialize ETA telemetry file for AI agent wait protocol
    RUNNER_ETA_FILE.write_text(json.dumps({
        "status": "running",
        "total_eta_sec": total_test_eta,
        "start_time": start_time,
        "elapsed_sec": 0.0,
        "remaining_eta_sec": total_test_eta,
        "slow_tests_total": len(slow_tests),
        "fast_tests_total": len(fast_tests),
        "completed": 0,
        "passed": 0,
        "failed": 0,
    }, indent=2), encoding="utf-8")

    free_pct, dynamic_workers = detect_cpu_freeness_and_workers(min_workers=4, max_cap=min(8, CPU_CORES))
    unified_worker_limit = dynamic_workers


    dirty_pkg_map: dict[str, list[dict[str, Any]]] = {}
    for t in dirty_tests:
        dirty_pkg_map.setdefault(t["package"], []).append(t)

    # 1. Package-level Git SHA & Done-Checking Cache Filter (BEFORE compilation)
    manifest = load_test_manifest()
    current_head = get_current_git_head_sha(repo_root)
    last_head = manifest.get("last_git_hash", "")
    changed_files = get_repo_changed_files_since_sha(last_head, current_head, repo_root)
    cached_pkg_map, active_pkg_map = separate_execution_packages(
        dirty_pkg_map, {}, manifest, changed_files, repo_root, force
    )
    passed_count += apply_cached_package_results(cached_pkg_map, tests)

    if not active_pkg_map:
        elapsed = round(time.monotonic() - start_time, 2)
        print(f"⚡ [Quad Runner] All {len(cached_pkg_map)} packages verified & cached (done checking on git {current_head[:8]}).\n", flush=True)
        manifest["last_git_hash"] = current_head
        save_test_manifest(manifest)
        inventory["summary"]["dirty"] = 0
        inventory["summary"]["cached"] = len(tests)
        atomic_write_json(TEST_INVENTORY_PATH, inventory)
        atomic_write_json(TEST_INVENTORY_CACHE_PATH, inventory)
        clear_repo_test_temp()
        out_msg = f"Passed {passed_count} tests in {format_duration(elapsed)} (all {len(cached_pkg_map)} packages done checking / cached)"
        return JobResult(
            name=name, cmd=["go", "test", "smart-incremental"], code=0,
            out=out_msg, err="", elapsed=elapsed, is_cached=True
        )

    # 2. Warmup Pre-compilation Phase (ONLY for changed packages)
    print(f"\n⚡ [Warmup Engine] Pre-compiling {len(active_pkg_map)} changed test packages into OS temp...", flush=True)
    t_warm = time.monotonic()
    bin_map, warm_err = precompile_test_binaries_warmup(list(active_pkg_map.keys()), repo_root, force=force)
    warm_sec = round(time.monotonic() - t_warm, 2)
    print(f"   Warmup complete in {format_duration(warm_sec)}.\n", flush=True)
    if warm_err:
        return JobResult(
            name=name, cmd=["go", "test", "warmup"], code=1,
            out="", err=warm_err, elapsed=warm_sec
        )

    # 3. Quad-Process Execution Phase (4 workers, 8 goroutines each)
    worker_limit = min(4, len(active_pkg_map)) if active_pkg_map else 1
    cached_info = f" [{len(cached_pkg_map)} packages done checking / cached]" if cached_pkg_map else ""
    print(f"⚡ [Quad Runner] Executing {len(active_pkg_map)} packages across {worker_limit} concurrent processes (8 goroutines each){cached_info}...\n", flush=True)

    with ThreadPoolExecutor(max_workers=worker_limit) as executor:
        futures = {
            executor.submit(run_precompiled_package_worker, pkg, bin_map.get(pkg, ""), b_tests, repo_root, timeout_sec): (pkg, b_tests)
            for pkg, b_tests in active_pkg_map.items()
        }

        for fut in as_completed(futures):
            pkg, b_tests = futures[fut]
            try:
                pkg_passed, pkg_failed, pkg_out, test_results = fut.result()
                passed_count += pkg_passed
                failed_count += pkg_failed
                if pkg_failed > 0:
                    error_outputs.append(f"[{pkg}] {pkg_out}")
                    record_package_test_failure_status(pkg, bin_map.get(pkg, ""), manifest, repo_root, current_head, len(b_tests))
                else:
                    record_package_test_success(pkg, bin_map.get(pkg, ""), manifest, repo_root, current_head, len(b_tests))
                for tid, res_info in test_results.items():
                    if tid in tests:
                        tests[tid]["duration_sec"] = res_info["elapsed"]
                        tests[tid]["last_status"] = res_info["status"]
                        tests[tid]["last_run_at"] = time.strftime("%Y-%m-%dT%H:%M:%S")
                        tests[tid]["needs_run"] = (res_info["status"] != "passed")
                        tests[tid]["code_hash"] = compute_file_hash(repo_root / tests[tid].get("target_file", ""))
                        tests[tid]["test_hash"] = compute_file_hash(repo_root / tests[tid].get("test_file", ""))
                        if res_info["elapsed"] >= slow_threshold:
                            tests[tid]["is_slow"] = True
                            tests[tid]["tier"] = "slow"
                        record_job_timing(f"GoTest:{tid}", res_info["elapsed"])
            except Exception as ex:
                failed_count += len(b_tests)
                error_outputs.append(f"[{pkg}] Test worker exception: {ex}")

            # Update live telemetry after each completed package/batch
            cur_elapsed = round(time.monotonic() - start_time, 1)
            rem = max(1.0, round(total_test_eta - cur_elapsed, 1))
            RUNNER_ETA_FILE.write_text(json.dumps({
                "status": "running",
                "total_eta_sec": total_test_eta,
                "elapsed_sec": cur_elapsed,
                "remaining_eta_sec": rem,
                "slow_tests_total": len(slow_tests),
                "fast_tests_total": len(fast_tests),
                "completed": passed_count + failed_count,
                "passed": passed_count,
                "failed": failed_count,
            }, indent=2), encoding="utf-8")

    manifest["last_git_hash"] = current_head
    save_test_manifest(manifest)


    elapsed = round(time.monotonic() - start_time, 2)
    RUNNER_ETA_FILE.write_text(json.dumps({
        "status": "completed",
        "total_eta_sec": total_test_eta,
        "elapsed_sec": elapsed,
        "remaining_eta_sec": 0.0,
        "slow_tests_total": len(slow_tests),
        "fast_tests_total": len(fast_tests),
        "completed": passed_count + failed_count,
        "passed": passed_count,
        "failed": failed_count,
    }, indent=2), encoding="utf-8")

    inventory["summary"]["dirty"] = failed_count
    inventory["summary"]["cached"] = len(tests) - failed_count
    atomic_write_json(TEST_INVENTORY_PATH, inventory)
    atomic_write_json(TEST_INVENTORY_CACHE_PATH, inventory)

    clear_repo_test_temp()
    if failed_count > 0:
        err_text = "\n".join(error_outputs)
        return JobResult(
            name=name, cmd=["go", "test", "smart-incremental"], code=1,
            out="", err=err_text, elapsed=elapsed
        )

    out_msg = f"Passed {passed_count} tests ({len(slow_tests)} slow, {len(fast_tests)} fast) in {format_duration(elapsed)} ({len(tests) - total_dirty} tests cached)"
    return JobResult(
        name=name, cmd=["go", "test", "smart-incremental"], code=0,
        out=out_msg, err="", elapsed=elapsed
    )


def resolve_command_binary(cmd: list[str]) -> list[str]:
    """Resolves executable path using shutil.which if needed, with GOPATH/bin and goreleaser fallbacks."""
    resolved = list(cmd)
    if not resolved:
        return resolved

    binary_name = cmd[0]
    binary_path = shutil.which(binary_name)
    if binary_path is not None:
        resolved[0] = binary_path
        return resolved

    # Fallback to GOPATH/bin for Go tools like goreleaser
    gopath = os.environ.get("GOPATH", "d:/go-path")
    ext = ".exe" if os.name == "nt" else ""
    candidate = Path(gopath) / "bin" / f"{binary_name}{ext}"
    if candidate.is_file():
        resolved[0] = str(candidate)
        return resolved

    # Fallback for goreleaser if not installed locally
    if binary_name == "goreleaser":
        return ["go", "run", "github.com/goreleaser/goreleaser/v2@latest"] + resolved[1:]

    return resolved


def build_timeout_result(
    name: str, cmd: list[str], exc: subprocess.TimeoutExpired, elapsed: float, cwd: str | None, env: dict | None
) -> JobResult:
    """Creates JobResult representing a timeout failure."""
    out = exc.stdout if isinstance(exc.stdout, str) else ""
    record_job_timing(name, elapsed)

    return JobResult(
        name=name, cmd=cmd, code="timeout", out=out, err=f"Job timed out after {exc.timeout}s",
        elapsed=elapsed, cwd=cwd, env_overrides=env,
    )


def build_error_result(
    name: str, cmd: list[str], exc: Exception, elapsed: float, cwd: str | None, env: dict | None
) -> JobResult:
    """Creates JobResult representing a subprocess error failure."""
    record_job_timing(name, elapsed)

    return JobResult(name=name, cmd=cmd, code=1, out="", err=str(exc), elapsed=elapsed, cwd=cwd, env_overrides=env)


def execute_subprocess(
    cmd: list[str], timeout_sec: int, env: dict[str, str] | None, cwd: str | None
) -> subprocess.CompletedProcess:
    """Executes subprocess synchronously with standard options."""
    resolved = resolve_command_binary(cmd)
    sub_env = dict(os.environ) if env is None else dict(env)
    sub_env.setdefault("PYTHONUNBUFFERED", "1")
    sub_env["GOMAXPROCS"] = str(CPU_CORES)
    sub_env["GOTMPDIR"] = str(REPO_BUILD_TEMP)
    sub_env["TMPDIR"] = str(REPO_TEST_TEMP)
    sub_env["TEMP"] = str(REPO_TEST_TEMP)
    sub_env["TMP"] = str(REPO_TEST_TEMP)
    sub_env["GOCACHE"] = str(REPO_GOCACHE_TEMP)
    res = subprocess.run(
        resolved, capture_output=True, text=True, encoding=DEFAULT_ENCODING,
        errors="replace", timeout=timeout_sec, env=sub_env, cwd=cwd,
    )

    return res


def build_success_result(
    name: str, cmd: list[str], res: subprocess.CompletedProcess, elapsed: float, cwd: str | None, env: dict | None
) -> JobResult:
    """Creates JobResult representing a completed subprocess."""
    record_job_timing(name, elapsed)

    return JobResult(
        name=name, cmd=cmd, code=res.returncode, out=res.stdout, err=res.stderr,
        elapsed=elapsed, cwd=cwd, env_overrides=env,
    )


def prepare_web_app_gate(name: str, cmd: list[str]) -> JobResult | None:
    """Prepares build artifacts for web app and packaging gates."""
    clear_repo_build_temp()
    for dist_dir in (REPO_ROOT / "dist", REPO_ROOT / "cli" / "dist"):
        robust_rmtree(dist_dir)
    if name == "Web App Build" and not (REPO_ROOT / "node_modules").is_dir():
        return JobResult(name=name, cmd=cmd, code=0, out="[SKIP] node_modules not installed.", err="", elapsed=0.01, is_cached=True)

    return None


def prepare_compile_or_build_gate(name: str, cmd: list[str]) -> JobResult | None:
    """Prepares directory artifacts for compile and build gates."""
    if name == "Go Compile Gate":
        clear_repo_build_temp()
        try:
            (REPO_ROOT / "bin" / "gitmap.exe").unlink(missing_ok=True)
        except OSError:
            pass
    elif name in ("Web App Build", "GoReleaser Snapshot Build"):
        return prepare_web_app_gate(name, cmd)

    return None


def ensure_smoke_binary_present(name: str) -> None:
    """Ensures bin/gitmap.exe is compiled before executing smoke test suites."""
    if name in ("E2E Smoke Suite", "History Purge Smoke", "History Pin Smoke"):
        bin_exe = REPO_ROOT / "bin" / "gitmap.exe"
        if not bin_exe.exists():
            subprocess.run(["go", "build", "-C", "cli", "-o", "../bin/gitmap.exe", "."], check=False)


def execute_job_and_record(cmd: list[str], t_sec: int, env: dict | None, cwd: str | None, name: str) -> JobResult:
    """Invokes subprocess and returns constructed JobResult."""
    start = time.monotonic()
    try:
        res = execute_subprocess(cmd, t_sec, env, cwd)
        return build_success_result(name, cmd, res, round(time.monotonic() - start, 2), cwd, env)
    except subprocess.TimeoutExpired as exc:
        return build_timeout_result(name, cmd, exc, round(time.monotonic() - start, 2), cwd, env)
    except Exception as exc:
        return build_error_result(name, cmd, exc, round(time.monotonic() - start, 2), cwd, env)


def run_job(
    name: str, cmd: list[str], timeout_sec: int, env: dict[str, str] | None = None, cwd: str | None = None
) -> JobResult:
    """Executes a single gate subprocess and records duration, return code, and streams."""
    early_res = prepare_compile_or_build_gate(name, cmd)
    if early_res is not None:
        return early_res
    ensure_smoke_binary_present(name)

    return execute_job_and_record(cmd, timeout_sec, env, cwd, name)


def extract_stack_or_error(res: JobResult) -> str:
    """Combines and returns non-empty stdout/stderr from a failed job."""
    out = (res.out or "").strip()
    err = (res.err or "").strip()
    if out and err:
        return f"{out}\n\n{err}"

    return err or out or "No output captured."


PATH_EXTS = "go|py|ts|tsx|js|jsx|json|sh|ps1|yml|yaml|md|sql|toml|mod|sum"
PATH_PATTERNS = [
    re.compile(r'(?:[a-zA-Z]:[\\/])?[a-zA-Z0-9_\-./\\]+\.(?:' + PATH_EXTS + r')(?::\d+(?::\d+)?)?'),
    re.compile(r'[a-zA-Z0-9_\-./\\]+\.(?:ts|tsx|js|jsx)\s*\(\d+(?:,\d+)?\)'),
    re.compile(r'File\s+"([^"]+\.(?:py|sh|ps1))",\s+line\s+(\d+)'),
]
IGNORED_PREFIXES = ("http://", "https://", "node_modules/", "vendor/", ".git/", ".tmp/", "go/pkg/mod/", "AppData/", "site-packages/")


def normalize_suspect_coord(raw: str) -> str:
    """Converts parens syntax to standard colon syntax."""
    clean = raw.strip().strip('"').strip("'").replace("\\", "/")
    paren_match = re.match(r'^(.*?)\s*\((\d+)(?:,(\d+))?\)$', clean)
    if paren_match:
        base, line, col = paren_match.groups()
        coord = f":{col}" if col else ""

        return f"{base}:{line}{coord}"

    return clean


def resolve_suspect_path(clean: str, cwd: str | None, root: Path) -> str:
    """Resolves suspect path against cwd and repo root."""
    path_no_coord = re.sub(r':\d+(?::\d+)?$', '', clean)
    root_str = str(root).replace("\\", "/") + "/"
    if clean.startswith(root_str):
        clean = clean[len(root_str):]
        path_no_coord = path_no_coord[len(root_str):]
    if cwd and not clean.startswith(f"{cwd}/"):
        if (root / cwd / path_no_coord).is_file():
            return f"{cwd}/{clean}"
    if (root / path_no_coord).is_file():
        return clean

    return clean if "." in clean else ""


def extract_regex_suspect_files(text: str, cwd: str | None, root: Path) -> list[str]:
    """Extracts raw candidate paths using regex patterns."""
    seen: set[str] = set()
    files: list[str] = []
    for pat in PATH_PATTERNS:
        for match in pat.findall(text):
            raw = match if isinstance(match, str) else f"{match[0]}:{match[1]}"
            if any(raw.startswith(p) or f"/{p}" in raw for p in IGNORED_PREFIXES):
                continue
            resolved = resolve_suspect_path(normalize_suspect_coord(raw), cwd, root)
            if resolved and resolved not in seen:
                seen.add(resolved)
                files.append(resolved)

    return files[:10]


def extract_suspect_fallback(gate_name: str, delta: set[str]) -> list[str]:
    """Fallback to gate relevant delta or tool scripts when error output has no files."""
    spec = GATE_SPECS.get(gate_name)
    if not spec:
        return []
    rel_delta = [p for p in delta if spec.matches_path(p)]
    if rel_delta:
        return rel_delta[:5]

    return (spec.tool_scripts + spec.configs)[:5]


def extract_failing_files(text: str, gate_name: str = "", cwd: str | None = None, root: Path | None = None, delta: set[str] | None = None) -> list[str]:
    """Multi-tiered suspect file extractor with semantic fallback."""
    actual_root = root or Path.cwd()
    found = extract_regex_suspect_files(text, cwd, actual_root)
    if not found and gate_name:
        return extract_suspect_fallback(gate_name, delta or set())

    return found


def format_suspect_files_bullet(files: list[str]) -> str:
    """Formats list of suspect files into formatted bullet string."""
    if not files:
        return "    (None detected in output)"

    return "\n".join(f"    • {f}" for f in files)


def format_banner_header(name: str, idx: int, total: int) -> str:
    """Formats ANSI failure banner top border and title."""
    header = (
        f"\n\033[1;91m================================================================\n"
        f"🚨 [IMMEDIATE FAILURE DETECTED] [{idx}/{total}] {name}\n"
        f"================================================================\033[0m\n"
    )

    return header


def format_banner_metadata(res: JobResult, files: list[str]) -> str:
    """Formats metadata lines for failure banner."""
    cmd_str = " ".join(res.cmd) if isinstance(res.cmd, list) else str(res.cmd)
    cwd_str = res.cwd or "."
    env_str = " ".join(f"{k}={v}" for k, v in res.env_overrides.items()) if res.env_overrides else "(default)"
    files_str = format_suspect_files_bullet(files)

    return (
        f"  Command       : {cmd_str}\n  Working Dir   : {cwd_str}\n  Env Overrides : {env_str}\n"
        f"  Exit Code     : {res.code} ({format_duration(res.elapsed)})\n  Failing Files :\n{files_str}\n"
        f"  Stream Log    : {normalize_repo_rel(CICD_ERRORS_LOG)}\n"
        f"  Stream JSON   : {normalize_repo_rel(CICD_ERRORS_JSON)}\n"
        f"  Stream Events : {normalize_repo_rel(CICD_EVENTS_JSONL)}\n\n"
    )


def format_failure_banner(res: JobResult, idx: int, total: int, files: list[str], err_text: str) -> str:
    """Formats full immediate ANSI failure banner for terminal stdout."""
    header = format_banner_header(res.name, idx, total)
    meta = format_banner_metadata(res, files)
    body = (
        f"{meta}\033[1mStack Trace / Failure Output:\033[0m\n"
        f"----------------------------------------------------------------\n{err_text}\n"
        f"\033[1;91m================================================================\033[0m\n"
    )

    return header + body


def print_immediate_failure_report(res: JobResult, idx: int, total: int) -> None:
    """Immediately prints full failure stack trace to terminal without waiting for suite completion."""
    err_text = extract_stack_or_error(res)
    suspect_files = extract_failing_files(err_text, res.name, res.cwd)
    banner = format_failure_banner(res, idx, total, suspect_files, err_text)
    print(banner, flush=True)


def format_error_log_entry(res: JobResult, suspect_files: list[str], ts: str, err_text: str) -> str:
    """Builds markdown entry for appending to errors.log."""
    cmd_str = " ".join(res.cmd) if isinstance(res.cmd, list) else str(res.cmd)
    files_line = f"- **Suspect Files**: {', '.join(suspect_files)}\n" if suspect_files else ""
    entry = (
        f"### [{ts}] FAIL: {res.name}\n- **Command**: `{cmd_str}`\n"
        f"- **Exit Code**: `{res.code}` ({format_duration(res.elapsed)})\n{files_line}"
        f"```text\n{strip_ansi(err_text)}\n```\n\n"
    )

    return entry


def format_run_log_line(res: JobResult) -> str:
    """Formats single line for appending to chronological run.log."""
    ts = time.strftime("%Y-%m-%d %H:%M:%S")
    status = "PASS" if res.is_success else "FAIL"
    cached_tag = " [CACHED]" if res.is_cached else ""
    cmd_str = " ".join(res.cmd) if isinstance(res.cmd, list) else str(res.cmd)
    line = f"[{ts}] [{status}{cached_tag}] {res.name} (code={res.code}, {format_duration(res.elapsed)}) | Cmd: {cmd_str}\n"
    if not res.is_success:
        err = extract_stack_or_error(res)
        line += f"  Error: {strip_ansi(err)}\n"

    return line


def direct_append_sync(file_path: Path, content: str, has_fsync: bool = False) -> None:
    """Appends content to file and flushes, syncing to disk only when requested."""
    file_path.parent.mkdir(parents=True, exist_ok=True)
    with DISK_WRITE_LOCK:
        try:
            with open(file_path, "a", encoding=DEFAULT_ENCODING) as fh:
                fh.write(content)
                fh.flush()
                if has_fsync:
                    os.fsync(fh.fileno())
        except OSError as err:
            sys.stderr.write(f"[WARN] Failed writing to {file_path}: {err}\n")


def emit_telemetry_event(
    event_type: str, session_dir: Path | None, payload: dict[str, Any], has_fsync: bool = False
) -> None:
    """Emits append-only real-time event to events.jsonl and changelog.log."""
    ts = time.strftime("%Y-%m-%dT%H:%M:%S")
    event_obj = {"timestamp": ts, "event": event_type, **payload}
    json_line = json.dumps(event_obj, separators=(",", ":")) + "\n"
    summary_msg = payload.get("message") or f"[{event_type.upper()}] {payload.get('name', '')}"
    changelog_line = f"[{ts}] {summary_msg}\n"
    targets = [(CICD_EVENTS_JSONL, json_line), (CICD_CHANGELOG_LOG, changelog_line)]
    if session_dir is not None:
        targets.extend([(session_dir / "events.jsonl", json_line), (session_dir / "changelog.log", changelog_line)])
    for path, text in targets:
        direct_append_sync(path, text, has_fsync=has_fsync)


def append_run_log_entry(res: JobResult, session_dir: Path | None) -> None:
    """Appends gate execution outcome to chronological run.log."""
    entry = format_run_log_line(res)
    targets = [CICD_RUN_LOG]
    if session_dir is not None:
        targets.append(session_dir / "run.log")
    has_fsync = not res.is_success
    for t in targets:
        direct_append_sync(t, entry, has_fsync=has_fsync)


def write_failure_markdown(entry: str, session_dir: Path | None) -> None:
    """Appends markdown failure entry to active errors.log files."""
    targets = [CICD_ERRORS_LOG]
    if session_dir is not None:
        targets.append(session_dir / "errors.log")
    for t in targets:
        direct_append_sync(t, entry, has_fsync=True)


def save_failure_trace_file(res: JobResult, err_text: str, session_dir: Path | None) -> str:
    """Saves individual failure trace to separate file and returns relative path."""
    if session_dir is None:
        return ""
    failures_dir = session_dir / "failed_tests"
    failures_dir.mkdir(parents=True, exist_ok=True)
    safe_name = "".join([c if c.isalnum() or c in ("-", "_") else "_" for c in res.name])
    failure_file = failures_dir / f"{safe_name}.log"
    failure_content = f"Test Name: {res.name}\nCommand: {res.cmd}\nCode: {res.code}\n\nStack Trace / Output:\n{err_text}"
    failure_file.write_text(failure_content, encoding=DEFAULT_ENCODING)

    return normalize_repo_rel(failure_file)


def append_failure_to_disk(res: JobResult, state: dict[str, Any], session_dir: Path | None) -> None:
    """Appends failure information to real-time logs in session and root temp."""
    ts = time.strftime("%Y-%m-%d %H:%M:%S")
    err_text = extract_stack_or_error(res)
    suspect_files = extract_failing_files(err_text)
    write_failure_markdown(format_error_log_entry(res, suspect_files, ts, err_text), session_dir)
    log_path = save_failure_trace_file(res, err_text, session_dir)
    errors = state.setdefault("errors_list", [])
    errors.append({"name": res.name, "cmd": res.cmd, "code": res.code, "elapsed": res.elapsed, "suspect_files": suspect_files, "error": strip_ansi(err_text), "log_path": log_path})
    atomic_write_json(CICD_ERRORS_JSON, errors, has_fsync=True)
    if session_dir is not None:
        atomic_write_json(session_dir / "errors.json", errors, has_fsync=True)


class TelemetryTracker:
    """Thread-safe multi-slot in-flight progress tracker and reporter."""

    def __init__(self, total_jobs: int, is_tty: bool, is_json: bool, show_all: bool, heartbeat_interval: float = DEFAULT_HEARTBEAT_INTERVAL):
        self.total_jobs = total_jobs
        self.is_tty = is_tty
        self.is_json = is_json
        self.show_all = show_all
        self.heartbeat_interval = heartbeat_interval
        self.active_jobs: dict[str, float] = {}
        self.completed_count = 0
        self.cached_count = 0
        self.lock = threading.Lock()
        self.last_print = time.monotonic()
        self.start_heartbeat()

    def start_heartbeat(self) -> None:
        """Spawns background heartbeat daemon."""
        self._stop_event = threading.Event()
        self._thread = threading.Thread(target=self._run_heartbeat, daemon=True)
        self._thread.start()

    def _run_heartbeat(self) -> None:
        while not self._stop_event.is_set():
            if self._stop_event.wait(timeout=self.heartbeat_interval):
                break
            with self.lock:
                act = len(self.active_jobs)
            if act > 0:
                self.tick(force=True)

    def stop_heartbeat(self) -> None:
        """Stops background heartbeat daemon."""
        if hasattr(self, "_stop_event"):
            self._stop_event.set()

    def start_job(self, name: str) -> None:
        """Registers a newly launched job in active tracking."""
        with self.lock:
            self.active_jobs[name] = time.monotonic()
        self.tick(force=False)

    def finish_job(self, name: str, is_cached: bool = False) -> None:
        """Removes job from active tracking and increments counter."""
        with self.lock:
            if name in self.active_jobs:
                del self.active_jobs[name]
            self.completed_count += 1
            if is_cached:
                self.cached_count += 1
        self.tick(force=False)

    def _format_active_summary(self, now: float) -> tuple[int, list[str], float]:
        with self.lock:
            active_list = list(self.active_jobs.items())
            count = len(active_list)
            items = [f"{n} ({format_duration(now - st)})" for n, st in active_list[:4]]
            if count > 4:
                items.append(f"+{count - 4} more")
            oldest_elapsed = max((now - st for _, st in active_list), default=0.0)

        return count, items, oldest_elapsed

    def _build_tick_message(self, act_count: int, items: list[str], elapsed: float) -> str:
        items_str = ", ".join(items)
        if self.completed_count == 0:
            return f"[IN-FLIGHT] {act_count} active: [{items_str}] | in-progress ({format_duration(elapsed)} elapsed)"

        pct = int(100.0 * self.completed_count / max(1, self.total_jobs))

        return f"[IN-FLIGHT] {act_count} active: [{items_str}] | {self.completed_count}/{self.total_jobs} done ({pct}%)"

    def tick(self, force: bool = False) -> None:
        """Emits progress heartbeat if interval has elapsed (strictly every 25s or more)."""
        if self.is_json:
            return
        now = time.monotonic()
        if not force and (now - self.last_print < self.heartbeat_interval):
            return
        self.last_print = now
        act_count, items, elapsed = self._format_active_summary(now)
        if act_count == 0 and self.completed_count == self.total_jobs:
            return
        msg = self._build_tick_message(act_count, items, elapsed)
        self._write_heartbeat(msg)

    def _write_heartbeat(self, msg: str) -> None:
        if self.is_tty:
            sys.stdout.write(f"\r{msg}\033[K")
            sys.stdout.flush()
        else:
            print(msg, flush=True)

    def clear_line(self) -> None:
        """Clears terminal line if in TTY mode."""
        if self.is_tty and not self.is_json:
            sys.stdout.write("\r\033[K")
            sys.stdout.flush()


def load_previous_state(force: bool) -> dict[str, Any]:
    """Loads previous state.json if present and not bypassed by force flag."""
    if not force and CICD_STATE_JSON.is_file():
        try:
            return json.loads(CICD_STATE_JSON.read_text(encoding=DEFAULT_ENCODING))
        except Exception as err:
            return {}

    return {}


def init_session_scaffolding(force: bool, resume: bool) -> tuple[Path, dict[str, Any]]:
    """Initializes session folder, latest pointer, and loads previous state."""
    CICD_RUNS_DIR.mkdir(parents=True, exist_ok=True)
    ts = time.strftime("%Y%m%d_%H%M%S")
    import uuid
    run_hash = hashlib.md5(str(uuid.uuid4()).encode()).hexdigest()[:8]
    session_dir = CICD_RUNS_DIR / f"{run_hash}-{ts}"
    session_dir.mkdir(parents=True, exist_ok=True)
    link_latest_session(session_dir, CICD_LATEST_DIR)
    prev_state = load_previous_state(force)
    init_files(session_dir, ts)

    return session_dir, prev_state


def init_empty_files(session_dir: Path, ts: str) -> None:
    """Writes empty headers to errors and run logs."""
    run_header = f"# CI/CD Run Log — {ts}\n\n"
    CICD_ERRORS_LOG.write_text(f"# CI/CD Failure Log — {ts}\n\n", encoding=DEFAULT_ENCODING)
    (session_dir / "errors.log").write_text(f"# CI/CD Failure Log — {ts}\n\n", encoding=DEFAULT_ENCODING)
    atomic_write_json(CICD_ERRORS_JSON, [])
    atomic_write_json(session_dir / "errors.json", [])
    CICD_RUN_LOG.write_text(run_header, encoding=DEFAULT_ENCODING)
    (session_dir / "run.log").write_text(run_header, encoding=DEFAULT_ENCODING)


def init_files(session_dir: Path, ts: str) -> None:
    """Initializes log files with session header in root and session dir."""
    init_empty_files(session_dir, ts)
    CICD_EVENTS_JSONL.write_text("", encoding=DEFAULT_ENCODING)
    (session_dir / "events.jsonl").write_text("", encoding=DEFAULT_ENCODING)
    CICD_CHANGELOG_LOG.write_text(f"# CI/CD Changelog — {ts}\n\n", encoding=DEFAULT_ENCODING)
    (session_dir / "changelog.log").write_text(f"# CI/CD Changelog — {ts}\n\n", encoding=DEFAULT_ENCODING)


def serialize_state(state: dict[str, Any]) -> dict[str, Any]:
    """Prepares serializable state dictionary without non-primitive objects."""
    clean = {k: list(v) if isinstance(v, (set, frozenset)) else v for k, v in state.items() if k != "results"}
    clean["results"] = [asdict(r) for r in state.get("results", [])]

    return clean


def persist_state(session_dir: Path, state: dict[str, Any]) -> None:
    """Atomically writes state.json to session dir and root."""
    state["updated_at"] = time.strftime("%Y-%m-%d %H:%M:%S")
    payload = serialize_state(state)
    atomic_write_json(session_dir / "state.json", payload)
    atomic_write_json(CICD_STATE_JSON, payload)


def update_summary_file(session_dir: Path, meta: dict[str, Any]) -> None:
    """Atomically writes summary.json to session dir and root."""
    atomic_write_json(session_dir / "summary.json", meta)
    atomic_write_json(CICD_SUMMARY_JSON, meta)


def compute_summary_meta(total: int, results: list, errors_list: list, is_finished: bool) -> dict[str, Any]:
    """Builds summary statistics dictionary for summary.json."""
    passed = sum(1 for r in results if r.is_success)
    failed = sum(1 for r in results if not r.is_success)
    status = "completed" if is_finished and failed == 0 else ("failed" if is_finished else "running")

    return {
        "status": status, "total_gates": total, "passed_gates": passed, "failed_gates": failed,
        "remaining_gates": max(0, total - passed - failed), "cached_gates": sum(1 for r in results if r.is_cached),
        "errors_count": len(errors_list), "active_failures": [e["name"] for e in errors_list],
    }


def update_cicd_summary(state: dict[str, Any], session_dir: Path, is_finished: bool = False) -> None:
    """Updates summary.json metadata with real-time status and counts."""
    meta = compute_summary_meta(state["total"], state.get("results", []), state.get("errors_list", []), is_finished)
    update_summary_file(session_dir, meta)


def get_cache_size_and_count(cache_dir: Path) -> tuple[int, int]:
    """Calculates total byte size and file count in cache directory."""
    if not cache_dir.is_dir():
        return 0, 0
    files = [f for f in cache_dir.iterdir() if f.is_file()]
    total_size = sum(f.stat().st_size for f in files)

    return total_size, len(files)


def get_cache_age_info(cache_dir: Path) -> tuple[str, float]:
    """Finds latest modification timestamp and elapsed seconds across cache files."""
    if not cache_dir.is_dir():
        return "N/A", 0.0
    files = [f for f in cache_dir.iterdir() if f.is_file()]
    if not files:
        return "N/A", 0.0
    latest_mtime = max(f.stat().st_mtime for f in files)
    elapsed = max(0.0, time.time() - latest_mtime)
    formatted = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(latest_mtime))

    return formatted, elapsed


def get_cache_hit_rate_summary() -> str:
    """Reads state and inventory cache to calculate observed hit rates."""
    parts = []
    if CACHE_STATE_FILE.is_file():
        res = load_previous_state(False).get("results", [])
        if res:
            cached = sum(1 for r in res if isinstance(r, dict) and r.get("is_cached"))
            parts.append(f"Gates: {cached}/{len(res)} ({int(100 * cached / len(res))}%)")
    if CACHE_INVENTORY_FILE.is_file():
        summ = load_raw_test_inventory(CACHE_INVENTORY_FILE).get("summary", {})
        if summ.get("total", 0) > 0:
            c, t = summ.get("cached", 0), summ.get("total", 1)
            parts.append(f"Tests: {c}/{t} ({int(100 * c / t)}%)")

    return " | ".join(parts) if parts else "No active run data"


def display_cache_stats() -> None:
    """Prints cache directory size, item counts, hit rates, and age."""
    total_bytes, count = get_cache_size_and_count(CACHE_DIR)
    latest_date, elapsed = get_cache_age_info(CACHE_DIR)
    hit_rates = get_cache_hit_rate_summary()
    print("================================================================")
    print("                CI/CD LOCAL RUNNER CACHE STATS                  ")
    print("================================================================")
    print(f"📁 Cache Directory   : {CACHE_DIR}")
    print(f"📦 Cached Files      : {count} items ({total_bytes / 1024:.1f} KB)")
    print(f"⏱️  Latest Cache Age : {latest_date} ({elapsed:.1f}s ago)")
    print(f"🎯 Cache Hit Rates   : {hit_rates}")
    print("================================================================")


def remove_path_safely(p: Path) -> None:
    """Removes single file or directory path safely."""
    if p.is_file():
        try:
            p.unlink()
        except OSError:
            pass
    elif p.is_dir():
        robust_rmtree(p)


def clear_cicd_cache() -> None:
    """Safely removes cache files, resets runner-eta.json, and clears memo cache."""
    FILE_HASH_MEMO_CACHE.clear()
    if CACHE_DIR.is_dir():
        for item in CACHE_DIR.iterdir():
            remove_path_safely(item)
    for legacy in (CICD_DIR / "state.json", CICD_DIR / "timings.json", CICD_DIR / "test-inventory.json", CICD_DIR / "last_run_cache.json"):
        remove_path_safely(legacy)
    if RUNNER_ETA_FILE.is_file():
        atomic_write_json(RUNNER_ETA_FILE, {"status": "idle", "eta_sec": 0, "active_jobs": []})
    print(f"✔ CI/CD cache cleared successfully ({CACHE_DIR}).")


def prune_session_runs_by_ttl(runs_dir: Path, ttl_sec: float) -> int:
    """Removes session directories older than TTL."""
    if not runs_dir.is_dir():
        return 0
    now = time.time()
    pruned = 0
    for p in runs_dir.iterdir():
        if p.is_dir() and (now - p.stat().st_mtime > ttl_sec):
            robust_rmtree(p)
            pruned += 1

    return pruned


def prune_session_runs_by_count(runs_dir: Path, keep_runs: int) -> int:
    """Keeps the most recent session directories, removing excess."""
    if not runs_dir.is_dir():
        return 0
    dirs = [p for p in runs_dir.iterdir() if p.is_dir()]
    if len(dirs) <= keep_runs:
        return 0
    dirs.sort(key=lambda p: p.stat().st_mtime, reverse=True)
    for old_dir in dirs[keep_runs:]:
        robust_rmtree(old_dir)

    return max(0, len(dirs) - keep_runs)


def prune_cicd_cache(max_size_mb: int = 25, keep_runs: int = 3, ttl_hours: int = 24) -> None:
    """Enforces size budget and purges old sessions and temporary artifacts."""
    ttl_sec = ttl_hours * 3600.0
    pruned_ttl = prune_session_runs_by_ttl(CICD_RUNS_DIR, ttl_sec)
    pruned_cnt = prune_session_runs_by_count(CICD_RUNS_DIR, keep_runs)
    clean_stale_temp_artifacts()
    total_bytes, _ = get_cache_size_and_count(CACHE_DIR)
    if total_bytes > max_size_mb * 1024 * 1024:
        clear_cicd_cache()
    print(f"✔ CI/CD cache pruned: removed {pruned_ttl + pruned_cnt} old run(s), budget: {max_size_mb}MB.")


def filter_batch_jobs(jobs: dict[str, Any], query: str) -> dict[str, Any]:
    """Filters dictionary of jobs by name substring query."""
    q = query.lower()
    matched = {k: v for k, v in jobs.items() if q in k.lower()}

    return matched


def filter_job_batches(batches: list[dict[str, Any]], query: str | None) -> list[dict[str, Any]]:
    """Returns filtered copy of batches matching user search query."""
    if not query:
        return batches
    filtered = []
    for b in batches:
        matched = filter_batch_jobs(b.get("jobs", {}), query)
        if matched:
            copy_b = dict(b)
            copy_b["jobs"] = matched
            filtered.append(copy_b)

    return filtered


def get_current_host_os() -> str:
    """Returns normalized host OS string: windows, darwin, or linux."""
    is_windows = sys.platform == "win32"
    if is_windows:
        return "windows"
    is_darwin = sys.platform == "darwin"
    if is_darwin:
        return "darwin"

    return "linux"


def is_gate_matching_host_os(gate_name: str, host_os: str) -> bool:
    """Checks if a platform-specific gate matches current host operating system."""
    name_lower = gate_name.lower()
    has_linux_mismatch = "(linux)" in name_lower and host_os != "linux"
    if has_linux_mismatch:
        return False
    has_darwin_mismatch = "(darwin)" in name_lower and host_os != "darwin"
    if has_darwin_mismatch:
        return False
    has_win_mismatch = "(windows)" in name_lower and host_os != "windows"
    if has_win_mismatch:
        return False

    return True


def filter_single_batch_by_host_os(b: dict[str, Any], host_os: str) -> dict[str, Any] | None:
    """Filters jobs within a single batch by host OS matching."""
    jobs = b.get("jobs", {})
    matched = {k: v for k, v in jobs.items() if is_gate_matching_host_os(k, host_os)}
    has_matches = bool(matched)
    if not has_matches:
        return None
    copy_b = dict(b)
    copy_b["jobs"] = matched

    return copy_b


def filter_batches_by_host_os(batches: list[dict[str, Any]], allow_all_os: bool) -> list[dict[str, Any]]:
    """Filters out cross-platform gates not matching current host OS unless allow_all_os is True."""
    is_bypass = allow_all_os or os.environ.get("CI") == "true"
    if is_bypass:
        return batches
    host_os = get_current_host_os()
    filtered = []
    for b in batches:
        filtered_batch = filter_single_batch_by_host_os(b, host_os)
        if filtered_batch is not None:
            filtered.append(filtered_batch)

    return filtered


def is_test_batch(batch_name: str) -> bool:
    """Detects whether an entire batch represents a test execution stage."""
    b_lower = batch_name.lower()
    return any(k in b_lower for k in ("smoke", "unit test", "coverage", "race"))


def is_test_job(job_name: str, command: Any) -> bool:
    """Identifies individual test execution suites, smoke tests, and script unit tests."""
    name_lower = job_name.lower()
    if any(k in name_lower for k in ("unit test", "smoke", "coverage", "race test")):
        return True
    cmd_str = str(command).lower()
    if "test_ci_scripts.py" in cmd_str:
        return True
    return False


def create_base_arg_parser() -> argparse.ArgumentParser:
    """Initializes ArgumentParser with basic description and epilog."""
    banner = (
        "Fast Multi-Worker Local CI/CD Runner with Incremental Caching & Telemetry.\n\n"
        "⚡ PRIORITY SHORTCUTS:\n"
        "  python 03-ai-scripts/06-cicd-local-runner.py run-smart        (or: smart, --smart, -s)\n"
        "  python 03-ai-scripts/06-cicd-local-runner.py run-incremental  (or: incremental, --incremental)\n"
        "  -> Checks Git changed files, builds ONLY changed packages into OS temp, and runs Quad Runner.\n"
    )

    return argparse.ArgumentParser(
        prog="python 03-ai-scripts/06-cicd-local-runner.py",
        description=banner,
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )


def add_execution_mode_arguments(parser: argparse.ArgumentParser) -> None:
    """Adds execution mode CLI arguments."""
    parser.add_argument(
        "--smart", "--incremental", "-s",
        dest="smart_mode",
        action="store_true",
        help="[PRIORITY] Execute smart incremental Go tests on changed packages only.",
    )
    parser.add_argument("--sync", dest="sync_mode", action="store_true", help="Run sequentially.")
    parser.add_argument("-w", "--workers", type=int, default=DEFAULT_WORKERS, help="Worker threads.")
    parser.add_argument("--io-workers", type=int, default=DEFAULT_IO_WORKERS, help="IO worker limit.")
    parser.add_argument("-t", "--timeout", type=int, default=DEFAULT_TIMEOUT_SEC, help="Job timeout.")
    parser.add_argument("--filter", type=str, default=None, help="Filter jobs by substring.")
    parser.add_argument(
        "--all-os",
        dest="all_os",
        action="store_true",
        help="Run cross-OS build-tag and vet gates for all operating systems (default: current host OS only).",
    )
    parser.add_argument(
        "--no-tests", "--skip-tests",
        dest="no_tests",
        action="store_true",
        help="Skip all test execution suites (standard development mode; tests disabled by default unless commanded by owner)."
    )
    parser.add_argument(
        "--run-tests", "--with-tests",
        dest="run_tests",
        action="store_true",
        help="Explicitly execute all unit test suites, integration tests, and coverage checks."
    )
    parser.add_argument(
        "--fast",
        dest="fast_mode",
        action="store_true",
        help="Execute only hot and warm tests based on heatmap, skipping cold tests."
    )
    parser.add_argument(
        "--pkg", "--package", "-p", "--target-file", "--file",
        dest="package_filter", nargs="*", default=[],
        help="Run specific Go test package based on code file path, code file name, or Go package name."
    )
    parser.add_argument(
        "positional_targets", nargs="*", default=[],
        help="Optional positional package names, code file paths, or file names to test."
    )


def add_reporting_arguments(parser: argparse.ArgumentParser) -> None:
    """Adds reporting and visibility CLI arguments."""
    parser.add_argument("--all-paths", "--all-passed", "--all", dest="show_all", action="store_true", help="Show all.")
    parser.add_argument("--failed", dest="failed_only", action="store_true", help="Show failed only.")
    parser.add_argument("-o", "--output-paths", dest="output_paths", nargs="*", type=str, default=[], help="Multiple output file paths.")
    parser.add_argument("--json", dest="json_mode", action="store_true", help="Output machine-readable JSON.")
    parser.add_argument("--heatmap", dest="show_heatmap", action="store_true", help="Display test suite heat map and exit.")
    parser.add_argument("--top", dest="top_n", type=int, default=25, help="Limit heat map display to top N entries (default: 25).")
    parser.add_argument("--eta-interval", type=int, default=120, help="Print ETA interval in seconds.")
    parser.add_argument("--heartbeat-interval", type=float, default=DEFAULT_HEARTBEAT_INTERVAL, help="Heartbeat interval in seconds for in-flight progress (default: 25.0).")
    parser.add_argument("--inventory-only", dest="inventory_only", action="store_true", help="Discover and catalog all tests into JSON manifest and exit.")


def add_caching_and_resume_arguments(parser: argparse.ArgumentParser) -> None:
    """Adds incremental caching and crash resumption arguments."""
    parser.add_argument("--force", "--fresh", "--clean", "--no-cache", dest="force_run", action="store_true", help="Run all.")
    parser.add_argument("--cache-stats", dest="cache_stats", action="store_true", help="Display cache statistics and exit.")
    parser.add_argument("--clear-cache", "--clean-cache", dest="clear_cache", action="store_true", help="Purge CI/CD cache and exit.")
    parser.add_argument("--prune-cache", dest="prune_cache", action="store_true", help="Enforce cache size budget and exit.")
    parser.add_argument("--clean-gocache", dest="clean_gocache", action="store_true", help="Purge Go build and test caches.")
    parser.add_argument("--resume", dest="resume_mode", action="store_true", help="Resume interrupted session.")
    parser.add_argument("--changed-only", "-c", "--recent", dest="changed_only", action="store_true", help="Scope linters to files changed in recent commits.")
    parser.add_argument("--commits", "-n", dest="commits", type=int, default=20, help="Commit window for changed files (default: 20).")


def normalize_smart_command_args(args: argparse.Namespace) -> None:
    """Normalizes positional smart/incremental command tokens into smart execution mode."""
    smart_tokens = {"smart", "run-smart", "incremental", "run-incremental", "smart-test", "test"}
    remaining_targets = []
    for t in getattr(args, "package_filter", []):
        if t.lower() in smart_tokens:
            args.smart_mode = True
        else:
            remaining_targets.append(t)
    args.package_filter = remaining_targets
    if getattr(args, "smart_mode", False) and not args.filter:
        args.filter = "Go Smart Incremental Tests"


def parse_args() -> argparse.Namespace:
    """Constructs CLI argument parser with comprehensive argument support."""
    parser = create_base_arg_parser()
    add_execution_mode_arguments(parser)
    add_reporting_arguments(parser)
    add_caching_and_resume_arguments(parser)

    args = parser.parse_args()
    raw_targets = list(getattr(args, "package_filter", []) or []) + list(getattr(args, "positional_targets", []) or [])
    args.package_filter = [t for t in raw_targets if t]
    normalize_smart_command_args(args)

    return args


def collect_tool_stats(repo_root: Path, scripts: list[str]) -> dict[str, list]:
    """Collects current file stats for tool scripts."""
    stats = {}
    for s in scripts:
        st = check_file_stat(repo_root, s)
        if st is not None:
            stats[s] = list(st)

    return stats


def collect_artifact_stats(repo_root: Path, artifacts: list[str]) -> dict[str, list]:
    """Collects current file stats for artifacts."""
    stats = {}
    for a in artifacts:
        st = check_file_stat(repo_root, a)
        if st is not None:
            stats[a] = list(st)

    return stats


def record_gate_success(state: dict[str, Any], name: str, res: JobResult, cmd_hash: str, spec: GateSpec, root: Path) -> None:
    """Records passing gate execution outcome in persistent state."""
    state.setdefault("gates", {})[name] = {
        "status": "PASSED", "code": 0, "cmd_hash": cmd_hash, "elapsed": res.elapsed,
        "is_cached": res.is_cached, "tool_stats": collect_tool_stats(root, spec.tool_scripts),
        "artifact_stats": collect_artifact_stats(root, spec.artifact_outputs + spec.artifact_inputs),
    }


def record_gate_failure(state: dict[str, Any], name: str, res: JobResult, cmd_hash: str) -> None:
    """Records failing gate execution outcome in persistent state."""
    err_text = extract_stack_or_error(res)
    state.setdefault("gates", {})[name] = {
        "status": "FAILED", "code": res.code, "cmd_hash": cmd_hash, "elapsed": res.elapsed,
        "error": strip_ansi(err_text),
    }


def handle_success_result(res: JobResult, cmd_hash: str, spec: GateSpec, state: dict, sdir: Path, tel: TelemetryTracker, root: Path) -> None:
    """Handles reporting and state update for a passed quality gate."""
    record_gate_success(state, res.name, res, cmd_hash, spec, root)
    evt = "gate_cached" if res.is_cached else "gate_passed"
    emit_telemetry_event(evt, sdir, {"name": res.name, "elapsed": res.elapsed})
    if tel.show_all and not state["is_json"]:
        tel.clear_line()
        tag = " [cached]" if res.is_cached else ""
        idx, total = state["counter"], state["total"]
        print(f"  [{idx:2d}/{total}] \033[1;92m✓ PASS\033[0m [{res.name}] ({format_duration(res.elapsed)}){tag}", flush=True)


def handle_failure_result(res: JobResult, cmd_hash: str, state: dict, sdir: Path, tel: TelemetryTracker) -> None:
    """Handles reporting and state update for a failed quality gate."""
    record_gate_failure(state, res.name, res, cmd_hash)
    append_failure_to_disk(res, state, sdir)
    emit_telemetry_event("gate_failed", sdir, {"name": res.name, "code": res.code, "elapsed": res.elapsed}, has_fsync=True)
    if not state["is_json"]:
        tel.clear_line()
        print_immediate_failure_report(res, state["counter"], state["total"])


def track_executed_gate(res: JobResult, state: dict[str, Any]) -> None:
    """Tracks non-cached executed gate in state."""
    if not res.is_cached:
        state.setdefault("freshly_executed", set()).add(res.name)


def dispatch_result_outcome(
    res: JobResult, cmd_hash: str, spec: GateSpec, state: dict, sdir: Path, tel: TelemetryTracker, root: Path
) -> None:
    """Dispatches result to either success or failure handler."""
    if res.is_success:
        handle_success_result(res, cmd_hash, spec, state, sdir, tel, root)
    else:
        handle_failure_result(res, cmd_hash, state, sdir, tel)


def handle_completed_result(
    res: JobResult, cmd_hash: str, spec: GateSpec, state: dict[str, Any], sdir: Path, tel: TelemetryTracker, root: Path
) -> None:
    """Processes completed JobResult, updates state, logs, and telemetry with thread synchronization."""
    with DISK_WRITE_LOCK:
        state["counter"] += 1
        state["results"].append(res)
        track_executed_gate(res, state)
        append_run_log_entry(res, sdir)
        dispatch_result_outcome(res, cmd_hash, spec, state, sdir, tel, root)
        tel.finish_job(res.name, is_cached=res.is_cached)


def check_single_gate_skip(name: str, cmd: Any, prev: dict, delta: set, root: Path, executed: set) -> tuple[bool, str, str, GateSpec, Any]:
    """Evaluates whether an individual gate meets caching conditions."""
    raw_cmd = cmd.get("cmd") if isinstance(cmd, dict) else cmd
    env = cmd.get("env") if isinstance(cmd, dict) else None
    cwd = cmd.get("cwd") if isinstance(cmd, dict) else None
    cmd_hash = compute_cmd_hash(raw_cmd, env, cwd)
    spec = GATE_SPECS.get(name, GateSpec(name))
    can_skip, reason = evaluate_gate_skip(spec, prev.get("gates", {}).get(name), cmd_hash, delta, root, executed)

    return can_skip, reason, cmd_hash, spec, raw_cmd


def build_cached_job_result(name: str, raw_cmd: Any, cmd: Any, reason: str) -> JobResult:
    """Creates JobResult for a skipped cached gate."""
    cwd = cmd.get("cwd") if isinstance(cmd, dict) else None
    env = cmd.get("env") if isinstance(cmd, dict) else None

    return JobResult(
        name=name, cmd=raw_cmd, code=0, out=f"[CACHED] {reason}", err="",
        elapsed=0.0, cwd=cwd, env_overrides=env, is_cached=True,
    )


def evaluate_batch_skips(
    items: list[tuple[str, Any]], state: dict, prev: dict, delta: set,
    root: Path, sdir: Path, tel: TelemetryTracker, args: argparse.Namespace | None = None
) -> list[tuple[str, Any, str, GateSpec]]:
    """Evaluates gates for skipping; records skipped gates and returns list of jobs to run."""
    to_run = []
    executed = state.get("freshly_executed", set())
    has_pkg_filter = bool(args and (getattr(args, "package_filter", None) or getattr(args, "smart_mode", False)))
    for name, cmd in items:
        if has_pkg_filter and name == "Go Smart Incremental Tests":
            raw = cmd.get("cmd") if isinstance(cmd, dict) else cmd
            env = cmd.get("env") if isinstance(cmd, dict) else None
            cwd = cmd.get("cwd") if isinstance(cmd, dict) else None
            to_run.append((name, cmd, compute_cmd_hash(raw, env, cwd), GATE_SPECS.get(name, GateSpec(name))))
            continue
        can_skip, reason, cmd_hash, spec, raw_cmd = check_single_gate_skip(name, cmd, prev, delta, root, executed)
        if can_skip:
            res = build_cached_job_result(name, raw_cmd, cmd, reason)
            handle_completed_result(res, cmd_hash, spec, state, sdir, tel, root)
        else:
            to_run.append((name, cmd, cmd_hash, spec))

    return to_run


def adapt_cmd_for_changed_only(name: str, raw_cmd: Any, args: argparse.Namespace) -> Any:
    """Appends changed-only flags to supported python checkers."""
    is_changed = getattr(args, "changed_only", False)
    if not is_changed or not isinstance(raw_cmd, list):
        return raw_cmd
    supported_gates = {"Nested If Linter", "Boolean & Enum Linter", "Relative Path Check"}
    if name in supported_gates:
        commits = getattr(args, "commits", 20)
        return list(raw_cmd) + ["--changed-only", "--commits", str(commits)]

    return raw_cmd


def submit_job_futures(executor: ThreadPoolExecutor, to_run: list, args: argparse.Namespace, telemetry: TelemetryTracker, root: Path) -> dict:
    """Submits active jobs to ThreadPoolExecutor and registers telemetry."""
    future_map = {}
    for name, cmd, cmd_hash, spec in to_run:
        telemetry.start_job(name)
        if (isinstance(cmd, dict) and cmd.get("type") == "smart_go_tests") or name == "Go Smart Incremental Tests":
            fut = executor.submit(
                run_smart_go_tests, name, args.timeout, args.workers, bool(args.force_run), root, telemetry, getattr(args, "package_filter", [])
            )
            future_map[fut] = (name, cmd_hash, spec, cmd)
            continue
        raw_cmd = cmd.get("cmd") if isinstance(cmd, dict) else cmd
        raw_cmd = adapt_cmd_for_changed_only(name, raw_cmd, args)
        env = {**os.environ, **cmd.get("env")} if isinstance(cmd, dict) and "env" in cmd else None
        cwd = cmd.get("cwd") if isinstance(cmd, dict) else None
        fut = executor.submit(run_job, name, raw_cmd, args.timeout, env, cwd)
        future_map[fut] = (name, cmd_hash, spec, cmd)

    return future_map


def wait_and_handle_batch_futures(fut_map: dict, state: dict, sdir: Path, tel: TelemetryTracker, root: Path) -> None:
    """Collects completed batch futures from executor and handles outcomes."""
    for fut in as_completed(fut_map):
        name, cmd_hash, spec, cmd = fut_map[fut]
        try:
            res = fut.result()
        except Exception as ex:
            raw = cmd.get("cmd") if isinstance(cmd, dict) else cmd
            res = JobResult(name=name, cmd=raw, code=1, out="", err=str(ex), elapsed=0.0)
        handle_completed_result(res, cmd_hash, spec, state, sdir, tel, root)


def execute_job_batch(
    batch: dict[str, Any], workers: int, is_sync: bool, args: argparse.Namespace, state: dict[str, Any],
    prev_state: dict[str, Any], repo_delta: set[str], root: Path, sdir: Path, tel: TelemetryTracker
) -> None:
    """Executes all jobs within a single batch with worker pool."""
    items = list(batch["jobs"].items())
    to_run = evaluate_batch_skips(items, state, prev_state, repo_delta, root, sdir, tel, args)
    if to_run:
        limit = batch.get("max_workers")
        batch_workers = 1 if is_sync else min(workers, limit or workers, len(to_run))
        with ThreadPoolExecutor(max_workers=batch_workers) as executor:
            fut_map = submit_job_futures(executor, to_run, args, tel, root)
            wait_and_handle_batch_futures(fut_map, state, sdir, tel, root)
    persist_state(sdir, state)
    update_cicd_summary(state, sdir, is_finished=False)


def format_segment_gate_bullet(idx: int, gate_name: str) -> str:
    """Formats an individual gate bullet for upfront listing."""
    return f"    [{idx:2d}] • {gate_name}"


def print_segment_summary(s_idx: int, batch: dict[str, Any]) -> None:
    """Prints single segment header and its list of queued gates."""
    b_name = batch.get("name", f"Segment {s_idx}")
    jobs = batch.get("jobs", {})
    limit = batch.get("max_workers")
    worker_desc = f"max {limit} workers" if limit else "parallel worker pool"
    print(f"\n  Segment {s_idx}: \033[1m{b_name}\033[0m ({len(jobs)} gates, {worker_desc})")
    for g_idx, gate_name in enumerate(jobs.keys(), 1):
        print(format_segment_gate_bullet(g_idx, gate_name))


def print_queued_segments_and_gates(batches: list[dict[str, Any]], workers_label: str) -> None:
    """Prints informational execution banner and all queued test segments and gates upfront."""
    total_jobs = sum(len(b["jobs"]) for b in batches)
    print("================================================================")
    print("           PARALLEL LOCAL CI/CD QUALITY GATE RUNNER             ")
    print("================================================================")
    print(f"🚀 Execution Mode          : {workers_label}")
    print(f"📋 Total Enqueued Segments : {len(batches)}")
    print(f"📋 Total Enqueued Gates    : {total_jobs}")
    print("----------------------------------------------------------------")
    print("📋 QUEUED TEST SEGMENTS & GATES:")
    for s_idx, batch in enumerate(batches, 1):
        print_segment_summary(s_idx, batch)
    print("\n================================================================\n", flush=True)


def build_json_payload(results: list[JobResult], counts: tuple[int, int, int, int], total_elapsed: float) -> dict[str, Any]:
    """Constructs dictionary structure for machine-readable JSON output."""
    total, passed, failed, timeout = counts
    payload = {
        "status": "passed" if failed == 0 and timeout == 0 else "failed",
        "total_elapsed_sec": total_elapsed,
        "summary": {"total": total, "passed": passed, "failed": failed, "timeout": timeout, "cached": sum(1 for r in results if r.is_cached)},
        "results": [asdict(r) for r in results],
    }

    return payload


def handle_json_output(args: argparse.Namespace, results: list[JobResult], counts: tuple[int, int, int, int], total_elapsed: float, has_failures: bool) -> int:
    """Emits JSON formatted output to stdout or specified file."""
    payload = build_json_payload(results, counts, total_elapsed)
    out_str = json.dumps(payload, indent=2)
    if getattr(args, "output_paths", []):
        for p in args.output_paths:
            Path(p).write_text(out_str, encoding=DEFAULT_ENCODING)
    else:
        print(out_str)
    exit_code = 1 if has_failures else 0

    return exit_code


def format_report_row(res: JobResult) -> str:
    """Formats single job result into summary table row."""
    status_str = "\033[1;92mPASS\033[0m" if res.is_success else "\033[1;91mFAIL\033[0m"
    cached_flag = " [cached]" if res.is_cached else ""
    dur_str = format_duration(res.elapsed)

    return f"  [{status_str}] {res.name:<32} ({dur_str}) exit={res.code}{cached_flag}"


def format_report_header(total: int, passed: int, failed: int, timeout: int, elapsed: float, cached: int) -> list[str]:
    """Builds top section lines for summary report."""
    return [
        "================================================================",
        "                     CI/CD EXECUTION REPORT                     ",
        "================================================================",
        f"  Total Gates   : {total}",
        f"  Passed        : {passed} ({cached} cached)",
        f"  Failed        : {failed}",
        f"  Timeouts      : {timeout}",
        f"  Total Duration: {format_duration(elapsed)}",
        "----------------------------------------------------------------",
    ]


def format_full_report(results: list[JobResult], total: int, passed: int, failed: int, timeout: int, elapsed: float, show_all: bool) -> str:
    """Renders formatted execution table of quality gate outcomes."""
    cached = sum(1 for r in results if r.is_cached)
    lines = format_report_header(total, passed, failed, timeout, elapsed, cached)
    for r in results:
        if show_all or not r.is_success:
            lines.append(format_report_row(r))
    lines.append("================================================================")

    return "\n".join(lines)


def format_remediation_banner_header(failed_count: int) -> list[str]:
    """Generates the header section for the AI remediation summary banner."""
    plural_suffix = "S" if failed_count != 1 else ""
    lines = [
        "\n\033[1;91m" + "=" * 70,
        f"🚨 AI AGENT REMEDIATION SUMMARY & LOG LOCATIONS ({failed_count} GATE FAILURE{plural_suffix})",
        "=" * 70 + "\033[0m",
    ]

    return lines


def format_log_locations_section(session_dir: Path | None) -> list[str]:
    """Formats exact relative log file locations for AI agent inspection."""
    sdir_rel = normalize_repo_rel(session_dir) if session_dir else ".ai-memory/cicd/latest"
    sdir_rel = sdir_rel.rstrip("/")
    lines = [
        "📂 \033[1mLog Files & Artifact Locations\033[0m:",
        f"  • Live Markdown Stream  : {normalize_repo_rel(CICD_ERRORS_LOG)}",
        f"  • Structured JSON Errors: {normalize_repo_rel(CICD_ERRORS_JSON)}",
        f"  • Live Event Stream     : {normalize_repo_rel(CICD_EVENTS_JSONL)}",
        f"  • Real-time Telemetry   : {normalize_repo_rel(CICD_SUMMARY_JSON)}",
        f"  • Incremental Cache     : {normalize_repo_rel(CICD_STATE_JSON)}",
        f"  • Full Chronological Log: {normalize_repo_rel(CICD_RUN_LOG)}",
        f"  • Session Run Directory : {sdir_rel}/",
        f"  • Individual Fail Logs  : {sdir_rel}/failed_tests/",
    ]

    return lines


def format_suspect_files_summary(errors_list: list[dict[str, Any]]) -> list[str]:
    """Aggregates and formats deduplicated suspect files across failing gates."""
    seen: set[str] = set()
    for err in errors_list:
        for f in err.get("suspect_files", []):
            seen.add(normalize_repo_rel(f))
    if not seen:
        return ["🔍 \033[1mSuspect Files to Inspect\033[0m: (None detected in error output)"]
    lines = ["🔍 \033[1mSuspect Files to Inspect\033[0m (deduplicated):"]
    for f in sorted(seen)[:12]:
        lines.append(f"  • {f}")

    return lines


def format_retest_commands_section(errors_list: list[dict[str, Any]]) -> list[str]:
    """Generates copy-pasteable targeted re-test commands for failed gates."""
    lines = ["🎯 \033[1mTargeted Single-Gate Re-Test Commands\033[0m:"]
    for err in errors_list:
        name = err.get("name", "")
        raw_cmd = err.get("cmd", "")
        cmd_str = " ".join(raw_cmd) if isinstance(raw_cmd, list) else str(raw_cmd)
        lines.append(f"  • Gate: \"{name}\"")
        lines.append(f"    Runner Filter: python 03-ai-scripts/06-cicd-local-runner.py --filter \"{name}\"")
        lines.append(f"    Direct Exec  : {cmd_str}")

    return lines


def format_agent_next_steps_section() -> list[str]:
    """Formats actionable step-by-step guidance for autonomous AI agents."""
    lines = [
        "🛠️  \033[1mAI Agent Remediation Protocol\033[0m:",
        "  1. Inspect Errors : Call view_file on .ai-memory/cicd/errors.log (or read errors.json)",
        "  2. Surgical Fix   : Edit suspect files complying with 02-spec/02-coding-guidelines/",
        "  3. Single Re-Test : Run targeted filter command above to confirm local fix",
        "  4. Suite Green    : Run python 03-ai-scripts/06-cicd-local-runner.py (unchanged gates skip in ~0.5ms)",
        "\033[1;91m" + "=" * 70 + "\033[0m",
    ]

    return lines


def format_single_trace_metadata(idx: int, total: int, err: dict[str, Any]) -> list[str]:
    """Formats metadata lines for a failed gate in final summary."""
    cmd = err.get("cmd", "")
    cmd_str = " ".join(cmd) if isinstance(cmd, list) else str(cmd)
    suspects = ", ".join(err.get("suspect_files", [])) or "(None detected)"

    return [
        f"\033[1;91m[{idx}/{total}] FAIL: {err.get('name', '')}\033[0m",
        f"  Command      : {cmd_str}",
        f"  Exit Code    : {err.get('code', '')} ({format_duration(err.get('elapsed', 0.0))})",
        f"  Suspect Files: {suspects}",
        "  --- Stack Trace & Error Output ---",
    ]


def format_single_trace_summary(idx: int, total: int, err: dict[str, Any]) -> list[str]:
    """Formats metadata and full stack trace for a single failed gate."""
    meta = format_single_trace_metadata(idx, total, err)
    trace = err.get("error", "No output captured.")

    return [*meta, trace, "----------------------------------------------------------------"]


def format_failing_traces_summary_section(errors_list: list[dict[str, Any]]) -> list[str]:
    """Summarizes each failing test along with its full stack trace and output."""
    if not errors_list:
        return []
    total = len(errors_list)
    lines = [
        "💥 \033[1;91mFAILING TESTS & FULL STACK TRACES SUMMARY\033[0m:",
        "================================================================",
    ]
    for i, err in enumerate(errors_list, 1):
        lines.extend(format_single_trace_summary(i, total, err))

    return lines


def format_ai_remediation_banner(state: dict[str, Any], session_dir: Path | None) -> str:
    """Assembles the full post-execution remediation banner for terminal output."""
    errs = state.get("errors_list", [])
    parts = [
        *format_remediation_banner_header(len(errs)), "",
        *format_log_locations_section(session_dir), "",
        *format_suspect_files_summary(errs), "",
        *format_failing_traces_summary_section(errs), "",
        *format_retest_commands_section(errs), "",
        *format_agent_next_steps_section(),
    ]

    return "\n".join(parts)


def print_failure_report(report: str, failed_count: int, state: dict[str, Any], session_dir: Path | None) -> int:
    """Prints execution table, failure summary, and AI remediation banner."""
    print(report)
    banner = format_ai_remediation_banner(state, session_dir)
    print(banner)
    print(f"\033[1;91m[FAILURE]\033[0m CI/CD quality gates failed with {failed_count} error(s).\n")

    return 1


def print_success_report(report: str, show_all: bool, total: int, cached: int, elapsed: float) -> int:
    """Prints success confirmation and returns exit code 0."""
    if show_all:
        print("\n" + report)
        print("\n\033[1;92m🎉 All quality gates passed successfully! Codebase is 100% green.\033[0m")
    else:
        executed = total - cached
        print(f"✔ All passed. ({total} gates [{cached} cached, {executed} executed] in {format_duration(elapsed)})")

    return 0


def handle_text_output(
    args: argparse.Namespace, results: list[JobResult], counts: tuple[int, int, int, int],
    total_elapsed: float, state: dict[str, Any], session_dir: Path | None
) -> int:
    """Renders text terminal report and writes to output file if requested."""
    total, passed, failed, timeout = counts
    cached = sum(1 for r in results if r.is_cached)
    if failed > 0 or timeout > 0:
        rep = format_full_report(results, total, passed, failed, timeout, total_elapsed, show_all=False)

        return print_failure_report(rep, failed + timeout, state, session_dir)
    rep = format_full_report(results, total, passed, failed, timeout, total_elapsed, show_all=True)

    return print_success_report(rep, args.show_all, total, cached, total_elapsed)


def setup_runner_state(total_jobs: int, is_json: bool, curr_head: str, curr_dirty: dict) -> dict[str, Any]:
    """Builds the initial state dictionary for runner execution."""
    state = {
        "counter": 0,
        "total": total_jobs,
        "results": [],
        "is_json": is_json,
        "head": curr_head,
        "dirty": curr_dirty,
        "gates": {},
        "freshly_executed": set(),
        "errors_list": [],
    }

    return state


def prepare_runner_temp_dirs(args: argparse.Namespace) -> None:
    """Cleans temporary directories and optionally Go cache if requested."""
    clear_repo_build_temp()
    clear_repo_test_temp()
    if getattr(args, "force_run", False) or getattr(args, "clean_gocache", False):
        clear_repo_go_cache()
    clear_stale_failures_log()
    clean_stale_temp_artifacts()
    prune_old_cicd_runs(keep_count=5)
    clean_coverage_artifacts()


def prepare_runner_context(args: argparse.Namespace, root: Path, total_jobs: int) -> tuple[Path, dict, set, TelemetryTracker, dict]:
    """Prepares directories, git delta, and initial state machine."""
    prepare_runner_temp_dirs(args)
    sdir, prev = init_session_scaffolding(args.force_run, args.resume_mode)
    curr_head, curr_dirty = get_head_commit(root), get_dirty_files_map(root)
    delta = compute_repo_delta(root, prev.get("head", ""), curr_head, prev.get("dirty", {}), curr_dirty)
    hb_interval = getattr(args, "heartbeat_interval", DEFAULT_HEARTBEAT_INTERVAL)
    telemetry = TelemetryTracker(total_jobs, sys.stdout.isatty(), bool(args.json_mode), args.show_all, heartbeat_interval=hb_interval)
    state = setup_runner_state(total_jobs, bool(args.json_mode), curr_head, curr_dirty)
    update_cicd_summary(state, sdir, is_finished=False)
    emit_telemetry_event("run_started", sdir, {"total_gates": total_jobs, "session": sdir.name})

    return sdir, prev, delta, telemetry, state


def execute_agent_group(
    agent_name: str,
    group_batches: list[dict[str, Any]],
    workers: int,
    is_sync: bool,
    args: argparse.Namespace,
    st: dict[str, Any],
    prev: dict[str, Any],
    delta: set[str],
    root: Path,
    sdir: Path,
    tel: TelemetryTracker,
) -> None:
    """Runs an ordered sequence of batches assigned to an autonomous pipeline section agent."""
    for b in group_batches:
        execute_job_batch(b, workers, is_sync, args, st, prev, delta, root, sdir, tel)


def run_parallel_agent_pipeline(
    batches: list[dict],
    workers: int,
    args: argparse.Namespace,
    st: dict,
    prev: dict,
    delta: set,
    root: Path,
    sdir: Path,
    tel: TelemetryTracker,
) -> None:
    """Runs autonomous pipeline section agents concurrently across logical CPU cores."""
    batch_map = {b["name"]: b for b in batches}

    agent1_batches = [b for b in [batch_map.get("Linters & AST Checks")] if b]
    agent2_batches = [
        b
        for b in [
            batch_map.get("Compile Gates"),
            batch_map.get("Packaging Gates"),
            batch_map.get("E2E Smoke Tests"),
        ]
        if b
    ]
    agent3_batches = [
        b
        for b in [
            batch_map.get("Smart Unit Tests & Coverage"),
            batch_map.get("Coverage Verification"),
            batch_map.get("Race Detection"),
        ]
        if b
    ]

    handled = {
        "Linters & AST Checks",
        "Compile Gates",
        "Packaging Gates",
        "E2E Smoke Tests",
        "Smart Unit Tests & Coverage",
        "Coverage Verification",
        "Race Detection",
    }
    other_batches = [b for b in batches if b["name"] not in handled]

    agent_groups = [
        ("Agent 1 (Linters & AST Checks)", agent1_batches),
        ("Agent 2 (Build & Packaging & Smoke Suite)", agent2_batches),
        ("Agent 3 (Unit Tests & Coverage & Race)", agent3_batches),
    ]
    if other_batches:
        agent_groups.append(("Agent 4 (Other Gates)", other_batches))

    active_groups = [(name, grp) for name, grp in agent_groups if grp]
    if len(active_groups) <= 1:
        for _, grp in active_groups:
            execute_agent_group("Sequential", grp, workers, False, args, st, prev, delta, root, sdir, tel)

        return

    for name, grp in active_groups:
        execute_agent_group(name, grp, workers, False, args, st, prev, delta, root, sdir, tel)


def run_batch_sequence(
    batches: list[dict],
    args: argparse.Namespace,
    st: dict,
    prev: dict,
    delta: set,
    root: Path,
    sdir: Path,
    tel: TelemetryTracker,
) -> None:
    """Executes all enqueued job batches with autonomous pipeline section agents."""
    workers = 1 if args.sync_mode else max(1, args.workers)
    try:
        if args.sync_mode or len(batches) <= 1:
            for b in batches:
                execute_job_batch(b, workers, args.sync_mode, args, st, prev, delta, root, sdir, tel)
        else:
            run_parallel_agent_pipeline(batches, workers, args, st, prev, delta, root, sdir, tel)
    finally:
        tel.stop_heartbeat()
        tel.clear_line()


def extract_runner_counts(total: int, results: list[JobResult]) -> tuple[int, int, int, int]:
    """Computes tuple of total passed failed and timeout counts."""
    pass_cnt = sum(1 for r in results if r.is_success)
    fail_cnt = sum(1 for r in results if not r.is_success and not r.is_timeout)
    time_cnt = sum(1 for r in results if r.is_timeout)

    return total, pass_cnt, fail_cnt, time_cnt


def show_upfront_plan_if_text(args: argparse.Namespace, batches: list[dict[str, Any]]) -> None:
    """Displays queued segments and gates banner when running in text mode."""
    if not args.json_mode:
        label = "Synchronous (1 worker)" if args.sync_mode else f"Parallel ({args.workers} workers, 3 Autonomous Section Agents)"
        print_queued_segments_and_gates(batches, label)


def emit_runner_report(args: argparse.Namespace, st: dict, counts: tuple, elapsed: float, sdir: Path) -> int:
    """Emits final runner report in either JSON or human-readable format."""
    if args.json_mode:
        has_failed = bool(counts[2] > 0 or counts[3] > 0)

        return handle_json_output(args, st["results"], counts, elapsed, has_failed)

    return handle_text_output(args, st["results"], counts, elapsed, st, sdir)


def cleanup_runner_temp_artifacts() -> None:
    """Performs post-runner temporary artifact cleaning."""
    clear_repo_test_temp()
    clean_stale_temp_artifacts()
    clean_coverage_artifacts()
    prune_old_cicd_runs(keep_count=5)


def execute_runner(args: argparse.Namespace, active_batches: list[dict[str, Any]], repo_root: Path) -> int:
    """Orchestrates test batch execution and report output generation."""
    total_jobs = sum(len(b["jobs"]) for b in active_batches)
    if total_jobs == 0:
        return 0
    show_upfront_plan_if_text(args, active_batches)
    sdir, prev, delta, tel, st = prepare_runner_context(args, repo_root, total_jobs)
    start = time.monotonic()
    run_batch_sequence(active_batches, args, st, prev, delta, repo_root, sdir, tel)
    elapsed = round(time.monotonic() - start, 2)
    persist_state(sdir, st)
    update_cicd_summary(st, sdir, is_finished=True)
    cleanup_runner_temp_artifacts()

    return emit_runner_report(args, st, extract_runner_counts(total_jobs, st["results"]), elapsed, sdir)


def load_cicd_timings(path: Path) -> dict[str, float]:
    """Loads historical CI/CD job timings from JSON file if available."""
    if not path.exists():
        return {}
    try:
        data = json.loads(path.read_text(encoding=DEFAULT_ENCODING))
        if isinstance(data, dict):
            return {k: float(v) for k, v in data.items()}
    except (json.JSONDecodeError, OSError) as err:
        sys.stderr.write(f"[WARN] Failed to load timings: {err}\n")

    return {}


def save_cicd_timings(path: Path, timings: dict[str, float]) -> None:
    """Persists recorded CI/CD job timings to the timing file."""
    path.parent.mkdir(parents=True, exist_ok=True)
    try:
        path.write_text(json.dumps(timings, indent=2), encoding=DEFAULT_ENCODING)
    except OSError as err:
        sys.stderr.write(f"[WARN] Failed to save timings: {err}\n")


def calculate_total_eta(active_batches: list[dict[str, Any]], timings: dict[str, float]) -> int:
    """Computes total estimated execution time across all active batches using historical and inventory data."""
    total_sec = 0.0
    for batch in active_batches:
        batch_sum = 0.0
        batch_max = 0.0
        for job_name in batch.get("jobs", {}):
            if job_name == "Go Smart Incremental Tests":
                inv = load_raw_test_inventory(TEST_INVENTORY_PATH)
                dirty_tests = [t for t in inv.get("tests", {}).values() if t.get("needs_run", True)]
                slow_threshold = float(os.environ.get("GITMAP_SLOW_TEST_THRESHOLD", "4.0"))
                slow_t = [t for t in dirty_tests if t.get("tier") in ("slow", "heavy") or t.get("is_slow", False) or float(t.get("duration_sec", 0.0)) >= slow_threshold]
                fast_t = [t for t in dirty_tests if t not in slow_t]
                job_est = (sum(float(t.get("duration_sec", 4.0)) for t in slow_t) / 8.0) + (sum(float(t.get("duration_sec", 0.005)) for t in fast_t) / 16.0)
                job_est = max(5.0, round(job_est, 1))
            else:
                job_est = timings.get(job_name, DEFAULT_JOB_ESTIMATE_SEC)
            batch_sum += job_est
            batch_max = max(batch_max, job_est)
        workers = batch.get("workers", DEFAULT_WORKERS)
        batch_est = max(batch_max, batch_sum / workers) if workers > 0 else batch_sum
        total_sec += batch_est

    return int(total_sec)


def run_eta_worker(interval_sec: int, total_est_sec: int, start_time: float, stop_event: threading.Event) -> None:
    """Background worker reporting remaining ETA every interval and updating live telemetry."""
    while not stop_event.is_set():
        if stop_event.wait(timeout=interval_sec):
            break
        elapsed = time.time() - start_time
        remaining = max(1, int(total_est_sec - elapsed))
        print(f"\n[ETA] Estimated remaining time: {format_duration(remaining)}\n", flush=True)
        if RUNNER_ETA_FILE.is_file():
            try:
                eta_data = json.loads(RUNNER_ETA_FILE.read_text(encoding="utf-8"))
                eta_data["remaining_eta_sec"] = remaining
                eta_data["elapsed_sec"] = round(elapsed, 1)
                atomic_write_json(RUNNER_ETA_FILE, eta_data)
            except Exception:
                pass


def start_eta_reporter(interval_sec: int, total_est_sec: int, stop_event: threading.Event) -> threading.Thread | None:
    """Spawns background ETA daemon thread if interval > 0."""
    if interval_sec <= 0:
        return None
    worker = threading.Thread(target=run_eta_worker, args=(interval_sec, total_est_sec, time.time(), stop_event), daemon=True)
    worker.start()

    return worker


def cleanup_pipeline_post_run(args: argparse.Namespace, timings: dict) -> None:
    """Performs post-pipeline timing persistence and temporary file cleanup."""
    timings.update(GLOBAL_TIMINGS)
    save_cicd_timings(CACHE_TIMINGS_FILE, timings)
    clear_repo_test_temp()
    if getattr(args, "clean_gocache", False):
        clear_repo_go_cache()
    clean_stale_temp_artifacts()
    clean_coverage_artifacts()
    prune_old_cicd_runs(keep_count=5)


def run_pipeline_with_eta(args: argparse.Namespace, batches: list, root: Path, timings: dict, total_est: int) -> int:
    """Runs test execution pipeline with background ETA reporter."""
    stop_event = threading.Event()
    start_eta_reporter(args.eta_interval, total_est, stop_event)
    try:
        return execute_runner(args, batches, root)
    finally:
        stop_event.set()
        cleanup_pipeline_post_run(args, timings)


def generate_git_changed_manifest_fallback(repo_root: Path, commits: int) -> None:
    """Generates git-changed-files.json using git diff when extractor script is missing."""
    manifest_path = repo_root / ".ai-memory" / "temp" / "git-changed-files.json"
    manifest_path.parent.mkdir(parents=True, exist_ok=True)
    head_sha = get_head_commit(repo_root)
    dirty_files = get_dirty_files_map(repo_root)
    cmd = ["git", "diff", "--name-only", f"HEAD~{commits}..HEAD"]
    res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, check=False)
    diff_lines = [p.strip() for p in res.stdout.splitlines() if p.strip()] if res.returncode == 0 else []
    all_paths = set(diff_lines) | set(dirty_files.keys())
    records = [{"path": p, "status": "working-tree" if p in dirty_files else "committed", "extension": Path(p).suffix} for p in sorted(all_paths)]
    payload = {"last_commit_hash": head_sha, "commit_window": commits, "total_files": len(records), "files": records}
    atomic_write_json(manifest_path, payload)


def run_extractor_if_present(extractor: Path, commits: int, is_force: bool, repo_root: Path) -> bool:
    """Executes extractor script if present and returns whether execution succeeded."""
    if not extractor.is_file():
        return False
    cmd = [sys.executable, str(extractor), "--commits", str(commits), "--quiet"]
    if is_force:
        cmd.append("--no-incremental")
    res = subprocess.run(cmd, cwd=str(repo_root), check=False)

    return res.returncode == 0


def ensure_manifest_if_changed_only(args: argparse.Namespace, repo_root: Path) -> None:
    """Pre-generates git-changed-files.json when running in changed-only mode."""
    if not getattr(args, "changed_only", False):
        return
    commits = getattr(args, "commits", 20)
    is_force = bool(getattr(args, "force_run", False))
    extractor = repo_root / "03-ai-scripts" / "27-git-changed-files.py"
    has_run = run_extractor_if_present(extractor, commits, is_force, repo_root)
    if not has_run:
        generate_git_changed_manifest_fallback(repo_root, commits)


def print_cached_run_banner(elapsed: float, ttl_sec: float, data: dict[str, Any]) -> None:
    """Prints formatted terminal banner for cached debounce run."""
    print("================================================================")
    print("Here is the result from the previous run.")
    print("================================================================")
    print(f"⏱️  Cached from previous run {elapsed:.1f}s ago (debounce TTL: {ttl_sec:.0f}s).")
    status_text = "PASSED (exit 0)" if data.get("exit_code") == 0 else f"FAILED (exit {data.get('exit_code')})"
    print(f"📋 Status: {status_text}")
    if data.get("summary"):
        print(f"📊 Summary: {data.get('summary')}")
    print("================================================================")


def check_recent_run_cache(cache_file: Path, signature: str, is_force: bool = False, normal_ttl: float = 15.0, force_ttl: float = 5.0) -> int | None:
    """Returns cached exit code if an identical runner run occurred within the debounce window."""
    if not cache_file.exists():
        return None
    ttl_sec = force_ttl if is_force else normal_ttl
    try:
        data = json.loads(cache_file.read_text(encoding=DEFAULT_ENCODING))
        elapsed = time.time() - float(data.get("timestamp", 0.0))
        if elapsed < ttl_sec and data.get("signature") == signature:
            print_cached_run_banner(elapsed, ttl_sec, data)
            return int(data.get("exit_code", 0))
    except Exception:
        return None
    return None


def save_recent_run_cache(cache_file: Path, signature: str, exit_code: int, summary: str = "") -> None:
    """Persists recent run result to debounce cache."""
    try:
        cache_file.parent.mkdir(parents=True, exist_ok=True)
        data = {
            "timestamp": time.time(),
            "signature": signature,
            "exit_code": exit_code,
            "summary": summary,
        }
        cache_file.write_text(json.dumps(data, indent=2), encoding=DEFAULT_ENCODING)
    except Exception:
        pass


def handle_cache_cli_commands(args: argparse.Namespace) -> bool:
    """Handles cache management CLI flags and exits early if matched."""
    if getattr(args, "cache_stats", False):
        display_cache_stats()
        return True
    if getattr(args, "clear_cache", False):
        clear_cicd_cache()
        return True
    if getattr(args, "prune_cache", False):
        prune_cicd_cache()
        return True

    return False


def filter_out_test_batches(batches: list[dict[str, Any]]) -> list[dict[str, Any]]:
    """Filters out all test jobs when --no-tests flag is provided."""
    filtered = []
    for b in batches:
        if is_test_batch(b.get("name", "")):
            continue
        jobs = {k: v for k, v in b.get("jobs", {}).items() if not is_test_job(k, v)}
        if jobs:
            b_copy = dict(b)
            b_copy["jobs"] = jobs
            filtered.append(b_copy)

    return filtered


def resolve_filtered_pipeline_batches(args: argparse.Namespace) -> list[dict[str, Any]]:
    """Resolves and filters job batches based on package and job filters."""
    if (getattr(args, "package_filter", None) or getattr(args, "smart_mode", False)) and not args.filter:
        batches = filter_job_batches(JOB_BATCHES, "Go Smart Incremental Tests")
    else:
        batches = filter_job_batches(JOB_BATCHES, args.filter)
    batches = filter_batches_by_host_os(batches, getattr(args, "all_os", False))
    if getattr(args, "no_tests", False):
        return filter_out_test_batches(batches)

    return batches


def check_debounce_and_manifest(args: argparse.Namespace, sig: str, root: Path) -> None:
    """Checks debounce cache and pre-generates changed manifest if requested."""
    cached_code = check_recent_run_cache(CACHE_DEBOUNCE_FILE, sig, is_force=bool(getattr(args, "force_run", False)))
    if cached_code is not None:
        sys.exit(cached_code)
    ensure_manifest_if_changed_only(args, root)


def execute_pipeline_entry(args: argparse.Namespace, root: Path, batches: list) -> int:
    """Loads timings and runs execution pipeline with ETA calculation."""
    timings = load_cicd_timings(CACHE_TIMINGS_FILE)

    return run_pipeline_with_eta(args, batches, root, timings, calculate_total_eta(batches, timings))


def check_inventory_only_mode(args: argparse.Namespace, root: Path) -> None:
    """Updates test inventory and exits early if --inventory-only flag set."""
    inventory = build_or_update_test_inventory(root, force=bool(args.force_run))
    if getattr(args, "inventory_only", False):
        print_inventory_summary(inventory)
        sys.exit(0)


def handle_heatmap_cli(args: argparse.Namespace, root: Path) -> bool:
    """Displays test heatmap and exits if --heatmap was specified."""
    if not getattr(args, "show_heatmap", False):
        return False
    try:
        sys.path.insert(0, str(Path(__file__).parent))
        inv_gen = importlib.import_module("33-test-inventory-generator")
        inv = inv_gen.build_test_inventory(root, commits=getattr(args, "commits", 50))
        inv_gen.display_heatmap_report(inv, top_n=getattr(args, "top_n", 25))
    except Exception as exc:
        sys.stderr.write(f"[WARN] Failed to generate heatmap report: {exc}\n")
    return True


def main() -> None:
    """Primary entry point for local CI/CD quality gate runner."""
    args = parse_args()
    if handle_cache_cli_commands(args):
        sys.exit(0)
    root = Path(__file__).resolve().parent.parent
    if handle_heatmap_cli(args, root):
        sys.exit(0)
    if getattr(args, "fast_mode", False):
        os.environ["GITMAP_FAST_TESTS"] = "1"
    sig = f"no_tests={getattr(args, 'no_tests', False)},run_tests={getattr(args, 'run_tests', False)},filter={getattr(args, 'filter', '') or ''},pkg={getattr(args, 'package_filter', [])},fast={getattr(args, 'fast_mode', False)}"
    check_debounce_and_manifest(args, sig, root)
    check_inventory_only_mode(args, root)
    batches = resolve_filtered_pipeline_batches(args)
    code = execute_pipeline_entry(args, root, batches)
    save_recent_run_cache(CACHE_DEBOUNCE_FILE, sig, code, f"Executed {len(batches)} batches | Result code: {code}")
    sys.exit(code)


if __name__ == "__main__":
    main()
