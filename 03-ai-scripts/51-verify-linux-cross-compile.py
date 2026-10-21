#!/usr/bin/env python3
"""51-verify-linux-cross-compile.py — Cross-platform Go Linux verification suite."""

from __future__ import annotations

import argparse
from dataclasses import asdict, dataclass
from importlib import import_module
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")
ExitCodeType = engine.ExitCodeType


@dataclass
class GateResult:
    """Represents the verification result of a single cross-compilation gate."""
    name: str
    command: str
    is_success: bool
    duration_sec: float
    output: str
    error: str


def resolve_repository_root() -> Path:
    """Returns repository root directory."""
    script_dir = Path(__file__).resolve().parent

    return script_dir.parent


def resolve_cli_directory() -> Path:
    """Returns cli module directory relative to repository root."""
    root = resolve_repository_root()

    return root / "cli"


def check_tool_version(cmd: list[str]) -> str:
    """Executes version check command returning trimmed output or empty string."""
    try:
        proc = subprocess.run(cmd, capture_output=True, text=True, check=False)
        if proc.returncode == 0:
            return proc.stdout.strip()
    except (FileNotFoundError, OSError):
        pass

    return ""


def check_is_docker_available() -> bool:
    """Checks whether docker executable is available on PATH."""
    has_bin = shutil.which("docker") is not None

    return has_bin


def run_cross_compile_gate(name: str, cmd: list[str], extra_env: dict[str, str], work_dir: Path) -> GateResult:
    """Executes a single Go cross-compilation gate and records timing and output."""
    env = os.environ.copy()
    env.update(extra_env)
    start = time.perf_counter()
    proc = subprocess.run(cmd, cwd=str(work_dir), env=env, capture_output=True, text=True, check=False)
    elapsed = time.perf_counter() - start
    is_success = (proc.returncode == 0)

    return GateResult(
        name=name,
        command=" ".join(cmd),
        is_success=is_success,
        duration_sec=round(elapsed, 3),
        output=proc.stdout.strip(),
        error=proc.stderr.strip(),
    )


def build_linux_cross_env() -> dict[str, str]:
    """Constructs environment mapping for pure Go Linux cross-compilation."""
    return {
        "GOOS": "linux",
        "GOARCH": "amd64",
        "CGO_ENABLED": "0",
    }


def execute_build_gate(cli_dir: Path, temp_dir: Path) -> GateResult:
    """Executes Linux binary cross-compilation gate."""
    out_path = temp_dir / "gitmap-linux"
    cmd = ["go", "build", "-o", str(out_path), "."]

    return run_cross_compile_gate("Linux Go Build", cmd, build_linux_cross_env(), cli_dir)


def execute_test_compile_gate(cli_dir: Path, temp_dir: Path) -> GateResult:
    """Executes Linux test compilation gate for cmdschedule."""
    out_path = temp_dir / "schedule.test"
    cmd = ["go", "test", "-c", "./cmdschedule", "-o", str(out_path)]

    return run_cross_compile_gate("Linux Test Compile (cmdschedule)", cmd, build_linux_cross_env(), cli_dir)


def execute_vet_gate(cli_dir: Path) -> GateResult:
    """Executes Linux go vet gate across cmdschedule and cmdmail."""
    cmd = ["go", "vet", "./cmdschedule", "./cmdmail"]

    return run_cross_compile_gate("Linux Go Vet (cmdschedule, cmdmail)", cmd, {"GOOS": "linux"}, cli_dir)


def execute_all_gates(cli_dir: Path) -> list[GateResult]:
    """Runs all required Linux cross-compilation gates sequentially in a temp directory."""
    results: list[GateResult] = []
    with tempfile.TemporaryDirectory(prefix="gitmap-linux-gate-") as tmp:
        temp_dir = Path(tmp)
        results.append(execute_build_gate(cli_dir, temp_dir))
        results.append(execute_test_compile_gate(cli_dir, temp_dir))
        results.append(execute_vet_gate(cli_dir))

    return results


def is_all_successful(results: list[GateResult]) -> bool:
    """Determines if all gates succeeded."""
    has_failed = any(not r.is_success for r in results)

    return not has_failed


def print_single_gate_report(res: GateResult) -> None:
    """Prints status line and optional error details for a single gate."""
    status_icon = "✔" if res.is_success else "✖"
    print(f"  {status_icon} [{res.name}] in {res.duration_sec:.2f}s")
    if not res.is_success and res.error:
        print(f"      Error: {res.error}")
    if not res.is_success and res.output:
        print(f"      Output: {res.output}")


def print_terminal_summary(results: list[GateResult], is_docker_active: bool, go_ver: str) -> None:
    """Prints standard user-facing report."""
    print("=" * 60)
    print("GitMap Linux Cross-Compilation Verification Gate")
    print(f"Go Toolchain:  {go_ver}")
    docker_status = "Available" if is_docker_active else "Absent (pure Go host gates enforced)"
    print(f"Docker Daemon: {docker_status}")
    print("-" * 60)
    for res in results:
        print_single_gate_report(res)
    print("-" * 60)
    has_passed = is_all_successful(results)
    conclusion = "ALL GATES PASSED (Hermetic Linux Compatibility Verified)" if has_passed else "FAILURES DETECTED"
    print(f"Conclusion:    {conclusion}")
    print("=" * 60)


def output_json_report(results: list[GateResult], is_docker_active: bool, go_ver: str) -> None:
    """Outputs machine-readable JSON summary for AI agents."""
    payload = {
        "go_version": go_ver,
        "is_docker_available": is_docker_active,
        "is_all_passed": is_all_successful(results),
        "gates": [asdict(r) for r in results],
    }
    print(json.dumps(payload, indent=2))


def parse_cli_args() -> argparse.Namespace:
    """Parses command line arguments."""
    parser = argparse.ArgumentParser(description="GitMap Linux Cross-Compilation Verification")
    parser.add_argument("--json", action="store_true", help="Output results in JSON format")
    parser.add_argument("--quiet", "-q", action="store_true", help="Print summary only")

    return parser.parse_args()


def main() -> int:
    """Main entrypoint orchestrating Linux verification."""
    args = parse_cli_args()
    cli_dir = resolve_cli_directory()
    go_ver = check_tool_version(["go", "version"])
    is_docker = check_is_docker_available()
    results = execute_all_gates(cli_dir)
    has_passed = is_all_successful(results)
    if args.json:
        output_json_report(results, is_docker, go_ver)

        return ExitCodeType.SUCCESS.value if has_passed else ExitCodeType.VIOLATIONS_FOUND.value
    if not args.quiet:
        print_terminal_summary(results, is_docker, go_ver)
    elif has_passed:
        print("✔ Linux cross-compilation verified successfully.")
    else:
        print("✖ Linux cross-compilation failed.")

    return ExitCodeType.SUCCESS.value if has_passed else ExitCodeType.VIOLATIONS_FOUND.value


if __name__ == "__main__":
    sys.exit(main())
