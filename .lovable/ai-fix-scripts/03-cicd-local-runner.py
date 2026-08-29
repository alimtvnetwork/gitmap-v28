#!/usr/bin/env python3
"""Auto-generated CI/CD local runner (03-cicd-local-runner.py).
Re-generate or execute directly: python .lovable/ai-fix-scripts/03-cicd-local-runner.py
"""
import concurrent.futures
import os
import re
import subprocess
import sys
import time

sys.stdout.reconfigure(encoding="utf-8")

ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
GITMAP_DIR = os.path.join(ROOT_DIR, "gitmap")


def run_bare_stderr_check():
    """Validates no bare fmt.Fprintln(os.Stderr, err) exists in gitmap/cmd."""
    cmd_dir = os.path.join(GITMAP_DIR, "cmd")
    pattern = re.compile(r"fmt\.Fprintln\(os\.Stderr,\s*err\)")
    offenders = []
    for root, _, files in os.walk(cmd_dir):
        for f in files:
            if f.endswith(".go") and not f.endswith("_test.go"):
                path = os.path.join(root, f)
                with open(path, "r", encoding="utf-8", errors="replace") as fh:
                    for idx, line in enumerate(fh, start=1):
                        if pattern.search(line):
                            rel = os.path.relpath(path, ROOT_DIR)
                            offenders.append(f"{rel}:{idx}:{line.strip()}")
    if offenders:
        return False, "Bare stderr error prints found in gitmap/cmd:\n" + "\n".join(offenders)
    return True, "No bare fmt.Fprintln(os.Stderr, err) found."


def run_legacy_refs_check():
    """Scans repository for forbidden legacy refs gitmap-v[567]\\b unless whitelisted."""
    pattern = re.compile(r"gitmap-v[567]\b")
    exclude_dirs = {".git", "node_modules", "dist", "build", "bin", ".next", ".gitmap", "vendor", "coverage"}
    offenders = []
    for root, dirs, files in os.walk(ROOT_DIR):
        dirs[:] = [d for d in dirs if d not in exclude_dirs]
        for f in files:
            if f == "check-legacy-refs.sh":
                continue
            path = os.path.join(root, f)
            try:
                with open(path, "r", encoding="utf-8", errors="ignore") as fh:
                    for idx, line in enumerate(fh, start=1):
                        if pattern.search(line) and "gitmap-legacy-ref-allow" not in line:
                            rel = os.path.relpath(path, ROOT_DIR)
                            offenders.append(f"{rel}:{idx}:{line.strip()}")
            except Exception:
                pass
    if offenders:
        return False, "Legacy references found:\n" + "\n".join(offenders)
    return True, "No forbidden legacy refs found."


def run_cmd(name, cmd, cwd):
    start = time.time()
    try:
        res = subprocess.run(cmd, cwd=cwd, shell=True, capture_output=True, text=True, encoding="utf-8", errors="replace")
        elapsed = time.time() - start
        if res.returncode == 0:
            return True, f"--- [PASS] {name} ({elapsed:.2f}s) ---"
        out = (res.stdout or "") + "\n" + (res.stderr or "")
        return False, f"--- [FAIL] {name} ({elapsed:.2f}s) ---\n{out.strip()}"
    except Exception as exc:
        elapsed = time.time() - start
        return False, f"--- [ERROR] {name} ({elapsed:.2f}s) ---\n{str(exc)}"


def run_changelog_sync_check():
    """Validates constants.Version has a matching heading in changelog.md."""
    const_file = os.path.join(GITMAP_DIR, "constants", "constants.go")
    changelog_file = os.path.join(ROOT_DIR, "changelog.md")
    version = None
    with open(const_file, "r", encoding="utf-8", errors="replace") as fh:
        for line in fh:
            m = re.search(r'^var Version\s*=\s*"([^"]+)"', line)
            if m:
                version = m.group(1)
                break
    if not version:
        return False, f"Could not parse var Version from {const_file}"
    escaped_ver = re.escape(version)
    pattern = re.compile(rf"^##\s+\[?v?{escaped_ver}\]?(?:\s|$)", re.MULTILINE)
    with open(changelog_file, "r", encoding="utf-8", errors="replace") as fh:
        content = fh.read()
    if not pattern.search(content):
        return False, f"CHANGELOG drift: constants.Version is {version} but no matching heading exists in changelog.md"
    return True, f"constants.Version {version} in sync with changelog.md"


def run_generate_drift_check():
    """Runs go generate ./... and verifies no files under gitmap/ were modified by it."""
    before_res = subprocess.run("git status --porcelain .", cwd=GITMAP_DIR, shell=True, capture_output=True, text=True, encoding="utf-8")
    res = subprocess.run("go generate ./...", cwd=GITMAP_DIR, shell=True, capture_output=True, text=True, encoding="utf-8")
    if res.returncode != 0:
        return False, f"go generate ./... failed:\n{res.stderr or res.stdout}"
    after_res = subprocess.run("git status --porcelain .", cwd=GITMAP_DIR, shell=True, capture_output=True, text=True, encoding="utf-8")
    if before_res.stdout != after_res.stdout:
        diff_res = subprocess.run("git diff .", cwd=GITMAP_DIR, shell=True, capture_output=True, text=True, encoding="utf-8")
        return False, f"go generate ./... modified files under gitmap/:\n{diff_res.stdout or diff_res.stderr}"
    return True, "All generated files in gitmap/ are in sync."


def main():
    print("=== CI/CD Local Runner (03-cicd-local-runner.py) ===")
    checks = [
        ("Legacy Refs Check", run_legacy_refs_check, None, None),
        ("Bare Stderr Check", run_bare_stderr_check, None, None),
        ("Error Management Linter", None, "python linter-scripts/check-error-management.py", ROOT_DIR),
        ("Boolean & Enum Linter", None, "python linter-scripts/check-enum-and-boolean.py", ROOT_DIR),
        ("Changelog Version Sync", run_changelog_sync_check, None, None),
        ("Generate Drift Check", run_generate_drift_check, None, None),
        ("Installer Smoke", None, "pwsh -File .github/scripts/smoke-installer.ps1 source", ROOT_DIR),
        ("Web Version Sync", None, "npx vitest run src/test/version-sync.test.ts", ROOT_DIR),
        ("Docs Site Build", None, "npm run build", ROOT_DIR),
        ("CLI Zero-Args Smoke", None, "go run .", GITMAP_DIR),
        ("Go Vet", None, "go vet ./...", GITMAP_DIR),
        ("Compile Gate", None, "go test -run=^$ ./... -count=1", GITMAP_DIR),
        ("Full Suite Lint", None, "golangci-lint run ./...", GITMAP_DIR),
    ]

    all_passed = True
    results = []
    batch_size = 3

    for i in range(0, len(checks), batch_size):
        batch = checks[i : i + batch_size]
        with concurrent.futures.ThreadPoolExecutor(max_workers=batch_size) as executor:
            futures = {}
            for name, fn, cmd, cwd in batch:
                if fn is not None:
                    futures[executor.submit(fn)] = name
                else:
                    futures[executor.submit(run_cmd, name, cmd, cwd)] = name

            for future in concurrent.futures.as_completed(futures):
                name = futures[future]
                res = future.result()
                passed = res[0]
                output = res[1]
                if not passed:
                    all_passed = False
                results.append((name, passed, output))

    for name, passed, output in sorted(results, key=lambda x: x[0]):
        status = "✅ PASS" if passed else "❌ FAIL"
        print(f"\n{status} [{name}]")
        print(output)

    if all_passed:
        print("\n[SUCCESS] All CI/CD local checks passed (exit 0)!")
        sys.exit(0)
    else:
        print("\n[FAIL] Some CI/CD checks failed (exit 1)!")
        sys.exit(1)


if __name__ == "__main__":
    main()
