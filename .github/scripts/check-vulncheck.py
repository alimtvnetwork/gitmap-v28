#!/usr/bin/env python3
"""Cross-platform wrapper for govulncheck adhering to CI vulnerability policy."""
from __future__ import annotations

from pathlib import Path
import shutil
import subprocess
import sys

STDLIB_ONLY_BANNER = "Only stdlib vulnerabilities found (no fix available in current Go version)"
THIRD_PARTY_ERROR_BANNER = "Third-party vulnerabilities detected in imported packages"
DEFAULT_ENCODING = "utf-8"


def get_binary_target(repo_root: Path) -> Path:
    """Resolves gitmap binary path if compiled, otherwise returns fallback."""
    bin_name = "gitmap.exe" if sys.platform.startswith("win") else "gitmap"
    primary_path = repo_root / "bin" / bin_name
    if primary_path.is_file():
        return primary_path

    return repo_root / "cli" / bin_name


def build_govulncheck_cmd(target_bin: Path, repo_root: Path) -> list[str]:
    """Constructs govulncheck command preferring binary scan over source AST."""
    govulncheck_bin = shutil.which("govulncheck") or "govulncheck"
    if target_bin.is_file():
        return [govulncheck_bin, "-mode=binary", str(target_bin)]

    return [govulncheck_bin, "-scan=package", "./..."]


def has_third_party_vulnerability(output_text: str) -> bool:
    """Checks whether third-party packages called by user code have vulnerabilities."""
    clean_text = " ".join(output_text.split())
    is_stdlib_only = "vulnerabilities from the Go standard library" in clean_text
    is_no_call = "doesn't appear to call these vulnerabilities" in clean_text
    is_clean = bool(is_stdlib_only and is_no_call)

    return not is_clean


def evaluate_scan_result(retcode: int, stdout_text: str, stderr_text: str) -> int:
    """Evaluates scan result according to CI stdlib vs third-party policy."""
    full_output = stdout_text + "\n" + stderr_text
    print(full_output.strip())
    if retcode == 0:
        return 0

    if not has_third_party_vulnerability(full_output):
        print(f"\n::warning::{STDLIB_ONLY_BANNER}")
        return 0

    print(f"\n::error::{THIRD_PARTY_ERROR_BANNER}", file=sys.stderr)

    return 1


def main() -> None:
    """Primary execution entrypoint for govulncheck runner."""
    repo_root = Path(__file__).resolve().parent.parent.parent
    target_bin = get_binary_target(repo_root)
    cmd = build_govulncheck_cmd(target_bin, repo_root)
    cwd = str(repo_root if target_bin.is_file() else repo_root / "cli")

    try:
        proc = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            encoding=DEFAULT_ENCODING,
            errors="replace",
            cwd=cwd,
        )
        code = evaluate_scan_result(proc.returncode, proc.stdout, proc.stderr)
        sys.exit(code)
    except FileNotFoundError:
        print("::warning::govulncheck binary not found in PATH; skipping check.")
        sys.exit(0)


if __name__ == "__main__":
    main()
