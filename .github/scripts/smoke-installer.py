#!/usr/bin/env python3
"""Cross-platform smoke test for gitmap installer.

Modes:
  source   Build gitmap from the current checkout into a tempdir, then run
           `<tempdir>/gitmap version` and assert it matches v$EXPECTED.
           Used by ci.yml on every PR — no network release dependency.

  release  Run gitmap/scripts/install.sh (or install.ps1 on Windows) against a
           published GitHub release (--version "v$EXPECTED" --no-discovery),
           then run the installed binary and assert. Used by release.yml.

Reads EXPECTED from env or falls back to gitmap/constants/constants.go.
Exits 0 on success, non-zero with diagnostic on failure.
"""
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
import time


def get_expected_version(repo_root: str) -> str:
    expected = os.environ.get("EXPECTED", "").strip()
    if expected:
        return expected.lstrip("v")

    constants_path = os.path.join(repo_root, "gitmap", "constants", "constants.go")
    if os.path.isfile(constants_path):
        try:
            with open(constants_path, "r", encoding="utf-8") as fh:
                for line in fh:
                    m = re.search(r'^(?:const|var)\s+Version\s*=\s*"([^"]+)"', line.strip())
                    if m:
                        return m.group(1).lstrip("v")
        except Exception:
            pass

    version_json = os.path.join(repo_root, "version.json")
    if os.path.isfile(version_json):
        try:
            with open(version_json, "r", encoding="utf-8") as fh:
                data = json.load(fh)
                return data.get("Version", data.get("version", "")).lstrip("v")
        except Exception:
            pass

    return ""


def load_deploy_manifest(repo_root: str):
    manifest_path = os.path.join(repo_root, "gitmap", "constants", "deploy-manifest.json")
    app_subdir = "gitmap-cli"
    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    legacy_subdirs = ["gitmap"]

    if os.path.isfile(manifest_path):
        try:
            with open(manifest_path, "r", encoding="utf-8") as fh:
                data = json.load(fh)
            app_subdir = data.get("appSubdir", app_subdir)
            if os.name == "nt":
                bin_name = data.get("binaryName", {}).get("windows", bin_name)
            else:
                bin_name = data.get("binaryName", {}).get("unix", bin_name)
            legacy_subdirs = data.get("legacyAppSubdirs", legacy_subdirs)
        except Exception:
            pass

    return app_subdir, bin_name, legacy_subdirs


def run_source_mode(repo_root: str, expected: str, workdir: str) -> str:
    print(f"▶ Building gitmap from source into {workdir}")
    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    bin_path = os.path.join(workdir, bin_name)
    gitmap_dir = os.path.join(repo_root, "gitmap")

    cmd = ["go", "build", "-o", bin_path, "."]
    res = subprocess.run(cmd, cwd=gitmap_dir, capture_output=True, text=True, encoding="utf-8")
    if res.returncode != 0:
        print(f"::error::go build failed (exit {res.returncode}):\n{res.stderr or res.stdout}", file=sys.stderr)
        sys.exit(3)

    return bin_path


def run_release_installer_with_retry(repo_root: str, expected: str, dest_dir: str, max_retries: int = 5, retry_delay_sec: int = 10) -> bool:
    is_windows = os.name == "nt"

    for attempt in range(1, max_retries + 1):
        print(f"▶ [Attempt {attempt}/{max_retries}] Running release installer for v{expected}...")
        if is_windows:
            script_path = os.path.join(repo_root, "gitmap", "scripts", "install.ps1")
            pwsh_bin = shutil.which("pwsh") or shutil.which("powershell") or "powershell"
            cmd = [
                pwsh_bin,
                "-File", script_path,
                "-Version", f"v{expected}",
                "-InstallDir", dest_dir,
                "-NoPath",
                "-NoDiscovery"
            ]
        else:
            script_path = os.path.join(repo_root, "gitmap", "scripts", "install.sh")
            cmd = [
                "bash", script_path,
                "--version", f"v{expected}",
                "--dir", dest_dir,
                "--no-path",
                "--no-discovery"
            ]

        res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, encoding="utf-8", errors="replace")
        if res.returncode == 0:
            print("  Installer finished successfully.")
            return True

        print(f"  Attempt {attempt} failed (exit {res.returncode}):\n{res.stdout}\n{res.stderr}")
        if attempt < max_retries:
            print(f"  Waiting {retry_delay_sec}s for release assets to propagate...")
            time.sleep(retry_delay_sec)

    return False


def locate_installed_binary(dest_dir: str, app_subdir: str, bin_name: str, legacy_subdirs: list) -> str:
    # 1. Primary candidate: dest/app_subdir/bin_name
    primary = os.path.join(dest_dir, app_subdir, bin_name)
    if os.path.isfile(primary):
        return primary

    # 2. Direct under dest
    direct = os.path.join(dest_dir, bin_name)
    if os.path.isfile(direct):
        return direct

    # 3. Legacy candidates
    for leg in legacy_subdirs:
        leg_path = os.path.join(dest_dir, leg, bin_name)
        if os.path.isfile(leg_path):
            return leg_path

    # 4. Recursive search
    for root, _, files in os.walk(dest_dir):
        for f in files:
            if f.lower() == bin_name.lower():
                return os.path.join(root, f)

    return ""


def run_release_mode(repo_root: str, expected: str, workdir: str) -> str:
    dest_dir = os.path.join(workdir, "install")
    os.makedirs(dest_dir, exist_ok=True)

    app_subdir, bin_name, legacy_subdirs = load_deploy_manifest(repo_root)

    success = run_release_installer_with_retry(repo_root, expected, dest_dir)
    if not success:
        print(f"::error::Release installer failed after retries for v{expected}", file=sys.stderr)
        sys.exit(3)

    bin_path = locate_installed_binary(dest_dir, app_subdir, bin_name, legacy_subdirs)
    if not bin_path or not os.path.isfile(bin_path):
        print(f"::error::Could not locate installed gitmap binary under {dest_dir}", file=sys.stderr)
        sys.exit(3)

    if os.name != "nt":
        os.chmod(bin_path, 0o755)

    return bin_path


def verify_binary_version(bin_path: str, expected: str):
    res = subprocess.run([bin_path, "version"], capture_output=True, text=True, encoding="utf-8")
    output = (res.stdout + res.stderr).strip()

    actual_line = ""
    for line in output.splitlines():
        line = line.strip()
        if re.match(r"^gitmap\s+v[0-9]", line):
            actual_line = line
            break

    print(f"▶ Actual output: {actual_line or output}")
    expected_line = f"gitmap v{expected}"

    if actual_line != expected_line:
        print(f"::error::Version mismatch!\n  expected: {expected_line}\n  actual:   {actual_line}", file=sys.stderr)
        sys.exit(4)

    print(f"✅ Installer smoke test passed: {actual_line}")


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    mode = sys.argv[1] if len(sys.argv) > 1 else "source"
    if mode not in ("source", "release"):
        print(f"::error::Unknown mode '{mode}' (expected 'source' or 'release')", file=sys.stderr)
        sys.exit(2)

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    expected = get_expected_version(repo_root)

    if not expected:
        print("::error::Could not determine expected version", file=sys.stderr)
        sys.exit(2)

    workdir = tempfile.mkdtemp(prefix="gitmap-smoke-")
    try:
        print(f"▶ Smoke mode:    {mode}")
        print(f"▶ Expected:      v{expected}")
        print(f"▶ Workdir:       {workdir}")

        if mode == "source":
            bin_path = run_source_mode(repo_root, expected, workdir)
        else:
            bin_path = run_release_mode(repo_root, expected, workdir)

        verify_binary_version(bin_path, expected)
    finally:
        shutil.rmtree(workdir, ignore_errors=True)


if __name__ == "__main__":
    main()
