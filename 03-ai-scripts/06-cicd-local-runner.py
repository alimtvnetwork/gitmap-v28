#!/usr/bin/env python3
"""
Fast Multi-Threaded Local CI/CD Runner
Executes repository quality gates in parallel using ThreadPoolExecutor and enforces zero-failure tolerance.

All Enums, Constants, and Utility Functions are imported directly from 02-shared-engine.py.
"""

from concurrent.futures import ThreadPoolExecutor
from importlib import import_module
from pathlib import Path
import subprocess
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
DEFAULT_MAX_WORKERS = engine.DEFAULT_MAX_WORKERS
CI_JOBS_MATRIX = engine.CI_JOBS_MATRIX
format_keys = engine.format_keys

def execute_ci_job(job_name: str, command: list[str]) -> tuple[str, bool, str]:
    """Executes a single validation check asynchronously."""
    try:
        res = subprocess.run(command, capture_output=True, text=True, encoding=DEFAULT_ENCODING, errors="replace")
        is_success = (res.returncode == 0)
        if is_success:
            return (job_name, True, res.stdout)

        return (job_name, False, res.stdout + LINE_SEPARATOR + res.stderr)
    except Exception as e:
        return (job_name, False, str(e))

def run_pipeline(
    jobs: dict[str, list[str]] | None = None,
    max_workers: int = DEFAULT_MAX_WORKERS
) -> int:
    """Dispatches all jobs concurrently and prints clean summary report."""
    target_jobs = jobs or CI_JOBS_MATRIX
    enqueued_names = format_keys(target_jobs)

    print("🚀 Running Local CI/CD Pipeline via ThreadPoolExecutor...")
    print(f"📋 Enqueued Jobs: {enqueued_names}{LINE_SEPARATOR}")

    results = []
    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(execute_ci_job, name, cmd) for name, cmd in target_jobs.items()]
        for f in futures:
            results.append(f.result())

    has_failures = False
    print(f"{LINE_SEPARATOR}================ FINAL SUMMARY ================")
    for name, is_success, log in results:
        if is_success:
            print(f"✅ {name}: PASSED")
        else:
            print(f"❌ {name}: FAILED")
            has_failures = True
            print(f"--- {name} LOG ---")
            print(log.strip())
            print(f"--------------------{LINE_SEPARATOR}")

    if has_failures:
        print(f"{LINE_SEPARATOR}❌ Pipeline failed.")
        return ExitCodeType.VIOLATIONS_FOUND.value

    print("🎉 All jobs passed successfully!")
    return ExitCodeType.SUCCESS.value

def main():
    sys.exit(run_pipeline())

if __name__ == "__main__":
    main()
