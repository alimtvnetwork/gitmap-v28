"""Tests for the ``--diff-prev`` flag and ``--diff-base`` numeric
shorthand.

Two layers under test:

* **Unit** — ``_normalize_diff_base`` is the single source of truth
  for "did the user type a shorthand?" and must round-trip every
  shape the CLI promises in its ``--help`` text. Bare ints expand to
  ``HEAD~N``; ``~N`` / ``^N`` get prefixed with ``HEAD``; everything
  else is returned verbatim so we never silently rewrite a real ref.

* **CLI** — argparse-level wiring: ``--diff-prev`` must be mutually
  exclusive with ``--diff-base`` and ``--changed-files`` (stacking
  two intake sources is a footgun we already guard against for
  --diff-base/--changed-files), and a non-numeric ``--diff-prev``
  value must fail fast with a clear message instead of being shipped
  to git as ``HEAD~origin/main``.

The end-to-end "actually run against a real git repo" path is
covered indirectly: once the CLI normalises the ref, every
downstream code path is the existing --diff-base codepath that's
already exercised by the diff-mode test suites. We assert the
normalisation, not git itself.
"""

from __future__ import annotations

import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from conftest_shim import load_placeholder_linter  # noqa: E402

LINTER = (Path(__file__).resolve().parent.parent
          / "check-placeholder-comments.py")


def _load_linter_module():
    """Return the linter module via the shared conftest shim. The shim
    pre-registers the module in ``sys.modules`` BEFORE executing it,
    which is required for ``@dataclass`` decorators inside the linter
    to resolve their own type hints (Python 3.13 dataclass internals
    look up ``cls.__module__`` in ``sys.modules``)."""
    return load_placeholder_linter()


def _run(*args: str, cwd: Path) -> tuple[int, str, str]:
    r = subprocess.run([sys.executable, str(LINTER), *args],
                       cwd=cwd, capture_output=True, text=True)
    return r.returncode, r.stdout, r.stderr


class NormalizeDiffBaseUnit(unittest.TestCase):
    """``_normalize_diff_base`` shape-by-shape contract."""

    def setUp(self) -> None:
        self.mod = _load_linter_module()

    def test_bare_int_expands_to_head_tilde(self) -> None:
        # The headline shorthand: `--diff-base 1` ≡ `HEAD~1`.
        self.assertEqual(self.mod._normalize_diff_base("1"), "HEAD~1")
        self.assertEqual(self.mod._normalize_diff_base("3"), "HEAD~3")
        self.assertEqual(self.mod._normalize_diff_base("42"), "HEAD~42")

    def test_zero_expands_to_head_itself(self) -> None:
        # `HEAD~0 == HEAD`. Useful for smoke-testing the diff plumbing
        # without changing scope; should round-trip to a string git
        # accepts rather than being treated as "no shorthand".
        self.assertEqual(self.mod._normalize_diff_base("0"), "HEAD~0")

    def test_tilde_n_and_caret_n_get_head_prefix(self) -> None:
        # Users who type `--diff-base ~2` clearly mean HEAD~2; same
        # for `^2` (first-parent walk). Anything else after the
        # operator (e.g. `~main`) is NOT shorthand and must pass
        # through verbatim — covered by test_passthrough_*.
        self.assertEqual(self.mod._normalize_diff_base("~2"), "HEAD~2")
        self.assertEqual(self.mod._normalize_diff_base("^1"), "HEAD^1")

    def test_passthrough_named_refs(self) -> None:
        # Real refs must NEVER be rewritten — git owns this grammar.
        for ref in ("HEAD", "main", "origin/main", "v1.2.3",
                    "abc1234", "HEAD~1", "HEAD^1", "feature/foo"):
            self.assertEqual(self.mod._normalize_diff_base(ref), ref,
                             f"ref {ref!r} should pass through unchanged")

    def test_passthrough_mixed_tilde_with_name(self) -> None:
        # `~main` is NOT a numeric shorthand — only `~<int>` is.
        # Same for `^foo`. These must hit git verbatim so git's own
        # error message ("unknown revision") surfaces unchanged.
        self.assertEqual(self.mod._normalize_diff_base("~main"), "~main")
        self.assertEqual(self.mod._normalize_diff_base("^foo"), "^foo")

    def test_whitespace_is_trimmed_before_classification(self) -> None:
        # CLI users sometimes paste with stray whitespace; the helper
        # should classify on the stripped value but the original is
        # only rewritten when it IS a shorthand. Empty stays empty.
        self.assertEqual(self.mod._normalize_diff_base("  2  "), "HEAD~2")
        self.assertEqual(self.mod._normalize_diff_base(""), "")


class DiffPrevCliWiring(unittest.TestCase):
    """End-to-end argparse: mutex + numeric validation on --diff-prev."""

    def test_diff_prev_requires_non_negative_int(self) -> None:
        # `--diff-prev abc` must fail at the CLI boundary — we don't
        # want git to receive `HEAD~abc` and produce a confusing
        # "unknown revision" error several frames later.
        with tempfile.TemporaryDirectory() as td:
            code, _, err = _run("--diff-prev", "abc",
                                "--root", td, cwd=Path(td))
        self.assertEqual(code, 2)
        self.assertIn("--diff-prev", err)
        self.assertIn("non-negative integer", err)

    def test_diff_prev_rejects_negative(self) -> None:
        # argparse will pass through "-1" as a value (it's not a
        # known flag); our own validator catches it. Keeps the error
        # surface single-sourced.
        with tempfile.TemporaryDirectory() as td:
            code, _, err = _run("--diff-prev", "-1",
                                "--root", td, cwd=Path(td))
        self.assertEqual(code, 2)
        # Either argparse rejects "-1" as an unknown option OR our
        # validator catches it; both are acceptable hard-fails. We
        # only assert the exit code so we're not coupled to which
        # layer fires first.

    def test_diff_prev_mutex_with_diff_base(self) -> None:
        # Stacking two intake sources is a footgun: which one wins?
        # We refuse rather than silently picking, matching the
        # existing --diff-base/--changed-files mutex.
        with tempfile.TemporaryDirectory() as td:
            code, _, err = _run("--diff-prev", "1",
                                "--diff-base", "origin/main",
                                "--root", td,
                                cwd=Path(td))
        self.assertEqual(code, 2)
        self.assertIn("--diff-prev", err)
        self.assertIn("--diff-base", err)

    def test_diff_prev_mutex_with_changed_files(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            (Path(td) / "list.txt").write_text("")
            code, _, err = _run("--diff-prev", "1",
                                "--changed-files", "list.txt",
                                "--root", td,
                                cwd=Path(td))
        self.assertEqual(code, 2)
        self.assertIn("--diff-prev", err)
        self.assertIn("--changed-files", err)


if __name__ == "__main__":  # pragma: no cover - manual runner
    unittest.main()