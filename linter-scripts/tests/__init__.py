# Marker so unittest discovery treats this directory as a package.
# Intentionally empty — the runner adds the parent directory to
# ``sys.path`` so each ``test_*.py`` can import the linter modules
# via the shim defined in ``conftest_shim.py``.