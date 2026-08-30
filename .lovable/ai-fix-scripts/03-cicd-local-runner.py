#!/usr/bin/env python3
"""Auto-generated CI/CD local runner with concurrent worker pool and log aggregation.
Do not edit manually. Re-generate by running:
python .lovable/ai-fix-scripts/03-cicd-local-runner.py --rebuild
"""
import concurrent.futures
import os
import subprocess
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

# ── Configurable Variables ──────────────────────────────────────────────────
BATCH_SIZE = 3       # Maximum jobs to run concurrently per batch
JOB_TIMEOUT_SEC = 300  # Maximum seconds before a single job is timed out

ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
GITMAP_DIR = os.path.join(ROOT_DIR, "gitmap")

# ── Environment Configuration ───────────────────────────────────────────────
os.environ.setdefault("CI", "true")
os.environ.setdefault("NODE_ENV", "test")
os.environ.setdefault("GOTOOLCHAIN", "local")


def run_cmd(name, cmd_str, cwd):
    start = time.monotonic()
    try:
        res = subprocess.run(
            cmd_str,
            cwd=cwd,
            shell=True,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            timeout=JOB_TIMEOUT_SEC
        )
        elapsed = round(time.monotonic() - start, 2)
        if res.returncode == 0:
            return name, cmd_str, 0, res.stdout, res.stderr, elapsed
        else:
            return name, cmd_str, res.returncode, res.stdout, res.stderr, elapsed
    except subprocess.TimeoutExpired as e:
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd_str, "timeout", e.stdout or "", f"Job timed out after {JOB_TIMEOUT_SEC}s", elapsed
    except Exception as e:
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd_str, 1, "", str(e), elapsed


def run_generate_drift_check():
    start = time.monotonic()
    name = "Generate Drift Check"
    cmd_str = "go generate ./... & git diff"
    try:
        before_res = subprocess.run("git status --porcelain .", cwd=GITMAP_DIR, shell=True, capture_output=True, text=True, encoding="utf-8", timeout=JOB_TIMEOUT_SEC)
        res = subprocess.run("go generate ./...", cwd=GITMAP_DIR, shell=True, capture_output=True, text=True, encoding="utf-8", timeout=JOB_TIMEOUT_SEC)
        elapsed = round(time.monotonic() - start, 2)
        if res.returncode != 0:
            return name, cmd_str, res.returncode, res.stdout, f"go generate ./... failed:\n{res.stderr or res.stdout}", elapsed
        after_res = subprocess.run("git status --porcelain .", cwd=GITMAP_DIR, shell=True, capture_output=True, text=True, encoding="utf-8", timeout=JOB_TIMEOUT_SEC)
        if before_res.stdout != after_res.stdout:
            diff_res = subprocess.run("git diff .", cwd=GITMAP_DIR, shell=True, capture_output=True, text=True, encoding="utf-8", timeout=JOB_TIMEOUT_SEC)
            return name, cmd_str, 1, diff_res.stdout, f"go generate ./... modified files under gitmap/:\n{diff_res.stdout or diff_res.stderr}", elapsed
        return name, cmd_str, 0, "All generated files in gitmap/ are in sync.", "", elapsed
    except subprocess.TimeoutExpired:
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd_str, "timeout", "", f"Job timed out after {JOB_TIMEOUT_SEC}s", elapsed
    except Exception as e:
        elapsed = round(time.monotonic() - start, 2)
        return name, cmd_str, 1, "", str(e), elapsed


def main():
    print("=== CI/CD Local Runner (03-cicd-local-runner.py) ===")
    stages = [
        # Stage 1: Static policy checks (parallel)
        [
            ("Legacy Refs Check", None, "python .github/scripts/check-legacy-refs.py", ROOT_DIR),
            ("Bare Stderr Check", None, "python .github/scripts/check-bare-stderr-err.py", ROOT_DIR),
            ("Cmd Naming Check", None, "python .github/scripts/check-cmd-naming.py gitmap/cmd", ROOT_DIR),
        ],
        # Stage 2: Constants & layout checks (parallel)
        [
            ("Constants Naming Check", None, "python .github/scripts/check-constants-naming.py", ROOT_DIR),
            ("Deploy Layout Check", None, "python .github/scripts/check-deploy-layout.py", ROOT_DIR),
            ("No Golden Allow Leak Check", None, "python .github/scripts/check-no-golden-allow-leak.py", ROOT_DIR),
        ],
        # Stage 3: Linters & scripts (parallel)
        [
            ("File Size Check", None, "python .github/scripts/file-size-check.py 200", ROOT_DIR),
            ("CI Scripts Unit Tests", None, "python .github/scripts/tests/test_ci_scripts.py", ROOT_DIR),
            ("Error Management Linter", None, "python linter-scripts/check-error-management.py", ROOT_DIR),
        ],
        # Stage 4: Guidelines, Sync, and Paths (parallel)
        [
            ("Boolean & Enum Linter", None, "python linter-scripts/check-enum-and-boolean.py", ROOT_DIR),
            ("Boolean Guidelines Linter", None, "python linter-scripts/check-boolean-guidelines.py", ROOT_DIR),
            ("Nested If Linter", None, "python linter-scripts/check-nested-ifs.py", ROOT_DIR),
            ("Changelog Version Sync", None, "python .github/scripts/check-changelog-version-sync.py", ROOT_DIR),
            ("Relative Path Check", None, "python linter-scripts/check-relative-paths.py", ROOT_DIR),
        ],
        # Stage 5: Web & Docs (parallel)
        [
            ("Web Version Sync", None, "npx vitest run src/test/version-sync.test.ts", ROOT_DIR),
            ("Docs Site Build", None, "npm run build", ROOT_DIR),
            ("Spell Check (misspell)", None, "python .github/scripts/misspell-changed.py", ROOT_DIR),
        ],
        # Stage 6: Code Generation Drift Check (sequential)
        [
            ("Generate Drift Check", run_generate_drift_check, None, None),
        ],
        # Stage 7: Go Vet & Smoke (parallel - separate workdirs)
        [
            ("Go Vet", None, "go vet -p=2 ./...", GITMAP_DIR),
            ("Installer Smoke", None, "python .github/scripts/smoke-installer.py source", ROOT_DIR),
        ],
        # Stage 8: CLI Zero-Args Smoke & Command Invocations (sequential)
        [
            ("CLI Zero-Args Smoke", None, "go run .", GITMAP_DIR),
            ("Command Invocations Gate", None, 'go test -p=2 -v -run "^(TestInstallCommandInvocations|TestPipelineCommandInvocations|TestTopLevelErrorLogsAndLogs|TestHelpFlagTrigger|TestBuildErrorLogsPayloadWithFallback)$" ./cmd', GITMAP_DIR),
        ],
        # Stage 9: Compile Gate (sequential)
        [
            ("Compile Gate", None, "go test -run=^$ -p=2 ./... -count=1", GITMAP_DIR),
        ],
        # Stage 10: Full Suite Linter (sequential)
        [
            ("Full Suite Lint", None, "golangci-lint run --concurrency=2 ./...", GITMAP_DIR),
        ],
    ]

    all_checks = [item for stage in stages for item in stage]
    total_jobs = len(all_checks)
    print(f"[INFO] Enqueued {total_jobs} quality gates across worker pool stages...\n")

    all_results = {}
    total_start = time.monotonic()

    for stage_idx, stage in enumerate(stages, 1):
        with concurrent.futures.ThreadPoolExecutor(max_workers=min(BATCH_SIZE, len(stage))) as executor:
            futures = {}
            for name, fn, cmd, cwd in stage:
                if fn is not None:
                    futures[executor.submit(fn)] = name
                else:
                    futures[executor.submit(run_cmd, name, cmd, cwd)] = name

            for future in concurrent.futures.as_completed(futures):
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
                    all_results[job_name] = (1, "", str(ex), 0, "")
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
        report_failures(failed_jobs, timeout_jobs)
        print(f"\n[FAILURE] CI/CD quality gates failed with {len(failed_jobs) + len(timeout_jobs)} error(s).")
        sys.exit(1)

    print(f"\n[SUCCESS] All {total_jobs} CI/CD quality gates passed (exit 0)!")
    sys.exit(0)


def print_failure_log(name, elapsed, out, err, cmd):
    print(f"\n[FAILURE LOG] Job: {name} (Duration: {elapsed}s)")
    print(f"Command: {cmd}")
    if out.strip():
        print(f"Stdout:\n{out.strip()}")
    if err.strip():
        print(f"Stderr:\n{err.strip()}")
    print("-" * 60)


def print_timeout_log(name, elapsed, err, cmd):
    print(f"\n[TIMEOUT LOG] Job: {name} (Duration: {elapsed}s)")
    print(f"Command: {cmd}")
    print(f"Reason: {err}")
    print("-" * 60)


def report_failures(failed_jobs, timeout_jobs):
    print("Detailed Failure Logs:")
    print("-" * 60)
    for name, elapsed, out, err, cmd in failed_jobs:
        print_failure_log(name, elapsed, out, err, cmd)
    for name, elapsed, err, cmd in timeout_jobs:
        print_timeout_log(name, elapsed, err, cmd)


if __name__ == "__main__":
    main()
