# Learned: Test Inventory & Atomic File Change Tracking

**Recorded:** 2026-09-12
**Scope:** Repository-wide test management, local CI runner, and multi-agent coordination

---

## Architecture Overview

To balance fast developer loops with reliable release validation, two core systems have been established:

1. **Centralized Test Inventory (`.lovable/test-inventory.json`):**
   - Automatically catalogs all test files, test suites, test functions, and package directories across Go, Python, and TypeScript.
   - Maintained via `python 03-ai-scripts/33-test-inventory-generator.py`.
   - Maps source targets to corresponding test suites so that affected tests can be precisely resolved without brute-force full scans.

2. **Atomic File Change Tracking (`.lovable/temp/recent-file-changes.json`):**
   - Cross-platform file locking via `.lovable/temp/recent-file-changes.lock` using `msvcrt` on Windows and `fcntl` on POSIX systems.
   - Appends distinct modified repository-relative paths across multi-agent turns.
   - Automatically computes associated test files from the inventory manifest.
   - Used by release orchestrator and test runner scripts to run target tests when explicitly authorized.

---

## Tooling Commands

```bash
# Scan repository and regenerate the test inventory manifest:
python 03-ai-scripts/33-test-inventory-generator.py

# Atomically record modified files under lock:
python 03-ai-scripts/33-test-inventory-generator.py --record <path1> <path2>

# Query recently modified files and their associated tests:
python 03-ai-scripts/33-test-inventory-generator.py --query-recent

# Clear recent changes after a verified release or test run:
python 03-ai-scripts/33-test-inventory-generator.py --clear
```

---

## Operational Rule

- Standard development prompts and coding guideline fixes MUST NOT execute test suites; they run `03-ai-scripts/06-cicd-local-runner.py --no-tests`.
- Only explicit release workflows (`03-ai-scripts/29-release-orchestrator.py` or explicit release prompts) or explicit user commands may execute tests.
