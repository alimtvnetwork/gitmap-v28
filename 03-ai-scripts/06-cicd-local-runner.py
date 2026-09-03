#!/usr/bin/env python3
"""Auto-generated CI/CD local runner with concurrent worker pool and log aggregation.
Do not edit manually. Re-generate by running:
python 03-ai-scripts/06-cicd-local-runner.py --rebuild
"""
from concurrent.futures import ThreadPoolExecutor, as_completed
import os
import shutil
import subprocess
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

# ── Configurable Variables ──────────────────────────────────────────────────
BATCH_SIZE      = 3    # Number of jobs to run concurrently (round-robin worker pool)
JOB_TIMEOUT_SEC = 300  # Maximum seconds before a single job is timed out

# ── Environment Configuration ───────────────────────────────────────────────
os.environ.setdefault("CI", "true")
os.environ.setdefault("NODE_ENV", "test")

# ── Job Definitions (Partitioned into Order-Dependent Sequential Batches) ───
JOB_BATCHES: list[dict[str, list[str]]] = [
    # Batch 1: Linters & AST Checks (order-independent, run concurrently)
    {
        "Spell Check (misspell)": [sys.executable, ".github/scripts/misspell-changed.py"],
        "Nested If Linter": [sys.executable, "linter-scripts/check-nested-ifs.py"],
        "Boolean & Enum Linter": [sys.executable, "linter-scripts/check-enum-and-boolean.py"],
        "Error Management Check": [sys.executable, "linter-scripts/check-error-management.py"],
        "Relative Path Check": [sys.executable, "linter-scripts/check-relative-paths.py"],
        "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
        "CLI Help Parity Check": [sys.executable, "03-ai-scripts/09-cli-help-auditor.py"],
        "Constants Registry AST Check": ["go", "test", "-C", "gitmap", "./constants/...", "-run", "TestTopLevelCmdRegistryMatchesAST", "-count=1"],
        "Constants Collision Check": ["go", "test", "-C", "gitmap", "./constants/...", "-run", "TestTopLevelCmdConstantsAreUnique", "-count=1"],
        "Helptext Parity Check": ["go", "test", "-C", "gitmap", "./helptext/...", "-count=1"],
    },
    # Batch 2: Compile & Packaging Gates
    {
        "Go Compile Gate": ["go", "build", "-C", "gitmap", "-o", "../bin/gitmap.exe", "."],
        "Web App Build": ["npm", "run", "build"],
    },
    # Batch 3: E2E Smoke Tests (depends on Go Compile Gate binary)
    {
        "E2E Smoke Suite": [sys.executable, ".github/scripts/e2e-cli-smoke.py", "bin/gitmap.exe"],
    },
]

def run_job(name: str, cmd: list[str]) -> tuple[str, list[str], int | str, str, str, float]:
    start = time.monotonic()
    resolved_cmd = list(cmd)
    binary_path = shutil.which(cmd[0])

    if binary_path is not None:
        resolved_cmd[0] = binary_path

    try:
        result = subprocess.run(
            resolved_cmd,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            timeout=JOB_TIMEOUT_SEC,
        )
        elapsed = round(time.monotonic() - start, 2)

        return name, cmd, result.returncode, result.stdout, result.stderr, elapsed
    except subprocess.TimeoutExpired as e:
        elapsed = round(time.monotonic() - start, 2)

        return name, cmd, "timeout", e.stdout or "", f"Job timed out after {JOB_TIMEOUT_SEC}s", elapsed
    except Exception as e:
        elapsed = round(time.monotonic() - start, 2)

        return name, cmd, 1, "", str(e), elapsed

def main() -> None:
    total_jobs = sum(len(b) for b in JOB_BATCHES)
    print(f"[INFO] Enqueued {total_jobs} quality gates across {BATCH_SIZE} concurrent workers ({len(JOB_BATCHES)} sequential batches)...\n")

    all_results = {}
    total_start = time.monotonic()

    # Execute each batch sequentially; jobs within a batch execute concurrently
    for batch_idx, batch_jobs in enumerate(JOB_BATCHES, start=1):
        job_items = list(batch_jobs.items())

        with ThreadPoolExecutor(max_workers=BATCH_SIZE) as executor:
            futures = {executor.submit(run_job, name, cmd): name for name, cmd in job_items}

            for future in as_completed(futures):
                try:
                    name, cmd, code, out, err, elapsed = future.result()
                    all_results[name] = (code, out, err, elapsed, cmd)

                    if code == 0:
                        print(f"  PASS [{name}] ({elapsed}s)")
                    elif code == "timeout":
                        print(f"  TIMEOUT [{name}] ({elapsed}s)")
                    else:
                        print(f"  FAIL [{name}] ({elapsed}s)")
                except Exception as ex:
                    job_name = futures[future]
                    all_results[job_name] = (1, "", str(ex), 0, batch_jobs.get(job_name, []))
                    print(f"  FAIL [{job_name}] (Exception: {ex})")

    total_elapsed = round(time.monotonic() - total_start, 2)

    # ── Final Consolidated Summary Report ──────────────────────────────────
    print("\n" + "=" * 60)
    print("           CI/CD EXECUTION SUMMARY REPORT")
    print("=" * 60)

    passed_jobs = []
    failed_jobs = []
    timeout_jobs = []

    for name, (code, out, err, elapsed, cmd) in all_results.items():
        if code == 0:
            passed_jobs.append((name, elapsed))
        elif code == "timeout":
            timeout_jobs.append((name, elapsed, err, cmd))
        else:
            failed_jobs.append((name, elapsed, out, err, cmd))

    print(f"Total: {total_jobs} | Passed: {len(passed_jobs)} | Failed: {len(failed_jobs)} | Timeouts: {len(timeout_jobs)} | Time: {total_elapsed}s\n")

    if failed_jobs or timeout_jobs:
        print("Detailed Failure Logs:")
        print("-" * 60)

        for name, elapsed, out, err, cmd in failed_jobs:
            print(f"\n[FAILURE LOG] Job: {name} (Duration: {elapsed}s)")
            print(f"Command: {' '.join(cmd)}")

            if out.strip():
                print(f"Stdout:\n{out.strip()}")

            if err.strip():
                print(f"Stderr:\n{err.strip()}")

            print("-" * 60)

        for name, elapsed, err, cmd in timeout_jobs:
            print(f"\n[TIMEOUT LOG] Job: {name} (Duration: {elapsed}s)")
            print(f"Command: {' '.join(cmd)}")
            print(f"Reason: {err}")
            print("-" * 60)

        print(f"\n[FAILURE] CI/CD quality gates failed with {len(failed_jobs) + len(timeout_jobs)} error(s).")
        sys.exit(1)
    else:
        print(f"\n[SUCCESS] All {total_jobs} CI/CD quality gates passed (exit 0)!")
        sys.exit(0)

if __name__ == "__main__":
    main()
