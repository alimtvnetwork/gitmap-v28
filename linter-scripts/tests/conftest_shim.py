"""Shared import shim for the linter-scripts/tests/ unit tests.

The linter scripts live as hyphenated, executable ``.py`` files
(``check-placeholder-comments.py`` etc.) so they can be invoked
directly from CI / npm scripts without a ``python -m`` wrapper.
That naming is hostile to ``import`` — Python identifiers can't
contain hyphens — so each test module imports the linter under
test through this shim instead of through ``import``.

The shim:

1. Resolves the linter's path relative to this file (so tests work
   from any CWD).
2. Builds an ``importlib`` spec with a Python-legal module name
   (underscores instead of hyphens).
3. Registers the module in ``sys.modules`` *before* ``exec_module``
   runs — this is required by ``dataclasses._is_type`` under
   Python 3.13, which looks the module up in ``sys.modules`` while
   processing each ``@dataclass`` decorator. Skipping the
   pre-registration step yields a confusing
   ``AttributeError: 'NoneType' object has no attribute '__dict__'``
   from deep inside the ``dataclasses`` module.
4. Caches the loaded module so repeated imports across test
   modules are free.

Usage::

    from conftest_shim import load_placeholder_linter
    chk = load_placeholder_linter()
    excerpt = chk._parse_unified_diff_post(diff_text)
"""

from __future__ import annotations

import importlib.util
import sys
from pathlib import Path
from types import ModuleType

_HERE = Path(__file__).resolve().parent
_SCRIPTS = _HERE.parent

_CACHE: dict[str, ModuleType] = {}


def _load(file_name: str, module_name: str) -> ModuleType:
    if module_name in _CACHE:
        return _CACHE[module_name]
    path = _SCRIPTS / file_name
    if not path.is_file():
        raise FileNotFoundError(
            f"conftest_shim: linter source missing at {path}"
        )
    spec = importlib.util.spec_from_file_location(module_name, path)
    if spec is None or spec.loader is None:
        raise ImportError(
            f"conftest_shim: spec_from_file_location returned None for {path}"
        )
    mod = importlib.util.module_from_spec(spec)
    # Pre-registration: see step 3 in the module docstring.
    sys.modules[module_name] = mod
    spec.loader.exec_module(mod)
    _CACHE[module_name] = mod
    return mod


def load_placeholder_linter() -> ModuleType:
    """Return the loaded ``check-placeholder-comments.py`` module."""
    return _load("check-placeholder-comments.py",
                 "check_placeholder_comments")


def load_tunable_constants_linter() -> ModuleType:
    """Return the loaded ``check-tunable-constants.py`` module."""
    return _load("check-tunable-constants.py",
                 "check_tunable_constants")