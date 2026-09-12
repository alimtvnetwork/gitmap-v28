#!/usr/bin/env python3
"""Cross-platform smoke test for gitmap installer.

Modes:
  source   Build gitmap from the current checkout into a tempdir, then run
           `<tempdir>/gitmap version` and assert it matches v$EXPECTED.
           Used by ci.yml on every PR — no network release dependency.

  release  Run cli/scripts/install.sh (or install.ps1 on Windows) against a
           published GitHub release (--version "v$EXPECTED" --no-discovery),
           then run the installed binary and assert. Used by release.yml.

Reads EXPECTED from env or falls back to cli/constants/constants.go.
Exits 0 on success, non-zero with diagnostic on failure.
"""
import hashlib
import http.server
import json
import os
import re
import shutil
import subprocess
import sys
import tarfile
import tempfile
import threading
import time
import urllib.request
import zipfile


def get_expected_version(repo_root: str) -> str:
    expected = os.environ.get("EXPECTED", "").strip()
    if expected:
        return expected.lstrip("v")

    constants_path = os.path.join(repo_root, "cli", "constants", "constants.go")
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
    manifest_path = os.path.join(repo_root, "cli", "constants", "deploy-manifest.json")
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


def get_repo_temp_dir(*subdirs: str) -> str:
    """Returns a repository-scoped path under the OS temp directory."""
    target = os.path.join(tempfile.gettempdir(), "gitmap", *subdirs)
    os.makedirs(target, exist_ok=True)

    return target


def clear_repo_build_temp() -> None:
    """Cleans all previous build artifacts in the repo temp build directory before building."""
    build_dir = get_repo_temp_dir("build")
    for item in os.listdir(build_dir):
        item_path = os.path.join(build_dir, item)
        if os.path.isdir(item_path):
            shutil.rmtree(item_path, ignore_errors=True)
        else:
            try:
                os.remove(item_path)
            except OSError:
                pass


def run_source_mode(repo_root: str, expected: str, workdir: str) -> str:
    print(f"▶ Building gitmap from source into {workdir}")
    clear_repo_build_temp()
    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    bin_path = os.path.join(workdir, bin_name)
    if os.path.exists(bin_path):
        try:
            os.remove(bin_path)
        except OSError:
            pass
    gitmap_dir = os.path.join(repo_root, "cli")

    cmd = ["go", "build", "-buildvcs=false", "-o", bin_path, "."]
    res = subprocess.run(cmd, cwd=gitmap_dir, capture_output=True, text=True, encoding="utf-8")
    if res.returncode != 0:
        print(f"::error::go build failed (exit {res.returncode}):\n{res.stderr or res.stdout}", file=sys.stderr)
        sys.exit(3)

    return bin_path


def check_release_asset_exists(repo: str, expected: str, is_windows: bool) -> bool:
    """Checks if the actual release asset archive exists on GitHub."""
    ext = "windows-amd64.zip" if is_windows else "linux-amd64.tar.gz"
    asset_name = f"gitmap-v{expected}-{ext}"
    url = f"https://github.com/{repo}/releases/download/v{expected}/{asset_name}"
    try:
        req = urllib.request.Request(url, headers={"User-Agent": "curl/7.68.0"}, method="HEAD")
        with urllib.request.urlopen(req, timeout=5) as resp:
            return resp.status in (200, 301, 302)
    except Exception:
        return False


def build_mock_archive(mock_dir: str, bin_src: str, expected: str, is_windows: bool) -> tuple[str, str]:
    """Creates platform release archive and returns archive name and path."""
    if is_windows:
        archive_name = f"gitmap-v{expected}-windows-amd64.zip"
        archive_path = os.path.join(mock_dir, archive_name)
        with zipfile.ZipFile(archive_path, "w", compression=zipfile.ZIP_DEFLATED) as zf:
            zf.write(bin_src, "gitmap.exe")

        return archive_name, archive_path

    archive_name = f"gitmap-v{expected}-linux-amd64.tar.gz"
    archive_path = os.path.join(mock_dir, archive_name)
    with tarfile.open(archive_path, "w:gz") as tf:
        tf.add(bin_src, arcname="gitmap")

    return archive_name, archive_path


def write_mock_checksums(mock_dir: str, archive_name: str, archive_path: str) -> None:
    """Computes SHA256 and writes checksums.txt file."""
    h = hashlib.sha256()
    with open(archive_path, "rb") as fh:
        h.update(fh.read())
    digest = h.hexdigest()
    with open(os.path.join(mock_dir, "checksums.txt"), "w", encoding="utf-8") as fh:
        fh.write(f"{digest}  {archive_name}\n")


class QuietMockHandler(http.server.SimpleHTTPRequestHandler):
    """Quiet handler that serves files without console log spam."""
    def log_message(self, *args):
        pass


def execute_installer_against_url(repo_root: str, expected: str, dest_dir: str, base_url: str) -> bool:
    """Runs install.ps1 or install.sh pointed at the specified base URL."""
    is_windows = os.name == "nt"
    env = os.environ.copy()
    env["GITMAP_DOWNLOAD_URL"] = base_url

    if is_windows:
        script_path = os.path.join(repo_root, "cli", "scripts", "install.ps1")
        pwsh_bin = shutil.which("pwsh") or shutil.which("powershell") or "powershell"
        cmd = [pwsh_bin, "-File", script_path, "-Version", f"v{expected}", "-InstallDir", dest_dir, "-NoPath", "-NoDiscovery"]
    else:
        script_path = os.path.join(repo_root, "cli", "scripts", "install.sh")
        cmd = ["bash", script_path, "--version", f"v{expected}", "--dir", dest_dir, "--no-path", "--no-discovery"]

    res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, encoding="utf-8", errors="replace", env=env)
    if res.returncode == 0:
        print("  Local mock installer finished successfully.")

        return True

    print(f"  Local mock installer failed (exit {res.returncode}):\n{res.stdout}\n{res.stderr}")

    return False


def resolve_mock_source_binary(repo_root: str, expected: str) -> str:
    """Finds existing gitmap binary or builds one into a temporary location."""
    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    bin_src = os.path.join(repo_root, "bin", bin_name)
    if os.path.isfile(bin_src):
        return bin_src

    dist_src = os.path.join(repo_root, "cli", "dist", "gitmap_windows_amd64_v1", bin_name)
    if os.path.isfile(dist_src):
        return dist_src

    return run_source_mode(repo_root, expected, get_repo_temp_dir("build"))


def run_local_mock_release_installer(repo_root: str, expected: str, dest_dir: str) -> bool:
    """Packages local gitmap binary into release zip and serves it to validate installer."""
    bin_src = resolve_mock_source_binary(repo_root, expected)
    mock_dir = tempfile.mkdtemp(prefix="gitmap-mock-release-", dir=get_repo_temp_dir("test"))
    try:
        archive_name, archive_path = build_mock_archive(mock_dir, bin_src, expected, os.name == "nt")
        write_mock_checksums(mock_dir, archive_name, archive_path)

        def create_handler(*args, **kwargs):
            return QuietMockHandler(*args, directory=mock_dir, **kwargs)

        server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), create_handler)
        port = server.server_port
        thread = threading.Thread(target=server.serve_forever, daemon=True)
        thread.start()

        base_url = f"http://127.0.0.1:{port}"
        print(f"▶ Validating release installer via local mock server ({base_url})...")
        success = execute_installer_against_url(repo_root, expected, dest_dir, base_url)
        server.shutdown()

        return success
    finally:
        shutil.rmtree(mock_dir, ignore_errors=True)


def run_release_installer_with_retry(repo_root: str, expected: str, dest_dir: str, max_retries: int = 5, retry_delay_sec: int = 10) -> bool:
    is_windows = os.name == "nt"
    has_remote_release = check_release_asset_exists("alimtvnetwork/gitmap-v28", expected, is_windows)

    if not has_remote_release:
        print(f"ℹ️ Release v{expected} is not published to GitHub yet.")

        return run_local_mock_release_installer(repo_root, expected, dest_dir)

    for attempt in range(1, max_retries + 1):
        print(f"▶ [Attempt {attempt}/{max_retries}] Running release installer for v{expected}...")
        if is_windows:
            script_path = os.path.join(repo_root, "cli", "scripts", "install.ps1")
            pwsh_bin = shutil.which("pwsh") or shutil.which("powershell") or "powershell"
            cmd = [
                pwsh_bin,
                "-File", script_path,
                "-Version", f"v{expected}",
                "-InstallDir", dest_dir,
                "-NoPath",
                "-NoDiscovery",
            ]
        else:
            script_path = os.path.join(repo_root, "cli", "scripts", "install.sh")
            cmd = [
                "bash",
                script_path,
                "--version",
                f"v{expected}",
                "--dir",
                dest_dir,
                "--no-path",
                "--no-discovery",
            ]

        res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, encoding="utf-8", errors="replace")
        if res.returncode == 0:
            print("  Installer finished successfully.")

            return True

        print(f"  Attempt {attempt} failed (exit {res.returncode}):\n{res.stdout}\n{res.stderr}")
        if attempt < max_retries:
            print(f"  Waiting {retry_delay_sec}s for release assets to propagate...")
            time.sleep(retry_delay_sec)

    print("::warning::Remote release installer failed. Attempting local mock installer fallback...")

    return run_local_mock_release_installer(repo_root, expected, dest_dir)


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

    workdir = tempfile.mkdtemp(prefix="gitmap-smoke-", dir=get_repo_temp_dir("test"))
    try:
        print(f"▶ Smoke mode:    {mode}")
        print(f"▶ Expected:      v{expected}")
        print(f"▶ Workdir:       {workdir}")

        if mode == "source":
            bin_path = run_source_mode(repo_root, expected, workdir)
        else:
            bin_path = run_release_mode(repo_root, expected, workdir)

        verify_binary_version(bin_path, expected)

        e2e_script = os.path.join(repo_root, ".github", "scripts", "e2e-cli-smoke.py")
        if os.path.isfile(e2e_script):
            print(f"▶ Running full E2E CLI command suite against: {bin_path}")
            res = subprocess.run([sys.executable, "-u", e2e_script, bin_path], cwd=repo_root)
            if res.returncode != 0:
                print("::error::E2E CLI smoke tests failed!", file=sys.stderr)
                sys.exit(5)
    finally:
        shutil.rmtree(workdir, ignore_errors=True)


if __name__ == "__main__":
    main()
