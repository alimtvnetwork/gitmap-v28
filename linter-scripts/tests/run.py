#!/usr/bin/env python3
"""Stdlib test runner for ``linter-scripts/tests/``.

Mirrors ``linters-cicd/tests/run.py`` so the two suites have the
same shape and CI-side ergonomics. Discovers every ``test_*.py``
in this directory and runs them with verbose output. Exits 0 on
success, 1 on any failure or error.

Usage::

    python3 linter-scripts/tests/run.py

Wired into ``package.json`` as ``npm run test:linter-scripts``.
"""

from __future__ import annotations

import sys
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
ROOT = HERE.parent  # linter-scripts/

# ``top_level_dir`` lets unittest's discovery resolve the tests
# package correctly even when invoked from outside the repo root.
sys.path.insert(0, str(ROOT))
sys.path.insert(0, str(HERE))


def main() -> int:
    loader = unittest.TestLoader()
    suite = loader.discover(start_dir=str(HERE), pattern="test_*.py",
                            top_level_dir=str(ROOT))
    runner = unittest.TextTestRunner(verbosity=2)
    result = runner.run(suite)
    return 0 if result.wasSuccessful() else 1


if __name__ == "__main__":
    sys.exit(main())