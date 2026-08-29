import os
import sys
import subprocess
import concurrent.futures
import time

# Force utf-8 stdout
sys.stdout.reconfigure(encoding='utf-8')

def run_command(name, cmd, cwd="."):
    print(f"--- [START] {name} ---")
    start = time.time()
    try:
        if cmd.startswith("docker ") or cmd.startswith("docker-compose "):
            print(f"Skipping docker command: {cmd}")
            return True, ""
            
        result = subprocess.run(cmd, cwd=cwd, shell=True, capture_output=True, text=True, encoding="utf-8", errors="replace")
        elapsed = time.time() - start
        if result.returncode == 0:
            print(f"--- [PASS] {name} ({elapsed:.2f}s) ---")
            return True, result.stdout
        else:
            print(f"--- [FAIL] {name} ({elapsed:.2f}s) ---")
            print(result.stdout)
            print(result.stderr)
            return False, result.stdout + "\n" + result.stderr
    except Exception as e:
        elapsed = time.time() - start
        print(f"--- [ERROR] {name} ({elapsed:.2f}s) ---")
        print(str(e))
        return False, str(e)

def main():
    gitmap_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", "..", "gitmap", "gitmap"))
    if not os.path.exists(gitmap_dir):
        gitmap_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", "gitmap"))

    # Running golangci-lint directly instead of using bash
    commands = [
        ("Go Vet", "go vet ./...", gitmap_dir),
        ("Full Suite Lint", "golangci-lint run ./...", gitmap_dir)
    ]

    success = True
    with concurrent.futures.ThreadPoolExecutor() as executor:
        futures = {executor.submit(run_command, name, cmd, cwd): name for name, cmd, cwd in commands}
        for future in concurrent.futures.as_completed(futures):
            passed, _ = future.result()
            if not passed:
                success = False

    if success:
        print("\n[SUCCESS] All CI/CD local checks passed!")
        sys.exit(0)
    else:
        print("\n[FAIL] CI/CD local checks failed. See output above.")
        sys.exit(1)

if __name__ == "__main__":
    main()
