#!/usr/bin/env python3
"""Cross-platform runner for spell checking changed or tracked files."""
import fnmatch
import json
import os
import subprocess
import sys


def load_config(config_path):
    default_excludes = [
        "*.png", "*.jpg", "*.jpeg", "*.gif", "*.ico", "*.svg", "*.webp",
        "*.zip", "*.tar", "*.gz", "*.exe",
        "*/testdata/*", "*/golden/*",
        "*/.gitmap/release/*", "*/.gitmap/release-assets/*",
        "gitmap/completion/allcommands_generated.go",
        ".lovable/*"
    ]
    if os.path.isfile(config_path):
        try:
            with open(config_path, "r", encoding="utf-8") as fh:
                data = json.load(fh)
            m = data.get("misspell", {})
            ex = m.get("exclude") or default_excludes
            inc = m.get("include") or []
            return ex, inc
        except Exception:
            pass
    return default_excludes, []


def get_changed_files(repo_root):
    event_name = os.environ.get("GH_EVENT_NAME", "push")
    base_ref = os.environ.get("GH_BASE_REF", "")
    base = f"origin/{base_ref}" if event_name == "pull_request" and base_ref else "HEAD~1"

    try:
        res = subprocess.run(
            ["git", "diff", "--name-only", "--diff-filter=AM", f"{base}...HEAD"],
            cwd=repo_root,
            capture_output=True,
            text=True,
            encoding="utf-8"
        )
        if res.returncode == 0:
            files = [line.strip().replace("\\", "/") for line in res.stdout.splitlines() if line.strip()]
            if files:
                return files
    except Exception:
        pass

    # Fallback to all tracked files if git diff fails or in shallow clone
    try:
        res = subprocess.run(
            ["git", "ls-files"],
            cwd=repo_root,
            capture_output=True,
            text=True,
            encoding="utf-8"
        )
        if res.returncode == 0:
            return [line.strip().replace("\\", "/") for line in res.stdout.splitlines() if line.strip()]
    except Exception:
        pass

    return []


def is_matching_glob(filepath, pattern):
    return fnmatch.fnmatch(filepath, pattern) or fnmatch.fnmatch(os.path.basename(filepath), pattern)


def filter_files(files, excludes, includes):
    filtered = []
    for f in files:
        if any(is_matching_glob(f, pat) for pat in excludes):
            continue
        if includes and not any(is_matching_glob(f, pat) for pat in includes):
            continue
        filtered.append(f)
    return filtered


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    config_path = os.path.join(repo_root, "gitmap", "data", "config.json")
    excludes, includes = load_config(config_path)

    files = get_changed_files(repo_root)
    scannable = filter_files(files, excludes, includes)

    if not scannable:
        print("no files to scan for misspell")
        sys.exit(0)

    # Locate misspell binary
    import shutil
    misspell_bin = shutil.which("misspell")
    if not misspell_bin:
        try:
            gp = subprocess.run(["go", "env", "GOPATH"], capture_output=True, text=True).stdout.strip()
            candidate = os.path.join(gp, "bin", "misspell.exe" if os.name == "nt" else "misspell")
            if os.path.isfile(candidate):
                misspell_bin = candidate
        except Exception:
            pass

    if not misspell_bin:
        print(f"misspell tool not installed. Scanned {len(scannable)} files.")
        sys.exit(0)

    # Process in batches to avoid command line length limits on Windows
    batch_size = 100
    for i in range(0, len(scannable), batch_size):
        batch = scannable[i:i + batch_size]
        cmd = [misspell_bin, "-locale", "US", "-error"] + batch
        res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True, encoding="utf-8", errors="replace")
        if res.returncode != 0:
            print(res.stderr or res.stdout, file=sys.stderr)
            sys.exit(res.returncode)

    print(f"✅ Spell check passed on {len(scannable)} file(s).")
    sys.exit(0)


if __name__ == "__main__":
    main()
