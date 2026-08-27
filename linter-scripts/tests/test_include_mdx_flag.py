"""End-to-end tests for the ``--include-mdx`` convenience flag.

The flag is a single-keystroke shortcut for "scan ``.md`` AND
``.mdx`` files". Semantics under test:

* **Default + flag** — without any explicit ``--extension``, the
  baseline stays ``md`` and the union becomes ``("md", "mdx")``.
* **Explicit ``--extension`` + flag** — the flag *adds* ``mdx`` to
  whatever the user already passed; it never replaces the baseline.
* **Idempotence** — passing both ``--include-mdx`` and
  ``--extension mdx`` is a no-op (mdx is deduped, segment unchanged).
* **Canonical equivalence** — ``--include-mdx`` and
  ``--extension md --extension mdx`` produce the same cache segment
  AND the same cache key, so a sentinel written by one form is hit
  by the other (the whole point of the canonicalisation).
* **Off by default** — the legacy ``.md``-only baseline is preserved
  for any invocation that omits the flag.

The first three are subprocess tests against the CLI so we exercise
argparse + the resolution block end-to-end. The equivalence test
leans on the cache-hit log line as a black-box signal that two
different invocations produced byte-identical sentinels.
"""

from __future__ import annotations

import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

LINTER = (Path(__file__).resolve().parent.parent
          / "check-placeholder-comments.py")


def _run(*args: str, cwd: Path) -> tuple[int, str, str]:
    r = subprocess.run([sys.executable, str(LINTER), *args],
                       cwd=cwd, capture_output=True, text=True)
    return r.returncode, r.stdout, r.stderr


def _make_repo(td: Path) -> tuple[Path, Path]:
    """spec/intro.md + spec/intro.mdx, both clean, plus a fresh
    cache dir. Mirrors the layout used by the cache-segregation
    suite so the two test files exercise the same fixture shape."""
    spec = td / "spec"
    spec.mkdir()
    (spec / "intro.md").write_text("# spec\nplain prose.\n")
    (spec / "intro.mdx").write_text("# spec\nplain prose.\n")
    cache = td / "cache"
    cache.mkdir()
    return spec, cache


def _segments(cache: Path) -> set[str]:
    return {p.parent.name for p in cache.rglob("*.pass")}


class IncludeMdxBaseline(unittest.TestCase):

    def test_flag_off_keeps_md_only_baseline(self) -> None:
        """No flag → segment is ``ext-md`` (the legacy default).
        Regression guard: a future refactor that flips the default
        to include mdx must trip this test."""
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, err = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache),
                cwd=Path(td),
            )
            self.assertEqual(rc, 0, f"stderr={err!r}")
            self.assertEqual(_segments(cache), {"ext-md"})

    def test_flag_on_unions_with_default(self) -> None:
        """``--include-mdx`` alone → ``ext-md+mdx``."""
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache), "--include-mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            self.assertEqual(_segments(cache), {"ext-md+mdx"})


class IncludeMdxComposition(unittest.TestCase):
    """Composition with explicit ``--extension`` flags."""

    def test_explicit_baseline_is_preserved(self) -> None:
        """``--extension txt --include-mdx`` → ``ext-mdx+txt``.
        The flag must NOT silently re-introduce ``md`` when the
        operator deliberately scoped the baseline elsewhere."""
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache),
                "--extension", "txt", "--include-mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            self.assertEqual(_segments(cache), {"ext-mdx+txt"})

    def test_idempotent_with_explicit_mdx(self) -> None:
        """When the operator passes ``--extension mdx`` they have
        explicitly scoped the baseline AWAY from ``md``; the
        ``--include-mdx`` flag then collapses to a no-op (mdx is
        already in the allowlist). Result: ``ext-mdx``, not
        ``ext-md+mdx`` — the flag never silently re-introduces ``md``
        when the operator opted out of it."""
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache),
                "--include-mdx", "--extension", "mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            # mdx (explicit, replaces the default ``md`` baseline)
            # + mdx (the convenience flag, deduped) → ext-mdx.
            self.assertEqual(_segments(cache), {"ext-mdx"})


class IncludeMdxCanonicalEquivalence(unittest.TestCase):
    """The flag must collapse to the same cache key as the long
    form, so a sentinel written by one invocation is HIT by the
    other. Cache-hit log line is the black-box witness."""

    def test_flag_form_hits_sentinel_from_long_form(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            # Write a sentinel via the long form.
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache),
                "--extension", "md", "--extension", "mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            # Now look it up via the convenience form. MUST hit.
            rc, out, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache), "--include-mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            self.assertIn("cache hit", out)
            self.assertIn("ext-md+mdx/", out)
            self.assertEqual(_segments(cache), {"ext-md+mdx"})

    def test_long_form_hits_sentinel_from_flag_form(self) -> None:
        """Symmetric direction: write via flag, hit via long form."""
        with tempfile.TemporaryDirectory() as td:
            spec, cache = _make_repo(Path(td))
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache), "--include-mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            rc, out, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--cache-dir", str(cache),
                "--extension", "mdx", "--extension", "md",
                cwd=Path(td),
            )
            self.assertEqual(rc, 0)
            self.assertIn("cache hit", out)


class IncludeMdxActuallyScansMdx(unittest.TestCase):
    """Wiring smoke test: the flag must reach ``iter_markdown_files``
    so an .mdx file is actually parsed for violations, not just
    used as a cache-segment label."""

    def test_mdx_violation_is_reported_with_flag(self) -> None:
        with tempfile.TemporaryDirectory() as td:
            spec, _ = _make_repo(Path(td))
            # Inject a guaranteed P-001 (empty TODO) into the .mdx.
            (spec / "intro.mdx").write_text(
                "# spec\n\n<!-- TODO: -->\n"
            )
            # Without the flag: .mdx is invisible, exit 0.
            rc, _, _ = _run(
                "--root", str(spec), "--repo-root", td,
                cwd=Path(td),
            )
            self.assertEqual(rc, 0,
                "baseline scan should ignore .mdx and pass")
            # With the flag: .mdx is scanned, P-001 fires, exit 1.
            rc, out, _ = _run(
                "--root", str(spec), "--repo-root", td,
                "--include-mdx",
                cwd=Path(td),
            )
            self.assertEqual(rc, 1,
                f"expected violation exit, got rc={rc}, out={out!r}")
            self.assertIn("P-001", out)
            self.assertIn("intro.mdx", out)


if __name__ == "__main__":
    unittest.main()