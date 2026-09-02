#!/usr/bin/env python3
"""run.py — Universal cross-platform build, deploy, and runner engine for gitmap (Windows, macOS, Linux).

Usage:
  python run.py                                    # pull, build, deploy
  python run.py --no-pull                          # skip git pull
  python run.py --force-pull                       # discard local changes + pull
  python run.py --no-deploy                        # skip deploy step
  python run.py -r scan                            # build + run 'gitmap scan'
  python run.py -r agy status                      # build + run 'gitmap agy status'
  python run.py fix-repo --all                     # forward to fix_repo.py
  python run.py setup                              # forward to setup.py
"""

from __future__ import annotations

import argparse
import datetime
import json
import os
import shutil
import subprocess
import sys
from pathlib import Path

if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8", errors="replace")
        sys.stderr.reconfigure(encoding="utf-8", errors="replace")
    except Exception:
        pass

REPO_ROOT = Path(__file__).resolve().parent
GITMAP_DIR = REPO_ROOT / "gitmap"
DATA_DIR = GITMAP_DIR / "data"


def write_step(step_name: str, msg: str) -> None:
    print(f"\n  \033[35m[{step_name}]\033[0m \033[37m{msg}\033[0m")
    print("  " + ("-" * 50))


def write_ok(msg: str) -> None:
    print(f"  \033[32mOK {msg}\033[0m")


def write_info(msg: str) -> None:
    print(f"  \033[36m->\033[0m \033[90m{msg}\033[0m")


def write_warn(msg: str) -> None:
    print(f"  \033[33m!! {msg}\033[0m")


def write_fail(msg: str) -> None:
    print(f"  \033[31mXX {msg}\033[0m", file=sys.stderr)


def get_default_deploy_dir() -> Path:
    if os.name == "nt":
        local_app = os.getenv("LOCALAPPDATA")
        if local_app:
            return Path(local_app) / "gitmap-cli"
        return Path.home() / "AppData" / "Local" / "gitmap-cli"
    return Path.home() / ".local" / "bin" / "gitmap-cli"


def load_deploy_config() -> tuple[Path, str, str]:
    cfg_file = GITMAP_DIR / "powershell.json"
    deploy_path = get_default_deploy_dir()
    bin_name = "gitmap.exe" if os.name == "nt" else "gitmap"
    build_out = "./bin"

    if cfg_file.is_file():
        try:
            with open(cfg_file, "r", encoding="utf-8") as f:
                data = json.load(f)
            dp = data.get("deployPath")
            if dp:
                deploy_path = Path(dp)
            bn = data.get("binaryName")
            if bn:
                bin_name = bn if (os.name != "nt" or bn.endswith(".exe")) else f"{bn}.exe"
            bo = data.get("buildOutput")
            if bo:
                build_out = bo
        except OSError:
            pass

    return deploy_path, bin_name, build_out


def save_deploy_config(deploy_path: Path, bin_name: str) -> None:
    cfg_file = GITMAP_DIR / "powershell.json"
    data = {
        "deployPath": str(deploy_path),
        "buildOutput": "./bin",
        "binaryName": bin_name,
        "goSource": "./gitmap",
        "copyData": True,
    }
    try:
        with open(cfg_file, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2)
    except OSError:
        pass


def pull_code(force: bool) -> bool:
    write_step("Git", "Checking for updates...")
    if force:
        write_info("Force pull requested — discarding local changes")
        subprocess.run(["git", "reset", "--hard", "HEAD"], cwd=str(REPO_ROOT), check=False)
        subprocess.run(["git", "clean", "-fd"], cwd=str(REPO_ROOT), check=False)

    res = subprocess.run(["git", "pull"], cwd=str(REPO_ROOT))
    if res.returncode == 0:
        write_ok("Git pull completed")
        return True
    write_warn("Git pull failed or has uncommitted changes; continuing with local files")
    return False


def build_binary(out_binary: Path) -> bool:
    write_step("Build", "Compiling Go binary...")

    # Run build-stamp
    stamp_script = REPO_ROOT / "scripts" / "build_stamp.py"
    if stamp_script.is_file():
        subprocess.run([sys.executable, str(stamp_script)], check=False)

    out_binary.parent.mkdir(parents=True, exist_ok=True)

    # Git metadata
    def git_out(args: list[str]) -> str:
        try:
            return subprocess.check_output(["git"] + args, cwd=str(REPO_ROOT), stderr=subprocess.DEVNULL, text=True).strip()
        except subprocess.SubprocessError:
            return ""

    build_commit = git_out(["rev-parse", "HEAD"])
    build_branch = git_out(["rev-parse", "--abbrev-ref", "HEAD"])
    build_repo = git_out(["remote", "get-url", "origin"])
    build_date = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    abs_repo_root = str(REPO_ROOT).replace("\\", "/")

    ldflags = (
        f"-X 'github.com/alimtvnetwork/gitmap-v28/gitmap/constants.RepoPath={abs_repo_root}' "
        f"-X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildCommit={build_commit}' "
        f"-X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildBranch={build_branch}' "
        f"-X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildRepo={build_repo}' "
        f"-X 'github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.BuildDate={build_date}'"
    )

    go_exe = shutil.which("go")
    if not go_exe:
        write_fail("Go compiler not found on PATH!")
        return False

    cmd = [go_exe, "build", "-ldflags", ldflags, "-o", str(out_binary), "."]
    res = subprocess.run(cmd, cwd=str(GITMAP_DIR))
    if res.returncode != 0:
        write_fail("Go build failed")
        return False

    write_ok(f"Binary built successfully -> {out_binary}")
    return True


def copy_data_directory(src_data: Path, dest_dir: Path) -> None:
    if not src_data.is_dir():
        return
    dest_data = dest_dir / "data"
    dest_data.mkdir(parents=True, exist_ok=True)
    for f in src_data.glob("*.json"):
        shutil.copy2(f, dest_data / f.name)
    write_ok(f"Copied data assets -> {dest_data}")


def deploy_binary(built_binary: Path, deploy_dir: Path, bin_name: str) -> Path:
    write_step("Deploy", f"Deploying to {deploy_dir}...")
    deploy_dir.mkdir(parents=True, exist_ok=True)

    dest_binary = deploy_dir / bin_name
    shutil.copy2(built_binary, dest_binary)
    try:
        dest_binary.chmod(0o755)
    except OSError:
        pass
    write_ok(f"Deployed {bin_name} to {deploy_dir}")

    # Create alias
    alias_name = "gm.exe" if os.name == "nt" else "gm"
    alias_dest = deploy_dir / alias_name
    try:
        shutil.copy2(dest_binary, alias_dest)
        write_ok(f"Created alias {alias_name}")
    except OSError:
        pass

    copy_data_directory(DATA_DIR, deploy_dir)
    save_deploy_config(deploy_dir, bin_name)

    # Sync to sibling if applicable on Windows
    if os.name == "nt" and "gitmap-cli" in deploy_dir.name:
        sibling = deploy_dir.parent / "gitmap"
        if sibling.is_dir():
            try:
                shutil.copy2(dest_binary, sibling / bin_name)
                copy_data_directory(DATA_DIR, sibling)
                write_ok(f"Synced to sibling deployment -> {sibling}")
            except OSError:
                pass

    return dest_binary


def main() -> int:
    # Forward common subcommands if passed as first arg
    if len(sys.argv) > 1:
        first = sys.argv[1]
        if first in ("fix-repo", "fix_repo"):
            return subprocess.run([sys.executable, str(REPO_ROOT / "fix_repo.py")] + sys.argv[2:]).returncode
        elif first == "setup":
            return subprocess.run([sys.executable, str(REPO_ROOT / "setup.py")] + sys.argv[2:]).returncode
        elif first == "init":
            return subprocess.run([sys.executable, str(REPO_ROOT / "init.py")] + sys.argv[2:]).returncode
        elif first in ("uninstall", "uninstall-quick"):
            return subprocess.run([sys.executable, str(REPO_ROOT / "uninstall_quick.py")] + sys.argv[2:]).returncode

    parser = argparse.ArgumentParser(description="Build, deploy, and run gitmap CLI cross-platform.", add_help=False)
    parser.add_argument("--no-pull", action="store_true", help="Skip git pull")
    parser.add_argument("--force-pull", action="store_true", help="Discard local changes and git pull")
    parser.add_argument("--no-deploy", action="store_true", help="Skip deploy step")
    parser.add_argument("--no-setup", action="store_true", help="Skip setup step")
    parser.add_argument("--deploy-path", help="Deploy target directory")
    parser.add_argument("-d", "--deploy", action="store_true", help="Explicit deploy flag")
    parser.add_argument("-r", "--run", action="store_true", help="Run gitmap after build")
    parser.add_argument("-t", "--test", action="store_true", help="Run Go tests")
    parser.add_argument("--debug-repo-detect", action="store_true", help="Debug diagnostics")
    parser.add_argument("--quiet", action="store_true", help="Quiet output")
    parser.add_argument("-h", "--help", action="store_true", help="Show help")

    args, remaining = parser.parse_known_args()

    if args.help:
        print(__doc__)
        return 0

    if not args.no_pull:
        pull_code(args.force_pull)

    deploy_dir, bin_name, build_out_dir = load_deploy_config()
    if args.deploy_path:
        deploy_dir = Path(args.deploy_path)

    out_binary = REPO_ROOT / build_out_dir / bin_name
    if not build_binary(out_binary):
        return 1

    active_binary = out_binary
    if not args.no_deploy:
        active_binary = deploy_binary(out_binary, deploy_dir, bin_name)

    if args.run or remaining:
        run_args = remaining
        if not run_args:
            # Default run action
            run_args = ["status"]
        write_step("Run", f"Executing: {active_binary.name} {' '.join(run_args)}")
        res = subprocess.run([str(active_binary)] + run_args)
        return res.returncode

    write_step("Done", "All operations completed successfully.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
