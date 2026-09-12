#!/usr/bin/env python3
"""Cross-platform unused linter baseline-diff runner and gate.

Executes golangci-lint with the 'unused' analyzer, normalizes results,
and validates against baseline to prevent any new dead code regressions.
Compatible with 06-cicd-local-runner.py, GitHub Actions, and standalone CLI use.
"""
import os
import subprocess
import sys


def resolve_default_paths():
    if os.environ.get("CI") == "true" and sys.platform != "win32":
        default_out = os.path.join("/tmp", "lint-unused-current", "report.json")
        default_base = os.path.join("/tmp", "lint-unused-baseline", "report.json")
    else:
        temp_dir = os.path.join(".lovable", "temp", "cicd")
        default_out = os.path.join(temp_dir, "lint-unused-current.json")
        default_base = os.path.join(temp_dir, "lint-unused-baseline.json")
    return default_out, default_base


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    os.environ.setdefault("LINTER", "unused")
    default_out, default_base = resolve_default_paths()

    if "CURRENT_OUT" not in os.environ:
        os.environ["CURRENT_OUT"] = default_out

    if "BASELINE" not in os.environ and os.path.isfile(default_base) and os.path.getsize(default_base) > 0:
        os.environ["BASELINE"] = default_base

    script_path = os.path.join(os.path.dirname(__file__), "check-single-linter-diff.py")
    args = sys.argv[1:]
    if not args:
        args = ["cli"]

    cmd = [sys.executable, script_path] + args
    res = subprocess.run(cmd)
    sys.exit(res.returncode)


if __name__ == "__main__":
    main()
